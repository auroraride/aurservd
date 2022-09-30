// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/enterpriseprice"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "time"
)

type enterpriseRiderService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    agent        *ent.Agent
    enterprise   *ent.Enterprise
    employeeInfo *model.Employee
}

func NewEnterpriseRider() *enterpriseRiderService {
    return &enterpriseRiderService{
        ctx: context.Background(),
    }
}

func NewEnterpriseRiderWithRider(r *ent.Rider) *enterpriseRiderService {
    s := NewEnterpriseRider()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewEnterpriseRiderWithModifier(m *model.Modifier) *enterpriseRiderService {
    s := NewEnterpriseRider()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewEnterpriseRiderWithEmployee(e *ent.Employee) *enterpriseRiderService {
    s := NewEnterpriseRider()
    if e != nil {
        s.employee = e
        s.employeeInfo = &model.Employee{
            ID:    e.ID,
            Name:  e.Name,
            Phone: e.Phone,
        }
        s.ctx = context.WithValue(s.ctx, "employee", s.employeeInfo)
    }
    return s
}

func NewEnterpriseRiderWithAgent(ag *ent.Agent, en *ent.Enterprise) *enterpriseRiderService {
    s := NewEnterpriseRider()
    s.agent = ag
    s.enterprise = en
    return s
}

// Create 新增骑手
func (s *enterpriseRiderService) Create(req *model.EnterpriseRiderCreateReq) model.EnterpriseRider {
    if s.agent != nil && (req.EnterpriseID == 0 || req.StationID == 0) {
        snag.Panic("缺失团签信息")
    }

    // 查询团签
    e := NewEnterprise().QueryX(req.EnterpriseID)

    // 判定代理商字段
    var ep *ent.EnterprisePrice
    if e.Agent {
        if req.Days == 0 || req.PriceID == 0 {
            snag.Panic("代理商必选骑士卡信息")
        }
        ep, _ = ent.Database.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.EnterpriseID(e.ID), enterpriseprice.ID(req.PriceID)).First(s.ctx)
        if ep == nil {
            snag.Panic("未找到价格信息")
        }
    }

    // 查询是否存在
    if ent.Database.Rider.QueryNotDeleted().Where(rider.Phone(req.Phone)).ExistX(s.ctx) {
        snag.Panic("此手机号已存在")
    }

    stat := NewEnterpriseStation().Query(req.StationID)
    var r *ent.Rider

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        var per *ent.Person
        // 创建person
        per, err = tx.Person.Create().SetName(req.Name).Save(s.ctx)
        if err != nil {
            return
        }

        // 创建rider
        r, err = tx.Rider.Create().SetPhone(req.Phone).SetEnterpriseID(req.EnterpriseID).SetStationID(req.StationID).SetPerson(per).Save(s.ctx)
        if err != nil {
            return
        }

        // 如果是代理, 创建待激活骑士卡
        // 代理商添加骑手订阅强制为新签
        if e.Agent {
            _, err = tx.Subscribe.Create().
                SetType(model.OrderTypeNewly).
                SetRiderID(r.ID).
                SetModel(ep.Model).
                SetRemaining(0).
                SetInitialDays(req.Days).
                SetStatus(model.SubscribeStatusInactive).
                SetCityID(ep.CityID).
                SetEnterpriseID(req.EnterpriseID).
                SetStationID(req.StationID).
                Save(s.ctx)
        }
        return
    })

    // 记录日志
    go logging.NewOperateLog().
        SetRef(r).
        SetModifier(s.modifier).
        SetAgent(s.agent).
        SetOperate(model.OperateRiderCreate).
        Send()

    return model.EnterpriseRider{
        ID:        r.ID,
        Name:      req.Name,
        Phone:     req.Phone,
        CreatedAt: r.CreatedAt.Format(carbon.DateTimeLayout),
        Station: model.EnterpriseStation{
            ID:   stat.ID,
            Name: stat.Name,
        },
    }
}

