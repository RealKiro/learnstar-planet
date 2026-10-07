// 年级战场（PK）服务：教室端同年级班级排行、教师端排行与本班统计、挑战发起。
//
// 忠实移植自 Laravel App\Services\PkService（leaderboard / myStats / challenge）与
// DisplayController::classroomPKLeaderboard。
//
// 有意差异：
//  1. challenge 在 Laravel 里另把 {challenger, target, challenged_at, expires_at, status} 写入
//     Cache（键 pk_challenge:<my_class>:<target_class>，7 天）。Go 端没有 Cache；全仓也**没有**
//     任何读取该键的接口或路径（挑战记录只写不读），故此处只返回同样的响应载荷、不落库——
//     对外可观测行为一致。若将来需要「待应战」列表，应补一张表而不是复活 Cache 语义。
//  2. 汇总口径差异：Laravel 的 metricsFor 用 `$students->avg(fn($s) => $s->pet->level ?? 0)`，
//     即「无宠物的学生按 0 级参与平均」；Go 端用一次宠物查询算出同一口径（无宠物 = 0 级）。
//  3. 教室端排除 code = 'DEMO' 的学校（Laravel whereHas school code != 'DEMO'），
//     教师端与 Laravel 一致**不**排除。
package services

import (
	"errors"
	"math"
	"sort"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// pkPeakPetLevel 「巅峰宠物」等级门槛（同 Laravel PkService::metricsFor 的 `level >= 8`）。
const pkPeakPetLevel = 8

// pkChallengeDays PK 挑战有效期（天，同 Laravel PkService::challenge 的 addDays(7)）。
const pkChallengeDays = 7

// PkService 年级战场服务。
type PkService struct {
	db    *gorm.DB
	scope *Scope
}

// NewPkService 创建 PK 服务。
func NewPkService(db *gorm.DB, scope *Scope) *PkService {
	return &PkService{db: db, scope: scope}
}

// PkMetrics 单个班级的 PK 指标（字段名同 Laravel PkService::metricsFor）。
type PkMetrics struct {
	TotalScore   int     `json:"totalScore"`
	StudentCount int     `json:"studentCount"`
	AvgLevel     float64 `json:"avgLevel"`
	PeakCount    int     `json:"peakCount"`
	WeekGrowth   int     `json:"weekGrowth"`
}

// PkClassRow 教室端（班级码）一行（字段名逐字同 Laravel classroomPKLeaderboard）。
type PkClassRow struct {
	Name         string  `json:"name"`
	TotalScore   int     `json:"totalScore"`
	StudentCount int     `json:"studentCount"`
	AvgLevel     float64 `json:"avgLevel"`
	PeakCount    int     `json:"peakCount"`
	WeekGrowth   int     `json:"weekGrowth"`
	IsOwn        bool    `json:"isOwn"`
}

// PkTeacherRow 教师端一行（字段名逐字同 Laravel PkService::leaderboard）。
type PkTeacherRow struct {
	ClassID      uint    `json:"class_id"`
	Name         string  `json:"name"`
	IsOwn        bool    `json:"isOwn"`
	TotalScore   int     `json:"totalScore"`
	StudentCount int     `json:"studentCount"`
	AvgLevel     float64 `json:"avgLevel"`
	PeakCount    int     `json:"peakCount"`
	WeekGrowth   int     `json:"weekGrowth"`
}

// PkStats 本班 PK 统计（字段名逐字同 Laravel PkService::myStats：metrics + rank）。
type PkStats struct {
	TotalScore   int     `json:"totalScore"`
	StudentCount int     `json:"studentCount"`
	AvgLevel     float64 `json:"avgLevel"`
	PeakCount    int     `json:"peakCount"`
	WeekGrowth   int     `json:"weekGrowth"`
	Rank         int     `json:"rank"`
}

// PkChallengeData 挑战发起结果的数据（data 字段名逐字同 Laravel PkService::challenge）。
type PkChallengeData struct {
	TargetClass string `json:"target_class"`
	ExpiresAt   string `json:"expires_at"`
}

// PkChallengeResult 挑战发起结果（Message 为 Laravel 原文，由处理器放进 message 字段）。
type PkChallengeResult struct {
	Message string
	Data    PkChallengeData
}

// LeaderboardForClass 教室端：同年级各班 PK 排行（按总积分降序，isOwn 标记本班）。
func (s *PkService) LeaderboardForClass(classID uint) ([]PkClassRow, error) {
	rows := []PkClassRow{}

	var myClass models.ClassRoom
	err := s.db.First(&myClass, classID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Laravel：班级不存在或没有年级 → 直接返回空数组（不报错）。
		return rows, nil
	}
	if err != nil {
		return nil, err
	}
	if myClass.Grade == "" {
		return rows, nil
	}

	gradeClasses, err := s.gradeClasses(myClass, true)
	if err != nil {
		return nil, err
	}

	for _, class := range gradeClasses {
		metrics, err := s.metricsFor(class.ID)
		if err != nil {
			return nil, err
		}
		rows = append(rows, PkClassRow{
			Name:         class.Name,
			TotalScore:   metrics.TotalScore,
			StudentCount: metrics.StudentCount,
			AvgLevel:     metrics.AvgLevel,
			PeakCount:    metrics.PeakCount,
			WeekGrowth:   metrics.WeekGrowth,
			IsOwn:        class.ID == myClass.ID,
		})
	}

	// Laravel sortByDesc('totalScore')：分值降序；PHP 8 排序稳定，同分保持原（id 升序）顺序。
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].TotalScore > rows[j].TotalScore })
	return rows, nil
}

