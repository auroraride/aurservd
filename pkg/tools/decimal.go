// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-20
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
	"math"

	"github.com/shopspring/decimal"
)

type decimalTool struct {
}

func NewDecimal() *decimalTool {
	return &decimalTool{}
}

// Sum returns f1 + f2
func (*decimalTool) Sum(f1, f2 float64) float64 {
	f, _ := decimal.NewFromFloat(f1).Add(decimal.NewFromFloat(f2)).Float64()
	return math.Round(f*100.00) / 100.0
}

// Sub returns f1 - f2
func (*decimalTool) Sub(f1, f2 float64) float64 {
	f, _ := decimal.NewFromFloat(f1).Sub(decimal.NewFromFloat(f2)).Float64()
	return math.Round(f*100.00) / 100.0
}

// Mul returns f1 × f2
func (*decimalTool) Mul(f1, f2 float64) float64 {
	f, _ := decimal.NewFromFloat(f1).Mul(decimal.NewFromFloat(f2)).Float64()
	return math.Round(f*100.0) / 100.0
}
