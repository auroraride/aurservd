// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package app

// http headers
const (
	HeaderDispositionType = "Content-Disposition"
	// HeaderContentType the ContentType Header
	HeaderContentType = "Content-Type"
	// HeaderCaptchaID 图片验证码ID
	HeaderCaptchaID = "X-Captcha-Id"
	// HeaderDeviceSerial 骑手设备序列号 (由此判定是否更换了设备)
	HeaderDeviceSerial = "X-Device-Serial"
	// HeaderDeviceType 骑手设备类型
	HeaderDeviceType = "X-Device-Type"
	// HeaderPushId 骑手设备推送ID
	HeaderPushId = "X-Push-Id"
	// HeaderRiderToken 骑手token
	HeaderRiderToken = "X-Rider-Token"
	// HeaderManagerToken 后台token
	HeaderManagerToken = "X-Manager-Token"
	// HeaderEmployeeToken 店员token
	HeaderEmployeeToken = "X-Employee-Token"
	// HeaderAgentToken 代理商token
	HeaderAgentToken = "X-Agent-Token"
	// HeaderPromotionToken 推广token
	HeaderPromotionToken = "X-Promotion-Token"
)
