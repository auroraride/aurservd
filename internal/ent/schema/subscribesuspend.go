package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// SubscribeSuspend holds the schema definition for the SubscribeSuspend entity.
type SubscribeSuspend struct {
    ent.Schema
}

// Annotations of the SubscribeSuspend.
func (SubscribeSuspend) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "subscribe_suspend"},
    }
}

// Fields of the SubscribeSuspend.
func (SubscribeSuspend) Fields() []ent.Field {
    return []ent.Field{
        field.Int("days").Default(0).Comment("暂停天数"),
        field.Time("start_at").Comment("开始时间"),
        field.Time("stop_at").Optional().Nillable().Comment("结束时间"),
    }
}

// Edges of the SubscribeSuspend.
func (SubscribeSuspend) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (SubscribeSuspend) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.Modifier{},
        SubscribeMixin{},
        CityMixin{},
        RiderMixin{},
    }
}

func (SubscribeSuspend) Indexes() []ent.Index {
    return []ent.Index{}
}
