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
    "github.com/auroraride/aurservd/app/model"
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
        field.String("remark").Nillable().Optional().Comment("备注"),
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
