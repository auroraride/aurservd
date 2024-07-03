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

type RiderPhoneDeviceMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m RiderPhoneDeviceMixin) Fields() []ent.Field {
	relate := field.Uint64("device_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m RiderPhoneDeviceMixin) Edges() []ent.Edge {
	e := edge.To("device", RiderPhoneDevice.Type).Unique().Field("device_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m RiderPhoneDeviceMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("device_id"))
	}
	return
}

// RiderPhoneDevice holds the schema definition for the RiderPhoneDevice entity.
type RiderPhoneDevice struct {
	ent.Schema
}

// Annotations of the RiderPhoneDevice.
func (RiderPhoneDevice) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "rider_phone_device"},
		entsql.WithComments(true),
	}
}

// Fields of the RiderPhoneDevice.
func (RiderPhoneDevice) Fields() []ent.Field {
	return []ent.Field{
		field.String("device_sn").Optional().Comment("设备编号"),
		field.String("model").Optional().Comment("手机型号"),
		field.String("brand").Optional().Comment("手机品牌"),
		field.String("os_version").Optional().Comment("系统版本"),
		field.String("os_name").Optional().Comment("系统名称"),
		field.Uint64("screen_width").Optional().Comment("屏幕宽度"),
		field.Uint64("screen_height").Optional().Comment("屏幕高度"),
		field.String("imei").Optional().Comment("IMEI"),
		field.Uint64("rider_id").Unique().Comment("骑手ID"),
	}
}

// Edges of the RiderPhoneDevice.
func (RiderPhoneDevice) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (RiderPhoneDevice) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
	}
}

func (RiderPhoneDevice) Indexes() []ent.Index {
	return []ent.Index{}
}
