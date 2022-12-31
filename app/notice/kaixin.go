// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-25
// Based on aurservd by liasica, magicrolan@qq.com.

package notice

import (
    "bytes"
    "fmt"
    "github.com/auroraride/adapter"
    "github.com/auroraride/adapter/codec"
    "github.com/auroraride/adapter/tcp"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/goccy/go-json"
    log "github.com/sirupsen/logrus"
)

var (
    cabinetSyncBytes  = []byte(`"type":"cabinet_sync"`)
    exchangeStepBytes = []byte(`"type":"exchange_step"`)
)

func kaixin() {
    s := tcp.NewServer(ar.Config.Adapter.Kaixin.TcpBind, log.StandardLogger(), &codec.HeaderLength{}, func(b []byte) {
        fmt.Println(string(b))

        switch true {
        case bytes.Contains(b, cabinetSyncBytes):
            req := new(adapter.Data[adapter.CabinetSyncData])
            err := json.Unmarshal(b, req)
            if err != nil {
                log.Errorf("同步消息解析失败: %v", err)
                return
            }
            service.NewCabinet().Sync(req.Value)

        case bytes.Contains(b, exchangeStepBytes):
            req := new(adapter.Data[adapter.ExchangeStepResult])
            err := json.Unmarshal(b, req)
            if err != nil {
                log.Errorf("同步消息解析失败: %v", err)
                return
            }
            service.NewIntelligentCabinet().ExchangeStep(req.Value)
        }
    })
    log.Fatal(s.Run())
}
