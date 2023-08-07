// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-02
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
	"context"

	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionearnings"
)

type promotionEarningsTask struct {
}

func NewPromotionEarnings() *promotionEarningsTask {
	return &promotionEarningsTask{}
}

func (t *promotionEarningsTask) Start() {

	go t.Do()

	c := cron.New()
	_, err := c.AddFunc("@daily", func() {
		zap.L().Info("开始执行 @daily[earnings] 定时任务")
		go t.Do()
	})
	if err != nil {
		zap.L().Fatal("@daily[earnings] 定时任务执行失败", zap.Error(err))
		return
	}
	c.Start()
}

func (*promotionEarningsTask) Do() {
	ctx := context.Background()

	// 查询收益明细 7天前的收益
	all, _ := ent.Database.PromotionEarnings.Query().Where(
		promotionearnings.CreatedAtGTE(carbon.Now().StartOfDay().AddDays(-7).Carbon2Time()),
		promotionearnings.CreatedAtLTE(carbon.Now().EndOfDay().AddDays(-7).Carbon2Time()),
		promotionearnings.Status(promotion.EarningsStatusUnsettled.Value()),
	).All(ctx)

	// 记录每个用户的总收益
	total := make(map[uint64]float64)

	// 遍历收益明细 计算总收益
	for _, item := range all {
		total[item.MemberID] += item.Amount
	}

	ent.WithTxPanic(ctx, func(tx *ent.Tx) (err error) {
		for k, v := range total {
			// 更新用户总收益
			_, err = tx.PromotionMember.UpdateOneID(k).SetBalance(v).AddFrozen(-v).Save(ctx)
			if err != nil {
				zap.L().Error("更新用户总收益失败", zap.Error(err))
				return err
			}
		}
		for _, v := range all {
			// 更新收益明细状态
			_, err = tx.PromotionEarnings.Update().Where(promotionearnings.ID(v.ID)).SetStatus(promotion.EarningsStatusSettled.Value()).Save(ctx)
			if err != nil {
				zap.L().Error("更新收益失败", zap.Error(err))
				return
			}
		}
		return
	})
}
