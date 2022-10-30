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
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "time"
)

// SubscribeAdditional 订阅额外时间
// 用于计算寄存和暂停
type SubscribeAdditional interface {
    *SubscribePause | *SubscribeSuspend

    GetRiderID() uint64

    // GetStartAt 获取开始日期
    GetStartAt() time.Time

    // GetEndAt 获取结束日期
    GetEndAt() time.Time

    // GetMaxDays 获取最大允许天数
    GetMaxDays() int

    // GetAdditionalDays 获取额外时间
    // 返回参数:
    // days 当前额外时间
    // overdue 超期时间
    // isCurrent 是否生效中
    GetAdditionalDays() (days, overdue, duplicate int, isCurrent bool)
}

type SubscribeAdditionalResult[T SubscribeAdditional] struct {
    Current              T
    TotalDays            int // 总额外天数
    CurrentDays          int // 当前实际天数
    CurrentOverdue       int // 当前超期天数
    CurrentDuplicateDays int // 当前重复天数
}

func (ss *SubscribeSuspend) GetRiderID() uint64 {
    return ss.RiderID
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

func (ss *SubscribeSuspend) GetAdditionalDays() (days, overdue, duplicate int, current bool) {
    days, overdue, current = subscribeAdditionalDays[*SubscribeSuspend](ss)
    return
}

func (sp *SubscribePause) GetRiderID() uint64 {
    return sp.RiderID
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

func (sp *SubscribePause) GetAdditionalDays() (days, overdue, duplicate int, current bool) {
    days, overdue, current = subscribeAdditionalDays[*SubscribePause](sp)

    var overdueDate time.Time
    if overdue > 0 {
        // 寄存开始计算日期
        // 寄存本应超期日期
        overdueDate = carbon.Now().StartOfDay().SubDays(overdue).Carbon2Time()
    }

    items, _ := sp.QuerySuspends().Order(Desc(subscribesuspend.FieldStartAt)).All(context.Background())
    for _, item := range items {
        // 计算实际超期寄存天数
        if !overdueDate.IsZero() {
            sat := tools.NewTime().PauseBeginning(item.StartAt)
            if overdueDate.Before(sat) {
                // 如果超期日期是在暂停之前, 此时 实际超期天数 = 暂停开始日期 - 超期日期
                dd := int((sat.Unix() - overdueDate.Unix()) / 86400)
                overdue -= dd
                days += dd
            } else {
                // 若超期日期在暂停日期之后, 此时 无超期天数
                days += overdue
                overdue = 0
            }
        }
        // 累积重复天数
        d, _, _, _ := item.GetAdditionalDays()
        duplicate += d
    }

    // 实际天数应该减去重复天数, 减去重复天数为实际发生天数
    days -= duplicate

    return
}

func (s *Subscribe) AdditionalItems() (SubscribePauses, SubscribeSuspends) {
    ctx := context.Background()
    pauses, _ := Database.SubscribePause.QueryNotDeleted().Where(subscribepause.SubscribeID(s.ID)).Order(Desc(subscribepause.FieldStartAt)).All(ctx)
    suspends, _ := Database.SubscribeSuspend.Query().Where(subscribesuspend.SubscribeID(s.ID)).Order(Desc(subscribesuspend.FieldStartAt)).All(ctx)
    return pauses, suspends
}

// SubscribeAdditionalCalculate 计算额外天数
func SubscribeAdditionalCalculate[T SubscribeAdditional](items []T) (data SubscribeAdditionalResult[T]) {
    data = SubscribeAdditionalResult[T]{}
    for _, item := range items {
        days, overdue, duplicate, isCurrent := item.GetAdditionalDays()

        if isCurrent {
            data.Current = item
            data.CurrentDays = days
            data.CurrentDuplicateDays = duplicate

            if overdue > 0 {
                data.CurrentOverdue = overdue
            }
        }

        data.TotalDays += days
    }
    return
}

func subscribeAdditionalDays[T SubscribeAdditional](item T) (days int, overdue int, current bool) {
    // 判定是否已结束, 若已结束按结束时间计算
    endAt := item.GetEndAt()
    if endAt.IsZero() {
        // 未结束使用当前时间计算
        current = true
        endAt = time.Now()
    }

    // 获取开始日期
    start := carbon.Time2Carbon(item.GetStartAt())
    // 如果开始日期不是0点, 则开始日期为第二天0点
    if start.Timestamp() != start.StartOfDay().Timestamp() {
        start = start.Tomorrow().StartOfDay()
    }

    // 获取计算截止日期
    end := carbon.Time2Carbon(endAt).StartOfDay()

    // 计算额外时间, 时间需要尽可能的少算, 额外时间 = 结束当天0点 - 开始当日24点(第二天0点)
    days = int(start.DiffAbsInDays(end))

    // 判定寄存时间是否超限, 寄存时间超限后会继续计费
    max := item.GetMaxDays()

    if max > 0 && days > max {
        overdue = days - max
        days = max
    }

    return
}

func (suo *SubscribeUpdateOne) UpdateTarget(cabinetID, storeID, employeeID *uint64) *SubscribeUpdateOne {
    suo.SetNillableCabinetID(cabinetID).
        SetNillableEmployeeID(employeeID).
        SetNillableStoreID(storeID)
    if cabinetID != nil {
        suo.ClearEmployeeID().ClearStoreID()
    }
    if storeID != nil || employeeID != nil {
        suo.ClearCabinetID()
    }
    return suo
}
