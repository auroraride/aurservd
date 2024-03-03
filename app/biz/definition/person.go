// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-11
// Based on aurservd by liasica, magicrolan@qq.com.

package definition

// PersonCertificationOcrClientRes 腾讯人身核验OCR参数
type PersonCertificationOcrClientRes struct {
	AppID   string `json:"appId,omitempty"`   // WBAppid
	UserId  string `json:"userId,omitempty"`  // 用户唯一标识
	OrderNo string `json:"orderNo,omitempty"` // 订单号
	Version string `json:"version,omitempty"` // 版本号
	Nonce   string `json:"nonce,omitempty"`   // 随机字符串
	Sign    string `json:"sign,omitempty"`    // 签名
}

// PersonCertificationOcrCloudReq 阿里云OCR签名请求参数
type PersonCertificationOcrCloudReq struct {
	Hash string `json:"hash" query:"hash" validate:"required"` // 文件 sha256 hex string
}

// PersonCertificationOcrCloudRes 阿里云OCR签名响应参数
type PersonCertificationOcrCloudRes struct {
	ContentType   string `json:"contentType"`
	Action        string `json:"action"`
	Date          string `json:"date"`
	Token         string `json:"token"`
	Nonce         string `json:"nonce"`
	Version       string `json:"version"`
	Authorization string `json:"authorization"`
}

// PersonCertificationFaceReq 实名认证人脸核身请求参数
type PersonCertificationFaceReq struct {
	OrderNo  string `json:"orderNo" query:"orderNo"`      // 订单号，用户使用OCR识别时不为空
	Identity string `json:"identity" validate:"required"` // 加密身份信息
}

// PersonCertificationFaceRes 实名认证人脸核身参数
type PersonCertificationFaceRes struct {
	PersonCertificationOcrClientRes
	FaceId      string `json:"faceId,omitempty"`
	Licence     string `json:"licence,omitempty"`
	BindedPhone string `json:"bindedPhone,omitempty"` // 已绑定的其他手机号
}

// PersonCertificationFaceResultReq 实名认证人脸核身结果请求参数
type PersonCertificationFaceResultReq struct {
	OrderNo string `json:"orderNo" query:"orderNo" validate:"required"` // 订单号
}

// PersonCertificationFaceResultRes 实名认证人脸核身结果
type PersonCertificationFaceResultRes struct {
	Success bool `json:"success"` // 是否成功
}
