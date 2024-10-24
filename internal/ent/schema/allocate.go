package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/internal"
)

type AllocateMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AllocateMixin) Fields() []ent.Field {
	relate := field.Uint64("allocate_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AllocateMixin) Edges() []ent.Edge {
	e := edge.To("allocate", Allocate.Type).Unique().Field("allocate_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AllocateMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("allocate_id"))
	}
	return
}

// Allocate holds the schema definition for the Allocate entity.
type Allocate struct {
	ent.Schema
}

// Annotations of the Allocate.
func (Allocate) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "allocate"},
		entsql.WithComments(true),
	}
}

// Fields of the Allocate.
func (Allocate) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").Values(model.SubscribeTypeBattery, model.SubscribeTypeEbike).Comment("分配类型"),
		field.Uint8("status").Comment("分配状态"),
		field.Time("time").Comment("分配时间"),
		field.String("model").Comment("电池型号"),
		field.Uint64("ebike_id").Optional().Nillable().Comment("电车ID"),
		field.Uint64("battery_id").Optional().Nillable().Comment("电池ID"),
	}
}

// Edges of the Allocate.
func (Allocate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("contract", Contract.Type).Unique(),

		edge.From("ebike", Asset.Type).Ref("ebike_allocates").Unique().Field("ebike_id"),
		edge.From("battery", Asset.Type).Ref("battery_allocates").Unique().Field("battery_id"),
	}
}

func (Allocate) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{DisableIndex: true, Optional: true},

		RiderMixin{},
		SubscribeMixin{Unique: true},

		internal.Modifier{},
		EmployeeMixin{Optional: true},

		CabinetMixin{Optional: true},
		StoreMixin{Optional: true},

		EbikeBrandMixin{Optional: true},

		// BatteryMixin{Optional: true},
		StationMixin{Optional: true},

		AgentMixin{Optional: true},
	}
}

func (Allocate) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("time"),
	}
}
