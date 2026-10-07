// 管理端运维服务（二）：批量账号、教师批量创建 / CSV 导入 / 模板下载、学校 LOGO 上传。
//
// 忠实移植自 Laravel App\Http\Controllers\Api\SchoolAdminController：
// batchResetPassword / batchDeleteAccounts / batchCreateTeachers / importTeachers /
// downloadTeacherTemplate / uploadLogo，以及 App\Services\AuthService::createTeacherAccounts
// （智能去重用户名/昵称）。
package services

import (
	"bytes"
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	// DefaultTeacherPassword 批量重置密码的默认值（同 Laravel batchResetPassword 的 'ls123456'）。
	DefaultTeacherPassword = "ls123456"
	// MinPasswordLen / MaxPasswordLen 密码长度限制（同 Laravel min:6|max:50）。
	MinPasswordLen = 6
	MaxPasswordLen = 50
	// maxLogoBytes LOGO 大小上限（同 Laravel max:2048，单位 KB）。
	maxLogoBytes = 2048 * 1024
)

// validClassTeacherRoles 班级教师角色白名单（同 Laravel in:head_teacher,...）。
var validClassTeacherRoles = map[string]bool{
	"head_teacher":    true,
	"co_teacher":      true,
	"subject_teacher": true,
	"grade_lead":      true,
	"admin_director":  true,
}

var emailShape = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// ============================================================
// 批量账号操作（教师/家长通用，role 区分）
// ============================================================

// BatchResetPassword 批量重置账号密码，返回实际被重置的账号数。
//
// 忠实移植自 Laravel batchResetPassword：只影响「本校 + role=teacher + 选中 ID」；
// 一个都没命中返回 404「未找到所选账号」；密码只哈希一次（同 Laravel 单个 Hash::make）。
//
// 有意差异（本批已收窄）：Go 的 users 表现在也有 `plain_password` 列，
// 重置时同步写入明文（与 Laravel batchResetPassword 的 password_hash + password_changed
// 语义一致；Laravel 该动作**不**写 plain_password，明文只在单条 resetTeacherPassword
// 与创建教师时写入——Go 端选择两条重置路径都写，见 README「本批有意差异」）。
func (o *AdminOps) BatchResetPassword(schoolID uint, ids []uint, newPassword string) (int64, error) {
	if len(ids) == 0 {
		return 0, ErrUnprocessable("参数错误")
	}
	if newPassword == "" {
		newPassword = DefaultTeacherPassword
	}
	if len([]rune(newPassword)) < MinPasswordLen || len([]rune(newPassword)) > MaxPasswordLen {
		return 0, ErrUnprocessable("参数错误")
	}

	query := o.db.Model(&models.User{}).
		Where("school_id = ? AND role = ? AND id IN ?", schoolID, "teacher", ids)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, ErrNotFound("未找到所选账号")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	// 注意：不能复用上面的 query 变量做 Updates——GORM 会因残留的 Model 上下文生成
	// `UPDATE users ... FROM users WHERE ...` 而报「ambiguous column name」（sqlite 实测）。
	updateQuery := o.db.Model(&models.User{}).
		Where("school_id = ? AND role = ? AND id IN ?", schoolID, "teacher", ids)
	if err := updateQuery.Updates(map[string]any{
		"password_hash":    string(hash),
		"plain_password":   newPassword,
		"password_changed": false,
	}).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// BatchDeleteAccounts 批量删除教师账号（硬删除，释放 username 唯一索引）。
//
// 与 Laravel batchDeleteAccounts 的差异：
//  1. Laravel 用 `$acc->delete()`（User 无 SoftDeletes → 实际也是硬删除），但不解除班级关联，
//     会留下指向不存在用户的 `class_rooms.teacher_id`；Go 端在删除前把 teacher_id 置空并清理
//     `class_room_teachers` 关联（同 disableTeacher 的做法），避免脏引用。
//  2. Laravel 的批量删除**没有**机器人账号守卫（只有单删 disableTeacher 有）；Go 端按项目决策 13
//     「机器人账号不可从后台删除」跳过 API 机器人账号，并在返回值里给出 protected_count。
func (o *AdminOps) BatchDeleteAccounts(schoolID uint, ids []uint) (deleted int64, protected int, err error) {
	if len(ids) == 0 {
		return 0, 0, ErrUnprocessable("参数错误")
	}

	var accounts []models.User
	if err := o.db.Where("school_id = ? AND role = ? AND id IN ?", schoolID, "teacher", ids).
		Find(&accounts).Error; err != nil {
		return 0, 0, err
	}
	if len(accounts) == 0 {
		return 0, 0, ErrNotFound("未找到所选账号")
	}

	deletable := make([]uint, 0, len(accounts))
	for _, acc := range accounts {
		if acc.IsAPIBot {
			protected++
			continue
		}
		deletable = append(deletable, acc.ID)
	}
	if len(deletable) == 0 {
		return 0, protected, ErrForbidden("API 机器人账号不可删除。如需停用，请在 .env 中设置 BOT_ENABLED=false 后重启")
	}

	err = o.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ClassRoom{}).Where("teacher_id IN ?", deletable).
			Update("teacher_id", nil).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id IN ?", deletable).Delete(&models.ClassRoomTeacher{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", deletable).Delete(&models.User{}).Error
	})
	if err != nil {
		return 0, protected, err
	}
	return int64(len(deletable)), protected, nil
}

