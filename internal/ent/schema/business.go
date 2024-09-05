package schema

import (
	"ariga.io/atlas/sql/postgres"
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

type BusinessMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m BusinessMixin) Fields() []ent.Field {
	f := field.Uint64("business_id")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m BusinessMixin) Edges() []ent.Edge {
	e := edge.To("business", Business.Type).Unique().Field("business_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m BusinessMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("business_id"))
	}
	return
}

// Business holds the schema definition for the Business entity.
type Business struct {
	ent.Schema
}

// Annotations of the Business.
func (Business) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "business"},
		entsql.WithComments(true),
	}
}

// Fields of the Business.
func (Business) Fields() []ent.Field {
	return []ent.Field{
		field.Other("type", model.BusinessTypeActive).SchemaType(map[string]string{
			dialect.Postgres: postgres.TypeCharVar,
		}).Comment("业务类型"),
		field.JSON("bin_info", &model.BinInfo{}).Optional().Comment("仓位信息"),
		field.String("asset_transfer_sn").Optional().Comment("出入库编码"),
		field.Uint64("rto_ebike_id").Optional().Nillable().Comment("以租代购车辆ID，生成后禁止修改"),
	}
}

// Edges of the Business.
func (Business) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("rto_ebike", Asset.Type).Unique().Field("rto_ebike_id"),
	}
}

func (Business) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		RiderMixin{},
		CityMixin{},
		SubscribeMixin{},

		EmployeeMixin{Optional: true},
		StoreMixin{Optional: true},
		PlanMixin{Optional: true},
		EnterpriseMixin{Optional: true},
		StationMixin{Optional: true},
		CabinetMixin{Optional: true},
		BatteryMixin{Optional: true},
		AgentMixin{Optional: true},
	}
}

func (Business) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("asset_transfer_sn"),
		index.Fields("rto_ebike_id"),
	}
}
