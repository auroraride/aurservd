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
    if s.entRider == nil {
        return nil
    }

    points := s.entRider.Points

    ratio := model.PointRatio
    max := ratio * float64(points)

    res := &model.PaymentConfigure{
        Points:  points,
        Ratio:   ratio,
        Coupons: make([]model.Coupon, 0),
    }

    // 获取骑手优惠券
    now := time.Now()
    coupons := NewCoupon().QueryEffective(s.rider.ID)

    // 计算优惠券
    var exclusive, stackable float64
    stackables := make(map[string]float64)
    for _, c := range coupons {
        isExclusive := c.Rule == model.CouponRuleExclusive.Value()
        cate := utils.Md5String(fmt.Sprintf("%d", c.TemplateID))
        amount := c.Amount
        r := model.Coupon{
            ID:        c.ID,
            Cate:      cate,
            Useable:   c.ExpiresAt.After(now),
            Amount:    amount,
            Name:      c.Name,
            ExpiredAt: c.ExpiresAt.Format("2006.1.2"),
            Code:      c.Code,
            Exclusive: isExclusive,
            Plans:     c.Plans,
            Cities:    c.Cities,
        }
        // 可叠加券
        if sc, ok := stackables[cate]; !ok || sc < amount {
            stackables[cate] = amount
        }
        // 互斥券, 计算最大值
        if isExclusive && exclusive < amount {
            exclusive = amount
        }

        res.Coupons = append(res.Coupons, r)
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
