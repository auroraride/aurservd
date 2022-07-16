// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-03
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// ExchangeCabinet 电柜换电
type ExchangeCabinet struct {
    Alternative bool                    `json:"alternative"` // 是否使用备选方案
    Info        *RiderCabinetOperating  `json:"info"`        // 换电信息
    Result      *RiderCabinetOperateRes `json:"result"`      // 换电结果
}

type ExchangeStoreReq struct {
    Code string `json:"code" validate:"required,startswith=STORE:"` // 二维码
}

type ExchangeStoreRes struct {
    Model     string `json:"model"`     // 电池型号
    StoreName string `json:"storeName"` // 门店名称
    Time      int64  `json:"time"`      // 时间戳
    UUID      string `json:"uuid"`      // 编码
}

type ExchangeOverview struct {
    Times int `json:"times"` // 换电次数
    Days  int `json:"days"`  // 换电天数
}

type ExchangeLogReq struct {
    PaginationReq
}

type ExchangeLogBinInfo struct {
    EmptyIndex int `json:"emptyIndex"` // 空电仓位
    FullIndex  int `json:"fullIndex"`  // 满电仓位
}

type ExchangeRiderListRes struct {
    ID      uint64             `json:"id"`
    Name    string             `json:"name"`    // 门店或电柜名称
    Type    string             `json:"type"`    // 门店或电柜
    Time    string             `json:"time"`    // 换电时间
    Success bool               `json:"success"` // 是否成功
    City    City               `json:"city"`    // 城市
    BinInfo ExchangeLogBinInfo `json:"binInfo"` // 仓位信息
}

type ExchangeListReq struct {
    PaginationReq

    Aimed   uint8   `json:"aimed" query:"aimed"`     // 筛选对象 0:全部 1:个签 2:团签
    Start   *string `json:"start" query:"start"`     // 筛选开始日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
    End     *string `json:"end" query:"end"`         // 筛选结束日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
    Keyword *string `json:"keyword" query:"keyword"` // 筛选骑手姓名或电话
}

type ExchangeEmployeeListRes struct {
    ID         uint64           `json:"id"`
    Name       string           `json:"name"`                 // 骑手姓名
    Phone      string           `json:"phone"`                // 骑手电话
    Time       string           `json:"time"`                 // 换电时间
    Model      string           `json:"model"`                // 电池型号
    Enterprise *EnterpriseBasic `json:"enterprise,omitempty"` // 团签企业, 个签无此字段
    Plan       *Plan            `json:"plan,omitempty"`       // 骑士卡, 团签无此字段
}

type ExchangeManagerListReq struct {
    ExchangeListReq
    Target      uint8  `json:"target" query:"target"`                         // 换电类别 0:全部 1:电柜 2:门店
    CityID      uint64 `json:"cityId" query:"cityId"`                         // 城市
    Employee    string `json:"employee" query:"employee"`                     // 筛选店员(手机号或姓名)
    Start       string `json:"start" query:"start"`                           // 开始时间
    End         string `json:"end" query:"end"`                               // 结束时间
    Alternative uint8  `json:"alternative" query:"alternative" enums:"0,1,2"` // 换电方案 0:全部 1:满电 2:非满电
    CabinetID   uint64 `json:"cabinetId" query:"cabinetId"`                   // 选择电柜ID
    StoreID     uint64 `json:"storeId" query:"storeId"`                       // 选择门店ID
    Serial      string `json:"serial" query:"serial"`                         // 电柜编号
    Status      uint8  `json:"status" query:"status" enums:"0,1,2"`           // 换电状态 0:全部 1:成功 2:失败
}

type ExchangeManagerListRes struct {
    ID          uint64            `json:"id"`
    Name        string            `json:"name"`                 // 骑手姓名
    Phone       string            `json:"phone"`                // 骑手电话
    Time        string            `json:"time"`                 // 换电时间
    Model       string            `json:"model"`                // 电池型号
    Alternative bool              `json:"alternative"`          // 换电方案 `true`非满电 `false`满电, 只有`true`的时候才显示为`非满电`
    Enterprise  *EnterpriseBasic  `json:"enterprise,omitempty"` // 团签企业, 个签无此字段
    Store       *Store            `json:"store,omitempty"`      // 门店, 电柜换电无此字段
    Cabinet     *CabinetBasicInfo `json:"cabinet,omitempty"`    // 电柜, 门店换电无此字段
    City        City              `json:"city"`                 // 城市
    Success     bool              `json:"success"`              // 是否成功
    Full        string            `json:"full,omitempty"`       // 满电仓位信息, 门店换电不存在
    Empty       string            `json:"empty,omitempty"`      // 空仓位信息, 门店换电不存在
    Error       string            `json:"error,omitempty"`      // 换电失败原因
}
