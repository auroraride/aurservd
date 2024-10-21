// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-10-21, by liasica

package helper

import (
	"time"

	"github.com/golang-module/carbon/v2"
)

// OverdueDays 计算逾期天数
// t1:开始时间, t2:结束时间
// 例如: t1 = 2024-10-21 00:08:00, t2 = 2024-10-22 00:00:00, 结果为逾期1天
func OverdueDays(t1, t2 time.Time) int {
	c1 := carbon.CreateFromStdTime(t1).StartOfDay()
	c2 := carbon.CreateFromStdTime(t2).StartOfDay()

	if c1.Gte(c2) {
		return 0
	}

	return int(c1.DiffInDays(c2))
}
