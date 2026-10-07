// 课表自动排课：单班规则排课（纯计算，不落库）+ 全校智能排课（可选落库）。
// 忠实移植自 Laravel App\Services\TimetableService 的
// generate / tryGenerate / placeBlock / sessionMatches / generateSchool
// （backend/app/Services/TimetableService.php 第 629-1015 行），
// 路由见 backend/routes/api.php 第 147-148 行。
//
// 有意差异（逐条）：
//  1. 随机源可注入：TimetableService 新增 seedSource 字段（SetSeedSource），
//     默认按时间种子；Laravel 直接用全局 mt_rand（shuffle）。算法本身逐句照做
//     （块序打散 → 候选先 shuffle 再按「该日该科目已有节数」稳定升序取第一个）。
//  2. 节次上午 / 下午判定：Laravel 用 `strtotime('1970-01-01 ' . start_time)`，解析
//     失败时落到 false（= 下午）；本实现无法解析 / 为空的 start_time 视为「未知」
//     （isAm 里无该键），sessionMatches 对未知时段一律放行（等价于 Laravel 里
//     `$isAm[$cell] ?? null === null` 的分支）。
//  3. entries 输出顺序：按 rules.days 的原始顺序，日内按 period_index 升序，
//     与 Laravel 的 $grid 插入顺序 + ksort 一致；rules.days 非升序时才与「全局 day 升序」不同。
//  4. 全校落库的事务范围：与 Laravel 相同（一个外层事务包住全部班级 save + 作废 pending 申请），
//     班级内层 save 借 GORM savepoint 实现。
package services

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// 重试上限：单班 60 次、全校 80 次（Laravel generate / generateSchool）。
const (
	generateMaxAttempts       = 60
	generateSchoolMaxAttempts = 80
)

// 排课警告文案（与 Laravel 逐字一致）。
const (
	generateNoPeriodsWarning   = "请先在「节次作息」中设置每日节次"
	generateNoDaysWarning      = "请至少选择一个上课日"
	generateNoWeeklyWarning    = "请为至少一个科目设置每周节数"
	generateRetryFailedWarning = "尝试多次均无法排出满足全部规则的课表，请减少节数、放宽连堂或时段限制"
	generateNoAssignmentWaring = "尚未登记任课：请先在各班「任课设置」中填写科目 → 教师"
	generateNoIntersectWarning = "任课表中的科目与规则配置无交集，请核对科目名称"
	generateSchoolRetryWarning = "尝试多次均无法排出满足全部规则的全校课表（教师冲突难以避免），请减少节数、放宽连堂或调整任课"
)

// GenerateRuleSubject 单个科目的排课规则（入参形状同 Laravel rules.subjects[]）。
//
// Weekly < 1 或 Name 为空串的条目在服务层被跳过（对应 Laravel `$weekly < 1` 分支）；
// MaxPerDay 为 0 时按 1 处理（Laravel `max(1, (int)($s['max_per_day'] ?? 1))`）；
// Session 非 any/am/pm 时按 any 处理。
type GenerateRuleSubject struct {
	Name          string `json:"name"`
	Weekly        int    `json:"weekly"`
	Double        bool   `json:"double"`
	Session       string `json:"session"`
	MaxPerDay     int    `json:"max_per_day"`
	ForbidPeriods []int  `json:"forbid_periods"`
}

// GenerateRules 排课规则。Days 为 nil 时按 Laravel `?? [1,2,3,4,5]` 取默认值；
// 非 nil 的空切片等价于「未选择任何上课日」。
type GenerateRules struct {
	Days     []int                 `json:"days"`
	Subjects []GenerateRuleSubject `json:"subjects"`
}

// GenerateEntry 生成的一条排课（week_type 恒为 all，teacher_name / room 恒为 null）。
type GenerateEntry struct {
	Weekday     int     `json:"weekday"`
	PeriodIndex int     `json:"period_index"`
	WeekType    string  `json:"week_type"`
	SubjectName string  `json:"subject_name"`
	TeacherName *string `json:"teacher_name"`
	Room        *string `json:"room"`
}

// GenerateResult 单班生成结果（与 Laravel generate 返回数组字段一致）。
type GenerateResult struct {
	Success  bool            `json:"success"`
	Warnings []string        `json:"warnings"`
	Entries  []GenerateEntry `json:"entries"`
}

