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
	"github.com/auroraride/aurservd/internal/ent/assetmaintenance"
	"github.com/auroraride/aurservd/internal/ent/assetmaintenancedetails"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// AssetMaintenanceQuery is the builder for querying AssetMaintenance entities.
type AssetMaintenanceQuery struct {
	config
	ctx                    *QueryContext
	order                  []assetmaintenance.OrderOption
	inters                 []Interceptor
	predicates             []predicate.AssetMaintenance
	withCabinet            *CabinetQuery
	withMaintainer         *MaintainerQuery
	withMaintenanceDetails *AssetMaintenanceDetailsQuery
	modifiers              []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AssetMaintenanceQuery builder.
func (amq *AssetMaintenanceQuery) Where(ps ...predicate.AssetMaintenance) *AssetMaintenanceQuery {
	amq.predicates = append(amq.predicates, ps...)
	return amq
}

// Limit the number of records to be returned by this query.
func (amq *AssetMaintenanceQuery) Limit(limit int) *AssetMaintenanceQuery {
	amq.ctx.Limit = &limit
	return amq
}

// Offset to start from.
func (amq *AssetMaintenanceQuery) Offset(offset int) *AssetMaintenanceQuery {
	amq.ctx.Offset = &offset
	return amq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (amq *AssetMaintenanceQuery) Unique(unique bool) *AssetMaintenanceQuery {
	amq.ctx.Unique = &unique
	return amq
}

// Order specifies how the records should be ordered.
func (amq *AssetMaintenanceQuery) Order(o ...assetmaintenance.OrderOption) *AssetMaintenanceQuery {
	amq.order = append(amq.order, o...)
	return amq
}

// QueryCabinet chains the current query on the "cabinet" edge.
func (amq *AssetMaintenanceQuery) QueryCabinet() *CabinetQuery {
	query := (&CabinetClient{config: amq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := amq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := amq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetmaintenance.Table, assetmaintenance.FieldID, selector),
			sqlgraph.To(cabinet.Table, cabinet.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, assetmaintenance.CabinetTable, assetmaintenance.CabinetColumn),
		)
		fromU = sqlgraph.SetNeighbors(amq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryMaintainer chains the current query on the "maintainer" edge.
func (amq *AssetMaintenanceQuery) QueryMaintainer() *MaintainerQuery {
	query := (&MaintainerClient{config: amq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := amq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := amq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetmaintenance.Table, assetmaintenance.FieldID, selector),
			sqlgraph.To(maintainer.Table, maintainer.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, assetmaintenance.MaintainerTable, assetmaintenance.MaintainerColumn),
		)
		fromU = sqlgraph.SetNeighbors(amq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryMaintenanceDetails chains the current query on the "maintenance_details" edge.
func (amq *AssetMaintenanceQuery) QueryMaintenanceDetails() *AssetMaintenanceDetailsQuery {
	query := (&AssetMaintenanceDetailsClient{config: amq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := amq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := amq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetmaintenance.Table, assetmaintenance.FieldID, selector),
			sqlgraph.To(assetmaintenancedetails.Table, assetmaintenancedetails.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, assetmaintenance.MaintenanceDetailsTable, assetmaintenance.MaintenanceDetailsColumn),
		)
		fromU = sqlgraph.SetNeighbors(amq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first AssetMaintenance entity from the query.
// Returns a *NotFoundError when no AssetMaintenance was found.
func (amq *AssetMaintenanceQuery) First(ctx context.Context) (*AssetMaintenance, error) {
	nodes, err := amq.Limit(1).All(setContextOp(ctx, amq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{assetmaintenance.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (amq *AssetMaintenanceQuery) FirstX(ctx context.Context) *AssetMaintenance {
	node, err := amq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first AssetMaintenance ID from the query.
// Returns a *NotFoundError when no AssetMaintenance ID was found.
func (amq *AssetMaintenanceQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = amq.Limit(1).IDs(setContextOp(ctx, amq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{assetmaintenance.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (amq *AssetMaintenanceQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := amq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single AssetMaintenance entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one AssetMaintenance entity is found.
// Returns a *NotFoundError when no AssetMaintenance entities are found.
func (amq *AssetMaintenanceQuery) Only(ctx context.Context) (*AssetMaintenance, error) {
	nodes, err := amq.Limit(2).All(setContextOp(ctx, amq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{assetmaintenance.Label}
	default:
		return nil, &NotSingularError{assetmaintenance.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (amq *AssetMaintenanceQuery) OnlyX(ctx context.Context) *AssetMaintenance {
	node, err := amq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only AssetMaintenance ID in the query.
// Returns a *NotSingularError when more than one AssetMaintenance ID is found.
// Returns a *NotFoundError when no entities are found.
func (amq *AssetMaintenanceQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = amq.Limit(2).IDs(setContextOp(ctx, amq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{assetmaintenance.Label}
	default:
		err = &NotSingularError{assetmaintenance.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (amq *AssetMaintenanceQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := amq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AssetMaintenances.
func (amq *AssetMaintenanceQuery) All(ctx context.Context) ([]*AssetMaintenance, error) {
	ctx = setContextOp(ctx, amq.ctx, ent.OpQueryAll)
	if err := amq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*AssetMaintenance, *AssetMaintenanceQuery]()
	return withInterceptors[[]*AssetMaintenance](ctx, amq, qr, amq.inters)
}

// AllX is like All, but panics if an error occurs.
func (amq *AssetMaintenanceQuery) AllX(ctx context.Context) []*AssetMaintenance {
	nodes, err := amq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of AssetMaintenance IDs.
func (amq *AssetMaintenanceQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if amq.ctx.Unique == nil && amq.path != nil {
		amq.Unique(true)
	}
	ctx = setContextOp(ctx, amq.ctx, ent.OpQueryIDs)
	if err = amq.Select(assetmaintenance.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (amq *AssetMaintenanceQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := amq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (amq *AssetMaintenanceQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, amq.ctx, ent.OpQueryCount)
	if err := amq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, amq, querierCount[*AssetMaintenanceQuery](), amq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (amq *AssetMaintenanceQuery) CountX(ctx context.Context) int {
	count, err := amq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (amq *AssetMaintenanceQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, amq.ctx, ent.OpQueryExist)
	switch _, err := amq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (amq *AssetMaintenanceQuery) ExistX(ctx context.Context) bool {
	exist, err := amq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AssetMaintenanceQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (amq *AssetMaintenanceQuery) Clone() *AssetMaintenanceQuery {
	if amq == nil {
		return nil
	}
	return &AssetMaintenanceQuery{
		config:                 amq.config,
		ctx:                    amq.ctx.Clone(),
		order:                  append([]assetmaintenance.OrderOption{}, amq.order...),
		inters:                 append([]Interceptor{}, amq.inters...),
		predicates:             append([]predicate.AssetMaintenance{}, amq.predicates...),
		withCabinet:            amq.withCabinet.Clone(),
		withMaintainer:         amq.withMaintainer.Clone(),
		withMaintenanceDetails: amq.withMaintenanceDetails.Clone(),
		// clone intermediate query.
		sql:       amq.sql.Clone(),
		path:      amq.path,
		modifiers: append([]func(*sql.Selector){}, amq.modifiers...),
	}
}

// WithCabinet tells the query-builder to eager-load the nodes that are connected to
// the "cabinet" edge. The optional arguments are used to configure the query builder of the edge.
func (amq *AssetMaintenanceQuery) WithCabinet(opts ...func(*CabinetQuery)) *AssetMaintenanceQuery {
	query := (&CabinetClient{config: amq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	amq.withCabinet = query
	return amq
}

// WithMaintainer tells the query-builder to eager-load the nodes that are connected to
// the "maintainer" edge. The optional arguments are used to configure the query builder of the edge.
func (amq *AssetMaintenanceQuery) WithMaintainer(opts ...func(*MaintainerQuery)) *AssetMaintenanceQuery {
	query := (&MaintainerClient{config: amq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	amq.withMaintainer = query
	return amq
}

// WithMaintenanceDetails tells the query-builder to eager-load the nodes that are connected to
// the "maintenance_details" edge. The optional arguments are used to configure the query builder of the edge.
func (amq *AssetMaintenanceQuery) WithMaintenanceDetails(opts ...func(*AssetMaintenanceDetailsQuery)) *AssetMaintenanceQuery {
	query := (&AssetMaintenanceDetailsClient{config: amq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	amq.withMaintenanceDetails = query
	return amq
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
//	client.AssetMaintenance.Query().
//		GroupBy(assetmaintenance.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (amq *AssetMaintenanceQuery) GroupBy(field string, fields ...string) *AssetMaintenanceGroupBy {
	amq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &AssetMaintenanceGroupBy{build: amq}
	grbuild.flds = &amq.ctx.Fields
	grbuild.label = assetmaintenance.Label
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
//	client.AssetMaintenance.Query().
//		Select(assetmaintenance.FieldCreatedAt).
//		Scan(ctx, &v)
func (amq *AssetMaintenanceQuery) Select(fields ...string) *AssetMaintenanceSelect {
	amq.ctx.Fields = append(amq.ctx.Fields, fields...)
	sbuild := &AssetMaintenanceSelect{AssetMaintenanceQuery: amq}
	sbuild.label = assetmaintenance.Label
	sbuild.flds, sbuild.scan = &amq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a AssetMaintenanceSelect configured with the given aggregations.
func (amq *AssetMaintenanceQuery) Aggregate(fns ...AggregateFunc) *AssetMaintenanceSelect {
	return amq.Select().Aggregate(fns...)
}

func (amq *AssetMaintenanceQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range amq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, amq); err != nil {
				return err
			}
		}
	}
	for _, f := range amq.ctx.Fields {
		if !assetmaintenance.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if amq.path != nil {
		prev, err := amq.path(ctx)
		if err != nil {
			return err
		}
		amq.sql = prev
	}
	return nil
}

func (amq *AssetMaintenanceQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*AssetMaintenance, error) {
	var (
		nodes       = []*AssetMaintenance{}
		_spec       = amq.querySpec()
		loadedTypes = [3]bool{
			amq.withCabinet != nil,
			amq.withMaintainer != nil,
			amq.withMaintenanceDetails != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*AssetMaintenance).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &AssetMaintenance{config: amq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(amq.modifiers) > 0 {
		_spec.Modifiers = amq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, amq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := amq.withCabinet; query != nil {
		if err := amq.loadCabinet(ctx, query, nodes, nil,
			func(n *AssetMaintenance, e *Cabinet) { n.Edges.Cabinet = e }); err != nil {
			return nil, err
		}
	}
	if query := amq.withMaintainer; query != nil {
		if err := amq.loadMaintainer(ctx, query, nodes, nil,
			func(n *AssetMaintenance, e *Maintainer) { n.Edges.Maintainer = e }); err != nil {
			return nil, err
		}
	}
	if query := amq.withMaintenanceDetails; query != nil {
		if err := amq.loadMaintenanceDetails(ctx, query, nodes,
			func(n *AssetMaintenance) { n.Edges.MaintenanceDetails = []*AssetMaintenanceDetails{} },
			func(n *AssetMaintenance, e *AssetMaintenanceDetails) {
				n.Edges.MaintenanceDetails = append(n.Edges.MaintenanceDetails, e)
			}); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (amq *AssetMaintenanceQuery) loadCabinet(ctx context.Context, query *CabinetQuery, nodes []*AssetMaintenance, init func(*AssetMaintenance), assign func(*AssetMaintenance, *Cabinet)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*AssetMaintenance)
	for i := range nodes {
		if nodes[i].CabinetID == nil {
			continue
		}
		fk := *nodes[i].CabinetID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(cabinet.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "cabinet_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (amq *AssetMaintenanceQuery) loadMaintainer(ctx context.Context, query *MaintainerQuery, nodes []*AssetMaintenance, init func(*AssetMaintenance), assign func(*AssetMaintenance, *Maintainer)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*AssetMaintenance)
	for i := range nodes {
		if nodes[i].MaintainerID == nil {
			continue
		}
		fk := *nodes[i].MaintainerID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(maintainer.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "maintainer_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (amq *AssetMaintenanceQuery) loadMaintenanceDetails(ctx context.Context, query *AssetMaintenanceDetailsQuery, nodes []*AssetMaintenance, init func(*AssetMaintenance), assign func(*AssetMaintenance, *AssetMaintenanceDetails)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uint64]*AssetMaintenance)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(assetmaintenancedetails.FieldMaintenanceID)
	}
	query.Where(predicate.AssetMaintenanceDetails(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(assetmaintenance.MaintenanceDetailsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.MaintenanceID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "maintenance_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (amq *AssetMaintenanceQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := amq.querySpec()
	if len(amq.modifiers) > 0 {
		_spec.Modifiers = amq.modifiers
	}
	_spec.Node.Columns = amq.ctx.Fields
	if len(amq.ctx.Fields) > 0 {
		_spec.Unique = amq.ctx.Unique != nil && *amq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, amq.driver, _spec)
}

func (amq *AssetMaintenanceQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(assetmaintenance.Table, assetmaintenance.Columns, sqlgraph.NewFieldSpec(assetmaintenance.FieldID, field.TypeUint64))
	_spec.From = amq.sql
	if unique := amq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if amq.path != nil {
		_spec.Unique = true
	}
	if fields := amq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, assetmaintenance.FieldID)
		for i := range fields {
			if fields[i] != assetmaintenance.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if amq.withCabinet != nil {
			_spec.Node.AddColumnOnce(assetmaintenance.FieldCabinetID)
		}
		if amq.withMaintainer != nil {
			_spec.Node.AddColumnOnce(assetmaintenance.FieldMaintainerID)
		}
	}
	if ps := amq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := amq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := amq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := amq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (amq *AssetMaintenanceQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(amq.driver.Dialect())
	t1 := builder.Table(assetmaintenance.Table)
	columns := amq.ctx.Fields
	if len(columns) == 0 {
		columns = assetmaintenance.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if amq.sql != nil {
		selector = amq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if amq.ctx.Unique != nil && *amq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range amq.modifiers {
		m(selector)
	}
	for _, p := range amq.predicates {
		p(selector)
	}
	for _, p := range amq.order {
		p(selector)
	}
	if offset := amq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := amq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (amq *AssetMaintenanceQuery) Modify(modifiers ...func(s *sql.Selector)) *AssetMaintenanceSelect {
	amq.modifiers = append(amq.modifiers, modifiers...)
	return amq.Select()
}

type AssetMaintenanceQueryWith string

var (
	AssetMaintenanceQueryWithCabinet            AssetMaintenanceQueryWith = "Cabinet"
	AssetMaintenanceQueryWithMaintainer         AssetMaintenanceQueryWith = "Maintainer"
	AssetMaintenanceQueryWithMaintenanceDetails AssetMaintenanceQueryWith = "MaintenanceDetails"
)

func (amq *AssetMaintenanceQuery) With(withEdges ...AssetMaintenanceQueryWith) *AssetMaintenanceQuery {
	for _, v := range withEdges {
		switch v {
		case AssetMaintenanceQueryWithCabinet:
			amq.WithCabinet()
		case AssetMaintenanceQueryWithMaintainer:
			amq.WithMaintainer()
		case AssetMaintenanceQueryWithMaintenanceDetails:
			amq.WithMaintenanceDetails()
		}
	}
	return amq
}

// AssetMaintenanceGroupBy is the group-by builder for AssetMaintenance entities.
type AssetMaintenanceGroupBy struct {
	selector
	build *AssetMaintenanceQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (amgb *AssetMaintenanceGroupBy) Aggregate(fns ...AggregateFunc) *AssetMaintenanceGroupBy {
	amgb.fns = append(amgb.fns, fns...)
	return amgb
}

// Scan applies the selector query and scans the result into the given value.
func (amgb *AssetMaintenanceGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, amgb.build.ctx, ent.OpQueryGroupBy)
	if err := amgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AssetMaintenanceQuery, *AssetMaintenanceGroupBy](ctx, amgb.build, amgb, amgb.build.inters, v)
}

func (amgb *AssetMaintenanceGroupBy) sqlScan(ctx context.Context, root *AssetMaintenanceQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(amgb.fns))
	for _, fn := range amgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*amgb.flds)+len(amgb.fns))
		for _, f := range *amgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*amgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := amgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// AssetMaintenanceSelect is the builder for selecting fields of AssetMaintenance entities.
type AssetMaintenanceSelect struct {
	*AssetMaintenanceQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ams *AssetMaintenanceSelect) Aggregate(fns ...AggregateFunc) *AssetMaintenanceSelect {
	ams.fns = append(ams.fns, fns...)
	return ams
}

// Scan applies the selector query and scans the result into the given value.
func (ams *AssetMaintenanceSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ams.ctx, ent.OpQuerySelect)
	if err := ams.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AssetMaintenanceQuery, *AssetMaintenanceSelect](ctx, ams.AssetMaintenanceQuery, ams, ams.inters, v)
}

func (ams *AssetMaintenanceSelect) sqlScan(ctx context.Context, root *AssetMaintenanceQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ams.fns))
	for _, fn := range ams.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ams.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ams.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ams *AssetMaintenanceSelect) Modify(modifiers ...func(s *sql.Selector)) *AssetMaintenanceSelect {
	ams.modifiers = append(ams.modifiers, modifiers...)
	return ams
}
