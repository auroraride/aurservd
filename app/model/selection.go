// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-19
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type SelectionPlanModelReq struct {
	PlanID uint64 `json:"planId" validate:"required" trans:"骑行卡ID" query:"planId"`
}

type SelectionCabinetModelReq struct {
	CabinetID uint64 `json:"cabinetId" validate:"required" trans:"电柜ID" query:"cabinetId"`
}

type SelectionIndexCabinetModelReq struct {
	CityID *uint64 `json:"cityId" query:"cityId"`
}
