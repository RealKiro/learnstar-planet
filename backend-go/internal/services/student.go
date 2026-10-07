// 教师端学生管理服务：批量导入（查重 + 跨班冲突防护）、增、改、删（均限定在管辖班级内）。
//
// 忠实移植自 Laravel App\Services\StudentService（教师端口径，带 classIds 作用域）：
// import / create / update / delete。
//
// 有意差异：
//  1. 删除语义：Laravel 的 Student 有 SoftDeletes（只打 deleted_at）；Go schema 的
//     `students.deleted_at` 是普通 *time.Time（不触发软删除），故与 Go 端既有
//     AdminDeleteStudent / BatchDeleteStudents 保持一致：**硬删除并级联清理**
//     宠物 / 积分 / 审计日志。
//  2. Laravel 的 list() 由本文件 List 实现（旧的 services.Teacher.Students 缺 search /
//     分页 / meta 且额外过滤了 status，已删除以免两份不一致的实现并存）。
//  3. 404 文案：Laravel 的 findOrFail 统一返回「资源不存在」；Go 端沿用既有更具体的
//     中文文案（如「学生不存在」），状态码一致。
package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// StudentService 教师端学生服务。
type StudentService struct {
	db *gorm.DB
}

// NewStudentService 创建教师端学生服务。
func NewStudentService(db *gorm.DB) *StudentService {
	return &StudentService{db: db}
}

// TeacherStudentImportRow 教师端导入的一行（同 Laravel students/*.name|class_name|gender|student_no）。
type TeacherStudentImportRow struct {
	Name      string
	ClassName string
	Gender    string
	StudentNo string
}

// TeacherStudentImportResult 教师端导入结果（message + imported_count + skipped 文案列表）。
// 说明：管理端批量导入另有一套 StudentImportResult（含 errors/warnings），二者语义不同。
type TeacherStudentImportResult struct {
	Message       string
	ImportedCount int
	Skipped       []string
}

// Import 批量导入学生：按班级名在 classIDs 范围内定位班级，同班同学号/同名跳过，
// 跨班同学号疑似转班跳过；班级名不在范围内或 name/class_name 为空的行静默忽略。
// 逐行处理、跳过项不中断整体导入（同 Laravel，未包事务）。
func (s *StudentService) Import(classIDs []uint, rows []TeacherStudentImportRow) (*TeacherStudentImportResult, error) {
	imported := 0
	skipped := []string{}

	if len(classIDs) > 0 {
		for _, data := range rows {
			if strings.TrimSpace(data.Name) == "" || strings.TrimSpace(data.ClassName) == "" {
				continue
			}
			var class models.ClassRoom
			err := s.db.Where("id IN ? AND name = ?", classIDs, data.ClassName).First(&class).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			if err != nil {
				return nil, err
			}

			name := strings.TrimSpace(data.Name)
			studentNo := strings.TrimSpace(data.StudentNo)

			// 查重：同班同学号（有学号时）或同班同名已存在 → 跳过。
			dupQuery := s.db.Model(&models.Student{}).Where("class_id = ?", class.ID)
			var dup int64
			if studentNo != "" {
				dupQuery = dupQuery.Where("student_no = ?", studentNo)
			} else {
				dupQuery = dupQuery.Where("name = ?", name)
			}
			if err := dupQuery.Count(&dup).Error; err != nil {
				return nil, err
			}
			if dup > 0 {
				label := name
				if studentNo != "" {
					label = fmt.Sprintf("%s（学号 %s）", name, studentNo)
				}
				skipped = append(skipped, label+"："+class.Name+" 已存在")
				continue
			}

			// 跨班冲突防护：同学号已在本校其他班级（active）→ 疑似转班，跳过并提示。
			if studentNo != "" {
				var cross []models.Student
				err := s.db.Where("student_no = ? AND class_id != ? AND status = ?", studentNo, class.ID, "active").
					Where("class_id IN (?)", s.db.Model(&models.ClassRoom{}).Select("id").Where("school_id = ?", class.SchoolID)).
					Find(&cross).Error
				if err != nil {
					return nil, err
				}
				if len(cross) > 0 {
					skipped = append(skipped, fmt.Sprintf(
						"%s（学号 %s）：已存在于 %s，如为转班请联系管理员使用批量转班",
						name, studentNo, classDisplayName(s.classNameMap(classIDs), cross[0].ClassID)))
					continue
				}
			}

			gender := data.Gender
			if gender == "" {
				gender = "未知"
			}
			student := models.Student{
				ClassID:   class.ID,
				Name:      name,
				Gender:    gender,
				StudentNo: studentNo,
				Status:    "active",
			}
			if err := s.db.Create(&student).Error; err != nil {
				return nil, err
			}
			imported++
		}
	}

	message := fmt.Sprintf("成功导入 %d 名学生", imported)
	if len(skipped) > 0 {
		message += fmt.Sprintf("，跳过 %d 条重复/冲突记录", len(skipped))
	}
	return &TeacherStudentImportResult{Message: message, ImportedCount: imported, Skipped: skipped}, nil
}

// classNameMap 返回班级 id → 名称（用于跳过文案里的班级名）。
func (s *StudentService) classNameMap(classIDs []uint) map[uint]string {
	out := map[uint]string{}
	if len(classIDs) == 0 {
		return out
	}
	var classes []models.ClassRoom
	if err := s.db.Where("id IN ?", classIDs).Find(&classes).Error; err != nil {
		return out
	}
	for _, c := range classes {
		out[c.ID] = c.Name
	}
	return out
}

