// 学生宠物图鉴收藏（pet_collections）：读取接口 + 换宠时的归档/恢复写入。
//
// 忠实移植自 Laravel App\Services\PetService::collection / switchPet 与
// App\Http\Controllers\Api\DisplayController::classroomSwitchPet、App\Models\PetCollection：
//   - 读取（GET /teacher/pets/:studentId/collection）：当前激活宠物补录进图鉴（firstOrCreate）；
//   - 教师端换宠（老师替学生换）：归档旧物种进度 → 恢复目标物种进度（切回时全保留）；
//   - 教室端换宠（学生自选）：归档旧物种进度 → 只换物种，等级/经验/心情保留。
//
// 两条路径的差异是 Laravel 的原生行为（见 PetService.php 第 200-226 行与
// DisplayController.php 第 1226-1238 行），不是移植差异。
package services

import (
	"errors"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// PetCollectionEntry 图鉴中的一条（字段与顺序逐字同 Laravel collection 的元素）。
type PetCollectionEntry struct {
	Species    string `json:"species"`
	Level      int    `json:"level"`
	Experience int    `json:"experience"`
	Mood       int    `json:"mood"`
	IsActive   bool   `json:"is_active"`
}

// PetCollectionView 学生图鉴视图（含已解锁槽位与班级当前系列）。
type PetCollectionView struct {
	StudentID     uint                 `json:"student_id"`
	StudentName   string               `json:"student_name"`
	TotalScore    int                  `json:"total_score"`
	UnlockSlots   int                  `json:"unlock_slots"`
	ClassSeries   *string              `json:"class_series"`
	ActiveSpecies *string              `json:"active_species"`
	Collection    []PetCollectionEntry `json:"collection"`
}

// Collection 学生宠物图鉴：先把当前激活宠物补进图鉴（firstOrCreate 语义），
// 再按 species 升序返回全部收藏。学生不在管辖班级 → 404（同 Laravel findOrFail）。
func (p *PetService) Collection(u *models.User, studentID uint) (*PetCollectionView, error) {
	student, err := p.scope.StudentInScope(u, studentID)
	if err != nil {
		return nil, err
	}

	activePet, err := p.activePetOrNil(student.ID)
	if err != nil {
		return nil, err
	}

	// 确保当前激活宠物在收藏中（Laravel PetCollection::firstOrCreate）。
	if activePet != nil && activePet.Species != "" {
		var count int64
		if err := p.db.Model(&models.PetCollection{}).
			Where("student_id = ? AND species = ?", student.ID, activePet.Species).
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count == 0 {
			entry := models.PetCollection{
				StudentID:  student.ID,
				Species:    activePet.Species,
				Level:      activePet.Level,
				Experience: activePet.Experience,
				Mood:       activePet.Mood,
				IsActive:   true,
			}
			if err := p.db.Create(&entry).Error; err != nil {
				return nil, err
			}
		}
	}

	var collections []models.PetCollection
	if err := p.db.Where("student_id = ?", student.ID).Order("species ASC").
		Find(&collections).Error; err != nil {
		return nil, err
	}

	view := &PetCollectionView{
		StudentID:   student.ID,
		StudentName: student.Name,
		TotalScore:  student.TotalScore,
		UnlockSlots: models.UnlockSlotsForScore(student.TotalScore),
		Collection:  make([]PetCollectionEntry, 0, len(collections)),
	}
	if activePet != nil {
		species := activePet.Species
		view.ActiveSpecies = &species
	}
	var class models.ClassRoom
	if err := p.db.First(&class, student.ClassID).Error; err == nil {
		if series, ok := class.SettingsMap()["pet_series"].(string); ok && series != "" {
			view.ClassSeries = &series
		}
	}
	for _, c := range collections {
		view.Collection = append(view.Collection, PetCollectionEntry{
			Species:    c.Species,
			Level:      c.Level,
			Experience: c.Experience,
			Mood:       c.Mood,
			IsActive:   c.IsActive,
		})
	}
	return view, nil
}

// activePetOrNil 取学生宠物；无宠物返回 (nil, nil)。
func (p *PetService) activePetOrNil(studentID uint) (*models.Pet, error) {
	var pet models.Pet
	if err := p.db.Where("student_id = ?", studentID).First(&pet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pet, nil
}

// upsertPetCollection 把某物种的进度写入图鉴，等价 Laravel `PetCollection::updateOrCreate`：
// 按 `(student_id, species)` 命中则覆盖 level/experience/mood/is_active，未命中则新建。
//
// 注意：`pet_collections.is_active` 带 `gorm:"default:true"`，GORM 的 Create 会跳过 `false`
// 零值而由数据库填 true（本仓既有坑位，见 README「惩罚规则标记」），故归档（false）时必须显式回写。
func upsertPetCollection(tx *gorm.DB, studentID uint, species string, level, experience, mood int, isActive bool) error {
	if species == "" {
		return nil
	}

	var row models.PetCollection
	err := tx.Where("student_id = ? AND species = ?", studentID, species).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		created := models.PetCollection{
			StudentID:  studentID,
			Species:    species,
			Level:      level,
			Experience: experience,
			Mood:       mood,
			IsActive:   isActive,
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		if !isActive {
			return tx.Model(&models.PetCollection{}).Where("id = ?", created.ID).
				Update("is_active", false).Error
		}
		return nil
	}
	if err != nil {
		return err
	}

	// 用 map 更新：GORM 的 struct Updates 同样会跳过 false 零值。
	return tx.Model(&models.PetCollection{}).Where("id = ?", row.ID).Updates(map[string]any{
		"level":      level,
		"experience": experience,
		"mood":       mood,
		"is_active":  isActive,
	}).Error
}

// findPetCollection 取某物种的图鉴记录；无记录返回 (nil, nil)。
func findPetCollection(tx *gorm.DB, studentID uint, species string) (*models.PetCollection, error) {
	if species == "" {
		return nil, nil
	}
	var row models.PetCollection
	err := tx.Where("student_id = ? AND species = ?", studentID, species).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}
