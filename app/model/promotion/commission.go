package promotion

type CommissionType uint8

const (
	CommissionDefault CommissionType = iota // 默认全局返佣方案
	CommissionCommon                        // 通用返佣方案
	CommissionCustom                        // 为个人自定义返佣方案
)

func (c CommissionType) Value() uint8 {
	v := uint8(c)
	return v
}

func (c CommissionType) String() string {
	switch c {
	case CommissionDefault:
		return "推广返佣(全局)"
	case CommissionCommon:
		return "通用返佣"
	case CommissionCustom:
		return "自定义"
	}
	return ""
}

// CommissionLimitedType 返佣次数类型
type CommissionLimitedType uint8

const (
	CommissionLimited   CommissionLimitedType = iota + 1 // 有限返佣
	CommissionUnlimited                                  // 无限返佣
)

type CommissionRuleKey string

const (
	FirstLevelNewSubscribeKey      CommissionRuleKey = "firstLevelNewSubscribe"      // 一级团员新签
	FirstLevelRenewalSubscribeKey  CommissionRuleKey = "firstLevelRenewalSubscribe"  // 一级团员续费
	SecondLevelNewSubscribeKey     CommissionRuleKey = "secondLevelNewSubscribe"     // 二级团员新签
	SecondLevelRenewalSubscribeKey CommissionRuleKey = "secondLevelRenewalSubscribe" // 二级团员续费
)

// MaxCommission 返佣最高金额
type MaxCommission struct {
	FirstLevelNewSubscribe      float64 `json:"firstLevelNewSubscribe"`      // 一级团员新签
	FirstLevelRenewalSubscribe  float64 `json:"firstLevelRenewalSubscribe"`  // 一级团员续费
	SecondLevelNewSubscribe     float64 `json:"secondLevelNewSubscribe"`     // 二级团员新签`
	SecondLevelRenewalSubscribe float64 `json:"secondLevelRenewalSubscribe"` // 二级团员续费
}

var CommissionRuleKeyNames = map[CommissionRuleKey]string{
	FirstLevelNewSubscribeKey:      "一级团员新签",
	SecondLevelNewSubscribeKey:     "二级团员新签",
	FirstLevelRenewalSubscribeKey:  "一级团员续费",
	SecondLevelRenewalSubscribeKey: "二级团员续费",
}

func (k CommissionRuleKey) String() string {
	if name, ok := CommissionRuleKeyNames[k]; ok {
		return name
	}
	return ""
}

func (k CommissionRuleKey) Value() string {
	return string(k)
}

// CommissionSelection 返佣方案筛选列表返回参数
type CommissionSelection struct {
	ID   uint64 `json:"id" `   // id
	Name string `json:"name" ` // 方案名称
}

// CommissionTaskSelect 返佣任务选择
type CommissionTaskSelect struct {
	Key  CommissionRuleKey `json:"key" `  // key
	Name string            `json:"name" ` // 名称
}

// CommissionCreateReq  创建返佣方案请求参数
type CommissionCreateReq struct {
	ID       uint64          `json:"id" param:"id"`                     // id
	Name     string          `json:"name" validate:"required"`          // 方案名称
	Rule     CommissionRule  `json:"rule" validate:"required"`          // 返佣规则
	Type     *CommissionType `json:"type" validate:"required"`          // 返佣类型 0:默认全局返佣方案 1:通用返佣方案 2:为个人自定义返佣方案
	MemberID *uint64         `json:"memberId"`                          // 会员id
	Desc     *string         `json:"desc"`                              // 返佣说明
	PlanID   []uint64        `json:"planId" validate:"required,unique"` // 骑士卡方案ID
}

// CommissionDetail 详情
type CommissionDetail struct {
	ID                   uint64            `json:"id" `                  // id
	Name                 string            `json:"name"`                 // 方案名称
	Rule                 CommissionRule    `json:"rule" `                // 返佣规则
	Type                 CommissionType    `json:"type" `                // 返佣类型 0:默认全局返佣方案 1:通用返佣方案 2:为个人自定义返佣方案
	MemberID             *uint64           `json:"memberId,omitempty"`   // 会员id
	Desc                 *string           `json:"desc"`                 // 返佣说明
	Enable               bool              `json:"enable" `              // 启用状态 false:禁用 true:启用
	AmountSum            float64           `json:"amountSum"`            // 佣金总额
	CreatedAt            string            `json:"createdAt"`            // 创建时间
	StartAt              string            `json:"startAt"`              // 开始时间
	EndAt                string            `json:"endAt"`                // 结束时间
	Plan                 []*CommissionPlan `json:"plan"`                 // 骑士卡方案
	FistNewNumSum        uint64            `json:"fistNewNumSum"`        // 一级团员新签人数
	FistNewAmountSum     float64           `json:"fistNewAmountSum"`     // 一级团员新签金额
	FistRenewNumSum      uint64            `json:"fistRenewNumSum"`      // 一级团员续费人数
	FistRenewAmountSum   float64           `json:"fistRenewAmountSum"`   // 一级团员续费金额
	SecondNewNumSum      uint64            `json:"secondNewNumSum"`      // 二级团员新签人数
	SecondNewAmountSum   float64           `json:"secondNewAmountSum"`   // 二级团员新签金额
	SecondRenewNumSum    uint64            `json:"secondRenewNumSum"`    // 二级团员续费人数
	SecondRenewAmountSum float64           `json:"secondRenewAmountSum"` // 二级团员续费金额

}

