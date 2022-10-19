package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/app/model"
)

type EbikeAllocateMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m EbikeAllocateMixin) Fields() []ent.Field {
    relate := field.Uint64("allocate_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m EbikeAllocateMixin) Edges() []ent.Edge {
    e := edge.To("allocate", EbikeAllocate.Type).Unique().Field("allocate_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m EbikeAllocateMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("allocate_id"))
    }
    return
}

// EbikeAllocate holds the schema definition for the EbikeAllocate entity.
type EbikeAllocate struct {
    ent.Schema
}

// Annotations of the EbikeAllocate.
func (EbikeAllocate) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "ebike_allocate"},
    }
}

// Fields of the EbikeAllocate.
func (EbikeAllocate) Fields() []ent.Field {
    return []ent.Field{
        field.Uint8("status").Comment("分配状态"),
        field.JSON("info", &model.EbikeAllocate{}).Comment("电车信息"),
        field.Time("time").Comment("分配时间"),
    }
}

// Edges of the EbikeAllocate.
func (EbikeAllocate) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("contract", Contract.Type).Unique(),
    }
}

func (EbikeAllocate) Mixin() []ent.Mixin {
    return []ent.Mixin{
        EmployeeMixin{Optional: true},
        StoreMixin{},
        EbikeMixin{},
        EbikeBrandMixin{},
        SubscribeMixin{Unique: true},
        RiderMixin{},
    }
}

func (EbikeAllocate) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("time"),
    }
}
