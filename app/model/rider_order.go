// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type RiderOrder struct {
    ID     uint64         `json:"id"`
    Type   uint           `json:"type"`
    Status uint8          `json:"status"`
    Payway uint8          `json:"payway"`
    PayAt  string         `json:"payAt"`
    Amount float64        `json:"amount"`
    Plan   Plan           `json:"plan"`
    City   City           `json:"city"`
    Models []BatteryModel `json:"models"`
}
