// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "database/sql/driver"

const (
	BusinessAimedAll        uint8 = iota // 全部
	BusinessAimedPersonal                // 个签
	BusinessAimedEnterprise              // 团签
)

type BusinessType string

const (
	BusinessTypeActive      BusinessType = "active"      // 激活
	BusinessTypePause       BusinessType = "pause"       // 寄存
	BusinessTypeContinue    BusinessType = "continue"    // 结束寄存
	BusinessTypeUnsubscribe BusinessType = "unsubscribe" // 退租
)

func (b BusinessType) String() string {
	return string(b)
}

func (b BusinessType) Value() (driver.Value, error) {
	return string(b), nil
}

func (b *BusinessType) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		return nil
	case string:
		*b = BusinessType(v)
	}
	return nil
}

func BusinessTypeText(bt string) string {
	return map[string]string{
		"active":      "激活",
		"pause":       "寄存",
		"continue":    "结束寄存",
		"unsubscribe": "退租",
	}[bt]
}

type BusinessIsRto uint8

const (
	BusinessIsRtoUnSend BusinessIsRto = iota // 未赠送
	BusinessIsRtoSend                        // 已赠送
)

func (s BusinessIsRto) Value() uint8 {
	return uint8(s)
}

func (s BusinessIsRto) String() string {
	switch s {
	case BusinessIsRtoUnSend:
		return "未赠送"
	case BusinessIsRtoSend:
		return "已赠送"
	}
	return " - "
}

type BusinessSubscribeID struct {
	SubscribeID uint64 `json:"subscribeId" validate:"required" trans:"订阅ID"`
}

type BusinessFilter struct {
	Goal         StoreCabiletGoal `json:"goal" query:"goal" enums:"0,1,2"`   // 查询目标, 0:不筛选 1:门店 2:电柜
	EmployeeID   uint64           `json:"employeeId" query:"employeeId"`     // 店员ID, 店员端请求忽略此参数
	EnterpriseID uint64           `json:"enterpriseId" query:"enterpriseId"` // 企业ID
	Aimed        uint8            `json:"aimed" query:"aimed"`               // 筛选业务对象 0:全部 1:个签 2:团签
	Start        *string          `json:"start" query:"start"`               // 筛选开始日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
	End          *string          `json:"end" query:"end"`                   // 筛选结束日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
	Keyword      *string          `json:"keyword" query:"keyword"`           // 筛选骑手姓名或电话
	StoreID      uint64           `json:"storeId" query:"storeId"`           // 筛选门店
	CabinetID    uint64           `json:"cabinetId" query:"cabinetId"`       // 筛选电柜
	CityID       uint64           `json:"cityId" query:"cityId"`             // 筛选城市
	// 筛选业务类别 active:激活 pause:寄存 continue:结束寄存 unsubscribe:退订
	Type *string `json:"type" enums:"active,pause,continue,unsubscribe" query:"type"`
}

type BusinessListReq struct {
	PaginationReq
	BusinessFilter
}

type BusinessExportReq struct {
	BusinessFilter
	Remark string `json:"remark" validate:"required" trans:"备注"`
}

type BusinessEmployeeListRes struct {
	ID                uint64             `json:"id"`
	Name              string             `json:"name"`                        // 骑手姓名
	Phone             string             `json:"phone"`                       // 骑手电话
	Type              string             `json:"type"`                        // 业务类别
	Time              string             `json:"time"`                        // 业务时间
	City              string             `json:"city"`                        // 城市
	Plan              *Plan              `json:"plan,omitempty"`              // 骑士卡, 团签无此字段
	Enterprise        *Enterprise        `json:"enterprise,omitempty"`        // 团签企业, 个签无此字段
	EnterpriseStation *EnterpriseStation `json:"enterpriseStation,omitempty"` // 站点
	IsRto             uint8              `json:"isRto"`                       // 是否已赠车
	Remark            string             `json:"remark"`                      // 备注
}

type BusinessListRes struct {
	BusinessEmployeeListRes
	Operator string            `json:"operator"`           // 操作人
	Employee *Employee         `json:"employee,omitempty"` // 店员, 可能为空
	Cabinet  *CabinetBasicInfo `json:"cabinet,omitempty"`  // 电柜, 可能为空
	Store    *Store            `json:"store,omitempty"`    // 门店, 可能为空
}

type BusinessPauseFilter struct {
	CityID  uint64 `json:"cityId" query:"cityId"`   // 城市ID
	RiderID uint64 `json:"riderId" query:"riderId"` // 骑手ID
	Status  uint8  `json:"status" query:"status"`   // 状态 0:全部 1:寄存中 2:已结束
	Overdue bool   `json:"overdue" query:"overdue"` // 逾期 false:全部 true:已逾期

	StartAscription uint8 `json:"startAscription" query:"startAscription" enums:"0,1,2"` // 寄存类别 0:全部 1:门店 2:电柜
	EndAscription   uint8 `json:"endAscription" query:"endAscription" enums:"0,1,2"`     // 取电类别 0:全部 1:门店 2:电柜

	StartDate string `json:"startDate" query:"startDate"` // 寄存时间段, 逗号分隔, 如2022-08-01,2022-08-07
	EndDate   string `json:"endDate" query:"endDate"`     // 取电时间段, 逗号分隔, 如2022-08-01,2022-08-07

	StartBy string `json:"startBy" query:"startBy"` // 寄存办理人
	EndBy   string `json:"endBy" query:"endBy"`     // 取电办理人

	StartTarget string `json:"startTarget" query:"startTarget"` // 寄存地点
	EndTarget   string `json:"endTarget" query:"endTarget"`     // 结束地点
}

type BusinessPauseList struct {
	PaginationReq
	BusinessPauseFilter
}

type BusinessPauseListRes struct {
	Status          string `json:"status"`          // 当前状态
	City            string `json:"city"`            // 城市
	Name            string `json:"name"`            // 骑手姓名
	Phone           string `json:"phone"`           // 骑手电话
	Plan            string `json:"subscribe"`       // 骑士卡
	Start           string `json:"start"`           // 寄存开始时间
	StartTarget     string `json:"startTarget"`     // 寄存地点
	StartAscription string `json:"startAscription"` // 寄存类别
	StartBy         string `json:"startBy"`         // 寄存办理
	End             string `json:"end"`             // 取电时间
	EndTarget       string `json:"endTarget"`       // 取电地点
	EndAscription   string `json:"endAscription"`   // 取电类别
	EndBy           string `json:"endBy"`           // 取电办理
	Days            int    `json:"days"`            // 寄存天数
	OverdueDays     int    `json:"overdueDays"`     // 超期天数
	Remaining       int    `json:"remaining"`       // 剩余天数
	SuspendDays     int    `json:"suspendDays"`     // 暂停扣费天数
}
