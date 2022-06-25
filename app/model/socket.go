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
    Speech string `json:"speech,omitempty"` // 播报消息
}

func (res *EmployeeSocketMessage) Bytes() []byte {
    b, _ := json.Marshal(res)
    return b
}

type RiderSocketMessage struct {
    SocketMessage

    Assistance *AssistanceRiderSocketMessage `json:"assistance,omitempty"` // 救援消息
}

func (res *RiderSocketMessage) Bytes() []byte {
    b, _ := json.Marshal(res)
    return b
}
