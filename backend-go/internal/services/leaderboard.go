// 排行榜服务：MVP 直接查库（总分/本周/宠物等级），不再引入 Redis。
package services

import (
	"sort"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// LeaderboardService 排行榜服务。
type LeaderboardService struct {
	db *gorm.DB
}

// NewLeaderboardService 创建排行榜服务。
func NewLeaderboardService(db *gorm.DB) *LeaderboardService {
	return &LeaderboardService{db: db}
}

// LeaderboardEntry 排行榜条目。
type LeaderboardEntry struct {
	Rank      int    `json:"rank"`
	StudentID uint   `json:"student_id"`
	Name      string `json:"name"`
	StudentNo string `json:"student_no"`
	Avatar    string `json:"avatar"`
	Score     int    `json:"score"`
	PetEmoji  string `json:"pet_emoji,omitempty"`
	PetName   string `json:"pet_name,omitempty"`
	PetLevel  int    `json:"pet_level,omitempty"`
}

// Total 班级总分排行榜。
func (l *LeaderboardService) Total(classID uint, limit int) ([]LeaderboardEntry, error) {
	limit = normalizeLimit(limit)

	var students []models.Student
	if err := l.db.Where("class_id = ? AND status = ?", classID, "active").
		Order("total_score DESC, id ASC").Limit(limit).Find(&students).Error; err != nil {
		return nil, err
	}

	out := make([]LeaderboardEntry, 0, len(students))
	for i, s := range students {
		out = append(out, LeaderboardEntry{
			Rank:      i + 1,
			StudentID: s.ID,
			Name:      s.Name,
			StudentNo: s.StudentNo,
			Avatar:    s.AvatarPath,
			Score:     s.TotalScore,
		})
	}
	return out, nil
}

// Weekly 班级本周积分排行榜（按本周新增积分求和排序）。
func (l *LeaderboardService) Weekly(classID uint, limit int) ([]LeaderboardEntry, error) {
	limit = normalizeLimit(limit)
	weekStart := util.StartOfWeek(util.Now())

	var students []models.Student
	if err := l.db.Where("class_id = ? AND status = ?", classID, "active").
		Find(&students).Error; err != nil {
		return nil, err
	}

	type row struct {
		StudentID uint
		Sum       int64
	}
	var rows []row
	if err := l.db.Model(&models.Score{}).
		Select("student_id, COALESCE(SUM(amount), 0) AS sum").
		Where("class_id = ? AND created_at >= ?", classID, weekStart).
		Group("student_id").Scan(&rows).Error; err != nil {
		return nil, err
	}

	weekly := make(map[uint]int, len(rows))
	for _, r := range rows {
		weekly[r.StudentID] = int(r.Sum)
	}

	sort.SliceStable(students, func(i, j int) bool {
		if weekly[students[i].ID] != weekly[students[j].ID] {
			return weekly[students[i].ID] > weekly[students[j].ID]
		}
		return students[i].ID < students[j].ID
	})

	out := make([]LeaderboardEntry, 0, limit)
	for i, s := range students {
		if i >= limit {
			break
		}
		out = append(out, LeaderboardEntry{
			Rank:      i + 1,
			StudentID: s.ID,
			Name:      s.Name,
			StudentNo: s.StudentNo,
			Avatar:    s.AvatarPath,
			Score:     weekly[s.ID],
		})
	}
	return out, nil
}

// PetLevel 班级宠物等级排行榜。
func (l *LeaderboardService) PetLevel(classID uint, limit int) ([]LeaderboardEntry, error) {
	limit = normalizeLimit(limit)

	var pets []models.Pet
	if err := l.db.Where("class_id = ?", classID).
		Order("level DESC, id ASC").Limit(limit).Find(&pets).Error; err != nil {
		return nil, err
	}

	students := l.studentsByPet(pets)

	out := make([]LeaderboardEntry, 0, len(pets))
	for i, pet := range pets {
		stage := pet.CurrentStage()
		st, ok := students[pet.StudentID]
		name, no, avatar := "", "", ""
		if ok {
			name, no, avatar = st.Name, st.StudentNo, st.AvatarPath
		}
		out = append(out, LeaderboardEntry{
			Rank:      i + 1,
			StudentID: pet.StudentID,
			Name:      name,
			StudentNo: no,
			Avatar:    avatar,
			Score:     pet.Level,
			PetEmoji:  stage.Emoji,
			PetName:   stage.Name,
			PetLevel:  pet.Level,
		})
	}
	return out, nil
}

func (l *LeaderboardService) studentsByPet(pets []models.Pet) map[uint]models.Student {
	ids := make([]uint, 0, len(pets))
	for _, p := range pets {
		ids = append(ids, p.StudentID)
	}

	out := map[uint]models.Student{}
	if len(ids) == 0 {
		return out
	}

	var students []models.Student
	if err := l.db.Where("id IN ?", ids).Find(&students).Error; err != nil {
		return out
	}
	for _, s := range students {
		out[s.ID] = s
	}
	return out
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}
