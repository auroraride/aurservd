// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-25
// Based on aurservd by liasica, magicrolan@qq.com.

package sync

import (
    "github.com/auroraride/adapter/defs/cabdef"
    "github.com/auroraride/adapter/sync"
    "github.com/auroraride/adapter/zlog"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/go-redis/redis/v9"
    "go.uber.org/zap"
)

func runKxcab() {
    logger := zlog.StandardLogger().GetLogger().WithOptions(zap.AddCallerSkip(-2))

    // TODO 后面更换成默认DB
    rdb := redis.NewClient(&redis.Options{
        Addr: ar.Config.Redis.Address,
    })

    // 同步电柜
    go func() {
        sync.New[cabdef.CabinetMessage](
            rdb,
            ar.Config.Environment,
            sync.StreamCabinet,
            func(data *cabdef.CabinetMessage) {
                service.NewCabinet().Sync(data)
            },
            logger,
        ).Run()
    }()

    // 同步换电步骤
    go func() {
        sync.New[cabdef.ExchangeStepMessage](
            rdb,
            ar.Config.Environment,
            sync.StreamExchange,
            func(data *cabdef.ExchangeStepMessage) {
                service.NewIntelligentCabinet().ExchangeStepSync(data)
            },
            logger,
        ).Run()
    }()

    // TODO 同步电池信息, 是否有必要?
    // service.NewBattery().Sync(msg.(*cabdef.BatteryMessage))
}
