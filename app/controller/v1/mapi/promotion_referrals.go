package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionReferrals struct {
}

var PromotionReferrals = new(promotionReferrals)

// ProgressList
//	@ID			ManagerProgressList
//	@Router		/manager/v1/promotion/progress/list [GET]
//	@Summary	PMA001  会员推荐进度列表
//	@Tags		[PM]推广管理接口
//	@Accept		json
//	@Produce	json
//	@Param		X-Manager-Token	header		string							true	"管理员校验token"
//	@Param		body			body		promotion.ReferralsProgressReq	true	"会员等级配置"
//	@Success	200				{object}	[]promotion.ReferralsProgressRes
func (r *promotionReferrals) ProgressList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.ReferralsProgressReq](c)
	return ctx.SendResponse(service.NewPromotionReferralsService().ReferralsProgressList(ctx, req))
}
