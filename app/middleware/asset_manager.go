package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	assetpermission "github.com/auroraride/aurservd/app/assetpermission"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
)

var (
	assetManagerSkipper = map[string]bool{
		"/manager/v2/asset/permission":  true,
		"/manager/v2/asset/user/signin": true,
	}
)

// AssetManagerMiddleware 后台中间件
func AssetManagerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p := c.Path()
			if assetManagerSkipper[p] {
				return next(c)
			}

			// 判定登录
			token := c.Request().Header.Get(app.HeaderAssetManagerToken)
			if token == "" {
				token = c.QueryParam("token")
			}
			id, err := cache.Get(context.Background(), token).Uint64()
			if err != nil {
				snag.Panic(snag.StatusUnauthorized)
			}
			s := biz.NewAssetManager()
			var m *ent.AssetManager
			m, err = s.GetAssetManagerById(id)
			if err != nil || m == nil {
				snag.Panic(snag.StatusUnauthorized)
			}

			// 延长token有效期
			s.ExtendTokenTime(m.ID, token)

			perms, _ := s.GetAssetPermissions(m)
			if !assetpermission.Contains(strings.ToUpper(c.Request().Method), p, perms) {
				snag.Panic(snag.StatusForbidden)
			}

			// 重载context
			return next(app.NewAssetManagerContext(c, m, &model.Modifier{
				ID:    m.ID,
				Name:  m.Name,
				Phone: m.Phone,
			}))
		}
	}
}
