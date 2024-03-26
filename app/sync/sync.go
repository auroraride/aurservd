// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-28
// Based on aurservd by liasica, magicrolan@qq.com.

package sync

import (
	"github.com/auroraride/adapter/defs/batdef"
	"github.com/auroraride/adapter/sync"
	"github.com/redis/go-redis/v9"

	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
)

func Run() {
	// TODO 后面更换成默认DB
	rdb := redis.NewClient(&redis.Options{
		Addr: ar.Config.Redis.Address,
	})

	// 同步电池流转
	go func() {
		sync.New[batdef.BatteryFlow](
			rdb,
			ar.Config.Environment,
			sync.StreamBatteryFlow,
			func(data []*batdef.BatteryFlow) {
				go service.NewBatteryBms().Sync(data)
			},
		).Run()
	}()

	// TODO 同步电池信息, 是否有必要?
	// service.NewBattery().Sync(msg.(*cabdef.BatteryMessage))
}
