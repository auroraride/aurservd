// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-08
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/labstack/echo/v4"
)

type enterprise struct{}

var Enterprise = new(enterprise)

// Battery
// @ID           RiderEnterpriseBattery
// @Router       /rider/v1/enterprise/battery [GET]
// @Summary      R3010 企业骑手获取可用电池
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        cityId  query  uint64  true  "城市ID"
// @Success      200  {object}  []string  "请求成功"
func (*enterprise) Battery(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.EnterprisePriceBatteryModelListReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).BatteryModels(req))
}

// Subscribe
// @ID           RiderEnterpriseSubscribe
// @Router       /rider/v1/enterprise/subscribe [POST]
// @Summary      R3011 企业骑手选择电池
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.EnterpriseRiderSubscribeChooseReq  true  "电池选择请求"
// @Success      200  {object}  model.EnterpriseRiderSubscribeChooseRes  "请求成功"
func (*enterprise) Subscribe(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.EnterpriseRiderSubscribeChooseReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).ChooseBatteryModel(req))
}

// SubscribeStatus
// @ID           RiderEnterpriseSubscribeStatus
// @Router       /rider/v1/enterprise/subscribe [GET]
// @Summary      R3012 企业骑手订阅激活状态
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        id  query  uint64  true  "订阅ID"
// @Success      200  {object}  bool  "TRUE已激活, FALSE未激活"
func (*enterprise) SubscribeStatus(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.EnterpriseRiderSubscribeStatusReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).SubscribeStatus(req))
}
