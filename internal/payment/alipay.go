// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

// FandAuthFreeze 线上资金授权冻结
func (c *alipayClient) FandAuthFreeze(pc *model.PaymentCache) (string, error) {
	cfg := ar.Config.Payment.AlipayAuthFreeze
	amount, subject, no, _ := pc.GetPaymentArgs()
	trade := FundAuthOrderAppFreeze{
		FundAuthOrderAppFreeze: alipay.FundAuthOrderAppFreeze{
			OutOrderNo:   no,
			OutRequestNo: no,
			OrderTitle:   subject,
			NotifyURL:    cfg.NotifyUrl,
			Amount:       fmt.Sprintf("%.2f", amount),
			ProductCode:  "PRE_AUTH_ONLINE",
			PayTimeout:   "20m",
		},
	}

	// 芝麻信用免押金
	if pc.CacheType == model.PaymentCacheTypeAliDepositFree {
		trade.ExtraParam = fmt.Sprintf(`{"category":"%s","serviceId":"%s"}`, cfg.Category, cfg.ServiceId)
		// 目前为纯免押
		trade.DepositProductMode = "DEPOSIT_ONLY"
	}

	res, err := c.FundAuthOrderAppFreeze(trade)
	return res, err
}

type FundAuthOrderAppFreeze struct {
	alipay.FundAuthOrderAppFreeze
	// 免押受理台模式，根据免押不同业务模式将开通受理台区分三种模式，商家可根据调用预授权冻结接口传入的参数决定该笔免押订单使用哪种受理台模式。不同受理台模式需要传入不同参数，其中：POSTPAY 表示后付金额已知，POSTPAY_UNCERTAIN 表示后付金额未知，DEPOSIT_ONLY 表示纯免押。
	DepositProductMode string `json:"deposit_product_mode,omitempty"`
}

func (c *alipayClient) FundAuthOrderAppFreeze(param FundAuthOrderAppFreeze) (result string, err error) {
	return c.EncodeParam(param)
}

// NotificationFandAuthFreeze 资金授权冻结回调
func (c *alipayClient) NotificationFandAuthFreeze(req *http.Request) *model.PaymentCache {
	err := req.ParseForm()
	if err != nil {
		zap.L().Error("资金授权冻结回调失败", zap.Error(err))
		return nil
	}

	var result *model.FandAuthFreezeNotification
	result, err = c.DecodeFandAuthFreezeNotification(req.Form)
	if err != nil {
		zap.L().Error("资金授权冻结回调解析失败", zap.Error(err))
		return nil
	}
	// 从缓存中获取订单数据
	pc := new(model.PaymentCache)
	out := result.OutOrderNo
	err = cache.Get(context.Background(), out).Scan(pc)
	if err != nil {
		zap.L().Error("从缓存获取订单信息失败", zap.Error(err))
		return nil
	}

	if result.Status == model.FandAuthFreezeStatusSuccess {
		switch pc.CacheType {
		case model.PaymentCacheTypePlan:
			// 预授权支付
			pc.Subscribe.TradeNo = result.OutOrderNo
		case model.PaymentCacheTypeAliDepositFree:
			// 芝麻信用免押金
			pc.DepositFree.TradeNo = result.OutOrderNo
		default:
			return nil
		}
		return pc
	}
	return nil
}

