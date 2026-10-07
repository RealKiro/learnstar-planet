// 币种服务：积分 → 钱包币兑换与钱包币消费（商城结算的最小依赖子集）。
package services

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"gorm.io/gorm"
)

// CurrencyService 多币种钱包与汇率结算。
type CurrencyService struct {
	db *gorm.DB
}

// NewCurrencyService 创建币种服务。
func NewCurrencyService(db *gorm.DB) *CurrencyService {
	return &CurrencyService{db: db}
}

// WalletCurrencies 支持的钱包币种（与 Laravel Wallet::currencies 一致）。
var WalletCurrencies = map[string]string{
	"science":     "科学币",
	"reading":     "读书币",
	"class_point": "班级积分",
}

// CurrencyLabel 返回币种中文名，未知币种回退为原始代码。
func CurrencyLabel(code string) string {
	if label, ok := WalletCurrencies[code]; ok {
		return label
	}
	return code
}

// defaultRates 积分 → 币种默认汇率（2:1 防通胀：2 积分 = 1 币）。
var defaultRates = []models.ExchangeRate{
	{Name: "积分 → 科学币", FromCurrency: "score", ToCurrency: "science", Rate: 0.5, IsActive: true},
	{Name: "积分 → 读书币", FromCurrency: "score", ToCurrency: "reading", Rate: 0.5, IsActive: true},
	{Name: "积分 → 体育币", FromCurrency: "score", ToCurrency: "class_point", Rate: 0.5, IsActive: true},
}

// ensureDefaultRates 惰性播种学校默认汇率（幂等，按 from+to 去重）。
func (s *CurrencyService) ensureDefaultRates(schoolID uint) error {
	for _, d := range defaultRates {
		var count int64
		if err := s.db.Model(&models.ExchangeRate{}).
			Where("school_id = ? AND from_currency = ? AND to_currency = ?",
				schoolID, d.FromCurrency, d.ToCurrency).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		rate := d
		rate.SchoolID = schoolID
		if err := s.db.Create(&rate).Error; err != nil {
			return err
		}
	}
	return nil
}

// getOrCreateWallet 获取或创建学生某币种钱包。
func (s *CurrencyService) getOrCreateWallet(tx *gorm.DB, studentID uint, currency string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := tx.Where("student_id = ? AND currency_type = ?", studentID, currency).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		wallet = models.Wallet{StudentID: studentID, CurrencyType: currency, Balance: 0}
		if err := tx.Create(&wallet).Error; err != nil {
			return nil, err
		}
		return &wallet, nil
	}
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// ExchangeResult 兑换结果。
type ExchangeResult struct {
	RemainingScore int `json:"remaining_score"`
	WalletBalance  int `json:"wallet_balance"`
}