// List 列举骑手
func (s *enterpriseRiderService) List(req *model.EnterpriseRiderListReq) *model.PaginationRes {
    q := ent.Database.Rider.
        Query().
        WithSubscribes(func(sq *ent.SubscribeQuery) {
            sq.Where(subscribe.StartAtNotNil()).Order(ent.Desc(subscribe.FieldCreatedAt))
        }).
        WithStation().
        Where(rider.EnterpriseID(req.EnterpriseID)).
        Order(ent.Desc(rider.FieldCreatedAt))
    if req.Keyword != nil {
        q.Where(
            rider.Or(
                rider.NameContainsFold(*req.Keyword),
                rider.PhoneContainsFold(*req.Keyword),
            ),
        )
    }

    // 筛选删除状态
    switch req.Deleted {
    case 1:
        q.Where(rider.DeletedAtNotNil())
        break
    case 2:
        q.Where(rider.DeletedAtIsNil())
        break
    }

    // 筛选订阅状态
    switch req.SubscribeStatus {
    case 1:
        q.Where(rider.HasSubscribesWith(subscribe.Status(model.SubscribeStatusUsing)))
        break
    case 2:
        q.Where(rider.HasSubscribesWith(subscribe.Status(model.SubscribeStatusUnSubscribed)))
        break
    case 3:
        q.Where(
            rider.Or(
                rider.HasSubscribesWith(subscribe.Status(model.SubscribeStatusInactive)),
                rider.Not(rider.HasSubscribes()),
            ),
        )
        break
    }
    tt := tools.NewTime()
    var rs, re time.Time
    if req.Start != nil {
        rs = tt.ParseDateStringX(*req.Start)
        q.Where(rider.HasSubscribesWith(subscribe.StartAtGTE(rs)))
    }
    if req.End != nil {
        re = tt.ParseDateStringX(*req.End)
        q.Where(rider.HasSubscribesWith(subscribe.StartAtLT(re.AddDate(0, 0, 1))))
    }
    return model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        func(item *ent.Rider) model.EnterpriseRider {
            res := model.EnterpriseRider{
                ID:        item.ID,
                Phone:     item.Phone,
                CreatedAt: item.CreatedAt.Format(carbon.DateTimeLayout),
                Name:      item.Name,
                Station: model.EnterpriseStation{
                    ID:   item.Edges.Station.ID,
                    Name: item.Edges.Station.Name,
                },
            }
            if item.Edges.Subscribes != nil {
                for i, sub := range item.Edges.Subscribes {
                    var days int
                    if i == 0 {
                        res.SubscribeStatus = sub.Status
                        res.Model = sub.Model
                    }
                    if sub.StartAt == nil {
                        continue
                    }
                    // 计算订阅使用天数
                    // 根据请求的时间范围计算时间周期
                    before := rs
                    after := re

                    // 如果请求日期为空或请求日期在开始日期之前
                    if before.IsZero() || before.Before(*sub.StartAt) {
                        before = *sub.StartAt
                    }

                    // 截止日期默认为当前日期或请求日期
                    if after.IsZero() {
                        after = time.Now()
                    }
                    // 如果订阅有结束日期并且结束日期在请求日期之前
                    if sub.EndAt != nil && after.After(*sub.EndAt) {
                        after = *sub.EndAt
                    }

                    days = tt.UsedDays(after, before)

                    // 总天数
                    res.Days += days
                    // 判断是否已结算
                    if sub.LastBillDate == nil {
                        res.Unsettled += days
                    } else {
                        res.Unsettled += tt.UsedDaysToNow(carbon.Time2Carbon(*sub.LastBillDate).StartOfDay().AddDay().Carbon2Time())
                    }
                }
            }
            // 已被删除
            if item.DeletedAt != nil {
                res.DeletedAt = item.DeletedAt.Format(carbon.DateTimeLayout)
            }
            return res
        },
    )
}

// BatteryModels 列出企业可用电压型号
func (s *enterpriseRiderService) BatteryModels(req *model.EnterprisePriceBatteryModelListReq) []string {
    if s.rider.EnterpriseID == nil {
        snag.Panic("非企业骑手")
    }

    items, _ := ent.Database.EnterprisePrice.QueryNotDeleted().
        Where(enterpriseprice.EnterpriseID(*s.rider.EnterpriseID), enterpriseprice.CityID(req.CityID)).
        All(s.ctx)

    res := make([]string, len(items))

    for i, item := range items {
        res[i] = item.Model
    }
    return res
}

// ChooseBatteryModel 选择电池型号
func (s *enterpriseRiderService) ChooseBatteryModel(req *model.EnterpriseRiderSubscribeChooseReq) model.EnterpriseRiderSubscribeChooseRes {
    // 查询骑手是否签约过
    if !NewContract().Effective(s.rider) {
        snag.Panic("请先签约")
    }

    enterpriseID := s.rider.EnterpriseID
    if enterpriseID == nil {
        snag.Panic("非企业骑手")
    }
    ep, _ := ent.Database.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.EnterpriseID(*enterpriseID), enterpriseprice.Model(req.Model)).First(s.ctx)
    if ep == nil {
        snag.Panic("未找到电池")
    }
    // 判断骑手是否有使用中的电池
    sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(
        subscribe.EnterpriseID(*s.rider.EnterpriseID),
        subscribe.RiderID(s.rider.ID),
        subscribe.StatusIn(model.SubscribeStatusUsing, model.SubscribeStatusInactive),
    ).Order(ent.Desc(subscribe.FieldCreatedAt)).First(s.ctx)
    if sub != nil && sub.StartAt != nil {
        snag.Panic("已被激活, 无法重新选择电池")
    }
    var err error
    if sub == nil {
        sub, err = ent.Database.Subscribe.Create().
            SetEnterpriseID(*s.rider.EnterpriseID).
            SetStationID(*s.rider.StationID).
            SetModel(ep.Model).
            SetCityID(ep.CityID).
            SetRiderID(s.rider.ID).
            Save(s.ctx)
    } else {
        sub, err = sub.Update().SetModel(ep.Model).Save(s.ctx)
        if err != nil {
            return model.EnterpriseRiderSubscribeChooseRes{}
        }
    }
    if err != nil {
        log.Error(err)
        snag.Panic("型号选择失败")
    }
    return model.EnterpriseRiderSubscribeChooseRes{
        Qrcode: fmt.Sprintf("SUBSCRIBE:%d", sub.ID),
    }
}

// SubscribeStatus 骑手激活电池状态
func (s *enterpriseRiderService) SubscribeStatus(req *model.EnterpriseRiderSubscribeStatusReq) bool {
    now := time.Now()
    for {
        sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(
            subscribe.ID(req.ID),
        ).First(s.ctx)
        if sub == nil {
            snag.Panic("未找到有效订阅")
        }
        if sub.StartAt != nil {
            return true
        }
        if time.Now().Sub(now).Seconds() >= 30 {
            return false
        }
        time.Sleep(1 * time.Second)
    }
}