// GenerateSchoolClass 全校生成结果里的班级摘要。
type GenerateSchoolClass struct {
	ClassID    uint   `json:"class_id"`
	ClassName  string `json:"class_name"`
	EntryCount int    `json:"entry_count"`
}

// GenerateSchoolResult 全校生成结果（字段是 classes，不是 entries）。
type GenerateSchoolResult struct {
	Success  bool                  `json:"success"`
	Warnings []string              `json:"warnings"`
	Classes  []GenerateSchoolClass `json:"classes"`
}

// generateBlock 一个课时块（连堂为 size=2）。
type generateBlock struct {
	subject   string
	size      int
	session   string
	maxPerDay int
	forbid    map[int]bool
	teacher   *string // 全校排课为教师姓名；单班为 nil（对应 Laravel 的 null）
}

// generateSubjectConfig 全校排课的科目配置（按科目名全局生效，重名以最后一次为准）。
type generateSubjectConfig struct {
	name      string
	weekly    int
	double    bool
	session   string
	maxPerDay int
	forbid    map[int]bool
}

// generateClassPlan 某班的待排块。
type generateClassPlan struct {
	id     uint
	name   string
	blocks []generateBlock
}

// generateClassResult 某班的排课结果。
type generateClassResult struct {
	id      uint
	name    string
	entries []GenerateEntry
}

// Generate 按规则为单个班级生成课表（纯计算，不落库）。
func (s *TimetableService) Generate(schoolID uint, rules GenerateRules) (*GenerateResult, error) {
	var periods []models.ClassPeriod
	if err := s.db.Where("school_id = ?", schoolID).Order("period_index ASC").Find(&periods).Error; err != nil {
		return nil, err
	}
	if len(periods) == 0 {
		return &GenerateResult{Success: false, Warnings: []string{generateNoPeriodsWarning}, Entries: []GenerateEntry{}}, nil
	}

	days := normalizeGenerateDays(rules.Days)
	if len(days) == 0 {
		return &GenerateResult{Success: false, Warnings: []string{generateNoDaysWarning}, Entries: []GenerateEntry{}}, nil
	}

	periodIndexes, isAm := generateGridMeta(periods)

	// 拆课时块
	warnings := []string{}
	blocks := []generateBlock{}
	for _, sub := range rules.Subjects {
		name := strings.TrimSpace(sub.Name)
		weekly := sub.Weekly
		if name == "" || weekly < 1 {
			continue
		}
		double := sub.Double
		session := normalizeGenerateSession(sub.Session)
		maxPerDay := normalizeGenerateMaxPerDay(sub.MaxPerDay, double)
		forbid := intSet(sub.ForbidPeriods)

		remaining := weekly
		if double {
			for remaining >= 2 {
				blocks = append(blocks, generateBlock{subject: name, size: 2, session: session, maxPerDay: maxPerDay, forbid: forbid})
				remaining -= 2
			}
		}
		for i := 0; i < remaining; i++ {
			blocks = append(blocks, generateBlock{subject: name, size: 1, session: session, maxPerDay: maxPerDay, forbid: forbid})
		}

		if weekly > maxPerDay*len(days) {
			warnings = append(warnings, fmt.Sprintf("科目「%s」每周 %d 节，超出每日上限 %d × %d 天，无法排下", name, weekly, maxPerDay, len(days)))
		}
	}

	capacity := len(days) * len(periodIndexes)
	needed := 0
	for _, b := range blocks {
		needed += b.size
	}
	if needed > capacity {
		warnings = append(warnings, fmt.Sprintf("所需节数 %d 超过可用格数 %d（%d 天 × %d 节），请调整规则", needed, capacity, len(days), len(periodIndexes)))
		return &GenerateResult{Success: false, Warnings: warnings, Entries: []GenerateEntry{}}, nil
	}
	if needed == 0 {
		warnings = append(warnings, generateNoWeeklyWarning)
		return &GenerateResult{Success: false, Warnings: warnings, Entries: []GenerateEntry{}}, nil
	}

	rng := s.newRand()
	for attempt := 0; attempt < generateMaxAttempts; attempt++ {
		if entries, ok := tryGenerateBlocks(rng, blocks, days, periodIndexes, isAm); ok {
			return &GenerateResult{Success: true, Warnings: warnings, Entries: entries}, nil
		}
	}

	warnings = append(warnings, generateRetryFailedWarning)
	return &GenerateResult{Success: false, Warnings: warnings, Entries: []GenerateEntry{}}, nil
}

