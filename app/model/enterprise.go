// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
	EnterpriseStatusLack         uint8 = iota // 未合作
	EnterpriseStatusCollaborated              // 合作中
	EnterpriseStatusSuspensed                 // 已暂停
)

const (
	EnterprisePaymentPrepay  uint8 = iota + 1 // 预付费
	EnterprisePaymentPostPay                  // 后付费
)

type EnterpriseContract struct {
	ID    uint64 `json:"id,omitempty"` // 合同ID, 请求`M9018 编辑合同`时携带此字段为编辑, 不携带为新增
	Start string `json:"start" validate:"required" trans:"合同开始日期"`
	End   string `json:"end" validate:"required" trans:"合同结束日期"`
	File  string `json:"file" validate:"required" trans:"合同文件"`
}

type EnterprisePrice struct {
	ID       uint64  `json:"id"`
	Model    string  `json:"model"`
	Price    float64 `json:"price"`
	CityID   uint64  `json:"cityId"`
	CityName string  `json:"cityName"`
}

type EnterpriseContractModifyReq struct {
	EnterpriseContract
	EnterpriseID uint64 `json:"enterpriseId" validate:"required" trans:"企业ID"`
}

type EnterprisePriceModifyReq struct {
	ID           uint64  `json:"id,omitempty"` // 价格ID, 请求`M9016 编辑价格`时携带此字段为编辑, 不携带为新增
	CityID       uint64  `json:"cityId" validate:"required" trans:"城市"`
	Model        string  `json:"model" validate:"required" trans:"电池型号"`
	Price        float64 `json:"price" validate:"required" trans:"单价(元/天)"`
	EnterpriseID uint64  `json:"enterpriseId" validate:"required" trans:"企业ID"`
	Intelligent  bool    `json:"intelligent"` // 是否智能电池
}

type EnterprisePriceWithCity struct {
	ID          uint64  `json:"id"`          // 价格ID
	Model       string  `json:"model"`       // 电池型号
	Price       float64 `json:"price"`       // 单价(元/天)
	City        City    `json:"city"`        // 城市
	Intelligent bool    `json:"intelligent"` // 是否智能电池
}

// Enterprise 企业基础字段
type Enterprise struct {
	ID    uint64 `json:"id"`    // 企业ID
	Name  string `json:"name"`  // 企业名称
	Agent bool   `json:"agent"` // 是否代理商模式
}

type EnterpriseContact struct {
	Name  string `json:"name"`  // 联系姓名
	Phone string `json:"phone"` // 联系电话
}

// EnterpriseDetail 企业详细字段
type EnterpriseDetail struct {
	Name         *string  `json:"name" validate:"required" trans:"团签名称"`
	CompanyName  *string  `json:"companyName" validate:"required" trans:"企业全称"`
	Status       *uint8   `json:"status" enums:"0,1,2" validate:"required,min=0,max=2" trans:"合作状态"` // 0:未合作 1:合作中 2:已暂停
	ContactName  *string  `json:"contactName" validate:"required" trans:"联系人"`
	ContactPhone *string  `json:"contactPhone" validate:"required" trans:"联系电话"`
	IdcardNumber *string  `json:"idcardNumber" validate:"required" trans:"身份证号"`
	CityID       *uint64  `json:"cityId" validate:"required" trans:"所在城市"`
	Address      *string  `json:"address" validate:"required" trans:"企业地址"`
	Payment      *uint8   `json:"payment" validate:"required,min=1,max=2" enums:"1,2" trans:"付费方式"` // 1:预付费 2:后付费
	Deposit      *float64 `json:"deposit" validate:"required" trans:"押金"`
	Agent        *bool    `json:"agent"`              // 代理商 `true`是 `false`否
	UseStore     *bool    `json:"useStore,omitempty"` // 可使用门店 `true`允许 `false`不允许 (仅代理商模式生效), 骑手是否可以使用门店进行激活和退租
	Days         *[]int   `json:"days,omitempty"`     // 代理商时间选项
}

type EnterpriseDetailWithID struct {
	*EnterpriseDetail
	ID uint64 `json:"id" param:"id" validate:"required" trans:"企业ID"`
}

