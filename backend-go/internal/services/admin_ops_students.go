// 管理端运维服务（三）：学生批量导入/删除/转班、班级批量创建、班级教师分配。
//
// 忠实移植自 Laravel App\Http\Controllers\Api\SchoolAdminController：
// importStudents / batchDeleteStudents / batchMoveStudents / batchCreateClasses /
// assignClassTeacher / removeClassTeacher，以及 App\Services\StudentService::import 的
// 查重与跨班冲突防护口径。
package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// 导入结果上限说明：与 Laravel 一致，逐行处理、跳过项不中断整体导入（整批不做事务）。
// 学生姓名会先 TrimSpace 再入库（Laravel 直接用原值，属轻微加固差异）。

// StudentImportRow 学生导入的一行（同 Laravel students/*.name/class_name/gender/student_no）。
type StudentImportRow struct {
	Name      string
	ClassName string
	Gender    string
	StudentNo string
}

// StudentImportCreated 导入成功的一行（同 Laravel created 元素结构）。
type StudentImportCreated struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	ClassName string `json:"class_name"`
	Gender    string `json:"gender"`
}

// StudentImportResult 导入结果：成功行 + 行内错误 + 非阻塞提醒。
// Laravel：errors 非空时整体返回 422，且**已成功的行仍然落库**（本实现保持一致）。
type StudentImportResult struct {
	Created  []StudentImportCreated
	Errors   []string
	Warnings []string
}

// normalizeGender 性别归一化（同 Laravel：男生/男→男，女生/女→女，其余→未知）。
func normalizeGender(g string) string {
	switch strings.TrimSpace(g) {
	case "男", "男生":
		return "男"
	case "女", "女生":
		return "女"
	default:
		return "未知"
	}
}

// assignDefaultPet 为新生自动分配默认宠物（同 Laravel assignDefaultPet：
// myth 系列随机物种、名称「{学生名}的萌宠」、等级 1 / 经验 0 / 心情 80）。
// 与 Admin 服务共用包级函数 assignDefaultPetFor（避免两份随机物种池逻辑）。
func (o *AdminOps) assignDefaultPet(student *models.Student) error {
	return assignDefaultPetFor(o.db, student)
}

// schoolClassSubquery 本校班级 ID 子查询。
// 注意：GORM 子查询必须写成 `IN (?)`，`IN ?` 无效。
func (o *AdminOps) schoolClassSubquery(schoolID uint) *gorm.DB {
	return o.db.Model(&models.ClassRoom{}).Select("id").Where("school_id = ?", schoolID)
}

