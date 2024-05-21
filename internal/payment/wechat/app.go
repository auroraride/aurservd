package wechat

import (
	"context"
	"errors"
	"io"
	"math"
	"time"

	"github.com/auroraride/adapter/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
)

type appClient struct {
	*core.Client
	*commonClient
}

func NewApp() *appClient {
	client := NewWechatClientWithConfig(ar.Config.Payment.Wechatpay)
	return &appClient{
		Client:       client,
		commonClient: newCommonClient(client),
	}
}

// AppPay APP支付
func (c *appClient) AppPay(pc *model.PaymentCache) (string, error) {
	amount, subject, no, _ := pc.GetPaymentArgs()
	cfg := ar.Config.Payment.Wechatpay

	svc := app.AppApiService{
		Client: c.Client,
	}

	resp, result, err := svc.PrepayWithRequestPayment(context.Background(), app.PrepayRequest{
		Appid:       core.String(cfg.AppID),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(subject),
		OutTradeNo:  core.String(no),
		TimeExpire:  core.Time(time.Now().Add(10 * time.Minute)),
		NotifyUrl:   core.String(cfg.NotifyUrl),
		Amount: &app.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(int64(math.Round(amount * 100))),
		},
	})

	if result == nil {
		zap.L().Error("微信App支付调用失败result为空")
		return "", err
	}

	if err != nil {
		b, _ := io.ReadAll(result.Response.Body)
		zap.L().Error("微信App支付调用失败", log.ResponseBody(b), zap.Error(err))
		return "", err
	}

	var out struct {
		*app.PrepayWithRequestPaymentResponse
		AppID string `json:"appId"`
	}

	out.PrepayWithRequestPaymentResponse = resp
	out.AppID = cfg.AppID

	b, _ := jsoniter.Marshal(out)

	return string(b), nil
}

func (c *appClient) Native(pc *model.PaymentCache) (string, error) {
	amount, subject, no, attach := pc.GetPaymentArgs()
	cfg := ar.Config.Payment.Wechatpay

	svc := native.NativeApiService{Client: c.Client}
	resp, _, err := svc.Prepay(context.Background(), native.PrepayRequest{
		Appid:       core.String(cfg.AppID),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(subject),
		OutTradeNo:  core.String(no),
		TimeExpire:  core.Time(time.Now().Add(10 * time.Minute)),
		NotifyUrl:   core.String(cfg.NotifyUrl),
		Attach:      core.String(attach),
		Amount: &native.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(int64(math.Round(amount * 100))),
		},
	})

	if err != nil {
		zap.L().Error("微信Native支付调用失败", zap.Error(err))
		return "", err
	}

	if resp.CodeUrl == nil {
		return "", errors.New("支付二维码获取失败")
	}

	return *resp.CodeUrl, nil
}
