package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

// Feedback holds the schema definition for the Feedback entity.
type Feedback struct {
	ent.Schema
}

// Fields of the Feedback.
func (Feedback) Fields() []ent.Field {
	return []ent.Field{
		field.String("content").Comment("反馈内容"),
		field.Uint8("type").Default(0).Comment("反馈类型"),
		field.JSON("url", []string{}).Optional().Comment("反馈图片"),
		field.String("name").Optional().Comment("姓名"),
		field.String("phone").Optional().Comment("电话"),
	}
}

// Edges of the Feedback.
func (Feedback) Edges() []ent.Edge {
	return []ent.Edge{}
}
func (Feedback) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		EnterpriseMixin{Optional: true},
		AgentMixin{Optional: true},
	}
}

func (Feedback) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "feedback"},
	}
}
