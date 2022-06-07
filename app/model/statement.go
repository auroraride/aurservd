// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-06
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type StatementBillReq struct {
    End string `json:"end" validate:"required" trans:"账单截止日期"`
    ID  uint64 `json:"id" validate:"required" trans:"企业ID"`
}
