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
        field.Uint64("city_id").Comment("城市ID"),
        field.Uint8("status").Default(1).Comment("订单状态 0未支付 1已支付 2申请退款 3已退款"),
        field.Uint8("payway").Immutable().Comment("支付方式 1支付宝 2微信"),
        field.Uint("type").Immutable().Comment("订单类型 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金"),
        field.String("out_trade_no").Unique().Immutable().Comment("交易订单号"),
        field.String("trade_no").Immutable().Comment("平台订单号"),
        field.Float("amount").Immutable().Comment("支付金额"),
        field.JSON("plan_detail", model.PlanItem{}).Optional().Comment("骑士卡详情"),
        field.JSON("refund", model.OrderRefund{}).Optional().Comment("退款详细"),
        field.Uint64("parent_id").Optional().Comment("续签所属订单ID"),
        field.Time("start_at").Optional().Nillable().Comment("开始时间"),
        field.Time("end_at").Optional().Nillable().Comment("结束时间"),
        field.Time("paused_at").Optional().Comment("当前是否暂停计费, 暂停计费时间"),
        field.Uint("days").Optional().Comment("骑士卡天数"),
    }
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("orders").Required().Unique().Field("rider_id").Comment("骑手ID"),
        edge.From("plan", Plan.Type).Ref("orders").Unique().Field("plan_id").Comment("骑士卡"),
        edge.To("commission", Commission.Type).Unique(),
        edge.To("children", Order.Type).From("parent").Field("parent_id").Unique().Comment("续签所属订单"),
        edge.From("city", City.Type).Ref("orders").Required().Unique().Field("city_id").Comment("城市"),
        edge.To("pauses", OrderPause.Type),
        edge.To("arrearages", OrderArrearage.Type),
        edge.To("alters", OrderAlter.Type),
    }
}

func (Order) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.SonyflakeIDMixin{},
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Order) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("trade_no", "out_trade_no"),
        index.Fields("status"),
        index.Fields("paused_at"),
        index.Fields("start_at", "end_at"),
    }
}
