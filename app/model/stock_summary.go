package model

// StockSummaryReq 物资统计数据入参
type StockSummaryReq struct {
	CabinetID    *uint64 `json:"cabinetID"`    // 电柜ID
	StoreID      *uint64 `json:"storeID"`      // 门店ID
	RiderID      *uint64 `json:"riderID"`      // 骑手ID
	BatteryID    *uint64 `json:"batteryID"`    // 电池ID
	EbikeID      *uint64 `json:"ebikeID"`      // 电车ID
	EnterpriseID *uint64 `json:"enterpriseID"` // 团签ID
	StationID    *uint64 `json:"stationID"`    // 站点ID
	Num          int     `json:"num"`          // 数量
	Date         string  `json:"date"`         // 日期
}
