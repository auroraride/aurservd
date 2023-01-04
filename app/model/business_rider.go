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