// Exchange 扣减学生积分并按学校汇率发放钱包币（单事务原子提交）。
// 扣分、审计日志、宠物经验同步与旧 Laravel CurrencyService::exchange 一致。
func (s *CurrencyService) Exchange(student *models.Student, toCurrency string, scoreAmount int, operatedBy uint) (*ExchangeResult, error) {
	if scoreAmount <= 0 {
		return nil, ErrBadRequest("兑换积分必须大于 0")
	}
	if student.TotalScore < scoreAmount {
		return nil, ErrBadRequest("积分不足，当前余额：" + strconv.Itoa(student.TotalScore))
	}

	var class models.ClassRoom
	if err := s.db.First(&class, student.ClassID).Error; err != nil {
		return nil, err
	}
	if err := s.ensureDefaultRates(class.SchoolID); err != nil {
		return nil, err
	}

	var rate models.ExchangeRate
	err := s.db.Where("school_id = ? AND from_currency = ? AND to_currency = ? AND is_active = ?",
		class.SchoolID, "score", toCurrency, true).First(&rate).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBadRequest("未找到可用的汇率配置")
	}
	if err != nil {
		return nil, err
	}

	toAmount := int(math.Round(float64(scoreAmount) * rate.Rate))
	if toAmount <= 0 {
		return nil, ErrBadRequest("兑换金额过小，无法兑换")
	}

	result := &ExchangeResult{}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		balanceBefore := student.TotalScore

		// 1) 扣减积分。
		newBalance := balanceBefore - scoreAmount
		if err := tx.Model(&models.Student{}).Where("id = ?", student.ID).
			Update("total_score", newBalance).Error; err != nil {
			return err
		}

		// 2) 扣分记录。
		reason := "兑换" + CurrencyLabel(toCurrency)
		score := models.Score{
			StudentID: student.ID,
			ClassID:   student.ClassID,
			Amount:    -scoreAmount,
			Reason:    reason,
			GivenBy:   operatedBy,
		}
		if err := tx.Create(&score).Error; err != nil {
			return err
		}

		// 3) 审计日志。
		if err := tx.Create(&models.ScoreLog{
			StudentID:     student.ID,
			ScoreID:       score.ID,
			BalanceBefore: balanceBefore,
			BalanceAfter:  newBalance,
			Description:   reason,
		}).Error; err != nil {
			return err
		}

		// 4) 扣减宠物经验（1:1）。
		if err := syncPetForDelta(tx, student.ID, newBalance, -scoreAmount); err != nil {
			return err
		}

		// 5) 增加钱包余额（显式计算，避免依赖 GORM Update 回写内存字段导致重复累加）。
		wallet, err := s.getOrCreateWallet(tx, student.ID, toCurrency)
		if err != nil {
			return err
		}
		walletNewBalance := wallet.Balance + toAmount
		if err := tx.Model(&models.Wallet{}).Where("id = ?", wallet.ID).Update("balance", walletNewBalance).Error; err != nil {
			return err
		}
		wallet.Balance = walletNewBalance

		// 6) 兑换日志。
		if err := tx.Create(&models.ExchangeLog{
			StudentID:    student.ID,
			FromCurrency: "score",
			ToCurrency:   toCurrency,
			FromAmount:   scoreAmount,
			ToAmount:     toAmount,
			OperatedBy:   &operatedBy,
		}).Error; err != nil {
			return err
		}

		student.TotalScore = newBalance
		result.RemainingScore = newBalance
		result.WalletBalance = wallet.Balance
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Spend 消费钱包币种（用于活动专属商城结算）。
func (s *CurrencyService) Spend(studentID uint, currency string, amount int, reason string) error {
	if amount <= 0 {
		return ErrBadRequest("消费数量必须大于 0")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		wallet, err := s.getOrCreateWallet(tx, studentID, currency)
		if err != nil {
			return err
		}
		if wallet.Balance < amount {
			return ErrBadRequest(CurrencyLabel(currency) + "余额不足，当前余额：" + strconv.Itoa(wallet.Balance))
		}

		if err := tx.Model(wallet).Update("balance", wallet.Balance-amount).Error; err != nil {
			return err
		}

		return tx.Create(&models.ExchangeLog{
			StudentID:    studentID,
			FromCurrency: currency,
			ToCurrency:   currency,
			FromAmount:   amount,
			ToAmount:     amount,
			OperatedBy:   nil,
		}).Error
	})
}

// RatesForSchool 返回学校汇率列表；首次访问（为空）惰性播种默认汇率。
// 与 Laravel CurrencyService::ratesForSchool 一致：默认多条统一播种，列表按 from/to 升序。
func (s *CurrencyService) RatesForSchool(schoolID uint) ([]models.ExchangeRate, error) {
	if err := s.ensureDefaultRates(schoolID); err != nil {
		return nil, err
	}
	var rates []models.ExchangeRate
	err := s.db.Where("school_id = ?", schoolID).
		Order("from_currency ASC").Order("to_currency ASC").
		Find(&rates).Error
	return rates, err
}

