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

type StationMixin struct {
	mixin.Schema
	DisableIndex bool
	Optional     bool
}

func (m StationMixin) Fields() []ent.Field {
	f := field.Uint64("station_id").Comment("站点ID")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m StationMixin) Edges() []ent.Edge {
	e := edge.To("station", EnterpriseStation.Type).Unique().Field("station_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m StationMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("station_id"))
	}
	return
}

// EnterpriseStation holds the schema definition for the EnterpriseStation entity.
type EnterpriseStation struct {
	ent.Schema
}

// Annotations of the EnterpriseStation.
func (EnterpriseStation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enterprise_station"},
		entsql.WithComments(true),
	}
}

// Fields of the EnterpriseStation.
func (EnterpriseStation) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("enterprise_id").Comment("企业ID"),
		field.String("name").Comment("站点名称"),
	}
}

// Edges of the EnterpriseStation.
func (EnterpriseStation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("enterprise", Enterprise.Type).Ref("stations").Unique().Required().Field("enterprise_id"),
		edge.From("agents", Agent.Type).Ref("stations"),

		edge.To("swap_putin_batteries", EnterpriseBatterySwap.Type),
		edge.To("swap_putout_batteries", EnterpriseBatterySwap.Type),

		edge.To("cabinets", Cabinet.Type),
		// edge.To("batteries", Battery.Type),
		edge.To("asset", Asset.Type),
		edge.To("stocks", Stock.Type),
	}
}

func (EnterpriseStation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		CityMixin{Optional: true},
	}
}

func (EnterpriseStation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("enterprise_id"),
	}
}
