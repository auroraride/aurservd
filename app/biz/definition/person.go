// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-11
// Based on aurservd by liasica, magicrolan@qq.com.

package definition

// PersonCertificationReq 实名认证请求
type PersonCertificationReq struct {
	Name         *string `json:"name" validate:"required" trans:"姓名"`
	IDCardNumber *string `json:"idCardNumber" validate:"required" trans:"身份证号"`
	MetaInfo     *string `json:"metaInfo" validate:"required" trans:"环境参数"`
}

// PersonCertification 实名认证ID
type PersonCertification struct {
	CertifyId string `json:"certifyId" query:"certifyId" validate:"required" trans:"实名认证ID"`
}

// PersonCertificationResultRes 实名认证结果
type PersonCertificationResultRes struct {
	Passed bool `json:"passed"` // 是否通过
}