// GenerateSchool 全校智能排课（依据任课表，教师冲突为硬约束）。
//
// commit=false 只验证可行性并返回摘要；commit=true 时在事务内落库全部班级的课表，
// 并作废相关班级的待审课表修改申请。
func (s *TimetableService) GenerateSchool(schoolID uint, rules GenerateRules, commit bool) (*GenerateSchoolResult, error) {
	var periods []models.ClassPeriod
	if err := s.db.Where("school_id = ?", schoolID).Order("period_index ASC").Find(&periods).Error; err != nil {
		return nil, err
	}
	if len(periods) == 0 {
		return &GenerateSchoolResult{Success: false, Warnings: []string{generateNoPeriodsWarning}, Classes: []GenerateSchoolClass{}}, nil
	}

	days := normalizeGenerateDays(rules.Days)
	if len(days) == 0 {
		return &GenerateSchoolResult{Success: false, Warnings: []string{generateNoDaysWarning}, Classes: []GenerateSchoolClass{}}, nil
	}

	periodIndexes, isAm := generateGridMeta(periods)

	// 任课表
	var assignments []models.TimetableTeacherAssignment
	if err := s.db.Where("school_id = ?", schoolID).Find(&assignments).Error; err != nil {
		return nil, err
	}
	if len(assignments) == 0 {
		return &GenerateSchoolResult{Success: false, Warnings: []string{generateNoAssignmentWaring}, Classes: []GenerateSchoolClass{}}, nil
	}
	teacherOf := map[uint]map[string]string{}
	for _, a := range assignments {
		if teacherOf[a.ClassID] == nil {
			teacherOf[a.ClassID] = map[string]string{}
		}
		teacherOf[a.ClassID][a.SubjectName] = a.TeacherName
	}

	// 科目配置（按名全局生效，重名以最后一次为准、位置保持首次出现）
	warnings := []string{}
	configOrder := []string{}
	configByName := map[string]generateSubjectConfig{}
	for _, sub := range rules.Subjects {
		name := strings.TrimSpace(sub.Name)
		weekly := sub.Weekly
		if name == "" || weekly < 1 {
			continue
		}
		cfg := generateSubjectConfig{
			name:      name,
			weekly:    weekly,
			double:    sub.Double,
			session:   normalizeGenerateSession(sub.Session),
			maxPerDay: normalizeGenerateMaxPerDay(sub.MaxPerDay, sub.Double),
			forbid:    intSet(sub.ForbidPeriods),
		}
		if _, seen := configByName[name]; !seen {
			configOrder = append(configOrder, name)
		}
		configByName[name] = cfg
	}
	if len(configOrder) == 0 {
		return &GenerateSchoolResult{Success: false, Warnings: []string{generateNoWeeklyWarning}, Classes: []GenerateSchoolClass{}}, nil
	}

	// 各班课时块（按 grade, name 顺序；未登记任课的班级跳过）
	var classes []models.ClassRoom
	if err := s.db.Where("school_id = ?", schoolID).Order("grade ASC, name ASC").Find(&classes).Error; err != nil {
		return nil, err
	}
	plans := []generateClassPlan{}
	for _, class := range classes {
		cid := class.ID
		if len(teacherOf[cid]) == 0 {
			continue
		}

		blocks := []generateBlock{}
		for _, name := range configOrder {
			cfg := configByName[name]
			teacher, assigned := teacherOf[cid][name]
			if !assigned {
				continue // 该班此科目未任课
			}

			if cfg.weekly > cfg.maxPerDay*len(days) {
				warnings = append(warnings, fmt.Sprintf("%s「%s」每周 %d 节超出每日上限 × %d 天", class.Name, name, cfg.weekly, len(days)))
			}

			teacherName := teacher
			remaining := cfg.weekly
			if cfg.double {
				for remaining >= 2 {
					blocks = append(blocks, generateBlock{subject: name, size: 2, session: cfg.session, maxPerDay: cfg.maxPerDay, forbid: cfg.forbid, teacher: &teacherName})
					remaining -= 2
				}
			}
			for i := 0; i < remaining; i++ {
				blocks = append(blocks, generateBlock{subject: name, size: 1, session: cfg.session, maxPerDay: cfg.maxPerDay, forbid: cfg.forbid, teacher: &teacherName})
			}
		}
		if len(blocks) > 0 {
			plans = append(plans, generateClassPlan{id: cid, name: class.Name, blocks: blocks})
		}
	}

	if len(plans) == 0 {
		warnings = append(warnings, generateNoIntersectWarning)
		return &GenerateSchoolResult{Success: false, Warnings: warnings, Classes: []GenerateSchoolClass{}}, nil
	}

	// 教师不可用时段（学校级）：排课时视为该教师已被占用
	busyTemplate := map[string]map[int]map[int]bool{}
	var unavailables []models.TimetableTeacherUnavailability
	if err := s.db.Where("school_id = ?", schoolID).Find(&unavailables).Error; err != nil {
		return nil, err
	}
	for _, u := range unavailables {
		if busyTemplate[u.TeacherName] == nil {
			busyTemplate[u.TeacherName] = map[int]map[int]bool{}
		}
		if busyTemplate[u.TeacherName][u.Weekday] == nil {
			busyTemplate[u.TeacherName][u.Weekday] = map[int]bool{}
		}
		busyTemplate[u.TeacherName][u.Weekday][u.PeriodIndex] = true
	}

	rng := s.newRand()
	for attempt := 0; attempt < generateSchoolMaxAttempts; attempt++ {
		teacherBusy := cloneTeacherBusy(busyTemplate)
		results := []generateClassResult{}
		ok := true

		for _, plan := range plans {
			grid := map[int]map[int]string{}
			perDay := map[int]map[string]int{}
			for _, d := range days {
				grid[d] = map[int]string{}
				perDay[d] = map[string]int{}
			}

			for _, block := range plan.blocks {
				if !placeBlock(rng, block, days, periodIndexes, isAm, grid, perDay, teacherBusy) {
					ok = false
					break
				}
			}
			if !ok {
				break
			}

			entries := []GenerateEntry{}
			for _, day := range days {
				for _, period := range sortedIntKeys(grid[day]) {
					subject := grid[day][period]
					entries = append(entries, GenerateEntry{
						Weekday: day, PeriodIndex: period, WeekType: "all", SubjectName: subject,
						TeacherName: optionalString(teacherOf[plan.id][subject]),
					})
				}
			}
			results = append(results, generateClassResult{id: plan.id, name: plan.name, entries: entries})
		}

		if !ok {
			continue
		}

		if commit {
			if err := s.commitSchoolResults(schoolID, results); err != nil {
				return nil, err
			}
		}

		summaries := make([]GenerateSchoolClass, 0, len(results))
		for _, res := range results {
			summaries = append(summaries, GenerateSchoolClass{
				ClassID: res.id, ClassName: res.name, EntryCount: len(res.entries),
			})
		}
		return &GenerateSchoolResult{Success: true, Warnings: warnings, Classes: summaries}, nil
	}

	warnings = append(warnings, generateSchoolRetryWarning)
	return &GenerateSchoolResult{Success: false, Warnings: warnings, Classes: []GenerateSchoolClass{}}, nil
}

