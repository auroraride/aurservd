// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-11
// Based on aurservd by liasica, magicrolan@qq.com.

package definition

// PersonCertificationOcrRes 腾讯人脸核验OCR参数
type PersonCertificationOcrRes struct {
	AppID   string `json:"appId"`   // WBAppid
	UserId  string `json:"userId"`  // 用户唯一标识
	OrderNo string `json:"orderNo"` // 订单号
	Version string `json:"version"` // 版本号
	Nonce   string `json:"nonce"`   // 随机字符串
	Sign    string `json:"sign"`    // 签名
}

// PersonCertificationFaceReq 实名认证人脸核身请求参数
type PersonCertificationFaceReq struct {
	OrderNo  string `json:"orderNo" query:"orderNo"`   // 订单号
	Identity string `json:"identity" query:"identity"` // 身份信息 // TODO: 加密传输
}

type PersonCertificationFaceRes struct {
}
