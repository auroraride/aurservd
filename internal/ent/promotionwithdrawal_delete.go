// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/promotionwithdrawal"
)

// PromotionWithdrawalDelete is the builder for deleting a PromotionWithdrawal entity.
type PromotionWithdrawalDelete struct {
	config
	hooks    []Hook
	mutation *PromotionWithdrawalMutation
}

// Where appends a list predicates to the PromotionWithdrawalDelete builder.
func (pwd *PromotionWithdrawalDelete) Where(ps ...predicate.PromotionWithdrawal) *PromotionWithdrawalDelete {
	pwd.mutation.Where(ps...)
	return pwd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pwd *PromotionWithdrawalDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, pwd.sqlExec, pwd.mutation, pwd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pwd *PromotionWithdrawalDelete) ExecX(ctx context.Context) int {
	n, err := pwd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pwd *PromotionWithdrawalDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(promotionwithdrawal.Table, sqlgraph.NewFieldSpec(promotionwithdrawal.FieldID, field.TypeUint64))
	if ps := pwd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pwd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pwd.mutation.done = true
	return affected, err
}

// PromotionWithdrawalDeleteOne is the builder for deleting a single PromotionWithdrawal entity.
type PromotionWithdrawalDeleteOne struct {
	pwd *PromotionWithdrawalDelete
}

// Where appends a list predicates to the PromotionWithdrawalDelete builder.
func (pwdo *PromotionWithdrawalDeleteOne) Where(ps ...predicate.PromotionWithdrawal) *PromotionWithdrawalDeleteOne {
	pwdo.pwd.mutation.Where(ps...)
	return pwdo
}

// Exec executes the deletion query.
func (pwdo *PromotionWithdrawalDeleteOne) Exec(ctx context.Context) error {
	n, err := pwdo.pwd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{promotionwithdrawal.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pwdo *PromotionWithdrawalDeleteOne) ExecX(ctx context.Context) {
	if err := pwdo.Exec(ctx); err != nil {
		panic(err)
	}
}