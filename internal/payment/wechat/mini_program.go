package wechat

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
)

type miniProgramClient struct {
	*core.Client
	*commonClient
}

func NewMiniProgram() *miniProgramClient {
	client := NewWechatClientWithConfig(ar.Config.Payment.Wechat)
	return &miniProgramClient{
		Client:       client,
		commonClient: newCommonClient(client),
	}
}

func (c *miniProgramClient) Miniprogram(appID, openID string, pc *model.PaymentCache) (*jsapi.PrepayWithRequestPaymentResponse, error) {
	amount, subject, no, _ := pc.GetPaymentArgs()
	cfg := ar.Config.Payment.Wechat

	svc := jsapi.JsapiApiService{Client: c.Client}
	resp, _, err := svc.PrepayWithRequestPayment(context.Background(), jsapi.PrepayRequest{
		Appid:       core.String(appID),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(subject),
		OutTradeNo:  core.String(no),
		TimeExpire:  core.Time(time.Now().Add(10 * time.Minute)),
		NotifyUrl:   core.String(cfg.NotifyUrl),
		Payer:       &jsapi.Payer{Openid: core.String(openID)},
		Amount: &jsapi.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(int64(math.Round(amount * 100))),
		},
	})

	if err != nil {
		zap.L().Error("微信Miniprogram支付调用失败", zap.Error(err))
		return nil, err
	}

	if resp.PrepayId == nil {
		return nil, errors.New("支付二维码获取失败")
	}

	return resp, nil
}
