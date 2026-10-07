// 公开字典接口：积分分类（对应 Laravel StudentController::scoreCategories）。
package handlers

import (
	"net/http"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// ScoreCategories GET /api/v1/common/score-categories（公开，无需登录）。
//
// 响应 `{"data":[{id,name,icon,sort}, …]}`：**不带 `message`**（与 Laravel
// `response()->json(['data' => $data])` 一致，其余已移植接口用的是仓库统一信封）。
// id / name 取自 services.CategoryLabels 这一唯一真源，图标与 sort 见 services.ScoreCategories。
func (h *Handlers) ScoreCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": services.ScoreCategories()})
}
