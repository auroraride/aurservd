package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/internal"
)

type PointLogMixin struct {
	mixin.Schema
	DisableIndex bool
	Optional     bool
}

func (m PointLogMixin) Fields() []ent.Field {
	relate := field.Uint64("point_log_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PointLogMixin) Edges() []ent.Edge {
	e := edge.To("Log", PointLog.Type).Unique().Field("log_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PointLogMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("log_id"))
	}
	return
}

// PointLog holds the schema definition for the PointLog entity.
type PointLog struct {
	ent.Schema
}

// Annotations of the PointLog.
func (PointLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "point_log"},
		entsql.WithComments(true),
	}
}

// Fields of the PointLog.
func (PointLog) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("type").Comment("变动类型"),
		field.Int64("points").Comment("变动数量"),
		field.Int64("after").Comment("变动结果"),
		field.String("reason").Optional().Nillable().Comment("原因"),
		field.JSON("attach", &model.PointLogAttach{}).Optional().Comment("其他信息"),
	}
}

// Edges of the PointLog.
func (PointLog) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PointLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		RiderMixin{},
		OrderMixin{Optional: true},
		internal.HookModifier{},
		internal.HookEmployee{},
	}
}

func (PointLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
	}
}
