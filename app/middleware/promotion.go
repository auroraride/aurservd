package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/snag"
)

func Promotion() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				me *ent.PromotionMember
			)

			// 查看是否是登录态
			token := strings.TrimSpace(c.Request().Header.Get(app.HeaderPromotionToken))
			if token != "" {
				// 查找登录用户
				me = service.NewPromotionMemberService().TokenVerify(token)
			}
			return next(app.NewPromotionContext(c, me))
		}
	}
}

func PromotionAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.(*app.PromotionContext)
			if ctx.Member == nil {
				snag.Panic(snag.StatusUnauthorized)
			}
			return next(ctx)
		}
	}
}
