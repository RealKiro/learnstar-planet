// 课表服务层测试：保存/快照/审批流/教师周课表/CSES 导出。
package services_test

import (
	"strings"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTimetableService(db *gorm.DB) *services.TimetableService {
	return services.NewTimetableService(db)
}

func strPtr(s string) *string { return &s }

// seedTimetable 准备学校 / 班级 / 教师。
func seedTimetable(t *testing.T, db *gorm.DB) (models.School, models.ClassRoom, models.User) {
	t.Helper()
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "t1")
	class, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	return school, class, teacher
}

func TestPaletteColor_StableAndInPalette(t *testing.T) {
	svc := newTimetableService(setupDB(t))

	first := svc.PaletteColor("语文")
	assert.Equal(t, first, svc.PaletteColor("语文"), "同名科目必须稳定同色")
	assert.True(t, strings.HasPrefix(first, "#"), "应返回色板中的十六进制色值")
	assert.NotEqual(t, svc.PaletteColor("语文"), svc.PaletteColor("数学"), "不同科目应取到不同颜色")
}

func TestSave_Bootstrap(t *testing.T) {
	db := setupDB(t)
	school, class, _ := seedTimetable(t, db)
	svc := newTimetableService(db)

	payload := services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{
			{Name: "语文"},
			{Name: "数学", Color: strPtr("#123456"), SimplifiedName: strPtr("数")},
		},
		Periods: []services.TimetablePeriodInput{
			{PeriodIndex: 1, Name: "第一节", StartTime: "08:00", EndTime: "08:45"},
			{PeriodIndex: 2, StartTime: "08:55", EndTime: "09:40"},
		},
		Entries: []services.TimetableEntryInput{
			{Weekday: 1, PeriodIndex: 1, SubjectName: "语文", TeacherName: strPtr("李老师"), Room: strPtr("101")},
			{Weekday: 1, PeriodIndex: 2, SubjectName: "数学", WeekType: "odd"},
		},
	}

	autoColored, err := svc.Save(class.ID, school.ID, payload)
	require.NoError(t, err)
	assert.Equal(t, 1, autoColored, "仅未指定颜色的语文需要自动配色")

	data, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	require.Len(t, data.Subjects, 2)
	assert.Equal(t, "语文", data.Subjects[0].Name)
	assert.True(t, strings.HasPrefix(data.Subjects[0].Color, "#"), "自动配色应写入库")
	assert.Equal(t, "数", data.Subjects[1].SimplifiedName)

	require.Len(t, data.Periods, 2)
	assert.Equal(t, "第一节", data.Periods[0].Name)
	assert.Equal(t, "第2节", data.Periods[1].Name, "未传 name 时按节次生成默认名")

	require.Len(t, data.Entries, 2)
	assert.Equal(t, 1, data.Entries[0].Weekday)
	require.NotNil(t, data.Entries[0].SubjectName)
	assert.Equal(t, "语文", *data.Entries[0].SubjectName)
	require.NotNil(t, data.Entries[0].TeacherName)
	assert.Equal(t, "李老师", *data.Entries[0].TeacherName)
	require.NotNil(t, data.Entries[0].Room)
	assert.Equal(t, "101", *data.Entries[0].Room)
	assert.Equal(t, "all", data.Entries[0].WeekType)
	assert.Equal(t, "odd", data.Entries[1].WeekType)
	assert.Nil(t, data.Entries[1].TeacherName, "未传教师应为 null")
}

func TestSave_KeepsExistingColorAndSimplified(t *testing.T) {
	db := setupDB(t)
	school, class, _ := seedTimetable(t, db)
	svc := newTimetableService(db)

	_, err := svc.Save(class.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "英语", Color: strPtr("#ABCDEF"), SimplifiedName: strPtr("英")}},
	})
	require.NoError(t, err)

	// 第二次不传 color / simplified_name：保留库中已有值，且不再计入自动配色。
	autoColored, err := svc.Save(class.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "英语"}},
	})
	require.NoError(t, err)
	assert.Equal(t, 0, autoColored)

	data, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	require.Len(t, data.Subjects, 1)
	assert.Equal(t, "#ABCDEF", data.Subjects[0].Color)
	assert.Equal(t, "英", data.Subjects[0].SimplifiedName)
}

