// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "bytes"
    "fmt"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/golang-module/carbon/v2"
    jsoniter "github.com/json-iterator/go"
    "os"
    "sync"
    "time"
)

type Logger struct {
    mu   sync.Mutex
    name string
}

type logData struct {
    Time   string `json:"time"`
    Times  int    `json:"times"`
    Result any    `json:"result"`
}

func NewLogger(name string) *Logger {
    return &Logger{
        mu:   sync.Mutex{},
        name: name,
    }
}

func (l *Logger) Write(times int, data any) {
    path := fmt.Sprintf("runtime/logs/%s/%s.log", l.name, time.Now().Format(carbon.DateLayout))
    _ = utils.NewFile(path).CreateDirectoryIfNotExist()

    buffer := &bytes.Buffer{}
    encoder := jsoniter.NewEncoder(buffer)
    encoder.SetEscapeHTML(false)
    _ = encoder.Encode(logData{
        Time:   time.Now().Format(carbon.DateTimeLayout),
        Times:  times,
        Result: data,
    })
    // buffer.WriteString("\n")

    // 写入日志文件
    l.mu.Lock()
    defer l.mu.Unlock()

    file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    _, _ = file.Write(buffer.Bytes())
    if err := file.Close(); err != nil {
        return
    }
}
