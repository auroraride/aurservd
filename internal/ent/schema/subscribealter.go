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

// SubscribeAlter holds the schema definition for the SubscribeAlter entity.
type SubscribeAlter struct {
    ent.Schema
}

// Annotations of the SubscribeAlter.
func (SubscribeAlter) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "subscribe_alter"},
        entsql.WithComments(true),
    }
}

// Fields of the SubscribeAlter.
func (SubscribeAlter) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("subscribe_id").Comment("订阅ID"),
        field.Int("days").Comment("更改天数"),
    }
}

// Edges of the SubscribeAlter.
func (SubscribeAlter) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("subscribe", Subscribe.Type).Ref("alters").Required().Unique().Field("subscribe_id").Comment("订阅"),
    }
}

func (SubscribeAlter) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
        RiderMixin{},
        ManagerMixin{Optional: true},
        EnterpriseMixin{Optional: true},
        AgentMixin{Optional: true},
    }
}

func (SubscribeAlter) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("subscribe_id"),
    }
}
