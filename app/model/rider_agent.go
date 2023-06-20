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
	AgentRiderStatusWillOverdue  = "willOverdue"
)

const WillOverdueNum = 3

type AgentRiderListReq struct {
	PaginationReq

	CityID    uint64 `json:"cityId" query:"cityId"`                                                         // 城市
	Keyword   string `json:"keyword" query:"keyword"`                                                       // 关键词
	StationID uint64 `json:"stationId" query:"stationId"`                                                   // 站点ID
	Status    string `json:"status" query:"status" enums:"inactive,using,overdue,unsubscribed,willOverdue"` // inactive:未激活 using:计费中 overdue:已超期 unsubscribed:已退租 willOverdue:即将超期
}

type AgentRider struct {
	ID          uint64 `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`        // 姓名
	Phone       string `json:"phone,omitempty"`       // 电话
	Status      string `json:"status,omitempty"`      // 状态
	Remaining   *int   `json:"remaining,omitempty"`   // 剩余天数
	City        *City  `json:"city,omitempty"`        // 城市
	Date        string `json:"date,omitempty"`        // 创建日期
	EndAt       string `json:"endAt,omitempty"`       // 退租日期
	StopAt      string `json:"stopAt,omitempty"`      // 到期日期
	StartAt     string `json:"startAt,omitempty"`     // 开始日期
	SubscribeID uint64 `json:"subscribeId,omitempty"` // 订阅ID
	Model       string `json:"model,omitempty"`       // 电池型号
	BatterySN   string `json:"batterySn,omitempty"`   // 电池sn
	Used        int    `json:"used,omitempty"`        // 使用天数
	IsAuthed    bool   `json:"isAuthed"`              // 是否实名认证 ture已实名 false未实名
	Intelligent bool   `json:"intelligent"`           // 是否智能套餐

	Station *EnterpriseStation `json:"station,omitempty"`    // 站点信息
	Ebike   *Ebike             `json:"ebikeBrand,omitempty"` // 电车品牌信息
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
	Logs []AgentRiderLog `json:"logs"`
}
