// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-25
// Based on aurservd by liasica, magicrolan@qq.com.

package bridge

import (
    "fmt"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/bridge"
    "github.com/auroraride/bridge/pb"
    log "github.com/sirupsen/logrus"
)

func cabinet() {
    log.Fatal(bridge.NewCabinet(log.StandardLogger()).RunServer(ar.Config.Bridge.Cabinet, func(data *pb.CabinetSyncRequest) {
        // 保存
        fmt.Println(data)
    }))
}
