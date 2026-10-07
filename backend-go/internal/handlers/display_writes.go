// 班级大屏（教室端）写操作 / 查询 + 教师端宠物系列与 PK 处理器。
//
// 对应路由（本批新增）：
//
//	GET  /api/v1/display/shop-items       — 本班在售商品
//	POST /api/v1/display/redeem           — 快捷兑换
//	POST /api/v1/display/transfer         — 学生间转赠
//	POST /api/v1/display/switch-series    — 整班切换宠物系列
//	POST /api/v1/display/pets/switch      — 学生换宠
//	GET  /api/v1/display/pk/leaderboard   — 同年级 PK 排行（教室端口径）
//	POST /api/v1/teacher/class/switch-series
//	GET  /api/v1/teacher/pk/leaderboard
//	POST /api/v1/teacher/pk/challenge
//	GET  /api/v1/teacher/pk/my-stats
//
// 说明：带 message 的写操作按 Laravel 原样返回中文提示（与 handlers/timetable.go 同一做法）；
// 无 message 的读接口沿用项目统一信封 {"data":...,"message":"ok"}。
package handlers

import (
	"net/http"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

// ============================================================
// 教室端（班级码）—— 商城 / 转赠 / 宠物系列 / PK
// ============================================================

// DisplayShopItems 教室端商品列表（本班级级、在售）。
func (h *Handlers) DisplayShopItems(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	items, err := h.shop.DisplayItems(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, items)
}

// DisplayRedeem 教室端快捷兑换（立即扣分 + 生成已批准兑换记录）。
func (h *Handlers) DisplayRedeem(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	var req struct {
		StudentID uint `json:"student_id" binding:"required"`
		ItemID    uint `json:"item_id" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.shop.DisplayRedeem(classID, req.StudentID, req.ItemID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// DisplayTransfer 学生间积分转赠。
func (h *Handlers) DisplayTransfer(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	var req struct {
		FromID uint `json:"from_id" binding:"required"`
		ToID   uint `json:"to_id" binding:"required"`
		Amount int  `json:"amount" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.display.QuickTransfer(classID, req.FromID, req.ToID, req.Amount)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// DisplaySwitchSeries 教室端整班切换宠物系列（每人扣 20 积分 + 发免费自选机会）。
func (h *Handlers) DisplaySwitchSeries(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	var req struct {
		SeriesID string `json:"series_id" binding:"required,max=50"`
	}
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.petSeries.SwitchSeriesForClassroom(classID, req.SeriesID)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result.Data, "message": result.Message})
}

// DisplaySwitchPet 教室端学生换宠（首次免费机会或按等级扣分）。
func (h *Handlers) DisplaySwitchPet(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	var req struct {
		StudentID  uint   `json:"student_id" binding:"required"`
		PetSpecies string `json:"pet_species" binding:"required,max=50"`
	}
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.pets.SwitchForClassroom(classID, req.StudentID, req.PetSpecies)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result.Data, "message": result.Message})
}

// DisplayPKLeaderboard 教室端同年级 PK 排行榜。
func (h *Handlers) DisplayPKLeaderboard(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}
	rows, err := h.pk.LeaderboardForClass(classID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// ============================================================
// 教师端 —— 宠物系列 / PK
// ============================================================

// TeacherSwitchSeries 教师端切换本班宠物系列。
func (h *Handlers) TeacherSwitchSeries(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		SeriesID string `json:"series_id" binding:"required,max=50"`
	}
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.petSeries.SwitchSeries(u, req.SeriesID)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result.Data, "message": result.Message})
}

// TeacherPKLeaderboard 教师端同年级 PK 排行榜。
func (h *Handlers) TeacherPKLeaderboard(c *gin.Context) {
	u := middleware.CurrentUser(c)
	rows, err := h.pk.Leaderboard(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// TeacherPKMyStats 教师端本班 PK 统计（含名次）。
func (h *Handlers) TeacherPKMyStats(c *gin.Context) {
	u := middleware.CurrentUser(c)
	stats, err := h.pk.MyStats(u)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, stats)
}

// TeacherPKChallenge 教师端发起 PK 挑战。
func (h *Handlers) TeacherPKChallenge(c *gin.Context) {
	u := middleware.CurrentUser(c)

	var req struct {
		TargetClassID uint `json:"target_class_id" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.pk.Challenge(u, req.TargetClassID)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result.Data, "message": result.Message})
}
