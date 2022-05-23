// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "context"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/mixin"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
)

// Modifier 修改或创建人
type Modifier struct {
    mixin.Schema
}

func (Modifier) Fields() []ent.Field {
    return []ent.Field{
        field.JSON("creator", &model.Modifier{}).Immutable().Optional().Comment("创建人"),
        field.JSON("last_modifier", &model.Modifier{}).Optional().Comment("最后修改人"),
        field.String("remark").Optional().Comment("备注"),
    }
}

func (Modifier) Hooks() []ent.Hook {
    return []ent.Hook{
        func(next ent.Mutator) ent.Mutator {
            return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
                ml, ok := m.(model.ModifierLogger)
                if !ok {
                    return nil, fmt.Errorf("unexpected audit-log call from mutation type %T", m)
                }
                mod := model.GetModifierFromContext(ctx)
                if mod != nil {
                    switch op := m.Op(); {
                    case op.Is(ent.OpCreate):
                        ml.SetCreator(mod)
                        ml.SetLastModifier(mod)
                    case op.Is(ent.OpUpdateOne | ent.OpUpdate):
                        ml.SetLastModifier(mod)
                    }
                }
                return next.Mutate(ctx, m)
            })
        },
    }
}
