package biz

import (
	"context"
	"time"

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

func (s *guideBiz) All() []*model.GuideDetail {
	// 按照sort字段降序排列
	items, err := s.orm.Query().Order(ent.Desc(guide.FieldSort)).All(context.Background())
	if err != nil {
		snag.Panic(err)
	}
	var res []*model.GuideDetail
	for _, item := range items {
		res = append(res, toGuideDetail(item))
	}
	return res
}

func (s *guideBiz) List(req *model.GuideListReq) *model.PaginationRes {
	// 分页按照sort字段降序排列
	query := s.orm.Query().Order(ent.Desc(guide.FieldSort))
	return model.ParsePaginationResponse(query, req.PaginationReq, func(item *ent.Guide) *model.GuideDetail {
		return toGuideDetail(item)
	})
}

func (s *guideBiz) Save(req *model.GuideSaveReq) *model.GuideDetail {
	data, err := s.orm.Create().
		SetName(req.Name).
		SetSort(req.Sort).
		SetAnswer(req.Answer).
		SetRemark(req.Remark).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		snag.Panic(err)
	}
	return toGuideDetail(data)
}

func (s *guideBiz) Modify(req *model.GuideModifyReq) {
	err := s.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetSort(req.Sort).
		SetAnswer(req.Answer).
		SetRemark(req.Remark).
		SetUpdatedAt(time.Now()).
		Exec(context.Background())
	if err != nil {
		snag.Panic(err)
	}
}

func (s *guideBiz) Get(id uint64) *model.GuideDetail {
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

func toGuideDetail(item *ent.Guide) *model.GuideDetail {
	return &model.GuideDetail{
		ID:        item.ID,
		Name:      item.Name,
		Sort:      item.Sort,
		Answer:    item.Answer,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
