// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "bufio"
    "bytes"
    "fmt"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/golang-module/carbon/v2"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
    "io"
    "io/fs"
    "io/ioutil"
    "net"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)

type BodyDumpConfig struct {
    WithRequestHeaders  bool
    WithResponseHeaders bool
    Skipper             map[string]bool
}

var DefaultBodyDumpConfig = BodyDumpConfig{
    WithRequestHeaders:  false,
    WithResponseHeaders: false,
}

type bodyDumpHandler func(echo.Context, []byte, []byte)

type bodyDumpResponseWriter struct {
    io.Writer
    http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
    w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
    w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
    return w.ResponseWriter.(http.Hijacker).Hijack()
}

func dump(handler bodyDumpHandler) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) (err error) {
            // Request
            var reqBody []byte
            if c.Request().Body != nil { // Read
                reqBody, _ = ioutil.ReadAll(c.Request().Body)
            }
            c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

            // Response
            resBody := new(bytes.Buffer)
            mw := io.MultiWriter(c.Response().Writer, resBody)
            writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
            c.Response().Writer = writer

            err = next(c)

            // Callback
            handler(c, reqBody, resBody.Bytes())

            return
        }
    }
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

func BodyDumpRawWithInterval(skipper map[string]bool) echo.MiddlewareFunc {
    return BodyDumpWithInterval(BodyDumpConfig{
        WithRequestHeaders:  false,
        WithResponseHeaders: false,
        Skipper:             skipper,
    })
}

func logBuffer(config BodyDumpConfig, c echo.Context, reqBody, resBody []byte) (buffer bytes.Buffer) {
    buffer.WriteString(fmt.Sprintf("%s [%s] %s", time.Now().Format("2006-01-02 15:04:05.00000"), c.Request().Method, c.Request().RequestURI))
    if config.WithRequestHeaders {
        buffer.WriteRune('\n')
        buffer.WriteString("-----Request Header-----")
        buffer.WriteRune('\n')
        for k, v := range c.Request().Header {
            buffer.WriteString(k)
            buffer.WriteString(": ")
            buffer.WriteString(strings.Join(v, ","))
            buffer.WriteRune('\n')
        }
    }
    if len(reqBody) > 0 {
        buffer.WriteString("\n[REQ] ")
        buffer.Write(reqBody)
    }
    if config.WithResponseHeaders {
        buffer.WriteRune('\n')
        buffer.WriteString(fmt.Sprintf("-----Response[%d] Header-----", c.Response().Status))
        buffer.WriteRune('\n')
        for k, v := range c.Response().Header() {
            buffer.WriteString(k)
            buffer.WriteString(": ")
            buffer.WriteString(strings.Join(v, ","))
            buffer.WriteRune('\n')
        }
    }
    if len(resBody) > 0 {
        buffer.WriteString("\n[RES] ")
        buffer.Write(resBody)
    }
    buffer.WriteRune('\n')
    return
}

// BodyDumpWithConfig 保存请求/返回日志
func BodyDumpWithConfig(config BodyDumpConfig) echo.MiddlewareFunc {
    return dump(func(c echo.Context, reqBody, resBody []byte) {
        b := logBuffer(config, c, reqBody, resBody)
        log.Info(b.String())
    })
}

// BodyDumpWithInterval 保存请求/返回日志(定时删除 7day)
func BodyDumpWithInterval(config BodyDumpConfig) echo.MiddlewareFunc {
    return dump(func(c echo.Context, reqBody, resBody []byte) {
        if config.Skipper[c.Path()] {
            return
        }
        now := time.Now()
        d := "runtime/logs/api"
        p := filepath.Join(d, fmt.Sprintf("%s.log", now.Format(carbon.DateLayout)))
        _ = utils.NewFile(p).CreateDirectoryIfNotExist()
        // 删除七天之前的日志文件
        _ = filepath.Walk(d, func(path string, info fs.FileInfo, err error) error {
            if info == nil || info.IsDir() {
                return nil
            }
            if now.Sub(info.ModTime()).Seconds() > 604800.0 {
                _ = os.Remove(path)
            }
            return nil
        })
        b := logBuffer(config, c, reqBody, resBody)
        f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
        if err != nil {
            log.Errorf("[BODY LOG] 文件打开失败: %s", err)
        }
        defer func(f *os.File) {
            _ = f.Close()
        }(f)

        _, _ = f.Write(b.Bytes())
    })
}
