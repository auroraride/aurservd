// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package common

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/internal/ali"
)

type oss struct {
}

var Oss = new(oss)

// Token
// @ID			AliyunOssToken
// @Router		/common/oss/token [get]
// @Summary		C3 获取阿里云oss临时凭证
// @Description	上传文件必须，单次获取有效时间为1个小时
// @Tags		Communal - 公共接口
// @Accept		json
// @Produce		json
// @Success		200	{object}	model.AliyunOssStsRes	"请求成功"
func (*oss) Token(c echo.Context) error {
	return app.Context(c).SendResponse(ali.NewOss().StsToken())
}
