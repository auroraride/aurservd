// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-11
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type person struct{}

var Person = new(person)

// CertificationOcr
// @ID           CertificationOcr
// @Router       /v2/certification/ocr [GET]
// @Summary      获取人脸核身OCR参数
// @Tags         Person - 实人
// @Accept       json
// @Produce      json
// @Success      200  {object}  definition.PersonCertificationOcrRes  "请求成功"
func (*person) CertificationOcr(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(biz.NewPerson().CertificationOcr(ctx.Rider))
}

// CertificationFace
// @ID           CertificationFace
// @Router       /v2/certification/face [GET]
// @Summary      获取人脸核身参数
// @Tags         Person - 实人
// @Accept       json
// @Produce      json
// @Param        orderNo query string true "订单编号"
// @Success      200  {object}  definition.PersonCertificationOcrRes  "请求成功"
func (*person) CertificationFace(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.PersonCertificationFaceReq](c)
	return ctx.SendResponse(biz.NewPerson().CertificationFace(ctx.Rider, req.OrderNo))
}
