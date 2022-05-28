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

// OrderPause holds the schema definition for the OrderPause entity.
type OrderPause struct {
    ent.Schema
}

// Annotations of the OrderPause.
func (OrderPause) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "order_pause"},
    }
}

// Fields of the OrderPause.
func (OrderPause) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Comment("骑手ID"),
        field.Uint64("order_id").Comment("订单ID"),
        field.Time("start_at").Comment("暂停开始时间"),
        field.Time("end_at").Optional().Comment("暂停结束时间"),
        field.Int("days").Optional().Comment("暂停天数"),
    }
}

// Edges of the OrderPause.
func (OrderPause) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("pauses").Required().Unique().Field("rider_id").Comment("骑手"),
        edge.From("order", Order.Type).Ref("pauses").Required().Unique().Field("order_id").Comment("订单"),
    }
}

func (OrderPause) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (OrderPause) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("start_at"),
        index.Fields("end_at"),
    }
}
