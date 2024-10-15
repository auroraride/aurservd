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

type PurchaseOrderMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PurchaseOrderMixin) Fields() []ent.Field {
	relate := field.Uint64("order_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PurchaseOrderMixin) Edges() []ent.Edge {
	e := edge.To("order", PurchaseOrder.Type).Unique().Field("order_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PurchaseOrderMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("order_id"))
	}
	return
}

// PurchaseOrder holds the schema definition for the PurchaseOrder entity.
type PurchaseOrder struct {
	ent.Schema
}

// Annotations of the PurchaseOrder.
func (PurchaseOrder) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "purchase_order"},
		entsql.WithComments(true),
	}
}

// Fields of the PurchaseOrder.
func (PurchaseOrder) Fields() []ent.Field {
	return []ent.Field{
		field.String("sn").Comment("车架号"),
		field.Enum("status").
			Values("pending", "staging", "ended", "cancelled", "refunded").
			Default("pending").
			Comment("状态, pending: 待支付, staging: 分期执行中, ended: 已完成, cancelled: 已取消, refunded: 已退款"),
		field.String("contract_url").Optional().Comment("合同URL"),
		field.Int("installment_stage").Default(0).Comment("当前分期阶段，从0开始"),
		field.Int("installment_total").Comment("分期总数"),
		field.JSON("installment_plan", model.GoodsPaymentPlan{}).Comment("分期方案"),
		field.Time("start_date").Optional().SchemaType(map[string]string{dialect.Postgres: "date"}).Comment("开始日期"),
		field.Time("next_date").Nillable().Optional().SchemaType(map[string]string{dialect.Postgres: "date"}).Comment("下次支付日期"),
		field.Strings("images").Optional().Comment("图片"),
	}
}

// Edges of the PurchaseOrder.
func (PurchaseOrder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("payments", PurchasePayment.Type),
	}
}

func (PurchaseOrder) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		RiderMixin{},
		GoodsMixin{},
		StoreMixin{Optional: true},
	}
}

func (PurchaseOrder) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sn"),
		index.Fields("status"),
		index.Fields("next_date"),
	}
}
