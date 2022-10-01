// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type BaseService struct {
    ctx      context.Context
    rider    *model.Rider
    modifier *model.Modifier
    employee *model.Employee
}

func newService(params ...any) (bs *BaseService) {
    bs = &BaseService{}
    ctx := context.Background()
    for _, param := range params {
        switch p := param.(type) {
        case *ent.Rider:
            bs.rider = &model.Rider{
                ID:    p.ID,
                Phone: p.Phone,
                Name:  p.Name,
            }
            ctx = context.WithValue(ctx, "rider", bs.rider)
        case *ent.Manager:
            bs.modifier = &model.Modifier{
                ID:    p.ID,
                Phone: p.Phone,
                Name:  p.Name,
            }
            ctx = context.WithValue(ctx, "modifier", bs.modifier)
        case *model.Modifier:
            bs.modifier = p
            ctx = context.WithValue(ctx, "modifier", bs.modifier)
        case *ent.Employee:
            bs.employee = &model.Employee{
                ID:    p.ID,
                Name:  p.Name,
                Phone: p.Phone,
            }
            ctx = context.WithValue(ctx, "employee", bs.employee)
        }
    }

    bs.ctx = ctx

    return
}
