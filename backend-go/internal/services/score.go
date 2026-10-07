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
	balanceBefore := student.TotalScore

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

	// 余额不为负（与教室端一致）。
	newBalance := maxInt(balanceBefore+amount, 0)
	if err := tx.Model(&models.Student{}).Where("id = ?", student.ID).
		Update("total_score", newBalance).Error; err != nil {
		return score, err
	}

	log := models.ScoreLog{
		StudentID:     student.ID,
		ScoreID:       score.ID,
		BalanceBefore: balanceBefore,
		BalanceAfter:  newBalance,
		Description:   reason,
	}
	if err := tx.Create(&log).Error; err != nil {
		return score, err
	}

	if err := syncPetForDelta(tx, student.ID, newBalance, amount); err != nil {
		return score, err
	}

	student.TotalScore = newBalance
	return score, nil
}

// syncPetForDelta 依据积分变动同步宠物经验并校正等级。
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
	if student.TotalScore < amount {
		return nil, ErrBadRequest(fmt.Sprintf("积分不足，当前余额：%d", student.TotalScore))
	}

	var score models.Score
	err := s.db.Transaction(func(tx *gorm.DB) error {
		balanceBefore := student.TotalScore

		score = models.Score{
			StudentID: student.ID,
			ClassID:   student.ClassID,
			Amount:    -amount,
			Reason:    "兑换消耗：" + reason,
			GivenBy:   spentBy,
		}
		if err := tx.Create(&score).Error; err != nil {
			return err
		}

		newBalance := balanceBefore - amount
		if err := tx.Model(&models.Student{}).Where("id = ?", student.ID).
			Update("total_score", newBalance).Error; err != nil {
			return err
		}

		log := models.ScoreLog{
			StudentID:     student.ID,
			ScoreID:       score.ID,
			BalanceBefore: balanceBefore,
			BalanceAfter:  newBalance,
			Description:   score.Reason,
		}
		if err := tx.Create(&log).Error; err != nil {
			return err
		}

		if err := syncPetForDelta(tx, student.ID, newBalance, -amount); err != nil {
			return err
		}

		student.TotalScore = newBalance
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 推送给班级大屏（score_update + is_spend，同 Laravel spendScore）。
	s.publishScoreUpdate(student, -amount, score.Reason, true)

	return &score, nil
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
func (s *ScoreService) Undo(original *models.Score, operatedBy uint) (*models.Score, error) {
	var undo models.Score
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var student models.Student
		if err := tx.First(&student, original.StudentID).Error; err != nil {
			return err
		}

		undoAmount := -original.Amount
		undo = models.Score{
			StudentID: student.ID,
			ClassID:   original.ClassID,
			Amount:    undoAmount,
			Reason:    "撤回操作（原：" + original.Reason + "）",
			GivenBy:   operatedBy,
		}
		if err := tx.Create(&undo).Error; err != nil {
			return err
		}

		newBalance := maxInt(student.TotalScore+undoAmount, 0)
		if err := tx.Model(&models.Student{}).Where("id = ?", student.ID).
			Update("total_score", newBalance).Error; err != nil {
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
