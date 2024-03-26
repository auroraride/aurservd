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

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/internal"
)

type PlanMixin struct {
	mixin.Schema
	DisableIndex bool
	Optional     bool
	Comment      string
}

func (m PlanMixin) Fields() []ent.Field {
	f := field.Uint64("plan_id").Comment(m.Comment)
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m PlanMixin) Edges() []ent.Edge {
	e := edge.To("plan", Plan.Type).Unique().Field("plan_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PlanMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("plan_id"))
	}
	return
}

// Plan holds the schema definition for the Plan entity.
type Plan struct {
	ent.Schema
}

// Annotations of the Plan.
func (Plan) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "plan"},
		entsql.WithComments(true),
	}
}

// Fields of the Plan.
func (Plan) Fields() []ent.Field {
	return []ent.Field{
		field.String("model").Optional().Comment("电池型号"),
		field.Bool("enable").Comment("是否启用"),
		field.Uint8("type").Default(model.PlanTypeBattery.Value()).Comment("骑士卡类别 1:单电 2:车加电"),
		field.String("name").Comment("骑士卡名称"),
		field.Time("start").Comment("有效期开始日期"),
		field.Time("end").Comment("有效期结束日期"),
		field.Float("price").Comment("骑士卡价格"),
		field.Uint("days").Comment("骑士卡天数"),
		field.Float("commission").Comment("提成"),
		field.Float("original").Optional().Comment("原价"),
		field.String("desc").Optional().Comment("优惠信息"),
		field.Uint64("parent_id").Optional().Nillable().Comment("父级"),
		field.Float("discount_newly").Default(0).Comment("新签减免"),
		field.Strings("notes").Optional().Comment("购买须知"),
		field.Bool("intelligent").Default(false).Comment("是否智能柜套餐"),
		field.Bool("deposit").Default(false).Comment("是否开启押金(只对V2版本接口有用)"),
		field.Float("deposit_amount").Optional().Comment("押金金额"),
		field.Bool("deposit_wechat_payscore").Optional().Default(false).Comment("微信支付分免押金"),
		field.Bool("deposit_alipay_auth_freeze").Optional().Default(false).Comment("预授权信用免押金"),
		field.Bool("deposit_contract").Optional().Default(false).Comment("合同免押金"),
		field.Bool("deposit_pay").Optional().Default(false).Comment("支付押金"),
	}
}

// Edges of the Plan.
func (Plan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cities", City.Type),
		edge.To("complexes", Plan.Type).From("parent").Field("parent_id").Unique(),

		edge.To("commissions", PromotionCommissionPlan.Type),
	}
}

func (Plan) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
		EbikeBrandMixin{Optional: true},
		AgreementMixin{Optional: true},
	}
}

func (Plan) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("model"),
		index.Fields("days"),
		index.Fields("enable"),
		index.Fields("start", "end"),
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}
