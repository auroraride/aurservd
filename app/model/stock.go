// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    StockTypeTransfer         uint8 = iota // 调拨 (出库入库)
    StockTypeRiderObtain                   // 骑手领取电池 (出库)
    StockTypeRiderPause                    // 骑手寄存电池 (入库)
    StockTypeRiderContinue                 // 骑手结束寄存电池 (出库)
    StockTypeRiderUnSubscribe              // 骑手归还电池 (入库)
)

const (
    StockTargetPlaform uint8 = iota
    StockTargetStore         // 调拨对象 - 门店
    StockTargetCabinet       // 调拨对象 - 电柜
)

func StockNumberOfRiderBusiness(typ uint8) (num int) {
    switch typ {
    case StockTypeRiderObtain, StockTypeRiderContinue:
        num = -1
        break
    case StockTypeRiderPause, StockTypeRiderUnSubscribe:
        num = 1
        break
    }
    return
}

type StockTransferReq struct {
    Model string `json:"model,omitempty"` // 电池型号 (和`物资名称`不能同时存在, 也不能同时为空)
    Name  string `json:"name,omitempty"`  // 物资名称 (和`电池型号`不能同时存在, 也不能同时为空)
    Num   int    `json:"num"`             // 调拨数量

    OutboundID     uint64 `json:"outboundId"`                   // 调出自 0:平台
    OutBoundTarget uint8  `json:"outBoundTarget" enums:"0,1,2"` // 调出目标 0:平台 1:门店 2:电柜
    InboundID      uint64 `json:"inboundId"`                    // 调入至 0:平台
    InboundTarget  uint8  `json:"inboundTarget" enums:"0,1,2"`  // 调入目标 0:平台 1:门店 2:电柜

    Force bool `swaggerignore:"true"` // 是否强制 (忽略电柜初始化)
}

type StockListReq struct {
    PaginationReq

    Name    *string `json:"name" query:"name"`       // 门店名称
    CityID  *uint64 `json:"cityId" query:"cityId"`   // 城市ID
    Start   *string `json:"start" query:"start"`     // 开始时间
    End     *string `json:"end" query:"end"`         // 结束时间
    StoreID *uint64 `json:"storeId" query:"storeId"` // 门店ID
}

type StockMaterial struct {
    Name      string `json:"name"`      // 物资名称
    Outbound  int    `json:"outbound"`  // 出库数量
    Inbound   int    `json:"inbound"`   // 入库数量
    Surplus   int    `json:"surplus"`   // 剩余
    Exception int    `json:"exception"` // 异常数量
}

type StockStoreMaterial struct {
    Model string `json:"model"` // 电池型号
    Name  string `json:"name"`  // 物资名称
    Num   int    `json:"num"`   // 物资数量
}

type StockListRes struct {
    Store        Store            `json:"store"`        // 门店
    City         City             `json:"city"`         // 城市
    BatteryTotal int              `json:"batteryTotal"` // 电池总数
    Batteries    []*StockMaterial `json:"batteries"`    // 电池详情
    Materials    []*StockMaterial `json:"materials"`    // 非电池物资详情
}

type StockOverview struct {
    Total     int `json:"total"`     // 电池总数
    Surplus   int `json:"surplus"`   // 库存电池
    Outbound  int `json:"outbound"`  // 电池出库数
    Inbound   int `json:"inbound"`   // 电池库存数
    Exception int `json:"exception"` // 电池异常数
}

type StockWithRiderReq struct {
    RiderID   uint64 `json:"riderId"`   // 骑手ID
    Model     string `json:"model"`     // 电池型号
    StockType uint8  `json:"stockType"` // 出入库类型

    StoreID    *uint64 `json:"storeId"`    // 门店ID
    EmployeeID *uint64 `json:"employeeId"` // 店员ID
    ManagerID  *uint64 `json:"managerId"`  // 管理员ID
    CabinetID  *uint64 `json:"cabinetId"`  // 电柜ID
}

type StockEmployeeOverviewBattery struct {
    Model    string `json:"model"`    // 电池型号
    Surplus  int    `json:"surplus"`  // 库存电池
    Outbound int    `json:"outbound"` // 今日出库
    Inbound  int    `json:"inbound"`  // 今日入库
}

type StockEmployeeOverviewMaterial struct {
    Name    string `json:"name"`    // 物资名称
    Surplus int    `json:"surplus"` // 库存数量
}

type StockEmployeeOverview struct {
    Batteries []*StockEmployeeOverviewBattery  `json:"batteries"` // 电池
    Materials []*StockEmployeeOverviewMaterial `json:"materials"` // 非电池物资
}

type StockEmployeeListReq struct {
    PaginationReq

    Start    *string `json:"start" query:"start"`       // 筛选开始日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
    End      *string `json:"end" query:"end"`           // 筛选结束日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
    Outbound bool    `json:"outbound" query:"outbound"` // 是否筛选出库, false(默认):入库 true:出库
}

type StockEmployeeListResItem struct {
    ID    uint64  `json:"id"`
    Type  uint8   `json:"type"`            // 出入库类型 0:调拨 1:激活 2:寄存 3:结束寄存 4:退租
    Name  string  `json:"name,omitempty"`  // 骑手姓名, 平台调拨此字段不存在
    Phone string  `json:"phone,omitempty"` // 骑手电话, 平台调拨此字段不存在
    Num   int     `json:"num"`             // 数量, 正数为入库, 负数为出库(前端显示正数)
    Model *string `json:"model"`           // 电池型号
    Time  string  `json:"time"`            // 时间
}

type StockEmployeeListRes struct {
    Today *int `json:"today,omitempty"` // 今日出库/入库数量, 仅第一页请求时会返回

    *PaginationRes
}
