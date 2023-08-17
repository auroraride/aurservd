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

type PromotionCommissionPlanMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionCommissionPlanMixin) Fields() []ent.Field {
	relate := field.Uint64("plan_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionCommissionPlanMixin) Edges() []ent.Edge {
	e := edge.To("plan", PromotionCommissionPlan.Type).Unique().Field("plan_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionCommissionPlanMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("plan_id"))
	}
	return
}

// PromotionCommissionPlan holds the schema definition for the PromotionCommissionPlan entity.
type PromotionCommissionPlan struct {
	ent.Schema
}

// Annotations of the PromotionCommissionPlan.
func (PromotionCommissionPlan) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_commission_plan"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionCommissionPlan.
func (PromotionCommissionPlan) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("commission_id").Optional(),
		field.Uint64("plan_id").Optional(),
	}
}

// Edges of the PromotionCommissionPlan.
func (PromotionCommissionPlan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("promotion_commission", PromotionCommission.Type).Ref("plans").Unique().Field("commission_id"),
		edge.From("plan", Plan.Type).Ref("commissions").Unique().Field("plan_id"),
	}
}

func (PromotionCommissionPlan) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		PromotionMemberMixin{Optional: true},
	}
}

func (PromotionCommissionPlan) Indexes() []ent.Index {
	return []ent.Index{}
}
