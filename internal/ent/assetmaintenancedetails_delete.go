// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/assetmaintenancedetails"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// AssetMaintenanceDetailsDelete is the builder for deleting a AssetMaintenanceDetails entity.
type AssetMaintenanceDetailsDelete struct {
	config
	hooks    []Hook
	mutation *AssetMaintenanceDetailsMutation
}

// Where appends a list predicates to the AssetMaintenanceDetailsDelete builder.
func (amdd *AssetMaintenanceDetailsDelete) Where(ps ...predicate.AssetMaintenanceDetails) *AssetMaintenanceDetailsDelete {
	amdd.mutation.Where(ps...)
	return amdd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (amdd *AssetMaintenanceDetailsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, amdd.sqlExec, amdd.mutation, amdd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (amdd *AssetMaintenanceDetailsDelete) ExecX(ctx context.Context) int {
	n, err := amdd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (amdd *AssetMaintenanceDetailsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(assetmaintenancedetails.Table, sqlgraph.NewFieldSpec(assetmaintenancedetails.FieldID, field.TypeUint64))
	if ps := amdd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, amdd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	amdd.mutation.done = true
	return affected, err
}

// AssetMaintenanceDetailsDeleteOne is the builder for deleting a single AssetMaintenanceDetails entity.
type AssetMaintenanceDetailsDeleteOne struct {
	amdd *AssetMaintenanceDetailsDelete
}

// Where appends a list predicates to the AssetMaintenanceDetailsDelete builder.
func (amddo *AssetMaintenanceDetailsDeleteOne) Where(ps ...predicate.AssetMaintenanceDetails) *AssetMaintenanceDetailsDeleteOne {
	amddo.amdd.mutation.Where(ps...)
	return amddo
}

// Exec executes the deletion query.
func (amddo *AssetMaintenanceDetailsDeleteOne) Exec(ctx context.Context) error {
	n, err := amddo.amdd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{assetmaintenancedetails.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (amddo *AssetMaintenanceDetailsDeleteOne) ExecX(ctx context.Context) {
	if err := amddo.Exec(ctx); err != nil {
		panic(err)
	}
}