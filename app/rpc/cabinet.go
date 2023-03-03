// Copyright (C) liasica. 2023-present.
//
// Created at 2023-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package rpc

import (
    "context"
    "github.com/auroraride/adapter"
    "github.com/auroraride/adapter/rpc"
    "github.com/auroraride/adapter/rpc/pb"
    "go.uber.org/zap"
    "sync"
)

var (
    // 客户端列表
    // brand + [intelligent] => pb.CabinetClient
    cabinets sync.Map
)

func CabinetKey(br adapter.CabinetBrand, intelligent bool) string {
    var sf string
    if !intelligent {
        sf = "-non-intelligent"
    }

    return br.RpcName() + sf
}

func GetCabinet(key string) pb.CabinetClient {
    if c, ok := cabinets.Load(key); ok {
        return c.(pb.CabinetClient)
    }

    addr := serverAddress(key)
    if addr == "" {
        return nil
    }

    c, err := rpc.NewClient(addr, pb.NewCabinetClient)
    if err != nil {
        zap.L().Error(key+" rpc连接失败", zap.Error(err))
    }

    cabinets.Store(key, c)

    return c
}

func CabinetBatch(key string, req *pb.CabinetBatchRequest) *pb.CabinetBatchResponse {
    c := GetCabinet(key)
    if c == nil {
        return nil
    }

    res, _ := c.Batch(context.Background(), req)
    return res
}
