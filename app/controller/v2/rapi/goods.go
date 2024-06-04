// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-06-03, by Jorjan

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type goods struct{}

var Goods = new(goods)

// List
// @ID		GoodsList
// @Router	/rider/v2/goods [GET]
// @Summary	商品列表
// @Tags	Goods - 商品
// @Accept	json
// @Produce	json
// @Param	query	query		definition.GoodsListForRiderReq	true	"门店列表请求参数"
// @Success	200		{object}	[]definition.GoodsDetail		"请求成功"
func (*goods) List(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.GoodsListForRiderReq](c)
	return ctx.SendResponse(biz.NewGoods().ListForRider(req))
}

// Detail
// @ID		GoodsDetail
// @Router	/rider/v2/goods/{id} [GET]
// @Summary	商品详情
// @Tags	Goods - 商品
// @Accept	json
// @Produce	json
// @Param	id	path		uint64					true	"商品ID"
// @Success	200	{object}	definition.GoodsDetail	"请求成功"
func (*goods) Detail(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewGoods().Detail(req.ID))
}
