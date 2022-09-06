// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type SuspendReq struct {
    ID     uint64 `json:"id" validate:"required" trans:"订阅ID"`
    Remark string `json:"remark" validate:"required" trans:"备注"`
}

type SuspendListReq struct {
    PaginationReq
    SuspendListFilter
}

type SuspendListFilter struct {
    CityID  uint64 `json:"cityId" query:"cityId"`   // 城市ID
    RiderID uint64 `json:"riderId" query:"riderId"` // 骑手ID
    Start   string `json:"start" query:"start"`     // 开始日期, 如 2022-08-01
    End     string `json:"end" query:"end"`         // 结束日期, 如 2022-08-07
}

type SuspendListRes struct {
    City            string `json:"city"`            // 城市
    Name            string `json:"name"`            // 姓名
    Phone           string `json:"phone"`           // 手机号
    Plan            string `json:"plan"`            // 骑士卡
    Status          string `json:"status"`          // 状态
    SubscribeDays   int    `json:"subscribeDays"`   // 剩余天数
    SubscribeStatus string `json:"subscribeStatus"` // 骑士卡状态
    Days            int    `json:"days"`            // 暂停天数
    Start           string `json:"start"`           // 开始时间
    StartBy         string `json:"startBy"`         // 开始操作人
    End             string `json:"end"`             // 结束时间
    EndBy           string `json:"endBy"`           // 结束操作人
    Reason          string `json:"reason"`          // 暂停原因
}
