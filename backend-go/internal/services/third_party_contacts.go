// 管理端通讯录拉取与导入（企微 / 当前学校所选第三方平台）。
//
// 忠实移植自 Laravel SchoolAdminController 的 wechatWorkContacts / importWechatWorkUsers /
// thirdPartyContacts / importThirdPartyUsers 四个动作（后两者中 thirdPartyContacts 走 provider，
// importThirdPartyUsers 直接复用企微导入逻辑）。
//
// 与 Laravel 的有意差异（逐条）：
//  1. 传输层 / 上游异常：Laravel `catch (\Throwable $e) → 400 {message: $e->getMessage()}`；
//     Go 端把错误包成 `ErrBadRequest`（400 + 原始中文文案），保持同样的状态码与文案。
//  2. 校验失败：Laravel 走 Validator → 422 `{message:"参数错误", errors:{...}}`（英文默认文案）；
//     Go 端沿用仓库约定：422「参数错误」+ `errors`（键名与 Laravel 一致，文案中文）。
//  3. 学校缺失：Laravel contacts 直接读 `$request->user()->school->id`（null 会 500），
//     import 显式 404「未找到学校」；Go 端两条路径统一返回 404「未找到学校」（不伪造 500）。
//  4. 空字符串代替 NULL：可空列（phone / email / gender）在 Go 端落空串（本仓既有约定），
//     Laravel 落 NULL；查重与展示口径不受影响。
package services

