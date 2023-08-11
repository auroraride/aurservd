package promotion

// Level 会员等级
type Level struct {
	ID              uint64   `json:"id" param:"id"`                       // id
	Level           *uint64  `json:"level" `                              // 会员等级
	GrowthValue     *uint64  `json:"growthValue" validate:"required"`     // 成长值
	CommissionRatio *float64 `json:"commissionRatio" validate:"required"` // 佣金比例
}

// LevelSelection 会员等级筛选列表返回参数
type LevelSelection struct {
	ID    uint64 `json:"id" `    // id
	Level uint64 `json:"level" ` // 会员等级
}
