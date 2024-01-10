// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"io"
	"net/http"

	"github.com/auroraride/adapter/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/payment"
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
func (*callback) ESignCallback(c echo.Context) (err error) {
	var b []byte
	b, err = io.ReadAll(c.Request().Body)
	if err != nil {
		zap.L().Error("合同回调内容读取失败", zap.Error(err))
		return
	}

	go service.NewContract().Notice(b)
	return c.JSON(http.StatusOK, map[string]int{"code": 200})
}

// AlipayCallback 支付宝回调
func (*callback) AlipayCallback(c echo.Context) (err error) {
	res := payment.NewAlipay().Notification(c.Request())

	zap.L().Info("支付宝支付缓存更新", log.JsonData(res))

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

// AlipayFandAuthFreeze 支付宝资金授权冻结回调
func (*callback) AlipayFandAuthFreeze(c echo.Context) (err error) {
	res := payment.NewAlipay().NotificationFandAuthFreeze(c.Request())

	zap.L().Info("支付宝资金授权冻结缓存更新", log.JsonData(res))

	service.NewOrder().DoPayment(res)
	return c.String(http.StatusOK, "success")
}
