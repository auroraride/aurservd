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
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type ManagerMixin struct {
    mixin.Schema
    Optional bool
}

func (mmm ManagerMixin) Fields() []ent.Field {
    f := field.Uint64("manager_id").Comment("管理人ID")
    if mmm.Optional {
        f.Optional().Nillable()
    }
    return []ent.Field{f}
}

func (mmm ManagerMixin) Edges() []ent.Edge {
    e := edge.To("manager", Manager.Type).Unique().Field("manager_id")
    if !mmm.Optional {
        e.Required()
    }
    return []ent.Edge{
        e,
    }
}

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
        field.String("phone").MaxLen(30).Comment("账户/手机号"),
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
        index.Fields("phone").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
        index.Fields("name").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
