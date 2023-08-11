// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-02
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
	"context"
	"errors"

	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionearnings"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/pkg/tools"
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

	dl := tools.NewDecimal()
	// 记录每个用户的总收益
	total := make(map[uint64]float64)

	// 遍历收益明细 计算总收益
	for _, item := range all {
		total[item.MemberID] = dl.Sum(total[item.MemberID], item.Amount)
	}

	ent.WithTxPanic(ctx, func(tx *ent.Tx) (err error) {
		for k, v := range total {

			// 查询用户
			member, _ := tx.PromotionMember.Query().Where(promotionmember.ID(k)).First(ctx)
			if member == nil {
				zap.L().Error("用户不存在更新冻结金额失败", zap.Uint64("会员id", k), zap.Float64("金额", v))
				continue
			}
			balance := dl.Sum(member.Balance, v)
			frozen := dl.Sub(member.Frozen, v)
			// 更新用户总收益
			_, err = tx.PromotionMember.UpdateOneID(k).SetBalance(balance).SetFrozen(frozen).Save(ctx)
			if err != nil {
				zap.L().Error("更新用户总收益失败", zap.Error(err), zap.Uint64("会员id", k), zap.Float64("金额", v))
			}
		}

		ids := make([]uint64, len(all))
		for i, v := range all {
			ids[i] = v.ID
		}

		if len(ids) == 0 {
			return
		}

		if _, err = tx.PromotionEarnings.Update().
			Where(promotionearnings.IDIn(ids...)).
			SetStatus(promotion.EarningsStatusSettled.Value()).
			SetRemark("定时任务结算").
			Save(ctx); err != nil {
			zap.L().Error("更新收益失败", zap.Error(err), zap.Uint64s("ids", ids))
			return errors.New("更新收益失败")
		}
		return
	})
}
