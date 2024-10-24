// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/promotionmembercommission"
)

// PromotionMemberCommissionDelete is the builder for deleting a PromotionMemberCommission entity.
type PromotionMemberCommissionDelete struct {
	config
	hooks    []Hook
	mutation *PromotionMemberCommissionMutation
}

// Where appends a list predicates to the PromotionMemberCommissionDelete builder.
func (pmcd *PromotionMemberCommissionDelete) Where(ps ...predicate.PromotionMemberCommission) *PromotionMemberCommissionDelete {
	pmcd.mutation.Where(ps...)
	return pmcd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pmcd *PromotionMemberCommissionDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, pmcd.sqlExec, pmcd.mutation, pmcd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pmcd *PromotionMemberCommissionDelete) ExecX(ctx context.Context) int {
	n, err := pmcd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pmcd *PromotionMemberCommissionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(promotionmembercommission.Table, sqlgraph.NewFieldSpec(promotionmembercommission.FieldID, field.TypeUint64))
	if ps := pmcd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pmcd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pmcd.mutation.done = true
	return affected, err
}

// PromotionMemberCommissionDeleteOne is the builder for deleting a single PromotionMemberCommission entity.
type PromotionMemberCommissionDeleteOne struct {
	pmcd *PromotionMemberCommissionDelete
}

// Where appends a list predicates to the PromotionMemberCommissionDelete builder.
func (pmcdo *PromotionMemberCommissionDeleteOne) Where(ps ...predicate.PromotionMemberCommission) *PromotionMemberCommissionDeleteOne {
	pmcdo.pmcd.mutation.Where(ps...)
	return pmcdo
}

// Exec executes the deletion query.
func (pmcdo *PromotionMemberCommissionDeleteOne) Exec(ctx context.Context) error {
	n, err := pmcdo.pmcd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{promotionmembercommission.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pmcdo *PromotionMemberCommissionDeleteOne) ExecX(ctx context.Context) {
	if err := pmcdo.Exec(ctx); err != nil {
		panic(err)
	}
}
