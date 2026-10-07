// 积分消费的安全守卫测试（幂等 / 原子认领 / 权威余额）。
//
// 背景：这一组缺陷在旧实现里都是「读-判-写不在同一个事务里」或「缺少状态判定」，
// 且经 `git show 0bb6135^:backend/app/Services/*.php` 核对，Laravel 原实现一模一样——
// 也就是说它们是历史设计缺口，本次是有意偏离「忠实移植」口径做的安全性修复。
//
// 覆盖：
//  1. 撤回幂等（同一条记录只能撤回一次，否则可反复点撤回刷分）
//  2. 兑换发放必须已批准（否则「免单」拿到商品且永久无法补扣）
//  3. 已扣款的兑换不能被改成 rejected（否则学生钱没了、记录显示拒绝）
//  4. 审批只扣一次款（条件认领 + 同一事务）
//  5. 余额一律以数据库权威值为准（内存快照陈旧时既不误拒、也不丢更新）
package services_test

import (
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================
// A. 余额以数据库权威值为准（消除读-改-写与「先查后扣」竞态）
// ============================================================

// TestSpendScoreUsesAuthoritativeBalance 内存快照陈旧时，按库里的权威余额判定与扣减。
func TestSpendScoreUsesAuthoritativeBalance(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	svc := services.NewScoreService(db)

	_, err := svc.GiveScore(&student, 100, "起始", 1, nil)
	require.NoError(t, err)

	// 绕过服务直接改库（模拟「另一个请求刚改过余额」），让内存快照立刻变陈旧。
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", student.ID).
		Update("total_score", 500).Error)

	// 陈旧快照上是 100（< 300），但权威余额是 500 → 必须成功扣款。
	// 旧实现会拿快照判定，这里会误报「积分不足」。
	_, err = svc.SpendScore(&student, 300, "测试", 1)
	require.NoError(t, err, "应按数据库权威余额判定，而不是调用方传进来的陈旧快照")

	var reloaded models.Student
	require.NoError(t, db.First(&reloaded, student.ID).Error)
	assert.Equal(t, 200, reloaded.TotalScore, "500-300=200；旧实现会写成 100-300 后钳 0")
}

// TestGiveScoreUsesAuthoritativeBalance 加分同样以权威余额为基准（不覆盖他人的并发写）。
func TestGiveScoreUsesAuthoritativeBalance(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	svc := services.NewScoreService(db)

	_, err := svc.GiveScore(&student, 10, "起始", 1, nil)
	require.NoError(t, err)
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", student.ID).
		Update("total_score", 100).Error)

	_, err = svc.GiveScore(&student, 5, "追加", 1, nil)
	require.NoError(t, err)

	var reloaded models.Student
	require.NoError(t, db.First(&reloaded, student.ID).Error)
	assert.Equal(t, 105, reloaded.TotalScore, "100+5=105；旧实现会写成 10+5=15（丢更新）")
}

// TestSpendScoreRejectsInsufficientWithoutWriting 余额不足时一个字节都不写（不再依赖事后钳 0）。
func TestSpendScoreRejectsInsufficientWithoutWriting(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	svc := services.NewScoreService(db)

	_, err := svc.SpendScore(&student, 50, "测试", 1)
	require.Error(t, err)
	assert.Equal(t, "积分不足，当前余额：0", appErrorOf(t, err, 400))

	var rows int64
	require.NoError(t, db.Model(&models.Score{}).Where("student_id = ?", student.ID).Count(&rows).Error)
	assert.Equal(t, int64(0), rows, "扣款失败不应留下积分流水")

	var reloaded models.Student
	require.NoError(t, db.First(&reloaded, student.ID).Error)
	assert.Equal(t, 0, reloaded.TotalScore)
}

// ============================================================
// B. 撤回幂等
// ============================================================

