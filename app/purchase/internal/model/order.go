package model

import (
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

var (
	PurchaseOrderFormula = "计算方法：违约时间超过7天，按照本期应支付金额的20%计算。"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待支付
	OrderStatusStaging   OrderStatus = "staging"   // 分期执行中
	OrderStatusEnded     OrderStatus = "ended"     // 已完成
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
	OrderStatusRefunded  OrderStatus = "refunded"  // 已退款
)

func (o OrderStatus) Value() string {
	return string(o)
}

func (o OrderStatus) String() string {
	switch o {
	case OrderStatusPending:
		return "待支付"
	case OrderStatusStaging:
		return "分期执行中"
	case OrderStatusEnded:
		return "已完成"
	case OrderStatusCancelled:
		return "已取消"
	case OrderStatusRefunded:
		return "已退款"
	default:
		return "未知"
	}
}

// RepayStatus 还款状态
type RepayStatus uint8

const (
	RepayStatusNormal  RepayStatus = iota + 1 // 正常
	RepayStatusOverdue                        // 逾期
)

func (r RepayStatus) Value() uint8 {
	return uint8(r)
}

func (r RepayStatus) String() string {
	switch r {
	case RepayStatusNormal:
		return "正常"
	case RepayStatusOverdue:
		return "逾期"
	default:
		return "未知"
	}
}

// PurchaseOrderCreateReq 创建订单请求
type PurchaseOrderCreateReq struct {
	GoodsID   uint64 `json:"goodsId" validate:"required"`   // 商品id
	PlanIndex *int   `json:"planIndex" validate:"required"` // 付款计划索引
}

// PurchaseOrderListReq 订单列表请求
type PurchaseOrderListReq struct {
	model.PaginationReq
	PurchaseOrderListFilter
}

type PurchaseOrderListFilter struct {
	Keyword     *string      `json:"keyword" query:"keyword"`         // 关键字
	ID          *uint64      `json:"id" query:"id"`                   // 订单编号
	Sn          *string      `json:"sn" query:"sn"`                   // 车架号
	Status      *OrderStatus `json:"status" query:"status"`           // 订单状态
	RepayStatus *RepayStatus `json:"repayStatus" query:"repayStatus"` // 还款状态
	StoreID     *uint64      `json:"storeId" query:"storeId"`         // 门店ID
	Start       *string      `json:"start" query:"start"`             // 开始时间
	End         *string      `json:"end" query:"end"`                 // 结束时间
	RiderID     *uint64      `json:"riderId" query:"riderId"`         // 骑手ID
}

type PurchaseOrderExportReq struct {
	PurchaseOrderListFilter
	Remark string `json:"remark" validate:"required" trans:"备注"`
}

// PurchaseOrderListRes 订单列表返回
type PurchaseOrderListRes struct {
	ID               uint64            `json:"id"`               // 订单编号
	Status           OrderStatus       `json:"status"`           // 订单状态 pending: 待支付, staging: 分期执行中, ended: 已完成, cancelled: 已取消, refunded: 已退款
	Goods            *definition.Goods `json:"goods,omitempty"`  // 商品信息
	Amount           float64           `json:"amount"`           // 订单金额
	PaidAmount       float64           `json:"paidAmount"`       // 已支付金额
	InstallmentTotal int               `json:"installmentTotal"` // 分期总数
	InstallmentStage int               `json:"installmentStage"` // 当前分期阶段
	RepayStatus      RepayStatus       `json:"repayStatus"`      // 还款状态 // 1-正常 2-逾期
	RiderName        string            `json:"riderName"`        // 骑手名称
	RiderPhone       string            `json:"riderPhone"`       // 骑手电话
	StoreID          uint64            `json:"storeId"`          // 门店ID
	StoreName        string            `json:"storeName"`        // 提车门店
	Sn               string            `json:"sn"`               // 车架号
	Color            string            `json:"color"`            // 车辆颜色
	ActiveName       *string           `json:"activeName"`       // 激活人名称
	ActivePhone      *string           `json:"activePhone"`      // 激活人电话
	CreatedAt        string            `json:"createdAt"`        // 创建时间
	StartDate        *string           `json:"startDate"`        // 激活时间
	Remark           string            `json:"remark"`           // 备注
	ContractUrl      *string           `json:"contractUrl"`      // 合同url
	Signed           bool              `json:"signed"`           // 是否签约 true:已签约 false:未签约
	PlanIndex        *int              `json:"planIndex"`        // 付款计划索引
	DocID            *string           `json:"docId"`            // 合同ID
	Formula          string            `json:"formula"`          // 违约说明
	InstallmentPlan  []float64         `json:"installmentPlan"`  // 分期方案
}

// PurchaseOrderDetail 订单详情
type PurchaseOrderDetail struct {
	PurchaseOrderListRes
	Payments []*PaymentDetail       `json:"payments"` // 分期订单数据（还款计划）
	Follows  []*PurchaseOrderFollow `json:"follows"`  // 订单跟进数据
}

// PurchaseOrderFollow 订单跟进信息
type PurchaseOrderFollow struct {
	ID        uint64          `json:"id"`        // 跟进ID
	Content   string          `json:"content"`   // 跟进内容
	Pics      []string        `json:"pics"`      // 跟进图片
	Modifier  *model.Modifier `json:"modifier"`  // 跟进人
	CreatedAt string          `json:"createdAt"` // 跟进时间
}

// PurchaseOrderFollowReq 购车订单跟进请求
type PurchaseOrderFollowReq struct {
	ID      uint64   `json:"id" validate:"required"`      // 订单编号
	Content string   `json:"content" validate:"required"` // 跟进内容
	Pics    []string `json:"pics" validate:"max=10"`      // 跟进图片
}

// PurchaseOrderActiveReq 购车订单激活请求
type PurchaseOrderActiveReq struct {
	ID        uint64 `json:"id" validate:"required"`        // 订单编号
	GoodsID   uint64 `json:"goodsId" validate:"required"`   // 商品ID
	PlanIndex *int   `json:"planIndex" validate:"required"` // 付款计划索引（分期期数）
	StoreID   uint64 `json:"storeId" validate:"required"`   // 门店ID
	Sn        string `json:"sn" validate:"required"`        // 车架号
	Color     string `json:"color" validate:"required"`     // 车辆颜色
	Remark    string `json:"remark"`                        // 备注
}
