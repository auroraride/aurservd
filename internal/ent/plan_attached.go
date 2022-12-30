// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-21
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/shopspring/decimal"
    "math"
)

func (pl *Plan) BasicInfo() *model.Plan {
    return &model.Plan{
        ID:          pl.ID,
        Name:        pl.Name,
        Days:        pl.Days,
        Intelligent: pl.Intelligent,
    }
}

func (pl *Plan) OverdueFee(remaining int) (fee float64, formula string) {
    fee, _ = decimal.NewFromFloat(pl.Price).Div(decimal.NewFromInt(int64(pl.Days))).Mul(decimal.NewFromInt(int64(remaining)).Neg()).Mul(decimal.NewFromFloat(1.24)).Float64()
    fee = math.Round(fee*100) / 100

    formula = fmt.Sprintf("(上次购买骑士卡价格 %.2f元 ÷ 天数 %d天) × 逾期天数 %d天 × 1.24 = 逾期费用 %.2f元", pl.Price, pl.Days, remaining, fee)
    return
}

func (pc *PlanCreate) Clone() (creator *PlanCreate) {
    mutation := new(PlanMutation)
    *mutation = *pc.mutation
    return &PlanCreate{
        config:   pc.config,
        mutation: mutation,
        hooks:    pc.hooks,
        conflict: pc.conflict,
    }
}
