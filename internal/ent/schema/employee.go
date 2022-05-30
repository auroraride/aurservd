package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type EmployeeMixin struct {
    mixin.Schema
}

func (EmployeeMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("employee_id").Optional().Comment("操作店员ID"),
    }
}

func (EmployeeMixin) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("employee", Employee.Type).Unique().Field("employee_id"),
    }
}

// Employee holds the schema definition for the Employee entity.
type Employee struct {
    ent.Schema
}

// Annotations of the Employee.
func (Employee) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "employee"},
    }
}

// Fields of the Employee.
func (Employee) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Comment("姓名"),
    }
}

// Edges of the Employee.
func (Employee) Edges() []ent.Edge {
    return []ent.Edge{
        // edge.To("subscribes", Subscribe.Type),
    }
}

func (Employee) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Employee) Indexes() []ent.Index {
    return []ent.Index{}
}
