// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-05
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type bill struct{}

var Bill = new(bill)

// Usage
// @ID           AgentBillUsage
// @Router       /agent/v1/bill/usage [GET]
// @Summary      A4001 使用明细
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        query  query   model.StatementUsageReq  true  "筛选项"
// @Success      200  {object}  model.Pagination{items=[]model.StatementUsageRes}  "请求成功"
func (*bill) Usage(c echo.Context) (err error) {
    ctx, req := app.AgentContextAndBinding[model.StatementUsageReq](c)
    req.StatementUsageFilter.ID = ctx.Enterprise.ID
    // req.StatementUsageFilter.ID = 42949672960
    return ctx.SendResponse(service.NewEnterpriseStatement().Usage(req))
}
