// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type battery struct{}

var Battery = new(battery)

// Selection
// @ID		AgentBatterySelection
// @Router	/agent/v1/battery/selection [GET]
// @Summary	AA001 电池选择搜索
// @Tags	[A]代理接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string					true	"代理校验token"
// @Param	query			query		model.BatterySearchReq	true	"筛选项"
// @Success	200				{object}	[]model.Battery
func (*battery) Selection(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.BatterySearchReq](c)
	// TODO 子代理
	return ctx.SendResponse(service.NewSelection().BatterySerialSearch(&model.BatterySearchReq{
		Serial:       req.Serial,
		EnterpriseID: &ctx.Enterprise.ID,
	}))
}

// Model
// @ID		AgentBatteryModel
// @Router	/agent/v1/battery/model [GET]
// @Summary	AA002 电池型号列表
// @Tags	[A]代理接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string						true	"代理校验token"
// @Param	query			query		model.BatteryModelListReq	false	"筛选项"
// @Success	200				{object}	model.ItemListRes
func (*battery) Model(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.BatteryModelListReq](c)
	return ctx.SendResponse(service.NewBatteryModel().List(req))
}

// List 电池列表
// @ID		AgentBatteryList
// @Router	/agent/v1/battery [GET]
// @Summary	AA003 电池列表
// @Tags	[A]代理接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string											true	"代理校验token"
// @Param	query			query		model.AssetListReq								false	"筛选项"
// @Success	200				{object}	model.ItemListRes{items=[]model.AssetListRes}	"请求成功"
func (*battery) List(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AssetListReq](c)
	return ctx.SendResponse(service.NewAsset().List(ctx.Request().Context(), req))
}
