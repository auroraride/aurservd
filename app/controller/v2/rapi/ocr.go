// Copyright (C) liasica. 2024-present.
//
// Created at 2024-02-28
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type ocr struct{}

var Ocr = new(ocr)

// Signature
// @ID		Signature
// @Router	/certification/ocr/signature [GET]
// @Summary	获取阿里云OCR签名
// @Tags	Person - 实人
// @Accept	json
// @Produce	json
// @Param	query	query		definition.PersonCertificationOcrSignatureRequest	true	"请求参数"
// @Success	200		{object}	definition.PersonCertificationOcrSignatureResponse	"请求成功"
func (*ocr) Signature(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.PersonCertificationOcrSignatureRequest](c)
	return ctx.SendResponse(biz.NewPerson().Signature(req.Hash))
}
