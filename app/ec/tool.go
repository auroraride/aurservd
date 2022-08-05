// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-31
// Based on aurservd by liasica, magicrolan@qq.com.

package ec

func Pointer[T any](i T) *T {
    return &i
}