// DeleteTeacher 删除单个教师账号（忠实移植 Laravel `SchoolAdminController::disableTeacher`）。
//
//   - 教师须在本校且 `role=teacher`，否则 404「教师不存在」；
//   - API 机器人账号 → 403「API 机器人账号不可删除。如需停用，请在 .env 中设置 BOT_ENABLED=false 后重启」；
//   - 解除班级关联：`class_rooms.teacher_id` 置空 + 删除 `class_room_teachers` 关联
//     （同 Laravel 的两条语句；Go 端放在同一事务里，避免中途失败留下脏引用）；
//   - 硬删除用户（Go 的 User 无软删除，与 Laravel `$teacher->delete()` 的实际效果一致）。
func (o *AdminOps) DeleteTeacher(schoolID, teacherID uint) error {
	var teacher models.User
	err := o.db.Where("school_id = ? AND role = ? AND id = ?", schoolID, "teacher", teacherID).
		First(&teacher).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound("教师不存在")
	}
	if err != nil {
		return err
	}
	if teacher.IsAPIBot {
		return ErrForbidden("API 机器人账号不可删除。如需停用，请在 .env 中设置 BOT_ENABLED=false 后重启")
	}

	return o.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ClassRoom{}).Where("teacher_id = ?", teacher.ID).
			Update("teacher_id", nil).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", teacher.ID).Delete(&models.ClassRoomTeacher{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.User{}, teacher.ID).Error
	})
}

// ============================================================
// 教师批量创建（含班级分配）
// ============================================================

// TeacherAssignmentInput 一名教师的班级分配（同 Laravel batch-create 的 assignments.*）。
type TeacherAssignmentInput struct {
	ClassID uint
	Role    string
	Subject string
}

// TeacherInput 创建教师的入参（字段名同 Laravel createTeacherAccounts 的数组键）。
type TeacherInput struct {
	Name        string
	Nickname    string
	Subject     string
	GradeTeam   string
	Phone       string
	Email       string
	Username    string
	Password    string
	Assignments []TeacherAssignmentInput
}

// TeacherCreated 创建结果（同 Laravel createTeacherAccounts 的返回结构，中文键不变）。
type TeacherCreated struct {
	ID              uint   `json:"id"`
	Username        string `json:"username"`
	Nickname        string `json:"nickname"`
	InitialPassword string `json:"initial_password"`
	Name            string `json:"name"`
}

