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

type PromotionGrowthMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionGrowthMixin) Fields() []ent.Field {
	relate := field.Uint64("growth_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionGrowthMixin) Edges() []ent.Edge {
	e := edge.To("growth", PromotionGrowth.Type).Unique().Field("growth_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionGrowthMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("growth_id"))
	}
	return
}

// PromotionGrowth holds the schema definition for the PromotionGrowth entity.
type PromotionGrowth struct {
	ent.Schema
}

// Annotations of the PromotionGrowth.
func (PromotionGrowth) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_growth"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionGrowth.
func (PromotionGrowth) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("status").Default(1).Comment("状态 1:有效 2:无效"),
		field.Uint64("growth_value").Comment("成长值"),
	}
}

// Edges of the PromotionGrowth.
func (PromotionGrowth) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PromotionGrowth) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		PromotionMemberMixin{Optional: true},
		PromotionLevelTaskMixin{Optional: true},
	}
}

func (PromotionGrowth) Indexes() []ent.Index {
	return []ent.Index{}
}
