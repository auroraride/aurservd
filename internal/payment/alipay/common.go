package alipay

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

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type commonClient struct {
	*alipay.Client
}

func newCommonClient(client *alipay.Client) *commonClient {
	return &commonClient{
		Client: client,
	}
}

func newAlipayClientWithConfig(cfg ar.AlipayConfig) *alipay.Client {
	client, err := alipay.New(cfg.Appid, cfg.PrivateKey, true)
	if err != nil {
		snag.Panic(err)
	}

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
	if cfg.EncryptKey != "" {
		err = client.SetEncryptKey(cfg.EncryptKey)
		if err != nil {
			snag.Panic(err)
		}
	}

	return client
}

// Refund 退款
func (c *commonClient) Refund(req *model.PaymentRefund) {
	result, err := c.Client.TradeRefund(context.Background(), alipay.TradeRefund{
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
func (c *commonClient) Notification(req *http.Request) *model.PaymentCache {
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
	if result.OutRequestNo != "" {
		out = result.OutRequestNo
	}
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
		pc.Refund.Time = carbon.Parse(result.GmtRefund).StdTime()
		return pc
	case model.PaymentCacheTypeDeposit:
		if result.TradeStatus == alipay.TradeStatusSuccess {
			pc.DepositCredit.TradeNo = result.TradeNo
		} else {
			return nil
		}
		return pc
	default:
		snag.Panic("unhandled default case")
	}
	return nil
}

// FandAuthFreeze 线上资金授权冻结
func (c *commonClient) FandAuthFreeze(pc *model.PaymentCache) (string, error) {
	cfg := ar.Config.Payment.AlipayAuthFreeze
	var no, subject, amount string

	switch {
	case pc.Subscribe != nil:
		no = pc.Subscribe.OutOrderNo
		subject = pc.Subscribe.Name
		amount = fmt.Sprintf("%.2f", pc.Subscribe.Amount)
	case pc.DepositCredit != nil:
		no = pc.DepositCredit.OutOrderNo
		subject = "押金"
		amount = fmt.Sprintf("%.2f", pc.DepositCredit.Amount)
	}

	trade := FundAuthOrderAppFreeze{
		FundAuthOrderAppFreeze: alipay.FundAuthOrderAppFreeze{
			OutOrderNo:   no,
			OutRequestNo: no,
			OrderTitle:   subject,
			NotifyURL:    cfg.FreezeNotifyUrl,
			Amount:       amount,
			ProductCode:  "PRE_AUTH_ONLINE",
			PayTimeout:   "20m",
		},
	}

	extraParam := definition.ExtraParam{}

	// 预授权信用免押金
	if pc.CacheType == model.PaymentCacheTypeDeposit {
		// 目前为纯免押
		trade.DepositProductMode = "DEPOSIT_ONLY"
	} else if pc.CacheType == model.PaymentCacheTypeAlipayAuthFreeze {

		// 当有押金的时候
		if pc.Subscribe.Deposit > 0 {
			trade.DepositProductMode = "POSTPAY"
			extraParam.CreditExtInfo = &definition.CreditExtInfo{
				// 租押分离
				AssessmentAmount: fmt.Sprintf("%.2f", pc.Subscribe.Deposit),
			}
			postPayments := definition.PostPayments{
				Name:        "租金",
				Amount:      fmt.Sprintf("%.2f", pc.Subscribe.Amount-pc.Subscribe.Deposit),
				Description: fmt.Sprintf("%.2f", pc.Subscribe.Amount-pc.Subscribe.Deposit) + "元/月",
			}
			postPaymentsString, err := json.Marshal(postPayments)
			if err != nil {
				return "", err
			}
			trade.PostPayments = string(postPaymentsString)
		}
	}

	if pc.DepositCredit != nil || pc.Subscribe.Deposit > 0 {
		extraParam.Category = cfg.Category
		extraParam.ServiceId = cfg.ServiceId
		extraParamString, err := json.Marshal(extraParam)
		if err != nil {
			return "", err
		}
		trade.ExtraParam = string(extraParamString)
	}

	res, err := c.FundAuthOrderAppFreeze(trade)
	if err != nil {
		zap.L().Error("支付宝预授权冻结请求失败", zap.Error(err))
		return "", err
	}
	zap.L().Info("支付宝预授权冻结请求成功", log.JsonData(res))
	return res, nil
}

type FundAuthOrderAppFreeze struct {
	alipay.FundAuthOrderAppFreeze
	// 免押受理台模式，根据免押不同业务模式将开通受理台区分三种模式，商家可根据调用预授权冻结接口传入的参数决定该笔免押订单使用哪种受理台模式。不同受理台模式需要传入不同参数，其中：POSTPAY 表示后付金额已知，POSTPAY_UNCERTAIN 表示后付金额未知，DEPOSIT_ONLY 表示纯免押。
	DepositProductMode string `json:"deposit_product_mode,omitempty"`
	// 后付费项目，有付费项目时需要传入该字段。不同受理台模式需要传入不同参数，后付费项目名称和计费说明需要通过校验规则，同时计费说明将展示在开通受理台上。
	PostPayments string `json:"post_payments,omitempty"`
}

// FundAuthOrderAppFreeze 这里复写alipay.FundAuthOrderAppFreeze方法是为了解决alipay.FundAuthOrderAppFreeze 扩展字段未更新问题
func (c *commonClient) FundAuthOrderAppFreeze(param FundAuthOrderAppFreeze) (result string, err error) {
	return c.EncodeParam(param)
}

// NotificationFandAuthFreeze 资金授权冻结回调
func (c *commonClient) NotificationFandAuthFreeze(req *http.Request) *model.PaymentCache {
	err := req.ParseForm()
	if err != nil {
		zap.L().Error("资金授权冻结回调失败", zap.Error(err))
		return nil
	}

	var result *definition.FandAuthFreezeNotification
	result, err = c.DecodeFandAuthFreezeNotification(req.Form)
	if err != nil {
		zap.L().Error("资金授权冻结回调解析失败", zap.Error(err))
		return nil
	}

	if result.Status == definition.FandAuthFreezeStatusSuccess {
		switch result.NotifyType {
		case definition.FandAuthNotifyType:
			// 从缓存中获取订单数据
			pc := new(model.PaymentCache)
			out := result.OutOrderNo
			err = cache.Get(context.Background(), out).Scan(pc)
			if err != nil {
				zap.L().Error("从缓存获取订单信息失败", zap.Error(err))
				return nil
			}
			switch pc.CacheType {
			case model.PaymentCacheTypeAlipayAuthFreeze:
				// 预授权支付
				pc.Subscribe.AuthNo = result.AuthNo
				pc.Subscribe.OutRequestNo = result.OutRequestNo
			case model.PaymentCacheTypeDeposit:
				pc.DepositCredit.AuthNo = result.AuthNo
				pc.DepositCredit.OutRequestNo = result.OutRequestNo
			default:
				return nil
			}
			return pc
		}

	}

	zap.L().Info("资金授权冻结回调", log.JsonData(result))
	return nil
}

// DecodeFandAuthFreezeNotification 解析线上资金授权冻结回调
func (c *commonClient) DecodeFandAuthFreezeNotification(values url.Values) (notification *definition.FandAuthFreezeNotification, err error) {

	if err = c.VerifySign(values); err != nil {
		return nil, err
	}

	notification = &definition.FandAuthFreezeNotification{}
	notification.AuthNo = values.Get("auth_no")
	notification.NotifyType = values.Get("notify_type")
	notification.FundAuthFreeze = values.Get("fund_auth_freeze")
	notification.OutOrderNo = values.Get("out_order_no")
	notification.OperationID = values.Get("operation_id")
	notification.OutRequestNo = values.Get("out_request_no")
	notification.OperationType = values.Get("operation_type")
	notification.Amount = values.Get("amount")
	notification.Status = definition.FandAuthFreezeStatus(values.Get("status"))
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
func (c *commonClient) FandAuthUnfreeze(refund *model.PaymentRefund, req *definition.FandAuthUnfreezeReq) error {
	if req.Remark == "" {
		req.Remark = "申请退款"
	}
	trade := alipay.FundAuthOrderUnfreeze{
		AuthNo:       req.AuthNo,
		OutRequestNo: tools.NewUnique().NewSN28(),
		Amount:       fmt.Sprintf("%.2f", req.Amount),
		Remark:       req.Remark,
		NotifyURL:    ar.Config.Payment.AlipayAuthFreeze.UnfreezeNotifyUrl,
	}

	if req.IsDeposit {
		trade.ExtraParam = `{"unfreezeBizInfo":{"bizComplete":true}}`
	}

	res, err := c.FundAuthOrderUnfreeze(context.Background(), trade)
	if err != nil || !res.IsSuccess() {
		zap.L().Error("资金授权解冻失败", zap.Error(err), zap.Error(res.Error), log.JsonData(trade))
		return err
	}

	refund.Request = true
	refund.Success = true
	refund.Time = time.Now()

	return nil
}

// NotificationFandAuthUnfreeze 解冻回调
func (c *commonClient) NotificationFandAuthUnfreeze(req *http.Request) (res *model.PaymentCache) {
	res = new(model.PaymentCache)
	err := req.ParseForm()
	if err != nil {
		zap.L().Error("资金授权解冻回调失败", zap.Error(err))
		return
	}

	result, err := c.DecodeFandAuthFreezeNotification(req.Form)
	if err != nil {
		zap.L().Error("资金授权解冻回调解析失败", zap.Error(err))
		return
	}

	// 从缓存中获取订单数据
	pc := new(model.PaymentCache)
	out := result.OutOrderNo
	err = cache.Get(context.Background(), out).Scan(pc)
	if err != nil {
		zap.L().Error("从缓存获取订单信息失败", zap.Error(err))
		return
	}

	if result.Status == definition.FandAuthFreezeStatusSuccess {
		pc.Refund.Success = true
		pc.Refund.Time = time.Now()
	}
	return pc
}

// AlipayTradePay 资金冻结转支付
func (c *commonClient) AlipayTradePay(req *definition.TradePay) (res *alipay.TradePayRsp, err error) {
	cfg := ar.Config.Payment.AlipayAuthFreeze

	// 如果下单时指定了免押受理台模式则必填
	businessParams := map[string]string{
		"deduction_subject": "DEPOSIT",
	}
	jsonBusinessParams, err := json.Marshal(businessParams)
	if err != nil {
		return nil, err
	}

	if req.AuthConfirmMode == "" {
		req.AuthConfirmMode = "NOT_COMPLETE"
	}

	trade := alipay.TradePay{
		AuthNo: req.AuthNo,
		Trade: alipay.Trade{
			OutTradeNo:     req.OutTradeNo,
			TotalAmount:    fmt.Sprintf("%.2f", req.TotalAmount),
			Subject:        req.Subject,
			ProductCode:    "PRE_AUTH_ONLINE",
			NotifyURL:      cfg.TradePayNotifyUrl,
			BusinessParams: jsonBusinessParams,
		},
		AuthConfirmMode: req.AuthConfirmMode,
	}

	res, err = c.TradePay(context.Background(), trade)
	if err != nil {
		return nil, err
	}

	if !res.IsSuccess() {
		return nil, res.Error
	}
	zap.L().Info("资金冻结转支付成功", log.JsonData(res))

	return res, nil
}

// NotificationTradePay 扣款完成回调通知
func (c *commonClient) NotificationTradePay(req *http.Request) (res *definition.OrderDepositFreezeToPayRes) {
	res = new(definition.OrderDepositFreezeToPayRes)
	err := req.ParseForm()
	if err != nil {
		zap.L().Error("支付宝回调失败", zap.Error(err))
		return
	}

	var result *alipay.Notification
	result, err = c.Client.DecodeNotification(req.Form)
	zap.L().Info("支付宝回调", log.JsonData(result))
	if err != nil {
		zap.L().Error("支付宝回调解析失败", zap.Error(err))
		return
	}
	res.TradeNo = result.TradeNo
	res.OutTradeNo = result.OutTradeNo
	return res
}

// AlipayFundAuthOperationDetailQuery 查询预授权订单
func (c *commonClient) AlipayFundAuthOperationDetailQuery(req definition.FundAuthOperationDetailReq) (result *alipay.FundAuthOperationDetailQueryRsp, err error) {
	trade := alipay.FundAuthOperationDetailQuery{
		OutOrderNo:   req.OutOrderNo,
		OutRequestNo: req.OutRequestNo,
	}
	result, err = c.FundAuthOperationDetailQuery(context.Background(), trade)
	if err != nil {
		zap.L().Error("查询预授权订单失败", zap.Error(err), log.Payload(req))
		return nil, err
	}
	if result == nil {
		zap.L().Error("查询预授权订单失败，返回结果为空", log.Payload(req))
		return nil, errors.New("查询预授权订单失败")
	}
	if !result.Error.IsSuccess() {
		zap.L().Error("查询预授权订单失败", zap.Error(result.Error), log.Payload(req))
		return nil, result.Error
	}

	zap.L().Info("查询预授权订单成功", log.JsonData(result))
	return result, nil
}

// AlipayFundAuthOperationCancel 取消预授权订单 订单为以下状态时可以取消订单：INIT（初始化）、AUTHORIZED（已创建）（此时一般为用户取消服务时使用）。
func (c *commonClient) AlipayFundAuthOperationCancel(authNo, outRequestNo string) (result *alipay.FundAuthOperationCancelRsp, err error) {
	trade := alipay.FundAuthOperationCancel{
		AuthNo:       authNo,
		OutRequestNo: outRequestNo,
		Remark:       "用户取消",
	}
	result, err = c.FundAuthOperationCancel(context.Background(), trade)
	if !result.Error.IsSuccess() {
		zap.L().Error("取消预授权免押订单失败", zap.Error(result.Error))
		return nil, result.Error
	}
	zap.L().Info("取消预授权免押订单成功", log.JsonData(result))
	return result, err
}
