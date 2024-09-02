package model

type TransferAssetDetail struct {
	ID            uint64 `json:"id"`                  // 物资ID
	Name          string `json:"name"`                // 物资名称
	Outbound      int    `json:"outbound"`            // 出库数量
	Inbound       int    `json:"inbound"`             // 入库数量
	Surplus       int    `json:"surplus"`             // 剩余
	Exception     int    `json:"exception,omitempty"` // 异常数量(电柜无)
	InTimeAt      string `json:"inTimeAt"`            // 入库时间
	InOperateName string `json:"inOperateName"`       // 入库人
}
