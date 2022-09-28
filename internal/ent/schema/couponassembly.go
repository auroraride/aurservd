package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type CouponAssemblyMixin struct {
    mixin.Schema
    Optional bool
}

func (m CouponAssemblyMixin) Fields() []ent.Field {
    relate := field.Uint64("assembly_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m CouponAssemblyMixin) Edges() []ent.Edge {
    e := edge.To("assembly", CouponAssembly.Type).Unique().Field("assembly_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

// CouponAssembly holds the schema definition for the CouponAssembly entity.
type CouponAssembly struct {
    ent.Schema
}

// Annotations of the CouponAssembly.
func (CouponAssembly) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "coupon_assembly"},
    }
}

// Fields of the CouponAssembly.
func (CouponAssembly) Fields() []ent.Field {
    return []ent.Field{}
}

// Edges of the CouponAssembly.
func (CouponAssembly) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (CouponAssembly) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.Modifier{},
        CouponTemplateMixin{},
    }
}

func (CouponAssembly) Indexes() []ent.Index {
    return []ent.Index{}
}
