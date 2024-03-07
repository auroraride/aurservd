// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-07, by lisicen

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/advert"
	"github.com/auroraride/aurservd/pkg/snag"
)

type advertBiz struct {
	orm *ent.AdvertClient
}

func NewAdvert() *advertBiz {
	return &advertBiz{
		orm: ent.Database.Advert,
	}
}

func (a *advertBiz) Get(id uint64) *definition.AdvertDetail {
	ctx := context.Background()
	item, err := a.orm.Get(ctx, id)
	if err != nil {
		snag.Panic(err)
	}
	return toAdvertDetail(item)
}

func (a *advertBiz) Create(req *definition.AdvertSaveReq) {
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

func (a *advertBiz) Modify(req *definition.AdvertModifyReq) {
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

func (a *advertBiz) Delete(id uint64) {
	ctx := context.Background()
	err := a.orm.DeleteOneID(id).Exec(ctx)
	if err != nil {
		snag.Panic(err)
	}
}

func (a *advertBiz) All() []*definition.AdvertDetail {
	ctx := context.Background()
	items, err := a.orm.Query().Order(ent.Desc(advert.FieldSort)).All(ctx)
	if err != nil {
		snag.Panic(err)
	}
	var res []*definition.AdvertDetail
	for _, item := range items {
		res = append(res, toAdvertDetail(item))
	}
	return res
}

func (a *advertBiz) List(req *definition.AdvertListReq) *model.PaginationRes {
	query := a.orm.Query().Order(ent.Desc(advert.FieldSort))
	return model.ParsePaginationResponse(query, req.PaginationReq, func(item *ent.Advert) *definition.AdvertDetail {
		return toAdvertDetail(item)
	})
}

func toAdvertDetail(item *ent.Advert) *definition.AdvertDetail {
	return &definition.AdvertDetail{
		ID:     item.ID,
		Name:   item.Name,
		Sort:   item.Sort,
		Image:  item.Image,
		Link:   item.Link,
		Remark: item.Remark,
	}
}
