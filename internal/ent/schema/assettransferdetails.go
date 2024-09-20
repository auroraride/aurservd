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

type AssetTransferDetailsMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetTransferDetailsMixin) Fields() []ent.Field {
	relate := field.Uint64("details_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetTransferDetailsMixin) Edges() []ent.Edge {
	e := edge.To("details", AssetTransferDetails.Type).Unique().Field("details_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetTransferDetailsMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("details_id"))
	}
	return
}

// AssetTransferDetails holds the schema definition for the AssetTransferDetails entity.
type AssetTransferDetails struct {
	ent.Schema
}

// Annotations of the AssetTransferDetails.
func (AssetTransferDetails) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_transfer_details"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetTransferDetails.
func (AssetTransferDetails) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("transfer_id").Optional().Comment("调拨ID"),
		field.Bool("is_in").Default(false).Comment("是否入库"),
		field.Uint64("in_operate_id").Optional().Comment("入库人id"),
		field.Uint8("in_operate_type").Optional().Comment("入库角色类型 1:资产后台 2:门店 3:代理 4:运维 5:电柜 6:骑手"),
		field.Time("in_time_at").Optional().Nillable().Comment("入库时间"),
		field.Uint64("asset_id").Optional().Comment("资产ID"),
	}
}

// Edges of the AssetTransferDetails.
func (AssetTransferDetails) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transfer", AssetTransfer.Type).Ref("transfer_details").Unique().Field("transfer_id"),

		// 入库关联操作人员
		edge.To("in_operate_asset_manager", AssetManager.Type).Unique().Field("in_operate_id"), // 资产后台
		edge.To("in_operate_employee", Employee.Type).Unique().Field("in_operate_id"),          // 门店店员
		edge.To("in_operate_agent", Agent.Type).Unique().Field("in_operate_id"),                // 代理
		edge.To("in_operate_maintainer", Maintainer.Type).Unique().Field("in_operate_id"),      // 运维
		edge.To("in_operate_cabinet", Cabinet.Type).Unique().Field("in_operate_id"),            // 电柜
		edge.To("in_operate_rider", Rider.Type).Unique().Field("in_operate_id"),                // 骑手
		edge.To("in_operate_manager", Manager.Type).Unique().Field("in_operate_id"),            // 业务后台

		// 关联资产
		edge.From("asset", Asset.Type).Ref("transfer_details").Unique().Field("asset_id"),
	}
}

func (AssetTransferDetails) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (AssetTransferDetails) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("asset_id"),
		index.Fields("in_operate_id"),
		index.Fields("in_operate_type"),
		index.Fields("transfer_id", "asset_id"),
	}
}
