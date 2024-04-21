package papi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type referralsProgress struct {
}

var ReferralsProgress = new(referralsProgress)

// List
// @ID		ReferralsProgressList
// @Router	/promotion/v1/referrals/progress [GET]
// @Summary	P8001 推荐进度
// @Tags	[P]推广接口
// @Accept	json
// @Produce	json
// @Param	X-Promotion-Token	header		string							true	"会员校验token"
// @Param	body				body		promotion.ReferralsProgressReq	true	"查询请求"
// @Success	200					{object}	promotion.ReferralsProgressRes	"请求成功"
func (*referralsProgress) List(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.ReferralsProgressReq](c)
	return ctx.SendResponse(service.NewPromotionReferralsService().ReferralsProgressList(ctx, &promotion.ReferralsProgressReq{
		PaginationReq:           req.PaginationReq,
		MemberID:                &ctx.Member.ID,
		ReferralsProgressFilter: req.ReferralsProgressFilter,
	}))
}
