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

type PromotionLevelMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionLevelMixin) Fields() []ent.Field {
	relate := field.Uint64("level_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionLevelMixin) Edges() []ent.Edge {
	e := edge.To("level", PromotionLevel.Type).Unique().Field("level_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionLevelMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("level_id"))
	}
	return
}

// PromotionLevel holds the schema definition for the PromotionLevel entity.
type PromotionLevel struct {
	ent.Schema
}

// Annotations of the PromotionLevel.
func (PromotionLevel) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_level"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionLevel.
func (PromotionLevel) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("level").Comment("会员等级"),
		field.Uint64("growth_value").Default(0).Comment("所需成长值"),
		field.Float("commission_ratio").Default(0).Comment("会员权益佣金比例"),
	}
}

// Edges of the PromotionLevel.
func (PromotionLevel) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PromotionLevel) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (PromotionLevel) Indexes() []ent.Index {
	return []ent.Index{}
}
