package schema

import (
	"context"
	"errors"
	"slices"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type GoodsMixin struct {
	mixin.Schema
	DisableIndex bool
	Optional     bool
}

func (m GoodsMixin) Fields() []ent.Field {
	f := field.Uint64("goods_id").Comment("商品ID")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m GoodsMixin) Edges() []ent.Edge {
	e := edge.To("goods", Goods.Type).Unique().Field("goods_id").Comment("商品ID")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m GoodsMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("goods_id"))
	}
	return
}

// Goods holds the schema definition for the Goods entity.
type Goods struct {
	ent.Schema
}

// Annotations of the Goods.
func (Goods) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "goods"},
		entsql.WithComments(true),
	}
}

// Fields of the Goods.
func (Goods) Fields() []ent.Field {
	return []ent.Field{
		field.String("sn").Comment("商品编号"),
		field.String("name").Comment("商品名称"),
		field.Uint8("type").Default(1).Comment("商品类别 1:电车"),
		field.Strings("lables").Optional().Comment("商品标签"),
		field.Float("price").Comment("商品价格"),
		field.Int("weight").Comment("商品权重"),
		field.String("head_pic").Comment("列表头图"),
		field.Strings("photos").Comment("商品图片"),
		field.Strings("intro").Comment("商品介绍"),
		field.Uint8("status").Default(0).Comment("商品状态 0下架 1上架"),
		field.JSON("installment", [][]float64{}).Optional().Comment("分期方案"),
	}
}

// Edges of the Goods.
func (Goods) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("stores", StoreGoods.Type),
	}
}

func (Goods) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Goods) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sn").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
		index.Fields("status"),
	}
}

// SortInstallment 对分期方案进行排序
func SortInstallment(items [][]float64) error {
	slices.SortFunc(items, func(i, j []float64) int {
		if len(i) < len(j) {
			return -1
		}
		return 1
	})
	n := 0
	for _, s := range items {
		l := len(s)
		if n == l {
			return errors.New("分期方案重复")
		}
		n = l
	}
	return nil
}

type GoodsMutation interface {
	Installment() ([][]float64, bool)
	SetInstallment(f [][]float64)
}

func (Goods) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if m.Op().Is(ent.OpCreate | ent.OpUpdateOne | ent.OpUpdate) {
					switch r := m.(type) {
					case GoodsMutation:
						installment, ok := r.Installment()
						if ok {
							if err := SortInstallment(installment); err != nil {
								return nil, err
							}
							r.SetInstallment(installment)
						}
					}
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}