// ImportStudents 批量导入学生（按模板班级名称匹配）。
//
// 口径逐条对齐 Laravel importStudents：
//   - 班级名不在本校「active」班级列表中 → 该行报错「班级「X」不存在」；
//   - 同班同学号已存在 → 报错并跳过；
//   - 同学号已在本校其他班级（active）→ 报错并提示走批量转班；
//   - 本校其他班级有同名学生 → 追加**非阻塞** warning；
//   - 成功行创建学生并自动分配默认宠物。
//
// 行号按入参数组下标 +1 计算（同 Laravel `$idx + 1`）。
func (o *AdminOps) ImportStudents(schoolID uint, rows []StudentImportRow) (*StudentImportResult, error) {
	if len(rows) == 0 {
		return nil, ErrUnprocessable("参数错误：students 至少 1 条")
	}

	var classes []models.ClassRoom
	if err := o.db.Where("school_id = ? AND status = ?", schoolID, "active").
		Find(&classes).Error; err != nil {
		return nil, err
	}
	classByName := make(map[string]uint, len(classes))
	classNames := make(map[uint]string, len(classes))
	for _, c := range classes {
		// 同名班级后者覆盖，同 Laravel `pluck('id', 'name')`。
		classByName[c.Name] = c.ID
		classNames[c.ID] = c.Name
	}

	result := &StudentImportResult{
		Created:  []StudentImportCreated{},
		Errors:   []string{},
		Warnings: []string{},
	}

	for idx, row := range rows {
		line := idx + 1
		className := strings.TrimSpace(row.ClassName)
		classID, ok := classByName[className]
		if !ok {
			result.Errors = append(result.Errors,
				fmt.Sprintf("第 %d 行：班级「%s」不存在", line, className))
			continue
		}

		name := strings.TrimSpace(row.Name)
		gender := normalizeGender(row.Gender)
		studentNo := strings.TrimSpace(row.StudentNo)

		if studentNo != "" {
			// 同班同学号查重
			var dup int64
			if err := o.db.Model(&models.Student{}).
				Where("class_id = ? AND student_no = ?", classID, studentNo).
				Count(&dup).Error; err != nil {
				return nil, err
			}
			if dup > 0 {
				result.Errors = append(result.Errors,
					fmt.Sprintf("第 %d 行：学号「%s」在班级「%s」已存在，已跳过", line, studentNo, className))
				continue
			}

			// 跨班冲突防护：同学号已在本校其他班级（active）→ 疑似转班
			var cross models.Student
			err := o.db.Where("student_no = ? AND class_id != ? AND status = ?", studentNo, classID, "active").
				Where("class_id IN (?)", o.schoolClassSubquery(schoolID)).
				First(&cross).Error
			if err == nil {
				result.Errors = append(result.Errors,
					fmt.Sprintf("第 %d 行：学号「%s」已存在于「%s」。若为同一学生转班，请使用学生管理的批量转班功能，不要在两个班重复创建",
						line, studentNo, classDisplayName(classNames, cross.ClassID)))
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}

		// 同名提醒（不阻塞）
		var sameName models.Student
		err := o.db.Where("name = ? AND class_id != ? AND status = ?", name, classID, "active").
			Where("class_id IN (?)", o.schoolClassSubquery(schoolID)).
			First(&sameName).Error
		if err == nil {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("第 %d 行：新导入的「%s」与「%s」现有学生同名。若为同一学生（转班），请删除本条并用批量转班；若为同名不同人，可忽略本提醒",
					line, name, classDisplayName(classNames, sameName.ClassID)))
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		student := models.Student{
			ClassID:   classID,
			Name:      name,
			Gender:    gender,
			StudentNo: studentNo,
			Status:    "active",
		}
		if err := o.db.Create(&student).Error; err != nil {
			return nil, err
		}
		if err := o.assignDefaultPet(&student); err != nil {
			return nil, err
		}
		result.Created = append(result.Created, StudentImportCreated{
			ID:        student.ID,
			Name:      student.Name,
			ClassName: className,
			Gender:    student.Gender,
		})
	}

	return result, nil
}

// classDisplayName 取班级显示名（找不到时回退「其他班级」，同 Laravel `?? '其他班级'`）。
func classDisplayName(names map[uint]string, classID uint) string {
	if name, ok := names[classID]; ok && name != "" {
		return name
	}
	return "其他班级"
}

// ParseStudentCSV 解析学生导入 CSV（表头别名 姓名/班级/性别/学号[/手机号] 或
// name/class_name/gender/student_no[/phone]）。
//
// 与 Laravel 的差异：Laravel 的 students/import 只接受 JSON 数组，不接受文件上传；
// Go 端在 JSON 契约之外**额外**支持 CSV 文件上传（任务要求用标准库解析）。
// 首行没有任何可识别表头时，按固定列序（姓名, 班级, 性别, 学号, 手机号）解析。
func (o *AdminOps) ParseStudentCSV(r io.Reader) ([]StudentImportRow, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, ErrUnprocessable("CSV 解析失败：" + err.Error())
	}
	if len(records) == 0 {
		return []StudentImportRow{}, nil
	}

	// 去掉 UTF-8 BOM
	if len(records[0]) > 0 {
		records[0][0] = strings.TrimPrefix(records[0][0], "\ufeff")
	}

	header := make([]string, 0, len(records[0]))
	hasHeader := false
	for _, cell := range records[0] {
		key := strings.TrimSpace(cell)
		header = append(header, key)
		switch key {
		case "姓名", "班级", "性别", "学号", "手机号",
			"name", "class_name", "gender", "student_no", "phone":
			hasHeader = true
		}
	}

	rows := make([]StudentImportRow, 0, len(records)-1)
	if !hasHeader {
		// 无表头：按固定列序解析全部行。
		for _, rec := range records {
			rows = append(rows, studentRowFromColumns([]string{cellAt(rec, 0), cellAt(rec, 1), cellAt(rec, 2), cellAt(rec, 3), cellAt(rec, 4)}))
		}
		return rows, nil
	}

	for _, rec := range records[1:] {
		cell := map[string]string{}
		for i, key := range header {
			if i < len(rec) {
				cell[key] = rec[i]
			}
		}
		name := firstNonEmptyValue(cell, "姓名", "name")
		className := firstNonEmptyValue(cell, "班级", "class_name")
		if strings.TrimSpace(name) == "" && strings.TrimSpace(className) == "" {
			continue
		}
		rows = append(rows, StudentImportRow{
			Name:      strings.TrimSpace(name),
			ClassName: strings.TrimSpace(className),
			Gender:    strings.TrimSpace(firstNonEmptyValue(cell, "性别", "gender")),
			StudentNo: strings.TrimSpace(firstNonEmptyValue(cell, "学号", "student_no")),
		})
	}
	return rows, nil
}

