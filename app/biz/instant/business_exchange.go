// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-25, by liasica

package instant

import (
	"github.com/auroraride/adapter/rpc/pb"
)

// 换电业务
type exchange struct {
	// 柜控消息管道
	pipeline map[string]chan *pb.CabinetExchangeResponse
}

func (s *businessServer) Exchange(req *pb.BusinessExchangeRequest, stream pb.Business_ExchangeServer) error {
	// minsoc := cache.Float64(model.SettingExchangeMinBatteryKey)
	// timeout := model.CabinetBusinessStepTimeout

	// adapterPB.CabinetExchangeRequest{
	// 	Uuid:   ,
	// 	Serial:  data.Data.Cabinet.Serial,
	// 	Battery: "XCE16A1024010001",
	// 	// Battery: "XCE16A1024010004",
	// 	// Battery: "XCE16A1024010002",
	// 	Expires: 12000000,
	// 	Timeout: 120,
	// 	Minsoc:  0,
	// }

	return nil
}

// func (s *businessServer) exchangeRequest() {
// 	// 连接服务器
// 	conn, err := grpc.Dial(":15002", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("net.Connect err: %v", err)
// 	}
// 	defer func(conn *grpc.ClientConn) {
// 		_ = conn.Close()
// 	}(conn)
//
// 	// 建立gRPC连接
// 	grpcClient := cpb.NewCabinetClient(conn)
//
// 	// 获取流
// 	ctx := metadata.NewOutgoingContext(context.Background(), metadata.MD{"user": {"rider", "185013545512"}})
// 	stream, err := grpcClient.Exchange(ctx, &cpb.CabinetExchangeRequest{
// 		Uuid:    data.Data.Uuid,
// 		Serial:  data.Data.Cabinet.Serial,
// 		Battery: "XCE16A1024010001",
// 		// Battery: "XCE16A1024010004",
// 		// Battery: "XCE16A1024010002",
// 		Expires: 12000000,
// 		Timeout: 120,
// 		Minsoc:  0,
// 	})
// 	if err != nil {
// 		log.Fatalf("Call Exchange err: %v", err)
// 	}
//
// 	var res *pb.CabinetExchangeResponse
// 	for {
// 		res, err = stream.Recv()
// 		if err == io.EOF {
// 			log.Println("请求已结束")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatalf("Conversations get stream err: %v", err)
// 		}
// 		// 打印返回值
// 		log.Printf("步骤结果 >> %#v", res)
// 	}
// }
