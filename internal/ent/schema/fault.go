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

type FaultMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m FaultMixin) Fields() []ent.Field {
	relate := field.Uint64("fault_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m FaultMixin) Edges() []ent.Edge {
	e := edge.To("fault", Fault.Type).Unique().Field("fault_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m FaultMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("fault_id"))
	}
	return
}

// Fault holds the schema definition for the Fault entity.
type Fault struct {
	ent.Schema
}

// Annotations of the Fault.
func (Fault) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fault"},
		entsql.WithComments(true),
	}
}

// Fields of the Fault.
func (Fault) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("status").Default(0).Comment("故障状态 0未处理 1已处理"),
		field.String("description").Optional().Comment("故障留言"),
		field.JSON("attachments", []string{}).Optional().Comment("附件"),
		field.Uint8("type").Default(1).Comment("故障类型 1:电柜故障 2:电池故障 3:车辆故障 4:其他"),
		field.Strings("fault").Optional().Comment("故障内容"),
		field.Uint64("ebike_id").Optional().Nillable().Comment("电车ID"),
		field.Uint64("battery_id").Optional().Nillable().Comment("电池ID"),
	}
}

// Edges of the Fault.
func (Fault) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ebike", Asset.Type).Unique().Field("ebike_id"),
		edge.To("battery", Asset.Type).Unique().Field("battery_id"),
	}
}

func (Fault) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
		CityMixin{},
		CabinetMixin{Optional: true},
		RiderMixin{Optional: true},
	}
}

func (Fault) Indexes() []ent.Index {
	return []ent.Index{}
}
