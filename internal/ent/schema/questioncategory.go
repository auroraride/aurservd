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

type QuestionCategoryMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m QuestionCategoryMixin) Fields() []ent.Field {
	relate := field.Uint64("category_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m QuestionCategoryMixin) Edges() []ent.Edge {
	e := edge.To("category", QuestionCategory.Type).Unique().Field("category_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m QuestionCategoryMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("category_id"))
	}
	return
}

// QuestionCategory holds the schema definition for the QuestionCategory entity.
type QuestionCategory struct {
	ent.Schema
}

// Annotations of the QuestionCategory.
func (QuestionCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "question_category"},
		entsql.WithComments(true),
	}
}

// Fields of the QuestionCategory.
func (QuestionCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Comment("名称"),
		field.Int("sort").Default(0).Comment("排序"),
	}
}

// Edges of the QuestionCategory.
func (QuestionCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("questions", Question.Type).Comment("问题列表"),
	}
}

func (QuestionCategory) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (QuestionCategory) Indexes() []ent.Index {
	return []ent.Index{}
}
