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

// Person holds the schema definition for the Person entity.
type Person struct {
    ent.Schema
}

// Annotations of the Person.
func (Person) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "person"},
    }
}

// Fields of the Person.
func (Person) Fields() []ent.Field {
    return []ent.Field{
        field.Uint8("status").Default(0).Comment("认证状态"),
        field.Bool("block").Default(false).Comment("封禁"),
        field.String("name").MaxLen(40).Comment("真实姓名"),
        field.String("ic_number").MaxLen(40).Unique().Comment("证件号码"),
        field.Uint8("ic_type").Default(1).Comment("证件类别"),
        field.String("ic_portrait").Comment("证件人像面"),
        field.String("ic_national").Comment("证件国徽面"),
        field.String("face_img").Comment("人脸照片"),
        field.JSON("face_verify_result", &model.FaceVerifyResult{}).Optional().Comment("人脸识别验证结果详情"),
        field.Time("result_at").Nillable().Optional().Comment("结果获取时间"),
    }
}

// Edges of the Person.
func (Person) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("rider", Rider.Type),
    }
}

func (Person) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.LastModify{},
    }
}

func (Person) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("status"),
    }
}
