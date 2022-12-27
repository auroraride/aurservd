package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "time"
)

// Role holds the schema definition for the Role entity.
type Role struct {
    ent.Schema
}

// Annotations of the Role.
func (Role) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "role"},
        entsql.WithComments(true),
    }
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Unique().Comment("角色"),
        field.Strings("permissions").Optional().Comment("权限列表"),
        field.Bool("buildin").Default(false).Comment("是否内置角色"),
        field.Bool("super").Default(false).Comment("是否超级管理员"),
        field.Time("created_at").Immutable().Default(time.Now),
    }
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("managers", Manager.Type),
    }
}

func (Role) Mixin() []ent.Mixin {
    return []ent.Mixin{}
}

func (Role) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("name"),
        index.Fields("buildin"),
    }
}
