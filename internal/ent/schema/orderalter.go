package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// OrderAlter holds the schema definition for the OrderAlter entity.
type OrderAlter struct {
    ent.Schema
}

// Annotations of the OrderAlter.
func (OrderAlter) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "order_alter"},
    }
}

// Fields of the OrderAlter.
func (OrderAlter) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Comment("骑手ID"),
        field.Uint64("order_id").Comment("订单ID"),
        field.Int("days").Comment("更改天数"),
        field.String("reason").Comment("更改原因"),
    }
}

// Edges of the OrderAlter.
func (OrderAlter) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("alters").Required().Unique().Field("rider_id").Comment("骑手"),
        edge.From("order", Order.Type).Ref("alters").Required().Unique().Field("order_id").Comment("订单"),
    }
}

func (OrderAlter) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (OrderAlter) Indexes() []ent.Index {
    return []ent.Index{}
}
