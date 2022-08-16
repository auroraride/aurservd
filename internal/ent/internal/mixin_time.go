// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "time"
)

// TimeMixin 时间字段
type TimeMixin struct {
    mixin.Schema
    DoNotIndexCreatedAt bool
}

func (TimeMixin) Fields() []ent.Field {
    return []ent.Field{
        // .SchemaType(map[string]string{dialect.Postgres: "timestamp without time zone"})
        field.Time("created_at").Immutable().Default(time.Now),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}

func (t TimeMixin) Indexes() []ent.Index {
    var list []ent.Index
    if !t.DoNotIndexCreatedAt {
        list = append(list, index.Fields("created_at"))
    }
    return list
}
