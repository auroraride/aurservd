package biz

import (
	"context"
	"errors"
	"time"

	"github.com/auroraride/adapter/log"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/payment"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/tools"
)

type orderBiz struct {
	orm *ent.OrderClient
	ctx context.Context
}

func NewOrderBiz() *orderBiz {
	return &orderBiz{
		orm: ent.Database.Order,
		ctx: context.Background(),
	}
}

// DepositFree 押金免支付订单
func (s *orderBiz) DepositFree(r *ent.Rider, req *definition.OrderDepositFreeReq) (result *definition.OrderDepositFreeRes, err error) {
	result = &definition.OrderDepositFreeRes{}

	no := tools.NewUnique().NewSN28()
	// 查询套餐是否存在
	p := service.NewPlan().QueryEffectiveWithID(req.PlanID)

	// 判定套餐是否支持免押金
	if !p.Deposit {
		return nil, errors.New("套餐不支持免押金")
	}

	if p.Deposit {
		var isDepositPayway bool
		for _, item := range p.DepositPayway {
			if item == req.Payway.Value() {
				isDepositPayway = true
				break
			}
		}
		if !isDepositPayway {
			return nil, errors.New("套餐不支持免押金支付方式")
		}
	}

	// 订单字段
	prepay := &model.PaymentCache{
		CacheType: model.PaymentCacheTypeAliDepositFree,
		DepositFree: &model.DepositFree{
			OutTradeNo: no,
			RiderID:    r.ID,
			Amount:     p.DepositAmount,
			Plan:       p.BasicInfo(),
			PlanID:     req.PlanID,
			Payway:     req.Payway.Value(),
		},
	}

	// 订单缓存
	err = cache.Set(s.ctx, no, prepay, 20*time.Minute).Err()
	if err != nil {
		zap.L().Error("获取免押订单信息失败", zap.Error(err))
		return nil, err
	}

	// 发起请求
	var str string
	switch req.Payway {
	case definition.PaywayAlipayDeposit:
		str, err = payment.NewAlipay().FandAuthFreeze(prepay)
		if err != nil {
			zap.L().Error("支付宝预授权支付请求失败", zap.Error(err))
			return nil, err
		}
	case definition.PaywayWechatDeposit:
		// TODO

	}
	result.Prepay = str
	return
}

func (s *orderBiz) DoPayment(pc *model.PaymentCache) {
	if pc == nil {
		return
	}
	switch pc.CacheType {
	case model.PaymentCacheTypePlan:
		// 预支付
		service.NewOrder().OrderPaid(pc.Subscribe)
	case model.PaymentCacheTypeAliDepositFree:
		// 押金
		s.Deposit(pc.DepositFree)
	}
}

// 押金创建订单
func (s *orderBiz) Deposit(trade *model.DepositFree) {
	// 查询押金订单是否存在
	if exists, err := ent.Database.Order.Query().Where(order.OutTradeNo(trade.OutTradeNo)).Exist(s.ctx); err == nil && exists {
		return
	}
	zap.L().Info("押金订单: "+trade.OutTradeNo, log.JsonData(trade))

	// 创建押金订单
	_, err := s.orm.Create().
		SetStatus(model.OrderStatusPaid).
		SetPayway(trade.Payway).
		SetType(model.OrderTypeDeposit).
		SetOutTradeNo(trade.OutTradeNo).
		SetTradeNo(trade.TradeNo).
		SetAmount(trade.Amount).
		SetTotal(trade.Amount).
		SetRiderID(trade.RiderID).
		Save(s.ctx)
	if err != nil {
		zap.L().Error("订单已支付, 但押金订单创建失败: "+trade.OutTradeNo, zap.Error(err))
		return
	}

}
