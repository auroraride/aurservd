// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-27
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
	"context"
	"fmt"

	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/order"
)

type orderTask struct {
}

func NewOrder() *orderTask {
	return &orderTask{}
}

func (t *orderTask) Start() {
	if !ar.Config.Task.Order {
		return
	}

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
// 处理预授权订单 当订阅到期时间为,自动发起预授权转支付
// 订单创建时间距离当前时间超过358天，自动发起预授权转支付
func (t *orderTask) Do() {
	ctx := context.Background()
	now := carbon.Now().StdTime()
	nextStart := carbon.Now().AddDay().StartOfDay().StdTime()

	o, _ := ent.Database.Order.QueryNotDeleted().Where(
		order.SubscribeEndAtLTE(nextStart),
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

	for _, v := range o {
		// 如果订阅到期时间为当天
		if v.SubscribeEndAt.Format("2006-01-02") == now.Format("2006-01-02") ||
			// 如果订单创建时间距离当前时间超过358天
			now.Sub(v.CreatedAt.Local()).Hours()/24 >= 358 {
			err := service.NewOrder().TradePay(v)
			if err != nil {
				zap.L().Error(fmt.Sprintf("定时预授权转支付失败 订单ID: %d, 用户id: %d", v.ID, v.RiderID), zap.Error(err))
				continue
			}
		}
	}
}
