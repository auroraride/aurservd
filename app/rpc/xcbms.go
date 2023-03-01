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
)

var (
    bmsXcClient pb.BatteryClient
)

func createXcClient() {
    var err error
    bmsXcClient, err = rpc.NewClient(ar.Config.Rpc.BmsXc.Server, pb.NewBatteryClient)
    if err != nil {
        zap.L().Error("xcbms rpc连接失败", zap.Error(err))
    }
}

func tryCreateXcClient() bool {
    if bmsXcClient == nil {
        createXcClient()
    }
    return bmsXcClient != nil
}

func XcBmsBatch(ctx context.Context, req *pb.BatteryBatchRequest) (res *pb.BatteryBatchResponse, err error) {
    if !tryCreateXcClient() {
        return
    }
    return bmsXcClient.Batch(ctx, req)
}

func XcBmsSample(ctx context.Context, req *pb.BatterySnRequest) (res *pb.BatterySampleResponse, err error) {
    if !tryCreateXcClient() {
        return
    }
    return bmsXcClient.Sample(ctx, req)
}

func XcBmsFaultList(ctx context.Context, req *pb.BatteryFaultListRequest) (res *pb.BatteryFaultListResponse, err error) {
    if !tryCreateXcClient() {
        return
    }
    return bmsXcClient.FaultList(ctx, req)
}

func XcBmsFaultOverview(ctx context.Context, req *pb.BatterySnRequest) (res *pb.BatteryFaultOverviewResponse, err error) {
    if !tryCreateXcClient() {
        return
    }
    return bmsXcClient.FaultOverview(ctx, req)
}

func XcBmsStatistics(ctx context.Context, req *pb.BatterySnRequest) (res *pb.BatteryStatisticsResponse, err error) {
    if !tryCreateXcClient() {
        return
    }
    return bmsXcClient.Statistics(ctx, req)
}

func XcBmsPosition(ctx context.Context, req *pb.BatteryPositionRequest) (res *pb.BatteryPositionResponse, err error) {
    if !tryCreateXcClient() {
        return
    }
    return bmsXcClient.Position(ctx, req)
}
