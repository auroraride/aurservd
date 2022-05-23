// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "time"
)

const (
    CabinetBinBatteryFault = "有电池无电压"
)

// CabinetBrand 电柜品牌
type CabinetBrand string

const (
    CabinetBrandKaixin  CabinetBrand = "KAIXIN"  // 凯信电柜
    CabinetBrandYundong CabinetBrand = "YUNDONG" // 云动电柜
)

func (cb CabinetBrand) Value() string {
    return string(cb)
}

func (cb CabinetBrand) String() string {
    switch cb {
    case CabinetBrandKaixin:
        return "凯信"
    case CabinetBrandYundong:
        return "云动"
    }
    return "Unknown"
}

// CabinetStatus 设备状态
const (
    CabinetStatusPending     uint8 = iota // 未投放
    CabinetStatusNormal                   // 运营中
    CabinetStatusMaintenance              // 维护中
)

// 设备健康状态
const (
    CabinetHealthStatusOffline uint8 = iota // 离线
    CabinetHealthStatusOnline               // 在线
    CabinetHealthStatusFault                // 故障
)

// Cabinet 电柜基础属性
type Cabinet struct {
    BranchID *uint64      `json:"branchId"`                                                      // 网点
    Status   uint8        `json:"status" enums:"0,1,2"`                                          // 电柜状态 0未投放 1运营中 2维护中
    Brand    CabinetBrand `json:"brand" validate:"required" trans:"品牌" enums:"KAIXIN,YUNDONG"` // KAIXIN(凯信) YUNDONG(云动)
    Serial   string       `json:"serial" validate:"required" trans:"电柜原始编码"`
    Name     string       `json:"name" validate:"required" trans:"电柜名称"`
    Doors    uint         `json:"doors" validate:"required" trans:"柜门数量"`
    Remark   *string      `json:"remark" trans:"备注"`
}

// CabinetCreateReq 电柜创建请求
type CabinetCreateReq struct {
    Cabinet
    Models []uint64 `json:"models" trans:"电池型号" validate:"required"`
}

// CabinetItem 电柜属性
type CabinetItem struct {
    ID uint64 `json:"id"` // 电柜ID
    Sn string `json:"sn"` // 平台编码
    Cabinet
    Models    []BatteryModel `json:"models"`              // 电池型号
    City      *City          `json:"city,omitempty"`      // 城市
    CreatedAt time.Time      `json:"createdAt,omitempty"` // 创建时间
}

// CabinetQueryReq 电柜查询请求
type CabinetQueryReq struct {
    PaginationReq

    Serial *string `json:"serial" query:"serial"`
    Name   *string `json:"name" query:"name"`
    CityId *uint64 `json:"cityId" query:"cityId"`
    Brand  *string `json:"brand" query:"brand"`
    Status *uint8  `json:"status" query:"status"`
}

// CabinetModifyReq 电柜修改请求
type CabinetModifyReq struct {
    ID       uint64        `json:"id" param:"id"`
    BranchID *uint64       `json:"branchId"`                                  // 网点
    Status   *uint8        `json:"status" enums:"0,1,2"`                      // 电柜状态 0未投放 1运营中 2维护中
    Brand    *CabinetBrand `json:"brand" trans:"品牌" enums:"KAIXIN,YUNDONG"` // KAIXIN(凯信) YUNDONG(云动)
    Serial   *string       `json:"serial" trans:"电柜原始编码"`
    Name     *string       `json:"name" trans:"电柜名称"`
    Doors    *uint         `json:"doors" trans:"柜门数量"`
    Remark   *string       `json:"remark" trans:"备注"`
    Models   *[]uint64     `json:"models" trans:"电池型号" validate:"required"`
}

// CabinetDeleteReq 电柜删除请求
type CabinetDeleteReq struct {
    ID uint64 `json:"id" param:"id"`
}

// CabinetBin 仓位详细信息
// 1000mA = 1A
// 1000mV = 1V
// (锁定状态 / 备注信息) 需要携带到下次的状态更新中
type CabinetBin struct {
    Name          string             `json:"name"`          // 柜门名称
    BatterySN     string             `json:"batterySN"`     // 电池序列号
    Full          bool               `json:"full"`          // 是否满电
    Battery       bool               `json:"battery"`       // 是否有电池
    Electricity   BatteryElectricity `json:"electricity"`   // 当前电量
    OpenStatus    bool               `json:"openStatus"`    // 是否开门
    DoorHealth    bool               `json:"doorHealth"`    // 是否锁仓 (柜门是否正常)
    Current       float64            `json:"current"`       // 充电电流(A)
    Voltage       float64            `json:"voltage"`       // 电压(V)
    ChargerErrors []string           `json:"chargerErrors"` // 故障信息
    Remark        string             `json:"remark"`        // 备注
}

// CabinetSnParamReq sn请求
type CabinetSnParamReq struct {
    Sn string `json:"sn" param:"sn" validate:"required"`
}

// CabinetDetailRes 电柜详细信息返回
type CabinetDetailRes struct {
    CabinetItem
    Bin []CabinetBin `json:"bins"` // 仓位信息
}

// CanUse 仓位是否可以换电
func (cb CabinetBin) CanUse() bool {
    return cb.Battery && cb.Electricity.IsBatteryFull() && !cb.OpenStatus && cb.DoorHealth && len(cb.ChargerErrors) == 0
}

// CabinetDoorOperate 柜门操作
type CabinetDoorOperate uint

const (
    CabinetDoorOperateOpen   CabinetDoorOperate = iota + 1 // 开仓
    CabinetDoorOperateLock                                 // 锁定(标记为故障)
    CabinetDoorOperateUnlock                               // 解锁(取消标记故障)
)

func (cdo CabinetDoorOperate) String() string {
    switch cdo {
    case CabinetDoorOperateOpen:
        return "开仓"
    case CabinetDoorOperateLock:
        return "锁定"
    case CabinetDoorOperateUnlock:
        return "解锁"
    }
    return ""
}

var CabinetDoorOperates = map[CabinetDoorOperate]map[CabinetBrand]string{
    CabinetDoorOperateOpen: {
        CabinetBrandKaixin:  "1",
        CabinetBrandYundong: "opendoor",
    },
    CabinetDoorOperateLock: {
        CabinetBrandKaixin:  "3",
        CabinetBrandYundong: "disabledoor",
    },
    CabinetDoorOperateUnlock: {
        CabinetBrandKaixin:  "4",
        CabinetBrandYundong: "enabledoor",
    },
}

// Value 获取柜门操作值
func (cdo CabinetDoorOperate) Value(brand CabinetBrand) (v string, ex bool) {
    v, ex = CabinetDoorOperates[cdo][brand]
    return
}

// CabinetDoorOperateReq 仓门操作
// TODO 给用户开仓逻辑
type CabinetDoorOperateReq struct {
    ID        *uint64             `json:"id" validate:"required"`        // 电柜ID
    Index     *int                `json:"index" validate:"required"`     // 仓门index
    Remark    string              `json:"remark"`                        // 操作原因
    Operation *CabinetDoorOperate `json:"operation" validate:"required"` // 操作方式 1:开仓 2:锁定(标记为故障) 3:解锁(取消标记故障)
    Phone     *string             `json:"phone"`                         // 骑手手机号
}
