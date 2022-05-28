package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// OrderArrearage holds the schema definition for the OrderArrearage entity.
type OrderArrearage struct {
    ent.Schema
}

// Annotations of the OrderArrearage.
func (OrderArrearage) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "order_arrearage"},
    }
}

// Fields of the OrderArrearage.
func (OrderArrearage) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Comment("骑手ID"),
        field.Uint64("order_id").Comment("订单ID"),
    }
}

// Edges of the OrderArrearage.
func (OrderArrearage) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("arrearages").Required().Unique().Field("rider_id").Comment("骑手"),
        edge.From("order", Order.Type).Ref("arrearages").Required().Unique().Field("order_id").Comment("订单"),
    }
}

func (OrderArrearage) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (OrderArrearage) Indexes() []ent.Index {
    return []ent.Index{}
}
