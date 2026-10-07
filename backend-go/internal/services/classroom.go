// 教师端「我的班级 / 切换班级 / 模式」服务。
// 忠实移植自 Laravel App\Services\DashboardService(classesFor/switchTo) 与
// App\Services\ClassroomMessagingService(getMode/setMode)。
//
// 口径差异说明：Laravel 用 class_room_teachers 关联表承载多教师多角色
// （head_teacher/co_teacher/subject_teacher/grade_lead/admin_director），
// Go 端重设计 schema 后只保留 class_rooms.teacher_id（班主任）这一条关联，
// 故非机器人教师的 role 统一返回 head_teacher；API 机器人仍返回 api_bot。
package services

import (
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// 教师端模式取值（classroom_display=班级大屏，teacher_manage=教师管理）。
const (
	ModeClassroomDisplay = "classroom_display"
	ModeTeacherManage    = "teacher_manage"
)

// ActiveClassSettingKey 当前激活班级的用户设置键。
const ActiveClassSettingKey = "active_class_id"

// DisplayModeSettingKey 当前模式的用户设置键。
const DisplayModeSettingKey = "display_mode"

// ClassroomService 教师端班级与模式服务。
type ClassroomService struct {
	db    *gorm.DB
	scope *Scope
}

// NewClassroomService 创建班级/模式服务。
func NewClassroomService(db *gorm.DB, scope *Scope) *ClassroomService {
	return &ClassroomService{db: db, scope: scope}
}

// MyClassView 我的班级条目。
type MyClassView struct {
	ClassID   uint   `json:"class_id"`
	ClassName string `json:"class_name"`
	Grade     string `json:"grade"`
	Role      string `json:"role"`
}

// MyClasses 教师关联班级列表。
// API 机器人返回本校全部启用班级（供外部系统枚举），普通教师返回本人担任班主任的班级。
func (s *ClassroomService) MyClasses(u *models.User) ([]MyClassView, error) {
	classes, err := s.scope.Classes(u)
	if err != nil {
		return nil, err
	}

	role := "head_teacher"
	if u.IsAPIBot {
		role = "api_bot"
	}

	views := []MyClassView{}
	for _, c := range classes {
		views = append(views, MyClassView{
			ClassID:   c.ID,
			ClassName: c.Name,
			Grade:     c.Grade,
			Role:      role,
		})
	}
	return views, nil
}

// SwitchTo 切换当前激活班级；未分配该班级返回 403。
func (s *ClassroomService) SwitchTo(u *models.User, classID uint) error {
	if err := s.scope.ClassInScope(u, classID); err != nil {
		// 越权（含班级不存在）统一映射为 403，与 Laravel DashboardService::switchTo 契约一致。
		return ErrForbidden("您未被分配到此班级")
	}
	return s.setSetting(u, ActiveClassSettingKey, classID)
}

// ModeView 当前模式与激活班级。
type ModeView struct {
	Mode          string `json:"mode"`
	ActiveClassID *uint  `json:"active_class_id"`
}

// GetMode 返回当前模式（默认 classroom_display）与激活班级。
func (s *ClassroomService) GetMode(u *models.User) (*ModeView, error) {
	return &ModeView{
		Mode:          u.SettingString(DisplayModeSettingKey, ModeClassroomDisplay),
		ActiveClassID: settingUintPtr(u.SettingUint(ActiveClassSettingKey)),
	}, nil
}

// SetModeResult 切换模式结果。
type SetModeResult struct {
	Message string   `json:"message"`
	Data    ModeView `json:"data"`
}

// SetMode 切换模式；附带 classID 且已分配时同步切换激活班级。
// 校验口径与 Laravel 控制器一致：mode 非法 → 422；非 classroom_display 时 password 必填 → 422
// （Laravel 只校验存在性、不校验密码本身，此处保持一致）。
func (s *ClassroomService) SetMode(u *models.User, mode string, classID uint, password string) (*SetModeResult, error) {
	if mode != ModeClassroomDisplay && mode != ModeTeacherManage {
		return nil, ErrUnprocessable("模式不合法")
	}
	if mode != ModeClassroomDisplay && password == "" {
		return nil, ErrUnprocessable("切换该模式需要输入密码")
	}

	if err := s.setSetting(u, DisplayModeSettingKey, mode); err != nil {
		return nil, err
	}

	if classID > 0 {
		// 已分配则同步切换激活班级；未分配则静默忽略（同 Laravel）。
		if err := s.scope.ClassInScope(u, classID); err == nil {
			if err := s.setSetting(u, ActiveClassSettingKey, classID); err != nil {
				return nil, err
			}
		}
	}

	label := "班级大屏"
	if mode == ModeTeacherManage {
		label = "教师管理"
	}

	return &SetModeResult{
		Message: "已切换为" + label + "模式",
		Data: ModeView{
			Mode:          mode,
			ActiveClassID: settingUintPtr(u.SettingUint(ActiveClassSettingKey)),
		},
	}, nil
}

// setSetting 合并写入单个用户设置项（只更新 settings 列）。
func (s *ClassroomService) setSetting(u *models.User, key string, value any) error {
	raw, err := u.WithSetting(key, value)
	if err != nil {
		return err
	}
	if err := s.db.Model(&models.User{}).Where("id = ?", u.ID).Update("settings", raw).Error; err != nil {
		return err
	}
	u.Settings = raw
	return nil
}

func settingUintPtr(v uint) *uint {
	if v == 0 {
		return nil
	}
	return &v
}
