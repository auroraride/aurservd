// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package payment

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/auroraride/adapter/log"
	"github.com/golang-module/carbon/v2"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
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
	err = client.LoadAppCertPublicKeyFromFile(cfg.AppPublicCert) // 加载应用公钥证书
	if err != nil {
		snag.Panic(err)
	}
	err = client.LoadAliPayRootCertFromFile(cfg.RootCert) // 加载支付宝根证书
	if err != nil {
		snag.Panic(err)
	}
	err = client.LoadAlipayCertPublicKeyFromFile(cfg.PublicCert) // 加载支付宝公钥证书
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

	err := c.Client.LoadAppCertPublicKeyFromFile(cfg.AppPublicCert) // 加载应用公钥证书
	if err != nil {
		snag.Panic(err)
	}
	err = c.Client.LoadAliPayRootCertFromFile(cfg.RootCert) // 加载支付宝根证书
	if err != nil {
		snag.Panic(err)
	}
	err = c.Client.LoadAlipayCertPublicKeyFromFile(cfg.PublicCert) // 加载支付宝公钥证书
	if err != nil {
		snag.Panic(err)
	}
}

// AppPay app支付
func (c *alipayClient) AppPay(pc *model.PaymentCache) (string, error) {
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

func (c *alipayClient) AppPayDemo() (string, string, error) {
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

func (c *alipayClient) Native(pc *model.PaymentCache) (string, error) {
	cfg := ar.Config.Payment.Alipay

	c.loadCerts()

	amount, subject, no, _ := pc.GetPaymentArgs()
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

	return res.QRCode, nil
}

// Refund 退款
func (c *alipayClient) Refund(req *model.PaymentRefund) {
	result, err := c.Client.TradeRefund(alipay.TradeRefund{
		TradeNo:      req.TradeNo,
		OutRequestNo: req.OutRefundNo,
		RefundAmount: fmt.Sprintf("%.2f", req.RefundAmount),
		RefundReason: req.Reason,
	})
	zap.L().Info(req.TradeNo+": 支付宝退款反馈", log.JsonData(result), zap.Error(err))

	if err != nil {
		snag.Panic("退款处理失败")
		return
	}
	if !result.Code.IsSuccess() {
		snag.Panic(result.Msg)
		return
	}

	req.Request = true
	if result.FundChange == "Y" {
		req.Success = true
		req.Time = time.Now()
	}
}

// Notification 支付宝回调
func (c *alipayClient) Notification(req *http.Request) *model.PaymentCache {
	err := req.ParseForm()
	if err != nil {
		zap.L().Error("支付宝回调失败", zap.Error(err))
		return nil
	}

	var result *alipay.Notification
	result, err = c.Client.DecodeNotification(req.Form)
	zap.L().Info("支付宝回调", log.JsonData(result), zap.Error(err))
	if err != nil {
		return nil
	}

	// 从缓存中获取订单数据
	pc := new(model.PaymentCache)
	out := result.OutTradeNo
	err = cache.Get(context.Background(), out).Scan(pc)
	if err != nil {
		zap.L().Error("从缓存获取订单信息失败", zap.Error(err))
		return nil
	}

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
		pc.Refund.Time = carbon.Parse(result.GmtRefund).ToStdTime()
		return pc
	}

	return nil
}
