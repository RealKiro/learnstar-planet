// 宠物服务：教师端班级宠物总览、详情、喂养、更名与切换 + 教室端（班级码）换宠。
// 图鉴（pet_collections）：教师端与教室端的换宠都会归档旧物种进度、并在切回该物种时恢复
// （同 Laravel PetService::switchPet / DisplayController::classroomSwitchPet）。
// 大屏事件：仅「喂养」发布 pet_update（与 Laravel PetService 一致——rename / switch 不发事件）。
package services

import (
	"errors"
	"fmt"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// PetService 宠物服务。
type PetService struct {
	db     *gorm.DB
	scope  *Scope
	events *DisplayEvents
}

// NewPetService 创建宠物服务（内部自带大屏事件发布器：pet_update 事件）。
func NewPetService(db *gorm.DB, scope *Scope) *PetService {
	return &PetService{db: db, scope: scope, events: NewDisplayEvents(db)}
}

// PetOverviewItem 班级宠物总览条目。
type PetOverviewItem struct {
	ID          uint   `json:"id"`
	StudentID   uint   `json:"student_id"`
	StudentName string `json:"student_name"`
	Name        string `json:"name"`
	Species     string `json:"species"`
	Level       int    `json:"level"`
	Exp         int    `json:"exp"`
	Mood        int    `json:"mood"`
	StageName   string `json:"stage_name"`
	Emoji       string `json:"emoji"`
}

// ClassOverview 返回教师可管理班级内的宠物总览。
func (p *PetService) ClassOverview(u *models.User) ([]PetOverviewItem, error) {
	classIDs, err := p.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return []PetOverviewItem{}, nil
	}

	var pets []models.Pet
	if err := p.db.Where("class_id IN ?", classIDs).Order("id ASC").Find(&pets).Error; err != nil {
		return nil, err
	}

	names := p.studentNames(pets)

	out := make([]PetOverviewItem, 0, len(pets))
	for _, pet := range pets {
		stage := pet.CurrentStage()
		out = append(out, PetOverviewItem{
			ID:          pet.ID,
			StudentID:   pet.StudentID,
			StudentName: names[pet.StudentID],
			Name:        pet.Name,
			Species:     pet.Species,
			Level:       pet.Level,
			Exp:         pet.Experience,
			Mood:        pet.Mood,
			StageName:   stage.Name,
			Emoji:       stage.Emoji,
		})
	}
	return out, nil
}

// PetFor 返回指定学生的宠物详情（含阶段信息），读取前先做心情衰减。
func (p *PetService) PetFor(u *models.User, studentID uint) (map[string]any, error) {
	student, err := p.scope.StudentInScope(u, studentID)
	if err != nil {
		return nil, err
	}

	pet, err := p.findPet(student.ID)
	if err != nil {
		return nil, err
	}

	now := util.Now()
	pet.DecayMood(now)
	if err := p.db.Save(pet).Error; err != nil {
		return nil, err
	}

	stage := pet.CurrentStage()
	return map[string]any{
		"id":           pet.ID,
		"name":         pet.Name,
		"species":      pet.Species,
		"level":        pet.Level,
		"exp":          pet.Experience,
		"mood":         pet.Mood,
		"emoji":        stage.Emoji,
		"stage_name":   stage.Name,
		"last_fed_at":  pet.LastFedAt,
		"student_id":   student.ID,
		"student_name": student.Name,
	}, nil
}