// CommissionDeleteReq 删除会员返佣方案请求参数
type CommissionDeleteReq struct {
	ID uint64 `json:"id" validate:"required" param:"id"` // id
}

// CommissionEnableReq 方案状态修改
type CommissionEnableReq struct {
	ID     uint64 `json:"id" validate:"required" param:"id"` // id
	Enable *bool  `json:"enable" validate:"required"`        // 启用状态 false:禁用 true:启用
}

// CommissionRule 会员返佣规则
type CommissionRule struct {
	NewUserCommission map[CommissionRuleKey]*CommissionRuleConfig `json:"newUserCommission"` // 新用户返佣规则
	RenewalCommission map[CommissionRuleKey]*CommissionRuleConfig `json:"renewalCommission"` // 续签返佣规则
}

// CommissionRuleConfig 返佣配置
type CommissionRuleConfig struct {
	Name        string                `json:"name"`        // 任务名称
	LimitedType CommissionLimitedType `json:"limitedType"` // 返佣次数类型 1:有限返佣 2:无限返佣
	OptionType  CommissionOptionType  `json:"optionType"`  // 选项类型 1:固定金额 2:比例
	Value       []float64             `json:"value"`       // 选项值
	Desc        string                `json:"desc"`        // 任务内容说明
}

// CommissionOptionType 表示选项类型的枚举
type CommissionOptionType int

const (
	FixedAmount CommissionOptionType = iota + 1 // 固定金额
	Percentage                                  // 比例
)

type CommissionCalculationType uint8

const (
	CommissionTypeNewlySigned CommissionCalculationType = iota + 1 // 新签
	CommissionTypeRenewal                                          // 续签
)

// CommissionCalculation 佣金计算
type CommissionCalculation struct {
	OrderID      uint64  `json:"orderID"`      // 订单id
	RiderID      uint64  `json:"riderID"`      // 骑手id
	ActualAmount float64 `json:"actualAmount"` // 实付金额
	PlanID       uint64  `json:"planID"`       // 骑士卡方案ID

	Type CommissionCalculationType `json:"type"` // 任务类型 1:新签 2:续签
}

// CommissionRuleRes 返佣规则返回参数
type CommissionRuleRes struct {
	Detail     []CommissionRuleDetail `json:"detail"`     // 返佣规则
	DetailDesc *string                `json:"detailDesc"` // 详细规则说明
}

// CommissionRuleDetail 返佣规则
type CommissionRuleDetail struct {
	Key    CommissionRuleKey `json:"key"`    // key
	Name   string            `json:"name"`   // 名称
	Ratio  float64           `json:"ratio"`  // 比例
	Amount uint64            `json:"amount"` // 金额
	Desc   string            `json:"desc"`   // 说明
}

type CommissionPlanSelectionReq struct {
	ID uint64 `json:"id" validate:"required" param:"id"` // id
}

type CommissionPlanListRes struct {
	CommissionID   uint64            `json:"commissionId"`   // 返佣方案ID
	CommissionName string            `json:"name"`           // 方案名称
	Type           CommissionType    `json:"type"`           // 返佣类型 0:默认通用返佣方案 1:通用返佣方案 2:为个人自定义返佣方案
	Plan           []*CommissionPlan `json:"plan,omitempty"` // 骑士卡方案
	CreatedAt      string            `json:"createdAt"`      // 创建时间
	Rule           *CommissionRule   `json:"rule"`           // 返佣规则
}

type CommissionPlan struct {
	ID     uint64  `json:"id"`     // 骑士卡ID
	Name   string  `json:"name"`   // 骑士卡名称
	Amount float64 `json:"amount"` // 骑士卡价格
	Enable bool    `json:"enable"` // 是否禁用 true 启用 false 禁用
	End    string  `json:"end"`    // 有效期结束日期
}

// CommissionMaxPlan 每个方案最大返佣金额
type CommissionMaxPlan struct {
	FirstNewAmount      float64 `json:"firstNewAmount"`      // 一级团员新签
	FirstRenewalAmount  float64 `json:"firstRenewalAmount"`  // 一级团员续费
	SecondNewAmount     float64 `json:"secondNewAmount"`     // 二级团员新签
	SecondRenewalAmount float64 `json:"secondRenewalAmount"` // 二级团员续费
}
