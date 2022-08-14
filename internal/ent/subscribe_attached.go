// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/internal/ent/subscribesuspend"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/golang-module/carbon/v2"
    "math"
    "time"
)

// SubscribeAdditional 订阅额外时间
// 用于计算寄存和暂停
type SubscribeAdditional interface {
    *SubscribePause | *SubscribeSuspend

    GetStartAt() time.Time // 开始日期
    GetEndAt() time.Time   // 结束日期
    GetMaxDays() int       // 最大天数
}

type SubscribeAdditionalResult[T SubscribeAdditional] struct {
    Current     T
    TotalDays   int // 总额外天数
    CurrentDays int // 当前额外天数
    OverdueDays int // 超期天数
}

func (ss *SubscribeSuspend) GetStartAt() time.Time {
    return ss.StartAt
}

func (ss *SubscribeSuspend) GetEndAt() time.Time {
    return ss.EndAt
}

func (ss *SubscribeSuspend) GetMaxDays() int {
    return 0
}

func (sp *SubscribePause) GetStartAt() time.Time {
    return sp.StartAt
}

func (sp *SubscribePause) GetEndAt() time.Time {
    return sp.EndAt
}

func (sp *SubscribePause) GetMaxDays() int {
    return cache.Int(model.SettingPauseMaxDays)
}

func (s *Subscribe) AdditionalItems() (SubscribePauses, SubscribeSuspends) {
    ctx := context.Background()
    pauses := s.Edges.Pauses
    if pauses == nil {
        pauses, _ = s.QueryPauses().Order(Desc(subscribepause.FieldStartAt)).All(ctx)
    }

    suspends := s.Edges.Suspends
    if suspends == nil {
        suspends, _ = s.QuerySuspends().Order(Desc(subscribesuspend.FieldStartAt)).All(ctx)
    }

    return pauses, suspends
}

// SubscribeAdditionalCalculate 计算额外天数
func SubscribeAdditionalCalculate[T SubscribeAdditional](items []T) (data SubscribeAdditionalResult[T]) {
    data = SubscribeAdditionalResult[T]{}
    for _, item := range items {
        end := item.GetEndAt()
        cur := false
        if end.IsZero() {
            // 未结束使用当前时间计算
            data.Current = item
            cur = true
        } else {
            // 已结束按结束时间计算
            end = time.Now()
        }

        // 计算额外时间, 时间需要尽可能的少算, 额外时间 = 结束当天0点 - 开始当日24点(第二天0点)
        days := int(math.Abs(float64(carbon.Time2Carbon(item.GetStartAt()).StartOfDay().AddDay().DiffInDays(carbon.Time2Carbon(end).StartOfDay()))))

        // 判定寄存时间是否超限, 寄存时间超限后会继续计费
        max := item.GetMaxDays()
        if max > 0 && days > max {
            data.OverdueDays = days - max
            days = max
        }

        // 如果是当前阶段计算当前额外时间
        if cur {
            data.CurrentDays = days
        }

        // 累加额外时间
        data.TotalDays += days
    }
    return
}
