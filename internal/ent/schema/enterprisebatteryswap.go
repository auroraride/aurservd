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

type EnterpriseBatterySwapMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m EnterpriseBatterySwapMixin) Fields() []ent.Field {
	relate := field.Uint64("swap_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m EnterpriseBatterySwapMixin) Edges() []ent.Edge {
	e := edge.To("swap", EnterpriseBatterySwap.Type).Unique().Field("swap_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m EnterpriseBatterySwapMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("swap_id"))
	}
	return
}

// EnterpriseBatterySwap holds the schema definition for the EnterpriseBatterySwap entity.
type EnterpriseBatterySwap struct {
	ent.Schema
}

// Annotations of the EnterpriseBatterySwap.
func (EnterpriseBatterySwap) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enterprise_battery_swap"},
		entsql.WithComments(true),
	}
}

// Fields of the EnterpriseBatterySwap.
func (EnterpriseBatterySwap) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("putin_battery_id").Comment("放入的电池ID"),
		field.String("putin_battery_sn").Comment("放入的电池编码"),
		field.Uint64("putin_enterprise_id").Optional().Nillable().Comment("放入归属团签ID, 空值是平台骑手放入"),
		field.Uint64("putin_station_id").Optional().Nillable().Comment("放入归属站点ID, 空值是平台骑手放入"),

		field.Uint64("putout_battery_id").Comment("取出的电池ID"),
		field.String("putout_battery_sn").Comment("取出的电池编码"),
		field.Uint64("putout_enterprise_id").Optional().Nillable().Comment("取出归属团签ID, 空值是从平台电柜取出"),
		field.Uint64("putout_station_id").Optional().Nillable().Comment("取出归属站点ID, 空值是从平台电柜取出"),
	}
}

// Edges of the EnterpriseBatterySwap.
func (EnterpriseBatterySwap) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("putin_battery", Battery.Type).Required().Unique().Field("putin_battery_id"),
		edge.From("putin_enterprise", Enterprise.Type).Ref("swap_putin_batteries").Unique().Field("putin_enterprise_id"),
		edge.From("putin_station", EnterpriseStation.Type).Ref("swap_putin_batteries").Unique().Field("putin_station_id"),

		edge.To("putout_battery", Battery.Type).Required().Unique().Field("putout_battery_id"),
		edge.From("putout_enterprise", Enterprise.Type).Ref("swap_putout_batteries").Unique().Field("putout_enterprise_id"),
		edge.From("putout_station", EnterpriseStation.Type).Ref("swap_putout_batteries").Unique().Field("putout_station_id"),
	}
}

func (EnterpriseBatterySwap) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},

		ExchangeMixin{},
		CabinetMixin{},
	}
}

func (EnterpriseBatterySwap) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("putin_battery_id"),
		index.Fields("putin_battery_sn"),
		index.Fields("putin_enterprise_id"),
		index.Fields("putin_station_id"),

		index.Fields("putout_battery_id"),
		index.Fields("putout_battery_sn"),
		index.Fields("putout_enterprise_id"),
		index.Fields("putout_station_id"),
	}
}
