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

// Order holds the schema definition for the Order entity.
type Order struct {
    ent.Schema
}

// Annotations of the Order.
func (Order) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "order"},
    }
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Comment("骑手ID"),
        field.Uint64("parent_id").Optional().Comment("父订单ID"),
        field.Uint64("subscribe_id").Optional().Comment("所属订阅ID"),
        field.Uint8("status").Default(1).Comment("订单状态 0未支付 1已支付 2申请退款 3已退款"),
        field.Uint8("payway").Immutable().Comment("支付方式 1支付宝 2微信"),
        field.Uint("type").Immutable().Comment("订单类型 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金"),
        field.String("out_trade_no").Immutable().Comment("交易订单号"),
        field.String("trade_no").Immutable().Comment("平台订单号"),
        field.Float("amount").Immutable().Comment("子订单金额(拆分项此条订单)"),
        field.Float("total").Immutable().Default(0).Comment("此次支付总金额(包含所有子订单的总支付)"),
        field.Time("refund_at").Optional().Nillable().Comment("退款时间"),
    }
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("orders").Required().Unique().Field("rider_id").Comment("骑手"),
        edge.From("subscribe", Subscribe.Type).Ref("orders").Unique().Field("subscribe_id").Comment("所属订阅"),
        edge.To("commission", Commission.Type).Unique(),
        edge.To("children", Order.Type).From("parent").Field("parent_id").Unique().Comment("子订单"),
        edge.To("refunds", OrderRefund.Type).Comment("退款"),
    }
}

func (Order) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        PlanMixin{Optional: true},
        CityMixin{Optional: true},
    }
}

func (Order) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("trade_no"),
        index.Fields("out_trade_no"),
        index.Fields("status"),
    }
}
