package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
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
func (subscribe) StoreModify(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.SubscribeStoreModifyReq](c)
	return ctx.SendResponse(biz.NewSubscribe().StoreModify(ctx.Rider, req))
}
