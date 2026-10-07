// Package models 定义 Go 后端的数据模型（重设计的精简 schema）。
// 与 Laravel 端不共享表结构，二者独立演进，互不影响。
package models

import "time"

// School 学校。
type School struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	Name             string `gorm:"size:200;not null" json:"name"`
	Code             string `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Address          string `gorm:"size:500" json:"address"`
	ContactPhone     string `gorm:"size:64" json:"contact_phone"`
	ContactEmail     string `gorm:"size:200" json:"contact_email"`
	LogoPath         string `gorm:"size:500" json:"logo_path"`
	Status           string `gorm:"size:20;default:active;index" json:"status"`
	ScoreRulesSeeded bool   `gorm:"default:false" json:"score_rules_seeded"`
	// Settings 学校级设置（JSON 文本列，语义同 Laravel schools.settings；
	// 班级码前缀 display_code_prefix 存于此）。
	// json:"-" ：列内存 JSON 文本，直接序列化会与 Laravel 的 settings 对象形状不一致，
	// 需要时经 SettingString / SettingsMap 读取（对外输出由各接口构造视图）。
	Settings  string    `gorm:"type:text" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User 教师 / 学校管理员账号。
type User struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	SchoolID        uint       `gorm:"index;not null" json:"school_id"`
	Role            string     `gorm:"size:20;not null;index" json:"role"` // school_admin | teacher
	Username        string     `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash    string     `gorm:"size:255;not null" json:"-"`
	Name            string     `gorm:"size:100;not null" json:"name"`
	Nickname        string     `gorm:"size:100" json:"nickname"`
	Subject         string     `gorm:"size:100" json:"subject"`
	GradeTeam       string     `gorm:"size:100" json:"grade_team"`
	AvatarPath      string     `gorm:"size:500" json:"avatar_path"`
	Phone           string     `gorm:"size:32" json:"phone"`
	Email           string     `gorm:"size:200" json:"email"`
	Status          string     `gorm:"size:20;default:active;index" json:"status"`
	IsAPIBot        bool       `gorm:"default:false" json:"is_api_bot"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	PasswordChanged bool       `gorm:"default:false" json:"password_changed"`
	// PlainPassword 明文密码（Laravel users.plain_password，供管理员「查看教师密码」使用）。
	// json:"-"：列存在但**不随 User 序列化输出**——Laravel 的 User 未把该字段放进 $hidden，
	// 因此 Laravel 的 user JSON 会带上它；Go 端刻意不扩大泄漏面，改由
	// GET /admin/teachers/:id/password 单点返回（见 README「本批有意差异」）。
	PlainPassword string `gorm:"size:100" json:"-"`
	// Settings 用户级设置（JSON 文本列，语义同 Laravel users.settings）。
	Settings  string    `gorm:"type:text" json:"settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsSchoolAdmin 判断是否学校管理员。
func (u *User) IsSchoolAdmin() bool { return u.Role == "school_admin" }

// IsTeacher 判断是否教师（含 API 机器人，其 role 也是 teacher）。
func (u *User) IsTeacher() bool { return u.Role == "teacher" }

// ClassRoom 班级。
type ClassRoom struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	SchoolID             uint       `gorm:"index;not null" json:"school_id"`
	Name                 string     `gorm:"size:100;not null" json:"name"`
	Grade                string     `gorm:"size:50" json:"grade"`
	Year                 string     `gorm:"size:20" json:"year"`
	TeacherID            *uint      `gorm:"index" json:"teacher_id"`
	MaxStudents          int        `gorm:"default:0" json:"max_students"`
	DisplayCode          string     `gorm:"size:64;index" json:"display_code"`
	DisplayCodeUpdatedAt *time.Time `json:"display_code_updated_at"`
	// Settings 班级级设置（JSON 文本列，语义同 Laravel class_rooms.settings；pet_series 存于此）。
	// json:"-" 同 School.Settings：对外输出由各接口构造视图，避免暴露仓内 JSON 文本。
	Settings  string    `gorm:"type:text" json:"-"`
	Status    string    `gorm:"size:20;default:active;index" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Student 学生。
