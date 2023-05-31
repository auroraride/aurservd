// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-08
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
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

// GetSubscribeInfo 团签骑手个人信息
// @ID           RiderEnterpriseGetSubscribeInfo
// @Router       /rider/v1/enterprise/subscribe/info [GET]
// @Summary      R3013 企业骑手订阅信息
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200  {object}  model.RiderSubscribeRsp  "请求成功"
func (*enterprise) GetSubscribeInfo(c echo.Context) (err error) {
	ctx := c.(*app.RiderContext)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).GetSubscribeInfo(ctx.Rider))
}

// ApplyAddSubscribeTime 骑手申请增加订阅时长
// @ID           RiderEnterpriseApplyAddSubscribeTime
// @Router       /rider/v1/enterprise/subscribe/add [POST]
// @Summary      R3014 企业骑手申请增加订阅时长
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.RiderSubscribeAddReq  true  "申请增加订阅时长请求"
// @Success      200  {object}  bool  "请求成功"
func (*enterprise) ApplyAddSubscribeTime(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.RiderSubscribeAddReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).AddSubscribeDays(req))
}

// ApplyList 申请列表
// @ID           RiderEnterpriseApplyList
// @Router       /rider/v1/enterprise/subscribe/alter/list [GET]
// @Summary      R3015 企业骑手申请列表
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query   model.ApplyReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.ApplyListRsp}  "请求成功"
func (*enterprise) ApplyList(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.ApplyReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).SubscribeAlterList(req))
}

// RiderEnterpriseInfo 骑手团签信息
// @ID           RiderEnterpriseInfo
// @Router       /rider/v1/enterprise/info [GET]
// @Summary      R3016 企业骑手团签信息
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.EnterproseInfoReq  true  "团签信息请求"
// @Success      200  {object}  model.EnterproseInfoRsp  "请求成功"
func (*enterprise) RiderEnterpriseInfo(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.EnterproseInfoReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).RiderEnterpriseInfo(req))
}

// JoinEnterprise 加入团签
// @ID           RiderEnterpriseJoin
// @Router       /rider/v1/enterprise/join [POST]
// @Summary      R3017 企业骑手加入团签
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        body  body     model.RiderJoinEnterpriseReq true  "加入团签请求"
// @Success      200  {object}  bool  "请求成功"
func (s *enterprise) JoinEnterprise(c echo.Context) error {
	ctx, req := app.RiderContextAndBinding[model.EnterproseInfoReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).JoinEnterprise(req))
}
