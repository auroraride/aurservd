// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-17
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "github.com/goccy/go-json"
)

type RiderSocketMessage struct {
    SocketMessage

    Assistance   *AssistanceSocketMessage `json:"assistance,omitempty"`   // 救援消息
    ContractSign *ContractSignReq         `json:"contractSign,omitempty"` // 签约消息
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

func (res *RiderSocketMessage) Bytes() []byte {
    b, _ := json.Marshal(res)
    return b
}
