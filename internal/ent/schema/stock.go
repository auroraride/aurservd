package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Stock holds the schema definition for the Stock entity.
type Stock struct {
    ent.Schema
}

// Annotations of the Stock.
func (Stock) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "stock"},
    }
}

// Fields of the Stock.
func (Stock) Fields() []ent.Field {
    return []ent.Field{
        field.String("sn").Comment("调拨编号"),
        field.Uint64("inbound_store_id").Optional().Nillable().Comment("所属门店ID"),
        field.Uint64("outbound_store_id").Optional().Nillable().Comment("调入自门店ID"),
        field.String("name").Comment("物资名称"),
        field.Float("voltage").Optional().Nillable().Comment("电池型号(电压)"),
        field.Int("num").Comment("物资数量: 正值调入 / 负值调出"),
    }
}

// Edges of the Stock.
func (Stock) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("inbound_store", Store.Type).Unique().Ref("inboundStocks").Field("inbound_store_id"),
        edge.From("outbound_store", Store.Type).Unique().Ref("outboundStocks").Field("outbound_store_id"),
    }
}

func (Stock) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{
            IndexCreator: true,
        },
    }
}

func (Stock) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("name"),
        index.Fields("voltage"),
        index.Fields("sn"),
    }
}
