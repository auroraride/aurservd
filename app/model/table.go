// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-21
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type Table interface {
    GetTableName() string
    GetID() uint64
}
