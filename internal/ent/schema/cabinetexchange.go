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

// CabinetExchange holds the schema definition for the CabinetExchange entity.
type CabinetExchange struct {
    ent.Schema
}

// Annotations of the CabinetExchange.
func (CabinetExchange) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "cabinet_exchange"},
    }
}

// Fields of the CabinetExchange.
func (CabinetExchange) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Comment("骑手ID"),
        field.Uint64("cabinet_id").Comment("电柜ID"),
        field.Bool("alternative").Default(false).Comment("是否备无满电选方案"),
        field.Uint("step").Comment("步骤"),
        field.Uint("status").Comment("状态"),
        field.Uint("bin_index").Comment("仓位Index"),
        field.JSON("bin", model.CabinetBin{}).Comment("仓位详情"),
    }
}

// Edges of the CabinetExchange.
func (CabinetExchange) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("exchanges").Required().Unique().Field("rider_id").Comment("骑手ID"),
        edge.From("cabinet", Cabinet.Type).Ref("exchanges").Required().Unique().Field("cabinet_id").Comment("电柜ID")}
}

func (CabinetExchange) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (CabinetExchange) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("bin_index"),
        index.Fields("status"),
        index.Fields("step"),
    }
}
