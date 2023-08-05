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

type PromotionLevelTaskMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionLevelTaskMixin) Fields() []ent.Field {
	relate := field.Uint64("task_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionLevelTaskMixin) Edges() []ent.Edge {
	e := edge.To("task", PromotionLevelTask.Type).Unique().Field("task_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionLevelTaskMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("task_id"))
	}
	return
}

// PromotionLevelTask holds the schema definition for the PromotionLevelTask entity.
type PromotionLevelTask struct {
	ent.Schema
}

// Annotations of the PromotionLevelTask.
func (PromotionLevelTask) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_level_task"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionLevelTask.
func (PromotionLevelTask) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("任务名称"),
		field.String("description").Comment("任务描述"),
		field.Uint8("type").Comment("完成条件 1: 签约 2:续费"),
		field.Uint64("growth_value").Default(0).Comment("任务成长值"),
		field.String("key").Optional().Comment("任务key"),
	}
}

// Edges of the PromotionLevelTask.
func (PromotionLevelTask) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PromotionLevelTask) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (PromotionLevelTask) Indexes() []ent.Index {
	return []ent.Index{}
}
