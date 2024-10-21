package definition

import "github.com/auroraride/aurservd/app/model"

// RiderChangePhoneReq 骑手修改手机号
type RiderChangePhoneReq struct {
	Phone   string `json:"phone" validate:"required"`                 // 手机号
	SmsId   string `json:"smsId" validate:"required" trans:"短信ID"`    // 短信ID
	SmsCode string `json:"smsCode" validate:"required" trans:"短信验证码"` // 短信验证码
}

// RiderAllocatedReq 骑手分配信息
type RiderAllocatedReq struct {
	Serial string `json:"serial" validate:"required" query:"serial"`
	Phone  string `json:"phone" validate:"required" query:"phone"`
}

// RiderSignupReq 骑手登录V2
type RiderSignupReq struct {
	Phone         string  `json:"phone,omitempty" validate:"required_if=SigninType 1" trans:"电话"`    // 电话
	SmsId         string  `json:"smsId,omitempty" validate:"required_if=SigninType 1" trans:"短信ID"`  // 短信ID
	SmsCode       string  `json:"smsCode,omitempty" validate:"required_if=SigninType 1" trans:"验证码"` // 验证码
	SigninType    *uint64 `json:"signinType"`                                                        // 登录类型 1:短信登录 2:授权登录
	AuthType      uint8   `json:"authType" validate:"required_if=SigninType 2"`                      // 授权类型 1:微信 2:支付宝
	EncryptedData string  `json:"encryptedData" validate:"required_if=AuthType 2 SigninType 2"`      // 支付宝获取手机加密数据
	AuthCode      string  `json:"authCode" validate:"required_if=AuthType 1 SigninType 2"`           // 授权码
}

// RiderSetMobPushReq 骑手推送ID
type RiderSetMobPushReq struct {
	PushId string `json:"pushId" validate:"required" trans:"骑手推送ID"`
}

// RiderSigninRes 骑手登录数据返回
type RiderSigninRes struct {
	ID                uint64                   `json:"id"`
	Phone             string                   `json:"phone"`                       // 电话
	Name              string                   `json:"name"`                        // 姓名, 实名认证后才会有
	Token             string                   `json:"token,omitempty"`             // 认证token
	IsNewDevice       bool                     `json:"isNewDevice"`                 // 是否新设备
	IsAuthed          bool                     `json:"isAuthed"`                    // 是否已认证
	IsContactFilled   bool                     `json:"isContactFilled"`             // 联系人是否添加
	Contact           *model.RiderContact      `json:"contact,omitempty"`           // 联系人
	Qrcode            string                   `json:"qrcode"`                      // 二维码
	Deposit           float64                  `json:"deposit"`                     // 需缴押金
	OrderNotActived   *bool                    `json:"orderNotActived,omitempty"`   // 是否存在未激活订单
	Subscribe         *model.Subscribe         `json:"subscribe,omitempty"`         // 骑士卡
	Enterprise        *model.Enterprise        `json:"enterprise,omitempty"`        // 所属企业
	UseStore          bool                     `json:"useStore"`                    // 是否可使用门店办理业务
	CabinetBusiness   bool                     `json:"cabinetBusiness"`             // 是否可以自主使用电柜办理业务
	EnterpriseContact *model.EnterpriseContact `json:"enterpriseContact,omitempty"` // 团签联系方式
	ExitEnterprise    bool                     `json:"exitEnterprise"`              // 判断能否退出团签
	Station           *model.EnterpriseStation `json:"station,omitempty"`           // 站点
	ContractDocID     string                   `json:"contractDocId,omitempty"`     // 签署合同编号
	Purchase          bool                     `json:"purchase"`                    // 待支付购车订单
}
