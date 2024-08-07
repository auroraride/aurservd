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

type AssetHistoryMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetHistoryMixin) Fields() []ent.Field {
	relate := field.Uint64("history_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetHistoryMixin) Edges() []ent.Edge {
	e := edge.To("history", AssetHistory.Type).Unique().Field("history_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetHistoryMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("history_id"))
	}
	return
}

// AssetHistory holds the schema definition for the AssetHistory entity.
type AssetHistory struct {
	ent.Schema
}

// Annotations of the AssetHistory.
func (AssetHistory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_history"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetHistory.
func (AssetHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("from_location_type").Optional().Comment("开始位置类型 1:仓库 2:门店 3:电柜 4:站点 5:骑手 6:运维"),
		field.Int("from_location_id").Optional().Comment("开始位置ID"),
		field.Int("to_location_type").Optional().Comment("目标位置类型 1:仓库 2:门店 3:电柜 4:站点 5:骑手 6:运维"),
		field.Int("to_location_id").Optional().Comment("目标位置ID"),
		field.Uint8("type").Optional().Comment("调拨类型 1:初始入库 2:平台调拨 3:门店调拨 4:代理调拨 5:运维调拨 6:系统业务自动调拨"),
	}
}

// Edges of the AssetHistory.
func (AssetHistory) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (AssetHistory) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
		AssetMixin{Optional: true},
	}
}

func (AssetHistory) Indexes() []ent.Index {
	return []ent.Index{}
}
