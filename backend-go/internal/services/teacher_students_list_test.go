package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// List 严格对齐 Laravel StudentService::list：
//   - 不过滤 status（停用/毕业学生同样列出）
//   - 可选 search 同时命中姓名与学号，且**必须只在管辖班级范围内**匹配
//     （GORM 原始字符串条件不会自动补括号，缺括号会因 AND 优先级更高而带出越权行）
//   - order by name / 每页 50 / meta 四项 / 追加 pet_* 字段
func TestStudentServiceListMatchesLaravel(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	outside := seedClass(t, db, school.ID) // 同校但不在管辖范围

	mkStudent := func(classID uint, name, studentNo, status string) models.Student {
		student := models.Student{ClassID: classID, Name: name, StudentNo: studentNo, Status: status, Gender: "未知"}
		require.NoError(t, db.Create(&student).Error)
		return student
	}

	aaa := mkStudent(class.ID, "aaa", "001", "active")
	mkStudent(class.ID, "bbb", "002", "inactive")
	mkStudent(class.ID, "ccc", "003", "active")
	// 越权行：另一个班级的同名学生，学号命中搜索词时**不得**出现。
	mkStudent(outside.ID, "outside", "X999", "active")
	target := mkStudent(class.ID, "zzz", "X999", "active")

	// 宠物字段（Laravel with('pet') 后 array_merge 三个 pet_* 字段）。
	require.NoError(t, db.Create(&models.Pet{
		StudentID: aaa.ID, ClassID: class.ID, Name: "aaa的萌宠",
		Species: "zhulong", Level: 3, Experience: 40, Mood: 80,
	}).Error)

	svc := services.NewStudentService(db)

	// 1) 无搜索：不过滤 status，按 name 升序，pet_* 追加字段齐全。
	page, err := svc.List([]uint{class.ID}, "", 1, 0)
	require.NoError(t, err)
	require.Len(t, page.Data, 4)
	assert.Equal(t, []string{"aaa", "bbb", "ccc", "zzz"},
		[]string{page.Data[0].Name, page.Data[1].Name, page.Data[2].Name, page.Data[3].Name},
		"Laravel orderBy('name') 且不按 status 过滤")
	assert.Equal(t, 4, int(page.Meta.Total))
	assert.Equal(t, 1, page.Meta.CurrentPage)
	assert.Equal(t, 1, page.Meta.LastPage)
	assert.Equal(t, 50, page.Meta.PerPage, "Laravel paginate(50) 固定每页 50")
	assert.Equal(t, "zhulong", page.Data[0].PetSpecies)
	assert.Equal(t, 3, page.Data[0].PetLevel)
	assert.Equal(t, "aaa的萌宠", page.Data[0].PetName)
	assert.Equal(t, "", page.Data[1].PetSpecies, "无宠物学生 pet_* 为空值")

	// 2) search 按学号命中：越权班级的同号学生必须被排除（括号优先级回归）。
	byNo, err := svc.List([]uint{class.ID}, "X999", 1, 0)
	require.NoError(t, err)
	require.Len(t, byNo.Data, 1, "search 命中必须在 class_id 范围内（AND/OR 括号）")
	assert.Equal(t, "zzz", byNo.Data[0].Name)
	assert.Equal(t, target.ID, byNo.Data[0].ID)
	assert.Equal(t, 1, int(byNo.Meta.Total))

	// 3) search 按姓名命中。
	byName, err := svc.List([]uint{class.ID}, "bb", 1, 0)
	require.NoError(t, err)
	require.Len(t, byName.Data, 1)
	assert.Equal(t, "bbb", byName.Data[0].Name)

	// 4) 分页：perPage 可注入（处理器固定传 50），末页 last_page 按 total 计算。
	paged, err := svc.List([]uint{class.ID}, "", 2, 3)
	require.NoError(t, err)
	assert.Len(t, paged.Data, 1)
	assert.Equal(t, 2, paged.Meta.LastPage)

	// 5) 无管辖班级 → 空列表而非报错。
	empty, err := svc.List(nil, "", 1, 0)
	require.NoError(t, err)
	assert.Empty(t, empty.Data)
	assert.Equal(t, 0, int(empty.Meta.Total))
}