func TestSave_SkipsInvalidEntries(t *testing.T) {
	db := setupDB(t)
	school, class, _ := seedTimetable(t, db)
	svc := newTimetableService(db)

	_, err := svc.Save(class.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
		Entries: []services.TimetableEntryInput{
			{Weekday: 1, PeriodIndex: 1, SubjectName: "不存在科目"},
			{Weekday: 0, PeriodIndex: 1, SubjectName: "语文"},
			{Weekday: 8, PeriodIndex: 1, SubjectName: "语文"},
			{Weekday: 1, PeriodIndex: 0, SubjectName: "语文"},
			{Weekday: 2, PeriodIndex: 1, SubjectName: "语文", WeekType: "invalid"},
		},
	})
	require.NoError(t, err)

	data, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	require.Len(t, data.Entries, 1, "非法行应被跳过")
	assert.Equal(t, 2, data.Entries[0].Weekday)
	assert.Equal(t, "all", data.Entries[0].WeekType, "非法周次回退为 all")
}

func TestSave_ReplacesEntries(t *testing.T) {
	db := setupDB(t)
	school, class, _ := seedTimetable(t, db)
	svc := newTimetableService(db)

	_, err := svc.Save(class.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}, {Name: "数学"}},
		Entries: []services.TimetableEntryInput{
			{Weekday: 1, PeriodIndex: 1, SubjectName: "语文"},
			{Weekday: 1, PeriodIndex: 2, SubjectName: "数学"},
		},
	})
	require.NoError(t, err)

	// 再次保存只留一条：旧排课整体被覆盖。
	_, err = svc.Save(class.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}, {Name: "数学"}},
		Entries: []services.TimetableEntryInput{
			{Weekday: 3, PeriodIndex: 4, SubjectName: "数学"},
		},
	})
	require.NoError(t, err)

	data, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	require.Len(t, data.Entries, 1)
	assert.Equal(t, 3, data.Entries[0].Weekday)
	require.NotNil(t, data.Entries[0].SubjectName)
	assert.Equal(t, "数学", *data.Entries[0].SubjectName)
}

func TestSave_ScopedPerClass(t *testing.T) {
	db := setupDB(t)
	school, class, teacher := seedTimetable(t, db)
	other, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "二班")
	svc := newTimetableService(db)

	payload := services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
		Entries:  []services.TimetableEntryInput{{Weekday: 1, PeriodIndex: 1, SubjectName: "语文"}},
	}
	_, err := svc.Save(class.ID, school.ID, payload)
	require.NoError(t, err)

	otherData, err := svc.Bootstrap(other.ID, school.ID)
	require.NoError(t, err)
	assert.Empty(t, otherData.Entries, "保存某班不应影响其他班级")

	thisData, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	assert.Len(t, thisData.Entries, 1)
}

func TestSubmitChange_ListForClass(t *testing.T) {
	db := setupDB(t)
	school, class, teacher := seedTimetable(t, db)
	svc := newTimetableService(db)

	payload := services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
		Entries: []services.TimetableEntryInput{
			{Weekday: 1, PeriodIndex: 1, SubjectName: "语文"},
			{Weekday: 2, PeriodIndex: 3, SubjectName: "语文"},
		},
	}
	change, err := svc.SubmitChange(class.ID, school.ID, teacher.ID, payload)
	require.NoError(t, err)
	assert.Equal(t, models.TimetableChangePending, change.Status)
	assert.Equal(t, 2, change.EntryCount)

	// 提交不直接生效：课表仍为空。
	data, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	assert.Empty(t, data.Entries, "申请未批准前不应改动课表")

	views, err := svc.ListChangesForClass(class.ID)
	require.NoError(t, err)
	require.Len(t, views, 1)
	assert.Equal(t, "pending", views[0].Status)
	assert.Equal(t, 2, views[0].EntryCount)
	assert.NotEmpty(t, views[0].CreatedAt)
	require.NotNil(t, views[0].RequesterName)
	assert.Equal(t, teacher.Name, *views[0].RequesterName)
	assert.Empty(t, views[0].ReviewedAt)
}