import (
	"fmt"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// ContactsImportTeacher 通讯录导入的教师一条（教师姓名必填，手机号 / 邮箱可空）。
type ContactsImportTeacher struct {
	Name   string  `json:"name"`
	Mobile *string `json:"mobile"`
	Phone  *string `json:"phone"`
	Email  *string `json:"email"`
}

// ContactsImportStudent 通讯录导入的学生一条（姓名与班级必填，性别可空）。
type ContactsImportStudent struct {
	Name    string  `json:"name"`
	ClassID *int    `json:"class_id"`
	Gender  *string `json:"gender"`
}

// ContactsImportRequest 通讯录导入请求体（同 Laravel `teachers` / `students` 数组）。
type ContactsImportRequest struct {
	Teachers []ContactsImportTeacher `json:"teachers"`
	Students []ContactsImportStudent `json:"students"`
}

// SkippedEntry 被跳过的条目（`{name, reason}`，同 Laravel skipped_teachers / skipped_students 元素）。
type SkippedEntry struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// ImportedTeacherAccount 导入后返回的教师账号（仅创建时可见初始密码，同 Laravel）。
type ImportedTeacherAccount struct {
	Name            string `json:"name"`
	Username        string `json:"username"`
	InitialPassword string `json:"initial_password"`
}

// ContactsImportData 导入结果数据（字段与顺序同 Laravel）。
type ContactsImportData struct {
	CreatedTeachers int                      `json:"created_teachers"`
	CreatedStudents int                      `json:"created_students"`
	SkippedTeachers []SkippedEntry           `json:"skipped_teachers"`
	SkippedStudents []SkippedEntry           `json:"skipped_students"`
	TeacherAccounts []ImportedTeacherAccount `json:"teacher_accounts"`
}

// ContactsImportResult 导入结果（message + data，同 Laravel 顶层结构）。
type ContactsImportResult struct {
	Message string             `json:"message"`
	Data    ContactsImportData `json:"data"`
}

// WechatWorkContacts 拉取企微通讯录（部门 + 成员），供前端预览导入。
func (s *ThirdParty) WechatWorkContacts(schoolID uint) (*ThirdPartyContacts, error) {
	school, err := s.SchoolByID(schoolID)
	if err != nil {
		return nil, err
	}
	if school == nil {
		return nil, ErrNotFound("未找到学校")
	}
	contacts, err := s.wechatWork.FetchContacts(school.ID)
	if err != nil {
		return nil, ErrBadRequest(err.Error())
	}
	return contacts, nil
}

// ThirdPartyContacts 拉取当前学校所选第三方平台的通讯录，并附带平台标识。
func (s *ThirdParty) ThirdPartyContacts(schoolID uint) (*ThirdPartyContacts, error) {
	school, err := s.SchoolByID(schoolID)
	if err != nil {
		return nil, err
	}
	if school == nil {
		return nil, ErrNotFound("未找到学校")
	}
	provider, err := s.manager.ProviderFor(school.SettingString("third_party_platform", ""))
	if err != nil {
		return nil, ErrBadRequest(err.Error())
	}
	contacts, err := provider.FetchContacts(school.ID)
	if err != nil {
		return nil, ErrBadRequest(err.Error())
	}
	contacts.Platform = provider.Key()
	return contacts, nil
}

// ImportContacts 从通讯录批量导入教师与学生（企微与第三方平台共用，同 Laravel
// importThirdPartyUsers 直接委派 importWechatWorkUsers）。
func (s *ThirdParty) ImportContacts(schoolID uint, req ContactsImportRequest) (*ContactsImportResult, error) {
	school, err := s.SchoolByID(schoolID)
	if err != nil {
		return nil, err
	}
	if school == nil {
		return nil, ErrNotFound("未找到学校")
	}

	if err := validateContactsImport(req); err != nil {
		return nil, err
	}

	skippedTeachers := make([]SkippedEntry, 0)
	skippedStudents := make([]SkippedEntry, 0)

	// ---- 教师：归一化手机号 + 请求内查重，重复导入不再生成 _2/_3 冗余账号 ----
	type normalizedTeacher struct {
		Name  string
		Phone string
		Email string
	}
	teachers := make([]normalizedTeacher, 0, len(req.Teachers))
	seenPhone := map[string]bool{}
	seenName := map[string]bool{}
	for _, t := range req.Teachers {
		// 通讯录用 mobile 字段，createTeacherAccounts 存 phone —— 统一映射，避免手机号丢失。
		rawPhone := ""
		if t.Phone != nil {
			rawPhone = *t.Phone
		} else if t.Mobile != nil {
			rawPhone = *t.Mobile
		}
		phone := stripPhoneSeparators(rawPhone)
		email := ""
		if t.Email != nil {
			email = *t.Email
		}
		if phone != "" {
			if seenPhone[phone] {
				continue
			}
			seenPhone[phone] = true
		} else if seenName[t.Name] {
			continue
		}
		seenName[t.Name] = true
		teachers = append(teachers, normalizedTeacher{Name: t.Name, Phone: phone, Email: email})
	}

	// ---- 库内查重：手机号已存在，或默认实名用户名（= 姓名）已存在 → 跳过 ----
	kept := make([]normalizedTeacher, 0, len(teachers))
	if len(teachers) > 0 {
		phoneList := make([]string, 0, len(teachers))
		nameList := make([]string, 0, len(teachers))
		for _, t := range teachers {
			if t.Phone != "" {
				phoneList = append(phoneList, t.Phone)
			}
			nameList = append(nameList, t.Name)
		}

		existingPhones := map[string]bool{}
		if len(phoneList) > 0 {
			var phones []string
			if err := s.db.Model(&models.User{}).
				Where("school_id = ? AND phone IN ?", school.ID, phoneList).
				Pluck("phone", &phones).Error; err != nil {
				return nil, err
			}
			for _, p := range phones {
				existingPhones[p] = true
			}
		}
		existingUsernames := map[string]bool{}
		if len(nameList) > 0 {
			var usernames []string
			if err := s.db.Model(&models.User{}).
				Where("school_id = ? AND username IN ?", school.ID, nameList).
				Pluck("username", &usernames).Error; err != nil {
				return nil, err
			}
			for _, u := range usernames {
				existingUsernames[u] = true
			}
		}

		for _, t := range teachers {
			if t.Phone != "" && existingPhones[t.Phone] {
				skippedTeachers = append(skippedTeachers, SkippedEntry{Name: t.Name, Reason: "手机号已存在"})
				continue
			}
			if existingUsernames[t.Name] {
				skippedTeachers = append(skippedTeachers, SkippedEntry{Name: t.Name, Reason: "同名账号已存在"})
				continue
			}
			kept = append(kept, t)
		}
	}

	teacherAccounts := make([]ImportedTeacherAccount, 0)
	createdTeachers := 0
	if len(kept) > 0 {
		inputs := make([]TeacherInput, 0, len(kept))
		for _, t := range kept {
			inputs = append(inputs, TeacherInput{Name: t.Name, Phone: t.Phone, Email: t.Email})
		}
		result, err := s.ops.CreateTeacherAccounts(school.ID, inputs)
		if err != nil {
			return nil, err
		}
		createdTeachers = len(result)
		for _, a := range result {
			teacherAccounts = append(teacherAccounts, ImportedTeacherAccount{
				Name:            a.Name,
				Username:        a.Username,
				InitialPassword: a.InitialPassword,
			})
		}
	}

	// ---- 学生：事务内创建，单条失败整体回滚；同班同名跳过 ----
	createdStudents := 0
	if len(req.Students) > 0 {
		err := s.db.Transaction(func(tx *gorm.DB) error {
			for _, stu := range req.Students {
				name := strings.TrimSpace(stu.Name)
				classID := uint(0)
				if stu.ClassID != nil && *stu.ClassID > 0 {
					classID = uint(*stu.ClassID)
				}
				var count int64
				if err := tx.Model(&models.Student{}).
					Where("class_id = ? AND name = ?", classID, name).
					Count(&count).Error; err != nil {
					return err
				}
				if count > 0 {
					skippedStudents = append(skippedStudents, SkippedEntry{Name: name, Reason: "该班级已有同名学生"})
					continue
				}
				gender := ""
				if stu.Gender != nil {
					gender = *stu.Gender
				}
				student := models.Student{ClassID: classID, Name: name, Gender: gender, Status: "active"}
				if err := tx.Create(&student).Error; err != nil {
					return err
				}
				if err := assignDefaultPetFor(tx, &student); err != nil {
					return err
				}
				createdStudents++
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	message := fmt.Sprintf("已导入 %d 名教师、%d 名学生", createdTeachers, createdStudents)
	if len(skippedTeachers) > 0 {
		message += fmt.Sprintf("，跳过已存在教师 %d 名", len(skippedTeachers))
	}
	if len(skippedStudents) > 0 {
		message += fmt.Sprintf("，跳过同名学生 %d 名", len(skippedStudents))
	}

	return &ContactsImportResult{
		Message: message,
		Data: ContactsImportData{
			CreatedTeachers: createdTeachers,
			CreatedStudents: createdStudents,
			SkippedTeachers: skippedTeachers,
			SkippedStudents: skippedStudents,
			TeacherAccounts: teacherAccounts,
		},
	}, nil
}

// validateContactsImport 等价 Laravel importWechatWorkUsers 的 Validator 规则
// （teachers.*.name required|string|max:50、teachers.*.mobile nullable|string|max:30、
// teachers.*.email nullable|string|max:120、students.*.name required|string|max:50、
// students.*.class_id required|integer、students.*.gender nullable|string|max:10）。
func validateContactsImport(req ContactsImportRequest) error {
	errs := map[string][]string{}
	for i, t := range req.Teachers {
		if strings.TrimSpace(t.Name) == "" {
			errs[fmt.Sprintf("teachers.%d.name", i)] = []string{"教师姓名必填"}
		} else if len([]rune(t.Name)) > 50 {
			errs[fmt.Sprintf("teachers.%d.name", i)] = []string{"教师姓名不能超过 50 字"}
		}
		if t.Mobile != nil && len([]rune(*t.Mobile)) > 30 {
			errs[fmt.Sprintf("teachers.%d.mobile", i)] = []string{"手机号不能超过 30 字"}
		}
		if t.Email != nil && len([]rune(*t.Email)) > 120 {
			errs[fmt.Sprintf("teachers.%d.email", i)] = []string{"邮箱不能超过 120 字"}
		}
	}
	for i, stu := range req.Students {
		if strings.TrimSpace(stu.Name) == "" {
			errs[fmt.Sprintf("students.%d.name", i)] = []string{"学生姓名必填"}
		} else if len([]rune(stu.Name)) > 50 {
			errs[fmt.Sprintf("students.%d.name", i)] = []string{"学生姓名不能超过 50 字"}
		}
		if stu.ClassID == nil {
			errs[fmt.Sprintf("students.%d.class_id", i)] = []string{"班级 ID 必填"}
		}
		if stu.Gender != nil && len([]rune(*stu.Gender)) > 10 {
			errs[fmt.Sprintf("students.%d.gender", i)] = []string{"性别不能超过 10 字"}
		}
	}
	if len(errs) > 0 {
		return NewValidationError(errs)
	}
	return nil
}

// stripPhoneSeparators 去掉手机号中的空白与连字符（同 Laravel `preg_replace('/[\s\-]/', ”, $phone)`）。
func stripPhoneSeparators(phone string) string {
	var b strings.Builder
	for _, r := range phone {
		if r == '-' || r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\v' || r == '\f' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
