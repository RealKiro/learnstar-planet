// 课表网格导出（CSV 版）：
//
//	GET /api/v1/admin/classes/:id/timetable/export-excel  （单班）
//	GET /api/v1/admin/timetable/export-excel              （全校）
//
// 网格布局逐字照抄 Laravel App\Services\TimetableService::excelGrid（第 1326-1379 行）：
//
//	第 1 行 = "<班级名> 课表"；第 2 行 = 空行；第 3 行 = 表头「节次 / 时间 / 星期X…」；
//	之后每个节次一行：「第N节」「HH:MM-HH:MM」+ 各星期单元格（多条排课以换行分隔、
//	科目行后接「教师 · 教室」元信息行）；单元格文本与 Laravel 完全一致。
//
// 上课日列数取「排课中出现的最大星期」与 5 的较大者（同上）。
//
// 有意差异：
//  1. Laravel 用 maatwebsite/excel 输出 .xlsx（单班一个工作表；全校每班一个工作表并做同名去重），
//     Go 端零新依赖改用 encoding/csv：单班文件名 <班级名>-课表.csv；全校文件名 全校课表.csv，
//     多班按 Laravel 的 sheet 顺序（school_id 内 grade ASC, name ASC）顺序拼接到同一个 CSV，
//     班级之间插入一个空行作分隔（CSV 无工作表概念，故不做 sheet 名去重）。
//  2. xlsx 的样式（加粗/自动换行/列宽）在 CSV 无对应概念，直接丢弃。
//  3. 空行（len(row) == 0）输出为仅含换行的空记录，占位与 xlsx 的空行一致。
package services

import (
	"errors"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// 课表导出文件名（Laravel：<班级名>-课表.xlsx / 全校课表.xlsx；CSV 版改后缀）。
const (
	timetableClassFileSuffix = "-课表"
	timetableSchoolFilename  = "全校课表"
	csvExt                   = ".csv"
)

// TimetableGrid 某班课表的导出网格。
type TimetableGrid struct {
	ClassName string
	Rows      [][]string
}

// TimetableExportFile 课表导出文件（CSV）。
type TimetableExportFile struct {
	Filename string
	Content  []byte
}

// ExcelGrid 某班课表的导出网格（同 Laravel TimetableService::excelGrid）。
// 注意：班级不存在时班级名取「班级」（同 `ClassRoom::...->value('name') ?: '班级'`）。
func (s *TimetableService) ExcelGrid(classID, schoolID uint) (*TimetableGrid, error) {
	className := "班级"
	var class models.ClassRoom
	err := s.db.Select("id", "name").First(&class, classID).Error
	if err == nil {
		if class.Name != "" {
			className = class.Name
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 科目 id → 名称（同 Subject::where('school_id')->pluck('name','id')）。
	var subjects []models.Subject
	if err := s.db.Where("school_id = ?", schoolID).Find(&subjects).Error; err != nil {
		return nil, err
	}
	nameByID := map[uint]string{}
	for _, subject := range subjects {
		nameByID[subject.ID] = subject.Name
	}

	// 单元格：weekday => period_index => 文本列表（单双周多条并存）。
	cells := map[int]map[int][]string{}
	maxDay := 5

	var entries []models.TimetableEntry
	if err := s.db.Where("class_id = ?", classID).
		Order("weekday ASC, period_index ASC").Find(&entries).Error; err != nil {
		return nil, err
	}
	for i := range entries {
		entry := &entries[i]
		if entry.SubjectID == nil {
			continue
		}
		subject := nameByID[*entry.SubjectID]
		if subject == "" {
			continue
		}
		if entry.Weekday > maxDay {
			maxDay = entry.Weekday
		}

		week := entry.WeekType
		if week == "" {
			week = "all"
		}
		label := ""
		switch week {
		case "all":
			label = ""
		case "odd":
			label = "（单周）"
		default:
			label = "（双周）"
		}

		lines := []string{subject + label}
		if meta := joinNonEmpty(" · ", filterFalsy([]string{entry.TeacherName, entry.Room})); meta != "" {
			lines = append(lines, meta)
		}

		if cells[entry.Weekday] == nil {
			cells[entry.Weekday] = map[int][]string{}
		}
		cells[entry.Weekday][entry.PeriodIndex] = append(cells[entry.Weekday][entry.PeriodIndex], joinLines(lines))
	}

	var periods []models.ClassPeriod
	if err := s.db.Where("school_id = ?", schoolID).
		Order("period_index ASC").Find(&periods).Error; err != nil {
		return nil, err
	}

	rows := [][]string{{className + " 课表"}}
	rows = append(rows, []string{})

	header := []string{"节次", "时间"}
	for day := 1; day <= maxDay; day++ {
		header = append(header, weekdayLabels[day])
	}
	rows = append(rows, header)

	for _, period := range periods {
		row := []string{
			"第" + itoa(period.PeriodIndex) + "节",
			period.StartTime + "-" + period.EndTime,
		}
		for day := 1; day <= maxDay; day++ {
			row = append(row, joinLines(cells[day][period.PeriodIndex]))
		}
		rows = append(rows, row)
	}

	return &TimetableGrid{ClassName: className, Rows: rows}, nil
}

// ExportClassFile 单班课表导出。
//
// 同 Laravel adminExportClassExcel：班级必须属于该管理员所在学校，否则 404「班级不存在」；
// 无任何排课 / 无节次时仍导出（只有标题、表头行）。
func (s *TimetableService) ExportClassFile(classID, schoolID uint) (*TimetableExportFile, error) {
	var class models.ClassRoom
	err := s.db.Where("id = ? AND school_id = ?", classID, schoolID).First(&class).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("班级不存在")
	}
	if err != nil {
		return nil, err
	}

	grid, err := s.ExcelGrid(classID, schoolID)
	if err != nil {
		return nil, err
	}

	return &TimetableExportFile{
		Filename: grid.ClassName + timetableClassFileSuffix + csvExt,
		Content:  csvText(grid.Rows),
	}, nil
}

// ExportSchoolFile 全校课表导出；本校无班级时返回 (nil, nil)（控制器 404「暂无班级可导出」）。
func (s *TimetableService) ExportSchoolFile(schoolID uint) (*TimetableExportFile, error) {
	var classes []models.ClassRoom
	if err := s.db.Select("id", "name").
		Where("school_id = ?", schoolID).
		Order("grade ASC, name ASC").Find(&classes).Error; err != nil {
		return nil, err
	}
	if len(classes) == 0 {
		return nil, nil
	}

	rows := [][]string{}
	for i := range classes {
		grid, err := s.ExcelGrid(classes[i].ID, schoolID)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			// 班级之间插入一个空行作分隔（对应 Laravel 的多工作表边界）。
			rows = append(rows, []string{})
		}
		rows = append(rows, grid.Rows...)
	}

	return &TimetableExportFile{
		Filename: timetableSchoolFilename + csvExt,
		Content:  csvText(rows),
	}, nil
}

// filterFalsy 等价 PHP array_filter（无回调）：剔除 ""、"0" 等假值。
func filterFalsy(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" && value != "0" {
			out = append(out, value)
		}
	}
	return out
}

// joinNonEmpty 拼接非空元素（PHP implode 的自然语义）。
func joinNonEmpty(sep string, values []string) string {
	out := ""
	for i, value := range values {
		if i > 0 {
			out += sep
		}
		out += value
	}
	return out
}

// joinLines 以换行拼接单元格内多行文本（PHP implode("\n", ...)）。
func joinLines(lines []string) string {
	out := ""
	for i, line := range lines {
		if i > 0 {
			out += "\n"
		}
		out += line
	}
	return out
}
