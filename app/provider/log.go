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
    Time     string `json:"time"`
    Response any    `json:"response"`
}

func NewLogger(name string) *Logger {
    return &Logger{
        mu:   sync.Mutex{},
        name: name,
    }
}

func (l *Logger) Write(data any) {
    path := fmt.Sprintf("runtime/logs/%s/%s.log", l.name, time.Now().Format(carbon.DateLayout))
    var buffer bytes.Buffer
    _ = utils.NewFile(path).CreateDirectoryIfNotExist()
    b, _ := jsoniter.Marshal(logData{
        Time:     time.Now().Format(carbon.DateTimeLayout),
        Response: data,
    })
    buffer.Write(b)
    buffer.WriteString("\n")

    // 写入日志文件
    l.mu.Lock()
    defer l.mu.Unlock()

    file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    _, _ = file.Write(buffer.Bytes())
    if err := file.Close(); err != nil {
        return
    }
}
