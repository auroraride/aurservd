// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
    "bytes"
    "github.com/auroraride/aurservd/internal/ar"
    jsoniter "github.com/json-iterator/go"
    "github.com/labstack/echo/v4"
    "net/http"
)

const (
    StatusOK                  = iota << 8 //  0x000 请求成功
    StatusError                           //  0x100 请求失败
    StatusUnauthorized                    //  0x200 需要认证 (需要登录)
    StatusForbidden                       //  0x300 没有权限
    StatusNotFound                        //  0x400 资源未获
    StatusInternalServerError             //  0x500 未知错误
    StatusRequireAuth                     //  0x600 需要实名
    StatusLocked                          //  0x700 需要验证 (更换设备需要人脸验证)
    StatusRequireContact                  //  0x800 需要联系人
    StatusRequestTimeout                  //  0x900 请求过期
)

type Response struct {
    echo.Context `json:"-"`

    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func NewResponse(c echo.Context) *Response {
    return &Response{
        Context: c,
    }
}

// SetData 设置数据
func (r *Response) SetData(data interface{}) *Response {
    r.Data = data
    return r
}

// SetMessage 设置消息
func (r *Response) SetMessage(message string) *Response {
    r.Message = message
    return r
}

// SetHeader 设置header
func (r *Response) SetHeader(key, val string) *Response {
    r.Response().Header().Set(key, val)
    return r
}

// SetHeaders 批量设置header
func (r *Response) SetHeaders(m map[string]string) *Response {
    for key, val := range m {
        r.Response().Header().Set(key, val)
    }
    return r
}

// NewError 错误
func (r *Response) Error(code int) *Response {
    r.Code = code
    return r
}

// Success 成功
func (r *Response) Success() *Response {
    r.Code = StatusOK
    return r
}

func (r *Response) Send() error {
    if r.Code == StatusOK && r.Message == "" {
        r.Message = "OK"
    }
    if r.Code == StatusOK && r.Data == nil {
        r.Data = ar.Map{"status": true}
    }
    buffer := &bytes.Buffer{}
    encoder := jsoniter.NewEncoder(buffer)
    encoder.SetEscapeHTML(false)
    _ = encoder.Encode(r)
    return r.JSONBlob(http.StatusOK, buffer.Bytes())
}
