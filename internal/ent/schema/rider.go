package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/internal"
)

type RiderMixin struct {
	mixin.Schema
	DisableIndex bool
	Optional     bool
}

func (m RiderMixin) Fields() []ent.Field {
	f := field.Uint64("rider_id").Comment("骑手ID")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m RiderMixin) Edges() []ent.Edge {
	e := edge.To("rider", Rider.Type).Unique().Field("rider_id").Comment("骑手")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m RiderMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("rider_id"))
	}
	return
}

// Rider holds the schema definition for the Rider entity.
type Rider struct {
	ent.Schema
}

// Annotations of the Rider.
func (Rider) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "rider"},
		entsql.WithComments(true),
	}
}

// Fields of the Rider.
func (Rider) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("person_id").Optional().Nillable().Comment("身份"),
		field.String("name").Optional().Comment("骑手姓名"),
		field.String("id_card_number").Optional().Comment("身份证号"),
		field.Uint64("enterprise_id").Optional().Nillable().Comment("所属企业"),
		field.String("phone").MaxLen(11).Comment("手机号"),
		field.JSON("contact", &model.RiderContact{}).Optional().Comment("紧急联系人"),
		field.Uint8("device_type").Optional().Comment("登录设备类型: 1iOS 2Android"),
		field.String("last_device").Optional().MaxLen(60).Comment("最近登录设备"),
		field.Bool("is_new_device").Default(false).Comment("是否新设备"),
		field.String("push_id").MaxLen(60).Optional().Comment("推送ID"),
		field.Time("last_signin_at").Nillable().Optional().Comment("最后登录时间"),
		field.Bool("blocked").Default(false).Comment("是否封禁骑手账号"),
		field.Int64("points").Default(0).Comment("骑手积分"),
		field.JSON("exchange_limit", model.RiderExchangeLimit{}).Optional().Comment("换电间隔配置"),
		field.JSON("exchange_frequency", model.RiderExchangeFrequency{}).Optional().Comment("换电频次配置"),
		field.Time("join_enterprise_at").Optional().Nillable().Comment("加入团签时间"),
	}
}

// Edges of the Rider.
func (Rider) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("person", Person.Type).Ref("riders").Unique().Field("person_id"),
		edge.From("enterprise", Enterprise.Type).Ref("riders").Unique().Field("enterprise_id"),

		edge.To("contracts", Contract.Type),
		edge.To("faults", CabinetFault.Type),
		edge.To("orders", Order.Type),

		edge.To("exchanges", Exchange.Type).Comment("换电记录"),
		edge.To("subscribes", Subscribe.Type).Comment("订阅"),

		edge.To("asset", Asset.Type),
		edge.To("stocks", Stock.Type),
		edge.To("followups", RiderFollowUp.Type),

		edge.To("battery", Asset.Type).Unique(),
		edge.To("battery_flows", BatteryFlow.Type),
	}
}

func (Rider) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		StationMixin{Optional: true},
	}
}

func (Rider) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("phone").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
		index.Fields("id_card_number"),
		index.Fields("person_id"),
		index.Fields("enterprise_id"),
		index.Fields("last_device"),
		index.Fields("push_id"),
	}
}
