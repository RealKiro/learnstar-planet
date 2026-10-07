// 管理端（本批 5 条）：学生分页列表 / 单建学生 / 班级详情 / 教师班级分配 / 查看教师密码。
//
// 忠实移植自 Laravel App\Http\Controllers\Api\SchoolAdminController 的
// listStudents / createStudent / show / assignTeacherClasses / getTeacherPassword。
// 与 Laravel 的有意差异（DTO 字段、软删除口径、班级可见范围）逐条写在方法注释与
// README「本批有意差异」中。
package services

import (
	cryptorand "crypto/rand"
	"errors"
	"math/big"
	"math/rand/v2"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ============================================================
// 学生分页列表（GET /admin/students）
// ============================================================

// AdminStudentQuery 学生列表筛选条件（对应 Laravel listStudents 的 query 参数）。
type AdminStudentQuery struct {
	Search  string
	ClassID uint
	Grade   string
	Status  string // 默认 active；all 表示不过滤
	PerPage int    // 默认 50
	Page    int    // 默认 1
}

// AdminStudentRow 学生列表的一行（学生基础字段 + class_name/class_grade/pet_* 追加字段）。
//
// 与 Laravel 的差异：Laravel 直接 `$s->toArray()` 输出，故还带 `deleted_at` 与嵌套的
// `class_room` / `pet` 关联对象；Go 端只输出显式字段（含 5 个追加字段，供前端列表页使用），
// 空串的可空列（student_no / avatar_path）序列化为 null，与 Laravel 的 null 对齐。
type AdminStudentRow struct {
	ID         uint      `json:"id"`
	ClassID    uint      `json:"class_id"`
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
	StudentNo  *string   `json:"student_no"`
	AvatarPath *string   `json:"avatar_path"`
	TotalScore int       `json:"total_score"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ClassName  *string   `json:"class_name"`
	ClassGrade *string   `json:"class_grade"`
	PetSpecies string    `json:"pet_species"`
	PetLevel   int       `json:"pet_level"`
	PetName    string    `json:"pet_name"`
}

// AdminStudentPage 分页结果（meta 字段同 Laravel LengthAwarePaginator 的四项）。
type AdminStudentPage struct {
	Data []AdminStudentRow `json:"data"`
	Meta AdminPageMeta     `json:"meta"`
}

// AdminPageMeta 分页元信息。
type AdminPageMeta struct {
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
}

// ListStudentsPaged 本校学生分页列表（筛选 search / class_id / grade / status）。
// 排序 `id desc`；每行追加班级名/年级与宠物物种/等级/名称（同 Laravel listStudents）。
func (a *Admin) ListStudentsPaged(schoolID uint, q AdminStudentQuery) (*AdminStudentPage, error) {
	if q.PerPage <= 0 {
		q.PerPage = 50
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Status == "" {
		q.Status = "active"
	}

	// 每次重建查询，避免 Count 的 SELECT 污染后续 Find（GORM 用法坑）。

	base := func() *gorm.DB {
		tx := a.db.Model(&models.Student{}).
			Where("class_id IN (?)", a.db.Model(&models.ClassRoom{}).Select("id").Where("school_id = ?", schoolID))
		if q.Search != "" {
			like := "%" + q.Search + "%"
			// 必须显式括号：GORM 不会为原始字符串条件补括号，
			// 否则 `AND` 优先级高于 `OR`，学号命中会把其他班级/其他学校的学生一并带出。
			tx = tx.Where("(name LIKE ? OR student_no LIKE ?)", like, like)
		}
		if q.ClassID != 0 {
			tx = tx.Where("class_id = ?", q.ClassID)
		}
		if q.Grade != "" {
			tx = tx.Where("class_id IN (?)", a.db.Model(&models.ClassRoom{}).Select("id").
				Where("school_id = ? AND grade = ?", schoolID, q.Grade))
		}
		if q.Status != "all" {
			tx = tx.Where("status = ?", q.Status)
		}
		return tx
	}

	var total int64
	if err := base().Count(&total).Error; err != nil {
		return nil, err
	}

	var students []models.Student
	if err := base().Order("id DESC").
		Offset((q.Page - 1) * q.PerPage).Limit(q.PerPage).
		Find(&students).Error; err != nil {
		return nil, err
	}

	rows, err := a.buildStudentRows(students)
	if err != nil {
		return nil, err
	}

	lastPage := int((total + int64(q.PerPage) - 1) / int64(q.PerPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return &AdminStudentPage{
		Data: rows,
		Meta: AdminPageMeta{CurrentPage: q.Page, LastPage: lastPage, Total: total, PerPage: q.PerPage},
	}, nil
}

// buildStudentRows 批量补齐班级与宠物字段（避免 N+1）。
func (a *Admin) buildStudentRows(students []models.Student) ([]AdminStudentRow, error) {
	rows := make([]AdminStudentRow, 0, len(students))
	if len(students) == 0 {
		return rows, nil
	}

	classIDs := make([]uint, 0, len(students))
	studentIDs := make([]uint, 0, len(students))
	for _, s := range students {
		classIDs = append(classIDs, s.ClassID)
		studentIDs = append(studentIDs, s.ID)
	}

	classNames := map[uint]string{}
	classGrades := map[uint]string{}
	var classes []models.ClassRoom
	if err := a.db.Where("id IN ?", classIDs).Find(&classes).Error; err != nil {
		return nil, err
	}
	for _, c := range classes {
		classNames[c.ID] = c.Name
		classGrades[c.ID] = c.Grade
	}

	pets := map[uint]models.Pet{}
	var petList []models.Pet
	if err := a.db.Where("student_id IN ?", studentIDs).Find(&petList).Error; err != nil {
		return nil, err
	}
	for _, p := range petList {
		pets[p.StudentID] = p
	}

	for _, s := range students {
		row := AdminStudentRow{
			ID:         s.ID,
			ClassID:    s.ClassID,
			Name:       s.Name,
			Gender:     s.Gender,
			StudentNo:  nullableString(s.StudentNo),
			AvatarPath: nullableString(s.AvatarPath),
			TotalScore: s.TotalScore,
			Status:     s.Status,
			CreatedAt:  s.CreatedAt,
			UpdatedAt:  s.UpdatedAt,
		}
		if name, ok := classNames[s.ClassID]; ok {
			row.ClassName = &name
		}
		if grade, ok := classGrades[s.ClassID]; ok {
			row.ClassGrade = &grade
		}
		if pet, ok := pets[s.ID]; ok {
			row.PetSpecies = pet.Species
			row.PetLevel = pet.Level
			row.PetName = pet.Name
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// nullableString 空串 → nil（对齐 Laravel 可空列在 JSON 中的 null）。
func nullableString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// ============================================================
// 单建学生（POST /admin/students）
// ============================================================

// CreateStudentInClass 本校班级内创建学生并自动分配默认宠物。
// 同 Laravel createStudent：`exists:class_rooms,id` 不通过（班级整体不存在）→ 422「参数错误」+ errors；
// 班级存在但不属于本校 → 404「班级不存在」（findOrFail）。
func (a *Admin) CreateStudentInClass(schoolID, classID uint, name, gender, studentNo string) (*models.Student, error) {
	var anyClass models.ClassRoom
	err := a.db.First(&anyClass, classID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, NewValidationError(map[string][]string{"class_id": {"所选班级不存在"}})
	}
	if err != nil {
		return nil, err
	}

	class, err := a.classInSchool(schoolID, classID)
	if err != nil {
		return nil, err
	}

	student := models.Student{
		ClassID:   class.ID,
		Name:      name,
		Gender:    normalizeGender(gender),
		StudentNo: studentNo,
		Status:    "active",
	}
	if err := a.db.Create(&student).Error; err != nil {
		return nil, err
	}
	// 自动分配萌宠（同 Laravel assignDefaultPet）。
	if err := assignDefaultPetFor(a.db, &student); err != nil {
		return nil, err
	}
	return &student, nil
}

// ============================================================
// 班级详情（GET /admin/classes/:id）
// ============================================================

// AdminClassDetail 班级详情视图：班级字段 + teacher 与 students 关联。
type AdminClassDetail struct {
	models.ClassRoom
	// Settings 班级设置对象（Laravel `class_rooms.settings` 是 array cast；前端班级页读 `settings.pet_series`）。
	// 嵌入的 ClassRoom.Settings 带 `json:"-"`，故此处的 `settings` 键不冲突。
	Settings map[string]any   `json:"settings"`
	Teacher  *models.User     `json:"teacher"`
	Students []models.Student `json:"students"`
}

// ClassDetail 取本校班级详情（含 teacher / students）；跨校或不存在 → 404「班级不存在」。
func (a *Admin) ClassDetail(schoolID, classID uint) (*AdminClassDetail, error) {
	class, err := a.classInSchool(schoolID, classID)
	if err != nil {
		return nil, err
	}

	settings := class.SettingsMap()
	if settings == nil {
		settings = map[string]any{}
	}
	detail := &AdminClassDetail{ClassRoom: *class, Settings: settings, Students: []models.Student{}}
	if class.TeacherID != nil {
		var teacher models.User
		err := a.db.First(&teacher, *class.TeacherID).Error
		if err == nil {
			detail.Teacher = &teacher
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if err := a.db.Where("class_id = ?", class.ID).Order("id ASC").Find(&detail.Students).Error; err != nil {
		return nil, err
	}
	return detail, nil
}

// ============================================================
// 教师班级分配（PUT /admin/teachers/:id/classes）
// ============================================================

// TeacherClassAssignment 一次班级分配（对应 Laravel assignments.*）。
type TeacherClassAssignment struct {
	ClassID uint
	Role    string
	Subject *string
}

// TeacherClassAssignmentResult 已同步的分配（字段同 Laravel synced 元素）。
type TeacherClassAssignmentResult struct {
	ClassID   uint    `json:"class_id"`
	ClassName string  `json:"class_name"`
	Role      string  `json:"role"`
	Subject   *string `json:"subject"`
}

// AssignTeacherClasses 以 **replace 语义** 设置某教师的班级分配：
// 逐条 upsert（class_room_id + user_id），role=head_teacher 时同步 class_rooms.teacher_id；
// 随后删除本次未提交的旧分配，被删除且原 role 为 head_teacher 的班级把 teacher_id 置空。
// 教师须在本校且 role=teacher，否则 404「教师不存在」；班级须在本校，否则 404「班级不存在」。
func (a *Admin) AssignTeacherClasses(schoolID, teacherID uint, assignments []TeacherClassAssignment) ([]TeacherClassAssignmentResult, error) {
	teacher, err := a.teacherInSchool(schoolID, teacherID)
	if err != nil {
		return nil, err
	}

	synced := make([]TeacherClassAssignmentResult, 0, len(assignments))
	err = a.db.Transaction(func(tx *gorm.DB) error {
		submitted := make([]uint, 0, len(assignments))
		for _, item := range assignments {
			var class models.ClassRoom
			if err := tx.Where("id = ? AND school_id = ?", item.ClassID, schoolID).First(&class).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrNotFound("班级不存在")
				}
				return err
			}
			if !validClassTeacherRoles[item.Role] {
				return ErrUnprocessable("参数错误")
			}

			subject := ""
			if item.Subject != nil {
				subject = *item.Subject
			}

			var existing models.ClassRoomTeacher
			findErr := tx.Where("class_room_id = ? AND user_id = ?", class.ID, teacher.ID).
				First(&existing).Error
			switch {
			case findErr == nil:
				if err := tx.Model(&models.ClassRoomTeacher{}).Where("id = ?", existing.ID).
					Updates(map[string]any{"role": item.Role, "subject": subject}).Error; err != nil {
					return err
				}
			case errors.Is(findErr, gorm.ErrRecordNotFound):
				if err := tx.Create(&models.ClassRoomTeacher{
					ClassRoomID: class.ID,
					UserID:      teacher.ID,
					Role:        item.Role,
					Subject:     subject,
				}).Error; err != nil {
					return err
				}
			default:
				return findErr
			}

			if item.Role == "head_teacher" {
				if err := tx.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
					Update("teacher_id", teacher.ID).Error; err != nil {
					return err
				}
			}

			submitted = append(submitted, class.ID)
			result := TeacherClassAssignmentResult{ClassID: class.ID, ClassName: class.Name, Role: item.Role}
			if item.Subject != nil {
				subjectCopy := *item.Subject
				result.Subject = &subjectCopy
			}
			synced = append(synced, result)
		}

		// 删除「本次未提交」的旧分配（submitted 为空时删除全部，同 Laravel whereNotIn(..., [0])）。
		oldQuery := tx.Where("user_id = ?", teacher.ID)
		if len(submitted) > 0 {
			oldQuery = oldQuery.Where("class_room_id NOT IN ?", submitted)
		}
		var removed []models.ClassRoomTeacher
		if err := oldQuery.Find(&removed).Error; err != nil {
			return err
		}
		for _, old := range removed {
			if old.Role == "head_teacher" {
				if err := tx.Model(&models.ClassRoom{}).Where("id = ?", old.ClassRoomID).
					Update("teacher_id", nil).Error; err != nil {
					return err
				}
			}
			if err := tx.Delete(&models.ClassRoomTeacher{}, old.ID).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return synced, nil
}

// ============================================================
// 查看教师密码（GET /admin/teachers/:id/password）
// ============================================================

// randomPassword 生成 n 位随机密码（字母 + 数字，同 Laravel str()->random(8)）。
func randomPassword(n int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	out := make([]byte, n)
	for i := range out {
		idx, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		out[i] = alphabet[idx.Int64()]
	}
	return string(out), nil
}

// TeacherPassword 查看教师明文密码；无明文记录时自动生成 8 位随机密码并写回
// （同 Laravel getTeacherPassword：这会重置该教师密码，并把 password_changed 置回 false）。
// 第二个返回值为 true 表示本次自动生成了新密码。
func (a *Admin) TeacherPassword(schoolID, teacherID uint) (string, bool, error) {
	teacher, err := a.teacherInSchool(schoolID, teacherID)
	if err != nil {
		return "", false, err
	}
	if teacher.PlainPassword != "" {
		return teacher.PlainPassword, false, nil
	}

	newPassword, err := randomPassword(8)
	if err != nil {
		return "", false, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", false, err
	}
	if err := a.db.Model(&models.User{}).Where("id = ?", teacher.ID).Updates(map[string]any{
		"password_hash":    string(hash),
		"plain_password":   newPassword,
		"password_changed": false,
	}).Error; err != nil {
		return "", false, err
	}
	return newPassword, true, nil
}

// ============================================================
// 内部工具
// ============================================================

// assignDefaultPetFor 为新生自动分配默认宠物：myth 系列随机物种、名称「{学生名}的萌宠」、
// 等级 1 / 经验 0 / 心情 80（同 Laravel assignDefaultPet）。
// 包级函数，供 Admin 与 AdminOps（批量导入）共用。
func assignDefaultPetFor(db *gorm.DB, student *models.Student) error {
	pool := models.SpeciesPoolForSeries("myth")
	species := "zhulong"
	if len(pool) > 0 {
		species = pool[rand.IntN(len(pool))]
	}
	return db.Create(&models.Pet{
		StudentID:  student.ID,
		ClassID:    student.ClassID,
		Name:       student.Name + "的萌宠",
		Species:    species,
		Level:      1,
		Experience: 0,
		Mood:       80,
	}).Error
}

// ClassExistsAnywhere 班级是否存在（不做学校过滤；对应 Laravel `exists:class_rooms,id` 规则）。
func (a *Admin) ClassExistsAnywhere(classID uint) (bool, error) {
	var count int64
	if err := a.db.Model(&models.ClassRoom{}).Where("id = ?", classID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
