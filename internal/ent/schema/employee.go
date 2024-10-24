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
	"github.com/google/uuid"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type EmployeeMixin struct {
	mixin.Schema
	Optional     bool
	Immutable    bool
	DisableIndex bool
}

func (m EmployeeMixin) Fields() []ent.Field {
	f := field.Uint64("employee_id").Comment("店员ID")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m EmployeeMixin) Edges() []ent.Edge {
	e := edge.To("employee", Employee.Type).Unique().Field("employee_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m EmployeeMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("employee_id"))
	}
	return
}

// Employee holds the schema definition for the Employee entity.
type Employee struct {
	ent.Schema
}

// Annotations of the Employee.
func (Employee) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "employee"},
		entsql.WithComments(true),
	}
}

// Fields of the Employee.
func (Employee) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("sn", uuid.New()).Optional().Unique(),
		field.String("name").Comment("姓名"),
		field.String("phone").Comment("电话"),
		field.Bool("enable").Default(true).Comment("启用状态"),
		field.String("password").Optional().Comment("密码"),
		field.Uint("limit").Default(0).Comment("限制范围(m)"),
		field.Uint64("duty_store_id").Optional().Nillable().Comment("上班门店ID"),
	}
}

// Edges of the Employee.
func (Employee) Edges() []ent.Edge {
	return []ent.Edge{
		// edge.To("subscribes", Subscribe.Type),
		edge.To("store", Store.Type).Unique(),
		edge.To("attendances", Attendance.Type),
		edge.To("stocks", Stock.Type),
		edge.To("exchanges", Exchange.Type),
		edge.To("commissions", Commission.Type),
		edge.To("assistances", Assistance.Type),
		edge.From("stores", Store.Type).Ref("employees"),
		edge.From("duty_store", Store.Type).Ref("duty_employees").Field("duty_store_id").Unique(),
	}
}

func (Employee) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		CityMixin{Optional: false},
		StoreGroupMixin{Optional: true},
	}
}

func (Employee) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("enable"),
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
		index.Fields("phone").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}
