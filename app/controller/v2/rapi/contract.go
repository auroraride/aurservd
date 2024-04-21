package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type contract struct{}

var Contract = new(contract)

// Sign
// @ID		ContractSign
// @Router	/rider/v2/contract/sign [POST]
// @Summary	签署合同
// @Tags	Contract - 合同
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		definition.ContractSignNewReq	true	"desc"
// @Success	200				{object}	definition.ContractSignNewRes	"请求成功"
func (*contract) Sign(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.ContractSignNewReq](c)
	return ctx.SendResponse(biz.NewContract().Sign(ctx.Rider, req))
}

// Create
// @ID		ContractCreate
// @Router	/rider/v2/contract/create [POST]
// @Summary	创建合同
// @Tags	Contract - 合同
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		definition.ContractCreateReq	true	"desc"
// @Success	200				{object}	definition.ContractCreateRes	"请求成功"
func (*contract) Create(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.ContractCreateReq](c)
	return ctx.SendResponse(biz.NewContract().Create(ctx.Rider, req))
}
