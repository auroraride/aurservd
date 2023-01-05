// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-25
// Based on aurservd by liasica, magicrolan@qq.com.

package notice

import (
    "github.com/auroraride/adapter/codec"
    "github.com/auroraride/adapter/defs/cabdef"
    "github.com/auroraride/adapter/message"
    "github.com/auroraride/adapter/tcp"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    log "github.com/sirupsen/logrus"
)

func kaixin() {
    s := tcp.NewServer(ar.Config.Adapter.Kaixin.TcpBind, log.StandardLogger(), &codec.HeaderLength{}, func(b []byte) {
        // fmt.Println(string(b))

        t, msg, err := message.Unpack(b)
        if err != nil {
            log.Errorf("同步消息解析失败: %v", err)
            return
        }

        switch t {

        case message.TypeCabkitSync:
            // 同步电柜
            service.NewCabinet().Sync(msg.(*cabdef.CabinetMessage))

        case message.TypeCabkitBattery:
            // TODO 同步电池信息, 是否有必要?
            // service.NewBattery().Sync(msg.(*cabdef.BatteryMessage))

        case message.TypeCabkitExchangeStep:
            // 换电步骤
            service.NewIntelligentCabinet().ExchangeStepSync(msg.(*cabdef.ExchangeStepMessage))
        }

    })
    log.Fatal(s.Run())
}