// TestUndoIsIdempotent 同一条记录只能被撤回一次（否则反复点「撤回」可刷分）。
func TestUndoIsIdempotent(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	svc := services.NewScoreService(db)

	_, err := svc.GiveScore(&student, 10, "加分", 1, nil)
	require.NoError(t, err)
	spend, err := svc.SpendScore(&student, 3, "消费", 1)
	require.NoError(t, err)

	undo, err := svc.Undo(spend, 1)
	require.NoError(t, err, "首次撤回应成功")
	require.NotNil(t, undo.UndoOfScoreID)
	assert.Equal(t, spend.ID, *undo.UndoOfScoreID, "反向流水必须记下被撤回的原记录")

	// 第二次撤回被拒；没有守卫时会把余额从 10 刷到 13。
	_, err = svc.Undo(spend, 1)
	require.Error(t, err, "同一条记录不允许被撤回两次")
	assert.Equal(t, "该积分记录已撤回", appErrorOf(t, err, 400))

	var reloaded models.Student
	require.NoError(t, db.First(&reloaded, student.ID).Error)
	assert.Equal(t, 10, reloaded.TotalScore, "被重复撤回应导致余额变成 13")

	var total int64
	require.NoError(t, db.Model(&models.Score{}).Where("student_id = ?", student.ID).Count(&total).Error)
	assert.Equal(t, int64(3), total, "流水 = 加分 + 消费 + 1 条撤回")

	var undone int64
	require.NoError(t, db.Model(&models.Score{}).Where("undo_of_score_id = ?", spend.ID).
		Count(&undone).Error)
	assert.Equal(t, int64(1), undone)
}

// TestUndoOfScoreIDUniqueIndex 库层兜底：并发双击时第二个反向流水会被唯一索引拦住。
func TestUndoOfScoreIDUniqueIndex(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)

	first := uint(42)
	require.NoError(t, db.Create(&models.Score{
		StudentID: student.ID, ClassID: class.ID, Amount: 1,
		Reason: "撤回操作（手工）", UndoOfScoreID: &first,
	}).Error)

	dup := uint(42)
	err := db.Create(&models.Score{
		StudentID: student.ID, ClassID: class.ID, Amount: 1,
		Reason: "撤回操作（手工重复）", UndoOfScoreID: &dup,
	}).Error
	require.Error(t, err, "undo_of_score_id 的唯一索引应拦住第二条撤回流水")
}

// ============================================================
// C. 兑换状态机
// ============================================================

// shopFixture 商城测试夹具：教师 + 班级 + 学生（含初始积分）+ 商品。
type shopFixture struct {
	DB      *gorm.DB
	Svc     *services.ShopService
	Scores  *services.ScoreService
	Teacher models.User
	Student models.Student
	Item    models.ShopItem
}

func newShopFixture(t *testing.T) shopFixture {
	t.Helper()
	db := setupDB(t)
	school := seedSchool(t, db)
	teacher := seedTeacher(t, db, school.ID, "teacher1")
	_, student := seedTeacherClass(t, db, school.ID, teacher.ID, "一班")

	scores := services.NewScoreService(db)
	svc := services.NewShopService(db, services.NewScope(db), scores, services.NewCurrencyService(db))
	_, err := scores.GiveScore(&student, 200, "测试", teacher.ID, nil)
	require.NoError(t, err)

	item, err := svc.CreateItem(&teacher, "铅笔", "", "stationery", "score", "", "", 100, 0)
	require.NoError(t, err)

	return shopFixture{DB: db, Svc: svc, Scores: scores, Teacher: teacher, Student: student, Item: *item}
}

func (f shopFixture) balance(t *testing.T) int {
	t.Helper()
	var s models.Student
	require.NoError(t, f.DB.First(&s, f.Student.ID).Error)
	return s.TotalScore
}

// TestShopApproveChargesExactlyOnce 审批只扣一次款（条件认领与结算同一事务）。
func TestShopApproveChargesExactlyOnce(t *testing.T) {
	f := newShopFixture(t)
	red, err := f.Svc.CreateRedemption(&f.Teacher, f.Student.ID, &f.Item)
	require.NoError(t, err)

	res, err := f.Svc.ApproveRedemption(&f.Teacher, red.ID)
	require.NoError(t, err)
	assert.Equal(t, 100, res.RemainingScore)

	_, err = f.Svc.ApproveRedemption(&f.Teacher, red.ID)
	require.Error(t, err, "重复审批必须被拒")
	assert.Equal(t, "该兑换已处理", appErrorOf(t, err, 400))

	assert.Equal(t, 100, f.balance(t), "200-100=100：只允许扣一次")

	var spendRows int64
	require.NoError(t, f.DB.Model(&models.Score{}).
		Where("student_id = ? AND amount < 0", f.Student.ID).Count(&spendRows).Error)
	assert.Equal(t, int64(1), spendRows, "扣分流水只应有一条")
}

