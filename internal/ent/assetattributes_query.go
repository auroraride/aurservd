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
	"github.com/auroraride/aurservd/internal/ent/assetattributes"
	"github.com/auroraride/aurservd/internal/ent/assetattributevalues"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// AssetAttributesQuery is the builder for querying AssetAttributes entities.
type AssetAttributesQuery struct {
	config
	ctx        *QueryContext
	order      []assetattributes.OrderOption
	inters     []Interceptor
	predicates []predicate.AssetAttributes
	withValues *AssetAttributeValuesQuery
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AssetAttributesQuery builder.
func (aaq *AssetAttributesQuery) Where(ps ...predicate.AssetAttributes) *AssetAttributesQuery {
	aaq.predicates = append(aaq.predicates, ps...)
	return aaq
}

// Limit the number of records to be returned by this query.
func (aaq *AssetAttributesQuery) Limit(limit int) *AssetAttributesQuery {
	aaq.ctx.Limit = &limit
	return aaq
}

// Offset to start from.
func (aaq *AssetAttributesQuery) Offset(offset int) *AssetAttributesQuery {
	aaq.ctx.Offset = &offset
	return aaq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (aaq *AssetAttributesQuery) Unique(unique bool) *AssetAttributesQuery {
	aaq.ctx.Unique = &unique
	return aaq
}

// Order specifies how the records should be ordered.
func (aaq *AssetAttributesQuery) Order(o ...assetattributes.OrderOption) *AssetAttributesQuery {
	aaq.order = append(aaq.order, o...)
	return aaq
}

// QueryValues chains the current query on the "values" edge.
func (aaq *AssetAttributesQuery) QueryValues() *AssetAttributeValuesQuery {
	query := (&AssetAttributeValuesClient{config: aaq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aaq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aaq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetattributes.Table, assetattributes.FieldID, selector),
			sqlgraph.To(assetattributevalues.Table, assetattributevalues.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, assetattributes.ValuesTable, assetattributes.ValuesColumn),
		)
		fromU = sqlgraph.SetNeighbors(aaq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first AssetAttributes entity from the query.
// Returns a *NotFoundError when no AssetAttributes was found.
func (aaq *AssetAttributesQuery) First(ctx context.Context) (*AssetAttributes, error) {
	nodes, err := aaq.Limit(1).All(setContextOp(ctx, aaq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{assetattributes.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (aaq *AssetAttributesQuery) FirstX(ctx context.Context) *AssetAttributes {
	node, err := aaq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first AssetAttributes ID from the query.
// Returns a *NotFoundError when no AssetAttributes ID was found.
func (aaq *AssetAttributesQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = aaq.Limit(1).IDs(setContextOp(ctx, aaq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{assetattributes.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (aaq *AssetAttributesQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := aaq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single AssetAttributes entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one AssetAttributes entity is found.
// Returns a *NotFoundError when no AssetAttributes entities are found.
func (aaq *AssetAttributesQuery) Only(ctx context.Context) (*AssetAttributes, error) {
	nodes, err := aaq.Limit(2).All(setContextOp(ctx, aaq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{assetattributes.Label}
	default:
		return nil, &NotSingularError{assetattributes.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (aaq *AssetAttributesQuery) OnlyX(ctx context.Context) *AssetAttributes {
	node, err := aaq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only AssetAttributes ID in the query.
// Returns a *NotSingularError when more than one AssetAttributes ID is found.
// Returns a *NotFoundError when no entities are found.
func (aaq *AssetAttributesQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = aaq.Limit(2).IDs(setContextOp(ctx, aaq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{assetattributes.Label}
	default:
		err = &NotSingularError{assetattributes.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (aaq *AssetAttributesQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := aaq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AssetAttributesSlice.
func (aaq *AssetAttributesQuery) All(ctx context.Context) ([]*AssetAttributes, error) {
	ctx = setContextOp(ctx, aaq.ctx, ent.OpQueryAll)
	if err := aaq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*AssetAttributes, *AssetAttributesQuery]()
	return withInterceptors[[]*AssetAttributes](ctx, aaq, qr, aaq.inters)
}

// AllX is like All, but panics if an error occurs.
func (aaq *AssetAttributesQuery) AllX(ctx context.Context) []*AssetAttributes {
	nodes, err := aaq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of AssetAttributes IDs.
func (aaq *AssetAttributesQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if aaq.ctx.Unique == nil && aaq.path != nil {
		aaq.Unique(true)
	}
	ctx = setContextOp(ctx, aaq.ctx, ent.OpQueryIDs)
	if err = aaq.Select(assetattributes.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (aaq *AssetAttributesQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := aaq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (aaq *AssetAttributesQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, aaq.ctx, ent.OpQueryCount)
	if err := aaq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, aaq, querierCount[*AssetAttributesQuery](), aaq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (aaq *AssetAttributesQuery) CountX(ctx context.Context) int {
	count, err := aaq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (aaq *AssetAttributesQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, aaq.ctx, ent.OpQueryExist)
	switch _, err := aaq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (aaq *AssetAttributesQuery) ExistX(ctx context.Context) bool {
	exist, err := aaq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AssetAttributesQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (aaq *AssetAttributesQuery) Clone() *AssetAttributesQuery {
	if aaq == nil {
		return nil
	}
	return &AssetAttributesQuery{
		config:     aaq.config,
		ctx:        aaq.ctx.Clone(),
		order:      append([]assetattributes.OrderOption{}, aaq.order...),
		inters:     append([]Interceptor{}, aaq.inters...),
		predicates: append([]predicate.AssetAttributes{}, aaq.predicates...),
		withValues: aaq.withValues.Clone(),
		// clone intermediate query.
		sql:       aaq.sql.Clone(),
		path:      aaq.path,
		modifiers: append([]func(*sql.Selector){}, aaq.modifiers...),
	}
}

// WithValues tells the query-builder to eager-load the nodes that are connected to
// the "values" edge. The optional arguments are used to configure the query builder of the edge.
func (aaq *AssetAttributesQuery) WithValues(opts ...func(*AssetAttributeValuesQuery)) *AssetAttributesQuery {
	query := (&AssetAttributeValuesClient{config: aaq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	aaq.withValues = query
	return aaq
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
//	client.AssetAttributes.Query().
//		GroupBy(assetattributes.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (aaq *AssetAttributesQuery) GroupBy(field string, fields ...string) *AssetAttributesGroupBy {
	aaq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &AssetAttributesGroupBy{build: aaq}
	grbuild.flds = &aaq.ctx.Fields
	grbuild.label = assetattributes.Label
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
//	client.AssetAttributes.Query().
//		Select(assetattributes.FieldCreatedAt).
//		Scan(ctx, &v)
func (aaq *AssetAttributesQuery) Select(fields ...string) *AssetAttributesSelect {
	aaq.ctx.Fields = append(aaq.ctx.Fields, fields...)
	sbuild := &AssetAttributesSelect{AssetAttributesQuery: aaq}
	sbuild.label = assetattributes.Label
	sbuild.flds, sbuild.scan = &aaq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a AssetAttributesSelect configured with the given aggregations.
func (aaq *AssetAttributesQuery) Aggregate(fns ...AggregateFunc) *AssetAttributesSelect {
	return aaq.Select().Aggregate(fns...)
}

func (aaq *AssetAttributesQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range aaq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, aaq); err != nil {
				return err
			}
		}
	}
	for _, f := range aaq.ctx.Fields {
		if !assetattributes.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if aaq.path != nil {
		prev, err := aaq.path(ctx)
		if err != nil {
			return err
		}
		aaq.sql = prev
	}
	return nil
}

func (aaq *AssetAttributesQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*AssetAttributes, error) {
	var (
		nodes       = []*AssetAttributes{}
		_spec       = aaq.querySpec()
		loadedTypes = [1]bool{
			aaq.withValues != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*AssetAttributes).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &AssetAttributes{config: aaq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(aaq.modifiers) > 0 {
		_spec.Modifiers = aaq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, aaq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := aaq.withValues; query != nil {
		if err := aaq.loadValues(ctx, query, nodes,
			func(n *AssetAttributes) { n.Edges.Values = []*AssetAttributeValues{} },
			func(n *AssetAttributes, e *AssetAttributeValues) { n.Edges.Values = append(n.Edges.Values, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (aaq *AssetAttributesQuery) loadValues(ctx context.Context, query *AssetAttributeValuesQuery, nodes []*AssetAttributes, init func(*AssetAttributes), assign func(*AssetAttributes, *AssetAttributeValues)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uint64]*AssetAttributes)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(assetattributevalues.FieldAttributeID)
	}
	query.Where(predicate.AssetAttributeValues(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(assetattributes.ValuesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.AttributeID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "attribute_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (aaq *AssetAttributesQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := aaq.querySpec()
	if len(aaq.modifiers) > 0 {
		_spec.Modifiers = aaq.modifiers
	}
	_spec.Node.Columns = aaq.ctx.Fields
	if len(aaq.ctx.Fields) > 0 {
		_spec.Unique = aaq.ctx.Unique != nil && *aaq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, aaq.driver, _spec)
}

func (aaq *AssetAttributesQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(assetattributes.Table, assetattributes.Columns, sqlgraph.NewFieldSpec(assetattributes.FieldID, field.TypeUint64))
	_spec.From = aaq.sql
	if unique := aaq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if aaq.path != nil {
		_spec.Unique = true
	}
	if fields := aaq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, assetattributes.FieldID)
		for i := range fields {
			if fields[i] != assetattributes.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := aaq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := aaq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := aaq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := aaq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (aaq *AssetAttributesQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(aaq.driver.Dialect())
	t1 := builder.Table(assetattributes.Table)
	columns := aaq.ctx.Fields
	if len(columns) == 0 {
		columns = assetattributes.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if aaq.sql != nil {
		selector = aaq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if aaq.ctx.Unique != nil && *aaq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range aaq.modifiers {
		m(selector)
	}
	for _, p := range aaq.predicates {
		p(selector)
	}
	for _, p := range aaq.order {
		p(selector)
	}
	if offset := aaq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := aaq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (aaq *AssetAttributesQuery) Modify(modifiers ...func(s *sql.Selector)) *AssetAttributesSelect {
	aaq.modifiers = append(aaq.modifiers, modifiers...)
	return aaq.Select()
}

type AssetAttributesQueryWith string

var (
	AssetAttributesQueryWithValues AssetAttributesQueryWith = "Values"
)

func (aaq *AssetAttributesQuery) With(withEdges ...AssetAttributesQueryWith) *AssetAttributesQuery {
	for _, v := range withEdges {
		switch v {
		case AssetAttributesQueryWithValues:
			aaq.WithValues()
		}
	}
	return aaq
}

// AssetAttributesGroupBy is the group-by builder for AssetAttributes entities.
type AssetAttributesGroupBy struct {
	selector
	build *AssetAttributesQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (aagb *AssetAttributesGroupBy) Aggregate(fns ...AggregateFunc) *AssetAttributesGroupBy {
	aagb.fns = append(aagb.fns, fns...)
	return aagb
}

// Scan applies the selector query and scans the result into the given value.
func (aagb *AssetAttributesGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, aagb.build.ctx, ent.OpQueryGroupBy)
	if err := aagb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AssetAttributesQuery, *AssetAttributesGroupBy](ctx, aagb.build, aagb, aagb.build.inters, v)
}

func (aagb *AssetAttributesGroupBy) sqlScan(ctx context.Context, root *AssetAttributesQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(aagb.fns))
	for _, fn := range aagb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*aagb.flds)+len(aagb.fns))
		for _, f := range *aagb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*aagb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := aagb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// AssetAttributesSelect is the builder for selecting fields of AssetAttributes entities.
type AssetAttributesSelect struct {
	*AssetAttributesQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (aas *AssetAttributesSelect) Aggregate(fns ...AggregateFunc) *AssetAttributesSelect {
	aas.fns = append(aas.fns, fns...)
	return aas
}

// Scan applies the selector query and scans the result into the given value.
func (aas *AssetAttributesSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, aas.ctx, ent.OpQuerySelect)
	if err := aas.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AssetAttributesQuery, *AssetAttributesSelect](ctx, aas.AssetAttributesQuery, aas, aas.inters, v)
}

func (aas *AssetAttributesSelect) sqlScan(ctx context.Context, root *AssetAttributesQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(aas.fns))
	for _, fn := range aas.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*aas.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := aas.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (aas *AssetAttributesSelect) Modify(modifiers ...func(s *sql.Selector)) *AssetAttributesSelect {
	aas.modifiers = append(aas.modifiers, modifiers...)
	return aas
}