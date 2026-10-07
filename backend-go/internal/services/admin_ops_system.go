// 管理端运维服务（四）：学年升级、系统诊断 / 状态 / 日志 / 修复。
//
// 忠实移植自 Laravel App\Http\Controllers\Api\SchoolAdminController：
// previewGradeUpgrade / executeGradeUpgrade / systemDiagnose / systemStatus / systemLogs /
// systemRepair。与 Laravel 的有意差异逐条写在方法注释里。
package services

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/database"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// ============================================================
// 学年升级（不可逆）
// ============================================================

// gradeUpgradeMap 年级升级映射（同 Laravel preview/execute 内的 $gradeMap）。
var gradeUpgradeMap = map[string]string{
	"一年级": "二年级",
	"二年级": "三年级",
	"三年级": "四年级",
	"四年级": "五年级",
	"五年级": "六年级",
}

// 预览/执行的提示文案（逐字对齐 Laravel）。
const (
	gradeUpgradePreviewNote = "六年级学生将标记为毕业，二至五年级学生随班级升级到下一年级。一年级新生需在升级后手动创建班级并导入。"
	gradeUpgradeExecuteNote = "六年级学生已标记为毕业。请创建新一年级班级并导入新生名单。" +
		"注意：企业微信/钉钉/飞书侧的部门与名单不会自动同步——请先在第三方平台完成升班调整，" +
		"再回来导入通讯录并核对班级映射；导入时同学号跨班会被拦截并提示，避免转班学生重复创建。"
)

// GradeUpgradeClass 可升级班级（同 Laravel upgrade_classes 元素）。
type GradeUpgradeClass struct {
	ClassID      uint   `json:"class_id"`
	ClassName    string `json:"class_name"`
	NewName      string `json:"new_name"`
	OldGrade     string `json:"old_grade"`
	NewGrade     string `json:"new_grade"`
	StudentCount int64  `json:"student_count"`
}

// GradeGraduateClass 毕业班级（同 Laravel graduate_classes 元素）。
type GradeGraduateClass struct {
	ClassID      uint   `json:"class_id"`
	ClassName    string `json:"class_name"`
	StudentCount int64  `json:"student_count"`
}

// GradeUpgradeSummary 预览汇总（同 Laravel summary 结构）。
type GradeUpgradeSummary struct {
	UpgradeClassCount    int    `json:"upgrade_class_count"`
	GraduateClassCount   int    `json:"graduate_class_count"`
	UpgradeStudentCount  int64  `json:"upgrade_student_count"`
	GraduateStudentCount int64  `json:"graduate_student_count"`
	Note                 string `json:"note"`
}

// GradeUpgradePreview 升级预览（dry-run）。
type GradeUpgradePreview struct {
	UpgradeClasses  []GradeUpgradeClass  `json:"upgrade_classes"`
	GraduateClasses []GradeGraduateClass `json:"graduate_classes"`
	Summary         GradeUpgradeSummary  `json:"summary"`
}

// GradeUpgradeResultData 执行结果（同 Laravel executeGradeUpgrade 的 data 结构）。
type GradeUpgradeResultData struct {
	UpgradedClasses   int    `json:"upgraded_classes"`
	ArchivedClasses   int    `json:"archived_classes"`
	GraduatedStudents int64  `json:"graduated_students"`
	Note              string `json:"note"`
}

