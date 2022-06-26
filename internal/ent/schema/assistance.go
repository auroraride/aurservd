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

// Assistance holds the schema definition for the Assistance entity.
type Assistance struct {
    ent.Schema
}

// Annotations of the Assistance.
func (Assistance) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "assistance"},
    }
}

// Fields of the Assistance.
func (Assistance) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("employee_id").Optional().Nillable().Comment("店员ID"),
        field.Uint64("order_id").Optional().Nillable().Comment("支付订单"),
        field.Uint8("status").Default(0).Comment("救援状态 0:待分配 1:已接单/已分配 2:已拒绝 3:救援失败 4:救援成功待支付 5:救援成功已支付"),
        field.Float("lng").Comment("经度"),
        field.Float("lat").Comment("纬度"),
        field.String("address").Comment("详细地址"),
        field.String("breakdown").Comment("故障"),
        field.String("breakdown_desc").Optional().Comment("故障描述"),
        field.Strings("breakdown_photos").Comment("故障照片"),
        field.String("cancel_reason").Optional().Nillable().Comment("取消原因"),
        field.String("cancel_reason_desc").Optional().Nillable().Comment("取消原因详细描述"),
        field.Float("distance").Optional().Comment("救援距离"),
        field.String("reason").Optional().Comment("救援原因"),
        field.String("detect_photo").Optional().Comment("检测照片"),
        field.String("joint_photo").Optional().Comment("与用户合影"),
        field.Float("cost").Default(0).Optional().Comment("本次救援费用"),
        field.String("refused_desc").Optional().Nillable().Comment("拒绝原因"),
        field.Time("pay_at").Optional().Nillable().Comment("支付时间"),
        field.Time("allocate_at").Optional().Nillable().Comment("分配时间"),
        field.Int("wait").Default(0).Comment("分配等待时间(s)"),
        field.String("free_reason").Optional().Nillable().Comment("免费理由"),
        field.String("fail_reason").Optional().Nillable().Comment("失败原因"),
        field.Time("process_at").Optional().Nillable().Comment("救援处理时间"),
        field.Float("price").Optional().Comment("救援费用单价 元/公里"),
        field.Int("navi_duration").Optional().Comment("路径导航规划时间 (s)"),
        field.Strings("navi_polylines").Optional().Comment("路径导航规划坐标组"),
    }
}

// Edges of the Assistance.
func (Assistance) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("order", Order.Type).Unique().Ref("assistance").Field("order_id"),
        edge.From("employee", Employee.Type).Unique().Ref("assistances").Field("employee_id"),
    }
}

func (Assistance) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        StoreMixin{Optional: true},
        RiderMixin{},
        SubscribeMixin{},
        CityMixin{},
    }
}

func (Assistance) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("status"),
    }
}
