// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/purchaseorder"
)

// PurchaseOrderDelete is the builder for deleting a PurchaseOrder entity.
type PurchaseOrderDelete struct {
	config
	hooks    []Hook
	mutation *PurchaseOrderMutation
}

// Where appends a list predicates to the PurchaseOrderDelete builder.
func (pod *PurchaseOrderDelete) Where(ps ...predicate.PurchaseOrder) *PurchaseOrderDelete {
	pod.mutation.Where(ps...)
	return pod
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pod *PurchaseOrderDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, pod.sqlExec, pod.mutation, pod.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pod *PurchaseOrderDelete) ExecX(ctx context.Context) int {
	n, err := pod.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pod *PurchaseOrderDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(purchaseorder.Table, sqlgraph.NewFieldSpec(purchaseorder.FieldID, field.TypeUint64))
	if ps := pod.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pod.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pod.mutation.done = true
	return affected, err
}

// PurchaseOrderDeleteOne is the builder for deleting a single PurchaseOrder entity.
type PurchaseOrderDeleteOne struct {
	pod *PurchaseOrderDelete
}

// Where appends a list predicates to the PurchaseOrderDelete builder.
func (podo *PurchaseOrderDeleteOne) Where(ps ...predicate.PurchaseOrder) *PurchaseOrderDeleteOne {
	podo.pod.mutation.Where(ps...)
	return podo
}

// Exec executes the deletion query.
func (podo *PurchaseOrderDeleteOne) Exec(ctx context.Context) error {
	n, err := podo.pod.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{purchaseorder.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (podo *PurchaseOrderDeleteOne) ExecX(ctx context.Context) {
	if err := podo.Exec(ctx); err != nil {
		panic(err)
	}
}
