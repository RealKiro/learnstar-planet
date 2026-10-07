// 教师端商城处理器。
package handlers

import (
	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// TeacherShopItems 返回商城商品列表（可按币种过滤，首次访问自动播种默认商品）。
func (h *Handlers) TeacherShopItems(c *gin.Context) {
	u := middleware.CurrentUser(c)
	items, err := h.shop.Items(u, c.Query("currency_type"))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, items)
}

// TeacherShopItemCreate 创建学校级商品。
func (h *Handlers) TeacherShopItemCreate(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Name         string `json:"name" binding:"required"`
		Description  string `json:"description"`
		Category     string `json:"category"`
		CostScore    int    `json:"cost_score" binding:"required"`
		CurrencyType string `json:"currency_type"`
		EventTag     string `json:"event_tag"`
		Stock        int    `json:"stock"`
		ImagePath    string `json:"image_path"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.CostScore < 1 {
		fail(c, services.ErrBadRequest("cost_score 必须大于 0"))
		return
	}

	item, err := h.shop.CreateItem(u, req.Name, req.Description, req.Category,
		req.CurrencyType, req.EventTag, req.ImagePath, req.CostScore, req.Stock)
	if err != nil {
		fail(c, err)
		return
	}
	okCreated(c, item)
}

// TeacherShopItemUpdate 更新商品。
func (h *Handlers) TeacherShopItemUpdate(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	item, err := h.shop.FindItem(u, id)
	if err != nil {
		fail(c, err)
		return
	}

	var req struct {
		Name         *string `json:"name"`
		Description  *string `json:"description"`
		Category     *string `json:"category"`
		CostScore    *int    `json:"cost_score"`
		CurrencyType *string `json:"currency_type"`
		EventTag     *string `json:"event_tag"`
		Stock        *int    `json:"stock"`
		ImagePath    *string `json:"image_path"`
		IsActive     *bool   `json:"is_active"`
	}
	if !bindJSON(c, &req) {
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.CostScore != nil {
		updates["cost_score"] = *req.CostScore
	}
	if req.CurrencyType != nil {
		updates["currency_type"] = *req.CurrencyType
	}
	if req.EventTag != nil {
		updates["event_tag"] = *req.EventTag
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}
	if req.ImagePath != nil {
		updates["image_path"] = *req.ImagePath
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	updated, err := h.shop.UpdateItem(item, updates)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, updated)
}

// TeacherShopItemDelete 删除商品。
func (h *Handlers) TeacherShopItemDelete(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	item, err := h.shop.FindItem(u, id)
	if err != nil {
		fail(c, err)
		return
	}
	if err := h.shop.DeleteItem(item); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "商品已删除")
}

// TeacherShopRedemptions 返回兑换记录列表。
func (h *Handlers) TeacherShopRedemptions(c *gin.Context) {
	u := middleware.CurrentUser(c)
	redemptions, err := h.shop.ListRedemptions(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, redemptions)
}

// TeacherShopRedemptionCreate 教师代学生发起兑换。
func (h *Handlers) TeacherShopRedemptionCreate(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		StudentID uint `json:"student_id" binding:"required"`
		ItemID    uint `json:"shop_item_id" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}

	item, err := h.shop.FindItem(u, req.ItemID)
	if err != nil {
		fail(c, err)
		return
	}
	if !item.IsActive {
		fail(c, services.ErrUnprocessable("该商品已下架"))
		return
	}

	redemption, err := h.shop.CreateRedemption(u, req.StudentID, item)
	if err != nil {
		fail(c, err)
		return
	}
	okCreated(c, redemption)
}

// TeacherShopRedemptionApprove 审批通过并结算。
func (h *Handlers) TeacherShopRedemptionApprove(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	result, err := h.shop.ApproveRedemption(u, id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// TeacherShopRedemptionReject 拒绝兑换。
func (h *Handlers) TeacherShopRedemptionReject(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	if err := h.shop.RejectRedemption(u, id); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "已拒绝兑换")
}

// TeacherShopRedemptionDeliver 标记已发放。
func (h *Handlers) TeacherShopRedemptionDeliver(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	if err := h.shop.DeliverRedemption(u, id); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "已标记为已发放")
}