// commitSchoolResults 事务内落库全校排课结果并作废相关班级的待审申请（同 Laravel generateSchool 的 commit 分支）。
func (s *TimetableService) commitSchoolResults(schoolID uint, results []generateClassResult) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txSvc := TimetableService{db: tx}

		classIDs := make([]uint, 0, len(results))
		for _, res := range results {
			subjectNames := []string{}
			for _, e := range res.entries {
				subjectNames = appendUniqueString(subjectNames, e.SubjectName)
			}
			subjects := make([]TimetableSubjectInput, 0, len(subjectNames))
			for _, name := range subjectNames {
				subjects = append(subjects, TimetableSubjectInput{Name: name})
			}
			entries := make([]TimetableEntryInput, 0, len(res.entries))
			for _, e := range res.entries {
				entries = append(entries, TimetableEntryInput{
					Weekday: e.Weekday, PeriodIndex: e.PeriodIndex, WeekType: e.WeekType,
					SubjectName: e.SubjectName, TeacherName: e.TeacherName, Room: e.Room,
				})
			}

			if _, err := txSvc.Save(res.id, schoolID, TimetablePayload{
				Subjects: subjects,
				Periods:  []TimetablePeriodInput{},
				Entries:  entries,
			}); err != nil {
				return err
			}
			classIDs = append(classIDs, res.id)
		}

		return tx.Model(&models.TimetableChangeRequest{}).
			Where("class_id IN ? AND status = ?", classIDs, models.TimetableChangePending).
			Updates(map[string]any{
				"status":      models.TimetableChangeRejected,
				"review_note": "管理员全校自动排课，本申请自动作废",
				"reviewed_at": util.Now(),
			}).Error
	})
}

