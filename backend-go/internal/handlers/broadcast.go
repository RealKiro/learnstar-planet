// 教师端广播处理器。
package handlers

import (
	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

// TeacherBroadcastsList 返回教师管辖班级的最近广播列表。
func (h *Handlers) TeacherBroadcastsList(c *gin.Context) {
	u := middleware.CurrentUser(c)
	broadcasts, err := h.broadcast.Recent(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, broadcasts)
}

// TeacherBroadcastSend 向目标班级发送广播。
func (h *Handlers) TeacherBroadcastSend(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Content  string `json:"content" binding:"required"`
		Type     string `json:"type"`
		ClassIDs []uint `json:"class_ids"`
		Voice    *bool  `json:"voice"`
		Loop     *bool  `json:"loop"`
		Duration *int   `json:"duration"`
	}
	if !bindJSON(c, &req) {
		return
	}

	// 校验 type 合法枚举（Laravel: in:banner,popup,fullscreen）。缺省默认 banner。
	switch req.Type {
	case "banner", "popup", "fullscreen", "":
	default:
		fail(c, services.ErrBadRequest("广播类型不合法"))
		return
	}
	if len(req.Content) > 500 {
		fail(c, services.ErrBadRequest("广播内容不能超过 500 字"))
		return
	}

	voice := true
	if req.Voice != nil {
		voice = *req.Voice
	}
	loop := false
	if req.Loop != nil {
		loop = *req.Loop
	}
	duration := 10
	if req.Duration != nil {
		duration = *req.Duration
		if duration < 0 || duration > 300 {
			fail(c, services.ErrBadRequest("展示时长需在 0-300 秒之间"))
			return
		}
	}

	sent, err := h.broadcast.Send(u, req.Content, req.Type, voice, loop, duration, req.ClassIDs)
	if err != nil {
		fail(c, err)
		return
	}
	if sent == 0 {
		fail(c, services.ErrBadRequest("没有可发送的班级"))
		return
	}
	okMessage(c, "广播已发送至 "+strconv.Itoa(sent)+" 个班级")
}

// TeacherBroadcastGet 返回单条广播。
func (h *Handlers) TeacherBroadcastGet(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	broadcast, err := h.broadcast.FindInScope(u, id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, broadcast)
}
