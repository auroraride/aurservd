package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionSetting struct {
}

var PromotionSetting = new(promotionSetting)

// Setting
// @ID		ManagerPromotionSetting
// @Router	/manager/v1/promotion/setting/{key} [GET]
// @Summary	获取推广配置
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		promotion.SettingReq	true	"查询请求"
// @Success	200				{object}	[]promotion.Setting
func (p *promotionSetting) Setting(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.SettingReq](c)
	return ctx.SendResponse(service.NewPromotionSettingService().Setting(req))
}

// Update
// @ID		ManagerPromotionSettingUpdate
// @Router	/manager/v1/promotion/setting/{key} [PUT]
// @Summary	更新推广配置
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string				true	"管理员校验token"
// @Param	body			body		promotion.Setting	true	"更新请求参数"
// @Success	200				{object}	model.StatusResponse
func (p *promotionSetting) Update(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.Setting](c)
	service.NewPromotionSettingService(ctx.Modifier).Update(req)
	return ctx.SendResponse()
}
