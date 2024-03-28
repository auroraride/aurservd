package definition

import (
	"github.com/auroraride/aurservd/app/model"
)

type AgreementUserType uint8

const (
	AgreementUserTypePersonal AgreementUserType = iota + 1 // 协议类型 个签
	AgreementUserTypeGroup                                 //  协议类型 团签
)

func (a AgreementUserType) Value() uint8 {
	return uint8(a)
}

// AgreementCreateReq 创建协议
type AgreementCreateReq struct {
	Name          string `json:"name" validate:"required"`          // 协议名称
	UserType      uint8  `json:"userType" validate:"required"`      // 用户类型 1:个签 2:团签
	Content       string `json:"content" validate:"required"`       // 内容
	ForceReadTime uint8  `json:"forceReadTime" validate:"required"` // 强制阅读时间
	IsDefault     *bool  `json:"isDefault" validate:"required"`     // 是否为默认协议
}

type AgreementModifyReq struct {
	model.IDParamReq
	AgreementCreateReq
}

// AgreementDetail 协议返回
type AgreementDetail struct {
	ID            uint64 `json:"id"`                  // ID
	Name          string `json:"name"`                // 协议名称
	Content       string `json:"content"`             // 内容
	UserType      uint8  `json:"userType,omitempty"`  // 用户类型 1:个签 2:团签
	ForceReadTime uint8  `json:"forceReadTime"`       // 强制阅读时间
	IsDefault     *bool  `json:"isDefault,omitempty"` // 是否为默认协议
	CreatedAt     string `json:"createdAt,omitempty"` // 创建时间
	URL           string `json:"url"`                 // URL
	Hash          string `json:"hash"`                // Hash
}

// AgreementSearchReq 协议搜索
type AgreementSearchReq struct {
	UserType *uint8 `json:"userType" query:"userType"` // 用户类型 1:个签 2:团签
}

// AgreementSelectionRes 协议选择
type AgreementSelectionRes struct {
	ID        uint64 `json:"id"`        // ID
	Name      string `json:"name"`      // 协议名称
	IsDefault bool   `json:"isDefault"` // 是否为默认协议 true:是 false:否
	Hash      string `json:"hash"`      // Hash
}

type AppAgreementDetailRes struct {
	Title   string `json:"title"`   // 标题
	Content string `json:"content"` // 内容
}
