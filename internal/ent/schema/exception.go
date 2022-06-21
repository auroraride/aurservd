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

// Exception holds the schema definition for the Exception entity.
type Exception struct {
    ent.Schema
}

// Annotations of the Exception.
func (Exception) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "exception"},
    }
}

// Fields of the Exception.
func (Exception) Fields() []ent.Field {
    return []ent.Field{
        field.Uint8("status").Default(model.ExceptionStatusPending).Comment("异常状态"),
        field.Uint64("store_id").Comment("门店ID"),
        field.String("name").Comment("物资名称"),
        field.String("model").Optional().Nillable().Comment("电池型号"),
        field.Int("num").Immutable().Comment("异常数量"),
        field.String("reason").Comment("异常原因"),
        field.String("description").Optional().Comment("异常描述"),
        field.JSON("attachments", []string{}).Optional().Comment("附件"),
    }
}

// Edges of the Exception.
func (Exception) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("store", Store.Type).Required().Unique().Ref("exceptions").Field("store_id"),
    }
}

func (Exception) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        CityMixin{},
        EmployeeMixin{},
    }
}

func (Exception) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("name"),
        index.Fields("model"),
        index.Fields("num"),
    }
}
