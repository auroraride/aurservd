package papi

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
//	@ID			PromotionSetting
//	@Router		/promotion/v1/setting/{key} [GET]
//	@Summary	P7001 获取推广配置
//	@Tags		[P]推广接口
//	@Accept		json
//	@Produce	json
//	@Param		X-Promotion-Token	header		string					true	"管理员校验token"
//	@Param		body				body		promotion.SettingReq	true	"查询请求"
//	@Success	200					{object}	[]promotion.Setting
func (p *promotionSetting) Setting(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.SettingReq](c)
	return ctx.SendResponse(service.NewPromotionSettingService().Setting(req))
}
