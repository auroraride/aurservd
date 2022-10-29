// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-29
// Based on aurservd by liasica, magicrolan@qq.com.

package env

type Env[T any] struct {
    Key   string `json:"key"`   // 键值
    Desc  string `json:"desc"`  // 描述
    Value T      `json:"value"` // 值
}
