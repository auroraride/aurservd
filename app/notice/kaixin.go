// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-25
// Based on aurservd by liasica, magicrolan@qq.com.

package notice

import (
    "github.com/auroraride/adapter"
    "github.com/auroraride/adapter/codec"
    "github.com/auroraride/adapter/tcp"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    log "github.com/sirupsen/logrus"
)

func kaixin() {
    s := tcp.NewServer(ar.Config.Adapter.Kaixin.TcpBind, log.StandardLogger(), &codec.HeaderLength{}, func(b []byte) {
        // fmt.Println(string(b))

        t, message, err := adapter.Unpack(b)
        if err != nil {
            log.Errorf("同步消息解析失败: %v", err)
            return
        }

        switch t {

        case adapter.TypeCabinet:
            // 同步电柜
            service.NewCabinet().Sync(message.(*adapter.CabinetMessage))

        case adapter.TypeBattery:
            // 同步电池信息
            service.NewBattery().Sync(message.(*adapter.BatteryMessage))

        case adapter.TypeExchangeStep:
            // 换电步骤
            service.NewIntelligentCabinet().ExchangeStep(message.(*adapter.ExchangeStepMessage))
        }

    })
    log.Fatal(s.Run())
}
