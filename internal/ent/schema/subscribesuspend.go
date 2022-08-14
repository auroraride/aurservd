package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/app/model"
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
        field.Uint64("subscribe_id").Comment("订阅ID"),
        field.Uint64("pause_id").Optional().Comment("寄存ID"),
        field.Int("days").Default(0).Comment("暂停天数"),
        field.Time("start_at").Comment("开始时间"),
        field.Time("end_at").Optional().Comment("结束时间"),
        field.String("end_reason").Optional().Comment("结束理由"),
        field.JSON("end_modifier", &model.Modifier{}).Optional().Comment("继续计费管理员信息"),
    }
}

// Edges of the SubscribeSuspend.
func (SubscribeSuspend) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("subscribe", Subscribe.Type).Ref("suspends").Required().Unique().Field("subscribe_id"),
        edge.From("pause", SubscribePause.Type).Ref("suspends").Unique().Field("pause_id"),
    }
}

func (SubscribeSuspend) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.Modifier{},
        CityMixin{},
        RiderMixin{},
    }
}

func (SubscribeSuspend) Indexes() []ent.Index {
    return []ent.Index{}
}
