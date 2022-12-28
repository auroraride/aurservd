// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-18
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
    "github.com/goccy/go-json"
    log "github.com/sirupsen/logrus"
)

type logTool struct {
}

func NewLog() *logTool {
    return &logTool{}
}

func (*logTool) Infof(format string, args ...interface{}) {
    data := make([]interface{}, len(args))
    for i, param := range args {
        b, _ := json.Marshal(param)
        data[i] = string(b)
    }
    log.Infof(format, data...)
}
