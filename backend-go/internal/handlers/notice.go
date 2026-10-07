// 教师端通知处理器。
package handlers

import (
	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// TeacherNoticesList 返回教师管辖班级的通知列表。
func (h *Handlers) TeacherNoticesList(c *gin.Context) {
	u := middleware.CurrentUser(c)
	notices, err := h.notice.List(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, notices)
}

// TeacherNoticeCreate 创建通知。
func (h *Handlers) TeacherNoticeCreate(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Type    string `json:"type"`
	}
	if !bindJSON(c, &req) {
		return
	}
	// 校验 type 合法枚举（Laravel: in:info,homework,event,urgent）。缺省默认 info。
	switch req.Type {
	case "info", "homework", "event", "urgent", "":
	default:
		fail(c, services.ErrBadRequest("通知类型不合法"))
		return
	}

	notice, err := h.notice.Create(u, req.Title, req.Content, req.Type)
	if err != nil {
		fail(c, err)
		return
	}
	okCreated(c, notice)
}

// TeacherNoticeUpdate 更新通知。
func (h *Handlers) TeacherNoticeUpdate(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	notice, err := h.notice.FindInScope(u, id)
	if err != nil {
		fail(c, err)
		return
	}

	var req struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
		Type    *string `json:"type"`
	}
	if !bindJSON(c, &req) {
		return
	}

	updates := map[string]any{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Type != nil {
		switch *req.Type {
		case "info", "homework", "event", "urgent":
		default:
			fail(c, services.ErrBadRequest("通知类型不合法"))
			return
		}
		updates["type"] = *req.Type
	}

	updated, err := h.notice.Update(notice, updates)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, updated)
}

// TeacherNoticePublish 发布通知。
func (h *Handlers) TeacherNoticePublish(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	notice, err := h.notice.FindInScope(u, id)
	if err != nil {
		fail(c, err)
		return
	}
	if _, err := h.notice.Publish(notice); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "通知已发布")
}

// TeacherNoticeUnpublish 撤回通知。
func (h *Handlers) TeacherNoticeUnpublish(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	notice, err := h.notice.FindInScope(u, id)
	if err != nil {
		fail(c, err)
		return
	}
	if _, err := h.notice.Unpublish(notice); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "通知已撤回")
}

// TeacherNoticeDelete 删除通知。
func (h *Handlers) TeacherNoticeDelete(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	notice, err := h.notice.FindInScope(u, id)
	if err != nil {
		fail(c, err)
		return
	}
	if err := h.notice.Delete(notice); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "通知已删除")
}
