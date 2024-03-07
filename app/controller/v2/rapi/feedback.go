package rapi

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/service"
	"github.com/labstack/echo/v4"
)

var Feedback = new(feedback)

type feedback struct{}

// Feedback
// @ID		RiderFeedback
// @Router	/rider/v2/feedback [POST]
// @Summary	意见反馈
// @Tags	Rider - 骑手接口
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string				true	"骑手校验token"
// @Param	body			body		model.FeedbackReq	true	"反馈内容"
// @Success	200				{object}	bool				"请求成功"
func (*feedback) Feedback(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.FeedbackReq](c)
	return ctx.SendResponse(biz.NewFeedback().RiderCreate(req, ctx.Rider))
}

// FeedbackImage
// @ID		RiderFeedbackImage
// @Router	/rider/v2/feedback/image [POST]
// @Summary	意见反馈上传图片
// @Tags	Rider - 图片接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string		true	"代理校验token"
// @Param	images			formData	file		true	"图片文件"
// @Success	200				{object}	[]string	"请求成功"
func (*feedback) FeedbackImage(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(service.NewFeedback().UploadImage(ctx))
}
