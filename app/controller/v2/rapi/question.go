// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-08, by lisicen

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
)

type question struct{}

var Question = new(question)

// All
// @ID		QuestionAll
// @Router	/rider/v2/question [GET]
// @Summary	常见问题列表
// @Tags	Question - 常见问题
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	[]definition.QuestionDetail	"请求成功"
func (a *question) All(c echo.Context) error {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(biz.NewQuestionBiz().All())
}
