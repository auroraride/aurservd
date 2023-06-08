// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/enterpriseprepayment"
	"github.com/auroraride/aurservd/internal/payment"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type prepaymentService struct {
	*BaseService
	ctx context.Context
	orm *ent.EnterprisePrepaymentClient
}

func NewPrepayment(params ...any) *prepaymentService {
	return &prepaymentService{
		BaseService: newService(params...),
		ctx:         context.Background(),
		orm:         ent.Database.EnterprisePrepayment,
	}
}

func (s *prepaymentService) Overview(en *ent.Enterprise) (res model.PrepaymentOverview) {
	res.Balance = en.Balance
	var result []struct {
		EnterpriseID uint64  `json:"enterprise_id"`
		Amount       float64 `json:"amount"`
		Times        int     `json:"times"`
	}
	_ = ent.Database.EnterprisePrepayment.
		Query().
		Where(enterpriseprepayment.EnterpriseID(en.ID)).
		GroupBy(enterpriseprepayment.FieldEnterpriseID).
		Aggregate(
			ent.As(ent.Sum(enterpriseprepayment.FieldAmount), "amount"),
			ent.As(ent.Count(), "times"),
		).
		Scan(s.ctx, &result)
	if len(result) == 0 {
		return
	}
	res.Times = result[0].Times
	res.Amount = result[0].Amount
	res.Cost = tools.NewDecimal().Sub(res.Amount, res.Balance)
	return
}

func (s *prepaymentService) List(enterpriseID uint64, req *model.PrepaymentListReq) *model.PaginationRes {
	q := s.orm.Query().
		Where(enterpriseprepayment.EnterpriseID(enterpriseID)).
		WithAgent().
		Order(ent.Desc(enterpriseprepayment.FieldCreatedAt))

	// 筛选时间段
	if req.Start != "" {
		q.Where(enterpriseprepayment.CreatedAtGTE(tools.NewTime().ParseDateStringX(req.Start)))
	}
	if req.End != "" {
		q.Where(enterpriseprepayment.CreatedAtLT(tools.NewTime().ParseNextDateStringX(req.End)))
	}

	// 筛选支付方式
	if req.Payway != 0 {
		q.Where(enterpriseprepayment.Payway(req.Payway))
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.EnterprisePrepayment) model.PrepaymentListRes {
		res := model.PrepaymentListRes{
			Amount: item.Amount,
			Time:   item.CreatedAt.Format(carbon.DateTimeLayout),
			Remark: item.Remark,
			Payway: item.Payway,
			Name:   "-",
		}
		if item.Creator != nil {
			res.Name = item.Creator.Name + "-" + item.Creator.Phone
		}
		if item.Edges.Agent != nil {
			res.Name = item.Edges.Agent.Name + "-" + item.Edges.Agent.Phone
		}
		return res
	})
}

// WechatMiniprogramPay 小程序储值
func (s *prepaymentService) WechatMiniprogramPay(ag *ent.Agent, req *model.AgentPrepayReq) (data *model.AgentPrepayRes) {
	defer func() {
		go func() {
			b, _ := jsoniter.Marshal(data)
			zap.L().Info("代理商小程序充值", zap.ByteString("data", b))
		}()
	}()

	amount := req.Amount
	// Development模式支付一分钱
	if ar.Config.Environment.IsDevelopment() {
		amount = 0.01
	}

	pc := &model.PaymentCache{
		CacheType: model.PaymentCacheTypeAgentPrepay,
		AgentPrepay: &model.PaymentAgentPrepay{
			AgentPrepay: &model.AgentPrepay{
				EnterpriseID: ag.EnterpriseID,
				Remark:       "代理商小程序储值",
				Amount:       amount,
				ID:           ag.ID,
				Name:         ag.Name,
				Phone:        ag.Phone,
			},
			Payway:     model.PaywayAgentWxMiniprogram,
			OutTradeNo: tools.NewUnique().NewSN(),
		},
	}

	// 缓存支付信息
	err := cache.Set(s.ctx, pc.AgentPrepay.OutTradeNo, pc, 20*time.Minute).Err()
	if err != nil {
		snag.Panic(err)
		zap.L().Error("代理商小程序支付缓存失败", zap.Error(err))
	}

	// 生成预支付订单
	var res *jsapi.PrepayWithRequestPaymentResponse
	res, err = payment.NewWechat().Miniprogram(ar.Config.WechatMiniprogram.Agent.AppID, req.OpenID, pc)
	if err != nil {
		snag.Panic(err)
		zap.L().Error("代理商小程序支付订单生成失败", zap.Error(err))
	}

	// 生成随机字符串并签名
	data = &model.AgentPrepayRes{
		PrepayId:  res.PrepayId,
		Appid:     res.Appid,
		TimeStamp: res.TimeStamp,
		NonceStr:  res.NonceStr,
		Package:   res.Package,
		SignType:  res.SignType,
		PaySign:   res.PaySign,
	}

	return
}

// UpdateBalance 更新代理商余额
func (s *prepaymentService) UpdateBalance(payway model.Payway, req *model.AgentPrepay) (balance float64, err error) {
	defer func() {
		if err != nil {
			zap.L().Error("更新代理商金额失败", zap.Error(err))
		}
	}()

	// 获取团签
	var e *ent.Enterprise
	e, err = NewEnterprise().Query(req.EnterpriseID)
	if err != nil {
		return
	}

	// 获取当前账单信息
	set := NewEnterpriseStatement().Current(e)

	// 充值前余额
	before := e.Balance

	// 事务处理
	err = ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
		// 创建预付费记录
		creator := tx.EnterprisePrepayment.Create().SetEnterpriseID(req.EnterpriseID).SetAmount(req.Amount).SetRemark(req.Remark).SetPayway(payway)
		if req.ID > 0 {
			creator.SetAgentID(req.ID)
		}
		_, err = creator.Save(s.ctx)
		if err != nil {
			return
		}

		td := tools.NewDecimal()

		// 更新余额
		// 账单表
		balance = td.Sum(e.Balance, req.Amount)
		_, err = tx.EnterpriseStatement.UpdateOne(set).Save(s.ctx)
		if err != nil {
			return
		}

		// 更新企业表
		_, err = tx.Enterprise.UpdateOne(e).SetBalance(balance).AddPrepaymentTotal(req.Amount).Save(s.ctx)
		return
	})

	// 记录日志
	go func() {
		l := logging.NewOperateLog().
			SetRef(e).
			SetModifier(s.modifier).
			SetDiff(fmt.Sprintf("余额%.2f元", before), fmt.Sprintf("余额%.2f元", balance)).
			SetRemark(req.Remark)
		if s.modifier != nil {
			l.SetModifier(s.modifier).SetOperate(model.OperateEnterprisePrepayment)
		}
		if req.ID > 0 {
			l.SetOperate(model.OperateAgentPrepay).SetOperator(&logging.Operator{
				Type:  model.OperatorTypeAgent,
				ID:    req.ID,
				Phone: req.Phone,
				Name:  req.Name,
			})
		}
		l.Send()
	}()

	return
}
