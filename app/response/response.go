// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package response

import (
    "github.com/labstack/echo/v4"
    "net/http"
)

const (
    StatusOK                  = iota << 8 // 请求成功 0x00
    StatusError                           // 请求失败 0x100
    StatusUnauthorized                    // 需要认证 0x200 (需要登录)
    StatusForbidden                       // 没有权限 0x300
    StatusNotFound                        // 资源未获 0x400
    StatusInternalServerError             // 未知错误 0x500
    StatusNotAcceptable                   // 需要实名 0x600
    StatusLocked                          // 需要验证 0x700 (更换设备需要人脸验证)
    StatusTooManyRequests                 // 触发限流 0x800
    StatusRequestTimeout                  // 请求过期 0x300
)

type response struct {
    echo.Context `json:"-"`

    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func New(c echo.Context) *response {
    return &response{
        Context: c,
    }
}

// SetData 设置数据
func (r *response) SetData(data interface{}) *response {
    r.Data = data
    return r
}

// SetMessage 设置消息
func (r *response) SetMessage(message string) *response {
    r.Message = message
    return r
}

// SetHeader 设置header
func (r *response) SetHeader(key, val string) *response {
    r.Response().Header().Set(key, val)
    return r
}

// SetHeaders 批量设置header
func (r *response) SetHeaders(m map[string]string) *response {
    for key, val := range m {
        r.Response().Header().Set(key, val)
    }
    return r
}

// Error 错误
func (r *response) Error(code int) *response {
    r.Code = code
    return r
}

// Success 成功
func (r *response) Success() *response {
    r.Code = StatusOK
    return r
}

func (r *response) Send() error {
    if r.Code == StatusOK && r.Message == "" {
        r.Message = "OK"
    }
    return r.JSON(http.StatusOK, r)
}
