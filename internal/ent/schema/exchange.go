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

// Exchange holds the schema definition for the Exchange entity.
type Exchange struct {
    ent.Schema
}

// Annotations of the Exchange.
func (Exchange) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "exchange"},
    }
}

// Fields of the Exchange.
func (Exchange) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Comment("骑手ID"),
        field.String("uuid").Unique(),
        field.Uint64("cabinet_id").Optional().Comment("电柜ID"),
        field.Bool("success").Default(true).Comment("是否成功"),
        field.JSON("detail", &model.ExchangeCabinet{}).Optional().Comment("电柜换电信息"),
    }
}

// Edges of the Exchange.
func (Exchange) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("cabinet", Cabinet.Type).Unique().Ref("exchanges").Field("cabinet_id"),
        edge.From("rider", Rider.Type).Unique().Required().Ref("exchanges").Field("rider_id"),
    }
}

func (Exchange) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        CityMixin{},
        EmployeeMixin{Optional: true},
        StoreMixin{Optional: true},
    }
}

func (Exchange) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("success"),
    }
}
