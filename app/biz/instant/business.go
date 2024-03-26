// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-25, by liasica

package instant

import (
	"fmt"
	"net"

	"github.com/auroraride/adapter/rpc/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/auroraride/aurservd/internal/ar"
)

type businessServer struct {
	pb.UnimplementedBusinessServer

	exchange *exchange
}

var businessInstance *businessServer

func NewBusinessServer() *businessServer {
	return businessInstance
}

// 启动业务RPC服务
func (s *businessServer) run() {
	businessInstance = &businessServer{
		exchange: &exchange{pipeline: make(map[string]chan *pb.CabinetExchangeResponse)},
	}
	bind := ar.Config.GRPC.Business
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		zap.L().Fatal("BIZ RPC启动失败", zap.Error(err))
		return
	}

	fmt.Printf("⇨ biz gRPC server started on %s\n", bind)

	serv := grpc.NewServer()
	pb.RegisterBusinessServer(
		serv,
		businessInstance,
	)
	err = serv.Serve(lis)
	if err != nil {
		zap.L().Fatal("RPC启动失败", zap.Error(err))
	}
}
