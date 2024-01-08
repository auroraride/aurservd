// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-23
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assistance struct{}

var Assistance = new(assistance)

// Breakdown
// @ID           RiderAssistanceBreakdown
// @Router       /rider/v1/assistance/breakdown [GET]
// @Summary      R5001 获取救援原因
// @Tags         Assistance - 救援
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200 {object}  []string  "请求成功"
func (*assistance) Breakdown(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(service.NewAssistance().Breakdown())
}

// Create
// @ID           RiderAssistanceCreate
// @Router       /rider/v1/assistance [POST]
// @Summary      R5002 发起救援
// @Tags         Assistance - 救援
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.AssistanceCreateReq  true  "救援参数"
// @Success      200 {object}   model.AssistanceCreateRes  "请求成功"
func (*assistance) Create(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.AssistanceCreateReq](c)
	return ctx.SendResponse(service.NewAssistanceWithRider(ctx.Rider).Create(req))
}

// Cancel
// @ID           RiderAssistanceCancel
// @Router       /rider/v1/assistance/cancel [POST]
// @Summary      R5003 取消救援
// @Tags         Assistance - 救援
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body    model.AssistanceCancelReq  true  "取消请求"
// @Success      200 {object}  model.StatusResponse  "请求成功"
func (*assistance) Cancel(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.AssistanceCancelReq](c)
	service.NewAssistanceWithRider(ctx.Rider).Cancel(req)
	return ctx.SendResponse()
}

// Current
// @ID           RiderAssistanceCurrent
// @Router       /rider/v1/assistance/current [GET]
// @Summary      R5004 当前救援
// @Tags         Assistance - 救援
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200 {object}  model.AssistanceSocketMessage  "救援信息, 救援不存在的时候返回data为null"
func (*assistance) Current(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(service.NewAssistance().CurrentMessage(ctx.Rider.ID))
}

// List
// @ID           RiderAssistanceList
// @Router       /rider/v1/assistance [GET]
// @Summary      R5005 救援列表
// @Tags         Assistance - 救援
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query   model.PaginationReq  true  "分页参数"
// @Success      200  {object}  model.Pagination{items=[]model.AssistanceSimpleListRes}  "请求成功"
func (*assistance) List(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.PaginationReq](c)
	return ctx.SendResponse(service.NewAssistanceWithRider(ctx.Rider).SimpleList(*req))
}
