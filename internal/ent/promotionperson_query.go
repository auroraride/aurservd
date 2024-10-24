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
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/promotionperson"
)

// PromotionPersonQuery is the builder for querying PromotionPerson entities.
type PromotionPersonQuery struct {
	config
	ctx        *QueryContext
	order      []promotionperson.OrderOption
	inters     []Interceptor
	predicates []predicate.PromotionPerson
	withMember *PromotionMemberQuery
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PromotionPersonQuery builder.
func (ppq *PromotionPersonQuery) Where(ps ...predicate.PromotionPerson) *PromotionPersonQuery {
	ppq.predicates = append(ppq.predicates, ps...)
	return ppq
}

// Limit the number of records to be returned by this query.
func (ppq *PromotionPersonQuery) Limit(limit int) *PromotionPersonQuery {
	ppq.ctx.Limit = &limit
	return ppq
}

// Offset to start from.
func (ppq *PromotionPersonQuery) Offset(offset int) *PromotionPersonQuery {
	ppq.ctx.Offset = &offset
	return ppq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (ppq *PromotionPersonQuery) Unique(unique bool) *PromotionPersonQuery {
	ppq.ctx.Unique = &unique
	return ppq
}

// Order specifies how the records should be ordered.
func (ppq *PromotionPersonQuery) Order(o ...promotionperson.OrderOption) *PromotionPersonQuery {
	ppq.order = append(ppq.order, o...)
	return ppq
}

// QueryMember chains the current query on the "member" edge.
func (ppq *PromotionPersonQuery) QueryMember() *PromotionMemberQuery {
	query := (&PromotionMemberClient{config: ppq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ppq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := ppq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(promotionperson.Table, promotionperson.FieldID, selector),
			sqlgraph.To(promotionmember.Table, promotionmember.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, promotionperson.MemberTable, promotionperson.MemberColumn),
		)
		fromU = sqlgraph.SetNeighbors(ppq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first PromotionPerson entity from the query.
// Returns a *NotFoundError when no PromotionPerson was found.
func (ppq *PromotionPersonQuery) First(ctx context.Context) (*PromotionPerson, error) {
	nodes, err := ppq.Limit(1).All(setContextOp(ctx, ppq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{promotionperson.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ppq *PromotionPersonQuery) FirstX(ctx context.Context) *PromotionPerson {
	node, err := ppq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first PromotionPerson ID from the query.
// Returns a *NotFoundError when no PromotionPerson ID was found.
func (ppq *PromotionPersonQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = ppq.Limit(1).IDs(setContextOp(ctx, ppq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{promotionperson.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (ppq *PromotionPersonQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := ppq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single PromotionPerson entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one PromotionPerson entity is found.
// Returns a *NotFoundError when no PromotionPerson entities are found.
func (ppq *PromotionPersonQuery) Only(ctx context.Context) (*PromotionPerson, error) {
	nodes, err := ppq.Limit(2).All(setContextOp(ctx, ppq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{promotionperson.Label}
	default:
		return nil, &NotSingularError{promotionperson.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ppq *PromotionPersonQuery) OnlyX(ctx context.Context) *PromotionPerson {
	node, err := ppq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only PromotionPerson ID in the query.
// Returns a *NotSingularError when more than one PromotionPerson ID is found.
// Returns a *NotFoundError when no entities are found.
func (ppq *PromotionPersonQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = ppq.Limit(2).IDs(setContextOp(ctx, ppq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{promotionperson.Label}
	default:
		err = &NotSingularError{promotionperson.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (ppq *PromotionPersonQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := ppq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PromotionPersons.
func (ppq *PromotionPersonQuery) All(ctx context.Context) ([]*PromotionPerson, error) {
	ctx = setContextOp(ctx, ppq.ctx, ent.OpQueryAll)
	if err := ppq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*PromotionPerson, *PromotionPersonQuery]()
	return withInterceptors[[]*PromotionPerson](ctx, ppq, qr, ppq.inters)
}

// AllX is like All, but panics if an error occurs.
func (ppq *PromotionPersonQuery) AllX(ctx context.Context) []*PromotionPerson {
	nodes, err := ppq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of PromotionPerson IDs.
func (ppq *PromotionPersonQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if ppq.ctx.Unique == nil && ppq.path != nil {
		ppq.Unique(true)
	}
	ctx = setContextOp(ctx, ppq.ctx, ent.OpQueryIDs)
	if err = ppq.Select(promotionperson.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ppq *PromotionPersonQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := ppq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ppq *PromotionPersonQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, ppq.ctx, ent.OpQueryCount)
	if err := ppq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, ppq, querierCount[*PromotionPersonQuery](), ppq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (ppq *PromotionPersonQuery) CountX(ctx context.Context) int {
	count, err := ppq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ppq *PromotionPersonQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, ppq.ctx, ent.OpQueryExist)
	switch _, err := ppq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (ppq *PromotionPersonQuery) ExistX(ctx context.Context) bool {
	exist, err := ppq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PromotionPersonQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ppq *PromotionPersonQuery) Clone() *PromotionPersonQuery {
	if ppq == nil {
		return nil
	}
	return &PromotionPersonQuery{
		config:     ppq.config,
		ctx:        ppq.ctx.Clone(),
		order:      append([]promotionperson.OrderOption{}, ppq.order...),
		inters:     append([]Interceptor{}, ppq.inters...),
		predicates: append([]predicate.PromotionPerson{}, ppq.predicates...),
		withMember: ppq.withMember.Clone(),
		// clone intermediate query.
		sql:       ppq.sql.Clone(),
		path:      ppq.path,
		modifiers: append([]func(*sql.Selector){}, ppq.modifiers...),
	}
}

// WithMember tells the query-builder to eager-load the nodes that are connected to
// the "member" edge. The optional arguments are used to configure the query builder of the edge.
func (ppq *PromotionPersonQuery) WithMember(opts ...func(*PromotionMemberQuery)) *PromotionPersonQuery {
	query := (&PromotionMemberClient{config: ppq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	ppq.withMember = query
	return ppq
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
//	client.PromotionPerson.Query().
//		GroupBy(promotionperson.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (ppq *PromotionPersonQuery) GroupBy(field string, fields ...string) *PromotionPersonGroupBy {
	ppq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &PromotionPersonGroupBy{build: ppq}
	grbuild.flds = &ppq.ctx.Fields
	grbuild.label = promotionperson.Label
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
//	client.PromotionPerson.Query().
//		Select(promotionperson.FieldCreatedAt).
//		Scan(ctx, &v)
func (ppq *PromotionPersonQuery) Select(fields ...string) *PromotionPersonSelect {
	ppq.ctx.Fields = append(ppq.ctx.Fields, fields...)
	sbuild := &PromotionPersonSelect{PromotionPersonQuery: ppq}
	sbuild.label = promotionperson.Label
	sbuild.flds, sbuild.scan = &ppq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a PromotionPersonSelect configured with the given aggregations.
func (ppq *PromotionPersonQuery) Aggregate(fns ...AggregateFunc) *PromotionPersonSelect {
	return ppq.Select().Aggregate(fns...)
}

func (ppq *PromotionPersonQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range ppq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, ppq); err != nil {
				return err
			}
		}
	}
	for _, f := range ppq.ctx.Fields {
		if !promotionperson.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if ppq.path != nil {
		prev, err := ppq.path(ctx)
		if err != nil {
			return err
		}
		ppq.sql = prev
	}
	return nil
}

func (ppq *PromotionPersonQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*PromotionPerson, error) {
	var (
		nodes       = []*PromotionPerson{}
		_spec       = ppq.querySpec()
		loadedTypes = [1]bool{
			ppq.withMember != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*PromotionPerson).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &PromotionPerson{config: ppq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(ppq.modifiers) > 0 {
		_spec.Modifiers = ppq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, ppq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := ppq.withMember; query != nil {
		if err := ppq.loadMember(ctx, query, nodes,
			func(n *PromotionPerson) { n.Edges.Member = []*PromotionMember{} },
			func(n *PromotionPerson, e *PromotionMember) { n.Edges.Member = append(n.Edges.Member, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (ppq *PromotionPersonQuery) loadMember(ctx context.Context, query *PromotionMemberQuery, nodes []*PromotionPerson, init func(*PromotionPerson), assign func(*PromotionPerson, *PromotionMember)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uint64]*PromotionPerson)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(promotionmember.FieldPersonID)
	}
	query.Where(predicate.PromotionMember(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(promotionperson.MemberColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.PersonID
		if fk == nil {
			return fmt.Errorf(`foreign-key "person_id" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "person_id" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (ppq *PromotionPersonQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ppq.querySpec()
	if len(ppq.modifiers) > 0 {
		_spec.Modifiers = ppq.modifiers
	}
	_spec.Node.Columns = ppq.ctx.Fields
	if len(ppq.ctx.Fields) > 0 {
		_spec.Unique = ppq.ctx.Unique != nil && *ppq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, ppq.driver, _spec)
}

func (ppq *PromotionPersonQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(promotionperson.Table, promotionperson.Columns, sqlgraph.NewFieldSpec(promotionperson.FieldID, field.TypeUint64))
	_spec.From = ppq.sql
	if unique := ppq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if ppq.path != nil {
		_spec.Unique = true
	}
	if fields := ppq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, promotionperson.FieldID)
		for i := range fields {
			if fields[i] != promotionperson.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := ppq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ppq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ppq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ppq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ppq *PromotionPersonQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(ppq.driver.Dialect())
	t1 := builder.Table(promotionperson.Table)
	columns := ppq.ctx.Fields
	if len(columns) == 0 {
		columns = promotionperson.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if ppq.sql != nil {
		selector = ppq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if ppq.ctx.Unique != nil && *ppq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range ppq.modifiers {
		m(selector)
	}
	for _, p := range ppq.predicates {
		p(selector)
	}
	for _, p := range ppq.order {
		p(selector)
	}
	if offset := ppq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ppq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ppq *PromotionPersonQuery) Modify(modifiers ...func(s *sql.Selector)) *PromotionPersonSelect {
	ppq.modifiers = append(ppq.modifiers, modifiers...)
	return ppq.Select()
}

type PromotionPersonQueryWith string

var (
	PromotionPersonQueryWithMember PromotionPersonQueryWith = "Member"
)

func (ppq *PromotionPersonQuery) With(withEdges ...PromotionPersonQueryWith) *PromotionPersonQuery {
	for _, v := range withEdges {
		switch v {
		case PromotionPersonQueryWithMember:
			ppq.WithMember()
		}
	}
	return ppq
}

// PromotionPersonGroupBy is the group-by builder for PromotionPerson entities.
type PromotionPersonGroupBy struct {
	selector
	build *PromotionPersonQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ppgb *PromotionPersonGroupBy) Aggregate(fns ...AggregateFunc) *PromotionPersonGroupBy {
	ppgb.fns = append(ppgb.fns, fns...)
	return ppgb
}

// Scan applies the selector query and scans the result into the given value.
func (ppgb *PromotionPersonGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ppgb.build.ctx, ent.OpQueryGroupBy)
	if err := ppgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PromotionPersonQuery, *PromotionPersonGroupBy](ctx, ppgb.build, ppgb, ppgb.build.inters, v)
}

func (ppgb *PromotionPersonGroupBy) sqlScan(ctx context.Context, root *PromotionPersonQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(ppgb.fns))
	for _, fn := range ppgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*ppgb.flds)+len(ppgb.fns))
		for _, f := range *ppgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*ppgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ppgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// PromotionPersonSelect is the builder for selecting fields of PromotionPerson entities.
type PromotionPersonSelect struct {
	*PromotionPersonQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (pps *PromotionPersonSelect) Aggregate(fns ...AggregateFunc) *PromotionPersonSelect {
	pps.fns = append(pps.fns, fns...)
	return pps
}

// Scan applies the selector query and scans the result into the given value.
func (pps *PromotionPersonSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pps.ctx, ent.OpQuerySelect)
	if err := pps.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PromotionPersonQuery, *PromotionPersonSelect](ctx, pps.PromotionPersonQuery, pps, pps.inters, v)
}

func (pps *PromotionPersonSelect) sqlScan(ctx context.Context, root *PromotionPersonQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(pps.fns))
	for _, fn := range pps.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*pps.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (pps *PromotionPersonSelect) Modify(modifiers ...func(s *sql.Selector)) *PromotionPersonSelect {
	pps.modifiers = append(pps.modifiers, modifiers...)
	return pps
}
