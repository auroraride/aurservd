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
        field.Time("start_at").Comment("暂停开始时间"),
        field.Time("end_at").Optional().Comment("暂停结束时间"),
        field.Int("days").Optional().Comment("暂停天数"),
        field.Uint64("end_employee_id").Optional().Nillable().Comment("结束寄存店员ID"),
        field.Bool("overdue").Default(false).Comment("是否超期"),
        field.JSON("end_modifier", &model.Modifier{}).Optional().Comment("结束寄存管理员信息"),
    }
}

// Edges of the SubscribePause.
func (SubscribePause) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("subscribe", Subscribe.Type).Ref("pauses").Required().Unique().Field("subscribe_id").Comment("订阅"),
        edge.To("end_employee", Employee.Type).Unique().Field("end_employee_id"),
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
        index.Fields("end_modifier").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
