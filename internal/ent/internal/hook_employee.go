// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-27
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "context"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/app/model"
)

type HookEmployee struct {
    mixin.Schema
}

type HookEmployeeMutator interface {
    SetEmployeeInfo(value *model.Employee)
}

type HookEmployeeIDMutator interface {
    SetEmployeeID(v uint64)
    EmployeeID() (uint64, bool)
}

// Fields of the PointLog.
func (HookEmployee) Fields() []ent.Field {
    return []ent.Field{
        field.JSON("employee_info", &model.Employee{}).Optional().Comment("店员"),
    }
}

func (h HookEmployee) Hooks() []ent.Hook {
    return []ent.Hook{
        func(next ent.Mutator) ent.Mutator {
            return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
                h.setEmployeeInfo(ctx, m)
                h.setEmployeeID(ctx, m)
                return next.Mutate(ctx, m)
            })
        },
    }
}

func (HookEmployee) setEmployeeInfo(ctx context.Context, m ent.Mutation) {
    ml, ok := m.(HookEmployeeMutator)
    if ok {
        value, _ := ctx.Value("employee").(*model.Employee)
        if value != nil {
            switch op := m.Op(); {
            case op.Is(ent.OpCreate):
                ml.SetEmployeeInfo(value)
            case op.Is(ent.OpUpdateOne | ent.OpUpdate):
                // TODO: 更新?
            }
        }
    }
}

// TODO 判定 nilableID
func (HookEmployee) setEmployeeID(ctx context.Context, m ent.Mutation) {
    ml, ok := m.(HookEmployeeIDMutator)
    if ok {
        eid, idOk := ml.EmployeeID()
        if eid == 0 || !idOk {
            value, _ := ctx.Value("employee").(*model.Employee)
            if value != nil {
                switch op := m.Op(); {
                case op.Is(ent.OpCreate):
                    ml.SetEmployeeID(value.ID)
                case op.Is(ent.OpUpdateOne | ent.OpUpdate):
                    // TODO: 更新?
                }
            }
        }
    }
}
