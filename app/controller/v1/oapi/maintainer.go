// Copyright (C) liasica. 2023-present.
//
// Created at 2023-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package oapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type maintainer struct{}

var Maintainer = new(maintainer)

// Signin
// @ID           AgentMaintainerSignin
// @Router       /maintainer/v1/maintainer/signin [POST]
// @Summary      O1001 运维登录
// @Tags         [O]运维接口
// @Accept       json
// @Produce      json
// @Param        X-Maintainer-Token  header  string  true  "运维校验token"
// @Param        body  body  model.MaintainerSigninReq  true  "请求参数"
// @Success      200  {object}  model.MaintainerSigninRes  "请求成功"
func (*maintainer) Signin(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[model.MaintainerSigninReq](c)
	return ctx.SendResponse(service.NewMaintainer().Signin(req))
}
