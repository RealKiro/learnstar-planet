// 班级大屏码生成服务。
//
// 规则（与 Laravel DisplayCodeService 逐条对齐）：`前缀` + 年级数字 + 班号数字，不补零。
// 前缀默认 LS，可在学校设置 display_code_prefix 自定义（2-4 个英文字母，本校统一，如 BJ / LE）。
// 例：一年级1班 → LS11，二年级3班 → LS23。确定性、不随机、同年级班号仅支持 1-9。
//
// 与 Laravel 的有意差异：原实现每处写码都同步维护 Cache 的「班级码 → 班级 ID」映射
// （Cache::put/forget），Go 端没有 Cache，改为登录时直接查库 + 确定性计算匹配
// （见 DisplayService.resolveCodeToClassID），故此处不再有缓存写入。
package services

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// DisplayCodeDefaultPrefix 默认班级码前缀。
const DisplayCodeDefaultPrefix = "LS"

// displayGradeMap 年级 → 数字（有序：PHP GRADE_MAP 的遍历顺序即子串匹配优先级）。
var displayGradeMap = [][2]string{
	{"一年级", "1"}, {"二年级", "2"}, {"三年级", "3"},
	{"四年级", "4"}, {"五年级", "5"}, {"六年级", "6"},
	{"七年级", "7"}, {"八年级", "8"}, {"九年级", "9"},
	{"高一", "10"}, {"高二", "11"}, {"高三", "12"},
}

// displayChineseNumbers 中文数字 → 阿拉伯数字。
var displayChineseNumbers = map[string]string{
	"一": "1", "二": "2", "三": "3", "四": "4", "五": "5",
	"六": "6", "七": "7", "八": "8", "九": "9", "十": "10",
}

var (
	reGradeDigit     = regexp.MustCompile(`(\d+)`)
	reClassNoParen   = regexp.MustCompile(`[（(]\s*(\d+)\s*[)）]\s*班`)
	reClassNoPlain   = regexp.MustCompile(`(\d+)\s*班`)
	reClassNoOrdinal = regexp.MustCompile(`第([一二三四五六七八九十]+)班`)
	reClassNoChinese = regexp.MustCompile(`([一二三四五六七八九十]+)班`)
	// reDisplayCodePattern 校验「已成形」的班级码前缀（大写字母 2-4 位）。
	reDisplayCodePattern = regexp.MustCompile(`^[A-Z]{2,4}$`)
	// reDisplayCodeInput 校验管理员提交的前缀入参（大小写字母 2-4 位）。
	reDisplayCodeInput = regexp.MustCompile(`^[A-Za-z]{2,4}$`)
)

// GradeDigit 年级 → 数字（'三年级'→'3'，不补零）。
// 仅接受 1-9 年级（与 Laravel 一致：暂不考虑 10+），无法解析或超出范围返回 ""。
func GradeDigit(grade string) string {
	if grade == "" {
		return ""
	}

	digit := ""
	for _, kv := range displayGradeMap {
		if strings.Contains(grade, kv[0]) {
			digit = kv[1]
			break
		}
	}
	if digit == "" {
		if m := reGradeDigit.FindStringSubmatch(grade); m != nil {
			digit = m[1]
		}
	}

	return digitsInRange(digit, 9)
}

// ClassNo 班级名 → 班号数字（不补零），支持多种格式。
// 仅接受 1-9 班号，无法解析或超出范围返回 ""。
func ClassNo(name string) string {
	num := ""
	switch {
	case reClassNoParen.MatchString(name):
		// 1) 标准：X年级（N）班 / X年级(N)班
		num = reClassNoParen.FindStringSubmatch(name)[1]
	case reClassNoPlain.MatchString(name):
		// 2) 无括号数字：一年级1班 / 一年级 1 班
		num = reClassNoPlain.FindStringSubmatch(name)[1]
	case reClassNoOrdinal.MatchString(name):
		// 3) 中文序数：一年级第一班
		num = ChineseNumberToArabic(reClassNoOrdinal.FindStringSubmatch(name)[1])
	case reClassNoChinese.MatchString(name):
		// 4) 中文数字：一年级一班
		num = ChineseNumberToArabic(reClassNoChinese.FindStringSubmatch(name)[1])
	}

	return digitsInRange(num, 9)
}

// ChineseNumberToArabic 中文数字 → 阿拉伯数字（'十'→10，'二十三'→23，未知→'0'）。
func ChineseNumberToArabic(cn string) string {
	if cn == "十" {
		return "10"
	}
	if idx := strings.Index(cn, "十"); idx >= 0 {
		tensPart := cn[:idx]
		onesPart := cn[idx+len("十"):]

		tens := 1
		if tensPart != "" {
			tens = atoiOrZero(displayChineseNumbers[tensPart])
		}
		ones := 0
		if onesPart != "" {
			ones = atoiOrZero(displayChineseNumbers[onesPart])
		}

		return strconv.Itoa(tens*10 + ones)
	}

	if v, ok := displayChineseNumbers[cn]; ok {
		return v
	}
	return "0"
}

// IsValidDisplayPrefix 判断是否为合法的「已成形」班级码前缀（大写 2-4 字母）。
func IsValidDisplayPrefix(prefix string) bool { return reDisplayCodePattern.MatchString(prefix) }

// IsValidDisplayPrefixInput 判断管理员提交的前缀入参是否合法（2-4 个英文字母）。
func IsValidDisplayPrefixInput(prefix string) bool { return reDisplayCodeInput.MatchString(prefix) }

// DisplayCodeService 班级码生成 / 批量重置服务。
type DisplayCodeService struct {
	db *gorm.DB
}