// Leaderboard 教师端：以教师可管辖的第一个班级为「本班」，列出同年级排行。
func (s *PkService) Leaderboard(u *models.User) ([]PkTeacherRow, error) {
	rows := []PkTeacherRow{}

	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return rows, nil
	}

	var myClass models.ClassRoom
	err = s.db.First(&myClass, classIDs[0]).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return rows, nil
	}
	if err != nil {
		return nil, err
	}

	// 教师端与 Laravel 一致：不排除 DEMO 学校。
	gradeClasses, err := s.gradeClasses(myClass, false)
	if err != nil {
		return nil, err
	}

	for _, class := range gradeClasses {
		metrics, err := s.metricsFor(class.ID)
		if err != nil {
			return nil, err
		}
		rows = append(rows, PkTeacherRow{
			ClassID:      class.ID,
			Name:         class.Name,
			IsOwn:        class.ID == myClass.ID,
			TotalScore:   metrics.TotalScore,
			StudentCount: metrics.StudentCount,
			AvgLevel:     metrics.AvgLevel,
			PeakCount:    metrics.PeakCount,
			WeekGrowth:   metrics.WeekGrowth,
		})
	}

	sort.SliceStable(rows, func(i, j int) bool { return rows[i].TotalScore > rows[j].TotalScore })
	return rows, nil
}

// MyStats 本班 PK 统计（含同年级内按总积分的排名；无可管辖班级时全 0）。
func (s *PkService) MyStats(u *models.User) (*PkStats, error) {
	stats := &PkStats{}

	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return stats, nil
	}

	metrics, err := s.metricsFor(classIDs[0])
	if err != nil {
		return nil, err
	}
	stats.TotalScore = metrics.TotalScore
	stats.StudentCount = metrics.StudentCount
	stats.AvgLevel = metrics.AvgLevel
	stats.PeakCount = metrics.PeakCount
	stats.WeekGrowth = metrics.WeekGrowth

	var myClass models.ClassRoom
	err = s.db.First(&myClass, classIDs[0]).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return stats, nil
	}
	if err != nil {
		return nil, err
	}

	// 排名：同年级各班总积分降序（PHP arsort，PHP 8 排序稳定 → 同分保持班级 id 升序），
	// 本班的下标 + 1 即名次（未命中保持 0）。
	gradeClasses, err := s.gradeClasses(myClass, false)
	if err != nil {
		return nil, err
	}
	type classScore struct {
		classID uint
		score   int
	}
	scores := make([]classScore, 0, len(gradeClasses))
	for _, class := range gradeClasses {
		var sum int64
		if err := s.db.Model(&models.Student{}).
			Where("class_id = ? AND status = ?", class.ID, "active").
			Select("COALESCE(SUM(total_score), 0)").Scan(&sum).Error; err != nil {
			return nil, err
		}
		scores = append(scores, classScore{classID: class.ID, score: int(sum)})
	}
	sort.SliceStable(scores, func(i, j int) bool { return scores[i].score > scores[j].score })
	for i, row := range scores {
		if row.classID == myClass.ID {
			stats.Rank = i + 1
			break
		}
	}

	return stats, nil
}

