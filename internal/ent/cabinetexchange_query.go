// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/cabinetexchange"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
)

// CabinetExchangeQuery is the builder for querying CabinetExchange entities.
type CabinetExchangeQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.CabinetExchange
	// eager-loading edges.
	withRider   *RiderQuery
	withCabinet *CabinetQuery
	modifiers   []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CabinetExchangeQuery builder.
func (ceq *CabinetExchangeQuery) Where(ps ...predicate.CabinetExchange) *CabinetExchangeQuery {
	ceq.predicates = append(ceq.predicates, ps...)
	return ceq
}

// Limit adds a limit step to the query.
func (ceq *CabinetExchangeQuery) Limit(limit int) *CabinetExchangeQuery {
	ceq.limit = &limit
	return ceq
}

// Offset adds an offset step to the query.
func (ceq *CabinetExchangeQuery) Offset(offset int) *CabinetExchangeQuery {
	ceq.offset = &offset
	return ceq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (ceq *CabinetExchangeQuery) Unique(unique bool) *CabinetExchangeQuery {
	ceq.unique = &unique
	return ceq
}

// Order adds an order step to the query.
func (ceq *CabinetExchangeQuery) Order(o ...OrderFunc) *CabinetExchangeQuery {
	ceq.order = append(ceq.order, o...)
	return ceq
}

// QueryRider chains the current query on the "rider" edge.
func (ceq *CabinetExchangeQuery) QueryRider() *RiderQuery {
	query := &RiderQuery{config: ceq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ceq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := ceq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(cabinetexchange.Table, cabinetexchange.FieldID, selector),
			sqlgraph.To(rider.Table, rider.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, cabinetexchange.RiderTable, cabinetexchange.RiderColumn),
		)
		fromU = sqlgraph.SetNeighbors(ceq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryCabinet chains the current query on the "cabinet" edge.
func (ceq *CabinetExchangeQuery) QueryCabinet() *CabinetQuery {
	query := &CabinetQuery{config: ceq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ceq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := ceq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(cabinetexchange.Table, cabinetexchange.FieldID, selector),
			sqlgraph.To(cabinet.Table, cabinet.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, cabinetexchange.CabinetTable, cabinetexchange.CabinetColumn),
		)
		fromU = sqlgraph.SetNeighbors(ceq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CabinetExchange entity from the query.
// Returns a *NotFoundError when no CabinetExchange was found.
func (ceq *CabinetExchangeQuery) First(ctx context.Context) (*CabinetExchange, error) {
	nodes, err := ceq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{cabinetexchange.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ceq *CabinetExchangeQuery) FirstX(ctx context.Context) *CabinetExchange {
	node, err := ceq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first CabinetExchange ID from the query.
// Returns a *NotFoundError when no CabinetExchange ID was found.
func (ceq *CabinetExchangeQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = ceq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{cabinetexchange.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (ceq *CabinetExchangeQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := ceq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single CabinetExchange entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one CabinetExchange entity is found.
// Returns a *NotFoundError when no CabinetExchange entities are found.
func (ceq *CabinetExchangeQuery) Only(ctx context.Context) (*CabinetExchange, error) {
	nodes, err := ceq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{cabinetexchange.Label}
	default:
		return nil, &NotSingularError{cabinetexchange.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ceq *CabinetExchangeQuery) OnlyX(ctx context.Context) *CabinetExchange {
	node, err := ceq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only CabinetExchange ID in the query.
// Returns a *NotSingularError when more than one CabinetExchange ID is found.
// Returns a *NotFoundError when no entities are found.
func (ceq *CabinetExchangeQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = ceq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{cabinetexchange.Label}
	default:
		err = &NotSingularError{cabinetexchange.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (ceq *CabinetExchangeQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := ceq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CabinetExchanges.
func (ceq *CabinetExchangeQuery) All(ctx context.Context) ([]*CabinetExchange, error) {
	if err := ceq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return ceq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (ceq *CabinetExchangeQuery) AllX(ctx context.Context) []*CabinetExchange {
	nodes, err := ceq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of CabinetExchange IDs.
func (ceq *CabinetExchangeQuery) IDs(ctx context.Context) ([]uint64, error) {
	var ids []uint64
	if err := ceq.Select(cabinetexchange.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ceq *CabinetExchangeQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := ceq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ceq *CabinetExchangeQuery) Count(ctx context.Context) (int, error) {
	if err := ceq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return ceq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (ceq *CabinetExchangeQuery) CountX(ctx context.Context) int {
	count, err := ceq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ceq *CabinetExchangeQuery) Exist(ctx context.Context) (bool, error) {
	if err := ceq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return ceq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (ceq *CabinetExchangeQuery) ExistX(ctx context.Context) bool {
	exist, err := ceq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CabinetExchangeQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ceq *CabinetExchangeQuery) Clone() *CabinetExchangeQuery {
	if ceq == nil {
		return nil
	}
	return &CabinetExchangeQuery{
		config:      ceq.config,
		limit:       ceq.limit,
		offset:      ceq.offset,
		order:       append([]OrderFunc{}, ceq.order...),
		predicates:  append([]predicate.CabinetExchange{}, ceq.predicates...),
		withRider:   ceq.withRider.Clone(),
		withCabinet: ceq.withCabinet.Clone(),
		// clone intermediate query.
		sql:    ceq.sql.Clone(),
		path:   ceq.path,
		unique: ceq.unique,
	}
}

// WithRider tells the query-builder to eager-load the nodes that are connected to
// the "rider" edge. The optional arguments are used to configure the query builder of the edge.
func (ceq *CabinetExchangeQuery) WithRider(opts ...func(*RiderQuery)) *CabinetExchangeQuery {
	query := &RiderQuery{config: ceq.config}
	for _, opt := range opts {
		opt(query)
	}
	ceq.withRider = query
	return ceq
}

// WithCabinet tells the query-builder to eager-load the nodes that are connected to
// the "cabinet" edge. The optional arguments are used to configure the query builder of the edge.
func (ceq *CabinetExchangeQuery) WithCabinet(opts ...func(*CabinetQuery)) *CabinetExchangeQuery {
	query := &CabinetQuery{config: ceq.config}
	for _, opt := range opts {
		opt(query)
	}
	ceq.withCabinet = query
	return ceq
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
//	client.CabinetExchange.Query().
//		GroupBy(cabinetexchange.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (ceq *CabinetExchangeQuery) GroupBy(field string, fields ...string) *CabinetExchangeGroupBy {
	grbuild := &CabinetExchangeGroupBy{config: ceq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ceq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ceq.sqlQuery(ctx), nil
	}
	grbuild.label = cabinetexchange.Label
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
//	client.CabinetExchange.Query().
//		Select(cabinetexchange.FieldCreatedAt).
//		Scan(ctx, &v)
//
func (ceq *CabinetExchangeQuery) Select(fields ...string) *CabinetExchangeSelect {
	ceq.fields = append(ceq.fields, fields...)
	selbuild := &CabinetExchangeSelect{CabinetExchangeQuery: ceq}
	selbuild.label = cabinetexchange.Label
	selbuild.flds, selbuild.scan = &ceq.fields, selbuild.Scan
	return selbuild
}

func (ceq *CabinetExchangeQuery) prepareQuery(ctx context.Context) error {
	for _, f := range ceq.fields {
		if !cabinetexchange.ValidColumn(f) {
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

func (ceq *CabinetExchangeQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*CabinetExchange, error) {
	var (
		nodes       = []*CabinetExchange{}
		_spec       = ceq.querySpec()
		loadedTypes = [2]bool{
			ceq.withRider != nil,
			ceq.withCabinet != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		return (*CabinetExchange).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		node := &CabinetExchange{config: ceq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
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

	if query := ceq.withRider; query != nil {
		ids := make([]uint64, 0, len(nodes))
		nodeids := make(map[uint64][]*CabinetExchange)
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
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "rider_id" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Rider = n
			}
		}
	}

	if query := ceq.withCabinet; query != nil {
		ids := make([]uint64, 0, len(nodes))
		nodeids := make(map[uint64][]*CabinetExchange)
		for i := range nodes {
			fk := nodes[i].CabinetID
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(cabinet.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "cabinet_id" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Cabinet = n
			}
		}
	}

	return nodes, nil
}

func (ceq *CabinetExchangeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ceq.querySpec()
	if len(ceq.modifiers) > 0 {
		_spec.Modifiers = ceq.modifiers
	}
	_spec.Node.Columns = ceq.fields
	if len(ceq.fields) > 0 {
		_spec.Unique = ceq.unique != nil && *ceq.unique
	}
	return sqlgraph.CountNodes(ctx, ceq.driver, _spec)
}

func (ceq *CabinetExchangeQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := ceq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (ceq *CabinetExchangeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cabinetexchange.Table,
			Columns: cabinetexchange.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: cabinetexchange.FieldID,
			},
		},
		From:   ceq.sql,
		Unique: true,
	}
	if unique := ceq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := ceq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cabinetexchange.FieldID)
		for i := range fields {
			if fields[i] != cabinetexchange.FieldID {
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
	if limit := ceq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ceq.offset; offset != nil {
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

func (ceq *CabinetExchangeQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(ceq.driver.Dialect())
	t1 := builder.Table(cabinetexchange.Table)
	columns := ceq.fields
	if len(columns) == 0 {
		columns = cabinetexchange.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if ceq.sql != nil {
		selector = ceq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if ceq.unique != nil && *ceq.unique {
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
	if offset := ceq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ceq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ceq *CabinetExchangeQuery) Modify(modifiers ...func(s *sql.Selector)) *CabinetExchangeSelect {
	ceq.modifiers = append(ceq.modifiers, modifiers...)
	return ceq.Select()
}

// CabinetExchangeGroupBy is the group-by builder for CabinetExchange entities.
type CabinetExchangeGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (cegb *CabinetExchangeGroupBy) Aggregate(fns ...AggregateFunc) *CabinetExchangeGroupBy {
	cegb.fns = append(cegb.fns, fns...)
	return cegb
}

// Scan applies the group-by query and scans the result into the given value.
func (cegb *CabinetExchangeGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := cegb.path(ctx)
	if err != nil {
		return err
	}
	cegb.sql = query
	return cegb.sqlScan(ctx, v)
}

func (cegb *CabinetExchangeGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range cegb.fields {
		if !cabinetexchange.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := cegb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cegb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (cegb *CabinetExchangeGroupBy) sqlQuery() *sql.Selector {
	selector := cegb.sql.Select()
	aggregation := make([]string, 0, len(cegb.fns))
	for _, fn := range cegb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(cegb.fields)+len(cegb.fns))
		for _, f := range cegb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(cegb.fields...)...)
}

// CabinetExchangeSelect is the builder for selecting fields of CabinetExchange entities.
type CabinetExchangeSelect struct {
	*CabinetExchangeQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (ces *CabinetExchangeSelect) Scan(ctx context.Context, v interface{}) error {
	if err := ces.prepareQuery(ctx); err != nil {
		return err
	}
	ces.sql = ces.CabinetExchangeQuery.sqlQuery(ctx)
	return ces.sqlScan(ctx, v)
}

func (ces *CabinetExchangeSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ces.sql.Query()
	if err := ces.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ces *CabinetExchangeSelect) Modify(modifiers ...func(s *sql.Selector)) *CabinetExchangeSelect {
	ces.modifiers = append(ces.modifiers, modifiers...)
	return ces
}