// RateInput 新增汇率入参。
type RateInput struct {
	Name         string
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

// CreateRate 新增学校级汇率，isActive 由调用方决定（教师端恒 true，管理员端可取请求值）。
func (s *CurrencyService) CreateRate(schoolID uint, in RateInput, isActive bool) (*models.ExchangeRate, error) {
	rate := models.ExchangeRate{
		SchoolID:     schoolID,
		Name:         in.Name,
		FromCurrency: in.FromCurrency,
		ToCurrency:   in.ToCurrency,
		Rate:         in.Rate,
		IsActive:     isActive,
	}
	if err := s.db.Create(&rate).Error; err != nil {
		return nil, err
	}
	// GORM 对带 default:true 的零值字段会先落默认值（false 会变 true），
	// 故创建后针对 is_active=false 显式回写（SQL 直更新，不经过零值省略逻辑）。
	if !isActive {
		if err := s.db.Model(&models.ExchangeRate{}).Where("id = ?", rate.ID).
			Update("is_active", false).Error; err != nil {
			return nil, err
		}
		rate.IsActive = false
	}
	return &rate, nil
}

// ListRates 仅查询学校汇率（管理员端，不播种默认）。
func (s *CurrencyService) ListRates(schoolID uint) ([]models.ExchangeRate, error) {
	var rates []models.ExchangeRate
	err := s.db.Where("school_id = ?", schoolID).
		Order("from_currency ASC").Order("to_currency ASC").
		Find(&rates).Error
	return rates, err
}

// UpdateRate 更新学校级汇率（仅 rate/is_active），越权或不存在返回 404。
func (s *CurrencyService) UpdateRate(schoolID, id uint, rate *float64, isActive *bool) (*models.ExchangeRate, error) {
	var r models.ExchangeRate
	err := s.db.Where("id = ? AND school_id = ?", id, schoolID).First(&r).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound("汇率不存在或不属于本学校")
		}
		return nil, err
	}

	updates := map[string]any{}
	if rate != nil {
		updates["rate"] = *rate
	}
	if isActive != nil {
		updates["is_active"] = *isActive
	}

	if len(updates) > 0 {
		if err := s.db.Model(&r).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&r, r.ID).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// CrossExchangeResult 交叉兑换结果。
type CrossExchangeResult struct {
	FromBalance int `json:"from_balance"`
	ToBalance   int `json:"to_balance"`
}

// CrossExchange 币种间兑换（学校汇率），仅操作钱包余额，不触碰积分与宠物经验。
// 与 Laravel CurrencyService::crossExchange 语义一致。
func (s *CurrencyService) CrossExchange(student *models.Student, from, to string, amount int, operatedBy uint) (*CrossExchangeResult, error) {
	if amount <= 0 {
		return nil, ErrBadRequest("兑换数量必须大于 0")
	}
	if from == to {
		return nil, ErrBadRequest("源币种和目标币种不能相同")
	}

	var class models.ClassRoom
	if err := s.db.First(&class, student.ClassID).Error; err != nil {
		return nil, err
	}
	if err := s.ensureDefaultRates(class.SchoolID); err != nil {
		return nil, err
	}

	var rate models.ExchangeRate
	err := s.db.Where("school_id = ? AND from_currency = ? AND to_currency = ? AND is_active = ?",
		class.SchoolID, from, to, true).First(&rate).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBadRequest("未找到可用的汇率配置")
	}
	if err != nil {
		return nil, err
	}

	fromWallet, err := s.getOrCreateWallet(s.db, student.ID, from)
	if err != nil {
		return nil, err
	}
	if fromWallet.Balance < amount {
		return nil, ErrBadRequest(CurrencyLabel(from) + "余额不足，当前余额：" + strconv.Itoa(fromWallet.Balance))
	}

	toAmount := int(math.Round(float64(amount) * rate.Rate))
	if toAmount <= 0 {
		return nil, ErrBadRequest("兑换金额过小，无法兑换")
	}

	result := &CrossExchangeResult{}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 扣减源币种（以事务内最新余额为准）。
		newFrom := fromWallet.Balance - amount
		if err := tx.Model(&models.Wallet{}).Where("id = ?", fromWallet.ID).Update("balance", newFrom).Error; err != nil {
			return err
		}

		// 增加目标币种（显式计算，避免依赖 GORM Update 回写内存字段导致重复累加）。
		toWallet, err := s.getOrCreateWallet(tx, student.ID, to)
		if err != nil {
			return err
		}
		newTo := toWallet.Balance + toAmount
		if err := tx.Model(&models.Wallet{}).Where("id = ?", toWallet.ID).Update("balance", newTo).Error; err != nil {
			return err
		}

		// 兑换日志。
		if err := tx.Create(&models.ExchangeLog{
			StudentID:    student.ID,
			FromCurrency: from,
			ToCurrency:   to,
			FromAmount:   amount,
			ToAmount:     toAmount,
			OperatedBy:   &operatedBy,
		}).Error; err != nil {
			return err
		}

		result.FromBalance = newFrom
		result.ToBalance = newTo
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// WalletView 教师端钱包列表条目。
type WalletView struct {
	StudentID    uint   `json:"student_id"`
	StudentName  string `json:"student_name"`
	CurrencyType string `json:"currency_type"`
	Balance      int    `json:"balance"`
}

