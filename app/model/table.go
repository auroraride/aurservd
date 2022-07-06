// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-21
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type Table interface {
    GetTableName() string
    GetID() uint64
}

type IDName interface {
    GetName() string
    GetID() uint64
}

type NilIDName struct {
}

func (d *NilIDName) GetID() uint64 {
    return 0
}

func (d *NilIDName) GetName() string {
    return ""
}
