// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-14
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type cabinet struct{}

var Cabinet = new(cabinet)

// Create
// @ID		CabinetCreate
// @Router	/manager/v1/cabinet [POST]
// @Summary	创建电柜
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string									true	"管理员校验token"
// @Param	body			body		model.CabinetCreateReq					true	"电柜数据"
// @Success	200				{object}	model.ItemRes{item=model.CabinetItem}	"请求成功"
func (*cabinet) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetCreateReq](c)

	return ctx.SendResponse(
		model.ItemRes{Item: service.NewCabinetWithModifier(ctx.Modifier).CreateCabinet(req)},
	)
}

// List
// @ID		CabinetList
// @Router	/manager/v1/cabinet [GET]
// @Summary	查询电柜
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string											true	"管理员校验token"
// @Param	query			query		model.CabinetQueryReq							true	"搜索参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.CabinetItem}	"请求成功"
func (*cabinet) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetQueryReq](c)
	return ctx.SendResponse(service.NewCabinetWithModifier(ctx.Modifier).List(req))
}

// Modify
// @ID		CabinetModify
// @Router	/manager/v1/cabinet/{id} [PUT]
// @Summary	编辑电柜
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.CabinetCreateReq	true	"电柜数据"
// @Param	id				path		int						true	"电柜ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*cabinet) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetModifyReq](c)
	service.NewCabinetWithModifier(ctx.Modifier).Modify(req)
	return ctx.SendResponse()
}

// Delete
// @ID		CabinetDelete
// @Router	/manager/v1/cabinet/{id} [DELETE]
// @Summary	删除电柜
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		int						true	"电柜ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*cabinet) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetDeleteReq](c)
	service.NewCabinetWithModifier(ctx.Modifier).Delete(req)

	return ctx.SendResponse()
}

// Detail
// @ID		CabinetDetail
// @Router	/manager/v1/cabinet/{id} [GET]
// @Summary	获取并更新电柜详细信息
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		int						true	"电柜ID"
// @Success	200				{object}	model.CabinetDetailRes	"请求成功"
func (*cabinet) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewCabinet().DetailFromID(req.ID))
}

// DoorOperate
// @ID		CabinetDoorOperate
// @Router	/manager/v1/cabinet/door-operate [POST]
// @Summary	仓位操作
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.CabinetDoorOperateReq	true	"柜门操作请求"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*cabinet) DoorOperate(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetDoorOperateReq](c)
	return ctx.SendResponse(
		model.StatusResponse{Status: service.NewCabinetMgrWithModifier(ctx.Modifier).BinOperate(logging.GetOperatorX(ctx.Modifier), req.ID, req)},
	)
}

// Reboot
// @ID		CabinetReboot
// @Router	/manager/v1/cabinet/reboot [POST]
// @Summary	重启电柜
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.IDPostReq			true	"重启请求"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*cabinet) Reboot(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDPostReq](c)

	return ctx.SendResponse(
		model.StatusResponse{Status: service.NewCabinetMgrWithModifier(ctx.Modifier).Reboot(req)},
	)
}

// Fault
// @ID		CabinetFault
// @Router	/manager/v1/cabinet/fault [GET]
// @Summary	故障列表
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Param	query			query		model.CabinetFaultListReq							false	"请求体"
// @Success	200				{object}	model.PaginationRes{items=[]model.CabinetFaultItem}	"请求成功"
func (*cabinet) Fault(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetFaultListReq](c)
	return ctx.SendResponse(
		service.NewCabinetFault().List(req),
	)
}

// FaultDeal
// @ID		CabinetFaultDeal
// @Router	/manager/v1/cabint/fault/{id} [PUT]
// @Summary	处理故障
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	id				path		int							true	"故障ID"
// @Param	body			body		model.CabinetFaultDealReq	true	"desc"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*cabinet) FaultDeal(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetFaultDealReq](c)
	service.NewCabinetFaultWithModifier(ctx.Modifier).Deal(req)
	return ctx.SendResponse()
}

// Data
// @ID		ManagerCabinetData
// @Router	/manager/v1/cabinet/data [GET]
// @Summary	电柜数据表
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Param	query			query		model.CabinetDataReq								false	"筛选数据"
// @Success	200				{object}	model.PaginationRes{items=[]model.CabinetDataRes}	"请求成功"
func (*cabinet) Data(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetDataReq](c)
	return ctx.SendResponse(service.NewCabinetWithModifier(ctx.Modifier).Data(req))
}

// Transfer
// @ID		ManagerCabinetTransfer
// @Router	/manager/v1/cabinet/transfer [POST]
// @Summary	初始化电柜调拨
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.CabinetTransferReq	true	"调拨数据"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*cabinet) Transfer(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetTransferReq](c)
	service.NewCabinetWithModifier(ctx.Modifier).Transfer(req)
	return ctx.SendResponse()
}

// Maintain
// @ID		ManagerCabinetMaintain
// @Router	/manager/v1/cabinet/maintain [POST]
// @Summary	电柜操作维护
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.CabinetMaintainReq	true	"请求参数"
// @Success	200				{object}	model.CabinetDetailRes		"请求成功"
func (*cabinet) Maintain(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetMaintainReq](c)
	return ctx.SendResponse(service.NewCabinetMgrWithModifier(ctx.Modifier).Maintain(logging.GetOperatorX(ctx.Modifier), req))
}

// OpenBind
// @ID			ManagerCabinetOpenBind
// @Router		/manager/v1/cabinet/openbind [POST]
// @Summary		开仓取电池并绑定骑手
// @Description	<仅智能电柜可用, 普通电柜无法请求, 判定标准: `intelligent = true`>
// @Tags		电柜
// @Accept		json
// @Produce		json
// @Param		X-Manager-Token	header		string						true	"管理员校验token"
// @Param		body			body		model.CabinetOpenBindReq	true	"操作请求"
// @Success		200				{object}	model.StatusResponse		"请求成功"
func (*cabinet) OpenBind(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetOpenBindReq](c)
	service.NewIntelligentCabinet(ctx.Modifier, ctx.Operator).OpenBind(req)
	return ctx.SendResponse()
}

// Deactivate
// @ID		CabinetDeactivate
// @Router	/manager/v1/cabinet/bin/deactivate [POST]
// @Summary	仓位逻辑禁用/启用
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.CabinetBinDeactivateReq	true	"柜门操作请求"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*cabinet) Deactivate(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetBinDeactivateReq](c)
	return ctx.SendResponse(
		model.StatusResponse{Status: service.NewCabinetMgrWithModifier(ctx.Modifier).BinOperate(logging.GetOperatorX(ctx.Modifier), req.ID, req)},
	)
}

// Interrupt
// @ID		ManagerCabinetInterrupt
// @Router	/manager/v1/cabinet/interrupt [POST]
// @Summary	中断电柜业务
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.CabinetInterruptRequest	true	"中断请求"
// @Success	200				{object}	pb.CabinetBizResponse			"请求成功"
func (*cabinet) Interrupt(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetInterruptRequest](c)
	return ctx.SendResponse(service.NewCabinetWithModifier(ctx.Modifier).Interrupt(logging.GetOperatorX(ctx.Modifier), req))
}
