// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-11-29
// Based on shiguangju by liasica, magicrolan@qq.com.

package logger

import (
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    log "github.com/sirupsen/logrus"
    "path"
    "time"
)

type Config struct {
    Color    bool   // 是否启用日志颜色
    Level    string // 日志等级
    Age      int    // 日志保存时间（小时）
    Json     bool   // 日志以json格式保存
    RootPath string
}

func LoadWithConfig(cfg Config) {
    rotateOptions := []rotatelogs.Option{
        rotatelogs.WithRotationTime(time.Hour * 24),
    }
    rotateOptions = append(rotateOptions, rotatelogs.WithMaxAge(time.Duration(cfg.Age)*time.Hour))
    // rotateOptions = append(rotateOptions, rotatelogs.ForceNewFile())
    w, err := rotatelogs.New(path.Join("runtime/logs", "%Y-%m-%d.log"), rotateOptions...)
    if err != nil {
        log.Errorf("rotatelogs 初始化失败: %v", err)
        panic(err)
    }

    consoleFormatter := LogFormat{EnableColor: cfg.Color, Console: true, RootPath: cfg.RootPath}
    fileFormatter := LogFormat{EnableColor: false, SaveJson: cfg.Json, RootPath: cfg.RootPath}
    log.AddHook(NewLocalHook(w, consoleFormatter, fileFormatter, GetLogLevel(cfg.Level)...))
    log.SetReportCaller(true)
}
