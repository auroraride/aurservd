package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

// BatteryModelNew holds the schema definition for the BatteryModelNew entity.
type BatteryModelNew struct {
	ent.Schema
}

// Annotations of the BatteryModelNew.
func (BatteryModelNew) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "battery_model_new"},
		entsql.WithComments(true),
	}
}

// Fields of the BatteryModelNew.
func (BatteryModelNew) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("type").Comment("电池类型 1智能电池 2非智能电池"),
		field.Uint("voltage").Comment("电压"),
		field.Uint("capacity").Comment("容量"),
		field.String("model").Unique().Comment("电池型号"),
	}
}

// Edges of the BatteryModelNew.
func (BatteryModelNew) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (BatteryModelNew) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (BatteryModelNew) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("model"),
	}
}
