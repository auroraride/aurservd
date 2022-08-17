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

type EnterpriseMixin struct {
    mixin.Schema
    Optional bool
}

func (m EnterpriseMixin) Fields() []ent.Field {
    f := field.Uint64("enterprise_id").Comment("企业ID")
    if m.Optional {
        f.Optional().Nillable()
    }
    return []ent.Field{f}
}

func (m EnterpriseMixin) Edges() []ent.Edge {
    e := edge.To("enterprise", Enterprise.Type).Unique().Field("enterprise_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

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
        field.String("name").Comment("团签名称"),
        field.String("company_name").Optional().Comment("企业全称"),
        field.Uint8("status").Comment("合作状态 0未合作 1合作中 2暂停合作"),
        field.String("contact_name").Comment("联系人姓名"),
        field.String("contact_phone").Comment("联系人电话"),
        field.String("idcard_number").Comment("身份证号码"),
        field.String("address").Comment("地址"),
        field.Uint8("payment").Comment("付费方式 1预付费 2后付费"),
        field.Float("deposit").Default(0).Comment("押金"),
        field.Float("balance").Default(0).Comment("账户余额"),
        field.Float("prepayment_total").Default(0).Comment("总储值金额 = 总金额 - 轧账金额(修改价格后自动轧账)"),
        field.Time("suspensed_at").Nillable().Optional().Comment("暂停合作时间"),
    }
}

// Edges of the Enterprise.
func (Enterprise) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("riders", Rider.Type),
        edge.To("contracts", EnterpriseContract.Type),
        edge.To("prices", EnterprisePrice.Type),
        edge.To("subscribes", Subscribe.Type),
        edge.To("statements", EnterpriseStatement.Type),
        edge.To("stations", EnterpriseStation.Type),
        edge.To("bills", EnterpriseBill.Type),
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
        index.Fields("payment"),
        index.Fields("balance"),
        index.Fields("name").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
        index.Fields("contact_name").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
        index.Fields("contact_phone").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
        index.Fields("idcard_number").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
