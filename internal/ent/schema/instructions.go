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

type InstructionsMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m InstructionsMixin) Fields() []ent.Field {
	relate := field.Uint64("instructions_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m InstructionsMixin) Edges() []ent.Edge {
	e := edge.To("instructions", Instructions.Type).Unique().Field("instructions_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m InstructionsMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("instructions_id"))
	}
	return
}

// Instructions holds the schema definition for the Instructions entity.
type Instructions struct {
	ent.Schema
}

// Annotations of the Instructions.
func (Instructions) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "instructions"},
		entsql.WithComments(true),
	}
}

// Fields of the Instructions.
func (Instructions) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").Comment("标题"),
		field.JSON("content", new(interface{})).Comment("内容"),
		field.String("key").Comment("标识"),
	}
}

// Edges of the Instructions.
func (Instructions) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Instructions) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Instructions) Indexes() []ent.Index {
	return []ent.Index{}
}