// Challenge 发起 PK 挑战：校验目标后返回 {target_class, expires_at}（挑战记录不落库，见文件头说明）。
func (s *PkService) Challenge(u *models.User, targetClassID uint) (*PkChallengeResult, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	var myClassID uint
	if len(classIDs) > 0 {
		myClassID = classIDs[0]
	}
	if myClassID == 0 || myClassID == targetClassID {
		return nil, ErrBadRequest("无效的挑战目标")
	}

	var target models.ClassRoom
	err = s.db.First(&target, targetClassID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("目标班级不存在")
	}
	if err != nil {
		return nil, err
	}

	return &PkChallengeResult{
		Message: "🚀 挑战已发起！",
		Data: PkChallengeData{
			TargetClass: target.Name,
			// Laravel now()->addDays(7)->toDateTimeString() → "2006-01-02 15:04:05"。
			ExpiresAt: util.Now().AddDate(0, 0, pkChallengeDays).Format("2006-01-02 15:04:05"),
		},
	}, nil
}

// metricsFor 计算单个班级的 PK 指标（总积分/人数/均级/巅峰数/本周增长）。
func (s *PkService) metricsFor(classID uint) (*PkMetrics, error) {
	var students []models.Student
	if err := s.db.Where("class_id = ? AND status = ?", classID, "active").
		Order("id ASC").Find(&students).Error; err != nil {
		return nil, err
	}

	metrics := &PkMetrics{StudentCount: len(students)}
	ids := make([]uint, 0, len(students))
	for _, st := range students {
		metrics.TotalScore += st.TotalScore
		ids = append(ids, st.ID)
	}

	levels := map[uint]int{}
	if len(ids) > 0 {
		var pets []models.Pet
		if err := s.db.Where("student_id IN ?", ids).Find(&pets).Error; err != nil {
			return nil, err
		}
		for _, pet := range pets {
			levels[pet.StudentID] = pet.Level
		}
	}

	levelSum := 0
	for _, st := range students {
		level := levels[st.ID] // 无宠物 = 0 级（同 Laravel `$s->pet->level ?? 0`）
		levelSum += level
		if level >= pkPeakPetLevel {
			metrics.PeakCount++
		}
	}
	if len(students) > 0 {
		// Laravel round($avgLevel, 1)：四舍五入到 1 位小数（PHP 半值远离零，同 math.Round）。
		metrics.AvgLevel = math.Round(float64(levelSum)/float64(len(students))*10) / 10
	}

	if len(ids) > 0 {
		var weekly int64
		if err := s.db.Model(&models.Score{}).
			Where("student_id IN ? AND created_at >= ?", ids, util.StartOfWeek(util.Now())).
			Select("COALESCE(SUM(amount), 0)").Scan(&weekly).Error; err != nil {
			return nil, err
		}
		metrics.WeekGrowth = int(weekly)
	}

	return metrics, nil
}

// gradeClasses 同年级班级（status = active，按 id ASC 保证顺序可复现）。
// excludeDemo 为真时排除 code = 'DEMO' 的学校（教室端；Laravel whereHas school code != 'DEMO'）。
//
// 年级比较：用 COALESCE(grade, 空串) 把 NULL 与空串视为同一档 —— 本 schema 的 grade 是
// 非空字符串列（无值即空串），对应 Laravel 的 null 年级。
func (s *PkService) gradeClasses(myClass models.ClassRoom, excludeDemo bool) ([]models.ClassRoom, error) {
	q := s.db.Where("COALESCE(grade, '') = ? AND status = ?", myClass.Grade, "active")
	if excludeDemo {
		schools := s.db.Model(&models.School{}).Select("id").Where("code <> ?", "DEMO")
		q = q.Where("school_id IN (?)", schools)
	}

	var classes []models.ClassRoom
	if err := q.Order("id ASC").Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}
