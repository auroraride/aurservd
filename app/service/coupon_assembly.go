// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-28
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/couponassembly"
	"github.com/golang-module/carbon/v2"
)

type couponAssemblyService struct {
	ctx          context.Context
	modifier     *model.Modifier
	rider        *ent.Rider
	employee     *ent.Employee
	employeeInfo *model.Employee
	orm          *ent.CouponAssemblyClient
}

func NewCouponAssembly() *couponAssemblyService {
	return &couponAssemblyService{
		ctx: context.Background(),
		orm: ent.Database.CouponAssembly,
	}
}

func NewCouponAssemblyWithModifier(m *model.Modifier) *couponAssemblyService {
	s := NewCouponAssembly()
	s.ctx = context.WithValue(s.ctx, "modifier", m)
	s.modifier = m
	return s
}

// List 获取发券记录
func (s *couponAssemblyService) List(req *model.CouponAssemblyListReq) *model.PaginationRes {
	q := s.orm.Query().Order(ent.Desc(couponassembly.FieldCreatedAt))
	if req.TemplateID != 0 {
		q.Where(couponassembly.TemplateID(req.TemplateID))
	}
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.CouponAssembly) model.CouponAssembly {
		return model.CouponAssembly{
			ID:      item.ID,
			Time:    item.CreatedAt.Format(carbon.DateTimeLayout),
			Remark:  item.Remark,
			Creator: item.Creator,
			Amount:  item.Amount,
			Target:  item.Target,
			Number:  item.Number,
			Name:    item.Name,
			Meta: model.CouponTemplateMeta{
				CouponTemplate: model.CouponTemplate{
					Rule:           item.Meta.Rule,
					CouponDuration: item.Meta.CouponDuration,
					Multiple:       item.Meta.Multiple,
				},
				Cities: item.Meta.Cities,
				Plans:  item.Meta.Plans,
			},
		}
	})
}
