package schema

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/internal"
	"github.com/auroraride/aurservd/pkg/snag"
)

type SubscribeMixin struct {
	mixin.Schema
	DisableIndex bool
	Optional     bool
	Unique       bool
}

func (m SubscribeMixin) Fields() []ent.Field {
	f := field.Uint64("subscribe_id")
	if m.Optional {
		f.Optional().Nillable()
	}
	if m.Unique {
		f.Unique()
	}
	return []ent.Field{f}
}

func (m SubscribeMixin) Edges() []ent.Edge {
	e := edge.To("subscribe", Subscribe.Type).Unique().Field("subscribe_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m SubscribeMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		i := index.Fields("subscribe_id")
		if m.Unique {
			i.Unique()
		}
		arr = append(arr, i)
	}
	return
}

// Subscribe holds the schema definition for the Subscribe entity.
type Subscribe struct {
	ent.Schema
}

// Annotations of the Subscribe.
func (Subscribe) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{
			Table: "subscribe",
		},
	}
}

// Fields of the Subscribe.
func (Subscribe) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("rider_id").Comment("骑手ID"),
		field.Uint64("initial_order_id").Optional().Comment("初始订单ID(开通订阅的初始订单), 团签用户无此字段"),
		field.Uint64("enterprise_id").Optional().Nillable().Comment("企业ID"),
		field.Uint8("status").Default(model.SubscribeStatusInactive).Comment("当前订阅状态"),
		field.Uint("type").Default(0).Immutable().Comment("订阅类型 0团签 1新签 2续签 3重签 4更改电池, 除0值外 其他值参考order.type"),
		field.String("model").Comment("电池型号"),
		// field.Int("days").Comment("总天数 = 骑士卡天数 + 改动天数 + 暂停天数 + 续费天数 + 已缴纳逾期滞纳金天数"),
		field.Int("initial_days").Optional().Comment("初始骑士卡天数, 个签和代理模式团签有此字段"),
		field.Int("alter_days").Default(0).Comment("改动天数"),
		field.Int("pause_days").Default(0).Comment("寄存天数"),
		field.Int("suspend_days").Default(0).Comment("暂停天数"),
		field.Int("renewal_days").Default(0).Comment("续期天数"),
		field.Int("overdue_days").Default(0).Comment("已缴纳逾期滞纳金天数"),
		field.Int("remaining").Default(0).Comment("剩余天数, 负数为逾期, 代理商骑手剩余时间根据agent_end_at计算"),
		field.Time("paused_at").Optional().Nillable().Comment("当前寄存时间"),
		field.Time("suspend_at").Optional().Nillable().Comment("当前暂停时间"),
		field.Time("start_at").Optional().Nillable().Comment("激活时间"),
		field.Time("end_at").Optional().Nillable().Comment("归还/团签结束时间"),
		field.Time("refund_at").Optional().Nillable().Comment("退款时间"),
		field.String("unsubscribe_reason").Optional().Comment("退租理由"),
		field.Time("last_bill_date").Optional().Nillable().Comment("上次结算日期(包含该日期)"),
		field.Bool("pause_overdue").Default(false).Comment("是否超期退租"),
		field.Time("agent_end_at").Optional().Nillable().Comment("代理商处到期日期"),
		field.String("formula").Optional().Nillable().Comment("计算公式"),
		field.Bool("need_contract").Default(false).Comment("是否需要签约"),
		field.Bool("intelligent").Default(false).Comment("是否智能柜套餐"),
	}
}

// Edges of the Subscribe.
func (Subscribe) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("rider", Rider.Type).Ref("subscribes").Required().Unique().Field("rider_id"),
		edge.From("enterprise", Enterprise.Type).Ref("subscribes").Unique().Field("enterprise_id"),

		edge.To("pauses", SubscribePause.Type),
		edge.To("suspends", SubscribeSuspend.Type),
		edge.To("alters", SubscribeAlter.Type),
		edge.To("orders", Order.Type),

		edge.To("initial_order", Order.Type).Unique().Field("initial_order_id").Comment("对应初始订单"),

		edge.To("bills", EnterpriseBill.Type),

		edge.To("battery", Battery.Type).Unique(),
	}
}

func (Subscribe) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		PlanMixin{Optional: true},
		EmployeeMixin{Optional: true},
		CityMixin{},
		StationMixin{Optional: true},
		StoreMixin{Optional: true},
		CabinetMixin{Optional: true},

		// 电车
		EbikeBrandMixin{Optional: true},
		EbikeMixin{Optional: true},
	}
}

func (Subscribe) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("rider_id"),
		index.Fields("enterprise_id"),
		index.Fields("paused_at"),
		index.Fields("last_bill_date"),
		index.Fields("start_at", "end_at"),
	}
}

func (Subscribe) Hooks() []ent.Hook {
	type intr interface {
		UnsubscribeReason() (r string, exists bool)
		SetUnsubscribeReason(s string)
		Status() (r uint8, exists bool)
	}
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if m.Op().Is(ent.OpUpdate | ent.OpUpdateOne) {
					if sub, ok := m.(intr); ok {
						if status, ok := sub.Status(); ok && status == model.SubscribeStatusUnSubscribed {
							if reason, _ := sub.UnsubscribeReason(); reason == "" {
								snag.Panic("退租理由必填")
							}
						}
					}
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}
