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

type AssetMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetMixin) Fields() []ent.Field {
	relate := field.Uint64("asset_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetMixin) Edges() []ent.Edge {
	e := edge.To("asset", Asset.Type).Unique().Field("asset_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("asset_id"))
	}
	return
}

// Asset holds the schema definition for the Asset entity.
type Asset struct {
	ent.Schema
}

// Annotations of the Asset.
func (Asset) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset"},
		entsql.WithComments(true),
	}
}

// Fields of the Asset.
func (Asset) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Unique(),
		field.Uint8("type").Comment("资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它"),
		field.String("name").Comment("资产名称"),
		field.String("sn").Optional().Comment("资产编号"),
		field.Uint8("status").Default(0).Comment("资产状态 1:库存中 2:配送中 3:使用中 4:故障 5:报废"),
		field.Bool("enable").Default(false).Comment("是否启用"),
		field.Uint8("locations_type").Optional().Comment("资产位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手"),
		field.Uint64("locations_id").Optional().Comment("资产位置ID"),
		field.Uint64("rto_rider_id").Optional().Nillable().Comment("以租代购骑手ID，生成后禁止修改"),
		field.Time("check_at").Optional().Nillable().Comment("盘点时间"),
		field.String("brand_name").Optional().Comment("品牌名称"),
		field.Uint64("subscribe_id").Optional().Nillable().Comment("订阅ID"),
		field.Int("ordinal").Optional().Nillable().Comment("仓位号"),
	}
}

// Edges of the Asset.
func (Asset) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("values", AssetAttributeValues.Type), // 关联资产属性

		edge.To("scrap_details", AssetScrapDetails.Type),             // 关联报废详情
		edge.To("transfer_details", AssetTransferDetails.Type),       // 关联调拨详情
		edge.To("maintenance_details", AssetMaintenanceDetails.Type), // 关联维护详情
		edge.To("check_details", AssetCheckDetails.Type),             // 关联盘点详情

		// edge.From("subscribe", Subscribe.Type).Unique().Ref("battery").Field("subscribe_id"), // 关联订阅

		edge.From("warehouse", Warehouse.Type).Unique().Ref("asset").Field("locations_id"),       // 关联仓库
		edge.From("store", Store.Type).Unique().Ref("asset").Field("locations_id"),               // 关联门店
		edge.From("cabinet", Cabinet.Type).Unique().Ref("asset").Field("locations_id"),           // 关联电柜
		edge.From("station", EnterpriseStation.Type).Ref("asset").Unique().Field("locations_id"), // 关联站点
		edge.From("rider", Rider.Type).Unique().Ref("asset").Field("locations_id"),               // 关联骑手
		edge.From("operator", Maintainer.Type).Unique().Ref("asset").Field("locations_id"),       // 关联运维

		edge.To("ebike_allocates", Allocate.Type),
		edge.To("battery_allocates", Allocate.Type),
		edge.To("rto_rider", Rider.Type).Unique().Field("rto_rider_id"),

		edge.From("battery_rider", Rider.Type).Ref("battery").Unique().Field("locations_id").Comment("所属骑手"),
	}
}

func (Asset) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
		EbikeBrandMixin{Optional: true},
		BatteryModelMixin{Optional: true},
		CityMixin{Optional: true},
		MaterialMixin{Optional: true},
	}
}

func (Asset) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("sn"),
		index.Fields("sn", "type"),
		index.Fields("locations_id"),
		index.Fields("locations_type"),
		index.Fields("locations_type", "locations_id"),
	}
}
