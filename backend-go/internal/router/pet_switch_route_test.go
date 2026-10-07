// 教师端宠物切换的 **HTTP 契约**回归（入参 / 出参与图鉴副作用）。
//
// 权威来源：Laravel App\Http\Controllers\Api\TeacherController::switchPet（第 711-730 行）
// 与 App\Services\PetService::switchPet（第 155-276 行）：
//
//	入参 `pet_species`（必填 ≤50）+ `pet_name`（可空，缺省回退 pet_species）；
//	出参顶层 `{message, data}`，`cost` / `free_pick_used` 只在「原有宠物」分支出现。
//
// 全部使用内存 SQLite + httptest，**不访问外网**。
package router_test

import (
	"net/http"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTeacherPetSwitchHTTPContract 首次分配 → 再次切换（扣分 + 图鉴归档/激活）。
func TestTeacherPetSwitchHTTPContract(t *testing.T) {
	f := newPeopleFixture(t)
	require.NoError(t, f.DB.Model(&models.Student{}).Where("id = ?", f.Student.ID).
		Update("total_score", 50).Error)
	path := "/api/v1/teacher/pets/" + itoa(int(f.Student.ID)) + "/switch"

	// ① 无宠物 + 未传 pet_name → 名称回退物种 id；data 不含 cost / free_pick_used。
	first := f.do(t, f.TeacherToken, http.MethodPost, path, `{"pet_species":"zhulong"}`)
	require.Equal(t, http.StatusOK, first.StatusCode)
	firstBody := decodeBody(t, first)
	assert.Equal(t, "已为您分配宠物「zhulong」", firstBody["message"])
	data, ok := firstBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "zhulong", data["pet_name"])
	assert.Equal(t, "zhulong", data["pet_species"])
	assert.Equal(t, float64(1), data["level"])
	assert.Equal(t, float64(0), data["experience"])
	assert.NotContains(t, data, "cost")
	assert.NotContains(t, data, "free_pick_used")

	// ② 有宠物 + pet_name → 扣分（Lv.1 → 5），旧物种归档、目标激活、last_switched_at 写入。
	second := f.do(t, f.TeacherToken, http.MethodPost, path, `{"pet_species":"qilin","pet_name":"麒麟"}`)
	require.Equal(t, http.StatusOK, second.StatusCode)
	secondBody := decodeBody(t, second)
	assert.Equal(t, "宠物已更换为「麒麟」（扣除 5 积分）", secondBody["message"])
	secondData, ok := secondBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "麒麟", secondData["pet_name"])
	assert.Equal(t, "qilin", secondData["pet_species"])
	assert.Equal(t, float64(5), secondData["cost"])
	assert.Equal(t, false, secondData["free_pick_used"])

	var student models.Student
	require.NoError(t, f.DB.First(&student, f.Student.ID).Error)
	assert.Equal(t, 45, student.TotalScore)

	var pet models.Pet
	require.NoError(t, f.DB.Where("student_id = ?", f.Student.ID).First(&pet).Error)
	assert.Equal(t, "qilin", pet.Species)
	require.NotNil(t, pet.LastSwitchedAt)

	var archived models.PetCollection
	require.NoError(t, f.DB.Where("student_id = ? AND species = ?", f.Student.ID, "zhulong").
		First(&archived).Error)
	assert.False(t, archived.IsActive)

	// ③ 图鉴接口能看到两条记录，且激活物种是 qilin。
	collection := f.do(t, f.TeacherToken, http.MethodGet,
		"/api/v1/teacher/pets/"+itoa(int(f.Student.ID))+"/collection", "")
	require.Equal(t, http.StatusOK, collection.StatusCode)
	cData, ok := decodeBody(t, collection)["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "qilin", cData["active_species"])
	entries, ok := cData["collection"].([]any)
	require.True(t, ok)
	require.Len(t, entries, 2)
	activeBySpecies := map[string]bool{}
	for _, e := range entries {
		row := e.(map[string]any)
		activeBySpecies[row["species"].(string)] = row["is_active"].(bool)
	}
	assert.False(t, activeBySpecies["zhulong"])
	assert.True(t, activeBySpecies["qilin"])

	// ④ 缺 pet_species → 422（Laravel validate 的 required）。
	bad := f.do(t, f.TeacherToken, http.MethodPost, path, `{"pet_name":"x"}`)
	assert.Equal(t, http.StatusUnprocessableEntity, bad.StatusCode)

	// ⑤ 同物种 → 422「当前已经是这只宠物啦」。
	same := f.do(t, f.TeacherToken, http.MethodPost, path, `{"pet_species":"qilin","pet_name":"其它"}`)
	require.Equal(t, http.StatusUnprocessableEntity, same.StatusCode)
	assert.Equal(t, "当前已经是这只宠物啦", decodeBody(t, same)["message"])
}
