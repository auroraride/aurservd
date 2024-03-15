// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type ReserveStatus uint8

const (
	ReserveStatusPending    ReserveStatus = iota // 已预约
	ReserveStatusProcessing                      // 进行中
	ReserveStatusSuccess                         // 已完成
	ReserveStatusFail                            // 已失败
	ReserveStatusTimeout                         // 已超时
	ReserveStatusCancel                          // 已取消
	ReserveStatusInvalid                         // 已失效
)

func (rs ReserveStatus) Value() uint8 {
	return uint8(rs)
}

func (rs ReserveStatus) String() string {
	switch rs {
	default:
		return "已预约"
	case ReserveStatusProcessing:
		return "进行中"
	case ReserveStatusSuccess:
		return "已完成"
	case ReserveStatusFail:
		return "已失败"
	case ReserveStatusTimeout:
		return "已超时"
	case ReserveStatusCancel:
		return "已取消"
	case ReserveStatusInvalid:
		return "已失效"
	}
}

type ReserveCreateReq struct {
	CabinetID uint64 `json:"cabinetId"`                                           // 电柜ID
	Business  string `json:"business"  enums:"active,pause,continue,unsubscribe"` // 业务选项 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
}

type ReserveUnfinishedRes struct {
	ID        uint64        `json:"id"`        // 预约ID
	CabinetID uint64        `json:"cabinetId"` // 电柜ID
	Fid       string        `json:"fid"`       // 设施ID
	Business  string        `json:"business"`  // 预约业务 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
	Time      string        `json:"time"`      // 预约时间
	Status    ReserveStatus `json:"status"`    // 状态 0:已预约 1:进行中
}

type ReserveListReq struct {
	PaginationReq
	ReserveListFilter
}

type ReserveListFilter struct {
	CityID    uint64 `json:"cityId" query:"cityId"`                                               // 城市ID
	RiderID   uint64 `json:"riderId" query:"riderId"`                                             // 骑手ID
	Start     string `json:"start" query:"start"`                                                 // 开始日期, 如 2022-08-01
	End       string `json:"end" query:"end"`                                                     // 结束日期, 如 2022-08-07
	CabinetID uint64 `json:"cabinetId" query:"cabinetId"`                                         // 电柜ID
	Business  string `json:"business" query:"business" enums:"active,pause,continue,unsubscribe"` // 预约业务 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
	Serial    string `json:"serial" query:"serial"`                                               // 电柜编号
}

type ReserveListRes struct {
	City        string `json:"city"`        // 城市
	Name        string `json:"name"`        // 姓名
	Phone       string `json:"phone"`       // 电话
	Business    string `json:"business"`    // 业务
	Status      string `json:"status"`      // 状态
	CabinetName string `json:"cabinetName"` // 电柜名称
	Serial      string `json:"serial"`      // 电柜编号
	CreatedAt   string `json:"createdAt"`   // 创建时间
}

type ReserveCabinetItem struct {
	Name     string `json:"name"`                                               // 姓名
	Phone    string `json:"phone"`                                              // 电话
	Business string `json:"business" enums:"active,pause,continue,unsubscribe"` // 业务 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
	Time     string `json:"time"`                                               // 预约时间
}
