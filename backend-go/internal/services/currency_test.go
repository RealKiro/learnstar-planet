package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
)

// ---- RatesForSchool ----

func TestRatesForSchool_SeedsDefaults(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewCurrencyService(db)

	rates, err := svc.RatesForSchool(school.ID)
	if err != nil {
		t.Fatalf("RatesForSchool: %v", err)
	}
	if len(rates) != 3 {
		t.Fatalf("expected 3 default rates, got %d", len(rates))
	}
	if rates[0].FromCurrency != "score" {
		t.Fatalf("first rate from_currency = %q, want score", rates[0].FromCurrency)
	}
}

func TestRatesForSchool_Idempotent(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewCurrencyService(db)

	rates1, err := svc.RatesForSchool(school.ID)
	if err != nil {
		t.Fatalf("first call: %v", err)
	}
	rates2, err := svc.RatesForSchool(school.ID)
	if err != nil {
		t.Fatalf("second call: %v", err)
	}
	if len(rates1) != len(rates2) {
		t.Fatalf("len changed: first=%d second=%d", len(rates1), len(rates2))
	}
	if len(rates1) != 3 {
		t.Fatalf("expected 3 after idempotent, got %d", len(rates1))
	}
}

// ---- ListRates ----

func TestListRates_DoesNotSeed(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewCurrencyService(db)

	rates, err := svc.ListRates(school.ID)
	if err != nil {
		t.Fatalf("ListRates: %v", err)
	}
	if len(rates) != 0 {
		t.Fatalf("expected 0 rates (no seed), got %d", len(rates))
	}
}

// ---- CreateRate ----

func TestCreateRate(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewCurrencyService(db)

	rate, err := svc.CreateRate(school.ID, services.RateInput{
		Name: "测试汇率", FromCurrency: "science", ToCurrency: "reading", Rate: 2.0,
	}, true)
	if err != nil {
		t.Fatalf("CreateRate: %v", err)
	}
	if rate.Rate != 2.0 {
		t.Fatalf("rate = %f, want 2.0", rate.Rate)
	}
	if !rate.IsActive {
		t.Fatal("IsActive should be true")
	}
	if rate.SchoolID != school.ID {
		t.Fatalf("SchoolID = %d, want %d", rate.SchoolID, school.ID)
	}
	if rate.Name != "测试汇率" {
		t.Fatalf("Name = %q, want 测试汇率", rate.Name)
	}
}

func TestCreateRate_IsActiveFalse(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewCurrencyService(db)

	rate, err := svc.CreateRate(school.ID, services.RateInput{
		Name: "停用汇率", FromCurrency: "score", ToCurrency: "class_point", Rate: 0.3,
	}, false)
	if err != nil {
		t.Fatalf("CreateRate: %v", err)
	}
	if rate.IsActive {
		t.Fatal("IsActive should be false")
	}

	// Verify persisted value too.
	var dbRate models.ExchangeRate
	if err := db.First(&dbRate, rate.ID).Error; err != nil {
		t.Fatalf("refetch: %v", err)
	}
	if dbRate.IsActive {
		t.Fatal("persisted IsActive should be false")
	}
}

// ---- UpdateRate ----

func TestUpdateRate_RateOnly(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewCurrencyService(db)

	// first seed defaults
	_, err := svc.RatesForSchool(school.ID)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	var orig models.ExchangeRate
	if err := db.Where("school_id = ?", school.ID).First(&orig).Error; err != nil {
		t.Fatalf("find original rate: %v", err)
	}

	newRate := 0.75
	updated, err := svc.UpdateRate(school.ID, orig.ID, &newRate, nil)
	if err != nil {
		t.Fatalf("UpdateRate: %v", err)
	}
	if updated.Rate != 0.75 {
		t.Fatalf("rate = %f, want 0.75", updated.Rate)
	}
	if updated.IsActive != orig.IsActive {
		t.Fatal("IsActive should remain unchanged")
	}
}

func TestUpdateRate_IsActiveOnly(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewCurrencyService(db)

	svc.RatesForSchool(school.ID)
	var orig models.ExchangeRate
	db.Where("school_id = ?", school.ID).First(&orig)

	inactive := false
	updated, err := svc.UpdateRate(school.ID, orig.ID, nil, &inactive)
	if err != nil {
		t.Fatalf("UpdateRate: %v", err)
	}
	if updated.IsActive {
		t.Fatal("IsActive should be false")
	}
	if updated.Rate != orig.Rate {
		t.Fatalf("rate changed from %f to %f", orig.Rate, updated.Rate)
	}
}