// Create 在管辖班级内添加学生（返回 nil 表示目标班级不在管辖范围 → 控制器映射 403）。
// 性别归一化为 男/女/未知（同 Laravel StudentService::create）。
func (s *StudentService) Create(classIDs []uint, name string, classID uint, gender, studentNo string) (*models.Student, error) {
	var class models.ClassRoom
	err := s.db.Where("id IN ?", classIDs).First(&class, classID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
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
	if err := s.db.Create(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

// Update 更新管辖班级内的学生（仅 name / gender / student_no）；越权或不存在 → 404。
func (s *StudentService) Update(classIDs []uint, id uint, updates map[string]any) (*models.Student, error) {
	student, err := s.findScoped(classIDs, id)
	if err != nil {
		return nil, err
	}
	if len(updates) > 0 {
		if err := s.db.Model(&models.Student{}).Where("id = ?", student.ID).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.findScoped(classIDs, id)
}

// Delete 删除管辖班级内的学生，返回被删学生（供响应文案）。
//
// 与 Laravel 的差异：Laravel 是软删除（SoftDeletes）；Go 端与既有
// AdminDeleteStudent 一致——硬删除并级联清理宠物 / 审计日志 / 积分记录。
func (s *StudentService) Delete(classIDs []uint, id uint) (*models.Student, error) {
	student, err := s.findScoped(classIDs, id)
	if err != nil {
		return nil, err
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("student_id = ?", student.ID).Delete(&models.Pet{}).Error; err != nil {
			return err
		}
		if err := tx.Where("student_id = ?", student.ID).Delete(&models.ScoreLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("student_id = ?", student.ID).Delete(&models.Score{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Student{}, student.ID).Error
	})
	if err != nil {
		return nil, err
	}
	return student, nil
}

// findScoped 在管辖班级范围内取学生（越权/不存在 → 404）。
func (s *StudentService) findScoped(classIDs []uint, id uint) (*models.Student, error) {
	var student models.Student
	err := s.db.Where("id = ? AND class_id IN ?", id, classIDs).First(&student).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("学生不存在")
	}
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// TeacherStudentRow 学生行：学生原字段 + 追加宠物字段
// （同 Laravel `array_merge($s->toArray(), [pet_species, pet_level, pet_name])`）。
type TeacherStudentRow struct {
	models.Student
	PetSpecies string `json:"pet_species"`
	PetLevel   int    `json:"pet_level"`
	PetName    string `json:"pet_name"`
}

// TeacherStudentPage 教师端学生分页结果（meta 四项同 Laravel LengthAwarePaginator）。
type TeacherStudentPage struct {
	Data []TeacherStudentRow `json:"data"`
	Meta AdminPageMeta       `json:"meta"`
}

// List 教师所带班级学生分页列表（同 Laravel StudentService::list）。
//
//	whereIn('class_id', $classIds) + with('classRoom:id,name,grade') + with('pet')
//	可选 search 命中姓名或学号；orderBy('name')；paginate(50)
//
// 注意两点与旧 Go 实现的区别（都是为了对齐 Laravel）：
//   - **不过滤 `status`**：Laravel 没有 status 条件，停用/毕业学生同样列出；
//   - 分页固定 50（Laravel `paginate(50)` 硬编码，`per_page` 参数不生效）。
func (s *StudentService) List(classIDs []uint, search string, page, perPage int) (*TeacherStudentPage, error) {
	if perPage <= 0 {
		perPage = 50
	}
	if page <= 0 {
		page = 1
	}

	meta := AdminPageMeta{CurrentPage: page, LastPage: 1, PerPage: perPage}
	rows := []TeacherStudentRow{}
	if len(classIDs) == 0 {
		return &TeacherStudentPage{Data: rows, Meta: meta}, nil
	}

	// 每次重建查询，避免 Count 的 SELECT 污染后续 Find。
	base := func() *gorm.DB {
		tx := s.db.Model(&models.Student{}).Where("class_id IN ?", classIDs)
		if search != "" {
			like := "%" + search + "%"
			// 括号必需：GORM 不会为原始字符串条件补括号，缺括号时
			// `AND` 优先级高于 `OR`，学号命中会带出其他班级/其他学校的学生。
			tx = tx.Where("(name LIKE ? OR student_no LIKE ?)", like, like)
		}
		return tx
	}

	var total int64
	if err := base().Count(&total).Error; err != nil {
		return nil, err
	}

	var students []models.Student
	if err := base().Order("name ASC").
		Offset((page - 1) * perPage).Limit(perPage).
		Find(&students).Error; err != nil {
		return nil, err
	}

	petByStudent := map[uint]models.Pet{}
	if len(students) > 0 {
		ids := make([]uint, 0, len(students))
		for _, student := range students {
			ids = append(ids, student.ID)
		}
		var pets []models.Pet
		if err := s.db.Where("student_id IN ?", ids).Find(&pets).Error; err != nil {
			return nil, err
		}
		for _, pet := range pets {
			petByStudent[pet.StudentID] = pet
		}
	}

	for _, student := range students {
		row := TeacherStudentRow{Student: student}
		if pet, ok := petByStudent[student.ID]; ok {
			row.PetSpecies = pet.Species
			row.PetLevel = pet.Level
			row.PetName = pet.Name
		}
		rows = append(rows, row)
	}

	if total > 0 {
		meta.LastPage = int((total + int64(perPage) - 1) / int64(perPage))
	}
	meta.Total = total
	return &TeacherStudentPage{Data: rows, Meta: meta}, nil
}
