package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type store struct{}

var Store = new(store)

// List 门店列表
// @ID		StoreList
// @Router	/rider/v2/store [GET]
// @Summary	门店列表
// @Tags	Store - 门店
// @Accept	json
// @Produce	json
// @Param	query	query		definition.StoreListReq		true	"门店列表请求参数"
// @Success	200		{object}	[]definition.StoreDetail	"请求成功"
func (*store) List(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.StoreListReq](c)
	return ctx.SendResponse(biz.NewStore().List(req))
}

// Detail 门店详情
// @ID		StoreDetail
// @Router	/rider/v2/store/{id} [GET]
// @Summary	门店详情
// @Tags	Store - 门店
// @Accept	json
// @Produce	json
// @Param	id		path		uint64						true	"订阅ID"
// @Param	query	query		definition.StoreDetailReq	true	"门店详情"
// @Success	200		{object}	definition.StoreDetail		"请求成功"
func (*store) Detail(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.StoreDetailReq](c)
	return ctx.SendResponse(biz.NewStore().Detail(req))
}

// StoreBySubscribe
// @ID		StoreBySubscribe
// @Router	/rider/v2/store/subscribe/{id} [GET]
// @Summary	根据订阅查询骑手激活门店信息
// @Tags	Store - 门店
// @Accept	json
// @Produce	json
// @Param	id		path		uint64						true	"订阅ID"
// @Param	query	query		definition.StoreDetailReq	true	"门店详情"
// @Success	200		{object}	definition.StoreDetail		"请求成功"
func (*store) StoreBySubscribe(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.StoreDetailReq](c)
	return ctx.SendResponse(biz.NewStore().StoreBySubscribe(ctx.Rider, req))
}
