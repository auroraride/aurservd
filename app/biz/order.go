package biz

import (
	"context"
	"time"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/payment"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type order struct {
	orm *ent.OrderClient
	ctx context.Context
}

func NewOrder() *order {
	return &order{
		orm: ent.Database.Order,
	}
}

// DepositFree 押金免支付订单
func (s *order) DepositFree(r *ent.Rider, req *definition.OrderDepositFreeReq) (result *definition.OrderDepositFreeRes) {
	no := tools.NewUnique().NewSN28()
	// 查询套餐是否存在
	p := service.NewPlan().QueryEffectiveWithID(req.PlanID)

	// 订单字段
	prepay := &model.PaymentCache{
		CacheType: model.PaymentCacheTypeDepositFree,
		Subscribe: &model.PaymentSubscribe{
			OrderType:    7,
			OutTradeNo:   no,
			RiderID:      r.ID,
			Name:         p.Name + "押金",
			Amount:       p.DepositAmount,
			Payway:       req.Payway,
			EbikeBrandID: p.BrandID,
			Plan:         p.BasicInfo(),
		},
	}
	// 订单缓存
	err := cache.Set(s.ctx, no, prepay, 20*time.Minute).Err()
	str, err := payment.NewAlipay().FandAuthFreeze(prepay)
	if err != nil {
		snag.Panic("支付宝预授权支付请求失败")
	}
	result.Prepay = str
	return
}
