// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type Setting struct {
    Desc  string      `json:"desc"`
    Value interface{} `json:"value"`
}
