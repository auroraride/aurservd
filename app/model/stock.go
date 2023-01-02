// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "errors"
    "strings"
)

const (
    StockNameEbike = "电车"
)

const (
    StockTypeTransfer         uint8 = iota // 调拨 (出库入库)
    StockTypeRiderActive                   // 骑手激活电池 (出库)
    StockTypeRiderPause                    // 骑手寄存电池 (入库)
    StockTypeRiderContinue                 // 骑手结束寄存电池 (出库)
    StockTypeRiderUnSubscribe              // 骑手归还电池 (入库)
)

var (
    StockTypesText = map[uint8]string{
        StockTypeTransfer:         "调拨",
        StockTypeRiderActive:      "激活",
        StockTypeRiderPause:       "寄存",
        StockTypeRiderContinue:    "结束寄存",
        StockTypeRiderUnSubscribe: "退租",
    }
)

const (
    StockTargetPlaform uint8 = iota
    StockTargetStore         // 调拨对象 - 门店
    StockTargetCabinet       // 调拨对象 - 电柜
)

// StockNumberOfRiderBusiness 出入库电池数量
func StockNumberOfRiderBusiness(typ uint8) (num int) {
    switch typ {
    case StockTypeRiderActive, StockTypeRiderContinue:
        num = -1
        break
    case StockTypeRiderPause, StockTypeRiderUnSubscribe:
        num = 1
        break
    }
    return
}

type StockTransferLoopper struct {
    Message string // 错误提示消息

    EbikeID   *uint64 // 电车ID
    BrandID   *uint64 // 电车型号ID
    BrandName *string // 电车型号
    EbikeSN   *string // 车架号
}

type StockTransferReq struct {
    Model string `json:"model,omitempty"` // 电池型号 (和 `物资名称` / `车架` 不能同时存在, 也不能同时为空)
    Name  string `json:"name,omitempty"`  // 物资名称 (和 `电池型号` / `车架` 不能同时存在, 也不能同时为空)
    Num   int    `json:"num"`             // 调拨数量

    Ebikes []string `json:"ebikes,omitempty"` // 车架编号 (和 `物资名称` / `车架` 不能同时存在, 也不能同时为空)

    Remark string `json:"remark" validate:"required" trans:"备注"`

    OutboundID     uint64 `json:"outboundId"`                   // 调出自 0:平台
    OutboundTarget uint8  `json:"outboundTarget" enums:"0,1,2"` // 调出目标 0:平台 1:门店 2:电柜
    InboundID      uint64 `json:"inboundId"`                    // 调入至 0:平台
    InboundTarget  uint8  `json:"inboundTarget" enums:"0,1,2"`  // 调入目标 0:平台 1:门店 2:电柜

    Force bool `swaggerignore:"true"` // 是否强制 (忽略电柜初始化)
}

func (req *StockTransferReq) IsToStore() bool {
    return req.InboundTarget == StockTargetStore
}

func (req *StockTransferReq) IsFromStore() bool {
    return req.OutboundTarget == StockTargetStore
}

func (req *StockTransferReq) IsToCabinet() bool {
    return req.InboundTarget == StockTargetCabinet
}

func (req *StockTransferReq) IsFromCabinet() bool {
    return req.OutboundTarget == StockTargetCabinet
}

func (req *StockTransferReq) IsToPlaform() bool {
    return req.InboundTarget == StockTargetPlaform
}

func (req *StockTransferReq) IsFromPlaform() bool {
    return req.OutboundTarget == StockTargetPlaform
}

func (req *StockTransferReq) Batchable() bool {
    return len(req.Ebikes) == 0
}

// ParticipateCabinet 是否有电柜参与
func (req *StockTransferReq) ParticipateCabinet() bool {
    return req.IsFromCabinet() || req.IsToCabinet()
}

// RealNumber 获取实际物资数量
func (req *StockTransferReq) RealNumber() int {
    if len(req.Ebikes) > 0 {
        return 1
    }
    return req.Num
}

func (req *StockTransferReq) RealName() string {
    if req.Model != "" {
        return req.Model
    }
    return req.Name
}

