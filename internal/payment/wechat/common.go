package wechat

import (
	"context"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/auroraride/adapter/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type commonClient struct {
	*core.Client
	notifyClient *notify.Handler
}

func newCommonClient(client *core.Client) *commonClient {
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(ar.Config.Payment.Wechat.MchID)
	notifyClient, err := notify.NewRSANotifyHandler(ar.Config.Payment.Wechat.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certVisitor))
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	return &commonClient{
		Client:       client,
		notifyClient: notifyClient,
	}
}

// NewWechatClientWithConfig 初始化微信支付客户端
func NewWechatClientWithConfig(cfg ar.WechatpayConfig) *core.Client {
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
	return client
}

// Refund 退款
func (c *commonClient) Refund(req *model.PaymentRefund) {
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
func (c *commonClient) Notification(req *http.Request) *model.PaymentCache {
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
func (c *commonClient) RefundNotification(req *http.Request) *model.PaymentCache {
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

// PurchaseNotification 购买微信支付回调
func (c *commonClient) PurchaseNotification(req *http.Request) *model.PurchaseNotificationRes {
	transaction := new(payments.Transaction)
	nq, err := c.notifyClient.ParseNotifyRequest(context.Background(), req, transaction)
	if err != nil {
		zap.L().Error("微信回调解析失败", zap.Error(err))
		return nil
	}

	b, _ := jsoniter.Marshal(transaction)
	nb, _ := jsoniter.Marshal(nq)
	zap.L().Info("微信支付回调", zap.ByteString("transaction", b), zap.ByteString("notifiy", nb))
	state := transaction.TradeState
	if *state != "SUCCESS" {
		return nil
	}
	out := transaction.OutTradeNo
	if out == nil {
		return nil
	}
	tranID := transaction.TransactionId
	if tranID == nil {
		return nil
	}
	return &model.PurchaseNotificationRes{
		OutTradeNo: *out,
		TradeNo:    *tranID,
		Amount:     tools.NewDecimal().Mul(float64(*transaction.Amount.Total), 0.01),
		Payway:     "wechat",
	}
}
