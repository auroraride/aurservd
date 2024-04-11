package rpc

import (
	"context"
	"errors"

	"github.com/liasica/edocseal/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = "10.10.10.240:1111"

// Sgin 签约
func Sgin(ctx context.Context, req *pb.ContractSignRequest) (string, error) {
	gc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("rpc连接失败", zap.Error(err))
		return "", errors.New("rpc连接失败")
	}
	c := pb.NewContractClient(gc)
	res, err := c.Sign(ctx, req)
	if err != nil {
		zap.L().Error("请求失败", zap.Error(err))
		return "", errors.New("请求失败")
	}

	if res.Status != pb.ContractSignStatus_SUCCESS {
		zap.L().Error("请求失败", zap.String("message", res.Message))
		return "", errors.New(res.Message)
	}

	return res.Url, nil
}

// Create 创建合同
func Create(ctx context.Context, templateId string, values map[string]*pb.ContractFromField) (request *pb.ContractCreateResponse, err error) {
	gc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("rpc连接失败", zap.Error(err))
		return nil, errors.New("rpc连接失败")
	}
	c := pb.NewContractClient(gc)
	request, err = c.Create(ctx, &pb.ContractCreateRequest{
		TemplateId: templateId,
		Values:     values,
	})
	if err != nil {
		zap.L().Error("请求失败", zap.Error(err))
		return nil, errors.New("请求失败")
	}

	return request, nil
}
