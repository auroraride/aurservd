package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/guide"
)

type guideBiz struct {
	orm *ent.GuideClient
	ctx context.Context
}

func NewGuide() *guideBiz {
	return &guideBiz{
		orm: ent.Database.Guide,
		ctx: context.Background(),
	}
}

func (s *guideBiz) All() ([]*definition.GuideDetail, error) {
	// 按照sort字段降序排列
	items, err := s.orm.Query().Order(ent.Desc(guide.FieldSort)).All(s.ctx)
	if err != nil {
		return nil, err
	}
	var res []*definition.GuideDetail
	for _, item := range items {
		res = append(res, toGuideDetail(item))
	}
	return res, nil
}

func (s *guideBiz) List(req *definition.GuideListReq) *model.PaginationRes {
	// 分页按照sort字段降序排列
	query := s.orm.Query().Order(ent.Desc(guide.FieldSort))
	return model.ParsePaginationResponse(query, req.PaginationReq, func(item *ent.Guide) *definition.GuideDetail {
		return toGuideDetail(item)
	})
}

func (s *guideBiz) Create(req *definition.GuideSaveReq) (*definition.GuideDetail, error) {
	data, err := s.orm.Create().
		SetName(req.Name).
		SetSort(req.Sort).
		SetAnswer(req.Answer).
		SetRemark(req.Remark).
		Save(s.ctx)
	if err != nil {
		return nil, err
	}
	return toGuideDetail(data), nil
}

func (s *guideBiz) Modify(req *definition.GuideModifyReq) error {
	err := s.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetSort(req.Sort).
		SetAnswer(req.Answer).
		SetRemark(req.Remark).
		Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *guideBiz) Detail(id uint64) (*definition.GuideDetail, error) {
	data, err := s.orm.Query().Where(guide.ID(id)).First(context.Background())
	if err != nil {
		return nil, err
	}
	return toGuideDetail(data), nil
}

func (s *guideBiz) Delete(id uint64) error {
	_, err := s.orm.SoftDeleteOneID(id).Save(s.ctx)
	if err != nil {
		return err
	}
	return nil
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
