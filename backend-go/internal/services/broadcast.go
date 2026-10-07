// 班级广播（大屏横幅/弹窗/全屏）服务：读写教师管辖范围内的广播，并向班级大屏发布 broadcast 事件。
// 忠实移植自 Laravel App\Services\BroadcastService（含大屏 SSE 事件推送）。
package services

import (
	"errors"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// BroadcastService 班级广播读写服务。
type BroadcastService struct {
	db     *gorm.DB
	scope  *Scope
	events *DisplayEvents
}

// NewBroadcastService 创建广播服务（内部自带大屏事件发布器：broadcast 事件）。
func NewBroadcastService(db *gorm.DB, scope *Scope) *BroadcastService {
	return &BroadcastService{db: db, scope: scope, events: NewDisplayEvents(db)}
}

// Recent 返回教师管辖班级最近 20 条广播（按创建时间倒序）。
func (s *BroadcastService) Recent(u *models.User) ([]models.Broadcast, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return []models.Broadcast{}, nil
	}

	var broadcasts []models.Broadcast
	if err := s.db.Where("class_id IN ?", classIDs).
		Order("created_at DESC, id DESC").Limit(20).Find(&broadcasts).Error; err != nil {
		return nil, err
	}
	return broadcasts, nil
}

// FindInScope 按作用域查找广播（越权即 404）。
func (s *BroadcastService) FindInScope(u *models.User, id uint) (*models.Broadcast, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	var broadcast models.Broadcast
	err = s.db.Where("id = ? AND class_id IN ?", id, classIDs).First(&broadcast).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("广播不存在或不可见")
	}
	if err != nil {
		return nil, err
	}
	return &broadcast, nil
}

// Send 向目标班级发送广播。
//
// 目标班级 = 传入的 class_ids（若有）∩ 教师可访问班级；
// 未指定 class_ids 时发送给全部可访问班级。
//
// 返回实际发送的班级数；返回 0 表示没有可发送的班级。
func (s *BroadcastService) Send(u *models.User, content, broadcastType string, voice, loop bool, duration int, classIDs []uint) (int, error) {
	accessible, err := s.scope.ClassIDs(u)
	if err != nil {
		return 0, err
	}

	var targetIDs []uint
	if len(classIDs) == 0 {
		targetIDs = accessible
	} else {
		allowed := make(map[uint]bool, len(accessible))
		for _, id := range accessible {
			allowed[id] = true
		}
		for _, id := range classIDs {
			if allowed[id] {
				targetIDs = append(targetIDs, id)
			}
		}
	}

	if len(targetIDs) == 0 {
		return 0, nil
	}

	sent := 0
	now := util.Now()
	for _, classID := range targetIDs {
		broadcast := models.Broadcast{
			ClassID:        classID,
			SchoolID:       u.SchoolID,
			TeacherID:      &u.ID,
			Content:        content,
			Type:           broadcastType,
			VoiceEnabled:   voice,
			LoopEnabled:    loop,
			DisplaySeconds: duration,
			Status:         "sent",
			SentAt:         &now,
		}
		if err := s.db.Create(&broadcast).Error; err != nil {
			return sent, err
		}

		// 推送事件到班级大屏（broadcast），字段同 Laravel BroadcastService::send。
		publishDisplayEvent(s.events, classID, DisplayEventBroadcast, DisplayBroadcastData{
			ID:             broadcast.ID,
			Type:           broadcast.Type,
			Content:        broadcast.Content,
			DisplaySeconds: broadcast.DisplaySeconds,
			VoiceEnabled:   broadcast.VoiceEnabled,
			CreatedAt:      isoTimeOrEmpty(&broadcast.CreatedAt),
		})

		sent++
	}

	return sent, nil
}
