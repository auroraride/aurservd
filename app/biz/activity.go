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
	"github.com/auroraride/aurservd/pkg/snag"
)

type activityBiz struct {
	orm *ent.ActivityClient
}

func NewActivity() *activityBiz {
	return &activityBiz{
		orm: ent.Database.Activity,
	}
}

func (a *activityBiz) Get(id uint64) *definition.ActivityDetail {
	ctx := context.Background()
	item, err := a.orm.Get(ctx, id)
	if err != nil {
		snag.Panic(err)
	}
	return toActivityDetail(item)
}

func (a *activityBiz) Create(req *definition.ActivityCreateReq) {
	ctx := context.Background()
	_, err := a.orm.Create().
		SetName(req.Name).
		SetSort(req.Sort).
		SetImage(req.Image).
		SetLink(req.Link).
		SetRemark(req.Remark).
		Save(ctx)
	if err != nil {
		snag.Panic(err)
	}
}

func (a *activityBiz) Modify(req *definition.ActivityModifyReq) {
	ctx := context.Background()
	_, err := a.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetSort(req.Sort).
		SetImage(req.Image).
		SetLink(req.Link).
		SetRemark(req.Remark).
		Save(ctx)
	if err != nil {
		snag.Panic(err)
	}
}

func (a *activityBiz) Delete(id uint64) {
	ctx := context.Background()
	err := a.orm.DeleteOneID(id).Exec(ctx)
	if err != nil {
		snag.Panic(err)
	}
}

func (a *activityBiz) All() []*definition.ActivityDetail {
	ctx := context.Background()
	items, err := a.orm.Query().Order(ent.Desc(activity.FieldSort)).All(ctx)
	if err != nil {
		snag.Panic(err)
	}
	var res []*definition.ActivityDetail
	for _, item := range items {
		res = append(res, toActivityDetail(item))
	}
	return res
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
			Remark: item.Remark,
		},
	}
}
