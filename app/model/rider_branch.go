// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-11
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type RiderBranchDetailReq struct {
    BranchID uint64 `json:"branchId" query:"branchId" validate:"required" trans:"网点ID"`
    Model    string `json:"model" query:"model" validate:"required" trans:"电柜型号"`
    Type     uint8  `json:"type" query:"type" validate:"required" enums:"1,2" trans:"类别"` // 1:门店 2:电柜
}

type RiderBranchDetailRes struct {
    Name string `json:"name"` // 门店或电柜名称
}
