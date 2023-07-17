// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
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
	"github.com/auroraride/aurservd/pkg/silk"
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

	stat := NewEnterpriseStation().QueryX(req.StationID)

	var (
		r   *ent.Rider
		sub *ent.Subscribe
	)

	r, _ = ent.Database.Rider.QueryNotDeleted().Where(rider.Phone(req.Phone)).First(s.ctx)

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		// 查询骑手和团签
		// 若骑手存在则删除原骑手信息并且新增骑手
		if r != nil {
			// 查询订阅信息
			sub, _ = NewSubscribe().QueryEffective(r.ID)
			if sub != nil && sub.EnterpriseID == nil {
				return errors.New("该骑手不能绑定,已有未完成的订单")
			}

			if sub.EnterpriseID != nil && sub.Status != model.SubscribeStatusInactive {
				err = tx.Subscribe.DeleteOne(sub).Exec(s.ctx)
				if err != nil {
					return errors.New("原未激活订阅处理失败")
				}
			}

			// 删除之前骑手信息
			_, err = tx.Rider.SoftDeleteOne(r).SetRemark("更改团签").Save(s.ctx)
			if err != nil {
				return errors.New("更改骑手失败")
			}

			// 删除旧的骑手登录信息
			NewRider().Signout(r)

			// 新增骑手信息
			r, err = s.CopyAndCreateRider(tx, r, &model.RiderConvert{
				EnterpriseID:     &req.EnterpriseID,
				StationID:        &req.StationID,
				Remark:           "代理转化骑手",
				Name:             req.Name,
				JoinEnterpriseAt: silk.Pointer(time.Now()),
			})
			if err != nil {
				return errors.New("转化骑手失败")
			}

		} else {
			// 未存在骑手创建骑手 并创建团签订阅信息
			var per *ent.Person
			// 创建person
			per, err = tx.Person.Create().SetName(req.Name).Save(s.ctx)
			if err != nil {
				return err
			}

			// 创建rider
			r, err = tx.Rider.Create().SetPhone(req.Phone).
				SetEnterpriseID(req.EnterpriseID).
				SetStationID(req.StationID).
				SetPerson(per).
				SetJoinEnterpriseAt(time.Now()).
				SetName(req.Name).
				Save(s.ctx)
			if err != nil {
				return err
			}
		}
		// 创建订阅信息
		err = tx.Subscribe.Create().
			SetRiderID(r.ID).
			SetModel(ep.Model).
			SetNillableBrandID(ep.BrandID).
			SetIntelligent(ep.Intelligent).
			SetRemaining(0).
			SetInitialDays(req.Days).
			SetStatus(model.SubscribeStatusInactive).
			SetCityID(ep.CityID).
			// 团签骑手无须签合同
			SetNeedContract(false).
			SetEnterpriseID(req.EnterpriseID).
			SetStationID(req.StationID).
			Exec(s.ctx)
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
			sq.Order(ent.Desc(subscribe.FieldCreatedAt))
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
func (s *enterpriseRiderService) BatteryModels(req *model.EnterprisePricePlanListReq) []string {
	if s.rider.EnterpriseID == nil {
		snag.Panic("非企业骑手")
	}

	items, _ := ent.Database.EnterprisePrice.QueryNotDeleted().
		Where(enterpriseprice.EnterpriseID(*s.rider.EnterpriseID), enterpriseprice.CityID(req.CityID), enterpriseprice.BrandIDIsNil()).
		All(s.ctx)

	res := make([]string, len(items))

	for i, item := range items {
		m := item.Model
		if item.Intelligent {
			// TODO 使用别的方式进行区分智能非智能而不是直接修改电池型号名称
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
		snag.Panic("请联系代理激活")
	}

	// TODO 使用别的方式进行区分智能非智能而不是直接修改电池型号名称
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

// RiderEnterpriseInfo 骑手团签信息
func (s *enterpriseRiderService) RiderEnterpriseInfo(req *model.EnterproseInfoReq, riderID uint64) *model.EnterproseInfoRsp {
	rsp := &model.EnterproseInfoRsp{
		IsJoin: true,
	}
	// 查询订阅信息
	sub, _ := NewSubscribe().QueryEffective(riderID)
	if sub != nil && (sub.EnterpriseID == nil || (sub.EnterpriseID != nil && sub.Status != model.SubscribeStatusInactive)) {
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

	price := NewEnterprisePrice().PriceList(req.EnterpriseId)
	rsp.PriceList = price
	if enterpriseInfo.Days != nil {
		rsp.Days = enterpriseInfo.Days
	}
	return rsp
}

// JoinEnterprise 加入团签
// TODO: 和 Create 方法实现的功能基本类似
func (s *enterpriseRiderService) JoinEnterprise(req *model.EnterpriseJoinReq, rid *ent.Rider) {
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
	if sub != nil && sub.EnterpriseID == nil {
		snag.Panic("有未完成的订单")
	}

	ep, _ := ent.Database.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.EnterpriseID(req.EnterpriseId), enterpriseprice.ID(req.PriceID)).First(s.ctx)
	if ep == nil {
		snag.Panic("未找到价格信息")
	}

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		// 若原有团签订阅未激活，直接删除
		if sub != nil && sub.EnterpriseID != nil && sub.Status != model.SubscribeStatusInactive {
			err = tx.Subscribe.DeleteOne(sub).Exec(s.ctx)
			if err != nil {
				return errors.New("原未激活订阅处理失败")
			}
		}

		_, err = ent.Database.Rider.Update().Where(rider.ID(rid.ID)).
			SetEnterpriseID(req.EnterpriseId).
			SetStationID(req.StationId).
			SetJoinEnterpriseAt(time.Now()).
			Save(s.ctx)
		if err != nil {
			snag.Panic("加入团签失败")
		}

		return tx.Subscribe.Create().
			SetRiderID(rid.ID).
			SetModel(ep.Model).
			SetNillableBrandID(ep.BrandID).
			SetIntelligent(ep.Intelligent).
			SetRemaining(0).
			SetInitialDays(req.Days).
			SetStatus(model.SubscribeStatusInactive).
			SetCityID(ep.CityID).
			SetNeedContract(false).
			SetEnterpriseID(req.EnterpriseId).
			SetStationID(req.StationId).
			Exec(s.ctx)
	})

}

