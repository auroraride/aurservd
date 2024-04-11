// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-08, by lisicen

package biz

import (
	"context"
	"errors"
	"strings"

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
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("请勿重复添加")
		}
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
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("请勿重复添加")
		}
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
	qs, err := b.orm.Query().Where(question.ID(id)).First(b.ctx)
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
	query := b.orm.Query().WithCategory(func(query *ent.QuestionCategoryQuery) {
		query.Order(ent.Desc(question.FieldSort))
	}).Order(ent.Desc(question.FieldSort))
	if req.Keyword != nil {
		query.Where(question.NameContains(*req.Keyword))
	}
	if req.CategoryID != nil {
		query.Where(question.CategoryIDEQ(*req.CategoryID))
	}

	return model.ParsePaginationResponse(query, req.PaginationReq, func(item *ent.Question) *definition.QuestionDetail {
		questionCommon := definition.QuestionCommon{
			Name:       item.Name,
			Sort:       item.Sort,
			CategoryID: item.CategoryID,
			Answer:     item.Answer,
		}
		if item.Edges.Category != nil {
			questionCommon.CategoryName = item.Edges.Category.Name
		} else {
			questionCommon.CategoryName = "其他"
		}

		return &definition.QuestionDetail{
			IDRes:          model.IDRes{ID: item.ID},
			QuestionCommon: questionCommon,
		}
	}), nil
}

// All 全部
func (b *questionBiz) All() ([]*definition.QuestionDetail, error) {
	query := b.orm.Query().WithCategory(
		func(query *ent.QuestionCategoryQuery) {
			query.Order(ent.Desc(question.FieldSort))
		},
	).Order(ent.Desc(question.FieldSort))
	items, err := query.All(b.ctx)
	if len(items) == 0 {
		return nil, err
	}
	var res []*definition.QuestionDetail
	for _, item := range items {
		data := &definition.QuestionDetail{
			IDRes: model.IDRes{ID: item.ID},
			QuestionCommon: definition.QuestionCommon{
				Name:   item.Name,
				Sort:   item.Sort,
				Answer: item.Answer,
			}}
		if item.Edges.Category != nil {
			data.QuestionCommon.CategoryName = item.Edges.Category.Name
			data.QuestionCommon.CategoryID = item.CategoryID
		}
		res = append(res, data)

	}
	return res, nil
}
