// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-01
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"
	"strconv"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/enterprisecontract"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
	"github.com/auroraride/aurservd/pkg/utils"
)

type agentService struct {
	*BaseService

	tokenCacheKey string
	orm           *ent.AgentClient
}

func NewAgent(params ...any) *agentService {
	return &agentService{
		BaseService:   newService(params...),
		tokenCacheKey: ar.Config.Environment.UpperString() + ":" + "AGENT:TOKEN",
		orm:           ent.Database.Agent,
	}
}

func (s *agentService) Query(id uint64) (*ent.Agent, error) {
	return s.orm.QueryNotDeleted().Where(agent.ID(id)).WithEnterprise().First(s.ctx)
}

func (s *agentService) QueryX(id uint64) *ent.Agent {
	ag, _ := s.Query(id)
	if ag == nil {
		snag.Panic("未找到代理账号")
	}
	return ag
}

func (s *agentService) Create(req *model.AgentCreateReq) {
	en := NewEnterprise().QueryX(req.EnterpriseID)
	if !en.Agent {
		snag.Panic("团签模式非代理商")
	}
	_, err := s.orm.Create().
		SetName(req.Name).
		SetPhone(req.Phone).
		SetEnterpriseID(req.EnterpriseID).
		Save(s.ctx)
	if err != nil {
		snag.Panic("添加失败")
	}
}

func (s *agentService) Modify(req *model.AgentModifyReq) {
	up := s.QueryX(req.ID).Update()
	if req.Name != "" {
		up.SetName(req.Name)
	}
	if req.Phone != "" {
		up.SetPhone(req.Phone)
	}
	_, err := up.Save(s.ctx)
	if err != nil {
		snag.Panic("修改失败")
	}
}

func (s *agentService) Delete(req *model.IDParamReq) {
	ag := s.QueryX(req.ID)
	_, _ = s.orm.SoftDeleteOne(ag).Save(s.ctx)
}

func (s *agentService) List(req *model.AgentListReq) *model.PaginationRes {
	// res []model.AgentMeta
	items, _ := s.orm.QueryNotDeleted().Where(agent.EnterpriseID(req.EnterpriseID)).All(s.ctx)

	out := make([]model.AgentMeta, len(items))
	for i, item := range items {
		out[i] = model.AgentMeta{
			ID:    item.ID,
			Name:  item.Name,
			Phone: item.Phone,
		}
	}
	return &model.PaginationRes{
		Pagination: model.Pagination{
			Current: 1,
			Pages:   1,
			Total:   len(items),
		},
		Items: out,
	}
}

// TokenVerify Token校验
// 返回携带 ent.Enterprise 和 ent.EnterpriseStations
func (s *agentService) TokenVerify(token string) (ag *ent.Agent, en *ent.Enterprise, stations ent.EnterpriseStations) {
	// 获取token对应ID
	id, _ := ar.Redis.HGet(s.ctx, s.tokenCacheKey, token).Uint64()
	if id <= 0 {
		return
	}

	// 反向校验token是否正确
	if ar.Redis.HGet(s.ctx, s.tokenCacheKey, strconv.FormatUint(id, 10)).Val() != token {
		return
	}

	// 获取agent
	ag, _ = s.orm.QueryNotDeleted().Where(agent.ID(id)).WithEnterprise(func(eq *ent.EnterpriseQuery) {
		eq.WithStations()
	}).WithStations().First(s.ctx)
	if ag == nil || ag.Edges.Enterprise == nil {
		return
	}

	en = ag.Edges.Enterprise
	stations = ag.Edges.Stations

	// 如果是超级管理员，拥有所有站点
	if ag.Super {
		stations = en.Edges.Stations
	}

	return
}

// 代理登录, 返回代理端资料
// 需要关联查询 enterprise
func (s *agentService) signin(ag *ent.Agent) *model.AgentSigninRes {
	idstr := strconv.FormatUint(ag.ID, 10)
	// 查询并删除旧token key
	exists := ar.Redis.HGet(s.ctx, s.tokenCacheKey, idstr).Val()
	if exists != "" {
		ar.Redis.HDel(s.ctx, s.tokenCacheKey, exists)
	}

	// 生成token
	token := utils.NewEcdsaToken()

	// 存储登录token和ID进行对应
	ar.Redis.HSet(s.ctx, s.tokenCacheKey, token, ag.ID)
	ar.Redis.HSet(s.ctx, s.tokenCacheKey, idstr, token)

	return &model.AgentSigninRes{
		Profile: s.Profile(ag, ag.Edges.Enterprise),
		Token:   token,
	}
}

