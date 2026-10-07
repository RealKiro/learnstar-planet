// 班级码登录（学生端 / 班级大屏统一入口）：POST /api/v1/auth/class/login。
//
// 忠实移植自 Laravel App\Http\Controllers\Api\AuthController::classLogin：
//  1. 班级码大写去首尾空格后先按 display_code 查；查不到再按确定性计算码全库匹配；
//  2. 无命中 → 401「班级码无效，请核对后重试」；命中多于一个 → 401「班级码在多所学校存在，请向班主任确认」；
//  3. 签发 class_<班级码>_<32 位随机串>，有效期 24 小时；
//  4. 响应 data = {token, class_id, class_name, grade, student_count}。
//
// 有意差异：
//  1. Laravel 把 class_token:<token> → class_id 存 Cache（24 小时）；Go 端无 Cache，
//     复用既有 display_tokens 表（models.ResolveDisplayToken 本就兼容 class_ 前缀），
//     故这里签发的 token 可直接用于全部 display/* 接口（含 SSE / poll）。
//  2. Laravel 该路由带 throttle:10,1；Go 端已按同参数挂 `middleware.Throttle`
//     （进程内固定窗口，仅单实例有效；与 display/login 的「同码失败锁定」机制并存）。
//  3. 响应逐字同 Laravel —— 只有 data、没有 message 字段（Go 其余接口统一加 message:"ok"）。
package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
)

// ClassLoginResult 班级码登录结果（字段名逐字同 Laravel classLogin 的 data）。
type ClassLoginResult struct {
	Token        string `json:"token"`
	ClassID      uint   `json:"class_id"`
	ClassName    string `json:"class_name"`
	Grade        string `json:"grade"`
	StudentCount int64  `json:"student_count"`
}

// ClassLogin 班级码登录：命中唯一班级后签发 24 小时有效的 class_ token。
func (d *DisplayService) ClassLogin(rawCode string) (*ClassLoginResult, error) {
	code := strings.ToUpper(strings.TrimSpace(rawCode))

	classes, err := d.resolveCodeToClasses(code)
	if err != nil {
		return nil, err
	}
	if len(classes) == 0 {
		// 注意与 DisplayService.Login 的「班级码无效，请检查后重试」不同字（Laravel 亦然）。
		return nil, NewAppError(http.StatusUnauthorized, "班级码无效，请核对后重试")
	}
	if len(classes) > 1 {
		return nil, NewAppError(http.StatusUnauthorized, "班级码在多所学校存在，请向班主任确认")
	}

	class := classes[0]
	token, err := models.NewClassTokenValue(code)
	if err != nil {
		return nil, err
	}

	// 过期行惰性清理：Laravel 由 Cache TTL 自动回收（同 DisplayService.issueToken）。
	// ⚠️ 写入与比较统一 UTC 域（SQLite 对 time 列做文本比较，见 services/auth.go 的时区纪律）。
	now := time.Now().UTC()
	_ = d.db.Where("expires_at < ?", now).Delete(&models.DisplayToken{}).Error
	if err := d.db.Create(&models.DisplayToken{
		Token:     token,
		ClassID:   class.ID,
		ClassName: class.Name,
		Grade:     class.Grade,
		ExpiresAt: now.Add(models.DisplayTokenTTL * time.Second),
	}).Error; err != nil {
		return nil, err
	}

	studentCount, err := d.activeStudentCount(class.ID)
	if err != nil {
		return nil, err
	}

	return &ClassLoginResult{
		Token:        token,
		ClassID:      class.ID,
		ClassName:    class.Name,
		Grade:        class.Grade,
		StudentCount: studentCount,
	}, nil
}

// resolveCodeToClasses 解析班级码 → 全部命中班级（可能跨校多命中，故返回切片而不是单条）。
// 遍历顺序固定 id ASC，保证多命中时的行为可复现（Laravel 未显式排序）。
func (d *DisplayService) resolveCodeToClasses(code string) ([]models.ClassRoom, error) {
	if code == "" {
		return nil, nil
	}

	var rooms []models.ClassRoom
	if err := d.db.Where("display_code = ?", code).Order("id ASC").Find(&rooms).Error; err != nil {
		return nil, err
	}
	if len(rooms) > 0 {
		return rooms, nil
	}

	// 数据库未刷新时按确定性计算码匹配（同 Laravel：遍历全部班级调 generate()）。
	var all []models.ClassRoom
	if err := d.db.Order("id ASC").Find(&all).Error; err != nil {
		return nil, err
	}
	matched := make([]models.ClassRoom, 0, 1)
	for i := range all {
		if d.codes.Generate(&all[i], nil) == code {
			matched = append(matched, all[i])
		}
	}
	return matched, nil
}
