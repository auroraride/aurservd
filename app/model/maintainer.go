package model

type MaintainerCabinetOperate string

const (
	MaintainerCabinetOperateOpenAll           MaintainerCabinetOperate = "open-all"           // 全部开仓
	MaintainerCabinetOperateMaintenance       MaintainerCabinetOperate = "maintenance"        // 维护
	MaintainerCabinetOperateMaintenanceCancel MaintainerCabinetOperate = "maintenance-cancel" // 取消维护
	MaintainerCabinetOperateInterrupt         MaintainerCabinetOperate = "interrupt"          // 中断维护
	MaintainerCabinetOperateReboot            MaintainerCabinetOperate = "reboot"             // 重启
)

// NeedMaintenance 是否需要维护
func (o MaintainerCabinetOperate) NeedMaintenance() bool {
	return o != MaintainerCabinetOperateMaintenance
}

type MaintainerBinOperate string

const (
	MaintainerBinOperateOpen    MaintainerBinOperate = "open"    // 开仓
	MaintainerBinOperateLock    MaintainerBinOperate = "lock"    // 锁定
	MaintainerBinOperateUnlock  MaintainerBinOperate = "unlock"  // 解锁
	MaintainerBinOperateDisable MaintainerBinOperate = "disable" // 禁用
	MaintainerBinOperateEnable  MaintainerBinOperate = "enable"  // 启用
)

type Maintainer struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`   // 姓名
	Enable bool   `json:"enable"` // 是否启用
	Phone  string `json:"phone"`  // 电话
	Cities []City `json:"cities"` // 城市列表
}

type MaintainerListReq struct {
	PaginationReq
	CityID  uint64 `json:"cityId" query:"cityId"`   // 城市
	Keyword string `json:"keyword" query:"keyword"` // 关键词
}

type MaintainerCreateReq struct {
	Name     string   `json:"name" validate:"required" trans:"姓名"`
	Phone    string   `json:"phone" validate:"required" trans:"电话"`
	Password string   `json:"password" validate:"required" trans:"密码"`
	CityIDs  []uint64 `json:"cityIds" validate:"required" trans:"城市列表"`
	Enable   *bool    `json:"enable"` // 是否启用，默认`true`
}

type MaintainerModifyReq struct {
	ID       uint64   `json:"id" param:"id" validate:"required"`
	Name     *string  `json:"name"`     // 姓名
	Phone    *string  `json:"phone"`    // 电话
	Password *string  `json:"password"` // 密码
	CityIDs  []uint64 `json:"cityIDs"`  // 城市列表
	Enable   *bool    `json:"enable"`   // 是否启用
}

type MaintainerSigninReq struct {
	Phone    string `json:"phone" validate:"required" trans:"电话"`
	Password string `json:"password" validate:"required" trans:"密码"`
}

type MaintainerSigninRes struct {
	*Maintainer
	Token string `json:"token"`
}

type MaintainerCabinetDetailReq struct {
	Serial string `json:"serial" param:"serial" validate:"required" trans:"电柜编码"`
}

type MaintainerCabinetDetailRes struct {
	*CabinetDetailRes
	Branch      Branch              `json:"branch"`      // 所属网点
	Maintenance AssetMaintenanceRes `json:"maintenance"` // 电柜维保信息
}

type MaintainerCabinetOperateReq struct {
	Serial  string                         `json:"serial" param:"serial" validate:"required" trans:"电柜编号"`
	Operate MaintainerCabinetOperate       `json:"operate" validate:"required" trans:"操作"` // 除`maintenance`之外所有操作均需要提前设电柜为维护状态。 open:开仓, open-all:全部开仓（暂时不做）, maintenance:电柜维护, maintenance-cancel:取消维护, interrupt:中断业务, reboot:重启（暂时不做）
	Lng     float64                        `json:"lng" validate:"required" trans:"经度"`
	Lat     float64                        `json:"lat" validate:"required" trans:"纬度"`
	Reason  string                         `json:"reason"`  // 操作原因（中断业务、维保失败 必填）
	Content string                         `json:"content"` // 维保内容（取消维护 必填）
	Details []AssetMaintenanceCreateDetail `json:"details"` // 维保使用配件
	Status  AssetMaintenanceStatus         `json:"status"`  // 维修状态 1:维修中 2:已维修 3:维修失败 4:已取消 5:暂停维护
	Mini    bool                           `json:"mini"`    // 是否小程序门店、代理调用电柜
}

type MaintainerBinOperateReq struct {
	Serial  string               `json:"serial" param:"serial" validate:"required" trans:"电柜编号"`
	Ordinal int                  `json:"ordinal" param:"ordinal" validate:"required" trans:"仓位序号"`
	Operate MaintainerBinOperate `json:"operate" validate:"required" trans:"操作"` // 所有操作均需要提前设电柜为维护状态。open:开仓, lock:锁定, unlock:解锁, disable:禁用, enable:启用
	Lng     float64              `json:"lng" validate:"required" trans:"经度"`
	Lat     float64              `json:"lat" validate:"required" trans:"纬度"`
	Reason  string               `json:"reason" validate:"required" trans:"操作原因"`
}

type MaintainerCabinetPauseReq struct {
	Serial string                 `json:"serial" param:"serial" validate:"required" trans:"电柜编号"`
	Status AssetMaintenanceStatus `json:"status" validate:"required" trans:"维修状态"` // 维修状态 1:维修中 2:已维修 3:维修失败 4:已取消 5:暂停维护
	Lng    float64                `json:"lng" validate:"required" trans:"经度"`
	Lat    float64                `json:"lat" validate:"required" trans:"纬度"`
}

type MaintainerCabinetListReq struct {
	PaginationReq
	Status  *CabinetStatus `json:"status" query:"status"`   // 电柜状态 0未投放 1运营中 2维护中
	ModelID *uint64        `json:"modelID" query:"modelID"` // 电池型号ID
	Keyword *string        `json:"keyword" query:"keyword"` // 关键词
}
