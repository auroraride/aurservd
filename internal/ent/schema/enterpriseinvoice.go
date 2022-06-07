package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// EnterpriseInvoice holds the schema definition for the EnterpriseInvoice entity.
type EnterpriseInvoice struct {
    ent.Schema
}

// Annotations of the EnterpriseInvoice.
func (EnterpriseInvoice) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "enterprise_invoice"},
    }
}

// Fields of the EnterpriseInvoice.
func (EnterpriseInvoice) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Comment("骑手ID"),
        field.Float("price").Comment("单价"),
        field.Uint64("statement_id").Optional().Comment("团签结账对账单ID"),
    }
}

// Edges of the EnterpriseInvoice.
func (EnterpriseInvoice) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Unique().Required().Ref("invoices").Field("rider_id"),
        edge.From("statement", EnterpriseStatement.Type).Unique().Ref("invoices").Field("statement_id"),
    }
}

func (EnterpriseInvoice) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        StationMixin{},
        EnterpriseMixin{},
    }
}

func (EnterpriseInvoice) Indexes() []ent.Index {
    return []ent.Index{}
}
