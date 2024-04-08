package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type contractTemplate struct{}

var ContractTemplate = new(contractTemplate)

// List
// @ID		ManagerContractTemplateList
// @Router	/manager/v1/contract/template [GET]
// @Summary	合同模板列表
// @Tags	合同
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string									true	"管理员校验token"
// @Success	200				{object}	[]definition.ContractTemplateListRes	"请求成功"
func (*contractTemplate) List(c echo.Context) (err error) {
	ctx := app.GetManagerContext(c)
	return ctx.SendResponse(biz.NewContractTemplate().List())
}

// Create
// @ID		ManagerContractTemplateCreate
// @Router	/manager/v1/contract/template [POST]
// @Summary	创建合同模板
// @Tags	合同
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string									true	"管理员校验token"
// @Param	body			body		definition.ContractTemplateCreateReq	true	"desc"
// @Success	200				{object}	model.StatusResponse					"请求成功"
func (*contractTemplate) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.ContractTemplateCreateReq](c)
	return ctx.SendResponse(biz.NewContractTemplate().Create(req))
}

// Modify
// @ID		ManagerContractTemplateModify
// @Router	/manager/v1/contract/template/{id} [PUT]
// @Summary	修改合同模板
// @Tags	合同
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string									true	"管理员校验token"
// @Param	body			body		definition.ContractTemplateModifyReq	true	"desc"
// @Success	200				{object}	model.StatusResponse					"请求成功"
func (*contractTemplate) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.ContractTemplateModifyReq](c)
	return ctx.SendResponse(biz.NewContractTemplate().Modify(req))
}