// Feed 喂养宠物（先衰减后 +20 心情，上限 100）。
func (p *PetService) Feed(u *models.User, studentID uint) (map[string]any, error) {
	student, err := p.scope.StudentInScope(u, studentID)
	if err != nil {
		return nil, err
	}

	pet, err := p.findPet(student.ID)
	if err != nil {
		return nil, err
	}

	now := util.Now()
	pet.DecayMood(now)
	pet.Feed(now)
	if err := p.db.Save(pet).Error; err != nil {
		return nil, err
	}

	// 推送给班级大屏（pet_update，type=feed）：字段/键序同 Laravel PetService::feed。
	publishDisplayEvent(p.events, student.ClassID, DisplayEventPetUpdate, DisplayPetUpdateData{
		StudentID:   student.ID,
		StudentName: student.Name,
		Type:        "feed",
		Mood:        pet.Mood,
		Level:       pet.Level,
		Experience:  pet.Experience,
	})

	return map[string]any{
		"message": fmt.Sprintf("已喂养「%s」", pet.Name),
		"mood":    pet.Mood,
		"level":   pet.Level,
	}, nil
}

// Rename 宠物更名，返回成功提示。
func (p *PetService) Rename(u *models.User, studentID uint, name string) (string, error) {
	student, err := p.scope.StudentInScope(u, studentID)
	if err != nil {
		return "", err
	}

	pet, err := p.findPet(student.ID)
	if err != nil {
		return "", err
	}

	pet.Name = name
	if err := p.db.Save(pet).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("宠物已更名为「%s」", name), nil
}

// TeacherSwitchPetData 教师端换宠的 data（字段与出现条件逐字同 Laravel PetService::switchPet：
// `cost` / `free_pick_used` 只在「原有宠物」分支出现）。
type TeacherSwitchPetData struct {
	PetName      string `json:"pet_name"`
	PetSpecies   string `json:"pet_species"`
	Level        int    `json:"level"`
	Experience   int    `json:"experience"`
	Cost         *int   `json:"cost,omitempty"`
	FreePickUsed *bool  `json:"free_pick_used,omitempty"`
}

// TeacherSwitchPetResult 教师端换宠结果（Message 逐字同 Laravel，由处理器放进 message 字段）。
type TeacherSwitchPetResult struct {
	Message string
	Data    TeacherSwitchPetData
}

