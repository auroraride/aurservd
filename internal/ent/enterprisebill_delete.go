// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/enterprisebill"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// EnterpriseBillDelete is the builder for deleting a EnterpriseBill entity.
type EnterpriseBillDelete struct {
	config
	hooks    []Hook
	mutation *EnterpriseBillMutation
}

// Where appends a list predicates to the EnterpriseBillDelete builder.
func (ebd *EnterpriseBillDelete) Where(ps ...predicate.EnterpriseBill) *EnterpriseBillDelete {
	ebd.mutation.Where(ps...)
	return ebd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ebd *EnterpriseBillDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, ebd.sqlExec, ebd.mutation, ebd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ebd *EnterpriseBillDelete) ExecX(ctx context.Context) int {
	n, err := ebd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ebd *EnterpriseBillDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(enterprisebill.Table, sqlgraph.NewFieldSpec(enterprisebill.FieldID, field.TypeUint64))
	if ps := ebd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ebd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ebd.mutation.done = true
	return affected, err
}

// EnterpriseBillDeleteOne is the builder for deleting a single EnterpriseBill entity.
type EnterpriseBillDeleteOne struct {
	ebd *EnterpriseBillDelete
}

// Where appends a list predicates to the EnterpriseBillDelete builder.
func (ebdo *EnterpriseBillDeleteOne) Where(ps ...predicate.EnterpriseBill) *EnterpriseBillDeleteOne {
	ebdo.ebd.mutation.Where(ps...)
	return ebdo
}

// Exec executes the deletion query.
func (ebdo *EnterpriseBillDeleteOne) Exec(ctx context.Context) error {
	n, err := ebdo.ebd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{enterprisebill.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ebdo *EnterpriseBillDeleteOne) ExecX(ctx context.Context) {
	if err := ebdo.Exec(ctx); err != nil {
		panic(err)
	}
}
