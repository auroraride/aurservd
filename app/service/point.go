// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-27
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/pointlog"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/golang-module/carbon/v2"
)

type pointService struct {
	ctx          context.Context
	modifier     *model.Modifier
	rider        *ent.Rider
	employee     *ent.Employee
	employeeInfo *model.Employee
	orm          *ent.PointLogClient
}

func NewPoint() *pointService {
	return &pointService{
		ctx: context.Background(),
		orm: ent.Database.PointLog,
	}
}

func NewPointWithModifier(m *model.Modifier) *pointService {
	s := NewPoint()
	s.ctx = context.WithValue(s.ctx, "modifier", m)
	s.modifier = m
	return s
}

// Modify 修改积分
func (s *pointService) Modify(req *model.PointModifyReq) error {
	r := NewRider().Query(req.RiderID)
	modify := req.Points
	if req.Type == model.PointLogTypeConsume && modify > 0 {
		modify = -modify
	}
	if req.Type == model.PointLogTypeAward && modify < 0 {
		modify = -modify
	}
	after := r.Points + modify
	if after < 0 {
		return errors.New("积分余额不能小于0")
	}
	return ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
		err = tx.Rider.UpdateOne(r).SetPoints(after).Exec(s.ctx)
		if err != nil {
			return
		}
		return tx.PointLog.Create().SetRiderID(req.RiderID).SetPoints(modify).SetReason(req.Reason).SetType(req.Type.Value()).SetAfter(after).Exec(s.ctx)
	})
}

// List 积分变动日志
func (s *pointService) List(req *model.PointLogListReq) *model.PaginationRes {
	q := s.orm.Query().Order(ent.Desc(pointlog.FieldCreatedAt)).WithRider()
	if req.RiderID == 0 {
		if req.Keyword != "" {
			q.Where(
				pointlog.HasRiderWith(rider.Or(
					rider.NameContainsFold(req.Keyword),
					rider.PhoneContainsFold(req.Keyword),
				)),
			)
		}
	} else {
		q.Where(pointlog.RiderID(req.RiderID))
	}
	if req.Type != 0 {
		q.Where(pointlog.Type(req.Type.Value()))
	}
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.PointLog) model.PointLogListRes {
		var mp, mm string
		if item.Attach != nil {
			if item.Attach.Plan != nil {
				mp = fmt.Sprintf("%s-%d天", item.Attach.Plan.Name, item.Attach.Plan.Days)
			}
		}

		if item.Modifier != nil {
			mm = item.Modifier.Name
		}

		if item.Edges.Rider != nil {
			mm = item.Edges.Rider.GetExportInfo()
		}

		return model.PointLogListRes{
			ID:       item.ID,
			Type:     model.PointLogType(item.Type).String(),
			Plan:     mp,
			Points:   item.Points,
			Reason:   item.Reason,
			After:    item.After,
			Modifier: mm,
			Time:     item.CreatedAt.Format(carbon.DateTimeLayout),
		}
	})
}

// Batch 批量发放积分
func (s *pointService) Batch(req *model.PointBatchReq) []string {
	riders, _, notfound := NewRider().QueryPhones(req.Phones)

	for i, p := range notfound {
		notfound[i] = fmt.Sprintf("未找到: %s", p)
	}

	for _, r := range riders {
		err := s.Modify(&model.PointModifyReq{
			RiderID: r.ID,
			Points:  req.Points,
			Reason:  req.Reason,
			Type:    req.Type,
		})
		if err != nil {
			notfound = append(notfound, fmt.Sprintf("%v: %s", err, r.Phone))
		}
	}

	return notfound
}

// Real 获取真实积分余额
// 真实积分 = 账户积分 - 预消耗
func (s *pointService) Real(r *ent.Rider) int64 {
	points := r.Points
	// 从缓存中获取
	x, _ := cache.Get(s.ctx, fmt.Sprintf("POINTS_%d", r.ID)).Int64()
	return points - x
}

// PreConsume 预消耗积分
func (s *pointService) PreConsume(r *ent.Rider, v int64) (last int64, err error) {
	points := r.Points
	// 从缓存中获取
	key := fmt.Sprintf("POINTS_%d", r.ID)
	x, _ := cache.Get(s.ctx, key).Int64()
	x += v
	if points-x < 0 {
		err = errors.New("积分余额不能为负")
		return
	}
	last = points - x
	err = cache.Set(s.ctx, key, last, 20*time.Minute).Err()
	return
}

// RemovePreConsume 移除预消耗的积分
func (s *pointService) RemovePreConsume(r *ent.Rider, v int64) {
	key := fmt.Sprintf("POINTS_%d", r.ID)
	x, _ := cache.Get(s.ctx, key).Int64()
	if x < v {
		return
	}
	x -= v
	ttl, _ := cache.TTL(s.ctx, key).Result()
	cache.Set(s.ctx, key, x, ttl)
}

func (s *pointService) Detail(r *ent.Rider) model.PointRes {
	points := r.Points
	x, _ := cache.Get(s.ctx, fmt.Sprintf("POINTS_%d", r.ID)).Int64()
	return model.PointRes{
		Points: points,
		Locked: x,
	}
}

// CalculateGift 计算赠送积分
func (s *pointService) CalculateGift(amount float64, cityID uint64) (points int64, proportion float64) {
	arr, err := GetSetting[[]model.SettingConsumePoint](model.SettingConsumePointKey)
	if err != nil {
		return
	}
	for _, set := range arr {
		if set.CityID == cityID {
			proportion = set.Proportion
			points = int64(math.Round(proportion * amount))
			if points < 0 {
				points = 0
			}
			return
		}
	}
	return
}
