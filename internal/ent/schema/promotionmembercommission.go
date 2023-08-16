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

type PromotionMemberCommissionMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionMemberCommissionMixin) Fields() []ent.Field {
	relate := field.Uint64("commission_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionMemberCommissionMixin) Edges() []ent.Edge {
	e := edge.To("commission", PromotionMemberCommission.Type).Unique().Field("commission_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionMemberCommissionMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("commission_id"))
	}
	return
}

// PromotionMemberCommission holds the schema definition for the PromotionMemberCommission entity.
type PromotionMemberCommission struct {
	ent.Schema
}

// Annotations of the PromotionMemberCommission.
func (PromotionMemberCommission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_member_commission"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionMemberCommission.
func (PromotionMemberCommission) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("member_id").Optional().Comment("会员ID"),
	}
}

// Edges of the PromotionMemberCommission.
func (PromotionMemberCommission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("member", PromotionMember.Type).Ref("commissions").Unique().Field("member_id"),
	}
}

func (PromotionMemberCommission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		PromotionCommissionMixin{},
	}
}

func (PromotionMemberCommission) Indexes() []ent.Index {
	return []ent.Index{}
}
