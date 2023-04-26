// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-05
// Based on aurservd by liasica, magicrolan@qq.com.

package rpc

import (
	"context"
	"sync"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/rpc"
	"github.com/auroraride/adapter/rpc/pb"
	"go.uber.org/zap"
)

var (
	// 电池客户端列表
	// brand => pb.BatteryClient
	batteryClients sync.Map
)

func GetBattery(brand adapter.BatteryBrand) pb.BatteryClient {
	key := brand.RpcName()
	if c, ok := batteryClients.Load(key); ok {
		return c.(pb.BatteryClient)
	}

	addr := serverAddress(key)
	if addr == "" {
		return nil
	}

	c, err := rpc.NewClient(addr, pb.NewBatteryClient)
	if err != nil {
		zap.L().Error(key+" rpc连接失败", zap.Error(err))
	}

	batteryClients.Store(key, c)

	return c
}

func BmsBatch(brand adapter.BatteryBrand, req *pb.BatteryBatchRequest) (res *pb.BatteryBatchResponse) {
	c := GetBattery(brand)
	if c == nil {
		return
	}

	res, _ = c.Batch(context.Background(), req)
	return
}

func BmsSample(brand adapter.BatteryBrand, req *pb.BatterySnRequest) (res *pb.BatterySampleResponse) {
	c := GetBattery(brand)
	if c == nil {
		return
	}

	res, _ = c.Sample(context.Background(), req)
	return
}

func BmsFaultList(brand adapter.BatteryBrand, req *pb.BatteryFaultListRequest) (res *pb.BatteryFaultListResponse) {
	c := GetBattery(brand)
	if c == nil {
		return
	}

	res, _ = c.FaultList(context.Background(), req)
	return
}

func BmsFaultOverview(brand adapter.BatteryBrand, req *pb.BatterySnRequest) (res *pb.BatteryFaultOverviewResponse) {
	c := GetBattery(brand)
	if c == nil {
		return
	}

	res, _ = c.FaultOverview(context.Background(), req)
	return
}

func BmsStatistics(brand adapter.BatteryBrand, req *pb.BatterySnRequest) (res *pb.BatteryStatisticsResponse) {
	c := GetBattery(brand)
	if c == nil {
		return
	}

	res, _ = c.Statistics(context.Background(), req)
	return
}

func BmsPosition(brand adapter.BatteryBrand, req *pb.BatteryPositionRequest) (res *pb.BatteryPositionResponse) {
	c := GetBattery(brand)
	if c == nil {
		return
	}

	res, _ = c.Position(context.Background(), req)
	return
}