// ValidateTeacherInputs 校验批量创建教师的入参（同 Laravel batch-create 的 validator）。
func (o *AdminOps) ValidateTeacherInputs(inputs []TeacherInput) error {
	if len(inputs) == 0 {
		return ErrUnprocessable("参数错误：teachers 至少 1 条")
	}
	for _, in := range inputs {
		name := strings.TrimSpace(in.Name)
		if name == "" {
			return ErrUnprocessable("参数错误：教师姓名必填")
		}
		if len([]rune(name)) > 50 {
			return ErrUnprocessable("参数错误：教师姓名不能超过 50 字")
		}
		if len([]rune(in.Nickname)) > 80 {
			return ErrUnprocessable("参数错误：昵称不能超过 80 字")
		}
		if len([]rune(in.Subject)) > 50 {
			return ErrUnprocessable("参数错误：科目不能超过 50 字")
		}
		if len([]rune(in.GradeTeam)) > 50 {
			return ErrUnprocessable("参数错误：年级团队不能超过 50 字")
		}
		if len([]rune(in.Phone)) > 30 {
			return ErrUnprocessable("参数错误：手机号不能超过 30 字")
		}
		if in.Email != "" && (len([]rune(in.Email)) > 100 || !emailShape.MatchString(in.Email)) {
			return ErrUnprocessable("参数错误：邮箱格式不正确")
		}
		if strings.TrimSpace(in.Username) != "" && len([]rune(in.Username)) > 50 {
			return ErrUnprocessable("参数错误：用户名不能超过 50 字")
		}
		if in.Password == "" {
			return ErrUnprocessable("参数错误：初始密码必填")
		}
		if len([]rune(in.Password)) < MinPasswordLen || len([]rune(in.Password)) > MaxPasswordLen {
			return ErrUnprocessable("参数错误：密码长度需在 6-50 位之间")
		}
		for _, a := range in.Assignments {
			if a.ClassID == 0 {
				return ErrUnprocessable("参数错误：班级 ID 必填")
			}
			if !validClassTeacherRoles[a.Role] {
				return ErrUnprocessable("参数错误：班级角色不合法")
			}
		}
	}
	return nil
}

// BatchCreateTeachers 批量创建教师账号并写入班级分配。
//
// 忠实移植自 Laravel batchCreateTeachers：
// 1) createTeacherAccounts 建号（用户名/昵称智能去重）；
// 2) 按教师**姓名**回查入参中的 assignments，逐条 upsert `class_room_teachers`；
// 3) role = head_teacher 时把 `class_rooms.teacher_id` 指向该教师（仅限本校班级）。
//
// 有意差异：Go 无拼音库，昵称默认取「姓名」本身（Laravel 为姓名拼音，PinyinService 缺失时
// 同样回退为姓名），去重规则仍为 `_2/_3` 递增。
func (o *AdminOps) BatchCreateTeachers(schoolID uint, inputs []TeacherInput) ([]TeacherCreated, error) {
	if err := o.ValidateTeacherInputs(inputs); err != nil {
		return nil, err
	}

	created, err := o.CreateTeacherAccounts(schoolID, inputs)
	if err != nil {
		return nil, err
	}

	// Laravel 用 `collect($teachers)->firstWhere('name', $teacher['name'])` 回查入参，
	// 因此同名教师的分配会落到第一个同名入参上（此处保持一致）。
	byName := map[string]TeacherInput{}
	for _, in := range inputs {
		if _, ok := byName[strings.TrimSpace(in.Name)]; !ok {
			byName[strings.TrimSpace(in.Name)] = in
		}
	}

	for _, t := range created {
		in, ok := byName[t.Name]
		if !ok || len(in.Assignments) == 0 {
			continue
		}
		if err := o.applyAssignments(schoolID, t.ID, in.Assignments); err != nil {
			return nil, err
		}
	}
	return created, nil
}

