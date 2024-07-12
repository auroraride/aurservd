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

type BatteryNewMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m BatteryNewMixin) Fields() []ent.Field {
	relate := field.Uint64("new_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m BatteryNewMixin) Edges() []ent.Edge {
	e := edge.To("new", BatteryNew.Type).Unique().Field("new_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m BatteryNewMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("new_id"))
	}
	return
}

// BatteryNew holds the schema definition for the BatteryNew entity.
type BatteryNew struct {
	ent.Schema
}

// Annotations of the BatteryNew.
func (BatteryNew) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "battery_new"},
		entsql.WithComments(true),
	}
}

// Fields of the BatteryNew.
func (BatteryNew) Fields() []ent.Field {
	return []ent.Field{
		field.String("sn").Unique().Comment("电池编号"),
		field.Uint64("enterprise_id").Optional().Nillable().Comment("所属团签"),
		field.Uint8("asset_locations_type").Default(1).Comment("资产位置类型 1:仓库 2:门店 3:电柜 4:站点 5:骑手 6:运维"),
		field.Uint64("asset_locations_id").Comment("资产位置ID"),
		field.String("asset_locations").Optional().Comment("资产位置"),
		field.String("brand").Comment("品牌"),
		field.Bool("enable").Default(false).Comment("是否启用"),
		field.String("model").Comment("电池型号"),
		field.Uint8("asset_status").Default(1).Comment("资产状态0:待入库 1:库存中 2:配送中 3:使用中 4:故障 5:报废"),
		field.Uint64("status").Default(1).Comment("电池状态 1:正常 2:故障 3:报废"),
		field.Uint8("scrap_reason_type").Optional().Comment("报废原因 1:丢失 2:损坏 3:其他"),
		field.Time("scrap_at").Optional().Comment("报废时间"),
		field.Uint64("operate_id").Optional().Comment("操作报废人员ID"),
		field.Uint64("operate_role").Optional().Comment("操作人员角色"),
		field.String("operate_user").Optional().Comment("操作人员"),
		field.Uint64("warehouse_id").Optional().Comment("仓库ID"),
	}
}

// Edges of the BatteryNew.
func (BatteryNew) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (BatteryNew) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
		CityMixin{Optional: true},
	}
}

func (BatteryNew) Indexes() []ent.Index {
	return []ent.Index{}
}
