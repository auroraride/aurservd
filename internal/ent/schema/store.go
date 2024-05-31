package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type StoreMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
	Prefix       string
}

func (m StoreMixin) prefield() (string, string) {
	if m.Prefix == "" {
		return "store_id", "store"
	}
	return fmt.Sprintf("%s_store_id", m.Prefix), fmt.Sprintf("%sStore", m.Prefix)
}

func (m StoreMixin) Fields() []ent.Field {
	pf, _ := m.prefield()
	f := field.Uint64(pf).Comment("门店ID")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m StoreMixin) Edges() []ent.Edge {
	pf, pn := m.prefield()
	e := edge.To(pn, Store.Type).Unique().Field(pf)
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m StoreMixin) Indexes() (arr []ent.Index) {
	pf, _ := m.prefield()
	if !m.DisableIndex {
		arr = append(arr, index.Fields(pf))
	}
	return
}

// Store holds the schema definition for the Store entity.
type Store struct {
	ent.Schema
}

// Annotations of the Store.
func (Store) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "store"},
		entsql.WithComments(true),
	}
}

// Fields of the Store.
func (Store) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("employee_id").Optional().Nillable().Comment("上班员工ID"),
		field.Uint64("branch_id").Comment("网点ID"),
		field.String("sn").Immutable().Unique().Comment("门店编号"),
		field.String("name").Comment("门店名称"),
		field.Uint8("status").Default(0).Comment("门店状态 0维护 1营业 2休息 3隐藏"),
		field.Float("lng").Comment("经度"),
		field.Float("lat").Comment("纬度"),
		field.String("address").Comment("详细地址"),
		field.Bool("ebike_obtain").Default(false).Comment("是否可以领取车辆(租车)"),
		field.Bool("ebike_repair").Default(false).Comment("是否可以维修车辆"),
		field.Bool("ebike_sale").Default(false).Comment("是否可以购买车辆"),
		field.Bool("rest").Default(false).Comment("是否拥有驿站"),
		field.String("business_hours").Optional().Comment("营业时间"),
		field.Strings("photos").Optional().Comment("门店照片"),
	}
}

// Edges of the Store.
func (Store) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("branch", Branch.Type).Ref("stores").Required().Unique().Field("branch_id"),
		edge.From("employee", Employee.Type).Ref("store").Unique().Field("employee_id"),

		edge.To("stocks", Stock.Type),
		edge.To("attendances", Attendance.Type),
		edge.To("exceptions", Exception.Type),
		edge.To("goods", StoreGoods.Type),
	}
}

func (Store) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		CityMixin{},
	}
}

func (Store) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("branch_id"),
		index.Fields("employee_id"),
		index.Fields("status"),
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}
