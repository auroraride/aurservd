package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Setting holds the schema definition for the Setting entity.
type Setting struct {
    ent.Schema
}

// Annotations of the Setting.
func (Setting) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "setting"},
    }
}

// Fields of the Setting.
func (Setting) Fields() []ent.Field {
    return []ent.Field{
        field.String("key").MaxLen(40).Unique().Comment("设置名"),
        field.JSON("val", model.Setting{}).Comment("设置值"),
    }
}

// Edges of the Setting.
func (Setting) Edges() []ent.Edge {
    return nil
}

func (Setting) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
    }
}

func (Setting) Indexes() []ent.Index {
    return nil
}
