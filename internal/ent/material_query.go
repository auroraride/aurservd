// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// MaterialQuery is the builder for querying Material entities.
type MaterialQuery struct {
	config
	ctx        *QueryContext
	order      []material.OrderOption
	inters     []Interceptor
	predicates []predicate.Material
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the MaterialQuery builder.
func (mq *MaterialQuery) Where(ps ...predicate.Material) *MaterialQuery {
	mq.predicates = append(mq.predicates, ps...)
	return mq
}

// Limit the number of records to be returned by this query.
func (mq *MaterialQuery) Limit(limit int) *MaterialQuery {
	mq.ctx.Limit = &limit
	return mq
}

// Offset to start from.
func (mq *MaterialQuery) Offset(offset int) *MaterialQuery {
	mq.ctx.Offset = &offset
	return mq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (mq *MaterialQuery) Unique(unique bool) *MaterialQuery {
	mq.ctx.Unique = &unique
	return mq
}

// Order specifies how the records should be ordered.
func (mq *MaterialQuery) Order(o ...material.OrderOption) *MaterialQuery {
	mq.order = append(mq.order, o...)
	return mq
}

// First returns the first Material entity from the query.
// Returns a *NotFoundError when no Material was found.
func (mq *MaterialQuery) First(ctx context.Context) (*Material, error) {
	nodes, err := mq.Limit(1).All(setContextOp(ctx, mq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{material.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (mq *MaterialQuery) FirstX(ctx context.Context) *Material {
	node, err := mq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Material ID from the query.
// Returns a *NotFoundError when no Material ID was found.
func (mq *MaterialQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = mq.Limit(1).IDs(setContextOp(ctx, mq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{material.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (mq *MaterialQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := mq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Material entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Material entity is found.
// Returns a *NotFoundError when no Material entities are found.
func (mq *MaterialQuery) Only(ctx context.Context) (*Material, error) {
	nodes, err := mq.Limit(2).All(setContextOp(ctx, mq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{material.Label}
	default:
		return nil, &NotSingularError{material.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (mq *MaterialQuery) OnlyX(ctx context.Context) *Material {
	node, err := mq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Material ID in the query.
// Returns a *NotSingularError when more than one Material ID is found.
// Returns a *NotFoundError when no entities are found.
func (mq *MaterialQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = mq.Limit(2).IDs(setContextOp(ctx, mq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{material.Label}
	default:
		err = &NotSingularError{material.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (mq *MaterialQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := mq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Materials.
func (mq *MaterialQuery) All(ctx context.Context) ([]*Material, error) {
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryAll)
	if err := mq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Material, *MaterialQuery]()
	return withInterceptors[[]*Material](ctx, mq, qr, mq.inters)
}

// AllX is like All, but panics if an error occurs.
func (mq *MaterialQuery) AllX(ctx context.Context) []*Material {
	nodes, err := mq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Material IDs.
func (mq *MaterialQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if mq.ctx.Unique == nil && mq.path != nil {
		mq.Unique(true)
	}
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryIDs)
	if err = mq.Select(material.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (mq *MaterialQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := mq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (mq *MaterialQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryCount)
	if err := mq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, mq, querierCount[*MaterialQuery](), mq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (mq *MaterialQuery) CountX(ctx context.Context) int {
	count, err := mq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (mq *MaterialQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryExist)
	switch _, err := mq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (mq *MaterialQuery) ExistX(ctx context.Context) bool {
	exist, err := mq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the MaterialQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (mq *MaterialQuery) Clone() *MaterialQuery {
	if mq == nil {
		return nil
	}
	return &MaterialQuery{
		config:     mq.config,
		ctx:        mq.ctx.Clone(),
		order:      append([]material.OrderOption{}, mq.order...),
		inters:     append([]Interceptor{}, mq.inters...),
		predicates: append([]predicate.Material{}, mq.predicates...),
		// clone intermediate query.
		sql:       mq.sql.Clone(),
		path:      mq.path,
		modifiers: append([]func(*sql.Selector){}, mq.modifiers...),
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Material.Query().
//		GroupBy(material.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (mq *MaterialQuery) GroupBy(field string, fields ...string) *MaterialGroupBy {
	mq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &MaterialGroupBy{build: mq}
	grbuild.flds = &mq.ctx.Fields
	grbuild.label = material.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//	}
//
//	client.Material.Query().
//		Select(material.FieldCreatedAt).
//		Scan(ctx, &v)
func (mq *MaterialQuery) Select(fields ...string) *MaterialSelect {
	mq.ctx.Fields = append(mq.ctx.Fields, fields...)
	sbuild := &MaterialSelect{MaterialQuery: mq}
	sbuild.label = material.Label
	sbuild.flds, sbuild.scan = &mq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a MaterialSelect configured with the given aggregations.
func (mq *MaterialQuery) Aggregate(fns ...AggregateFunc) *MaterialSelect {
	return mq.Select().Aggregate(fns...)
}

func (mq *MaterialQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range mq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, mq); err != nil {
				return err
			}
		}
	}
	for _, f := range mq.ctx.Fields {
		if !material.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if mq.path != nil {
		prev, err := mq.path(ctx)
		if err != nil {
			return err
		}
		mq.sql = prev
	}
	return nil
}

func (mq *MaterialQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Material, error) {
	var (
		nodes = []*Material{}
		_spec = mq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Material).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Material{config: mq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(mq.modifiers) > 0 {
		_spec.Modifiers = mq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, mq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (mq *MaterialQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := mq.querySpec()
	if len(mq.modifiers) > 0 {
		_spec.Modifiers = mq.modifiers
	}
	_spec.Node.Columns = mq.ctx.Fields
	if len(mq.ctx.Fields) > 0 {
		_spec.Unique = mq.ctx.Unique != nil && *mq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, mq.driver, _spec)
}

func (mq *MaterialQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(material.Table, material.Columns, sqlgraph.NewFieldSpec(material.FieldID, field.TypeUint64))
	_spec.From = mq.sql
	if unique := mq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if mq.path != nil {
		_spec.Unique = true
	}
	if fields := mq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, material.FieldID)
		for i := range fields {
			if fields[i] != material.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := mq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := mq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := mq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := mq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (mq *MaterialQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(mq.driver.Dialect())
	t1 := builder.Table(material.Table)
	columns := mq.ctx.Fields
	if len(columns) == 0 {
		columns = material.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if mq.sql != nil {
		selector = mq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if mq.ctx.Unique != nil && *mq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range mq.modifiers {
		m(selector)
	}
	for _, p := range mq.predicates {
		p(selector)
	}
	for _, p := range mq.order {
		p(selector)
	}
	if offset := mq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := mq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (mq *MaterialQuery) Modify(modifiers ...func(s *sql.Selector)) *MaterialSelect {
	mq.modifiers = append(mq.modifiers, modifiers...)
	return mq.Select()
}

type MaterialQueryWith string

var ()

func (mq *MaterialQuery) With(withEdges ...MaterialQueryWith) *MaterialQuery {
	for _, v := range withEdges {
		switch v {
		}
	}
	return mq
}

// MaterialGroupBy is the group-by builder for Material entities.
type MaterialGroupBy struct {
	selector
	build *MaterialQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (mgb *MaterialGroupBy) Aggregate(fns ...AggregateFunc) *MaterialGroupBy {
	mgb.fns = append(mgb.fns, fns...)
	return mgb
}

// Scan applies the selector query and scans the result into the given value.
func (mgb *MaterialGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mgb.build.ctx, ent.OpQueryGroupBy)
	if err := mgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*MaterialQuery, *MaterialGroupBy](ctx, mgb.build, mgb, mgb.build.inters, v)
}

func (mgb *MaterialGroupBy) sqlScan(ctx context.Context, root *MaterialQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(mgb.fns))
	for _, fn := range mgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*mgb.flds)+len(mgb.fns))
		for _, f := range *mgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*mgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// MaterialSelect is the builder for selecting fields of Material entities.
type MaterialSelect struct {
	*MaterialQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ms *MaterialSelect) Aggregate(fns ...AggregateFunc) *MaterialSelect {
	ms.fns = append(ms.fns, fns...)
	return ms
}

// Scan applies the selector query and scans the result into the given value.
func (ms *MaterialSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ms.ctx, ent.OpQuerySelect)
	if err := ms.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*MaterialQuery, *MaterialSelect](ctx, ms.MaterialQuery, ms, ms.inters, v)
}

func (ms *MaterialSelect) sqlScan(ctx context.Context, root *MaterialQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ms.fns))
	for _, fn := range ms.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ms.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ms.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ms *MaterialSelect) Modify(modifiers ...func(s *sql.Selector)) *MaterialSelect {
	ms.modifiers = append(ms.modifiers, modifiers...)
	return ms
}
