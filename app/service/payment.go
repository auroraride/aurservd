// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/pkg/utils"
    "math"
    "time"
)

type paymentService struct {
    *BaseService
}

func NewPayment(params ...any) *paymentService {
    return &paymentService{
        BaseService: newService(params...),
    }
}

// Configure 获取支付配置
func (s *paymentService) Configure() *model.PaymentConfigure {
    if s.Rider() == nil {
        return nil
    }

    points := s.Rider().Points

    ratio := 0.01
    max := ratio * float64(points)

    res := &model.PaymentConfigure{
        Points:  points,
        Ratio:   ratio,
        Coupons: make([]model.Coupon, 0),
    }

    // 获取骑手优惠券
    now := time.Now()
    coupons := NewCoupon().QueryRiderNotUsed(s.rider.ID)

    // 计算优惠券
    var exclusive, stackable float64
    stackables := make(map[string]float64)
    for _, coupon := range coupons {
        isExclusive := coupon.Rule == model.CouponRuleExclusive.Value()
        cate := utils.Md5String(fmt.Sprintf("%d", coupon.TemplateID))
        amount := coupon.Amount
        res.Coupons = append(res.Coupons, model.Coupon{
            Cate:      cate,
            Useable:   coupon.ExpiresAt.After(now),
            Amount:    amount,
            Name:      coupon.Name,
            ExpiredAt: coupon.ExpiresAt.Format("2006.1.2"),
            Code:      coupon.Code,
            Exclusive: isExclusive,
        })
        // 可叠加券
        if sc, ok := stackables[cate]; !ok || sc < amount {
            stackables[cate] = amount
        }
        // 互斥券, 计算最大值
        if isExclusive && exclusive < amount {
            exclusive = amount
        }
    }

    // 可叠加券值
    for _, sc := range stackables {
        stackable += sc
    }

    if stackable > exclusive {
        max += stackable
    } else {
        max += exclusive
    }

    max = math.Round(max*100.00) / 100.00
    res.MaxDiscount = max

    return res
}
