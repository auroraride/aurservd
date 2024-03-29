// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-17
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type WalletOverview struct {
	Balance float64 `json:"balance"` // 账户余额
	Points  int64   `json:"points"`  // 积分数量
	Coupons int     `json:"coupons"` // 可使用优惠券数量
	Deposit float64 `json:"deposit"` // 已缴纳押金

	DepositType *uint8 `json:"depositType,omitempty"` // 押金类型  1:芝麻免押 2:微信支付分免押 3:合同押金 4:支付押金
}
