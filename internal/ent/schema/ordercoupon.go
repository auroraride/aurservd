package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type OrderCouponMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m OrderCouponMixin) Fields() []ent.Field {
    relate := field.Uint64("coupon_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m OrderCouponMixin) Edges() []ent.Edge {
    e := edge.To("coupon", OrderCoupon.Type).Unique().Field("coupon_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m OrderCouponMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("coupon_id"))
    }
    return
}

// OrderCoupon holds the schema definition for the OrderCoupon entity.
type OrderCoupon struct {
    ent.Schema
}

// Annotations of the OrderCoupon.
func (OrderCoupon) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "order_coupon"},
    }
}

// Fields of the OrderCoupon.
func (OrderCoupon) Fields() []ent.Field {
    return []ent.Field{}
}

// Edges of the OrderCoupon.
func (OrderCoupon) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (OrderCoupon) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (OrderCoupon) Indexes() []ent.Index {
    return []ent.Index{}
}
