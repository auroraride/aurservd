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

type PromotionReferralsMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionReferralsMixin) Fields() []ent.Field {
	relate := field.Uint64("referrals_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionReferralsMixin) Edges() []ent.Edge {
	e := edge.To("referrals", PromotionReferrals.Type).Unique().Field("referrals_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionReferralsMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("referrals_id"))
	}
	return
}

// PromotionReferrals holds the schema definition for the PromotionReferrals entity.
type PromotionReferrals struct {
	ent.Schema
}

// Annotations of the PromotionReferrals.
func (PromotionReferrals) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_referrals"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionReferrals.
func (PromotionReferrals) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("referring_member_id").Optional().Nillable().Comment("推广者id"),
		field.Uint64("referred_member_id").Optional().Unique().Comment("被推广者ID"),
	}
}

// Edges of the PromotionReferrals.
func (PromotionReferrals) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("referring_member", PromotionMember.Type).Ref("referring").Unique().Field("referring_member_id"),
		edge.From("referred_member", PromotionMember.Type).Ref("referred").Unique().Field("referred_member_id"),
	}
}

func (PromotionReferrals) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},

		RiderMixin{Optional: true},
		SubscribeMixin{Optional: true},
	}
}

func (PromotionReferrals) Indexes() []ent.Index {
	return []ent.Index{}
}
