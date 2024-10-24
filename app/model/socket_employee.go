// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-17
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import jsoniter "github.com/json-iterator/go"

type EmployeeSocketMessage struct {
	SocketMessage
	Speech          string  `json:"speech,omitempty"`          // 播报消息
	AssistanceID    *uint64 `json:"assistanceId,omitempty"`    // 救援ID
	EbikeAllocateID *uint64 `json:"ebikeAllocateId,omitempty"` // 电车分配ID
}

func (res *EmployeeSocketMessage) Bytes() []byte {
	b, _ := jsoniter.Marshal(res)
	return b
}
