package definition

import "github.com/auroraride/aurservd/app/model"

// QuestionCommon 常见问题通用字段
type QuestionCommon struct {
	Name         string `json:"name" validate:"required"`   // 问题名称
	Sort         int    `json:"sort" validate:"required"`   // 排序
	CategoryID   uint64 `json:"categoryId,omitempty"`       // 分类ID
	CategoryName string `json:"categoryName,omitempty"`     // 分类名称
	Answer       string `json:"answer" validate:"required"` // 答案
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
	Keyword    *string `json:"keyword" query:"keyword"`                                 // 关键字
	CategoryID *uint64 `json:"categoryId" query:"categoryId" trans:"分类ID,0表示其他，不传表示全部"` // 分类ID
}
