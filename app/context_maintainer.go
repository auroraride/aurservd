package app

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/internal/ent"
)

type MaintainerContext struct {
	*BaseContext

	Maintainer *ent.Maintainer
	Cities     ent.Cities
}

func NewMaintainerContext(c echo.Context, m *ent.Maintainer) *MaintainerContext {
	ctx := &MaintainerContext{
		BaseContext: Context(c),
		Maintainer:  m,
	}
	if m != nil {
		ctx.Cities = m.Edges.Cities
	}
	return ctx
}

func MaintainerContextAndBinding[T any](c echo.Context) (*MaintainerContext, *T) {
	return ContextBindingX[MaintainerContext, T](c)
}

// CityIDs 城市ID列表
func (c *MaintainerContext) CityIDs() (ids []uint64) {
	for _, city := range c.Cities {
		ids = append(ids, city.ID)
	}
	return
}