// PreviewGradeUpgrade 预览学年升级（dry-run，不写库）。
// 只取 status=active 的班级：六年级 → 毕业清单；一~五年级 → 升级清单（新班级名为年级前缀替换结果）。
func (o *AdminOps) PreviewGradeUpgrade(schoolID uint) (*GradeUpgradePreview, error) {
	var classes []models.ClassRoom
	if err := o.db.Where("school_id = ? AND status = ?", schoolID, "active").
		Order("id ASC").Find(&classes).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(classes))
	for _, c := range classes {
		ids = append(ids, c.ID)
	}
	counts, err := o.opsStudentCounts(ids)
	if err != nil {
		return nil, err
	}

	upgrade := []GradeUpgradeClass{}
	graduate := []GradeGraduateClass{}
	var upgradeStudents, graduateStudents int64

	for _, c := range classes {
		if c.Grade == "六年级" {
			studentCount := counts[c.ID]
			graduate = append(graduate, GradeGraduateClass{
				ClassID:      c.ID,
				ClassName:    c.Name,
				StudentCount: studentCount,
			})
			graduateStudents += studentCount
			continue
		}
		if newGrade, ok := gradeUpgradeMap[c.Grade]; ok {
			studentCount := counts[c.ID]
			upgrade = append(upgrade, GradeUpgradeClass{
				ClassID:      c.ID,
				ClassName:    c.Name,
				NewName:      strings.ReplaceAll(c.Name, c.Grade, newGrade),
				OldGrade:     c.Grade,
				NewGrade:     newGrade,
				StudentCount: studentCount,
			})
			upgradeStudents += studentCount
		}
	}

	return &GradeUpgradePreview{
		UpgradeClasses:  upgrade,
		GraduateClasses: graduate,
		Summary: GradeUpgradeSummary{
			UpgradeClassCount:    len(upgrade),
			GraduateClassCount:   len(graduate),
			UpgradeStudentCount:  upgradeStudents,
			GraduateStudentCount: graduateStudents,
			Note:                 gradeUpgradePreviewNote,
		},
	}, nil
}

