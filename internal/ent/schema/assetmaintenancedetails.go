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

type AssetMaintenanceDetailsMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetMaintenanceDetailsMixin) Fields() []ent.Field {
	relate := field.Uint64("details_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetMaintenanceDetailsMixin) Edges() []ent.Edge {
	e := edge.To("details", AssetMaintenanceDetails.Type).Unique().Field("details_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetMaintenanceDetailsMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("details_id"))
	}
	return
}

// AssetMaintenanceDetails holds the schema definition for the AssetMaintenanceDetails entity.
type AssetMaintenanceDetails struct {
	ent.Schema
}

// Annotations of the AssetMaintenanceDetails.
func (AssetMaintenanceDetails) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_maintenance_details"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetMaintenanceDetails.
func (AssetMaintenanceDetails) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the AssetMaintenanceDetails.
func (AssetMaintenanceDetails) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (AssetMaintenanceDetails) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (AssetMaintenanceDetails) Indexes() []ent.Index {
	return []ent.Index{}
}
