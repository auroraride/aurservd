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
	"github.com/auroraride/aurservd/internal/ent/manager"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/riderfollowup"
)

// RiderFollowUpQuery is the builder for querying RiderFollowUp entities.
type RiderFollowUpQuery struct {
	config
	ctx         *QueryContext
	order       []riderfollowup.OrderOption
	inters      []Interceptor
	predicates  []predicate.RiderFollowUp
	withManager *ManagerQuery
	withRider   *RiderQuery
	modifiers   []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the RiderFollowUpQuery builder.
func (rfuq *RiderFollowUpQuery) Where(ps ...predicate.RiderFollowUp) *RiderFollowUpQuery {
	rfuq.predicates = append(rfuq.predicates, ps...)
	return rfuq
}

// Limit the number of records to be returned by this query.
func (rfuq *RiderFollowUpQuery) Limit(limit int) *RiderFollowUpQuery {
	rfuq.ctx.Limit = &limit
	return rfuq
}

// Offset to start from.
func (rfuq *RiderFollowUpQuery) Offset(offset int) *RiderFollowUpQuery {
	rfuq.ctx.Offset = &offset
	return rfuq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rfuq *RiderFollowUpQuery) Unique(unique bool) *RiderFollowUpQuery {
	rfuq.ctx.Unique = &unique
	return rfuq
}

// Order specifies how the records should be ordered.
func (rfuq *RiderFollowUpQuery) Order(o ...riderfollowup.OrderOption) *RiderFollowUpQuery {
	rfuq.order = append(rfuq.order, o...)
	return rfuq
}

// QueryManager chains the current query on the "manager" edge.
func (rfuq *RiderFollowUpQuery) QueryManager() *ManagerQuery {
	query := (&ManagerClient{config: rfuq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rfuq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rfuq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(riderfollowup.Table, riderfollowup.FieldID, selector),
			sqlgraph.To(manager.Table, manager.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, riderfollowup.ManagerTable, riderfollowup.ManagerColumn),
		)
		fromU = sqlgraph.SetNeighbors(rfuq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRider chains the current query on the "rider" edge.
func (rfuq *RiderFollowUpQuery) QueryRider() *RiderQuery {
	query := (&RiderClient{config: rfuq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rfuq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rfuq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(riderfollowup.Table, riderfollowup.FieldID, selector),
			sqlgraph.To(rider.Table, rider.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, riderfollowup.RiderTable, riderfollowup.RiderColumn),
		)
		fromU = sqlgraph.SetNeighbors(rfuq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first RiderFollowUp entity from the query.
// Returns a *NotFoundError when no RiderFollowUp was found.
func (rfuq *RiderFollowUpQuery) First(ctx context.Context) (*RiderFollowUp, error) {
	nodes, err := rfuq.Limit(1).All(setContextOp(ctx, rfuq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{riderfollowup.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rfuq *RiderFollowUpQuery) FirstX(ctx context.Context) *RiderFollowUp {
	node, err := rfuq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first RiderFollowUp ID from the query.
// Returns a *NotFoundError when no RiderFollowUp ID was found.
func (rfuq *RiderFollowUpQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = rfuq.Limit(1).IDs(setContextOp(ctx, rfuq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{riderfollowup.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rfuq *RiderFollowUpQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := rfuq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single RiderFollowUp entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one RiderFollowUp entity is found.
// Returns a *NotFoundError when no RiderFollowUp entities are found.
func (rfuq *RiderFollowUpQuery) Only(ctx context.Context) (*RiderFollowUp, error) {
	nodes, err := rfuq.Limit(2).All(setContextOp(ctx, rfuq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{riderfollowup.Label}
	default:
		return nil, &NotSingularError{riderfollowup.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rfuq *RiderFollowUpQuery) OnlyX(ctx context.Context) *RiderFollowUp {
	node, err := rfuq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only RiderFollowUp ID in the query.
// Returns a *NotSingularError when more than one RiderFollowUp ID is found.
// Returns a *NotFoundError when no entities are found.
func (rfuq *RiderFollowUpQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = rfuq.Limit(2).IDs(setContextOp(ctx, rfuq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{riderfollowup.Label}
	default:
		err = &NotSingularError{riderfollowup.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rfuq *RiderFollowUpQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := rfuq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of RiderFollowUps.
func (rfuq *RiderFollowUpQuery) All(ctx context.Context) ([]*RiderFollowUp, error) {
	ctx = setContextOp(ctx, rfuq.ctx, ent.OpQueryAll)
	if err := rfuq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*RiderFollowUp, *RiderFollowUpQuery]()
	return withInterceptors[[]*RiderFollowUp](ctx, rfuq, qr, rfuq.inters)
}

// AllX is like All, but panics if an error occurs.
func (rfuq *RiderFollowUpQuery) AllX(ctx context.Context) []*RiderFollowUp {
	nodes, err := rfuq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of RiderFollowUp IDs.
func (rfuq *RiderFollowUpQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if rfuq.ctx.Unique == nil && rfuq.path != nil {
		rfuq.Unique(true)
	}
	ctx = setContextOp(ctx, rfuq.ctx, ent.OpQueryIDs)
	if err = rfuq.Select(riderfollowup.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rfuq *RiderFollowUpQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := rfuq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rfuq *RiderFollowUpQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, rfuq.ctx, ent.OpQueryCount)
	if err := rfuq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, rfuq, querierCount[*RiderFollowUpQuery](), rfuq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (rfuq *RiderFollowUpQuery) CountX(ctx context.Context) int {
	count, err := rfuq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rfuq *RiderFollowUpQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, rfuq.ctx, ent.OpQueryExist)
	switch _, err := rfuq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (rfuq *RiderFollowUpQuery) ExistX(ctx context.Context) bool {
	exist, err := rfuq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the RiderFollowUpQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rfuq *RiderFollowUpQuery) Clone() *RiderFollowUpQuery {
	if rfuq == nil {
		return nil
	}
	return &RiderFollowUpQuery{
		config:      rfuq.config,
		ctx:         rfuq.ctx.Clone(),
		order:       append([]riderfollowup.OrderOption{}, rfuq.order...),
		inters:      append([]Interceptor{}, rfuq.inters...),
		predicates:  append([]predicate.RiderFollowUp{}, rfuq.predicates...),
		withManager: rfuq.withManager.Clone(),
		withRider:   rfuq.withRider.Clone(),
		// clone intermediate query.
		sql:       rfuq.sql.Clone(),
		path:      rfuq.path,
		modifiers: append([]func(*sql.Selector){}, rfuq.modifiers...),
	}
}

// WithManager tells the query-builder to eager-load the nodes that are connected to
// the "manager" edge. The optional arguments are used to configure the query builder of the edge.
func (rfuq *RiderFollowUpQuery) WithManager(opts ...func(*ManagerQuery)) *RiderFollowUpQuery {
	query := (&ManagerClient{config: rfuq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rfuq.withManager = query
	return rfuq
}

// WithRider tells the query-builder to eager-load the nodes that are connected to
// the "rider" edge. The optional arguments are used to configure the query builder of the edge.
func (rfuq *RiderFollowUpQuery) WithRider(opts ...func(*RiderQuery)) *RiderFollowUpQuery {
	query := (&RiderClient{config: rfuq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rfuq.withRider = query
	return rfuq
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
//	client.RiderFollowUp.Query().
//		GroupBy(riderfollowup.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (rfuq *RiderFollowUpQuery) GroupBy(field string, fields ...string) *RiderFollowUpGroupBy {
	rfuq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &RiderFollowUpGroupBy{build: rfuq}
	grbuild.flds = &rfuq.ctx.Fields
	grbuild.label = riderfollowup.Label
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
//	client.RiderFollowUp.Query().
//		Select(riderfollowup.FieldCreatedAt).
//		Scan(ctx, &v)
func (rfuq *RiderFollowUpQuery) Select(fields ...string) *RiderFollowUpSelect {
	rfuq.ctx.Fields = append(rfuq.ctx.Fields, fields...)
	sbuild := &RiderFollowUpSelect{RiderFollowUpQuery: rfuq}
	sbuild.label = riderfollowup.Label
	sbuild.flds, sbuild.scan = &rfuq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a RiderFollowUpSelect configured with the given aggregations.
func (rfuq *RiderFollowUpQuery) Aggregate(fns ...AggregateFunc) *RiderFollowUpSelect {
	return rfuq.Select().Aggregate(fns...)
}

func (rfuq *RiderFollowUpQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range rfuq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, rfuq); err != nil {
				return err
			}
		}
	}
	for _, f := range rfuq.ctx.Fields {
		if !riderfollowup.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if rfuq.path != nil {
		prev, err := rfuq.path(ctx)
		if err != nil {
			return err
		}
		rfuq.sql = prev
	}
	return nil
}

func (rfuq *RiderFollowUpQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*RiderFollowUp, error) {
	var (
		nodes       = []*RiderFollowUp{}
		_spec       = rfuq.querySpec()
		loadedTypes = [2]bool{
			rfuq.withManager != nil,
			rfuq.withRider != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*RiderFollowUp).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &RiderFollowUp{config: rfuq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(rfuq.modifiers) > 0 {
		_spec.Modifiers = rfuq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, rfuq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := rfuq.withManager; query != nil {
		if err := rfuq.loadManager(ctx, query, nodes, nil,
			func(n *RiderFollowUp, e *Manager) { n.Edges.Manager = e }); err != nil {
			return nil, err
		}
	}
	if query := rfuq.withRider; query != nil {
		if err := rfuq.loadRider(ctx, query, nodes, nil,
			func(n *RiderFollowUp, e *Rider) { n.Edges.Rider = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (rfuq *RiderFollowUpQuery) loadManager(ctx context.Context, query *ManagerQuery, nodes []*RiderFollowUp, init func(*RiderFollowUp), assign func(*RiderFollowUp, *Manager)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*RiderFollowUp)
	for i := range nodes {
		fk := nodes[i].ManagerID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(manager.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "manager_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (rfuq *RiderFollowUpQuery) loadRider(ctx context.Context, query *RiderQuery, nodes []*RiderFollowUp, init func(*RiderFollowUp), assign func(*RiderFollowUp, *Rider)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*RiderFollowUp)
	for i := range nodes {
		fk := nodes[i].RiderID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(rider.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "rider_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (rfuq *RiderFollowUpQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rfuq.querySpec()
	if len(rfuq.modifiers) > 0 {
		_spec.Modifiers = rfuq.modifiers
	}
	_spec.Node.Columns = rfuq.ctx.Fields
	if len(rfuq.ctx.Fields) > 0 {
		_spec.Unique = rfuq.ctx.Unique != nil && *rfuq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, rfuq.driver, _spec)
}

func (rfuq *RiderFollowUpQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(riderfollowup.Table, riderfollowup.Columns, sqlgraph.NewFieldSpec(riderfollowup.FieldID, field.TypeUint64))
	_spec.From = rfuq.sql
	if unique := rfuq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if rfuq.path != nil {
		_spec.Unique = true
	}
	if fields := rfuq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, riderfollowup.FieldID)
		for i := range fields {
			if fields[i] != riderfollowup.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if rfuq.withManager != nil {
			_spec.Node.AddColumnOnce(riderfollowup.FieldManagerID)
		}
		if rfuq.withRider != nil {
			_spec.Node.AddColumnOnce(riderfollowup.FieldRiderID)
		}
	}
	if ps := rfuq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rfuq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rfuq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rfuq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rfuq *RiderFollowUpQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rfuq.driver.Dialect())
	t1 := builder.Table(riderfollowup.Table)
	columns := rfuq.ctx.Fields
	if len(columns) == 0 {
		columns = riderfollowup.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if rfuq.sql != nil {
		selector = rfuq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if rfuq.ctx.Unique != nil && *rfuq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range rfuq.modifiers {
		m(selector)
	}
	for _, p := range rfuq.predicates {
		p(selector)
	}
	for _, p := range rfuq.order {
		p(selector)
	}
	if offset := rfuq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rfuq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (rfuq *RiderFollowUpQuery) Modify(modifiers ...func(s *sql.Selector)) *RiderFollowUpSelect {
	rfuq.modifiers = append(rfuq.modifiers, modifiers...)
	return rfuq.Select()
}

type RiderFollowUpQueryWith string

var (
	RiderFollowUpQueryWithManager RiderFollowUpQueryWith = "Manager"
	RiderFollowUpQueryWithRider   RiderFollowUpQueryWith = "Rider"
)

func (rfuq *RiderFollowUpQuery) With(withEdges ...RiderFollowUpQueryWith) *RiderFollowUpQuery {
	for _, v := range withEdges {
		switch v {
		case RiderFollowUpQueryWithManager:
			rfuq.WithManager()
		case RiderFollowUpQueryWithRider:
			rfuq.WithRider()
		}
	}
	return rfuq
}

// RiderFollowUpGroupBy is the group-by builder for RiderFollowUp entities.
type RiderFollowUpGroupBy struct {
	selector
	build *RiderFollowUpQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rfugb *RiderFollowUpGroupBy) Aggregate(fns ...AggregateFunc) *RiderFollowUpGroupBy {
	rfugb.fns = append(rfugb.fns, fns...)
	return rfugb
}

// Scan applies the selector query and scans the result into the given value.
func (rfugb *RiderFollowUpGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rfugb.build.ctx, ent.OpQueryGroupBy)
	if err := rfugb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*RiderFollowUpQuery, *RiderFollowUpGroupBy](ctx, rfugb.build, rfugb, rfugb.build.inters, v)
}

func (rfugb *RiderFollowUpGroupBy) sqlScan(ctx context.Context, root *RiderFollowUpQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(rfugb.fns))
	for _, fn := range rfugb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*rfugb.flds)+len(rfugb.fns))
		for _, f := range *rfugb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*rfugb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rfugb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// RiderFollowUpSelect is the builder for selecting fields of RiderFollowUp entities.
type RiderFollowUpSelect struct {
	*RiderFollowUpQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (rfus *RiderFollowUpSelect) Aggregate(fns ...AggregateFunc) *RiderFollowUpSelect {
	rfus.fns = append(rfus.fns, fns...)
	return rfus
}

// Scan applies the selector query and scans the result into the given value.
func (rfus *RiderFollowUpSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rfus.ctx, ent.OpQuerySelect)
	if err := rfus.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*RiderFollowUpQuery, *RiderFollowUpSelect](ctx, rfus.RiderFollowUpQuery, rfus, rfus.inters, v)
}

func (rfus *RiderFollowUpSelect) sqlScan(ctx context.Context, root *RiderFollowUpQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(rfus.fns))
	for _, fn := range rfus.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*rfus.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rfus.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (rfus *RiderFollowUpSelect) Modify(modifiers ...func(s *sql.Selector)) *RiderFollowUpSelect {
	rfus.modifiers = append(rfus.modifiers, modifiers...)
	return rfus
}
