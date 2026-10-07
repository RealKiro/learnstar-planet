// 班级大屏（教室端）写操作：学生间积分转赠。
//
// 移植自 Laravel App\Http\Controllers\Api\DisplayController::quickTransfer。
//
// 说明：
//  1. 两笔积分变动都走既有 ScoreService（等价 Laravel 的 scoreService->giveScore）：
//     含积分记录 + 学生总分 max(0) 钳制 + 宠物经验同步 + score_logs 审计 + 两条 score_update 事件。
//  2. 校验顺序同 Laravel：参数校验（422）→ 学生存在性（404「学生不存在」）→ 余额（400「积分不足」）。
//     有意差异：Laravel 的 `to_id different:from_id` 与 `amount min:1|max:100` 由 Validator 产出
//     英文校验消息，Go 端用与项目其他接口一致的中文 422 文案（状态码一致）。
//  3. Laravel 用 try/catch 把异常统一成 500「转赠失败」；Go 端沿用项目统一的 fail() 兜底
//     （500「服务器内部错误」），未复制该文案。
package services

import (
	"fmt"

	"gorm.io/gorm"
)

// QuickTransferResult 转赠结果（字段名逐字同 Laravel quickTransfer 的 data）。
type QuickTransferResult struct {
	FromName string `json:"from_name"`
	ToName   string `json:"to_name"`
	Amount   int    `json:"amount"`
}

// QuickTransfer 学生间积分转赠（同班两个学生之间，单次 1..100 分）。
func (d *DisplayService) QuickTransfer(classID, fromID, toID uint, amount int) (*QuickTransferResult, error) {
	if fromID == toID {
		return nil, ErrUnprocessable("不能转赠给自己")
	}
	if amount < 1 {
		return nil, ErrUnprocessable("转赠积分至少 1 分")
	}
	if amount > 100 {
		return nil, ErrUnprocessable("单次转赠不能超过 100 分")
	}

	from, err := d.classStudent(classID, fromID)
	if err != nil {
		return nil, err
	}
	to, err := d.classStudent(classID, toID)
	if err != nil {
		return nil, err
	}

	if from.TotalScore < amount {
		return nil, ErrBadRequest("积分不足")
	}

	// 操作人：班级教师，缺失时兜底 user_id = 1（Laravel `$teacherId ?: 1`）。
	operator := classTeacherID(d.db, classID)
	if operator == 0 {
		operator = 1
	}

	outReason := "转赠给 " + to.Name
	inReason := fmt.Sprintf("来自 %s 的转赠", from.Name)

	// ⚠️ 两笔变动必须在同一个事务里：原先各走一次 ScoreService.GiveScore（各自独立提交），
	// 第二步失败就会「扣了没到账」。转出侧另用**条件扣减**而不是「加减 + 钳 0」——否则并发下
	// 余额不足时会被钳到 0 并照常给接收方加分（凭空造分）。
	err = d.db.Transaction(func(tx *gorm.DB) error {
		balanceBefore, newBalance, ok, err := deductScoreAtomic(tx, from.ID, amount)
		if err != nil {
			return err
		}
		if !ok {
			return ErrBadRequest("积分不足")
		}
		if _, err := recordScoreTx(tx, from, -amount, outReason, operator, nil, balanceBefore, newBalance); err != nil {
			return err
		}
		_, err = giveScoreTx(tx, to, amount, inReason, operator, nil)
		return err
	})
	if err != nil {
		return nil, err
	}

	// 事务提交后逐条发布 score_update（同 GiveScore 的口径：逐个学生发布、is_spend=false）。
	d.scores.publishScoreUpdate(from, -amount, outReason, false)
	d.scores.publishScoreUpdate(to, amount, inReason, false)

	return &QuickTransferResult{FromName: from.Name, ToName: to.Name, Amount: amount}, nil
}