// studentRowFromColumns 按固定列序构造一行（姓名, 班级, 性别, 学号[, 手机号]）。
func studentRowFromColumns(cols []string) StudentImportRow {
	return StudentImportRow{
		Name:      strings.TrimSpace(cellAt(cols, 0)),
		ClassName: strings.TrimSpace(cellAt(cols, 1)),
		Gender:    strings.TrimSpace(cellAt(cols, 2)),
		StudentNo: strings.TrimSpace(cellAt(cols, 3)),
	}
}

// cellAt 安全取列。
func cellAt(cols []string, i int) string {
	if i < len(cols) {
		return cols[i]
	}
	return ""
}

// BatchDeleteStudents 批量删除学生，返回实际删除条数。
//
// 与 Laravel batchDeleteStudents 的差异：Laravel 的 Student 有 SoftDeletes，删除只是打
// `deleted_at`；Go schema 的 `students.deleted_at` 是普通 `*time.Time`（GORM 不会触发软删除），
// 一致于 Go 端既有的单条删除语义：**硬删除并级联清理**宠物/积分/审计日志。
func (o *AdminOps) BatchDeleteStudents(schoolID uint, ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, ErrUnprocessable("参数错误")
	}
	var students []models.Student
	if err := o.db.Where("id IN ?", ids).
		Where("class_id IN (?)", o.schoolClassSubquery(schoolID)).
		Find(&students).Error; err != nil {
		return 0, err
	}
	if len(students) == 0 {
		return 0, nil
	}
	studentIDs := make([]uint, 0, len(students))
	for _, s := range students {
		studentIDs = append(studentIDs, s.ID)
	}

	err := o.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("student_id IN ?", studentIDs).Delete(&models.Pet{}).Error; err != nil {
			return err
		}
		if err := tx.Where("student_id IN ?", studentIDs).Delete(&models.ScoreLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("student_id IN ?", studentIDs).Delete(&models.Score{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", studentIDs).Delete(&models.Student{}).Error
	})
	if err != nil {
		return 0, err
	}
	return int64(len(studentIDs)), nil
}

// BatchMoveStudents 批量把学生转入目标班级，返回 (条数, 目标班级名)。
// 目标班级必须属于本校（否则 404）；只移动本校班级下的学生（跨校 ID 静默忽略，同 Laravel）。
//
// 有意差异：Laravel 只改 `students.class_id`，不动宠物；Go 端同样**不迁移宠物**（保持逐字一致）。
func (o *AdminOps) BatchMoveStudents(schoolID uint, ids []uint, targetClassID uint) (int64, string, error) {
	if len(ids) == 0 || targetClassID == 0 {
		return 0, "", ErrUnprocessable("参数错误")
	}

	var target models.ClassRoom
	if err := o.db.Where("id = ? AND school_id = ?", targetClassID, schoolID).
		First(&target).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, "", ErrNotFound("班级不存在")
		}
		return 0, "", err
	}

	res := o.db.Model(&models.Student{}).
		Where("id IN ?", ids).
		Where("class_id IN (?)", o.schoolClassSubquery(schoolID)).
		Update("class_id", target.ID)
	if res.Error != nil {
		return 0, "", res.Error
	}
	return res.RowsAffected, target.Name, nil
}

