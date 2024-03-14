// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-08, by lisicen

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/activity"
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
	item, err := a.orm.Get(a.ctx, id)
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
		SetNillableIndex(req.Index).
		SetNillableRemark(req.Remark).
		Save(ctx)
	if err != nil {
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
		SetNillableIndex(req.Index).
		SetNillableRemark(req.Remark).
		Save(a.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *activityBiz) Delete(id uint64) error {
	err := a.orm.SoftDeleteOneID(id).Exec(a.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *activityBiz) All() ([]*definition.ActivityDetail, error) {
	items, err := a.orm.Query().Order(ent.Desc(activity.FieldSort)).All(a.ctx)
	if err != nil {
		return nil, err
	}
	var res []*definition.ActivityDetail
	for _, item := range items {
		res = append(res, toActivityDetail(item))
	}
	return res, nil
}

func (a *activityBiz) List(req *model.PaginationReq) *model.PaginationRes {
	query := a.orm.Query().Order(ent.Desc(activity.FieldSort))
	return model.ParsePaginationResponse(query, *req, func(item *ent.Activity) *definition.ActivityDetail {
		return toActivityDetail(item)
	})
}

func toActivityDetail(item *ent.Activity) *definition.ActivityDetail {
	return &definition.ActivityDetail{
		ID: item.ID,
		ActivityReqCommon: definition.ActivityReqCommon{
			Name:   item.Name,
			Sort:   item.Sort,
			Image:  item.Image,
			Link:   item.Link,
			Remark: &item.Remark,
			Popup:  &item.Popup,
			Index:  &item.Index,
		},
	}
}
