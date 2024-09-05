package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/auroraride/adapter"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type BatteryFlowMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m BatteryFlowMixin) Fields() []ent.Field {
	relate := field.Uint64("flow_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m BatteryFlowMixin) Edges() []ent.Edge {
	e := edge.To("flow", BatteryFlow.Type).Unique().Field("flow_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m BatteryFlowMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("flow_id"))
	}
	return
}

// BatteryFlow holds the schema definition for the BatteryFlow entity.
type BatteryFlow struct {
	ent.Schema
}

// Annotations of the BatteryFlow.
func (BatteryFlow) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "battery_flow"},
		entsql.WithComments(true),
	}
}

// Fields of the BatteryFlow.
func (BatteryFlow) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("battery_id").Comment("电池ID"),
		field.String("sn").Comment("电池编号"),
		field.Float("soc").Default(-1).Comment("容量, -1代表未查询到"),
		field.Uint64("rider_id").Optional().Comment("骑手ID"),
		field.Uint64("cabinet_id").Optional().Comment("电柜ID"),
		field.String("serial").Optional().Comment("电柜编号"),
		field.Int("ordinal").Optional().Comment("仓位序号, 从1开始"),
		field.Other("geom", &adapter.Geometry{}).
			Optional().
			SchemaType(map[string]string{dialect.Postgres: "geometry(POINT, 4326)"}).
			Comment("坐标"),
		field.String("remark").Optional().Nillable().Comment("备注信息"),
	}
}

// Edges of the BatteryFlow.
func (BatteryFlow) Edges() []ent.Edge {
	return []ent.Edge{
		// edge.From("battery", Asset.Type).Field("battery_id").
		// 	Required().Unique().Ref("flows").
		// 	Annotations(
		// 		entsql.WithComments(true),
		// 	).
		// 	Comment("所属电池"),

		edge.From("cabinet", Cabinet.Type).Field("cabinet_id").
			Unique().Ref("battery_flows").
			Annotations(
				entsql.WithComments(true),
			).
			Comment("所属电柜"),

		edge.From("rider", Rider.Type).Field("rider_id").
			Unique().Ref("battery_flows").
			Annotations(
				entsql.WithComments(true),
			).
			Comment("所属骑手"),
	}
}

func (BatteryFlow) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		SubscribeMixin{Optional: true},
	}
}

func (BatteryFlow) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sn").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}
