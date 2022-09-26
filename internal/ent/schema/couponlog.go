package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// CouponLog holds the schema definition for the CouponLog entity.
type CouponLog struct {
    ent.Schema
}

// Annotations of the CouponLog.
func (CouponLog) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "coupon_log"},
    }
}

// Fields of the CouponLog.
func (CouponLog) Fields() []ent.Field {
    return []ent.Field{}
}

// Edges of the CouponLog.
func (CouponLog) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (CouponLog) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (CouponLog) Indexes() []ent.Index {
    return []ent.Index{}
}
