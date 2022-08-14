package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Reserve holds the schema definition for the Reserve entity.
type Reserve struct {
    ent.Schema
}

// Annotations of the Reserve.
func (Reserve) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "reserve"},
    }
}

// Fields of the Reserve.
func (Reserve) Fields() []ent.Field {
    return []ent.Field{
        field.Uint8("status").Default(model.ReserveStatusPending.Value()).Comment("预约状态"),
        field.String("type").Comment("业务类型"),
    }
}

// Edges of the Reserve.
func (Reserve) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (Reserve) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
        CabinetMixin{},
        RiderMixin{},
        CityMixin{},
        BusinessMixin{Optional: true},
    }
}

func (Reserve) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("type"),
    }
}
