// 课表进阶服务层测试：CSV 导入（预览 / 落库 / 覆盖 / 报错）、不可用时段 replace、
// 冲突检查（教师跨班 + 单双周 + 不可用）、任课设置 replace。
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

const timetableCsvHeader = "年级,班级,星期,第几节,开始时间,结束时间,科目,周次,教师,教室"

// timetableCsv 拼接 CSV 文本。
func timetableCsv(rows ...string) []byte {
	all := append([]string{timetableCsvHeader}, rows...)
	return []byte(strings.Join(all, "\n") + "\n")
}

// newTimetableAdvancedFixture 学校 + 一个带年级的班级。
func newTimetableAdvancedFixture(t *testing.T, db *gorm.DB) (models.School, models.ClassRoom) {
	t.Helper()
	school := seedSchool(t, db)
	class := makeClass(t, db, school.ID, "一年级（1）班", "一年级")
	return school, class
}

// CSV dry_run：只解析校验不落库，返回汇总与错误行。
func TestTimetableImportCsvDryRun(t *testing.T) {
	db := setupDB(t)
	school, _ := newTimetableAdvancedFixture(t, db)
	svc := newTimetableService(db)

	content := timetableCsv(
		"一年级,一年级（1）班,星期一,1,08:00,08:45,语文,每周,张老师,101",
		"一年级,一年级（1）班,星期二,2,08:55,09:40,数学,单周,李老师,102",
		"一年级,一年级（1）班,星期八,0,,,,", // 非法行
	)

	summary, err := svc.ImportFromCsv(content, school.ID, true)
	require.NoError(t, err)
	assert.True(t, summary.DryRun)
	assert.False(t, summary.Imported)
	assert.Equal(t, 3, summary.TotalRows)
	assert.Equal(t, 2, summary.PeriodCount)
	require.Len(t, summary.Classes, 1)
	assert.Equal(t, "一年级", summary.Classes[0].Grade)
	assert.Equal(t, "一年级（1）班", summary.Classes[0].Name)
	assert.True(t, summary.Classes[0].Found)
	assert.Equal(t, 2, summary.Classes[0].EntryCount)
	require.Len(t, summary.Errors, 1)
	assert.Equal(t, "第4行：星期 / 第几节 / 科目无效（星期=星期八 节次=0 科目=空），已跳过", summary.Errors[0])
	assert.Nil(t, summary.SubjectsAutoColored, "预览模式不返回自动配色数")

	// 未落库。
	var entryCount, subjectCount, periodCount int64
	require.NoError(t, db.Model(&models.TimetableEntry{}).Count(&entryCount).Error)
	require.NoError(t, db.Model(&models.Subject{}).Count(&subjectCount).Error)
	require.NoError(t, db.Model(&models.ClassPeriod{}).Count(&periodCount).Error)
	assert.Equal(t, int64(0), entryCount)
	assert.Equal(t, int64(0), subjectCount)
	assert.Equal(t, int64(0), periodCount)
}

