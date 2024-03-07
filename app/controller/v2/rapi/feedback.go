package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/service"
)

var Feedback = new(feedback)

type feedback struct{}

// Create
// @ID		Create
// @Router	/rider/v2/feedback [POST]
// @Summary	意见反馈
// @Tags	Feedback - 反馈
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	body			body		definition.FeedbackReq	true	"反馈内容"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*feedback) Create(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.FeedbackReq](c)
	biz.NewFeedback().RiderCreate(req, ctx.Rider)
	return ctx.SendResponse()
}

// UploadImage
// @ID		UploadImage
// @Router	/rider/v2/feedback/image [POST]
// @Summary	意见反馈上传图片
// @Tags	Feedback - 反馈
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string		true	"代理校验token"
// @Param	images			formData	file		true	"图片文件"
// @Success	200				{object}	[]string	"请求成功"
func (*feedback) UploadImage(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(service.NewFeedback().UploadImage(ctx))
}
