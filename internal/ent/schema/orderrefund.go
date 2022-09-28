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

// OrderRefund holds the schema definition for the OrderRefund entity.
type OrderRefund struct {
    ent.Schema
}

// Annotations of the OrderRefund.
func (OrderRefund) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "order_refund"},
    }
}

// Fields of the OrderRefund.
func (OrderRefund) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("order_id").Comment("订单ID"),
        field.Uint8("status").Comment("退款状态"),
        field.Float("amount").Comment("退款金额"),
        field.String("out_refund_no").Unique().Comment("退款订单编号"),
        field.String("reason").Comment("退款理由"),
        field.Time("refund_at").Optional().Nillable().Comment("退款成功时间"),
    }
}

// Edges of the OrderRefund.
func (OrderRefund) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("order", Order.Type).Ref("refund").Required().Unique().Field("order_id").Comment("订单"),
    }
}

func (OrderRefund) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (OrderRefund) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("order_id"),
        index.Fields("status"),
        index.Fields("out_refund_no"),
    }
}
