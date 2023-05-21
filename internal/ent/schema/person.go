package schema

import (
	"context"

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
		entsql.WithComments(true),
	}
}

// Fields of the Person.
func (Person) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("status").Default(0).Comment("认证状态"),
		field.Bool("banned").Default(false).Comment("是否封禁身份"),
		field.String("name").MaxLen(40).Comment("真实姓名"),
		field.String("id_card_number").Optional().MaxLen(40).Unique().Comment("证件号码"),
		field.Uint8("id_card_type").Default(1).Comment("证件类别"),
		field.String("id_card_portrait").Optional().Comment("证件人像面"),
		field.String("id_card_national").Optional().Comment("证件国徽面"),
		field.String("auth_face").Optional().Comment("实名认证人脸照片"),
		field.JSON("auth_result", &model.FaceVerifyResult{}).Optional().Comment("实名认证结果详情"),
		field.Time("auth_at").Nillable().Optional().Comment("实名认证结果获取时间"),
		field.String("esign_account_id").Optional().Comment("E签宝账户ID"),
		field.String("baidu_verify_token").Optional().Comment("百度人脸verify_token"),
		field.String("baidu_log_id").Optional().Comment("百度人脸log_id"),
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
		internal.Modifier{},
	}
}

func (Person) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}

func (Person) Hooks() []ent.Hook {
	type person interface {
		Name() (r string, exists bool)
		IDCardNumber() (r string, exists bool)
		ClearEsignAccountID()
	}
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if m.Op().Is(ent.OpUpdate | ent.OpUpdateOne) {
					if p, ok := m.(person); ok {
						do := false
						if _, exists := p.Name(); exists {
							do = true
						}
						if _, exists := p.IDCardNumber(); exists {
							do = true
						}
						if do {
							p.ClearEsignAccountID()
						}
					}
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}
