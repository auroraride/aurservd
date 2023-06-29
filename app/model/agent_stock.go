package model

type AgentStockDetailReq struct {
	PaginationReq
	AgentStockDetailFilter
}

type AgentStockDetailFilter struct {
	Materials    string `json:"materials" query:"materials"`         // 查询物资类别, 默认为电池, 逗号分隔 battery:电池 ebike:电车 others:其他物资
	Serial       string `json:"serial" query:"serial"`               // 电柜编号
	CityID       uint64 `json:"cityId" query:"cityId"`               // 城市ID
	CabinetID    uint64 `json:"cabinetId" query:"cabinetId"`         // 电柜ID
	Start        string `json:"start" query:"start"`                 // 开始时间
	End          string `json:"end" query:"end"`                     // 结束时间
	Positive     bool   `json:"positive" query:"positive"`           // 是否正序(默认倒序)
	Type         uint8  `json:"type" query:"type" enums:"0,1,2,3,4"` // 调拨类型, 0:调拨 1:激活 2:寄存 3:结束寄存 4:退租
	Keyword      string `json:"keyword" query:"keyword"`             // 查询电柜编号、车架号、电池编码
	EnterpriseID uint64 `json:"enterpriseId" query:"enterpriseId"`   // 企业ID (代理该参数不生效)
	StationID    uint64 `json:"stationId" query:"stationId"`         // 站点ID
	EbikeSN      string `json:"ebikeSn" query:"ebikeSn"`             // 车架号
	BatterySN    string `json:"batterySn" query:"batterySn"`         // 电池编码
	Model        string `json:"model" query:"model"`                 // 电池型号
}

type AgentStockDetailRes struct {
	ID       uint64       `json:"id"`              // 调拨ID
	Sn       string       `json:"sn"`              // 调拨编号
	City     string       `json:"city"`            // 城市
	Inbound  TransferInfo `json:"inbound"`         // 调入
	Outbound TransferInfo `json:"outbound"`        // 调出
	Name     string       `json:"name"`            // 物资
	Num      int          `json:"num"`             // 数量
	Type     string       `json:"type"`            // 类型
	Operator string       `json:"operator"`        // 操作人
	Time     string       `json:"time"`            // 时间
	Remark   string       `json:"remark"`          // 备注信息
	Rider    string       `json:"rider,omitempty"` // 骑手, 仅业务发生有此字段
}
type TransferInfo struct {
	PlatformName   string `json:"platformName"`   // 平台名称
	EnterpriseName string `json:"enterpriseName"` // 团签名称
	StationName    string `json:"stationName"`    // 站点名称
	CabinetName    string `json:"cabinetName"`    // 电柜名称
	CabinetSerial  string `json:"cabinetSerial"`  // 电柜编号
	RiderName      string `json:"riderName"`      // 骑手姓名
	RiderPhone     string `json:"riderPhone"`     // 骑手电话
}
