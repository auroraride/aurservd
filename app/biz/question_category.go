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
	"github.com/auroraride/aurservd/internal/ent/questioncategory"
)

type questionCategoryBiz struct {
	orm *ent.QuestionCategoryClient
	ctx context.Context
}

func NewQuestionCategoryBiz() *questionCategoryBiz {
	return &questionCategoryBiz{
		orm: ent.Database.QuestionCategory,
		ctx: context.Background(),
	}
}

// Detail 详情
func (q *questionCategoryBiz) Detail(id uint64) (*definition.QuestionCategoryDetail, error) {
	item, err := q.orm.QueryNotDeleted().Where(questioncategory.ID(id)).First(q.ctx)
	if err != nil {
		return nil, err
	}
	return &definition.QuestionCategoryDetail{
		IDRes: model.IDRes{ID: item.ID},
		QuestionCategoryCommon: definition.QuestionCategoryCommon{
			Name:   item.Name,
			Sort:   item.Sort,
			Remark: item.Remark,
		},
	}, nil
}

// Create 创建
func (q *questionCategoryBiz) Create(req *definition.QuestionCategoryCreateReq) error {
	_, err := q.orm.Create().
		SetName(req.Name).
		SetSort(req.Sort).
		SetRemark(req.Remark).
		Save(q.ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("请勿重复添加")
		}
		return err
	}
	return nil
}

// Modify 修改
func (q *questionCategoryBiz) Modify(req *definition.QuestionCategoryModifyReq) error {
	_, err := q.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetSort(req.Sort).
		SetRemark(req.Remark).
		Save(q.ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("请勿重复添加")
		}
		return err
	}
	return nil
}

// Delete 删除
func (q *questionCategoryBiz) Delete(id uint64) error {
	// 删除分类要将分类下的问题移动到其他分类
	err := ent.Database.Question.Update().Where(question.CategoryID(id)).
		SetCategoryID(0).Exec(q.ctx)
	if err != nil {
		return err
	}
	err = q.orm.SoftDeleteOneID(id).Exec(q.ctx)
	if err != nil {
		return err
	}
	return nil
}

// List 列表
func (q *questionCategoryBiz) List(req *definition.QuestionCategoryListReq) (*model.PaginationRes, error) {
	query := q.orm.QueryNotDeleted().Order(ent.Desc(questioncategory.FieldSort))
	return model.ParsePaginationResponse(query, req.PaginationReq, func(item *ent.QuestionCategory) *definition.QuestionCategoryDetail {
		return &definition.QuestionCategoryDetail{
			IDRes: model.IDRes{ID: item.ID},
			QuestionCategoryCommon: definition.QuestionCategoryCommon{
				Name:   item.Name,
				Sort:   item.Sort,
				Remark: item.Remark,
			},
		}
	}), nil
}

// All 获取所有
func (q *questionCategoryBiz) All() ([]*definition.QuestionCategoryDetail, error) {
	items, err := q.orm.QueryNotDeleted().Order(ent.Desc(questioncategory.FieldSort)).All(q.ctx)
	if err != nil {
		return nil, err
	}
	var res []*definition.QuestionCategoryDetail
	for _, item := range items {
		res = append(res, &definition.QuestionCategoryDetail{
			IDRes: model.IDRes{ID: item.ID},
			QuestionCategoryCommon: definition.QuestionCategoryCommon{
				Name:   item.Name,
				Sort:   item.Sort,
				Remark: item.Remark,
			},
		})
	}

	res = append(res, &definition.QuestionCategoryDetail{
		IDRes: model.IDRes{ID: 0},
		QuestionCategoryCommon: definition.QuestionCategoryCommon{
			Name: "其他",
			Sort: 0,
		},
	})

	return res, nil
}
