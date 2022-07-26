// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    OrderTypeNewly      uint = iota + 1 // 新签, 需要计算业绩, 当退订时间超出设定时间间隔后视为新签
    OrderTypeRenewal                    // 续签, 无需计算业绩
    OrderTypeAgain                      // 重签, 无需计算业绩
    OrderTypeTransform                  // 更改电池, 相当于续签 无需计算业绩 TODO 更改电池逻辑
    OrderTypeAssistance                 // 救援
    OrderTypeFee                        // 滞纳金
    OrderTypeDeposit                    // 押金
)

var (
    // OrderSubscribeTypes 骑手骑士卡订单类型
    OrderSubscribeTypes = []uint{OrderTypeNewly, OrderTypeAgain, OrderTypeRenewal, OrderTypeTransform}
)

const (
    OrderPaywayManual uint8 = iota // 后台手动调整
    OrderPaywayAlipay              // 支付宝支付
    OrderPaywayWechat              // 微信支付
)

const (
    OrderStatusPending       uint8 = iota // 未支付
    OrderStatusPaid                       // 已支付
    OrderStatusRefundPending              // 申请退款, 退款后业绩订单需要删除
    OrderStatusRefundSuccess              // 已退款
    OrderStatusRefundRefused              // 退款被拒绝
)

// OrderCreateReq 订单创建请求
type OrderCreateReq struct {
    PlanID    uint64 `json:"planId" validate:"required" trans:"套餐ID"`
    Payway    uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2"`              // 1支付宝 2微信
    OrderType uint   `json:"orderType" validate:"required" trans:"订单类型" enums:"1,2,3,4,5,6,7"` // 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金

    CityID uint64 `json:"cityId"` // 城市ID, 新签必填
    Model  string `json:"model"`  // 用户所选电池型号, 新签必填
}

// OrderCreateRes 订单创建返回
type OrderCreateRes struct {
    Prepay     string `json:"prepay"`     // 预支付字符串
    OutTradeNo string `json:"outTradeNo"` // 交易编码
}

// OrderListReq 订单列表请求
type OrderListReq struct {
    PaginationReq

    Type       *uint   `json:"type" query:"type"`                   // 订单类型 1:新签 2:续签 3:重签 4:更改电池 5:救援 6:滞纳金 7:押金
    CityID     *uint64 `json:"cityId" query:"cityId"`               // 城市ID
    Keyword    *string `json:"keyword" query:"keyword"`             // 骑手姓名
    Start      *string `json:"start" query:"start"`                 // 时间起始, 格式为: 2022-01-01
    End        *string `json:"end" query:"end"`                     // 时间结束, 格式为: 2022-01-01
    StoreName  *string `json:"storeName" query:"storeName"`         // 门店名字
    Model      *string `json:"model" query:"model"`                 // 电池型号
    Days       *int    `json:"days" query:"days"`                   // 骑士卡时长(搜索大于等于)
    Refund     uint8   `json:"refund" query:"refund"`               // 退款查询 0:查询全部 1:查询未申请退款 2:查询已申请退款(包含退款中/已退款/已拒绝)
    EmployeeID uint64  `json:"employeeId" query:"employeeId"`       // 店员ID筛选
    Payway     *uint8  `json:"payway" query:"payway" enums:"0,1,2"` // 支付方式 0:手动 1:支付宝 2:微信, 不携带此参数为获取全部
}

type OrderEmployeeListReq struct {
    PaginationReq

    Aimed   uint8   `json:"aimed" query:"aimed"`     // 筛选对象 0:全部 1:个签 2:团签
    Start   *string `json:"start" query:"start"`     // 筛选开始日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
    End     *string `json:"end" query:"end"`         // 筛选结束日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
    Keyword *string `json:"keyword" query:"keyword"` // 筛选骑手姓名或电话
    Type    *string `json:"type" query:"type"`       // 筛选订单类别
}

type OrderStatusReq struct {
    OutTradeNo string `json:"outTradeNo" query:"outTradeNo"` // 订单编号
}

type OrderStatusRes struct {
    OutTradeNo string `json:"outTradeNo"` // 订单编号
    Paid       bool   `json:"paid"`       // 是否支付
}
