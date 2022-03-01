package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// BranchContract holds the schema definition for the BranchContract entity.
type BranchContract struct {
    ent.Schema
}

// Annotations of the BranchContract.
func (BranchContract) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "branch_contract"},
    }
}

// Fields of the BranchContract.
func (BranchContract) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("branch_id").Comment("网点ID"),
        field.String("landlord_name").Comment("房东姓名"),
        field.String("id_card_number").Comment("房东身份证"),
        field.String("phone").Comment("房东手机号"),
        field.String("bank_number").Comment("房东银行卡号"),
        field.Float("pledge").Comment("押金"),
        field.Float("rent").Comment("租金"),
        field.Uint("lease").Comment("租期月数"),
        field.Float("electricity_pledge").Comment("电费押金"),
        field.Float("electricity").Comment("电费单价"),
        field.Float("area").Comment("网点面积"),
        field.String("start_time").Comment("租期开始时间"),
        field.String("end_time").Comment("租期结束时间"),
        field.String("file").Comment("合同文件"),
        field.Strings("sheets").Comment("底单"),
    }
}

// Edges of the BranchContract.
func (BranchContract) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("branch", Branch.Type).
            Ref("contracts").
            Unique().
            Required().
            Field("branch_id"),
    }
}

func (BranchContract) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Creator{},
        internal.LastModifier{},
    }
}

func (BranchContract) Indexes() []ent.Index {
    return nil
}
