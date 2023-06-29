package model

import (
	"time"
)

// StockSummaryParams 物资统计数据入参
type StockSummaryParams struct {
	CabinetID    *uint64   `json:"cabinetID"`    // 电柜ID
	StoreID      *uint64   `json:"storeID"`      // 门店ID
	RiderID      *uint64   `json:"riderID"`      // 骑手ID
	EnterpriseID *uint64   `json:"enterpriseID"` // 团签ID
	StationID    *uint64   `json:"stationID"`    // 站点ID
	Model        string    `json:"model"`        // 电池型号/电车型号/其他物资类型
	Material     string    `json:"material"`     // 物资类型
	Num          int       `json:"num"`          // 物资出\入库数量
	Date         time.Time `json:"date"`         // 日期
	StockNum     int       `json:"stockNum"`     // 物资总数量
	Type         uint8     `json:"type"`         // 调拨类型
}

// BatterySummary 团签电池汇总
type BatterySummary struct {
	BatteryTotal        int `json:"batteryTotal"`        // 电池总数
	StationBatteryTotal int `json:"stationBatteryTotal"` // 站点电池数
	CabinetBatteryTotal int `json:"cabinetBatteryTotal"` // 电柜电池数
	RiderBatteryTotal   int `json:"riderBatteryTotal"`   // 骑手电池数
}

// EbikeSummary 电车汇总
type EbikeSummary struct {
	EbikeTotal int `json:"ebikeTotal"` // 电车总数
	// 站点电车数
	StationEbikeTotal int `json:"stationEbikeTotal"`
	// 骑手电车数
	RiderEbikeTotal int `json:"riderEbikeTotal"`
}

// StockSummaryReq 电池物资统计请求参数
type StockSummaryReq struct {
	EnterpriseID uint64 `json:"enterpriseID" query:"enterpriseId"` // 团签ID
}

type BatteryStockSummaryRsp struct {
	Overview BatterySummary       `json:"overview"` // 概览
	Group    []*BatteryStockGroup `json:"group"`    // 分组
}

// BatteryStockGroup 电池物资分组
type BatteryStockGroup struct {
	Model string `json:"model"`
	BatterySummary
}

type EbikeStockSummaryRsp struct {
	Overview EbikeSummary       `json:"overview"` // 概览
	Group    []*EbikeStockGroup `json:"group"`    // 分组
}

// EbikeStockGroup 电车物资分组统计
type EbikeStockGroup struct {
	Model string `json:"model"`
	EbikeSummary
}
