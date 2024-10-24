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
	"github.com/auroraride/aurservd/pkg/silk"
)

type enterprise struct{}

var Enterprise = new(enterprise)

// Battery
// @ID		EnterpriseBattery
// @Router	/rider/v1/enterprise/battery [GET]
// @Summary	R3010 企业骑手获取可用电池
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string		true	"骑手校验token"
// @Param	cityId			query		uint64		true	"城市ID"
// @Success	200				{object}	[]string	"请求成功"
func (*enterprise) Battery(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.EnterprisePricePlanListReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).BatteryModels(req))
}

// Subscribe
// @ID		EnterpriseSubscribe
// @Router	/rider/v1/enterprise/subscribe [POST]
// @Summary	R3011 企业骑手选择电池
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string									true	"骑手校验token"
// @Param	body			body		model.EnterpriseRiderSubscribeChooseReq	true	"电池选择请求"
// @Success	200				{object}	model.EnterpriseRiderSubscribeChooseRes	"请求成功"
func (*enterprise) Subscribe(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.EnterpriseRiderSubscribeChooseReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).ChooseBatteryModel(req))
}

// SubscribeStatus
// @ID		EnterpriseSubscribeStatus
// @Router	/rider/v1/enterprise/subscribe [GET]
// @Summary	R3012 企业骑手订阅激活状态
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string	true	"骑手校验token"
// @Param	id				query		uint64	true	"订阅ID"
// @Success	200				{object}	bool	"TRUE已激活, FALSE未激活"
func (*enterprise) SubscribeStatus(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.EnterpriseRiderSubscribeStatusReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).SubscribeStatus(req))
}

// SubscribeAlter
// @ID		EnterpriseSubscribeAlter
// @Router	/rider/v1/enterprise/subscribe/alter [POST]
// @Summary	R3013 加时申请
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		model.SubscribeAlterRiderReq	true	"申请增加订阅时长请求"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*enterprise) SubscribeAlter(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.SubscribeAlterRiderReq](c)
	service.NewSubscribeAlter(ctx.Rider).AlterDays(ctx.Rider, req)
	return ctx.SendResponse()
}

// SubscribeAlterList
// @ID		EnterpriseSubscribeAlterList
// @Router	/rider/v1/enterprise/subscribe/alter/list [GET]
// @Summary	R3014 加时申请列表
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string															true	"骑手校验token"
// @Param	query			query		model.SubscribeListRiderReq										true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.SubscribeAlterApplyListRes}	"请求成功"
func (*enterprise) SubscribeAlterList(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.SubscribeListRiderReq](c)
	return ctx.SendResponse(service.NewSubscribeAlter().List(&model.SubscribeAlterListReq{
		SubscribeAlterFilter: model.SubscribeAlterFilter{
			PaginationReq: req.PaginationReq,
			Start:         req.Start,
			End:           req.End,
			Status:        req.Status,
			RiderID:       silk.Pointer(ctx.Rider.ID),
		},
	}))
}

// JoinEnterprise
// @ID		EnterpriseJoinEnterprise
// @Router	/rider/v1/enterprise/join [POST]
// @Summary	R3015 企业骑手加入团签
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	body	body		model.EnterproseInfoReq	true	"加入团签请求"
// @Success	200		{object}	bool					"请求成功"
func (s *enterprise) JoinEnterprise(c echo.Context) error {
	ctx, req := app.RiderContextAndBinding[model.EnterpriseJoinReq](c)
	service.NewEnterpriseRider().JoinEnterprise(req, ctx.Rider)
	return ctx.SendResponse()
}

// RiderEnterpriseInfo
// @ID		EnterpriseRiderEnterpriseInfo
// @Router	/rider/v1/enterprise/info [GET]
// @Summary	R3016 骑手团签信息
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	query			query		model.EnterproseInfoReq	true	"团签信息请求"
// @Success	200				{object}	model.EnterproseInfoRsp	"请求成功"
func (*enterprise) RiderEnterpriseInfo(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.EnterproseInfoReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithRider(ctx.Rider).RiderEnterpriseInfo(req, ctx.Rider.ID))
}

// ExitEnterprise
// @ID		EnterpriseExitEnterprise
// @Router	/rider/v1/enterprise/exit [POST]
// @Summary	R3017 退出团签
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*enterprise) ExitEnterprise(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	service.NewEnterpriseRider().ExitEnterprise(ctx.Rider)
	return ctx.SendResponse()
}
