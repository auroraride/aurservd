// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-18
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "fmt"

type ImportRiderReq struct {
    Name  string `json:"name" validate:"required" trans:"姓名"`
    Phone string `json:"phone" validate:"required" trans:"电话"`
    Plan  string `json:"plan" validate:"required" trans:"订阅"`
    Days  string `json:"days" validate:"required" trans:"订阅天数"`
    Model string `json:"model" validate:"required" trans:"电池型号"`
    City  string `json:"city" validate:"required" trans:"城市"`
    Store string `json:"store" validate:"required" trans:"激活门店"`
    End   string `json:"end" validate:"required,datetime=2006-01-02" trans:"结束日期"`

    EmployeeID uint64 `json:"employeeId"` // 选择店员ID
}

func (r *ImportRiderReq) String() string {
    return fmt.Sprintf("%s:%s -> 订阅:%s, 订阅天数:%s, 电池型号:%s, 城市:%s, 激活门店:%s, 结束日期:%s",
        r.Name,
        r.Phone,
        r.Plan,
        r.Days,
        r.Model,
        r.City,
        r.Store,
        r.End,
    )
}
