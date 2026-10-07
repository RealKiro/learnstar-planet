// 课表进阶服务：CSV 批量导入、教师不可用时段、冲突检查、任课设置。
// 忠实移植自 Laravel App\Services\TimetableService 的
// importFromCsv / listUnavailabilities / saveUnavailabilities / checkConflicts /
// listAssignments / saveAssignments（routes/api.php 第 98-99、145、150-152 行）。
//
// 有意差异（逐条）：
//  1. CSV 仅接受 UTF-8：Laravel 用 mb_check_encoding + mb_convert_encoding 自动转 GBK，
//     Go 标准库无 GBK 解码（与既有的教师导入实现口径一致）→ 非 UTF-8 直接 422。
//  2. 任课设置 / 不可用时段对「同键重复行」去重（后者覆盖前者）：Laravel 直插会撞唯一索引报 500。
//  3. 导入时按班级依次调用 Save（每次自带事务），不做整批大事务（避免依赖 savepoint）。
//  4. 节次插入顺序按 period_index 升序（Laravel 按首次出现顺序；落库结果等价）。
//  5. 行号口径：按物理行号计数（表头为第 1 行），与 Laravel 一致；未做空行折叠。
package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// 课表 CSV 模板列别名（中文优先，同 Laravel importFromCsv 的 $aliases）。
var timetableCsvAliases = map[string][]string{
	"grade":     {"年级", "grade"},
	"class":     {"班级", "班级名称", "class", "class_name"},
	"weekday":   {"星期", "weekday"},
	"period":    {"第几节", "节次", "period", "period_index"},
	"start":     {"开始时间", "start", "start_time"},
	"end":       {"结束时间", "end", "end_time"},
	"subject":   {"科目", "subject"},
	"week_type": {"周次", "week_type"},
	"teacher":   {"教师", "老师", "teacher"},
	"room":      {"教室", "room"},
}

// timetableCsvRequiredCols CSV 必需列（缺一即报错）。
var timetableCsvRequiredCols = []string{"grade", "class", "weekday", "period", "subject"}

// CsvImportClassRow 导入汇总里的班级行。
type CsvImportClassRow struct {
	Grade      string `json:"grade"`
	Name       string `json:"name"`
	Found      bool   `json:"found"`
	EntryCount int    `json:"entry_count"`
}

// CsvImportSummary CSV 导入结果（字段与 Laravel importFromCsv 返回数组逐字一致）。
type CsvImportSummary struct {
	TotalRows           int                 `json:"total_rows"`
	Classes             []CsvImportClassRow `json:"classes"`
	PeriodCount         int                 `json:"period_count"`
	Errors              []string            `json:"errors"`
	DryRun              bool                `json:"dry_run"`
	Imported            bool                `json:"imported"`
	SubjectsAutoColored *int                `json:"subjects_auto_colored,omitempty"`
}

// csvImportClass CSV 中出现的班级及其待导入格子。
type csvImportClass struct {
	grade    string
	name     string
	rowLabel string
	classID  uint
	entries  []TimetableEntryInput
	subjects []string
}

