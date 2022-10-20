// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    AllocateExpiration = 600 // 分配过期时间(s), 超过一定时间不签约后分配失效, 需要重新分配
)

type AllocateStatus uint8

const (
    AllocateStatusPending AllocateStatus = iota + 1 // 未激活
    AllocateStatusSigned                            // 已签约
    AllocateStatusVoid                              // 已作废
)

func (s AllocateStatus) Value() uint8 {
    return uint8(s)
}

type AllocateEmployeeListReq struct {
    PaginationReq
    Status AllocateStatus `json:"status" query:"status"` // 签约状态 1:未签约(默认) 2:已签约
}

type AllocateDetail struct {
    Rider  Rider          `json:"rider" bson:"rider"` // 骑手信息
    ID     uint64         `json:"id"`
    Type   string         `json:"type" enums:"battery,ebike"` // 分配类型 battery:单电 ebike:车电
    Status AllocateStatus `json:"status"`                     // 1:未激活 2:已签约 3:已作废
    Time   string         `json:"time"`                       // 分配时间
    Model  string         `json:"model"`                      // 电池型号
    Ebike  *Ebike         `json:"ebike,omitempty"`            // 电车信息
}

type AllocateCreateReq struct {
    EbikeID     *uint64 `json:"ebikeId"` // 电车ID
    SubscribeID uint64  `json:"subscribeId" validate:"required" trans:"订阅ID"`

    StoreID    *uint64 `swaggerignore:"true"`
    EmployeeID *uint64 `swaggerignore:"true"`
    CabinetID  *uint64 `swaggerignore:"true"`
}
