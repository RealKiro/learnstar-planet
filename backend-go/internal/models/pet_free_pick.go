// 「免费自选宠物」机会的持久化载体。
//
// 有意差异（与 Laravel 的 Cache 口径）：Laravel 整班切换宠物系列后用
// Cache::put("pet_free_pick:<student_id>", 1, 3 天) 发放、用 Cache::has 判定、
// 用 Cache::forget 消费；Go 端没有 Cache/Redis，改为数据库表 pet_free_picks
// （与 display_tokens 的同类改写一致），过期即视为未持有——语义等价，不做定时清理。
package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// PetFreePickTTL 「免费自选」机会有效期（同 Laravel PetSeriesService::FREE_PICK_DAYS = 3 天）。
const PetFreePickTTL = 3 * 24 * time.Hour

// PetFreePick 一名学生当前持有的「免费自选宠物」机会（每人最多一条）。
type PetFreePick struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StudentID uint      `gorm:"uniqueIndex;not null" json:"student_id"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GrantPetFreePick 发放/续期一次免费自选机会（同 Laravel Cache::put，重复发放即刷新有效期）。
func GrantPetFreePick(db *gorm.DB, studentID uint, now time.Time) error {
	if studentID == 0 {
		return nil
	}

	row := PetFreePick{StudentID: studentID}
	err := db.Where("student_id = ?", studentID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(&PetFreePick{
			StudentID: studentID,
			ExpiresAt: now.Add(PetFreePickTTL),
		}).Error
	}
	if err != nil {
		return err
	}

	return db.Model(&PetFreePick{}).Where("id = ?", row.ID).
		Update("expires_at", now.Add(PetFreePickTTL)).Error
}

// HasPetFreePick 判断学生是否持有未过期的免费自选机会（同 Laravel Cache::has）。
func HasPetFreePick(db *gorm.DB, studentID uint, now time.Time) (bool, error) {
	var row PetFreePick
	err := db.Where("student_id = ? AND expires_at > ?", studentID, now).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConsumePetFreePick 消费免费自选机会（同 Laravel Cache::forget；不存在也不报错）。
func ConsumePetFreePick(db *gorm.DB, studentID uint) error {
	return db.Where("student_id = ?", studentID).Delete(&PetFreePick{}).Error
}
