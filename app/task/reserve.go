// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/cache"
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
			<-ticker.C
			service.NewReserve().Timeout()
		}
	}
}
