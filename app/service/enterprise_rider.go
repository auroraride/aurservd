// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/enterpriseprice"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribealter"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type enterpriseRiderService struct {
	ctx      context.Context
	modifier *model.Modifier
	rider    *ent.Rider
	agent    *ent.Agent
}

func NewEnterpriseRider() *enterpriseRiderService {
	return &enterpriseRiderService{
		ctx: context.Background(),
	}
}

func NewEnterpriseRiderWithRider(r *ent.Rider) *enterpriseRiderService {
	s := NewEnterpriseRider()
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, r)
	s.rider = r
	return s
}

func NewEnterpriseRiderWithModifier(m *model.Modifier) *enterpriseRiderService {
	s := NewEnterpriseRider()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
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
				SetIntelligent(ep.Intelligent).
				SetRemaining(0).
				SetInitialDays(req.Days).
				SetStatus(model.SubscribeStatusInactive).
				SetCityID(ep.CityID).
				// 团签骑手无须签合同 (2022-10-25)
				SetNeedContract(false).
				SetEnterpriseID(req.EnterpriseID).
				SetStationID(req.StationID).
				SetAgentEndAt(tools.NewTime().WillEnd(time.Now(), req.Days)).
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

// CreateByAgent 代理商小程序新增骑手
func (s *enterpriseRiderService) CreateByAgent(req *model.EnterpriseRiderCreateReq, ag *ent.Agent, sts ent.EnterpriseStations) {
	req.EnterpriseID = ag.EnterpriseID
	riderInfo, _ := ent.Database.Rider.Query().Where(rider.Phone(req.Phone)).First(s.ctx)

	// 判断代理商是否有该站点
	for _, v := range sts {
		if v.ID == req.StationID {
			break
		}
		snag.Panic("代理商没有绑定该站点")
	}

	if riderInfo != nil {
		// 查询订阅信息
		subscribeInfo, _ := NewSubscribe().QueryEffective(riderInfo.ID)
		if subscribeInfo != nil || riderInfo.EnterpriseID != nil {
			snag.Panic("该骑手不能绑定,已有团签或者已有未完成的订单")
		}
		// 更新rider
		if ent.Database.Rider.UpdateOne(riderInfo).
			SetEnterpriseID(req.EnterpriseID).
			SetStationID(req.StationID).
			Exec(s.ctx) != nil {
			snag.Panic("更新骑手失败")
		}
		return
	}

	// 创建rider
	if ent.Database.Rider.Create().SetPhone(req.Phone).
		SetEnterpriseID(req.EnterpriseID).
		SetStationID(req.StationID).
		SetName(req.Name).
		Exec(s.ctx) != nil {
		snag.Panic("创建骑手失败")
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
	case 2:
		q.Where(rider.DeletedAtIsNil())
	}

	// 筛选订阅状态
	switch req.SubscribeStatus {
	case 1:
		q.Where(rider.HasSubscribesWith(subscribe.Status(model.SubscribeStatusUsing)))
	case 2:
		q.Where(rider.HasSubscribesWith(subscribe.Status(model.SubscribeStatusUnSubscribed)))
	case 3:
		q.Where(
			rider.Or(
				rider.HasSubscribesWith(subscribe.Status(model.SubscribeStatusInactive)),
				rider.Not(rider.HasSubscribes()),
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
		m := item.Model
		if item.Intelligent {
			m += "「智能」"
		}
		res[i] = m
	}
	return res
}

// ChooseBatteryModel 团签选择电池型号
func (s *enterpriseRiderService) ChooseBatteryModel(req *model.EnterpriseRiderSubscribeChooseReq) model.EnterpriseRiderSubscribeChooseRes {
	e, _ := s.rider.QueryEnterprise().First(s.ctx)
	if e == nil {
		snag.Panic("非团签骑手")
	}
	if e.Agent {
		snag.Panic("代理骑手无法使用该功能")
	}

	ep, _ := ent.Database.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.EnterpriseID(e.ID), enterpriseprice.Model(req.Model)).First(s.ctx)
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
			SetNeedContract(true).
			SetIntelligent(ep.Intelligent).
			Save(s.ctx)
	} else {
		sub, err = sub.Update().SetIntelligent(ep.Intelligent).SetModel(ep.Model).Save(s.ctx)
		if err != nil {
			return model.EnterpriseRiderSubscribeChooseRes{}
		}
	}

	if err != nil {
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
		if time.Since(now).Seconds() >= 30 {
			return false
		}
		time.Sleep(1 * time.Second)
	}
}

// AddSubscribeDays 骑手申请增加团签订阅时长
func (s *enterpriseRiderService) AddSubscribeDays(req *model.RiderSubscribeAddReq, rid *ent.Rider) {
	// 查询骑手申请是否有未审批的
	info, _ := ent.Database.SubscribeAlter.QueryNotDeleted().Where(subscribealter.RiderID(rid.ID), subscribealter.StatusNEQ(model.SubscribeAlterUnreviewed)).First(s.ctx)
	if info != nil {
		snag.Panic("存在未审批的申请")
	}

	// 查询骑手是否有团签
	first, err := ent.Database.Rider.Query().Where(rider.ID(rid.ID), rider.HasSubscribesWith(subscribe.StatusIn(model.SubscribeNotUnSubscribed()...))).WithSubscribes().First(s.ctx)
	if err != nil {
		snag.Panic("团签状态异常")
	}
	sid := first.Edges.Subscribes[0].ID
	// 增加记录
	_, err = ent.Database.SubscribeAlter.Create().
		SetRiderID(rid.ID).
		SetEnterpriseID(*rid.EnterpriseID).
		SetSubscribeID(sid).
		SetDays(req.Days).
		SetStatus(model.SubscribeAlterUnreviewed).
		Save(s.ctx)
	if err != nil {
		snag.Panic("申请失败")
	}
}

// SubscribeAlterList 申请列表
func (s *enterpriseRiderService) SubscribeAlterList(req *model.SubscribeAlterApplyReq, rid *ent.Rider) *model.PaginationRes {
	q := ent.Database.SubscribeAlter.QueryNotDeleted().Where(subscribealter.RiderID(rid.ID)).Order(ent.Desc(subscribealter.FieldCreatedAt))

	if req.Start != nil && req.End != nil {
		rs := tools.NewTime().ParseDateStringX(*req.Start)
		re := tools.NewTime().ParseDateStringX(*req.End)
		q = q.Where(subscribealter.CreatedAtGTE(rs), subscribealter.CreatedAtLTE(re))
	}

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.SubscribeAlter) model.SubscribeAlterApplyListRsp {
			return model.SubscribeAlterApplyListRsp{
				ID:         item.ID,
				Days:       item.Days,
				ApplyTime:  item.CreatedAt.Format(carbon.DateTimeLayout), // 申请时间
				ReviewTime: item.UpdatedAt.Format(carbon.DateTimeLayout), // 审批时间
				Status:     item.Status,                                  // 审批状态
			}
		})
}

