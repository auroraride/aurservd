package model

type InstallmentDetail struct {
	Num    int       `json:"num"`    // 分期期数
	Prices []float64 `json:"prices"` // 每期价格
}
