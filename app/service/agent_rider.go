// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

type agentRiderService struct {
	*BaseService
}

func NewAgentRider(params ...any) *agentRiderService {
	return &agentRiderService{
		BaseService: newService(params...),
	}
}
