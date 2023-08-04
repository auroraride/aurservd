package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent/internal"
)

type PromotionEarningsMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionEarningsMixin) Fields() []ent.Field {
	relate := field.Uint64("earnings_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionEarningsMixin) Edges() []ent.Edge {
	e := edge.To("earnings", PromotionEarnings.Type).Unique().Field("earnings_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionEarningsMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("earnings_id"))
	}
	return
}

// PromotionEarnings holds the schema definition for the PromotionEarnings entity.
type PromotionEarnings struct {
	ent.Schema
}

// Annotations of the PromotionEarnings.
func (PromotionEarnings) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_earnings"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionEarnings.
func (PromotionEarnings) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("status").Default(promotion.EarningsStatusUnsettled.Value()).Comment("收益状态 0:未结算 1:已结算 2:已取消"),
		field.Float("amount").Default(0).Comment("收益金额"),
		field.String("commission_rule_key").Optional().Comment("返佣任务类型"),
	}
}

// Edges of the PromotionEarnings.
func (PromotionEarnings) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PromotionEarnings) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		PromotionCommissionMixin{},
		PromotionMemberMixin{},
		RiderMixin{},
	}
}

func (PromotionEarnings) Indexes() []ent.Index {
	return []ent.Index{}
}
