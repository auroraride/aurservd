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

type ActivityMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m ActivityMixin) Fields() []ent.Field {
	relate := field.Uint64("activity_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m ActivityMixin) Edges() []ent.Edge {
	e := edge.To("activity", Activity.Type).Unique().Field("activity_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m ActivityMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("activity_id"))
	}
	return
}

// Activity holds the schema definition for the Activity entity.
type Activity struct {
	ent.Schema
}

// Annotations of the Activity.
func (Activity) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "activity"},
		entsql.WithComments(true),
	}
}

// Fields of the Activity.
func (Activity) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("名称"),
		field.String("image").Comment("图片"),
		field.String("link").Comment("连接"),
		field.Int("sort").Default(0).Comment("排序"),
	}
}

// Edges of the Activity.
func (Activity) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Activity) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Activity) Indexes() []ent.Index {
	return []ent.Index{}
}
