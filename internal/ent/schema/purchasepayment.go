package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type PurchasePaymentMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PurchasePaymentMixin) Fields() []ent.Field {
	relate := field.Uint64("payment_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PurchasePaymentMixin) Edges() []ent.Edge {
	e := edge.To("payment", PurchasePayment.Type).Unique().Field("payment_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PurchasePaymentMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("payment_id"))
	}
	return
}

// PurchasePayment holds the schema definition for the PurchasePayment entity.
type PurchasePayment struct {
	ent.Schema
}

// Annotations of the PurchasePayment.
func (PurchasePayment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "purchase_payment"},
		entsql.WithComments(true),
	}
}

// Fields of the PurchasePayment.
func (PurchasePayment) Fields() []ent.Field {
	return []ent.Field{
		field.String("out_trade_no").Comment("交易订单号（自行生成）"),
		field.Int("index").Comment("分期序号"),
		field.Enum("status").Values("obligation", "paid", "canceled").Default("obligation").Comment("支付状态，obligation: 待付款, paid: 已支付, canceled: 已取消"),
		field.Enum("payway").Optional().Values("alipay", "wechat", "cash").Comment("支付方式，alipay: 支付宝, wechat: 微信, cash: 现金"),
		field.Float("total").Comment("支付金额"),
		field.Float("amount").Comment("账单金额"),
		field.Float("forfeit").Default(0).Comment("滞纳金"),
		field.Time("billing_date").Comment("账单日期"),
		field.Time("payment_time").Optional().Nillable().Comment("支付时间"),
		field.String("trade_no").Optional().Comment("平台订单号（微信或支付宝）"),
		field.Uint64("order_id").Optional().Comment("订单id"),
	}
}

// Edges of the PurchasePayment.
func (PurchasePayment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", PurchaseOrder.Type).Ref("payments").Unique().Field("order_id"),
	}
}

func (PurchasePayment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		RiderMixin{},
		GoodsMixin{},
	}
}

func (PurchasePayment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("payway"),
		index.Fields("index"),
		index.Fields("trade_no"),
	}
}
