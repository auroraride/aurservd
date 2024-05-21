package rapi

import (
	"net/http"

	"github.com/auroraride/adapter/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/payment/alipay"
)

type callback struct{}

var Callback = new(callback)

// AlipayFandAuthFreeze 支付宝资金授权冻结回调
func (*callback) AlipayFandAuthFreeze(c echo.Context) (err error) {
	res := alipay.NewApp().NotificationFandAuthFreeze(c.Request())
	zap.L().Info("支付宝资金授权冻结缓存更新", log.JsonData(res))
	biz.NewOrderBiz().DoPayment(res)
	return c.String(http.StatusOK, "success")
}

// AlipayFandAuthUnfreeze 支付宝资金授权解冻
func (*callback) AlipayFandAuthUnfreeze(c echo.Context) (err error) {
	res := alipay.NewApp().NotificationFandAuthUnfreeze(c.Request())
	zap.L().Info("支付宝资金授权解冻回调", log.JsonData(res))
	service.NewOrder().RefundSuccess(res.Refund)
	return c.String(http.StatusOK, "success")
}

// AlipayTradePay 扣款完成回调通知
func (*callback) AlipayTradePay(c echo.Context) (err error) {
	res := alipay.NewApp().NotificationTradePay(c.Request())
	zap.L().Info("支付宝扣款完成回调通知", log.JsonData(res))
	err = biz.NewOrderBiz().DoPaymentFreezeToPay(res)
	if err != nil {
		zap.L().Error("支付宝资金授权解冻回调失败", zap.Error(err))
	}
	return c.String(http.StatusOK, "success")
}

// AlipayMiniProgramPay 支付宝小程序支付回调
func (*callback) AlipayMiniProgramPay(c echo.Context) (err error) {
	res := alipay.NewMiniProgram().Notification(c.Request())
	zap.L().Info("支付宝小程序支付回调", log.JsonData(res))
	biz.NewOrderBiz().DoPayment(res)
	return c.String(http.StatusOK, "success")
}
