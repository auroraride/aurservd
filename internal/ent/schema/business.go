package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Business holds the schema definition for the Business entity.
type Business struct {
    ent.Schema
}

// Annotations of the Business.
func (Business) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "business"},
    }
}

// Fields of the Business.
func (Business) Fields() []ent.Field {
    return []ent.Field{
        field.Enum("type").Values("active", "pause", "continue", "unsubscribe").Comment("业务类型"),
        field.JSON("bin_info", &ec.BinInfo{}).Optional().Comment("仓位信息"),
    }
}

// Edges of the Business.
func (Business) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (Business) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        RiderMixin{},
        CityMixin{},
        SubscribeMixin{},

        EmployeeMixin{Optional: true},
        StoreMixin{Optional: true},
        PlanMixin{Optional: true},
        EnterpriseMixin{Optional: true},
        StationMixin{Optional: true},
        CabinetMixin{Optional: true},
    }
}

func (Business) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("type"),
    }
}
