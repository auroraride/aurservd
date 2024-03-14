package definition

import "github.com/auroraride/aurservd/app/model"

// QuestionCategoryCommon 通用
type QuestionCategoryCommon struct {
	Name   string `json:"name" validate:"required"` // 分类名称
	Sort   uint64 `json:"sort" validate:"required"` // 排序
	Remark string `json:"remark,omitempty"`         //  备注
}

// QuestionCategoryDetail 详情
type QuestionCategoryDetail struct {
	model.IDRes
	QuestionCategoryCommon
}

// QuestionCategoryCreateReq 创建请求
type QuestionCategoryCreateReq struct {
	QuestionCategoryCommon
}

// QuestionCategoryModifyReq 修改请求
type QuestionCategoryModifyReq struct {
	model.IDParamReq
	QuestionCategoryCommon
}

// QuestionCategoryListReq 列表请求
type QuestionCategoryListReq struct {
	model.PaginationReq
}