// Validate 校验
func (req *StockTransferReq) Validate() error {
    var mm, meb, mo int
    // 非智能电池
    if req.Model != "" {
        req.Model = strings.ToUpper(req.Model)
        mm = 1
    }
    // 其他物资
    if req.Name != "" {
        mo = 1
    }
    // 电车
    if len(req.Ebikes) > 0 {
        meb = 1
    }

    // 运算判定物资是否正确(互斥不为空)
    v := mm + meb + mo
    if v != 1 {
        return errors.New("物资选项错误")
    }

    // 校验电柜是否可参与
    if req.Model == "" && req.ParticipateCabinet() {
        return errors.New("电柜无法参与非电池物资调拨")
    }

    // 校验物资数量
    if req.RealNumber() < 1 {
        return errors.New("调拨物资数量计数错误")
    }

    if req.InboundID == 0 && req.OutboundID == 0 {
        return errors.New("平台之间无法调拨物资")
    }

    if req.IsToCabinet() && req.IsFromCabinet() {
        return errors.New("电柜之间无法调拨")
    }

    if ((req.IsToStore() || req.IsToCabinet()) && req.InboundID == 0) || (req.IsToPlaform() && req.InboundID != 0) {
        return errors.New("调入参数错误")
    }

    if ((req.IsFromStore() || req.IsFromCabinet()) && req.OutboundID == 0) || (req.IsFromPlaform() && req.OutboundID != 0) {
        return errors.New("调出参数错误")
    }

    return nil
}

type StockOverviewReq struct {
    Goal      StoreCabiletGoal `json:"goal" query:"goal" enums:"0,1,2"` // 查询目标, 0:不筛选 1:门店(默认) 2:电柜
    CabinetID uint64           `json:"cabinetId" query:"cabinetId"`     // 电柜ID, 仅goal为2的时候生效
    StoreID   uint64           `json:"storeId" query:"storeId"`         // 门店ID, 仅goal为1的时候生效
    CityID    uint64           `json:"cityId" query:"cityId"`           // 城市ID
    Start     string           `json:"start" query:"start"`             // 开始时间
    End       string           `json:"end" query:"end"`                 // 结束时间
}

// StockBatteryOverviewRes 电池概览
type StockBatteryOverviewRes struct {
    Model        string `json:"model"`         // 电池型号
    Num          int    `json:"num"`           // 总数
    STransfer    int    `json:"s_transfer"`    // 门店-平台调拨数量
    SOutboundnum int    `json:"s_outboundnum"` // 门店-出库数量
    SInboundnum  int    `json:"s_inboundnum"`  // 门店-入库数量
    SActive      int    `json:"s_active"`      // 门店-激活
    SPause       int    `json:"s_pause"`       // 门店-寄存
    SContinue    int    `json:"s_continue"`    // 门店-结束寄存
    SUnsubscribe int    `json:"s_unsubscribe"` // 门店-退租
    CTransfer    int    `json:"c_transfer"`    // 电柜-平台调拨数量
    COutboundnum int    `json:"c_outboundnum"` // 电柜-出库数量
    CInboundnum  int    `json:"c_inboundnum"`  // 电柜-入库数量
    CActive      int    `json:"c_active"`      // 电柜-激活
    CPause       int    `json:"c_pause"`       // 电柜-寄存
    CContinue    int    `json:"c_continue"`    // 电柜-结束寄存
    CUnsubscribe int    `json:"c_unsubscribe"` // 电柜-退租
}

type StockListReq struct {
    PaginationReq

    Name    *string `json:"name" query:"name"`       // 门店名称
    CityID  *uint64 `json:"cityId" query:"cityId"`   // 城市ID
    Start   *string `json:"start" query:"start"`     // 开始时间
    End     *string `json:"end" query:"end"`         // 结束时间
    StoreID *uint64 `json:"storeId" query:"storeId"` // 门店ID

    EbikeBrandID uint64 `json:"ebikeBrandId" query:"ebikeBrandId"` // 电车型号ID
    Model        string `json:"model" query:"model"`               // 电池型号
    Keyword      string `json:"keyword" query:"keyword"`           // 其他物资名称
}

type StockMaterial struct {
    Name      string `json:"name"`                // 物资名称
    Outbound  int    `json:"outbound"`            // 出库数量
    Inbound   int    `json:"inbound"`             // 入库数量
    Surplus   int    `json:"surplus"`             // 剩余
    Exception int    `json:"exception,omitempty"` // 异常数量(电柜无)
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
    EbikeTotal   int              `json:"ebikeTotal"`   // 电车总数
    Ebikes       []*StockMaterial `json:"ebikes"`       // 电车
}

