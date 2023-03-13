// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "fmt"
    "github.com/auroraride/adapter"
    "sort"
)

const (
    CabinetBinBatteryFault = "有电池无电压"
)

type CabinetStatus uint8

const (
    CabinetStatusPending     CabinetStatus = iota // 未投放
    CabinetStatusNormal                           // 运营中
    CabinetStatusMaintenance                      // 维护中
)

func (cs CabinetStatus) String() string {
    switch cs {
    case CabinetStatusPending:
        return "未投放"
    case CabinetStatusNormal:
        return "运营中"
    case CabinetStatusMaintenance:
        return "维护中"
    }
    return "未知"
}

func (cs CabinetStatus) Value() uint8 {
    return uint8(cs)
}

// 设备健康状态
const (
    CabinetHealthStatusOffline uint8 = iota // 离线
    CabinetHealthStatusOnline               // 在线
    CabinetHealthStatusFault                // 故障
)

// Cabinet 电柜基础属性
type Cabinet struct {
    BranchID    *uint64              `json:"branchId"`                                                              // 网点
    Status      CabinetStatus        `json:"status" enums:"0,1,2"`                                                  // 电柜状态 0未投放 1运营中 2维护中
    Brand       adapter.CabinetBrand `json:"brand" validate:"required" trans:"品牌" enums:"KAIXIN,YUNDONG,TUOBANG"` // KAIXIN(凯信) YUNDONG(云动) TUOBANG(拓邦)
    Serial      string               `json:"serial" validate:"required" trans:"电柜编码"`
    Name        string               `json:"name" validate:"required" trans:"电柜名称"`
    Doors       int                  `json:"doors"` // 柜门数量
    Remark      *string              `json:"remark" trans:"备注"`
    Health      *uint8               `json:"health"`      // 在线状态 0离线 1在线 2故障
    Intelligent bool                 `json:"intelligent"` // 是否智能柜
}

type CabinetBasicInfo struct {
    ID     uint64               `json:"id"`
    Brand  adapter.CabinetBrand `json:"brand" enums:"KAIXIN,YUNDONG,TUOBANG"` // 品牌: KAIXIN(凯信) YUNDONG(云动) TUOBANG(拓邦)
    Serial string               `json:"serial"`                               // 电柜编码
    Name   string               `json:"name"`                                 // 电柜名称
}

type CabinetListByDistanceRes struct {
    CabinetBasicInfo
    Status uint8 `json:"status"` // 电柜状态 0未投放 1运营中 2维护中
    Health uint8 `json:"health"` // 在线状态 0离线 1在线 2故障
}

// CabinetCreateReq 电柜创建请求
type CabinetCreateReq struct {
    Cabinet
    Models  []string `json:"models" trans:"电池型号" validate:"required"`
    SimSn   string   `json:"simSn,omitempty"`   // SIM卡号
    SimDate string   `json:"simDate,omitempty"` // SIM卡到期日期, 例: 2022-06-01
}

// CabinetItem 电柜属性
type CabinetItem struct {
    ID uint64 `json:"id"` // 电柜ID
    Sn string `json:"sn"` // 平台编码
    Cabinet
    Models      []string `json:"models"`              // 电池型号
    City        *City    `json:"city,omitempty"`      // 城市
    CreatedAt   string   `json:"createdAt,omitempty"` // 创建时间
    SimSn       string   `json:"simSn,omitempty"`     // SIM卡号
    SimDate     string   `json:"simDate,omitempty"`   // SIM卡到期日期, 例: 2022-06-01
    Transferred bool     `json:"transferred"`         // 是否初始化过调拨
    BatteryNum  int      `json:"batteryNum"`          // 电池数量
    Intelligent bool     `json:"intelligent"`         // 是否智能柜
}

