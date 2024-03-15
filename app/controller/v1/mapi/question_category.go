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

type questionCategory struct {
}

var QuestionCategory = new(questionCategory)

// Create
// @ID		Question_categoryCreate
// @Router	/manager/v1/question/category [POST]
// @Summary	创建问题分类
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string									true	"管理员校验token"
// @Param	body			body		definition.QuestionCategoryCreateReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse					"请求成功"
func (*questionCategory) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.QuestionCategoryCreateReq](c)
	return ctx.SendResponse(biz.NewQuestionCategoryBiz().Create(req))
}

// Modify
// @ID		Question_categoryModify
// @Router	/manager/v1/question/category [PUT]
// @Summary	修改问题分类
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string									true	"管理员校验token"
// @Param	body			body		definition.QuestionCategoryModifyReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse					"请求成功"
func (*questionCategory) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.QuestionCategoryModifyReq](c)
	return ctx.SendResponse(biz.NewQuestionCategoryBiz().Modify(req))
}

// Detail
// @ID		Question_categoryDetail
// @Router	/manager/v1/question/category/{id} [GET]
// @Summary	问题分类详情
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		string					true	"分类ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*questionCategory) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewQuestionCategoryBiz().Detail(req.ID))
}

// List
// @ID		Question_categoryList
// @Router	/manager/v1/question/category [GET]
// @Summary	问题分类分页列表
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// @Param	query			query		definition.QuestionCategoryListReq	true	"请求参数"
// @Success	200				{object}	[]definition.QuestionCategoryDetail	"请求成功"
func (*questionCategory) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.QuestionCategoryListReq](c)
	return ctx.SendResponse(biz.NewQuestionCategoryBiz().List(req))
}

// Delete
// @ID		Question_categoryDelete
// @Router	/manager/v1/question/category/{id} [DELETE]
// @Summary	删除问题分类
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		model.IDParamReq		true	"分类ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*questionCategory) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewQuestionCategoryBiz().Delete(req.ID))
}
