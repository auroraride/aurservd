// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package payment

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "github.com/smartwalle/alipay/v3"
    "net/http"
)

var _alipay *alipayClient

type alipayClient struct {
    *alipay.Client
    notifyClient *alipay.Client
}

func newAlipayClient() *alipay.Client {
    cfg := ar.Config.Payment.Alipay
    client, err := alipay.New(cfg.Appid, cfg.PrivateKey, true)
    if err != nil {
        snag.Panic(err)
    }
    return client
}

func newNotifyClient() *alipay.Client {
    cfg := ar.Config.Payment.Alipay
    client := newAlipayClient()
    var err error
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

// AppPay app支付
func (c *alipayClient) AppPay(prepay *model.OrderCache) (string, error) {
    cfg := ar.Config.Payment.Alipay
    trade := alipay.TradeAppPay{
        Trade: alipay.Trade{
            TotalAmount: fmt.Sprintf("%.2f", prepay.Amount),
            NotifyURL:   cfg.NotifyUrl,
            Subject:     prepay.Name,
            OutTradeNo:  prepay.OutTradeNo,
        },
        TimeExpire: prepay.Expire.Format(carbon.DateTimeLayout),
    }
    return c.TradeAppPay(trade)
}

// Notification 支付宝回调
func (c *alipayClient) Notification(req *http.Request, rep http.ResponseWriter) *model.OrderCache {
    result, err := c.notifyClient.GetTradeNotification(req)
    if err == nil {
        // 从缓存中获取订单数据
        out := result.OutTradeNo
        trade := new(model.OrderCache)
        err = cache.Get(context.Background(), "ORDER_"+out).Scan(trade)
        if err != nil {
            log.Error(err)
        }
        trade.TradeNo = result.TradeNo
        return trade
    } else {
        log.Error(err)
    }
    alipay.AckNotification(rep)
    return nil
}
