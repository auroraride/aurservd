// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/contracttemplate"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ContractTemplateDelete is the builder for deleting a ContractTemplate entity.
type ContractTemplateDelete struct {
	config
	hooks    []Hook
	mutation *ContractTemplateMutation
}

// Where appends a list predicates to the ContractTemplateDelete builder.
func (ctd *ContractTemplateDelete) Where(ps ...predicate.ContractTemplate) *ContractTemplateDelete {
	ctd.mutation.Where(ps...)
	return ctd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ctd *ContractTemplateDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, ctd.sqlExec, ctd.mutation, ctd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ctd *ContractTemplateDelete) ExecX(ctx context.Context) int {
	n, err := ctd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ctd *ContractTemplateDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(contracttemplate.Table, sqlgraph.NewFieldSpec(contracttemplate.FieldID, field.TypeUint64))
	if ps := ctd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ctd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ctd.mutation.done = true
	return affected, err
}

// ContractTemplateDeleteOne is the builder for deleting a single ContractTemplate entity.
type ContractTemplateDeleteOne struct {
	ctd *ContractTemplateDelete
}

// Where appends a list predicates to the ContractTemplateDelete builder.
func (ctdo *ContractTemplateDeleteOne) Where(ps ...predicate.ContractTemplate) *ContractTemplateDeleteOne {
	ctdo.ctd.mutation.Where(ps...)
	return ctdo
}

// Exec executes the deletion query.
func (ctdo *ContractTemplateDeleteOne) Exec(ctx context.Context) error {
	n, err := ctdo.ctd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{contracttemplate.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ctdo *ContractTemplateDeleteOne) ExecX(ctx context.Context) {
	if err := ctdo.Exec(ctx); err != nil {
		panic(err)
	}
}