// CSV 正式导入：覆盖该班课表、写科目/节次、作废待审申请。
func TestTimetableImportCsvCommit(t *testing.T) {
	db := setupDB(t)
	school, class := newTimetableAdvancedFixture(t, db)
	svc := newTimetableService(db)

	// 先留一条旧排课与一份待审申请，验证「整班覆盖 + 申请自动作废」。
	_, err := svc.Save(class.ID, school.ID, services.TimetablePayload{
		Subjects: []services.TimetableSubjectInput{{Name: "英语"}},
		Entries:  []services.TimetableEntryInput{{Weekday: 5, PeriodIndex: 5, SubjectName: "英语"}},
	})
	require.NoError(t, err)
	change := models.TimetableChangeRequest{
		SchoolID: school.ID, ClassID: class.ID, RequestedBy: 1,
		Payload: "{}", EntryCount: 1, Status: models.TimetableChangePending,
	}
	require.NoError(t, db.Create(&change).Error)

	content := timetableCsv(
		"一年级,一年级（1）班,星期一,1,08:00,08:45,语文,每周,张老师,101",
		"一年级,一年级（1）班,星期二,2,08:55,09:40,数学,单周,李老师,102",
	)
	summary, err := svc.ImportFromCsv(content, school.ID, false)
	require.NoError(t, err)
	assert.True(t, summary.Imported)
	require.NotNil(t, summary.SubjectsAutoColored)
	assert.Equal(t, 2, *summary.SubjectsAutoColored, "两个新科目自动配色")
	assert.Empty(t, summary.Errors)

	data, err := svc.Bootstrap(class.ID, school.ID)
	require.NoError(t, err)
	require.Len(t, data.Entries, 2, "旧排课被整班覆盖")
	assert.Len(t, data.Periods, 2)
	assert.GreaterOrEqual(t, len(data.Subjects), 3, "旧科目保留，新科目补入")

	byWeekday := map[int]services.TimetableEntryView{}
	for _, entry := range data.Entries {
		byWeekday[entry.Weekday] = entry
	}
	monday := byWeekday[1]
	require.NotNil(t, monday.SubjectName)
	assert.Equal(t, "语文", *monday.SubjectName)
	assert.Equal(t, 1, monday.PeriodIndex)
	assert.Equal(t, "all", monday.WeekType)
	require.NotNil(t, monday.TeacherName)
	assert.Equal(t, "张老师", *monday.TeacherName)
	require.NotNil(t, monday.Room)
	assert.Equal(t, "101", *monday.Room)

	tuesday := byWeekday[2]
	require.NotNil(t, tuesday.SubjectName)
	assert.Equal(t, "数学", *tuesday.SubjectName)
	assert.Equal(t, "odd", tuesday.WeekType, "「单周」→ odd")

	// 待审申请已自动作废。
	var reloaded models.TimetableChangeRequest
	require.NoError(t, db.First(&reloaded, change.ID).Error)
	assert.Equal(t, models.TimetableChangeRejected, reloaded.Status)
	require.NotNil(t, reloaded.ReviewNote)
	assert.Equal(t, "管理员批量导入课表，本申请自动作废", *reloaded.ReviewNote)
	require.NotNil(t, reloaded.ReviewedAt)
}

// CSV 中班级不存在：报错并跳过该班（不自动建班、不落库）。
func TestTimetableImportCsvMissingClass(t *testing.T) {
	db := setupDB(t)
	school, _ := newTimetableAdvancedFixture(t, db)
	svc := newTimetableService(db)

	content := timetableCsv("三年级,三年级（1）班,星期一,1,08:00,08:45,语文,每周,张老师,101")
	summary, err := svc.ImportFromCsv(content, school.ID, false)
	require.NoError(t, err)
	assert.Equal(t, 1, summary.TotalRows)
	require.Len(t, summary.Classes, 1)
	assert.False(t, summary.Classes[0].Found)
	assert.Equal(t, 1, summary.Classes[0].EntryCount)
	assert.Contains(t, summary.Errors, "班级「三年级/三年级（1）班」不存在，该班 1 行未导入（请先在班级列表创建）")

	var entryCount int64
	require.NoError(t, db.Model(&models.TimetableEntry{}).Count(&entryCount).Error)
	assert.Equal(t, int64(0), entryCount, "未匹配班级不落库")

	var classCount int64
	require.NoError(t, db.Model(&models.ClassRoom{}).Where("name = ?", "三年级（1）班").Count(&classCount).Error)
	assert.Equal(t, int64(0), classCount, "不自动建班")
}

