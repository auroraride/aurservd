// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "bytes"
    "fmt"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/goccy/go-json"
    "github.com/golang-module/carbon/v2"
    "os"
    "sync"
    "time"
)

type Logger struct {
    mu   sync.Mutex
    name string
}

func NewLogger(name string) *Logger {
    return &Logger{
        mu:   sync.Mutex{},
        name: name,
    }
}

func (l *Logger) Write(message any) {
    buffer := &bytes.Buffer{}
    buffer.WriteString(time.Now().Format(carbon.TimeLayout))
    buffer.WriteString(" ")
    switch message.(type) {
    case string:
        buffer.WriteString(message.(string))
        break
    case []byte:
        buffer.Write(message.([]byte))
        break
    default:
        b := &bytes.Buffer{}
        encoder := json.NewEncoder(b)
        encoder.SetEscapeHTML(false)
        _ = encoder.Encode(message)
        buffer.Write(b.Bytes())
        break
    }
    if buffer.Bytes()[len(buffer.Bytes())-1] != '\n' {
        buffer.WriteRune('\n')
    }

    // 写入文件
    path := fmt.Sprintf("runtime/logs/%s/%s.log", l.name, time.Now().Format(carbon.DateLayout))
    _ = utils.NewFile(path).CreateDirectoryIfNotExist()

    file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

    defer func(file *os.File) {
        _ = file.Close()
    }(file)

    _, _ = file.Write(buffer.Bytes())
}
