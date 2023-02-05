// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

type batteryXcService struct {
    *BaseService
}

func NewBatteryXc(params ...any) *batteryXcService {
    return &batteryXcService{
        BaseService: newService(params...),
    }
}