// applyAssignments 写入班级分配（upsert 关联 + head_teacher 回写 class_rooms.teacher_id）。
func (o *AdminOps) applyAssignments(schoolID, userID uint, assignments []TeacherAssignmentInput) error {
	for _, a := range assignments {
		var class models.ClassRoom
		if err := o.db.First(&class, a.ClassID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUnprocessable("参数错误：班级不存在")
			}
			return err
		}
		if err := o.upsertClassRoomTeacher(class.ID, userID, a.Role, a.Subject); err != nil {
			return err
		}
		if a.Role == "head_teacher" && class.SchoolID == schoolID {
			// Laravel：ClassRoom::where('id',..)->where('school_id',..)->update(['teacher_id'=>..])
			if err := o.db.Model(&models.ClassRoom{}).
				Where("id = ? AND school_id = ?", class.ID, schoolID).
				Update("teacher_id", userID).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// upsertClassRoomTeacher 等价 Laravel `ClassRoomTeacher::updateOrCreate(['class_room_id','user_id'], [...])`。
func (o *AdminOps) upsertClassRoomTeacher(classID, userID uint, role, subject string) error {
	var existing models.ClassRoomTeacher
	err := o.db.Where("class_room_id = ? AND user_id = ?", classID, userID).First(&existing).Error
	if err == nil {
		return o.db.Model(&models.ClassRoomTeacher{}).Where("id = ?", existing.ID).
			Updates(map[string]any{"role": role, "subject": subject}).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return o.db.Create(&models.ClassRoomTeacher{
		ClassRoomID: classID,
		UserID:      userID,
		Role:        role,
		Subject:     subject,
	}).Error
}

// CreateTeacherAccounts 建号（同 Laravel AuthService::createTeacherAccounts）。
// username 默认 = 姓名（同名追加 _2/_3…），nickname 默认 = 姓名（同上）；
// 密码为空时用 'ls123456'；逐个创建，任一失败即整体失败。
func (o *AdminOps) CreateTeacherAccounts(schoolID uint, inputs []TeacherInput) ([]TeacherCreated, error) {
	created := make([]TeacherCreated, 0, len(inputs))
	for _, in := range inputs {
		name := strings.TrimSpace(in.Name)
		if name == "" {
			continue
		}

		username := strings.TrimSpace(in.Username)
		if username == "" {
			u, err := o.uniqueValue(name, func(candidate string) (bool, error) {
				// username 在 Go schema 中是全局唯一索引（Laravel 为全校唯一且 username 列唯一），
				// 故此处按全局判重，避免落库撞唯一索引。
				var count int64
				if err := o.db.Model(&models.User{}).Where("username = ?", candidate).Count(&count).Error; err != nil {
					return false, err
				}
				return count > 0, nil
			})
			if err != nil {
				return nil, err
			}
			username = u
		}

		nickname := strings.TrimSpace(in.Nickname)
		if nickname == "" {
			n, err := o.uniqueValue(name, func(candidate string) (bool, error) {
				var count int64
				if err := o.db.Model(&models.User{}).
					Where("school_id = ? AND nickname = ?", schoolID, candidate).
					Count(&count).Error; err != nil {
					return false, err
				}
				return count > 0, nil
			})
			if err != nil {
				return nil, err
			}
			nickname = n
		}

		password := in.Password
		if password == "" {
			password = DefaultTeacherPassword
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		user := models.User{
			SchoolID:        schoolID,
			Role:            "teacher",
			Username:        username,
			PasswordHash:    string(hash),
			PlainPassword:   password,
			Name:            name,
			Nickname:        nickname,
			Subject:         in.Subject,
			GradeTeam:       in.GradeTeam,
			Phone:           in.Phone,
			Email:           in.Email,
			Status:          "active",
			PasswordChanged: false,
		}
		if err := o.db.Create(&user).Error; err != nil {
			return nil, err
		}

		created = append(created, TeacherCreated{
			ID:              user.ID,
			Username:        username,
			Nickname:        nickname,
			InitialPassword: password,
			Name:            user.Name,
		})
	}
	return created, nil
}

// uniqueValue 通用去重：候选被占用时依次追加 _2/_3/…（同 Laravel AuthService::makeUnique，
// 含 999 次后的随机后缀兜底）。
func (o *AdminOps) uniqueValue(base string, exists func(string) (bool, error)) (string, error) {
	candidate := base
	for i := 2; i <= 999; i++ {
		occupied, err := exists(candidate)
		if err != nil {
			return "", err
		}
		if !occupied {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s_%d", base, i)
	}
	// 极端兜底：超过 999 次重名时使用随机后缀。
	suffix, err := randomHex(4)
	if err != nil {
		return "", err
	}
	return base + "_" + suffix, nil
}

// randomHex 生成 n 字节的随机十六进制串。
func randomHex(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// ============================================================
// 教师 CSV 导入 / 模板下载
// ============================================================

// TeacherImportRow 教师导入的一行（预览结构与 Laravel importTeachers 的 preview 元素一致）。
type TeacherImportRow struct {
	Name      string `json:"name"`
	GradeTeam string `json:"grade_team"`
	Subject   string `json:"subject"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

// TeacherImportResult 导入结果；dry_run 时返回 preview，正式导入时返回 created。
// 该结构与 Laravel 一致，**不套 data 信封**（直接作为响应体顶层）。
type TeacherImportResult struct {
	Message string             `json:"message"`
	Total   int                `json:"total"`
	Preview []TeacherImportRow `json:"preview,omitempty"`
	Created []TeacherCreated   `json:"created,omitempty"`
}

// TeacherTemplateCSV 教师导入模板（带 UTF-8 BOM，同 Laravel downloadTeacherTemplate）。
func (o *AdminOps) TeacherTemplateCSV() string {
	// Laravel：BOM + 4 行 fputcsv（中文表头、英文字段名、两行示例），行尾 "\n"。
	return "\ufeff姓名,年级团队,科目,密码,手机号\n" +
		"name,grade_team,subject,password,phone\n" +
		"张老师,三年级团队,语文,star123456,13800138000\n" +
		"李老师,三年级团队,数学,,\n"
}

// ImportTeachersCSV 解析 CSV 并（在 dry_run=false 时）批量创建教师账号。
//
// 忠实移植自 Laravel importTeachers + parseTeacherFile + detectCsvDelimiter：
// 表头取首行，逐行按 `name/姓名`、`grade_team/年级团队/所属年级团队`、`subject/科目`、
// `password/密码`、`phone/手机号` 取值；姓名为空的行跳过；dry_run 默认 true（预览）。
//
// 有意差异：
//  1. Excel（xlsx/xls）不解析——Laravel 走 PhpSpreadsheet，Go 端无该依赖，直接 422 提示另存为 CSV；
//     （Laravel 在未装 PhpSpreadsheet 时 parseExcelFile 也返回空数组，同属「不支持」）；
//  2. 仅接受 UTF-8：Laravel 用 mb_detect_encoding + mb_convert_encoding 自动转 GBK，Go 端标准库
//     无 GBK 解码，检测到非 UTF-8 时返回 422 提示另存为 UTF-8。
func (o *AdminOps) ImportTeachersCSV(schoolID uint, filename string, content []byte, dryRun bool) (*TeacherImportResult, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".csv", ".txt", "":
		// 允许（无扩展名时按 CSV 解析）
	case ".xlsx", ".xls":
		return nil, ErrUnprocessable("Excel 文件暂不支持，请另存为 CSV（UTF-8）后重试")
	default:
		return nil, ErrUnprocessable("文件上传失败")
	}

	rows, err := parseTeacherCSV(content)
	if err != nil {
		return nil, err
	}

	if !dryRun {
		created := make([]TeacherCreated, 0, len(rows))
		for _, row := range rows {
			password := row.Password
			if password == "" {
				password = DefaultTeacherPassword
			}
			result, err := o.CreateTeacherAccounts(schoolID, []TeacherInput{{
				Name:      row.Name,
				GradeTeam: row.GradeTeam,
				Subject:   row.Subject,
				Password:  password,
				Phone:     row.Phone,
			}})
			if err != nil {
				return nil, err
			}
			created = append(created, result...)
		}
		return &TeacherImportResult{
			Message: fmt.Sprintf("已导入 %d 名教师", len(created)),
			Total:   len(created),
			Created: created,
		}, nil
	}

	return &TeacherImportResult{
		Message: fmt.Sprintf("预览模式：共 %d 条数据", len(rows)),
		Total:   len(rows),
		Preview: rows,
	}, nil
}

// parseTeacherCSV 解析教师导入 CSV（去 BOM、检测分隔符、按表头映射字段）。
func parseTeacherCSV(content []byte) ([]TeacherImportRow, error) {
	content = bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF})
	if len(content) == 0 {
		return []TeacherImportRow{}, nil
	}
	if !utf8.Valid(content) {
		return nil, ErrUnprocessable("文件编码不是 UTF-8，请另存为 UTF-8 后重试")
	}

	firstLine := string(content)
	if idx := strings.IndexAny(firstLine, "\r\n"); idx >= 0 {
		firstLine = firstLine[:idx]
	}
	reader := csv.NewReader(bytes.NewReader(content))
	reader.Comma = detectCSVDelimiter(firstLine)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, ErrUnprocessable("CSV 解析失败：" + err.Error())
	}
	if len(records) == 0 {
		return []TeacherImportRow{}, nil
	}

	header := records[0]
	rows := make([]TeacherImportRow, 0, len(records)-1)
	for _, rec := range records[1:] {
		cell := map[string]string{}
		for i, key := range header {
			if i < len(rec) {
				cell[strings.TrimSpace(key)] = rec[i]
			}
		}
		name := strings.TrimSpace(firstNonEmptyValue(cell, "name", "姓名"))
		if name == "" {
			continue
		}
		rows = append(rows, TeacherImportRow{
			Name:      name,
			GradeTeam: strings.TrimSpace(firstNonEmptyValue(cell, "grade_team", "年级团队", "所属年级团队")),
			Subject:   strings.TrimSpace(firstNonEmptyValue(cell, "subject", "科目")),
			Password:  strings.TrimSpace(firstNonEmptyValue(cell, "password", "密码")),
			Phone:     strings.TrimSpace(firstNonEmptyValue(cell, "phone", "手机号")),
		})
	}
	return rows, nil
}

// detectCSVDelimiter 等价 Laravel detectCsvDelimiter：制表符 / 逗号 / 分号中取出现次数最多者，
// 全为 0 时回退逗号。
func detectCSVDelimiter(line string) rune {
	best := ','
	bestCount := 0
	for _, d := range []rune{'\t', ',', ';'} {
		if count := strings.Count(line, string(d)); count > bestCount {
			bestCount = count
			best = d
		}
	}
	return best
}

// firstNonEmptyValue 按顺序取第一个非空（去空白后）的值。
func firstNonEmptyValue(cell map[string]string, keys ...string) string {
	for _, k := range keys {
		if v, ok := cell[k]; ok && strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// ============================================================
// 学校 LOGO 上传
// ============================================================

// SaveSchoolLogo 保存学校 LOGO 原图并把相对路径写入 schools.logo_path。
//
// 校验同 Laravel uploadLogo（image + mimes:jpeg,png,gif,webp + max:2048）。
//
// 有意差异：Go 端**不做图片缩放/裁剪**（Laravel 用 intervention/image 处理），直接保存原文件；
// 落盘目录为 <UploadRoot>/schools/，写库路径为 <UploadURLPrefix>/schools/<文件名>。
func (o *AdminOps) SaveSchoolLogo(schoolID uint, filename string, content []byte) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return "", ErrUnprocessable("参数错误：LOGO 仅支持 jpeg/png/gif/webp 格式")
	}
	if len(content) == 0 {
		return "", ErrUnprocessable("参数错误：LOGO 文件为空")
	}
	if len(content) > maxLogoBytes {
		return "", ErrUnprocessable("参数错误：LOGO 不能超过 2MB")
	}
	if !strings.HasPrefix(http.DetectContentType(content), "image/") {
		return "", ErrUnprocessable("参数错误：LOGO 必须是图片文件")
	}

	var school models.School
	if err := o.db.First(&school, schoolID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound("学校不存在")
		}
		return "", err
	}

	dir := filepath.Join(o.UploadRoot, SchoolLogoSubdir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	suffix, err := randomHex(6)
	if err != nil {
		return "", err
	}
	name := fmt.Sprintf("logo_%s_%s%s", time.Now().Format("20060102150405"), suffix, ext)
	if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
		return "", err
	}

	logoPath := o.UploadURLPrefix + "/" + SchoolLogoSubdir + "/" + name
	if err := o.db.Model(&models.School{}).Where("id = ?", schoolID).
		Update("logo_path", logoPath).Error; err != nil {
		return "", err
	}
	return logoPath, nil
}
