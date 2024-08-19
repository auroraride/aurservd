package app

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
)

// AssetManagerContext 管理员上下文
type AssetManagerContext struct {
	*BaseContext
	AssetManager *ent.AssetManager
	Modifier     *model.Modifier
}

// NewAssetManagerContext 新建管理员上下文
func NewAssetManagerContext(c echo.Context, mgr *ent.AssetManager, m *model.Modifier) *AssetManagerContext {
	return &AssetManagerContext{
		BaseContext:  Context(c),
		AssetManager: mgr,
		Modifier:     m,
	}
}

// GetAssetManagerContext 获取管理端上下文
func GetAssetManagerContext(c echo.Context) *AssetManagerContext {
	return c.(*AssetManagerContext)
}

// AssetManagerContextAndBinding 管理端上下文绑定数据
func AssetManagerContextAndBinding[T any](c echo.Context) (*AssetManagerContext, *T) {
	return ContextBindingX[AssetManagerContext, T](c)
}
