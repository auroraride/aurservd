package model

// GuideDetail 引导详情
type GuideDetail struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Sort      uint8  `json:"sort"`
	Answer    string `json:"answer"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// GuideSaveReq 引导保存请求
type GuideSaveReq struct {
	Name   string `json:"name" validate:"required"`
	Sort   uint8  `json:"sort" validate:"required"`
	Answer string `json:"answer" validate:"required"`
	Remark string `json:"remark"`
}

// GuideModifyReq 引导修改请求
type GuideModifyReq struct {
	IDParamReq
	Name   string `json:"name" validate:"required"`
	Sort   uint8  `json:"sort" validate:"required"`
	Answer string `json:"answer" validate:"required"`
	Remark string `json:"remark"`
}

// GuideListReq 引导列表请求
type GuideListReq struct {
	PaginationReq
}
