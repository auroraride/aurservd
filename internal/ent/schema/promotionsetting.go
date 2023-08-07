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

type PromotionSettingMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionSettingMixin) Fields() []ent.Field {
	relate := field.Uint64("setting_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionSettingMixin) Edges() []ent.Edge {
	e := edge.To("setting", PromotionSetting.Type).Unique().Field("setting_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionSettingMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("setting_id"))
	}
	return
}

// PromotionSetting holds the schema definition for the PromotionSetting entity.
type PromotionSetting struct {
	ent.Schema
}

// Annotations of the PromotionSetting.
func (PromotionSetting) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_setting"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionSetting.
func (PromotionSetting) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").Optional().Comment("标题"),
		field.Text("content").Optional().Comment("内容"),
		field.String("key").Comment("设置项"),
	}
}

// Edges of the PromotionSetting.
func (PromotionSetting) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PromotionSetting) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.Modifier{},
	}
}

func (PromotionSetting) Indexes() []ent.Index {
	return []ent.Index{}
}
