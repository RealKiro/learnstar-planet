// 班级宠物系列服务：教室端整班切换系列（每人扣 20 积分）+ 教师端班级系列切换。
//
// 忠实移植自：
//   - Laravel App\Services\PetSeriesService::switchSeries（教师端 POST teacher/class/switch-series）
//   - Laravel App\Http\Controllers\Api\DisplayController::classroomSwitchSeries（教室端 POST display/switch-series）
//
// 有意差异：
//  1. 免费自选机会：Laravel 用 Cache::put("pet_free_pick:<student_id>", 1, 3 天)，
//     Go 端无 Cache → 落库为 pet_free_picks 表（models.GrantPetFreePick），语义等价（过期即失效）。
//  2. 教室端的「每人扣 20 积分」在 Laravel 是直接改 student.total_score 后 save：
//     不写 scores / score_logs 审计、不发 score_update 大屏事件、不同步宠物经验。
//     Go 端逐条保持一致（不为原实现没有的事件「补发」），仅把整批写入放进一个事务。
//  3. 教师端没有 20 积分扣费（同 Laravel）。
//  4. 说明：「未选宠物的学生按新系列重抽」在原实现中并不存在 —— 整班切换只写班级设置 +
//     发放免费自选机会；无宠物学生的随机分配发生在 initialData（且固定取 myth 池，与班级
//     当前系列无关）、以及学生自己调 switch-pet 时。故此处不额外实现重抽。
package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// classroomSeriesSwitchCost 教室端整班切换系列的每人积分成本（同 Laravel `$costPerStudent = 20`）。
const classroomSeriesSwitchCost = 20

// PetSeriesService 班级宠物系列服务。
type PetSeriesService struct {
	db    *gorm.DB
	scope *Scope
}

// NewPetSeriesService 创建宠物系列服务。
func NewPetSeriesService(db *gorm.DB, scope *Scope) *PetSeriesService {
	return &PetSeriesService{db: db, scope: scope}
}

// ClassroomSwitchSeriesData 教室端整班切换的 data（字段名逐字同 Laravel）。
type ClassroomSwitchSeriesData struct {
	SeriesID         string `json:"series_id"`
	CostPerStudent   int    `json:"cost_per_student"`
	AffectedStudents int    `json:"affected_students"`
	FreePickGranted  bool   `json:"free_pick_granted"`
}

// TeacherSwitchSeriesData 教师端切换系列的 data（字段名逐字同 Laravel）。
type TeacherSwitchSeriesData struct {
	SeriesID        string `json:"series_id"`
	ClassID         uint   `json:"class_id"`
	FreePickGranted bool   `json:"free_pick_granted"`
	GrantedStudents int    `json:"granted_students"`
}

// ClassroomSwitchSeriesResult 教室端整班切换结果。
type ClassroomSwitchSeriesResult struct {
	Message string
	Data    ClassroomSwitchSeriesData
}

// TeacherSwitchSeriesResult 教师端切换结果。
type TeacherSwitchSeriesResult struct {
	Message string
	Data    TeacherSwitchSeriesData
}

