package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PUT /admin/classes/:id 的 settings 语义：`pet_series` 合法则**按键合并**进现有 settings
// （保留其它键，同 Laravel `$settings['pet_series'] = …; $class->settings = $settings`），
// 取值白名单逐字同 Laravel `in:cosmic,pokemon,cute,treasure,mythic,all`；
// 班级对外视图把 settings 输出为对象（前端班级页读 `settings.pet_series`）。
func TestAdminUpdateClassPetSeries(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	admin := services.NewAdmin(db)

	class := models.ClassRoom{
		SchoolID: school.ID, Name: "一班", Grade: "1", Year: "2026", Status: "active",
		Settings: `{"display_code_prefix":"ZZ"}`,
	}
	require.NoError(t, db.Create(&class).Error)

	// 白名单校验：合法值通过、非法值拒绝。
	assert.True(t, services.ValidClassRoomSeries("all"))
	assert.True(t, services.ValidClassRoomSeries("mythic"))
	assert.False(t, services.ValidClassRoomSeries("myth"), "旧词表外的系列 id 应被拒绝（同 Laravel）")
	assert.False(t, services.ValidClassRoomSeries(""))

	updated, err := admin.UpdateClassWithSettings(school.ID, class.ID,
		map[string]any{"name": "一班（新）"},
		map[string]any{"pet_series": "pokemon"})
	require.NoError(t, err)
	assert.Equal(t, "一班（新）", updated.Name)
	assert.Equal(t, "pokemon", updated.SettingString("pet_series", ""))
	assert.Equal(t, "ZZ", updated.SettingString("display_code_prefix", ""), "其它 settings 键必须保留（合并而非替换）")

	view := services.ClassRoomViewOf(updated)
	assert.Equal(t, "pokemon", view.Settings["pet_series"])
	assert.Equal(t, "一班（新）", view.Name)

	// 不传 settingsPatch → settings 保持不变。
	again, err := admin.UpdateClass(school.ID, class.ID, map[string]any{"grade": "2"})
	require.NoError(t, err)
	assert.Equal(t, "pokemon", again.SettingString("pet_series", ""))
	assert.Equal(t, "ZZ", again.SettingString("display_code_prefix", ""))

	// 班级详情把 settings 输出为对象。
	detail, err := admin.ClassDetail(school.ID, class.ID)
	require.NoError(t, err)
	assert.Equal(t, "pokemon", detail.Settings["pet_series"])
	assert.Equal(t, "ZZ", detail.Settings["display_code_prefix"])

	// 空 settings → 视图输出空对象而非 null。
	empty := models.ClassRoom{SchoolID: school.ID, Name: "二班", Status: "active"}
	require.NoError(t, db.Create(&empty).Error)
	emptyView := services.ClassRoomViewOf(&empty)
	require.NotNil(t, emptyView.Settings)
	assert.Empty(t, emptyView.Settings)
}
