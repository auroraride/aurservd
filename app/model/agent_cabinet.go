// Created at 2023-05-29

package model

// AgentCabinetDetailReq 代理电柜详情请求
type AgentCabinetDetailReq struct {
	Serial string `json:"serial" param:"serial" validate:"required"` // 电柜编号
}

// AgentCabinetDetailRes 代理电柜详情
type AgentCabinetDetailRes struct {
	Serial  string             `json:"serial"`
	Name    string             `json:"name"`                 // 电柜名称
	Status  uint8              `json:"status" enums:"1,2"`   // 投放状态, 1:运营中 2:维护中
	Health  uint8              `json:"health" enums:"0,1,2"` // 健康状态, 0:离线 1:在线 2:故障
	Station string             `json:"station"`              // 所属站点
	Bins    []*AgentCabinetBin `json:"bins,omitempty"`       // 仓位信息
}

type AgentCabinetBin struct {
	Ordinal   int        `json:"ordinal"`             // 仓位序号, 从`1`开始
	Usable    bool       `json:"usable"`              // 是否可用 `false`时,禁止代理操作该仓位
	BatterySN string     `json:"batterySn,omitempty"` // 当前电池编码 (无电池的时候为空)
	Soc       BatterySoc `json:"soc,omitempty"`       // 当前电量 (无电池时为空)
}
