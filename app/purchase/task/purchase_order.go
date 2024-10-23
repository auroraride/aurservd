package task

import (
	"context"

	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/purchasepayment"
)

type purchaseOrderTask struct {
}

func NewPurchaseOrder() *purchaseOrderTask {
	return &purchaseOrderTask{}
}

func (t *purchaseOrderTask) Start() {
	go t.Do()
	c := cron.New()
	_, err := c.AddFunc("@daily", func() {
		zap.L().Info("开始执行 @daily[purchaseOrderTask] 定时任务")
		t.Do()
	})
	if err != nil {
		zap.L().Fatal("@daily[purchaseOrderTask] 定时任务执行失败", zap.Error(err))
		return
	}
	c.Start()
}

// Do 检查当前分期账单是否逾期
func (*purchaseOrderTask) Do() {
	now := carbon.Now().StartOfDay().AddDays(-8).StdTime()
	payments, _ := ent.Database.PurchasePayment.QueryNotDeleted().Where(
		purchasepayment.StatusEQ(purchasepayment.StatusObligation),
		purchasepayment.BillingDateEQ(now),
	).All(context.Background())
	for _, v := range payments {
		if v.Forfeit != 0 {
			continue
		}
		err := v.Update().
			SetForfeit(v.Amount * 0.2).
			SetTotal(v.Amount * 1.2).
			Exec(context.Background())
		if err != nil {
			zap.L().Error("更新逾期账单失败", zap.Error(err))
			continue
		}
	}
}
