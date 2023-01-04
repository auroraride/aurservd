// Copyright (C) liasica. 2022-present.
//
// Created at 2022-04-20
// Based on api by liasica, magicrolan@qq.com.

package app

import (
    "bytes"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/pkg/snag"
    jsoniter "github.com/json-iterator/go"
    "net/http"
)

type Response struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

// processParam 处理参数
func (r *Response) processParam(param any) {
    switch v := param.(type) {
    case string:
        r.Message = v
        break
    case int:
        r.Code = v
    default:
        if v != nil {
            r.Data = param
        }
        break
    }
}

// SendResponse 发送响应
func (c *BaseContext) SendResponse(params ...any) error {
    r := &Response{
        Message: "ok",
        Code:    int(snag.StatusOK),
    }

    for _, param := range params {
        r.processParam(param)
    }

    if r.Code == int(snag.StatusOK) && r.Data == nil {
        r.Data = model.StatusResponse{Status: true}
    }

    buffer := &bytes.Buffer{}
    encoder := jsoniter.NewEncoder(buffer)
    encoder.SetEscapeHTML(false)
    _ = encoder.Encode(r)
 
    return c.JSONBlob(http.StatusOK, buffer.Bytes())
}
