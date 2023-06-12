// Created at 2023-06-12

package model

// AgentStatisticsOverviewRes 代理商小程序首页数据
type AgentStatisticsOverviewRes struct {
	AgentStatisticsOverviewAmount
	AgentStatisticsOverviewRider
	AgentStatisticsOverviewBattery
}

// AgentStatisticsOverviewAmount 团签金额汇总
type AgentStatisticsOverviewAmount struct {
	// 余额
	Balance float64 `json:"balance"`
	// 昨日扣费
	YesterdayDeduction float64 `json:"yesterdayDeduction"`
	// 日均扣费
	AverageDeduction float64 `json:"averageDeduction"`
}

// AgentStatisticsOverviewRider 团签骑手数据汇总
type AgentStatisticsOverviewRider struct {
	RiderTotal         int `json:"riderTotal"`         // 骑手总数
	NewRiderTotal      int `json:"newRiderTotal"`      // 新签骑手数
	QuitRiderTotal     int `json:"quitRiderTotal"`     // 退租骑手数
	BillingRiderTotal  int `json:"billingRiderTotal"`  // 计费骑手数
	ExpiringRiderTotal int `json:"expiringRiderTotal"` // 临期骑手数
	OverTimeRiderTotal int `json:"overTimeRiderTotal"` // 骑手加时审核数
}

// AgentStatisticsOverviewBattery 团签电池汇总
type AgentStatisticsOverviewBattery struct {
	// 电池总数
	BatteryTotal int `json:"batteryTotal"`
	// 站点电池数
	StationBatteryTotal int `json:"stationBatteryTotal"`
	// 电柜电池数
	CabinetBatteryTotal int `json:"cabinetBatteryTotal"`
	// 骑手电池数
	RiderBatteryTotal int `json:"riderBatteryTotal"`
}
