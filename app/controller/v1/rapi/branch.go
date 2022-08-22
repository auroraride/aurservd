// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/baidu"
    "github.com/labstack/echo/v4"
    "math"
)

type branch struct{}

var Branch = new(branch)

// List
// @ID           RiderBranchList
// @Router       /rider/v1/branch [GET]
// @Summary      R2001 列举网点
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query   model.BranchWithDistanceReq  true  "根据距离获取网点请求参数"
// @Success      200  {object}  []model.BranchWithDistanceRes  "请求成功"
func (*branch) List(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.BranchWithDistanceReq](c)
    return ctx.SendResponse(service.NewBranchWithRider(ctx.Rider).ListByDistanceRider(req))
}

// Riding
// @ID           RiderBranchRiding
// @Router       /rider/v1/branch/riding [GET]
// @Summary      R2002 网点骑行规划时间
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query  model.BranchRidingReq  true  "desc"
// @Success      200  {object}  model.BranchRidingRes  "请求成功"
func (*branch) Riding(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.BranchRidingReq](c)
    seconds, _, _ := baidu.NewMap().RidingPlanX(req.Origin, req.Destination)
    return ctx.SendResponse(model.BranchRidingRes{Minutes: math.Round(float64(seconds) / 60.0)})
}

// Facility
// @ID           RiderBranchFacility
// @Router       /rider/v1/branch/facility/{fid} [GET]
// @Summary      R2004 设施详情
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        fid  path  string  true  "设置标识"
// @Param        lng  query  float64  true  "经度"
// @Param        lat  query  float64  true  "纬度"
// @Success      200 {object}  model.BranchFacilityRes  "请求成功"
func (*branch) Facility(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.BranchFacilityReq](c)
    return ctx.SendResponse(service.NewBranchWithRider(ctx.Rider).Facility(req))
}
