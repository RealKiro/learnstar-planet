// 积分服务：事务内完成积分记录 + 余额更新 + 审计日志 + 宠物经验同步。
package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// ScoreService 积分操作服务。
type ScoreService struct {
	db     *gorm.DB
	events *DisplayEvents
}

// NewScoreService 创建积分服务（内部自带大屏事件发布器：积分变动的 score_update 事件）。
func NewScoreService(db *gorm.DB) *ScoreService {
	return &ScoreService{db: db, events: NewDisplayEvents(db)}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GiveScore 给单个学生加/减分（amount 可为负）。返回新积分记录。
func (s *ScoreService) GiveScore(student *models.Student, amount int, reason string, givenBy uint, scoreRuleID *uint) (*models.Score, error) {
	var score models.Score
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		score, err = giveScoreTx(tx, student, amount, reason, givenBy, scoreRuleID)
		return err
	})
	if err != nil {
		return nil, err
	}

	// 推送给班级大屏（score_update）：事务提交后发布，失败不影响主流程
	//（Laravel 在事务闭包内发布，事件可能早于提交可见；Go 端放在提交后，语义更严谨）。
	s.publishScoreUpdate(student, amount, reason, false)

	return &score, nil
}

// giveScoreTx 在指定事务内完成一次加减分。
func giveScoreTx(tx *gorm.DB, student *models.Student, amount int, reason string, givenBy uint, scoreRuleID *uint) (models.Score, error) {
	// ⚠️ 余额走 SQL 侧原子读-改-写（addScoreAtomic），不再「从内存快照算好新值再整值写回」：
	// 后者在并发下会丢更新（两条请求各自从同一份旧余额算出同一个新余额，后写覆盖先写）。
	balanceBefore, newBalance, err := addScoreAtomic(tx, student.ID, amount)
	if err != nil {
		return models.Score{}, err
	}
	return recordScoreTx(tx, student, amount, reason, givenBy, scoreRuleID, balanceBefore, newBalance)
}

// recordScoreTx 落一条积分记录 + 审计日志 + 同步宠物经验，并回写内存余额。
//
// 抽出来的理由：加分、消费、（转赠的）转出转入这几条链路只有「余额如何变」不同
// （原子加减 / 条件扣减），而记录与副作用完全一致——避免每条链路各自复制一遍
// scores + score_logs + 宠物同步（复制出来的副本最容易日后改漏一处）。
func recordScoreTx(tx *gorm.DB, student *models.Student, amount int, reason string, givenBy uint, scoreRuleID *uint, balanceBefore, balanceAfter int) (models.Score, error) {
	score := models.Score{
		StudentID:   student.ID,
		ClassID:     student.ClassID,
		ScoreRuleID: scoreRuleID,
		Amount:      amount,
		Reason:      reason,
		GivenBy:     givenBy,
	}
	if err := tx.Create(&score).Error; err != nil {
		return score, err
	}

	if err := tx.Create(&models.ScoreLog{
		StudentID:     student.ID,
		ScoreID:       score.ID,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   reason,
	}).Error; err != nil {
		return score, err
	}

	if err := syncPetForDelta(tx, student.ID, balanceAfter, amount); err != nil {
		return score, err
	}

	student.TotalScore = balanceAfter
	return score, nil
}

// currentScore 读取学生当前积分（事务内的权威值）。
func currentScore(tx *gorm.DB, studentID uint) (int, error) {
	var s models.Student
	if err := tx.Select("total_score").Where("id = ?", studentID).First(&s).Error; err != nil {
		return 0, err
	}
	return s.TotalScore, nil
}

// addScoreAtomic 原子地把学生积分加上 delta（可为负），负余额钳制为 0，返回（变更前, 变更后）。
//
// 为什么必须交给 SQL：并发下「读出来算好再整值写回」会丢更新——两条请求各自基于同一份
// 旧余额算出同一个新值，后写的覆盖先写的；「扣到负数就钳 0」这条不变式同理会被绕过。
// 这里用 `total_score = total_score + ?` 让数据库做主，再单独把负数钳成 0。
// 钳制不写 MAX()/GREATEST()：SQLite 是标量 MAX(a,b)、MySQL/Postgres 是 GREATEST(a,b)，
// 拆成「先加减、再把负数置 0」两步即可跨三库可移植（同事务内执行，语义等价）。
func addScoreAtomic(tx *gorm.DB, studentID uint, delta int) (before, after int, err error) {
	if before, err = currentScore(tx, studentID); err != nil {
		return 0, 0, err
	}
	if err = tx.Model(&models.Student{}).Where("id = ?", studentID).
		Update("total_score", gorm.Expr("total_score + ?", delta)).Error; err != nil {
		return 0, 0, err
	}
	if err = tx.Model(&models.Student{}).Where("id = ? AND total_score < 0", studentID).
		Update("total_score", 0).Error; err != nil {
		return 0, 0, err
	}
	if after, err = currentScore(tx, studentID); err != nil {
		return 0, 0, err
	}
	return before, after, nil
}

