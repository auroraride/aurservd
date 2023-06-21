// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-21
// Based on aurservd by liasica, magicrolan@qq.com.

package service

type enterprisePriceService struct {
	*BaseService
}

func NewEnterprisePrice(params ...any) *enterprisePriceService {
	return &enterprisePriceService{
		BaseService: newService(params...),
	}
}
