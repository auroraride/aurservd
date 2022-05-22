package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Manager holds the schema definition for the Manager entity.
type Manager struct {
    ent.Schema
}

// Annotations of the Manager.
func (Manager) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "manager"},
    }
}

// Fields of the Manager.
func (Manager) Fields() []ent.Field {
    return []ent.Field{
        field.String("phone").MaxLen(30).Unique().Comment("账户/手机号"),
        field.String("name").MaxLen(30).Comment("姓名"),
        field.String("password").Comment("密码"),
        field.Time("last_signin_at").Nillable().Optional().Comment("最后登录时间"),
    }
}

// Edges of the Manager.
func (Manager) Edges() []ent.Edge {
    return nil
}

func (Manager) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Manager) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("phone"),
    }
}
