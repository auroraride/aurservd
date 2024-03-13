package definition

import "github.com/auroraride/aurservd/app/model"

// QuestionCategoryCommon 通用
type QuestionCategoryCommon struct {
	Name   string `json:"name" validate:"required"`
	Sort   int    `json:"sort" validate:"required"`
	Remark string `json:"remark,omitempty"`
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
