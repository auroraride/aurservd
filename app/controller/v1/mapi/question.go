// Copyright (C) liasica. 2022-present.
//
// Created at 2024-03-12
// Based on aurservd by lisicen

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type question struct {
}

var Question = new(question)

// Create
// @ID		QuestionCreate
// @Router	/manager/v1/question [POST]
// @Summary	创建常见问题
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.QuestionCreateReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (a *question) Create(c echo.Context) error {
	ctx, req := app.ManagerContextAndBinding[definition.QuestionCreateReq](c)
	return ctx.SendResponse(biz.NewQuestionWithModifierBiz(ctx.Modifier).Create(req))
}

// Modify
// @ID		QuestionModify
// @Router	/manager/v1/question/:id [PUT]
// @Summary	修改常见问题
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.QuestionModifyReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (a *question) Modify(c echo.Context) error {
	ctx, req := app.ManagerContextAndBinding[definition.QuestionModifyReq](c)
	return ctx.SendResponse(biz.NewQuestionWithModifierBiz(ctx.Modifier).Modify(req))
}

// Detail
// @ID		QuestionDetail
// @Router	/manager/v1/question/{id} [GET]
// @Summary	常见问题详情
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	id				path		string						true	"分类ID"
// @Success	200				{object}	definition.QuestionDetail	"请求成功"
func (a *question) Detail(c echo.Context) error {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewQuestionBiz().Detail(req.ID))
}

// List
// @ID		QuestionList
// @Router	/manager/v1/question [GET]
// @Summary	常见问题分页列表
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	query			query		definition.QuestionListReq	true	"请求参数"
// @Success	200				{object}	[]definition.QuestionDetail	"请求成功"
func (a *question) List(c echo.Context) error {
	ctx, req := app.ManagerContextAndBinding[definition.QuestionListReq](c)
	return ctx.SendResponse(biz.NewQuestionBiz().List(req))
}

// Delete
// @ID		QuestionDelete
// @Router	/manager/v1/question/{id} [DELETE]
// @Summary	删除常见问题
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		string					true	"分类ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (a *question) Delete(c echo.Context) error {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewQuestionWithModifierBiz(ctx.Modifier).Delete(req.ID))
}
