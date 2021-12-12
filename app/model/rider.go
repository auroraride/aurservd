// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// RiderSignupReq 骑手登录请求数据
type RiderSignupReq struct {
    Phone   string `json:"phone" validate:"required"`
    SmsId   string `json:"smsId" validate:"required"`
    SmsCode string `json:"smsCode" validate:"required"`
}

// RiderSigninRes 骑手登录数据返回
type RiderSigninRes struct {
    Id          uint64 `json:"id"`
    Token       string `json:"token"`
    IsNewDevice bool   `json:"isNewDevice"`
}
