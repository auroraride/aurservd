// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/assetscrap"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// AssetScrapDelete is the builder for deleting a AssetScrap entity.
type AssetScrapDelete struct {
	config
	hooks    []Hook
	mutation *AssetScrapMutation
}

// Where appends a list predicates to the AssetScrapDelete builder.
func (asd *AssetScrapDelete) Where(ps ...predicate.AssetScrap) *AssetScrapDelete {
	asd.mutation.Where(ps...)
	return asd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (asd *AssetScrapDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, asd.sqlExec, asd.mutation, asd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (asd *AssetScrapDelete) ExecX(ctx context.Context) int {
	n, err := asd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (asd *AssetScrapDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(assetscrap.Table, sqlgraph.NewFieldSpec(assetscrap.FieldID, field.TypeUint64))
	if ps := asd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, asd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	asd.mutation.done = true
	return affected, err
}

// AssetScrapDeleteOne is the builder for deleting a single AssetScrap entity.
type AssetScrapDeleteOne struct {
	asd *AssetScrapDelete
}

// Where appends a list predicates to the AssetScrapDelete builder.
func (asdo *AssetScrapDeleteOne) Where(ps ...predicate.AssetScrap) *AssetScrapDeleteOne {
	asdo.asd.mutation.Where(ps...)
	return asdo
}

// Exec executes the deletion query.
func (asdo *AssetScrapDeleteOne) Exec(ctx context.Context) error {
	n, err := asdo.asd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{assetscrap.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (asdo *AssetScrapDeleteOne) ExecX(ctx context.Context) {
	if err := asdo.Exec(ctx); err != nil {
		panic(err)
	}
}