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
	"github.com/auroraride/aurservd/internal/ent/promotionbankcard"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/promotionwithdrawal"
)

// PromotionBankCardQuery is the builder for querying PromotionBankCard entities.
type PromotionBankCardQuery struct {
	config
	ctx             *QueryContext
	order           []promotionbankcard.OrderOption
	inters          []Interceptor
	predicates      []predicate.PromotionBankCard
	withMember      *PromotionMemberQuery
	withWithdrawals *PromotionWithdrawalQuery
	modifiers       []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PromotionBankCardQuery builder.
func (pbcq *PromotionBankCardQuery) Where(ps ...predicate.PromotionBankCard) *PromotionBankCardQuery {
	pbcq.predicates = append(pbcq.predicates, ps...)
	return pbcq
}

// Limit the number of records to be returned by this query.
func (pbcq *PromotionBankCardQuery) Limit(limit int) *PromotionBankCardQuery {
	pbcq.ctx.Limit = &limit
	return pbcq
}

// Offset to start from.
func (pbcq *PromotionBankCardQuery) Offset(offset int) *PromotionBankCardQuery {
	pbcq.ctx.Offset = &offset
	return pbcq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (pbcq *PromotionBankCardQuery) Unique(unique bool) *PromotionBankCardQuery {
	pbcq.ctx.Unique = &unique
	return pbcq
}

// Order specifies how the records should be ordered.
func (pbcq *PromotionBankCardQuery) Order(o ...promotionbankcard.OrderOption) *PromotionBankCardQuery {
	pbcq.order = append(pbcq.order, o...)
	return pbcq
}

// QueryMember chains the current query on the "member" edge.
func (pbcq *PromotionBankCardQuery) QueryMember() *PromotionMemberQuery {
	query := (&PromotionMemberClient{config: pbcq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pbcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pbcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(promotionbankcard.Table, promotionbankcard.FieldID, selector),
			sqlgraph.To(promotionmember.Table, promotionmember.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, promotionbankcard.MemberTable, promotionbankcard.MemberColumn),
		)
		fromU = sqlgraph.SetNeighbors(pbcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryWithdrawals chains the current query on the "withdrawals" edge.
func (pbcq *PromotionBankCardQuery) QueryWithdrawals() *PromotionWithdrawalQuery {
	query := (&PromotionWithdrawalClient{config: pbcq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pbcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pbcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(promotionbankcard.Table, promotionbankcard.FieldID, selector),
			sqlgraph.To(promotionwithdrawal.Table, promotionwithdrawal.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, promotionbankcard.WithdrawalsTable, promotionbankcard.WithdrawalsColumn),
		)
		fromU = sqlgraph.SetNeighbors(pbcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first PromotionBankCard entity from the query.
// Returns a *NotFoundError when no PromotionBankCard was found.
func (pbcq *PromotionBankCardQuery) First(ctx context.Context) (*PromotionBankCard, error) {
	nodes, err := pbcq.Limit(1).All(setContextOp(ctx, pbcq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{promotionbankcard.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pbcq *PromotionBankCardQuery) FirstX(ctx context.Context) *PromotionBankCard {
	node, err := pbcq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first PromotionBankCard ID from the query.
// Returns a *NotFoundError when no PromotionBankCard ID was found.
func (pbcq *PromotionBankCardQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = pbcq.Limit(1).IDs(setContextOp(ctx, pbcq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{promotionbankcard.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (pbcq *PromotionBankCardQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := pbcq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single PromotionBankCard entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one PromotionBankCard entity is found.
// Returns a *NotFoundError when no PromotionBankCard entities are found.
func (pbcq *PromotionBankCardQuery) Only(ctx context.Context) (*PromotionBankCard, error) {
	nodes, err := pbcq.Limit(2).All(setContextOp(ctx, pbcq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{promotionbankcard.Label}
	default:
		return nil, &NotSingularError{promotionbankcard.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pbcq *PromotionBankCardQuery) OnlyX(ctx context.Context) *PromotionBankCard {
	node, err := pbcq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only PromotionBankCard ID in the query.
// Returns a *NotSingularError when more than one PromotionBankCard ID is found.
// Returns a *NotFoundError when no entities are found.
func (pbcq *PromotionBankCardQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = pbcq.Limit(2).IDs(setContextOp(ctx, pbcq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{promotionbankcard.Label}
	default:
		err = &NotSingularError{promotionbankcard.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (pbcq *PromotionBankCardQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := pbcq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PromotionBankCards.
func (pbcq *PromotionBankCardQuery) All(ctx context.Context) ([]*PromotionBankCard, error) {
	ctx = setContextOp(ctx, pbcq.ctx, ent.OpQueryAll)
	if err := pbcq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*PromotionBankCard, *PromotionBankCardQuery]()
	return withInterceptors[[]*PromotionBankCard](ctx, pbcq, qr, pbcq.inters)
}

// AllX is like All, but panics if an error occurs.
func (pbcq *PromotionBankCardQuery) AllX(ctx context.Context) []*PromotionBankCard {
	nodes, err := pbcq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of PromotionBankCard IDs.
func (pbcq *PromotionBankCardQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if pbcq.ctx.Unique == nil && pbcq.path != nil {
		pbcq.Unique(true)
	}
	ctx = setContextOp(ctx, pbcq.ctx, ent.OpQueryIDs)
	if err = pbcq.Select(promotionbankcard.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pbcq *PromotionBankCardQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := pbcq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pbcq *PromotionBankCardQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, pbcq.ctx, ent.OpQueryCount)
	if err := pbcq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, pbcq, querierCount[*PromotionBankCardQuery](), pbcq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (pbcq *PromotionBankCardQuery) CountX(ctx context.Context) int {
	count, err := pbcq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pbcq *PromotionBankCardQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, pbcq.ctx, ent.OpQueryExist)
	switch _, err := pbcq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (pbcq *PromotionBankCardQuery) ExistX(ctx context.Context) bool {
	exist, err := pbcq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PromotionBankCardQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pbcq *PromotionBankCardQuery) Clone() *PromotionBankCardQuery {
	if pbcq == nil {
		return nil
	}
	return &PromotionBankCardQuery{
		config:          pbcq.config,
		ctx:             pbcq.ctx.Clone(),
		order:           append([]promotionbankcard.OrderOption{}, pbcq.order...),
		inters:          append([]Interceptor{}, pbcq.inters...),
		predicates:      append([]predicate.PromotionBankCard{}, pbcq.predicates...),
		withMember:      pbcq.withMember.Clone(),
		withWithdrawals: pbcq.withWithdrawals.Clone(),
		// clone intermediate query.
		sql:       pbcq.sql.Clone(),
		path:      pbcq.path,
		modifiers: append([]func(*sql.Selector){}, pbcq.modifiers...),
	}
}

// WithMember tells the query-builder to eager-load the nodes that are connected to
// the "member" edge. The optional arguments are used to configure the query builder of the edge.
func (pbcq *PromotionBankCardQuery) WithMember(opts ...func(*PromotionMemberQuery)) *PromotionBankCardQuery {
	query := (&PromotionMemberClient{config: pbcq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	pbcq.withMember = query
	return pbcq
}

// WithWithdrawals tells the query-builder to eager-load the nodes that are connected to
// the "withdrawals" edge. The optional arguments are used to configure the query builder of the edge.
func (pbcq *PromotionBankCardQuery) WithWithdrawals(opts ...func(*PromotionWithdrawalQuery)) *PromotionBankCardQuery {
	query := (&PromotionWithdrawalClient{config: pbcq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	pbcq.withWithdrawals = query
	return pbcq
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
//	client.PromotionBankCard.Query().
//		GroupBy(promotionbankcard.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (pbcq *PromotionBankCardQuery) GroupBy(field string, fields ...string) *PromotionBankCardGroupBy {
	pbcq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &PromotionBankCardGroupBy{build: pbcq}
	grbuild.flds = &pbcq.ctx.Fields
	grbuild.label = promotionbankcard.Label
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
//	client.PromotionBankCard.Query().
//		Select(promotionbankcard.FieldCreatedAt).
//		Scan(ctx, &v)
func (pbcq *PromotionBankCardQuery) Select(fields ...string) *PromotionBankCardSelect {
	pbcq.ctx.Fields = append(pbcq.ctx.Fields, fields...)
	sbuild := &PromotionBankCardSelect{PromotionBankCardQuery: pbcq}
	sbuild.label = promotionbankcard.Label
	sbuild.flds, sbuild.scan = &pbcq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a PromotionBankCardSelect configured with the given aggregations.
func (pbcq *PromotionBankCardQuery) Aggregate(fns ...AggregateFunc) *PromotionBankCardSelect {
	return pbcq.Select().Aggregate(fns...)
}

func (pbcq *PromotionBankCardQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range pbcq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, pbcq); err != nil {
				return err
			}
		}
	}
	for _, f := range pbcq.ctx.Fields {
		if !promotionbankcard.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if pbcq.path != nil {
		prev, err := pbcq.path(ctx)
		if err != nil {
			return err
		}
		pbcq.sql = prev
	}
	return nil
}

func (pbcq *PromotionBankCardQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*PromotionBankCard, error) {
	var (
		nodes       = []*PromotionBankCard{}
		_spec       = pbcq.querySpec()
		loadedTypes = [2]bool{
			pbcq.withMember != nil,
			pbcq.withWithdrawals != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*PromotionBankCard).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &PromotionBankCard{config: pbcq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(pbcq.modifiers) > 0 {
		_spec.Modifiers = pbcq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, pbcq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := pbcq.withMember; query != nil {
		if err := pbcq.loadMember(ctx, query, nodes, nil,
			func(n *PromotionBankCard, e *PromotionMember) { n.Edges.Member = e }); err != nil {
			return nil, err
		}
	}
	if query := pbcq.withWithdrawals; query != nil {
		if err := pbcq.loadWithdrawals(ctx, query, nodes,
			func(n *PromotionBankCard) { n.Edges.Withdrawals = []*PromotionWithdrawal{} },
			func(n *PromotionBankCard, e *PromotionWithdrawal) {
				n.Edges.Withdrawals = append(n.Edges.Withdrawals, e)
			}); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (pbcq *PromotionBankCardQuery) loadMember(ctx context.Context, query *PromotionMemberQuery, nodes []*PromotionBankCard, init func(*PromotionBankCard), assign func(*PromotionBankCard, *PromotionMember)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*PromotionBankCard)
	for i := range nodes {
		fk := nodes[i].MemberID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(promotionmember.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "member_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (pbcq *PromotionBankCardQuery) loadWithdrawals(ctx context.Context, query *PromotionWithdrawalQuery, nodes []*PromotionBankCard, init func(*PromotionBankCard), assign func(*PromotionBankCard, *PromotionWithdrawal)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uint64]*PromotionBankCard)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(promotionwithdrawal.FieldAccountID)
	}
	query.Where(predicate.PromotionWithdrawal(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(promotionbankcard.WithdrawalsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.AccountID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "account_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (pbcq *PromotionBankCardQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := pbcq.querySpec()
	if len(pbcq.modifiers) > 0 {
		_spec.Modifiers = pbcq.modifiers
	}
	_spec.Node.Columns = pbcq.ctx.Fields
	if len(pbcq.ctx.Fields) > 0 {
		_spec.Unique = pbcq.ctx.Unique != nil && *pbcq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, pbcq.driver, _spec)
}

func (pbcq *PromotionBankCardQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(promotionbankcard.Table, promotionbankcard.Columns, sqlgraph.NewFieldSpec(promotionbankcard.FieldID, field.TypeUint64))
	_spec.From = pbcq.sql
	if unique := pbcq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if pbcq.path != nil {
		_spec.Unique = true
	}
	if fields := pbcq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, promotionbankcard.FieldID)
		for i := range fields {
			if fields[i] != promotionbankcard.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if pbcq.withMember != nil {
			_spec.Node.AddColumnOnce(promotionbankcard.FieldMemberID)
		}
	}
	if ps := pbcq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := pbcq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := pbcq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := pbcq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (pbcq *PromotionBankCardQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(pbcq.driver.Dialect())
	t1 := builder.Table(promotionbankcard.Table)
	columns := pbcq.ctx.Fields
	if len(columns) == 0 {
		columns = promotionbankcard.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if pbcq.sql != nil {
		selector = pbcq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if pbcq.ctx.Unique != nil && *pbcq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range pbcq.modifiers {
		m(selector)
	}
	for _, p := range pbcq.predicates {
		p(selector)
	}
	for _, p := range pbcq.order {
		p(selector)
	}
	if offset := pbcq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := pbcq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (pbcq *PromotionBankCardQuery) Modify(modifiers ...func(s *sql.Selector)) *PromotionBankCardSelect {
	pbcq.modifiers = append(pbcq.modifiers, modifiers...)
	return pbcq.Select()
}

type PromotionBankCardQueryWith string

var (
	PromotionBankCardQueryWithMember      PromotionBankCardQueryWith = "Member"
	PromotionBankCardQueryWithWithdrawals PromotionBankCardQueryWith = "Withdrawals"
)

func (pbcq *PromotionBankCardQuery) With(withEdges ...PromotionBankCardQueryWith) *PromotionBankCardQuery {
	for _, v := range withEdges {
		switch v {
		case PromotionBankCardQueryWithMember:
			pbcq.WithMember()
		case PromotionBankCardQueryWithWithdrawals:
			pbcq.WithWithdrawals()
		}
	}
	return pbcq
}

// PromotionBankCardGroupBy is the group-by builder for PromotionBankCard entities.
type PromotionBankCardGroupBy struct {
	selector
	build *PromotionBankCardQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pbcgb *PromotionBankCardGroupBy) Aggregate(fns ...AggregateFunc) *PromotionBankCardGroupBy {
	pbcgb.fns = append(pbcgb.fns, fns...)
	return pbcgb
}

// Scan applies the selector query and scans the result into the given value.
func (pbcgb *PromotionBankCardGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pbcgb.build.ctx, ent.OpQueryGroupBy)
	if err := pbcgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PromotionBankCardQuery, *PromotionBankCardGroupBy](ctx, pbcgb.build, pbcgb, pbcgb.build.inters, v)
}

func (pbcgb *PromotionBankCardGroupBy) sqlScan(ctx context.Context, root *PromotionBankCardQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(pbcgb.fns))
	for _, fn := range pbcgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*pbcgb.flds)+len(pbcgb.fns))
		for _, f := range *pbcgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*pbcgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pbcgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// PromotionBankCardSelect is the builder for selecting fields of PromotionBankCard entities.
type PromotionBankCardSelect struct {
	*PromotionBankCardQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (pbcs *PromotionBankCardSelect) Aggregate(fns ...AggregateFunc) *PromotionBankCardSelect {
	pbcs.fns = append(pbcs.fns, fns...)
	return pbcs
}

// Scan applies the selector query and scans the result into the given value.
func (pbcs *PromotionBankCardSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pbcs.ctx, ent.OpQuerySelect)
	if err := pbcs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PromotionBankCardQuery, *PromotionBankCardSelect](ctx, pbcs.PromotionBankCardQuery, pbcs, pbcs.inters, v)
}

func (pbcs *PromotionBankCardSelect) sqlScan(ctx context.Context, root *PromotionBankCardQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(pbcs.fns))
	for _, fn := range pbcs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*pbcs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pbcs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (pbcs *PromotionBankCardSelect) Modify(modifiers ...func(s *sql.Selector)) *PromotionBankCardSelect {
	pbcs.modifiers = append(pbcs.modifiers, modifiers...)
	return pbcs
}
