// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-04
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type PrepaymentOverview struct {
    Times   int     `json:"times"`   // 充值次数
    Balance float64 `json:"balance"` // 当前余额
    Amount  float64 `json:"amount"`  // 总充值金额
    Cost    float64 `json:"cost"`    // 已使用金额
}

type PrepaymentListReq struct {
    PaginationReq

    EnterpriseID uint64 `json:"enterpriseId" query:"enterpriseId"`
}

type PrepaymentListRes struct {
    Amount float64 `json:"amount"` // 金额
    Name   string  `json:"name"`   // 操作人
    Time   string  `json:"time"`   // 时间
    Remark string  `json:"remark"` // 备注信息
}
