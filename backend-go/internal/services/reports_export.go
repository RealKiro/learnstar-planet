// 教师端报表导出（CSV 版）：GET /api/v1/teacher/reports/export/{type}。
//
// 忠实移植自 Laravel App\Http\Controllers\Api\TeacherController::exportReport
// + App\Services\ReportService::export（第 125-153 行）
// + App\Exports\{ScoresExport,PetsExport,AttendanceExport}：
//   - 作用域校验：class_id 缺省取管辖的第一个班级，(int) 转换后必须落在管辖范围内，否则 403「无权限」；
//   - 文件名主体 <班级名>-<Ymd-His>，班级不存在时班级名取「未知班级」；
//   - 列名（表头）与每一列的数据映射逐字照抄三个 Export 类；
//   - 未知 type 返回 (nil, nil)，由控制器按 Laravel 返回「导出类型 X 不支持，可选: scores, pets, attendance」。
//
// 有意差异：
//  1. Laravel 用 maatwebsite/excel 输出 .xlsx（含工作表标题 WithTitle）；Go 端零新依赖改用
//     encoding/csv 输出 CSV：文件名后缀 xlsx → csv，开头带 UTF-8 BOM，表头即 WithHeadings，
//     xlsx 的工作表名（WithTitle）在 CSV 无对应概念、直接丢弃。
//  2. AttendanceExport 的「签到时间」列在 Laravel 读的是属性 check_in_time，而 attendances 表并无该列
//     （$record->check_in_time 恒为 null → (string) null 为空串），故此处同样输出恒空串，
//     不改为 sign_in_at，以保证与源码逐字一致。
//  3. 学生 / 考勤记录查询补了 Order("id ASC")：Laravel 未显式排序（隐式为插入顺序），结果集一致。
package services

