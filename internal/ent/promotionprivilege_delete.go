// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/promotionprivilege"
)

// PromotionPrivilegeDelete is the builder for deleting a PromotionPrivilege entity.
type PromotionPrivilegeDelete struct {
	config
	hooks    []Hook
	mutation *PromotionPrivilegeMutation
}

// Where appends a list predicates to the PromotionPrivilegeDelete builder.
func (ppd *PromotionPrivilegeDelete) Where(ps ...predicate.PromotionPrivilege) *PromotionPrivilegeDelete {
	ppd.mutation.Where(ps...)
	return ppd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ppd *PromotionPrivilegeDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, ppd.sqlExec, ppd.mutation, ppd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ppd *PromotionPrivilegeDelete) ExecX(ctx context.Context) int {
	n, err := ppd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ppd *PromotionPrivilegeDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(promotionprivilege.Table, sqlgraph.NewFieldSpec(promotionprivilege.FieldID, field.TypeUint64))
	if ps := ppd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ppd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ppd.mutation.done = true
	return affected, err
}

// PromotionPrivilegeDeleteOne is the builder for deleting a single PromotionPrivilege entity.
type PromotionPrivilegeDeleteOne struct {
	ppd *PromotionPrivilegeDelete
}

// Where appends a list predicates to the PromotionPrivilegeDelete builder.
func (ppdo *PromotionPrivilegeDeleteOne) Where(ps ...predicate.PromotionPrivilege) *PromotionPrivilegeDeleteOne {
	ppdo.ppd.mutation.Where(ps...)
	return ppdo
}

// Exec executes the deletion query.
func (ppdo *PromotionPrivilegeDeleteOne) Exec(ctx context.Context) error {
	n, err := ppdo.ppd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{promotionprivilege.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ppdo *PromotionPrivilegeDeleteOne) ExecX(ctx context.Context) {
	if err := ppdo.Exec(ctx); err != nil {
		panic(err)
	}
}
