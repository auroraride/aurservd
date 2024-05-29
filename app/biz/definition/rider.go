package definition

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
