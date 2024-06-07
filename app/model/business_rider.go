// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-03
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import jsoniter "github.com/json-iterator/go"

type BusinessSubscribeReq struct {
	ID        uint64  `json:"id" validate:"required" trans:"订阅ID"`
	StoreID   *uint64 `json:"storeId" trans:"门店ID"`
	CabinetID *uint64 `json:"cabinetId" trans:"电柜ID"`
	AgentID   *uint64 `json:"agentId"` // 代理商ID

	RefundDeposit *bool    `json:"refundDeposit" trans:"是否退押金"` // 是否退押金(后台使用)
	DepositAmount *float64 `json:"depositAmount"`               // 退押金金额(后台使用)
	Remark        *string  `json:"remark"`                      // 备注

	Rto *bool `json:"rto"` // 管理员强制退租 - 是否参与以租代购

	RtoRemark    *string `json:"rtoRemark"`                     // 以租代购备注
	EbikeStoreID *uint64 `json:"ebikeStoreID" trans:"电车退租门店ID"` // 强制退租电车选择门店ID
	BatStoreID   *uint64 `json:"batStoreID" trans:"电池退租门店ID"`   // 强制退租电池选择门店ID
}

type BusinessCabinetReq struct {
	ID     uint64 `json:"id" validate:"required" trans:"订阅ID"`
	Serial string `json:"serial" validate:"required" trans:"电柜编码"`
}

type BusinessCabinetStatus struct {
	UUID  string `json:"uuid"`  // 操作ID, 使用此参数轮询获取状态
	Index int    `json:"index"` // 仓位Index, +1是仓位号
}

type BusinessCabinetStatusReq struct {
	UUID string `json:"uuid" validate:"required" query:"uuid" trans:"操作ID"`
}

type BusinessCabinetStatusRes struct {
	Success bool   `json:"success"` // 是否成功
	Stop    bool   `json:"stop"`    // 是否终止
	Message string `json:"message"` // 失败消息
}

func (r *BusinessCabinetStatusRes) MarshalBinary() ([]byte, error) {
	return jsoniter.Marshal(r)
}

func (r *BusinessCabinetStatusRes) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, r)
}

type BusinessPauseInfoRes struct {
	Start     string `json:"start,omitempty"`   // 寄存开始日期, 若此字段和`end`都不存在或为空时, 前端`寄存开始日期和结束寄存日期`显示未生效
	End       string `json:"end,omitempty"`     // 寄存结束日期
	Days      int    `json:"days"`              // 寄存天数
	Overdue   int    `json:"overdue,omitempty"` // 超期天数, 当此字段不存在时或为空时, 前端不显示`超出单词最长寄存时长`
	Remaining int    `json:"remaining"`         // 剩余天数
}

type BusinessRiderServiceDoReq struct {
	Rto  bool         `json:"rto"` // 是否满足并参与以租代购
	Type BusinessType `json:"bt"`  // 业务类型
}
