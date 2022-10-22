// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-18
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "fmt"

type ImportRiderFromExcel struct {
    Name  string `json:"name"`  // 姓名
    Phone string `json:"phone"` // 电话
    Plan  string `json:"plan"`  // 订阅
    Days  string `json:"days"`  // 订阅天数
    Model string `json:"model"` // 电池型号
    City  string `json:"city"`  // 城市
    Store string `json:"store"` // 激活门店
    End   string `json:"end"`   // 结束日期
}

func (r *ImportRiderFromExcel) String() string {
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

type ImportRiderCreateReq struct {
    Name       string `json:"name" validate:"required" trans:"姓名"`
    Phone      string `json:"phone" validate:"required" trans:"电话"`
    PlanID     uint64 `json:"planId" validate:"required" trans:"骑行卡ID"`
    CityID     uint64 `json:"cityId" validate:"required" trans:"城市ID"`
    StoreID    uint64 `json:"storeId" validate:"required" trans:"门店ID"`
    EmployeeID uint64 `json:"employeeId" validate:"required" trans:"店员ID"`
    End        string `json:"end" validate:"required,datetime=2006-01-02" trans:"结束日期"`
}
