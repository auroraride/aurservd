// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type branch struct {
}

var Branch = new(branch)

// Add 新增网点
func (*branch) Add(c echo.Context) (err error) {
    req := new(model.Branch)
    app.GetManagerContext(c).BindValidate(req)
    service.NewBranch().Add(req)
    return app.NewResponse(c).Send()
}
