package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
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
        field.String("name").Immutable().Comment("名称"),
        field.Float("amount").Comment("金额"),
        field.String("code").Unique().Comment("券码"),
        field.Time("expired_at").Comment("过期时间"),
        field.Time("used_at").Optional().Comment("使用时间"),
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
        internal.Modifier{},
        RiderMixin{Optional: true},
        CouponAssemblyMixin{},
        CouponTemplateMixin{},
        OrderMixin{Optional: true},
        PlanMixin{Optional: true},
    }
}

func (Coupon) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("name").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
