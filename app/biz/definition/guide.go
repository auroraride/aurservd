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
	Name   string `json:"name" validate:"required" trans:"名称"`
	Sort   uint8  `json:"sort" validate:"required" trans:"排序"`
	Answer string `json:"answer" validate:"required" trans:"答案"`
	Remark string `json:"remark,omitempty" trans:"备注"`
}

// GuideModifyReq 引导修改请求
type GuideModifyReq struct {
	model.IDParamReq
	Name   string `json:"name" validate:"required" trans:"名称"`
	Sort   uint8  `json:"sort" validate:"required" trans:"排序"`
	Answer string `json:"answer" validate:"required" trans:"答案"`
	Remark string `json:"remark,omitempty" trans:"备注"`
}

// GuideListReq 引导列表请求
type GuideListReq struct {
	model.PaginationReq
}
