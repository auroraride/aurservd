// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-11
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type person struct{}

var Person = new(person)

func (*person) Certification(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.PersonCertificationReq](c)
	return ctx.SendResponse(biz.NewPerson().Certification(req))
}

func (*person) CertificationCallback(c echo.Context) (err error) {
	fmt.Println(c.Request().RequestURI)
	return c.String(http.StatusOK, "OK")
}

func (*person) CertificationResult(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.PersonCertification](c)
	return ctx.SendResponse(biz.NewPerson().CertificationResult(req))
}