// ImportFromCsv 从 CSV 文本批量导入课表。
//
// 模板列：年级,班级,星期,第几节,开始时间,结束时间,科目,周次,教师,教室
//   - 按「年级 + 班级」定位 ClassRoom，找不到的行报错跳过（不自动建班）；
//   - 科目 / 节次为学校级 upsert（同节次时间冲突时取首次出现并告警）；
//   - 某班出现在 CSV 中即整体覆盖该班课表（先清后插），不影响未出现的班级；
//   - dryRun=true 只做解析校验不落库。
func (s *TimetableService) ImportFromCsv(content []byte, schoolID uint, dryRun bool) (*CsvImportSummary, error) {
	content = bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF})
	if !utf8.Valid(content) {
		return nil, ErrUnprocessable("文件编码不是 UTF-8，请另存为 UTF-8 后重试")
	}

	header, dataLines := parseTimetableCsv(string(content))
	col := map[string]int{}
	for key, names := range timetableCsvAliases {
		for i, h := range header {
			if containsString(names, h) {
				col[key] = i
				break
			}
		}
	}

	errors := []string{}
	ready := true
	for _, key := range timetableCsvRequiredCols {
		if _, ok := col[key]; !ok {
			ready = false
		}
	}
	if !ready {
		errors = append(errors, "CSV 缺少必需列（年级 / 班级 / 星期 / 第几节 / 科目），请使用模板")
	}

	classList := []*csvImportClass{}
	byKey := map[string]*csvImportClass{}
	periods := map[int]TimetablePeriodInput{}
	totalRows := 0

	if ready {
		get := func(cells []string, key string) string {
			idx, ok := col[key]
			if !ok || idx >= len(cells) {
				return ""
			}
			return strings.TrimSpace(cells[idx])
		}

		for i, line := range dataLines {
			cells := parseCSVLine(line)
			rowNo := i + 2 // +1 表头，+1 从 1 计数

			grade := get(cells, "grade")
			className := get(cells, "class")
			if grade == "" && className == "" {
				continue // 空行
			}
			totalRows++

			subjectName := get(cells, "subject")
			rawWeekday := get(cells, "weekday")
			rawPeriod := get(cells, "period")
			weekday := parseWeekday(rawWeekday)
			periodIndex := atoiOrZero(rawPeriod)

			if subjectName == "" || weekday == nil || periodIndex < 1 {
				errors = append(errors, fmt.Sprintf(
					"第%d行：星期 / 第几节 / 科目无效（星期=%s 节次=%s 科目=%s），已跳过",
					rowNo, orEmptyPlaceholder(rawWeekday), orEmptyPlaceholder(rawPeriod), orEmptyPlaceholder(subjectName)))
				continue
			}

			key := grade + "|" + className
			item, ok := byKey[key]
			if !ok {
				item = &csvImportClass{
					grade:    grade,
					name:     className,
					rowLabel: fmt.Sprintf("%s/%s", orDash(grade), orDash(className)),
				}
				byKey[key] = item
				classList = append(classList, item)
			}

			teacher := get(cells, "teacher")
			room := get(cells, "room")
			item.subjects = append(item.subjects, subjectName)
			item.entries = append(item.entries, TimetableEntryInput{
				Weekday:     *weekday,
				PeriodIndex: periodIndex,
				WeekType:    parseWeekType(get(cells, "week_type")),
				SubjectName: subjectName,
				TeacherName: optionalString(teacher),
				Room:        optionalString(room),
			})

			start := normalizeTime(get(cells, "start"))
			end := normalizeTime(get(cells, "end"))
			if start != "" && end != "" {
				if prev, exists := periods[periodIndex]; exists && (prev.StartTime != start || prev.EndTime != end) {
					errors = append(errors, fmt.Sprintf(
						"第%d行：第%d节时间（%s–%s）与之前出现的不一致，以首次为准", rowNo, periodIndex, start, end))
				} else {
					periods[periodIndex] = TimetablePeriodInput{
						PeriodIndex: periodIndex,
						Name:        "第" + itoa(periodIndex) + "节",
						StartTime:   start,
						EndTime:     end,
					}
				}
			} else if _, exists := periods[periodIndex]; !exists {
				errors = append(errors, fmt.Sprintf(
					"第%d行：第%d节缺少开始 / 结束时间，导入后该节将无作息时间", rowNo, periodIndex))
			}
		}

		// 解析班级 id（按「年级 + 班级」精确匹配，找不到不自动建班）
		grades := uniqueStrings(classList, func(c *csvImportClass) string { return c.grade })
		names := uniqueStrings(classList, func(c *csvImportClass) string { return c.name })
		var matched []models.ClassRoom
		if len(grades) > 0 && len(names) > 0 {
			if err := s.db.Where("school_id = ? AND grade IN ? AND name IN ?", schoolID, grades, names).
				Find(&matched).Error; err != nil {
				return nil, err
			}
		}
		for _, item := range classList {
			for i := range matched {
				if matched[i].Grade == item.grade && matched[i].Name == item.name {
					item.classID = matched[i].ID
					break
				}
			}
			if item.classID == 0 {
				errors = append(errors, fmt.Sprintf(
					"班级「%s」不存在，该班 %d 行未导入（请先在班级列表创建）", item.rowLabel, len(item.entries)))
			}
		}
	}

	classRows := make([]CsvImportClassRow, 0, len(classList))
	for _, item := range classList {
		classRows = append(classRows, CsvImportClassRow{
			Grade:      item.grade,
			Name:       item.name,
			Found:      item.classID != 0,
			EntryCount: len(item.entries),
		})
	}

	summary := &CsvImportSummary{
		TotalRows:   totalRows,
		Classes:     classRows,
		PeriodCount: len(periods),
		Errors:      errors,
		DryRun:      dryRun,
		Imported:    false,
	}

	// 预览模式只返回解析结果；正式导入时匹配到的班级照常导入（错误行与未匹配班级已在 errors 中列明）。
	if dryRun {
		return summary, nil
	}

	periodPayload := make([]TimetablePeriodInput, 0, len(periods))
	for _, idx := range sortedPeriodIndexes(periods) {
		periodPayload = append(periodPayload, periods[idx])
	}

	autoColored := 0
	importedClassIDs := []uint{}
	for _, item := range classList {
		if item.classID == 0 {
			continue
		}
		colored, err := s.Save(item.classID, schoolID, TimetablePayload{
			Subjects: subjectInputs(item.subjects),
			Periods:  periodPayload,
			Entries:  item.entries,
		})
		if err != nil {
			return nil, err
		}
		autoColored += colored
		importedClassIDs = append(importedClassIDs, item.classID)
	}

	// 受影响班级的待审申请自动作废（导入即覆盖）。
	if len(importedClassIDs) > 0 {
		if err := s.db.Model(&models.TimetableChangeRequest{}).
			Where("class_id IN ? AND status = ?", importedClassIDs, models.TimetableChangePending).
			Updates(map[string]any{
				"status":      models.TimetableChangeRejected,
				"review_note": "管理员批量导入课表，本申请自动作废",
				"reviewed_at": util.Now(),
			}).Error; err != nil {
			return nil, err
		}
	}

	summary.Imported = true
	total := autoColored
	summary.SubjectsAutoColored = &total

	return summary, nil
}

