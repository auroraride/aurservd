// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/planintroduce"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// PlanIntroduceDelete is the builder for deleting a PlanIntroduce entity.
type PlanIntroduceDelete struct {
	config
	hooks    []Hook
	mutation *PlanIntroduceMutation
}

// Where appends a list predicates to the PlanIntroduceDelete builder.
func (pid *PlanIntroduceDelete) Where(ps ...predicate.PlanIntroduce) *PlanIntroduceDelete {
	pid.mutation.Where(ps...)
	return pid
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pid *PlanIntroduceDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, pid.sqlExec, pid.mutation, pid.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pid *PlanIntroduceDelete) ExecX(ctx context.Context) int {
	n, err := pid.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pid *PlanIntroduceDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(planintroduce.Table, sqlgraph.NewFieldSpec(planintroduce.FieldID, field.TypeUint64))
	if ps := pid.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pid.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pid.mutation.done = true
	return affected, err
}

// PlanIntroduceDeleteOne is the builder for deleting a single PlanIntroduce entity.
type PlanIntroduceDeleteOne struct {
	pid *PlanIntroduceDelete
}

// Where appends a list predicates to the PlanIntroduceDelete builder.
func (pido *PlanIntroduceDeleteOne) Where(ps ...predicate.PlanIntroduce) *PlanIntroduceDeleteOne {
	pido.pid.mutation.Where(ps...)
	return pido
}

// Exec executes the deletion query.
func (pido *PlanIntroduceDeleteOne) Exec(ctx context.Context) error {
	n, err := pido.pid.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{planintroduce.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pido *PlanIntroduceDeleteOne) ExecX(ctx context.Context) {
	if err := pido.Exec(ctx); err != nil {
		panic(err)
	}
}