// deductScoreAtomic 原子「校验 + 扣减」：条件更新一次完成，余额不足则一个字节都不写。
// ok=false 表示余额不足（before = 不足时的权威余额，供报错文案使用）。
// 这是消除「先查后扣」竞态的关键：把余额判定放进 UPDATE 的 WHERE 里，并发下不会穿仓。
func deductScoreAtomic(tx *gorm.DB, studentID uint, amount int) (before, after int, ok bool, err error) {
	if before, err = currentScore(tx, studentID); err != nil {
		return 0, 0, false, err
	}
	res := tx.Model(&models.Student{}).
		Where("id = ? AND total_score >= ?", studentID, amount).
		Update("total_score", gorm.Expr("total_score - ?", amount))
	if res.Error != nil {
		return 0, 0, false, res.Error
	}
	if res.RowsAffected == 0 {
		return before, before, false, nil
	}
	if after, err = currentScore(tx, studentID); err != nil {
		return 0, 0, false, err
	}
	return before, after, true, nil
}

// syncPetForDelta 依据积分变动同步宠物经验并校正等级。
//
// ⚠️ 这里的「读宠物 → 改字段 → Save 整行」看似是读-改-写竞态，但在所有调用点都**天然串行**：
// 调用前，同一个事务已经用原子 UPDATE 改过该学生的 students 行（addScoreAtomic /
// deductScoreAtomic），并发请求会在「学生行」上排队（SQLite 全局写锁；MySQL/Postgres 行锁），
// 于是本函数读宠物的时刻必然晚于前一个事务提交。**新调用点务必保持「先动学生行、再动宠物」的顺序。**
// （真正没有保护的是 PetService.Feed / Rename：它们既不在事务里、也不动学生行——属外观层的
// mood / 名字更新，并发双击可能丢一次 +20，影响仅限展示，故未改。）
func syncPetForDelta(tx *gorm.DB, studentID uint, newBalance, amount int) error {
	var pet models.Pet
	err := tx.Where("student_id = ?", studentID).First(&pet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil // 学生尚无宠物
	}
	if err != nil {
		return err
	}

	if amount > 0 {
		pet.AddExperience(amount)
	} else {
		pet.RemoveExperience(maxInt(-amount, 0))
	}
	pet.SyncLevelWithScore(newBalance)
	return tx.Save(&pet).Error
}

// SpendScore 消耗积分（兑换奖励时调用），金额为正表示扣减。
func (s *ScoreService) SpendScore(student *models.Student, amount int, reason string, spentBy uint) (*models.Score, error) {
	var score models.Score
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		score, err = spendScoreTx(tx, student, amount, reason, spentBy)
		return err
	})
	if err != nil {
		return nil, err
	}

	// 推送给班级大屏（score_update + is_spend，同 Laravel spendScore）。
	s.publishScoreUpdate(student, -amount, score.Reason, true)

	return &score, nil
}

// spendScoreTx 在指定事务内完成一次积分消费（金额为正表示扣减）。
func spendScoreTx(tx *gorm.DB, student *models.Student, amount int, reason string, spentBy uint) (models.Score, error) {
	// ⚠️ 「校验余额」与「扣减」合并为一次条件更新（deductScoreAtomic）：并发下既不会穿仓，
	// 也不会把人家的消费覆盖掉。余额不足则一个字节都不写、直接报错。
	balanceBefore, newBalance, ok, err := deductScoreAtomic(tx, student.ID, amount)
	if err != nil {
		return models.Score{}, err
	}
	if !ok {
		return models.Score{}, ErrBadRequest(fmt.Sprintf("积分不足，当前余额：%d", balanceBefore))
	}
	return recordScoreTx(tx, student, -amount, "兑换消耗："+reason, spentBy, nil, balanceBefore, newBalance)
}

