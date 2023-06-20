// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-20
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type bike struct{}

var Bike = new(bike)

// Unallocated
// @ID           AgentBikeUnallocated
// @Router       /agent/v1/bike/unallocated [GET]
// @Summary      AB003 搜索未分配车辆
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        query  query  model.EbikeAgentSearchReq  true  "请求参数"
// @Success      200  {object}  model.Ebike  "请求成功"
func (*bike) Unallocated(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.EbikeAgentSearchReq](c)
	// TODO 二级代理校验站点权限
	return ctx.SendResponse(service.NewEbike().Unallocated(&model.EbikeUnallocatedParams{
		Keyword:   req.Keyword,
		StationID: req.StationID,
	}))
}
