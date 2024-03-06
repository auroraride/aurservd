// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-27
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/labstack/echo/v4"
)

type point struct{}

var Point = new(point)

// Modify
// @ID		ManagerPointModify
// @Router	/manager/v1/point/modify [POST]
// @Summary	M7016 修改骑手积分
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.PointModifyReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*point) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.PointModifyReq](c)
	err = service.NewPointWithModifier(ctx.Modifier).Modify(req)
	if err != nil {
		snag.Panic(err)
	}
	return ctx.SendResponse()
}

// Log
// @ID		ManagerPointLog
// @Router	/manager/v1/point/log [GET]
// @Summary	M7017 积分变动日志
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Param	query			query		model.PointLogListReq								false	"筛选选项"
// @Success	200				{object}	model.PaginationRes{items=[]model.PointLogListRes}	"请求成功"
func (*point) Log(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.PointLogListReq](c)
	return ctx.SendResponse(service.NewPointWithModifier(ctx.Modifier).List(req))
}

// Batch
// @ID		ManagerPointBatch
// @Router	/manager/v1/point/batch [POST]
// @Summary	M7018 批量变动积分
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string				true	"管理员校验token"
// @Param	body			body		model.PointBatchReq	true	"请求参数"
// @Success	200				{object}	[]string			"失败列表"
func (*point) Batch(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.PointBatchReq](c)
	return ctx.SendResponse(service.NewPointWithModifier(ctx.Modifier).Batch(req))
}
