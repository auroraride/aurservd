package definition

type ContractSignReq struct {
	SubscribeID uint64 `json:"subscribeId" validate:"required"` // 订阅ID
	Seal        string `json:"seal" validate:"required"`        // 签名Base64
}

type ContractAddTemplateReq struct {
	SubscribeID uint64 `json:"subscribeId" validate:"required"` // 订阅ID
}

type ContractAddTemplateRes struct {
	Url string `json:"url"` // 合同链接
}

type CreateContractReq struct {
	Url      string `json:"url"`      // 合同链接
	Name     string `json:"name"`     // 合同名称
	UserType uint8  `json:"userType"` // 用户类型 1:个签 2:团签
	Remark   string `json:"remark"`   // 备注
}
