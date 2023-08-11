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

type PromotionBankCardMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionBankCardMixin) Fields() []ent.Field {
	relate := field.Uint64("card_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionBankCardMixin) Edges() []ent.Edge {
	e := edge.To("card", PromotionBankCard.Type).Unique().Field("card_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionBankCardMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("card_id"))
	}
	return
}

// PromotionBankCard holds the schema definition for the PromotionBankCard entity.
type PromotionBankCard struct {
	ent.Schema
}

// Annotations of the PromotionBankCard.
func (PromotionBankCard) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_bank_card"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionBankCard.
func (PromotionBankCard) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("member_id").Optional().Comment("会员ID"),
		field.String("card_no").Comment("银行卡号"),
		field.String("bank").Optional().Comment("银行名称"),
		field.Bool("is_default").Default(false).Comment("是否是默认银行卡"),
		field.String("bank_logo_url").Optional().Comment("银行卡logo"),
		field.String("province").Optional().Comment("省份"),
		field.String("city").Optional().Comment("城市"),
	}
}

// Edges of the PromotionBankCard.
func (PromotionBankCard) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("member", PromotionMember.Type).Ref("cards").Unique().Field("member_id"),
		edge.To("withdrawals", PromotionWithdrawal.Type),
	}
}

func (PromotionBankCard) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (PromotionBankCard) Indexes() []ent.Index {
	return []ent.Index{}
}
