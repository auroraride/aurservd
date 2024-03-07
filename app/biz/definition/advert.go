package definition

import "github.com/auroraride/aurservd/app/model"

// AdvertDetail 广告详情
type AdvertDetail struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Sort   int    `json:"sort"`
	Image  string `json:"image"`
	Link   string `json:"link"`
	Remark string `json:"remark,omitempty"`
}

// AdvertReqCommon 广告请求公共字段
type AdvertReqCommon struct {
	Name   string `json:"name" validate:"required" trans:"名称"`
	Sort   int    `json:"sort" validate:"required" trans:"排序"`
	Image  string `json:"image" validate:"required" trans:"图片"`
	Link   string `json:"link" validate:"required" trans:"连接"`
	Remark string `json:"remark,omitempty"`
}

// AdvertSaveReq 广告保存请求
type AdvertSaveReq struct {
	AdvertReqCommon
}

// AdvertModifyReq 广告修改请求
type AdvertModifyReq struct {
	model.IDParamReq
	AdvertReqCommon
}

// AdvertListReq 广告列表请求
type AdvertListReq struct {
	model.PaginationReq
}
