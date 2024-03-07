package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/guide"
	"github.com/auroraride/aurservd/pkg/snag"
)

type guideBiz struct {
	orm *ent.GuideClient
}

func NewGuide() *guideBiz {
	return &guideBiz{
		orm: ent.Database.Guide,
	}
}

func (s *guideBiz) All() []*definition.GuideDetail {
	// 按照sort字段降序排列
	items, err := s.orm.Query().Order(ent.Desc(guide.FieldSort)).All(context.Background())
	if err != nil {
		snag.Panic(err)
	}
	var res []*definition.GuideDetail
	for _, item := range items {
		res = append(res, toGuideDetail(item))
	}
	return res
}

func (s *guideBiz) List(req *definition.GuideListReq) *model.PaginationRes {
	// 分页按照sort字段降序排列
	query := s.orm.Query().Order(ent.Desc(guide.FieldSort))
	return model.ParsePaginationResponse(query, req.PaginationReq, func(item *ent.Guide) *definition.GuideDetail {
		return toGuideDetail(item)
	})
}

func (s *guideBiz) Save(req *definition.GuideSaveReq) *definition.GuideDetail {
	data, err := s.orm.Create().
		SetName(req.Name).
		SetSort(req.Sort).
		SetAnswer(req.Answer).
		SetRemark(req.Remark).
		Save(context.Background())
	if err != nil {
		snag.Panic(err)
	}
	return toGuideDetail(data)
}

func (s *guideBiz) Modify(req *definition.GuideModifyReq) {
	err := s.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetSort(req.Sort).
		SetAnswer(req.Answer).
		SetRemark(req.Remark).
		Exec(context.Background())
	if err != nil {
		snag.Panic(err)
	}
}

func (s *guideBiz) Get(id uint64) *definition.GuideDetail {
	data, err := s.orm.Query().Where(guide.ID(id)).First(context.Background())
	if err != nil {
		snag.Panic(err)
	}
	return toGuideDetail(data)
}

func (s *guideBiz) Delete(id uint64) bool {
	_, err := s.orm.Delete().Where(guide.ID(id)).Exec(context.Background())
	if err != nil {
		snag.Panic(err)
	}
	return true
}

func toGuideDetail(item *ent.Guide) *definition.GuideDetail {
	return &definition.GuideDetail{
		ID:     item.ID,
		Name:   item.Name,
		Sort:   item.Sort,
		Answer: item.Answer,
		Remark: item.Remark,
	}
}
