package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type StockSummaryMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m StockSummaryMixin) Fields() []ent.Field {
	relate := field.Uint64("summary_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m StockSummaryMixin) Edges() []ent.Edge {
	e := edge.To("summary", StockSummary.Type).Unique().Field("summary_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m StockSummaryMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("summary_id"))
	}
	return
}

// StockSummary holds the schema definition for the StockSummary entity.
type StockSummary struct {
	ent.Schema
}

// Annotations of the StockSummary.
func (StockSummary) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "stock_summary"},
		entsql.WithComments(true),
	}
}

// Fields of the StockSummary.
func (StockSummary) Fields() []ent.Field {
	return []ent.Field{
		field.String("date").NotEmpty().Comment("日期"),
		field.Int("battery_num").Default(0).Comment("电池总数"),
		field.Int("battery_outbound_num").Default(0).Comment("电池出库总数"),
		field.Int("battery_inbound_num").Default(0).Comment("电池入库总数"),
		field.Int("bike_num").Default(0).Comment("电车总数"),
		field.Int("bike_outbound_num").Default(0).Comment("电车出库总数"),
		field.Int("bike_inbound_num").Default(0).Comment("电车入库总数"),
		field.Int("cabinet_battery_num").Default(0).Comment("电柜电池总数"),
		field.Int("rider_battery_num").Default(0).Comment("骑手电池总数"),
	}
}

// Edges of the StockSummary.
func (StockSummary) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (StockSummary) Mixin() []ent.Mixin {
	return []ent.Mixin{
		EnterpriseMixin{Optional: true},
		StationMixin{Optional: true},
		StoreMixin{Optional: true},
		RiderMixin{Optional: true},
		CabinetMixin{Optional: true},
	}
}

func (StockSummary) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("date"),
	}
}
