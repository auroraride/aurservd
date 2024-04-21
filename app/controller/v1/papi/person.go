package papi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionPerson struct {
}

var PromotionPerson = new(promotionPerson)

// RealName  实名认证
// @ID		PromotionRealNameAuth
// @Router	/promotion/v1/auth/realname [POST]
// @Summary	P6001 实名认证
// @Tags	[P]推广接口
// @Accept	json
// @Produce	json
// @Param	X-Promotion-Token	header		string						true	"会员校验token"
// @Param	body				body		promotion.RealNameAuthReq	true	"请求参数"
// @Success	200					{object}	promotion.RealNameAuthRes
func (p *promotionPerson) RealName(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.RealNameAuthReq](c)
	return ctx.SendResponse(service.NewPromotionPersonService().RealNameAuth(ctx.Member, req))
}
