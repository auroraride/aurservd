// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
	"context"
	"time"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/internal/ent/subscribepause"
	"github.com/auroraride/aurservd/internal/ent/subscribesuspend"
)

// SubscribeAdditional 订阅额外时间
// 用于计算寄存和暂停
type SubscribeAdditional interface {
	*SubscribePause | *SubscribeSuspend

	// GetStartAt 获取开始日期
	GetStartAt() time.Time

	// GetEndAt 获取结束日期
	GetEndAt() time.Time

	// GetMaxDays 获取最大允许天数
	GetMaxDays() int
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
	// return cache.Int(model.SettingPauseMaxDays)
	// TODO 若修改时间
	// TODO 寄存结束后需要标记固化寄存天数防止重复计算
	return 30
}

type SubscribeAdditionalItem struct {
	Start         carbon.Carbon
	Stop          carbon.Carbon
	OriginStop    carbon.Carbon // 原结束日期
	Current       bool          // 是否当前生效中
	MaxDays       int           // 最大天数 0:不限制
	PastDays      int           // 总过天数
	OverdueDays   int           // 过期天数
	Days          int           // 有效额外天数
	DuplicateDays int           // 重复天数
}

func NewSubscribeAdditionalItem[T SubscribeAdditional](item T) *SubscribeAdditionalItem {
	// 获取开始日期
	start := carbon.CreateFromStdTime(item.GetStartAt())
	// 如果开始日期不是0点, 则开始日期为第二天0点
	if start.Timestamp() != start.StartOfDay().Timestamp() {
		start = start.Tomorrow().StartOfDay()
	}

	// 判定是否已结束, 若已结束按结束时间计算
	var (
		current bool

		endAt   = item.GetEndAt()
		maxDays = item.GetMaxDays()
	)

	stop := carbon.CreateFromStdTime(endAt)
	if endAt.IsZero() {
		// 未结束使用当前时间计算
		current = true
		// data.Current = true
		stop = carbon.Now()
	}

	// 获取计算截止日期, 当日0点
	stop = stop.StartOfDay()

	// data.Stop = stop
	return &SubscribeAdditionalItem{
		Start:      start,
		Stop:       stop,
		OriginStop: stop,
		Current:    current,
		MaxDays:    maxDays,
	}
}

// CalculateDays 计算天数
func (sa *SubscribeAdditionalItem) CalculateDays() *SubscribeAdditionalItem {
	var (
		overdueDays   int
		duplicateDays int
	)

	// 计算额外时间, 时间需要尽可能的少算, 额外时间 = 结束当天0点 - 开始当日24点(第二天0点)
	pastDays := int(sa.Start.DiffAbsInDays(sa.Stop))
	days := pastDays

	if sa.MaxDays > 0 && pastDays > sa.MaxDays {
		overdueDays = pastDays - sa.MaxDays
		days = sa.MaxDays

		// 当寄存中生效的时候暂停扣费会产生重复天数, 需要计算并记录该部分重复的天数
		// 计算暂停计费和寄存之间的重复天数
		duplicateDays = int(sa.OriginStop.DiffInDays(sa.Stop))
	}

	sa.PastDays = pastDays
	sa.OverdueDays = overdueDays
	sa.Days = days
	sa.DuplicateDays = duplicateDays

	return sa
}

type SubscribeSuspendResult struct {
	Current *SubscribeSuspend // 当前暂停

	Days        int // 总额外天数
	CurrentDays int // 当前实际天数
}

type SubscribePauseResult struct {
	Current *SubscribePause // 当前寄存

	Days                 int // 总额外天数
	CurrentDays          int // 当前实际天数
	CurrentOverdueDays   int // 当前超期天数
	CurrentDuplicateDays int // 当前重复天数
}

// GetAdditionalItems 获取额外时间集
func (s *Subscribe) GetAdditionalItems() (pr SubscribePauseResult, sr SubscribeSuspendResult) {
	var (
		// pauses   []*SubscribeAdditionalItem // TODO: pauses是做啥用的??? 忘记了
		suspends []*SubscribeAdditionalItem
	)
	ctx := context.Background()
	subPauses, _ := Database.SubscribePause.QueryNotDeleted().Where(subscribepause.SubscribeID(s.ID)).Order(Desc(subscribepause.FieldStartAt)).All(ctx)
	subSuspends, _ := Database.SubscribeSuspend.Query().Where(subscribesuspend.SubscribeID(s.ID)).Order(Desc(subscribesuspend.FieldStartAt)).All(ctx)

	// 暂停数据
	for _, u := range subSuspends {
		item := NewSubscribeAdditionalItem(u).CalculateDays()
		suspends = append(suspends, item)
		sr.Days += item.Days
		if item.Current {
			sr.Current = u
			sr.CurrentDays = item.Days
		}
	}

	// 计算寄存和暂停的时间
	for _, p := range subPauses {
		item := NewSubscribeAdditionalItem(p)
		// 寄存时间 到 结束寄存时间 之间若有 暂停, 则将结束寄存时间计算为暂停开始时间.
		// 若时间重叠, 则标记结束寄存结束时间
		for _, u := range suspends {
			// 若暂停开始时间晚于寄存时间并且早于寄存结束时间, 则将暂停时间计算为该阶段的寄存结束时间
			if u.Start.DiffInSeconds(item.Start) >= 0 && item.Stop.DiffInSeconds(u.Start) >= 0 {
				item.Stop = u.Start
			}
		}
		item.CalculateDays()
		pr.Days += item.Days
		if item.Current {
			pr.Current = p
			pr.CurrentDays = item.Days
			pr.CurrentOverdueDays = item.OverdueDays
			pr.CurrentDuplicateDays = item.DuplicateDays
		}
		// pauses = append(pauses, item)
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
