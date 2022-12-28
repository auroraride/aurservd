// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package payment

import (
    "context"
    "errors"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/goccy/go-json"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "github.com/smartwalle/alipay/v3"
    "net/http"
    "time"
)

var _alipay *alipayClient

type alipayClient struct {
    *alipay.Client
}

func newAlipayClient() *alipay.Client {
    cfg := ar.Config.Payment.Alipay
    client, err := alipay.New(cfg.Appid, cfg.PrivateKey, true)
    if err != nil {
        snag.Panic(err)
    }
    err = client.LoadAppPublicCertFromFile(cfg.AppPublicCert) // 加载应用公钥证书
    if err != nil {
        snag.Panic(err)
    }
    err = client.LoadAliPayRootCertFromFile(cfg.RootCert) // 加载支付宝根证书
    if err != nil {
        snag.Panic(err)
    }
    err = client.LoadAliPayPublicCertFromFile(cfg.PublicCert) // 加载支付宝公钥证书
    if err != nil {
        snag.Panic(err)
    }
    return client
}

func NewAlipay() *alipayClient {
    return _alipay
}

func (c *alipayClient) loadCerts() {
    cfg := ar.Config.Payment.Alipay

    err := c.Client.LoadAppPublicCertFromFile(cfg.AppPublicCert) // 加载应用公钥证书
    if err != nil {
        snag.Panic(err)
    }
    err = c.Client.LoadAliPayRootCertFromFile(cfg.RootCert) // 加载支付宝根证书
    if err != nil {
        snag.Panic(err)
    }
    err = c.Client.LoadAliPayPublicCertFromFile(cfg.PublicCert) // 加载支付宝公钥证书
    if err != nil {
        snag.Panic(err)
    }
}

// AppPay app支付
func (c *alipayClient) AppPay(pc *model.PaymentCache) (string, error) {
    cfg := ar.Config.Payment.Alipay
    amount, subject, no := pc.GetPaymentArgs()
    trade := alipay.TradeAppPay{
        Trade: alipay.Trade{
            TotalAmount: fmt.Sprintf("%.2f", amount),
            NotifyURL:   cfg.NotifyUrl,
            Subject:     subject,
            OutTradeNo:  no,
        },
        TimeExpire: time.Now().Add(10 * time.Minute).Format(carbon.DateTimeLayout),
    }
    return c.TradeAppPay(trade)
}

func (c *alipayClient) AppPayDemo() (string, string, error) {
    cfg := ar.Config.Payment.Alipay
    no := tools.NewUnique().NewSN()
    trade := alipay.TradeAppPay{
        Trade: alipay.Trade{
            TotalAmount: fmt.Sprintf("%.2f", 0.01),
            NotifyURL:   cfg.NotifyUrl,
            Subject:     "测试支付",
            OutTradeNo:  no,
        },
        TimeExpire: time.Now().Add(10 * time.Minute).Format(carbon.DateTimeLayout),
    }

    s, err := c.TradeAppPay(trade)
    return s, no, err
}

func (c *alipayClient) Native(pc *model.PaymentCache) (string, error) {
    cfg := ar.Config.Payment.Alipay

    c.loadCerts()

    amount, subject, no := pc.GetPaymentArgs()
    trade := alipay.TradePreCreate{
        Trade: alipay.Trade{
            TotalAmount: fmt.Sprintf("%.2f", amount),
            NotifyURL:   cfg.NotifyUrl,
            Subject:     subject,
            OutTradeNo:  no,
        },
    }
    res, err := c.TradePreCreate(trade)
    if err != nil {
        return "", err
    }

    if !res.IsSuccess() {
        return "", errors.New("支付宝二维码生成失败")
    }

    return res.Content.QRCode, nil
}

// Refund 退款
func (c *alipayClient) Refund(req *model.PaymentRefund) {
    result, err := c.Client.TradeRefund(alipay.TradeRefund{
        TradeNo:      req.TradeNo,
        OutRequestNo: req.OutRefundNo,
        RefundAmount: fmt.Sprintf("%.2f", req.RefundAmount),
        RefundReason: req.Reason,
    })
    b, _ := json.MarshalIndent(result, "", "  ")
    log.Infof("[%s]支付宝退款反馈\n%s, err: %v", req.TradeNo, b, err)

    if err != nil {
        snag.Panic("退款处理失败")
        return
    }
    if !result.Content.Code.IsSuccess() {
        snag.Panic(result.Content.Msg)
        return
    }

    req.Request = true
    if result.Content.FundChange == "Y" {
        req.Success = true
        req.Time = carbon.Parse(result.Content.GmtRefundPay).Carbon2Time()
    }
    return
}

// Notification 支付宝回调
func (c *alipayClient) Notification(req *http.Request) *model.PaymentCache {
    result, err := c.Client.GetTradeNotification(req)
    b, _ := json.MarshalIndent(result, "", "  ")
    log.Infof("支付宝反馈\n%s", b)
    if err != nil {
        log.Error(err)
        return nil
    }

    // 从缓存中获取订单数据
    pc := new(model.PaymentCache)
    out := result.OutTradeNo
    err = cache.Get(context.Background(), out).Scan(pc)
    if err != nil {
        log.Errorf("从缓存获取订单信息失败: %v", err)
        return nil
    }

    b, _ = json.MarshalIndent(pc, "", "  ")
    log.Infof("获取到支付宝支付缓存: %s", b)

    switch pc.CacheType {
    case model.PaymentCacheTypePlan:
        if result.TradeStatus == alipay.TradeStatusSuccess {
            pc.Subscribe.TradeNo = result.TradeNo
        } else {
            return nil
        }
        return pc
    case model.PaymentCacheTypeOverdueFee:
        if result.TradeStatus == alipay.TradeStatusSuccess {
            pc.OverDueFee.TradeNo = result.TradeNo
        } else {
            return nil
        }
        return pc
    case model.PaymentCacheTypeAssistance:
        if result.TradeStatus == alipay.TradeStatusSuccess {
            pc.Assistance.TradeNo = result.TradeNo
        } else {
            return nil
        }
        return pc
    case model.PaymentCacheTypeRefund:
        pc.Refund.Success = true
        pc.Refund.Request = true
        pc.Refund.Time = carbon.Parse(result.GmtRefund).Carbon2Time()
        return pc
    }

    return nil
}
