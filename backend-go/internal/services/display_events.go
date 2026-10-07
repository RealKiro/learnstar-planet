// 班级大屏事件总线（事件发布/消费 + 轮询降级的数据源）。
//
// 忠实移植自 Laravel App\Services\DisplayEventService：
//
//	publish(classId, type, data) 追加一条 {id, type, data, created_at} 事件；
//	consume(classId, sinceId) 只返回 id > sinceId 的事件（sinceId 为空则返回全部）；
//	clear(classId) 清空该班事件；
//	TTL 10 分钟（600 秒）、每班最多保留 200 条。
//
// 有意差异（本批统一约束）：
//  1. 存储：Laravel 用 Cache（Redis List / File Cache 降级），Go 端无 Cache → 改为数据库表
//     display_events（模型 models.DisplayEvent）。
//  2. 事件 id：Laravel 的 id 来自独立的 Cache 计数器键（display:events:<class_id>:counter，
//     同样 600 秒 TTL），并**不**等于列表下标；Go 端用「每班 seq 自增」表达同义语义——
//     seq = 该班当前未过期事件的最大 seq + 1，故计数器随事件一起过期、不留独立状态。
//     客户端回传的 last_event_id（Laravel 的 $e['id']）在 Go 端即 seq，参数名保持一致。
//  3. clear()：Laravel 的 Cache::forget 只清列表、计数器仍在（下次 id 继续往下走）；
//     Go 端 seq 从存量事件推导，故 clear 之后的下一条事件 seq 重新从 1 开始。客户端只做
//     「seq > 上次值」比较，与 Laravel 计数器 10 分钟过期后的重置属同一类边界行为。
//  4. 过期清理：Laravel 由 Cache TTL 自动回收，Go 端在 publish/consume 时惰性删除
//     （不引入定时任务）。
package services

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// 大屏事件类型（同 Laravel DisplayEventService 注释：score_update|broadcast|notice|pet_update|refresh）。
const (
	// DisplayEventScoreUpdate 积分变动。
	DisplayEventScoreUpdate = "score_update"
	// DisplayEventBroadcast 广播（banner/popup/fullscreen）。
	DisplayEventBroadcast = "broadcast"
	// DisplayEventNotice 通知公告。
	DisplayEventNotice = "notice"
	// DisplayEventPetUpdate 宠物状态变动（喂食）。
	DisplayEventPetUpdate = "pet_update"
	// DisplayEventRefresh 班级码刷新，通知已连接大屏重载。
	DisplayEventRefresh = "refresh"
)

// DisplayEvents 大屏事件发布/消费服务。
type DisplayEvents struct {
	db *gorm.DB
}

// NewDisplayEvents 创建大屏事件服务。
func NewDisplayEvents(db *gorm.DB) *DisplayEvents {
	return &DisplayEvents{db: db}
}

