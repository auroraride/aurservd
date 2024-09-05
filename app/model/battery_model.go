package model

// BatteryModel 电池型号
type BatteryModel struct {
	ID    uint64 `json:"id,omitempty"`
	Model string `json:"model"` // 电池型号
}

// BatteryModelReq 电池型号创建请求
type BatteryModelReq struct {
	Model string `json:"model"` // 电池型号(POST), 例如60V30AH
}

type BatteryModelType uint8

const (
	BatteryModelTypeSmart    BatteryModelType = iota + 1 // 智能电池
	BatteryModelTypeNonSmart                             // 非智能电池
)

func (b BatteryModelType) String() string {
	switch b {
	case BatteryModelTypeSmart:
		return "智能电池"
	case BatteryModelTypeNonSmart:
		return "非智能电池"
	default:
		return "未知"
	}
}

func (b BatteryModelType) Value() uint8 {
	return uint8(b)
}

type SelectModelsReq struct {
	Type *BatteryModelType `json:"type" query:"type"` // 1智能电池 2非智能电池
}

type BatteryModelListReq struct {
	PaginationReq
	Type *BatteryModelType `json:"type" query:"type"` // 电池类型 1智能电池 2非智能电池
}

// BatteryModelDetail 电池型号信息
type BatteryModelDetail struct {
	ID       uint64           `json:"id"`       // 电池型号ID
	Model    string           `json:"model"`    // 电池型号
	Type     BatteryModelType `json:"type"`     // 电池类型
	Voltage  uint             `json:"voltage"`  // 电压
	Capacity uint             `json:"capacity"` // 容量
}

// BatteryModelCreateReq 创建
type BatteryModelCreateReq struct {
	Type     BatteryModelType `json:"type" validate:"required" enums:"1,2" trans:"电池类型"` // 电池类型
	Voltage  uint             `json:"voltage" validate:"required" trans:"电池电压"`          // 电池电压
	Capacity uint             `json:"capacity" validate:"required" trans:"电池容量"`         // 电池容量
}

// BatteryModelModifyReq 修改
type BatteryModelModifyReq struct {
	IDParamReq
	BatteryModelCreateReq
}
