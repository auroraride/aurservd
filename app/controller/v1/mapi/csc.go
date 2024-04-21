// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/pkg/snag"
)

type csc struct{}

var Csc = new(csc)

// BatchReminder
// @ID		CscBatchReminder
// @Router	/manager/v1/csc/irv [POST]
// @Summary	时光驹催费工具
// @Tags	催费
// @Accept	mpfd
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	file			formData	file					true	"外呼文件"
// @Success	200				{object}	model.ShiguangjuIVRRes	"请求成功"
func (*csc) BatchReminder(c echo.Context) (err error) {
	file, err := c.FormFile("file")
	if err != nil {
		snag.Panic(err)
	}

	return app.ContextX[app.ManagerContext](c).SendResponse(service.NewCSC().BatchReminder(file))
}
