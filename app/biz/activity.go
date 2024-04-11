// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-08, by lisicen

package biz

import (
	"context"
	"errors"
	"strings"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/internal/ent/activity"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/tools"
)

type activityBiz struct {
	orm *ent.ActivityClient
	ctx context.Context
}

func NewActivity() *activityBiz {
	return &activityBiz{
		orm: ent.Database.Activity,
		ctx: context.Background(),
	}
}

func (a *activityBiz) Detail(id uint64) (*definition.ActivityDetail, error) {
	item, err := a.orm.Query().Where(activity.ID(id)).First(a.ctx)
	if err != nil {
		return nil, err
	}
	return toActivityDetail(item), nil
}

func (a *activityBiz) Create(req *definition.ActivityCreateReq) error {
	ctx := context.Background()
	_, err := a.orm.Create().
		SetName(req.Name).
		SetSort(req.Sort).
		SetImage(req.Image).
		SetLink(req.Link).
		SetNillablePopup(req.Popup).
		SetNillableHome(req.Home).
		SetNillableRemark(req.Remark).
		SetNillableEnable(req.Enable).
		SetIntroduction(req.Introduction).
		Save(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("请勿重复添加")
		}
		return err
	}
	return nil
}

func (a *activityBiz) Modify(req *definition.ActivityModifyReq) error {
	_, err := a.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetSort(req.Sort).
		SetImage(req.Image).
		SetLink(req.Link).
		SetNillablePopup(req.Popup).
		SetNillableHome(req.Home).
		SetNillableRemark(req.Remark).
		SetNillableEnable(req.Enable).
		SetIntroduction(req.Introduction).
		Save(a.ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("请勿重复添加")
		}
		return err
	}
	return nil
}

func (a *activityBiz) All() ([]*definition.ActivityDetail, error) {
	items, err := a.orm.Query().
		Where(activity.EnableEQ(true)).
		Order(ent.Desc(activity.FieldSort)).All(a.ctx)
	if err != nil {
		return nil, err
	}
	var res []*definition.ActivityDetail
	for _, item := range items {
		res = append(res, toActivityDetail(item))
	}
	return res, nil
}

func (a *activityBiz) List(req *definition.ActivityListReq) *model.PaginationRes {
	query := a.orm.Query().Order(ent.Desc(activity.FieldSort))
	if req.Keyword != nil {
		query.Where(activity.NameContains(*req.Keyword))
	}

	if req.Enable != nil {
		query.Where(activity.EnableEQ(*req.Enable))
	}

	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		query.Where(
			activity.UpdatedAtGTE(start),
			activity.UpdatedAtLTE(end),
		)
	}

	return model.ParsePaginationResponse(query, req.PaginationReq, func(item *ent.Activity) *definition.ActivityDetail {
		return toActivityDetail(item)
	})
}

func toActivityDetail(item *ent.Activity) *definition.ActivityDetail {
	return &definition.ActivityDetail{
		ID: item.ID,
		Activity: definition.Activity{
			Name:         item.Name,
			Sort:         item.Sort,
			Image:        item.Image,
			Link:         item.Link,
			Remark:       item.Remark,
			Popup:        &item.Popup,
			Home:         &item.Home,
			Enable:       item.Enable,
			UpdatedAt:    item.UpdatedAt.Format(carbon.DateTimeLayout),
			Introduction: item.Introduction,
		},
	}
}
