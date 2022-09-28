package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Attendance holds the schema definition for the Attendance entity.
type Attendance struct {
    ent.Schema
}

// Annotations of the Attendance.
func (Attendance) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "attendance"},
    }
}

// Fields of the Attendance.
func (Attendance) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("store_id").Comment("门店ID"),
        field.Uint64("employee_id").Comment("店员ID"),
        field.JSON("inventory", []model.AttendanceInventory{}).Optional().Comment("物资盘点"),
        field.String("photo").Optional().Nillable().Comment("上班照片"),
        field.Bool("duty").Comment("是否上班盘点"),
        field.Time("date").Comment("日期"),
        field.Float("lng").Optional().Nillable().Comment("经度"),
        field.Float("lat").Optional().Nillable().Comment("纬度"),
        field.String("address").Optional().Nillable().Comment("详细地址"),
        field.Float("distance").Optional().Nillable().Comment("打卡距离"),
    }
}

// Edges of the Attendance.
func (Attendance) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("store", Store.Type).Ref("attendances").Required().Unique().Field("store_id"),
        edge.From("employee", Employee.Type).Ref("attendances").Required().Unique().Field("employee_id"),
    }
}

func (Attendance) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Attendance) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("date", "duty"),
        index.Fields("employee_id"),
        index.Fields("store_id"),
    }
}
