// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-05
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/defs/xcdef"
	"github.com/auroraride/adapter/rpc/pb"

	"github.com/auroraride/aurservd/pkg/silk"
)

type BmsBattery struct {
	Status BatteryStatus `json:"status"` // 状态, 0:静置 1:充电 2:放电 3:异常(此时faults字段存在)
	Soc    uint32        `json:"soc"`    // 剩余容量, 单位1%

	Charge    *bool         `json:"charge"`           // 充电是否开启
	DisCharge *bool         `json:"disCharge"`        // 放电是否开启
	Faults    *xcdef.Faults `json:"faults,omitempty"` // 故障列表, `0`:总压低, `1`:总压高, `2`:单体低, `3`:单体高, `6`:放电过流, `7`:充电过流, `8`:SOC低, `11`:充电高温, `12`:充电低温, `13`:放电高温, `14`:放电低温, `15`:短路, `16`:MOS高温
}

func NewBmsBattery(hb *pb.BatteryHeartbeat) (item *BmsBattery) {
	item = new(BmsBattery)

	item.Soc = hb.Soc

	if hb.MosStatus != nil {
		ms := new(xcdef.MosStatus)
		ms.FromBytes(hb.MosStatus)
		item.Charge = silk.Pointer(ms.CanCharge())
		item.DisCharge = silk.Pointer(ms.CanDisCharge())
	}

	switch {
	default:
		item.Status = BatteryStatusIdle
	case hb.Current > 0:
		item.Status = BatteryStatusCharging
	case hb.Current < 0:
		item.Status = BatteryStatusDisCharging
	}

	if hb.Faults != nil {
		faults := new(xcdef.Faults)
		faults.FromBytes(hb.Faults)
		item.Faults = faults
		item.Status = BatteryStatusFault
	}

	return
}

type BatterySNRequest struct {
	SN string `json:"sn" param:"sn" validate:"required"` // 电池编码
}

type BatteryBmsDetail struct {
	*BmsBattery
	// 电池编号
	Sn string `json:"sn,omitempty"`
	// 入库时间
	CreatedAt string `json:"createdAt"`
	// 最后通讯时间
	UpdatedAt string `json:"updatedAt"`
	// 电池总压 (V)
	Voltage float64 `json:"voltage,omitempty"`
	// 电流 (A, 充电为正, 放电为负)
	Current float64 `json:"current,omitempty"`
	// 健康度 单位1%
	Soh uint8 `json:"soh,omitempty"`
	// 是否在电柜
	InCabinet bool `json:"inCabinet,omitempty"`
	// 剩余容量 (单位AH)
	Capacity float64 `json:"capacity,omitempty"`
	// 最大单体电压 (mV)
	MonMaxVoltage uint16 `json:"monMaxVoltage,omitempty"`
	// 最大单体电压位置 (第x串)
	MonMaxVoltagePos uint8 `json:"monMaxVoltagePos,omitempty"`
	// 最小单体电压 (mV)
	MonMinVoltage uint16 `json:"monMinVoltage,omitempty"`
	// 最小单体电压位置 (第x串)
	MonMinVoltagePos uint8 `json:"monMinVoltagePos,omitempty"`
	// 最大温度 (单位1℃)
	MaxTemp uint16 `json:"maxTemp,omitempty"`
	// 最小温度 (单位1℃)
	MinTemp uint16 `json:"minTemp,omitempty"`
	// MOS状态 (Bit0表示充电, Bit1表示放电, 此字段无法判定电池是否充放电状态)
	MosStatus *xcdef.MosStatus `json:"mosStatus,omitempty"`
	// 单体电压 (24个单体电压, 单位mV)
	MonVoltage *xcdef.MonVoltage `json:"monVoltage,omitempty"`
	// 电池温度 (4个电池温度传感器, 单位1℃)
	Temp *xcdef.Temperature `json:"temp,omitempty"`
	// MOS温度 (1个MOS温度传感器, 单位1℃)
	MosTemp uint16 `json:"mosTemp,omitempty"`
	// 环境温度 (1个环境温度传感器, 单位1℃)
	EnvTemp uint16 `json:"envTemp,omitempty"`
	// 坐标
	Geom *adapter.Geometry `json:"geom,omitempty"`
	// GPS定位状态 (0=未定位 1=GPS定位 4=LBS定位)
	Gps xcdef.GPSStatus `json:"gps,omitempty"`
	// 4G通讯信号强度 (0-100 百分比形式)
	Strength uint8 `json:"strength,omitempty"`
	// 电池包循环次数 (80%累加一次)
	Cycles uint16 `json:"cycles,omitempty"`
	// 本次充电时长
	ChargingTime uint64 `json:"chargingTime,omitempty"`
	// 本次放电时长
	DisChargingTime uint64 `json:"disChargingTime,omitempty"`
	// 本次使用时长
	UsingTime uint64 `json:"usingTime,omitempty"`
	// 总充电时长
	TotalChargingTime uint64 `json:"totalChargingTime,omitempty"`
	// 总放电时长
	TotalDisChargingTime uint64 `json:"totalDisChargingTime,omitempty"`
	// 总使用时长
	TotalUsingTime uint64 `json:"totalUsingTime,omitempty"`
	// 当前位置
	BelongsTo string `json:"belongsTo,omitempty"`
	// 功率 (Kw)
	Power float64 `json:"power,omitempty"`
	// BMS软件版本
	SoftVersion uint32 `json:"softVersion,omitempty"`
	// BMS硬件版本
	HardVersion uint32 `json:"hardVersion,omitempty"`
	// 4G软件版本
	Soft4gVersion uint32 `json:"soft4GVersion,omitempty"`
	// 4G硬件版本
	Hard4gVersion uint32 `json:"hard4GVersion,omitempty"`
	// 4G板SN
	Sn4g uint64 `json:"sn4G,omitempty"`
	// SIM卡ICCID
	Iccid string `json:"iccid,omitempty"`
	// 电池是否在线
	Online bool `json:"online"`
	// 故障统计, 需要将13种故障都显示出来, 若无返回则数量为0 (`0`:总压低, `1`:总压高, `2`:单体低, `3`:单体高, `6`:放电过流, `7`:充电过流, `8`:SOC低, `11`:充电高温, `12`:充电低温, `13`:放电高温, `14`:放电低温, `15`:短路, `16`:MOS高温)
	FaultsOverview []*pb.BatteryFaultOverview `json:"faultsOverview"`
}