type EnterpriseListReq struct {
	PaginationReq
	CityID         *uint64 `json:"cityId" query:"cityId"`                 // 城市ID
	ContactKeyword *string `json:"contactKeyword" query:"contactKeyword"` // 联系人 姓名/电话/身份证 关键词
	Name           *string `json:"name" query:"name"`                     // 公司名称
	Status         *uint8  `json:"status" query:"status"`                 // 合作状态
	Payment        *uint8  `json:"payment" query:"payment" enums:"1,2"`   // 支付方式 1预付费 2后付费
	Start          *string `json:"start" query:"start"`                   // 合同到期时间晚于
	End            *string `json:"end" query:"end"`                       // 合同到期时间早于
	// StatementStart *string `json:"statementStart" query:"statementStart"` // 计费时间早于
	// StatementEnd   *string `json:"statementEnd" query:"statementEnd"`     // 计费时间晚于
	Agent *bool `json:"agent"` // 代理商 `true`是 `false`否
}

type EnterpriseRes struct {
	ID             uint64                    `json:"id"`                       // 企业ID
	Balance        float64                   `json:"balance"`                  // 可用余额
	Unsettlement   int                       `json:"unsettlement"`             // 未结算天数, 预付费企业此字段强制为0
	Name           string                    `json:"name"`                     // 团签名称
	CompanyName    string                    `json:"companyName"`              // 企业全称
	Status         uint8                     `json:"status" enums:"0,1,2" `    // 合作状态 0:未合作 1:已合作 2:已暂停
	ContactName    string                    `json:"contactName"`              // 联系人
	ContactPhone   string                    `json:"contactPhone"`             // 联系电话
	IdcardNumber   string                    `json:"idcardNumber"`             // 身份证号
	Address        string                    `json:"address"`                  // 企业地址
	Payment        uint8                     `json:"payment"`                  // 付费方式 1:预付费 2:后付费
	Deposit        float64                   `json:"deposit"`                  // 押金
	Riders         int                       `json:"riders"`                   // 骑手数量
	Contracts      []EnterpriseContract      `json:"contracts,omitempty"`      // 合同
	Prices         []EnterprisePriceWithCity `json:"prices,omitempty"`         // 价格列表
	City           City                      `json:"city"`                     // 城市
	StatementStart string                    `json:"statementStart,omitempty"` // 账单开始日期
	Agent          bool                      `json:"agent"`                    // 代理商 `true`是 `false`否
	UseStore       *bool                     `json:"useStore,omitempty"`       // 可使用门店 `true`允许 `false`不允许
	Days           *[]int                    `json:"days,omitempty"`           // 代理商时间选项
}

type EnterprisePrepaymentReq struct {
	ID     uint64  `json:"id" validate:"required" param:"id" trans:"企业ID"`
	Remark string  `json:"remark" validate:"required" trans:"备注"`
	Amount float64 `json:"amount" validate:"required" trans:"金额"`
}

type EnterpriseStationCreateReq struct {
	EnterpriseID uint64 `json:"enterpriseId" validate:"required"  trans:"企业ID"`
	Name         string `json:"name" validate:"required" trans:"站点名称"`
}

type EnterpriseStationModifyReq struct {
	Name string `json:"name" validate:"required" trans:"站点名称"`
	ID   uint64 `json:"id" validate:"required" param:"id" trans:"站点ID"`
}

type EnterpriseStation struct {
	ID   uint64 `json:"id"`   // 站点ID
	Name string `json:"name"` // 站点名称
}

type EnterpriseStationListReq struct {
	EnterpriseID uint64 `json:"enterpriseId" query:"enterpriseId" trans:"企业ID"`
}

type EnterprisePriceBatteryModelListReq struct {
	CityID uint64 `json:"cityId" validate:"required" query:"cityId" trans:"城市ID"`
}

type EnterpriseRiderSubscribeChooseReq struct {
	Model string `json:"model" validate:"required" trans:"电池型号"`
}

type EnterpriseRiderSubscribeChooseRes struct {
	Qrcode string `json:"qrcode"` // 二维码, 格式为SUBSCRIBE:订阅ID, 后续使用订阅ID请求状态
}

type EnterpriseRiderSubscribeStatusReq struct {
	ID uint64 `json:"id" validate:"required" query:"id" trans:"订阅ID"`
}
