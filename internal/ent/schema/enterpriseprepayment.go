package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// EnterprisePrepayment holds the schema definition for the EnterprisePrepayment entity.
type EnterprisePrepayment struct {
    ent.Schema
}

// Annotations of the EnterprisePrepayment.
func (EnterprisePrepayment) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "enterprise_prepayment"},
    }
}

// Fields of the EnterprisePrepayment.
func (EnterprisePrepayment) Fields() []ent.Field {
    return []ent.Field{
        field.Float("amount").Immutable().Comment("预付金额"),
    }
}

// Edges of the EnterprisePrepayment.
func (EnterprisePrepayment) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (EnterprisePrepayment) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        EnterpriseMixin{Optional: false},
    }
}

func (EnterprisePrepayment) Indexes() []ent.Index {
    return []ent.Index{}
}