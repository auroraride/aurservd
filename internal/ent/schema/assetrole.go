package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// AssetRole holds the schema definition for the AssetRole entity.
type AssetRole struct {
	ent.Schema
}

// Annotations of the AssetRole.
func (AssetRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_role"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetRole.
func (AssetRole) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().Comment("角色"),
		field.Strings("permissions").Optional().Comment("权限列表"),
		field.Bool("buildin").Default(false).Comment("是否内置角色"),
		field.Bool("super").Default(false).Comment("是否超级管理员"),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

// Edges of the AssetRole.
func (AssetRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("asset_managers", AssetManager.Type),
	}
}

func (AssetRole) Mixin() []ent.Mixin {
	return []ent.Mixin{}
}

func (AssetRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("buildin"),
	}
}
