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

type PromotionWithdrawalMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionWithdrawalMixin) Fields() []ent.Field {
	relate := field.Uint64("withdrawal_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionWithdrawalMixin) Edges() []ent.Edge {
	e := edge.To("withdrawal", PromotionWithdrawal.Type).Unique().Field("withdrawal_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionWithdrawalMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("withdrawal_id"))
	}
	return
}

// PromotionWithdrawal holds the schema definition for the PromotionWithdrawal entity.
type PromotionWithdrawal struct {
	ent.Schema
}

// Annotations of the PromotionWithdrawal.
func (PromotionWithdrawal) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_withdrawal"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionWithdrawal.
func (PromotionWithdrawal) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("status").Default(0).Comment("提现状态 0:待审核 1:成功 2:失败"),
		field.Float("apply_amount").Default(0).Comment("提现申请金额"),
		field.Float("amount").Default(0).Comment("提现金额"),
		field.Float("fee").Default(0).Comment("提现手续费"),
		field.Uint8("method").Comment("提现方式 1:银行卡"),
		field.Uint64("account_id").Optional().Comment("提现账号ID"),
		field.Time("apply_time").Optional().Comment("申请时间"),
		field.Time("review_time").Optional().Comment("审核时间"),
	}
}

// Edges of the PromotionWithdrawal.
func (PromotionWithdrawal) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cards", PromotionBankCard.Type).Ref("withdrawals").Unique().Field("account_id"),
	}
}

func (PromotionWithdrawal) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		PromotionMemberMixin{},
	}
}

func (PromotionWithdrawal) Indexes() []ent.Index {
	return []ent.Index{}
}
