// 教师端班级大屏消息处理器（大屏聚合数据 + 消息轮询 / 发送）。
// 路由路径逐字对齐 Laravel backend/routes/api.php 第 171-175 行；
// 行为对齐 TeacherController::classroomDisplay / pollClassroomMessages / sendClassroomMessage。
package handlers

import (
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/gin-gonic/gin"
)

// 课堂消息长度上限（同 Laravel content max:500 / title max:200）。
const (
	classroomMessageContentMaxRunes = 500
	classroomMessageTitleMaxRunes   = 200
)

// classroomTargetClassID 解析目标班级：优先 class_id 查询参数，缺省取用户设置里的激活班级。
func classroomTargetClassID(c *gin.Context, u *models.User) uint {
	if raw := strings.TrimSpace(c.Query("class_id")); raw != "" {
		return queryID(c, "class_id")
	}
	return u.SettingUint(services.ActiveClassSettingKey)
}

// TeacherClassroomDisplay 班级大屏聚合数据（学生宠物总览 + 广播 + 通知 + 最近积分）。
func (h *Handlers) TeacherClassroomDisplay(c *gin.Context) {
	u := middleware.CurrentUser(c)

	data, err := h.messaging.Display(u, classroomTargetClassID(c, u))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// TeacherClassroomMessages 大屏轮询增量消息（since 之后的广播与通知）。
func (h *Handlers) TeacherClassroomMessages(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var since *time.Time
	if raw := strings.TrimSpace(c.Query("since")); raw != "" {
		parsed, err := parseLooseTime(raw)
		if err != nil {
			fail(c, services.ErrUnprocessable("since 时间格式不正确"))
			return
		}
		since = &parsed
	}

	data, err := h.messaging.Poll(classroomTargetClassID(c, u), since)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// classroomMessageRequest 发送课堂消息的请求体。
type classroomMessageRequest struct {
	ClassID        *uint   `json:"class_id"`
	Content        string  `json:"content"`
	Title          *string `json:"title"`
	Type           string  `json:"type"`
	DisplaySeconds *int    `json:"display_seconds"`
	Voice          *bool   `json:"voice"`
}

// TeacherClassroomMessageSend 发送班级消息（广播或通知）。
//
// 越权检查先于参数校验（保持 Laravel 的历史响应顺序：越权 403 优先于 422）。
func (h *Handlers) TeacherClassroomMessageSend(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req classroomMessageRequest
	if !bindJSON(c, &req) {
		return
	}

	accessible, err := h.messaging.AccessibleClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}

	classID := uint(0)
	if len(accessible) > 0 {
		classID = accessible[0]
	}
	if req.ClassID != nil {
		classID = *req.ClassID
	}

	// 权限预检先于参数校验。
	if !containsUint(accessible, classID) {
		fail(c, services.ErrForbidden("无权限"))
		return
	}

	if classID == 0 {
		fail(c, services.ErrUnprocessable("class_id 必填"))
		return
	}
	if !services.IsClassroomMessageType(req.Type) {
		fail(c, services.ErrUnprocessable("type 取值不合法"))
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		fail(c, services.ErrUnprocessable("content 必填"))
		return
	}
	if utf8.RuneCountInString(req.Content) > classroomMessageContentMaxRunes {
		fail(c, services.ErrUnprocessable("content 不能超过 500 字"))
		return
	}
	if req.Title != nil && utf8.RuneCountInString(*req.Title) > classroomMessageTitleMaxRunes {
		fail(c, services.ErrUnprocessable("title 不能超过 200 字"))
		return
	}

	displaySeconds := 10
	if req.DisplaySeconds != nil {
		displaySeconds = *req.DisplaySeconds
		if displaySeconds < 3 || displaySeconds > 300 {
			fail(c, services.ErrUnprocessable("display_seconds 需在 3-300 之间"))
			return
		}
	}
	voice := true
	if req.Voice != nil {
		voice = *req.Voice
	}

	result, err := h.messaging.Send(u, accessible, classID, req.Type, services.ClassroomMessageInput{
		Content:        req.Content,
		Title:          req.Title,
		DisplaySeconds: displaySeconds,
		Voice:          voice,
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// parseLooseTime 宽松解析 since（RFC3339 / "2006-01-02 15:04:05" / "2006-01-02"，业务时区）。
func parseLooseTime(raw string) (time.Time, error) {
	layouts := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}
	var lastErr error
	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, raw, util.Loc)
		if err == nil {
			return parsed, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

// containsUint 判断 uint 切片是否包含目标值。
func containsUint(list []uint, value uint) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
