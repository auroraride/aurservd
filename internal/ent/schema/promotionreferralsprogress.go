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

type PromotionReferralsProgressMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionReferralsProgressMixin) Fields() []ent.Field {
	relate := field.Uint64("progress_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionReferralsProgressMixin) Edges() []ent.Edge {
	e := edge.To("progress", PromotionReferralsProgress.Type).Unique().Field("progress_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionReferralsProgressMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("progress_id"))
	}
	return
}

// PromotionReferralsProgress holds the schema definition for the PromotionReferralsProgress entity.
type PromotionReferralsProgress struct {
	ent.Schema
}

// Annotations of the PromotionReferralsProgress.
func (PromotionReferralsProgress) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_referrals_progress"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionReferralsProgress.
func (PromotionReferralsProgress) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("referring_member_id").Optional().Comment("推广者id"),
		field.Uint64("referred_member_id").Optional().Comment("被推广者ID<骑手>"),
		field.String("name").Optional().Comment("姓名"),
		field.Uint8("status").Nillable().Default(0).Comment("状态  0: 邀请中 1:邀请成功 2:邀请失败"),
	}
}

// Edges of the PromotionReferralsProgress.
func (PromotionReferralsProgress) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PromotionReferralsProgress) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.Modifier{},
		RiderMixin{Optional: true},
	}
}

func (PromotionReferralsProgress) Indexes() []ent.Index {
	return []ent.Index{}
}
