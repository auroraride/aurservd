// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-17
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"github.com/auroraride/aurservd/app/model"
)

type walletService struct {
	*BaseService
}

func NewWallet(params ...any) *walletService {
	return &walletService{
		BaseService: newService(params...),
	}
}

func (s *walletService) Overview() model.WalletOverview {
	// 获取当前订阅信息
	or := NewRider().DepositOrder(s.entRider.ID)
	var deposit float64
	var depositType uint8
	if or != nil {
		deposit = or.Amount
	}
	effective, _ := NewSubscribe().QueryEffective(s.rider.ID)
	if effective != nil {
		depositType = effective.DepositType
	}
	return model.WalletOverview{
		Balance:     0,
		Points:      s.entRider.Points,
		Coupons:     len(NewCoupon().QueryEffective(s.rider.ID)),
		Deposit:     deposit,
		DepositType: model.DepositType(depositType),
	}
}
