package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PUT|POST /admin/school 的 settings 语义同 Laravel：提供即**整体替换**（不是按键合并），
// 未提供保持原值；GET /admin/school 必须把 settings 作为对象输出（前端「学校设置」页依赖）。
func TestAdminUpdateSchoolSettingsReplace(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	admin := services.NewAdmin(db)

	// 预置一份含班级码前缀的设置。
	require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).
		Update("settings", `{"display_code_prefix":"XX","third_party_platform":"wechat_work"}`).Error)

	// 未提供 settings → 保持原值（只改 name）。
	updated, err := admin.UpdateSchool(school.ID, map[string]any{"name": "新校名"}, nil)
	require.NoError(t, err)
	assert.Equal(t, "新校名", updated.Name)
	assert.Equal(t, "XX", updated.SettingString("display_code_prefix", ""))
	assert.Equal(t, "wechat_work", updated.SettingString("third_party_platform", ""))

	// 提供 settings → 整体替换：旧键消失、新键生效（Laravel fill 的替换语义）。
	replacement := map[string]any{
		"third_party_platform":          "dingtalk",
		"enabled_third_party_platforms": []any{"dingtalk", "wechat_work"},
	}
	updated, err = admin.UpdateSchool(school.ID, map[string]any{}, &replacement)
	require.NoError(t, err)
	assert.Equal(t, "dingtalk", updated.SettingString("third_party_platform", ""))
	assert.Equal(t, "", updated.SettingString("display_code_prefix", ""), "旧键应被整体替换掉")

	settings := updated.SettingsMap()
	enabled, ok := settings["enabled_third_party_platforms"].([]any)
	require.True(t, ok, "数组设置应能读回：%#v", settings["enabled_third_party_platforms"])
	assert.Equal(t, []any{"dingtalk", "wechat_work"}, enabled)

	// 视图把 settings 解析成对象（不是原始 JSON 文本）。
	view := services.SchoolViewOf(updated)
	assert.Equal(t, "dingtalk", view.Settings["third_party_platform"])
	assert.Equal(t, "新校名", view.Name)

	// 空设置 → 视图输出空对象而非 null。
	empty := map[string]any{}
	cleared, err := admin.UpdateSchool(school.ID, map[string]any{}, &empty)
	require.NoError(t, err)
	view = services.SchoolViewOf(cleared)
	require.NotNil(t, view.Settings)
	assert.Empty(t, view.Settings)
	assert.Equal(t, "", cleared.SettingString("third_party_platform", ""))
}