// WalletsFor 返回教师所带班级学生的全部钱包。
// classIDs 为教师可管辖班级集合；学生需为 active 状态。
func (s *CurrencyService) WalletsFor(classIDs []uint) ([]WalletView, error) {
	if len(classIDs) == 0 {
		return []WalletView{}, nil
	}

	var students []models.Student
	if err := s.db.Where("class_id IN ? AND status = ?", classIDs, "active").Find(&students).Error; err != nil {
		return nil, err
	}
	studentName := make(map[uint]string, len(students))
	for _, st := range students {
		studentName[st.ID] = st.Name
	}

	var wallets []models.Wallet
	if err := s.db.Where("student_id IN ?", studentIDsOf(students)).Find(&wallets).Error; err != nil {
		return nil, err
	}

	views := make([]WalletView, 0, len(wallets))
	for _, w := range wallets {
		views = append(views, WalletView{
			StudentID:    w.StudentID,
			StudentName:  studentName[w.StudentID],
			CurrencyType: w.CurrencyType,
			Balance:      w.Balance,
		})
	}
	return views, nil
}

// PageMeta 分页元信息。
type PageMeta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	Total       int `json:"total"`
}

// ExchangeLogView 兑换日志条目（含学生名与学号）。
type ExchangeLogView struct {
	ID           uint      `json:"id"`
	StudentID    uint      `json:"student_id"`
	StudentName  string    `json:"student_name"`
	StudentNo    string    `json:"student_no"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	FromAmount   int       `json:"from_amount"`
	ToAmount     int       `json:"to_amount"`
	OperatedBy   *uint     `json:"operated_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// LogsResult 兑换日志分页结果（Laravel paginate(20) 形态）。
type LogsResult struct {
	Data []ExchangeLogView `json:"data"`
	Meta PageMeta          `json:"meta"`
}

// LogsFor 返回教师所带班级学生的兑换日志（分页，每页 20，按时间倒序）。
func (s *CurrencyService) LogsFor(classIDs []uint, page int) (*LogsResult, error) {
	const perPage = 20
	if page < 1 {
		page = 1
	}

	var students []models.Student
	if err := s.db.Where("class_id IN ?", classIDs).Find(&students).Error; err != nil {
		return nil, err
	}
	ids := studentIDsOf(students)
	studentMeta := make(map[uint][2]string, len(students)) // id → {name, student_no}
	for _, st := range students {
		studentMeta[st.ID] = [2]string{st.Name, st.StudentNo}
	}

	base := s.db.Model(&models.ExchangeLog{}).Where("student_id IN ?", ids)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	lastPage := int((total + perPage - 1) / perPage)
	if lastPage < 1 {
		lastPage = 1
	}

	var logs []models.ExchangeLog
	if err := base.Order("created_at DESC").
		Offset((page - 1) * perPage).Limit(perPage).
		Find(&logs).Error; err != nil {
		return nil, err
	}

	views := make([]ExchangeLogView, 0, len(logs))
	for _, l := range logs {
		meta, ok := studentMeta[l.StudentID]
		name := "已删除学生"
		if ok {
			name = meta[0]
		}
		no := ""
		if ok {
			no = meta[1]
		}
		views = append(views, ExchangeLogView{
			ID:           l.ID,
			StudentID:    l.StudentID,
			StudentName:  name,
			StudentNo:    no,
			FromCurrency: l.FromCurrency,
			ToCurrency:   l.ToCurrency,
			FromAmount:   l.FromAmount,
			ToAmount:     l.ToAmount,
			OperatedBy:   l.OperatedBy,
			CreatedAt:    l.CreatedAt,
		})
	}

	return &LogsResult{
		Data: views,
		Meta: PageMeta{CurrentPage: page, LastPage: lastPage, Total: int(total)},
	}, nil
}

// studentIDsOf 提取学生 ID 列表。
func studentIDsOf(students []models.Student) []uint {
	ids := make([]uint, 0, len(students))
	for _, st := range students {
		ids = append(ids, st.ID)
	}
	return ids
}
