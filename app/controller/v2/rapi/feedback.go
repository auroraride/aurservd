package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/model"
)

var Feedback = new(feedback)

type feedback struct{}

// Create
// @ID		FeedbackCreate
// @Router	/rider/v2/feedback [POST]
// @Summary	意见反馈
// @Tags	Feedback - 反馈
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	body			body		model.FeedbackReq		true	"反馈内容"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*feedback) Create(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.FeedbackReq](c)
	biz.NewFeedback().RiderCreate(req, ctx.Rider)
	return ctx.SendResponse()
}
