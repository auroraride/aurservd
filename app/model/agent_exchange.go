// Created at 2023-06-12

package model

type AgentExchangeListReq struct {
	PaginationReq
	Start     *string `json:"start" query:"start"`         // 筛选开始日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
	End       *string `json:"end" query:"end"`             // 筛选结束日期, 格式为yyyy-mm-dd, 例如: 2022-06-01
	Keyword   *string `json:"keyword" query:"keyword"`     // 筛选骑手姓名或电话
	CabinetID uint64  `json:"cabinetId" query:"cabinetId"` // 选择电柜ID
	Model     string  `json:"model" query:"model"`         // 电池型号
}
