package definition

import "github.com/auroraride/aurservd/app/model"

type GuideDetail struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Sort   uint8  `json:"sort"`
	Answer string `json:"answer"`
	Remark string `json:"remark,omitempty"`
}

// GuideSaveReq 引导保存请求
type GuideSaveReq struct {
	Name   string `json:"name" validate:"required"`
	Sort   uint8  `json:"sort" validate:"required"`
	Answer string `json:"answer" validate:"required"`
	Remark string `json:"remark,omitempty"`
}

// GuideModifyReq 引导修改请求
type GuideModifyReq struct {
	model.IDParamReq
	Name   string `json:"name" validate:"required"`
	Sort   uint8  `json:"sort" validate:"required"`
	Answer string `json:"answer" validate:"required"`
	Remark string `json:"remark,omitempty"`
}

// GuideListReq 引导列表请求
type GuideListReq struct {
	model.PaginationReq
}
