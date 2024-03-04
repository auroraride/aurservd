package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionGrowth struct {
}

var PromotionGrowth = new(promotionGrowth)

// List
//	@ID			ManagerPromotionGrowthList
//	@Router		/manager/v1/promotion/growth/{id} [GET]
//	@Summary	PM5001 会员成长值列表
//	@Tags		[PM]推广管理接口
//	@Accept		json
//	@Produce	json
//	@Param		X-Manager-Token	header		string				true	"管理员校验token"
//	@Param		body			body		promotion.GrowthReq	true	"请求参数"
//	@Success	200				{object}	[]promotion.GrowthRes
func (l *promotionGrowth) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.GrowthReq](c)
	return ctx.SendResponse(service.NewPromotionGrowthService().List(req))
}
