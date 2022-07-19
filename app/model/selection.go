// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-19
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type SelectionPlanModelReq struct {
    PlanID uint64 `json:"planId" validate:"required" trans:"骑行卡ID" query:"planId"`
}
