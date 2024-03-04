// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-28
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/labstack/echo/v4"
)

type contract struct{}

var Contract = new(contract)

// List
//	@ID			ManagerContractList
//	@Router		/manager/v1/contract [GET]
//	@Summary	M7022 合同列表
//	@Tags		[M]管理接口
//	@Accept		json
//	@Produce	json
//	@Param		X-Manager-Token	header		string												true	"管理员校验token"
//	@Param		qery			query		model.ContractListReq								false	"筛选选项"
//	@Success	200				{object}	model.PaginationRes{items=[]model.ContractListRes}	"请求成功"
func (*contract) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ContractListReq](c)
	return ctx.SendResponse(service.NewContract().List(req))
}
