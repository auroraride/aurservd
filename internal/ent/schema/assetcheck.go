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

type AssetCheckMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetCheckMixin) Fields() []ent.Field {
	relate := field.Uint64("check_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetCheckMixin) Edges() []ent.Edge {
	e := edge.To("check", AssetCheck.Type).Unique().Field("check_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetCheckMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("check_id"))
	}
	return
}

// AssetCheck holds the schema definition for the AssetCheck entity.
type AssetCheck struct {
	ent.Schema
}

// Annotations of the AssetCheck.
func (AssetCheck) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_check"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetCheck.
func (AssetCheck) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("status").Optional().Comment("盘点状态 1:待盘点 2:已盘点"),
		field.Uint("battery_num").Optional().Comment("应盘点电池数量"),
		field.Uint("battery_num_real").Optional().Comment("实盘电池数量"),
		field.Uint("ebike_num").Optional().Comment("应盘电车数量"),
		field.Uint("ebike_num_real").Optional().Comment("实盘电车数量"),
		field.Uint64("operate_id").Optional().Comment("盘点人id"),
		field.Uint8("operate_type").Optional().Comment("盘点角色类型 1:门店 3:代理 6:资产后台"),
		field.Uint8("locations_type").Optional().Comment("盘点位置类型 1:仓库 2:门店 3:代理"),
		field.Uint64("locations_id").Optional().Comment("盘点位置id"),
		field.Time("start_at").Optional().Nillable().Comment("盘点开始时间"),
		field.Time("end_at").Optional().Nillable().Comment("盘点结束时间"),
	}
}

// Edges of the AssetCheck.
func (AssetCheck) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("check_details", AssetCheckDetails.Type),
		edge.To("operate_asset_manager", AssetManager.Type).Unique().Field("operate_id"), // 资产后台
		edge.To("operate_employee", Employee.Type).Unique().Field("operate_id"),          // 门店
		edge.To("operate_agent", Agent.Type).Unique().Field("operate_id"),                // 代理
		edge.To("warehouse", Warehouse.Type).Unique().Field("locations_id"),              // 关联仓库
		edge.To("store", Store.Type).Unique().Field("locations_id"),                      // 关联门店
		edge.To("station", EnterpriseStation.Type).Unique().Field("locations_id"),        // 关联站点
	}
}

func (AssetCheck) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (AssetCheck) Indexes() []ent.Index {
	return []ent.Index{}
}
