// 商城服务：商品管理与兑换状态机（pending → approved/rejected → delivered）。
// 忠实移植自 Laravel App\Services\ShopService 的结算规则与默认商品清单。
package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// ShopService 商城业务服务。
type ShopService struct {
	db       *gorm.DB
	scope    *Scope
	scores   *ScoreService
	currency *CurrencyService
}

// NewShopService 创建商城服务。
func NewShopService(db *gorm.DB, scope *Scope, scores *ScoreService, currency *CurrencyService) *ShopService {
	return &ShopService{db: db, scope: scope, scores: scores, currency: currency}
}

// Items 返回学校级（class_id=null）+ 教师管辖班级级商品，可按币种过滤。
// 首次访问且学校级商品为空时自动播种默认商品（幂等）。
func (s *ShopService) Items(u *models.User, currencyType string) ([]models.ShopItem, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	q := s.db.Where("(school_id = ? AND class_id IS NULL) OR class_id IN ?", u.SchoolID, classIDs)
	if currencyType != "" {
		q = q.Where("currency_type = ?", currencyType)
	}

	var items []models.ShopItem
	if err := q.Order("category ASC, cost_score ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	if len(items) == 0 && u.SchoolID != 0 {
		if err := s.seedDefaultItems(u.SchoolID); err != nil {
			return nil, err
		}
		// 与 Laravel 一致：播种后返回学校级默认商品（不带币种过滤）。
		var defaults []models.ShopItem
		if err := s.db.Where("school_id = ? AND class_id IS NULL", u.SchoolID).
			Order("category ASC, cost_score ASC").Find(&defaults).Error; err != nil {
			return nil, err
		}
		return defaults, nil
	}

	return items, nil
}

