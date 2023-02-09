// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type battery struct{}

var Battery = new(battery)

// ListModels
// @ID           BatteryModels
// @Router       /manager/v1/battery/model [GET]
// @Summary      M4001 获取电池型号
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  model.ItemListRes{items=[]model.BatteryModel}  "请求成功"
func (*battery) ListModels(c echo.Context) (err error) {
    ctx := app.ContextX[app.ManagerContext](c)
    return ctx.SendResponse(service.NewBatteryModel().List())
}

// CreateModel
// @ID           BatteryCreateModel
// @Router       /manager/v1/battery/model [POST]
// @Summary      M4002 创建电池型号
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.BatteryModelReq  true  "电池型号数据"
// @Success      200  {object}  model.ItemRes{item=model.BatteryModel}  "请求成功"
func (*battery) CreateModel(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BatteryModelReq](c)
    return ctx.SendResponse(model.ItemRes{Item: service.NewBatteryModelWithModifier(ctx.Modifier).CreateModel(req)})
}

// DeleteModel
// @ID           ManagerBatteryDeleteModel
// @Router       /manager/v1/battery/model [DELETE]
// @Summary      M4003 删除电池型号
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.BatteryModelReq  true  "电池型号数据"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*battery) DeleteModel(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BatteryModelReq](c)
    service.NewBatteryModelWithModifier(ctx.Modifier).Delete(req)
    return ctx.SendResponse()
}

// List
// @ID           ManagerBatteryList
// @Router       /manager/v1/battery [GET]
// @Summary      M4004 电池列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token header string true "管理员校验token"
// @Param        qery query   model.BatteryListReq false "筛选条件"
// @Success      200 {object} model.PaginationRes{items=[]model.BatteryListRes} "请求成功"
func (*battery) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BatteryListReq](c)
    return ctx.SendResponse(service.NewBattery(ctx.Modifier).List(req))
}

// Create
// @ID           ManagerBatteryCreate
// @Router       /manager/v1/battery [POST]
// @Summary      M4005 添加电池
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token header string true "管理员校验token"
// @Param        body body    model.BatteryCreateReq true "电池信息"
// @Success      200 {object} model.StatusResponse "请求成功"
func (*battery) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BatteryCreateReq](c)
    service.NewBattery(ctx.Modifier).Create(req)
    return ctx.SendResponse()
}

// BatchCreate
// @ID           ManagerBatteryBatchCreate
// @Router       /manager/v1/battery/batch [POST]
// @Summary      M4006 批量导入电池
// @Description  参考 [MI007 批量导入电车]
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token header string true "管理员校验token"
// @Param        file  formData  file  true  "电池信息"
// @Success      200 {object} model.StatusResponse "请求成功"
func (*battery) BatchCreate(c echo.Context) (err error) {
    ctx := app.ContextX[app.ManagerContext](c)
    return ctx.SendResponse(service.NewBattery(ctx.Modifier).BatchCreate(ctx.Context))
}

// Modify
// @ID           ManagerBatteryModify
// @Router       /manager/v1/battery/{id} [PUT]
// @Summary      M4007 修改电池
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token header string true "管理员校验token"
// @Param        id path uint64 true "电池ID"
// @Param        body body model.BatteryModifyReq true "修改信息"
// @Success      200 {object} model.StatusResponse "请求成功"
func (*battery) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BatteryModifyReq](c)
    service.NewBattery(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}

// Bind
// @ID           ManagerBatteryBind
// @Router       /manager/v1/battery/bind [POST]
// @Summary      M4008 绑定骑手
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.BatteryBind  true  "绑定参数"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*battery) Bind(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BatteryBind](c)
    service.NewBattery(ctx.Modifier).BindRequest(req)
    return ctx.SendResponse()
}

// Unbind
// @ID           ManagerBatteryUnbind
// @Router       /manager/v1/battery/unbind [POST]
// @Summary      M4009 解绑骑手
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.BatteryUnbindRequest  true  "解绑参数"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*battery) Unbind(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BatteryUnbindRequest](c)
    service.NewBattery(ctx.Modifier).UnbindRequest(req)
    return ctx.SendResponse()
}

// Detail
// @ID           ManagerBatteryDetail
// @Router       /manager/v1/battery/xc/{sn} [GET]
// @Summary      M4010 电池详情
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        sn  path  string  true  "电池编号"
// @Success      200  {object}  model.XcBatteryDetail  "请求成功"
func (*battery) Detail(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.XcBatterySNRequest](c)
    return ctx.SendResponse(service.NewBatteryXc().Detail(req))
}

// Statistics
// @ID           ManagerBatteryStatistics
// @Router       /manager/v1/battery/xc/statistics/{sn} [GET]
// @Summary      M4011 电池数据
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        sn  path  string  true  "电池编号"
// @Success      200  {object}  model.XcBatteryStatistics  "请求成功"
func (*battery) Statistics(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.XcBatterySNRequest](c)
    return ctx.SendResponse(service.NewBatteryXc().Statistics(req))
}
