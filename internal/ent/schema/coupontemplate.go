package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type CouponTemplateMixin struct {
    mixin.Schema
    Optional bool
}

func (m CouponTemplateMixin) Fields() []ent.Field {
    relate := field.Uint64("template_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m CouponTemplateMixin) Edges() []ent.Edge {
    e := edge.To("template", CouponTemplate.Type).Unique().Field("template_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

// CouponTemplate holds the schema definition for the CouponTemplate entity.
type CouponTemplate struct {
    ent.Schema
}

// Annotations of the CouponTemplate.
func (CouponTemplate) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "coupon_template"},
    }
}

// Fields of the CouponTemplate.
func (CouponTemplate) Fields() []ent.Field {
    return []ent.Field{
        field.Bool("enable").Default(true).Comment("是否启用"),
        field.String("name").Comment("名称"),
        field.JSON("meta", &model.CouponTemplateMeta{}).Comment("详情"),
    }
}

// Edges of the CouponTemplate.
func (CouponTemplate) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (CouponTemplate) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
    }
}

func (CouponTemplate) Indexes() []ent.Index {
    return []ent.Index{}
}
