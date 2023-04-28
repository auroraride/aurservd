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
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/internal"
)

type CouponMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
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

func (m CouponMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("coupon_id"))
	}
	return
}

// Coupon holds the schema definition for the Coupon entity.
type Coupon struct {
	ent.Schema
}

// Annotations of the Coupon.
func (Coupon) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "coupon"},
		entsql.WithComments(true),
	}
}

// Fields of the Coupon.
func (Coupon) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("template_id"),
		field.Uint64("order_id").Optional().Nillable().Comment("订单ID"),
		field.String("name").Immutable().Comment("名称"),
		field.Uint8("rule").Comment("使用规则"),
		field.Bool("multiple").Default(false).Comment("该券是否可叠加"),
		field.Float("amount").Comment("金额"),
		field.String("code").Unique().Comment("券码"),
		field.Time("expires_at").Optional().Comment("过期时间"),
		field.Time("used_at").Optional().Nillable().Comment("使用时间"),
		field.JSON("duration", &model.CouponDuration{}).Comment("有效期规则"),
		field.JSON("plans", []*model.Plan{}).Optional().Comment("可用骑士卡"),
		field.JSON("cities", []model.City{}).Optional().Comment("可用城市"),
	}
}

// Edges of the Coupon.
func (Coupon) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("template", CouponTemplate.Type).Required().Unique().Ref("coupons").Field("template_id"),
		edge.From("order", Order.Type).Ref("coupons").Unique().Field("order_id"),
	}
}

func (Coupon) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.Modifier{},

		RiderMixin{Optional: true},
		CouponAssemblyMixin{},
		PlanMixin{Optional: true, Comment: "实际使用骑士卡"},
	}
}

func (Coupon) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
		index.Fields("order_id"),
		index.Fields("template_id"),
		index.Fields("multiple"),
		index.Fields("expires_at"),
		index.Fields("used_at"),
	}
}
