// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package snag

import "net/http"

type Error struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func (e Error) Error() string {
    return e.Message
}

func NewError(params ...interface{}) *Error {
    out := &Error{
        Message: "server error",
        Code:    http.StatusBadRequest,
    }

    for _, param := range params {
        switch param.(type) {
        case string:
            out.Message = param.(string)
            break
        case error:
            out.Message = param.(error).Error()
            break
        case int:
            out.Code = param.(int)
            break
        default:
            out.Data = param
            break
        }
    }

    return out
}
