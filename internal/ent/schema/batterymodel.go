package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// BatteryModel holds the schema definition for the BatteryModel entity.
type BatteryModel struct {
    ent.Schema
}

// Annotations of the BatteryModel.
func (BatteryModel) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "battery_model"},
    }
}

// Fields of the BatteryModel.
func (BatteryModel) Fields() []ent.Field {
    return []ent.Field{
        field.String("model").Unique().Comment("型号"),
        field.Bool("enable").Default(true).Comment("是否启用"),
    }
}

// Edges of the BatteryModel.
func (BatteryModel) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("cabinets", Cabinet.Type).Ref("bms"),
        edge.From("plans", Plan.Type).Ref("pms"),
    }
}

func (BatteryModel) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (BatteryModel) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("model").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
