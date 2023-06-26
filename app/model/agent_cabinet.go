// Created at 2023-05-29

package model

// AgentCabinetPermission 代理电柜操作权限
type AgentCabinetPermission uint

const (
	AgentCabinetPermissionView       AgentCabinetPermission = 1 << iota // 查看
	AgentCabinetPermissionMaintain                                      // 维护
	AgentCabinetPermissionReboot                                        // 重启
	AgentCabinetPermissionOpenBin                                       // 开仓（包含满电仓和空仓）
	AgentCabinetPermissionDisableBin                                    // 禁用仓位
	AgentCabinetPermissionEnableBin                                     // 启用仓位
)

var AgentCabinetPermissionAll = AgentCabinetPermissionView |
	AgentCabinetPermissionMaintain |
	AgentCabinetPermissionReboot |
	AgentCabinetPermissionOpenBin |
	AgentCabinetPermissionDisableBin |
	AgentCabinetPermissionEnableBin

// AgentCabinetDetailReq 代理电柜详情请求
type AgentCabinetDetailReq struct {
	Serial string   `json:"serial" param:"serial" validate:"required"` // 电柜编号
	Lng    *float64 `json:"lng" query:"lng"`                           // 经度
	Lat    *float64 `json:"lat" query:"lat"`                           // 纬度
}

// AgentCabinet 代理电柜详情
type AgentCabinet struct {
	Serial     string                 `json:"serial"`
	Name       string                 `json:"name"`                             // 电柜名称
	Status     uint8                  `json:"status" enums:"1,2"`               // 投放状态, 1:运营中 2:维护中
	Health     uint8                  `json:"health" enums:"0,1,2"`             // 健康状态, 0:离线 1:在线 2:故障
	Address    string                 `json:"address"`                          // 地址
	Lng        float64                `json:"lng"`                              // 电柜经度
	Lat        float64                `json:"lat"`                              // 电柜纬度
	Station    EnterpriseStation      `json:"station"`                          // 所属站点
	Bins       []*AgentCabinetBin     `json:"bins"`                             // 仓位信息
	Models     []string               `json:"models"`                           // 电池型号
	Permission AgentCabinetPermission `json:"permission" enums:"1,2,4,8,16,32"` // 权限，<a target="_blank" href="https://rqsz38umzux.feishu.cn/wiki/space/7248526754619539460?ccm_open_type=lark_wiki_spaceLink">查看WIKI</a>
	Distance   *float64               `json:"distance"`                         // 距离（仅携带`lng`和`lat`时有该字段）
}

type AgentCabinetBin struct {
	Ordinal   int        `json:"ordinal"`             // 仓位序号, 从`1`开始
	Usable    bool       `json:"usable"`              // 是否可用 `false`时,禁止代理操作该仓位
	BatterySN string     `json:"batterySn,omitempty"` // 当前电池编码 (无电池的时候为空)
	Soc       BatterySoc `json:"soc,omitempty"`       // 当前电量 (无电池时为空)
}

type AgentCabinetListReq struct {
	PaginationReq
	Lng       *float64 `json:"lng" query:"lng"`             // 经度
	Lat       *float64 `json:"lat" query:"lat"`             // 纬度
	StationID *uint64  `json:"stationId" query:"stationId"` // 站点
	Model     *string  `json:"model" query:"model"`         // 型号
	Serial    *string  `json:"serial" query:"serial"`       // 编号
}
