// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-01
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/agent"
    "github.com/auroraride/aurservd/internal/ent/enterprisecontract"
    "github.com/auroraride/aurservd/internal/ent/enterprisestation"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/golang-module/carbon/v2"
    "github.com/rs/xid"
    "time"
)

type agentService struct {
    cacheKeyPrefix string
    ctx            context.Context
    modifier       *model.Modifier
    agent          *ent.Agent
    enterprise     *ent.Enterprise
    orm            *ent.AgentClient
}

func NewAgent() *agentService {
    return &agentService{
        cacheKeyPrefix: "AGENT_",
        ctx:            context.Background(),
        orm:            ent.Database.Agent,
    }
}

func NewAgentWithAgent(ag *ent.Agent, en *ent.Enterprise) *agentService {
    s := NewAgent()
    s.agent = ag
    s.enterprise = en
    return s
}

func NewAgentWithModifier(m *model.Modifier) *agentService {
    s := NewAgent()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
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
    pass, _ := utils.PasswordGenerate(req.Password)
    _, err := s.orm.Create().
        SetName(req.Name).
        SetPhone(req.Phone).
        SetPassword(pass).
        SetEnterpriseID(req.EnterpriseID).
        Save(s.ctx)
    if err != nil {
        snag.Panic("添加失败")
    }
}

func (s *agentService) Modify(req *model.AgentModifyReq) {
    up := s.QueryX(req.ID).Update()
    if req.Password != "" {
        pass, _ := utils.PasswordGenerate(req.Password)
        up.SetPassword(pass)
    }
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
    // res []model.AgentListRes
    items, _ := s.orm.QueryNotDeleted().Where(agent.EnterpriseID(req.EnterpriseID)).All(s.ctx)

    out := make([]model.AgentListRes, len(items))
    for i, item := range items {
        out[i] = model.AgentListRes{
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

func (s *agentService) tokenKey(id uint64) string {
    return fmt.Sprintf("%s%d", s.cacheKeyPrefix, id)
}

func (s *agentService) Signin(req *model.AgentSigninReq) model.AgentSigninRes {
    ag, _ := s.orm.QueryNotDeleted().Where(agent.Phone(req.Phone)).WithEnterprise().First(s.ctx)
    if ag.Edges.Enterprise == nil {
        snag.Panic("登录失败")
    }
    en := ag.Edges.Enterprise
    if !en.Agent {
        snag.Panic("非代理商")
    }
    if ag == nil || !utils.PasswordCompare(req.Password, ag.Password) {
        snag.Panic("账号或密码错误")
    }

    // 生成token
    token := xid.New().String() + utils.RandTokenString()
    key := s.tokenKey(ag.ID)

    // 删除旧的token
    if old := cache.Get(s.ctx, key).Val(); old != "" {
        cache.Del(s.ctx, key)
        cache.Del(s.ctx, old)
    }

    s.ExtendTokenTime(ag.ID, token)

    return model.AgentSigninRes{
        Profile: s.Profile(ag, en),
        Token:   token,
    }
}

// ExtendTokenTime 延长登录有效期
func (s *agentService) ExtendTokenTime(id uint64, token string) {
    ctx := context.Background()
    cache.Set(ctx, s.tokenKey(id), token, 7*24*time.Hour)
    cache.Set(ctx, token, id, 7*24*time.Hour)
}

func (s *agentService) Profile(ag *ent.Agent, en *ent.Enterprise) model.AgentProfile {

    // 查询合同
    today := carbon.Now().StartOfDay().Carbon2Time()
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
            subscribe.EndAtGTE(carbon.Now().StartOfDay().Carbon2Time()),
        ),
    ).Count(s.ctx)

    yt := carbon.Yesterday().StartOfDay().Carbon2Time()
    td := carbon.Now().StartOfDay().Carbon2Time()
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
        cost = tools.NewDecimal().Sum(pm[srv.PriceKey(sub.CityID, sub.Model)].Price, cost)
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
        ID:        ag.ID,
        Phone:     ag.Phone,
        Name:      ag.Name,
        Contract:  cf,
        Balance:   en.Balance,
        Riders:    riders,
        Yesterday: cost,
        Using:     using,
        Stations:  stations,
        Prices:    prices,
        Days:      en.Days,
    }
}