func TestApproveChange_AppliesAndSupersedes(t *testing.T) {
	db := setupDB(t)
	school, class, teacher := seedTimetable(t, db)
	admin := seedTeacher(t, db, school.ID, "admin1")
	svc := newTimetableService(db)

	mk := func(weekday int) *models.TimetableChangeRequest {
		change, err := svc.SubmitChange(class.ID, school.ID, teacher.ID, services.TimetablePayload{
			Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
			Entries:  []services.TimetableEntryInput{{Weekday: weekday, PeriodIndex: 1, SubjectName: "语文"}},
		})
		require.NoError(t, err)
		return change
	}
	first := mk(1)
	second := mk(2)

	applied, err := svc.ApproveChange(first.ID, admin.ID, strPtr("同意"))
	require.NoError(t, err)
	assert.True(t, applied)

	// 通过第一个申请：快照立即生效，其余 pending 自动作废。
	data, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	require.Len(t, data.Entries, 1)
	assert.Equal(t, 1, data.Entries[0].Weekday)

	var secondRow models.TimetableChangeRequest
	require.NoError(t, db.First(&secondRow, second.ID).Error)
	assert.Equal(t, models.TimetableChangeRejected, secondRow.Status)
	assert.NotNil(t, secondRow.ReviewNote)

	// 幂等：同一申请不能再通过。
	applied, err = svc.ApproveChange(first.ID, admin.ID, nil)
	require.NoError(t, err)
	assert.False(t, applied)

	views, err := svc.ListChangesForClass(class.ID)
	require.NoError(t, err)
	require.Len(t, views, 2)
	assert.Equal(t, "approved", views[1].Status)
	require.NotNil(t, views[1].ReviewerName)
	assert.Equal(t, admin.Name, *views[1].ReviewerName)
}

func TestApproveChange_NotFound(t *testing.T) {
	db := setupDB(t)
	svc := newTimetableService(db)

	applied, err := svc.ApproveChange(999, 1, nil)
	require.NoError(t, err)
	assert.False(t, applied)
}

func TestRejectChange(t *testing.T) {
	db := setupDB(t)
	school, class, teacher := seedTimetable(t, db)
	admin := seedTeacher(t, db, school.ID, "admin1")
	svc := newTimetableService(db)

	change, err := svc.SubmitChange(class.ID, school.ID, teacher.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
		Entries:  []services.TimetableEntryInput{{Weekday: 1, PeriodIndex: 1, SubjectName: "语文"}},
	})
	require.NoError(t, err)

	ok1, err := svc.RejectChange(change.ID, admin.ID, strPtr("请调整后再提交"))
	require.NoError(t, err)
	assert.True(t, ok1)

	data, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	assert.Empty(t, data.Entries, "驳回不应改动课表")

	// 幂等：已驳回不能再驳回，也不能再通过。
	ok2, err := svc.RejectChange(change.ID, admin.ID, nil)
	require.NoError(t, err)
	assert.False(t, ok2)

	applied, err := svc.ApproveChange(change.ID, admin.ID, nil)
	require.NoError(t, err)
	assert.False(t, applied)
}

func TestApproveChange_AppliesSubjectChanges(t *testing.T) {
	db := setupDB(t)
	school, class, teacher := seedTimetable(t, db)
	admin := seedTeacher(t, db, school.ID, "admin1")
	svc := newTimetableService(db)

	change, err := svc.SubmitChange(class.ID, school.ID, teacher.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "体育"}, {Name: "美术"}},
		Periods:  []services.TimetablePeriodInput{{PeriodIndex: 1, Name: "第一节", StartTime: "08:00", EndTime: "08:45"}},
		Entries: []services.TimetableEntryInput{
			{Weekday: 5, PeriodIndex: 1, SubjectName: "体育"},
			{Weekday: 5, PeriodIndex: 1, SubjectName: "美术", WeekType: "even"},
		},
	})
	require.NoError(t, err)

	applied, err := svc.ApproveChange(change.ID, admin.ID, nil)
	require.NoError(t, err)
	assert.True(t, applied)

	data, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	assert.Len(t, data.Subjects, 2, "申请中的科目应一并写入")
	assert.Len(t, data.Periods, 1)
	assert.Len(t, data.Entries, 2, "同格不同周次的两条排课应并存")
}

