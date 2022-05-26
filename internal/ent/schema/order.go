package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/app/model"
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
        field.Uint64("plan_id").Optional().Comment("骑士卡ID"),
        field.Uint8("status").Default(1).Comment("订单状态 0未支付 1已支付 2申请退款 3已退款"),
        field.Uint8("payway").Immutable().Comment("支付方式 1支付宝 2微信"),
        field.Uint("type").Immutable().Comment("订单类型 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金"),
        field.String("out_trade_no").Unique().Immutable().Comment("交易订单号"),
        field.String("trade_no").Immutable().Comment("平台订单号"),
        field.Float("amount").Immutable().Comment("支付金额"),
        field.JSON("plan_detail", model.PlanItem{}).Optional().Comment("骑士卡详情"),
        field.Uint64("parent_id").Optional().Comment("续签/更改电池接续从属订单ID(上个订单)"),
        field.JSON("subordinate", model.OrderSubordinate{}).Optional().Comment("接续订单属性"),
    }
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("orders").Required().Unique().Field("rider_id"),
        edge.From("plan", Plan.Type).Ref("orders").Unique().Field("plan_id"),
        edge.To("commission", Commission.Type).Unique(),
    }
}

func (Order) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Order) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("trade_no"),
        index.Fields("status"),
    }
}
