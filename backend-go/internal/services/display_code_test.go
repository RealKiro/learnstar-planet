// 班级大屏码生成服务测试：确定性、年级/班号解析、前缀回退、兜底序号与批量重置。
package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// seedCodeClass 创建指定年级 / 班级名的班级。
func seedCodeClass(t *testing.T, db *gorm.DB, schoolID uint, grade, name string) models.ClassRoom {
	t.Helper()
	class := models.ClassRoom{SchoolID: schoolID, Grade: grade, Name: name, Status: "active"}
	require.NoError(t, db.Create(&class).Error)
	return class
}

func TestGradeDigit(t *testing.T) {
	cases := map[string]string{
		"一年级": "1",
		"三年级": "3",
		"九年级": "9",
		// 高中年级数字为 10/11/12，超出 1-9 上限 → 无法编码（与原实现一致）。
		"高一":  "",
		"高二":  "",
		"高三":  "",
		"3年级": "3",
		"":    "",
		"十年级": "",
		"幼儿园": "",
		"年级":  "",
	}
	for grade, want := range cases {
		assert.Equal(t, want, services.GradeDigit(grade), "grade=%q", grade)
	}
}

func TestClassNoFormats(t *testing.T) {
	cases := map[string]string{
		"一年级（3）班": "3",
		"一年级(3)班": "3",
		"一年级 1 班": "1",
		"一年级1班":   "1",
		"一年级第一班":  "1",
		"一年级第三班":  "3",
		"一年级一班":   "1",
		"三年级十班":   "", // 10 > 9
		"一年级第十班":  "", // 10 > 9
		"一年级12班":  "",
		"三年级":     "",
		"（0）班":    "",
	}
	for name, want := range cases {
		assert.Equal(t, want, services.ClassNo(name), "name=%q", name)
	}
}

func TestChineseNumberToArabic(t *testing.T) {
	assert.Equal(t, "10", services.ChineseNumberToArabic("十"))
	assert.Equal(t, "12", services.ChineseNumberToArabic("十二"))
	assert.Equal(t, "23", services.ChineseNumberToArabic("二十三"))
	assert.Equal(t, "3", services.ChineseNumberToArabic("三"))
	assert.Equal(t, "0", services.ChineseNumberToArabic("廿"))
}

// 班级码为确定性：同一班级每次结果一致；默认前缀 LS。
func TestGenerateDeterministicDefaultPrefix(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedCodeClass(t, db, school.ID, "一年级", "一年级（1）班")
	svc := services.NewDisplayCodeService(db)

	first := svc.Generate(&class, nil)
	assert.Equal(t, "LS11", first)
	assert.Equal(t, first, svc.Generate(&class, nil), "同一班级必须稳定输出同一班级码")

	second := seedCodeClass(t, db, school.ID, "二年级", "二年级3班")
	assert.Equal(t, "LS23", svc.Generate(&second, nil))
}

// 学校自定义前缀生效；非法前缀回退 LS。
func TestGenerateSchoolPrefixFallback(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedCodeClass(t, db, school.ID, "一年级", "一年级（1）班")
	svc := services.NewDisplayCodeService(db)

	assert.Equal(t, "LS", svc.SchoolPrefix(school.ID), "未配置时回退默认前缀")

	// 自定义前缀（小写会被转大写）。
	require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).
		Update("settings", `{"display_code_prefix":"bj"}`).Error)
	assert.Equal(t, "BJ", svc.SchoolPrefix(school.ID))
	assert.Equal(t, "BJ11", svc.Generate(&class, nil))

	// 非法前缀（含数字 / 超长 / 空）回退 LS。
	for _, bad := range []string{"B1", "LSXLSX", "L", "", "12"} {
		require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).
			Update("settings", `{"display_code_prefix":"`+bad+`"}`).Error)
		assert.Equal(t, "LS", svc.SchoolPrefix(school.ID), "bad prefix=%q", bad)
	}
}

// 显式传空前缀表示「不加前缀」（Laravel 的 $prefix ?? ... 语义：空串不回退学校前缀）。
func TestGenerateExplicitEmptyPrefix(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedCodeClass(t, db, school.ID, "一年级", "一年级（1）班")
	svc := services.NewDisplayCodeService(db)

	empty := ""
	assert.Equal(t, "11", svc.Generate(&class, &empty))
	custom := "BJ"
	assert.Equal(t, "BJ11", svc.Generate(&class, &custom))
}

// 年级无法解析返回空串。
func TestGenerateUnparsableGrade(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedCodeClass(t, db, school.ID, "幼儿园", "幼儿园（1）班")
	svc := services.NewDisplayCodeService(db)
	assert.Equal(t, "", svc.Generate(&class, nil))
}