func TestAdminSave_RejectsPendingChanges(t *testing.T) {
	db := setupDB(t)
	school, class, teacher := seedTimetable(t, db)
	svc := newTimetableService(db)

	change, err := svc.SubmitChange(class.ID, school.ID, teacher.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
		Entries:  []services.TimetableEntryInput{{Weekday: 1, PeriodIndex: 1, SubjectName: "语文"}},
	})
	require.NoError(t, err)

	require.NoError(t, svc.AdminSave(class.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "数学"}},
		Entries:  []services.TimetableEntryInput{{Weekday: 4, PeriodIndex: 2, SubjectName: "数学"}},
	}))

	var row models.TimetableChangeRequest
	require.NoError(t, db.First(&row, change.ID).Error)
	assert.Equal(t, models.TimetableChangeRejected, row.Status, "管理员直改后旧申请自动作废")
	require.NotNil(t, row.ReviewNote)
	assert.Contains(t, *row.ReviewNote, "管理员已直接修改")

	data, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	require.Len(t, data.Entries, 1)
	assert.Equal(t, 4, data.Entries[0].Weekday, "管理员直改即时生效")
}

func TestListChangesForSchool_FilterAndClassInfo(t *testing.T) {
	db := setupDB(t)
	school, class, teacher := seedTimetable(t, db)
	admin := seedTeacher(t, db, school.ID, "admin1")
	svc := newTimetableService(db)

	first, err := svc.SubmitChange(class.ID, school.ID, teacher.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
		Entries:  []services.TimetableEntryInput{{Weekday: 1, PeriodIndex: 1, SubjectName: "语文"}},
	})
	require.NoError(t, err)
	_, err = svc.SubmitChange(class.ID, school.ID, teacher.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
		Entries:  []services.TimetableEntryInput{{Weekday: 2, PeriodIndex: 1, SubjectName: "语文"}},
	})
	require.NoError(t, err)
	_, err = svc.RejectChange(first.ID, admin.ID, nil)
	require.NoError(t, err)

	all, err := svc.ListChangesForSchool(school.ID, "")
	require.NoError(t, err)
	require.Len(t, all, 2)
	assert.Equal(t, first.ID+1, all[0].ID, "最新在前")
	require.NotNil(t, all[0].ClassName)
	assert.Equal(t, class.Name, *all[0].ClassName)

	pending, err := svc.ListChangesForSchool(school.ID, "pending")
	require.NoError(t, err)
	require.Len(t, pending, 1)
	assert.Equal(t, models.TimetableChangePending, pending[0].Status)

	rejected, err := svc.ListChangesForSchool(school.ID, "rejected")
	require.NoError(t, err)
	require.Len(t, rejected, 1)
	assert.Equal(t, first.ID, rejected[0].ID)
	require.NotNil(t, rejected[0].ReviewedAt)
}

func TestListChangesForSchool_OtherSchoolIsolated(t *testing.T) {
	db := setupDB(t)
	school, class, teacher := seedTimetable(t, db)
	svc := newTimetableService(db)

	_, err := svc.SubmitChange(class.ID, school.ID, teacher.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
	})
	require.NoError(t, err)

	other, err := svc.ListChangesForSchool(school.ID+1, "")
	require.NoError(t, err)
	assert.Empty(t, other, "其他学校不应看到本校申请")
}