type BatteryStatistics struct {
	DateHour    []string  `json:"dateHour,omitempty"`    // 整点
	Voltage     []float64 `json:"voltage,omitempty"`     // 电压(V)
	BatTemp     []float64 `json:"batTemp,omitempty"`     // 电池温度(℃)
	MosTemp     []float64 `json:"mosTemp,omitempty"`     // Mos温度(℃)
	EnvTemp     []float64 `json:"envTemp,omitempty"`     // 环境温度(℃)
	Soc         []float64 `json:"soc,omitempty"`         // 剩余容量(%)
	Strength    []float64 `json:"strength,omitempty"`    // 信号强度(%)
	Charging    []float64 `json:"charging,omitempty"`    // 充电电流(A)
	DisCharging []float64 `json:"disCharging,omitempty"` // 放电电流(A)
}

type BatteryPositionReq struct {
	SN    string `json:"sn" param:"sn" validate:"required"` // 电池编码
	Start string `json:"start" query:"start"`               // 开始时间 (精确到秒, 格式为: yyyy-mm-dd hh:mm:ss, 默认6小时前)
	End   string `json:"end" query:"end"`                   // 结束时间 (精确到秒, 格式为: yyyy-mm-dd hh:mm:ss, 默认当前时间)
}

type BatteryPositionRes struct {
	Start      string               `json:"start,omitempty"`      // 开始时间
	End        string               `json:"end,omitempty"`        // 结束时间
	Stationary []*BatteryStationary `json:"stationary,omitempty"` // 静止位置信息
	Positions  []*BatteryPosition   `json:"positions,omitempty"`  // 所有位置信息
}

type BatteryStationary struct {
	InCabinet bool    `json:"inCabinet,omitempty"` // 是否在电柜中
	Duration  int64   `json:"duration,omitempty"`  // 停留时间 (秒)
	StartAt   string  `json:"startAt,omitempty"`   // 开始时间 (格式为: yyyy-mm-dd hh:mm:ss)
	EndAt     string  `json:"endAt,omitempty"`     // 结束时间 (格式为: yyyy-mm-dd hh:mm:ss)
	StartSoc  uint32  `json:"startSoc,omitempty"`  // 开始电量
	EndSoc    *uint32 `json:"endSoc,omitempty"`    // 结束电量, 可能为空
	Lng       float64 `json:"lng,omitempty"`       // 经度
	Lat       float64 `json:"lat,omitempty"`       // 纬度
}

type BatteryPosition struct {
	InCabinet  bool    `json:"inCabinet,omitempty"`  // 是否在电柜中
	Stationary bool    `json:"stationary,omitempty"` // 是否停留
	Soc        uint32  `json:"soc,omitempty"`        // 电量
	Lng        float64 `json:"lng,omitempty"`        // 经度
	Lat        float64 `json:"lat,omitempty"`        // 纬度
	Voltage    float64 `json:"voltage,omitempty"`    // 电压
	Gsm        uint32  `json:"gsm,omitempty"`        // GSM信号强度
	At         string  `json:"at,omitempty"`         // 时间 (格式为: yyyy-mm-dd hh:mm:ss)
}

type BatteryFaultReq struct {
	PaginationReq
	SN    *string              `json:"sn" query:"sn"`       // 电池编号
	Start string               `json:"start" query:"start"` // 开始日期 (格式为: yyyy-mm-dd)
	End   string               `json:"end" query:"end"`     // 结束日期 (格式为: yyyy-mm-dd)
	Fault *pb.BatteryFaultType `json:"fault" query:"fault"` // 选择故障, `0`:总压低, `1`:总压高, `2`:单体低, `3`:单体高, `6`:放电过流, `7`:充电过流, `8`:SOC低, `11`:充电高温, `12`:充电低温, `13`:放电高温, `14`:放电低温, `15`:短路, `16`:MOS高温
}

type BatteryFaultRes struct {
	Sn      string              `json:"sn,omitempty"`
	Fault   pb.BatteryFaultType `json:"fault,omitempty"`
	BeginAt string              `json:"begin_at,omitempty"`
	EndAt   string              `json:"end_at,omitempty"`
}

// BatteryTrackReq 电池轨迹
type BatteryTrackReq struct {
	Points []BatteryPoint `json:"points"` // 轨迹点
}

// BatteryPoint 电池轨迹点
type BatteryPoint struct {
	Latitude  float64 `json:"latitude" validate:"required"`  // 纬度
	Longitude float64 `json:"longitude" validate:"required"` // 经度
	LocTime   int32   `json:"locTime" validate:"required"`   // 时间UNIX时间戳（秒级）
}

type BatteryTrackRes struct {
	Points []BatteryPoint `json:"points"` // 轨迹点
}
