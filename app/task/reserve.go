// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/cache"
    "time"
)

type reserveTask struct {
    max time.Duration
}

func NewReserve() *reserveTask {
    return &reserveTask{
        max: time.Duration(cache.Int(model.SettingReserveDurationKey)),
    }
}

func (t *reserveTask) Start() {
    if ar.Config.Task.Reserve {
        ticker := time.NewTicker(1 * time.Minute)
        for {
            select {
            case <-ticker.C:
                service.NewReserve().Timeout()
                break
            }
        }
    }

}
