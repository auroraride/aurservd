// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

type allocateBatteryService struct {
	*BaseService
}

func NewAllocateBattery(params ...any) *allocateBatteryService {
	return &allocateBatteryService{
		BaseService: newService(params...),
	}
}
