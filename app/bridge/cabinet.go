// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-25
// Based on aurservd by liasica, magicrolan@qq.com.

package bridge

import (
    "fmt"
    "github.com/auroraride/adapter/codec"
    "github.com/auroraride/adapter/tcp"
    "github.com/auroraride/aurservd/internal/ar"
    log "github.com/sirupsen/logrus"
)

func cabinet() {
    addr := ar.Config.Adapter.Cabinet
    s := tcp.NewServer(addr, log.StandardLogger(), &codec.HeaderLength{}, func(b []byte) {
        fmt.Println(string(b))
    })
    log.Fatal(s.Run())
}
