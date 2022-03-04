// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// ShiguangjuIVRRes 时光驹外呼催缴返回
type ShiguangjuIVRRes struct {
    Items []*ShiguangjuIVRItem `json:"items"`
}

// ShiguangjuIVRItem 时光驹外呼催缴数据
type ShiguangjuIVRItem struct {
    Name    string `json:"name"`    // 姓名
    Phone   string `json:"phone"`   // 手机号
    Product string `json:"product"` // 套餐
    Status  bool   `json:"status"`  // 外呼状态
}
