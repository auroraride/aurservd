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

type CouponMixin struct {
    mixin.Schema
    Optional bool
}

func (m CouponMixin) Fields() []ent.Field {
    f := field.Uint64("coupon_id")
    if m.Optional {
        f.Optional().Nillable()
    }
    return []ent.Field{
        f,
    }
}

func (m CouponMixin) Edges() []ent.Edge {
    e := edge.To("coupon", Coupon.Type).Unique().Field("coupon_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

// Coupon holds the schema definition for the Coupon entity.
type Coupon struct {
    ent.Schema
}

// Annotations of the Coupon.
func (Coupon) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "coupon"},
    }
}

// Fields of the Coupon.
func (Coupon) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Comment("名称"),
        field.Int("total").Comment("总数"),
        field.Uint8("expired_type").Comment("过期类型"),
        field.Uint8("rule").Comment("优惠券规则, 1:互斥 2:叠加"),
    }
}

// Edges of the Coupon.
func (Coupon) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("cities", City.Type),
        edge.To("plans", Plan.Type),
    }
}

func (Coupon) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Coupon) Indexes() []ent.Index {
    return []ent.Index{}
}
