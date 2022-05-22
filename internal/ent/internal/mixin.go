// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/10
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "context"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "fmt"
    "github.com/auroraride/aurservd/internal/ent/hook"
    "github.com/sony/sonyflake"
)

// DeleteMixin 删除字段
type DeleteMixin struct {
    mixin.Schema
}

func (DeleteMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Time("deleted_at").Nillable().Optional(),
    }
}

func (DeleteMixin) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("deleted_at"),
    }
}

type SonyflakeIDMixin struct {
    mixin.Schema
}

func (SonyflakeIDMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("id"),
    }
}

// Hooks of the Mixin.
func (SonyflakeIDMixin) Hooks() []ent.Hook {
    return []ent.Hook{
        hook.On(IDHook(), ent.OpCreate),
    }
}

func IDHook() ent.Hook {
    sf := sonyflake.NewSonyflake(sonyflake.Settings{})
    type IDSetter interface {
        SetID(uint64)
    }
    return func(next ent.Mutator) ent.Mutator {
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            is, ok := m.(IDSetter)
            if !ok {
                return nil, fmt.Errorf("unexpected mutation %T", m)
            }
            id, err := sf.NextID()
            if err != nil {
                return nil, err
            }
            is.SetID(id)
            return next.Mutate(ctx, m)
        })
    }
}
