// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/10
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "time"
)

type TimeMixin struct {
    mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Time("created_at").Immutable().Default(time.Now),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}

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