package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assetAttributes struct{}

var AssetAttributes = new(assetAttributes)

// List
// @ID		AssetAttributesList
// @Router	/manager/v2/asset/attributes [GET]
// @Summary	资产属性列表
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	query					query		model.AssetAttributesListReq	true	"查询参数"
// @Success	200						{object}	model.AssetAttributesListRes	"请求成功"
func (*assetAttributes) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetAttributesListReq](c)
	return ctx.SendResponse(service.NewAssetAttributes().List(ctx.Request().Context(), req))
}