// SwitchSeriesForClassroom 教室端整班切换系列：全班每人扣 20 积分 + 各发一次免费自选机会。
func (s *PetSeriesService) SwitchSeriesForClassroom(classID uint, seriesID string) (*ClassroomSwitchSeriesResult, error) {
	if !models.IsValidSeries(seriesID) {
		return nil, ErrUnprocessable("无效的系列ID")
	}

	class, err := s.findClass(classID)
	if err != nil {
		return nil, err
	}

	students, err := s.activeStudents(classID)
	if err != nil {
		return nil, err
	}
	if len(students) == 0 {
		return nil, ErrBadRequest("班级没有活跃学生")
	}

	// 积分检查：任一学生不足则整体拒绝（前 3 名 + 总数，文案逐字同 Laravel）。
	names := make([]string, 0, 3)
	insufficient := 0
	for _, st := range students {
		if st.TotalScore < classroomSeriesSwitchCost {
			insufficient++
			if len(names) < 3 {
				names = append(names, st.Name)
			}
		}
	}
	if insufficient > 0 {
		more := ""
		if insufficient > 3 {
			more = fmt.Sprintf("等%d人", insufficient)
		}
		return nil, ErrBadRequest(fmt.Sprintf("积分不足：%s%s 每人需要 %d 积分",
			strings.Join(names, "、"), more, classroomSeriesSwitchCost))
	}

	now := util.Now()
	err = s.db.Transaction(func(tx *gorm.DB) error {
		for _, st := range students {
			// ⚠️ 条件更新一次性完成「校验 + 扣减」：并发下不会穿仓，也不会丢更新。
			_, _, ok, err := deductScoreAtomic(tx, st.ID, classroomSeriesSwitchCost)
			if err != nil {
				return err
			}
			if !ok {
				return ErrBadRequest(fmt.Sprintf("积分不足：%s 每人需要 %d 积分", st.Name, classroomSeriesSwitchCost))
			}
		}
		if err := s.saveSeries(tx, class, seriesID); err != nil {
			return err
		}
		for _, st := range students {
			if err := models.GrantPetFreePick(tx, st.ID, now); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &ClassroomSwitchSeriesResult{
		Message: fmt.Sprintf("已切换至「%s」系列：全班 %d 人各扣 %d 积分，并各获一次免费自选该系列宠物的机会",
			seriesID, len(students), classroomSeriesSwitchCost),
		Data: ClassroomSwitchSeriesData{
			SeriesID:         seriesID,
			CostPerStudent:   classroomSeriesSwitchCost,
			AffectedStudents: len(students),
			FreePickGranted:  true,
		},
	}, nil
}

// TeacherSwitchSeries 教师端切换本班宠物系列（不扣积分，仅写设置 + 发免费自选机会）。
func (s *PetSeriesService) SwitchSeries(u *models.User, seriesID string) (*TeacherSwitchSeriesResult, error) {
	if !models.IsValidSeries(seriesID) {
		return nil, ErrUnprocessable("无效的系列ID，可选值：" + strings.Join(models.SeriesIDs, ", "))
	}

	// Laravel PetSeriesService 取「教师可管辖班级的第一个」，没有则 400。
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return nil, ErrBadRequest("没有可管理的班级")
	}

	class, err := s.findClass(classIDs[0])
	if err != nil {
		return nil, err
	}

	students, err := s.activeStudents(class.ID)
	if err != nil {
		return nil, err
	}

	now := util.Now()
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.saveSeries(tx, class, seriesID); err != nil {
			return err
		}
		for _, st := range students {
			if err := models.GrantPetFreePick(tx, st.ID, now); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &TeacherSwitchSeriesResult{
		Message: fmt.Sprintf("已切换系列为「%s」，全班 %d 人各获一次免费自选该系列宠物的机会",
			seriesID, len(students)),
		Data: TeacherSwitchSeriesData{
			SeriesID:        seriesID,
			ClassID:         class.ID,
			FreePickGranted: true,
			GrantedStudents: len(students),
		},
	}, nil
}

// saveSeries 把 pet_series 合并写入班级 settings（保留其他键，同 Laravel settings 覆盖式写入的结果）。
func (s *PetSeriesService) saveSeries(tx *gorm.DB, class *models.ClassRoom, seriesID string) error {
	settings, err := class.WithSetting("pet_series", seriesID)
	if err != nil {
		return err
	}
	if err := tx.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
		Update("settings", settings).Error; err != nil {
		return err
	}
	class.Settings = settings
	return nil
}

// findClass 取班级（不存在 → 404，同 Laravel findOrFail）。
func (s *PetSeriesService) findClass(classID uint) (*models.ClassRoom, error) {
	var class models.ClassRoom
	err := s.db.First(&class, classID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("班级不存在")
	}
	if err != nil {
		return nil, err
	}
	return &class, nil
}

// activeStudents 班级内活跃学生（按 id ASC，便于结果可复现）。
func (s *PetSeriesService) activeStudents(classID uint) ([]models.Student, error) {
	var students []models.Student
	if err := s.db.Where("class_id = ? AND status = ?", classID, "active").
		Order("id ASC").Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

// ============================================================
// 教师端班级信息卡（GET /teacher/class）
// 移植自 Laravel PetSeriesService::classInfo。
// ============================================================

// ClassInfoView 班级信息卡（字段与顺序逐字同 Laravel classInfo 的返回数组）。
type ClassInfoView struct {
	ID           uint           `json:"id"`
	Name         string         `json:"name"`
	Grade        string         `json:"grade"`
	StudentCount int64          `json:"student_count"`
	TotalScore   int            `json:"total_score"`
	ClassPoints  int            `json:"class_points"`
	Settings     map[string]any `json:"settings"`
	DisplayCode  string         `json:"display_code"`
}

// ClassInfo 当前（首个可管辖）班级的信息卡；没有可管理的班级 → 400「没有可管理的班级」
// （Laravel 抛 DomainException，控制器映射 400）。
func (s *PetSeriesService) ClassInfo(u *models.User) (*ClassInfoView, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return nil, ErrBadRequest("没有可管理的班级")
	}

	class, err := s.findClass(classIDs[0])
	if err != nil {
		return nil, err
	}

	var studentCount int64
	if err := s.db.Model(&models.Student{}).
		Where("class_id = ? AND status = ?", class.ID, "active").
		Count(&studentCount).Error; err != nil {
		return nil, err
	}

	var totalScore int
	if err := s.db.Model(&models.Student{}).
		Where("class_id = ? AND status = ?", class.ID, "active").
		Select("COALESCE(SUM(total_score), 0)").Scan(&totalScore).Error; err != nil {
		return nil, err
	}

	// class_points：settings.class_points ?? 0（Laravel 直接 (int) 转换）。
	settings := class.SettingsMap()
	classPoints := 0
	switch v := settings["class_points"].(type) {
	case float64:
		classPoints = int(v)
	case int:
		classPoints = v
	}
	// 空 settings 在 Laravel 的 JSON 里是 null，这里对齐为 nil。
	if len(settings) == 0 {
		settings = nil
	}

	return &ClassInfoView{
		ID:           class.ID,
		Name:         class.Name,
		Grade:        class.Grade,
		StudentCount: studentCount,
		TotalScore:   totalScore,
		ClassPoints:  classPoints,
		Settings:     settings,
		DisplayCode:  NewDisplayCodeService(s.db).Generate(class, nil),
	}, nil
}
