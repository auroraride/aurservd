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
	return []ent.Field{
		field.String("reason").Comment("原因"),
		field.String("content").Comment("内容"),
		field.Uint8("status").Default(1).Comment("维修状态 1:维护中 2:已维修 3:维修失败 4:已取消"),
		field.Uint8("cabinet_status").Default(0).Comment("电柜状态 1:维护中 2:暂停维护"),
	}
}

// Edges of the AssetMaintenance.
func (AssetMaintenance) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("maintenance_details", AssetMaintenanceDetails.Type),
	}
}

func (AssetMaintenance) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
		CabinetMixin{Optional: true},
		MaintainerMixin{Optional: true},
	}
}

func (AssetMaintenance) Indexes() []ent.Index {
	return []ent.Index{}
}
