package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type BatteryModelMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m BatteryModelMixin) Fields() []ent.Field {
	relate := field.Uint64("model_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m BatteryModelMixin) Edges() []ent.Edge {
	e := edge.To("model", BatteryModel.Type).Unique().Field("model_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m BatteryModelMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("model_id"))
	}
	return
}

// BatteryModel holds the schema definition for the BatteryModel entity.
type BatteryModel struct {
	ent.Schema
}

// Annotations of the BatteryModel.
func (BatteryModel) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "battery_model"},
		entsql.WithComments(true),
	}
}

// Fields of the BatteryModel.
func (BatteryModel) Fields() []ent.Field {
	return []ent.Field{
		field.String("model").Unique().Comment("型号"),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Uint8("type").Default(1).Comment("电池类型 1智能电池 2非智能电池"),
		field.Uint("voltage").Optional().Comment("电压"),
		field.Uint("capacity").Optional().Comment("容量"),
	}
}

// Edges of the BatteryModel.
func (BatteryModel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cabinets", Cabinet.Type).Ref("models"),
	}
}

func (BatteryModel) Mixin() []ent.Mixin {
	return []ent.Mixin{}
}

func (BatteryModel) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("model").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}
