// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-08, by lisicen

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
)

type questionCategory struct{}

var QuestionCategory = new(questionCategory)

// All
// @ID		QuestionCategoryAll
// @Router	/rider/v2/question/category [GET]
// @Summary	问题分类列表
// @Tags	QuestionCategory - 问题分类
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string								true	"骑手校验token"
// @Success	200				{object}	[]definition.QuestionCategoryDetail	"请求成功"
func (a *questionCategory) All(c echo.Context) error {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(biz.NewQuestionCategoryBiz().All())
}
