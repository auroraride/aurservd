// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-17
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import "github.com/auroraride/aurservd/app/model"

type walletService struct {
    *BaseService
}

func NewWallet(params ...any) *walletService {
    return &walletService{
        BaseService: newService(params...),
    }
}

func (s *walletService) Overview() model.WalletOverview {
    return model.WalletOverview{
        Balance: 0,
        Points:  s.entRider.Points,
        Coupons: len(NewCoupon().QueryEffective(s.rider.ID)),
        Deposit: NewRider().DepositPaid(s.rider.ID).Deposit,
    }
}
