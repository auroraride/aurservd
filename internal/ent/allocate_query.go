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
	"github.com/auroraride/aurservd/internal/ent/allocate"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/contract"
	"github.com/auroraride/aurservd/internal/ent/ebike"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
)

// AllocateQuery is the builder for querying Allocate entities.
type AllocateQuery struct {
	config
	limit         *int
	offset        *int
	unique        *bool
	order         []OrderFunc
	fields        []string
	predicates    []predicate.Allocate
	withRider     *RiderQuery
	withEmployee  *EmployeeQuery
	withStore     *StoreQuery
	withEbike     *EbikeQuery
	withBrand     *EbikeBrandQuery
	withSubscribe *SubscribeQuery
	withCabinet   *CabinetQuery
	withContract  *ContractQuery
	modifiers     []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AllocateQuery builder.
func (aq *AllocateQuery) Where(ps ...predicate.Allocate) *AllocateQuery {
	aq.predicates = append(aq.predicates, ps...)
	return aq
}

// Limit adds a limit step to the query.
func (aq *AllocateQuery) Limit(limit int) *AllocateQuery {
	aq.limit = &limit
	return aq
}

// Offset adds an offset step to the query.
func (aq *AllocateQuery) Offset(offset int) *AllocateQuery {
	aq.offset = &offset
	return aq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (aq *AllocateQuery) Unique(unique bool) *AllocateQuery {
	aq.unique = &unique
	return aq
}

// Order adds an order step to the query.
func (aq *AllocateQuery) Order(o ...OrderFunc) *AllocateQuery {
	aq.order = append(aq.order, o...)
	return aq
}

// QueryRider chains the current query on the "rider" edge.
func (aq *AllocateQuery) QueryRider() *RiderQuery {
	query := &RiderQuery{config: aq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocate.Table, allocate.FieldID, selector),
			sqlgraph.To(rider.Table, rider.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, allocate.RiderTable, allocate.RiderColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEmployee chains the current query on the "employee" edge.
func (aq *AllocateQuery) QueryEmployee() *EmployeeQuery {
	query := &EmployeeQuery{config: aq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocate.Table, allocate.FieldID, selector),
			sqlgraph.To(employee.Table, employee.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, allocate.EmployeeTable, allocate.EmployeeColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryStore chains the current query on the "store" edge.
func (aq *AllocateQuery) QueryStore() *StoreQuery {
	query := &StoreQuery{config: aq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocate.Table, allocate.FieldID, selector),
			sqlgraph.To(store.Table, store.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, allocate.StoreTable, allocate.StoreColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEbike chains the current query on the "ebike" edge.
func (aq *AllocateQuery) QueryEbike() *EbikeQuery {
	query := &EbikeQuery{config: aq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocate.Table, allocate.FieldID, selector),
			sqlgraph.To(ebike.Table, ebike.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, allocate.EbikeTable, allocate.EbikeColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryBrand chains the current query on the "brand" edge.
func (aq *AllocateQuery) QueryBrand() *EbikeBrandQuery {
	query := &EbikeBrandQuery{config: aq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocate.Table, allocate.FieldID, selector),
			sqlgraph.To(ebikebrand.Table, ebikebrand.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, allocate.BrandTable, allocate.BrandColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QuerySubscribe chains the current query on the "subscribe" edge.
func (aq *AllocateQuery) QuerySubscribe() *SubscribeQuery {
	query := &SubscribeQuery{config: aq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocate.Table, allocate.FieldID, selector),
			sqlgraph.To(subscribe.Table, subscribe.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, allocate.SubscribeTable, allocate.SubscribeColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryCabinet chains the current query on the "cabinet" edge.
func (aq *AllocateQuery) QueryCabinet() *CabinetQuery {
	query := &CabinetQuery{config: aq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocate.Table, allocate.FieldID, selector),
			sqlgraph.To(cabinet.Table, cabinet.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, allocate.CabinetTable, allocate.CabinetColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryContract chains the current query on the "contract" edge.
func (aq *AllocateQuery) QueryContract() *ContractQuery {
	query := &ContractQuery{config: aq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocate.Table, allocate.FieldID, selector),
			sqlgraph.To(contract.Table, contract.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, allocate.ContractTable, allocate.ContractColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Allocate entity from the query.
// Returns a *NotFoundError when no Allocate was found.
func (aq *AllocateQuery) First(ctx context.Context) (*Allocate, error) {
	nodes, err := aq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{allocate.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (aq *AllocateQuery) FirstX(ctx context.Context) *Allocate {
	node, err := aq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Allocate ID from the query.
// Returns a *NotFoundError when no Allocate ID was found.
func (aq *AllocateQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = aq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{allocate.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (aq *AllocateQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := aq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Allocate entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Allocate entity is found.
// Returns a *NotFoundError when no Allocate entities are found.
func (aq *AllocateQuery) Only(ctx context.Context) (*Allocate, error) {
	nodes, err := aq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{allocate.Label}
	default:
		return nil, &NotSingularError{allocate.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (aq *AllocateQuery) OnlyX(ctx context.Context) *Allocate {
	node, err := aq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Allocate ID in the query.
// Returns a *NotSingularError when more than one Allocate ID is found.
// Returns a *NotFoundError when no entities are found.
func (aq *AllocateQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = aq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{allocate.Label}
	default:
		err = &NotSingularError{allocate.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (aq *AllocateQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := aq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Allocates.
func (aq *AllocateQuery) All(ctx context.Context) ([]*Allocate, error) {
	if err := aq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return aq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (aq *AllocateQuery) AllX(ctx context.Context) []*Allocate {
	nodes, err := aq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Allocate IDs.
func (aq *AllocateQuery) IDs(ctx context.Context) ([]uint64, error) {
	var ids []uint64
	if err := aq.Select(allocate.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (aq *AllocateQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := aq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (aq *AllocateQuery) Count(ctx context.Context) (int, error) {
	if err := aq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return aq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (aq *AllocateQuery) CountX(ctx context.Context) int {
	count, err := aq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (aq *AllocateQuery) Exist(ctx context.Context) (bool, error) {
	if err := aq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return aq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (aq *AllocateQuery) ExistX(ctx context.Context) bool {
	exist, err := aq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AllocateQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (aq *AllocateQuery) Clone() *AllocateQuery {
	if aq == nil {
		return nil
	}
	return &AllocateQuery{
		config:        aq.config,
		limit:         aq.limit,
		offset:        aq.offset,
		order:         append([]OrderFunc{}, aq.order...),
		predicates:    append([]predicate.Allocate{}, aq.predicates...),
		withRider:     aq.withRider.Clone(),
		withEmployee:  aq.withEmployee.Clone(),
		withStore:     aq.withStore.Clone(),
		withEbike:     aq.withEbike.Clone(),
		withBrand:     aq.withBrand.Clone(),
		withSubscribe: aq.withSubscribe.Clone(),
		withCabinet:   aq.withCabinet.Clone(),
		withContract:  aq.withContract.Clone(),
		// clone intermediate query.
		sql:    aq.sql.Clone(),
		path:   aq.path,
		unique: aq.unique,
	}
}

// WithRider tells the query-builder to eager-load the nodes that are connected to
// the "rider" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AllocateQuery) WithRider(opts ...func(*RiderQuery)) *AllocateQuery {
	query := &RiderQuery{config: aq.config}
	for _, opt := range opts {
		opt(query)
	}
	aq.withRider = query
	return aq
}

// WithEmployee tells the query-builder to eager-load the nodes that are connected to
// the "employee" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AllocateQuery) WithEmployee(opts ...func(*EmployeeQuery)) *AllocateQuery {
	query := &EmployeeQuery{config: aq.config}
	for _, opt := range opts {
		opt(query)
	}
	aq.withEmployee = query
	return aq
}

// WithStore tells the query-builder to eager-load the nodes that are connected to
// the "store" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AllocateQuery) WithStore(opts ...func(*StoreQuery)) *AllocateQuery {
	query := &StoreQuery{config: aq.config}
	for _, opt := range opts {
		opt(query)
	}
	aq.withStore = query
	return aq
}

// WithEbike tells the query-builder to eager-load the nodes that are connected to
// the "ebike" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AllocateQuery) WithEbike(opts ...func(*EbikeQuery)) *AllocateQuery {
	query := &EbikeQuery{config: aq.config}
	for _, opt := range opts {
		opt(query)
	}
	aq.withEbike = query
	return aq
}

// WithBrand tells the query-builder to eager-load the nodes that are connected to
// the "brand" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AllocateQuery) WithBrand(opts ...func(*EbikeBrandQuery)) *AllocateQuery {
	query := &EbikeBrandQuery{config: aq.config}
	for _, opt := range opts {
		opt(query)
	}
	aq.withBrand = query
	return aq
}

// WithSubscribe tells the query-builder to eager-load the nodes that are connected to
// the "subscribe" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AllocateQuery) WithSubscribe(opts ...func(*SubscribeQuery)) *AllocateQuery {
	query := &SubscribeQuery{config: aq.config}
	for _, opt := range opts {
		opt(query)
	}
	aq.withSubscribe = query
	return aq
}

// WithCabinet tells the query-builder to eager-load the nodes that are connected to
// the "cabinet" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AllocateQuery) WithCabinet(opts ...func(*CabinetQuery)) *AllocateQuery {
	query := &CabinetQuery{config: aq.config}
	for _, opt := range opts {
		opt(query)
	}
	aq.withCabinet = query
	return aq
}

// WithContract tells the query-builder to eager-load the nodes that are connected to
// the "contract" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AllocateQuery) WithContract(opts ...func(*ContractQuery)) *AllocateQuery {
	query := &ContractQuery{config: aq.config}
	for _, opt := range opts {
		opt(query)
	}
	aq.withContract = query
	return aq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Creator *model.Modifier `json:"creator,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Allocate.Query().
//		GroupBy(allocate.FieldCreator).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (aq *AllocateQuery) GroupBy(field string, fields ...string) *AllocateGroupBy {
	grbuild := &AllocateGroupBy{config: aq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return aq.sqlQuery(ctx), nil
	}
	grbuild.label = allocate.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Creator *model.Modifier `json:"creator,omitempty"`
//	}
//
//	client.Allocate.Query().
//		Select(allocate.FieldCreator).
//		Scan(ctx, &v)
func (aq *AllocateQuery) Select(fields ...string) *AllocateSelect {
	aq.fields = append(aq.fields, fields...)
	selbuild := &AllocateSelect{AllocateQuery: aq}
	selbuild.label = allocate.Label
	selbuild.flds, selbuild.scan = &aq.fields, selbuild.Scan
	return selbuild
}

func (aq *AllocateQuery) prepareQuery(ctx context.Context) error {
	for _, f := range aq.fields {
		if !allocate.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if aq.path != nil {
		prev, err := aq.path(ctx)
		if err != nil {
			return err
		}
		aq.sql = prev
	}
	return nil
}

func (aq *AllocateQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Allocate, error) {
	var (
		nodes       = []*Allocate{}
		_spec       = aq.querySpec()
		loadedTypes = [8]bool{
			aq.withRider != nil,
			aq.withEmployee != nil,
			aq.withStore != nil,
			aq.withEbike != nil,
			aq.withBrand != nil,
			aq.withSubscribe != nil,
			aq.withCabinet != nil,
			aq.withContract != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Allocate).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Allocate{config: aq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(aq.modifiers) > 0 {
		_spec.Modifiers = aq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, aq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := aq.withRider; query != nil {
		if err := aq.loadRider(ctx, query, nodes, nil,
			func(n *Allocate, e *Rider) { n.Edges.Rider = e }); err != nil {
			return nil, err
		}
	}
	if query := aq.withEmployee; query != nil {
		if err := aq.loadEmployee(ctx, query, nodes, nil,
			func(n *Allocate, e *Employee) { n.Edges.Employee = e }); err != nil {
			return nil, err
		}
	}
	if query := aq.withStore; query != nil {
		if err := aq.loadStore(ctx, query, nodes, nil,
			func(n *Allocate, e *Store) { n.Edges.Store = e }); err != nil {
			return nil, err
		}
	}
	if query := aq.withEbike; query != nil {
		if err := aq.loadEbike(ctx, query, nodes, nil,
			func(n *Allocate, e *Ebike) { n.Edges.Ebike = e }); err != nil {
			return nil, err
		}
	}
	if query := aq.withBrand; query != nil {
		if err := aq.loadBrand(ctx, query, nodes, nil,
			func(n *Allocate, e *EbikeBrand) { n.Edges.Brand = e }); err != nil {
			return nil, err
		}
	}
	if query := aq.withSubscribe; query != nil {
		if err := aq.loadSubscribe(ctx, query, nodes, nil,
			func(n *Allocate, e *Subscribe) { n.Edges.Subscribe = e }); err != nil {
			return nil, err
		}
	}
	if query := aq.withCabinet; query != nil {
		if err := aq.loadCabinet(ctx, query, nodes, nil,
			func(n *Allocate, e *Cabinet) { n.Edges.Cabinet = e }); err != nil {
			return nil, err
		}
	}
	if query := aq.withContract; query != nil {
		if err := aq.loadContract(ctx, query, nodes, nil,
			func(n *Allocate, e *Contract) { n.Edges.Contract = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (aq *AllocateQuery) loadRider(ctx context.Context, query *RiderQuery, nodes []*Allocate, init func(*Allocate), assign func(*Allocate, *Rider)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*Allocate)
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
func (aq *AllocateQuery) loadEmployee(ctx context.Context, query *EmployeeQuery, nodes []*Allocate, init func(*Allocate), assign func(*Allocate, *Employee)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*Allocate)
	for i := range nodes {
		if nodes[i].EmployeeID == nil {
			continue
		}
		fk := *nodes[i].EmployeeID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(employee.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "employee_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (aq *AllocateQuery) loadStore(ctx context.Context, query *StoreQuery, nodes []*Allocate, init func(*Allocate), assign func(*Allocate, *Store)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*Allocate)
	for i := range nodes {
		if nodes[i].StoreID == nil {
			continue
		}
		fk := *nodes[i].StoreID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(store.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "store_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (aq *AllocateQuery) loadEbike(ctx context.Context, query *EbikeQuery, nodes []*Allocate, init func(*Allocate), assign func(*Allocate, *Ebike)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*Allocate)
	for i := range nodes {
		if nodes[i].EbikeID == nil {
			continue
		}
		fk := *nodes[i].EbikeID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(ebike.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "ebike_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (aq *AllocateQuery) loadBrand(ctx context.Context, query *EbikeBrandQuery, nodes []*Allocate, init func(*Allocate), assign func(*Allocate, *EbikeBrand)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*Allocate)
	for i := range nodes {
		if nodes[i].BrandID == nil {
			continue
		}
		fk := *nodes[i].BrandID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(ebikebrand.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "brand_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (aq *AllocateQuery) loadSubscribe(ctx context.Context, query *SubscribeQuery, nodes []*Allocate, init func(*Allocate), assign func(*Allocate, *Subscribe)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*Allocate)
	for i := range nodes {
		fk := nodes[i].SubscribeID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(subscribe.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "subscribe_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (aq *AllocateQuery) loadCabinet(ctx context.Context, query *CabinetQuery, nodes []*Allocate, init func(*Allocate), assign func(*Allocate, *Cabinet)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*Allocate)
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
func (aq *AllocateQuery) loadContract(ctx context.Context, query *ContractQuery, nodes []*Allocate, init func(*Allocate), assign func(*Allocate, *Contract)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uint64]*Allocate)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
	}
	query.Where(predicate.Contract(func(s *sql.Selector) {
		s.Where(sql.InValues(allocate.ContractColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.AllocateID
		if fk == nil {
			return fmt.Errorf(`foreign-key "allocate_id" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "allocate_id" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (aq *AllocateQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := aq.querySpec()
	if len(aq.modifiers) > 0 {
		_spec.Modifiers = aq.modifiers
	}
	_spec.Node.Columns = aq.fields
	if len(aq.fields) > 0 {
		_spec.Unique = aq.unique != nil && *aq.unique
	}
	return sqlgraph.CountNodes(ctx, aq.driver, _spec)
}

func (aq *AllocateQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := aq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (aq *AllocateQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   allocate.Table,
			Columns: allocate.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: allocate.FieldID,
			},
		},
		From:   aq.sql,
		Unique: true,
	}
	if unique := aq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := aq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, allocate.FieldID)
		for i := range fields {
			if fields[i] != allocate.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := aq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := aq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := aq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := aq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (aq *AllocateQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(aq.driver.Dialect())
	t1 := builder.Table(allocate.Table)
	columns := aq.fields
	if len(columns) == 0 {
		columns = allocate.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if aq.sql != nil {
		selector = aq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if aq.unique != nil && *aq.unique {
		selector.Distinct()
	}
	for _, m := range aq.modifiers {
		m(selector)
	}
	for _, p := range aq.predicates {
		p(selector)
	}
	for _, p := range aq.order {
		p(selector)
	}
	if offset := aq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := aq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (aq *AllocateQuery) Modify(modifiers ...func(s *sql.Selector)) *AllocateSelect {
	aq.modifiers = append(aq.modifiers, modifiers...)
	return aq.Select()
}

// AllocateGroupBy is the group-by builder for Allocate entities.
type AllocateGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (agb *AllocateGroupBy) Aggregate(fns ...AggregateFunc) *AllocateGroupBy {
	agb.fns = append(agb.fns, fns...)
	return agb
}

// Scan applies the group-by query and scans the result into the given value.
func (agb *AllocateGroupBy) Scan(ctx context.Context, v any) error {
	query, err := agb.path(ctx)
	if err != nil {
		return err
	}
	agb.sql = query
	return agb.sqlScan(ctx, v)
}

func (agb *AllocateGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range agb.fields {
		if !allocate.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := agb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := agb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (agb *AllocateGroupBy) sqlQuery() *sql.Selector {
	selector := agb.sql.Select()
	aggregation := make([]string, 0, len(agb.fns))
	for _, fn := range agb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(agb.fields)+len(agb.fns))
		for _, f := range agb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(agb.fields...)...)
}

// AllocateSelect is the builder for selecting fields of Allocate entities.
type AllocateSelect struct {
	*AllocateQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (as *AllocateSelect) Scan(ctx context.Context, v any) error {
	if err := as.prepareQuery(ctx); err != nil {
		return err
	}
	as.sql = as.AllocateQuery.sqlQuery(ctx)
	return as.sqlScan(ctx, v)
}

func (as *AllocateSelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := as.sql.Query()
	if err := as.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (as *AllocateSelect) Modify(modifiers ...func(s *sql.Selector)) *AllocateSelect {
	as.modifiers = append(as.modifiers, modifiers...)
	return as
}