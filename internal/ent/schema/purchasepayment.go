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
		field.Enum("payway").Values("alipay", "wechat", "cash").Comment("支付方式，alipay: 支付宝, wechat: 微信, cash: 现金"),
		field.Int("index").Comment("支付序号"),
		field.Float("total").Comment("支付金额"),
		field.Float("amount").Comment("分期金额"),
		field.Float("forfeit").Default(0).Comment("滞纳金"),
		field.Time("paid_date").SchemaType(map[string]string{dialect.Postgres: "date"}).Comment("支付日期"),
		field.String("trade_no").Optional().Comment("平台订单号（微信或支付宝）"),
	}
}

// Edges of the PurchasePayment.
func (PurchasePayment) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PurchasePayment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		RiderMixin{},
		PurchaseCommodityMixin{},
		PurchaseOrderMixin{},
	}
}

func (PurchasePayment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("payway"),
		index.Fields("index"),
		index.Fields("trade_no"),
	}
}
