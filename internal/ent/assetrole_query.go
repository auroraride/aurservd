// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assetrole"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// AssetRoleQuery is the builder for querying AssetRole entities.
type AssetRoleQuery struct {
	config
	ctx               *QueryContext
	order             []assetrole.OrderOption
	inters            []Interceptor
	predicates        []predicate.AssetRole
	withAssetManagers *AssetManagerQuery
	modifiers         []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AssetRoleQuery builder.
func (arq *AssetRoleQuery) Where(ps ...predicate.AssetRole) *AssetRoleQuery {
	arq.predicates = append(arq.predicates, ps...)
	return arq
}

// Limit the number of records to be returned by this query.
func (arq *AssetRoleQuery) Limit(limit int) *AssetRoleQuery {
	arq.ctx.Limit = &limit
	return arq
}

// Offset to start from.
func (arq *AssetRoleQuery) Offset(offset int) *AssetRoleQuery {
	arq.ctx.Offset = &offset
	return arq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (arq *AssetRoleQuery) Unique(unique bool) *AssetRoleQuery {
	arq.ctx.Unique = &unique
	return arq
}

// Order specifies how the records should be ordered.
func (arq *AssetRoleQuery) Order(o ...assetrole.OrderOption) *AssetRoleQuery {
	arq.order = append(arq.order, o...)
	return arq
}

// QueryAssetManagers chains the current query on the "asset_managers" edge.
func (arq *AssetRoleQuery) QueryAssetManagers() *AssetManagerQuery {
	query := (&AssetManagerClient{config: arq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := arq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := arq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetrole.Table, assetrole.FieldID, selector),
			sqlgraph.To(assetmanager.Table, assetmanager.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, assetrole.AssetManagersTable, assetrole.AssetManagersColumn),
		)
		fromU = sqlgraph.SetNeighbors(arq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first AssetRole entity from the query.
// Returns a *NotFoundError when no AssetRole was found.
func (arq *AssetRoleQuery) First(ctx context.Context) (*AssetRole, error) {
	nodes, err := arq.Limit(1).All(setContextOp(ctx, arq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{assetrole.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (arq *AssetRoleQuery) FirstX(ctx context.Context) *AssetRole {
	node, err := arq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first AssetRole ID from the query.
// Returns a *NotFoundError when no AssetRole ID was found.
func (arq *AssetRoleQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = arq.Limit(1).IDs(setContextOp(ctx, arq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{assetrole.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (arq *AssetRoleQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := arq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single AssetRole entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one AssetRole entity is found.
// Returns a *NotFoundError when no AssetRole entities are found.
func (arq *AssetRoleQuery) Only(ctx context.Context) (*AssetRole, error) {
	nodes, err := arq.Limit(2).All(setContextOp(ctx, arq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{assetrole.Label}
	default:
		return nil, &NotSingularError{assetrole.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (arq *AssetRoleQuery) OnlyX(ctx context.Context) *AssetRole {
	node, err := arq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only AssetRole ID in the query.
// Returns a *NotSingularError when more than one AssetRole ID is found.
// Returns a *NotFoundError when no entities are found.
func (arq *AssetRoleQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = arq.Limit(2).IDs(setContextOp(ctx, arq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{assetrole.Label}
	default:
		err = &NotSingularError{assetrole.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (arq *AssetRoleQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := arq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AssetRoles.
func (arq *AssetRoleQuery) All(ctx context.Context) ([]*AssetRole, error) {
	ctx = setContextOp(ctx, arq.ctx, ent.OpQueryAll)
	if err := arq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*AssetRole, *AssetRoleQuery]()
	return withInterceptors[[]*AssetRole](ctx, arq, qr, arq.inters)
}

// AllX is like All, but panics if an error occurs.
func (arq *AssetRoleQuery) AllX(ctx context.Context) []*AssetRole {
	nodes, err := arq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of AssetRole IDs.
func (arq *AssetRoleQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if arq.ctx.Unique == nil && arq.path != nil {
		arq.Unique(true)
	}
	ctx = setContextOp(ctx, arq.ctx, ent.OpQueryIDs)
	if err = arq.Select(assetrole.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (arq *AssetRoleQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := arq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (arq *AssetRoleQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, arq.ctx, ent.OpQueryCount)
	if err := arq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, arq, querierCount[*AssetRoleQuery](), arq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (arq *AssetRoleQuery) CountX(ctx context.Context) int {
	count, err := arq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (arq *AssetRoleQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, arq.ctx, ent.OpQueryExist)
	switch _, err := arq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (arq *AssetRoleQuery) ExistX(ctx context.Context) bool {
	exist, err := arq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AssetRoleQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (arq *AssetRoleQuery) Clone() *AssetRoleQuery {
	if arq == nil {
		return nil
	}
	return &AssetRoleQuery{
		config:            arq.config,
		ctx:               arq.ctx.Clone(),
		order:             append([]assetrole.OrderOption{}, arq.order...),
		inters:            append([]Interceptor{}, arq.inters...),
		predicates:        append([]predicate.AssetRole{}, arq.predicates...),
		withAssetManagers: arq.withAssetManagers.Clone(),
		// clone intermediate query.
		sql:       arq.sql.Clone(),
		path:      arq.path,
		modifiers: append([]func(*sql.Selector){}, arq.modifiers...),
	}
}

// WithAssetManagers tells the query-builder to eager-load the nodes that are connected to
// the "asset_managers" edge. The optional arguments are used to configure the query builder of the edge.
func (arq *AssetRoleQuery) WithAssetManagers(opts ...func(*AssetManagerQuery)) *AssetRoleQuery {
	query := (&AssetManagerClient{config: arq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	arq.withAssetManagers = query
	return arq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.AssetRole.Query().
//		GroupBy(assetrole.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (arq *AssetRoleQuery) GroupBy(field string, fields ...string) *AssetRoleGroupBy {
	arq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &AssetRoleGroupBy{build: arq}
	grbuild.flds = &arq.ctx.Fields
	grbuild.label = assetrole.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.AssetRole.Query().
//		Select(assetrole.FieldName).
//		Scan(ctx, &v)
func (arq *AssetRoleQuery) Select(fields ...string) *AssetRoleSelect {
	arq.ctx.Fields = append(arq.ctx.Fields, fields...)
	sbuild := &AssetRoleSelect{AssetRoleQuery: arq}
	sbuild.label = assetrole.Label
	sbuild.flds, sbuild.scan = &arq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a AssetRoleSelect configured with the given aggregations.
func (arq *AssetRoleQuery) Aggregate(fns ...AggregateFunc) *AssetRoleSelect {
	return arq.Select().Aggregate(fns...)
}

func (arq *AssetRoleQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range arq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, arq); err != nil {
				return err
			}
		}
	}
	for _, f := range arq.ctx.Fields {
		if !assetrole.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if arq.path != nil {
		prev, err := arq.path(ctx)
		if err != nil {
			return err
		}
		arq.sql = prev
	}
	return nil
}

func (arq *AssetRoleQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*AssetRole, error) {
	var (
		nodes       = []*AssetRole{}
		_spec       = arq.querySpec()
		loadedTypes = [1]bool{
			arq.withAssetManagers != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*AssetRole).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &AssetRole{config: arq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(arq.modifiers) > 0 {
		_spec.Modifiers = arq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, arq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := arq.withAssetManagers; query != nil {
		if err := arq.loadAssetManagers(ctx, query, nodes,
			func(n *AssetRole) { n.Edges.AssetManagers = []*AssetManager{} },
			func(n *AssetRole, e *AssetManager) { n.Edges.AssetManagers = append(n.Edges.AssetManagers, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (arq *AssetRoleQuery) loadAssetManagers(ctx context.Context, query *AssetManagerQuery, nodes []*AssetRole, init func(*AssetRole), assign func(*AssetRole, *AssetManager)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uint64]*AssetRole)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(assetmanager.FieldRoleID)
	}
	query.Where(predicate.AssetManager(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(assetrole.AssetManagersColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.RoleID
		if fk == nil {
			return fmt.Errorf(`foreign-key "role_id" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "role_id" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (arq *AssetRoleQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := arq.querySpec()
	if len(arq.modifiers) > 0 {
		_spec.Modifiers = arq.modifiers
	}
	_spec.Node.Columns = arq.ctx.Fields
	if len(arq.ctx.Fields) > 0 {
		_spec.Unique = arq.ctx.Unique != nil && *arq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, arq.driver, _spec)
}

func (arq *AssetRoleQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(assetrole.Table, assetrole.Columns, sqlgraph.NewFieldSpec(assetrole.FieldID, field.TypeUint64))
	_spec.From = arq.sql
	if unique := arq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if arq.path != nil {
		_spec.Unique = true
	}
	if fields := arq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, assetrole.FieldID)
		for i := range fields {
			if fields[i] != assetrole.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := arq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := arq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := arq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := arq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (arq *AssetRoleQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(arq.driver.Dialect())
	t1 := builder.Table(assetrole.Table)
	columns := arq.ctx.Fields
	if len(columns) == 0 {
		columns = assetrole.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if arq.sql != nil {
		selector = arq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if arq.ctx.Unique != nil && *arq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range arq.modifiers {
		m(selector)
	}
	for _, p := range arq.predicates {
		p(selector)
	}
	for _, p := range arq.order {
		p(selector)
	}
	if offset := arq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := arq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (arq *AssetRoleQuery) Modify(modifiers ...func(s *sql.Selector)) *AssetRoleSelect {
	arq.modifiers = append(arq.modifiers, modifiers...)
	return arq.Select()
}

type AssetRoleQueryWith string

var (
	AssetRoleQueryWithAssetManagers AssetRoleQueryWith = "AssetManagers"
)

func (arq *AssetRoleQuery) With(withEdges ...AssetRoleQueryWith) *AssetRoleQuery {
	for _, v := range withEdges {
		switch v {
		case AssetRoleQueryWithAssetManagers:
			arq.WithAssetManagers()
		}
	}
	return arq
}

// AssetRoleGroupBy is the group-by builder for AssetRole entities.
type AssetRoleGroupBy struct {
	selector
	build *AssetRoleQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (argb *AssetRoleGroupBy) Aggregate(fns ...AggregateFunc) *AssetRoleGroupBy {
	argb.fns = append(argb.fns, fns...)
	return argb
}

// Scan applies the selector query and scans the result into the given value.
func (argb *AssetRoleGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, argb.build.ctx, ent.OpQueryGroupBy)
	if err := argb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AssetRoleQuery, *AssetRoleGroupBy](ctx, argb.build, argb, argb.build.inters, v)
}

func (argb *AssetRoleGroupBy) sqlScan(ctx context.Context, root *AssetRoleQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(argb.fns))
	for _, fn := range argb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*argb.flds)+len(argb.fns))
		for _, f := range *argb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*argb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := argb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// AssetRoleSelect is the builder for selecting fields of AssetRole entities.
type AssetRoleSelect struct {
	*AssetRoleQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ars *AssetRoleSelect) Aggregate(fns ...AggregateFunc) *AssetRoleSelect {
	ars.fns = append(ars.fns, fns...)
	return ars
}

// Scan applies the selector query and scans the result into the given value.
func (ars *AssetRoleSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ars.ctx, ent.OpQuerySelect)
	if err := ars.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AssetRoleQuery, *AssetRoleSelect](ctx, ars.AssetRoleQuery, ars, ars.inters, v)
}

func (ars *AssetRoleSelect) sqlScan(ctx context.Context, root *AssetRoleQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ars.fns))
	for _, fn := range ars.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ars.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ars.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ars *AssetRoleSelect) Modify(modifiers ...func(s *sql.Selector)) *AssetRoleSelect {
	ars.modifiers = append(ars.modifiers, modifiers...)
	return ars
}