// 班级大屏事件链路处理器：SSE 长连接、轮询降级、大屏加减分。
//
// 忠实移植自 Laravel App\Http\Controllers\Api\DisplayController 的
// sse / poll / quickScore / classroomGiveScore / classroomBatchGiveScore。
//
// 有意差异：
//  1. token 校验走 middleware.DisplayAuth（与其余大屏接口一致），失效即 401 JSON；
//     Laravel 的 sse 对无效 token 返回 200 + `event: error` 帧（前端 EventSource 自行跳登录页）。
//  2. Laravel 在 handler 内用 $request->input() 读参（query + body 均可），Go 端沿用
//     query 参数（last_event_id，参数名与 Laravel 一致）。
//  3. 成功响应统一走本项目的 {"data":...,"message":"ok"} 信封；scores/give 与
//     scores/batch-give 的消息文案逐字沿用 Laravel（加分成功/减分成功/批量操作完成…）。
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"github.com/gin-gonic/gin"
)

// sseTiming 是 SSE 循环的时序参数（对应 Laravel DisplayController::SSE_HEARTBEAT_INTERVAL
// = 10 秒、SSE_MAX_EXECUTION = 55 秒，以及空闲 2 秒 / 有事件 500 毫秒的轮询间隔）。
// 声明为变量是为了让测试覆写成极短值，避免 SSE 测试长时间挂起。
var sseTiming = struct {
	Heartbeat    time.Duration
	MaxExecution time.Duration
	IdleSleep    time.Duration
	BusySleep    time.Duration
}{
	Heartbeat:    10 * time.Second,
	MaxExecution: 55 * time.Second,
	IdleSleep:    2 * time.Second,
	BusySleep:    500 * time.Millisecond,
}

// ============================================================
// 显示端 API — SSE 长连接
// ============================================================

// DisplaySSE 大屏事件推送长连接。
//
// 实现同 Laravel：循环消费事件 → 有新事件立即推送 → 每 10 秒发一次心跳 →
// 最长运行 55 秒后优雅断开（发 `event: reconnect`，浏览器 EventSource 自动重连）。
// 差异：客户端断开由 c.Request.Context().Done() 检测并立即结束循环（PHP 无法感知断开，
// 只能空转到 55 秒）。
func (h *Handlers) DisplaySSE(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	lastEventID := queryInt(c, "last_event_id", 0)

	// SSE 必需的头（逐项同 Laravel：Content-Type / Cache-Control / Connection / X-Accel-Buffering）。
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Writer.WriteHeaderNow()

	ctx := c.Request.Context()
	start := time.Now()
	var lastHeartbeat time.Time

	for {
		if time.Since(start) >= sseTiming.MaxExecution {
			break
		}
		// 客户端已断开：立即结束，不再占用连接。
		select {
		case <-ctx.Done():
			return
		default:
		}

		events, err := h.display.ConsumeDisplayEvents(classID, &lastEventID)
		if err != nil {
			// 消费失败不打断长连接（poll 是降级通道，SSE 下一轮会重试）。
			events = nil
		}

		for _, ev := range events {
			// 帧格式同 Laravel：id / event / data 三行 + 空行结束。
			fmt.Fprintf(c.Writer, "id: %d\nevent: %s\ndata: %s\n\n", ev.ID, ev.Type, string(ev.Data))
			if ev.ID > lastEventID {
				lastEventID = ev.ID
			}
		}

		// 心跳（保持连接活跃）：首轮即发，之后每 Heartbeat 一次（同 Laravel $lastHeartbeat = 0）。
		if lastHeartbeat.IsZero() || time.Since(lastHeartbeat) >= sseTiming.Heartbeat {
			fmt.Fprintf(c.Writer, "event: heartbeat\ndata: {\"time\":\"%s\"}\n\n",
				util.Now().Format(time.RFC3339))
			lastHeartbeat = time.Now()
		}
		c.Writer.Flush()

		sleep := sseTiming.IdleSleep
		if len(events) > 0 {
			sleep = sseTiming.BusySleep
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(sleep):
		}
	}

	// 优雅关闭：发送重连指示。
	fmt.Fprint(c.Writer, "event: reconnect\ndata: {}\n\n")
	c.Writer.Flush()
}

// ============================================================
// 显示端 API — 轮询降级
// ============================================================

// DisplayPoll 轮询降级端点（SSE 不可用时使用），返回结构与 Laravel poll 一致：
// data.events / data.last_event_id / data.server_time。
func (h *Handlers) DisplayPoll(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	lastEventID := queryInt(c, "last_event_id", 0)
	var since *int
	if lastEventID != 0 {
		// Laravel `$lastEventId ?: null`：0 视为「未提供」→ 返回全部事件。
		since = &lastEventID
	}

	events, err := h.display.ConsumeDisplayEvents(classID, since)
	if err != nil {
		// Laravel 只记 warning 并返回空列表（降级通道不应因消费失败而 500）。
		events = []services.DisplayEventView{}
	}

	maxID := lastEventID
	for _, ev := range events {
		if ev.ID > maxID {
			maxID = ev.ID
		}
	}

	ok(c, gin.H{
		"events":        events,
		"last_event_id": maxID,
		"server_time":   util.Now().Format(time.RFC3339),
	})
}

// ============================================================
// 大屏快捷加减分
// ============================================================

// DisplayQuickScore 大屏快捷加减分（分值限于 ±1/±3/±5）。
func (h *Handlers) DisplayQuickScore(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	var req struct {
		StudentID *uint `json:"student_id"`
		Amount    *int  `json:"amount"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.StudentID == nil || req.Amount == nil {
		fail(c, services.ErrUnprocessable("student_id 与 amount 均为必填"))
		return
	}

	data, err := h.display.QuickScore(classID, *req.StudentID, *req.Amount)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, data)
}

// ============================================================
// 教室端加减分
// ============================================================

// DisplayGiveScore 教室端单个学生加减分（班级码模式单次 |points| < 30）。
func (h *Handlers) DisplayGiveScore(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	var req struct {
		StudentID *uint   `json:"student_id"`
		Points    *int    `json:"points"`
		Reason    *string `json:"reason"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.StudentID == nil || req.Points == nil {
		fail(c, services.ErrUnprocessable("student_id 与 points 均为必填"))
		return
	}
	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	result, err := h.display.ClassroomGiveScore(classID, *req.StudentID, *req.Points, reason)
	if err != nil {
		fail(c, err)
		return
	}

	// 文案同 Laravel：`($amount > 0 ? '加' : '减') . '分成功'`。
	message := "减分成功"
	if result.Points > 0 {
		message = "加分成功"
	}
	c.JSON(http.StatusOK, gin.H{"message": message, "data": result})
}

// DisplayBatchGiveScore 教室端批量加减分（单次 1~50 名学生）。
func (h *Handlers) DisplayBatchGiveScore(c *gin.Context) {
	classID, valid := displayClassID(c)
	if !valid {
		return
	}

	var req struct {
		StudentIDs []uint  `json:"student_ids"`
		Points     *int    `json:"points"`
		Reason     *string `json:"reason"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Points == nil {
		fail(c, services.ErrUnprocessable("points 为必填"))
		return
	}
	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	result, err := h.display.ClassroomBatchGiveScore(classID, req.StudentIDs, *req.Points, reason)
	if err != nil {
		fail(c, err)
		return
	}

	// 文案同 Laravel：'批量操作完成，处理了 N 名学生'。
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("批量操作完成，处理了 %d 名学生", result.Count),
		"data":    result,
	})
}
