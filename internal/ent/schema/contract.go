package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
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
        field.Uint8("status").Default(0).Comment("状态"), // 0-草稿 1-签署中 2-完成 3-撤销 5-过期（签署截至日期到期后触发） 7-拒签
        field.Uint64("rider_id").Comment("骑手"),
        field.String("flow_id").MaxLen(40).Unique().NotEmpty().Comment("E签宝流程ID"),
        field.String("sn").MaxLen(20).Unique().NotEmpty().Comment("合同编码"),
        field.JSON("files", []string{}).Optional().Comment("合同链接"),
        // TODO 其他诸如电池型号、ID、团签ID之类的信息字段
    }
}

// Edges of the Contract.
func (Contract) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("contract").Required().Unique().Field("rider_id"),
    }
}

func (Contract) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.LastModify{},
    }
}

func (Contract) Indexes() []ent.Index {
    return nil
}
