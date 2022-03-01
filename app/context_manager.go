// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/labstack/echo/v4"
)

// ManagerContext 管理员上下文
type ManagerContext struct {
    *Context

    Manager *ent.Manager
}

func GetManagerContext(c echo.Context) *ManagerContext  {
    return c.(*ManagerContext)
}