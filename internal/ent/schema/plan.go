package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type PlanMixin struct {
    mixin.Schema
    Optional bool
}

func (m PlanMixin) Fields() []ent.Field {
    f := field.Uint64("plan_id").Comment("骑士卡ID")
    if m.Optional {
        f.Optional().Nillable()
    }
    return []ent.Field{f}
}

func (m PlanMixin) Edges() []ent.Edge {
    e := edge.To("plan", Plan.Type).Unique().Field("plan_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

// Plan holds the schema definition for the Plan entity.
type Plan struct {
    ent.Schema
}

// Annotations of the Plan.
func (Plan) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "plan"},
    }
}

// Fields of the Plan.
func (Plan) Fields() []ent.Field {
    return []ent.Field{
        field.Bool("enable").Comment("是否启用"),
        field.String("name").Comment("骑士卡名称"),
        field.Other("start", model.Date{}).SchemaType(map[string]string{dialect.Postgres: "date"}).Comment("有效期开始日期"),
        field.Other("end", model.Date{}).SchemaType(map[string]string{dialect.Postgres: "date"}).Comment("有效期结束日期"),
        field.Float("price").Comment("骑士卡价格"),
        field.Uint("days").Comment("骑士卡天数"),
        field.Float("commission").Comment("提成"),
        field.Float("original").Optional().Comment("原价"),
        field.String("desc").Optional().Comment("优惠信息"),
        field.Uint64("parent_id").Optional().Nillable().Comment("父级"),
    }
}

// Edges of the Plan.
func (Plan) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("pms", BatteryModel.Type),
        edge.To("cities", City.Type),
        edge.To("complexes", Plan.Type).
            From("parent").
            Field("parent_id").
            Unique(),
    }
}

func (Plan) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Plan) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("days"),
        index.Fields("enable"),
        index.Fields("start", "end"),
        index.Fields("name").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