// tryGenerateBlocks 单轮贪心：随机块序 + 候选位置按「同科目当日记数少者优先」，失败返回 false。
func tryGenerateBlocks(rng *rand.Rand, blocks []generateBlock, days, periodIndexes []int, isAm map[int]bool) ([]GenerateEntry, bool) {
	shuffled := make([]generateBlock, len(blocks))
	copy(shuffled, blocks)
	rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	grid := map[int]map[int]string{}
	perDay := map[int]map[string]int{}
	for _, d := range days {
		grid[d] = map[int]string{}
		perDay[d] = map[string]int{}
	}

	for _, block := range shuffled {
		if !placeBlock(rng, block, days, periodIndexes, isAm, grid, perDay, nil) {
			return nil, false
		}
	}

	entries := []GenerateEntry{}
	for _, day := range days {
		for _, period := range sortedIntKeys(grid[day]) {
			entries = append(entries, GenerateEntry{
				Weekday: day, PeriodIndex: period, WeekType: "all", SubjectName: grid[day][period],
			})
		}
	}
	return entries, true
}

// placeBlock 放置一个课时块（贪心 + 候选评分；单班 / 全校共用，全校同时占教师时段）。
func placeBlock(rng *rand.Rand, block generateBlock, days, periodIndexes []int, isAm map[int]bool, grid map[int]map[int]string, perDay map[int]map[string]int, teacherBusy map[string]map[int]map[int]bool) bool {
	type candidate struct {
		day   int
		cells []int
	}
	candidates := []candidate{}

	for _, day := range days {
		for i := 0; i < len(periodIndexes); i++ {
			if i+block.size > len(periodIndexes) {
				continue
			}
			cells := periodIndexes[i : i+block.size]

			ok := true
			for _, cell := range cells {
				if _, used := grid[day][cell]; used {
					ok = false
					break
				}
				if block.forbid[cell] {
					ok = false
					break
				}
				if block.teacher != nil && blockTeacherBusy(teacherBusy, *block.teacher, day, cell) {
					ok = false // 教师冲突：该教师此时段已在其他班上课
					break
				}
				if !sessionMatches(block.session, isAm, cell) {
					ok = false
					break
				}
			}
			// 连堂不允许跨上 / 下午（如第 4 节上午 + 第 5 节下午）
			if ok && block.size == 2 {
				first, firstKnown := isAm[cells[0]]
				second, secondKnown := isAm[cells[1]]
				if firstKnown && secondKnown && first != second {
					ok = false
				}
			}
			if ok && perDay[day][block.subject]+block.size > block.maxPerDay {
				ok = false
			}
			if ok {
				candidates = append(candidates, candidate{day: day, cells: append([]int{}, cells...)})
			}
		}
	}

	if len(candidates) == 0 {
		return false
	}

	rng.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	sort.SliceStable(candidates, func(i, j int) bool {
		return perDay[candidates[i].day][block.subject] < perDay[candidates[j].day][block.subject]
	})

	best := candidates[0]
	for _, cell := range best.cells {
		grid[best.day][cell] = block.subject
		perDay[best.day][block.subject]++

		if block.teacher != nil && teacherBusy != nil {
			if teacherBusy[*block.teacher] == nil {
				teacherBusy[*block.teacher] = map[int]map[int]bool{}
			}
			if teacherBusy[*block.teacher][best.day] == nil {
				teacherBusy[*block.teacher][best.day] = map[int]bool{}
			}
			teacherBusy[*block.teacher][best.day][cell] = true
		}
	}

	return true
}

