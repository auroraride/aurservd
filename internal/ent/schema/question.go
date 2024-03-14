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

type QuestionMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m QuestionMixin) Fields() []ent.Field {
	relate := field.Uint64("question_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m QuestionMixin) Edges() []ent.Edge {
	e := edge.To("question", Question.Type).Unique().Field("question_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m QuestionMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("question_id"))
	}
	return
}

// Question holds the schema definition for the Question entity.
type Question struct {
	ent.Schema
}

// Annotations of the Question.
func (Question) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "question"},
		entsql.WithComments(true),
	}
}

// Fields of the Question.
func (Question) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Comment("问题名称"),
		field.Int("sort").Default(0).Comment("排序"),
		field.String("answer").NotEmpty().Comment("答案"),
		field.Uint64("category_id").Optional().Comment("分类ID"),
	}
}

// Edges of the Question.
func (Question) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", QuestionCategory.Type).Ref("questions").
			Unique().Field("category_id"),
	}
}

func (Question) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Question) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("category_id"),
	}
}