// CabinetQueryReq 电柜查询请求
type CabinetQueryReq struct {
    PaginationReq

    Serial      *string               `json:"serial" query:"serial"`           // 电柜编号
    Name        *string               `json:"name" query:"name"`               // 电柜名称
    CityID      *uint64               `json:"cityId" query:"cityId"`           // 城市ID
    Brand       *adapter.CabinetBrand `json:"brand" query:"brand"`             // 电柜型号 KAIXIN(凯信) YUNDONG(云动) TUOBANG(拓邦)
    Status      *uint8                `json:"status" query:"status"`           // 电柜状态
    Model       *string               `json:"model" query:"model"`             // 电池型号
    Online      uint8                 `json:"online" query:"online"`           // 在线状态
    Intelligent uint8                 `json:"intelligent" query:"intelligent"` // 是否智能柜 0:全部 1:是 2:否
}

// CabinetModifyReq 电柜修改请求
type CabinetModifyReq struct {
    ID          uint64                `json:"id" param:"id"`
    BranchID    *uint64               `json:"branchId"`                                  // 网点
    Status      *CabinetStatus        `json:"status" enums:"0,1,2"`                      // 电柜状态 0未投放 1运营中 2维护中
    Brand       *adapter.CabinetBrand `json:"brand" trans:"品牌" enums:"KAIXIN,YUNDONG"` // KAIXIN(凯信) YUNDONG(云动) TUOBANG(拓邦)
    Serial      *string               `json:"serial" trans:"电柜原始编码"`
    Name        *string               `json:"name" trans:"电柜名称"`
    Doors       *uint                 `json:"doors" trans:"柜门数量"`
    Remark      *string               `json:"remark" trans:"备注"`
    Models      *[]string             `json:"models" trans:"电池型号"`
    SimSn       *string               `json:"simSn,omitempty"`       // SIM卡号
    SimDate     *string               `json:"simDate,omitempty"`     // SIM卡到期日期, 例: 2022-06-01
    Intelligent *bool                 `json:"intelligent,omitempty"` // 是否智能柜
}

// CabinetDeleteReq 电柜删除请求
type CabinetDeleteReq struct {
    ID uint64 `json:"id" param:"id"`
}

type CabinetBinsMap map[int]*CabinetBin

type CabinetBins []*CabinetBin

// Sort 仓位根据使用条件排序
func (cbs CabinetBins) Sort() {
    sort.Slice(cbs, func(i, j int) bool {
        return cbs[i].Electricity.Value()+cbs[i].Voltage*0.1 > cbs[j].Electricity.Value()+cbs[j].Voltage*0.1
    })
}

// MaxEmpty 获取满电和空仓
func (cbs CabinetBins) MaxEmpty() (max *CabinetBin, empty *CabinetBin) {
    cbs.Sort()

    for _, bin := range cbs {
        if !bin.DoorHealth {
            continue
        }
        if bin.Battery && max == nil {
            max = bin
        }
        // 2022年08月06日16:28:26 曹博文反馈说凯信的仓位开开发现有电池,查询日志所得: 15:43:58时CH7208KXHD220408027电柜的2号仓电池不在位但电流和电压不为空
        // 2023年02月08日15:29:38 曹博文让改成最低电压小于45
        if !bin.Battery && bin.Voltage < 45 && bin.Electricity == 0 && empty == nil {
            empty = bin
        }
        if max != nil && empty != nil {
            return
        }
    }

    return
}

// CabinetBin 仓位详细信息
// 1000mA = 1A
// 1000mV = 1V
// (锁定状态 / 备注信息) 需要携带到下次的状态更新中
type CabinetBin struct {
    Index         int        `json:"index"`                   // 仓位index (从0开始)
    Name          string     `json:"name"`                    // 柜门名称
    BatterySN     string     `json:"batterySN"`               // 电池序列号
    Full          bool       `json:"full"`                    // 是否满电
    Battery       bool       `json:"battery"`                 // 是否有电池
    Electricity   BatterySoc `json:"electricity"`             // 当前电量
    OpenStatus    bool       `json:"openStatus"`              // 是否开门
    DoorHealth    bool       `json:"doorHealth"`              // 是否锁仓 (柜门是否正常)
    Current       float64    `json:"current"`                 // 充电电流(A)
    Voltage       float64    `json:"voltage"`                 // 电压(V)
    ChargerErrors []string   `json:"chargerErrors,omitempty"` // 故障信息
    Remark        string     `json:"remark,omitempty"`        // 备注
}

