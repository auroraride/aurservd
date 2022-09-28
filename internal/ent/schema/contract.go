package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Contract holds the schema definition for the Contract entity.
type Contract struct {
    ent.Schema
}

// Annotations of the Contract.
func (Contract) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "contract"},
    }
}

// Fields of the Contract.
func (Contract) Fields() []ent.Field {
    return []ent.Field{
        field.Uint8("status").Default(0).Comment("状态"), // 0草稿 1签署中 2完成 3撤销 4过期(签署截至日期到期后触发) 5拒签
        field.Uint64("rider_id").Comment("骑手"),
        field.String("flow_id").MaxLen(64).Unique().NotEmpty().Comment("E签宝流程ID"),
        field.String("sn").MaxLen(64).Unique().NotEmpty().Comment("合同编码"),
        field.JSON("files", []string{}).Optional().Comment("合同链接"),
        field.Bool("effective").Default(true).Comment("是否有效"), // TODO 需要实现逻辑, 当用户退租之后触发合同失效, 需要重新签订?
    }
}

// Edges of the Contract.
func (Contract) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("contracts").Required().Unique().Field("rider_id"),
    }
}

func (Contract) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Contract) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("rider_id"),
        index.Fields("status", "effective"),
    }
}
