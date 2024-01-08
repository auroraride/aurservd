// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-17
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type wallet struct{}

var Wallet = new(wallet)

// Overview
// @ID           RiderWalletOverview
// @Router       /rider/v1/wallet/overview [GET]
// @Summary      R9001 钱包概览
// @Tags         Wallet - 钱包
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200 {object}  model.WalletOverview  "请求成功"
func (*wallet) Overview(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(service.NewWallet(ctx.Rider).Overview())
}

// PointLog
// @ID           RiderWalletPointLog
// @Router       /rider/v1/wallet/pointlog [GET]
// @Summary      R9002 积分日志
// @Tags         Wallet - 钱包
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query  model.PaginationReq  false  "分页选项"
// @Success      200 {object}  model.PaginationRes{items=[]model.PointLogListRes}  "请求成功"
func (*wallet) PointLog(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.PaginationReq](c)
	return ctx.SendResponse(service.NewPoint().List(&model.PointLogListReq{
		PaginationReq: *req,
		RiderID:       ctx.Rider.ID,
	}))
}

// Points
// @ID           RiderWalletPoints
// @Router       /rider/v1/wallet/points [GET]
// @Summary      R9003 积分详情
// @Tags         Wallet - 钱包
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token header string true "骑手校验token"
// @Success      200 {object} model.PointRes "请求成功"
func (*wallet) Points(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(service.NewPoint().Detail(ctx.Rider))
}

// Coupons
// @ID           RiderWalletCoupons
// @Router       /rider/v1/wallet/coupons [GET]
// @Summary      R9004 优惠券列表
// @Tags         Wallet - 钱包
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token header string true "骑手校验token"
// @Param        query query  model.CouponRiderListReq  false  "筛选条件"
// @Success      200 {object} []model.CouponRider "请求成功"
func (*wallet) Coupons(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.CouponRiderListReq](c)
	return ctx.SendResponse(service.NewCouponWithRider(ctx.Rider).RiderList(req))
}
