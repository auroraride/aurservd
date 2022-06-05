package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// EnterpriseContract holds the schema definition for the EnterpriseContract entity.
type EnterpriseContract struct {
    ent.Schema
}

// Annotations of the EnterpriseContract.
func (EnterpriseContract) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "enterprise_contract"},
    }
}

// Fields of the EnterpriseContract.
func (EnterpriseContract) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("enterprise_id"),
        field.Time("start").SchemaType(map[string]string{dialect.Postgres: "date"}).Comment("合同开始时间"),
        field.Time("end").SchemaType(map[string]string{dialect.Postgres: "date"}).Comment("合同结束时间"),
        field.String("file").Comment("合同文件"),
    }
}

// Edges of the EnterpriseContract.
func (EnterpriseContract) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("enterprise", Enterprise.Type).
            Ref("contracts").
            Unique().
            Required().
            Field("enterprise_id"),
    }
}

func (EnterpriseContract) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (EnterpriseContract) Indexes() []ent.Index {
    return []ent.Index{}
}
