package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Enterprise holds the schema definition for the Enterprise entity.
type Enterprise struct {
    ent.Schema
}

// Annotations of the Enterprise.
func (Enterprise) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "enterprise"},
    }
}

// Fields of the Enterprise.
func (Enterprise) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Comment("企业名称"),
        field.Uint8("status").Comment("合作状态 0未合作 1合作中 2暂停合作"),
        field.String("contact_name").Comment("联系人姓名"),
        field.String("contact_phone").Comment("联系人电话"),
        field.String("idcard_number").Comment("身份证号码"),
        field.String("address").Comment("地址"),
        field.Uint8("payment").Comment("付费方式 1预付费 2后付费"),
        field.Float("deposit").Default(0).Comment("押金"),
        field.Float("balance").Default(0).Comment("账户余额"),
        field.Float("arrearage").Default(0).Comment("欠费金额"),
    }
}

// Edges of the Enterprise.
func (Enterprise) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("riders", Rider.Type),
        edge.To("contracts", EnterpriseContract.Type),
        edge.To("prices", EnterprisePrice.Type),
        edge.To("subscribes", Subscribe.Type),
    }
}

func (Enterprise) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        CityMixin{},
    }
}

func (Enterprise) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("name").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
