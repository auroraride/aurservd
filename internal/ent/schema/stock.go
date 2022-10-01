package schema

import (
    "context"
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/entc/integration/ent/hook"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/app/model"
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
        field.Uint8("type").Default(0).Comment("类型 0:调拨 1:领取电池 2:寄存电池 3:结束寄存 4:归还电池"),
        field.Uint64("store_id").Optional().Nillable().Comment("入库至 或 出库自 门店ID"),
        field.Uint64("cabinet_id").Optional().Nillable().Comment("入库至 或 出库自 电柜ID"),
        field.Uint64("rider_id").Optional().Nillable().Comment("对应骑手ID"),
        field.Uint64("employee_id").Optional().Nillable().Comment("操作店员ID"),
        field.String("name").Comment("物资名称"),
        field.String("model").Optional().Nillable().Comment("电池型号"),
        field.Int("num").Immutable().Comment("物资数量: 正值调入 / 负值调出"),
        field.Enum("material").Values("battery", "frame", "others").Comment("物资种类"),
    }
}

// Edges of the Stock.
func (Stock) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("store", Store.Type).Unique().Ref("stocks").Field("store_id"),
        edge.From("cabinet", Cabinet.Type).Unique().Ref("stocks").Field("cabinet_id"),
        edge.From("rider", Rider.Type).Unique().Ref("stocks").Field("rider_id"),
        edge.From("employee", Employee.Type).Unique().Ref("stocks").Field("employee_id"),

        edge.To("spouse", Stock.Type).Unique(),
    }
}

func (Stock) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{
            IndexCreator: true,
        },
        CityMixin{Optional: true},
        SubscribeMixin{Optional: true},

        // 电车
        EbikeMixin{Optional: true},
    }
}

func (Stock) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("store_id"),
        index.Fields("cabinet_id"),
        index.Fields("rider_id"),
        index.Fields("employee_id"),
        index.Fields("name"),
        index.Fields("model"),
        index.Fields("sn"),
    }
}

func (Stock) Hooks() []ent.Hook {
    type intr interface {
        Num() (r int, exists bool)
        StoreID() (r uint64, exists bool)
        GetType() (r uint8, exists bool)
    }

    return []ent.Hook{
        hook.On(
            func(next ent.Mutator) ent.Mutator {
                return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
                    if st, ok := m.(intr); ok {
                        n, _ := st.Num()
                        if t, ok := st.GetType(); ok {
                            // 骑手业务是否满足调拨电池数量要求
                            x := model.StockNumberOfRiderBusiness(t)
                            if x != 0 && n != x {
                                snag.Panic("电池数量错误")
                            }
                        }
                    }
                    return next.Mutate(ctx, m)
                })
            },
            ent.OpCreate,
        ),
    }
}
