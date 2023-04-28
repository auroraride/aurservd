// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-15
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
	"context"
	"fmt"

	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/ebike"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
)

type ExportInfoData interface {
	GetExportInfo() string
}

func (r *Rider) GetExportInfo() string {
	p := r.Edges.Person
	str := ""
	if p != nil {
		str = p.Name + " - "
	}
	return str + r.Phone
}

func (c *City) GetExportInfo() string {
	return c.Name
}

func (e *Employee) GetExportInfo() string {
	return e.Name
}

func (pl *Plan) GetExportInfo() string {
	return fmt.Sprintf("%s - %d天", pl.Name, pl.Days)
}

func (e *Enterprise) GetExportInfo() string {
	return e.Name
}

func (c *Cabinet) GetExportInfo() string {
	return c.Name + " - " + c.Serial
}

func (s *Store) GetExportInfo() string {
	return s.Name
}

func (e *Ebike) GetExportInfo() string {
	return e.Sn
}

func (eb *EbikeBrand) GetExportInfo() string {
	return eb.Name
}

func (b *Battery) GetExportInfo() string {
	return b.Sn
}

type ExportInfo struct {
	id    uint64
	table string
}

func NewExportInfo(id uint64, t string) *ExportInfo {
	return &ExportInfo{
		id:    id,
		table: t,
	}
}

func (ei *ExportInfo) GetExportInfoData() string {
	ctx := context.Background()
	var m ExportInfoData
	switch ei.table {
	case rider.Table:
		m, _ = Database.Rider.QueryNotDeleted().Where(rider.ID(ei.id)).WithPerson().First(ctx)
	case city.Table:
		m, _ = Database.City.QueryNotDeleted().Where(city.ID(ei.id)).First(ctx)
	case employee.Table:
		m, _ = Database.Employee.QueryNotDeleted().Where(employee.ID(ei.id)).First(ctx)
	case plan.Table:
		m, _ = Database.Plan.QueryNotDeleted().Where(plan.ID(ei.id)).First(ctx)
	case enterprise.Table:
		m, _ = Database.Enterprise.QueryNotDeleted().Where(enterprise.ID(ei.id)).First(ctx)
	case cabinet.Table:
		m, _ = Database.Cabinet.QueryNotDeleted().Where(cabinet.ID(ei.id)).First(ctx)
	case store.Table:
		m, _ = Database.Store.QueryNotDeleted().Where(store.ID(ei.id)).First(ctx)
	case ebike.Table:
		m, _ = Database.Ebike.Query().Where(ebike.ID(ei.id)).First(ctx)
	case ebikebrand.Table:
		m, _ = Database.EbikeBrand.Query().Where(ebikebrand.ID(ei.id)).First(ctx)
	case battery.Table:
		m, _ = Database.Battery.Query().Where(battery.ID(ei.id)).First(ctx)
	}
	if m == nil {
		return ""
	}
	return m.GetExportInfo()
}