// ClassCreated 批量创建的班级（同 Laravel batchCreateClasses 的 data 元素）。
type ClassCreated struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// BatchCreateClasses 批量创建班级：按「年级（N）班」命名，序号从该年级现有班级数 +1 起。
//
// 逐条对齐 Laravel batchCreateClasses：同名且 active 的班级已存在则跳过；序号按**该年级
// 全部班级数**（含非 active）计算；年级必填、count 限 1..20。
func (o *AdminOps) BatchCreateClasses(schoolID uint, grade string, count int, year string) ([]ClassCreated, error) {
	grade = strings.TrimSpace(grade)
	if grade == "" {
		return nil, ErrUnprocessable("参数错误：年级必填")
	}
	if count < 1 || count > 20 {
		return nil, ErrUnprocessable("参数错误：数量需在 1-20 之间")
	}

	var existingCount int64
	if err := o.db.Model(&models.ClassRoom{}).
		Where("school_id = ? AND grade = ?", schoolID, grade).
		Count(&existingCount).Error; err != nil {
		return nil, err
	}

	created := []ClassCreated{}
	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("%s（%d）班", grade, int(existingCount)+i)

		var dup int64
		if err := o.db.Model(&models.ClassRoom{}).
			Where("school_id = ? AND name = ? AND status = ?", schoolID, name, "active").
			Count(&dup).Error; err != nil {
			return nil, err
		}
		if dup > 0 {
			continue
		}

		class := models.ClassRoom{
			SchoolID: schoolID,
			Name:     name,
			Grade:    grade,
			Year:     year,
			Status:   "active",
		}
		if err := o.db.Create(&class).Error; err != nil {
			return nil, err
		}
		created = append(created, ClassCreated{ID: class.ID, Name: class.Name})
	}
	return created, nil
}

// AssignClassTeacher 分配班级教师（同 Laravel assignClassTeacher）。
//
// role 默认 subject_teacher；role = head_teacher 时把 class_rooms.teacher_id 指向该教师。
// 无论角色为何都会 upsert `class_room_teachers`（对齐 Laravel 的 updateOrCreate）。
func (o *AdminOps) AssignClassTeacher(schoolID, classID, teacherID uint, role string) (*models.ClassRoom, error) {
	if teacherID == 0 {
		return nil, ErrUnprocessable("参数错误")
	}
	if role == "" {
		role = "subject_teacher"
	}
	if !validClassTeacherRoles[role] {
		return nil, ErrUnprocessable("参数错误")
	}

	var class models.ClassRoom
	if err := o.db.Where("id = ? AND school_id = ?", classID, schoolID).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound("班级不存在")
		}
		return nil, err
	}

	// 教师必须属于本校且 role=teacher（同 Laravel findOrFail → 404）。
	var teacher models.User
	if err := o.db.Where("id = ? AND school_id = ? AND role = ?", teacherID, schoolID, "teacher").
		First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound("教师不存在")
		}
		return nil, err
	}

	if role == "head_teacher" {
		if err := o.db.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
			Update("teacher_id", teacher.ID).Error; err != nil {
			return nil, err
		}
	}
	if err := o.upsertClassRoomTeacher(class.ID, teacher.ID, role, ""); err != nil {
		return nil, err
	}

	var fresh models.ClassRoom
	if err := o.db.First(&fresh, class.ID).Error; err != nil {
		return nil, err
	}
	return &fresh, nil
}

// RemoveClassTeacher 移除班级教师关联（同 Laravel removeClassTeacher）：
// 删除 `class_room_teachers` 关联；若该教师正是 class_rooms.teacher_id 则一并置空。
func (o *AdminOps) RemoveClassTeacher(schoolID, classID, teacherID uint) error {
	if teacherID == 0 {
		return ErrUnprocessable("参数错误")
	}

	var class models.ClassRoom
	if err := o.db.Where("id = ? AND school_id = ?", classID, schoolID).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound("班级不存在")
		}
		return err
	}

	if err := o.db.Where("class_room_id = ? AND user_id = ?", class.ID, teacherID).
		Delete(&models.ClassRoomTeacher{}).Error; err != nil {
		return err
	}
	if class.TeacherID != nil && *class.TeacherID == teacherID {
		if err := o.db.Model(&models.ClassRoom{}).Where("id = ?", class.ID).
			Update("teacher_id", nil).Error; err != nil {
			return err
		}
	}
	return nil
}
