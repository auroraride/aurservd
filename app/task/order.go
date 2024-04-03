// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-27
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
	"context"
	"fmt"
	"time"

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
	// 处理预授权订单
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

// Do
// 处理预授权订单 当订阅到期时间为0时，自动发起预授权转支付
// 假如当前订阅未到期 查询当前订阅时间是否超过360天 如果超过360天则自动发起预授权转支付
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
		),
		order.TradePayAtIsNil(),
		order.PaywayIn(model.OrderPaywayAlipayAuthFreeze),
		order.Status(model.OrderStatusPaid),
	).
		WithSubscribe().
		WithPlan().
		All(ctx)

	if len(o) == 0 {
		return
	}

	now := time.Now()
	for _, v := range o {
		// 订阅到期时间为0  或者订阅开始时间到现在超过360天 (提前2天发起预授权转支付)
		if v.Edges.Subscribe != nil && (v.Edges.Subscribe.Remaining == 0 || now.Sub(v.Edges.Subscribe.StartAt.Local()).Hours()/24 >= 358) {
			err := service.NewOrder().TradePay(v)
			if err != nil {
				zap.L().Error(fmt.Sprintf("定时预授权转支付失败 订单ID: %d, 用户id: %d", v.ID, v.RiderID), zap.Error(err))
				continue
			}
		}
	}
}
