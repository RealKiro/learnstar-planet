// 班级通知（公告）服务：读写教师管辖范围内的通知，并在发布时向班级大屏推送 notice 事件。
// 忠实移植自 Laravel App\Services\NoticeService / ClassroomMessagingService 的 notice 分支。
package services

import (
	"errors"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// NoticeService 通知读写服务。
type NoticeService struct {
	db     *gorm.DB
	scope  *Scope
	events *DisplayEvents
}

// NewNoticeService 创建通知服务（内部自带大屏事件发布器：notice 事件）。
func NewNoticeService(db *gorm.DB, scope *Scope) *NoticeService {
	return &NoticeService{db: db, scope: scope, events: NewDisplayEvents(db)}
}

// List 返回教师管辖班级内的全部通知（按创建时间倒序）。
func (s *NoticeService) List(u *models.User) ([]models.Notice, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return []models.Notice{}, nil
	}

	var notices []models.Notice
	if err := s.db.Where("class_id IN ?", classIDs).
		Order("created_at DESC, id DESC").Find(&notices).Error; err != nil {
		return nil, err
	}
	return notices, nil
}

// Create 新建通知，落在教师首个管辖班级，默认未发布。
func (s *NoticeService) Create(u *models.User, title, content, noticeType string) (*models.Notice, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return nil, ErrUnprocessable("您暂无可管辖的班级")
	}
	if noticeType == "" {
		noticeType = "info"
	}

	notice := models.Notice{
		ClassID:     classIDs[0],
		SchoolID:    u.SchoolID,
		Title:       title,
		Content:     content,
		Type:        noticeType,
		PublishedBy: u.ID,
		IsPublished: false,
	}
	if err := s.db.Create(&notice).Error; err != nil {
		return nil, err
	}
	return &notice, nil
}

// FindInScope 按作用域查找通知（越权即 404）。
func (s *NoticeService) FindInScope(u *models.User, id uint) (*models.Notice, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	var notice models.Notice
	err = s.db.Where("id = ? AND class_id IN ?", id, classIDs).First(&notice).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("通知不存在或不可见")
	}
	if err != nil {
		return nil, err
	}
	return &notice, nil
}

// Update 更新通知的可编辑字段。
func (s *NoticeService) Update(notice *models.Notice, updates map[string]any) (*models.Notice, error) {
	if err := s.db.Model(notice).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(notice, notice.ID).Error; err != nil {
		return nil, err
	}
	return notice, nil
}

// Publish 发布通知并记录发布时间，同时向班级大屏推送 notice 事件
// （等价 Laravel ClassroomMessagingService::send 的 notice 分支 + NoticeService::publish 的广播语义）。
func (s *NoticeService) Publish(notice *models.Notice) (*models.Notice, error) {
	now := util.Now()
	if err := s.db.Model(notice).Updates(map[string]any{
		"is_published": true,
		"published_at": now,
	}).Error; err != nil {
		return nil, err
	}
	notice.IsPublished = true
	notice.PublishedAt = &now

	// 推送给班级大屏（notice）：字段/键序同 Laravel ClassroomMessagingService。
	publishDisplayEvent(s.events, notice.ClassID, DisplayEventNotice, DisplayNoticeData{
		ID:          notice.ID,
		Title:       notice.Title,
		Content:     notice.Content,
		Type:        notice.Type,
		PublishedAt: isoTimeOrEmpty(notice.PublishedAt),
	})

	return notice, nil
}

// Unpublish 撤回（取消发布）通知。
func (s *NoticeService) Unpublish(notice *models.Notice) (*models.Notice, error) {
	if err := s.db.Model(notice).Update("is_published", false).Error; err != nil {
		return nil, err
	}
	notice.IsPublished = false
	return notice, nil
}

// Delete 删除通知。
func (s *NoticeService) Delete(notice *models.Notice) error {
	return s.db.Delete(notice).Error
}
