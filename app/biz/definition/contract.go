package definition

// ContractSignNewReq 合同签署请求
type ContractSignNewReq struct {
	SubscribeID uint64 `json:"subscribeId" validate:"required"` // 订阅ID
	Seal        string `json:"seal" validate:"required"`        // 签名Base64
	DocId       string `json:"docId" validate:"required"`       // 合同ID
}

// ContractSignNewRes 合同签署响应
type ContractSignNewRes struct {
	Link string `json:"link"` // 合同链接
}

// ContractCreateRes 合同创建响应
type ContractCreateRes struct {
	Link         string `json:"link,omitempty"`  // 合同链接
	DocId        string `json:"docId,omitempty"` // 合同ID
	Effective    bool   `json:"effective"`       // 是否存在生效中的合同, 若返回值为true则代表无需签合同
	NeedRealName bool   `json:"needRealName"`    // 是否需要实名认证   true:需要  false:不需要
}

// ContractCreateReq 合同创建请求
type ContractCreateReq struct {
	SubscribeID uint64 `json:"subscribeId" validate:"required"` // 订阅ID
}

// ContractCreateRPCReq 合同创建RPC请求
type ContractCreateRPCReq struct {
	TemplateId string // 合同模板ID
	Addr       string // rpc地址
	ExpiresAt  int64  // 过期时间
	Idcard     string // 身份证
}
