// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type reserve struct{}

var Reserve = new(reserve)

// Unfinished
// @ID           RiderReserveUnfinished
// @Router       /rider/v1/reserve [GET]
// @Summary      R8001 获取未完成预约
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200 {object}  model.ReserveUnfinishedRes  "请求成功, 预约不存在时为`null`"
func (*reserve) Unfinished(c echo.Context) (err error) {
    ctx := app.ContextX[app.RiderContext](c)
    return ctx.SendResponse(service.NewReserve().RiderUnfinishedDetail(ctx.Rider.ID))
}

// Create
// @ID           RiderReserveCreate
// @Router       /rider/v1/reserve [POST]
// @Summary      R8002 创建预约
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.ReserveCreateReq  true  "预约信息"
// @Success      200 {object}   model.ReserveUnfinishedRes  "请求成功"
func (*reserve) Create(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.ReserveCreateReq](c)
    return ctx.SendResponse(service.NewReserveWithRider(ctx.Rider).Create(req))
}

// Cancel
// @ID           ManagerReserveCancel
// @Router       /rider/v1/reserve/{id} [DELETE]
// @Summary      R8003 取消预约
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        id  path  uint64  true  "预约ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*reserve) Cancel(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.IDParamReq](c)
    service.NewReserveWithRider(ctx.Rider).Cancel(req)
    return ctx.SendResponse()
}