// ExecuteGradeUpgrade 执行学年升级（事务内，**不可逆**，同 Laravel executeGradeUpgrade）。
//
// 步骤（在同一事务内）：
//  1. 六年级班级：active 学生 → status='graduated'；班级 status='archived'；
//  2. 一~五年级班级：grade 改为下一年级，班级名做年级前缀替换。
//
// ⚠️ 与 Laravel 一致：**不做幂等保护**，重复调用会继续把已升级的班级再升一级
// （归档后的六年级班级因不再 active 而不会被二次处理）。
func (o *AdminOps) ExecuteGradeUpgrade(schoolID uint) (*GradeUpgradeResultData, error) {
	var classes []models.ClassRoom
	if err := o.db.Where("school_id = ? AND status = ?", schoolID, "active").
		Order("id ASC").Find(&classes).Error; err != nil {
		return nil, err
	}

	upgraded := 0
	archived := 0
	var graduatedStudents int64

	err := o.db.Transaction(func(tx *gorm.DB) error {
		// 1. 六年级：学生标记毕业 + 班级归档
		for _, class := range classes {
			if class.Grade != "六年级" {
				continue
			}
			res := tx.Model(&models.Student{}).
				Where("class_id = ? AND status = ?", class.ID, "active").
				Update("status", "graduated")
			if res.Error != nil {
				return res.Error
			}
			graduatedStudents += res.RowsAffected
			if err := tx.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
				Update("status", "archived").Error; err != nil {
				return err
			}
			archived++
		}

		// 2. 其他年级：升级年级并重命名（基于步骤 1 之前的同一批对象，同 Laravel）
		for _, class := range classes {
			newGrade, ok := gradeUpgradeMap[class.Grade]
			if !ok {
				continue
			}
			newName := strings.ReplaceAll(class.Name, class.Grade, newGrade)
			if err := tx.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
				Updates(map[string]any{"grade": newGrade, "name": newName}).Error; err != nil {
				return err
			}
			upgraded++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &GradeUpgradeResultData{
		UpgradedClasses:   upgraded,
		ArchivedClasses:   archived,
		GraduatedStudents: graduatedStudents,
		Note:              gradeUpgradeExecuteNote,
	}, nil
}

// ============================================================
// 系统诊断 / 状态 / 日志 / 修复
// ============================================================

// DiagnoseItem 一条诊断项（同 Laravel systemDiagnose 的 data 元素）。
type DiagnoseItem struct {
	Item   string `json:"item"`
	Status string `json:"status"` // ok | missing | error | fixable | skipped
	Detail string `json:"detail,omitempty"`
}

// DiagnoseResult 诊断结果（handler 会展开为 data / has_issues / message）。
type DiagnoseResult struct {
	Items     []DiagnoseItem
	HasIssues bool
	Message   string
}

// Diagnose 检查数据库表结构关键字段是否齐全。
//
// 与 Laravel systemDiagnose 的差异：
//   - `third_party_bindings 表`：Go 端未移植第三方登录绑定（无此表），标记为 **skipped** 而非
//     missing —— 否则「系统状态」会永远报异常，而 repair 也无法修复它；
//   - `教师账号密码状态`：与 Laravel 同口径 —— users 现已有 `plain_password` 列（创建 / 重置 /
//     批量重置 / 改密 / 机器人播种都会写入明文），缺明文的账号计入 `fixable` 并给出重置提示；
//   - 表结构查询走 GORM Migrator；异常状态（error）保留在词表中，但本方法只在驱动层查询
//     失败时才可能产生（当前 sqlite/mysql/postgres 三种驱动下不会）。
func (o *AdminOps) Diagnose() (*DiagnoseResult, error) {
	items := []DiagnoseItem{}

	for _, col := range []string{"nickname", "subject", "grade_team"} {
		ok := o.hasColumn(&models.User{}, col)
		items = append(items, DiagnoseItem{
			Item:   "users 表." + col + " 字段",
			Status: diagnoseStatus(ok),
		})
	}

	for _, col := range []string{"display_code", "display_code_updated_at"} {
		ok := o.hasColumn(&models.ClassRoom{}, col)
		items = append(items, DiagnoseItem{
			Item:   "class_rooms 表." + col + " 字段",
			Status: diagnoseStatus(ok),
		})
	}

	hasTeachersTable := o.db.Migrator().HasTable("class_room_teachers")
	items = append(items, DiagnoseItem{
		Item:   "class_room_teachers 表",
		Status: diagnoseStatus(hasTeachersTable),
	})
	if hasTeachersTable {
		ok := o.db.Migrator().HasColumn("class_room_teachers", "subject")
		items = append(items, DiagnoseItem{
			Item:   "class_room_teachers.subject 字段",
			Status: diagnoseStatus(ok),
		})
	}

	if o.db.Migrator().HasTable("third_party_bindings") {
		items = append(items, DiagnoseItem{Item: "third_party_bindings 表", Status: "ok"})
	} else {
		items = append(items, DiagnoseItem{
			Item:   "third_party_bindings 表",
			Status: "skipped",
			Detail: "Go 端未移植第三方登录与账号绑定，无此表（不影响其他功能）",
		})
	}

	// 5. 教师账号明文密码状态（同 Laravel systemDiagnose：缺明文的账号计入 fixable）
	if o.hasColumn(&models.User{}, "plain_password") {
		emptyPwdCount := o.countTeachersWithoutPlainPassword()
		item := DiagnoseItem{Item: "教师账号密码状态", Status: "ok"}
		if emptyPwdCount > 0 {
			item.Status = "fixable"
			item.Detail = fmt.Sprintf("%d 个教师账号缺明文密码（第三方自动注册历史问题），请在教师列表逐个重置", emptyPwdCount)
		}
		items = append(items, item)
	}

	missing := 0
	fixable := 0
	for _, item := range items {
		switch item.Status {
		case "missing":
			missing++
		case "fixable":
			fixable++
		}
	}

	message := "系统状态正常"
	if missing > 0 {
		message = "检测到数据库结构缺失，可执行修复"
	} else if fixable > 0 {
		message = "检测到可修复的数据问题，可执行修复"
	}

	return &DiagnoseResult{
		Items:     items,
		HasIssues: missing > 0 || fixable > 0,
		Message:   message,
	}, nil
}

// hasColumn 判断模型对应表是否存在某列。
func (o *AdminOps) hasColumn(model any, column string) bool {
	return o.db.Migrator().HasColumn(model, column)
}

// diagnoseStatus 列/表存在 → ok，否则 missing（同 Laravel hasColumn 的判定）。
func diagnoseStatus(ok bool) string {
	if ok {
		return "ok"
	}
	return "missing"
}

// StatusVersion 版本信息。
//
// 与 Laravel systemStatus 的差异：Laravel 报 php / laravel 版本；Go 端报 go 运行时版本与
// 后端标识（backend = backend-go），并追加 db_driver / timezone 便于排查。
// app_env / app_debug / app_url 读取同名环境变量，缺省值与 Laravel 默认一致
// （production / false / http://localhost）。
type StatusVersion struct {
	Go       string `json:"go"`
	Backend  string `json:"backend"`
	AppEnv   string `json:"app_env"`
	AppDebug string `json:"app_debug"`
	AppURL   string `json:"app_url"`
	DBDriver string `json:"db_driver"`
	Timezone string `json:"timezone"`
}

// StatusData 系统状态（同上 data 结构）。
type StatusData struct {
	Version        StatusVersion `json:"version"`
	Migrations     []any         `json:"migrations"`
	MigrationCount int           `json:"migration_count"`
	Tables         []string      `json:"tables"`
	Note           string        `json:"note"`
}

// Status 返回版本信息与库表清单。
//
// 与 Laravel 的差异：Go 用 AutoMigrate 建表，**没有 migrations 记录表**，故以表清单
// （tables）替代迁移记录，`migrations` 恒为空数组、`migration_count` 恒为 0。
func (o *AdminOps) Status() (*StatusData, error) {
	tables, err := o.db.Migrator().GetTables()
	if err != nil {
		return nil, err
	}

	return &StatusData{
		Version: StatusVersion{
			Go:       runtime.Version(),
			Backend:  "backend-go",
			AppEnv:   envOrDefault("APP_ENV", "production"),
			AppDebug: envOrDefault("APP_DEBUG", "false"),
			AppURL:   envOrDefault("APP_URL", "http://localhost"),
			DBDriver: o.db.Dialector.Name(),
			Timezone: util.Loc.String(),
		},
		Migrations:     []any{},
		MigrationCount: 0,
		Tables:         tables,
		Note:           "Go 端以 GORM AutoMigrate 建表，没有 Laravel 的 migrations 记录表：migrations 恒为空，实际表清单见 tables。",
	}, nil
}

// envOrDefault 读取环境变量，缺省时返回默认值。
func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// LogEntry 一行日志（同 Laravel systemLogs 的 entries 元素）。
type LogEntry struct {
	Line  string `json:"line"`
	Level string `json:"level"`
}

// SystemLogsData 日志查询结果（同 Laravel systemLogs 的 data 结构，含说明性 message）。
type SystemLogsData struct {
	Entries    []LogEntry `json:"entries"`
	Total      int        `json:"total,omitempty"`
	Path       string     `json:"path"`
	Exists     bool       `json:"exists"`
	Size       int64      `json:"size,omitempty"`
	ModifiedAt string     `json:"modified_at,omitempty"`
	Message    string     `json:"message,omitempty"`
}

// logLevelPattern 同 Laravel：`\.(ERROR|WARNING|INFO|DEBUG|CRITICAL|ALERT|EMERGENCY|NOTICE):`
var logLevelPattern = regexp.MustCompile(`\.(ERROR|WARNING|INFO|DEBUG|CRITICAL|ALERT|EMERGENCY|NOTICE):`)

// SystemLogs 读取日志文件末 N 行并按等级过滤（时间逆序，同 Laravel systemLogs）。
//
// ⚠️ 与 Laravel 的有意差异：Laravel 固定读 `storage/logs/laravel.log`；Go 端日志走 stdout，
// **没有日志文件**。因此这里读取环境变量 `LOG_FILE` 指定的文件：文件存在时返回其末 N 行；
// 未配置或文件不存在时返回空列表 + 一条说明性 message（**不伪造日志**）。
// path 由 handler 从环境变量解析后传入，便于测试注入临时文件。
func (o *AdminOps) SystemLogs(path string, lines int, levelFilter string) *SystemLogsData {
	if lines <= 0 {
		lines = 200
	}
	data := &SystemLogsData{Entries: []LogEntry{}, Path: path}

	if path == "" {
		data.Message = "Go 端日志输出到 stdout，没有日志文件。如需在后台查看日志，请设置环境变量 LOG_FILE 指向日志文件路径（如 /app/logs/app.log）后重启。"
		return data
	}
	info, err := os.Stat(path)
	if err != nil {
		data.Message = "日志文件不存在：" + path + "。Go 端默认把日志打到 stdout，可用 LOG_FILE 环境变量指定日志文件。"
		return data
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		data.Message = "日志文件读取失败：" + err.Error()
		return data
	}

	all := strings.Split(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n")
	if len(all) > lines {
		all = all[len(all)-lines:]
	}

	filter := strings.ToUpper(strings.TrimSpace(levelFilter))
	entries := make([]LogEntry, 0, len(all))
	for _, line := range all {
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			continue
		}
		level := "INFO"
		if m := logLevelPattern.FindStringSubmatch(line); m != nil {
			level = m[1]
		}
		if filter != "" && level != filter {
			continue
		}
		entries = append(entries, LogEntry{Line: line, Level: level})
	}
	// 时间逆序
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	data.Entries = entries
	data.Total = len(entries)
	data.Exists = true
	data.Size = info.Size()
	data.ModifiedAt = info.ModTime().Format("2006-01-02 15:04:05")
	return data
}

// RepairResult 修复结果（同 Laravel systemRepair 的响应结构）。
type RepairResult struct {
	Message               string `json:"message"`
	Output                string `json:"output"`
	PendingPasswordResets int64  `json:"pending_password_resets"`
}

// Repair 执行系统修复：跑一遍 GORM AutoMigrate（幂等）。
//
// 与 Laravel systemRepair 的差异：Laravel 调 `artisan migrate --force`；Go 端以 AutoMigrate
// 代替迁移（幂等，无 schema 变更则不产生操作）。`pending_password_resets` 与提示文案
// **同 Laravel 口径**：统计 role=teacher 且明文密码为空的账号数。
func (o *AdminOps) Repair() (*RepairResult, error) {
	if err := database.Migrate(o.db); err != nil {
		return nil, err
	}

	emptyPwdCount := o.countTeachersWithoutPlainPassword()
	message := "数据库迁移已完成"
	if emptyPwdCount > 0 {
		message += fmt.Sprintf("；检测到 %d 个教师账号缺少明文密码记录（第三方自动注册历史问题），请在「教师管理 → 密码」中逐个重置为默认密码", emptyPwdCount)
	}

	return &RepairResult{
		Message:               message,
		Output:                "已执行 GORM AutoMigrate（幂等：无 schema 变更时不产生操作）",
		PendingPasswordResets: emptyPwdCount,
	}, nil
}

// countTeachersWithoutPlainPassword 统计 role=teacher 且 plain_password 为 NULL/空串的账号数
// （同 Laravel systemDiagnose/systemRepair 的统计口径）。
func (o *AdminOps) countTeachersWithoutPlainPassword() int64 {
	var count int64
	// 条件必须显式加括号：GORM 不会为原始字符串条件补括号，
	// 缺括号时 `AND` 优先级高于 `OR`，会把非教师账号也算进来。
	if err := o.db.Model(&models.User{}).Where("role = ?", "teacher").
		Where("(plain_password IS NULL OR plain_password = ?)", "").
		Count(&count).Error; err != nil {
		return 0
	}
	return count
}
