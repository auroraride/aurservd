// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/enterpriseprice"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "time"
)

type enterpriseRiderService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
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
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// Create 新增骑手
func (s *enterpriseRiderService) Create(req *model.EnterpriseRiderCreateReq) model.EnterpriseRider {
    // 查询是否存在
    if ar.Ent.Rider.QueryNotDeleted().Where(rider.Phone(req.Phone)).ExistX(s.ctx) {
        snag.Panic("此手机号已存在")
    }

    stat := NewEnterpriseStation().Query(req.StationID)

    tx, _ := ar.Ent.Tx(s.ctx)
    // 创建person
    p, err := tx.Person.Create().SetName(req.Name).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 创建rider
    var r *ent.Rider
    r, err = tx.Rider.Create().SetPhone(req.Phone).SetEnterpriseID(req.EnterpriseID).SetStationID(req.StationID).SetPerson(p).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _ = tx.Commit()

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
    q := ar.Ent.Rider.
        QueryNotDeleted().
        WithPerson().
        WithSubscribes(func(sq *ent.SubscribeQuery) {
            sq.Where(subscribe.StartAtNotNil()).Order(ent.Desc(subscribe.FieldCreatedAt))
        }).
        WithStation().
        Where(rider.EnterpriseID(req.EnterpriseID))
    if req.Keyword != nil {
        q.Where(
            rider.Or(
                rider.HasPersonWith(person.NameContainsFold(*req.Keyword)),
                rider.PhoneContainsFold(*req.Keyword),
            ),
        )
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
            p := item.Edges.Person
            res := model.EnterpriseRider{
                ID:        item.ID,
                Phone:     item.Phone,
                CreatedAt: item.CreatedAt.Format(carbon.DateTimeLayout),
                Station: model.EnterpriseStation{
                    ID:   item.Edges.Station.ID,
                    Name: item.Edges.Station.Name,
                },
            }
            if p != nil {
                res.Name = p.Name
            }
            if item.Edges.Subscribes != nil {
                for i, sub := range item.Edges.Subscribes {
                    var days int
                    if i == 0 {
                        res.Voltage = sub.Voltage
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

                    days = tt.DiffDaysOfNextDay(after, before)

                    // 总天数
                    res.Days += days
                    // 判断是否已结算
                    if sub.StatementID == nil {
                        res.Unsettled += days
                    }
                }
            }
            return res
        },
    )
}

// ListVoltage 列出可用电压型号
func (s *enterpriseRiderService) ListVoltage(req *model.EnterprisePriceVoltageListReq) []model.EnterprisePriceVoltageListRes {
    if s.rider.EnterpriseID == nil {
        snag.Panic("非企业骑手")
    }
    items, _ := ar.Ent.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.EnterpriseID(*s.rider.EnterpriseID), enterpriseprice.CityID(req.CityID)).All(s.ctx)
    res := make([]model.EnterprisePriceVoltageListRes, len(items))
    for i, item := range items {
        res[i] = model.EnterprisePriceVoltageListRes{
            Voltage: item.Voltage,
            ID:      item.ID,
        }
    }
    return res
}

// ChooseVoltage 选择电池
func (s *enterpriseRiderService) ChooseVoltage(req *model.EnterpriseRiderSubscribeChooseReq) model.EnterpriseRiderSubscribeChooseRes {
    enterpriseID := s.rider.EnterpriseID
    if enterpriseID == nil {
        snag.Panic("非企业骑手")
    }
    ep, _ := ar.Ent.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.EnterpriseID(*enterpriseID)).First(s.ctx)
    if ep == nil {
        snag.Panic("未找到电池")
    }
    // 判断骑手是否有使用中的电池
    sub, _ := ar.Ent.Subscribe.QueryNotDeleted().Where(
        subscribe.EnterpriseID(*s.rider.EnterpriseID),
    ).Order(ent.Desc(subscribe.FieldCreatedAt)).First(s.ctx)
    if sub != nil && sub.StartAt != nil {
        snag.Panic("已被激活, 无法重新选择电池")
    }
    var err error
    if sub == nil {
        sub, err = ar.Ent.Subscribe.Create().
            SetEnterpriseID(*s.rider.EnterpriseID).
            SetStationID(*s.rider.StationID).
            SetVoltage(ep.Voltage).
            SetCityID(ep.CityID).
            SetRiderID(s.rider.ID).
            Save(s.ctx)
    } else {
        sub, err = sub.Update().SetVoltage(ep.Voltage).Save(s.ctx)
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
        sub, _ := ar.Ent.Subscribe.QueryNotDeleted().Where(
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