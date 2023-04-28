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

type RiderBelongMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m RiderBelongMixin) Fields() []ent.Field {
    relate := field.Uint64("belong_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m RiderBelongMixin) Edges() []ent.Edge {
    e := edge.To("belong", RiderBelong.Type).Unique().Field("belong_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m RiderBelongMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("belong_id"))
    }
    return
}

// RiderBelong holds the schema definition for the RiderBelong entity.
type RiderBelong struct {
    ent.Schema
}

// Annotations of the RiderBelong.
func (RiderBelong) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "rider_belong"},
        entsql.WithComments(true),
    }
}

// Fields of the RiderBelong.
func (RiderBelong) Fields() []ent.Field {
    return []ent.Field{}
}

// Edges of the RiderBelong.
func (RiderBelong) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (RiderBelong) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (RiderBelong) Indexes() []ent.Index {
    return []ent.Index{}
}
