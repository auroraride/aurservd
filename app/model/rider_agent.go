// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-03
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
	AgentRiderStatusInactive     = "inactive"
	AgentRiderStatusUsing        = "using"
	AgentRiderStatusOverdue      = "overdue"
	AgentRiderStatusUnsubscribed = "unsubscribed"
	AgentRiderStatusWillOverdue  = "will_overdue"
)

type AgentRiderListReq struct {
	PaginationReq

	CityID    uint64 `json:"cityId" query:"cityId"`                                             // 城市
	Keyword   string `json:"keyword" query:"keyword"`                                           // 关键词
	StationID uint64 `json:"stationId" query:"stationId"`                                       // 站点ID
	Status    string `json:"status" query:"status" enums:"inactive,using,overdue,unsubscribed"` // inactive:未激活 using:计费中 overdue:已超期 unsubscribed:已退租
}

type AgentRider struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`                // 姓名
	Phone     string `json:"phone,omitempty"`     // 电话
	Status    string `json:"status"`              // 状态
	Remaining *int   `json:"remaining,omitempty"` // 剩余天数
	// City        *City  `json:"city,omitempty"`      // 城市
	// Date        string `json:"date"`            // 创建日期
	Station string `json:"station"` // 站点
	// EndAt       string `json:"endAt"`           // 退租日期
	StopAt  string `json:"stopAt"`  // 到期日期
	StartAt string `json:"startAt"` // 开始日期
	// SubscribeID uint64 `json:"subscribeId"`     // 订阅ID
	Model string `json:"model,omitempty"` // 电池型号
	// Used  int    `json:"used"`            // 使用天数
}

type AgentRiderLogReq struct {
	ID           uint64 `json:"id" query:"id" validate:"required" trans:"骑手ID"`
	EnterpriseID uint64 `json:"enterpriseId" query:"enterpriseId"`
}

type AgentRiderLog struct {
	Days int    `json:"days"` // 延长天数
	Time string `json:"time"` // 操作时间
	Name string `json:"name"` // 操作人
}

type AgentRiderDetail struct {
	AgentRider
	Logs []AgentRiderLog `json:"logs,omitempty"`
}