// DisplayEventView 单条事件的对外结构（字段名同 Laravel：id / type / data / created_at）。
type DisplayEventView struct {
	ID        int             `json:"id"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	CreatedAt string          `json:"created_at"`
}

// ============================================================
// 事件载荷（字段名/键序与 Laravel 各发布点逐字对齐）
// ============================================================

// DisplayScoreUpdateData score_update 事件载荷（Laravel ScoreService::giveScore / spendScore）。
type DisplayScoreUpdateData struct {
	StudentID     uint   `json:"student_id"`
	StudentName   string `json:"student_name"`
	StudentNo     string `json:"student_no"`
	Amount        int    `json:"amount"`
	Reason        string `json:"reason"`
	TotalScore    int    `json:"total_score"`
	PetLevel      *int   `json:"pet_level"`
	PetExperience *int   `json:"pet_experience"`
	PetMood       *int   `json:"pet_mood"`
	// IsSpend 仅「兑换消耗」路径出现（Laravel spendScore 的 'is_spend' => true）。
	IsSpend bool `json:"is_spend,omitempty"`
}

// DisplayBroadcastData broadcast 事件载荷（Laravel BroadcastService / ClassroomMessagingService）。
type DisplayBroadcastData struct {
	ID             uint   `json:"id"`
	Type           string `json:"type"`
	Content        string `json:"content"`
	DisplaySeconds int    `json:"display_seconds"`
	VoiceEnabled   bool   `json:"voice_enabled"`
	CreatedAt      string `json:"created_at"`
}

// DisplayNoticeData notice 事件载荷（Laravel ClassroomMessagingService 的 notice 分支）。
type DisplayNoticeData struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Type        string `json:"type"`
	PublishedAt string `json:"published_at"`
}

// DisplayPetUpdateData pet_update 事件载荷（Laravel PetService::feed）。
type DisplayPetUpdateData struct {
	StudentID   uint   `json:"student_id"`
	StudentName string `json:"student_name"`
	Type        string `json:"type"`
	Mood        int    `json:"mood"`
	Level       int    `json:"level"`
	Experience  int    `json:"experience"`
}

// DisplayRefreshData refresh 事件载荷。
// 教师端为 {old_code, new_code}（Laravel DisplayController::refreshDisplayCode），
// 管理员端只有 {new_code}（SchoolAdminController::refreshDisplayCode）。
type DisplayRefreshData struct {
	OldCode *string `json:"old_code,omitempty"`
	NewCode string  `json:"new_code"`
}

// ============================================================
// 发布 / 消费 / 清空
// ============================================================

// Publish 发布一条事件到班级频道：裁剪过期事件 → 取该班下一个 seq → 落库 → 保留最近 200 条。
func (e *DisplayEvents) Publish(classID uint, eventType string, data any) error {
	if classID == 0 {
		return nil
	}

	payload, err := marshalEventData(data)
	if err != nil {
		return err
	}

	return e.db.Transaction(func(tx *gorm.DB) error {
		if err := purgeExpiredEvents(tx, classID); err != nil {
			return err
		}

		seq, err := nextEventSeq(tx, classID)
		if err != nil {
			return err
		}

		row := models.DisplayEvent{
			ClassID:   classID,
			Seq:       seq,
			Type:      eventType,
			Data:      payload,
			CreatedAt: util.Now(),
		}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}

		return trimEvents(tx, classID)
	})
}

// Consume 消费自 since 之后的新事件（since 为 nil 时返回全部；过期事件不返回并被惰性删除）。
func (e *DisplayEvents) Consume(classID uint, since *int) ([]DisplayEventView, error) {
	if err := purgeExpiredEvents(e.db, classID); err != nil {
		return nil, err
	}

	q := e.db.Where("class_id = ?", classID)
	if since != nil {
		q = q.Where("seq > ?", *since)
	}

	var rows []models.DisplayEvent
	if err := q.Order("seq ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]DisplayEventView, 0, len(rows))
	for _, r := range rows {
		out = append(out, DisplayEventView{
			ID:        r.Seq,
			Type:      r.Type,
			Data:      json.RawMessage(r.Data),
			CreatedAt: r.CreatedAt.In(util.Loc).Format(time.RFC3339),
		})
	}
	return out, nil
}

// Clear 清除班级全部事件（同 Laravel DisplayEventService::clear）。
func (e *DisplayEvents) Clear(classID uint) error {
	return e.db.Where("class_id = ?", classID).Delete(&models.DisplayEvent{}).Error
}

// ============================================================
// 内部辅助
// ============================================================

// purgeExpiredEvents 删除该班 created_at 早于 10 分钟的事件（惰性清理，同 Laravel Cache TTL）。
func purgeExpiredEvents(db *gorm.DB, classID uint) error {
	cutoff := util.Now().Add(-models.DisplayEventTTL)
	return db.Where("class_id = ? AND created_at < ?", classID, cutoff).
		Delete(&models.DisplayEvent{}).Error
}

// nextEventSeq 该班下一个事件序号 = 当前最大 seq + 1（无事件时为 1，同 Laravel 计数器首次自增）。
func nextEventSeq(db *gorm.DB, classID uint) (int, error) {
	maxSeq := 0
	if err := db.Model(&models.DisplayEvent{}).Where("class_id = ?", classID).
		Select("COALESCE(MAX(seq), 0)").Scan(&maxSeq).Error; err != nil {
		return 0, err
	}
	return maxSeq + 1, nil
}

// trimEvents 只保留该班最近 200 条事件（超出即删除更早的，同 Laravel array_slice(-200)）。
func trimEvents(db *gorm.DB, classID uint) error {
	// 多取 1 条：len == 上限+1 说明超限，此时 seqs[上限] 即「第 201 新」的序号，其及更早的全部删除。
	keep := models.DisplayEventMaxPerClass + 1
	seqs := make([]int, 0, keep)
	if err := db.Model(&models.DisplayEvent{}).Where("class_id = ?", classID).
		Order("seq DESC").Limit(keep).Pluck("seq", &seqs).Error; err != nil {
		return err
	}
	if len(seqs) <= models.DisplayEventMaxPerClass {
		return nil
	}
	cutoff := seqs[models.DisplayEventMaxPerClass]
	return db.Where("class_id = ? AND seq <= ?", classID, cutoff).
		Delete(&models.DisplayEvent{}).Error
}

// marshalEventData 序列化事件载荷（等价 PHP json_encode：不转义 < > & 等 HTML 字符）。
func marshalEventData(data any) (string, error) {
	if data == nil {
		return "null", nil
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		return "", err
	}
	return string(bytes.TrimRight(buf.Bytes(), "\n")), nil
}

// publishDisplayEvent 尽力发布大屏事件：失败不影响主流程（同 Laravel 各发布点的 try/catch + warning）。
func publishDisplayEvent(events *DisplayEvents, classID uint, eventType string, data any) {
	if events == nil || classID == 0 {
		return
	}
	_ = events.Publish(classID, eventType, data)
}

// isoTimeOrEmpty 把可空时间格式化为 ISO8601（等价 Carbon toIso8601String；nil 输出空串）。
func isoTimeOrEmpty(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(util.Loc).Format(time.RFC3339)
}
