// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/promotiongrowth"
)

// PromotionGrowthDelete is the builder for deleting a PromotionGrowth entity.
type PromotionGrowthDelete struct {
	config
	hooks    []Hook
	mutation *PromotionGrowthMutation
}

// Where appends a list predicates to the PromotionGrowthDelete builder.
func (pgd *PromotionGrowthDelete) Where(ps ...predicate.PromotionGrowth) *PromotionGrowthDelete {
	pgd.mutation.Where(ps...)
	return pgd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pgd *PromotionGrowthDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, pgd.sqlExec, pgd.mutation, pgd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pgd *PromotionGrowthDelete) ExecX(ctx context.Context) int {
	n, err := pgd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pgd *PromotionGrowthDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(promotiongrowth.Table, sqlgraph.NewFieldSpec(promotiongrowth.FieldID, field.TypeUint64))
	if ps := pgd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pgd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pgd.mutation.done = true
	return affected, err
}

// PromotionGrowthDeleteOne is the builder for deleting a single PromotionGrowth entity.
type PromotionGrowthDeleteOne struct {
	pgd *PromotionGrowthDelete
}

// Where appends a list predicates to the PromotionGrowthDelete builder.
func (pgdo *PromotionGrowthDeleteOne) Where(ps ...predicate.PromotionGrowth) *PromotionGrowthDeleteOne {
	pgdo.pgd.mutation.Where(ps...)
	return pgdo
}

// Exec executes the deletion query.
func (pgdo *PromotionGrowthDeleteOne) Exec(ctx context.Context) error {
	n, err := pgdo.pgd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{promotiongrowth.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pgdo *PromotionGrowthDeleteOne) ExecX(ctx context.Context) {
	if err := pgdo.Exec(ctx); err != nil {
		panic(err)
	}
}
