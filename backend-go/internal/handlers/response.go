// 统一响应与请求解析辅助函数。
package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// ok 返回统一成功信封 {"data":..., "message":"ok"}。
func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data, "message": "ok"})
}

// okCreated 返回 201 成功信封。
func okCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data, "message": "ok"})
}

// okMessage 返回仅带提示文案的成功信封。
func okMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"data": nil, "message": message})
}

// fail 将业务错误映射为对应状态码；未知错误统一 500。
// ValidationError（服务层抛出的字段校验失败）渲染为 Laravel 风格的 422 + errors。
func fail(c *gin.Context, err error) {
	if ve, isValidation := services.AsValidationError(err); isValidation {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "参数错误", "errors": ve.Errors})
		return
	}
	if ae, isApp := services.AsAppError(err); isApp {
		c.JSON(ae.Status, gin.H{"message": ae.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
}

// bindJSON 解析 JSON 请求体，失败时返回 422。
func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "请求参数格式错误"})
		return false
	}
	return true
}

// paramID 解析路径参数为正整数 ID，失败时写 400 并返回 false。
func paramID(c *gin.Context, key string) (uint, bool) {
	n, err := strconv.ParseUint(c.Param(key), 10, 32)
	if err != nil || n == 0 {
		fail(c, services.ErrBadRequest("无效的 ID"))
		return 0, false
	}
	return uint(n), true
}

// queryID 解析查询参数为 ID，失败返回 0。
func queryID(c *gin.Context, key string) uint {
	n, _ := strconv.ParseUint(c.Query(key), 10, 32)
	return uint(n)
}

// queryInt 解析查询参数为整数，缺省或非法时返回默认值。
func queryInt(c *gin.Context, key string, def int) int {
	v := c.Query(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// validationFailed 返回 Laravel 风格的 422 校验失败响应（message="参数错误" + errors）。
// 用于非 binding 的结构化校验（字段键与 Laravel validator 的键一致，如 assignments.0.class_id）。
func validationFailed(c *gin.Context, errs map[string][]string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "参数错误", "errors": errs})
}

// readJSONBody 解析 JSON 请求体；空体视为空对象（Laravel 的 $request->all() 语义），
// 非法 JSON 返回 422「请求参数格式错误」。
func readJSONBody(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "请求参数格式错误"})
		return false
	}
	return true
}

// itoa 十进制整数转字符串（用于拼接 Laravel 风格的 errors 键与响应文案）。
func itoa(n int) string { return strconv.Itoa(n) }
