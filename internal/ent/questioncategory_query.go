// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/question"
	"github.com/auroraride/aurservd/internal/ent/questioncategory"
)

// QuestionCategoryQuery is the builder for querying QuestionCategory entities.
type QuestionCategoryQuery struct {
	config
	ctx           *QueryContext
	order         []questioncategory.OrderOption
	inters        []Interceptor
	predicates    []predicate.QuestionCategory
	withQuestions *QuestionQuery
	modifiers     []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the QuestionCategoryQuery builder.
func (qcq *QuestionCategoryQuery) Where(ps ...predicate.QuestionCategory) *QuestionCategoryQuery {
	qcq.predicates = append(qcq.predicates, ps...)
	return qcq
}

// Limit the number of records to be returned by this query.
func (qcq *QuestionCategoryQuery) Limit(limit int) *QuestionCategoryQuery {
	qcq.ctx.Limit = &limit
	return qcq
}

// Offset to start from.
func (qcq *QuestionCategoryQuery) Offset(offset int) *QuestionCategoryQuery {
	qcq.ctx.Offset = &offset
	return qcq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (qcq *QuestionCategoryQuery) Unique(unique bool) *QuestionCategoryQuery {
	qcq.ctx.Unique = &unique
	return qcq
}

// Order specifies how the records should be ordered.
func (qcq *QuestionCategoryQuery) Order(o ...questioncategory.OrderOption) *QuestionCategoryQuery {
	qcq.order = append(qcq.order, o...)
	return qcq
}

// QueryQuestions chains the current query on the "questions" edge.
func (qcq *QuestionCategoryQuery) QueryQuestions() *QuestionQuery {
	query := (&QuestionClient{config: qcq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := qcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := qcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(questioncategory.Table, questioncategory.FieldID, selector),
			sqlgraph.To(question.Table, question.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, questioncategory.QuestionsTable, questioncategory.QuestionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(qcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first QuestionCategory entity from the query.
// Returns a *NotFoundError when no QuestionCategory was found.
func (qcq *QuestionCategoryQuery) First(ctx context.Context) (*QuestionCategory, error) {
	nodes, err := qcq.Limit(1).All(setContextOp(ctx, qcq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{questioncategory.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (qcq *QuestionCategoryQuery) FirstX(ctx context.Context) *QuestionCategory {
	node, err := qcq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first QuestionCategory ID from the query.
// Returns a *NotFoundError when no QuestionCategory ID was found.
func (qcq *QuestionCategoryQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = qcq.Limit(1).IDs(setContextOp(ctx, qcq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{questioncategory.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (qcq *QuestionCategoryQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := qcq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single QuestionCategory entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one QuestionCategory entity is found.
// Returns a *NotFoundError when no QuestionCategory entities are found.
func (qcq *QuestionCategoryQuery) Only(ctx context.Context) (*QuestionCategory, error) {
	nodes, err := qcq.Limit(2).All(setContextOp(ctx, qcq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{questioncategory.Label}
	default:
		return nil, &NotSingularError{questioncategory.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (qcq *QuestionCategoryQuery) OnlyX(ctx context.Context) *QuestionCategory {
	node, err := qcq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only QuestionCategory ID in the query.
// Returns a *NotSingularError when more than one QuestionCategory ID is found.
// Returns a *NotFoundError when no entities are found.
func (qcq *QuestionCategoryQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = qcq.Limit(2).IDs(setContextOp(ctx, qcq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{questioncategory.Label}
	default:
		err = &NotSingularError{questioncategory.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (qcq *QuestionCategoryQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := qcq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of QuestionCategories.
func (qcq *QuestionCategoryQuery) All(ctx context.Context) ([]*QuestionCategory, error) {
	ctx = setContextOp(ctx, qcq.ctx, "All")
	if err := qcq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*QuestionCategory, *QuestionCategoryQuery]()
	return withInterceptors[[]*QuestionCategory](ctx, qcq, qr, qcq.inters)
}

// AllX is like All, but panics if an error occurs.
func (qcq *QuestionCategoryQuery) AllX(ctx context.Context) []*QuestionCategory {
	nodes, err := qcq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of QuestionCategory IDs.
func (qcq *QuestionCategoryQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if qcq.ctx.Unique == nil && qcq.path != nil {
		qcq.Unique(true)
	}
	ctx = setContextOp(ctx, qcq.ctx, "IDs")
	if err = qcq.Select(questioncategory.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (qcq *QuestionCategoryQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := qcq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (qcq *QuestionCategoryQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, qcq.ctx, "Count")
	if err := qcq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, qcq, querierCount[*QuestionCategoryQuery](), qcq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (qcq *QuestionCategoryQuery) CountX(ctx context.Context) int {
	count, err := qcq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (qcq *QuestionCategoryQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, qcq.ctx, "Exist")
	switch _, err := qcq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (qcq *QuestionCategoryQuery) ExistX(ctx context.Context) bool {
	exist, err := qcq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the QuestionCategoryQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (qcq *QuestionCategoryQuery) Clone() *QuestionCategoryQuery {
	if qcq == nil {
		return nil
	}
	return &QuestionCategoryQuery{
		config:        qcq.config,
		ctx:           qcq.ctx.Clone(),
		order:         append([]questioncategory.OrderOption{}, qcq.order...),
		inters:        append([]Interceptor{}, qcq.inters...),
		predicates:    append([]predicate.QuestionCategory{}, qcq.predicates...),
		withQuestions: qcq.withQuestions.Clone(),
		// clone intermediate query.
		sql:  qcq.sql.Clone(),
		path: qcq.path,
	}
}

// WithQuestions tells the query-builder to eager-load the nodes that are connected to
// the "questions" edge. The optional arguments are used to configure the query builder of the edge.
func (qcq *QuestionCategoryQuery) WithQuestions(opts ...func(*QuestionQuery)) *QuestionCategoryQuery {
	query := (&QuestionClient{config: qcq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	qcq.withQuestions = query
	return qcq
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
//	client.QuestionCategory.Query().
//		GroupBy(questioncategory.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (qcq *QuestionCategoryQuery) GroupBy(field string, fields ...string) *QuestionCategoryGroupBy {
	qcq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &QuestionCategoryGroupBy{build: qcq}
	grbuild.flds = &qcq.ctx.Fields
	grbuild.label = questioncategory.Label
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
//	client.QuestionCategory.Query().
//		Select(questioncategory.FieldCreatedAt).
//		Scan(ctx, &v)
func (qcq *QuestionCategoryQuery) Select(fields ...string) *QuestionCategorySelect {
	qcq.ctx.Fields = append(qcq.ctx.Fields, fields...)
	sbuild := &QuestionCategorySelect{QuestionCategoryQuery: qcq}
	sbuild.label = questioncategory.Label
	sbuild.flds, sbuild.scan = &qcq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a QuestionCategorySelect configured with the given aggregations.
func (qcq *QuestionCategoryQuery) Aggregate(fns ...AggregateFunc) *QuestionCategorySelect {
	return qcq.Select().Aggregate(fns...)
}

func (qcq *QuestionCategoryQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range qcq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, qcq); err != nil {
				return err
			}
		}
	}
	for _, f := range qcq.ctx.Fields {
		if !questioncategory.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if qcq.path != nil {
		prev, err := qcq.path(ctx)
		if err != nil {
			return err
		}
		qcq.sql = prev
	}
	return nil
}

func (qcq *QuestionCategoryQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*QuestionCategory, error) {
	var (
		nodes       = []*QuestionCategory{}
		_spec       = qcq.querySpec()
		loadedTypes = [1]bool{
			qcq.withQuestions != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*QuestionCategory).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &QuestionCategory{config: qcq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(qcq.modifiers) > 0 {
		_spec.Modifiers = qcq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, qcq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := qcq.withQuestions; query != nil {
		if err := qcq.loadQuestions(ctx, query, nodes,
			func(n *QuestionCategory) { n.Edges.Questions = []*Question{} },
			func(n *QuestionCategory, e *Question) { n.Edges.Questions = append(n.Edges.Questions, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (qcq *QuestionCategoryQuery) loadQuestions(ctx context.Context, query *QuestionQuery, nodes []*QuestionCategory, init func(*QuestionCategory), assign func(*QuestionCategory, *Question)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uint64]*QuestionCategory)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(question.FieldCategoryID)
	}
	query.Where(predicate.Question(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(questioncategory.QuestionsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.CategoryID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "category_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (qcq *QuestionCategoryQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := qcq.querySpec()
	if len(qcq.modifiers) > 0 {
		_spec.Modifiers = qcq.modifiers
	}
	_spec.Node.Columns = qcq.ctx.Fields
	if len(qcq.ctx.Fields) > 0 {
		_spec.Unique = qcq.ctx.Unique != nil && *qcq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, qcq.driver, _spec)
}

func (qcq *QuestionCategoryQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(questioncategory.Table, questioncategory.Columns, sqlgraph.NewFieldSpec(questioncategory.FieldID, field.TypeUint64))
	_spec.From = qcq.sql
	if unique := qcq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if qcq.path != nil {
		_spec.Unique = true
	}
	if fields := qcq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, questioncategory.FieldID)
		for i := range fields {
			if fields[i] != questioncategory.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := qcq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := qcq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := qcq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := qcq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (qcq *QuestionCategoryQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(qcq.driver.Dialect())
	t1 := builder.Table(questioncategory.Table)
	columns := qcq.ctx.Fields
	if len(columns) == 0 {
		columns = questioncategory.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if qcq.sql != nil {
		selector = qcq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if qcq.ctx.Unique != nil && *qcq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range qcq.modifiers {
		m(selector)
	}
	for _, p := range qcq.predicates {
		p(selector)
	}
	for _, p := range qcq.order {
		p(selector)
	}
	if offset := qcq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := qcq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (qcq *QuestionCategoryQuery) Modify(modifiers ...func(s *sql.Selector)) *QuestionCategorySelect {
	qcq.modifiers = append(qcq.modifiers, modifiers...)
	return qcq.Select()
}

type QuestionCategoryQueryWith string

var (
	QuestionCategoryQueryWithQuestions QuestionCategoryQueryWith = "Questions"
)

func (qcq *QuestionCategoryQuery) With(withEdges ...QuestionCategoryQueryWith) *QuestionCategoryQuery {
	for _, v := range withEdges {
		switch v {
		case QuestionCategoryQueryWithQuestions:
			qcq.WithQuestions()
		}
	}
	return qcq
}

// QuestionCategoryGroupBy is the group-by builder for QuestionCategory entities.
type QuestionCategoryGroupBy struct {
	selector
	build *QuestionCategoryQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (qcgb *QuestionCategoryGroupBy) Aggregate(fns ...AggregateFunc) *QuestionCategoryGroupBy {
	qcgb.fns = append(qcgb.fns, fns...)
	return qcgb
}

// Scan applies the selector query and scans the result into the given value.
func (qcgb *QuestionCategoryGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, qcgb.build.ctx, "GroupBy")
	if err := qcgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*QuestionCategoryQuery, *QuestionCategoryGroupBy](ctx, qcgb.build, qcgb, qcgb.build.inters, v)
}

func (qcgb *QuestionCategoryGroupBy) sqlScan(ctx context.Context, root *QuestionCategoryQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(qcgb.fns))
	for _, fn := range qcgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*qcgb.flds)+len(qcgb.fns))
		for _, f := range *qcgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*qcgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := qcgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// QuestionCategorySelect is the builder for selecting fields of QuestionCategory entities.
type QuestionCategorySelect struct {
	*QuestionCategoryQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (qcs *QuestionCategorySelect) Aggregate(fns ...AggregateFunc) *QuestionCategorySelect {
	qcs.fns = append(qcs.fns, fns...)
	return qcs
}

// Scan applies the selector query and scans the result into the given value.
func (qcs *QuestionCategorySelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, qcs.ctx, "Select")
	if err := qcs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*QuestionCategoryQuery, *QuestionCategorySelect](ctx, qcs.QuestionCategoryQuery, qcs, qcs.inters, v)
}

func (qcs *QuestionCategorySelect) sqlScan(ctx context.Context, root *QuestionCategoryQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(qcs.fns))
	for _, fn := range qcs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*qcs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := qcs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (qcs *QuestionCategorySelect) Modify(modifiers ...func(s *sql.Selector)) *QuestionCategorySelect {
	qcs.modifiers = append(qcs.modifiers, modifiers...)
	return qcs
}