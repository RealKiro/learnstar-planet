// AI 配置与对话模型。
//
// 移植自 Laravel 迁移 2026_07_21_000001_create_ai_tables / 2026_07_22_000001_add_providers_to_ai_settings /
// 2026_07_22_000002_add_provider_to_ai_conversations / 2026_08_05_000001_add_billing_fields_to_ai_conversations，
// 字段与默认值逐条对齐：
//   - ai_settings 每校一行（school_id），默认 enabled=false / provider=openai / model=gpt-3.5-turbo /
//     max_tokens=2000 / tokens_used=0 / tokens_limit=1000000，providers 为可空 JSON。
//   - ai_conversations 为对话流水（school_id + 可选 class_id/student_id/student_name/provider + question/answer
//   - tokens_used/prompt_tokens/completion_tokens + cost/currency + status）。
package models

import "time"

// AISetting 学校级 AI 配置（每校一行，Laravel 端用 firstOrCreate 惰性创建）。
type AISetting struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	SchoolID uint   `gorm:"index;not null" json:"school_id"`
	Enabled  bool   `gorm:"default:false" json:"enabled"`
	Provider string `gorm:"size:255;default:openai" json:"provider"`
	APIKey   string `gorm:"size:255" json:"api_key"`
	APIBase  string `gorm:"size:500" json:"api_base"`
	Model    string `gorm:"size:255;default:gpt-3.5-turbo" json:"model"`
	// MaxTokens 单次对话最大输出 token（Laravel max_tokens，默认 2000，接口校验 100-32000）。
	MaxTokens int `gorm:"default:2000" json:"max_tokens"`
	// TokensUsed 已消耗 token 累计（每次对话按上游返回的 tokens_used 累加）。
	TokensUsed int `gorm:"default:0" json:"tokens_used"`
	// TokensLimit token 上限（默认 1000000）。注意：Laravel 只在界面展示该值，**不做超限拒绝**，
	// Go 端同样不拦截（见 services 包注释）。
	TokensLimit int `gorm:"default:1000000" json:"tokens_limit"`
	// Providers 多供应商配置（JSON 文本列，语义同 Laravel ai_settings.providers）：
	// 数组，每项含 id/label/api_key/api_base/model/models/model_map/is_active/单价与计数器等。
	// json:"-" 同 School.Settings：列内存 JSON 文本，对外输出由各接口解码后构造视图。
	Providers string    `gorm:"type:text" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AIConversation 一条 AI 对话流水（教师端与大屏端共用一张表）。
type AIConversation struct {
	ID       uint  `gorm:"primaryKey" json:"id"`
	SchoolID uint  `gorm:"index;not null" json:"school_id"`
	ClassID  *uint `gorm:"index" json:"class_id"`
	// StudentID / StudentName 仅大屏端可能填写（Laravel 的 disp_/class_ token 载荷里没有学生信息，
	// 故实际落库为 student_id=NULL、student_name='匿名'；教师端为 student_name='教师'）。
	StudentID   *uint   `gorm:"index" json:"student_id"`
	StudentName *string `gorm:"size:100" json:"student_name"`
	Provider    *string `gorm:"size:50" json:"provider"`
	Question    string  `gorm:"type:text;not null" json:"question"`
	Answer      *string `gorm:"type:text" json:"answer"`
	TokensUsed  int     `gorm:"default:0" json:"tokens_used"`
	// PromptTokens / CompletionTokens 上游返回的输入/输出 token 拆分。
	PromptTokens     int       `gorm:"default:0" json:"prompt_tokens"`
	CompletionTokens int       `gorm:"default:0" json:"completion_tokens"`
	Cost             float64   `gorm:"type:decimal(12,6);default:0" json:"cost"`
	Currency         string    `gorm:"size:8;default:CNY" json:"currency"`
	Status           string    `gorm:"size:20;default:completed" json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