type Student struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	ClassID    uint   `gorm:"index;not null" json:"class_id"`
	Name       string `gorm:"size:100;not null" json:"name"`
	Gender     string `gorm:"size:10" json:"gender"`
	StudentNo  string `gorm:"size:64;index" json:"student_no"`
	AvatarPath string `gorm:"size:500" json:"avatar_path"`
	// ParentID 家长用户 ID（Laravel students.parent_id）。企微请假同步据此匹配学生：
	// 家长企微绑定（third_party_bindings.platform=wechat_work）→ 该用户名下的活跃学生。
	// ⚠️ Laravel 已于 2026_08_06_000006 迁移把家长角色移除并清空该列，故两边都只在
	// 「家长数据被外部写入」时才会命中（见 README「企微回调与请假同步」小节）。
	ParentID   *uint      `gorm:"index" json:"parent_id"`
	TotalScore int        `gorm:"default:0" json:"total_score"`
	Status     string     `gorm:"size:20;default:active;index" json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"-"`
}

// ScoreRule 积分规则（学校级共享 + 班级级自定义）。
type ScoreRule struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SchoolID   *uint     `gorm:"index" json:"school_id"`
	ClassID    *uint     `gorm:"index" json:"class_id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Amount     int       `gorm:"not null" json:"amount"`
	Category   string    `gorm:"size:32;default:custom;index" json:"category"`
	IsPositive bool      `gorm:"default:true" json:"is_positive"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	SortOrder  int       `gorm:"default:0" json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Score 一条积分变动记录（正为加分、负为减分）。
type Score struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	StudentID   uint   `gorm:"index;not null" json:"student_id"`
	ClassID     uint   `gorm:"index;not null" json:"class_id"`
	ScoreRuleID *uint  `gorm:"index" json:"score_rule_id"`
	Amount      int    `gorm:"not null" json:"amount"`
	Reason      string `gorm:"size:500" json:"reason"`
	GivenBy     uint   `gorm:"index" json:"given_by"`
	// UndoOfScoreID 指向被撤回的原积分记录（仅「撤回」产生的反向流水有值）。
	// uniqueIndex 由数据库层保证「一条记录只能被撤回一次」——撤回是加积分操作，
	// 缺这条约束时可反复点撤回把积分刷上去（见 services/score.go Undo 的幂等守卫）。
	UndoOfScoreID *uint     `gorm:"uniqueIndex" json:"undo_of_score_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ScoreLog 积分审计日志（记录余额变化前后值）。
type ScoreLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	StudentID     uint      `gorm:"index;not null" json:"student_id"`
	ScoreID       uint      `gorm:"index;not null" json:"score_id"`
	BalanceBefore int       `gorm:"not null" json:"balance_before"`
	BalanceAfter  int       `gorm:"not null" json:"balance_after"`
	Description   string    `gorm:"size:500" json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Notice 班级通知/公告。
type Notice struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ClassID     uint       `gorm:"index;not null" json:"class_id"`
	SchoolID    uint       `gorm:"index;not null" json:"school_id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Content     string     `gorm:"type:text" json:"content"`
	Type        string     `gorm:"size:20;default:info;index" json:"type"`
	PublishedBy uint       `gorm:"index" json:"published_by"`
	IsPublished bool       `gorm:"default:false" json:"is_published"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ShopItem 商城商品（学校级共享 + 班级级自定义）。
type ShopItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ClassID      *uint     `gorm:"index" json:"class_id"`
	SchoolID     uint      `gorm:"index;not null" json:"school_id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Description  string    `gorm:"size:500" json:"description"`
	Category     string    `gorm:"size:32;default:stationery;index" json:"category"`
	CostScore    int       `gorm:"not null" json:"cost_score"`
	CurrencyType string    `gorm:"size:32;default:score;index" json:"currency_type"`
	EventTag     string    `gorm:"size:50" json:"event_tag"`
	Stock        int       `gorm:"default:0" json:"stock"`
	ImagePath    string    `gorm:"size:500" json:"image_path"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ShopRedemption 商城兑换记录（pending → approved/rejected → delivered）。
type ShopRedemption struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	StudentID   uint       `gorm:"index;not null" json:"student_id"`
	ShopItemID  uint       `gorm:"index;not null" json:"shop_item_id"`
	ClassID     uint       `gorm:"index;not null" json:"class_id"`
	Cost        int        `gorm:"not null" json:"cost"`
	Status      string     `gorm:"size:20;default:pending;index" json:"status"`
	ApprovedBy  *uint      `json:"approved_by"`
	ApprovedAt  *time.Time `json:"approved_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Attendance 考勤记录（date 以 "2006-01-02" 形式存储，便于跨库比较）。
