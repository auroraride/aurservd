package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type asset struct{}

var Asset = new(asset)

// Count
// @ID		AssetCount
// @Router	/manager/v2/masset/count [GET]
// @Summary	查询有效的资产数量
// @Tags	Asset - 资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string				true	"管理员校验token"
// @Param	query			query		model.AssetFilter	true	"查询参数"
// @Success	200				{object}	model.AssetNumRes	"请求成功"
func (*asset) Count(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetFilter](c)
	return ctx.SendResponse(service.NewAsset().Count(ctx.Request().Context(), req))
}
