package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/snag"
)

func Maintainer() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				m *ent.Maintainer
			)

			// 查看是否是登录态
			token := strings.TrimSpace(c.Request().Header.Get(app.HeaderMaintainerToken))
			if token != "" {
				// 查找登录用户
				m = service.NewMaintainer().TokenVerify(token)
			}
			return next(app.NewMaintainerContext(c, m))
		}
	}
}

func MaintainerAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.(*app.MaintainerContext)
			if ctx.Maintainer == nil {
				snag.Panic(snag.StatusUnauthorized)
			}
			return next(ctx)
		}
	}
}
