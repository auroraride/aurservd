// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
	OrderTypeNewly           uint = iota + 1 // 新签, 需要计算业绩, 当退订时间超出设定时间间隔后视为新签
	OrderTypeRenewal                         // 续签, 无需计算业绩
	OrderTypeAgain                           // 重签, 无需计算业绩
	OrderTypeTransform                       // 更改电池, 相当于续签 无需计算业绩 TODO 更改电池逻辑
	OrderTypeAssistance                      // 救援
	OrderTypeFee                             // 滞纳金
	OrderTypeDeposit                         // 押金
	OrderTypeAgentPrepayment                 // 代理充值
)

const (
	OrderPaywayManual           uint8 = iota // 后台手动调整
	OrderPaywayAlipay                        // 支付宝支付
	OrderPaywayWechat                        // 微信支付
	OrderPaywayAlipayAuthFreeze              // 支付宝预授权支付
	OrderPaywayWechatDeposit                 // 微信支付分支付
	OrderPaywayAlipayDeposit                 // 支付宝芝麻信用分支付
)

const (
	OrderStatusPending       uint8 = iota // 未支付
	OrderStatusPaid                       // 已支付
	OrderStatusRefundPending              // 申请退款, 退款后业绩订单需要删除
	OrderStatusRefundSuccess              // 已退款
	OrderStatusRefundRefused              // 退款被拒绝
)

var (
	// OrderSubscribeTypes 骑手骑士卡订单类型
	OrderSubscribeTypes = []uint{OrderTypeNewly, OrderTypeAgain, OrderTypeRenewal, OrderTypeTransform}

	OrderTypes = map[uint]string{
		OrderTypeNewly:      "新签",
		OrderTypeRenewal:    "续签",
		OrderTypeAgain:      "重签",
		OrderTypeTransform:  "更改电池",
		OrderTypeAssistance: "救援",
		OrderTypeFee:        "滞纳金",
		OrderTypeDeposit:    "押金",
	}

	OrderStatuses = map[uint8]string{
		OrderStatusPending:       "未支付",
		OrderStatusPaid:          "已支付",
		OrderStatusRefundPending: "申请退款",
		OrderStatusRefundSuccess: "已退款",
		OrderStatusRefundRefused: "退款被拒绝",
	}

	OrderPayways = map[uint8]string{
		OrderPaywayManual:           "后台手动调整",
		OrderPaywayAlipay:           "支付宝支付",
		OrderPaywayWechat:           "微信支付",
		OrderPaywayAlipayAuthFreeze: "支付宝预授权支付",
	}
)

// OrderCreateReq 订单创建请求
type OrderCreateReq struct {
	PlanID    uint64 `json:"planId" validate:"required" trans:"套餐ID"`
	Payway    uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2,3"`            // 1支付宝 2微信 3支付宝预授权
	OrderType uint   `json:"orderType" validate:"required" trans:"订单类型" enums:"1,2,3,4,5,6,7"` // 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金

	CityID uint64 `json:"cityId"` // 城市ID, 新签必填

	Point   bool     `json:"point"`   // 是否使用积分
	Coupons []uint64 `json:"coupons"` // 优惠券
}

// OrderCreateRes 订单创建返回
type OrderCreateRes struct {
	Prepay     string `json:"prepay"`     // 预支付字符串
	OutTradeNo string `json:"outTradeNo"` // 交易编码
}

type OrderListFilter struct {
	RiderID    *uint64 `json:"riderId,omitempty" query:"riderId"`             // 骑手ID
	Type       *uint   `json:"type,omitempty" query:"type"`                   // 订单类型 1:新签 2:续签 3:重签 4:更改电池 5:救援 6:滞纳金 7:押金
	CityID     *uint64 `json:"cityId,omitempty" query:"cityId"`               // 城市ID
	Keyword    *string `json:"keyword,omitempty" query:"keyword"`             // 骑手姓名
	Start      *string `json:"start,omitempty" query:"start"`                 // 时间起始, 格式为: 2022-01-01
	End        *string `json:"end,omitempty" query:"end"`                     // 时间结束, 格式为: 2022-01-01
	StoreName  *string `json:"storeName,omitempty" query:"storeName"`         // 门店名字
	Model      *string `json:"model,omitempty" query:"model"`                 // 电池型号
	Days       *int    `json:"days,omitempty" query:"days"`                   // 骑士卡时长(搜索大于等于)
	Refund     *uint8  `json:"refund,omitempty" query:"refund"`               // 退款查询 0:查询全部 1:查询未申请退款 2:查询已申请退款(包含退款中/已退款/已拒绝)
	EmployeeID *uint64 `json:"employeeId,omitempty" query:"employeeId"`       // 店员ID筛选
	Payway     *uint8  `json:"payway,omitempty" query:"payway" enums:"0,1,2"` // 支付方式 0:手动 1:支付宝 2:微信, 不携带此参数为获取全部
	TradeNo    *string `json:"tradeNo,omitempty" query:"tradeNo"`             // 平台单号
}

// OrderListReq 订单列表请求
type OrderListReq struct {
	PaginationReq
	OrderListFilter
}

type OrderListExport struct {
	OrderListFilter
	Remark string `json:"remark" validate:"required" trans:"备注"`
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

type OrderAgent struct {
	ID             uint64 `json:"id"`             // 代理ID
	Name           string `json:"name"`           // 代理姓名
	Phone          string `json:"phone"`          // 代理电话
	EnterpriseID   uint64 `json:"enterpriseId"`   // 团签ID
	EnterpriseName string `json:"enterpriseName"` // 团签名称
}

// Order 骑手订单
type Order struct {
	ID            uint64        `json:"id"`                 // 订单ID
	Type          uint          `json:"type"`               // 订单类型 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金
	Status        uint8         `json:"status"`             // 订单状态 0未支付 1已支付 2申请退款 3已退款 4退款被拒绝
	Payway        uint8         `json:"payway"`             // 支付方式 1支付宝 2微信
	PayAt         string        `json:"payAt"`              // 支付时间
	Amount        float64       `json:"amount"`             // 支付金额
	OutTradeNo    string        `json:"outTradeNo"`         // 订单编号
	TradeNo       string        `json:"tradeNo"`            // 订单编号 (支付平台)
	City          City          `json:"city"`               // 城市
	Rider         Rider         `json:"rider"`              // 骑手
	Plan          *Plan         `json:"plan,omitempty"`     // 骑士卡, 非骑士卡订阅订单无此字段 (可为空)
	Model         string        `json:"model,omitempty"`    // 电池型号 (可为空)
	Store         *Store        `json:"store,omitempty"`    // 门店 (可为空)
	Employee      *Employee     `json:"employee,omitempty"` // 店员 (可为空)
	Refund        *Refund       `json:"refund,omitempty"`   // 退款详情 (可为空)
	PointAmount   float64       `json:"pointAmount"`        // 积分抵扣金额
	DiscountNewly float64       `json:"discountNewly"`      // 新签优惠
	CouponAmount  float64       `json:"couponAmount"`       // 优惠券抵扣金额
	Coupons       []CouponRider `json:"coupons,omitempty"`  // 使用的优惠券
	Ebike         *Ebike        `json:"ebike"`              // 车辆详情
	Agent         *OrderAgent   `json:"agent"`              // 代理信息
}
