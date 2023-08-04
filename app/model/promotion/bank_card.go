package promotion

const BankCardLimit = 5 // 银行卡上限数量

// BankCardReq 设置银行卡
type BankCardReq struct {
	CardNo string `json:"cardNo" validate:"required" param:"cardNo"` // 银行卡号
}

// BankCardRes 银行卡返回信息
type BankCardRes struct {
	ID          uint64 `json:"id"`          // id
	CardNo      string `json:"cardNo"`      // 银行卡号
	IsDefault   bool   `json:"isDefault"`   // 是否默认
	BankLogoURL string `json:"bankLogoUrl"` // 银行logo
	Bank        string `json:"bank"`        // 银行名称
}
