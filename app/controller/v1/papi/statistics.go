// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package papi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type statistics struct{}

var Statistics = new(statistics)

// TemaOverview
// @ID           PromotionStatisticsTemaOverview
// @Router       /promotion/v1/statistics/team/overview [GET]
// @Summary      P5001 首页我的团队统计
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "代理校验token"
// @Param        body  body  promotion.StatisticsReq  true  "请求参数"
// @Success      200  {object}  promotion.StatisticsTeamRes  "请求成功"
func (*statistics) TemaOverview(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.StatisticsReq](c)
	return ctx.SendResponse(service.NewPromotionStatisticsService().Team(ctx.Member, req))
}

// EarningsOverview
// @ID           PromotionStatisticsEarningsOverview
// @Router       /promotion/v1/statistics/earnings/overview [GET]
// @Summary      P5002 首页我的收益统计
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "代理校验token"
// @Param        body  body  promotion.StatisticsReq  true  "请求参数"
// @Success      200  {object}  promotion.StatisticsEarningsRes  "请求成功"
func (*statistics) EarningsOverview(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.StatisticsReq](c)
	return ctx.SendResponse(service.NewPromotionStatisticsService().Earnings(ctx.Member, req))
}

// RecordOverview
// @ID           PromotionStatisticsRecordOverview
// @Router       /promotion/v1/statistics/record/overview [GET]
// @Summary      P5003 首页邀请战绩统计
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "代理校验token"
// @Success      200  {object}   promotion.StatisticsRecordRes  "请求成功"
func (*statistics) RecordOverview(c echo.Context) (err error) {
	ctx := app.ContextX[app.PromotionContext](c)
	return ctx.SendResponse(service.NewPromotionStatisticsService().Record(ctx.Member))
}

// WalletOverview
// @ID           PromotionStatisticsWalletOverview
// @Router       /promotion/v1/statistics/wallet/overview [GET]
// @Summary      P5004 我的账户-统计
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "代理校验token"
// @Success      200  {object}   promotion.StatisticsWalletRes  "请求成功"
func (*statistics) WalletOverview(c echo.Context) (err error) {
	ctx := app.ContextX[app.PromotionContext](c)
	return ctx.SendResponse(service.NewPromotionStatisticsService().Wallet(ctx.Member))
}

// EarningsDetail
// @ID           PromotionStatisticsEarningsDetail
// @Router       /promotion/v1/statistics/earnings/detail [GET]
// @Summary      P5005 我的钱包-收益报表
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "代理校验token"
// @Param        body  body  promotion.StatisticsReq  true  "请求参数"
// @Success      200  {object}   promotion.StatisticsEarningsDetailRes  "请求成功"
func (*statistics) EarningsDetail(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.StatisticsReq](c)
	return ctx.SendResponse(service.NewPromotionStatisticsService().EarningsDetail(ctx.Member, req))
}

// Team
// @ID           PromotionStatisticsTeam
// @Router       /promotion/v1/statistics/team [GET]
// @Summary      P5006 我的团队-统计
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "代理校验token"
// @Success      200  {object}   []promotion.StatisticsTeamRes  "请求成功"
func (*statistics) Team(c echo.Context) (err error) {
	ctx := app.ContextX[app.PromotionContext](c)
	return ctx.SendResponse(service.NewPromotionStatisticsService().MyTeamStatistics(ctx.Member))
}

// TeamGrowth
// @ID           PromotionStatisticsTeamGrowth
// @Router       /promotion/v1/statistics/team/growth [GET]
// @Summary      P5007 我的团队-增长趋势
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "代理校验token"
// @Param        body  body  promotion.StatisticsReq  true  "请求参数"
// @Success      200  {object}   []promotion.StatisticsTeamGrowthTrendRes  "请求成功"
func (*statistics) TeamGrowth(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.StatisticsReq](c)
	return ctx.SendResponse(service.NewPromotionStatisticsService().TeamGrowth(ctx.Member, req))
}