// sessionMatches 时段限定匹配：any 恒真；am/pm 需已知该节上 / 下午，未知（无时间）则放行。
func sessionMatches(session string, isAm map[int]bool, cell int) bool {
	if session == "any" {
		return true
	}
	am, known := isAm[cell]
	if !known {
		return true
	}
	if session == "am" {
		return am
	}
	return !am
}

// normalizeGenerateDays 归一化上课日：取默认值（nil → 1-5）、过滤非 1-7、按首次出现去重。
func normalizeGenerateDays(raw []int) []int {
	if raw == nil {
		raw = []int{1, 2, 3, 4, 5}
	}
	days := []int{}
	seen := map[int]bool{}
	for _, d := range raw {
		if d < 1 || d > 7 || seen[d] {
			continue
		}
		seen[d] = true
		days = append(days, d)
	}
	return days
}

// normalizeGenerateSession 归一化时段限定：非 any/am/pm 一律回退 any。
func normalizeGenerateSession(raw string) string {
	if raw == "am" || raw == "pm" {
		return raw
	}
	return "any"
}

// normalizeGenerateMaxPerDay 归一化每日上限：最小 1；连堂至少 2。
func normalizeGenerateMaxPerDay(raw int, double bool) int {
	maxPerDay := raw
	if maxPerDay < 1 {
		maxPerDay = 1
	}
	if double && maxPerDay < 2 {
		maxPerDay = 2 // 连堂本身占同日 2 节
	}
	return maxPerDay
}

// generateGridMeta 节次索引列表与「上午」标记（时间未知 / 无法解析的节次不入表）。
func generateGridMeta(periods []models.ClassPeriod) ([]int, map[int]bool) {
	periodIndexes := make([]int, 0, len(periods))
	isAm := map[int]bool{}
	for _, p := range periods {
		periodIndexes = append(periodIndexes, p.PeriodIndex)
		if am, known := periodIsAm(p.StartTime); known {
			isAm[p.PeriodIndex] = am
		}
	}
	return periodIndexes, isAm
}

// periodIsAm 判定某节次是否上午：start_time 早于 12:00 为上午；无法解析 / 空串视为未知。
func periodIsAm(startTime string) (bool, bool) {
	normalized := normalizeTime(startTime)
	if normalized == "" || len(normalized) != 5 {
		return false, false
	}
	hour := atoiOrZero(normalized[0:2])
	minute := atoiOrZero(normalized[3:5])
	if hour > 23 || minute > 59 {
		return false, false
	}
	return hour*60+minute < 12*60, true
}

// intSet 转成集合（等价 Laravel `array_unique(array_map('intval', ...))` 的成员判断口径）。
func intSet(values []int) map[int]bool {
	set := map[int]bool{}
	for _, v := range values {
		set[v] = true
	}
	return set
}

// sortedIntKeys 整数键升序（等价 PHP ksort）。
func sortedIntKeys(m map[int]string) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// appendUniqueString 按首次出现顺序追加去重（等价 array_unique + array_values）。
func appendUniqueString(list []string, value string) []string {
	for _, item := range list {
		if item == value {
			return list
		}
	}
	return append(list, value)
}

// cloneTeacherBusy 深拷贝教师占用表（每次重试独立；Go 的 map 是引用类型，不能直接赋值）。
func cloneTeacherBusy(src map[string]map[int]map[int]bool) map[string]map[int]map[int]bool {
	dst := make(map[string]map[int]map[int]bool, len(src))
	for teacher, byDay := range src {
		days := make(map[int]map[int]bool, len(byDay))
		for day, cells := range byDay {
			periods := make(map[int]bool, len(cells))
			for period, busy := range cells {
				periods[period] = busy
			}
			days[day] = periods
		}
		dst[teacher] = days
	}
	return dst
}

// blockTeacherBusy 教师该时段是否已被占用。
func blockTeacherBusy(teacherBusy map[string]map[int]map[int]bool, teacher string, day, cell int) bool {
	if teacherBusy == nil {
		return false
	}
	return teacherBusy[teacher][day][cell]
}

// newRand 本次排课使用的随机源：注入了种子源则取其中下一个种子，否则按时间种子。
func (s *TimetableService) newRand() *rand.Rand {
	if s.seedSource != nil {
		return rand.New(rand.NewSource(s.seedSource()))
	}
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
