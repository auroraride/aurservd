package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Guide holds the schema definition for the Guide entity.
type Guide struct {
	ent.Schema
}

// Fields of the Guide.
func (Guide) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("名称"),
		field.Uint8("sort").Default(0).Comment("排序"),
		field.String("answer").Comment("答案"),
		field.String("remark").Optional().Comment("备注"),
		field.Time("created_at").Immutable(),
		field.Time("updated_at"),
	}
}

// Edges of the Guide.
func (Guide) Edges() []ent.Edge {
	return []ent.Edge{}
}
