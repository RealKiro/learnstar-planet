// 课表自动排课测试：单班规则排课（纯计算）+ 全校智能排课（预览 / 落库 / 教师冲突）。
// 全部为不变式测试：成功场景多轮断言（≥20 轮），不依赖某个具体随机结果。
package services_test

import (
	"sort"
	"strconv"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// newSeededTimetableService 课表服务 + 确定性但逐次变化的种子源（断言不变式而非具体随机结果）。
func newSeededTimetableService(db *gorm.DB) *services.TimetableService {
	seed := int64(0)
	svc := services.NewTimetableService(db)
	svc.SetSeedSource(func() int64 {
		seed++
		return seed * 7919
	})
	return svc
}

// seedClassPeriods 建学校级节次作息（节次名自动为「第 N 节」）。
func seedClassPeriods(t *testing.T, db *gorm.DB, schoolID uint, specs []services.TimetablePeriodInput) {
	t.Helper()
	for _, spec := range specs {
		period := models.ClassPeriod{
			SchoolID:    schoolID,
			PeriodIndex: spec.PeriodIndex,
			Name:        "第" + strconv.Itoa(spec.PeriodIndex) + "节",
			StartTime:   spec.StartTime,
			EndTime:     spec.EndTime,
		}
		require.NoError(t, db.Create(&period).Error)
	}
}

// amPmPeriods 六节：第 1-4 节上午（第 4 节 11:45 结束），第 5-6 节下午。
func amPmPeriods() []services.TimetablePeriodInput {
	return []services.TimetablePeriodInput{
		{PeriodIndex: 1, StartTime: "08:00", EndTime: "08:45"},
		{PeriodIndex: 2, StartTime: "08:55", EndTime: "09:40"},
		{PeriodIndex: 3, StartTime: "10:00", EndTime: "10:45"},
		{PeriodIndex: 4, StartTime: "11:00", EndTime: "11:45"},
		{PeriodIndex: 5, StartTime: "13:30", EndTime: "14:15"},
		{PeriodIndex: 6, StartTime: "14:25", EndTime: "15:10"},
	}
}

// cellsBySubject 某班课表格子中指定科目占用的节次（升序）。
func cellsBySubject(grid map[int]map[int]string, day int, subject string) []int {
	cells := []int{}
	for period, name := range grid[day] {
		if name == subject {
			cells = append(cells, period)
		}
	}
	sort.Ints(cells)
	return cells
}

// loadEntries 读取某班落库的排课（按星期 / 节次升序）。
func loadEntries(t *testing.T, db *gorm.DB, classID uint) []models.TimetableEntry {
	t.Helper()
	var entries []models.TimetableEntry
	require.NoError(t, db.Where("class_id = ?", classID).Order("weekday ASC, period_index ASC").Find(&entries).Error)
	return entries
}

// ① 单班可行规则：只检验不变式，多轮（25 轮）验证算法稳定。
func TestTimetableGenerateSingleClassInvariants(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	seedClassPeriods(t, db, school.ID, amPmPeriods())
	svc := newSeededTimetableService(db)

	rules := services.GenerateRules{
		Days: []int{1, 2, 3, 4, 5},
		Subjects: []services.GenerateRuleSubject{
			{Name: "语文", Weekly: 4, Double: true, MaxPerDay: 2},
			{Name: "数学", Weekly: 4, Session: "am", MaxPerDay: 1, ForbidPeriods: []int{1, 1}},
			{Name: "英语", Weekly: 3, Session: "pm", MaxPerDay: 2},
		},
	}

	const rounds = 25
	for round := 0; round < rounds; round++ {
		result, err := svc.Generate(school.ID, rules)
		require.NoError(t, err)
		require.True(t, result.Success, "第 %d 轮应排出课表，实际 warnings=%v", round, result.Warnings)
		assert.Empty(t, result.Warnings)
		require.Len(t, result.Entries, 11, "语文 4 + 数学 4 + 英语 3")

		grid := map[int]map[int]string{}
		perSubject := map[string]int{}
		perSubjectDay := map[string]map[int]int{}
		for _, entry := range result.Entries {
			assert.Equal(t, "all", entry.WeekType)
			assert.Nil(t, entry.TeacherName, "单班排课的教师名恒为 null")
			assert.Nil(t, entry.Room, "单班排课的教室恒为 null")
			assert.True(t, entry.Weekday >= 1 && entry.Weekday <= 5, "只排在上课日内")
			assert.True(t, entry.PeriodIndex >= 1 && entry.PeriodIndex <= 6)

			if grid[entry.Weekday] == nil {
				grid[entry.Weekday] = map[int]string{}
			}
			_, duplicated := grid[entry.Weekday][entry.PeriodIndex]
			require.False(t, duplicated, "同格被重复占用：星期 %d 第 %d 节", entry.Weekday, entry.PeriodIndex)
			grid[entry.Weekday][entry.PeriodIndex] = entry.SubjectName

			perSubject[entry.SubjectName]++
			if perSubjectDay[entry.SubjectName] == nil {
				perSubjectDay[entry.SubjectName] = map[int]int{}
			}
			perSubjectDay[entry.SubjectName][entry.Weekday]++
		}
		assert.Equal(t, 4, perSubject["语文"], "每科目排课节数 = weekly")
		assert.Equal(t, 4, perSubject["数学"], "每科目排课节数 = weekly")
		assert.Equal(t, 3, perSubject["英语"], "每科目排课节数 = weekly")

		// max_per_day 不被突破
		for subject, byDay := range perSubjectDay {
			for day, count := range byDay {
				limit := map[string]int{"语文": 2, "数学": 1, "英语": 2}[subject]
				assert.LessOrEqual(t, count, limit, "%s 星期 %d 超出每日上限", subject, day)
			}
		}

		// session 限定 + forbid_periods 从未被占用（上午 = 第 1-4 节，下午 = 第 5-6 节）
		for _, entry := range result.Entries {
			switch entry.SubjectName {
			case "数学":
				assert.Contains(t, []int{2, 3, 4}, entry.PeriodIndex, "数学限上午且禁排第 1 节")
			case "英语":
				assert.Contains(t, []int{5, 6}, entry.PeriodIndex, "英语限下午")
			}
		}

		// 连堂：语文每天 0 或 2 节，2 节必同日相邻且同一上 / 下午
		for day := 1; day <= 5; day++ {
			cells := cellsBySubject(grid, day, "语文")
			require.Contains(t, []int{0, 2}, len(cells), "连堂科目每天 0 或 2 节，实际 %v", cells)
			if len(cells) == 2 {
				assert.Equal(t, cells[0]+1, cells[1], "连堂两节必须相邻")
				assert.Equal(t, cells[0] <= 4, cells[1] <= 4, "连堂不得跨上 / 下午")
			}
		}
	}

	// 单班 generate 绝不落库
	var subjectCount, entryCount int64
	require.NoError(t, db.Model(&models.Subject{}).Count(&subjectCount).Error)
	require.NoError(t, db.Model(&models.TimetableEntry{}).Count(&entryCount).Error)
	assert.Equal(t, int64(0), subjectCount, "单班排课不写科目")
	assert.Equal(t, int64(0), entryCount, "单班排课不写排课")
}

// ② 单班边界与不可行：各类 warning 文案逐字对齐 + 不可能排下时 60 次重试后失败。
func TestTimetableGenerateSingleClassGuards(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newSeededTimetableService(db)

	weeklyOne := []services.GenerateRuleSubject{{Name: "语文", Weekly: 1}}

	// 无节次
	noPeriods, err := svc.Generate(school.ID, services.GenerateRules{Days: []int{1}, Subjects: weeklyOne})
	require.NoError(t, err)
	assert.False(t, noPeriods.Success)
	assert.Equal(t, []string{"请先在「节次作息」中设置每日节次"}, noPeriods.Warnings)
	assert.NotNil(t, noPeriods.Entries)
	assert.Empty(t, noPeriods.Entries)

	seedClassPeriods(t, db, school.ID, []services.TimetablePeriodInput{
		{PeriodIndex: 1, StartTime: "08:00", EndTime: "08:45"},
		{PeriodIndex: 2, StartTime: "08:55", EndTime: "09:40"},
	})

	// 未选上课日：空数组与「全为非法值」等价
	noDays, err := svc.Generate(school.ID, services.GenerateRules{Days: []int{}, Subjects: weeklyOne})
	require.NoError(t, err)
	assert.False(t, noDays.Success)
	assert.Equal(t, []string{"请至少选择一个上课日"}, noDays.Warnings)

	invalidDays, err := svc.Generate(school.ID, services.GenerateRules{Days: []int{0, 8}, Subjects: weeklyOne})
	require.NoError(t, err)
	assert.Equal(t, []string{"请至少选择一个上课日"}, invalidDays.Warnings)

	// 上课日缺省 → 默认 1-5
	defaultDays, err := svc.Generate(school.ID, services.GenerateRules{
		Subjects: []services.GenerateRuleSubject{{Name: "语文", Weekly: 2}},
	})
	require.NoError(t, err)
	require.True(t, defaultDays.Success, "缺省上课日应回退到 1-5：%v", defaultDays.Warnings)
	require.Len(t, defaultDays.Entries, 2)
	assert.NotEqual(t, defaultDays.Entries[0].Weekday, defaultDays.Entries[1].Weekday, "max_per_day 默认 1 → 分两天")

	// 所有科目 weekly < 1（含空名）→ 请为至少一个科目设置每周节数
	empty, err := svc.Generate(school.ID, services.GenerateRules{
		Days:     []int{1, 2},
		Subjects: []services.GenerateRuleSubject{{Name: "语文", Weekly: 0}, {Name: "  ", Weekly: 3}},
	})
	require.NoError(t, err)
	assert.False(t, empty.Success)
	assert.Equal(t, []string{"请为至少一个科目设置每周节数"}, empty.Warnings)

	// 所需节数超容量（并带出「每周节数超出每日上限」的前置 warning）
	over, err := svc.Generate(school.ID, services.GenerateRules{
		Days:     []int{1},
		Subjects: []services.GenerateRuleSubject{{Name: "语文", Weekly: 5}},
	})
	require.NoError(t, err)
	assert.False(t, over.Success)
	require.Len(t, over.Warnings, 2)
	assert.Equal(t, "科目「语文」每周 5 节，超出每日上限 1 × 1 天，无法排下", over.Warnings[0])
	assert.Equal(t, "所需节数 5 超过可用格数 2（1 天 × 2 节），请调整规则", over.Warnings[1])
	assert.Empty(t, over.Entries)

	// 连堂 + 禁排全部节次 → 无候选格，60 次重试后失败
	blocked, err := svc.Generate(school.ID, services.GenerateRules{
		Days:     []int{1},
		Subjects: []services.GenerateRuleSubject{{Name: "语文", Weekly: 2, Double: true, ForbidPeriods: []int{1, 2}}},
	})
	require.NoError(t, err)
	assert.False(t, blocked.Success)
	assert.Equal(t, []string{"尝试多次均无法排出满足全部规则的课表，请减少节数、放宽连堂或时段限制"}, blocked.Warnings)

	// session=am 但全校只有下午节次 → 同样落到重试失败
	otherSchool := seedOtherSchool(t, db)
	seedClassPeriods(t, db, otherSchool.ID, []services.TimetablePeriodInput{
		{PeriodIndex: 1, StartTime: "13:30", EndTime: "14:15"},
		{PeriodIndex: 2, StartTime: "14:25", EndTime: "15:10"},
	})
	pmOnly, err := svc.Generate(otherSchool.ID, services.GenerateRules{
		Days:     []int{1},
		Subjects: []services.GenerateRuleSubject{{Name: "数学", Weekly: 1, Session: "am"}},
	})
	require.NoError(t, err)
	assert.False(t, pmOnly.Success)
	assert.Equal(t, []string{"尝试多次均无法排出满足全部规则的课表，请减少节数、放宽连堂或时段限制"}, pmOnly.Warnings)
}

// ③ 全校排课的前置校验与配额 warning。
func TestTimetableGenerateSchoolGuards(t *testing.T) {
	db := setupDB(t)
	svc := newSeededTimetableService(db)

	// 无节次
	school := seedSchool(t, db)
	weeklyOne := []services.GenerateRuleSubject{{Name: "语文", Weekly: 1}}
	noPeriods, err := svc.GenerateSchool(school.ID, services.GenerateRules{Days: []int{1}, Subjects: weeklyOne}, false)
	require.NoError(t, err)
	assert.False(t, noPeriods.Success)
	assert.Equal(t, []string{"请先在「节次作息」中设置每日节次"}, noPeriods.Warnings)
	assert.NotNil(t, noPeriods.Classes)
	assert.Empty(t, noPeriods.Classes)

	seedClassPeriods(t, db, school.ID, amPmPeriods())

	// 未选上课日
	noDays, err := svc.GenerateSchool(school.ID, services.GenerateRules{Days: []int{}, Subjects: weeklyOne}, false)
	require.NoError(t, err)
	assert.False(t, noDays.Success)
	assert.Equal(t, []string{"请至少选择一个上课日"}, noDays.Warnings)

	// 尚未登记任课
	noAssignment, err := svc.GenerateSchool(school.ID, services.GenerateRules{Days: []int{1}, Subjects: weeklyOne}, false)
	require.NoError(t, err)
	assert.False(t, noAssignment.Success)
	assert.Equal(t, []string{"尚未登记任课：请先在各班「任课设置」中填写科目 → 教师"}, noAssignment.Warnings)

	// 无科目配置
	classA := makeClass(t, db, school.ID, "一班", "一年级")
	require.NoError(t, svc.SaveAssignments(classA.ID, school.ID, []services.TimetableAssignmentRow{
		{SubjectName: "语文", TeacherName: "张老师"},
	}))
	noConfig, err := svc.GenerateSchool(school.ID, services.GenerateRules{
		Days:     []int{1},
		Subjects: []services.GenerateRuleSubject{{Name: "语文", Weekly: 0}},
	}, false)
	require.NoError(t, err)
	assert.False(t, noConfig.Success)
	assert.Equal(t, []string{"请为至少一个科目设置每周节数"}, noConfig.Warnings)

	// 任课与规则无交集
	noIntersect, err := svc.GenerateSchool(school.ID, services.GenerateRules{
		Days:     []int{1},
		Subjects: []services.GenerateRuleSubject{{Name: "地理", Weekly: 2}},
	}, false)
	require.NoError(t, err)
	assert.False(t, noIntersect.Success)
	assert.Equal(t, []string{"任课表中的科目与规则配置无交集，请核对科目名称"}, noIntersect.Warnings)

	// 某班每周节数超出每日上限 → 班级 warning + 80 次重试后全校失败
	overload, err := svc.GenerateSchool(school.ID, services.GenerateRules{
		Days:     []int{1, 2},
		Subjects: []services.GenerateRuleSubject{{Name: "语文", Weekly: 4, MaxPerDay: 1}},
	}, false)
	require.NoError(t, err)
	assert.False(t, overload.Success)
	require.Len(t, overload.Warnings, 2)
	assert.Equal(t, "一班「语文」每周 4 节超出每日上限 × 2 天", overload.Warnings[0])
	assert.Equal(t, "尝试多次均无法排出满足全部规则的全校课表（教师冲突难以避免），请减少节数、放宽连堂或调整任课", overload.Warnings[1])
	assert.Empty(t, overload.Classes)
}

// newGenerateSchoolFixture 全校排课夹具：一班（语文 + 数学）、二班（语文）、三班（无任课），
// 另含张老师第 1 天第 1 节不可用。
func newGenerateSchoolFixture(t *testing.T, db *gorm.DB, svc *services.TimetableService) (models.School, models.ClassRoom, models.ClassRoom, models.ClassRoom, services.GenerateRules) {
	t.Helper()
	school := seedSchool(t, db)
	seedClassPeriods(t, db, school.ID, []services.TimetablePeriodInput{
		{PeriodIndex: 1, StartTime: "08:00", EndTime: "08:45"},
		{PeriodIndex: 2, StartTime: "08:55", EndTime: "09:40"},
		{PeriodIndex: 3, StartTime: "13:30", EndTime: "14:15"},
	})

	classA := makeClass(t, db, school.ID, "一班", "一年级")
	classB := makeClass(t, db, school.ID, "二班", "一年级")
	classC := makeClass(t, db, school.ID, "三班", "一年级")

	require.NoError(t, svc.SaveAssignments(classA.ID, school.ID, []services.TimetableAssignmentRow{
		{SubjectName: "语文", TeacherName: "张老师"},
		{SubjectName: "数学", TeacherName: "李老师"},
	}))
	require.NoError(t, svc.SaveAssignments(classB.ID, school.ID, []services.TimetableAssignmentRow{
		{SubjectName: "语文", TeacherName: "王老师"},
	}))
	require.NoError(t, svc.SaveUnavailabilities(school.ID, "张老师", []services.UnavailabilityCell{{Weekday: 1, PeriodIndex: 1}}))

	rules := services.GenerateRules{
		Days: []int{1, 2, 3},
		Subjects: []services.GenerateRuleSubject{
			{Name: "语文", Weekly: 2, MaxPerDay: 1},
			{Name: "数学", Weekly: 2, MaxPerDay: 1, Session: "am"},
		},
	}
	return school, classA, classB, classC, rules
}

// makePendingChange 建一条待审课表修改申请。
func makePendingChange(t *testing.T, db *gorm.DB, schoolID, classID uint) models.TimetableChangeRequest {
	t.Helper()
	request := models.TimetableChangeRequest{
		SchoolID: schoolID, ClassID: classID, RequestedBy: 1,
		Payload: "{}", EntryCount: 1, Status: models.TimetableChangePending,
	}
	require.NoError(t, db.Create(&request).Error)
	return request
}

// ④ commit=false 只验证可行性：不写排课 / 科目，不作废待审申请。
func TestTimetableGenerateSchoolPreviewDoesNotCommit(t *testing.T) {
	db := setupDB(t)
	svc := newSeededTimetableService(db)
	school, classA, classB, classC, rules := newGenerateSchoolFixture(t, db, svc)

	legacySubject := models.Subject{SchoolID: school.ID, Name: "旧科目", Color: "#123456"}
	require.NoError(t, db.Create(&legacySubject).Error)
	require.NoError(t, db.Create(&models.TimetableEntry{
		ClassID: classA.ID, Weekday: 5, PeriodIndex: 1, WeekType: "all",
		SubjectID: &legacySubject.ID, TeacherName: "旧老师",
	}).Error)
	pending := makePendingChange(t, db, school.ID, classA.ID)

	result, err := svc.GenerateSchool(school.ID, rules, false)
	require.NoError(t, err)
	require.True(t, result.Success, "排课应可行：%v", result.Warnings)
	assert.Empty(t, result.Warnings)
	require.Len(t, result.Classes, 2, "三班未登记任课，不进摘要")
	assert.Equal(t, classA.ID, result.Classes[0].ClassID)
	assert.Equal(t, "一班", result.Classes[0].ClassName)
	assert.Equal(t, 4, result.Classes[0].EntryCount)
	assert.Equal(t, classB.ID, result.Classes[1].ClassID)
	assert.Equal(t, "二班", result.Classes[1].ClassName)
	assert.Equal(t, 2, result.Classes[1].EntryCount)

	// 排课未落库
	entriesA := loadEntries(t, db, classA.ID)
	require.Len(t, entriesA, 1)
	assert.Equal(t, 5, entriesA[0].Weekday)
	assert.Equal(t, 1, entriesA[0].PeriodIndex)
	assert.Empty(t, loadEntries(t, db, classB.ID))
	assert.Empty(t, loadEntries(t, db, classC.ID))

	var subjectCount int64
	require.NoError(t, db.Model(&models.Subject{}).Count(&subjectCount).Error)
	assert.Equal(t, int64(1), subjectCount, "预览不新增科目")

	var reloaded models.TimetableChangeRequest
	require.NoError(t, db.First(&reloaded, pending.ID).Error)
	assert.Equal(t, models.TimetableChangePending, reloaded.Status)
	assert.Nil(t, reloaded.ReviewedAt)
}

// ⑤ commit=true：事务内落库全部班级 + 作废相关班级待审申请；教师名与不可用时段约束生效。
func TestTimetableGenerateSchoolCommitPersists(t *testing.T) {
	db := setupDB(t)
	svc := newSeededTimetableService(db)
	school, classA, classB, classC, rules := newGenerateSchoolFixture(t, db, svc)

	legacySubject := models.Subject{SchoolID: school.ID, Name: "旧科目", Color: "#123456"}
	require.NoError(t, db.Create(&legacySubject).Error)
	require.NoError(t, db.Create(&models.TimetableEntry{
		ClassID: classA.ID, Weekday: 4, PeriodIndex: 3, WeekType: "all", SubjectID: &legacySubject.ID,
	}).Error)
	require.NoError(t, db.Create(&models.TimetableEntry{
		ClassID: classC.ID, Weekday: 5, PeriodIndex: 2, WeekType: "all",
		SubjectID: &legacySubject.ID, TeacherName: "旧老师",
	}).Error)

	pendingA := makePendingChange(t, db, school.ID, classA.ID)
	pendingB := makePendingChange(t, db, school.ID, classB.ID)
	pendingC := makePendingChange(t, db, school.ID, classC.ID)

	result, err := svc.GenerateSchool(school.ID, rules, true)
	require.NoError(t, err)
	require.True(t, result.Success, "排课应可行：%v", result.Warnings)
	require.Len(t, result.Classes, 2)

	// 每班落库条数与摘要一致、同班同格不重复、周次恒 all
	for _, summary := range result.Classes {
		entries := loadEntries(t, db, summary.ClassID)
		require.Len(t, entries, summary.EntryCount, "摘要条数应与落库一致")
		seen := map[string]bool{}
		for _, entry := range entries {
			key := strconv.Itoa(entry.Weekday) + "|" + strconv.Itoa(entry.PeriodIndex)
			assert.False(t, seen[key], "同班同格重复：%s", key)
			seen[key] = true
			assert.Equal(t, "all", entry.WeekType)
			require.NotNil(t, entry.SubjectID, "科目必须已 upsert 并解析出 id")
			assert.Empty(t, entry.Room)
		}
	}

	// 一班：语文 2 节（张老师，避开不可用时段）/ 数学 2 节（李老师，限上午）
	dataA, err := svc.Bootstrap(classA.ID, school.ID)
	require.NoError(t, err)
	require.Len(t, dataA.Entries, 4, "一班的旧排课被整班覆盖")
	perSubject := map[string]int{}
	for _, entry := range dataA.Entries {
		require.NotNil(t, entry.SubjectName)
		require.NotNil(t, entry.TeacherName)
		perSubject[*entry.SubjectName]++
		if *entry.SubjectName == "语文" {
			assert.Equal(t, "张老师", *entry.TeacherName)
			assert.False(t, entry.Weekday == 1 && entry.PeriodIndex == 1, "被标记不可用的格子永不被占用")
		} else {
			assert.Equal(t, "数学", *entry.SubjectName)
			assert.Equal(t, "李老师", *entry.TeacherName)
			assert.Contains(t, []int{1, 2}, entry.PeriodIndex, "数学限上午")
		}
	}
	assert.Equal(t, map[string]int{"语文": 2, "数学": 2}, perSubject)

	// 二班：语文 2 节 → 王老师
	dataB, err := svc.Bootstrap(classB.ID, school.ID)
	require.NoError(t, err)
	require.Len(t, dataB.Entries, 2)
	for _, entry := range dataB.Entries {
		require.NotNil(t, entry.SubjectName)
		assert.Equal(t, "语文", *entry.SubjectName)
		require.NotNil(t, entry.TeacherName)
		assert.Equal(t, "王老师", *entry.TeacherName)
	}

	// 未参与排课的班级不受影响
	entriesC := loadEntries(t, db, classC.ID)
	require.Len(t, entriesC, 1)
	assert.Equal(t, 5, entriesC[0].Weekday)

	// 一班 / 二班待审申请作废；三班（未参与排课）保持 pending
	for _, id := range []uint{pendingA.ID, pendingB.ID} {
		var reloaded models.TimetableChangeRequest
		require.NoError(t, db.First(&reloaded, id).Error)
		assert.Equal(t, models.TimetableChangeRejected, reloaded.Status)
		require.NotNil(t, reloaded.ReviewNote)
		assert.Equal(t, "管理员全校自动排课，本申请自动作废", *reloaded.ReviewNote)
		assert.NotNil(t, reloaded.ReviewedAt)
	}
	var stillPending models.TimetableChangeRequest
	require.NoError(t, db.First(&stillPending, pendingC.ID).Error)
	assert.Equal(t, models.TimetableChangePending, stillPending.Status)
	assert.Nil(t, stillPending.ReviewedAt)
}

// ⑥ 同一教师带两个班：可行时同格不得冲突（多轮），不可行时返回失败且不落库。
func TestTimetableGenerateSchoolSharedTeacher(t *testing.T) {
	db := setupDB(t)
	svc := newSeededTimetableService(db)
	school := seedSchool(t, db)
	seedClassPeriods(t, db, school.ID, []services.TimetablePeriodInput{
		{PeriodIndex: 1, StartTime: "08:00", EndTime: "08:45"},
		{PeriodIndex: 2, StartTime: "08:55", EndTime: "09:40"},
	})
	classA := makeClass(t, db, school.ID, "一班", "一年级")
	classB := makeClass(t, db, school.ID, "二班", "一年级")
	for _, class := range []models.ClassRoom{classA, classB} {
		require.NoError(t, svc.SaveAssignments(class.ID, school.ID, []services.TimetableAssignmentRow{
			{SubjectName: "语文", TeacherName: "张老师"},
		}))
	}

	// 可行：2 天 × 2 节共 4 格，正好够同一教师的 4 节 —— 两个班同格必冲突
	feasible := services.GenerateRules{
		Days:     []int{1, 2},
		Subjects: []services.GenerateRuleSubject{{Name: "语文", Weekly: 2, MaxPerDay: 1}},
	}
	for round := 0; round < 5; round++ {
		result, err := svc.GenerateSchool(school.ID, feasible, true)
		require.NoError(t, err)
		require.True(t, result.Success, "第 %d 轮应成功：%v", round, result.Warnings)

		teacherSlots := map[string]bool{}
		total := 0
		for _, class := range []models.ClassRoom{classA, classB} {
			entries := loadEntries(t, db, class.ID)
			require.Len(t, entries, 2)
			for _, entry := range entries {
				assert.Equal(t, "张老师", entry.TeacherName)
				key := strconv.Itoa(entry.Weekday) + "|" + strconv.Itoa(entry.PeriodIndex)
				assert.False(t, teacherSlots[key], "同一教师在两个班同格：%s", key)
				teacherSlots[key] = true
				total++
			}
		}
		assert.Equal(t, 4, total)
	}

	// 不可行：每天只有 1 节，同一教师无法同时带两个班 → 失败且不落库
	require.NoError(t, db.Where("school_id = ?", school.ID).Delete(&models.ClassPeriod{}).Error)
	seedClassPeriods(t, db, school.ID, []services.TimetablePeriodInput{
		{PeriodIndex: 1, StartTime: "08:00", EndTime: "08:45"},
	})
	require.NoError(t, db.Where("class_id IN ?", []uint{classA.ID, classB.ID}).Delete(&models.TimetableEntry{}).Error)
	var subjectCountBefore int64
	require.NoError(t, db.Model(&models.Subject{}).Count(&subjectCountBefore).Error)

	infeasible := services.GenerateRules{
		Days:     []int{1, 2},
		Subjects: []services.GenerateRuleSubject{{Name: "语文", Weekly: 2, MaxPerDay: 1}},
	}
	result, err := svc.GenerateSchool(school.ID, infeasible, false)
	require.NoError(t, err)
	assert.False(t, result.Success, "教师冲突不可避免时应返回失败而不是产生冲突")
	assert.Equal(t, []string{"尝试多次均无法排出满足全部规则的全校课表（教师冲突难以避免），请减少节数、放宽连堂或调整任课"}, result.Warnings)
	assert.NotNil(t, result.Classes)
	assert.Empty(t, result.Classes)

	committed, err := svc.GenerateSchool(school.ID, infeasible, true)
	require.NoError(t, err)
	assert.False(t, committed.Success)
	for _, class := range []models.ClassRoom{classA, classB} {
		assert.Empty(t, loadEntries(t, db, class.ID), "失败时不得落库")
	}
	var subjectCountAfter int64
	require.NoError(t, db.Model(&models.Subject{}).Count(&subjectCountAfter).Error)
	assert.Equal(t, subjectCountBefore, subjectCountAfter, "失败时不新增科目")
}
