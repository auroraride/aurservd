// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package payment

import (
    "context"
    "errors"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
    "github.com/wechatpay-apiv3/wechatpay-go/core"
    "github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
    "github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
    "github.com/wechatpay-apiv3/wechatpay-go/core/notify"
    "github.com/wechatpay-apiv3/wechatpay-go/core/option"
    "github.com/wechatpay-apiv3/wechatpay-go/services/payments"
    "github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
    "github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
    "github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
    "github.com/wechatpay-apiv3/wechatpay-go/utils"
    "io/ioutil"
    "math"
    "net/http"
    "time"
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

func (c *wechatClient) AppPayDemo() (string, string, error) {
    cfg := ar.Config.Payment.Wechat

    svc := app.AppApiService{
        Client: c.Client,
    }

    no := tools.NewUnique().NewSN28()

    resp, result, err := svc.PrepayWithRequestPayment(context.Background(), app.PrepayRequest{
        Appid:       core.String(cfg.AppID),
        Mchid:       core.String(cfg.MchID),
        Description: core.String("测试支付"),
        OutTradeNo:  core.String(no),
        TimeExpire:  core.Time(time.Now().Add(10 * time.Minute)),
        NotifyUrl:   core.String(cfg.NotifyUrl),
        Amount: &app.Amount{
            Currency: core.String("CNY"),
            Total:    core.Int64(int64(math.Round(1))),
        },
    })

    if err != nil {
        b, _ := ioutil.ReadAll(result.Response.Body)
        log.Errorf("微信支付调用失败: %#v, %s", err, string(b))
        return "", "", err
    }

    var out struct {
        *app.PrepayWithRequestPaymentResponse
        AppID string `json:"appId"`
    }

    out.PrepayWithRequestPaymentResponse = resp
    out.AppID = cfg.AppID

    b, _ := jsoniter.Marshal(out)

    return string(b), no, nil
}

// AppPay APP支付
func (c *wechatClient) AppPay(pc *model.PaymentCache) (string, error) {
    amount, subject, no := pc.GetPaymentArgs()
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

func (c *wechatClient) Native(pc *model.PaymentCache) (string, error) {
    amount, subject, no := pc.GetPaymentArgs()
    cfg := ar.Config.Payment.Wechat

    svc := native.NativeApiService{Client: c.Client}
    resp, result, err := svc.Prepay(context.Background(), native.PrepayRequest{
        Appid:       core.String(cfg.AppID),
        Mchid:       core.String(cfg.MchID),
        Description: core.String(subject),
        OutTradeNo:  core.String(no),
        TimeExpire:  core.Time(time.Now().Add(10 * time.Minute)),
        NotifyUrl:   core.String(cfg.NotifyUrl),
        Amount: &native.Amount{
            Currency: core.String("CNY"),
            Total:    core.Int64(int64(math.Round(amount * 100))),
        },
    })

    if err != nil {
        b, _ := ioutil.ReadAll(result.Response.Body)
        log.Errorf("微信支付调用失败: %#v, %s", err, string(b))
        return "", err
    }

    if resp.CodeUrl == nil {
        return "", errors.New("支付二维码获取失败")
    }

    return *resp.CodeUrl, nil
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
        b, _ := ioutil.ReadAll(result.Response.Body)
        log.Errorf("微信退款调用失败: %#v, %s", err, string(b))
        snag.Panic("退款处理失败")
    }

    req.Request = true

    s := resp.Status
    if *s == refunddomestic.STATUS_SUCCESS {
        req.Success = true
        req.Time = *resp.SuccessTime
    }
    return
}

// Notification 微信支付回调
func (c *wechatClient) Notification(req *http.Request) *model.PaymentCache {
    transaction := new(payments.Transaction)
    nq, err := c.notifyClient.ParseNotifyRequest(context.Background(), req, transaction)
    if err != nil {
        log.Error(err)
        return nil
    }

    b, _ := jsoniter.MarshalIndent(transaction, "", "  ")
    nb, _ := jsoniter.MarshalIndent(nq, "", "  ")
    log.Infof("微信支付回调反馈\n%s\n%s", b, nb)

    pc := new(model.PaymentCache)
    // 从缓存中获取订单数据
    out := transaction.OutTradeNo
    err = cache.Get(context.Background(), *out).Scan(pc)
    if err != nil {
        log.Errorf("从缓存获取订单信息失败: %v", err)
        return nil
    }

    b, _ = jsoniter.MarshalIndent(pc, "", "  ")
    log.Infof("获取到微信支付回调缓存: %s", b)

    state := transaction.TradeState
    if *state != "SUCCESS" {
        return nil
    }

    switch pc.CacheType {
    case model.PaymentCacheTypePlan:
        pc.Subscribe.TradeNo = *(transaction.TransactionId)
        break
    case model.PaymentCacheTypeOverdueFee:
        pc.OverDueFee.TradeNo = *(transaction.TransactionId)
        break
    case model.PaymentCacheTypeAssistance:
        pc.Assistance.TradeNo = *(transaction.TransactionId)
        break
    default:
        return nil
    }

    b, _ = jsoniter.MarshalIndent(pc, "", "  ")
    log.Infof("微信支付缓存更新: %s", b)

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
        log.Error(err)
        return nil
    }

    b, _ := jsoniter.MarshalIndent(transaction, "", "  ")
    nb, _ := jsoniter.MarshalIndent(nq, "", "  ")
    log.Infof("微信退款回调反馈\n%s\n%s", b, nb)

    pc := new(model.PaymentCache)

    // 从缓存中获取订单数据
    err = cache.Get(context.Background(), transaction.OutTradeNo).Scan(pc)
    if err != nil {
        log.Errorf("从缓存获取订单信息失败: %v", err)
        return nil
    }

    b, _ = jsoniter.MarshalIndent(pc, "", "  ")
    log.Infof("获取到微信退款回调缓存: %s", b)

    if transaction.RefundStatus != "SUCCESS" {
        return nil
    }

    pc.Refund.Success = true
    pc.Refund.Request = true
    pc.Refund.Time = transaction.SuccessTime

    b, _ = jsoniter.MarshalIndent(pc, "", "  ")
    log.Infof("微信退款缓存更新: %s", b)

    return pc
}
