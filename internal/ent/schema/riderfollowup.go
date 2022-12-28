package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// RiderFollowUp holds the schema definition for the RiderFollowUp entity.
type RiderFollowUp struct {
    ent.Schema
}

// Annotations of the RiderFollowUp.
func (RiderFollowUp) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "rider_follow_up"},
        entsql.WithComments(true),
    }
}

// Fields of the RiderFollowUp.
func (RiderFollowUp) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Comment("骑手ID"),
    }
}

// Edges of the RiderFollowUp.
func (RiderFollowUp) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Unique().Required().Ref("followups").Field("rider_id"),
    }
}

func (RiderFollowUp) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
        ManagerMixin{},
    }
}

func (RiderFollowUp) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("rider_id"),
    }
}
