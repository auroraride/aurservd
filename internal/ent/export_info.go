// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-15
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/enterprise"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/internal/ent/rider"
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
    return fmt.Sprintf("%s - %d", pl.Name, pl.Days)
}

func (e *Enterprise) GetExportInfo() string {
    return e.Name
}

func (c *Cabinet) GetExportInfo() string {
    return c.Name + " - " + c.Serial
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
        break
    case city.Table:
        m, _ = Database.City.QueryNotDeleted().Where(city.ID(ei.id)).First(ctx)
        break
    case employee.Table:
        m, _ = Database.Employee.QueryNotDeleted().Where(employee.ID(ei.id)).First(ctx)
        break
    case plan.Table:
        m, _ = Database.Plan.QueryNotDeleted().Where(plan.ID(ei.id)).First(ctx)
        break
    case enterprise.Table:
        m, _ = Database.Enterprise.QueryNotDeleted().Where(enterprise.ID(ei.id)).First(ctx)
        break
    case cabinet.Table:
        m, _ = Database.Cabinet.QueryNotDeleted().Where(cabinet.ID(ei.id)).First(ctx)
        break
    }
    if m == nil {
        return ""
    }
    return m.GetExportInfo()
}