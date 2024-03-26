// Copyright (C) liasica. 2023-present.
//
// Created at 2023-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package rpc

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/log"
	"github.com/auroraride/adapter/rpc"
	"github.com/auroraride/adapter/rpc/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/auroraride/aurservd/app/model"
)

var (
	// 客户端列表
	// brand + [intelligent] => pb.CabinetClient
	cabinetClients sync.Map
)

func CabinetKey(br adapter.CabinetBrand, intelligent bool) string {
	var sf string
	if !intelligent {
		sf = "-non-intelligent"
	}

	return br.RpcName() + sf
}

func GetCabinet(key string) pb.CabinetClient {
	if c, ok := cabinetClients.Load(key); ok {
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

	cabinetClients.Store(key, c)

	return c
}

func CabinetSync(key string, req *pb.CabinetSyncRequest) *pb.CabinetSyncResponse {
	c := GetCabinet(key)
	if c == nil {
		return nil
	}

	res, _ := c.Sync(context.Background(), req)
	return res
}

func CabinetBiz(key string, req *pb.CabinetBizRequest) *pb.CabinetBizResponse {
	c := GetCabinet(key)
	if c == nil {
		return nil
	}

	res, _ := c.Biz(context.Background(), req)
	return res
}

func CabinetInterrupt(key string, req *pb.CabinetInterruptRequest) *pb.CabinetBizResponse {
	c := GetCabinet(key)
	if c == nil {
		return nil
	}

	res, _ := c.Interrupt(context.Background(), req)
	return res
}

// BinOperateResultFunc 步骤操作结果回调
// 第一个参数：步骤结果
// 第二个参数：是否终止
type BinOperateResultFunc func(*pb.CabinetExchangeResponse, bool)

// CabinetExchange 换电
func CabinetExchange(key string, user *adapter.User, req *pb.CabinetExchangeRequest, bor BinOperateResultFunc) (err error) {
	zap.L().Info("换电请求", user.ZapField(), log.Payload(req))

	c := GetCabinet(key)
	if c == nil {
		return errors.New("无法连接到柜控服务器")
	}
	var stream pb.Cabinet_ExchangeClient
	stream, err = c.Exchange(metadata.NewOutgoingContext(context.Background(), metadata.MD{"user": {user.Type.String(), user.ID}}), req)
	if err != nil {
		return
	}

	var (
		res      *pb.CabinetExchangeResponse
		stop     bool
		laststep uint32
		message  string
	)
	for {
		res, err = stream.Recv()

		// 连接结束
		if err == io.EOF {
			stop = true
			zap.L().Info("换电请求已结束", user.ZapField(), user.ZapField(), log.Payload(res))
			// 判定返回步骤是否最后一步，若非最后一步代表换电失败
			if res != nil {
				if res.Step < 4 {
					err = errors.New("网络异常，换电失败")
				} else {
					err = nil
				}
			}
		}

		if err != nil {
			zap.L().Error("换电请求失败", user.ZapField(), zap.Error(err))
			stop = true
			message = err.Error()
		}

		if res != nil {
			laststep = res.Step
			stop = res.Step == model.ExchangeStepPutOut.Uint32()
		}

		if res == nil {
			stop = true
			res = &pb.CabinetExchangeResponse{Success: false, Step: laststep + 1, Message: message}
		}

		// 步骤结果回调
		bor(res, stop)

		if stop {
			return
		}
	}
}
