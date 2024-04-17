package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type subscribe struct{}

var Subscribe = new(subscribe)

// StoreModify
// @ID		SubscribeStoreModify
// @Router	/rider/v2/subscribe/store [PUT]
// @Summary	车电套餐修改激活门店
// @Tags	Subscribe - 订阅
// @Accept	json
// @Produce	json
// @Param	body	body		definition.SubscribeStoreModifyReq	true	"请求详情"
// @Success	200		{object}	model.StatusResponse				"请求成功"
func (*subscribe) StoreModify(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.SubscribeStoreModifyReq](c)
	return ctx.SendResponse(biz.NewSubscribe().StoreModify(ctx.Rider, req))
}

// SubscribeStatus
// @ID		SubscribeStatus
// @Router	/rider/v2/subscribe/status [GET]
// @Summary	查询订阅是否激活
// @Tags	Subscribe - 订阅
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string	true	"骑手校验token"
// @Param	id				query		uint64	true	"订阅ID"
// @Success	200				{object}	bool	"TRUE已激活, FALSE未激活"
func (*subscribe) SubscribeStatus(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.IDQueryReq](c)
	return ctx.SendResponse(biz.NewSubscribe().SubscribeStatus(ctx.Rider, &model.EnterpriseRiderSubscribeStatusReq{ID: req.ID}))
}
