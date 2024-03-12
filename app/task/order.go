// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-27
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
	"context"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
)

type orderTask struct {
}

func NewOrder() *orderTask {
	return &orderTask{}
}

func (t *orderTask) Start() {
	c := cron.New()
	_, err := c.AddFunc("0 9 * * *", func() {
		zap.L().Info("开始执行 @daily[orderTradePay] 定时任务")
		t.Do()
	})
	if err != nil {
		zap.L().Fatal("@daily[orderTradePay] 定时任务执行失败", zap.Error(err))
		return
	}
	c.Start()
}

func (t *orderTask) Do() {
	ctx := context.Background()
	// 查询订阅到期的
	o, _ := ent.Database.Order.QueryNotDeleted().Where(
		order.HasSubscribeWith(
			// 未退款
			subscribe.RefundAtIsNil(),
			// 未结束
			subscribe.EndAtIsNil(),
			// 已开始
			subscribe.StartAtNotNil(),
			// 非企业
			subscribe.EnterpriseIDIsNil(),
			// 剩余0天
			subscribe.RemainingEQ(0),
		),
		order.TradePayAtIsNil(),
		order.PaywayIn(model.OrderPaywayAlipayAuthFreeze, model.OrderPaywayWechatDeposit),
		// 查询订单和关联订单金额
		// order.Type()
	).
		WithSubscribe().
		WithPlan().
		All(ctx)

	if len(o) == 0 {
		return
	}

	for _, v := range o {
		service.NewOrder().TradePay(v)
	}
}
