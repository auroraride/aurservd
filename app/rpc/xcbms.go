// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-05
// Based on aurservd by liasica, magicrolan@qq.com.

package rpc

import (
    "github.com/auroraride/adapter/defs/xcdef/proto"
    "github.com/auroraride/adapter/rpc"
    "github.com/auroraride/aurservd/internal/ar"
    "go.uber.org/zap"
    "google.golang.org/grpc"
)

var (
    XcbmsClient proto.BatteryClient
)

func createXcClient() {
    err := rpc.NewClient(ar.Config.Rpc.Xcbms.Server, func(conn *grpc.ClientConn) {
        XcbmsClient = proto.NewBatteryClient(conn)
    })
    if err != nil {
        zap.L().Error("xcbms rpc连接失败", zap.Error(err))
    }
}
