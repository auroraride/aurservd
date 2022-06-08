package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// SubscribePause holds the schema definition for the SubscribePause entity.
type SubscribePause struct {
    ent.Schema
}

// Annotations of the SubscribePause.
func (SubscribePause) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "subscribe_pause"},
    }
}

// Fields of the SubscribePause.
func (SubscribePause) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("subscribe_id").Comment("订阅ID"),
        field.Time("start_at").Comment("暂停开始时间"),
        field.Time("end_at").Optional().Comment("暂停结束时间"),
        field.Int("days").Optional().Comment("暂停天数"),
    }
}

// Edges of the SubscribePause.
func (SubscribePause) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("subscribe", Subscribe.Type).Ref("pauses").Required().Unique().Field("subscribe_id").Comment("订阅"),
    }
}

func (SubscribePause) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
        RiderMixin{},
        EmployeeMixin{Optional: true},
    }
}

func (SubscribePause) Indexes() []ent.Index {
    return []ent.Index{}
}