func TestForTeacher(t *testing.T) {
	db := setupDB(t)
	school, class, teacher := seedTimetable(t, db)
	second, _ := seedTeacherClass(t, db, school.ID, teacher.ID, "二班")
	svc := newTimetableService(db)

	_, err := svc.Save(class.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
		Periods:  []services.TimetablePeriodInput{{PeriodIndex: 1, StartTime: "08:00", EndTime: "08:45"}},
		Entries: []services.TimetableEntryInput{
			{Weekday: 1, PeriodIndex: 1, SubjectName: "语文", TeacherName: strPtr("李老师"), Room: strPtr("101")},
			{Weekday: 2, PeriodIndex: 1, SubjectName: "语文", TeacherName: strPtr("王老师")},
		},
	})
	require.NoError(t, err)

	// 另一班的同一位教师课程
	_, err = svc.Save(second.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "语文"}},
		Entries: []services.TimetableEntryInput{
			{Weekday: 3, PeriodIndex: 1, SubjectName: "语文", TeacherName: strPtr("李老师")},
			{Weekday: 4, PeriodIndex: 1, SubjectName: "语文", TeacherName: strPtr("赵老师")},
		},
	})
	require.NoError(t, err)

	schedule, err := svc.ForTeacher(school.ID, "李老师")
	require.NoError(t, err)
	assert.Equal(t, "李老师", schedule.TeacherName)
	require.Len(t, schedule.Subjects, 1)
	require.Len(t, schedule.Periods, 1)
	require.Len(t, schedule.Entries, 2, "仅返回该教师的课")
	assert.Equal(t, class.ID, schedule.Entries[0].ClassID)
	require.NotNil(t, schedule.Entries[0].ClassName)
	assert.Equal(t, class.Name, *schedule.Entries[0].ClassName)
	assert.Equal(t, 3, schedule.Entries[1].Weekday)
	require.NotNil(t, schedule.Entries[1].ClassName)
	assert.Equal(t, second.Name, *schedule.Entries[1].ClassName)

	none, err := svc.ForTeacher(school.ID, "查无此人")
	require.NoError(t, err)
	assert.Empty(t, none.Entries)
}

func TestToCses(t *testing.T) {
	db := setupDB(t)
	school, class, _ := seedTimetable(t, db)
	svc := newTimetableService(db)

	_, err := svc.Save(class.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{
			{Name: "语文", SimplifiedName: strPtr("语")},
			{Name: "数学"},
		},
		Periods: []services.TimetablePeriodInput{
			{PeriodIndex: 1, StartTime: "08:00", EndTime: "08:45"},
			{PeriodIndex: 2, StartTime: "08:55", EndTime: "09:40"},
		},
		Entries: []services.TimetableEntryInput{
			{Weekday: 1, PeriodIndex: 1, SubjectName: "语文", TeacherName: strPtr("李老师"), Room: strPtr("101")},
			{Weekday: 1, PeriodIndex: 2, SubjectName: "数学", WeekType: "odd"},
		},
	})
	require.NoError(t, err)

	yaml, err := svc.ToCses(class.ID, school.ID)
	require.NoError(t, err)

	assert.True(t, strings.HasSuffix(yaml, "\n"))
	assert.Contains(t, yaml, "version: 1")
	assert.Contains(t, yaml, "subjects:")
	assert.Contains(t, yaml, "  - name: \"语文\"")
	assert.Contains(t, yaml, "    simplified_name: \"语\"")
	assert.Contains(t, yaml, "schedules:")
	assert.Contains(t, yaml, "  - name: \"星期一\"\n    enable_day: 1\n    weeks: all")
	assert.Contains(t, yaml, "  - name: \"星期一·单周\"\n    enable_day: 1\n    weeks: odd")
	assert.Contains(t, yaml, "      - subject: \"语文\"")
	assert.Contains(t, yaml, "        start_time: \"08:00\"")
	assert.Contains(t, yaml, "        end_time: \"08:45\"")
	assert.Contains(t, yaml, "        teacher: \"李老师\"")
	assert.Contains(t, yaml, "        room: \"101\"")

	// 单双周分组顺序固定为 all → odd → even，且星期一排在星期二之前。
	assert.Less(t, strings.Index(yaml, "weeks: all"), strings.Index(yaml, "weeks: odd"))
}

func TestToCses_EmptyTimetable(t *testing.T) {
	db := setupDB(t)
	school, class, _ := seedTimetable(t, db)
	svc := newTimetableService(db)

	yaml, err := svc.ToCses(class.ID, school.ID)
	require.NoError(t, err)
	assert.Contains(t, yaml, "subjects:\n  []")
	assert.Contains(t, yaml, "schedules:\n  []")
}

func TestToCses_EscapesQuotes(t *testing.T) {
	db := setupDB(t)
	school, class, _ := seedTimetable(t, db)
	svc := newTimetableService(db)

	_, err := svc.Save(class.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: `语"文\`}},
	})
	require.NoError(t, err)

	yaml, err := svc.ToCses(class.ID, school.ID)
	require.NoError(t, err)
	assert.Contains(t, yaml, `  - name: "语\"文\\"`, "双引号与反斜杠都应转义")
}
