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

type AssetCheckDetailsMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetCheckDetailsMixin) Fields() []ent.Field {
	relate := field.Uint64("details_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetCheckDetailsMixin) Edges() []ent.Edge {
	e := edge.To("details", AssetCheckDetails.Type).Unique().Field("details_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetCheckDetailsMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("details_id"))
	}
	return
}

// AssetCheckDetails holds the schema definition for the AssetCheckDetails entity.
type AssetCheckDetails struct {
	ent.Schema
}

// Annotations of the AssetCheckDetails.
func (AssetCheckDetails) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_check_details"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetCheckDetails.
func (AssetCheckDetails) Fields() []ent.Field {
	return []ent.Field{
		field.String("sn").Optional().Comment("资产SN"),
		field.Uint64("asset_id").Optional().Comment("资产ID"),
		field.Uint64("check_id").Optional().Comment("盘点ID"),
		field.Uint64("real_locations_id").Optional().Comment("实际位置ID"),
		field.Uint8("real_locations_type").Optional().Comment("实际位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手"),
		field.Uint64("locations_id").Optional().Comment("原位置ID"),
		field.Uint8("locations_type").Optional().Comment("原位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手"),
		field.Uint8("status").Optional().Default(0).Comment("处理状态 0:未处理 1:已入库 2:已出库 3:已报废"),
		field.Uint8("result").Optional().Default(0).Comment("盘点结果 0:未盘点 1:正常 2:亏 3:盈"),
		field.Uint64("operate_id").Optional().Comment("操作人id"),
		field.Time("operate_at").Optional().Nillable().Comment("处理时间"),
	}
}

// Edges of the AssetCheckDetails.
func (AssetCheckDetails) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("asset", Asset.Type).Ref("check_details").Unique().Field("asset_id"),
		edge.From("check", AssetCheck.Type).Ref("check_details").Unique().Field("check_id"),

		edge.To("warehouse", Warehouse.Type).Unique().Field("locations_id"),       // 关联仓库
		edge.To("store", Store.Type).Unique().Field("locations_id"),               // 关联门店
		edge.To("cabinet", Cabinet.Type).Unique().Field("locations_id"),           // 关联电柜
		edge.To("station", EnterpriseStation.Type).Unique().Field("locations_id"), // 关联站点
		edge.To("rider", Rider.Type).Unique().Field("locations_id"),               // 关联骑手
		edge.To("operator", Maintainer.Type).Unique().Field("locations_id"),       // 关联运维

		edge.To("real_warehouse", Warehouse.Type).Unique().Field("real_locations_id"),       // 关联仓库
		edge.To("real_store", Store.Type).Unique().Field("real_locations_id"),               // 关联门店
		edge.To("real_cabinet", Cabinet.Type).Unique().Field("real_locations_id"),           // 关联电柜
		edge.To("real_station", EnterpriseStation.Type).Unique().Field("real_locations_id"), // 关联站点
		edge.To("real_rider", Rider.Type).Unique().Field("real_locations_id"),               // 关联骑手
		edge.To("real_operator", Maintainer.Type).Unique().Field("real_locations_id"),       // 关联运维
	}
}

func (AssetCheckDetails) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
		MaintainerMixin{Optional: true},
	}
}

func (AssetCheckDetails) Indexes() []ent.Index {
	return []ent.Index{}
}
