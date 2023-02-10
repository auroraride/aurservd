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

type TaskMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m TaskMixin) Fields() []ent.Field {
    relate := field.Uint64("task_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m TaskMixin) Edges() []ent.Edge {
    e := edge.To("task", Task.Type).Unique().Field("task_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m TaskMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("task_id"))
    }
    return
}

// Task holds the schema definition for the Task entity.
type Task struct {
    ent.Schema
}

// Annotations of the Task.
func (Task) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "task"},
        entsql.WithComments(true),
    }
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
    return []ent.Field{}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (Task) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},

        CabinetMixin{},
        RiderMixin{},
        ExchangeMixin{Optional: true},
    }
}

func (Task) Indexes() []ent.Index {
    return []ent.Index{}
}
