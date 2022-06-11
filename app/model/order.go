// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    OrderTypeNewly     uint = iota + 1 // 新签, 需要计算业绩
    OrderTypeRenewal                   // 续签, 无需计算业绩
    OrderTypeAgain                     // 重签, 相当于新签, 需要判定是否计算业绩
    OrderTypeTransform                 // 更改电池, 相当于续签 无需计算业绩 TODO 更改电池逻辑
    OrderTypeRescue                    // 救援
    OrderTypeFee                       // 滞纳金
    OrderTypeDeposit                   // 押金
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
    CityID    uint64 `json:"cityId" validate:"required" trans:"城市ID"`
    PlanID    uint64 `json:"planId" validate:"required" trans:"套餐ID"`
    Payway    uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2"`              // 1支付宝 2微信
    OrderType uint   `json:"orderType" validate:"required" trans:"订单类型" enums:"1,2,3,4,5,6,7"` // 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金
}

// OrderCreateRes 订单创建返回
type OrderCreateRes struct {
    Prepay     string `json:"prepay"`     // 预支付字符串
    OutTradeNo string `json:"outTradeNo"` // 交易编码
}

// OrderNotActived 骑手未激活订单信息
type OrderNotActived struct {
    ID      uint64         `json:"id"`                 // 订单编号
    Amount  float64        `json:"amount"`             // 骑士卡金额
    Deposit float64        `json:"deposit"`            // 押金, 若押金为0则押金一行不显示
    Total   float64        `json:"total"`              // 总金额, 总金额为 amount + deposit
    Payway  uint8          `json:"payway" enums:"1,2"` // 支付方式 1支付宝 2微信
    Plan    PlanItem       `json:"plan"`               // 骑行卡详情
    City    City           `json:"city"`               // 所属城市
    Models  []BatteryModel `json:"models"`             // 可用电池型号, 显示为`72V30AH`即Voltage(V)+Capacity(AH), 逗号分隔
    Time    string         `json:"time"`               // 支付时间
}

// OrderListReq 订单列表请求
type OrderListReq struct {
    PaginationReq

    RiderID      *uint64  `json:"riderId" query:"riderId"`           // 骑手ID
    Type         *uint    `json:"type" query:"type"`                 // 订单类型 1:新签 2:续签 3:重签 4:更改电池 5:救援 6:滞纳金 7:押金
    CityID       *uint64  `json:"cityId" query:"cityId"`             // 城市ID
    RiderName    *string  `json:"riderName" query:"riderName"`       // 骑手姓名
    RiderPhone   *string  `json:"riderPhone" query:"riderPhone"`     // 骑手电话
    Start        *string  `json:"start" query:"start"`               // 时间起始, 格式为: 2022-01-01
    End          *string  `json:"end" query:"end"`                   // 时间结束, 格式为: 2022-01-01
    EmployeeName *string  `json:"employeeName" query:"employeeName"` // 店员名字
    StoreName    *string  `json:"storeName" query:"storeName"`       // 门店名字
    Voltage      *float64 `json:"voltage" query:"voltage"`           // 电压
    Days         *int     `json:"days" query:"days"`                 // 骑士卡时长(搜索大于等于)
    Refund       uint8    `json:"refund" query:"refund"`             // 退款查询 0:查询全部 1:查询未申请退款 2:查询已申请退款(包含退款中/已退款/已拒绝)
}
