// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    EbikeAllocateExpiration = 600 // 分配过期时间(s), 超过一定时间不签约后分配失效, 需要重新分配
)

type EbikeAllocateStatus uint8

const (
    EbikeAllocateStatusPending EbikeAllocateStatus = iota + 1 // 未激活
    EbikeAllocateStatusSigned                                 // 已签约
    EbikeAllocateStatusVoid                                   // 已作废
)

func (s EbikeAllocateStatus) Value() uint8 {
    return uint8(s)
}

type EbikeAllocate struct {
    Rider Rider  `json:"rider" bson:"rider"` // 骑手信息
    Ebike *Ebike `json:"ebike" bson:"ebike"` // 电车信息
    Model string `json:"model" bson:"model"` // 电池型号
}

type EbikeAllocateReq struct {
    EbikeID     uint64 `json:"ebikeId" validate:"required" trans:"电车ID"`
    SubscribeID uint64 `json:"subscribeId" validate:"required" trans:"订阅ID"`
}

type EbikeAllocateEmployeeListReq struct {
    PaginationReq
    Status EbikeAllocateStatus `json:"status" query:"status"` // 签约状态 1:未签约(默认) 2:已签约
}

type EbikeAllocateDetail struct {
    *EbikeAllocate
    ID     uint64              `json:"id"`
    Status EbikeAllocateStatus `json:"status"` // 1:未激活 2:已签约 3:已作废
    Time   string              `json:"time"`   // 分配时间
}
