// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/payment"
    jsoniter "github.com/json-iterator/go"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
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
    go service.NewContract().Notice(c.Request())
    return c.JSON(http.StatusOK, map[string]int{"code": 200})
}

// AlipayCallback 支付宝回调
func (*callback) AlipayCallback(c echo.Context) (err error) {
    res := payment.NewAlipay().Notification(c.Request())

    b, _ := jsoniter.MarshalIndent(res, "", "  ")
    log.Infof("支付宝支付缓存更新: %s", b)

    service.NewOrder().DoPayment(res)
    return c.String(http.StatusOK, "success")
}

// WechatPayCallback 微信回调
func (*callback) WechatPayCallback(c echo.Context) (err error) {
    res := payment.NewWechat().Notification(c.Request())
    service.NewOrder().DoPayment(res)
    return c.JSON(http.StatusOK, ar.Map{"code": "SUCCESS", "message": "成功"})
}

// WechatRefundCallback 微信退款回调
func (*callback) WechatRefundCallback(c echo.Context) (err error) {
    res := payment.NewWechat().RefundNotification(c.Request())
    service.NewOrder().DoPayment(res)
    return c.JSON(http.StatusOK, ar.Map{"code": "SUCCESS", "message": "成功"})
}
