package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type WarehouseMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
	Prefix       string
}

func (m WarehouseMixin) prefield() (string, string) {
	if m.Prefix == "" {
		return "warehouse_id", "warehouse"
	}
	return fmt.Sprintf("%s_warehouse_id", m.Prefix), fmt.Sprintf("%sWarehouse", m.Prefix)
}

func (m WarehouseMixin) Fields() []ent.Field {
	pf, _ := m.prefield()
	f := field.Uint64(pf).Comment("仓库ID")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m WarehouseMixin) Edges() []ent.Edge {
	pf, pn := m.prefield()
	e := edge.To(pn, Warehouse.Type).Unique().Field(pf)
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m WarehouseMixin) Indexes() (arr []ent.Index) {
	pf, _ := m.prefield()
	if !m.DisableIndex {
		arr = append(arr, index.Fields(pf))
	}
	return
}

// Warehouse holds the schema definition for the Warehouse entity.
type Warehouse struct {
	ent.Schema
}

// Annotations of the Warehouse.
func (Warehouse) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "warehouse"},
		entsql.WithComments(true),
	}
}

// Fields of the Warehouse.
func (Warehouse) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("仓库名称"),
		field.Float("lng").Comment("经度"),
		field.Float("lat").Comment("纬度"),
		field.String("address").Optional().Comment("详细地址"),
		field.String("sn").Immutable().Unique().Comment("仓库编号"),
		field.Uint64("asset_manager_id").Optional().Nillable().Comment("上班仓管员ID"),
	}
}

// Edges of the Warehouse.
func (Warehouse) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("asset_manager", AssetManager.Type).Ref("warehouse").Unique().Field("asset_manager_id"),
		edge.To("asset_managers", AssetManager.Type),
		edge.To("asset", Asset.Type),
	}
}

func (Warehouse) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
		CityMixin{},
	}
}

func (Warehouse) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}
