// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-08, by lisicen

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/question"
)

type questionBiz struct {
	orm *ent.QuestionClient
	ctx context.Context
}

func NewQuestionBiz() *questionBiz {
	return &questionBiz{
		orm: ent.Database.Question,
		ctx: context.Background(),
	}
}

// Create 创建
func (b *questionBiz) Create(req *definition.QuestionCreateReq) error {
	_, err := b.orm.Create().
		SetName(req.Name).
		SetSort(req.Sort).
		SetCategoryID(req.CategoryID).
		SetAnswer(req.Answer).
		Save(b.ctx)
	if err != nil {
		return err
	}
	return nil
}

// Modify 修改
func (b *questionBiz) Modify(req *definition.QuestionModifyReq) error {
	_, err := b.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetSort(req.Sort).
		SetCategoryID(req.CategoryID).
		SetAnswer(req.Answer).
		Save(b.ctx)
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (b *questionBiz) Delete(id uint64) error {
	err := b.orm.DeleteOneID(id).Exec(b.ctx)
	if err != nil {
		return err
	}
	return nil
}

// Detail 详情
func (b *questionBiz) Detail(id uint64) (*definition.QuestionDetail, error) {
	qs, err := b.orm.Get(b.ctx, id)
	if err != nil {
		return nil, err
	}
	return &definition.QuestionDetail{
		IDRes: model.IDRes{ID: qs.ID},
		QuestionCommon: definition.QuestionCommon{
			Name:       qs.Name,
			Sort:       qs.Sort,
			CategoryID: qs.CategoryID,
			Answer:     qs.Answer,
		},
	}, nil
}

// List 列表
func (b *questionBiz) List(req *definition.QuestionListReq) (*model.PaginationRes, error) {
	query := b.orm.Query().WithCategory()
	if req.Keyword != nil {
		query = query.Where(question.NameContains(*req.Keyword))
	}
	if req.CategoryID != nil {
		query = query.Where(question.CategoryIDEQ(*req.CategoryID))
	}

	return model.ParsePaginationResponse(query, req.PaginationReq, func(item *ent.Question) *definition.QuestionDetail {
		return &definition.QuestionDetail{
			IDRes: model.IDRes{ID: item.ID},
			QuestionCommon: definition.QuestionCommon{
				Name:         item.Name,
				Sort:         item.Sort,
				CategoryID:   item.CategoryID,
				CategoryName: item.Edges.Category.Name,
				Answer:       item.Answer,
			},
		}
	}), nil
}

// All 全部
func (b *questionBiz) All() ([]*definition.QuestionDetail, error) {
	query := b.orm.Query().WithCategory()
	items, err := query.All(b.ctx)
	if err != nil {
		return nil, err
	}
	var res []*definition.QuestionDetail
	for _, item := range items {
		res = append(res, &definition.QuestionDetail{
			IDRes: model.IDRes{ID: item.ID},
			QuestionCommon: definition.QuestionCommon{
				Name:         item.Name,
				Sort:         item.Sort,
				CategoryID:   item.CategoryID,
				CategoryName: item.Edges.Category.Name,
				Answer:       item.Answer,
			},
		})
	}
	return res, nil
}