// 班号解析不到时用同年级按 id 排序的序号；超过 9 无法编码。
func TestFallbackClassNoAndOverflow(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewDisplayCodeService(db)

	first := seedCodeClass(t, db, school.ID, "三年级", "三年级甲班")
	second := seedCodeClass(t, db, school.ID, "三年级", "三年级乙班")
	// 另一个年级不参与序号计算。
	seedCodeClass(t, db, school.ID, "四年级", "四年级甲班")

	assert.Equal(t, "1", svc.FallbackClassNo(&first))
	assert.Equal(t, "2", svc.FallbackClassNo(&second))
	assert.Equal(t, "LS31", svc.Generate(&first, nil))
	assert.Equal(t, "LS32", svc.Generate(&second, nil))

	// 同年级第 10 个班无法编码（班号 > 9）。
	var tenth models.ClassRoom
	for i := 3; i <= 10; i++ {
		tenth = seedCodeClass(t, db, school.ID, "三年级", "三年级附加班"+string(rune('A'+i)))
	}
	assert.Equal(t, "10", svc.FallbackClassNo(&tenth))
	assert.Equal(t, "", svc.Generate(&tenth, nil))
}

// 批量重置：regenerated / skipped / conflicts 计数与落库。
func TestRegenerateAll(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewDisplayCodeService(db)

	ok1 := seedCodeClass(t, db, school.ID, "一年级", "一年级（1）班")
	ok2 := seedCodeClass(t, db, school.ID, "一年级", "一年级（2）班")
	seedCodeClass(t, db, school.ID, "学前班", "学前班（1）班")             // 年级无法解析 → skipped
	conflict := seedCodeClass(t, db, school.ID, "一年级", "一年级（1）班") // 与 ok1 同码 → conflict

	result, err := svc.RegenerateAll(&school.ID, "BJ")
	require.NoError(t, err)
	assert.Equal(t, 2, result.Regenerated)
	assert.Equal(t, 2, result.Skipped)
	assert.Equal(t, map[uint]string{conflict.ID: "BJ11"}, result.Conflicts)

	var reloaded1, reloaded2 models.ClassRoom
	require.NoError(t, db.First(&reloaded1, ok1.ID).Error)
	require.NoError(t, db.First(&reloaded2, ok2.ID).Error)
	assert.Equal(t, "BJ11", reloaded1.DisplayCode)
	assert.Equal(t, "BJ12", reloaded2.DisplayCode)
	require.NotNil(t, reloaded1.DisplayCodeUpdatedAt)

	// 冲突班级不写库。
	var reloadedConflict models.ClassRoom
	require.NoError(t, db.First(&reloadedConflict, conflict.ID).Error)
	assert.Equal(t, "", reloadedConflict.DisplayCode)

	// 再次执行结果一致（确定性 / 幂等）。
	again, err := svc.RegenerateAll(&school.ID, "BJ")
	require.NoError(t, err)
	assert.Equal(t, result.Regenerated, again.Regenerated)
	assert.Equal(t, result.Skipped, again.Skipped)
	assert.Equal(t, result.Conflicts, again.Conflicts)
}

// 空前缀按学校设置；非法前缀回退 LS；未指定学校时用默认前缀。
func TestRegenerateAllPrefixResolution(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewDisplayCodeService(db)
	seedCodeClass(t, db, school.ID, "一年级", "一年级（1）班")

	require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).
		Update("settings", `{"display_code_prefix":"LE"}`).Error)

	result, err := svc.RegenerateAll(&school.ID, "")
	require.NoError(t, err)
	require.Equal(t, 1, result.Regenerated)
	var class models.ClassRoom
	require.NoError(t, db.Where("school_id = ?", school.ID).First(&class).Error)
	assert.Equal(t, "LE11", class.DisplayCode)

	// 非法前缀回退 LS。
	_, err = svc.RegenerateAll(&school.ID, "bad-prefix")
	require.NoError(t, err)
	require.NoError(t, db.Where("school_id = ?", school.ID).First(&class).Error)
	assert.Equal(t, "LS11", class.DisplayCode)

	// 未指定学校 + 非法前缀：仍回退 LS。
	all, err := svc.RegenerateAll(nil, "1234")
	require.NoError(t, err)
	assert.Equal(t, 1, all.Regenerated)
	require.NoError(t, db.Where("school_id = ?", school.ID).First(&class).Error)
	assert.Equal(t, "LS11", class.DisplayCode)
}
