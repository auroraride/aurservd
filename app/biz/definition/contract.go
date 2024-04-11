package definition

// ContractSignNewReq 合同签署请求
type ContractSignNewReq struct {
	SubscribeID uint64 `json:"subscribeId" validate:"required"` // 订阅ID
	Seal        string `json:"seal" validate:"required"`        // 签名Base64
	DocId       string `json:"docId" validate:"required"`       // 合同ID
}

type ContractCreateRes struct {
	Link      string `json:"link"`      // 合同链接
	DocId     string `json:"docId"`     // 合同ID
	Effective bool   `json:"effective"` // 是否存在生效中的合同, 若返回值为true则代表无需签合同
}

type ContractCreateReq struct {
	SubscribeID uint64 `json:"subscribeId" validate:"required"` // 订阅ID
}
