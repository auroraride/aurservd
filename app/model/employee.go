// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type EmployeeItem struct {
    ID    uint64 `json:"id"`
    Name  string `json:"name"`  // 店员名称
    Phone string `json:"phone"` // 店员电话
}