// StockBusinessReq 业务库存调整请求
type StockBusinessReq struct {
    RiderID   uint64 `json:"riderId"`   // 骑手ID
    Model     string `json:"model"`     // 电池型号
    StockType uint8  `json:"stockType"` // 出入库类型
    CityID    uint64 `json:"cityId"`    // 城市

    StoreID     *uint64 `json:"storeId"`     // 门店ID
    EmployeeID  *uint64 `json:"employeeId"`  // 店员ID
    CabinetID   *uint64 `json:"cabinetId"`   // 电柜ID
    SubscribeID *uint64 `json:"subscribeId"` // 订阅ID

    Ebike   *EbikeBusinessInfo `json:"ebike"`   // 电车信息
    Battery *Battery           `json:"battery"` // 电池信息
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

type StockCabinetListReq struct {
    PaginationReq

    Serial    string `json:"serial" query:"serial"`       // 电柜编号
    CityID    uint64 `json:"cityId" query:"cityId"`       // 城市ID
    CabinetID uint64 `json:"cabinetId" query:"cabinetId"` // 电柜ID
    Start     string `json:"start" query:"start"`         // 开始时间
    End       string `json:"end" query:"end"`             // 结束时间
}

type StockCabinetListRes struct {
    ID        uint64           `json:"id"`        // 电柜ID
    Serial    string           `json:"serial"`    // 电柜编号
    City      City             `json:"city"`      // 城市
    Name      string           `json:"name"`      // 电柜名称
    Batteries []*StockMaterial `json:"batteries"` // 电池详情
}

type StockDetailFilter struct {
    Goal      StoreCabiletGoal `json:"goal" query:"goal" enums:"0,1,2,3"`   // 查询目标 0:不筛选 1:门店(默认) 2:电柜 3:平台
    Materials string           `json:"materials" query:"materials"`         // 查询物资类别, 默认为电池, 逗号分隔 battery:电池 ebike:电车 others:其他物资
    Serial    string           `json:"serial" query:"serial"`               // 电柜编号
    CityID    uint64           `json:"cityId" query:"cityId"`               // 城市ID
    CabinetID uint64           `json:"cabinetId" query:"cabinetId"`         // 电柜ID
    StoreID   uint64           `json:"storeId" query:"storeId"`             // 门店ID
    Start     string           `json:"start" query:"start"`                 // 开始时间
    End       string           `json:"end" query:"end"`                     // 结束时间
    Positive  bool             `json:"positive" query:"positive"`           // 是否正序(默认倒序)
    Type      uint8            `json:"type" query:"type" enums:"0,1,2,3,4"` // 调拨类型, 0:调拨 1:激活 2:寄存 3:结束寄存 4:退租
}

type StockDetailReq struct {
    PaginationReq
    StockDetailFilter
}

type StockDetailExportReq struct {
    StockDetailFilter
    Remark string `json:"remark" validate:"required" trans:"备注"`
}

type StockDetailRes struct {
    ID       uint64 `json:"id"`              // 调拨ID
    Sn       string `json:"sn"`              // 调拨编号
    City     string `json:"city"`            // 城市
    Inbound  string `json:"inbound"`         // 调入
    Outbound string `json:"outbound"`        // 调出
    Name     string `json:"name"`            // 物资
    Num      int    `json:"num"`             // 数量
    Type     string `json:"type"`            // 类型
    Operator string `json:"operator"`        // 操作人
    Time     string `json:"time"`            // 时间
    Remark   string `json:"remark"`          // 备注信息
    Rider    string `json:"rider,omitempty"` // 骑手, 仅业务发生有此字段
}

type StockInventoryMapData map[uint64]map[string]map[string]StockInventory

type StockInventoryReq struct {
    IDs      []uint64         `json:"ids"`
    Goal     StoreCabiletGoal `json:"goal"`
    Material string           `json:"material"`
    Name     string           `json:"name"`
}

type StockInventory struct {
    CabinetID uint64 `json:"cabinet_id,omitempty"`
    StoreID   uint64 `json:"store_id,omitempty"`
    Name      string `json:"name"`
    Num       int    `json:"num"`
    Material  string `json:"material"`
}
