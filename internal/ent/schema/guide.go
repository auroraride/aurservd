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

type GuideMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m GuideMixin) Fields() []ent.Field {
	relate := field.Uint64("guide_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m GuideMixin) Edges() []ent.Edge {
	e := edge.To("guide", Guide.Type).Unique().Field("guide_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m GuideMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("guide_id"))
	}
	return
}

// Guide holds the schema definition for the Guide entity.
type Guide struct {
	ent.Schema
}

// Annotations of the Guide.
func (Guide) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "guide"},
		entsql.WithComments(true),
	}
}

// Fields of the Guide.
func (Guide) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("名称"),
		field.Uint8("sort").Default(0).Comment("排序"),
		field.String("answer").Comment("答案"),
	}
}

// Edges of the Guide.
func (Guide) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Guide) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Guide) Indexes() []ent.Index {
	return []ent.Index{}
}
