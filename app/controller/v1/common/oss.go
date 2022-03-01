// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package common

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/labstack/echo/v4"
)

type oss struct {
}

var Oss = new(oss)

func (*oss) Token(c echo.Context) error {
    return app.NewResponse(c).SetData(ali.NewOss().StsToken()).Send()
}