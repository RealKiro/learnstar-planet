package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
)

// 首次访问播种 22 条默认学校级商品，且幂等不重复播种。
func TestShopDefaultItemsSeededIdempotent(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	scope := services.NewScope(db)
	scores := services.NewScoreService(db)
	currency := services.NewCurrencyService(db)
	svc := services.NewShopService(db, scope, scores, currency)

	items, err := svc.Items(&teacher, "")
	if err != nil {
		t.Fatalf("items: %v", err)
	}
	if len(items) != 22 {
		t.Fatalf("default item count = %d, want 22", len(items))
	}

	// 再次访问不重复播种。
	again, err := svc.Items(&teacher, "")
	if err != nil {
		t.Fatalf("items second: %v", err)
	}
	if len(again) != 22 {
		t.Fatalf("duplicate seeding: second call returned %d, want 22", len(again))
	}
}

// 默认商品含积分充值类与积分兑换类。
func TestShopDefaultItemsContent(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewShopService(db, services.NewScope(db),
		services.NewScoreService(db), services.NewCurrencyService(db))

	items, err := svc.Items(&teacher, "")
	if err != nil {
		t.Fatalf("items: %v", err)
	}
	hasPoints := false
	allSchoolLevel := true
	for _, it := range items {
		if it.Category == "points" {
			hasPoints = true
		}
		if it.ClassID != nil {
			allSchoolLevel = false
		}
		if it.SchoolID != school.ID {
			t.Fatalf("item school_id = %d, want %d", it.SchoolID, school.ID)
		}
	}
	if !hasPoints {
		t.Fatal("default items should include points (积分充值) category")
	}
	if !allSchoolLevel {
		t.Fatal("default items should all be school-level (class_id null)")
	}
}

// Items 作用域：仅学校级 + 本班商品，别班班级级商品不可见。
func TestShopItemsScoped(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	t1 := seedTeacher(t, db, school.ID, "teacher1")
	t2 := seedTeacher(t, db, school.ID, "teacher2")
	class1, _ := seedTeacherClass(t, db, school.ID, t1.ID, "一班")
	seedTeacherClass(t, db, school.ID, t2.ID, "二班")

	svc := services.NewShopService(db, services.NewScope(db),
		services.NewScoreService(db), services.NewCurrencyService(db))

	// 触发播种后追加一个一班班级级商品。
	if _, err := svc.Items(&t1, ""); err != nil {
		t.Fatalf("seed: %v", err)
	}
	classItem := models.ShopItem{ClassID: &class1.ID, SchoolID: school.ID,
		Name: "班级专属", Category: "stationery", CostScore: 50, CurrencyType: "score", IsActive: true}
	if err := db.Create(&classItem).Error; err != nil {
		t.Fatalf("create class item: %v", err)
	}

	items, err := svc.Items(&t1, "")
	if err != nil {
		t.Fatalf("items t1: %v", err)
	}
	if len(items) != 23 {
		t.Fatalf("t1 sees %d items, want 23 (22 school + 1 own class)", len(items))
	}

	// t2 看不见一班的班级级商品。
	items2, err := svc.Items(&t2, "")
	if err != nil {
		t.Fatalf("items t2: %v", err)
	}
	for _, it := range items2 {
		if it.ID == classItem.ID {
			t.Fatalf("t2 should not see another class's item")
		}
	}
}

// 币种过滤：仅返回匹配的币种商品（播种后首次返回默认集不受过滤影响）。
func TestShopItemsCurrencyFilter(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	// 先播种再查询，确保过滤生效。
	svc := services.NewShopService(db, services.NewScope(db),
		services.NewScoreService(db), services.NewCurrencyService(db))
	if _, err := svc.Items(&teacher, ""); err != nil {
		t.Fatalf("seed: %v", err)
	}

	science, err := svc.Items(&teacher, "science")
	if err != nil {
		t.Fatalf("items science: %v", err)
	}
	for _, it := range science {
		if it.CurrencyType != "science" {
			t.Fatalf("science filter returned currency %q", it.CurrencyType)
		}
	}
}

// CreateItem 默认补全 category/currency。
func TestShopCreateItemDefaults(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")

	svc := services.NewShopService(db, services.NewScope(db),
		services.NewScoreService(db), services.NewCurrencyService(db))
	item, err := svc.CreateItem(&teacher, "新商品", "描述", "", "", "", "", 120, 5)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if item.Category != "physical" {
		t.Fatalf("category = %q, want physical", item.Category)
	}
	if item.CurrencyType != "score" {
		t.Fatalf("currency = %q, want score", item.CurrencyType)
	}
	if item.SchoolID != school.ID || item.ClassID != nil {
		t.Fatalf("item scope wrong: school=%d class=%v", item.SchoolID, item.ClassID)
	}
	if !item.IsActive {
		t.Fatal("new item should be active")
	}
}