// CabinetBinRemark 仓位备注
type CabinetBinRemark struct {
    Index  int    `json:"index"`  // 仓位index
    Remark string `json:"remark"` // 备注信息
}

// CabinetSnParamReq sn请求
type CabinetSnParamReq struct {
    Sn string `json:"sn" param:"sn" validate:"required"`
}

// CabinetDetailRes 电柜详细信息返回
type CabinetDetailRes struct {
    CabinetItem
    Bin      []CabinetBin         `json:"bins"`     // 仓位信息
    StockNum int                  `json:"stockNum"` // 库存电池
    Reserves []ReserveCabinetItem `json:"reserves"` // 当前预约
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

var CabinetDoorOperates = map[CabinetDoorOperate]map[adapter.CabinetBrand]string{
    CabinetDoorOperateOpen: {
        adapter.CabinetBrandKaixin:  "1",
        adapter.CabinetBrandYundong: "opendoor",
    },
    CabinetDoorOperateLock: {
        adapter.CabinetBrandKaixin:  "3",
        adapter.CabinetBrandYundong: "disabledoor",
    },
    CabinetDoorOperateUnlock: {
        adapter.CabinetBrandKaixin:  "4",
        adapter.CabinetBrandYundong: "enabledoor",
    },
}

// Value 获取柜门操作值
func (cdo CabinetDoorOperate) Value(brand adapter.CabinetBrand) (v string, ex bool) {
    v, ex = CabinetDoorOperates[cdo][brand]
    return
}

type CabinetDoorOperatorRole uint8

const (
    CabinetDoorOperatorRoleManager CabinetDoorOperatorRole = iota + 1 // 后台人员
    CabinetDoorOperatorRoleRider                                      // 骑手
)

func (cdor CabinetDoorOperatorRole) String() string {
    switch cdor {
    case CabinetDoorOperatorRoleManager:
        return "后台人员"
    case CabinetDoorOperatorRoleRider:
        return "骑手"
    }
    return "未知"
}

// CabinetDoorOperator 柜门操作人
type CabinetDoorOperator struct {
    ID    uint64                  `json:"id"`    // 用户ID
    Role  CabinetDoorOperatorRole `json:"role"`  // 角色
    Name  string                  `json:"name"`  // 姓名
    Phone string                  `json:"phone"` // 手机号
}

// CabinetDoorOperateReq 仓门操作
type CabinetDoorOperateReq struct {
    ID        uint64              `json:"id" validate:"required"`        // 电柜ID
    Index     *int                `json:"index" validate:"required"`     // 仓门index
    Remark    string              `json:"remark"`                        // 操作原因
    Operation *CabinetDoorOperate `json:"operation" validate:"required"` // 操作方式 1:开仓 2:锁定(标记为故障) 3:解锁(取消标记故障)
}

// BinInfo 任务电柜仓位信息
type BinInfo struct {
    Index       int        `json:"index" bson:"index"`             // 仓位index
    Electricity BatterySoc `json:"electricity" bson:"electricity"` // 电量
    Voltage     float64    `json:"voltage" bson:"voltage"`         // 电压(V)
}

func (b *BinInfo) String() string {
    return fmt.Sprintf(
        "%d号仓, 电压: %.2fV, 电流: %2.fA",
        b.Index+1,
        b.Voltage,
        b.Electricity,
    )
}

type YundongDeployInfo struct {
    SN       string  `json:"-"`
    AreaCode string  `json:"-"`
    Address  string  `json:"address"`
    Lat      float64 `json:"lat"`
    Lng      float64 `json:"lng"`
    Name     string  `json:"name"`
    Phone    string  `json:"phone"`
    Contact  string  `json:"contact"`
    City     string  `json:"city"`
}

type CabinetDataReq struct {
    PaginationReq

    Status uint8                `json:"status" enums:"0,1,2,3" query:"status"`      // 电柜状态 0:全部 1:在线 2:离线 3:锁仓
    Votage float64              `json:"votage" query:"votage"`                      // 电压型号筛选
    Name   string               `json:"name" query:"name"`                          // 电柜名称
    Serial string               `json:"serial" query:"serial"`                      // 电柜编号
    CityID uint64               `json:"cityId" query:"cityId"`                      // 城市
    Brand  adapter.CabinetBrand `json:"brand" enums:"KAIXIN,YUNDONG" query:"brand"` // 品牌 KAIXIN(凯信) YUNDONG(云动) TUOBANG(拓邦)
}

const (
    CabinetDataBinStatusEmpty    uint8 = iota // 无电池
    CabinetDataBinStatusCharging              // 充电中
    CabinetDataBinStatusFull                  // 满电
    CabinetDataBinStatusLock                  // 锁仓
)

type CabinetDataBin struct {
    Status uint8  `json:"status" enums:"0,1,2,3"` // 状态 0:无电池 1:有电池未满电 2:满电 3:锁仓, 状态值大的优先显示
    Remark string `json:"remark,omitempty"`       // 备注信息, 仅状态`3 锁仓`有此字段
}

type CabinetDataRes struct {
    ID         uint64               `json:"id"`
    Name       string               `json:"name"`       // 名称
    Serial     string               `json:"serial"`     // 编号
    Model      string               `json:"model"`      // 型号
    Brand      adapter.CabinetBrand `json:"brand"`      // 品牌 KAIXIN(凯信) YUNDONG(云动) TUOBANG(拓邦)
    Online     bool                 `json:"online"`     // 是否在线
    BinNum     int                  `json:"binNum"`     // 仓位数量
    BatteryNum int                  `json:"batteryNum"` // 电池数量
    EmptyNum   int                  `json:"emptyNum"`   // 空仓数量
    LockNum    int                  `json:"lockNum"`    // 锁仓数
    FullNum    int                  `json:"fullNum"`    // 满电数
    Bins       []CabinetDataBin     `json:"bins"`       // 仓位信息
    StockNum   int                  `json:"stockNum"`   // 库存电池
}

type CabinetSimNotice struct {
    Serial string `json:"serial"`
    Name   string `json:"name"`
    City   string `json:"city"`
    Sim    string `json:"sim"`
    End    string `json:"end"`
}

type CabinetTransferReq struct {
    CabinetID uint64 `json:"cabinetId" validate:"required" trans:"电柜ID"`
    Model     string `json:"model" validate:"required" trans:"型号"`
    Num       int    `json:"num" validate:"required" trans:"数量"`
}

type CabinetMaintainReq struct {
    ID       uint64 `json:"id" validate:"required"`       // 电柜ID
    Maintain *bool  `json:"maintain" validate:"required"` // 是否维护
}

type CabinetSerialQueryReq struct {
    Serial string `json:"serial" query:"serial" validate:"required" trans:"电柜编码"`
}

type CabinetOpenBindReq struct {
    Phone     string `json:"phone" validate:"required" trans:"骑手电话"`
    ID        uint64 `json:"id" validate:"required" trans:"电柜ID"`
    Index     *int   `json:"index" validate:"required" trans:"仓门序号"`
    Remark    string `json:"remark" validate:"required" trans:"操作原因"` // 可手动输入, 预留项: `吞电池` `柜门未关好` `仓位不足`
    BatterySN string `json:"batterySn" validate:"required" trans:"电池编码"`
}

type CabinetRpcBatchRequestItem struct {
    Serial      string `json:"serial"`
    Intelligent bool   `json:"intelligent"`
}
