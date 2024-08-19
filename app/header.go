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
	// HeaderAgentToken 代理token
	HeaderAgentToken = "X-Agent-Token"
	// HeaderMaintainerToken 运维token
	HeaderMaintainerToken = "X-Maintainer-Token"
	// HeaderPromotionToken 推广token
	HeaderPromotionToken = "X-Promotion-Token"
	// HeaderToastVisible 自动显示toast（APP 2.x+）
	HeaderToastVisible = "X-Toast-Visible"
	// HeaderAssetManagerToken 仓库后台token
	HeaderAssetManagerToken = "X-AssetManager-Token"
	// HeaderWarestoreToken 库管token
	HeaderWarestoreToken = "X-Warestore-Token"
)
