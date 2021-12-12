package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type RiderContact struct {
    Name  string
    Phone string
}

// Rider holds the schema definition for the Rider entity.
type Rider struct {
    ent.Schema
}

// Annotations of the Rider.
func (Rider) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "rider"},
    }
}

// Fields of the Rider.
func (Rider) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("person_id").Optional().Nillable().Comment("实人"),
        // field.Uint8("status").Default(0).Comment("业务状态"),
        field.String("phone").MaxLen(11).Unique().Comment("手机号"),
        field.JSON("contact", &RiderContact{}).Optional().Comment("紧急联系人"),
        field.Uint8("device_type").Comment("登录设备类型: 1iOS 2Android"),
        field.String("device_sn").MaxLen(60).Unique().Comment("登录设备ID"),
        field.String("device_push_id").MaxLen(60).Unique().Optional().Nillable().Comment("登录设备推送ID"),
    }
}

// Edges of the Rider.
func (Rider) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("person", Person.Type).Ref("rider").Unique().Field("person_id"),
    }
}

func (Rider) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.LastModify{},
    }
}

func (Rider) Indexes() []ent.Index {
    return nil
}
