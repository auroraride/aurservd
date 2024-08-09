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

type AssetTransferMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetTransferMixin) Fields() []ent.Field {
	relate := field.Uint64("transfer_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetTransferMixin) Edges() []ent.Edge {
	e := edge.To("transfer", AssetTransfer.Type).Unique().Field("transfer_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetTransferMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("transfer_id"))
	}
	return
}

// AssetTransfer holds the schema definition for the AssetTransfer entity.
type AssetTransfer struct {
	ent.Schema
}

// Annotations of the AssetTransfer.
func (AssetTransfer) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_transfer"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetTransfer.
func (AssetTransfer) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("status").Default(0).Comment("调拨状态 1:配送中 2:待入库 3:已入库 4:已取消"),
		field.String("sn").Unique().Comment("调拨单号"),
		field.Uint8("from_location_type").Optional().Nillable().Comment("开始位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手"),
		field.Uint64("from_location_id").Optional().Nillable().Comment("开始位置ID"),
		field.Uint8("to_location_type").Optional().Comment("目标位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手"),
		field.Uint64("to_location_id").Optional().Comment("目标位置ID"),
		field.Uint("out_num").Optional().Comment("调出数量"),
		field.Uint("in_num").Optional().Comment("调入数量"),
		field.Uint64("out_operate_id").Optional().Nillable().Comment("出库人id"),
		field.Uint8("out_operate_type").Optional().Nillable().Comment("出库角色类型 1:资产后台 2:门店 3:代理 4:运维 5:电柜 6:骑手"),
		field.Time("out_time_at").Optional().Nillable().Comment("出库时间"),
		field.String("reason").Optional().Comment("调拨事由"),
		field.Uint8("type").Optional().Comment("调拨类型 1:初始入库 2:调拨 3:激活 4:寄存 5:取消寄存 6:退租"),
	}
}

// Edges of the AssetTransfer.
func (AssetTransfer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("transfer_details", AssetTransferDetails.Type),
		// 调拨from位置
		edge.To("from_location_store", Store.Type).Unique().Field("from_location_id"),               // 关联门店
		edge.To("from_location_cabinet", Cabinet.Type).Unique().Field("from_location_id"),           // 关联电柜
		edge.To("from_location_station", EnterpriseStation.Type).Unique().Field("from_location_id"), // 关联站点
		edge.To("from_location_rider", Rider.Type).Unique().Field("from_location_id"),               // 关联骑手
		edge.To("from_location_operator", Maintainer.Type).Unique().Field("from_location_id"),       // 关联运维
		edge.To("from_location_warehouse", Warehouse.Type).Unique().Field("from_location_id"),       // 关联仓库
		// 调拨to位置
		edge.To("to_location_store", Store.Type).Unique().Field("to_location_id"),               // 关联门店
		edge.To("to_location_cabinet", Cabinet.Type).Unique().Field("to_location_id"),           // 关联电柜
		edge.To("to_location_station", EnterpriseStation.Type).Unique().Field("to_location_id"), // 关联站点
		edge.To("to_location_rider", Rider.Type).Unique().Field("to_location_id"),               // 关联骑手
		edge.To("to_location_operator", Maintainer.Type).Unique().Field("to_location_id"),       // 关联运维
		edge.To("to_location_warehouse", Warehouse.Type).Unique().Field("to_location_id"),       // 关联仓库
		// 出库关联操作人员
		edge.To("out_operate_manager", Manager.Type).Unique().Field("out_operate_id"),       // 资产后台
		edge.To("out_operate_store", Store.Type).Unique().Field("out_operate_id"),           // 门店
		edge.To("out_operate_agent", Agent.Type).Unique().Field("out_operate_id"),           // 代理
		edge.To("out_operate_maintainer", Maintainer.Type).Unique().Field("out_operate_id"), // 运维
		edge.To("out_operate_cabinet", Cabinet.Type).Unique().Field("out_operate_id"),       // 电柜
		edge.To("out_operate_rider", Rider.Type).Unique().Field("out_operate_id"),           // 骑手
	}
}

func (AssetTransfer) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (AssetTransfer) Indexes() []ent.Index {
	return []ent.Index{}
}
