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

type PromotionMemberMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionMemberMixin) Fields() []ent.Field {
	relate := field.Uint64("member_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionMemberMixin) Edges() []ent.Edge {
	e := edge.To("member", PromotionMember.Type).Unique().Field("member_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionMemberMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("member_id"))
	}
	return
}

// PromotionMember holds the schema definition for the PromotionMember entity.
type PromotionMember struct {
	ent.Schema
}

// Annotations of the PromotionMember.
func (PromotionMember) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_member"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionMember.
func (PromotionMember) Fields() []ent.Field {
	return []ent.Field{
		field.String("phone").Comment("会员手机号").Unique(),
		field.String("name").Optional().Comment("会员姓名"),
		field.Float("balance").Default(0).Comment("钱包余额"),
		field.Float("frozen").Default(0).Comment("钱包冻结金额"),
		field.Uint64("total_growth_value").Default(0).Comment("总成长值"),
		field.Uint64("current_growth_value").Default(0).Comment("当前等级成长值"),
		field.Bool("enable").Default(true).Comment("是否启用"),
		field.Uint64("person_id").Optional().Nillable().Comment("实名认证ID"),
		field.String("avatar_url").Optional().Comment("头像url"),
	}
}

// Edges of the PromotionMember.
func (PromotionMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("referring", PromotionReferrals.Type),
		edge.To("referred", PromotionReferrals.Type).Unique(),

		edge.From("person", PromotionPerson.Type).Ref("member").Unique().Field("person_id"),
		edge.To("cards", PromotionBankCard.Type),
	}
}

func (PromotionMember) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		RiderMixin{Optional: true},
		PromotionLevelMixin{Optional: true},
		PromotionCommissionMixin{Optional: true},
	}
}

func (PromotionMember) Indexes() []ent.Index {
	return []ent.Index{}
}
