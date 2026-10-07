// 管理端运维服务（一）：服务骨架、全校报表、大屏登录日志。
//
// 忠实移植自 Laravel App\Http\Controllers\Api\SchoolAdminController 的对应动作：
// schoolOverview / reportsByGrade / reportsByClass / displayLoginLogs。
// 与 Laravel 的有意差异逐条写在方法注释里。
package services

import (
	"math"
	"os"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// 上传目录（映射 Laravel 的 storage/app/public 与公开访问前缀 /storage）。
const (
	// DefaultUploadRoot 默认上传落盘根目录（相对启动目录）。
	DefaultUploadRoot = "storage/app/uploads"
	// DefaultUploadURLPrefix 默认上传访问前缀（写入 schools.logo_path 的前缀）。
	DefaultUploadURLPrefix = "/storage/app/uploads"
	// SchoolLogoSubdir 学校 LOGO 子目录。
	SchoolLogoSubdir = "schools"
)

// AdminOps 管理端运维服务：报表 / 批量账号 / CSV 导入 / 学年升级 / 系统诊断 / 学校 LOGO。
type AdminOps struct {
	db *gorm.DB

	// UploadRoot 上传文件落盘根目录（默认 storage/app/uploads）。
	// 测试可指向 t.TempDir()，避免污染工作区。
	UploadRoot string
	// UploadURLPrefix 对外可见的访问前缀（默认 /storage/app/uploads）。
	UploadURLPrefix string
}

// NewAdminOps 创建管理端运维服务。
// NewAdminOps 创建管理端运维服务。
// 上传落盘目录可用环境变量 UPLOAD_DIR 覆盖（默认 storage/app/uploads）；
// 对外访问前缀固定为 /storage/app/uploads，映射 Laravel 的 storage/app/public + /storage。
func NewAdminOps(db *gorm.DB) *AdminOps {
	ops := &AdminOps{db: db, UploadRoot: DefaultUploadRoot, UploadURLPrefix: DefaultUploadURLPrefix}
	if dir := os.Getenv("UPLOAD_DIR"); dir != "" {
		ops.UploadRoot = dir
	}
	return ops
}

// ============================================================
// 公共小工具
// ============================================================

// opsClassIDs 返回本校全部班级 ID（含已归档班级）。
// 口径同 Laravel `ClassRoom::where('school_id', $school->id)->pluck('id')`。
func (o *AdminOps) opsClassIDs(schoolID uint) ([]uint, error) {
	var ids []uint
	if err := o.db.Model(&models.ClassRoom{}).Where("school_id = ?", schoolID).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// opsClassNames 返回本校「班级 ID → 班级名」。
func (o *AdminOps) opsClassNames(schoolID uint) (map[uint]string, error) {
	var classes []models.ClassRoom
	if err := o.db.Where("school_id = ?", schoolID).Find(&classes).Error; err != nil {
		return nil, err
	}
	names := make(map[uint]string, len(classes))
	for _, c := range classes {
		names[c.ID] = c.Name
	}
	return names, nil
}

// opsStudentCounts 返回「班级 ID → 学生数」，等价 Laravel `withCount('students')`
// （不带 status 过滤：毕业学生也计入，同 Laravel 关系计数）。
func (o *AdminOps) opsStudentCounts(classIDs []uint) (map[uint]int64, error) {
	counts := map[uint]int64{}
	if len(classIDs) == 0 {
		return counts, nil
	}
	var rows []struct {
		ClassID uint
		Total   int64
	}
	if err := o.db.Model(&models.Student{}).
		Select("class_id, COUNT(*) AS total").
		Where("class_id IN ?", classIDs).
		Group("class_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		counts[r.ClassID] = r.Total
	}
	return counts, nil
}

// opsScoreSumRange 统计 [from, to) 区间内的积分合计（无班级时返回 0）。
// 等价 Laravel `whereBetween('created_at', [$start, $end])` 的 sum('amount')。
func (o *AdminOps) opsScoreSumRange(classIDs []uint, from, to time.Time) (int, error) {
	if len(classIDs) == 0 {
		return 0, nil
	}
	var total int
	if err := o.db.Model(&models.Score{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("class_id IN ? AND created_at >= ? AND created_at < ?", classIDs, from, to).
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// opsRound1 等价 PHP round($v, 1)。
func opsRound1(v float64) float64 { return math.Round(v*10) / 10 }

// ============================================================
// 报表
// ============================================================

// OverviewData 全校概览（Laravel schoolOverview 的 data 结构）。
type OverviewData struct {
	ClassCount        int64   `json:"class_count"`
	TeacherCount      int64   `json:"teacher_count"`
	StudentCount      int64   `json:"student_count"`
	MonthlyScore      int     `json:"monthly_score"`
	LastMonthScore    int     `json:"last_month_score"`
	ScoreTrendPercent float64 `json:"score_trend_percent"`
	MonthLabel        string  `json:"month_label"`
}

// SchoolOverview 全校概览：班级/教师/学生数 + 本月与上月积分与环比。
// class_count 统计本校全部班级（含归档，Laravel 未加 status 过滤）；
// teacher_count 仅统计 role=teacher 且 status=active；
// student_count 统计本校班级下全部学生（含毕业，Laravel 未加 status 过滤）。
func (o *AdminOps) SchoolOverview(schoolID uint) (*OverviewData, error) {
	var classCount, teacherCount, studentCount int64
	if err := o.db.Model(&models.ClassRoom{}).Where("school_id = ?", schoolID).
		Count(&classCount).Error; err != nil {
		return nil, err
	}
	if err := o.db.Model(&models.User{}).
		Where("school_id = ? AND role = ? AND status = ?", schoolID, "teacher", "active").
		Count(&teacherCount).Error; err != nil {
		return nil, err
	}
	classIDs, err := o.opsClassIDs(schoolID)
	if err != nil {
		return nil, err
	}
	if len(classIDs) > 0 {
		if err := o.db.Model(&models.Student{}).Where("class_id IN ?", classIDs).
			Count(&studentCount).Error; err != nil {
			return nil, err
		}
	}

	now := util.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, util.Loc)
	nextMonthStart := monthStart.AddDate(0, 1, 0)
	lastMonthStart := monthStart.AddDate(0, -1, 0)

	monthlyScore, err := o.opsScoreSumRange(classIDs, monthStart, nextMonthStart)
	if err != nil {
		return nil, err
	}
	lastMonthScore, err := o.opsScoreSumRange(classIDs, lastMonthStart, monthStart)
	if err != nil {
		return nil, err
	}

	// Laravel：上月 > 0 时算环比；否则本月 > 0 记 100，否则 0。
	trend := 0.0
	switch {
	case lastMonthScore > 0:
		trend = opsRound1(float64(monthlyScore-lastMonthScore) / math.Max(float64(lastMonthScore), 1) * 100)
	case monthlyScore > 0:
		trend = 100
	}

	return &OverviewData{
		ClassCount:        classCount,
		TeacherCount:      teacherCount,
		StudentCount:      studentCount,
		MonthlyScore:      monthlyScore,
		LastMonthScore:    lastMonthScore,
		ScoreTrendPercent: trend,
		MonthLabel:        now.Format("2006-01"),
	}, nil
}

// GradeReportRow 按年级汇总的一行（Laravel reportsByGrade 元素结构）。
type GradeReportRow struct {
	Grade        string  `json:"grade"`
	ClassCount   int     `json:"class_count"`
	StudentCount int     `json:"student_count"`
	AvgScore     float64 `json:"avg_score"`
	TotalScore   int     `json:"total_score"`
}

// ReportsByGrade 按年级汇总：班级数 / 学生数 / 平均分 / 总分。
// 年级为空时归入「未分年级」；分组顺序按班级首次出现顺序（Laravel groupBy 保序）。
func (o *AdminOps) ReportsByGrade(schoolID uint) ([]GradeReportRow, error) {
	var classes []models.ClassRoom
	if err := o.db.Where("school_id = ?", schoolID).Order("id ASC").Find(&classes).Error; err != nil {
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

	type gradeGroup struct {
		grade    string
		classIDs []uint
	}
	var groups []gradeGroup
	index := map[string]int{}
	for _, c := range classes {
		grade := c.Grade
		if grade == "" {
			grade = "未分年级"
		}
		if i, ok := index[grade]; ok {
			groups[i].classIDs = append(groups[i].classIDs, c.ID)
			continue
		}
		index[grade] = len(groups)
		groups = append(groups, gradeGroup{grade: grade, classIDs: []uint{c.ID}})
	}

	rows := make([]GradeReportRow, 0, len(groups))
	for _, g := range groups {
		var totalScore int
		if err := o.db.Model(&models.Score{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("class_id IN ?", g.classIDs).
			Scan(&totalScore).Error; err != nil {
			return nil, err
		}
		var studentCount int64
		for _, id := range g.classIDs {
			studentCount += counts[id]
		}
		avg := 0.0
		if studentCount > 0 {
			avg = opsRound1(float64(totalScore) / float64(studentCount))
		}
		rows = append(rows, GradeReportRow{
			Grade:        g.grade,
			ClassCount:   len(g.classIDs),
			StudentCount: int(studentCount),
			AvgScore:     avg,
			TotalScore:   totalScore,
		})
	}
	return rows, nil
}

// ClassReportRow 按班级汇总的一行（Laravel reportsByClass 元素结构）。
type ClassReportRow struct {
	ClassID      uint    `json:"class_id"`
	ClassName    string  `json:"class_name"`
	Grade        string  `json:"grade"`
	TeacherName  *string `json:"teacher_name"`
	StudentCount int     `json:"student_count"`
	MonthlyScore int     `json:"monthly_score"`
}

// ReportsByClass 按班级汇总：班主任 / 学生数 / 本月积分。
//
// monthly_score 的口径同 Laravel：`whereMonth('created_at', 当前月)`——**只看月份数字，
// 不看年份**。该条件无法用跨库通用 SQL 表达（sqlite 是 strftime、MySQL 是 MONTH()），
// 故本方法一次性取出本校班级的 (class_id, amount, created_at) 三列后在 Go 侧按月份聚合。
func (o *AdminOps) ReportsByClass(schoolID uint) ([]ClassReportRow, error) {
	var classes []models.ClassRoom
	if err := o.db.Where("school_id = ?", schoolID).Order("id ASC").Find(&classes).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(classes))
	teacherIDs := make([]uint, 0, len(classes))
	for _, c := range classes {
		ids = append(ids, c.ID)
		if c.TeacherID != nil {
			teacherIDs = append(teacherIDs, *c.TeacherID)
		}
	}

	counts, err := o.opsStudentCounts(ids)
	if err != nil {
		return nil, err
	}

	teacherNames := map[uint]string{}
	if len(teacherIDs) > 0 {
		var teachers []models.User
		if err := o.db.Where("id IN ?", teacherIDs).
			Select("id, name").Find(&teachers).Error; err != nil {
			return nil, err
		}
		for _, t := range teachers {
			teacherNames[t.ID] = t.Name
		}
	}

	monthly := map[uint]int{}
	if len(ids) > 0 {
		var scores []struct {
			ClassID   uint
			Amount    int
			CreatedAt time.Time
		}
		if err := o.db.Model(&models.Score{}).
			Select("class_id, amount, created_at").
			Where("class_id IN ?", ids).
			Scan(&scores).Error; err != nil {
			return nil, err
		}
		month := util.Now().Month()
		for _, s := range scores {
			if s.CreatedAt.In(util.Loc).Month() == month {
				monthly[s.ClassID] += s.Amount
			}
		}
	}

	rows := make([]ClassReportRow, 0, len(classes))
	for _, c := range classes {
		var teacherName *string
		if c.TeacherID != nil {
			if name, ok := teacherNames[*c.TeacherID]; ok {
				teacherName = strPtr(name)
			}
		}
		rows = append(rows, ClassReportRow{
			ClassID:      c.ID,
			ClassName:    c.Name,
			Grade:        c.Grade,
			TeacherName:  teacherName,
			StudentCount: int(counts[c.ID]),
			MonthlyScore: monthly[c.ID],
		})
	}
	return rows, nil
}

// ============================================================
// 班级码登录日志（大屏）
// ============================================================

// DisplayLoginLogFilter 大屏登录日志筛选条件（同 Laravel displayLoginLogs 的 query 参数）。
type DisplayLoginLogFilter struct {
	ClassID uint   // class_id：精确匹配
	IP      string // ip：模糊匹配
	Date    string // date：按日过滤（2006-01-02）
	Page    int
}

// DisplayLoginLogItem 登录日志条目。
type DisplayLoginLogItem struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	ClassCode string `json:"class_code"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	LoginAt   string `json:"login_at"`
}

// DisplayLoginLogResult 分页结果（Laravel paginate(50) 形态：data + meta）。
type DisplayLoginLogResult struct {
	Data []DisplayLoginLogItem `json:"data"`
	Meta PageMeta              `json:"meta"`
}

// DisplayLoginLogs 本校班级的大屏（班级码）登录日志，按时间倒序分页（每页 50）。
func (o *AdminOps) DisplayLoginLogs(schoolID uint, f DisplayLoginLogFilter) (*DisplayLoginLogResult, error) {
	const perPage = 50

	classIDs, err := o.opsClassIDs(schoolID)
	if err != nil {
		return nil, err
	}
	names, err := o.opsClassNames(schoolID)
	if err != nil {
		return nil, err
	}

	base := o.db.Model(&models.DisplayLoginLog{})
	if len(classIDs) == 0 {
		base = base.Where("1 = 0")
	} else {
		base = base.Where("class_id IN ?", classIDs)
	}
	if f.ClassID > 0 {
		base = base.Where("class_id = ?", f.ClassID)
	}
	if f.IP != "" {
		base = base.Where("ip_address LIKE ?", "%"+f.IP+"%")
	}
	if f.Date != "" {
		// Laravel whereDate('created_at', 'Y-m-d')：按日比较，等价于 [当日 00:00, 次日 00:00)。
		if day, err := time.ParseInLocation("2006-01-02", f.Date, util.Loc); err == nil {
			base = base.Where("created_at >= ? AND created_at < ?", day, day.AddDate(0, 0, 1))
		}
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	lastPage := int((total + perPage - 1) / perPage)
	if lastPage < 1 {
		lastPage = 1
	}
	page := f.Page
	if page < 1 {
		page = 1
	}

	var logs []models.DisplayLoginLog
	if err := base.Order("created_at DESC, id DESC").
		Offset((page - 1) * perPage).Limit(perPage).
		Find(&logs).Error; err != nil {
		return nil, err
	}

	items := make([]DisplayLoginLogItem, 0, len(logs))
	for _, l := range logs {
		items = append(items, DisplayLoginLogItem{
			ID:        l.ID,
			ClassName: names[l.ClassID],
			ClassCode: l.ClassCode,
			IPAddress: l.IPAddress,
			UserAgent: l.UserAgent,
			LoginAt:   l.CreatedAt.In(util.Loc).Format("2006-01-02 15:04:05"),
		})
	}

	return &DisplayLoginLogResult{
		Data: items,
		Meta: PageMeta{CurrentPage: page, LastPage: lastPage, Total: int(total)},
	}, nil
}