// TestShopDeliverRequiresApproved 只有已批准（已扣款）的兑换才能标记发放。
func TestShopDeliverRequiresApproved(t *testing.T) {
	f := newShopFixture(t)
	red, err := f.Svc.CreateRedemption(&f.Teacher, f.Student.ID, &f.Item)
	require.NoError(t, err)

	// pending → delivered 必须被拒（否则是「免单」路径）。
	err = f.Svc.DeliverRedemption(&f.Teacher, red.ID)
	require.Error(t, err)
	assert.Equal(t, "该兑换尚未批准，不能标记发放", appErrorOf(t, err, 400))

	var row models.ShopRedemption
	require.NoError(t, f.DB.First(&row, red.ID).Error)
	assert.Equal(t, "pending", row.Status, "被拒后状态不得改变")
	assert.Nil(t, row.DeliveredAt)
	assert.Equal(t, 200, f.balance(t), "未批准就不该扣款")

	// 正规流程：approve → deliver。
	_, err = f.Svc.ApproveRedemption(&f.Teacher, red.ID)
	require.NoError(t, err)
	require.NoError(t, f.Svc.DeliverRedemption(&f.Teacher, red.ID))

	require.NoError(t, f.DB.First(&row, red.ID).Error)
	assert.Equal(t, "delivered", row.Status)
	assert.NotNil(t, row.DeliveredAt, "delivered_at 应落库（此前从未写入）")
}

// TestShopRejectRequiresPending 只有 pending 的兑换才能被拒绝（已扣款的不许改成 rejected）。
func TestShopRejectRequiresPending(t *testing.T) {
	f := newShopFixture(t)

	// 回归保护：pending 仍然可以拒绝，且不扣款。
	red1, err := f.Svc.CreateRedemption(&f.Teacher, f.Student.ID, &f.Item)
	require.NoError(t, err)
	require.NoError(t, f.Svc.RejectRedemption(&f.Teacher, red1.ID))
	var row1 models.ShopRedemption
	require.NoError(t, f.DB.First(&row1, red1.ID).Error)
	assert.Equal(t, "rejected", row1.Status)
	assert.Equal(t, 200, f.balance(t), "拒绝不结算")

	// 已批准（已扣款）的兑换不允许被改成 rejected（否则学生钱没了、记录显示拒绝）。
	red2, err := f.Svc.CreateRedemption(&f.Teacher, f.Student.ID, &f.Item)
	require.NoError(t, err)
	_, err = f.Svc.ApproveRedemption(&f.Teacher, red2.ID)
	require.NoError(t, err)
	require.Equal(t, 100, f.balance(t))

	err = f.Svc.RejectRedemption(&f.Teacher, red2.ID)
	require.Error(t, err)
	assert.Equal(t, "该兑换已处理，不能拒绝", appErrorOf(t, err, 400))

	var row2 models.ShopRedemption
	require.NoError(t, f.DB.First(&row2, red2.ID).Error)
	assert.Equal(t, "approved", row2.Status, "状态不得被改成 rejected")
	assert.Equal(t, 100, f.balance(t), "也不存在退款——钱与记录必须一致")
}

// countAllScores 统计当前积分流水条数（用于断言「整笔失败不留痕」）。
func countAllScores(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var n int64
	require.NoError(t, db.Model(&models.Score{}).Count(&n).Error)
	return n
}

