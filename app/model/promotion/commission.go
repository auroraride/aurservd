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
	ID       uint64          `json:"id" param:"id"`            // id
	Name     string          `json:"name" validate:"required"` // 方案名称
	Rule     CommissionRule  `json:"rule" validate:"required"` // 返佣规则
	Type     *CommissionType `json:"type" validate:"required"` // 返佣类型 0:默认全局返佣方案 1:通用返佣方案 2:为个人自定义返佣方案
	MemberID *uint64         `json:"memberId"`                 // 会员id
	Desc     *string         `json:"desc" validate:"required"` // 返佣说明
}

// CommissionDetail 详情
type CommissionDetail struct {
	ID        uint64         `json:"id" `                // id
	Name      string         `json:"name"`               // 方案名称
	Rule      CommissionRule `json:"rule" `              // 返佣规则
	Type      CommissionType `json:"type" `              // 返佣类型 0:默认全局返佣方案 1:通用返佣方案 2:为个人自定义返佣方案
	MemberID  *uint64        `json:"memberId,omitempty"` // 会员id
	Desc      *string        `json:"desc"`               // 返佣说明
	Enable    bool           `json:"enable" `            // 启用状态 false:禁用 true:启用
	UseCount  uint64         `json:"useCount"`           // 使用次数
	AmountSum float64        `json:"amountSum"`          // 佣金总额
	CreatedAt string         `json:"createdAt"`          // 创建时间
	StartAt   string         `json:"startAt"`            // 开始时间
	EndAt     string         `json:"endAt"`              // 结束时间
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
	Ratio       []float64             `json:"ratio"`       // 返佣比例列表
	Desc        string                `json:"desc"`        // 任务内容说明
}

type CommissionCalculationType uint8

const (
	CommissionTypeNewlySigned CommissionCalculationType = iota + 1 // 新签
	CommissionTypeRenewal                                          // 续签
)

// CommissionCalculation 佣金计算
type CommissionCalculation struct {
	RiderID        uint64                    `json:"riderID"`        // 骑手id
	CommissionBase float64                   `json:"commissionBase"` // 金额
	Type           CommissionCalculationType `json:"type"`           // 任务类型 1:新签 2:续签
}

// CommissionRuleRes 返佣规则返回参数
type CommissionRuleRes struct {
	Detail     []CommissionRuleDetail `json:"detail"`     // 返佣规则
	DetailDesc string                 `json:"detailDesc"` // 详细规则说明
}

// CommissionRuleDetail 返佣规则
type CommissionRuleDetail struct {
	Key   CommissionRuleKey `json:"key"`   // key
	Name  string            `json:"name"`  // 名称
	Ratio float64           `json:"ratio"` // 比例
	Desc  string            `json:"desc"`  // 说明
}