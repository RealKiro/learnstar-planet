// 管理端：全校积分规则与商品管理（学校级共享）+ 教师端我的班级/模式。
package handlers

import (
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// ---------------------------------- 管理端：积分规则 ----------------------------------

// AdminScoreRules 全校规则列表（学校级 + 教师创建的班级级）。
func (h *Handlers) AdminScoreRules(c *gin.Context) {
	u := middleware.CurrentUser(c)
	rules, err := h.rules.ListForSchool(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rules)
}

// AdminScoreRuleCreate 创建学校级规则（class_id = null）。
func (h *Handlers) AdminScoreRuleCreate(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Name       string `json:"name"`
		Amount     int    `json:"amount"`
		Category   string `json:"category"`
		IsPositive *bool  `json:"is_positive"`
		IsActive   *bool  `json:"is_active"`
	}
	if !bindJSON(c, &req) {
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" || len([]rune(name)) > 50 {
		fail(c, services.ErrUnprocessable("规则名称必填且不超过 50 字"))
		return
	}
	if req.Amount == 0 {
		fail(c, services.ErrUnprocessable("分值不能为 0"))
		return
	}

	isPositive := true
	if req.IsPositive != nil {
		isPositive = *req.IsPositive
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	rule, err := h.rules.CreateForSchool(u.SchoolID, name, req.Amount, req.Category, isPositive, isActive)
	if err != nil {
		fail(c, err)
		return
	}
	okCreated(c, rule)
}

// AdminScoreRuleUpdate 更新本校学校级规则（班级级 ID 一律 404）。
func (h *Handlers) AdminScoreRuleUpdate(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	rule, err := h.rules.FindSchoolLevel(u.SchoolID, id)
	if err != nil {
		fail(c, err)
		return
	}

	var req struct {
		Name       *string `json:"name"`
		Amount     *int    `json:"amount"`
		Category   *string `json:"category"`
		IsPositive *bool   `json:"is_positive"`
		IsActive   *bool   `json:"is_active"`
		SortOrder  *int    `json:"sort_order"`
	}
	if !bindJSON(c, &req) {
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" || len([]rune(name)) > 50 {
			fail(c, services.ErrUnprocessable("规则名称必填且不超过 50 字"))
			return
		}
		updates["name"] = name
	}
	if req.Amount != nil {
		if *req.Amount == 0 {
			fail(c, services.ErrUnprocessable("分值不能为 0"))
			return
		}
		updates["amount"] = *req.Amount
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.IsPositive != nil {
		updates["is_positive"] = *req.IsPositive
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	updated, err := h.rules.Update(rule, updates)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, updated)
}

// AdminScoreRuleDelete 删除本校学校级规则。
func (h *Handlers) AdminScoreRuleDelete(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	rule, err := h.rules.FindSchoolLevel(u.SchoolID, id)
	if err != nil {
		fail(c, err)
		return
	}
	if err := h.rules.Delete(rule); err != nil {
		fail(c, err)
		return
	}
	okMessage(c, "规则已删除")
}

// ---------------------------------- 管理端：商品 ----------------------------------

// AdminShopItems 全校商品列表（学校级 + 各班级级），可按币种过滤。
func (h *Handlers) AdminShopItems(c *gin.Context) {
	u := middleware.CurrentUser(c)
	items, err := h.shop.ItemsForSchool(u.SchoolID, c.Query("currency_type"))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, items)
}

// AdminShopItemCreate 创建学校级商品（class_id = null，同步所有班级）。
func (h *Handlers) AdminShopItemCreate(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		Category     string `json:"category"`
		CostScore    int    `json:"cost_score"`
		CurrencyType string `json:"currency_type"`
		EventTag     string `json:"event_tag"`
		Stock        int    `json:"stock"`
		ImagePath    string `json:"image_path"`
		IsActive     *bool  `json:"is_active"`
	}
	if !bindJSON(c, &req) {
		return
	}

	if err := validateShopItemIn(req.Name, req.CostScore, req.Stock, len([]rune(req.Description))); err != nil {
		fail(c, err)
		return
	}

	item, err := h.shop.CreateSchoolItem(u.SchoolID, services.SchoolItemInput{
		Name: req.Name, Description: req.Description, Category: req.Category,
		CostScore: req.CostScore, CurrencyType: req.CurrencyType, EventTag: req.EventTag,
		Stock: req.Stock, ImagePath: req.ImagePath, IsActive: req.IsActive,
	})
	if err != nil {
		fail(c, err)
		return
	}
	okCreated(c, item)
}

// AdminShopItemUpdate 更新本校学校级商品（班级级 ID 一律 404）。
func (h *Handlers) AdminShopItemUpdate(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	item, err := h.shop.FindSchoolItem(u.SchoolID, id)
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
		name := strings.TrimSpace(*req.Name)
		if name == "" || len([]rune(name)) > 100 {
			fail(c, services.ErrUnprocessable("商品名称必填且不超过 100 字"))
			return
		}
		updates["name"] = name
	}
	if req.CostScore != nil {
		if *req.CostScore < 1 {
			fail(c, services.ErrUnprocessable("所需积分至少为 1"))
			return
		}
		updates["cost_score"] = *req.CostScore
	}
	if req.Stock != nil {
		if *req.Stock < 0 {
			fail(c, services.ErrUnprocessable("库存不能为负数"))
			return
		}
		updates["stock"] = *req.Stock
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.CurrencyType != nil {
		updates["currency_type"] = *req.CurrencyType
	}
	if req.EventTag != nil {
		updates["event_tag"] = *req.EventTag
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

// AdminShopItemDelete 删除本校学校级商品。
func (h *Handlers) AdminShopItemDelete(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	item, err := h.shop.FindSchoolItem(u.SchoolID, id)
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

// validateShopItemIn 校验商品创建入参（对齐 Laravel 校验规则）。
func validateShopItemIn(name string, costScore, stock, descLen int) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" || len([]rune(trimmed)) > 100 {
		return services.ErrUnprocessable("商品名称必填且不超过 100 字")
	}
	if costScore < 1 {
		return services.ErrUnprocessable("所需积分至少为 1")
	}
	if stock < 0 {
		return services.ErrUnprocessable("库存不能为负数")
	}
	if descLen > 500 {
		return services.ErrUnprocessable("商品描述不超过 500 字")
	}
	return nil
}

// ---------------------------------- 教师端：我的班级 / 模式 ----------------------------------

// TeacherMyClasses 教师关联班级列表（机器人返回本校全部启用班级）。
func (h *Handlers) TeacherMyClasses(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classes, err := h.classroom.MyClasses(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, classes)
}

// TeacherSwitchClass 切换当前激活班级（未分配 → 403）。
func (h *Handlers) TeacherSwitchClass(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		ClassID uint `json:"class_id"`
	}
	if !bindJSON(c, &req) {
		return
	}

	if err := h.classroom.SwitchTo(u, req.ClassID); err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, gin.H{"message": "已切换", "data": gin.H{"active_class_id": req.ClassID}})
}

// TeacherGetMode 当前模式（默认 classroom_display）与激活班级。
func (h *Handlers) TeacherGetMode(c *gin.Context) {
	u := middleware.CurrentUser(c)
	mode, err := h.classroom.GetMode(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, mode)
}

// TeacherSetMode 切换模式（mode 非法或非 classroom_display 时缺密码 → 422）。
func (h *Handlers) TeacherSetMode(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		Mode     string `json:"mode"`
		ClassID  uint   `json:"class_id"`
		Password string `json:"password"`
	}
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.classroom.SetMode(u, req.Mode, req.ClassID, req.Password)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, result)
}
