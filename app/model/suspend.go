// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type SuspendReq struct {
    ID     uint64 `json:"id" validate:"required" trans:"订阅ID"`
    Remark string `json:"remark" validate:"required" trans:"备注"`
}
