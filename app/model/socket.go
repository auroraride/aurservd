// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "encoding/json"

type SocketBinaryMessage interface {
    Bytes() []byte
}

type SocketMessage struct {
    Error string `json:"error"`
}

func (res *SocketMessage) Bytes() []byte {
    b, _ := json.Marshal(res)
    return b
}

type EmployeeSocketMessage struct {
    SocketMessage
    Speech          string  `json:"speech,omitempty"`          // 播报消息
    AssistanceID    *uint64 `json:"assistanceId,omitempty"`    // 救援ID
    EbikeAllocateID *uint64 `json:"ebikeAllocateId,omitempty"` // 电车分配ID
}

func (res *EmployeeSocketMessage) Bytes() []byte {
    b, _ := json.Marshal(res)
    return b
}

type RiderSocketMessage struct {
    SocketMessage

    Assistance *AssistanceSocketMessage `json:"assistance,omitempty"` // 救援消息
}

func (res *RiderSocketMessage) Bytes() []byte {
    b, _ := json.Marshal(res)
    return b
}

type AssistanceSocketMessage struct {
    ID        uint64       `json:"id"`
    Status    uint8        `json:"status"`             // 状态 0:待分配 1:已分配 2:已拒绝 3:已失败 4:待支付 5:已支付
    Store     *StoreLngLat `json:"store,omitempty"`    // 门店信息
    Employee  *Employee    `json:"employee,omitempty"` // 店员信息
    Rider     LngLat       `json:"rider"`              // 骑手坐标
    Seconds   int          `json:"seconds"`            // 距离分配等待时间
    Breakdown string       `json:"breakdown"`          // 故障原因
}
