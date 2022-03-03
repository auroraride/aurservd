// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// StatusResponse 默认返回
type StatusResponse struct {
    Status bool `json:"status"` // 默认接口成功返回
}

// Modifier 修改人
type Modifier struct {
    ID    uint64 `json:"id,omitempty"`
    Name  string `json:"name,omitempty"`
    Phone string `json:"phone,omitempty"`
}

// SmsReq 短信请求
type SmsReq struct {
    Phone       string `json:"phone" validate:"required"`       // 手机号
    CaptchaCode string `json:"captchaCode" validate:"required"` // captcha 验证码
}

// SmsRes 短信发送返回
type SmsRes struct {
    Id string `json:"id"` // 任务ID
}

// AliyunOssStsRes 阿里云oss临时凭证
type AliyunOssStsRes struct {
    AccessKeySecret string `json:"accessKeySecret,omitempty"`
    Expiration      string `json:"expiration,omitempty"`
    AccessKeyId     string `json:"accessKeyId,omitempty"`
    StsToken        string `json:"stsToken,omitempty"`
    Bucket          string `json:"bucket,omitempty"`
    Region          string `json:"region,omitempty"`
}
