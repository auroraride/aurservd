// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/pkg/cache"
    "time"
)

type reserveTask struct {
    max time.Duration
}

func NewReserve() *reserveTask {
    return &reserveTask{
        max: time.Duration(cache.Int(model.SettingReserveDuration)),
    }
}

func (t *reserveTask) Start() {
    for {
        service.NewReserve().Timeout()
        // 每隔一分钟检查一次
        time.Sleep(1 * time.Minute)
    }
}
