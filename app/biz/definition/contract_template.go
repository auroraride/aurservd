package definition

import "github.com/auroraride/aurservd/app/model"

// ContractTemplateFields 模版公共字段
var ContractTemplateFields = []string{"name", "idcard", "address", "phone", "city", "riderSign", "riderContact", "aurDate", "riderDate", "payMonth"}

// FieldsUserMap 用户字段映射
var FieldsUserMap = map[uint8][]string{
	// 个签
	1: {"payRider"},
	// 团签
	2: {"payEnt", "entName", "entContact", "entStationentStation", "payerEnt"},
}

// FieldsSubMap 订阅字段映射
var FieldsSubMap = map[uint8][]string{
	// 单电
	1: {"schemaEbike", "ebikeScheme1", "ebikeBrand", "ebikeColor", "ebikeSN", "ebikePlate", "ebikeBattery", "ebikeModel", "ebikeScheme1Start", "ebikeScheme1Stop", "ebikeScheme1Price", "ebikeScheme1PayMonth", "ebikeScheme1PayTotal"},
	// 车电
	2: {"schemaBattery", "batteryModel", "batteryStart", "batteryStop", "batteryPrice", "batteryPayMonth", "batteryPayTotal"},
}

type ContractTemplate struct {
	ID        uint64 `json:"id"`        // ID
	Url       string `json:"url"`       // 合同链接
	Name      string `json:"name"`      // 合同名称
	UserType  uint8  `json:"userType"`  // 用户类型 1:个签 2:团签
	SubType   uint8  `json:"subType"`   // 订阅类型 1:单电 2:车电
	Enable    bool   `json:"enable"`    // 是否启用
	Sn        string `json:"sn"`        // 合同编号
	Remark    string `json:"remark"`    // 备注
	CreatedAt string `json:"createdAt"` // 创建时间
}

// ContractTemplateCreateReq 创建合同模板请求
type ContractTemplateCreateReq struct {
	Url      string  `json:"url" validate:"required"`                   // 合同链接
	Name     string  `json:"name" validate:"required" `                 // 合同名称
	UserType uint8   `json:"userType" validate:"required" enums:"1,2" ` // 用户类型 1:个签 2:团签
	SubType  uint8   `json:"subType" validate:"required" enums:"1,2"`   // 订阅类型 1:单电 2:车电
	Remark   *string `json:"remark"`                                    // 备注
	Enable   *bool   `json:"enable"`                                    // 是否启用
}

// ContractTemplateListRes 合同模板列表返回
type ContractTemplateListRes struct {
	ContractTemplate
}

// ContractTemplateModifyReq 修改合同模板请求
type ContractTemplateModifyReq struct {
	model.IDParamReq
	Name     string  `json:"name" validate:"required" `                 // 合同名称
	UserType uint8   `json:"userType" validate:"required" enums:"1,2" ` // 用户类型 1:个签 2:团签
	SubType  uint8   `json:"subType" validate:"required" enums:"1,2"`   // 订阅类型 1:单电 2:车电
	Remark   *string `json:"remark"`                                    // 备注
	Enable   *bool   `json:"enable"`                                    // 是否启用
}
