// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type plan struct{}

var Plan = new(plan)

// Create
// @ID           PlanCreate
// @Router       /manager/v1/plan [POST]
// @Summary      M6001 创建骑士卡
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.PlanCreateReq  true  "骑行卡信息"
// @Success      200  {object}  model.PlanListRes  "请求成功"
func (*plan) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanCreateReq](c)
    return ctx.SendResponse(service.NewPlanWithModifier(ctx.Modifier).Create(req))
}

// UpdateEnable
// @ID           PlanUpdateEnable
// @Router       /manager/v1/plan/{id} [PUT]
// @Summary      M6002 上下架骑士卡
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  int  true  "骑士卡ID"
// @Param        body  body     model.PlanEnableModifyReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*plan) UpdateEnable(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanEnableModifyReq](c)
    service.NewPlanWithModifier(ctx.Modifier).UpdateEnable(req)
    return ctx.SendResponse()
}

// Delete
// @ID           PlanDelete
// @Router       /manager/v1/plan/{id} [DELETE]
// @Summary      M6003 删除骑士卡
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  int  true  "骑士卡ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*plan) Delete(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
    service.NewPlanWithModifier(ctx.Modifier).Delete(req)
    return ctx.SendResponse()
}

// List
// @ID           PlanList
// @Router       /manager/v1/plan [GET]
// @Summary      M6004 列举骑士卡
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.PlanListReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.PlanListRes}  "请求成功"
func (*plan) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanListReq](c)
    return ctx.SendResponse(service.NewPlan().List(req))
}

// func (*plan) Modify(c echo.Context) (err error) {
//     ctx, req := app.ManagerContextAndBinding[model.PlanModifyReq](c)
//     return ctx.SendResponse(service.NewPlanWithModifier(ctx.Modifier).Modify(req))
// }

// IntroduceNotset
// @ID           ManagerPlanIntroduceNotset
// @Router       /manager/v1/plan/introduce/notset [GET]
// @Summary      M6005 获取未设定介绍的车电型号
// @Description  介绍分两种: 1.单电 2.车电. 但无论哪种, 电池型号必选
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []model.PlanIntroduceOption  "未设定列表"
func (*plan) IntroduceNotset(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewPlanIntroduce().Notset())
}

// IntroduceList
// @ID           ManagerPlanIntroduceList
// @Router       /manager/v1/plan/introduce [GET]
// @Summary      M6006 获取车电介绍
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []model.PlanIntroduce  "请求成功"
func (*plan) IntroduceList(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewPlanIntroduce().List())
}

// IntroduceCreate
// @ID           ManagerPlanIntroduceCreate
// @Router       /manager/v1/plan/introduce [POST]
// @Summary      M6007 创建车电介绍
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.PlanIntroduceCreateReq  true  "介绍详情"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*plan) IntroduceCreate(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanIntroduceCreateReq](c)
    service.NewPlanIntroduce(ctx.Modifier).Create(req)
    return ctx.SendResponse()
}

// IntroduceModify
// @ID           ManagerPlanIntroduceModify
// @Router       /manager/v1/plan/introduce [PUT]
// @Summary      M6008 修改车电介绍
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id     path  uint64  true  "介绍ID"
// @Param        image  body  string  true  "介绍图片"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*plan) IntroduceModify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanIntroduceModifyReq](c)
    service.NewPlanIntroduce(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}

// Time
// @ID           ManagerPlanTime
// @Router       /manager/v1/plan/time [POST]
// @Summary      M6009 修改骑士卡有效期
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.PlanModifyTimeReq  true  "骑士卡和有效期"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*plan) Time(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanModifyTimeReq](c)
    service.NewPlanWithModifier(ctx.Modifier).ModifyTime(req)
    return ctx.SendResponse()
}
