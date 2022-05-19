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
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/hook"
    "github.com/sony/sonyflake"
    "time"
)

// TimeMixin 时间字段
type TimeMixin struct {
    mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Time("created_at").Immutable().Default(time.Now),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}

func (TimeMixin) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("created_at"),
    }
}

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

// LastModifier 上次修改人
type LastModifier struct {
    mixin.Schema
}

func (LastModifier) Fields() []ent.Field {
    return []ent.Field{
        field.JSON("last_modifier", &model.Modifier{}).Optional().Comment("最后修改人"),
        field.String("remark").Optional().Comment("备注"),
    }
}

// Creator 创建人
type Creator struct {
    mixin.Schema
}

func (Creator) Fields() []ent.Field {
    return []ent.Field{
        field.JSON("creator", &model.Modifier{}).Optional().Comment("创建人"),
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
