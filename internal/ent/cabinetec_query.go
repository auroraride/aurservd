// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/cabinetec"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// CabinetEcQuery is the builder for querying CabinetEc entities.
type CabinetEcQuery struct {
	config
	ctx        *QueryContext
	order      []cabinetec.OrderOption
	inters     []Interceptor
	predicates []predicate.CabinetEc
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CabinetEcQuery builder.
func (ceq *CabinetEcQuery) Where(ps ...predicate.CabinetEc) *CabinetEcQuery {
	ceq.predicates = append(ceq.predicates, ps...)
	return ceq
}

// Limit the number of records to be returned by this query.
func (ceq *CabinetEcQuery) Limit(limit int) *CabinetEcQuery {
	ceq.ctx.Limit = &limit
	return ceq
}

// Offset to start from.
func (ceq *CabinetEcQuery) Offset(offset int) *CabinetEcQuery {
	ceq.ctx.Offset = &offset
	return ceq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (ceq *CabinetEcQuery) Unique(unique bool) *CabinetEcQuery {
	ceq.ctx.Unique = &unique
	return ceq
}

// Order specifies how the records should be ordered.
func (ceq *CabinetEcQuery) Order(o ...cabinetec.OrderOption) *CabinetEcQuery {
	ceq.order = append(ceq.order, o...)
	return ceq
}

// First returns the first CabinetEc entity from the query.
// Returns a *NotFoundError when no CabinetEc was found.
func (ceq *CabinetEcQuery) First(ctx context.Context) (*CabinetEc, error) {
	nodes, err := ceq.Limit(1).All(setContextOp(ctx, ceq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{cabinetec.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ceq *CabinetEcQuery) FirstX(ctx context.Context) *CabinetEc {
	node, err := ceq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first CabinetEc ID from the query.
// Returns a *NotFoundError when no CabinetEc ID was found.
func (ceq *CabinetEcQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = ceq.Limit(1).IDs(setContextOp(ctx, ceq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{cabinetec.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (ceq *CabinetEcQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := ceq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single CabinetEc entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one CabinetEc entity is found.
// Returns a *NotFoundError when no CabinetEc entities are found.
func (ceq *CabinetEcQuery) Only(ctx context.Context) (*CabinetEc, error) {
	nodes, err := ceq.Limit(2).All(setContextOp(ctx, ceq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{cabinetec.Label}
	default:
		return nil, &NotSingularError{cabinetec.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ceq *CabinetEcQuery) OnlyX(ctx context.Context) *CabinetEc {
	node, err := ceq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only CabinetEc ID in the query.
// Returns a *NotSingularError when more than one CabinetEc ID is found.
// Returns a *NotFoundError when no entities are found.
func (ceq *CabinetEcQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = ceq.Limit(2).IDs(setContextOp(ctx, ceq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{cabinetec.Label}
	default:
		err = &NotSingularError{cabinetec.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (ceq *CabinetEcQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := ceq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CabinetEcs.
func (ceq *CabinetEcQuery) All(ctx context.Context) ([]*CabinetEc, error) {
	ctx = setContextOp(ctx, ceq.ctx, "All")
	if err := ceq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*CabinetEc, *CabinetEcQuery]()
	return withInterceptors[[]*CabinetEc](ctx, ceq, qr, ceq.inters)
}

// AllX is like All, but panics if an error occurs.
func (ceq *CabinetEcQuery) AllX(ctx context.Context) []*CabinetEc {
	nodes, err := ceq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of CabinetEc IDs.
func (ceq *CabinetEcQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if ceq.ctx.Unique == nil && ceq.path != nil {
		ceq.Unique(true)
	}
	ctx = setContextOp(ctx, ceq.ctx, "IDs")
	if err = ceq.Select(cabinetec.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ceq *CabinetEcQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := ceq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ceq *CabinetEcQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, ceq.ctx, "Count")
	if err := ceq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, ceq, querierCount[*CabinetEcQuery](), ceq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (ceq *CabinetEcQuery) CountX(ctx context.Context) int {
	count, err := ceq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ceq *CabinetEcQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, ceq.ctx, "Exist")
	switch _, err := ceq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (ceq *CabinetEcQuery) ExistX(ctx context.Context) bool {
	exist, err := ceq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CabinetEcQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ceq *CabinetEcQuery) Clone() *CabinetEcQuery {
	if ceq == nil {
		return nil
	}
	return &CabinetEcQuery{
		config:     ceq.config,
		ctx:        ceq.ctx.Clone(),
		order:      append([]cabinetec.OrderOption{}, ceq.order...),
		inters:     append([]Interceptor{}, ceq.inters...),
		predicates: append([]predicate.CabinetEc{}, ceq.predicates...),
		// clone intermediate query.
		sql:  ceq.sql.Clone(),
		path: ceq.path,
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
//	client.CabinetEc.Query().
//		GroupBy(cabinetec.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (ceq *CabinetEcQuery) GroupBy(field string, fields ...string) *CabinetEcGroupBy {
	ceq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &CabinetEcGroupBy{build: ceq}
	grbuild.flds = &ceq.ctx.Fields
	grbuild.label = cabinetec.Label
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
//	client.CabinetEc.Query().
//		Select(cabinetec.FieldCreatedAt).
//		Scan(ctx, &v)
func (ceq *CabinetEcQuery) Select(fields ...string) *CabinetEcSelect {
	ceq.ctx.Fields = append(ceq.ctx.Fields, fields...)
	sbuild := &CabinetEcSelect{CabinetEcQuery: ceq}
	sbuild.label = cabinetec.Label
	sbuild.flds, sbuild.scan = &ceq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a CabinetEcSelect configured with the given aggregations.
func (ceq *CabinetEcQuery) Aggregate(fns ...AggregateFunc) *CabinetEcSelect {
	return ceq.Select().Aggregate(fns...)
}

func (ceq *CabinetEcQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range ceq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, ceq); err != nil {
				return err
			}
		}
	}
	for _, f := range ceq.ctx.Fields {
		if !cabinetec.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if ceq.path != nil {
		prev, err := ceq.path(ctx)
		if err != nil {
			return err
		}
		ceq.sql = prev
	}
	return nil
}

func (ceq *CabinetEcQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*CabinetEc, error) {
	var (
		nodes = []*CabinetEc{}
		_spec = ceq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*CabinetEc).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &CabinetEc{config: ceq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(ceq.modifiers) > 0 {
		_spec.Modifiers = ceq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, ceq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (ceq *CabinetEcQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ceq.querySpec()
	if len(ceq.modifiers) > 0 {
		_spec.Modifiers = ceq.modifiers
	}
	_spec.Node.Columns = ceq.ctx.Fields
	if len(ceq.ctx.Fields) > 0 {
		_spec.Unique = ceq.ctx.Unique != nil && *ceq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, ceq.driver, _spec)
}

func (ceq *CabinetEcQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(cabinetec.Table, cabinetec.Columns, sqlgraph.NewFieldSpec(cabinetec.FieldID, field.TypeUint64))
	_spec.From = ceq.sql
	if unique := ceq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if ceq.path != nil {
		_spec.Unique = true
	}
	if fields := ceq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cabinetec.FieldID)
		for i := range fields {
			if fields[i] != cabinetec.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := ceq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ceq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ceq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ceq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ceq *CabinetEcQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(ceq.driver.Dialect())
	t1 := builder.Table(cabinetec.Table)
	columns := ceq.ctx.Fields
	if len(columns) == 0 {
		columns = cabinetec.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if ceq.sql != nil {
		selector = ceq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if ceq.ctx.Unique != nil && *ceq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range ceq.modifiers {
		m(selector)
	}
	for _, p := range ceq.predicates {
		p(selector)
	}
	for _, p := range ceq.order {
		p(selector)
	}
	if offset := ceq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ceq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ceq *CabinetEcQuery) Modify(modifiers ...func(s *sql.Selector)) *CabinetEcSelect {
	ceq.modifiers = append(ceq.modifiers, modifiers...)
	return ceq.Select()
}

type CabinetEcQueryWith string

var ()

func (ceq *CabinetEcQuery) With(withEdges ...CabinetEcQueryWith) *CabinetEcQuery {
	for _, v := range withEdges {
		switch v {
		}
	}
	return ceq
}

// CabinetEcGroupBy is the group-by builder for CabinetEc entities.
type CabinetEcGroupBy struct {
	selector
	build *CabinetEcQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (cegb *CabinetEcGroupBy) Aggregate(fns ...AggregateFunc) *CabinetEcGroupBy {
	cegb.fns = append(cegb.fns, fns...)
	return cegb
}

// Scan applies the selector query and scans the result into the given value.
func (cegb *CabinetEcGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cegb.build.ctx, "GroupBy")
	if err := cegb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CabinetEcQuery, *CabinetEcGroupBy](ctx, cegb.build, cegb, cegb.build.inters, v)
}

func (cegb *CabinetEcGroupBy) sqlScan(ctx context.Context, root *CabinetEcQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(cegb.fns))
	for _, fn := range cegb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*cegb.flds)+len(cegb.fns))
		for _, f := range *cegb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*cegb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cegb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// CabinetEcSelect is the builder for selecting fields of CabinetEc entities.
type CabinetEcSelect struct {
	*CabinetEcQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ces *CabinetEcSelect) Aggregate(fns ...AggregateFunc) *CabinetEcSelect {
	ces.fns = append(ces.fns, fns...)
	return ces
}

// Scan applies the selector query and scans the result into the given value.
func (ces *CabinetEcSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ces.ctx, "Select")
	if err := ces.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CabinetEcQuery, *CabinetEcSelect](ctx, ces.CabinetEcQuery, ces, ces.inters, v)
}

func (ces *CabinetEcSelect) sqlScan(ctx context.Context, root *CabinetEcQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ces.fns))
	for _, fn := range ces.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ces.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ces.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ces *CabinetEcSelect) Modify(modifiers ...func(s *sql.Selector)) *CabinetEcSelect {
	ces.modifiers = append(ces.modifiers, modifiers...)
	return ces
}