package promotion

// SettingReq 推广设置请求
type SettingReq struct {
	Key SettingKey `json:"key"  validate:"required" param:"key"` // 键 GROWTH_VALUE_TEXT 成长值说明  WITHDRAWAL_TEXT 提现说明
}

// Setting 推广设置
type Setting struct {
	Key     SettingKey `json:"key"  validate:"required" param:"key"` // 键 GROWTH_VALUE_TEXT 成长值说明  WITHDRAWAL_TEXT 提现说明
	Title   string     `json:"title" `                               // 标题
	Context string     `json:"context" `                             // 内容
}

type SettingKey string

func (s SettingKey) Value() string {
	return string(s)
}

const (
	SettingGrowthValueText SettingKey = "GROWTH_VALUE_TEXT" // 成长值说明
	SettingWithdrawalText  SettingKey = "WITHDRAWAL_TEXT"   // 提现说明
)

var Settings = map[SettingKey]Setting{
	SettingGrowthValueText: {
		Key:     SettingKey(SettingGrowthValueText.Value()),
		Title:   "成长值说明",
		Context: "",
	},
	SettingWithdrawalText: {
		Key:     SettingKey(SettingWithdrawalText.Value()),
		Title:   "提现说明",
		Context: "",
	},
}
