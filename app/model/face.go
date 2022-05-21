// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-10
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// FaceAuthUrlResponse 人脸实名URL返回
type FaceAuthUrlResponse struct {
    Url string `json:"url"` // 实名认证跳转链接
}

// FaceResultReq 获取实名或人脸识别结果
type FaceResultReq struct {
    Token string `json:"token" param:"token"` // 实名或人脸识别token
}

// FaceAuthContext 实名认证上下文
type FaceAuthContext struct {
    RiderID uint64 // 骑手ID
    Serial  string // 设备序列号
}