type Attendance struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	ClassID       uint       `gorm:"index;not null" json:"class_id"`
	StudentID     uint       `gorm:"index;not null" json:"student_id"`
	TeacherID     *uint      `json:"teacher_id"`
	Date          string     `gorm:"size:10;index;not null" json:"date"`
	Status        string     `gorm:"size:20;default:present;index" json:"status"`
	Source        string     `gorm:"size:20;default:auto;index" json:"source"`
	Remark        string     `gorm:"size:500" json:"remark"`
	LeaveRecordID *uint      `json:"leave_record_id"`
	SignInAt      *time.Time `json:"sign_in_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Wallet 学生多币种钱包。
type Wallet struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	StudentID    uint      `gorm:"index;not null;uniqueIndex:idx_wallet_student_currency" json:"student_id"`
	CurrencyType string    `gorm:"size:32;not null;uniqueIndex:idx_wallet_student_currency" json:"currency_type"`
	Balance      int       `gorm:"default:0" json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ExchangeRate 学校级币种汇率（如 2 积分 = 1 科学币 → rate 0.5）。
type ExchangeRate struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SchoolID     uint      `gorm:"index;not null" json:"school_id"`
	Name         string    `gorm:"size:100" json:"name"`
	FromCurrency string    `gorm:"size:32;not null" json:"from_currency"`
	ToCurrency   string    `gorm:"size:32;not null" json:"to_currency"`
	Rate         float64   `gorm:"not null" json:"rate"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ExchangeLog 币种兑换记录。
type ExchangeLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	StudentID    uint      `gorm:"index;not null" json:"student_id"`
	FromCurrency string    `gorm:"size:32;not null" json:"from_currency"`
	ToCurrency   string    `gorm:"size:32;not null" json:"to_currency"`
	FromAmount   int       `gorm:"not null" json:"from_amount"`
	ToAmount     int       `gorm:"not null" json:"to_amount"`
	OperatedBy   *uint     `json:"operated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Broadcast 教室大屏实时广播（banner/popup/fullscreen）。
type Broadcast struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	SchoolID       uint       `gorm:"index;not null" json:"school_id"`
	ClassID        uint       `gorm:"index;not null" json:"class_id"`
	TeacherID      *uint      `gorm:"index" json:"teacher_id"`
	Content        string     `gorm:"size:500;not null" json:"content"`
	Type           string     `gorm:"size:20;default:banner" json:"type"`
	VoiceEnabled   bool       `gorm:"default:true" json:"voice_enabled"`
	LoopEnabled    bool       `gorm:"default:false" json:"loop_enabled"`
	DisplaySeconds int        `gorm:"default:10" json:"display_seconds"`
	Status         string     `gorm:"size:20;default:sent;index" json:"status"`
	SentAt         *time.Time `json:"sent_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ClassRoomTeacher 班级-教师关联（多教师多角色）。
// 移植自 Laravel `class_room_teachers`（迁移 2026_07_13_000003_create_class_room_teachers.php
// + 2026_07_24_000001_add_subject_to_class_room_teachers.php）。
// role ∈ head_teacher / co_teacher / subject_teacher / grade_lead / admin_director；
// head_teacher 的关联在写入时同步回写 class_rooms.teacher_id（同 Laravel
// ClassRoomTeacher::updateOrCreate + ClassRoom.teacher_id 的联动）。
// (class_room_id, user_id) 唯一：与 Laravel updateOrCreate 的查找键一致。
type ClassRoomTeacher struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ClassRoomID uint      `gorm:"index;not null;uniqueIndex:idx_class_room_teacher" json:"class_room_id"`
	UserID      uint      `gorm:"index;not null;uniqueIndex:idx_class_room_teacher" json:"user_id"`
	Role        string    `gorm:"size:32;default:subject_teacher" json:"role"`
	Subject     string    `gorm:"size:100" json:"subject"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