// CSV 缺必需列 / 非 UTF-8：直接给出明确错误。
func TestTimetableImportCsvInvalidInput(t *testing.T) {
	db := setupDB(t)
	school, _ := newTimetableAdvancedFixture(t, db)
	svc := newTimetableService(db)

	summary, err := svc.ImportFromCsv([]byte("年级,班级\n一年级,一年级（1）班\n"), school.ID, true)
	require.NoError(t, err)
	require.Len(t, summary.Errors, 1)
	assert.Equal(t, "CSV 缺少必需列（年级 / 班级 / 星期 / 第几节 / 科目），请使用模板", summary.Errors[0])
	assert.Equal(t, 0, summary.TotalRows)
	assert.Empty(t, summary.Classes)

	_, err = svc.ImportFromCsv([]byte{0xEF, 0xBB, 0xBF, 0xFF, 0xFE, 0x00}, school.ID, true)
	assert.Equal(t, "文件编码不是 UTF-8，请另存为 UTF-8 后重试", appErrorOf(t, err, 422))
}

// 不可用时段：replace 语义（先清后插）+ 非法格子跳过 + 不串教师。
func TestTimetableUnavailabilitiesReplace(t *testing.T) {
	db := setupDB(t)
	school, _ := newTimetableAdvancedFixture(t, db)
	svc := newTimetableService(db)

	require.NoError(t, svc.SaveUnavailabilities(school.ID, "张老师", []services.UnavailabilityCell{
		{Weekday: 1, PeriodIndex: 1},
		{Weekday: 2, PeriodIndex: 3},
		{Weekday: 0, PeriodIndex: 5}, // 非法星期跳过
		{Weekday: 3, PeriodIndex: 0}, // 非法节次跳过
		{Weekday: 2, PeriodIndex: 3}, // 重复格子去重
	}))
	require.NoError(t, svc.SaveUnavailabilities(school.ID, "李老师", []services.UnavailabilityCell{{Weekday: 3, PeriodIndex: 3}}))

	rows, err := svc.ListUnavailabilities(school.ID)
	require.NoError(t, err)
	require.Len(t, rows, 3)

	// 整体覆盖张老师：只剩 1 条，李老师不受影响。
	require.NoError(t, svc.SaveUnavailabilities(school.ID, "张老师", []services.UnavailabilityCell{{Weekday: 5, PeriodIndex: 5}}))
	rows, err = svc.ListUnavailabilities(school.ID)
	require.NoError(t, err)
	require.Len(t, rows, 2)

	zhang := 0
	li := 0
	for _, row := range rows {
		switch row.TeacherName {
		case "张老师":
			zhang++
			assert.Equal(t, 5, row.Weekday)
			assert.Equal(t, 5, row.PeriodIndex)
		case "李老师":
			li++
			assert.Equal(t, 3, row.Weekday)
			assert.Equal(t, 3, row.PeriodIndex)
		}
	}
	assert.Equal(t, 1, zhang)
	assert.Equal(t, 1, li)

	// 清空。
	require.NoError(t, svc.SaveUnavailabilities(school.ID, "张老师", nil))
	rows, err = svc.ListUnavailabilities(school.ID)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, "李老师", rows[0].TeacherName)
}

// conflictEntry 构造冲突检查输入。
func conflictEntry(weekday, period int, subject, weekType, teacher string) services.TimetableConflictEntry {
	return services.TimetableConflictEntry{
		Weekday: weekday, PeriodIndex: period, SubjectName: subject,
		WeekType: weekType, TeacherName: teacher,
	}
}

