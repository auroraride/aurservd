// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    BusinessAimedAll        uint8 = iota // 全部
    BusinessAimedPersonal                // 个签
    BusinessAimedEnterprise              // 团签
)

type BusinessSubscribeID struct {
    SubscribeID uint64 `json:"subscribeId" validate:"required" trans:"订阅ID"`
}

type BusinessListReq struct {
    PaginationReq

    EmployeeID   uint64  `json:"employeeId" query:"employeeId"`     // 店员ID, 店员端请求忽略此参数
    EnterpriseID uint64  `json:"enterpriseId" query:"enterpriseId"` // 企业ID
    Aimed        uint8   `json:"aimed" query:"aimed"`               // 筛选业务对象 0:全部 1:个签 2:团签
    Start        *string `json:"start" query:"start"`               // 筛选开始日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
    End          *string `json:"end" query:"end"`                   // 筛选结束日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
    Keyword      *string `json:"keyword" query:"keyword"`           // 筛选骑手姓名或电话
    // 筛选业务类别 active:激活 pause:寄存 continue:结束寄存 unsubscribe:退订
    Type *string `json:"type" enums:"active,pause,continue,unsubscribe" query:"type"`
}

type BusinessEmployeeListRes struct {
    ID         uint64           `json:"id"`
    Name       string           `json:"name"`                 // 骑手姓名
    Phone      string           `json:"phone"`                // 骑手电话
    Type       string           `json:"type"`                 // 业务类别
    Time       string           `json:"time"`                 // 业务时间
    Plan       *Plan            `json:"plan,omitempty"`       // 骑士卡, 团签无此字段
    Enterprise *EnterpriseBasic `json:"enterprise,omitempty"` // 团签企业, 个签无此字段
}

type BusinessListRes struct {
    BusinessEmployeeListRes
    Employee *Employee         `json:"employee,omitempty"` // 店员, 可能为空
    Cabinet  *CabinetBasicInfo `json:"cabinet,omitempty"`  // 电柜, 可能为空
    Store    *Store            `json:"store,omitempty"`    // 门店, 可能为空
}
