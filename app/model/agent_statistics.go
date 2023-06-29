// Created at 2023-06-12

package model

// AgentStatisticsOverviewRes 代理商小程序首页数据
type AgentStatisticsOverviewRes struct {
	AgentStatisticsOverviewAmount `json:"amount"`
	AgentStatisticsOverviewRider  `json:"rider"`
	BatterySummary                `json:"batterySummary"`
	EbikeSummary                  `json:"ebikeSummary"`
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
	RiderTotal          int `json:"riderTotal"`          // 骑手总数
	BillingRiderTotal   int `json:"billingRiderTotal"`   // 计费骑手数
	ExpiringRiderTotal  int `json:"expiringRiderTotal"`  // 临期骑手数
	SubscribeAlterTotal int `json:"subscribeAlterTotal"` // 加时待审核数
	ActivationTotal     int `json:"activationTotal"`     // 累计激活
	UnSubscribeTotal    int `json:"unSubscribeTotal"`    // 累计退租

	RiderOnlyBatteryTotal     int `json:"riderOnlyBatteryTotal"`     // 单电骑手
	RiderBatteryAndEbikeTotal int `json:"riderBatteryAndEbikeTotal"` // 车+电骑手
}
