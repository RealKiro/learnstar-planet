// 教师端报表服务：积分趋势、宠物等级分布、学生进度。
// 忠实移植自 Laravel App\Services\ReportService（scoreTrend / petDistribution / studentProgress），
// 所有查询都经 Scope（等价 TeacherClassScope）限定班级作用域。
package services

import (
	"errors"
	"sort"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// reportMaxTrendDays 积分趋势最多回溯天数（同 Laravel clamp 上限）。
const reportMaxTrendDays = 365

// reportProgressHistoryLimit 单人进度返回的历史条数（近 50 条）。
const reportProgressHistoryLimit = 50

// reportClassProgressLimit 全班进度每人取的最近积分条数（近 10 条）。
const reportClassProgressLimit = 10

// reportTrendWindow 涨跌趋势判定窗口（前 5 条 vs 之后 5 条）。
const reportTrendWindow = 5

// ReportService 教师端报表服务。
type ReportService struct {
	db    *gorm.DB
	scope *Scope
}

// NewReportService 创建报表服务。
func NewReportService(db *gorm.DB, scope *Scope) *ReportService {
	return &ReportService{db: db, scope: scope}
}

// ScoreTrendSeries 折线图单条序列。
type ScoreTrendSeries struct {
	Label string `json:"label"`
	Data  []int  `json:"data"`
}

// ScoreTrendData 积分趋势（labels + 双序列：得分 / 扣分）。
type ScoreTrendData struct {
	Labels   []string           `json:"labels"`
	Datasets []ScoreTrendSeries `json:"datasets"`
}

// ScoreTrend 近 N 天得分 / 扣分日趋势（双序列，供前端折线图）。
//
// 口径同 Laravel：days 收敛到 1..365；起点为今日零点往前 days-1 天；
// labels 形如 "09/13"；得分取当日 amount > 0 之和，扣分取当日 amount < 0 之和的绝对值。
func (s *ReportService) ScoreTrend(u *models.User, days int) (*ScoreTrendData, error) {
	if days < 1 {
		days = 1
	}
	if days > reportMaxTrendDays {
		days = reportMaxTrendDays
	}

	start := util.StartOfDay(util.Now()).AddDate(0, 0, -(days - 1))

	positive := map[string]int{}
	negative := map[string]int{}

	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) > 0 {
		var rows []models.Score
		if err := s.db.Select("amount", "created_at").
			Where("class_id IN ? AND created_at >= ?", classIDs, start).
			Find(&rows).Error; err != nil {
			return nil, err
		}
		for i := range rows {
			amount := rows[i].Amount
			if amount == 0 {
				continue
			}
			key := rows[i].CreatedAt.In(util.Loc).Format("2006-01-02")
			if amount > 0 {
				positive[key] += amount
			} else {
				negative[key] += -amount
			}
		}
	}

	result := &ScoreTrendData{
		Labels: []string{},
		Datasets: []ScoreTrendSeries{
			{Label: "得分", Data: []int{}},
			{Label: "扣分", Data: []int{}},
		},
	}
	for d := 0; d < days; d++ {
		day := start.AddDate(0, 0, d)
		key := day.Format("2006-01-02")
		result.Labels = append(result.Labels, day.Format("01/02"))
		result.Datasets[0].Data = append(result.Datasets[0].Data, positive[key])
		result.Datasets[1].Data = append(result.Datasets[1].Data, negative[key])
	}
	return result, nil
}

// PetLevelDistributionRow 宠物等级分布的一行。
type PetLevelDistributionRow struct {
	Level     int    `json:"level"`
	Count     int    `json:"count"`
	StageName string `json:"stage_name"`
}

// PetDistribution 宠物等级分布（按 level 升序分组，含阶段名）。
func (s *ReportService) PetDistribution(u *models.User) ([]PetLevelDistributionRow, error) {
	rows := []PetLevelDistributionRow{}

	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return rows, nil
	}

	var pets []models.Pet
	if err := s.db.Where("class_id IN ?", classIDs).Order("id ASC").Find(&pets).Error; err != nil {
		return nil, err
	}

	byLevel := map[int]*PetLevelDistributionRow{}
	levels := make([]int, 0, len(pets))
	for i := range pets {
		pet := &pets[i]
		if row, ok := byLevel[pet.Level]; ok {
			row.Count++
			continue
		}
		stage := pet.CurrentStage()
		byLevel[pet.Level] = &PetLevelDistributionRow{
			Level:     pet.Level,
			Count:     1,
			StageName: stage.Name,
		}
		levels = append(levels, pet.Level)
	}

	sort.Ints(levels)
	for _, level := range levels {
		rows = append(rows, *byLevel[level])
	}
	return rows, nil
}

// StudentProgressStudent 单人进度里的学生基本信息。
type StudentProgressStudent struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	TotalScore int    `json:"total_score"`
}

// StudentProgressDetail 单人进度（基本信息 + 近 50 条积分历史）。
type StudentProgressDetail struct {
	Student StudentProgressStudent `json:"student"`
	History []models.Score         `json:"history"`
}

// StudentProgressRow 全班进度的一行（每人最近 10 条积分 + 涨跌趋势）。
type StudentProgressRow struct {
	StudentID   uint   `json:"student_id"`
	StudentName string `json:"student_name"`
	Scores      []int  `json:"scores"`
	Trend       string `json:"trend"`
	Change      int    `json:"change"`
}

// StudentProgress 学生进度。
//
// studentID > 0 时返回该生（必须在管辖班级内，否则 404）近 50 条历史；
// 否则返回全班 active 学生每人近 10 条积分与涨跌趋势
// （change = 前 5 条之和 - 之后 5 条之和；> 5 up，< -5 down，其余 stable）。
func (s *ReportService) StudentProgress(u *models.User, studentID uint) (any, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	if studentID > 0 {
		var student models.Student
		err := s.db.Where("id = ? AND class_id IN ?", studentID, classIDs).First(&student).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound("学生不存在或不在管辖范围")
		}
		if err != nil {
			return nil, err
		}

		history := []models.Score{}
		if err := s.db.Where("student_id = ?", student.ID).
			Order("created_at DESC, id DESC").
			Limit(reportProgressHistoryLimit).Find(&history).Error; err != nil {
			return nil, err
		}

		return &StudentProgressDetail{
			Student: StudentProgressStudent{
				ID:         student.ID,
				Name:       student.Name,
				TotalScore: student.TotalScore,
			},
			History: history,
		}, nil
	}

	rows := []StudentProgressRow{}
	if len(classIDs) == 0 {
		return rows, nil
	}

	var students []models.Student
	if err := s.db.Where("class_id IN ? AND status = ?", classIDs, "active").
		Order("id ASC").Find(&students).Error; err != nil {
		return nil, err
	}

	for _, student := range students {
		amounts := []int{}
		if err := s.db.Model(&models.Score{}).
			Where("student_id = ?", student.ID).
			Order("created_at DESC, id DESC").
			Limit(reportClassProgressLimit).
			Pluck("amount", &amounts).Error; err != nil {
			return nil, err
		}
		if amounts == nil {
			amounts = []int{}
		}

		change := 0
		for i, amount := range amounts {
			if i < reportTrendWindow {
				change += amount
			} else {
				change -= amount
			}
		}

		trend := "stable"
		switch {
		case change > 5:
			trend = "up"
		case change < -5:
			trend = "down"
		}

		rows = append(rows, StudentProgressRow{
			StudentID:   student.ID,
			StudentName: student.Name,
			Scores:      amounts,
			Trend:       trend,
			Change:      change,
		})
	}
	return rows, nil
}
