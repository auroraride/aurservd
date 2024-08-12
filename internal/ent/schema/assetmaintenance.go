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

type AssetMaintenanceMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetMaintenanceMixin) Fields() []ent.Field {
	relate := field.Uint64("maintenance_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetMaintenanceMixin) Edges() []ent.Edge {
	e := edge.To("maintenance", AssetMaintenance.Type).Unique().Field("maintenance_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetMaintenanceMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("maintenance_id"))
	}
	return
}

// AssetMaintenance holds the schema definition for the AssetMaintenance entity.
type AssetMaintenance struct {
	ent.Schema
}

// Annotations of the AssetMaintenance.
func (AssetMaintenance) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_maintenance"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetMaintenance.
func (AssetMaintenance) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the AssetMaintenance.
func (AssetMaintenance) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (AssetMaintenance) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (AssetMaintenance) Indexes() []ent.Index {
	return []ent.Index{}
}
