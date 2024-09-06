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

type BatteryModelListReq struct {
	PaginationReq
}

// BatteryModelDetail 电池型号信息
type BatteryModelDetail struct {
	ID       uint64 `json:"id"`       // 电池型号ID
	Model    string `json:"model"`    // 电池型号
	Voltage  uint   `json:"voltage"`  // 电压
	Capacity uint   `json:"capacity"` // 容量
}

// BatteryModelCreateReq 创建
type BatteryModelCreateReq struct {
	Voltage  uint `json:"voltage" validate:"required" trans:"电池电压"`  // 电池电压
	Capacity uint `json:"capacity" validate:"required" trans:"电池容量"` // 电池容量
}

// BatteryModelModifyReq 修改
type BatteryModelModifyReq struct {
	IDParamReq
	BatteryModelCreateReq
}
