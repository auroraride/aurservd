package promotion

import "github.com/auroraride/aurservd/app/model"

// GrowthStatus 成长值状态
type GrowthStatus uint8

const (
	GrowthStatusValid   GrowthStatus = iota + 1 // 有效
	GrowthStatusInvalid                         // 无效
)

func (m GrowthStatus) Value() uint8 {
	return uint8(m)
}

type GrowthReq struct {
	model.PaginationReq
	GrowthFilter
	ID *uint64 `json:"id" param:"id"` // id
}

type GrowthFilter struct {
	Keyword     *string `json:"keyword" query:"keyword"`           // 关键词 手机号/姓名
	Status      *uint8  `json:"status" query:"status" enums:"1,2"` // 状态 1:有效 2:无效
	LevelTaskID *uint64 `json:"levelTaskId" query:"levelTaskId"`   // 等级任务id
	Start       *string `json:"start" query:"start"`               // 开始日期
	End         *string `json:"end" query:"end"`                   // 结束日期
}

// GrowthRes 成长值列表返回参数
type GrowthRes struct {
	GrowthDetail
	Phone string `json:"phone" ` // 手机号
	Name  string `json:"name" `  // 姓名
}

// GrowthDetail 成长值详情
type GrowthDetail struct {
	ID            uint64 `json:"id" `         // id
	Status        uint8  `json:"status"`      // 状态 1:有效 2:无效
	LevelTaskName string `json:"levelTaskId"` // 等级任务名称
	GrowthValue   uint64 `json:"growthValue"` // 成长值
	CreatedAt     string `json:"createdAt"`   // 创建时间
	Remark        string `json:"remark"`      // 备注
}

type GrowthCreateReq struct {
	ID          uint64 `json:"id" `                             // id
	MemberID    uint64 `json:"memberId" validate:"required"`    // 会员id
	Status      uint8  `json:"status"  enums:"1,2"`             // 状态 1:有效 2:无效
	GrowthValue uint64 `json:"growthValue" validate:"required"` // 成长值
	TaksID      uint64 `json:"taksId" validate:"required"`      // 任务id
}