// RiderEnterpriseInfo 骑手团签信息
func (s *enterpriseRiderService) RiderEnterpriseInfo(req *model.EnterproseInfoReq, riderID uint64) *model.EnterproseInfoRsp {
	rsp := &model.EnterproseInfoRsp{
		IsJoin: true,
	}
	// 查询订阅信息
	sub, _ := NewSubscribe().QueryEffective(riderID)
	if sub == nil {
		rsp.IsJoin = false
	}
	// 查询团签信息
	enterpriseInfo := NewEnterprise().QueryX(req.EnterpriseId)
	if enterpriseInfo == nil {
		snag.Panic("未找到企业信息")
	}
	// 查询站点信息
	stationInfo := ent.Database.EnterpriseStation.Query().Where(enterprisestation.IDEQ(req.StationId),
		enterprisestation.EnterpriseIDEQ(req.EnterpriseId)).FirstX(s.ctx)
	if stationInfo == nil {
		snag.Panic("未找到站点信息")
	}
	rsp.StationName = stationInfo.Name
	rsp.EnterproseName = enterpriseInfo.Name
	return rsp
}

// JoinEnterprise 加入团签
func (s *enterpriseRiderService) JoinEnterprise(req *model.EnterproseInfoReq, rid *ent.Rider) {
	// 判断团签是否存在或者站点是否存在
	// 查询团签信息
	if NewEnterprise().QueryX(req.EnterpriseId) == nil {
		snag.Panic("未找到企业信息")
	}
	// 查询站点信息
	if ent.Database.EnterpriseStation.Query().Where(enterprisestation.IDEQ(req.StationId),
		enterprisestation.EnterpriseIDEQ(req.EnterpriseId)).FirstX(s.ctx) == nil {
		snag.Panic("未找到站点信息")
	}
	// 判断骑手是否有未完成的订单
	// 查询订阅信息
	sub, _ := NewSubscribe().QueryEffective(rid.ID)
	if sub != nil {
		snag.Panic("有未完成的订单")
	}
	_, err := ent.Database.Rider.Update().Where(rider.ID(rid.ID)).
		SetEnterpriseID(req.EnterpriseId).
		SetStationID(req.StationId).
		Save(s.ctx)
	if err != nil {
		snag.Panic("加入团签失败")
	}
}
