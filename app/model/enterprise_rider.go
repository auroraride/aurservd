// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// EnterpriseRiderCreateReq 企业骑手创建请求
type EnterpriseRiderCreateReq struct {
	EnterpriseID uint64 `json:"enterpriseId" validate:"required" trans:"企业ID"`
	StationID    uint64 `json:"stationId" validate:"required" trans:"站点ID"`
	Name         string `json:"name" validate:"required" trans:"姓名"`
	Phone        string `json:"phone" validate:"required,phone" trans:"电话号"`

	Days    int    `json:"days,omitempty"`    // 天数, 仅代理商需要
	PriceID uint64 `json:"priceId,omitempty"` // 价格ID, 仅代理商需要
}

type EnterpriseRider struct {
	ID              uint64            `json:"id"`
	Name            string            `json:"name"`                // 姓名
	Phone           string            `json:"phone"`               // 电话
	Days            int               `json:"days"`                // 总天数
	Unsettled       int               `json:"unsettled"`           // 未结算天数
	Model           string            `json:"model,omitempty"`     // 可用电池型号, 骑手未开通订阅则此字段不存在
	CreatedAt       string            `json:"createdAt"`           // 添加时间
	Station         EnterpriseStation `json:"station"`             // 站点
	DeletedAt       string            `json:"deletedAt,omitempty"` // 删除时间, 已被删除的用户才会有此字段
	SubscribeStatus uint8             `json:"subscribeStatus"`     // 订阅状态
}

type EnterpriseRiderListReq struct {
	PaginationReq

	EnterpriseID    uint64  `json:"enterpriseId" validate:"required" query:"enterpriseId" trans:"企业ID"`
	Keyword         *string `json:"keyword" query:"keyword"`                                 // 搜索关键词
	Start           *string `json:"start" query:"start"`                                     // 使用开始时间
	End             *string `json:"end" query:"end"`                                         // 使用结束时间
	Deleted         uint8   `json:"deleted" query:"deleted" enums:"0,1,2"`                   // 筛选删除 0:全部 1:已删除 2:未删除
	SubscribeStatus uint8   `json:"subscribeStatus" query:"subscribeStatus" enums:"0,1,2,3"` // 筛选订阅状态 0:全部 1:计费中 2:已退租 3:未激活
}

// EnterpriseRiderInviteReq 邀请骑手
type EnterpriseRiderInviteReq struct {
	// 站点id
	StationID uint64 `json:"stationID" param:"stationID"  validate:"required"`
}
