package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// EnterprisePrice holds the schema definition for the EnterprisePrice entity.
type EnterprisePrice struct {
    ent.Schema
}

// Annotations of the EnterprisePrice.
func (EnterprisePrice) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "enterprise_price"},
    }
}

// Fields of the EnterprisePrice.
func (EnterprisePrice) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("enterprise_id"),
        field.Float("price").Comment("单价 元/天"),
        field.Float("voltage").Comment("可用电池电压型号"),
    }
}

// Edges of the EnterprisePrice.
func (EnterprisePrice) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("enterprise", Enterprise.Type).
            Ref("prices").
            Unique().
            Required().
            Field("enterprise_id"),
    }
}

func (EnterprisePrice) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        CityMixin{Optional: false},
    }
}

func (EnterprisePrice) Indexes() []ent.Index {
    return []ent.Index{}
}