// TestUndoRestoresEffectiveDeltaWhenClamped 撤回按「实际生效的变动额」回补，而不是账面金额。
//
// 场景：余额 10 时被判 -30，实际只扣到 0（「余额不为负」把负数钳成 0）；此时撤回必须只补回
// +10。按账面 -Amount 撤回会凭空多补 20 分（余额变 30）——这正是修正前的行为。
func TestUndoRestoresEffectiveDeltaWhenClamped(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	class := seedClass(t, db, school.ID)
	student := seedStudent(t, db, class.ID)
	svc := services.NewScoreService(db)

	_, err := svc.GiveScore(&student, 10, "起始", 1, nil)
	require.NoError(t, err)

	clamped, err := svc.GiveScore(&student, -30, "惩罚", 1, nil)
	require.NoError(t, err)
	assert.Equal(t, -30, clamped.Amount, "流水仍按账面记 -30")

	var afterGive models.Student
	require.NoError(t, db.First(&afterGive, student.ID).Error)
	require.Equal(t, 0, afterGive.TotalScore, "确认钳制确实发生：10-30 → 0")

	undo, err := svc.Undo(clamped, 1)
	require.NoError(t, err)

	var reloaded models.Student
	require.NoError(t, db.First(&reloaded, student.ID).Error)
	assert.Equal(t, 10, reloaded.TotalScore, "应只补回实际生效的 10 分；按账面 -30 撤回会变成 30")

	var log models.ScoreLog
	require.NoError(t, db.Where("score_id = ?", undo.ID).Order("id DESC").First(&log).Error)
	assert.Equal(t, 0, log.BalanceBefore)
	assert.Equal(t, 10, log.BalanceAfter, "撤回也要落审计行（此前只写 scores、不写 score_logs）")
}

// TestQuickTransferIsAtomicAndStrict 转赠：两笔变动同一事务落库，文案与审计齐全，
// 余额不足整笔失败且不留痕（旧实现两笔各自独立提交，第二步失败会「扣了没到账」）。
func TestQuickTransferIsAtomicAndStrict(t *testing.T) {
	db := setupDB(t)
	f := newDisplayFixture(t, db)

	from := models.Student{ClassID: f.Class.ID, Name: "甲", StudentNo: "1", TotalScore: 50, Status: "active"}
	to := models.Student{ClassID: f.Class.ID, Name: "乙", StudentNo: "2", TotalScore: 10, Status: "active"}
	require.NoError(t, db.Create(&from).Error)
	require.NoError(t, db.Create(&to).Error)

	res, err := f.Svc.QuickTransfer(f.Class.ID, from.ID, to.ID, 30)
	require.NoError(t, err)
	assert.Equal(t, "甲", res.FromName)
	assert.Equal(t, "乙", res.ToName)
	assert.Equal(t, 30, res.Amount)

	var a, b models.Student
	require.NoError(t, db.First(&a, from.ID).Error)
	require.NoError(t, db.First(&b, to.ID).Error)
	assert.Equal(t, 20, a.TotalScore)
	assert.Equal(t, 40, b.TotalScore)

	var out, in models.Score
	require.NoError(t, db.Where("student_id = ?", from.ID).Order("id DESC").First(&out).Error)
	require.NoError(t, db.Where("student_id = ?", to.ID).Order("id DESC").First(&in).Error)
	assert.Equal(t, -30, out.Amount)
	assert.Equal(t, "转赠给 乙", out.Reason)
	assert.Equal(t, 30, in.Amount)
	assert.Equal(t, "来自 甲 的转赠", in.Reason)

	var outLog, inLog models.ScoreLog
	require.NoError(t, db.Where("score_id = ?", out.ID).First(&outLog).Error)
	require.NoError(t, db.Where("score_id = ?", in.ID).First(&inLog).Error)
	assert.Equal(t, 50, outLog.BalanceBefore)
	assert.Equal(t, 20, outLog.BalanceAfter)
	assert.Equal(t, 10, inLog.BalanceBefore)
	assert.Equal(t, 40, inLog.BalanceAfter)

	before := countAllScores(t, db)
	_, err = f.Svc.QuickTransfer(f.Class.ID, from.ID, to.ID, 30)
	require.Error(t, err)
	assert.Equal(t, "积分不足", appErrorOf(t, err, 400))
	assert.Equal(t, before, countAllScores(t, db), "余额不足不应写入任何流水")
}
