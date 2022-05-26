// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package payment

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
    "github.com/wechatpay-apiv3/wechatpay-go/core"
    "github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
    "github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
    "github.com/wechatpay-apiv3/wechatpay-go/core/notify"
    "github.com/wechatpay-apiv3/wechatpay-go/core/option"
    "github.com/wechatpay-apiv3/wechatpay-go/services/payments"
    "github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
    "github.com/wechatpay-apiv3/wechatpay-go/utils"
    "io/ioutil"
    "math"
    "net/http"
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
        log.Fatal(err)
    }

    ctx := context.Background()
    opts := []core.ClientOption{
        option.WithWechatPayAutoAuthCipher(cfg.MchID, cfg.MchCertificateSerialNumber, mchPrivateKey, cfg.MchAPIv3Key),
    }
    client, err := core.NewClient(ctx, opts...)
    if err != nil {
        log.Fatal(err)
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
func (c *wechatClient) AppPay(prepay *model.OrderCache) (string, error) {
    cfg := ar.Config.Payment.Wechat

    svc := app.AppApiService{
        Client: c.Client,
    }

    resp, result, err := svc.PrepayWithRequestPayment(context.Background(), app.PrepayRequest{
        Appid:       core.String(cfg.AppID),
        Mchid:       core.String(cfg.MchID),
        Description: core.String(prepay.Name),
        OutTradeNo:  core.String(prepay.OutTradeNo),
        TimeExpire:  core.Time(prepay.Expire),
        NotifyUrl:   core.String(cfg.NotifyUrl),
        Amount: &app.Amount{
            Currency: core.String("CNY"),
            Total:    core.Int64(int64(math.Round(prepay.Amount * 100))),
        },
    })

    if err != nil {
        b, _ := ioutil.ReadAll(result.Response.Body)
        log.Errorf("微信支付调用失败: %#v, %s", err, string(b))
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

type WechatPayResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

// Notification 微信支付回调
func (c *wechatClient) Notification(req *http.Request) (*WechatPayResponse, *model.OrderCache) {
    transaction := new(payments.Transaction)
    _, err := c.notifyClient.ParseNotifyRequest(context.Background(), req, transaction)
    if err != nil {
        log.Error(err)
        return &WechatPayResponse{
            Code:    WechatPayCodeFail,
            Message: WechatPayCodeFail,
        }, nil
    }

    // 从缓存中获取订单数据
    out := transaction.OutTradeNo
    trade := new(model.OrderCache)
    err = ar.Cache.Get(context.Background(), "ORDER_"+*out).Scan(trade)
    if err != nil {
        log.Error(err)
    }
    trade.TradeNo = *(transaction.TransactionId)

    return &WechatPayResponse{
        Code:    WechatPayCodeSuccess,
        Message: WechatPayCodeSuccess,
    }, trade
}
