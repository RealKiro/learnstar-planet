// 教师端服务：管辖班级、学生列表与仪表盘统计。
package services

import (
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// Teacher 教师端业务服务。
type Teacher struct {
	db    *gorm.DB
	scope *Scope
}

// NewTeacher 创建教师端服务。
func NewTeacher(db *gorm.DB, scope *Scope) *Teacher {
	return &Teacher{db: db, scope: scope}
}

// Dashboard 返回教师首页统计：班级数、学生数、积分汇总。
func (t *Teacher) Dashboard(u *models.User) (map[string]any, error) {
	classIDs, err := t.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	var classCount, studentCount int64
	if len(classIDs) > 0 {
		if err := t.db.Model(&models.ClassRoom{}).Where("id IN ?", classIDs).Count(&classCount).Error; err != nil {
			return nil, err
		}
		if err := t.db.Model(&models.Student{}).Where("class_id IN ? AND status = ?", classIDs, "active").
			Count(&studentCount).Error; err != nil {
			return nil, err
		}
	}

	summary, err := NewScoreService(t.db).Summary(classIDs)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"class_count":   classCount,
		"student_count": studentCount,
		"score":         summary,
	}, nil
}
