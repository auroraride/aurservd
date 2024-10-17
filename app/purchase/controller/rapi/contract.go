package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/purchase/internal/model"
	"github.com/auroraride/aurservd/app/purchase/internal/service"
)

type contract struct{}

var Contract = new(contract)

// Sign
// @ID		ContractSign
// @Router	/rider/v2/purchase/contract/sign [POST]
// @Summary	签署合同
// @Tags	Contract - 合同
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		model.ContractSignNewReq	true	"desc"
// @Success	200				{object}	model.ContractSignNewRes	"请求成功"
func (*contract) Sign(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.ContractSignNewReq](c)
	return ctx.SendResponse(service.NewContract().Sign(ctx.Request().Context(), ctx.Rider, req))
}

// Create
// @ID		ContractCreate
// @Router	/rider/v2/purchase/contract/create [POST]
// @Summary	创建合同
// @Tags	Contract - 合同
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		model.ContractCreateReq	true	"desc"
// @Success	200				{object}	model.ContractCreateRes	"请求成功"
func (*contract) Create(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.ContractCreateReq](c)
	return ctx.SendResponse(service.NewContract().Create(ctx.Request().Context(), ctx.Rider, req))
}

// Detail
// @ID		ContractDetail
// @Router	/rider/v2/purchase/contract/{docId} [GET]
// @Summary	查看合同
// @Tags	Contract - 合同
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	docId			path		string							true	"合同ID"
// @Success	200				{object}	model.ContractDetailRes	"请求成功"
func (*contract) Detail(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.ContractDetailReq](c)
	return ctx.SendResponse(service.NewContract().Detail(ctx.Request().Context(), ctx.Rider, req))
}