// Signin 登录
func (s *agentService) Signin(req *model.AgentSigninReq) *model.AgentSigninRes {
	switch req.SigninType {
	case model.SigninTypeSms:
		// 校验短信
		NewSms().VerifyCodeX(req.Phone, req.SmsId, req.Code)
	case model.SigninTypeAuth:
		// 获取手机号
		req.Phone = NewminiProgram().GetPhoneNumber(req.Code)
	}
	ag, err := s.orm.QueryNotDeleted().Where(agent.Phone(req.Phone)).WithEnterprise().First(s.ctx)
	if err != nil || ag.Edges.Enterprise == nil {
		snag.Panic("账号不存在")
	}
	en := ag.Edges.Enterprise
	if !en.Agent {
		snag.Panic("非代理商")
	}
	return s.signin(ag)
}

// Profile 代理商资料
func (s *agentService) Profile(ag *ent.Agent, en *ent.Enterprise) model.AgentProfile {
	// 查询合同
	today := carbon.Now().StartOfDay().StdTime()
	cr, _ := ent.Database.EnterpriseContract.QueryNotDeleted().
		Where(
			enterprisecontract.EnterpriseID(en.ID),
			enterprisecontract.StartLTE(today),
			enterprisecontract.EndGTE(today),
		).
		First(s.ctx)

	var cf string
	if cr != nil {
		cf = fmt.Sprintf("https://cdn.auroraride.com/%s", cr.File)
	}

	riders, _ := ent.Database.Rider.QueryNotDeleted().Where(rider.EnterpriseID(en.ID)).Count(s.ctx)
	using, _ := ent.Database.Subscribe.QueryNotDeleted().Where(
		subscribe.EnterpriseID(en.ID),
		subscribe.StartAtNotNil(),
		subscribe.Or(
			subscribe.EndAtIsNil(),
			subscribe.EndAtGTE(carbon.Now().StartOfDay().StdTime()),
		),
	).Count(s.ctx)

	yt := carbon.Yesterday().StartOfDay().StdTime()
	td := carbon.Now().StartOfDay().StdTime()
	subs, _ := ent.Database.Subscribe.QueryNotDeleted().Where(
		subscribe.EnterpriseID(en.ID),
		// 启用时间早于今日零点
		subscribe.StartAtLT(td),
		// 结束时间晚于昨日零点或未结束
		subscribe.Or(
			subscribe.EndAtIsNil(),
			subscribe.EndAtGTE(yt),
		),
	).All(s.ctx)

	srv := NewEnterprise()
	pm := srv.GetPrices(en)
	var cost float64
	for _, sub := range subs {
		cost = tools.NewDecimal().Sum(pm[srv.PriceKey(sub.CityID, sub.Model, sub.BrandID)].Price, cost)
	}

	prices := make([]model.EnterprisePrice, 0)
	for _, price := range pm {
		prices = append(prices, price)
	}

	// stations
	stations := make([]model.EnterpriseStation, 0)
	sitems, _ := ent.Database.EnterpriseStation.QueryNotDeleted().Where(enterprisestation.EnterpriseID(en.ID)).All(s.ctx)
	for _, sitem := range sitems {
		stations = append(stations, model.EnterpriseStation{
			ID:   sitem.ID,
			Name: sitem.Name,
		})
	}

	return model.AgentProfile{
		Enterprise: model.Enterprise{
			ID:    en.ID,
			Name:  en.Name,
			Agent: true,
		},
		ID:             ag.ID,
		Phone:          ag.Phone,
		Name:           ag.Name,
		Contract:       cf,
		Balance:        en.Balance,
		Riders:         riders,
		Yesterday:      cost,
		Using:          using,
		Stations:       stations,
		Prices:         prices,
		Days:           en.Days,
		RechargeAmount: en.RechargeAmount,
		Distance:       en.Distance,
		CompanyName:    en.CompanyName,
		ContactName:    en.ContactName,
		ContactPhone:   en.ContactPhone,
	}
}
