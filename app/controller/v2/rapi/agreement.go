package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/model"
)

type agreement struct{}

var Agreement = new(agreement)

// QueryAgreementByEnterprisePriceID
// @ID		AgreementByEnterprisePriceID
// @Router	/rider/v2/agreement/enterprise/price/{id} [GET]
// @Summary	根据企业价格ID查询协议
// @Tags	Agreement - 协议
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	id				path		uint64					true	"企业价格ID"
// @Success	200				{object}	definition.Agreement	"请求成功"
func (*agreement) QueryAgreementByEnterprisePriceID(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewAgreement().QueryAgreementByEnterprisePriceID(req.ID))
}
