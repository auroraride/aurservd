package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

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
        field.Uint64("group_id").Optional().Nillable().Comment("团队"),
        // field.Uint8("status").Default(0).Comment("业务状态"),
        field.String("phone").MaxLen(11).Unique().Comment("手机号"),
        field.JSON("contact", &model.RiderContact{}).Optional().Comment("紧急联系人"),
        field.Uint8("device_type").Comment("登录设备类型: 1iOS 2Android"),
        field.String("last_device").MaxLen(60).Unique().Comment("上次登录设备ID"),
        field.String("last_face").Optional().Nillable().Comment("上次登录人脸"),
        field.String("push_id").MaxLen(60).Unique().Optional().Nillable().Comment("推送ID"),
        field.Time("last_signin_at").Nillable().Optional().Comment("最后登录时间"),
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
