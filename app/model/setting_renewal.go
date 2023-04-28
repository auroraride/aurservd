// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-30
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
	"github.com/auroraride/aurservd/pkg/cache"
)

// RecentSubscribePastDays 距上次订阅已过去天数
type RecentSubscribePastDays float64

func NewRecentSubscribePastDays(days int) RecentSubscribePastDays {
	return RecentSubscribePastDays(days)
}

// Commission 是否需要计算佣金
func (rspd RecentSubscribePastDays) Commission() bool {
	return rspd >= RecentSubscribePastDays(cache.Float64(SettingRenewalKey))
}

func (rspd RecentSubscribePastDays) Value() int {
	return int(rspd)
}