// 教师跨班冲突：同格子且周次重叠才算；同班格子与其他学校不计入。
func TestTimetableCheckConflictsTeacher(t *testing.T) {
	db := setupDB(t)
	school, target := newTimetableAdvancedFixture(t, db)
	other := makeClass(t, db, school.ID, "二班", "一年级")
	svc := newTimetableService(db)

	subject := models.Subject{SchoolID: school.ID, Name: "语文", Color: "#5B8FF9"}
	require.NoError(t, db.Create(&subject).Error)

	otherEntry := models.TimetableEntry{
		ClassID: other.ID, Weekday: 1, PeriodIndex: 1, WeekType: "all",
		SubjectID: &subject.ID, TeacherName: "张老师",
	}
	require.NoError(t, db.Create(&otherEntry).Error)

	// 本班已存在的同格子排课不参与自身冲突判定。
	require.NoError(t, db.Create(&models.TimetableEntry{
		ClassID: target.ID, Weekday: 1, PeriodIndex: 1, WeekType: "all",
		SubjectID: &subject.ID, TeacherName: "张老师",
	}).Error)

	// 他校同教师同格子不参与本校区冲突判定。
	otherSchool := seedOtherSchool(t, db)
	foreignClass := makeClass(t, db, otherSchool.ID, "外校一班", "一年级")
	require.NoError(t, db.Create(&models.TimetableEntry{
		ClassID: foreignClass.ID, Weekday: 1, PeriodIndex: 1, WeekType: "all",
		SubjectID: &subject.ID, TeacherName: "张老师",
	}).Error)

	conflicts, err := svc.CheckConflicts(school.ID, target.ID, []services.TimetableConflictEntry{
		conflictEntry(1, 1, "数学", "all", "张老师"),
	})
	require.NoError(t, err)
	require.Len(t, conflicts, 1, "只统计其他班级在本校的排课")

	conflict := conflicts[0]
	assert.Equal(t, "teacher", conflict.Type)
	assert.Equal(t, "张老师", conflict.TeacherName)
	assert.Equal(t, 1, conflict.Weekday)
	assert.Equal(t, 1, conflict.PeriodIndex)
	assert.Equal(t, "all", conflict.WeekType)
	assert.Equal(t, "数学", conflict.SubjectName)
	require.NotNil(t, conflict.OtherClassID)
	assert.Equal(t, other.ID, *conflict.OtherClassID)
	require.NotNil(t, conflict.OtherClassName)
	assert.Equal(t, "二班", *conflict.OtherClassName)
	require.NotNil(t, conflict.OtherSubjectName)
	assert.Equal(t, "语文", *conflict.OtherSubjectName)
	assert.Equal(t, "星期一第1节「数学」与 二班「语文」冲突（同一教师）", conflict.Message)

	// 周次重叠规则：all 与任何都冲突；odd 与 even 不冲突。
	require.NoError(t, db.Model(&models.TimetableEntry{}).Where("id = ?", otherEntry.ID).Update("week_type", "odd").Error)

	same, err := svc.CheckConflicts(school.ID, target.ID, []services.TimetableConflictEntry{
		conflictEntry(1, 1, "数学", "even", "张老师"),
	})
	require.NoError(t, err)
	assert.Empty(t, same, "odd 与 even 不重叠")

	overlap, err := svc.CheckConflicts(school.ID, target.ID, []services.TimetableConflictEntry{
		conflictEntry(1, 1, "数学", "all", "张老师"),
		conflictEntry(1, 1, "英语", "odd", "张老师"),
	})
	require.NoError(t, err)
	require.Len(t, overlap, 2, "all 与 odd 都命中 odd 班级排课")
	assert.Equal(t, "星期一单周第1节「英语」与 二班「语文」冲突（同一教师）", overlap[1].Message)

	// 无教师名 → 不报冲突。
	empty, err := svc.CheckConflicts(school.ID, target.ID, []services.TimetableConflictEntry{
		conflictEntry(1, 1, "数学", "all", ""),
	})
	require.NoError(t, err)
	assert.Empty(t, empty)
	assert.NotNil(t, empty, "返回空数组而不是 nil")
}

