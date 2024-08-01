package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assetTransfer struct{}

var AssetTransfer = new(assetTransfer)

// Transfer
// @ID		AssetTransfer
// @Router	/manager/v2/asset/transfer [POST]
// @Summary	资产调拨
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.AssetTransferCreateReq	true	"调拨参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*assetTransfer) Transfer(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetTransferCreateReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().Transfer(ctx.Request().Context(), req, ctx.Modifier))
}
