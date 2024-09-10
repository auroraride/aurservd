// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-22, by aurb

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type StoreGroupMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m StoreGroupMixin) Fields() []ent.Field {
	relate := field.Uint64("group_id").Comment("城市ID")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m StoreGroupMixin) Edges() []ent.Edge {
	e := edge.To("group", StoreGroup.Type).Unique().Field("group_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m StoreGroupMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("group_id"))
	}
	return
}

// StoreGroup holds the schema definition for the StoreGroup entity.
type StoreGroup struct {
	ent.Schema
}

// Annotations of the StoreGroup.
func (StoreGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "store_group"},
		entsql.WithComments(true),
	}
}

// Fields of the StoreGroup.
func (StoreGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("名称"),
	}
}

// Edges of the StoreGroup.
func (StoreGroup) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (StoreGroup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (StoreGroup) Indexes() []ent.Index {
	return []ent.Index{}
}
