// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-28, by aurb

package wapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type cabinet struct{}

var Cabinet = new(cabinet)

// Detail
// @ID		CabinetDetail
// @Router	/warestore/v2/cabinet/{serial} [GET]
// @Summary 获取电柜详情
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	serial				path		string								true	"电柜编号"
// @Param	X-Warestore-Token	header		string								true	"运维校验token"
// @Success	200					{object}	model.MaintainerCabinetDetailRes	"请求成功"
func (*cabinet) Detail(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[model.MaintainerCabinetDetailReq](c)
	return ctx.SendResponse(biz.NewCabinet().Detail(definition.AssetSignInfo{Employee: ctx.Employee}, req.Serial))
}

// Operate
// @ID		CabinetOperate
// @Router	/warestore/v2/cabinet/{serial} [POST]
// @Summary 电柜操作
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string								true	"运维校验token"
// @Param	serial				path		string								true	"电柜编号"
// @Param	body				body		model.MaintainerCabinetOperateReq	true	"请求参数"
// @Success	200					{object}	model.StatusResponse				"请求成功"
func (*cabinet) Operate(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[model.MaintainerCabinetOperateReq](c)
	return ctx.SendResponse(biz.NewCabinet().Operate(definition.AssetSignInfo{Employee: ctx.Employee}, req))
}

// BinOperate
// @ID		CabinetBinOperate
// @Router	/warestore/v2/cabinet/{serial}/{ordinal} [POST]
// @Summary 仓位操作
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"运维校验token"
// @Param	serial				path		string							true	"电柜编号"
// @Param	ordinal				path		int								true	"仓位序号，从1开始"
// @Param	body				body		model.MaintainerBinOperateReq	true	"请求参数"
// @Success	200					{object}	model.StatusResponse			"请求成功"
func (*cabinet) BinOperate(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[model.MaintainerBinOperateReq](c)
	return ctx.SendResponse(biz.NewCabinet().BinOperate(definition.AssetSignInfo{Employee: ctx.Employee}, req))
}