// BatchGive 批量加分，单事务原子提交，返回成功人数。
func (s *ScoreService) BatchGive(students []*models.Student, amount int, reason string, givenBy uint, scoreRuleID *uint) (int, error) {
	count := 0
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, student := range students {
			if _, err := giveScoreTx(tx, student, amount, reason, givenBy, scoreRuleID); err != nil {
				return err
			}
			count++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	// 批量加分是逐个学生调用 giveScore（同 Laravel batchGiveScore），故逐个发布 score_update。
	for _, student := range students {
		s.publishScoreUpdate(student, amount, reason, false)
	}

	return count, nil
}

// Undo 撤回一条积分记录（创建等额反向记录，与 Laravel undoScore 一致）。
//
// ⚠️ 幂等守卫：撤回是**加回**积分，所以同一条记录必须只能被撤回一次——原先没有任何守卫，
// 反复点「撤回」就能把积分刷上去（Laravel undoScore 亦如此）。这里两层防护：
//  1. 事务内先查 `undo_of_score_id = 本记录` 是否已存在，存在即 400；
//  2. Score.UndoOfScoreID 带 uniqueIndex，库层兜住并发双击（第二个插入失败并整体回滚）。
func (s *ScoreService) Undo(original *models.Score, operatedBy uint) (*models.Score, error) {
	var undo models.Score
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var alreadyUndone int64
		if err := tx.Model(&models.Score{}).
			Where("undo_of_score_id = ?", original.ID).Count(&alreadyUndone).Error; err != nil {
			return err
		}
		if alreadyUndone > 0 {
			return ErrBadRequest("该积分记录已撤回")
		}

		var student models.Student
		if err := tx.First(&student, original.StudentID).Error; err != nil {
			return err
		}

		// ⚠️ 撤回金额取「实际生效」的变动额，而不是账面金额：原记录若被「余额不为负」钳制过
		// （例如余额 10 时扣 30，实际只扣到 0），按 -Amount 撤回会凭空多补 20 分。
		// score_logs 里存着权威的 before/after，用它算实际变动额即可精确回补。
		// 找不到审计行时回退为 -Amount（老数据；撤回流水本身此前也不写审计行）。
		undoAmount := -original.Amount
		var audit models.ScoreLog
		if auditErr := tx.Where("score_id = ?", original.ID).Order("id DESC").First(&audit).Error; auditErr == nil {
			undoAmount = -(audit.BalanceAfter - audit.BalanceBefore)
		}
		originalID := original.ID
		undo = models.Score{
			StudentID:     student.ID,
			ClassID:       original.ClassID,
			Amount:        undoAmount,
			Reason:        "撤回操作（原：" + original.Reason + "）",
			GivenBy:       operatedBy,
			UndoOfScoreID: &originalID,
		}
		if err := tx.Create(&undo).Error; err != nil {
			return err
		}

		// 余额同样走 SQL 侧原子读-改-写（见 addScoreAtomic）。
		balanceBefore, newBalance, err := addScoreAtomic(tx, student.ID, undoAmount)
		if err != nil {
			return err
		}

		// 为撤回补一条审计行：让「每条 scores 都有配对的 before/after」这条不变式完整
		// （撤回原先只写 scores、不写 score_logs），也是上面「按实际变动额撤回」能递归生效的前提。
		if err := tx.Create(&models.ScoreLog{
			StudentID:     student.ID,
			ScoreID:       undo.ID,
			BalanceBefore: balanceBefore,
			BalanceAfter:  newBalance,
			Description:   undo.Reason,
		}).Error; err != nil {
			return err
		}

		if err := syncPetForDelta(tx, student.ID, newBalance, undoAmount); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &undo, nil
}

// FindScoredForTeacher 教师可见范围内的单条积分记录（越权 ID 返回 404）。
func (s *ScoreService) FindScoredForTeacher(classIDs []uint, id uint) (*models.Score, error) {
	var score models.Score
	err := s.db.Where("id = ? AND class_id IN ?", id, classIDs).First(&score).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("积分记录不存在或不可见")
	}
	if err != nil {
		return nil, err
	}
	return &score, nil
}

// ScoreSummary 多班级积分汇总（累计/今日/本周）。
type ScoreSummary struct {
	Total    int `json:"total"`
	Today    int `json:"today"`
	ThisWeek int `json:"this_week"`
}

// Summary 汇总多个班级的累计/今日/本周积分。
func (s *ScoreService) Summary(classIDs []uint) (*ScoreSummary, error) {
	if len(classIDs) == 0 {
		return &ScoreSummary{}, nil
	}

	now := util.Now()
	total, err := s.sumAmount(classIDs, time.Time{})
	if err != nil {
		return nil, err
	}
	today, err := s.sumAmount(classIDs, util.StartOfDay(now))
	if err != nil {
		return nil, err
	}
	week, err := s.sumAmount(classIDs, util.StartOfWeek(now))
	if err != nil {
		return nil, err
	}

	return &ScoreSummary{Total: total, Today: today, ThisWeek: week}, nil
}

func (s *ScoreService) sumAmount(classIDs []uint, from time.Time) (int, error) {
	q := s.db.Model(&models.Score{}).Where("class_id IN ?", classIDs)
	if !from.IsZero() {
		q = q.Where("created_at >= ?", from)
	}

	var sum int64
	if err := q.Select("COALESCE(SUM(amount), 0)").Scan(&sum).Error; err != nil {
		return 0, err
	}
	return int(sum), nil
}

