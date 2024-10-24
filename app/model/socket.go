// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import jsoniter "github.com/json-iterator/go"

type SocketBinaryMessage interface {
	Bytes() []byte
}

type SocketMessage struct {
	Error string `json:"error"`
}

func (res *SocketMessage) Bytes() []byte {
	b, _ := jsoniter.Marshal(res)
	return b
}
