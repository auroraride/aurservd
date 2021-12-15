// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package rorre

import "errors"

type Error struct {
    error
    Data interface{}
}

func NewError(params ...interface{}) *Error {
    out := &Error {
        error: errors.New("请求失败"),
    }
    if len(params) > 0 {
        e := params[0]
        switch e.(type) {
        case error:
            out.error = e.(error)
        case string:
            out.error = errors.New(e.(string))
        }
    }
    if len(params) > 1 {
        out.Data = params[1]
    }
    return out
}
