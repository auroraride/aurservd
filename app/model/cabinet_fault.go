// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type CabinetFaultListReq struct {
    PaginationReq

    CityID      *uint64 `json:"cityId" query:"cityID"`
    CabinetName *string `json:"cabinetName" query:"cabinetName"`
    Serial      *string `json:"serial" query:"serial"`
    Status      *uint   `json:"status" query:"status"`
    Start       *string `json:"start" query:"start"`
    End         *string `json:"end" query:"end"`
}
