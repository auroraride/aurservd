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
// @Summary      M5001 创建电柜
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
// @Summary      M5002 查询电柜
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.CabinetQueryReq  true  "搜索参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.CabinetItem}  "请求成功"
func (*cabinet) Query(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CabinetQueryReq](c)
    return ctx.SendResponse(service.NewCabinet().Query(req))
}

// Modify
// @ID           CabinetModify
// @Router       /manager/v1/cabinet/{id} [PUT]
// @Summary      M5003 编辑电柜
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
// @Summary      M5004 删除电柜
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