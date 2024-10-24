package alipay

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/auroraride/adapter/log"
	"github.com/golang-module/carbon/v2"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/tools"
)

type appClient struct {
	*alipay.Client
	*commonClient
}

func NewApp() *appClient {
	client := newAlipayClientWithConfig(ar.Config.Payment.Alipay)
	return &appClient{
		Client:       client,
		commonClient: newCommonClient(client),
	}
}

// AppPay app支付
func (c *appClient) AppPay(pc *model.PaymentCache) (string, error) {
	cfg := ar.Config.Payment.Alipay
	amount, subject, no, _ := pc.GetPaymentArgs()
	trade := alipay.TradeAppPay{
		Trade: alipay.Trade{
			TotalAmount: fmt.Sprintf("%.2f", amount),
			NotifyURL:   cfg.NotifyUrl,
			Subject:     subject,
			OutTradeNo:  no,
			TimeExpire:  time.Now().Add(10 * time.Minute).Format(carbon.DateTimeLayout),
		},
	}
	return c.TradeAppPay(trade)
}

func (c *appClient) AppPayDemo() (string, string, error) {
	cfg := ar.Config.Payment.Alipay
	no := tools.NewUnique().NewSN()
	trade := alipay.TradeAppPay{
		Trade: alipay.Trade{
			TotalAmount: fmt.Sprintf("%.2f", 0.01),
			NotifyURL:   cfg.NotifyUrl,
			Subject:     "测试支付",
			OutTradeNo:  no,
			TimeExpire:  time.Now().Add(10 * time.Minute).Format(carbon.DateTimeLayout),
		},
	}

	s, err := c.TradeAppPay(trade)
	return s, no, err
}

func (c *appClient) Native(pc *model.PaymentCache) (string, error) {
	cfg := ar.Config.Payment.Alipay

	amount, subject, no, _ := pc.GetPaymentArgs()
	trade := alipay.TradePreCreate{
		Trade: alipay.Trade{
			TotalAmount: fmt.Sprintf("%.2f", amount),
			NotifyURL:   cfg.NotifyUrl,
			Subject:     subject,
			OutTradeNo:  no,
		},
	}
	res, err := c.TradePreCreate(context.Background(), trade)
	if err != nil {
		return "", err
	}

	if !res.IsSuccess() {
		return "", errors.New("支付宝二维码生成失败")
	}

	return res.QRCode, nil
}

// AppPayPurchase 购买商品
func (c *appClient) AppPayPurchase(req *model.PurchasePayReq) (string, error) {
	notifyUrl := ar.Config.Payment.PurchaseAlipayNotifyUrl
	trade := alipay.TradeAppPay{
		Trade: alipay.Trade{
			TotalAmount: fmt.Sprintf("%.2f", req.Amount),
			NotifyURL:   notifyUrl,
			Subject:     req.Subject,
			OutTradeNo:  req.OutTradeNo,
			TimeExpire:  time.Now().Add(10 * time.Minute).Format(carbon.DateTimeLayout),
		},
	}
	zap.L().Info("支付宝支付", log.JsonData(trade))
	return c.TradeAppPay(trade)
}
