// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-05
// Based on aurservd by liasica, magicrolan@qq.com.

package rpc

import (
    "context"
    "github.com/auroraride/adapter/rpc"
    "github.com/auroraride/adapter/rpc/pb"
    "github.com/auroraride/aurservd/internal/ar"
    "go.uber.org/zap"
    "google.golang.org/grpc"
)

var (
    xcBmsClient pb.BatteryClient
)

func createXcClient() {
    err := rpc.NewClient(ar.Config.Rpc.Xcbms.Server, func(conn *grpc.ClientConn) {
        xcBmsClient = pb.NewBatteryClient(conn)
    })
    if err != nil {
        zap.L().Error("xcbms rpc连接失败", zap.Error(err))
    }
}

func XcBmsBatch(ctx context.Context, req *pb.BatteryBatchRequest) (res *pb.BatteryBatchResponse, err error) {
    if xcBmsClient == nil {
        return
    }
    return xcBmsClient.Batch(ctx, req)
}

func XcBmsSample(ctx context.Context, req *pb.BatterySnRequest) (res *pb.BatterySampleResponse, err error) {
    if xcBmsClient == nil {
        return
    }
    return xcBmsClient.Sample(ctx, req)
}

func XcBmsFaultList(ctx context.Context, req *pb.BatteryFaultListRequest) (res *pb.BatteryFaultListResponse, err error) {
    if xcBmsClient == nil {
        return
    }
    return xcBmsClient.FaultList(ctx, req)
}

func XcBmsFaultOverview(ctx context.Context, req *pb.BatterySnRequest) (res *pb.BatteryFaultOverviewResponse, err error) {
    if xcBmsClient == nil {
        return
    }
    return xcBmsClient.FaultOverview(ctx, req)
}

func XcBmsStatistics(ctx context.Context, req *pb.BatterySnRequest) (res *pb.BatteryStatisticsResponse, err error) {
    if xcBmsClient == nil {
        return
    }
    return xcBmsClient.Statistics(ctx, req)
}