// 兑换审批（积分类）：扣积分并置为 approved。
func TestShopApproveScoreSettlement(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	class, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	_ = class

	scope := services.NewScope(db)
	scores := services.NewScoreService(db)
	currency := services.NewCurrencyService(db)
	svc := services.NewShopService(db, scope, scores, currency)

	// 先给学生 150 积分。
	if _, err := scores.GiveScore(&student, 150, "测试加分", teacher.ID, nil); err != nil {
		t.Fatalf("give score: %v", err)
	}

	item, err := svc.CreateItem(&teacher, "免作业一次", "免交作业", "privilege", "score", "", "", 100, 0)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}

	red, err := svc.CreateRedemption(&teacher, student.ID, item)
	if err != nil {
		t.Fatalf("create redemption: %v", err)
	}
	if red.Status != "pending" || red.Cost != 100 {
		t.Fatalf("redemption init wrong: %+v", red)
	}

	res, err := svc.ApproveRedemption(&teacher, red.ID)
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	if res.RemainingScore != 50 {
		t.Fatalf("remaining = %d, want 50", res.RemainingScore)
	}
	if res.Message != "已批准兑换，扣除 100 积分" {
		t.Fatalf("message = %q", res.Message)
	}

	// 学生余额刷新。
	var st models.Student
	if err := db.First(&st, student.ID).Error; err != nil {
		t.Fatalf("reload student: %v", err)
	}
	if st.TotalScore != 50 {
		t.Fatalf("student total = %d, want 50", st.TotalScore)
	}

	// 兑换记录状态。
	var dbRed models.ShopRedemption
	if err := db.First(&dbRed, red.ID).Error; err != nil {
		t.Fatalf("reload redemption: %v", err)
	}
	if dbRed.Status != "approved" || dbRed.ApprovedBy == nil || dbRed.ApprovedAt == nil {
		t.Fatalf("approved redemption wrong: %+v", dbRed)
	}
}

// 重复审批同一兑换应报错（该兑换已处理）。
func TestShopApproveTwiceFails(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	scores := services.NewScoreService(db)
	svc := services.NewShopService(db, services.NewScope(db), scores,
		services.NewCurrencyService(db))
	if _, err := scores.GiveScore(&student, 300, "测试", teacher.ID, nil); err != nil {
		t.Fatalf("give: %v", err)
	}

	item, err := svc.CreateItem(&teacher, "铅笔", "", "stationery", "score", "", "", 100, 0)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	red, err := svc.CreateRedemption(&teacher, student.ID, item)
	if err != nil {
		t.Fatalf("create red: %v", err)
	}
	if _, err := svc.ApproveRedemption(&teacher, red.ID); err != nil {
		t.Fatalf("first approve: %v", err)
	}
	if _, err := svc.ApproveRedemption(&teacher, red.ID); err == nil {
		t.Fatal("expected error approving an already-processed redemption")
	} else if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}

// 积分不足时审批报错。
func TestShopApproveInsufficientScore(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	svc := services.NewShopService(db, services.NewScope(db),
		services.NewScoreService(db), services.NewCurrencyService(db))

	item, err := svc.CreateItem(&teacher, "贵重品", "", "stationery", "score", "", "", 1000, 0)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	red, err := svc.CreateRedemption(&teacher, student.ID, item)
	if err != nil {
		t.Fatalf("create red: %v", err)
	}
	if _, err := svc.ApproveRedemption(&teacher, red.ID); err == nil {
		t.Fatal("expected insufficient-score error")
	} else if e, ok := services.AsAppError(err); !ok || e.Status != 400 {
		t.Fatalf("expected 400, got %v", err)
	}
}

// 积分充值类商品审批：按汇率发放钱包币（2 积分 = 1 币）。
func TestShopApproveWalletSettlement(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	class, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")
	_ = class

	scores := services.NewScoreService(db)
	currency := services.NewCurrencyService(db)
	svc := services.NewShopService(db, services.NewScope(db), scores, currency)

	if _, err := scores.GiveScore(&student, 100, "测试", teacher.ID, nil); err != nil {
		t.Fatalf("give: %v", err)
	}

	// 班级积分充值：cost 30 积分 → 15 班级积分。
	item, err := svc.CreateItem(&teacher, "班级积分 +15", "兑换 15 班级积分", "points", "class_point", "", "", 30, 0)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}

	red, err := svc.CreateRedemption(&teacher, student.ID, item)
	if err != nil {
		t.Fatalf("create red: %v", err)
	}
	res, err := svc.ApproveRedemption(&teacher, red.ID)
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	if res.RemainingScore != 70 {
		t.Fatalf("remaining = %d, want 70", res.RemainingScore)
	}

	// 核对钱包余额 = 15。
	var wallet models.Wallet
	if err := db.Where("student_id = ? AND currency_type = ?", student.ID, "class_point").
		First(&wallet).Error; err != nil {
		t.Fatalf("load wallet: %v", err)
	}
	if wallet.Balance != 15 {
		t.Fatalf("wallet balance = %d, want 15", wallet.Balance)
	}
}

