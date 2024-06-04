package alipay

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
)

type miniProgramClient struct {
	*alipay.Client
	*commonClient
}

func NewMiniProgram() *miniProgramClient {
	client := newAlipayClientWithConfig(ar.Config.Payment.AlipayMiniprogram)
	return &miniProgramClient{
		Client:       client,
		commonClient: newCommonClient(client),
	}
}

// Trade  统一收单交易创建接口
func (c *miniProgramClient) Trade(pc *model.PaymentCache) (string, error) {
	cfg := ar.Config.Payment.AlipayMiniprogram
	amount, subject, no, _ := pc.GetPaymentArgs()
	trade := alipay.TradeCreate{
		Trade: alipay.Trade{
			TotalAmount: fmt.Sprintf("%.2f", amount),
			NotifyURL:   cfg.NotifyUrl,
			Subject:     subject,
			OutTradeNo:  no,
			ProductCode: "JSAPI_PAY",
		},
		OpAppId:     cfg.Appid,
		BuyerOpenId: *pc.Subscribe.BuyerOpenId,
	}
	res, err := c.TradeCreate(context.Background(), trade)
	if err != nil {
		return "", err
	}
	if !res.IsSuccess() {
		return "", errors.New("支付宝小程序支付失败")
	}
	return res.TradeNo, nil
}

// GetPhoneNumber 获取手机号
func (c *miniProgramClient) GetPhoneNumber(code string) (string, error) {
	phoneNumber, err := c.DecodePhoneNumber(context.Background(), code)
	if err != nil {
		return "", err
	}
	if phoneNumber.Code != "10000" {
		return "", errors.New(phoneNumber.Msg)
	}
	return phoneNumber.Mobile, nil
}

// GetOpenid 获取openid
func (c *miniProgramClient) GetOpenid(code string) (string, error) {
	auth, err := c.SystemOauthToken(context.Background(), alipay.SystemOauthToken{
		GrantType: "authorization_code",
		Code:      code,
	})

	if err != nil {
		zap.L().Error("获取支付宝openid失败", zap.Error(err))
		return "", err
	}

	if auth.OpenId == "" {
		return "", errors.New("获取openid失败")
	}
	return auth.OpenId, nil
}