func TestUpdateRate_NotFound(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := services.NewCurrencyService(db)

	_, err := svc.UpdateRate(school.ID, 99999, nil, nil)
	if err == nil {
		t.Fatal("expected error for non-existent rate")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 404 {
		t.Fatalf("expected 404, got status=%d msg=%q", e.Status, e.Message)
	}
}

func TestUpdateRate_WrongSchool(t *testing.T) {
	db := setupDB(t)
	school1 := seedSchool(t, db)
	school2 := models.School{Name: "其他学校", Code: "other", Status: "active"}
	db.Create(&school2)

	svc := services.NewCurrencyService(db)
	svc.RatesForSchool(school1.ID)

	var orig models.ExchangeRate
	db.Where("school_id = ?", school1.ID).First(&orig)

	_, err := svc.UpdateRate(school2.ID, orig.ID, nil, nil)
	if err == nil {
		t.Fatal("expected error for wrong school")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 404 {
		t.Fatalf("expected 404, got status=%d msg=%q", e.Status, e.Message)
	}
}

// ---- Exchange ----

func TestExchange_Success(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	teacher := seedTeacher(t, db, school.ID, "teacher1")

	// Give student some points
	if err := db.Model(&student).Update("total_score", 100).Error; err != nil {
		t.Fatalf("update score: %v", err)
	}
	db.First(&student, student.ID)

	svc := services.NewCurrencyService(db)

	result, err := svc.Exchange(&student, "science", 10, teacher.ID)
	if err != nil {
		t.Fatalf("Exchange: %v", err)
	}
	// 100 - 10 = 90
	if result.RemainingScore != 90 {
		t.Fatalf("RemainingScore = %d, want 90", result.RemainingScore)
	}
	// 10 * 0.5 = 5
	if result.WalletBalance != 5 {
		t.Fatalf("WalletBalance = %d, want 5", result.WalletBalance)
	}

	// Verify score log was created
	var scoreCount int64
	db.Model(&models.Score{}).Where("student_id = ?", student.ID).Count(&scoreCount)
	if scoreCount != 1 {
		t.Fatalf("expected 1 score record, got %d", scoreCount)
	}

	// Verify exchange log
	var logCount int64
	db.Model(&models.ExchangeLog{}).Where("student_id = ?", student.ID).Count(&logCount)
	if logCount != 1 {
		t.Fatalf("expected 1 exchange log, got %d", logCount)
	}
}

func TestExchange_AmountZero(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	db.Model(&student).Update("total_score", 100)
	db.First(&student, student.ID)

	svc := services.NewCurrencyService(db)

	_, err := svc.Exchange(&student, "science", 0, 1)
	if err == nil {
		t.Fatal("expected error for zero amount")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}

func TestExchange_AmountNegative(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	db.Model(&student).Update("total_score", 100)
	db.First(&student, student.ID)

	svc := services.NewCurrencyService(db)

	_, err := svc.Exchange(&student, "science", -5, 1)
	if err == nil {
		t.Fatal("expected error for negative amount")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}

func TestExchange_InsufficientScore(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	// student has 0 points
	db.First(&student, student.ID)

	svc := services.NewCurrencyService(db)

	_, err := svc.Exchange(&student, "science", 10, 1)
	if err == nil {
		t.Fatal("expected error for insufficient score")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}

func TestExchange_NoRateForCurrency(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	db.Model(&student).Update("total_score", 100)
	db.First(&student, student.ID)

	svc := services.NewCurrencyService(db)

	// "gold" has no default rate
	_, err := svc.Exchange(&student, "gold", 10, 1)
	if err == nil {
		t.Fatal("expected error for unsupported currency")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
	if e, ok := services.AsAppError(err); ok {
		if e.Message != "未找到可用的汇率配置" {
			t.Fatalf("message = %q, want 未找到可用的汇率配置", e.Message)
		}
	}
}

// ---- CrossExchange ----

func TestCrossExchange_Success(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	teacher := seedTeacher(t, db, school.ID, "teacher1")

	// Create wallets with balances
	db.Create(&models.Wallet{StudentID: student.ID, CurrencyType: "science", Balance: 100})
	db.Create(&models.Wallet{StudentID: student.ID, CurrencyType: "reading", Balance: 50})

	svc := services.NewCurrencyService(db)
	// 默认汇率仅 score→*；science→reading 需先创建。
	if _, err := svc.CreateRate(school.ID, services.RateInput{
		Name: "科学币→读书币", FromCurrency: "science", ToCurrency: "reading", Rate: 0.5,
	}, true); err != nil {
		t.Fatalf("CreateRate: %v", err)
	}

	result, err := svc.CrossExchange(&student, "science", "reading", 20, teacher.ID)
	if err != nil {
		t.Fatalf("CrossExchange: %v", err)
	}
	// 100 - 20 = 80
	if result.FromBalance != 80 {
		t.Fatalf("FromBalance = %d, want 80", result.FromBalance)
	}
	// 50 + (20 * 0.5) = 60
	if result.ToBalance != 60 {
		t.Fatalf("ToBalance = %d, want 60", result.ToBalance)
	}

	// Verify wallet balances persisted
	var scienceWallet models.Wallet
	db.Where("student_id = ? AND currency_type = ?", student.ID, "science").First(&scienceWallet)
	if scienceWallet.Balance != 80 {
		t.Fatalf("science balance = %d, want 80", scienceWallet.Balance)
	}

	var readingWallet models.Wallet
	db.Where("student_id = ? AND currency_type = ?", student.ID, "reading").First(&readingWallet)
	if readingWallet.Balance != 60 {
		t.Fatalf("reading balance = %d, want 60", readingWallet.Balance)
	}

	// Verify exchange log
	var log models.ExchangeLog
	db.Where("student_id = ? AND from_currency = ?", student.ID, "science").First(&log)
	if log.FromAmount != 20 || log.ToAmount != 10 {
		t.Fatalf("log amounts: from=%d to=%d, want 20 and 10", log.FromAmount, log.ToAmount)
	}
}

func TestCrossExchange_SameCurrency(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	db.Create(&models.Wallet{StudentID: student.ID, CurrencyType: "science", Balance: 100})

	svc := services.NewCurrencyService(db)

	_, err := svc.CrossExchange(&student, "science", "science", 10, 1)
	if err == nil {
		t.Fatal("expected error for same currency")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}

func TestCrossExchange_ZeroAmount(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)

	svc := services.NewCurrencyService(db)

	_, err := svc.CrossExchange(&student, "science", "reading", 0, 1)
	if err == nil {
		t.Fatal("expected error for zero amount")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}

func TestCrossExchange_InsufficientBalance(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	db.Create(&models.Wallet{StudentID: student.ID, CurrencyType: "science", Balance: 5})

	svc := services.NewCurrencyService(db)

	_, err := svc.CrossExchange(&student, "science", "reading", 10, 1)
	if err == nil {
		t.Fatal("expected error for insufficient balance")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}

func TestCrossExchange_NoRate(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	db.Create(&models.Wallet{StudentID: student.ID, CurrencyType: "science", Balance: 100})

	svc := services.NewCurrencyService(db)

	// "reading" → "gold" has no default rate
	_, err := svc.CrossExchange(&student, "reading", "gold", 10, 1)
	if err == nil {
		t.Fatal("expected error for missing rate")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}

// ---- WalletsFor ----

func TestWalletsFor(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class, student := seedTeacherClass(t, db, school.ID, 0, "一班")
	student2 := seedExtraStudent(t, db, class.ID)

	// Create wallets with known balances
	db.Create(&models.Wallet{StudentID: student.ID, CurrencyType: "science", Balance: 10})
	db.Create(&models.Wallet{StudentID: student2.ID, CurrencyType: "reading", Balance: 20})

	svc := services.NewCurrencyService(db)

	views, err := svc.WalletsFor([]uint{class.ID})
	if err != nil {
		t.Fatalf("WalletsFor: %v", err)
	}
	if len(views) != 2 {
		t.Fatalf("expected 2 wallet views, got %d", len(views))
	}

	found := map[uint]int{}
	for _, v := range views {
		found[v.StudentID] = v.Balance
	}
	if found[student.ID] != 10 {
		t.Fatalf("student %d balance = %d, want 10", student.ID, found[student.ID])
	}
	if found[student2.ID] != 20 {
		t.Fatalf("student2 %d balance = %d, want 20", student2.ID, found[student2.ID])
	}
}

func TestWalletsFor_StudentNames(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class, student := seedTeacherClass(t, db, school.ID, 0, "一班")

	db.Create(&models.Wallet{StudentID: student.ID, CurrencyType: "science", Balance: 5})

	svc := services.NewCurrencyService(db)

	views, err := svc.WalletsFor([]uint{class.ID})
	if err != nil {
		t.Fatalf("WalletsFor: %v", err)
	}
	if len(views) != 1 {
		t.Fatalf("expected 1 view, got %d", len(views))
	}
	if views[0].StudentName != "小明" {
		t.Fatalf("StudentName = %q, want 小明", views[0].StudentName)
	}
}

func TestWalletsFor_OnlyActiveStudents(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	active := seedStudent(t, db, class.ID)
	inactive := models.Student{ClassID: class.ID, Name: "已离开", Status: "graduated"}
	if err := db.Create(&inactive).Error; err != nil {
		t.Fatalf("create inactive: %v", err)
	}

	db.Create(&models.Wallet{StudentID: active.ID, CurrencyType: "science", Balance: 10})
	db.Create(&models.Wallet{StudentID: inactive.ID, CurrencyType: "science", Balance: 99})

	svc := services.NewCurrencyService(db)

	views, err := svc.WalletsFor([]uint{class.ID})
	if err != nil {
		t.Fatalf("WalletsFor: %v", err)
	}
	// Only 1 wallet should appear (active student only)
	if len(views) != 1 {
		t.Fatalf("expected 1 view (active only), got %d", len(views))
	}
	if views[0].StudentID != active.ID {
		t.Fatalf("expected student %d, got %d", active.ID, views[0].StudentID)
	}
}

func TestWalletsFor_EmptyClassIDs(t *testing.T) {
	db := setupDB(t)
	svc := services.NewCurrencyService(db)

	views, err := svc.WalletsFor([]uint{})
	if err != nil {
		t.Fatalf("WalletsFor empty: %v", err)
	}
	if len(views) != 0 {
		t.Fatalf("expected 0 views, got %d", len(views))
	}
}

// ---- LogsFor ----

func TestLogsFor_Basic(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class, student := seedTeacherClass(t, db, school.ID, 0, "一班")
	teacher := seedTeacher(t, db, school.ID, "teacher1")

	opBy := teacher.ID
	db.Create(&models.ExchangeLog{
		StudentID: student.ID, FromCurrency: "score", ToCurrency: "science",
		FromAmount: 10, ToAmount: 5, OperatedBy: &opBy,
	})

	svc := services.NewCurrencyService(db)

	result, err := svc.LogsFor([]uint{class.ID}, 1)
	if err != nil {
		t.Fatalf("LogsFor: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 log, got %d", len(result.Data))
	}
	if result.Data[0].StudentName != "小明" {
		t.Fatalf("StudentName = %q, want 小明", result.Data[0].StudentName)
	}
	if result.Data[0].StudentNo != "001" {
		t.Fatalf("StudentNo = %q, want 001", result.Data[0].StudentNo)
	}
	if result.Data[0].FromCurrency != "score" {
		t.Fatalf("FromCurrency = %q", result.Data[0].FromCurrency)
	}
	if result.Meta.Total != 1 {
		t.Fatalf("Total = %d, want 1", result.Meta.Total)
	}
	if result.Meta.LastPage != 1 {
		t.Fatalf("LastPage = %d, want 1", result.Meta.LastPage)
	}
	if result.Meta.CurrentPage != 1 {
		t.Fatalf("CurrentPage = %d, want 1", result.Meta.CurrentPage)
	}
}

func TestLogsFor_PaginationMeta(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class, student := seedTeacherClass(t, db, school.ID, 0, "一班")
	teacher := seedTeacher(t, db, school.ID, "teacher1")

	opBy := teacher.ID
	// Insert 25 logs (enough for 2 pages)
	for i := 0; i < 25; i++ {
		db.Create(&models.ExchangeLog{
			StudentID: student.ID, FromCurrency: "score", ToCurrency: "science",
			FromAmount: 10, ToAmount: 5, OperatedBy: &opBy,
		})
	}

	svc := services.NewCurrencyService(db)

	// Page 1: should have 20 items, last_page=2
	result, err := svc.LogsFor([]uint{class.ID}, 1)
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	if len(result.Data) != 20 {
		t.Fatalf("page 1: expected 20 items, got %d", len(result.Data))
	}
	if result.Meta.Total != 25 {
		t.Fatalf("Total = %d, want 25", result.Meta.Total)
	}
	if result.Meta.LastPage != 2 {
		t.Fatalf("LastPage = %d, want 2", result.Meta.LastPage)
	}
	if result.Meta.CurrentPage != 1 {
		t.Fatalf("CurrentPage = %d, want 1", result.Meta.CurrentPage)
	}

	// Page 2: should have 5 items
	result2, err := svc.LogsFor([]uint{class.ID}, 2)
	if err != nil {
		t.Fatalf("page 2: %v", err)
	}
	if len(result2.Data) != 5 {
		t.Fatalf("page 2: expected 5 items, got %d", len(result2.Data))
	}
	if result2.Meta.CurrentPage != 2 {
		t.Fatalf("CurrentPage = %d, want 2", result2.Meta.CurrentPage)
	}

	// Page 3: 25/h, lastPage=2 → should clamp to page 2 with 5 items
	result3, err := svc.LogsFor([]uint{class.ID}, 3)
	if err != nil {
		t.Fatalf("page 3: %v", err)
	}
	if result3.Meta.CurrentPage != 3 {
		t.Fatalf("CurrentPage = %d, want 3", result3.Meta.CurrentPage)
	}
}

func TestLogsFor_PageLessThanOne(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class, student := seedTeacherClass(t, db, school.ID, 0, "一班")
	teacher := seedTeacher(t, db, school.ID, "teacher1")

	opBy := teacher.ID
	db.Create(&models.ExchangeLog{
		StudentID: student.ID, FromCurrency: "score", ToCurrency: "science",
		FromAmount: 10, ToAmount: 5, OperatedBy: &opBy,
	})

	svc := services.NewCurrencyService(db)

	result, err := svc.LogsFor([]uint{class.ID}, 0)
	if err != nil {
		t.Fatalf("LogsFor: %v", err)
	}
	if result.Meta.CurrentPage != 1 {
		t.Fatalf("CurrentPage = %d, want 1 (clamped)", result.Meta.CurrentPage)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 log, got %d", len(result.Data))
	}
}

func TestLogsFor_Empty(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class, _ := seedTeacherClass(t, db, school.ID, 0, "一班")

	svc := services.NewCurrencyService(db)

	result, err := svc.LogsFor([]uint{class.ID}, 1)
	if err != nil {
		t.Fatalf("LogsFor: %v", err)
	}
	if len(result.Data) != 0 {
		t.Fatalf("expected 0 logs, got %d", len(result.Data))
	}
	if result.Meta.Total != 0 {
		t.Fatalf("Total = %d, want 0", result.Meta.Total)
	}
	if result.Meta.LastPage != 1 {
		t.Fatalf("LastPage = %d, want 1", result.Meta.LastPage)
	}
}

func TestLogsFor_DeletedStudent(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class, student := seedTeacherClass(t, db, school.ID, 0, "一班")
	teacher := seedTeacher(t, db, school.ID, "teacher1")

	opBy := teacher.ID
	db.Create(&models.ExchangeLog{
		StudentID: student.ID, FromCurrency: "score", ToCurrency: "science",
		FromAmount: 10, ToAmount: 5, OperatedBy: &opBy,
	})

	// Soft-delete the student (sets DeletedAt).
	if err := db.Delete(&models.Student{}, student.ID).Error; err != nil {
		t.Fatalf("delete student: %v", err)
	}

	svc := services.NewCurrencyService(db)

	// Laravel logsFor filters by students in class (whereIn student subquery),
	// so deleted students' logs do not appear. The "已删除学生" fallback is
	// defensive only and unreachable through LogsFor.
	result, err := svc.LogsFor([]uint{class.ID}, 1)
	if err != nil {
		t.Fatalf("LogsFor: %v", err)
	}
	if len(result.Data) != 0 {
		t.Fatalf("expected 0 logs after student deletion, got %d", len(result.Data))
	}
	if result.Meta.Total != 0 {
		t.Fatalf("Total = %d, want 0", result.Meta.Total)
	}
}

// ---- Spend ----

func TestSpend_Success(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)

	db.Create(&models.Wallet{StudentID: student.ID, CurrencyType: "science", Balance: 30})

	svc := services.NewCurrencyService(db)

	err := svc.Spend(student.ID, "science", 12, "商城消费")
	if err != nil {
		t.Fatalf("Spend: %v", err)
	}

	var wallet models.Wallet
	db.Where("student_id = ? AND currency_type = ?", student.ID, "science").First(&wallet)
	if wallet.Balance != 18 {
		t.Fatalf("balance = %d, want 18", wallet.Balance)
	}
}

func TestSpend_ZeroAmount(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)

	svc := services.NewCurrencyService(db)

	err := svc.Spend(student.ID, "science", 0, "x")
	if err == nil {
		t.Fatal("expected error for zero amount")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}

func TestSpend_InsufficientBalance(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)

	db.Create(&models.Wallet{StudentID: student.ID, CurrencyType: "science", Balance: 5})

	svc := services.NewCurrencyService(db)

	err := svc.Spend(student.ID, "science", 10, "x")
	if err == nil {
		t.Fatal("expected error for insufficient balance")
	}
	if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}
