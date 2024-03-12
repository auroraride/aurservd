package rapi

import (
	"net/http"

	"github.com/auroraride/adapter/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/internal/payment"
)

type callback struct{}

var Callback = new(callback)

// AlipayFandAuthFreeze 支付宝资金授权冻结回调
func (*callback) AlipayFandAuthFreeze(c echo.Context) (err error) {
	res := payment.NewAlipay().NotificationFandAuthFreeze(c.Request())
	zap.L().Info("支付宝资金授权冻结缓存更新", log.JsonData(res))
	biz.NewOrderBiz().DoPayment(res)
	return c.String(http.StatusOK, "success")
}

// AlipayTradePay 扣款完成回调通知
func (*callback) AlipayTradePay(c echo.Context) (err error) {
	res := payment.NewAlipay().NotificationTradePay(c.Request())
	zap.L().Info("支付宝扣款完成回调通知", log.JsonData(res))
	biz.NewOrderBiz().DoPaymentFreezeToPay(res)
	return c.String(http.StatusOK, "success")
}

// 解冻
