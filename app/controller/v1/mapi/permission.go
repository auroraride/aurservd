// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-05
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    perm "github.com/auroraride/aurservd/app/permission"
    "github.com/labstack/echo/v4"
    "sort"
    "strings"
)

type permission struct{}

var Permission = new(permission)

// List
// @ID           ManagerPermissionList
// @Router       /manager/v1/permission [GET]
// @Summary      MD001 权限列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []perm.Group  "请求成功"
func (*permission) List(c echo.Context) (err error) {
    ctx := app.Context(c)
    items := make([]*perm.Group, 0)
    for _, g := range perm.Groups {
        items = append(items, g)
    }

    sort.Slice(items, func(i, j int) bool {
        return strings.Compare(items[i].Name, items[j].Name) < 0
    })

    return ctx.SendResponse(items)
}
