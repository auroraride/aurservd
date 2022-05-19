// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// RiderTokenPermission 骑手token权限, 以此判定登陆后动作
type RiderTokenPermission uint8

const (
    RiderTokenPermissionCommon    RiderTokenPermission = iota // 普通权限
    RiderTokenPermissionAuth                                  // 需要实名验证
    RiderTokenPermissionNewDevice                             // 更换设备需要人脸验证
)

// RiderContext TODO 骑手上下文
type RiderContext struct {
}

// RiderSignupReq 骑手登录请求数据
type RiderSignupReq struct {
    Phone   string `json:"phone" validate:"required"`
    SmsId   string `json:"smsId" validate:"required"`
    SmsCode string `json:"smsCode" validate:"required"`
}

// RiderSigninRes 骑手登录数据返回
type RiderSigninRes struct {
    Id              uint64        `json:"id"`
    Token           string        `json:"token"`
    IsNewDevice     bool          `json:"isNewDevice"`
    IsAuthed        bool          `json:"isAuthed"`        // 是否已认证
    IsContactFilled bool          `json:"isContactFilled"` // 联系人是否添加
    Contact         *RiderContact `json:"contact,omitempty"`
    Qrcode          string        `json:"qrcode"`
}

// RiderContact 紧急联系人
type RiderContact struct {
    Name     string `json:"name" validate:"required" trans:"联系人姓名"`
    Phone    string `json:"phone" validate:"required,phone" trans:"联系人电话"`
    Relation string `json:"relation" validate:"required" trans:"关系"`
}

// RiderSampleInfo 骑手简单信息
type RiderSampleInfo struct {
    ID    uint64 `json:"id"`    // 骑手ID
    Name  string `json:"name"`  // 骑手姓名
    Phone string `json:"phone"` // 骑手电话
}