// FindItem 按作用域查找商品：学校级或本班班级级（越权即 404）。
func (s *ShopService) FindItem(u *models.User, id uint) (*models.ShopItem, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	var item models.ShopItem
	err = s.db.Where("id = ? AND ((school_id = ? AND class_id IS NULL) OR class_id IN ?)",
		id, u.SchoolID, classIDs).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("商品不存在或不可见")
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// CreateItem 创建学校级商品（class_id=null）。
func (s *ShopService) CreateItem(u *models.User, name, description, category, currencyType, eventTag, imagePath string, costScore, stock int) (*models.ShopItem, error) {
	if category == "" {
		category = "physical"
	}
	if currencyType == "" {
		currencyType = "score"
	}

	item := models.ShopItem{
		ClassID:      nil,
		SchoolID:     u.SchoolID,
		Name:         name,
		Description:  description,
		Category:     category,
		CostScore:    costScore,
		CurrencyType: currencyType,
		EventTag:     eventTag,
		Stock:        stock,
		ImagePath:    imagePath,
		IsActive:     true,
	}
	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateItem 更新商品字段。
func (s *ShopService) UpdateItem(item *models.ShopItem, updates map[string]any) (*models.ShopItem, error) {
	if err := s.db.Model(item).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(item, item.ID).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteItem 删除商品。
func (s *ShopService) DeleteItem(item *models.ShopItem) error {
	return s.db.Delete(item).Error
}

// RedemptionView 兑换记录列表视图（含学生与商品摘要）。
type RedemptionView struct {
	ID           uint       `json:"id"`
	StudentID    uint       `json:"student_id"`
	StudentName  string     `json:"student_name"`
	StudentNo    string     `json:"student_no"`
	ShopItemID   uint       `json:"shop_item_id"`
	ItemName     string     `json:"item_name"`
	Cost         int        `json:"cost"`
	CurrencyType string     `json:"currency_type"`
	EventTag     string     `json:"event_tag"`
	Category     string     `json:"category"`
	Status       string     `json:"status"`
	ClassID      uint       `json:"class_id"`
	ApprovedBy   *uint      `json:"approved_by"`
	ApprovedAt   *time.Time `json:"approved_at"`
	DeliveredAt  *time.Time `json:"delivered_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ListRedemptions 返回教师管辖班级内的兑换记录（含学生与商品摘要，按时间倒序）。
func (s *ShopService) ListRedemptions(u *models.User) ([]RedemptionView, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	if len(classIDs) == 0 {
		return []RedemptionView{}, nil
	}

	var redemptions []models.ShopRedemption
	if err := s.db.Where("class_id IN ?", classIDs).
		Order("created_at DESC, id DESC").Find(&redemptions).Error; err != nil {
		return nil, err
	}
	if len(redemptions) == 0 {
		return []RedemptionView{}, nil
	}

	studentIDs := make([]uint, 0, len(redemptions))
	itemIDs := make([]uint, 0, len(redemptions))
	for _, r := range redemptions {
		studentIDs = append(studentIDs, r.StudentID)
		itemIDs = append(itemIDs, r.ShopItemID)
	}

	var students []models.Student
	if err := s.db.Where("id IN ?", studentIDs).Find(&students).Error; err != nil {
		return nil, err
	}
	studentMap := make(map[uint]models.Student, len(students))
	for _, st := range students {
		studentMap[st.ID] = st
	}

	var items []models.ShopItem
	if err := s.db.Where("id IN ?", itemIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	itemMap := make(map[uint]models.ShopItem, len(items))
	for _, it := range items {
		itemMap[it.ID] = it
	}

	out := make([]RedemptionView, 0, len(redemptions))
	for _, r := range redemptions {
		v := RedemptionView{
			ID:          r.ID,
			StudentID:   r.StudentID,
			ShopItemID:  r.ShopItemID,
			Cost:        r.Cost,
			Status:      r.Status,
			ClassID:     r.ClassID,
			ApprovedBy:  r.ApprovedBy,
			ApprovedAt:  r.ApprovedAt,
			DeliveredAt: r.DeliveredAt,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
		if st, ok := studentMap[r.StudentID]; ok {
			v.StudentName = st.Name
			v.StudentNo = st.StudentNo
		}
		if it, ok := itemMap[r.ShopItemID]; ok {
			v.ItemName = it.Name
			v.CurrencyType = it.CurrencyType
			v.EventTag = it.EventTag
			v.Category = it.Category
		}
		out = append(out, v)
	}
	return out, nil
}

// CreateRedemption 教师代学生发起兑换（class_id 存学生所在班级，学校级商品亦如此）。
func (s *ShopService) CreateRedemption(u *models.User, studentID uint, item *models.ShopItem) (*models.ShopRedemption, error) {
	student, err := s.scope.StudentInScope(u, studentID)
	if err != nil {
		return nil, err
	}

	redemption := models.ShopRedemption{
		StudentID:  student.ID,
		ShopItemID: item.ID,
		ClassID:    student.ClassID,
		Cost:       item.CostScore,
		Status:     "pending",
	}
	if err := s.db.Create(&redemption).Error; err != nil {
		return nil, err
	}
	return &redemption, nil
}

// ApproveResult 兑换审批结果。
type ApproveResult struct {
	Message        string `json:"message"`
	RemainingScore int    `json:"remaining_score"`
	PetLevel       *int   `json:"pet_level"`
}

// consumeStockTx 消耗一件库存（仅对「限量商品」生效：stock > 0）。
//
// 口径（与前端既有的「0=∞」展示一致）：**stock == 0 视为不限量**，既不扣减也不下架——
// 系统内置的默认商品（免作业、冰淇淋等）都是 0，必须保持可无限兑换。
// stock > 0 才扣减，扣到 0 时**自动下架**（is_active=false）：所有购买路径都要求 is_active，
// 于是「售罄」自然表现为「买不到」，无需新增列。并发下用条件更新，不会超卖。
func (s *ShopService) consumeStockTx(tx *gorm.DB, itemID uint) error {
	var item models.ShopItem
	if err := tx.Select("stock").Where("id = ?", itemID).First(&item).Error; err != nil {
		return err
	}
	if item.Stock <= 0 {
		return nil // 不限量
	}

	res := tx.Model(&models.ShopItem{}).
		Where("id = ? AND stock > 0", itemID).
		Update("stock", gorm.Expr("stock - 1"))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrBadRequest("商品库存不足")
	}

	// 售罄 → 自动下架（0 在约定里表示「不限量」，故必须同时下架，否则会被当成无限库存）。
	return tx.Model(&models.ShopItem{}).
		Where("id = ? AND stock <= 0", itemID).
		Update("is_active", false).Error
}

// ApproveRedemption 审批通过并结算（状态机同 Laravel approveRedemption，但**并发与原子性**更强）。
//
// ⚠️ 与 Laravel 的有意差异（安全性修复）：
//  1. 「认领状态」用条件更新（WHERE status='pending'）并与结算放进**同一个事务**。
//     原实现（Laravel 同款）是「事务外读 status → 各自结算 → 再更新 status」，两处都会出事：
//     并发双击审批会各扣一次款；结算成功但状态更新失败则记录仍是 pending，重试再扣一次。
//  2. 结算改用事务内版本（spendScoreTx / exchangeTx / spendTx），因此失败会整体回滚，
//     状态退回 pending，不会留下「扣了款但记录未落地」。
func (s *ShopService) ApproveRedemption(u *models.User, id uint) (*ApproveResult, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}

	var redemption models.ShopRedemption
	err = s.db.Where("id = ? AND class_id IN ?", id, classIDs).First(&redemption).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("兑换记录不存在或不可见")
	}
	if err != nil {
		return nil, err
	}
	if redemption.Status != "pending" {
		return nil, ErrBadRequest("该兑换已处理")
	}

	var student models.Student
	if err := s.db.First(&student, redemption.StudentID).Error; err != nil {
		return nil, err
	}
	var item models.ShopItem
	if err := s.db.First(&item, redemption.ShopItemID).Error; err != nil {
		return nil, err
	}
	// ⚠️ 已下架的商品不再结算（含「售罄自动下架」）：否则卖出前就存在的 pending 记录会绕过
	// 库存判定——因为 stock == 0 在约定里表示「不限量」，无法与「售罄」区分，只能靠 is_active。
	// 教师若确需结算，可先把商品重新上架（必要时补库存）再审批。
	if !item.IsActive {
		return nil, ErrBadRequest("该商品已下架，无法结算")
	}

	itemName := item.Name
	if itemName == "" {
		itemName = "未知物品"
	}
	cost := redemption.Cost
	currency := item.CurrencyType
	if currency == "" {
		currency = "score"
	}

	// 积分充值类要走汇率：默认汇率播种必须在事务外（它自己写 exchange_rates，事务内用另一条
	// 连接写会与写锁互等）。
	var schoolID uint
	useExchange := item.Category == "points" && isWalletCurrency(currency)
	if useExchange {
		var class models.ClassRoom
		if err := s.db.First(&class, student.ClassID).Error; err != nil {
			return nil, err
		}
		if err := s.currency.ensureDefaultRates(class.SchoolID); err != nil {
			return nil, err
		}
		schoolID = class.SchoolID
	}

	var result *ApproveResult
	var spent *models.Score
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1) 原子认领：抢到的才结算；抢不到说明已被处理（并发下只有一次生效）。
		claim := tx.Model(&models.ShopRedemption{}).
			Where("id = ? AND class_id IN ? AND status = ?", id, classIDs, "pending").
			Updates(map[string]any{
				"status":      "approved",
				"approved_by": u.ID,
				"approved_at": util.Now(),
			})
		if claim.Error != nil {
			return claim.Error
		}
		if claim.RowsAffected == 0 {
			return ErrBadRequest("该兑换已处理")
		}

		// 2) 结算（与认领同一事务）：先按限量扣库存，再扣款/发货。
		if err := s.consumeStockTx(tx, item.ID); err != nil {
			return err
		}
		switch {
		case useExchange:
			// 积分充值类：扣积分 + 扣宠物经验 → 按汇率发放钱包币。
			if _, err := s.currency.exchangeTx(tx, &student, schoolID, currency, cost, u.ID); err != nil {
				return err
			}
		case currency == "score":
			// 积分兑换：扣积分 + 扣宠物经验。
			sc, err := spendScoreTx(tx, &student, cost, itemName, u.ID)
			if err != nil {
				return err
			}
			spent = &sc
		default:
			// 钱包币兑换：只扣钱包余额。
			if err := s.currency.spendTx(tx, student.ID, currency, cost); err != nil {
				return err
			}
		}

		// 3) 特权奖励自动发班级事件通知（已发布）。
		if item.Category == "privilege" {
			now := util.Now()
			notice := models.Notice{
				ClassID:     student.ClassID,
				SchoolID:    u.SchoolID,
				Title:       "特权奖励：" + itemName,
				Content:     fmt.Sprintf("%s 使用 %d %s 兑换了「%s」", student.Name, cost, currencyUnit(currency), itemName),
				Type:        "event",
				PublishedBy: u.ID,
				IsPublished: true,
				PublishedAt: &now,
			}
			if err := tx.Create(&notice).Error; err != nil {
				return err
			}
		}

		// 4) 回读权威余额与宠物等级。
		if err := tx.First(&student, student.ID).Error; err != nil {
			return err
		}
		var pet models.Pet
		var petLevel *int
		if err := tx.Where("student_id = ?", student.ID).First(&pet).Error; err == nil {
			petLevel = &pet.Level
		}

		result = &ApproveResult{
			Message:        fmt.Sprintf("已批准兑换，扣除 %d %s", cost, currencyUnit(currency)),
			RemainingScore: student.TotalScore,
			PetLevel:       petLevel,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 事务提交后再发大屏事件（与 ScoreService.SpendScore 同一口径：提交后发布）。
	if spent != nil {
		s.scores.publishScoreUpdate(&student, -cost, spent.Reason, true)
	}

	return result, nil
}

// RejectRedemption 拒绝兑换（不结算）。
//
// ⚠️ 守卫（有意差异）：只允许 pending → rejected，且用条件更新原子完成。原先没有任何状态判定
// （Laravel 同款），已 approve（已扣款）的记录也能被改成 rejected 且不退款——学生钱没了、
// 记录却显示拒绝。
func (s *ShopService) RejectRedemption(u *models.User, id uint) error {
	redemption, err := s.findRedemption(u, id)
	if err != nil {
		return err
	}
	res := s.db.Model(&models.ShopRedemption{}).
		Where("id = ? AND status = ?", redemption.ID, "pending").
		Update("status", "rejected")
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrBadRequest("该兑换已处理，不能拒绝")
	}
	return nil
}

// DeliverRedemption 标记为已发放。
//
// ⚠️ 守卫（有意差异）：只允许 approved → delivered。原先没有任何状态判定（Laravel 同款），
// pending（还没扣款）的记录也能直接置为 delivered → 学生「免单」拿到商品，且此后 approve
// 会因 status≠pending 被永久拒绝、无法补扣。顺带补写此前从未落库的 delivered_at。
func (s *ShopService) DeliverRedemption(u *models.User, id uint) error {
	redemption, err := s.findRedemption(u, id)
	if err != nil {
		return err
	}
	res := s.db.Model(&models.ShopRedemption{}).
		Where("id = ? AND status = ?", redemption.ID, "approved").
		Updates(map[string]any{"status": "delivered", "delivered_at": util.Now()})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrBadRequest("该兑换尚未批准，不能标记发放")
	}
	return nil
}

func (s *ShopService) findRedemption(u *models.User, id uint) (*models.ShopRedemption, error) {
	classIDs, err := s.scope.ClassIDs(u)
	if err != nil {
		return nil, err
	}
	var redemption models.ShopRedemption
	err = s.db.Where("id = ? AND class_id IN ?", id, classIDs).First(&redemption).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("兑换记录不存在或不可见")
	}
	if err != nil {
		return nil, err
	}
	return &redemption, nil
}

// isWalletCurrency 判断是否为支持的钱包币种。
func isWalletCurrency(code string) bool {
	_, ok := WalletCurrencies[code]
	return ok
}

// currencyUnit 返回结算文案中的货币单位（积分或钱包币中文名）。
func currencyUnit(currency string) string {
	if currency == "score" {
		return "积分"
	}
	return CurrencyLabel(currency)
}

// defaultItem 默认商品模板项。
type defaultItem struct {
	Name         string
	Description  string
	Category     string
	CostScore    int
	CurrencyType string
}

// defaultItems 首次访问播种的学校级默认商品（22 条，忠实移植自 ShopService::seedDefaultItems）。
var defaultItems = []defaultItem{
	// 积分充值类（按 2:1 汇率）
	{"班级积分 +10", "兑换 10 班级积分", "points", 20, "class_point"},
	{"科学币 +5", "兑换 5 科学币", "points", 10, "science"},
	{"读书币 +5", "兑换 5 读书币", "points", 10, "reading"},
	{"体育币 +5", "兑换 5 体育币", "points", 10, "class_point"},
	// 小商品 ≈100
	{"铅笔", "标准 HB 铅笔一支", "stationery", 100, "score"},
	{"橡皮擦", "4B 橡皮擦一块", "stationery", 100, "score"},
	{"草稿纸", "A4 草稿纸 10 张", "stationery", 100, "score"},
	{"免罚站一次", "免除一次罚站", "privilege", 100, "score"},
	{"免罚跑步一次", "免除一次罚跑步", "privilege", 100, "score"},
	// 中商品 120~150
	{"便利贴", "彩色便利贴一本", "stationery", 120, "score"},
	{"黑色圆珠笔", "0.5mm 黑色圆珠笔一支", "stationery", 150, "score"},
	{"蓝色圆珠笔", "0.5mm 蓝色圆珠笔一支", "stationery", 150, "score"},
	{"红色圆珠笔", "红色批改用笔一支", "stationery", 150, "score"},
	{"香蕉", "新鲜香蕉一根 🍌（请勿乱扔果皮）", "food", 150, "score"},
	{"免做卫生一次", "免除一次值日卫生", "privilege", 150, "score"},
	// 大商品 ≈200
	{"练习本", "方格练习本一本", "stationery", 180, "score"},
	{"苹果", "新鲜苹果一个 🍎（请勿乱扔果皮）", "food", 180, "score"},
	{"饮料", "矿泉水/饮料一瓶 🧃", "food", 200, "score"},
	{"牛奶", "纯牛奶一盒 🥛", "food", 200, "score"},
	{"集体观影", "全班集体观影一次", "activity", 200, "score"},
	{"免作业一次", "免交一次作业", "privilege", 200, "score"},
	{"3D打印作品", "3D 打印小作品一件", "physical", 200, "score"},
}

// seedDefaultItems 幂等播种学校级默认商品（仅当该校尚无学校级商品时）。
func (s *ShopService) seedDefaultItems(schoolID uint) error {
	var count int64
	if err := s.db.Model(&models.ShopItem{}).
		Where("school_id = ? AND class_id IS NULL", schoolID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, d := range defaultItems {
			item := models.ShopItem{
				ClassID:      nil,
				SchoolID:     schoolID,
				Name:         d.Name,
				Description:  d.Description,
				Category:     d.Category,
				CostScore:    d.CostScore,
				CurrencyType: d.CurrencyType,
				Stock:        0,
				IsActive:     true,
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ============================================================
// 管理员端：全校商品（学校级 class_id = null = 同步所有班级；班级级 = 仅该班级）
// 移植自 SchoolAdminController::adminListShopItems 等四个动作。
// ============================================================

// SchoolItemInput 管理员端商品输入（IsActive 为 nil 时默认 true）。
type SchoolItemInput struct {
	Name         string
	Description  string
	Category     string
	CostScore    int
	CurrencyType string
	EventTag     string
	Stock        int
	ImagePath    string
	IsActive     *bool
}

// AdminItemView 管理员端商品视图（附班级名与 scope 标注）。
type AdminItemView struct {
	ID           uint      `json:"id"`
	ClassID      *uint     `json:"class_id"`
	SchoolID     uint      `json:"school_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	CostScore    int       `json:"cost_score"`
	CurrencyType string    `json:"currency_type"`
	EventTag     string    `json:"event_tag"`
	Stock        int       `json:"stock"`
	ImagePath    string    `json:"image_path"`
	IsActive     bool      `json:"is_active"`
	ClassName    *string   `json:"class_name"`
	Scope        string    `json:"scope"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ItemsForSchool 管理员端商品列表：本校全部商品（学校级 + 各班级级），可按币种过滤。
func (s *ShopService) ItemsForSchool(schoolID uint, currencyType string) ([]AdminItemView, error) {
	q := s.db.Where("school_id = ?", schoolID)
	if currencyType != "" {
		q = q.Where("currency_type = ?", currencyType)
	}

	var items []models.ShopItem
	if err := q.Order("category ASC, cost_score ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	classNames := map[uint]string{}
	classIDs := []uint{}
	for _, it := range items {
		if it.ClassID != nil {
			classIDs = append(classIDs, *it.ClassID)
		}
	}
	if len(classIDs) > 0 {
		var classes []models.ClassRoom
		if err := s.db.Where("id IN ?", classIDs).Find(&classes).Error; err != nil {
			return nil, err
		}
		for _, c := range classes {
			classNames[c.ID] = c.Name
		}
	}

	views := []AdminItemView{}
	for _, it := range items {
		v := AdminItemView{
			ID: it.ID, ClassID: it.ClassID, SchoolID: it.SchoolID,
			Name: it.Name, Description: it.Description, Category: it.Category,
			CostScore: it.CostScore, CurrencyType: it.CurrencyType, EventTag: it.EventTag,
			Stock: it.Stock, ImagePath: it.ImagePath, IsActive: it.IsActive,
			Scope: "school", CreatedAt: it.CreatedAt, UpdatedAt: it.UpdatedAt,
		}
		if it.ClassID != nil {
			v.Scope = "class"
			if name, ok := classNames[*it.ClassID]; ok {
				v.ClassName = &name
			} else {
				// 与 Laravel 一致：班级已不存在时用 #id 兜底。
				fallback := "#" + itoa(int(*it.ClassID))
				v.ClassName = &fallback
			}
		}
		views = append(views, v)
	}
	return views, nil
}

// CreateSchoolItem 创建学校级商品（class_id = null，同步所有班级）。
func (s *ShopService) CreateSchoolItem(schoolID uint, in SchoolItemInput) (*models.ShopItem, error) {
	category := in.Category
	if category == "" {
		category = "physical"
	}
	currencyType := in.CurrencyType
	if currencyType == "" {
		currencyType = "score"
	}
	isActive := true
	if in.IsActive != nil {
		isActive = *in.IsActive
	}

	item := models.ShopItem{
		ClassID:      nil,
		SchoolID:     schoolID,
		Name:         in.Name,
		Description:  in.Description,
		Category:     category,
		CostScore:    in.CostScore,
		CurrencyType: currencyType,
		EventTag:     in.EventTag,
		Stock:        in.Stock,
		ImagePath:    in.ImagePath,
		IsActive:     isActive,
	}
	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}
	// is_active 带 default:true，显式 false 会被 GORM 跳过，需创建后回写。
	if err := forceFalseBool(s.db, &models.ShopItem{}, item.ID, "is_active", isActive); err != nil {
		return nil, err
	}
	item.IsActive = isActive
	return &item, nil
}

// FindSchoolItem 取本校学校级商品（班级级或跨校 ID 一律 404）。
func (s *ShopService) FindSchoolItem(schoolID, id uint) (*models.ShopItem, error) {
	var item models.ShopItem
	err := s.db.Where("id = ? AND school_id = ? AND class_id IS NULL", id, schoolID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound("商品不存在或不可见")
		}
		return nil, err
	}
	return &item, nil
}

// ============================================================
// 教室端（班级码）商城：商品列表 + 快捷兑换
// 移植自 Laravel DisplayController::quickShopItems / quickRedeem。
// ============================================================

// DisplayShopItem 教室端商品条目（字段逐字同 Laravel quickShopItems 的 get([...])）。
type DisplayShopItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CostScore   int    `json:"cost_score"`
	Stock       int    `json:"stock"`
	Category    string `json:"category"`
}

// DisplayRedeemResult 教室端兑换结果（字段名逐字同 Laravel quickRedeem 的 data）。
type DisplayRedeemResult struct {
	StudentName string `json:"student_name"`
	ItemName    string `json:"item_name"`
	Cost        int    `json:"cost"`
	TotalScore  int    `json:"total_score"`
}

// DisplayItems 教室端商品列表：只返回「本班级级 + 在售」商品。
// 与教师端 Items 的口径不同（后者含学校级商品、且空表时会播种默认商品）：Laravel
// quickShopItems 只查 class_id = 本班 AND is_active = true，此处逐条保持一致。
func (s *ShopService) DisplayItems(classID uint) ([]DisplayShopItem, error) {
	var items []models.ShopItem
	if err := s.db.Where("class_id = ? AND is_active = ?", classID, true).
		Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	out := make([]DisplayShopItem, 0, len(items))
	for _, it := range items {
		out = append(out, DisplayShopItem{
			ID:          it.ID,
			Name:        it.Name,
			Description: it.Description,
			CostScore:   it.CostScore,
			Stock:       it.Stock,
			Category:    it.Category,
		})
	}
	return out, nil
}

// DisplayRedeem 教室端快捷兑换：立即扣分（ScoreService.SpendScore，含审计日志 + score_update 事件）
// 并生成一条 status = approved 的兑换记录。
//
// 有意差异/忠实点：Laravel quickRedeem **不校验也不扣减库存**（stock 只做展示，与教师端审批路径
// 的结算口径不同），此处保持一致。
// ⚠️ 扣分原因与教师端审批路径**统一为「兑换消耗：<商品名>」**：原先这里传 "兑换："+商品名，
// 与 SpendScore 自带的前缀叠成「兑换消耗：兑换：<商品名>」（Laravel 同款）——同一个动作在
// 两个入口留下两种审计文案，属无意义的不一致，已统一。
func (s *ShopService) DisplayRedeem(classID, studentID, itemID uint) (*DisplayRedeemResult, error) {
	var student models.Student
	err := s.db.Where("class_id = ? AND id = ?", classID, studentID).First(&student).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("学生或商品不存在")
	}
	if err != nil {
		return nil, err
	}

	var item models.ShopItem
	err = s.db.Where("id = ? AND class_id = ? AND is_active = ?", itemID, classID, true).
		First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("学生或商品不存在")
	}
	if err != nil {
		return nil, err
	}

	if student.TotalScore < item.CostScore {
		return nil, ErrBadRequest("积分不足")
	}

	// 操作人：班级教师，缺失时兜底 user_id = 1（Laravel `$teacherId ?: 1`）。
	operator := classTeacherID(s.db, classID)
	if operator == 0 {
		operator = 1
	}

	// ⚠️ 「扣款 + 落兑换记录」必须同一事务：原先是两笔独立提交，落库失败会留下「已扣款但无
	// 兑换记录」的孤儿扣款（Laravel quickRedeem 同款问题）。
	now := util.Now()
	approvedBy := operator
	var spent models.Score
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 限量商品先扣库存（0 = 不限量；扣到 0 自动下架）。失败即整笔回滚，不会扣分。
		if err := s.consumeStockTx(tx, item.ID); err != nil {
			return err
		}
		var err error
		spent, err = spendScoreTx(tx, &student, item.CostScore, item.Name, operator)
		if err != nil {
			return err
		}
		return tx.Create(&models.ShopRedemption{
			StudentID:  student.ID,
			ShopItemID: item.ID,
			ClassID:    classID,
			Cost:       item.CostScore,
			Status:     "approved",
			ApprovedBy: &approvedBy,
			ApprovedAt: &now,
		}).Error
	})
	if err != nil {
		return nil, err
	}

	// 事务提交后发大屏事件（同 ScoreService.SpendScore 的口径）。
	s.scores.publishScoreUpdate(&student, -item.CostScore, spent.Reason, true)

	return &DisplayRedeemResult{
		StudentName: student.Name,
		ItemName:    item.Name,
		Cost:        item.CostScore,
		TotalScore:  student.TotalScore,
	}, nil
}
