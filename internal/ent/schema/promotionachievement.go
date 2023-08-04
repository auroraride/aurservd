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

type PromotionAchievementMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionAchievementMixin) Fields() []ent.Field {
	relate := field.Uint64("achievement_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionAchievementMixin) Edges() []ent.Edge {
	e := edge.To("achievement", PromotionAchievement.Type).Unique().Field("achievement_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionAchievementMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("achievement_id"))
	}
	return
}

// PromotionAchievement holds the schema definition for the PromotionAchievement entity.
type PromotionAchievement struct {
	ent.Schema
}

// Annotations of the PromotionAchievement.
func (PromotionAchievement) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_achievement"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionAchievement.
func (PromotionAchievement) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("成就名称"),
		field.Uint8("type").Comment("成就类型 1:邀请成就 2:收益成就"),
		field.String("icon").Comment("成就图标"),
		field.Uint64("condition").Comment("完成条件"),
	}
}

// Edges of the PromotionAchievement.
func (PromotionAchievement) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PromotionAchievement) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (PromotionAchievement) Indexes() []ent.Index {
	return []ent.Index{}
}
