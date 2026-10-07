// 管理端服务：学校信息、班级 CRUD、学生 CRUD、教师 CRUD / 重置密码。
// 所有操作均以 schoolID 为边界，防止管理员跨校操作。
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Admin 管理端业务服务。
type Admin struct {
	db *gorm.DB
}

// NewAdmin 创建管理端服务。
func NewAdmin(db *gorm.DB) *Admin {
	return &Admin{db: db}
}

// School 返回指定学校信息。
func (a *Admin) School(schoolID uint) (*models.School, error) {
	var school models.School
	if err := a.db.First(&school, schoolID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound("学校不存在")
		}
		return nil, err
	}
	return &school, nil
}

// SchoolView 学校对外视图：字段同 Laravel 的 School JSON，并把 `settings` 文本列
// 解析为**对象**输出（Laravel `schools.settings` 是 array cast，JSON 里就是对象）。
//
// 为什么需要它：Go 的 `models.School.Settings` 带 `json:"-"`（同为模型层的
// `SettingString`/`SettingsMap` 读取），直接把模型序列化会丢掉 settings，
// 前端「学校设置」页就读不到 `third_party_platform` / `enabled_third_party_platforms`。
type SchoolView struct {
	ID               uint           `json:"id"`
	Name             string         `json:"name"`
	Code             string         `json:"code"`
	Address          string         `json:"address"`
	ContactPhone     string         `json:"contact_phone"`
	ContactEmail     string         `json:"contact_email"`
	LogoPath         string         `json:"logo_path"`
	Status           string         `json:"status"`
	ScoreRulesSeeded bool           `json:"score_rules_seeded"`
	Settings         map[string]any `json:"settings"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// SchoolViewOf 构造学校对外视图；settings 为空时输出 `{}`（Laravel 为 `null`，
// 前端两处读取都写作 `s.settings || {}`，取 `{}` 更省心且不改变语义）。
func SchoolViewOf(school *models.School) SchoolView {
	settings := school.SettingsMap()
	if settings == nil {
		settings = map[string]any{}
	}
	return SchoolView{
		ID:               school.ID,
		Name:             school.Name,
		Code:             school.Code,
		Address:          school.Address,
		ContactPhone:     school.ContactPhone,
		ContactEmail:     school.ContactEmail,
		LogoPath:         school.LogoPath,
		Status:           school.Status,
		ScoreRulesSeeded: school.ScoreRulesSeeded,
		Settings:         settings,
		CreatedAt:        school.CreatedAt,
		UpdatedAt:        school.UpdatedAt,
	}
}

// UpdateSchool 更新学校可编辑字段并返回最新实体。
//
// settings 为 nil 表示请求体里**未提供**该键 → 保持原值；
// 非 nil 表示**整体替换** settings 列（同 Laravel `$school->fill($request->only([...,'settings']))`，
// 是替换而不是按键合并——前端提交的是它读回来的整份 settings，故该语义与 Laravel 一致）。
func (a *Admin) UpdateSchool(schoolID uint, updates map[string]any, settings *map[string]any) (*models.School, error) {
	if settings != nil {
		encoded, err := json.Marshal(*settings)
		if err != nil {
			return nil, err
		}
		updates["settings"] = string(encoded)
	}
	if len(updates) > 0 {
		if err := a.db.Model(&models.School{}).Where("id = ?", schoolID).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return a.School(schoolID)
}

// ListClasses 返回学校下全部班级（含已停用，按 ID 升序）。
func (a *Admin) ListClasses(schoolID uint) ([]models.ClassRoom, error) {
	var classes []models.ClassRoom
	if err := a.db.Where("school_id = ?", schoolID).Order("id ASC").Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

// CreateClass 创建班级；teacherID 若提供则校验为本校教师。
func (a *Admin) CreateClass(schoolID uint, name, grade, year string, teacherID *uint, maxStudents int) (*models.ClassRoom, error) {
	if name == "" {
		return nil, ErrUnprocessable("班级名称不能为空")
	}
	if teacherID != nil {
		if err := a.ensureTeacherInSchool(schoolID, *teacherID); err != nil {
			return nil, err
		}
	}

	class := models.ClassRoom{
		SchoolID:    schoolID,
		Name:        name,
		Grade:       grade,
		Year:        year,
		TeacherID:   teacherID,
		MaxStudents: maxStudents,
		Status:      "active",
	}
	if err := a.db.Create(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

// UpdateClass 更新班级可编辑字段并返回最新实体。
func (a *Admin) UpdateClass(schoolID, classID uint, updates map[string]any) (*models.ClassRoom, error) {
	return a.UpdateClassWithSettings(schoolID, classID, updates, nil)
}

// classRoomSeries 班级宠物系列的可写值（逐字同 Laravel `update` 的 `in:cosmic,pokemon,cute,treasure,mythic,all`）。
// 注意：这些是**旧系列 id**，`SpeciesPoolForSeries` 对它们返回空池 → 换宠时视为不限制（同 Laravel 注释）。
var classRoomSeries = []string{"all", "cosmic", "pokemon", "cute", "treasure", "mythic"}

// ValidClassRoomSeries 判断是否为合法班级宠物系列值。
func ValidClassRoomSeries(value string) bool {
	for _, allowed := range classRoomSeries {
		if value == allowed {
			return true
		}
	}
	return false
}

// UpdateClassWithSettings 更新班级可编辑字段，并可**合并** settings（Laravel `update`：
// 先 `fill($request->only([...]))`，再在 `$request->has('pet_series')` 时把该键并进现有 settings 后整体保存）。
//
// settingsPatch 为 nil 表示不动 settings；非 nil 时与现有 settings **按键合并**（保留其它键）。
func (a *Admin) UpdateClassWithSettings(schoolID, classID uint, updates map[string]any, settingsPatch map[string]any) (*models.ClassRoom, error) {
	class, err := a.classInSchool(schoolID, classID)
	if err != nil {
		return nil, err
	}
	if v, ok := updates["teacher_id"]; ok && v != nil {
		if id, ok := v.(uint); ok {
			if err := a.ensureTeacherInSchool(schoolID, id); err != nil {
				return nil, err
			}
		}
	}
	if len(settingsPatch) > 0 {
		merged := class.SettingsMap()
		for key, value := range settingsPatch {
			merged[key] = value
		}
		encoded, err := json.Marshal(merged)
		if err != nil {
			return nil, err
		}
		updates["settings"] = string(encoded)
	}
	if len(updates) > 0 {
		if err := a.db.Model(&models.ClassRoom{}).Where("id = ?", classID).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return a.classInSchool(schoolID, classID)
}

// ClassRoomView 班级对外视图：字段同 Laravel 的 ClassRoom JSON，并把 `settings` 文本列解析为**对象**输出
// （Laravel `class_rooms.settings` 是 array cast）。前端班级管理页读 `settings.pet_series`。
type ClassRoomView struct {
	ID                   uint           `json:"id"`
	SchoolID             uint           `json:"school_id"`
	Name                 string         `json:"name"`
	Grade                string         `json:"grade"`
	Year                 string         `json:"year"`
	TeacherID            *uint          `json:"teacher_id"`
	MaxStudents          int            `json:"max_students"`
	DisplayCode          string         `json:"display_code"`
	DisplayCodeUpdatedAt *time.Time     `json:"display_code_updated_at"`
	Settings             map[string]any `json:"settings"`
	Status               string         `json:"status"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

// ClassRoomViewOf 构造班级对外视图（settings 为空时输出 `{}`）。
func ClassRoomViewOf(class *models.ClassRoom) ClassRoomView {
	settings := class.SettingsMap()
	if settings == nil {
		settings = map[string]any{}
	}
	return ClassRoomView{
		ID:                   class.ID,
		SchoolID:             class.SchoolID,
		Name:                 class.Name,
		Grade:                class.Grade,
		Year:                 class.Year,
		TeacherID:            class.TeacherID,
		MaxStudents:          class.MaxStudents,
		DisplayCode:          class.DisplayCode,
		DisplayCodeUpdatedAt: class.DisplayCodeUpdatedAt,
		Settings:             settings,
		Status:               class.Status,
		CreatedAt:            class.CreatedAt,
		UpdatedAt:            class.UpdatedAt,
	}
}

// DeleteClass 删除班级（软删除依赖数据库支持；若存在学生则拒绝，避免孤儿数据）。
func (a *Admin) DeleteClass(schoolID, classID uint) error {
	if _, err := a.classInSchool(schoolID, classID); err != nil {
		return err
	}
	var studentCount int64
	if err := a.db.Model(&models.Student{}).Where("class_id = ?", classID).Count(&studentCount).Error; err != nil {
		return err
	}
	if studentCount > 0 {
		return ErrUnprocessable(fmt.Sprintf("班级内仍有 %d 名学生，请先清空或转班", studentCount))
	}
	return a.db.Where("id = ? AND school_id = ?", classID, schoolID).Delete(&models.ClassRoom{}).Error
}

// CreateStudent 创建学生；同班同名跳过（返回已有记录），同校同学号拦截。
func (a *Admin) CreateStudent(schoolID, classID uint, name, gender, studentNo string) (*models.Student, error) {
	if name == "" {
		return nil, ErrUnprocessable("学生姓名不能为空")
	}
	if _, err := a.classInSchool(schoolID, classID); err != nil {
		return nil, err
	}

	// 同班同名跳过（与 Laravel 导入口径一致）。
	var existing models.Student
	err := a.db.Where("class_id = ? AND name = ?", classID, name).First(&existing).Error
	if err == nil {
		return &existing, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 同学号跨班拦截。
	if studentNo != "" {
		var count int64
		if err := a.db.Model(&models.Student{}).
			Where("student_no = ? AND class_id IN (?)", studentNo,
				a.db.Model(&models.ClassRoom{}).Select("id").Where("school_id = ?", schoolID)).
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, ErrUnprocessable(fmt.Sprintf("学号 %s 已被本校其他学生使用", studentNo))
		}
	}

	student := models.Student{
		ClassID:   classID,
		Name:      name,
		Gender:    gender,
		StudentNo: studentNo,
		Status:    "active",
	}
	if err := a.db.Create(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

// UpdateStudent 更新学生可编辑字段并返回最新实体。
func (a *Admin) UpdateStudent(schoolID, studentID uint, updates map[string]any) (*models.Student, error) {
	if _, err := a.studentInSchool(schoolID, studentID); err != nil {
		return nil, err
	}
	if len(updates) > 0 {
		if err := a.db.Model(&models.Student{}).Where("id = ?", studentID).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return a.studentInSchool(schoolID, studentID)
}

// DeleteStudent 删除学生（软删除，同时清理其积分审计与宠物）。
func (a *Admin) DeleteStudent(schoolID, studentID uint) error {
	student, err := a.studentInSchool(schoolID, studentID)
	if err != nil {
		return err
	}
	return a.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("student_id = ?", student.ID).Delete(&models.Pet{}).Error; err != nil {
			return err
		}
		if err := tx.Where("student_id = ?", student.ID).Delete(&models.ScoreLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("student_id = ?", student.ID).Delete(&models.Score{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Student{}, student.ID).Error
	})
}

// ListTeachers 返回学校下全部教师（role=teacher）。
func (a *Admin) ListTeachers(schoolID uint) ([]models.User, error) {
	var teachers []models.User
	if err := a.db.Where("school_id = ? AND role = ?", schoolID, "teacher").
		Order("id ASC").Find(&teachers).Error; err != nil {
		return nil, err
	}
	return teachers, nil
}

// CreateTeacher 创建教师账号（用户名全校唯一）。
func (a *Admin) CreateTeacher(schoolID uint, username, password, name, phone string) (*models.User, error) {
	if username == "" || password == "" || name == "" {
		return nil, ErrUnprocessable("用户名、密码、姓名为必填项")
	}
	if len(password) < 6 {
		return nil, ErrUnprocessable("密码至少 6 位")
	}

	var count int64
	if err := a.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUnprocessable("用户名已存在")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	teacher := models.User{
		SchoolID:      schoolID,
		Role:          "teacher",
		Username:      username,
		PasswordHash:  string(hash),
		PlainPassword: password,
		Name:          name,
		Nickname:      name,
		Phone:         phone,
		Status:        "active",
	}
	if err := a.db.Create(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

// UpdateTeacher 更新教师可编辑字段并返回最新实体。
func (a *Admin) UpdateTeacher(schoolID, userID uint, updates map[string]any) (*models.User, error) {
	teacher, err := a.teacherInSchool(schoolID, userID)
	if err != nil {
		return nil, err
	}
	if len(updates) > 0 {
		if err := a.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return a.teacherInSchool(schoolID, teacher.ID)
}

// ResetPassword 重置教师密码。
func (a *Admin) ResetPassword(schoolID, userID uint, newPassword string) error {
	if _, err := a.teacherInSchool(schoolID, userID); err != nil {
		return err
	}
	if len(newPassword) < 6 {
		return ErrUnprocessable("密码至少 6 位")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return a.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]any{
		"password_hash":    string(hash),
		"plain_password":   newPassword,
		"password_changed": false,
	}).Error
}

// classInSchool 校验班级归属本校。
func (a *Admin) classInSchool(schoolID, classID uint) (*models.ClassRoom, error) {
	var class models.ClassRoom
	err := a.db.Where("id = ? AND school_id = ?", classID, schoolID).First(&class).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("班级不存在")
	}
	if err != nil {
		return nil, err
	}
	return &class, nil
}

// studentInSchool 校验学生归属本校。
func (a *Admin) studentInSchool(schoolID, studentID uint) (*models.Student, error) {
	var student models.Student
	err := a.db.Where("id = ? AND class_id IN (?)", studentID,
		a.db.Model(&models.ClassRoom{}).Select("id").Where("school_id = ?", schoolID)).
		First(&student).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("学生不存在")
	}
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// teacherInSchool 校验教师归属本校且角色为 teacher。
func (a *Admin) teacherInSchool(schoolID, userID uint) (*models.User, error) {
	var teacher models.User
	err := a.db.Where("id = ? AND school_id = ? AND role = ?", userID, schoolID, "teacher").
		First(&teacher).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("教师不存在")
	}
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// ensureTeacherInSchool 校验指定 ID 是否为该校教师（用作班主任）。
func (a *Admin) ensureTeacherInSchool(schoolID, userID uint) error {
	_, err := a.teacherInSchool(schoolID, userID)
	if err != nil {
		return ErrUnprocessable("班主任不存在或非本校教师")
	}
	return nil
}