// NewDisplayCodeService 创建班级码服务。
func NewDisplayCodeService(db *gorm.DB) *DisplayCodeService {
	return &DisplayCodeService{db: db}
}

// SchoolPrefix 学校自定义前缀（settings.display_code_prefix），非法/未配置回退默认 LS。
func (s *DisplayCodeService) SchoolPrefix(schoolID uint) string {
	var school models.School
	if err := s.db.First(&school, schoolID).Error; err != nil {
		// Laravel 用 School::find()，学校不存在时取不到 settings → 同样回退默认前缀。
		return DisplayCodeDefaultPrefix
	}

	prefix := strings.ToUpper(school.SettingString("display_code_prefix", ""))
	if IsValidDisplayPrefix(prefix) {
		return prefix
	}
	return DisplayCodeDefaultPrefix
}

// Generate 生成班级码（前缀 + 年级 + 班号，如 LS11）。
// prefix 为 nil 时取学校前缀（同 Laravel 的 `$prefix ?? self::schoolPrefix(...)`：
// 传空串表示「不加前缀」，不触发学校前缀回退）。
// 年级无法解析、班号超限（>9，含兜底序号）返回 ""——同年级超过 10 个班无法编码。
func (s *DisplayCodeService) Generate(room *models.ClassRoom, prefix *string) string {
	grade := GradeDigit(room.Grade)
	if grade == "" {
		return ""
	}

	class := ClassNo(room.Name)
	if class == "" {
		class = s.FallbackClassNo(room)
	}
	if !isAllDigits(class) || atoiOrZero(class) > 9 {
		return ""
	}

	if prefix == nil {
		p := s.SchoolPrefix(room.SchoolID)
		prefix = &p
	}

	return *prefix + grade + class
}

// FallbackClassNo 班号兜底：班级名无法解析班号时，用同校同年级内按 ID 排序的序号（1 起）。
func (s *DisplayCodeService) FallbackClassNo(room *models.ClassRoom) string {
	var ids []uint
	if err := s.db.Model(&models.ClassRoom{}).
		Where("school_id = ? AND grade = ?", room.SchoolID, room.Grade).
		Order("id ASC").Pluck("id", &ids).Error; err != nil {
		return "1"
	}

	for i, id := range ids {
		if id == room.ID {
			return strconv.Itoa(i + 1)
		}
	}
	return "1"
}

// DisplayCodeRegenerateResult 批量重生成结果（字段名同 Laravel regenerateAll 返回值）。
type DisplayCodeRegenerateResult struct {
	Regenerated int             `json:"regenerated"`
	Skipped     int             `json:"skipped"`
	Conflicts   map[uint]string `json:"conflicts"`
}

// RegenerateAll 批量重生成班级码（确定性：同一班级每次结果一致）。
//
// schoolID 为 nil 时处理全校库；prefix 为空时按学校取默认前缀（未指定学校则 LS），
// 非法前缀回退 LS。同校内重复码（同年级同班号多条）计入 conflicts 并跳过、不写库。
// 遍历顺序：Laravel 未显式排序，Go 端固定按 id ASC 以保证结果可复现。
func (s *DisplayCodeService) RegenerateAll(schoolID *uint, prefix string) (*DisplayCodeRegenerateResult, error) {
	q := s.db.Model(&models.ClassRoom{})
	if schoolID != nil {
		q = q.Where("school_id = ?", *schoolID)
	}

	if prefix == "" {
		if schoolID != nil {
			prefix = s.SchoolPrefix(*schoolID)
		} else {
			prefix = DisplayCodeDefaultPrefix
		}
	}
	if !IsValidDisplayPrefix(prefix) {
		prefix = DisplayCodeDefaultPrefix
	}

	var rooms []models.ClassRoom
	if err := q.Order("id ASC").Find(&rooms).Error; err != nil {
		return nil, err
	}

	result := &DisplayCodeRegenerateResult{Conflicts: map[uint]string{}}
	// used 校维度的已用码：schoolID → code → 班级 ID。
	used := map[uint]map[string]uint{}

	for i := range rooms {
		room := &rooms[i]

		grade := GradeDigit(room.Grade)
		if grade == "" {
			result.Skipped++
			continue
		}

		class := ClassNo(room.Name)
		if class == "" {
			class = s.FallbackClassNo(room)
		}
		// 班号超限（同年级 >10 个班）无法编码，跳过并计数。
		if !isAllDigits(class) || atoiOrZero(class) > 9 {
			result.Skipped++
			continue
		}

		code := prefix + grade + class

		if used[room.SchoolID] == nil {
			used[room.SchoolID] = map[string]uint{}
		}
		if owner, ok := used[room.SchoolID][code]; ok && owner != room.ID {
			result.Conflicts[room.ID] = code
			result.Skipped++
			continue
		}
		used[room.SchoolID][code] = room.ID

		now := util.Now()
		room.DisplayCode = code
		room.DisplayCodeUpdatedAt = &now
		if err := s.db.Save(room).Error; err != nil {
			return nil, err
		}
		result.Regenerated++
	}

	return result, nil
}

// digitsInRange 校验 raw 全为 ASCII 数字且数值在 1..max 之间；否则返回 ""（保留原始串，不补零）。
func digitsInRange(raw string, max int) string {
	if !isAllDigits(raw) {
		return ""
	}
	n := atoiOrZero(raw)
	if n < 1 || n > max {
		return ""
	}
	return raw
}

// isAllDigits 等价 PHP ctype_digit（空串返回 false）。
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func atoiOrZero(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
