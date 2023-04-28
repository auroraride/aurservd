// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/labstack/echo/v4"
)

type contract struct {
}

var Contract = new(contract)

// Sign
// @ID           RiderContractSign
// @Router       /rider/v1/contract/sign [POST]
// @Summary      R3003 签署合同
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.ContractSignReq  true  "desc"
// @Success      200  {object}  model.ContractSignRes  "请求成功"
func (*contract) Sign(c echo.Context) error {
	ctx, req := app.RiderContextAndBinding[model.ContractSignReq](c)
	return ctx.SendResponse(service.NewContractWithRider(ctx.Rider).Sign(req))
}

// SignResult
// @ID           RiderContractSignResult
// @Router       /rider/v1/constract/{sn} [GET]
// @Summary      R3004 合同签署结果
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        sn  path  string  true  "合同编号"
// @Success      200 {object}  model.StatusResponse  "请求成功"
func (*contract) SignResult(c echo.Context) error {
	ctx, req := app.RiderContextAndBinding[model.ContractSignResultReq](c)
	return ctx.SendResponse(service.NewContract().Result(ctx.Rider, req.Sn))
}
