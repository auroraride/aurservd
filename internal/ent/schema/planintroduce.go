package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type PlanIntroduceMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m PlanIntroduceMixin) Fields() []ent.Field {
    relate := field.Uint64("introduce_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m PlanIntroduceMixin) Edges() []ent.Edge {
    e := edge.To("introduce", PlanIntroduce.Type).Unique().Field("introduce_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m PlanIntroduceMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("introduce_id"))
    }
    return
}

// PlanIntroduce holds the schema definition for the PlanIntroduce entity.
type PlanIntroduce struct {
    ent.Schema
}

// Annotations of the PlanIntroduce.
func (PlanIntroduce) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "plan_introduce"},
    }
}

// Fields of the PlanIntroduce.
func (PlanIntroduce) Fields() []ent.Field {
    return []ent.Field{
        field.String("image"),
    }
}

// Edges of the PlanIntroduce.
func (PlanIntroduce) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (PlanIntroduce) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{DisableIndex: true},
        BatteryModelMixin{},
        EbikeBrandMixin{Optional: true},
    }
}

func (PlanIntroduce) Indexes() []ent.Index {
    return []ent.Index{}
}
