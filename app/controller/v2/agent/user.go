// Copyright (C) liasica. 2023-present.
//
// Created at 2023-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package agent

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type user struct{}

var User = new(user)

func (*user) Signin(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentSigninReq](c)
	return ctx.SendResponse(service.NewAgent().Signin(req))
}
