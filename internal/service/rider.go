// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

type rider struct {
}

func NewRiderService() *rider {
    return &rider{}
}

func (*rider) SendSignSms(phone string, captchaCode string) {

}