// 审批越权：t2 无法审批一班学生的兑换。
func TestShopApproveOutOfScope(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	t1 := seedTeacher(t, db, school.ID, "teacher1")
	t2 := seedTeacher(t, db, school.ID, "teacher2")
	_, student := seedTeacherClass(t, db, school.ID, t1.ID, "一班")
	seedTeacherClass(t, db, school.ID, t2.ID, "二班")

	scores := services.NewScoreService(db)
	svc := services.NewShopService(db, services.NewScope(db), scores,
		services.NewCurrencyService(db))
	if _, err := scores.GiveScore(&student, 200, "测试", t1.ID, nil); err != nil {
		t.Fatalf("give: %v", err)
	}

	item, err := svc.CreateItem(&t1, "铅笔", "", "stationery", "score", "", "", 100, 0)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	red, err := svc.CreateRedemption(&t1, student.ID, item)
	if err != nil {
		t.Fatalf("create red: %v", err)
	}
	if _, err := svc.ApproveRedemption(&t2, red.ID); err == nil {
		t.Fatal("expected out-of-scope error")
	} else if e, ok := services.AsAppError(err); !ok || e.Status != 404 {
		t.Fatalf("expected 404, got %v", err)
	}
}

// 拒绝与发放状态机。
func TestShopRejectAndDeliver(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	scores := services.NewScoreService(db)
	svc := services.NewShopService(db, services.NewScope(db), scores,
		services.NewCurrencyService(db))
	if _, err := scores.GiveScore(&student, 200, "测试", teacher.ID, nil); err != nil {
		t.Fatalf("give: %v", err)
	}

	item, err := svc.CreateItem(&teacher, "铅笔", "", "stationery", "score", "", "", 100, 0)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	red, err := svc.CreateRedemption(&teacher, student.ID, item)
	if err != nil {
		t.Fatalf("create red: %v", err)
	}

	// 拒绝。
	if err := svc.RejectRedemption(&teacher, red.ID); err != nil {
		t.Fatalf("reject: %v", err)
	}
	var r1 models.ShopRedemption
	if err := db.First(&r1, red.ID).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if r1.Status != "rejected" {
		t.Fatalf("status = %q, want rejected", r1.Status)
	}
	// 拒绝后积分未扣。
	var st models.Student
	if err := db.First(&st, student.ID).Error; err != nil {
		t.Fatalf("student: %v", err)
	}
	if st.TotalScore != 200 {
		t.Fatalf("reject should not spend, total = %d", st.TotalScore)
	}

	// 再建一笔并发放。
	red2, err := svc.CreateRedemption(&teacher, student.ID, item)
	if err != nil {
		t.Fatalf("create red2: %v", err)
	}
	if _, err := svc.ApproveRedemption(&teacher, red2.ID); err != nil {
		t.Fatalf("approve: %v", err)
	}
	if err := svc.DeliverRedemption(&teacher, red2.ID); err != nil {
		t.Fatalf("deliver: %v", err)
	}
	var r2 models.ShopRedemption
	if err := db.First(&r2, red2.ID).Error; err != nil {
		t.Fatalf("reload red2: %v", err)
	}
	if r2.Status != "delivered" {
		t.Fatalf("status = %q, want delivered", r2.Status)
	}
}

// ListRedemptions 返回教师管辖班级内记录，并带学生/商品摘要。
func TestShopListRedemptions(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	scores := services.NewScoreService(db)
	svc := services.NewShopService(db, services.NewScope(db), scores,
		services.NewCurrencyService(db))
	if _, err := scores.GiveScore(&student, 200, "测试", teacher.ID, nil); err != nil {
		t.Fatalf("give: %v", err)
	}

	item, err := svc.CreateItem(&teacher, "苹果", "新鲜苹果", "food", "score", "", "", 180, 0)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if _, err := svc.CreateRedemption(&teacher, student.ID, item); err != nil {
		t.Fatalf("create red: %v", err)
	}

	views, err := svc.ListRedemptions(&teacher)
	if err != nil {
		t.Fatalf("list redemptions: %v", err)
	}
	if len(views) != 1 {
		t.Fatalf("got %d redemptions, want 1", len(views))
	}
	v := views[0]
	if v.StudentName != "小明" || v.ItemName != "苹果" {
		t.Fatalf("view summary wrong: %+v", v)
	}
	if v.Cost != 180 || v.CurrencyType != "score" || v.Category != "food" {
		t.Fatalf("view item fields wrong: %+v", v)
	}
	if v.Status != "pending" {
		t.Fatalf("view status = %q, want pending", v.Status)
	}
}
