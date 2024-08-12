package model

type TransferAssetDetail struct {
	Name      string `json:"name"`                // 物资名称
	Outbound  int    `json:"outbound"`            // 出库数量
	Inbound   int    `json:"inbound"`             // 入库数量
	Surplus   int    `json:"surplus"`             // 剩余
	Exception int    `json:"exception,omitempty"` // 异常数量(电柜无)
}