// 不可用时段冲突。
func TestTimetableCheckConflictsUnavailable(t *testing.T) {
	db := setupDB(t)
	school, target := newTimetableAdvancedFixture(t, db)
	svc := newTimetableService(db)

	require.NoError(t, svc.SaveUnavailabilities(school.ID, "张老师", []services.UnavailabilityCell{{Weekday: 3, PeriodIndex: 2}}))

	conflicts, err := svc.CheckConflicts(school.ID, target.ID, []services.TimetableConflictEntry{
		conflictEntry(3, 2, "英语", "all", "张老师"),
	})
	require.NoError(t, err)
	require.Len(t, conflicts, 1)

	conflict := conflicts[0]
	assert.Equal(t, "unavailable", conflict.Type)
	assert.Equal(t, "张老师", conflict.TeacherName)
	assert.Equal(t, 3, conflict.Weekday)
	assert.Equal(t, 2, conflict.PeriodIndex)
	assert.Equal(t, "英语", conflict.SubjectName)
	assert.Nil(t, conflict.OtherClassID)
	assert.Nil(t, conflict.OtherClassName)
	assert.Nil(t, conflict.OtherSubjectName)
	assert.Equal(t, "星期三 第2节「英语」：教师 张老师 此时段已被标记为不可用", conflict.Message)

	// 科目名为空时文案回退为「课程」。
	fallback, err := svc.CheckConflicts(school.ID, target.ID, []services.TimetableConflictEntry{
		conflictEntry(3, 2, "", "all", "张老师"),
	})
	require.NoError(t, err)
	require.Len(t, fallback, 1)
	assert.Contains(t, fallback[0].Message, "第2节「课程」")
}

// 任课设置：replace 语义 + 空值跳过 + 读取排序。
func TestTimetableAssignmentsReplace(t *testing.T) {
	db := setupDB(t)
	school, class := newTimetableAdvancedFixture(t, db)
	svc := newTimetableService(db)

	require.NoError(t, svc.SaveAssignments(class.ID, school.ID, []services.TimetableAssignmentRow{
		{SubjectName: "数学", TeacherName: "李老师"},
		{SubjectName: "语文", TeacherName: "张老师"},
		{SubjectName: "", TeacherName: "空科目"},
		{SubjectName: "空教师", TeacherName: " "},
		{SubjectName: "数学", TeacherName: "重复科目"},
	}))

	rows, err := svc.ListAssignments(class.ID)
	require.NoError(t, err)
	require.Len(t, rows, 2, "空值跳过、同科目去重")
	got := map[string]string{}
	for _, row := range rows {
		got[row.SubjectName] = row.TeacherName
	}
	assert.Equal(t, "李老师", got["数学"])
	assert.Equal(t, "张老师", got["语文"])

	// replace：以本次提交为准。
	require.NoError(t, svc.SaveAssignments(class.ID, school.ID, []services.TimetableAssignmentRow{
		{SubjectName: "英语", TeacherName: "王老师"},
	}))
	rows, err = svc.ListAssignments(class.ID)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, "英语", rows[0].SubjectName)
	assert.Equal(t, "王老师", rows[0].TeacherName)

	// 清空。
	require.NoError(t, svc.SaveAssignments(class.ID, school.ID, nil))
	rows, err = svc.ListAssignments(class.ID)
	require.NoError(t, err)
	assert.Empty(t, rows)

	// 另一班的任课不受影响。
	other := makeClass(t, db, school.ID, "二班", "一年级")
	require.NoError(t, svc.SaveAssignments(other.ID, school.ID, []services.TimetableAssignmentRow{
		{SubjectName: "体育", TeacherName: "赵老师"},
	}))
	require.NoError(t, svc.SaveAssignments(class.ID, school.ID, []services.TimetableAssignmentRow{
		{SubjectName: "英语", TeacherName: "王老师"},
	}))
	otherRows, err := svc.ListAssignments(other.ID)
	require.NoError(t, err)
	require.Len(t, otherRows, 1)
	assert.Equal(t, "体育", otherRows[0].SubjectName)
}

// ClassInSchool：管理员越权防护。
func TestTimetableClassInSchool(t *testing.T) {
	db := setupDB(t)
	school, class := newTimetableAdvancedFixture(t, db)
	otherSchool := seedOtherSchool(t, db)
	foreign := makeClass(t, db, otherSchool.ID, "外校一班", "一年级")
	svc := newTimetableService(db)

	ok, err := svc.ClassInSchool(school.ID, class.ID)
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = svc.ClassInSchool(school.ID, foreign.ID)
	require.NoError(t, err)
	assert.False(t, ok)

	ok, err = svc.ClassInSchool(school.ID, 99999)
	require.NoError(t, err)
	assert.False(t, ok)
}
