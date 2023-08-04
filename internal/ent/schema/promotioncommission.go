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

type PromotionCommissionMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionCommissionMixin) Fields() []ent.Field {
	relate := field.Uint64("commission_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionCommissionMixin) Edges() []ent.Edge {
	e := edge.To("commission", PromotionCommission.Type).Unique().Field("commission_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionCommissionMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("commission_id"))
	}
	return
}

// PromotionCommission holds the schema definition for the PromotionCommission entity.
type PromotionCommission struct {
	ent.Schema
}

// Annotations of the PromotionCommission.
func (PromotionCommission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_commission"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionCommission.
func (PromotionCommission) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("type").Default(1).Nillable().Comment("类型 0:全局 1:通用 2:个人"),
		field.String("name").Comment("方案名"),
		field.JSON("rule", &promotion.CommissionRule{}).Comment("返佣方案规则"),
		field.Bool("enable").Default(true).Comment("启用状态 0:禁用 1:启用"),
		field.Float("amount_sum").Nillable().Default(0).Comment("累计返佣金额"),
		field.Text("desc").Optional().Nillable().Comment("返佣说明"),
		field.JSON("history_id", []uint64{}).Optional().Comment("历史记录id"),
		field.Time("start_at").Optional().Nillable().Comment("开始时间"),
		field.Time("end_at").Optional().Nillable().Comment("结束时间"),
	}
}

// Edges of the PromotionCommission.
func (PromotionCommission) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PromotionCommission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		PromotionMemberMixin{Optional: true},
	}
}

func (PromotionCommission) Indexes() []ent.Index {
	return []ent.Index{}
}
