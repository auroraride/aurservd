// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package payment

import (
	"context"
	"errors"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"

	"github.com/auroraride/aurservd/pkg/snag"

	"github.com/auroraride/adapter/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/cache"
)

const (
	WechatPayCodeSuccess = "SUCCESS"
	WechatPayCodeFail    = "FAIL"
)

var _wechat *wechatClient

type wechatClient struct {
	*core.Client
	notifyClient *notify.Handler
}

func newWechatClient() *wechatClient {
	cfg := ar.Config.Payment.Wechat
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.PrivateKeyPath)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	ctx := context.Background()
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(cfg.MchID, cfg.MchCertificateSerialNumber, mchPrivateKey, cfg.MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	// 获取商户号对应的微信支付平台证书访问器
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(cfg.MchID)

	return &wechatClient{
		Client:       client,
		notifyClient: notify.NewNotifyHandler(cfg.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certVisitor)),
	}
}

func NewWechat() *wechatClient {
	return _wechat
}

// AppPay APP支付
func (c *wechatClient) AppPay(pc *model.PaymentCache) (string, error) {
	amount, subject, no, _ := pc.GetPaymentArgs()
	cfg := ar.Config.Payment.Wechat

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

func (c *wechatClient) Native(pc *model.PaymentCache) (string, error) {
	amount, subject, no, attach := pc.GetPaymentArgs()
	cfg := ar.Config.Payment.Wechat

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

func (c *wechatClient) Miniprogram(appID, openID string, pc *model.PaymentCache) (*jsapi.PrepayWithRequestPaymentResponse, error) {
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

type WechatPayResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Refund 退款
func (c *wechatClient) Refund(req *model.PaymentRefund) {
	svc := refunddomestic.RefundsApiService{Client: c.Client}
	cfg := ar.Config.Payment.Wechat
	resp, result, err := svc.Create(context.Background(),
		refunddomestic.CreateRequest{
			TransactionId: core.String(req.TradeNo),
			OutRefundNo:   core.String(req.OutRefundNo),
			Reason:        core.String(req.Reason),
			NotifyUrl:     core.String(cfg.RefundUrl),
			Amount: &refunddomestic.AmountReq{
				Currency: core.String("CNY"),
				Refund:   core.Int64(int64(math.Round(req.RefundAmount * 100))),
				Total:    core.Int64(int64(math.Round(req.Total * 100))),
			},
		},
	)
	if err != nil {
		b, _ := io.ReadAll(result.Response.Body)
		zap.L().Error("微信退款调用失败", log.ResponseBody(b), zap.Error(err))
		snag.Panic("退款处理失败: " + err.Error())
	}

	req.Request = true

	s := resp.Status
	if *s == refunddomestic.STATUS_SUCCESS {
		req.Success = true
		req.Time = *resp.SuccessTime
	}
}

// Notification 微信支付回调
func (c *wechatClient) Notification(req *http.Request) *model.PaymentCache {
	transaction := new(payments.Transaction)
	nq, err := c.notifyClient.ParseNotifyRequest(context.Background(), req, transaction)
	if err != nil {
		zap.L().Error("微信回调解析失败", zap.Error(err))
		return nil
	}

	b, _ := jsoniter.Marshal(transaction)
	nb, _ := jsoniter.Marshal(nq)
	zap.L().Info("微信支付回调反馈", zap.ByteString("transaction", b), zap.ByteString("notifiy", nb))

	pc := new(model.PaymentCache)
	// 从缓存中获取订单数据
	out := transaction.OutTradeNo
	err = cache.Get(context.Background(), *out).Scan(pc)
	if err != nil {
		zap.L().Error("从缓存获取订单信息失败", zap.Error(err))
		return nil
	}

	b, _ = jsoniter.Marshal(pc)
	zap.L().Info("获取到微信支付回调缓存", zap.ByteString("transaction", b))

	state := transaction.TradeState
	if *state != "SUCCESS" {
		return nil
	}

	tranID := *(transaction.TransactionId)
	switch pc.CacheType {
	case model.PaymentCacheTypePlan:
		pc.Subscribe.TradeNo = tranID
	case model.PaymentCacheTypeOverdueFee:
		pc.OverDueFee.TradeNo = tranID
	case model.PaymentCacheTypeAssistance:
		pc.Assistance.TradeNo = tranID
	case model.PaymentCacheTypeAgentPrepay:
		pc.AgentPrepay.TradeNo = tranID
	case model.PaymentCacheTypeDeposit:
		pc.DepositCredit.TradeNo = tranID
	default:
		return nil
	}

	b, _ = jsoniter.Marshal(pc)
	zap.L().Info("微信支付缓存更新", zap.ByteString("transaction", b))

	return pc
}

type WechatRefundTransaction struct {
	Mchid               string    `json:"mchid"`
	TransactionId       string    `json:"transaction_id"`
	OutTradeNo          string    `json:"out_trade_no"`
	RefundId            string    `json:"refund_id"`
	OutRefundNo         string    `json:"out_refund_no"`
	RefundStatus        string    `json:"refund_status"`
	SuccessTime         time.Time `json:"success_time"`
	UserReceivedAccount string    `json:"user_received_account"`
	Amount              struct {
		Total       int `json:"total"`
		Refund      int `json:"refund"`
		PayerTotal  int `json:"payer_total"`
		PayerRefund int `json:"payer_refund"`
	} `json:"amount"`
}

// RefundNotification 微信退款回调
func (c *wechatClient) RefundNotification(req *http.Request) *model.PaymentCache {
	transaction := new(WechatRefundTransaction)
	nq, err := c.notifyClient.ParseNotifyRequest(context.Background(), req, transaction)
	if err != nil {
		zap.L().Error("微信退款回调解析失败", zap.Error(err))
		return nil
	}

	b, _ := jsoniter.Marshal(transaction)
	nb, _ := jsoniter.Marshal(nq)
	zap.L().Info("微信支付回调反馈", zap.ByteString("transaction", b), zap.ByteString("notifiy", nb))

	pc := new(model.PaymentCache)

	// 从缓存中获取订单数据
	err = cache.Get(context.Background(), transaction.OutRefundNo).Scan(pc)
	if err != nil {
		zap.L().Error("从缓存获取订单信息失败", zap.Error(err))
		return nil
	}

	b, _ = jsoniter.Marshal(pc)
	zap.L().Info("获取到微信退款回调缓存", zap.ByteString("transaction", b))

	if transaction.RefundStatus != "SUCCESS" {
		return nil
	}

	pc.Refund.Success = true
	pc.Refund.Request = true
	pc.Refund.Time = transaction.SuccessTime

	b, _ = jsoniter.Marshal(pc)
	zap.L().Info("微信支付缓存更新", zap.ByteString("transaction", b))

	return pc
}
