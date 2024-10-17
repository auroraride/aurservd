package rapi

import (
	"net/http"

	"github.com/auroraride/adapter/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/purchase/internal/service"
	"github.com/auroraride/aurservd/internal/payment/alipay"
	"github.com/auroraride/aurservd/internal/payment/wechat"
)

type callback struct{}

var Callback = new(callback)

// PurchaseAlipay 购买商品回调
func (*callback) PurchaseAlipay(c echo.Context) (err error) {
	res := alipay.NewApp().PurchaseNotification(c.Request())
	zap.L().Info("支付宝购买商品回调", log.JsonData(res))
	err = service.NewPayment().DoPayment(res)
	if err != nil {
		zap.L().Error("支付宝购买DoPayment错误", zap.Error(err))
	}
	return c.String(http.StatusOK, "success")
}

// PurchaseWechat 购买商品回调
func (*callback) PurchaseWechat(c echo.Context) (err error) {
	res := wechat.NewApp().PurchaseNotification(c.Request())
	zap.L().Info("微信购买商品回调", log.JsonData(res))
	err = service.NewPayment().DoPayment(res)
	if err != nil {
		zap.L().Error("微信DoPayment错误", zap.Error(err))
	}
	return c.String(http.StatusOK, "success")
}
