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
    rider    *ent.Rider
    modifier *model.Modifier
    employee *model.Employee
}

func NewService(params ...any) (bs *BaseService) {
    bs = &BaseService{}
    ctx := context.Background()
    for _, param := range params {
        switch p := param.(type) {
        case *ent.Rider:
            ctx = context.WithValue(ctx, "rider", p)
            bs.rider = p
        case *ent.Manager:
            m := &model.Modifier{
                ID:    p.ID,
                Name:  p.Name,
                Phone: p.Phone,
            }
            bs.modifier = m
            ctx = context.WithValue(ctx, "modifier", m)
        case *model.Modifier:
            bs.modifier = p
            ctx = context.WithValue(ctx, "modifier", p)
        case *ent.Employee:
            e := &model.Employee{
                ID:    p.ID,
                Name:  p.Name,
                Phone: p.Phone,
            }
            bs.employee = e
            ctx = context.WithValue(ctx, "employee", e)
        }
    }

    bs.ctx = ctx

    return
}
