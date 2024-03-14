package definition

import "github.com/auroraride/aurservd/app/model"

// QuestionCommon 常见问题通用字段
type QuestionCommon struct {
	Name         string `json:"name" validate:"required"`
	Sort         int    `json:"sort" validate:"required"`
	CategoryID   uint64 `json:"category_id,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
	Answer       string `json:"answer" validate:"required"`
}

// QuestionDetail 问题详情
type QuestionDetail struct {
	model.IDRes
	QuestionCommon
}

// QuestionCreateReq 创建请求
type QuestionCreateReq struct {
	QuestionCommon
}

// QuestionModifyReq 修改请求
type QuestionModifyReq struct {
	model.IDParamReq
	QuestionCommon
}

// QuestionListReq 列表请求
type QuestionListReq struct {
	model.PaginationReq
	Keyword    *string `json:"keywords,omitempty"`
	CategoryID *uint64 `json:"category_id,omitempty" trans:"分类ID,0表示其他，不传表示全部"`
}
