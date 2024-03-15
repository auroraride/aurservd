package definition

import "github.com/auroraride/aurservd/app/model"

type FaultType uint8

const (
	FaultTypeCabinet FaultType = iota + 1 // 电柜故障
	FaultTypeBattery                      // 电池故障
	FaultTypeEbike                        // 车辆故障
	FaultTypeOther                        // 其他
)

func (f FaultType) Value() uint8 {
	return uint8(f)
}

func (f FaultType) String() string {
	switch f {
	case FaultTypeCabinet:
		return "电柜故障"
	case FaultTypeBattery:
		return "电池故障"
	case FaultTypeEbike:
		return "车辆故障"
	case FaultTypeOther:
		return "其他"
	}
	return ""
}

type Fault struct {
	ID          uint64      `json:"id"`          // 故障ID
	City        string      `json:"city"`        // 城市
	Type        FaultType   `json:"type"`        // 故障类型 1:电柜故障 2:电池故障 3:车辆故障 4:其他
	DeviceNo    string      `json:"deviceNo"`    // 设备编号
	FaultCause  []string    `json:"faultCause"`  // 故障原因
	Status      uint8       `json:"status"`      // 故障状态 0未处理 1已处理
	Description string      `json:"description"` // 故障描述
	Attachments []string    `json:"attachments"` // 附件
	Rider       model.Rider `json:"rider"`       // 骑手
	CreatedAt   string      `json:"createdAt"`   // 创建时间
	Remark      string      `json:"remark"`      // 备注
}

// FaultListReq 故障列表请求
type FaultListReq struct {
	model.PaginationReq
	City     *string `json:"city" query:"city"`         // 城市
	DeviceNo *string `json:"deviceNo" query:"deviceNo"` // 设备编号
	Status   *uint8  `json:"status" query:"status"`     // 故障状态 0未处理 1已处理
	Type     *uint8  `json:"type" query:"type"`         // 故障类型 1:电柜故障 2:电池故障 3:车辆故障 4:其他
	Start    *string `json:"start" query:"start"`       // 开始时间
	End      *string `json:"end" query:"end"`           // 结束时间
	Keyword  *string `json:"keyword" query:"keyword"`   // 关键字 用户姓名、手机号
}

// FaultCreateReq 创建
type FaultCreateReq struct {
	DeviceNo    string    `json:"deviceNo" validate:"required"`                 // 设备编号
	FaultCause  []string  `json:"faultCause" validate:"required" trans:"故障"`    // 故障内容
	Description string    `json:"description" validate:"required" trans:"故障描述"` // 故障描述
	Type        FaultType `json:"type" validate:"required" trans:"故障类型"`        // 故障类型 1:电柜故障 2:电池故障 3:车辆故障 4:其他
	Attachments []string  `json:"attachments" validate:"max=3"`                 // 附件
	CityID      uint64    `json:"cityId" validate:"required"`                   // 城市ID
}

// FaultModifyStatusReq 修改状态
type FaultModifyStatusReq struct {
	model.IDParamReq
	Status uint8  `json:"status"` // 故障状态 0未处理 1已处理
	Remark string `json:"remark"` // 备注
}

// FaultCauseReq 故障原因请求
type FaultCauseReq struct {
	Key string `json:"key" param:"key"` // 键 EBIKE_FAULT: 车辆故障 BATTERY_FAULT: 电池故障 OTHER_FAULT: 其他故障 CABINET_FAULT: 电柜故障
}
