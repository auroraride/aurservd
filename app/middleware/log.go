// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "bytes"
    "fmt"
    "github.com/labstack/echo/v4"
    mw "github.com/labstack/echo/v4/middleware"
    log "github.com/sirupsen/logrus"
    "strings"
)

type BodyDumpConfig struct {
    WithRequestHeaders  bool
    WithResponseHeaders bool
}

var DefaultBodyDumpConfig = BodyDumpConfig{
    WithRequestHeaders:  false,
    WithResponseHeaders: false,
}

// BodyDump 以默认规则保存请求/返回日志
func BodyDump() echo.MiddlewareFunc {
    return BodyDumpWithConfig(DefaultBodyDumpConfig)
}

// BodyDumpRaw 以默认规则保存请求/返回日志
func BodyDumpRaw() echo.MiddlewareFunc {
    return BodyDumpWithConfig(BodyDumpConfig{
        WithRequestHeaders:  true,
        WithResponseHeaders: true,
    })
}

// BodyDumpWithConfig 保存请求/返回日志
func BodyDumpWithConfig(config BodyDumpConfig) echo.MiddlewareFunc {
    return mw.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
        var buffer bytes.Buffer
        buffer.WriteString(fmt.Sprintf("[%s] %s", c.Request().Method, c.Request().RequestURI))
        buffer.WriteRune('\n')
        buffer.WriteString("-----REQ-----")
        buffer.WriteRune('\n')
        if config.WithRequestHeaders {
            for k, v := range c.Request().Header {
                buffer.WriteString(k)
                buffer.WriteString(": ")
                buffer.WriteString(strings.Join(v, ","))
                buffer.WriteRune('\n')
            }
        }
        if len(reqBody) > 0 {
            buffer.WriteRune('\n')
            buffer.Write(reqBody)
        }
        buffer.WriteRune('\n')
        buffer.WriteString(fmt.Sprintf("-----RES[%d]-----", c.Response().Status))
        buffer.WriteRune('\n')
        if config.WithResponseHeaders {
            for k, v := range c.Response().Header() {
                buffer.WriteString(k)
                buffer.WriteString(": ")
                buffer.WriteString(strings.Join(v, ","))
                buffer.WriteRune('\n')
            }
        }
        if len(resBody) > 0 {
            buffer.WriteRune('\n')
            buffer.Write(resBody)
        }
        buffer.WriteRune('\n')
        log.Info(buffer.String())
    })
}
