// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-14
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type cabinet struct{}

var Cabinet = new(cabinet)

// Create
// @ID           CabinetCreate
// @Router       /manager/v1/cabinet [POST]
// @Summary      M50001 创建电柜
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.CabinetCreateReq  true  "电柜数据"
// @Success      200  {object}  model.ItemRes{item=model.CabinetItem}  "请求成功"
func (*cabinet) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CabinetCreateReq](c)

    return ctx.SendResponse(
        model.ItemRes{Item: service.NewCabinet().CreateCabinet(ctx.Modifier, req)},
    )
}

// Query
// @ID           CabinetQuery
// @Router       /manager/v1/cabinet [GET]
// @Summary      M50002 查询电柜
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.CabinetQueryReq  true  "搜索参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.CabinetItem}  "请求成功"
func (*cabinet) Query(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CabinetQueryReq](c)
    return ctx.SendResponse(service.NewCabinet().List(req))
}

// Modify
// @ID           CabinetModify
// @Router       /manager/v1/cabinet/{id} [PUT]
// @Summary      M50003 编辑电柜
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.CabinetCreateReq  true  "电柜数据"
// @Param        id    path  int  true  "电柜ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*cabinet) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CabinetModifyReq](c)
    service.NewCabinet().Modify(req)
    return ctx.SendResponse()
}

// Delete
// @ID           CabinetDelete
// @Router       /manager/v1/cabinet/{id} [DELETE]
// @Summary      M50004 删除电柜
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id    path  int  true  "电柜ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*cabinet) Delete(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CabinetDeleteReq](c)
    service.NewCabinet().Delete(ctx.Modifier, req)

    return ctx.SendResponse()
}

// Detail
// @ID           CabinetDetail
// @Router       /manager/v1/cabinet/{id} [GET]
// @Summary      M50005 获取并更新电柜详细信息
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id    path  int  true  "电柜ID"
// @Success      200  {object}  model.CabinetDetailRes  "请求成功"
func (*cabinet) Detail(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
    return ctx.SendResponse(service.NewCabinet().Detail(req.ID))
}

// DoorOperate
// @ID           CabinetDoorOperate
// @Router       /manager/v1/cabinet/door-operate [POST]
// @Summary      M50006 柜门操作
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.CabinetDoorOperateReq  true  "柜门操作请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*cabinet) DoorOperate(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CabinetDoorOperateReq](c)
    return ctx.SendResponse(
        model.StatusResponse{Status: service.NewCabinet().DoorOperate(ctx.Modifier, req)},
    )
}

// Reboot
// @ID           CabinetReboot
// @Router       /manager/v1/cabinet/reboot [POST]
// @Summary      M50007 重启电柜
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.IDPostReq  true  "重启请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*cabinet) Reboot(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDPostReq](c)
    return ctx.SendResponse(
        model.StatusResponse{Status: service.NewCabinet().Reboot(ctx.Modifier, req)},
    )
}

// Fault
// @ID           CabinetFault
// @Router       /manager/v1/cabinet/fault [GET]
// @Summary      M50008 故障列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query  model.CabinetFaultListReq  false  "请求体"
// @Success      200  {object}  model.PaginationRes{items=[]model.CabinetFaultItem}  "请求成功"
func (*cabinet) Fault(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CabinetFaultListReq](c)
    return ctx.SendResponse(
        service.NewCabinetFault().List(req),
    )
}

// FaultDeal
// @ID           CabinetFaultDeal
// @Router       /manager/v1/fault/{id} [PUT]
// @Summary      M50009 处理故障
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id    path  int  true  "故障ID"
// @Param        body  body  model.CabinetFaultDealReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*cabinet) FaultDeal(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CabinetFaultDealReq](c)
    service.NewCabinetFault().Deal(ctx.Modifier, req)
    return ctx.SendResponse()
}
