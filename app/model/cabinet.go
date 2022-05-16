// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "time"

// CabinetBrand 电柜品牌
type CabinetBrand string

const (
    CabinetBrandKaixin  CabinetBrand = "KAIXIN"  // 凯信电柜
    CabinetBrandYundong CabinetBrand = "YUNDONG" // 云动电柜
)

func (cb CabinetBrand) String() string {
    return string(cb)
}

// CabinetStatus 设备状态
type CabinetStatus uint

const (
    CabinetStatusPending     CabinetStatus = iota // 未投放
    CabinetStatusOK                               // 运营中
    CabinetStatusMaintenance                      // 维护中
)

// Cabinet 电柜基础属性
type Cabinet struct {
    BranchID *uint64       `json:"branchId"`                                                    // 网点
    Status   CabinetStatus `json:"status" enums:"0,1,2"`                                        // 电柜状态 0未投放 1运营中 2维护中
    Brand    CabinetBrand  `json:"brand" validate:"required" trans:"品牌" enums:"KAIXIN,YUNDONG"` // KAIXIN(凯信) YUNDONG(云动)
    Serial   string        `json:"serial" validate:"required" trans:"电柜原始编码"`
    Name     string        `json:"name" validate:"required" trans:"电柜名称"`
    Doors    uint          `json:"doors" validate:"required" trans:"柜门数量"`
    Remark   *string       `json:"remark" trans:"备注"`
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
    Models    []BatteryModel `json:"models"`    // 电池型号
    City      City           `json:"city"`      // 城市
    CreatedAt time.Time      `json:"createdAt"` // 创建时间
}

// CabinetQueryReq 电柜查询请求
type CabinetQueryReq struct {
    PaginationReq

    Serial *string        `json:"serial" query:"serial"`
    Name   *string        `json:"name" query:"name"`
    CityId *uint64        `json:"cityId" query:"cityId"`
    Brand  *string        `json:"brand" query:"brand"`
    Status *CabinetStatus `json:"status" query:"status"`
}

// CabinetModifyReq 电柜修改请求
type CabinetModifyReq struct {
    ID       uint64         `json:"id" param:"id"`
    BranchID *uint64        `json:"branchId"`                                // 网点
    Status   *CabinetStatus `json:"status" enums:"0,1,2"`                    // 电柜状态 0未投放 1运营中 2维护中
    Brand    *CabinetBrand  `json:"brand" trans:"品牌" enums:"KAIXIN,YUNDONG"` // KAIXIN(凯信) YUNDONG(云动)
    Serial   *string        `json:"serial" trans:"电柜原始编码"`
    Name     *string        `json:"name" trans:"电柜名称"`
    Doors    *uint          `json:"doors" trans:"柜门数量"`
    Remark   *string        `json:"remark" trans:"备注"`
    Models   *[]uint64      `json:"models" trans:"电池型号" validate:"required"`
}

// CabinetDeleteReq 电柜删除请求
type CabinetDeleteReq struct {
    ID uint64 `json:"id" param:"id"`
}

// CabinetBin 仓位详细信息
// TODO: (锁定状态 / 备注信息) 需要携带到下次的状态更新中
type CabinetBin struct {
    Name          string  `json:"name"`          // 柜门名称
    SN            string  `json:"sn"`            // 序列号
    Locked        bool    `json:"locked"`        // 是否锁定
    Full          bool    `json:"full"`          // 是否满电
    Status        bool    `json:"status"`        // 是否有电池
    Electricity   float64 `json:"electricity"`   // 当前电量
    OpenStatus    uint    `json:"openStatus"`    // 开门状态 0关闭 1开启
    DoorHealth    uint    `json:"doorHealth"`    // 柜门状态 0正常 1故障
    Current       float64 `json:"current"`       // 充电电流(A) 1000mA = 1A
    Voltage       float64 `json:"voltage"`       // 电压(V) 1000mV = 1V
    ChargerStatus uint    `json:"chargerStatus"` // 充电器状态 0:正常 1:电池充电过慢 2:电池充电过快 3:220V 丢失/充电器损坏 4:充电器状态错误 5:电池未连接到充电器 6:行程开关故障 7:充电触点接触不良 8:电池无法充满 9:电池无法充电 10:充电器通讯故障 11:行程开关接触不良 12:已取出，未解绑
    Remark        string  `json:"remark"`        // 备注
}
