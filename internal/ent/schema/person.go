package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
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
        field.String("id_card_number").MaxLen(40).Unique().Comment("证件号码"),
        field.Uint8("id_card_type").Default(1).Comment("证件类别"),
        field.String("id_card_portrait").Comment("证件人像面"),
        field.String("id_card_national").Comment("证件国徽面"),
        field.String("auth_face").Comment("实名认证人脸照片"),
        field.JSON("auth_result", &model.FaceVerifyResult{}).Optional().Comment("实名认证结果详情"),
        field.Time("auth_at").Nillable().Optional().Comment("实名认证结果获取时间"),
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
        internal.LastModifier{},
    }
}

func (Person) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("status"),
        index.Fields("name").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
