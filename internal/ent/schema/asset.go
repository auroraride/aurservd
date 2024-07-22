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
		field.Uint8("type").Comment("资产类型 1:电车 2:电池 3:其它"),
		field.String("name").Comment("资产名称"),
		field.String("sn").Comment("资产编号"),
		field.Uint8("status").Comment("资产状态 0:待入库 1:库存中 2:配送中 3:使用中 4:故障 5:报废"),
		field.Bool("enable").Default(false).Comment("是否启用"),
		field.Uint8("locations_type").Default(1).Comment("资产位置类型 1:仓库 2:门店 3:电柜 4:站点 5:骑手 6:运维"),
		field.Uint64("locations_id").Comment("资产位置ID"),
		field.Uint8("scrap_reason_type").Optional().Comment("报废原因 1:丢失 2:损坏 3:其他"),
		field.Time("scrap_at").Optional().Comment("报废时间"),
		field.Uint64("scrap_operate_id").Optional().Comment("操作报废人员ID"),
		field.Uint64("brand_id").Optional().Comment("品牌ID"),
		field.Uint64("model_id").Optional().Comment("型号ID"),
		field.Uint64("rto_rider_id").Optional().Nillable().Comment("以租代购骑手ID，生成后禁止修改"),
		field.Time("inventory_at").Optional().Comment("盘点时间"),
	}
}

// Edges of the Asset.
func (Asset) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Asset) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Asset) Indexes() []ent.Index {
	return []ent.Index{}
}
