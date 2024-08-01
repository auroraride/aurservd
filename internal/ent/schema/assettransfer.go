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
		field.Uint8("status").Comment("调拨状态 1:配送中 2:已入库 3:已取消"),
		field.String("sn").Unique().Comment("调拨单号"),
		field.Uint8("from_location_type").Optional().Comment("开始位置类型 1:仓库 2:门店 3:站点 4:运维"),
		field.Uint64("from_location_id").Optional().Comment("开始位置ID"),
		field.Uint8("to_location_type").Optional().Comment("目标位置类型 1:仓库 2:门店 3:站点 4:运维"),
		field.Uint64("to_location_id").Optional().Comment("目标位置ID"),
		field.Uint("out_num").Optional().Comment("调出数量"),
		field.Uint("in_num").Optional().Comment("调入数量"),
		field.Uint64("out_user_id").Optional().Comment("出库人id"),
		field.Uint8("out_role_type").Optional().Comment("出库角色类型 1:资产后台管理员 2:运维人员 3:代理管理员 4:门店店员"),
		field.Uint64("in_user_id").Optional().Comment("入库人id"),
		field.Uint8("in_role_type").Optional().Comment("入库角色类型 1:资产后台管理员 2:运维人员 3:代理管理员 4:门店店员"),
		field.Time("out_time_at").Optional().Comment("出库时间"),
		field.Time("in_time_at").Optional().Comment("入库时间"),
		field.Uint8("transfer_type").Optional().Comment("调拨类型 1:初始入库 2:平台调拨 3:门店调拨 4:代理调拨 5:运维调拨 6:系统业务自动调拨"),
		field.String("reason").Optional().Comment("调拨事由"),
	}
}

// Edges of the AssetTransfer.
func (AssetTransfer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("details", AssetTransferDetails.Type),
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