// RecentScore 最近积分记录（含学生姓名）。
type RecentScore struct {
	ID          uint      `json:"id"`
	StudentName string    `json:"student_name"`
	Amount      int       `json:"amount"`
	Reason      string    `json:"reason"`
	CreatedAt   time.Time `json:"created_at"`
}

// Recent 多班级最近积分记录（默认 20 条）。
func (s *ScoreService) Recent(classIDs []uint, limit int) ([]RecentScore, error) {
	if len(classIDs) == 0 {
		return []RecentScore{}, nil
	}
	if limit <= 0 {
		limit = 20
	}

	var scores []models.Score
	if err := s.db.Where("class_id IN ?", classIDs).
		Order("created_at DESC, id DESC").Limit(limit).Find(&scores).Error; err != nil {
		return nil, err
	}

	ids := make([]uint, 0, len(scores))
	for _, sc := range scores {
		ids = append(ids, sc.StudentID)
	}

	nameMap := map[uint]string{}
	if len(ids) > 0 {
		var students []models.Student
		if err := s.db.Where("id IN ?", ids).Find(&students).Error; err != nil {
			return nil, err
		}
		for _, st := range students {
			nameMap[st.ID] = st.Name
		}
	}

	out := make([]RecentScore, 0, len(scores))
	for _, sc := range scores {
		out = append(out, RecentScore{
			ID:          sc.ID,
			StudentName: nameMap[sc.StudentID],
			Amount:      sc.Amount,
			Reason:      sc.Reason,
			CreatedAt:   sc.CreatedAt,
		})
	}
	return out, nil
}

// History 单个学生的积分历史（默认 20 条）。
func (s *ScoreService) History(studentID uint, limit int) ([]models.Score, error) {
	if limit <= 0 {
		limit = 20
	}
	var scores []models.Score
	err := s.db.Where("student_id = ?", studentID).
		Order("created_at DESC, id DESC").Limit(limit).Find(&scores).Error
	return scores, err
}

// ============================================================
// 大屏事件发布（score_update）
// ============================================================

// publishScoreUpdate 发布 score_update 大屏事件（载荷字段/键序同 Laravel ScoreService）。
// student 上的 total_score 已由 giveScoreTx 更新为事务后的余额；宠物字段在提交后重新读取。
func (s *ScoreService) publishScoreUpdate(student *models.Student, amount int, reason string, isSpend bool) {
	if student == nil {
		return
	}

	petLevel, petExp, petMood := petEventFields(s.petOf(student.ID))
	publishDisplayEvent(s.events, student.ClassID, DisplayEventScoreUpdate, DisplayScoreUpdateData{
		StudentID:     student.ID,
		StudentName:   student.Name,
		StudentNo:     student.StudentNo,
		Amount:        amount,
		Reason:        reason,
		TotalScore:    student.TotalScore,
		PetLevel:      petLevel,
		PetExperience: petExp,
		PetMood:       petMood,
		IsSpend:       isSpend,
	})
}

// petOf 取学生宠物（无宠物返回 nil，对应 Laravel 的 $student->pet?->xxx）。
func (s *ScoreService) petOf(studentID uint) *models.Pet {
	var pet models.Pet
	if err := s.db.Where("student_id = ?", studentID).First(&pet).Error; err != nil {
		return nil
	}
	return &pet
}

// petEventFields 把宠物实体转成事件载荷里的可空字段（无宠物时三个字段均为 null）。
func petEventFields(pet *models.Pet) (*int, *int, *int) {
	if pet == nil {
		return nil, nil, nil
	}
	level, exp, mood := pet.Level, pet.Experience, pet.Mood
	return &level, &exp, &mood
}

// ============================================================
// 按规则加减分（POST /teacher/scores/give-by-rule/:ruleId）
// 移植自 Laravel ScoreService::giveScoreByRule（内部走 giveScore，
// 故事务/余额钳制/审计日志/宠物经验/大屏事件的语义与既有 GiveScore 完全一致）。
// ============================================================

// GiveScoreByRule 按规则给单个学生加减分：金额与原因取自规则（reason = 规则名），
// 并记录 score_rule_id。
func (s *ScoreService) GiveScoreByRule(student *models.Student, rule *models.ScoreRule, givenBy uint) (*models.Score, error) {
	if rule == nil {
		return nil, ErrBadRequest("规则不存在")
	}
	var ruleID *uint
	if rule.ID != 0 {
		id := rule.ID
		ruleID = &id
	}
	return s.GiveScore(student, rule.Amount, rule.Name, givenBy, ruleID)
}
