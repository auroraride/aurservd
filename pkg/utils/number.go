// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import "math"

type number struct {
}

func NewNumber() *number {
	return &number{}
}

func (*number) Decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}
