package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// SubscribePause holds the schema definition for the SubscribePause entity.
type SubscribePause struct {
    ent.Schema
}

// Annotations of the SubscribePause.
func (SubscribePause) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "subscribe_pause"},
    }
}

// Fields of the SubscribePause.
func (SubscribePause) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("subscribe_id").Comment("订阅ID"),
        field.Time("start_at").Comment("寄存开始时间"),
        field.Time("end_at").Optional().Comment("寄存结束时间"),
        field.Int("days").Optional().Comment("寄存天数 = 天数差 - 重复天数"),
        field.Uint64("end_employee_id").Optional().Nillable().Comment("结束寄存店员ID"),
        field.Int("overdue_days").Default(0).Comment("超期天数"),
        field.JSON("end_modifier", &model.Modifier{}).Optional().Comment("结束寄存管理员信息"),
        field.Bool("pause_overdue").Default(false).Comment("是否超期退租"),
        field.Int("suspend_days").Default(0).Comment("重复天数, 寄存过程中暂停扣费天数"),
    }
}

// Edges of the SubscribePause.
func (SubscribePause) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("subscribe", Subscribe.Type).Ref("pauses").Required().Unique().Field("subscribe_id").Comment("订阅"),
        edge.To("end_employee", Employee.Type).Unique().Field("end_employee_id"),
        edge.To("suspends", SubscribeSuspend.Type),
    }
}

func (SubscribePause) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{IndexCreator: true},
        RiderMixin{},
        EmployeeMixin{Optional: true},
        CityMixin{Optional: true},
        StoreMixin{Optional: true},
        StoreMixin{Optional: true, Prefix: "end"},
        CabinetMixin{Optional: true},
        CabinetMixin{Optional: true, Prefix: "end"},
    }
}

func (SubscribePause) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("subscribe_id"),
        index.Fields("end_modifier").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
