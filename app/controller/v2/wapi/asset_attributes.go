package wapi

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
// @Router	/warestore/v2/assets/attributes [GET]
// @Summary	资产属性列表
// @Tags	Assets - 资产管理
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"管理员校验token"
// @Param	query				query		model.AssetAttributesListReq	true	"查询参数"
// @Success	200					{object}	model.AssetAttributesListRes	"请求成功"
func (*assetAttributes) List(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[model.AssetAttributesListReq](c)
	return ctx.SendResponse(service.NewAssetAttributes().List(ctx.Request().Context(), req))
}
