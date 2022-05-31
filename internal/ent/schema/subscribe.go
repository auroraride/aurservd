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

// Subscribe holds the schema definition for the Subscribe entity.
type Subscribe struct {
    ent.Schema
}

// Annotations of the Subscribe.
func (Subscribe) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{
            Table: "subscribe",
        },
    }
}

// Fields of the Subscribe.
func (Subscribe) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Comment("骑手ID"),
        field.Uint64("initial_order_id").Comment("初始订单ID(开通订阅的初始订单)"),
        field.Uint("type").Immutable().Comment("订阅类型 1新签 2续签 3重签 4更改电池"),
        field.Float("voltage").Comment("可用电压型号"),
        field.Uint("days").Comment("骑士卡天数"),
        field.Uint("alter_days").Default(0).Comment("改动天数"),
        field.Uint("pause_days").Default(0).Comment("暂停天数"),
        field.Time("paused_at").Optional().Nillable().Comment("当前是否暂停计费, 暂停计费时间"),
        field.Time("start_at").Optional().Nillable().Comment("激活时间"),
        field.Time("end_at").Optional().Nillable().Comment("归还时间"),
        field.Time("refund_at").Optional().Nillable().Comment("退款时间"),
    }
}

// Edges of the Subscribe.
func (Subscribe) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("subscribes").Required().Unique().Field("rider_id").Comment("骑手"),

        edge.To("pauses", SubscribePause.Type),
        edge.To("alters", SubscribeAlter.Type),
        edge.To("orders", Order.Type),

        edge.To("initial_order", Order.Type).Unique().Required().Field("initial_order_id").Comment("对应初始订单"),
    }
}

func (Subscribe) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        PlanMixin{},
        EmployeeMixin{},
        CityMixin{},
    }
}

func (Subscribe) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("paused_at"),
        index.Fields("start_at", "end_at"),
    }
}
