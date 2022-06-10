// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-10
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type RiderMgrDepositReq struct {
    Amount float64 `json:"amount"` // 押金金额
    ID     uint64  `json:"id" validate:"required" trans:"骑手ID"`
}

type RiderMgrModifyReq struct {
    ID         uint64            `json:"id" validate:"required" trans:"骑手ID"`
    Phone      *string           `json:"phone"`                      // 手机号
    AuthStatus *PersonAuthStatus `json:"authStatus" enums:"0,1,2,3"` // 认证状态 0:未认证 1:认证中 2:已认证 3:认证失败
    Contact    *RiderContact     `json:"contact"`                    // 联系人
}