// Switch 教师端切换/首次分配宠物。逐条移植自 Laravel PetService::switchPet：
//
//	同物种守卫(422) → 类别限制(422) → 免费自选机会 → 按等级扣分(不足 400) →
//	旧物种进度归档(`is_active=false`) → 恢复目标物种进度（无记录则初始 1/0/80）→
//	目标物种标记激活(`is_active=true`) → 免费自选机会消费即删除。
//
// 扣除积分的口径与 Laravel 相同：直接改 `students.total_score` 后保存（不写 scores /
// score_logs 审计、不发 score_update 事件、不按积分重算等级），仅把整批写入放进一个事务。
//
// 与 Laravel 的既有差异（沿用）：免费自选机会由 Cache 键改为 `pet_free_picks` 表。
func (p *PetService) Switch(u *models.User, studentID uint, species, name string) (*TeacherSwitchPetResult, error) {
	student, err := p.scope.StudentInScope(u, studentID)
	if err != nil {
		return nil, err
	}

	var pet models.Pet
	petErr := p.db.Where("student_id = ?", student.ID).First(&pet).Error
	hasPet := petErr == nil
	if petErr != nil && !errors.Is(petErr, gorm.ErrRecordNotFound) {
		return nil, petErr
	}

	// ===== 同物种切换守卫（不扣费、不消耗免费机会） =====
	if hasPet && pet.Species == species {
		return nil, ErrUnprocessable("当前已经是这只宠物啦")
	}

	// ===== 类别限制：只能在本班当前类别内更换，不能跨类别领养 =====
	// 注 1：旧系列 id（cosmic/cute/all 等）在 SpeciesPoolForSeries 返回空池 → 视为不限制。
	// 注 2：Laravel 教师端此处回显的是**原始系列 id**（`'只能领养当前类别「' . $classSeries . '」…'`），
	// 而教室端回显 seriesLabel 的中文标签 —— 两端文案本就不同，此处照抄教师端口径。
	if series := p.classPetSeries(student.ClassID); series != "" {
		pool := models.SpeciesPoolForSeries(series)
		if len(pool) > 0 && !containsString(pool, species) {
			return nil, ErrUnprocessable("只能领养当前类别「" + series + "」的宠物，不能跨类别领养")
		}
	}

	now := util.Now()
	result := &TeacherSwitchPetResult{}
	err = p.db.Transaction(func(tx *gorm.DB) error {
		// ===== 无宠物：创建新宠物并收入图鉴（同 Laravel 第 248-275 行） =====
		if !hasPet {
			created := models.Pet{
				StudentID:  student.ID,
				ClassID:    student.ClassID,
				Name:       name,
				Species:    species,
				Level:      1,
				Experience: 0,
				Mood:       80,
			}
			if err := tx.Create(&created).Error; err != nil {
				return err
			}
			if err := upsertPetCollection(tx, student.ID, species, 1, 0, 80, true); err != nil {
				return err
			}
			pet = created
			result.Message = fmt.Sprintf("已为您分配宠物「%s」", created.Name)
			result.Data = TeacherSwitchPetData{
				PetName:    created.Name,
				PetSpecies: created.Species,
				Level:      created.Level,
				Experience: created.Experience,
			}
			return nil
		}

		// ===== 免费自选：整班切换后的一次机会（限当前类别，免费） =====
		usedFreePick, err := models.HasPetFreePick(tx, student.ID, now)
		if err != nil {
			return err
		}
		cost := 0
		if !usedFreePick {
			// 1) 先扣积分（等级越高越贵；免费自选不扣）——不足直接拒绝，不留脏数据。
			cost = models.SwitchCost(pet.Level)
			if student.TotalScore < cost {
				return ErrBadRequest(fmt.Sprintf("积分不足，更换宠物需 %d 积分", cost))
			}
			student.TotalScore -= cost
			if err := tx.Model(&models.Student{}).Where("id = ?", student.ID).
				Update("total_score", student.TotalScore).Error; err != nil {
				return err
			}
		}

		// ===== 目标物种的图鉴进度（切换前先查，切回时恢复） =====
		collection, err := findPetCollection(tx, student.ID, species)
		if err != nil {
			return err
		}

		// 2) 保存当前宠物进度到图鉴（进度全保留，旧物种标记为非激活）。
		if err := upsertPetCollection(tx, student.ID, pet.Species, pet.Level, pet.Experience, pet.Mood, false); err != nil {
			return err
		}

		// 3) 恢复目标物种进度（无图鉴记录 → 初始形态；有则等级/经验/心情全保留）。
		pet.Species = species
		if collection != nil {
			pet.Level = collection.Level
			pet.Experience = collection.Experience
			pet.Mood = collection.Mood
		} else {
			pet.Level = 1
			pet.Experience = 0
			pet.Mood = 80
		}
		pet.Name = name
		switchedAt := now
		pet.LastSwitchedAt = &switchedAt
		if err := tx.Save(&pet).Error; err != nil {
			return err
		}

		// 4) 目标物种标记激活。
		if err := upsertPetCollection(tx, student.ID, species, pet.Level, pet.Experience, pet.Mood, true); err != nil {
			return err
		}

		// 5) 免费自选机会使用即失效（同 Laravel Cache::forget）。
		if usedFreePick {
			if err := models.ConsumePetFreePick(tx, student.ID); err != nil {
				return err
			}
		}

		costValue, freePickValue := cost, usedFreePick
		result.Message = fmt.Sprintf("宠物已更换为「%s」（扣除 %d 积分）", pet.Name, cost)
		if usedFreePick {
			result.Message = "✅ 已使用整班切换的免费自选机会！"
		}
		result.Data = TeacherSwitchPetData{
			PetName:      pet.Name,
			PetSpecies:   pet.Species,
			Level:        pet.Level,
			Experience:   pet.Experience,
			Cost:         &costValue,
			FreePickUsed: &freePickValue,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// classPetSeries 读取班级设置里的 pet_series（未设置/班级不存在 → 空串，视为不限制）。
func (p *PetService) classPetSeries(classID uint) string {
	var class models.ClassRoom
	if err := p.db.First(&class, classID).Error; err != nil {
		return ""
	}
	return class.SettingString("pet_series", "")
}

func (p *PetService) findPet(studentID uint) (*models.Pet, error) {
	var pet models.Pet
	err := p.db.Where("student_id = ?", studentID).First(&pet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("该学生还没有宠物")
	}
	if err != nil {
		return nil, err
	}
	return &pet, nil
}

func (p *PetService) studentNames(pets []models.Pet) map[uint]string {
	ids := make([]uint, 0, len(pets))
	for _, pet := range pets {
		ids = append(ids, pet.StudentID)
	}

	out := map[uint]string{}
	if len(ids) == 0 {
		return out
	}

	var students []models.Student
	if err := p.db.Where("id IN ?", ids).Find(&students).Error; err != nil {
		return out
	}
	for _, st := range students {
		out[st.ID] = st.Name
	}
	return out
}

// ============================================================
// 教室端（班级码）换宠：POST /api/v1/display/pets/switch
// 移植自 Laravel DisplayController::classroomSwitchPet。
// ============================================================

// ClassroomSwitchPetData 教室端换宠结果（字段名逐字同 Laravel；free_pick_used 只在该分支出现）。
type ClassroomSwitchPetData struct {
	PetEmoji     string `json:"pet_emoji"`
	PetName      string `json:"pet_name"`
	PetSpecies   string `json:"pet_species"`
	TotalScore   int    `json:"total_score"`
	Cost         int    `json:"cost"`
	FreePickUsed *bool  `json:"free_pick_used,omitempty"`
}

// ClassroomSwitchPetResult 教室端换宠结果（Message 为 Laravel 原文，由处理器放进 message 字段）。
type ClassroomSwitchPetResult struct {
	Message string
	Data    ClassroomSwitchPetData
}

// SwitchForClassroom 学生自行换宠：同物种守卫 → 班级系列限制 → 免费自选或按等级扣分 →
// 归档旧物种进度 → 只换 species（等级/经验/心情保留）→ 目标物种标记激活。
//
// 「保留等级和经验，只换种类」是 Laravel `DisplayController::classroomSwitchPet` 的原生行为
// （与教师端 `PetService::switchPet` 的「恢复目标物种进度」口径**不同**，两端各自照抄）。
//
// 有意差异：
//  1. 免费自选机会：Laravel 查 Cache 的 pet_free_pick:<id>，Go 端读 pet_free_picks 表
//     （见 models/pet_free_pick.go），消费即删除，语义等价。
//  2. 扣分：Laravel 直接改 student.total_score 后 save（不写 scores / score_logs 审计、不发
//     score_update 事件、不同步宠物经验），此处保持一致，仅把扣分、归档与换种放进同一事务。
//  3. 无宠物分支：Laravel 只建 pets、**不写图鉴**（同 DisplayController 第 1268-1288 行），此处照抄。
func (p *PetService) SwitchForClassroom(classID, studentID uint, newSpecies string) (*ClassroomSwitchPetResult, error) {
	var class models.ClassRoom
	err := p.db.First(&class, classID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("班级不存在")
	}
	if err != nil {
		return nil, err
	}

	var student models.Student
	err = p.db.Where("class_id = ? AND id = ?", classID, studentID).First(&student).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("学生不存在")
	}
	if err != nil {
		return nil, err
	}

	var pet models.Pet
	petErr := p.db.Where("student_id = ?", student.ID).First(&pet).Error
	if petErr != nil && !errors.Is(petErr, gorm.ErrRecordNotFound) {
		return nil, petErr
	}
	hasPet := petErr == nil

	// 同物种守卫：不扣分、不消耗免费机会（Laravel 先判此项，再判系列限制）。
	if hasPet && pet.Species == newSpecies {
		return nil, ErrUnprocessable("当前已经是这只宠物啦")
	}

	// 类别限制：只能在本班当前类别内更换，不能跨类别领养（旧系列 id 的空池视为不限制）。
	series, _ := class.SettingsMap()["pet_series"].(string)
	if series != "" {
		pool := models.SpeciesPoolForSeries(series)
		if len(pool) > 0 && !containsString(pool, newSpecies) {
			return nil, ErrUnprocessable("只能领养当前类别「" + models.SeriesLabel(series) + "」的宠物，不能跨类别领养")
		}
	}

	if hasPet {
		now := util.Now()
		usedFreePick, err := models.HasPetFreePick(p.db, student.ID, now)
		if err != nil {
			return nil, err
		}

		cost := 0
		if !usedFreePick {
			cost = models.SwitchCost(pet.Level)
			if student.TotalScore < cost {
				return nil, ErrBadRequest(fmt.Sprintf("积分不足，更换宠物需 %d 积分", cost))
			}
		}

		err = p.db.Transaction(func(tx *gorm.DB) error {
			if !usedFreePick {
				if err := tx.Model(&models.Student{}).Where("id = ?", student.ID).
					Update("total_score", student.TotalScore-cost).Error; err != nil {
					return err
				}
				student.TotalScore -= cost
			}

			// 旧物种进度存入图鉴（保留等级/经验/心情，标记为非激活）。
			if err := upsertPetCollection(tx, student.ID, pet.Species, pet.Level, pet.Experience, pet.Mood, false); err != nil {
				return err
			}

			pet.Species = newSpecies
			if err := tx.Save(&pet).Error; err != nil {
				return err
			}

			// 目标物种标记激活（课堂端同样记录到图鉴，避免「已拥有却显示未拥有」）。
			if err := upsertPetCollection(tx, student.ID, newSpecies, pet.Level, pet.Experience, pet.Mood, true); err != nil {
				return err
			}

			if usedFreePick {
				return models.ConsumePetFreePick(tx, student.ID)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		stage := pet.CurrentStage()
		freePickUsed := usedFreePick
		message := fmt.Sprintf("✅ 已更换为「%s」，扣除 %d 积分", pet.Name, cost)
		if usedFreePick {
			message = "✅ 已使用整班切换的免费自选机会！"
		}
		return &ClassroomSwitchPetResult{
			Message: message,
			Data: ClassroomSwitchPetData{
				PetEmoji:     stage.Emoji,
				PetName:      pet.Name,
				PetSpecies:   newSpecies,
				TotalScore:   student.TotalScore,
				Cost:         cost,
				FreePickUsed: &freePickUsed,
			},
		}, nil
	}

	// 无宠物：创建新宠物（Laravel 的命名口径是「<学生名>的伙伴」，与 initialData 自动分配的
	// 「<学生名>的萌宠」不同字，此处保持一致）。
	newPet := models.Pet{
		StudentID:  student.ID,
		ClassID:    student.ClassID,
		Name:       student.Name + "的伙伴",
		Species:    newSpecies,
		Level:      1,
		Experience: 0,
		Mood:       80,
	}
	if err := p.db.Create(&newPet).Error; err != nil {
		return nil, err
	}

	stage := newPet.CurrentStage()
	return &ClassroomSwitchPetResult{
		Message: "🎉 新宠物已诞生！",
		Data: ClassroomSwitchPetData{
			PetEmoji:   stage.Emoji,
			PetName:    newPet.Name,
			PetSpecies: newSpecies,
			TotalScore: student.TotalScore,
			Cost:       0,
		},
	}, nil
}

// containsString 判断字符串是否在切片内（等价 PHP in_array($needle, $haystack, true)）。
func containsString(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}