// ExitEnterprise 退出团签
func (s *enterpriseRiderService) ExitEnterprise(r *ent.Rider) {
	if r.EnterpriseID == nil {
		snag.Panic("非团签骑手")
	}

	// 查询订阅
	sub, _ := NewSubscribe().QueryEffective(r.ID)

	if sub != nil && sub.Status != model.SubscribeStatusInactive {
		snag.Panic("骑士卡使用中无法转换")
	}

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		// 删除之前骑手信息
		_, err = tx.Rider.SoftDeleteOne(r).SetRemark("骑手退出团签").Save(s.ctx)
		if err != nil {
			return errors.New("删除骑手失败")
		}

		// 新增骑手信息
		_, err = s.CopyAndCreateRider(tx, r, &model.RiderConvert{
			Remark: "骑手退出团签",
		})

		if err != nil {
			return
		}

		// 删除未激活的订阅信息
		if sub != nil {
			err = tx.Subscribe.DeleteOne(sub).Exec(s.ctx)
		}

		return
	})

	// 删除旧的骑手登录信息
	NewRider().Signout(r)
}

// CopyAndCreateRider 复制并创建骑手信息
func (s *enterpriseRiderService) CopyAndCreateRider(tx *ent.Tx, r *ent.Rider, params *model.RiderConvert) (*ent.Rider, error) {
	name := params.Name
	if r.Name != "" {
		name = r.Name
	}
	return tx.Rider.Create().
		SetRemark(params.Remark).
		SetPhone(r.Phone).
		SetContact(r.Contact).
		SetDeviceType(r.DeviceType).
		SetLastDevice(r.LastDevice).
		SetIsNewDevice(r.IsNewDevice).
		SetNillableLastFace(r.LastFace).
		SetPushID(r.PushID).
		SetNillableLastSigninAt(r.LastSigninAt).
		SetBlocked(r.Blocked).
		SetNillablePersonID(r.PersonID).
		SetPoints(r.Points).
		SetName(name).
		SetIDCardNumber(r.IDCardNumber).
		SetExchangeLimit(r.ExchangeLimit).
		SetExchangeFrequency(r.ExchangeFrequency).
		SetNillableEnterpriseID(params.EnterpriseID).
		SetNillableStationID(params.StationID).
		SetNillableJoinEnterpriseAt(params.JoinEnterpriseAt).
		Save(s.ctx)
}
