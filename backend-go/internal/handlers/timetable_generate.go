// 课表自动排课处理器：规则排课（单班纯计算预览 / 全校智能排课 commit=true 落库）。
// 路由路径逐字对齐 Laravel backend/routes/api.php 第 147-148 行；
// 入参校验逐条对齐 TimetableController::generateRules()（第 346-361 行）与 generateSchool 的 commit。
package handlers

import (
	"strconv"
	"strings"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/middleware"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
)

// generateRequest 排课请求体（rules 必填；commit 仅全校排课使用，缺省 false）。
type generateRequest struct {
	Rules  *generateRulesPayload `json:"rules"`
	Commit any                   `json:"commit"`
}

// generateRulesPayload rules.days / rules.subjects 原样承接（null 与缺省同样处理）。
type generateRulesPayload struct {
	Days     *[]int                    `json:"days"`
	Subjects *[]generateSubjectPayload `json:"subjects"`
}

// generateSubjectPayload 单个科目的规则入参（nullable 字段用指针与「缺省」区分）。
type generateSubjectPayload struct {
	Name          *string `json:"name"`
	Weekly        *int    `json:"weekly"`
	Double        *bool   `json:"double"`
	Session       *string `json:"session"`
	MaxPerDay     *int    `json:"max_per_day"`
	ForbidPeriods *[]int  `json:"forbid_periods"`
}

// AdminTimetableGenerate 按规则为单个班级生成课表（纯计算，不落库；前端预览后再调保存接口）。
func (h *Handlers) AdminTimetableGenerate(c *gin.Context) {
	u := middleware.CurrentUser(c)

	rules, _, bok := bindGenerateRequest(c)
	if !bok {
		return
	}

	result, err := h.timetable.Generate(u.SchoolID, rules)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// AdminTimetableGenerateSchool 全校智能排课（依据任课表，教师冲突硬约束；commit=true 落库）。
func (h *Handlers) AdminTimetableGenerateSchool(c *gin.Context) {
	u := middleware.CurrentUser(c)

	rules, commit, bok := bindGenerateRequest(c)
	if !bok {
		return
	}

	result, err := h.timetable.GenerateSchool(u.SchoolID, rules, commit)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// bindGenerateRequest 解析并校验排课规则；校验失败时已写出 422，返回 ok=false。
// 第三个返回值对应 Laravel `$request->boolean('commit', false)`。
func bindGenerateRequest(c *gin.Context) (services.GenerateRules, bool, bool) {
	var req generateRequest
	if !bindJSON(c, &req) {
		return services.GenerateRules{}, false, false
	}

	commit, commitOK := laravelBoolean(req.Commit)
	if !commitOK {
		fail(c, services.ErrUnprocessable("commit 需为布尔值"))
		return services.GenerateRules{}, false, false
	}

	if req.Rules == nil {
		fail(c, services.ErrUnprocessable("rules 必填"))
		return services.GenerateRules{}, false, false
	}
	if req.Rules.Days == nil || len(*req.Rules.Days) == 0 {
		fail(c, services.ErrUnprocessable("rules.days 必填且至少需要一个上课日"))
		return services.GenerateRules{}, false, false
	}
	for _, day := range *req.Rules.Days {
		if day < 1 || day > 7 {
			fail(c, services.ErrUnprocessable("rules.days.* 需在 1-7 之间"))
			return services.GenerateRules{}, false, false
		}
	}
	if req.Rules.Subjects == nil || len(*req.Rules.Subjects) == 0 {
		fail(c, services.ErrUnprocessable("rules.subjects 必填且至少需要一个科目"))
		return services.GenerateRules{}, false, false
	}

	rules := services.GenerateRules{Days: *req.Rules.Days, Subjects: []services.GenerateRuleSubject{}}
	for i, sub := range *req.Rules.Subjects {
		prefix := "rules.subjects." + strconv.Itoa(i)

		if sub.Name == nil || trimSpace(*sub.Name) == "" {
			fail(c, services.ErrUnprocessable(prefix+".name 必填"))
			return services.GenerateRules{}, false, false
		}
		if runeLen(*sub.Name) > 50 {
			fail(c, services.ErrUnprocessable(prefix+".name 不能超过 50 字"))
			return services.GenerateRules{}, false, false
		}
		if sub.Weekly == nil {
			fail(c, services.ErrUnprocessable(prefix+".weekly 必填"))
			return services.GenerateRules{}, false, false
		}
		if *sub.Weekly < 1 || *sub.Weekly > 35 {
			fail(c, services.ErrUnprocessable(prefix+".weekly 需在 1-35 之间"))
			return services.GenerateRules{}, false, false
		}

		item := services.GenerateRuleSubject{
			Name:   trimSpace(*sub.Name),
			Weekly: *sub.Weekly,
		}
		if sub.Double != nil {
			item.Double = *sub.Double
		}
		if sub.Session != nil {
			session := trimSpace(*sub.Session)
			if session != "" && session != "any" && session != "am" && session != "pm" {
				fail(c, services.ErrUnprocessable(prefix+".session 取值需为 any/am/pm"))
				return services.GenerateRules{}, false, false
			}
			item.Session = session
		}
		if sub.MaxPerDay != nil {
			if *sub.MaxPerDay < 1 || *sub.MaxPerDay > 8 {
				fail(c, services.ErrUnprocessable(prefix+".max_per_day 需在 1-8 之间"))
				return services.GenerateRules{}, false, false
			}
			item.MaxPerDay = *sub.MaxPerDay
		}
		if sub.ForbidPeriods != nil {
			for _, period := range *sub.ForbidPeriods {
				if period < 1 || period > 30 {
					fail(c, services.ErrUnprocessable(prefix+".forbid_periods.* 需在 1-30 之间"))
					return services.GenerateRules{}, false, false
				}
			}
			item.ForbidPeriods = *sub.ForbidPeriods
		}

		rules.Subjects = append(rules.Subjects, item)
	}

	return rules, commit, true
}

// laravelBoolean 按 Laravel `boolean` 校验 + `$request->boolean()` 口径解析：
// 接受 true/false、1/0、"1"/"0"（null 视为 false）；其它取值返回 ok=false。
func laravelBoolean(value any) (bool, bool) {
	switch v := value.(type) {
	case nil:
		return false, true
	case bool:
		return v, true
	case float64:
		if v == 0 {
			return false, true
		}
		if v == 1 {
			return true, true
		}
		return false, false
	case string:
		switch strings.TrimSpace(v) {
		case "0":
			return false, true
		case "1":
			return true, true
		}
		return false, false
	default:
		return false, false
	}
}
