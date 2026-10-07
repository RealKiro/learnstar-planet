package handlers

import (
	"net/http"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// validCurrency 校验币种代码合法（与 Laravel in:score,science,reading,class_point 一致）。
func validCurrency(code string) bool {
	switch code {
	case "score", "science", "reading", "class_point":
		return true
	}
	return false
}

// ============================================================
// 教师端：汇率配置与钱包/兑换记录查询
// ============================================================

// TeacherCurrencyRates 教师端汇率列表（首次访问惰性播种默认汇率）。
func (h *Handlers) TeacherCurrencyRates(c *gin.Context) {
	u := middleware.CurrentUser(c)
	rates, err := h.currency.RatesForSchool(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rates)
}

// TeacherCurrencyRateCreate 新增学校级汇率（教师端恒启用）。
func (h *Handlers) TeacherCurrencyRateCreate(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req struct {
		Name         string  `json:"name" binding:"required"`
		FromCurrency string  `json:"from_currency" binding:"required"`
		ToCurrency   string  `json:"to_currency" binding:"required"`
		Rate         float64 `json:"rate" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if len(req.Name) > 100 {
		fail(c, services.ErrUnprocessable("汇率名称不能超过 100 字"))
		return
	}
	if !validCurrency(req.FromCurrency) || !validCurrency(req.ToCurrency) {
		fail(c, services.ErrUnprocessable("币种不合法"))
		return
	}
	if req.Rate < 0.01 {
		fail(c, services.ErrUnprocessable("汇率必须大于等于 0.01"))
		return
	}

	rate, err := h.currency.CreateRate(u.SchoolID, services.RateInput{
		Name:         req.Name,
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         req.Rate,
	}, true)
	if err != nil {
		fail(c, err)
		return
	}
	okCreated(c, rate)
}

// TeacherCurrencyRateUpdate 更新学校级汇率。
func (h *Handlers) TeacherCurrencyRateUpdate(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	var req struct {
		Rate     *float64 `json:"rate"`
		IsActive *bool    `json:"is_active"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Rate != nil && *req.Rate < 0.01 {
		fail(c, services.ErrUnprocessable("汇率必须大于等于 0.01"))
		return
	}

	rate, err := h.currency.UpdateRate(u.SchoolID, id, req.Rate, req.IsActive)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rate)
}

// TeacherCurrencyWallets 教师所带班级学生的全部钱包。
func (h *Handlers) TeacherCurrencyWallets(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}
	items, err := h.currency.WalletsFor(classIDs)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, items)
}

// TeacherCurrencyLogs 兑换记录（分页，每页 20）。
func (h *Handlers) TeacherCurrencyLogs(c *gin.Context) {
	u := middleware.CurrentUser(c)
	classIDs, err := h.scope.ClassIDs(u)
	if err != nil {
		fail(c, err)
		return
	}
	result, err := h.currency.LogsFor(classIDs, queryInt(c, "page", 1))
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// TeacherCurrencyExchange 积分兑换为钱包币。
func (h *Handlers) TeacherCurrencyExchange(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req struct {
		StudentID  uint   `json:"student_id" binding:"required"`
		ToCurrency string `json:"to_currency" binding:"required"`
		Amount     int    `json:"amount" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Amount < 1 {
		fail(c, services.ErrUnprocessable("兑换数量必须大于 0"))
		return
	}
	student, err := h.scope.StudentInScope(u, req.StudentID)
	if err != nil {
		fail(c, err)
		return
	}
	result, err := h.currency.Exchange(student, req.ToCurrency, req.Amount, u.ID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// TeacherCurrencyCrossExchange 币种间兑换（钱包间）。
func (h *Handlers) TeacherCurrencyCrossExchange(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req struct {
		StudentID    uint   `json:"student_id" binding:"required"`
		FromCurrency string `json:"from_currency" binding:"required"`
		ToCurrency   string `json:"to_currency" binding:"required"`
		Amount       int    `json:"amount" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Amount < 1 {
		fail(c, services.ErrUnprocessable("兑换数量必须大于 0"))
		return
	}
	student, err := h.scope.StudentInScope(u, req.StudentID)
	if err != nil {
		fail(c, err)
		return
	}
	result, err := h.currency.CrossExchange(student, req.FromCurrency, req.ToCurrency, req.Amount, u.ID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// ============================================================
// 管理员端：汇率配置
// ============================================================

// AdminExchangeRates 管理员汇率列表（不播种默认）。
func (h *Handlers) AdminExchangeRates(c *gin.Context) {
	u := middleware.CurrentUser(c)
	rates, err := h.currency.ListRates(u.SchoolID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rates)
}

// AdminExchangeRateCreate 管理员新增汇率（可指定 is_active）。
func (h *Handlers) AdminExchangeRateCreate(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req struct {
		FromCurrency string  `json:"from_currency" binding:"required"`
		ToCurrency   string  `json:"to_currency" binding:"required"`
		Rate         float64 `json:"rate" binding:"required"`
		IsActive     *bool   `json:"is_active"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if !validCurrency(req.FromCurrency) || !validCurrency(req.ToCurrency) {
		fail(c, services.ErrUnprocessable("币种不合法"))
		return
	}
	if req.Rate < 0.01 {
		fail(c, services.ErrUnprocessable("汇率必须大于等于 0.01"))
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	rate, err := h.currency.CreateRate(u.SchoolID, services.RateInput{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         req.Rate,
	}, isActive)
	if err != nil {
		fail(c, err)
		return
	}
	okCreated(c, rate)
}

// AdminExchangeRateUpdate 管理员更新汇率（同教师端）。
func (h *Handlers) AdminExchangeRateUpdate(c *gin.Context) {
	h.TeacherCurrencyRateUpdate(c)
}