// ============================================================
// 教师不可用时段（replace 语义）
// ============================================================

// TeacherUnavailabilityRow 不可用时段行。
type TeacherUnavailabilityRow struct {
	TeacherName string `json:"teacher_name"`
	Weekday     int    `json:"weekday"`
	PeriodIndex int    `json:"period_index"`
}

// UnavailabilityCell 保存不可用时段时的格子。
type UnavailabilityCell struct {
	Weekday     int `json:"weekday"`
	PeriodIndex int `json:"period_index"`
}

// ListUnavailabilities 全校教师不可用时段（平铺列表，按教师 / 星期 / 节次排序）。
func (s *TimetableService) ListUnavailabilities(schoolID uint) ([]TeacherUnavailabilityRow, error) {
	var rows []models.TimetableTeacherUnavailability
	if err := s.db.Where("school_id = ?", schoolID).
		Order("teacher_name ASC, weekday ASC, period_index ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	views := make([]TeacherUnavailabilityRow, 0, len(rows))
	for _, r := range rows {
		views = append(views, TeacherUnavailabilityRow{
			TeacherName: r.TeacherName,
			Weekday:     r.Weekday,
			PeriodIndex: r.PeriodIndex,
		})
	}
	return views, nil
}

// SaveUnavailabilities 整体保存某教师的不可用时段（先清后插）。
func (s *TimetableService) SaveUnavailabilities(schoolID uint, teacherName string, cells []UnavailabilityCell) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("school_id = ? AND teacher_name = ?", schoolID, teacherName).
			Delete(&models.TimetableTeacherUnavailability{}).Error; err != nil {
			return err
		}

		seen := map[int]bool{}
		for _, c := range cells {
			if c.Weekday < 1 || c.Weekday > 7 || c.PeriodIndex < 1 {
				continue
			}
			key := c.Weekday*1000 + c.PeriodIndex
			if seen[key] {
				continue // 同键重复行去重（否则会撞唯一索引）
			}
			seen[key] = true
			if err := tx.Create(&models.TimetableTeacherUnavailability{
				SchoolID:    schoolID,
				TeacherName: teacherName,
				Weekday:     c.Weekday,
				PeriodIndex: c.PeriodIndex,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ============================================================
// 冲突检查
// ============================================================

// TimetableConflictEntry 冲突检查输入的一格（某班编辑中的排课）。
type TimetableConflictEntry struct {
	Weekday     int    `json:"weekday"`
	PeriodIndex int    `json:"period_index"`
	SubjectName string `json:"subject_name"`
	WeekType    string `json:"week_type"`
	TeacherName string `json:"teacher_name"`
}

// TimetableConflictRow 冲突明细（type = teacher / unavailable）。
type TimetableConflictRow struct {
	Type             string  `json:"type"`
	TeacherName      string  `json:"teacher_name"`
	Weekday          int     `json:"weekday"`
	PeriodIndex      int     `json:"period_index"`
	WeekType         string  `json:"week_type"`
	SubjectName      string  `json:"subject_name"`
	OtherClassID     *uint   `json:"other_class_id"`
	OtherClassName   *string `json:"other_class_name"`
	OtherSubjectName *string `json:"other_subject_name"`
	Message          string  `json:"message"`
}

// CheckConflicts 冲突检查：给定某班当前编辑中的排课，检查
// ① 教师冲突——同一教师在其他班级同一时段（周次重叠）已有课；
// ② 不可用冲突——教师被标记了不可用时段。
// 返回冲突明细列表（空数组 = 无冲突）。
func (s *TimetableService) CheckConflicts(schoolID, classID uint, entries []TimetableConflictEntry) ([]TimetableConflictRow, error) {
	conflicts := []TimetableConflictRow{}

	teachers := make([]string, 0, len(entries))
	seenTeacher := map[string]bool{}
	for _, e := range entries {
		t := strings.TrimSpace(e.TeacherName)
		if t == "" || seenTeacher[t] {
			continue
		}
		seenTeacher[t] = true
		teachers = append(teachers, t)
	}
	if len(teachers) == 0 {
		return conflicts, nil
	}

	// 其他班级相关教师的全部排课（班级限定在本校）。
	var others []models.TimetableEntry
	classSub := s.db.Model(&models.ClassRoom{}).Select("id").Where("school_id = ?", schoolID)
	if err := s.db.Where("class_id <> ? AND teacher_name IN ?", classID, teachers).
		Where("class_id IN (?)", classSub).
		Find(&others).Error; err != nil {
		return nil, err
	}

	var classList []models.ClassRoom
	if err := s.db.Where("school_id = ?", schoolID).Find(&classList).Error; err != nil {
		return nil, err
	}
	classNameByID := map[uint]string{}
	for _, c := range classList {
		classNameByID[c.ID] = c.Name
	}

	var subjects []models.Subject
	if err := s.db.Where("school_id = ?", schoolID).Find(&subjects).Error; err != nil {
		return nil, err
	}
	subjectNameByID := map[uint]string{}
	for _, sub := range subjects {
		subjectNameByID[sub.ID] = sub.Name
	}

	var unavailRows []models.TimetableTeacherUnavailability
	if err := s.db.Where("school_id = ? AND teacher_name IN ?", schoolID, teachers).
		Find(&unavailRows).Error; err != nil {
		return nil, err
	}
	unavail := map[string]bool{}
	for _, u := range unavailRows {
		unavail[fmt.Sprintf("%s|%d|%d", u.TeacherName, u.Weekday, u.PeriodIndex)] = true
	}

	seen := map[string]bool{}
	for _, e := range entries {
		teacher := strings.TrimSpace(e.TeacherName)
		weekType := normalizeWeekType(e.WeekType)
		if teacher == "" || e.Weekday < 1 || e.Weekday > 7 || e.PeriodIndex < 1 {
			continue
		}

		// ① 与其他班的教师冲突（周次重叠才算：all 与任何都重叠；odd/even 仅同类重叠）
		for i := range others {
			o := &others[i]
			otherWeekType := normalizeWeekType(o.WeekType)
			if o.TeacherName != teacher ||
				o.Weekday != e.Weekday ||
				o.PeriodIndex != e.PeriodIndex ||
				!weekTypeOverlaps(weekType, otherWeekType) {
				continue
			}

			dedupeKey := fmt.Sprintf("%s|%d|%d|%s|%d", teacher, e.Weekday, e.PeriodIndex, weekType, o.ClassID)
			if seen[dedupeKey] {
				continue
			}
			seen[dedupeKey] = true

			otherClassName, ok := classNameByID[o.ClassID]
			if !ok {
				otherClassName = "其他班级"
			}
			otherSubjectName := ""
			if o.SubjectID != nil {
				otherSubjectName = subjectNameByID[*o.SubjectID]
			}
			otherClassID := o.ClassID

			subjectLabel := e.SubjectName
			if subjectLabel == "" {
				subjectLabel = "课程"
			}
			conflicts = append(conflicts, TimetableConflictRow{
				Type:             "teacher",
				TeacherName:      teacher,
				Weekday:          e.Weekday,
				PeriodIndex:      e.PeriodIndex,
				WeekType:         weekType,
				SubjectName:      e.SubjectName,
				OtherClassID:     &otherClassID,
				OtherClassName:   strPtr(otherClassName),
				OtherSubjectName: optionalString(otherSubjectName),
				Message: fmt.Sprintf("%s%s第%d节「%s」与 %s「%s」冲突（同一教师）",
					weekdayLabels[e.Weekday],
					weekTypeLabel(weekType),
					e.PeriodIndex,
					subjectLabel,
					otherClassName,
					otherSubjectName),
			})
		}

		// ② 教师不可用时段
		if unavail[fmt.Sprintf("%s|%d|%d", teacher, e.Weekday, e.PeriodIndex)] {
			dedupeKey := fmt.Sprintf("u|%s|%d|%d|%s", teacher, e.Weekday, e.PeriodIndex, weekType)
			if !seen[dedupeKey] {
				seen[dedupeKey] = true
				subjectLabel := e.SubjectName
				if subjectLabel == "" {
					subjectLabel = "课程"
				}
				conflicts = append(conflicts, TimetableConflictRow{
					Type:             "unavailable",
					TeacherName:      teacher,
					Weekday:          e.Weekday,
					PeriodIndex:      e.PeriodIndex,
					WeekType:         weekType,
					SubjectName:      e.SubjectName,
					OtherClassID:     nil,
					OtherClassName:   nil,
					OtherSubjectName: nil,
					Message: fmt.Sprintf("%s 第%d节「%s」：教师 %s 此时段已被标记为不可用",
						weekdayLabels[e.Weekday], e.PeriodIndex, subjectLabel, teacher),
				})
			}
		}
	}

	return conflicts, nil
}

// weekTypeOverlaps 周次是否重叠：all 与任何都重叠；odd/even 仅同类重叠。
func weekTypeOverlaps(a, b string) bool {
	if a == "all" || b == "all" {
		return true
	}
	return a == b
}

// weekTypeLabel 冲突文案里的单双周前缀（all 无前缀，同 Laravel）。
func weekTypeLabel(weekType string) string {
	switch weekType {
	case "all":
		return ""
	case "odd":
		return "单周"
	default:
		return "双周"
	}
}

// normalizeWeekType 归一化周次取值（空值视为 all）。
func normalizeWeekType(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if timetableWeekTypes[trimmed] {
		return trimmed
	}
	return "all"
}

// ============================================================
// 任课设置（班级 × 科目 → 教师，replace 语义）
// ============================================================

// TimetableAssignmentRow 任课行。
type TimetableAssignmentRow struct {
	SubjectName string `json:"subject_name"`
	TeacherName string `json:"teacher_name"`
}

// ClassInSchool 判断班级是否存在且属于该学校（管理员越权防护，等价 Laravel resolveAdminClass）。
func (s *TimetableService) ClassInSchool(schoolID, classID uint) (bool, error) {
	var count int64
	if err := s.db.Model(&models.ClassRoom{}).
		Where("id = ? AND school_id = ?", classID, schoolID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListAssignments 某班任课列表（按科目名升序）。
func (s *TimetableService) ListAssignments(classID uint) ([]TimetableAssignmentRow, error) {
	var rows []models.TimetableTeacherAssignment
	if err := s.db.Where("class_id = ?", classID).
		Order("subject_name ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	views := make([]TimetableAssignmentRow, 0, len(rows))
	for _, r := range rows {
		views = append(views, TimetableAssignmentRow{
			SubjectName: r.SubjectName,
			TeacherName: r.TeacherName,
		})
	}
	return views, nil
}

// SaveAssignments 整体保存某班任课（先清后插；空科目 / 空教师跳过）。
func (s *TimetableService) SaveAssignments(classID, schoolID uint, rows []TimetableAssignmentRow) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("class_id = ?", classID).Delete(&models.TimetableTeacherAssignment{}).Error; err != nil {
			return err
		}

		seen := map[string]bool{}
		for _, row := range rows {
			subject := strings.TrimSpace(row.SubjectName)
			teacher := strings.TrimSpace(row.TeacherName)
			if subject == "" || teacher == "" {
				continue
			}
			if seen[subject] {
				continue // 同科目重复行去重（否则会撞唯一索引）
			}
			seen[subject] = true
			if err := tx.Create(&models.TimetableTeacherAssignment{
				SchoolID:    schoolID,
				ClassID:     classID,
				SubjectName: subject,
				TeacherName: teacher,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ============================================================
// CSV 解析辅助
// ============================================================

// parseTimetableCsv 解析 CSV：去 BOM 由调用方完成，首行为表头，返回 [表头, 数据行]。
func parseTimetableCsv(content string) ([]string, []string) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, nil
	}
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return nil, nil
	}

	header := parseCSVLine(lines[0])
	for i := range header {
		header[i] = strings.TrimSpace(header[i])
	}
	return header, lines[1:]
}

// parseCSVLine 按 CSV 规则解析单行（等价 PHP str_getcsv）。
func parseCSVLine(line string) []string {
	reader := csv.NewReader(strings.NewReader(line))
	reader.FieldsPerRecord = -1
	fields, err := reader.Read()
	if err != nil || fields == nil {
		return []string{}
	}
	return fields
}

// parseWeekday 星期解析：支持 1-7 / 周一~周日 / 星期一~星期日（含「天」）。
func parseWeekday(raw string) *int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	if isAllDigits(raw) {
		n := atoiOrZero(raw)
		if n >= 1 && n <= 7 {
			return &n
		}
		return nil
	}
	weekdayChars := []struct {
		ch rune
		n  int
	}{
		{'一', 1}, {'二', 2}, {'三', 3}, {'四', 4},
		{'五', 5}, {'六', 6}, {'日', 7}, {'天', 7},
	}
	for _, item := range weekdayChars {
		if strings.ContainsRune(raw, item.ch) {
			n := item.n
			return &n
		}
	}
	return nil
}

// timetableTimeRe1 / timetableTimeRe2 时间解析（等价 Laravel normalizeTime 的两条正则）。
var (
	timetableTimeRe1 = regexp.MustCompile(`^(\d{1,2}):(\d{2})$`)
	timetableTimeRe2 = regexp.MustCompile(`^(\d{2})(\d{2})$`)
)

// parseWeekType 周次解析：每周 / all（默认）、单周 / 单 / odd、双周 / 双 / even。
func parseWeekType(raw string) string {
	if strings.Contains(raw, "双") || strings.Contains(raw, "even") {
		return "even"
	}
	if strings.Contains(raw, "单") || strings.Contains(raw, "odd") {
		return "odd"
	}
	return "all"
}

// normalizeTime 时间归一化：8:00 / 08:00 / 0800 → HH:MM，非法返回空串。
func normalizeTime(raw string) string {
	raw = strings.TrimSpace(raw)
	if m := timetableTimeRe1.FindStringSubmatch(raw); m != nil {
		return fmt.Sprintf("%02d:%s", atoiOrZero(m[1]), m[2])
	}
	if m := timetableTimeRe2.FindStringSubmatch(raw); m != nil {
		return fmt.Sprintf("%02d:%s", atoiOrZero(m[1]), m[2])
	}
	return ""
}

// subjectInputs 把科目名列表转为保存载荷（保持首次出现的顺序做 sort_order）。
func subjectInputs(names []string) []TimetableSubjectInput {
	seen := map[string]bool{}
	result := []TimetableSubjectInput{}
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		result = append(result, TimetableSubjectInput{Name: name})
	}
	return result
}

// sortedPeriodIndexes 节次升序（保证落库顺序稳定）。
func sortedPeriodIndexes(periods map[int]TimetablePeriodInput) []int {
	indexes := make([]int, 0, len(periods))
	for idx := range periods {
		indexes = append(indexes, idx)
	}
	sort.Ints(indexes)
	return indexes
}

// uniqueStrings 按首次出现顺序去重提取字段值。
func uniqueStrings(items []*csvImportClass, pick func(*csvImportClass) string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, item := range items {
		value := pick(item)
		if seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

// orDash 空值显示为 '-'（同 Laravel sprintf('%s/%s', $grade ?: '-', $name ?: '-')）。
func orDash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

// orEmptyPlaceholder 空值显示为 '空'（CSV 错误行文案）。
func orEmptyPlaceholder(value string) string {
	if strings.TrimSpace(value) == "" {
		return "空"
	}
	return value
}

// optionalString 空串返回 nil，否则返回指针。
func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
