// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type battery struct{}

var Battery = new(battery)

// ListModels
// @ID		BatteryModels
// @Router	/manager/v1/battery/model [GET]
// @Summary	获取电池型号
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string											true	"管理员校验token"
// @Param	query			query		model.SelectModelsReq							false	"筛选条件"
// @Success	200				{object}	model.ItemListRes{items=[]model.BatteryModel}	"请求成功"
func (*battery) ListModels(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.SelectModelsReq](c)
	return ctx.SendResponse(service.NewBatteryModel().ListModel(req))
}

// CreateModel
// @ID		BatteryCreateModel
// @Router	/manager/v1/battery/model [POST]
// @Summary	创建电池型号
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string									true	"管理员校验token"
// @Param	body			body		model.BatteryModelReq					true	"电池型号数据"
// @Success	200				{object}	model.ItemRes{item=model.BatteryModel}	"请求成功"
// func (*battery) CreateModel(c echo.Context) (err error) {
// 	ctx, req := app.ManagerContextAndBinding[model.BatteryModelReq](c)
// 	return ctx.SendResponse(model.ItemRes{Item: service.NewBatteryModelWithModifier(ctx.Modifier).CreateModel(req)})
// }

// DeleteModel
// @ID		ManagerBatteryDeleteModel
// @Router	/manager/v1/battery/model [DELETE]
// @Summary	删除电池型号
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.BatteryModelReq	true	"电池型号数据"
// @Success	200				{object}	model.StatusResponse	"请求成功"
// func (*battery) DeleteModel(c echo.Context) (err error) {
// 	ctx, req := app.ManagerContextAndBinding[model.BatteryModelReq](c)
// 	service.NewBatteryModelWithModifier(ctx.Modifier).Delete(req)
// 	return ctx.SendResponse()
// }

// List
// @ID		ManagerBatteryList
// @Router	/manager/v1/battery [GET]
// @Summary	电池列表
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Param	qery			query		model.BatteryListReq								false	"筛选条件"
// @Success	200				{object}	model.PaginationRes{items=[]model.BatteryListRes}	"请求成功"
// func (*battery) List(c echo.Context) (err error) {
// 	ctx, req := app.ManagerContextAndBinding[model.BatteryListReq](c)
// 	return ctx.SendResponse(service.NewBattery(ctx.Modifier).List(req))
// }

// Create
// @ID		ManagerBatteryCreate
// @Router	/manager/v1/battery [POST]
// @Summary	添加电池
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.BatteryCreateReq	true	"电池信息"
// @Success	200				{object}	model.StatusResponse	"请求成功"
// func (*battery) Create(c echo.Context) (err error) {
// 	ctx, req := app.ManagerContextAndBinding[model.BatteryCreateReq](c)
// 	service.NewBattery(ctx.Modifier).Create(req)
// 	return ctx.SendResponse()
// }

// BatchCreate
// @ID			ManagerBatteryBatchCreate
// @Router		/manager/v1/battery/batch [POST]
// @Summary		批量导入电池
// @Description	参考 [MI007 批量导入电车]
// @Tags		电池
// @Accept		json
// @Produce		json
// @Param		X-Manager-Token	header		string					true	"管理员校验token"
// @Param		file			formData	file					true	"电池信息"
// @Success		200				{object}	model.StatusResponse	"请求成功"
// func (*battery) BatchCreate(c echo.Context) (err error) {
// 	ctx := app.ContextX[app.ManagerContext](c)
// 	return ctx.SendResponse(service.NewBattery(ctx.Modifier).BatchCreate(ctx.Context))
// }

// Modify
// @ID		ManagerBatteryModify
// @Router	/manager/v1/battery/{id} [PUT]
// @Summary	修改电池
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"电池ID"
// @Param	body			body		model.BatteryModifyReq	true	"修改信息"
// @Success	200				{object}	model.StatusResponse	"请求成功"
// func (*battery) Modify(c echo.Context) (err error) {
// 	ctx, req := app.ManagerContextAndBinding[model.BatteryModifyReq](c)
// 	service.NewBattery(ctx.Modifier).Modify(req)
// 	return ctx.SendResponse()
// }

// Bind
// @ID		ManagerBatteryBind
// @Router	/manager/v1/battery/bind [POST]
// @Summary	绑定骑手
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.BatteryBind		true	"绑定参数"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*battery) Bind(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BatteryBind](c)
	return ctx.SendResponse(service.NewBattery(ctx.Modifier, ctx.Operator).BindRequest(req))
}

// Unbind
// @ID		ManagerBatteryUnbind
// @Router	/manager/v1/battery/unbind [POST]
// @Summary	解绑骑手
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.BatteryUnbindRequest	true	"解绑参数"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*battery) Unbind(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BatteryUnbindRequest](c)
	return ctx.SendResponse(service.NewBattery(ctx.Modifier, ctx.Operator).Unbind(req))
}

// Detail
// @ID		ManagerBatteryDetail
// @Router	/manager/v1/battery/xc/{sn} [GET]
// @Summary	电池详情
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	sn				path		string					true	"电池编号"
// @Success	200				{object}	model.BatteryBmsDetail	"请求成功"
func (*battery) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BatterySNRequest](c)
	return ctx.SendResponse(service.NewBatteryBms().Detail(req))
}

// Statistics
// @ID		ManagerBatteryStatistics
// @Router	/manager/v1/battery/xc/statistics/{sn} [GET]
// @Summary	电池数据
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	sn				path		string					true	"电池编号"
// @Success	200				{object}	model.BatteryStatistics	"请求成功"
func (*battery) Statistics(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BatterySNRequest](c)
	return ctx.SendResponse(service.NewBatteryBms().Statistics(req))
}

// Position
// @ID		ManagerBatteryPosition
// @Router	/manager/v1/battery/xc/position/{sn} [GET]
// @Summary	电池位置
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	sn				path		string						true	"电池编号"
// @Param	start			query		string						false	"开始时间 (精确到秒, 默认6小时前, 格式为: yyyy-mm-dd hh:mm:ss)"
// @Param	end				query		string						false	"结束时间 (精确到秒, 默认当前时间, 格式为: yyyy-mm-dd hh:mm:ss)"
// @Success	200				{object}	model.BatteryPositionRes	"请求成功"
func (*battery) Position(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BatteryPositionReq](c)
	return ctx.SendResponse(service.NewBatteryBms().Position(req))
}

// Fault
// @ID		ManagerBatteryFault
// @Router	/manager/v1/battery/xc/fault [GET]
// @Summary	电池故障列表
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Param	query			query		model.BatteryFaultReq								false	"请求参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.BatteryFaultRes}	"请求成功"
func (*battery) Fault(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BatteryFaultReq](c)
	return ctx.SendResponse(service.NewBatteryBms().FaultList(req))
}

// TrackRectify
// @ID		ManagerBatteryTrackRectify
// @Router	/manager/v1/battery/track/rectify [POST]
// @Summary	电池轨迹纠偏
// @Tags	电池
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.BatteryTrackReq	true	"轨迹点"
// @Success	200				{object}	model.BatteryTrackRes	"请求成功"
func (*battery) TrackRectify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BatteryTrackReq](c)
	return ctx.SendResponse(service.NewBatteryBms().TrackRectify(req))
}
