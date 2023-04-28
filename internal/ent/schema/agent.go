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

type AgentMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AgentMixin) Fields() []ent.Field {
	f := field.Uint64("agent_id")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m AgentMixin) Edges() []ent.Edge {
	e := edge.To("agent", Agent.Type).Unique().Field("agent_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AgentMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("agent_id"))
	}
	return
}

// Agent holds the schema definition for the Agent entity.
type Agent struct {
	ent.Schema
}

// Annotations of the Agent.
func (Agent) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "agent"},
		entsql.WithComments(true),
	}
}

// Fields of the Agent.
func (Agent) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("phone").Unique(),
		field.String("password"),
	}
}

// Edges of the Agent.
func (Agent) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Agent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
		EnterpriseMixin{},
	}
}

func (Agent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("phone"),
	}
}
