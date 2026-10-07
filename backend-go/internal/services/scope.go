// 教师班级范围服务：限定教师只能操作本校可管辖的班级。
package services

import (
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// Scope 计算教师可管辖的班级范围。
type Scope struct {
	db *gorm.DB
}

// NewScope 创建范围服务。
func NewScope(db *gorm.DB) *Scope {
	return &Scope{db: db}
}

// ClassIDs 返回教师可管辖的班级 ID 列表。
// API 机器人（IsAPIBot）放行本校全部班级；普通教师仅其担任班主任的班级。
func (s *Scope) ClassIDs(u *models.User) ([]uint, error) {
	q := s.db.Model(&models.ClassRoom{}).Where("school_id = ? AND status = ?", u.SchoolID, "active")

	if !u.IsAPIBot {
		q = q.Where("teacher_id = ?", u.ID)
	}

	var ids []uint
	if err := q.Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// Classes 返回教师可管辖的班级实体列表。
func (s *Scope) Classes(u *models.User) ([]models.ClassRoom, error) {
	q := s.db.Where("school_id = ? AND status = ?", u.SchoolID, "active")

	if !u.IsAPIBot {
		q = q.Where("teacher_id = ?", u.ID)
	}

	var classes []models.ClassRoom
	if err := q.Order("id ASC").Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

// StudentInScope 判断学生是否在教师管辖范围内，是则返回学生。
func (s *Scope) StudentInScope(u *models.User, studentID uint) (*models.Student, error) {
	classIDs, err := s.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	var student models.Student
	err = s.db.Where("id = ? AND class_id IN ?", studentID, classIDs).First(&student).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound("学生不存在或不在管辖范围")
		}
		return nil, err
	}
	return &student, nil
}

// ClassInScope 校验班级是否在教师管辖范围内。
func (s *Scope) ClassInScope(u *models.User, classID uint) error {
	classIDs, err := s.ClassIDs(u)
	if err != nil {
		return err
	}
	for _, id := range classIDs {
		if id == classID {
			return nil
		}
	}
	return ErrNotFound("班级不存在或不在管辖范围")
}