// DecodeFandAuthFreezeNotification 解析线上资金授权冻结回调
func (c *alipayClient) DecodeFandAuthFreezeNotification(values url.Values) (notification *model.FandAuthFreezeNotification, err error) {
	if err = c.VerifySign(values); err != nil {
		return nil, err
	}

	notification = &model.FandAuthFreezeNotification{}
	notification.AuthNo = values.Get("auth_no")
	notification.NotifyType = values.Get("notify_type")
	notification.FundAuthFreeze = values.Get("fund_auth_freeze")
	notification.OutOrderNo = values.Get("out_order_no")
	notification.OperationID = values.Get("operation_id")
	notification.OutRequestNo = values.Get("out_request_no")
	notification.OperationType = values.Get("operation_type")
	notification.Amount = values.Get("amount")
	notification.Status = model.FandAuthFreezeStatus(values.Get("status"))
	notification.GmtCreate = values.Get("gmt_create")
	notification.GmtTrans = values.Get("gmt_trans")
	notification.PayerLogonID = values.Get("payer_logon_id")
	notification.PayerUserID = values.Get("payer_user_id")
	notification.PayeeLogonID = values.Get("payee_logon_id")
	notification.PayeeUserID = values.Get("payee_user_id")
	notification.TotalFreezeAmount = values.Get("total_freeze_amount")
	notification.TotalUnfreezeAmount = values.Get("total_unfreeze_amount")
	notification.TotalPayAmount = values.Get("total_pay_amount")
	notification.RestAmount = values.Get("rest_amount")
	notification.CreditAmount = values.Get("credit_amount")
	notification.FundAmount = values.Get("fund_amount")
	notification.TotalFreezeCreditAmount = values.Get("total_freeze_credit_amount")
	notification.TotalFreezeFundAmount = values.Get("total_freeze_fund_amount")
	notification.TotalUnfreezeCreditAmount = values.Get("total_unfreeze_credit_amount")
	notification.TotalUnfreezeFundAmount = values.Get("total_unfreeze_fund_amount")
	notification.TotalPayCreditAmount = values.Get("total_pay_credit_amount")
	notification.TotalPayFundAmount = values.Get("total_pay_fund_amount")
	notification.RestCreditAmount = values.Get("rest_credit_amount")
	notification.RestFundAmount = values.Get("rest_fund_amount")
	notification.PreAuthType = values.Get("pre_auth_type")
	notification.CreditMerchantExt = values.Get("credit_merchant_ext")
	return notification, nil

}

// FandAuthUnfreeze 资金授权解冻
func (c *alipayClient) FandAuthUnfreeze(req *model.PaymentRefund) {
	trade := alipay.FundAuthOrderUnfreeze{
		AuthNo:       req.TradeNo,
		OutRequestNo: req.OutRefundNo,
		Amount:       fmt.Sprintf("%.2f", req.RefundAmount),
		Remark:       req.Reason,
	}

	res, err := c.FundAuthOrderUnfreeze(trade)
	if err != nil || !res.IsSuccess() {
		snag.Panic("资金授权解冻失败")
	}

	req.Request = true
	req.Success = true
	req.Time = time.Now()
}

// AlipayTradePay 资金冻结转支付
func (c *alipayClient) AlipayTradePay(req *model.TradePay) (res *alipay.TradePayRsp, err error) {
	cfg := ar.Config.Payment.AlipayAuthFreeze

	// 如果下单时指定了免押受理台模式则必填
	businessParams := map[string]string{
		"deduction_subject": "DEPOSIT",
	}
	jsonBusinessParams, err := json.Marshal(businessParams)
	if err != nil {
		return nil, err
	}

	trade := alipay.TradePay{
		AuthNo: req.AuthNo,
		Trade: alipay.Trade{
			OutTradeNo:     req.OutTradeNo,
			TotalAmount:    fmt.Sprintf("%.2f", req.TotalAmount),
			Subject:        req.Subject,
			ProductCode:    "PRE_AUTH_ONLINE",
			NotifyURL:      cfg.NotifyUrl,
			BusinessParams: jsonBusinessParams,
		},
		AuthConfirmMode: "COMPLETE", // COMPLETE：转交易支付完成结束预授权

	}

	res, err = c.TradePay(trade)
	if err != nil {
		return nil, err
	}

	if !res.IsSuccess() {
		return nil, res.Error
	}

	return res, nil
}

// 扣款完成回调通知
func (c *alipayClient) NotificationTradePay(req *http.Request) *model.PaymentCache {
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

// AlipayFundAuthOperationDetailQuery 查询查询芝麻免押订单
func (c *alipayClient) AlipayFundAuthOperationDetailQuery(outOrderNo, outRequestNo string) (result *alipay.FundAuthOperationDetailQueryRsp, err error) {
	trade := alipay.FundAuthOperationDetailQuery{
		OutOrderNo:   outOrderNo,
		OutRequestNo: outRequestNo,
	}
	result, err = c.FundAuthOperationDetailQuery(trade)
	return result, err
}

// AlipayFundAuthOperationCancelRequest 取消芝麻免押订单 订单为以下状态时可以取消订单：INIT（初始化）、AUTHORIZED（已创建）（此时一般为用户取消服务时使用）。
func (c *alipayClient) AlipayFundAuthOperationCancelRequest(outOrderNo, outRequestNo string) (result *alipay.FundAuthOperationCancelRsp, err error) {
	trade := alipay.FundAuthOperationCancel{
		OutOrderNo:   outOrderNo,
		OutRequestNo: outRequestNo,
		Remark:       "用户取消",
	}
	result, err = c.FundAuthOperationCancel(trade)
	return result, err
}
