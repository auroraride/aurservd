// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import "errors"

type Error struct {
    error
}

func NewError(err interface{}) *Error {
    switch err.(type) {
    case error:
        return &Error{err.(error)}
    case string:
        return &Error{errors.New(err.(string))}
    }
    return &Error{errors.New("请求失败")}
}
