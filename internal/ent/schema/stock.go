package schema

import (
    "context"
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/app/model"
    pEnt "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/hook"
    "github.com/auroraride/aurservd/internal/ent/internal"
    "github.com/auroraride/aurservd/pkg/snag"
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
        field.Uint8("type").Default(0).Comment("类型 0:调拨 1:领取电池 2:寄存电池 3:归还电池"),
        field.Uint64("store_id").Optional().Nillable().Comment("入库至 或 出库自 门店ID"),
        field.Uint64("rider_id").Optional().Nillable().Comment("对应骑手ID"),
        field.Uint64("employee_id").Optional().Nillable().Comment("操作店员ID"),
        field.String("name").Comment("物资名称"),
        field.Float("voltage").Optional().Nillable().Comment("电池型号(电压)"),
        field.Int("num").Immutable().Comment("物资数量: 正值调入 / 负值调出"),
    }
}

// Edges of the Stock.
func (Stock) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("store", Store.Type).Unique().Ref("stocks").Field("store_id"),
        edge.From("rider", Rider.Type).Unique().Ref("stocks").Field("rider_id"),
        edge.From("employee", Employee.Type).Unique().Ref("stocks").Field("employee_id"),
    }
}

func (Stock) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{
            IndexCreator: true,
        },
        ManagerMixin{Optional: true},
    }
}

func (Stock) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("name"),
        index.Fields("voltage"),
        index.Fields("sn"),
    }
}

func (Stock) Hooks() []ent.Hook {
    return []ent.Hook{
        hook.On(func(next ent.Mutator) ent.Mutator {
            return hook.StockFunc(func(ctx context.Context, m *pEnt.StockMutation) (ent.Value, error) {
                n, _ := m.Num()
                _, sok := m.StoreID()
                if t, ok := m.GetType(); ok {
                    if !sok || n != model.StockNumberOfRiderBusiness(t) {
                        snag.Panic("电池数量错误")
                    }
                }
                return next.Mutate(ctx, m)
            })
        }, ent.OpCreate),
    }
}
