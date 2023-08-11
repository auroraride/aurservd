package promotion

type StatisticsReq struct {
	Start *string `json:"start"  query:"start"` // 开始日期
	End   *string `json:"end"  query:"end"`     // 结束日期
}

// StatisticsEarningsRes  我的收益
type StatisticsEarningsRes struct {
	StatisticsEarningsDetail
}

// StatisticsEarningsDetailRes 收益明细
type StatisticsEarningsDetailRes struct {
	FirstLevelEarnings  StatisticsEarningsDetail `json:"firstLevelEarnings"`  // 一级收益
	SecondLevelEarnings StatisticsEarningsDetail `json:"secondLevelEarnings"` // 二级收益
}

// StatisticsEarningsDetail 收益明细
type StatisticsEarningsDetail struct {
	TotalEarnings        float64 `json:"totalEarnings"`        // 佣金总收益
	TotalNewSignEarnings float64 `json:"totalNewSignEarnings"` // 新签佣金总收益
	TotalRenewalEarnings float64 `json:"totalRenewalEarnings"` // 续费佣金总收益
}

// StatisticsTeamRes 团队
type StatisticsTeamRes struct {
	StatisticsTeamDetail
}

// StatisticsTeamDetailRes 我的团队统计明细返回
type StatisticsTeamDetailRes struct {
	StatisticsTeamDetail
}

// StatisticsTeamDetail 我的团队明细
type StatisticsTeamDetail struct {
	Level        uint64 `json:"level,omitempty"` // 级别
	TotalTeam    uint64 `json:"totalTeam"`       // 累计团队人数
	TotalNewSign uint64 `json:"totalNewSign"`    // 累计新签人数
	TotalRenewal uint64 `json:"totalRenewal"`    // 累计续费人数
}

// StatisticsRecordRes 我的战绩
type StatisticsRecordRes struct {
	TotalInvitation uint64  `json:"totalInvitation"` // 累计邀请人数
	TotalEarnings   float64 `json:"totalEarnings"`   // 累计收益
}

// StatisticsWalletRes 钱包
type StatisticsWalletRes struct {
	TotalBalance      float64 `json:"totalBalance"`      // 总余额 = 可提现余额 + 冻结余额
	Balance           float64 `json:"balance"`           // 可提现余额
	FreezeBalance     float64 `json:"freezeBalance"`     // 冻结余额
	TotalEarnings     float64 `json:"totalEarnings"`     // 累计收益
	TodayEarnings     float64 `json:"todayEarnings"`     // 今日收益
	YesterdayEarnings float64 `json:"yesterdayEarnings"` // 昨日收益
	StatisticsWithdrawalRes
}

// StatisticsWithdrawalRes 提现统计
type StatisticsWithdrawalRes struct {
	// 累计提现
	TotalWithdrawal float64 `json:"totalWithdrawal"`
}

// StatisticsTeamGrowthTrendRes 团队增长趋势统计
type StatisticsTeamGrowthTrendRes struct {
	// 月份
	Month string `json:"month"`
	// 一级人数
	FirstLevelNum uint64 `json:"firstLevelNum"`
	// 二级人数
	SecondLevelNum uint64 `json:"secondLevelNum"`
}
