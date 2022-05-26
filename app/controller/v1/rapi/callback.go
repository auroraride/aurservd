// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/payment"
    "github.com/labstack/echo/v4"
    "net/http"
)

type callback struct{}

var Callback = new(callback)

type callbackReq struct {
    Type  string `query:"type"`
    Token string `query:"token"`
    State string `query:"state"`
}

func (*callback) RiderCallback(c echo.Context) error {
    req := new(callbackReq)
    _ = c.Bind(req)
    return nil
}

// ESignCallback E签宝回调
func (*callback) ESignCallback(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]int{"code": 200})
}

// AlipayCallback 支付宝回调
func (*callback) AlipayCallback(c echo.Context) (err error) {
    service.NewOrder().OrderPaid(payment.NewAlipay().Notification(c.Request(), c.Response()))
    return nil
}

// WechatpayCallback 微信回调
func (*callback) WechatpayCallback(c echo.Context) (err error) {
    res, trade := payment.NewWechat().Notification(c.Request())
    go service.NewOrder().OrderPaid(trade)
    return c.JSON(http.StatusOK, res)
}