import (
	"errors"
	"strings"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// ReportExport 报表导出文件（CSV）。
type ReportExport struct {
	// Filename 文件名主体（同 Laravel，仅把 xlsx 换成 csv）。
	Filename string
	// Content CSV 内容（含 UTF-8 BOM）。
	Content []byte
}

// 三个导出类型的文件名后缀（同 ReportService::export 的三个分支）。
const (
	reportExportScoresSuffix     = "-积分报表"
	reportExportPetsSuffix       = "-宠物报表"
	reportExportAttendanceSuffix = "-考勤报表"
)

// ExportReport 生成报表导出文件。
//
// exportType ∈ scores|pets|attendance；其余类型返回 (nil, nil) 由控制器拼「不支持」文案。
// classIDRaw/classIDProvided 对应 `$request->input('class_id', ...)`：未提供时取管辖的第一个班级，
// 提供了则按 PHP `(int)` 语义转换（非数字 → 0 → 越权 403）。
func (s *ReportService) ExportReport(
	u *models.User,
	exportType string,
	classIDRaw string,
	classIDProvided bool,
	date string,
) (*ReportExport, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	classID := uint(0)
	if classIDProvided {
		classID = uint(phpCastInt(classIDRaw))
	} else if len(classIDs) > 0 {
		classID = classIDs[0]
	}
	if !containsClassID(classIDs, classID) {
		return nil, ErrForbidden("无权限")
	}

	className := "未知班级"
	var class models.ClassRoom
	switch err := s.db.First(&class, classID).Error; {
	case err == nil:
		className = class.Name
	case errors.Is(err, gorm.ErrRecordNotFound):
		// ClassRoom::find 为空 → optional($class)->name ?? '未知班级'
	default:
		return nil, err
	}

	base := className + "-" + util.Now().Format("20060102-150405")

	switch exportType {
	case "scores":
		rows, err := s.scoreExportRows(classID)
		if err != nil {
			return nil, err
		}
		return &ReportExport{Filename: base + reportExportScoresSuffix + ".csv", Content: csvText(rows)}, nil
	case "pets":
		rows, err := s.petExportRows(classID)
		if err != nil {
			return nil, err
		}
		return &ReportExport{Filename: base + reportExportPetsSuffix + ".csv", Content: csvText(rows)}, nil
	case "attendance":
		rows, err := s.attendanceExportRows(classID, date)
		if err != nil {
			return nil, err
		}
		return &ReportExport{Filename: base + reportExportAttendanceSuffix + ".csv", Content: csvText(rows)}, nil
	default:
		return nil, nil
	}
}

// scoreExportRows 积分报表（ScoresExport：姓名/学号/总积分/获得积分/扣除积分/宠物名/宠物等级）。
func (s *ReportService) scoreExportRows(classID uint) ([][]string, error) {
	students, petByStudent, err := s.exportStudents(classID)
	if err != nil {
		return nil, err
	}

	positive, negative, err := s.studentScoreSums(studentIDsOf(students))
	if err != nil {
		return nil, err
	}

	rows := [][]string{{"姓名", "学号", "总积分", "获得积分", "扣除积分", "宠物名", "宠物等级"}}
	for i := range students {
		pet := petByStudent[students[i].ID]

		petName := ""
		petLevel := 0
		if pet != nil {
			petName = pet.Name
			petLevel = pet.Level
		}

		negativeSum := negative[students[i].ID]
		if negativeSum < 0 {
			negativeSum = -negativeSum
		}

		rows = append(rows, []string{
			students[i].Name,
			students[i].StudentNo,
			itoa(students[i].TotalScore),
			itoa(positive[students[i].ID]),
			itoa(negativeSum),
			petName,
			itoa(petLevel),
		})
	}
	return rows, nil
}

// petExportRows 宠物报表（PetsExport：学生姓名/宠物名/宠物系列/等级/进化阶段/经验值）。
func (s *ReportService) petExportRows(classID uint) ([][]string, error) {
	students, petByStudent, err := s.exportStudents(classID)
	if err != nil {
		return nil, err
	}

	rows := [][]string{{"学生姓名", "宠物名", "宠物系列", "等级", "进化阶段", "经验值"}}
	for i := range students {
		petName := "无"
		species := ""
		level := 0
		stage := "未孵化"
		experience := 0
		if pet := petByStudent[students[i].ID]; pet != nil {
			petName = pet.Name
			species = pet.Species
			level = pet.Level
			stage = pet.CurrentStage().Name
			experience = pet.Experience
		}

		rows = append(rows, []string{
			students[i].Name,
			petName,
			species,
			itoa(level),
			stage,
			itoa(experience),
		})
	}
	return rows, nil
}

// attendanceExportRows 考勤报表（AttendanceExport：姓名/状态/签到时间/来源）。
func (s *ReportService) attendanceExportRows(classID uint, date string) ([][]string, error) {
	// 同 Laravel：`$students->pluck('id', 'name')` —— 以姓名为键（同名只留一条，保留首次出现的位置）。
	var students []models.Student
	if err := s.db.Select("id", "name").
		Where("class_id = ? AND status = ?", classID, "active").
		Order("id ASC").Find(&students).Error; err != nil {
		return nil, err
	}

	names := []string{}
	studentIDByName := map[string]uint{}
	for i := range students {
		if _, seen := studentIDByName[students[i].Name]; !seen {
			names = append(names, students[i].Name)
		}
		studentIDByName[students[i].Name] = students[i].ID
	}

	query := s.db.Where("class_id = ?", classID)
	if date != "" {
		// 等价 Laravel `whereDate('created_at', $date)`：按「业务时区下的当天」过滤。
		// 有意差异：Laravel 走 DB 端 DATE()/strftime()；Go 端用 [当日零点, 次日零点) 区间，
		// 因为纯 Go sqlite 驱动写入的是带时区偏移的 RFC3339 串，SQLite 的 date() 会先换算成 UTC 再取日，
		// 与「按本地日比较」的原语义不符（MySQL DATETIME 无时区，不做换算）。
		day, err := time.ParseInLocation("2006-01-02", leadingDate(date), util.Loc)
		if err != nil {
			// 非法日期与任何行都不匹配（同 whereDate 与非法值比较的结果）。
			query = query.Where("1 = 0")
		} else {
			query = query.Where("created_at >= ? AND created_at < ?", day, day.AddDate(0, 0, 1))
		}
	}
	var records []models.Attendance
	if err := query.Order("id ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	// `$records->keyBy('student_id')`：同一学生多条时后写入的覆盖先前的。
	byStudent := map[uint]models.Attendance{}
	for _, record := range records {
		byStudent[record.StudentID] = record
	}

	rows := [][]string{{"姓名", "状态", "签到时间", "来源"}}
	for _, name := range names {
		status := "未记录"
		checkIn := ""
		source := ""
		if record, ok := byStudent[studentIDByName[name]]; ok {
			status = attendanceStatusLabel(record.Status)
			// check_in：Laravel 读不存在的属性 check_in_time，恒为空串（见文件头注释）。
			source = record.Source
		}
		rows = append(rows, []string{name, status, checkIn, source})
	}
	return rows, nil
}

// leadingDate 取日期串的日期部分（`2026-10-03 12:00:00` → `2026-10-03`），
// 同 Laravel whereDate 对带时间部分的入参只取日期语义。
func leadingDate(value string) string {
	if idx := strings.IndexAny(value, " T"); idx > 0 {
		return value[:idx]
	}
	return value
}

// exportStudents 取班级 active 学生（同 ScoresExport / PetsExport 的 with('pet')），
// 并返回按 student_id 索引的宠物表（hasOne，重复时取首条，与 Eloquent 一致）。
func (s *ReportService) exportStudents(classID uint) ([]models.Student, map[uint]*models.Pet, error) {
	students := []models.Student{}
	if err := s.db.Where("class_id = ? AND status = ?", classID, "active").
		Order("id ASC").Find(&students).Error; err != nil {
		return nil, nil, err
	}

	petByStudent := map[uint]*models.Pet{}
	if ids := studentIDsOf(students); len(ids) > 0 {
		pets := []models.Pet{}
		if err := s.db.Where("student_id IN ?", ids).Order("id ASC").Find(&pets).Error; err != nil {
			return nil, nil, err
		}
		for i := range pets {
			if _, seen := petByStudent[pets[i].StudentID]; !seen {
				petByStudent[pets[i].StudentID] = &pets[i]
			}
		}
	}
	return students, petByStudent, nil
}

// studentScoreSums 汇总每个学生的正/负积分（同 ScoresExport 两条 SUM 查询：> 0 与 < 0）。
func (s *ReportService) studentScoreSums(studentIDs []uint) (map[uint]int, map[uint]int, error) {
	positive := map[uint]int{}
	negative := map[uint]int{}
	if len(studentIDs) == 0 {
		return positive, negative, nil
	}

	var rows []models.Score
	if err := s.db.Select("student_id", "amount").
		Where("student_id IN ?", studentIDs).Find(&rows).Error; err != nil {
		return nil, nil, err
	}
	for i := range rows {
		switch amount := rows[i].Amount; {
		case amount > 0:
			positive[rows[i].StudentID] += amount
		case amount < 0:
			negative[rows[i].StudentID] += amount
		}
	}
	return positive, negative, nil
}

// attendanceStatusLabel 考勤状态中文名（同 AttendanceExport::statusLabel，未知状态原样返回）。
func attendanceStatusLabel(status string) string {
	switch status {
	case "present":
		return "出勤"
	case "late":
		return "迟到"
	case "leave":
		return "请假"
	case "absent":
		return "缺席"
	default:
		return status
	}
}

// containsClassID 判断班级 ID 是否在管辖范围内。
func containsClassID(classIDs []uint, classID uint) bool {
	for _, id := range classIDs {
		if id == classID {
			return true
		}
	}
	return false
}

// phpCastInt 等价 PHP `(int) $string`：忽略首尾空白，取前导可选符号与连续数字，其余为 0。
// （PHP 的 (int) "12abc" === 12，strconv.Atoi 会失败，故单独实现。）
func phpCastInt(value string) int {
	s := strings.TrimSpace(value)
	i := 0
	negative := false
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		negative = s[i] == '-'
		i++
	}
	start := i
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if start == i {
		return 0
	}

	n := 0
	for _, digit := range s[start:i] {
		n = n*10 + int(digit-'0')
		if n < 0 { // 溢出兜底
			return 0
		}
	}
	if negative {
		return -n
	}
	return n
}
