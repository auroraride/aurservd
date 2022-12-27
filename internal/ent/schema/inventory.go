package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Inventory holds the schema definition for the Inventory entity.
type Inventory struct {
    ent.Schema
}

// Annotations of the Inventory.
func (Inventory) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "inventory"},
        entsql.WithComments(true),
    }
}

// Fields of the Inventory.
func (Inventory) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Unique().Comment("物资名称"),
        field.Bool("count").Comment("是否需要盘点"),
        field.Bool("transfer").Comment("是否可调拨"),
        field.Bool("purchase").Comment("是否可采购"),
    }
}

// Edges of the Inventory.
func (Inventory) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (Inventory) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Inventory) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("count"),
        index.Fields("transfer"),
        index.Fields("purchase"),
    }
}
