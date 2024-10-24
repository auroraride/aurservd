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
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/assetcheck"
	"github.com/auroraride/aurservd/internal/ent/assetcheckdetails"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
)

// AssetCheckQuery is the builder for querying AssetCheck entities.
type AssetCheckQuery struct {
	config
	ctx                     *QueryContext
	order                   []assetcheck.OrderOption
	inters                  []Interceptor
	predicates              []predicate.AssetCheck
	withCheckDetails        *AssetCheckDetailsQuery
	withOperateAssetManager *AssetManagerQuery
	withOperateEmployee     *EmployeeQuery
	withOperateAgent        *AgentQuery
	withWarehouse           *WarehouseQuery
	withStore               *StoreQuery
	withStation             *EnterpriseStationQuery
	modifiers               []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AssetCheckQuery builder.
func (acq *AssetCheckQuery) Where(ps ...predicate.AssetCheck) *AssetCheckQuery {
	acq.predicates = append(acq.predicates, ps...)
	return acq
}

// Limit the number of records to be returned by this query.
func (acq *AssetCheckQuery) Limit(limit int) *AssetCheckQuery {
	acq.ctx.Limit = &limit
	return acq
}

// Offset to start from.
func (acq *AssetCheckQuery) Offset(offset int) *AssetCheckQuery {
	acq.ctx.Offset = &offset
	return acq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (acq *AssetCheckQuery) Unique(unique bool) *AssetCheckQuery {
	acq.ctx.Unique = &unique
	return acq
}

// Order specifies how the records should be ordered.
func (acq *AssetCheckQuery) Order(o ...assetcheck.OrderOption) *AssetCheckQuery {
	acq.order = append(acq.order, o...)
	return acq
}

// QueryCheckDetails chains the current query on the "check_details" edge.
func (acq *AssetCheckQuery) QueryCheckDetails() *AssetCheckDetailsQuery {
	query := (&AssetCheckDetailsClient{config: acq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := acq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := acq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetcheck.Table, assetcheck.FieldID, selector),
			sqlgraph.To(assetcheckdetails.Table, assetcheckdetails.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, assetcheck.CheckDetailsTable, assetcheck.CheckDetailsColumn),
		)
		fromU = sqlgraph.SetNeighbors(acq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryOperateAssetManager chains the current query on the "operate_asset_manager" edge.
func (acq *AssetCheckQuery) QueryOperateAssetManager() *AssetManagerQuery {
	query := (&AssetManagerClient{config: acq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := acq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := acq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetcheck.Table, assetcheck.FieldID, selector),
			sqlgraph.To(assetmanager.Table, assetmanager.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, assetcheck.OperateAssetManagerTable, assetcheck.OperateAssetManagerColumn),
		)
		fromU = sqlgraph.SetNeighbors(acq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryOperateEmployee chains the current query on the "operate_employee" edge.
func (acq *AssetCheckQuery) QueryOperateEmployee() *EmployeeQuery {
	query := (&EmployeeClient{config: acq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := acq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := acq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetcheck.Table, assetcheck.FieldID, selector),
			sqlgraph.To(employee.Table, employee.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, assetcheck.OperateEmployeeTable, assetcheck.OperateEmployeeColumn),
		)
		fromU = sqlgraph.SetNeighbors(acq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryOperateAgent chains the current query on the "operate_agent" edge.
func (acq *AssetCheckQuery) QueryOperateAgent() *AgentQuery {
	query := (&AgentClient{config: acq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := acq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := acq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetcheck.Table, assetcheck.FieldID, selector),
			sqlgraph.To(agent.Table, agent.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, assetcheck.OperateAgentTable, assetcheck.OperateAgentColumn),
		)
		fromU = sqlgraph.SetNeighbors(acq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryWarehouse chains the current query on the "warehouse" edge.
func (acq *AssetCheckQuery) QueryWarehouse() *WarehouseQuery {
	query := (&WarehouseClient{config: acq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := acq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := acq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetcheck.Table, assetcheck.FieldID, selector),
			sqlgraph.To(warehouse.Table, warehouse.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, assetcheck.WarehouseTable, assetcheck.WarehouseColumn),
		)
		fromU = sqlgraph.SetNeighbors(acq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryStore chains the current query on the "store" edge.
func (acq *AssetCheckQuery) QueryStore() *StoreQuery {
	query := (&StoreClient{config: acq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := acq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := acq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetcheck.Table, assetcheck.FieldID, selector),
			sqlgraph.To(store.Table, store.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, assetcheck.StoreTable, assetcheck.StoreColumn),
		)
		fromU = sqlgraph.SetNeighbors(acq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryStation chains the current query on the "station" edge.
func (acq *AssetCheckQuery) QueryStation() *EnterpriseStationQuery {
	query := (&EnterpriseStationClient{config: acq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := acq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := acq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(assetcheck.Table, assetcheck.FieldID, selector),
			sqlgraph.To(enterprisestation.Table, enterprisestation.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, assetcheck.StationTable, assetcheck.StationColumn),
		)
		fromU = sqlgraph.SetNeighbors(acq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first AssetCheck entity from the query.
// Returns a *NotFoundError when no AssetCheck was found.
func (acq *AssetCheckQuery) First(ctx context.Context) (*AssetCheck, error) {
	nodes, err := acq.Limit(1).All(setContextOp(ctx, acq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{assetcheck.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (acq *AssetCheckQuery) FirstX(ctx context.Context) *AssetCheck {
	node, err := acq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first AssetCheck ID from the query.
// Returns a *NotFoundError when no AssetCheck ID was found.
func (acq *AssetCheckQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = acq.Limit(1).IDs(setContextOp(ctx, acq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{assetcheck.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (acq *AssetCheckQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := acq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single AssetCheck entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one AssetCheck entity is found.
// Returns a *NotFoundError when no AssetCheck entities are found.
func (acq *AssetCheckQuery) Only(ctx context.Context) (*AssetCheck, error) {
	nodes, err := acq.Limit(2).All(setContextOp(ctx, acq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{assetcheck.Label}
	default:
		return nil, &NotSingularError{assetcheck.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (acq *AssetCheckQuery) OnlyX(ctx context.Context) *AssetCheck {
	node, err := acq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only AssetCheck ID in the query.
// Returns a *NotSingularError when more than one AssetCheck ID is found.
// Returns a *NotFoundError when no entities are found.
func (acq *AssetCheckQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = acq.Limit(2).IDs(setContextOp(ctx, acq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{assetcheck.Label}
	default:
		err = &NotSingularError{assetcheck.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (acq *AssetCheckQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := acq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AssetChecks.
func (acq *AssetCheckQuery) All(ctx context.Context) ([]*AssetCheck, error) {
	ctx = setContextOp(ctx, acq.ctx, ent.OpQueryAll)
	if err := acq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*AssetCheck, *AssetCheckQuery]()
	return withInterceptors[[]*AssetCheck](ctx, acq, qr, acq.inters)
}

// AllX is like All, but panics if an error occurs.
func (acq *AssetCheckQuery) AllX(ctx context.Context) []*AssetCheck {
	nodes, err := acq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of AssetCheck IDs.
func (acq *AssetCheckQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if acq.ctx.Unique == nil && acq.path != nil {
		acq.Unique(true)
	}
	ctx = setContextOp(ctx, acq.ctx, ent.OpQueryIDs)
	if err = acq.Select(assetcheck.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (acq *AssetCheckQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := acq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (acq *AssetCheckQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, acq.ctx, ent.OpQueryCount)
	if err := acq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, acq, querierCount[*AssetCheckQuery](), acq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (acq *AssetCheckQuery) CountX(ctx context.Context) int {
	count, err := acq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (acq *AssetCheckQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, acq.ctx, ent.OpQueryExist)
	switch _, err := acq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (acq *AssetCheckQuery) ExistX(ctx context.Context) bool {
	exist, err := acq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AssetCheckQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (acq *AssetCheckQuery) Clone() *AssetCheckQuery {
	if acq == nil {
		return nil
	}
	return &AssetCheckQuery{
		config:                  acq.config,
		ctx:                     acq.ctx.Clone(),
		order:                   append([]assetcheck.OrderOption{}, acq.order...),
		inters:                  append([]Interceptor{}, acq.inters...),
		predicates:              append([]predicate.AssetCheck{}, acq.predicates...),
		withCheckDetails:        acq.withCheckDetails.Clone(),
		withOperateAssetManager: acq.withOperateAssetManager.Clone(),
		withOperateEmployee:     acq.withOperateEmployee.Clone(),
		withOperateAgent:        acq.withOperateAgent.Clone(),
		withWarehouse:           acq.withWarehouse.Clone(),
		withStore:               acq.withStore.Clone(),
		withStation:             acq.withStation.Clone(),
		// clone intermediate query.
		sql:       acq.sql.Clone(),
		path:      acq.path,
		modifiers: append([]func(*sql.Selector){}, acq.modifiers...),
	}
}

// WithCheckDetails tells the query-builder to eager-load the nodes that are connected to
// the "check_details" edge. The optional arguments are used to configure the query builder of the edge.
func (acq *AssetCheckQuery) WithCheckDetails(opts ...func(*AssetCheckDetailsQuery)) *AssetCheckQuery {
	query := (&AssetCheckDetailsClient{config: acq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	acq.withCheckDetails = query
	return acq
}

// WithOperateAssetManager tells the query-builder to eager-load the nodes that are connected to
// the "operate_asset_manager" edge. The optional arguments are used to configure the query builder of the edge.
func (acq *AssetCheckQuery) WithOperateAssetManager(opts ...func(*AssetManagerQuery)) *AssetCheckQuery {
	query := (&AssetManagerClient{config: acq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	acq.withOperateAssetManager = query
	return acq
}

// WithOperateEmployee tells the query-builder to eager-load the nodes that are connected to
// the "operate_employee" edge. The optional arguments are used to configure the query builder of the edge.
func (acq *AssetCheckQuery) WithOperateEmployee(opts ...func(*EmployeeQuery)) *AssetCheckQuery {
	query := (&EmployeeClient{config: acq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	acq.withOperateEmployee = query
	return acq
}

// WithOperateAgent tells the query-builder to eager-load the nodes that are connected to
// the "operate_agent" edge. The optional arguments are used to configure the query builder of the edge.
func (acq *AssetCheckQuery) WithOperateAgent(opts ...func(*AgentQuery)) *AssetCheckQuery {
	query := (&AgentClient{config: acq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	acq.withOperateAgent = query
	return acq
}

// WithWarehouse tells the query-builder to eager-load the nodes that are connected to
// the "warehouse" edge. The optional arguments are used to configure the query builder of the edge.
func (acq *AssetCheckQuery) WithWarehouse(opts ...func(*WarehouseQuery)) *AssetCheckQuery {
	query := (&WarehouseClient{config: acq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	acq.withWarehouse = query
	return acq
}

// WithStore tells the query-builder to eager-load the nodes that are connected to
// the "store" edge. The optional arguments are used to configure the query builder of the edge.
func (acq *AssetCheckQuery) WithStore(opts ...func(*StoreQuery)) *AssetCheckQuery {
	query := (&StoreClient{config: acq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	acq.withStore = query
	return acq
}

// WithStation tells the query-builder to eager-load the nodes that are connected to
// the "station" edge. The optional arguments are used to configure the query builder of the edge.
func (acq *AssetCheckQuery) WithStation(opts ...func(*EnterpriseStationQuery)) *AssetCheckQuery {
	query := (&EnterpriseStationClient{config: acq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	acq.withStation = query
	return acq
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
//	client.AssetCheck.Query().
//		GroupBy(assetcheck.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (acq *AssetCheckQuery) GroupBy(field string, fields ...string) *AssetCheckGroupBy {
	acq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &AssetCheckGroupBy{build: acq}
	grbuild.flds = &acq.ctx.Fields
	grbuild.label = assetcheck.Label
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
//	client.AssetCheck.Query().
//		Select(assetcheck.FieldCreatedAt).
//		Scan(ctx, &v)
func (acq *AssetCheckQuery) Select(fields ...string) *AssetCheckSelect {
	acq.ctx.Fields = append(acq.ctx.Fields, fields...)
	sbuild := &AssetCheckSelect{AssetCheckQuery: acq}
	sbuild.label = assetcheck.Label
	sbuild.flds, sbuild.scan = &acq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a AssetCheckSelect configured with the given aggregations.
func (acq *AssetCheckQuery) Aggregate(fns ...AggregateFunc) *AssetCheckSelect {
	return acq.Select().Aggregate(fns...)
}

func (acq *AssetCheckQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range acq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, acq); err != nil {
				return err
			}
		}
	}
	for _, f := range acq.ctx.Fields {
		if !assetcheck.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if acq.path != nil {
		prev, err := acq.path(ctx)
		if err != nil {
			return err
		}
		acq.sql = prev
	}
	return nil
}

func (acq *AssetCheckQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*AssetCheck, error) {
	var (
		nodes       = []*AssetCheck{}
		_spec       = acq.querySpec()
		loadedTypes = [7]bool{
			acq.withCheckDetails != nil,
			acq.withOperateAssetManager != nil,
			acq.withOperateEmployee != nil,
			acq.withOperateAgent != nil,
			acq.withWarehouse != nil,
			acq.withStore != nil,
			acq.withStation != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*AssetCheck).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &AssetCheck{config: acq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(acq.modifiers) > 0 {
		_spec.Modifiers = acq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, acq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := acq.withCheckDetails; query != nil {
		if err := acq.loadCheckDetails(ctx, query, nodes,
			func(n *AssetCheck) { n.Edges.CheckDetails = []*AssetCheckDetails{} },
			func(n *AssetCheck, e *AssetCheckDetails) { n.Edges.CheckDetails = append(n.Edges.CheckDetails, e) }); err != nil {
			return nil, err
		}
	}
	if query := acq.withOperateAssetManager; query != nil {
		if err := acq.loadOperateAssetManager(ctx, query, nodes, nil,
			func(n *AssetCheck, e *AssetManager) { n.Edges.OperateAssetManager = e }); err != nil {
			return nil, err
		}
	}
	if query := acq.withOperateEmployee; query != nil {
		if err := acq.loadOperateEmployee(ctx, query, nodes, nil,
			func(n *AssetCheck, e *Employee) { n.Edges.OperateEmployee = e }); err != nil {
			return nil, err
		}
	}
	if query := acq.withOperateAgent; query != nil {
		if err := acq.loadOperateAgent(ctx, query, nodes, nil,
			func(n *AssetCheck, e *Agent) { n.Edges.OperateAgent = e }); err != nil {
			return nil, err
		}
	}
	if query := acq.withWarehouse; query != nil {
		if err := acq.loadWarehouse(ctx, query, nodes, nil,
			func(n *AssetCheck, e *Warehouse) { n.Edges.Warehouse = e }); err != nil {
			return nil, err
		}
	}
	if query := acq.withStore; query != nil {
		if err := acq.loadStore(ctx, query, nodes, nil,
			func(n *AssetCheck, e *Store) { n.Edges.Store = e }); err != nil {
			return nil, err
		}
	}
	if query := acq.withStation; query != nil {
		if err := acq.loadStation(ctx, query, nodes, nil,
			func(n *AssetCheck, e *EnterpriseStation) { n.Edges.Station = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (acq *AssetCheckQuery) loadCheckDetails(ctx context.Context, query *AssetCheckDetailsQuery, nodes []*AssetCheck, init func(*AssetCheck), assign func(*AssetCheck, *AssetCheckDetails)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uint64]*AssetCheck)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(assetcheckdetails.FieldCheckID)
	}
	query.Where(predicate.AssetCheckDetails(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(assetcheck.CheckDetailsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.CheckID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "check_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (acq *AssetCheckQuery) loadOperateAssetManager(ctx context.Context, query *AssetManagerQuery, nodes []*AssetCheck, init func(*AssetCheck), assign func(*AssetCheck, *AssetManager)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*AssetCheck)
	for i := range nodes {
		fk := nodes[i].OperateID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(assetmanager.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "operate_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (acq *AssetCheckQuery) loadOperateEmployee(ctx context.Context, query *EmployeeQuery, nodes []*AssetCheck, init func(*AssetCheck), assign func(*AssetCheck, *Employee)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*AssetCheck)
	for i := range nodes {
		fk := nodes[i].OperateID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(employee.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "operate_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (acq *AssetCheckQuery) loadOperateAgent(ctx context.Context, query *AgentQuery, nodes []*AssetCheck, init func(*AssetCheck), assign func(*AssetCheck, *Agent)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*AssetCheck)
	for i := range nodes {
		fk := nodes[i].OperateID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(agent.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "operate_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (acq *AssetCheckQuery) loadWarehouse(ctx context.Context, query *WarehouseQuery, nodes []*AssetCheck, init func(*AssetCheck), assign func(*AssetCheck, *Warehouse)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*AssetCheck)
	for i := range nodes {
		fk := nodes[i].LocationsID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(warehouse.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "locations_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (acq *AssetCheckQuery) loadStore(ctx context.Context, query *StoreQuery, nodes []*AssetCheck, init func(*AssetCheck), assign func(*AssetCheck, *Store)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*AssetCheck)
	for i := range nodes {
		fk := nodes[i].LocationsID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(store.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "locations_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (acq *AssetCheckQuery) loadStation(ctx context.Context, query *EnterpriseStationQuery, nodes []*AssetCheck, init func(*AssetCheck), assign func(*AssetCheck, *EnterpriseStation)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*AssetCheck)
	for i := range nodes {
		fk := nodes[i].LocationsID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(enterprisestation.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "locations_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (acq *AssetCheckQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := acq.querySpec()
	if len(acq.modifiers) > 0 {
		_spec.Modifiers = acq.modifiers
	}
	_spec.Node.Columns = acq.ctx.Fields
	if len(acq.ctx.Fields) > 0 {
		_spec.Unique = acq.ctx.Unique != nil && *acq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, acq.driver, _spec)
}

func (acq *AssetCheckQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(assetcheck.Table, assetcheck.Columns, sqlgraph.NewFieldSpec(assetcheck.FieldID, field.TypeUint64))
	_spec.From = acq.sql
	if unique := acq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if acq.path != nil {
		_spec.Unique = true
	}
	if fields := acq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, assetcheck.FieldID)
		for i := range fields {
			if fields[i] != assetcheck.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if acq.withOperateAssetManager != nil {
			_spec.Node.AddColumnOnce(assetcheck.FieldOperateID)
		}
		if acq.withOperateEmployee != nil {
			_spec.Node.AddColumnOnce(assetcheck.FieldOperateID)
		}
		if acq.withOperateAgent != nil {
			_spec.Node.AddColumnOnce(assetcheck.FieldOperateID)
		}
		if acq.withWarehouse != nil {
			_spec.Node.AddColumnOnce(assetcheck.FieldLocationsID)
		}
		if acq.withStore != nil {
			_spec.Node.AddColumnOnce(assetcheck.FieldLocationsID)
		}
		if acq.withStation != nil {
			_spec.Node.AddColumnOnce(assetcheck.FieldLocationsID)
		}
	}
	if ps := acq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := acq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := acq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := acq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (acq *AssetCheckQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(acq.driver.Dialect())
	t1 := builder.Table(assetcheck.Table)
	columns := acq.ctx.Fields
	if len(columns) == 0 {
		columns = assetcheck.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if acq.sql != nil {
		selector = acq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if acq.ctx.Unique != nil && *acq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range acq.modifiers {
		m(selector)
	}
	for _, p := range acq.predicates {
		p(selector)
	}
	for _, p := range acq.order {
		p(selector)
	}
	if offset := acq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := acq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (acq *AssetCheckQuery) Modify(modifiers ...func(s *sql.Selector)) *AssetCheckSelect {
	acq.modifiers = append(acq.modifiers, modifiers...)
	return acq.Select()
}

type AssetCheckQueryWith string

var (
	AssetCheckQueryWithCheckDetails        AssetCheckQueryWith = "CheckDetails"
	AssetCheckQueryWithOperateAssetManager AssetCheckQueryWith = "OperateAssetManager"
	AssetCheckQueryWithOperateEmployee     AssetCheckQueryWith = "OperateEmployee"
	AssetCheckQueryWithOperateAgent        AssetCheckQueryWith = "OperateAgent"
	AssetCheckQueryWithWarehouse           AssetCheckQueryWith = "Warehouse"
	AssetCheckQueryWithStore               AssetCheckQueryWith = "Store"
	AssetCheckQueryWithStation             AssetCheckQueryWith = "Station"
)

func (acq *AssetCheckQuery) With(withEdges ...AssetCheckQueryWith) *AssetCheckQuery {
	for _, v := range withEdges {
		switch v {
		case AssetCheckQueryWithCheckDetails:
			acq.WithCheckDetails()
		case AssetCheckQueryWithOperateAssetManager:
			acq.WithOperateAssetManager()
		case AssetCheckQueryWithOperateEmployee:
			acq.WithOperateEmployee()
		case AssetCheckQueryWithOperateAgent:
			acq.WithOperateAgent()
		case AssetCheckQueryWithWarehouse:
			acq.WithWarehouse()
		case AssetCheckQueryWithStore:
			acq.WithStore()
		case AssetCheckQueryWithStation:
			acq.WithStation()
		}
	}
	return acq
}

// AssetCheckGroupBy is the group-by builder for AssetCheck entities.
type AssetCheckGroupBy struct {
	selector
	build *AssetCheckQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (acgb *AssetCheckGroupBy) Aggregate(fns ...AggregateFunc) *AssetCheckGroupBy {
	acgb.fns = append(acgb.fns, fns...)
	return acgb
}

// Scan applies the selector query and scans the result into the given value.
func (acgb *AssetCheckGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, acgb.build.ctx, ent.OpQueryGroupBy)
	if err := acgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AssetCheckQuery, *AssetCheckGroupBy](ctx, acgb.build, acgb, acgb.build.inters, v)
}

func (acgb *AssetCheckGroupBy) sqlScan(ctx context.Context, root *AssetCheckQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(acgb.fns))
	for _, fn := range acgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*acgb.flds)+len(acgb.fns))
		for _, f := range *acgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*acgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := acgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// AssetCheckSelect is the builder for selecting fields of AssetCheck entities.
type AssetCheckSelect struct {
	*AssetCheckQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (acs *AssetCheckSelect) Aggregate(fns ...AggregateFunc) *AssetCheckSelect {
	acs.fns = append(acs.fns, fns...)
	return acs
}

// Scan applies the selector query and scans the result into the given value.
func (acs *AssetCheckSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, acs.ctx, ent.OpQuerySelect)
	if err := acs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AssetCheckQuery, *AssetCheckSelect](ctx, acs.AssetCheckQuery, acs, acs.inters, v)
}

func (acs *AssetCheckSelect) sqlScan(ctx context.Context, root *AssetCheckQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(acs.fns))
	for _, fn := range acs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*acs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := acs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (acs *AssetCheckSelect) Modify(modifiers ...func(s *sql.Selector)) *AssetCheckSelect {
	acs.modifiers = append(acs.modifiers, modifiers...)
	return acs
}
