// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package snag

import "errors"

type Error struct {
    error
    Data interface{}
}

func NewError(param interface{}) *Error {
    out := &Error{
        error: errors.New("请求失败"),
    }
    switch param.(type) {
    case string:
        out.error = errors.New(param.(string))
    case error:
        out.error = param.(error)
    default:
        out.Data = param
    }
    return out
}
