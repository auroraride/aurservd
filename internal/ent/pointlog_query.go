// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/pointlog"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
)

// PointLogQuery is the builder for querying PointLog entities.
type PointLogQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.PointLog
	withRider  *RiderQuery
	withOrder  *OrderQuery
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PointLogQuery builder.
func (plq *PointLogQuery) Where(ps ...predicate.PointLog) *PointLogQuery {
	plq.predicates = append(plq.predicates, ps...)
	return plq
}

// Limit adds a limit step to the query.
func (plq *PointLogQuery) Limit(limit int) *PointLogQuery {
	plq.limit = &limit
	return plq
}

// Offset adds an offset step to the query.
func (plq *PointLogQuery) Offset(offset int) *PointLogQuery {
	plq.offset = &offset
	return plq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (plq *PointLogQuery) Unique(unique bool) *PointLogQuery {
	plq.unique = &unique
	return plq
}

// Order adds an order step to the query.
func (plq *PointLogQuery) Order(o ...OrderFunc) *PointLogQuery {
	plq.order = append(plq.order, o...)
	return plq
}

// QueryRider chains the current query on the "rider" edge.
func (plq *PointLogQuery) QueryRider() *RiderQuery {
	query := &RiderQuery{config: plq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := plq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := plq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(pointlog.Table, pointlog.FieldID, selector),
			sqlgraph.To(rider.Table, rider.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, pointlog.RiderTable, pointlog.RiderColumn),
		)
		fromU = sqlgraph.SetNeighbors(plq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryOrder chains the current query on the "order" edge.
func (plq *PointLogQuery) QueryOrder() *OrderQuery {
	query := &OrderQuery{config: plq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := plq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := plq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(pointlog.Table, pointlog.FieldID, selector),
			sqlgraph.To(order.Table, order.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, pointlog.OrderTable, pointlog.OrderColumn),
		)
		fromU = sqlgraph.SetNeighbors(plq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first PointLog entity from the query.
// Returns a *NotFoundError when no PointLog was found.
func (plq *PointLogQuery) First(ctx context.Context) (*PointLog, error) {
	nodes, err := plq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{pointlog.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (plq *PointLogQuery) FirstX(ctx context.Context) *PointLog {
	node, err := plq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first PointLog ID from the query.
// Returns a *NotFoundError when no PointLog ID was found.
func (plq *PointLogQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = plq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{pointlog.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (plq *PointLogQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := plq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single PointLog entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one PointLog entity is found.
// Returns a *NotFoundError when no PointLog entities are found.
func (plq *PointLogQuery) Only(ctx context.Context) (*PointLog, error) {
	nodes, err := plq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{pointlog.Label}
	default:
		return nil, &NotSingularError{pointlog.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (plq *PointLogQuery) OnlyX(ctx context.Context) *PointLog {
	node, err := plq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only PointLog ID in the query.
// Returns a *NotSingularError when more than one PointLog ID is found.
// Returns a *NotFoundError when no entities are found.
func (plq *PointLogQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = plq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{pointlog.Label}
	default:
		err = &NotSingularError{pointlog.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (plq *PointLogQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := plq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PointLogs.
func (plq *PointLogQuery) All(ctx context.Context) ([]*PointLog, error) {
	if err := plq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return plq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (plq *PointLogQuery) AllX(ctx context.Context) []*PointLog {
	nodes, err := plq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of PointLog IDs.
func (plq *PointLogQuery) IDs(ctx context.Context) ([]uint64, error) {
	var ids []uint64
	if err := plq.Select(pointlog.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (plq *PointLogQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := plq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (plq *PointLogQuery) Count(ctx context.Context) (int, error) {
	if err := plq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return plq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (plq *PointLogQuery) CountX(ctx context.Context) int {
	count, err := plq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (plq *PointLogQuery) Exist(ctx context.Context) (bool, error) {
	if err := plq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return plq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (plq *PointLogQuery) ExistX(ctx context.Context) bool {
	exist, err := plq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PointLogQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (plq *PointLogQuery) Clone() *PointLogQuery {
	if plq == nil {
		return nil
	}
	return &PointLogQuery{
		config:     plq.config,
		limit:      plq.limit,
		offset:     plq.offset,
		order:      append([]OrderFunc{}, plq.order...),
		predicates: append([]predicate.PointLog{}, plq.predicates...),
		withRider:  plq.withRider.Clone(),
		withOrder:  plq.withOrder.Clone(),
		// clone intermediate query.
		sql:    plq.sql.Clone(),
		path:   plq.path,
		unique: plq.unique,
	}
}

// WithRider tells the query-builder to eager-load the nodes that are connected to
// the "rider" edge. The optional arguments are used to configure the query builder of the edge.
func (plq *PointLogQuery) WithRider(opts ...func(*RiderQuery)) *PointLogQuery {
	query := &RiderQuery{config: plq.config}
	for _, opt := range opts {
		opt(query)
	}
	plq.withRider = query
	return plq
}

// WithOrder tells the query-builder to eager-load the nodes that are connected to
// the "order" edge. The optional arguments are used to configure the query builder of the edge.
func (plq *PointLogQuery) WithOrder(opts ...func(*OrderQuery)) *PointLogQuery {
	query := &OrderQuery{config: plq.config}
	for _, opt := range opts {
		opt(query)
	}
	plq.withOrder = query
	return plq
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
//	client.PointLog.Query().
//		GroupBy(pointlog.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (plq *PointLogQuery) GroupBy(field string, fields ...string) *PointLogGroupBy {
	grbuild := &PointLogGroupBy{config: plq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := plq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return plq.sqlQuery(ctx), nil
	}
	grbuild.label = pointlog.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
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
//	client.PointLog.Query().
//		Select(pointlog.FieldCreatedAt).
//		Scan(ctx, &v)
func (plq *PointLogQuery) Select(fields ...string) *PointLogSelect {
	plq.fields = append(plq.fields, fields...)
	selbuild := &PointLogSelect{PointLogQuery: plq}
	selbuild.label = pointlog.Label
	selbuild.flds, selbuild.scan = &plq.fields, selbuild.Scan
	return selbuild
}

func (plq *PointLogQuery) prepareQuery(ctx context.Context) error {
	for _, f := range plq.fields {
		if !pointlog.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if plq.path != nil {
		prev, err := plq.path(ctx)
		if err != nil {
			return err
		}
		plq.sql = prev
	}
	return nil
}

func (plq *PointLogQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*PointLog, error) {
	var (
		nodes       = []*PointLog{}
		_spec       = plq.querySpec()
		loadedTypes = [2]bool{
			plq.withRider != nil,
			plq.withOrder != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*PointLog).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &PointLog{config: plq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(plq.modifiers) > 0 {
		_spec.Modifiers = plq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, plq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := plq.withRider; query != nil {
		if err := plq.loadRider(ctx, query, nodes, nil,
			func(n *PointLog, e *Rider) { n.Edges.Rider = e }); err != nil {
			return nil, err
		}
	}
	if query := plq.withOrder; query != nil {
		if err := plq.loadOrder(ctx, query, nodes, nil,
			func(n *PointLog, e *Order) { n.Edges.Order = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (plq *PointLogQuery) loadRider(ctx context.Context, query *RiderQuery, nodes []*PointLog, init func(*PointLog), assign func(*PointLog, *Rider)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*PointLog)
	for i := range nodes {
		fk := nodes[i].RiderID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
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
func (plq *PointLogQuery) loadOrder(ctx context.Context, query *OrderQuery, nodes []*PointLog, init func(*PointLog), assign func(*PointLog, *Order)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*PointLog)
	for i := range nodes {
		if nodes[i].OrderID == nil {
			continue
		}
		fk := *nodes[i].OrderID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(order.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "order_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (plq *PointLogQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := plq.querySpec()
	if len(plq.modifiers) > 0 {
		_spec.Modifiers = plq.modifiers
	}
	_spec.Node.Columns = plq.fields
	if len(plq.fields) > 0 {
		_spec.Unique = plq.unique != nil && *plq.unique
	}
	return sqlgraph.CountNodes(ctx, plq.driver, _spec)
}

func (plq *PointLogQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := plq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (plq *PointLogQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   pointlog.Table,
			Columns: pointlog.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: pointlog.FieldID,
			},
		},
		From:   plq.sql,
		Unique: true,
	}
	if unique := plq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := plq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, pointlog.FieldID)
		for i := range fields {
			if fields[i] != pointlog.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := plq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := plq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := plq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := plq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (plq *PointLogQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(plq.driver.Dialect())
	t1 := builder.Table(pointlog.Table)
	columns := plq.fields
	if len(columns) == 0 {
		columns = pointlog.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if plq.sql != nil {
		selector = plq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if plq.unique != nil && *plq.unique {
		selector.Distinct()
	}
	for _, m := range plq.modifiers {
		m(selector)
	}
	for _, p := range plq.predicates {
		p(selector)
	}
	for _, p := range plq.order {
		p(selector)
	}
	if offset := plq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := plq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (plq *PointLogQuery) Modify(modifiers ...func(s *sql.Selector)) *PointLogSelect {
	plq.modifiers = append(plq.modifiers, modifiers...)
	return plq.Select()
}

// PointLogGroupBy is the group-by builder for PointLog entities.
type PointLogGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (plgb *PointLogGroupBy) Aggregate(fns ...AggregateFunc) *PointLogGroupBy {
	plgb.fns = append(plgb.fns, fns...)
	return plgb
}

// Scan applies the group-by query and scans the result into the given value.
func (plgb *PointLogGroupBy) Scan(ctx context.Context, v any) error {
	query, err := plgb.path(ctx)
	if err != nil {
		return err
	}
	plgb.sql = query
	return plgb.sqlScan(ctx, v)
}

func (plgb *PointLogGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range plgb.fields {
		if !pointlog.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := plgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := plgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (plgb *PointLogGroupBy) sqlQuery() *sql.Selector {
	selector := plgb.sql.Select()
	aggregation := make([]string, 0, len(plgb.fns))
	for _, fn := range plgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(plgb.fields)+len(plgb.fns))
		for _, f := range plgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(plgb.fields...)...)
}

// PointLogSelect is the builder for selecting fields of PointLog entities.
type PointLogSelect struct {
	*PointLogQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (pls *PointLogSelect) Scan(ctx context.Context, v any) error {
	if err := pls.prepareQuery(ctx); err != nil {
		return err
	}
	pls.sql = pls.PointLogQuery.sqlQuery(ctx)
	return pls.sqlScan(ctx, v)
}

func (pls *PointLogSelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := pls.sql.Query()
	if err := pls.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (pls *PointLogSelect) Modify(modifiers ...func(s *sql.Selector)) *PointLogSelect {
	pls.modifiers = append(pls.modifiers, modifiers...)
	return pls
}