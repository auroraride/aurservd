package definition

import "github.com/auroraride/aurservd/app/model"

type ContractTemplateAimed uint8

const (
	ContractTemplateAimedPersonal   ContractTemplateAimed = iota + 1 // 个签
	ContractTemplateAimedEnterprise                                  // 团签
)

func (c ContractTemplateAimed) Value() uint8 {
	return uint8(c)
}

type ContractTemplate struct {
	ID        uint64                `json:"id"`        // ID
	Url       string                `json:"url"`       // 合同链接
	Name      string                `json:"name"`      // 合同名称
	Aimed     ContractTemplateAimed `json:"aimed"`     // 用户类型 1:个签 2:团签
	PlanType  model.PlanType        `json:"planType"`  // 订阅类型 1:单电 2:车电
	Enable    bool                  `json:"enable"`    // 是否启用
	Hash      string                `json:"hash"`      // 合同hash
	Remark    string                `json:"remark"`    // 备注
	CreatedAt string                `json:"createdAt"` // 创建时间
}

// ContractTemplateCreateReq 创建合同模板请求
type ContractTemplateCreateReq struct {
	Url      string                `json:"url" validate:"required"`   // 合同链接
	Name     string                `json:"name" validate:"required" ` // 合同名称
	Aimed    ContractTemplateAimed `json:"aimed"`                     // 用户类型 1:个签 2:团签
	PlanType model.PlanType        `json:"planType"`                  // 订阅类型 1:单电 2:车电
	Remark   *string               `json:"remark"`                    // 备注
	Enable   *bool                 `json:"enable"`                    // 是否启用
	Hash     string                `json:"hash" validate:"required"`  // 合同hash
}

// ContractTemplateListRes 合同模板列表返回
type ContractTemplateListRes struct {
	ContractTemplate
}

// ContractTemplateModifyReq 修改合同模板请求
type ContractTemplateModifyReq struct {
	model.IDParamReq
	Name     string                `json:"name" validate:"required" ` // 合同名称
	Aimed    ContractTemplateAimed `json:"aimed"`                     // 用户类型 1:个签 2:团签
	PlanType model.PlanType        `json:"planType"`                  // 订阅类型 1:单电 2:车电
	Remark   *string               `json:"remark"`                    // 备注
	Enable   *bool                 `json:"enable"`                    // 是否启用
}
