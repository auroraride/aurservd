// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-17
// Based on aurservd by liasica, magicrolan@qq.com.

package service

type pointLogService struct {
    *BaseService
}

func NewPointLog(params ...any) *pointLogService {
    return &pointLogService{
        BaseService: newService(params...),
    }
}

func (s *pointLogService) Create() {

}