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
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
)

// WarehouseQuery is the builder for querying Warehouse entities.
type WarehouseQuery struct {
	config
	ctx                     *QueryContext
	order                   []warehouse.OrderOption
	inters                  []Interceptor
	predicates              []predicate.Warehouse
	withCity                *CityQuery
	withBelongAssetManagers *AssetManagerQuery
	withDutyAssetManagers   *AssetManagerQuery
	withAsset               *AssetQuery
	modifiers               []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the WarehouseQuery builder.
func (wq *WarehouseQuery) Where(ps ...predicate.Warehouse) *WarehouseQuery {
	wq.predicates = append(wq.predicates, ps...)
	return wq
}

// Limit the number of records to be returned by this query.
func (wq *WarehouseQuery) Limit(limit int) *WarehouseQuery {
	wq.ctx.Limit = &limit
	return wq
}

// Offset to start from.
func (wq *WarehouseQuery) Offset(offset int) *WarehouseQuery {
	wq.ctx.Offset = &offset
	return wq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (wq *WarehouseQuery) Unique(unique bool) *WarehouseQuery {
	wq.ctx.Unique = &unique
	return wq
}

// Order specifies how the records should be ordered.
func (wq *WarehouseQuery) Order(o ...warehouse.OrderOption) *WarehouseQuery {
	wq.order = append(wq.order, o...)
	return wq
}

// QueryCity chains the current query on the "city" edge.
func (wq *WarehouseQuery) QueryCity() *CityQuery {
	query := (&CityClient{config: wq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := wq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := wq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(warehouse.Table, warehouse.FieldID, selector),
			sqlgraph.To(city.Table, city.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, warehouse.CityTable, warehouse.CityColumn),
		)
		fromU = sqlgraph.SetNeighbors(wq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryBelongAssetManagers chains the current query on the "belong_asset_managers" edge.
func (wq *WarehouseQuery) QueryBelongAssetManagers() *AssetManagerQuery {
	query := (&AssetManagerClient{config: wq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := wq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := wq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(warehouse.Table, warehouse.FieldID, selector),
			sqlgraph.To(assetmanager.Table, assetmanager.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, warehouse.BelongAssetManagersTable, warehouse.BelongAssetManagersPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(wq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryDutyAssetManagers chains the current query on the "duty_asset_managers" edge.
func (wq *WarehouseQuery) QueryDutyAssetManagers() *AssetManagerQuery {
	query := (&AssetManagerClient{config: wq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := wq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := wq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(warehouse.Table, warehouse.FieldID, selector),
			sqlgraph.To(assetmanager.Table, assetmanager.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, warehouse.DutyAssetManagersTable, warehouse.DutyAssetManagersColumn),
		)
		fromU = sqlgraph.SetNeighbors(wq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryAsset chains the current query on the "asset" edge.
func (wq *WarehouseQuery) QueryAsset() *AssetQuery {
	query := (&AssetClient{config: wq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := wq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := wq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(warehouse.Table, warehouse.FieldID, selector),
			sqlgraph.To(asset.Table, asset.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, warehouse.AssetTable, warehouse.AssetColumn),
		)
		fromU = sqlgraph.SetNeighbors(wq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Warehouse entity from the query.
// Returns a *NotFoundError when no Warehouse was found.
func (wq *WarehouseQuery) First(ctx context.Context) (*Warehouse, error) {
	nodes, err := wq.Limit(1).All(setContextOp(ctx, wq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{warehouse.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (wq *WarehouseQuery) FirstX(ctx context.Context) *Warehouse {
	node, err := wq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Warehouse ID from the query.
// Returns a *NotFoundError when no Warehouse ID was found.
func (wq *WarehouseQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = wq.Limit(1).IDs(setContextOp(ctx, wq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{warehouse.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (wq *WarehouseQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := wq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Warehouse entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Warehouse entity is found.
// Returns a *NotFoundError when no Warehouse entities are found.
func (wq *WarehouseQuery) Only(ctx context.Context) (*Warehouse, error) {
	nodes, err := wq.Limit(2).All(setContextOp(ctx, wq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{warehouse.Label}
	default:
		return nil, &NotSingularError{warehouse.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (wq *WarehouseQuery) OnlyX(ctx context.Context) *Warehouse {
	node, err := wq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Warehouse ID in the query.
// Returns a *NotSingularError when more than one Warehouse ID is found.
// Returns a *NotFoundError when no entities are found.
func (wq *WarehouseQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = wq.Limit(2).IDs(setContextOp(ctx, wq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{warehouse.Label}
	default:
		err = &NotSingularError{warehouse.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (wq *WarehouseQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := wq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Warehouses.
func (wq *WarehouseQuery) All(ctx context.Context) ([]*Warehouse, error) {
	ctx = setContextOp(ctx, wq.ctx, ent.OpQueryAll)
	if err := wq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Warehouse, *WarehouseQuery]()
	return withInterceptors[[]*Warehouse](ctx, wq, qr, wq.inters)
}

// AllX is like All, but panics if an error occurs.
func (wq *WarehouseQuery) AllX(ctx context.Context) []*Warehouse {
	nodes, err := wq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Warehouse IDs.
func (wq *WarehouseQuery) IDs(ctx context.Context) (ids []uint64, err error) {
	if wq.ctx.Unique == nil && wq.path != nil {
		wq.Unique(true)
	}
	ctx = setContextOp(ctx, wq.ctx, ent.OpQueryIDs)
	if err = wq.Select(warehouse.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (wq *WarehouseQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := wq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (wq *WarehouseQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, wq.ctx, ent.OpQueryCount)
	if err := wq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, wq, querierCount[*WarehouseQuery](), wq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (wq *WarehouseQuery) CountX(ctx context.Context) int {
	count, err := wq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (wq *WarehouseQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, wq.ctx, ent.OpQueryExist)
	switch _, err := wq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (wq *WarehouseQuery) ExistX(ctx context.Context) bool {
	exist, err := wq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the WarehouseQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (wq *WarehouseQuery) Clone() *WarehouseQuery {
	if wq == nil {
		return nil
	}
	return &WarehouseQuery{
		config:                  wq.config,
		ctx:                     wq.ctx.Clone(),
		order:                   append([]warehouse.OrderOption{}, wq.order...),
		inters:                  append([]Interceptor{}, wq.inters...),
		predicates:              append([]predicate.Warehouse{}, wq.predicates...),
		withCity:                wq.withCity.Clone(),
		withBelongAssetManagers: wq.withBelongAssetManagers.Clone(),
		withDutyAssetManagers:   wq.withDutyAssetManagers.Clone(),
		withAsset:               wq.withAsset.Clone(),
		// clone intermediate query.
		sql:       wq.sql.Clone(),
		path:      wq.path,
		modifiers: append([]func(*sql.Selector){}, wq.modifiers...),
	}
}

// WithCity tells the query-builder to eager-load the nodes that are connected to
// the "city" edge. The optional arguments are used to configure the query builder of the edge.
func (wq *WarehouseQuery) WithCity(opts ...func(*CityQuery)) *WarehouseQuery {
	query := (&CityClient{config: wq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	wq.withCity = query
	return wq
}

// WithBelongAssetManagers tells the query-builder to eager-load the nodes that are connected to
// the "belong_asset_managers" edge. The optional arguments are used to configure the query builder of the edge.
func (wq *WarehouseQuery) WithBelongAssetManagers(opts ...func(*AssetManagerQuery)) *WarehouseQuery {
	query := (&AssetManagerClient{config: wq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	wq.withBelongAssetManagers = query
	return wq
}

// WithDutyAssetManagers tells the query-builder to eager-load the nodes that are connected to
// the "duty_asset_managers" edge. The optional arguments are used to configure the query builder of the edge.
func (wq *WarehouseQuery) WithDutyAssetManagers(opts ...func(*AssetManagerQuery)) *WarehouseQuery {
	query := (&AssetManagerClient{config: wq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	wq.withDutyAssetManagers = query
	return wq
}

// WithAsset tells the query-builder to eager-load the nodes that are connected to
// the "asset" edge. The optional arguments are used to configure the query builder of the edge.
func (wq *WarehouseQuery) WithAsset(opts ...func(*AssetQuery)) *WarehouseQuery {
	query := (&AssetClient{config: wq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	wq.withAsset = query
	return wq
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
//	client.Warehouse.Query().
//		GroupBy(warehouse.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (wq *WarehouseQuery) GroupBy(field string, fields ...string) *WarehouseGroupBy {
	wq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &WarehouseGroupBy{build: wq}
	grbuild.flds = &wq.ctx.Fields
	grbuild.label = warehouse.Label
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
//	client.Warehouse.Query().
//		Select(warehouse.FieldCreatedAt).
//		Scan(ctx, &v)
func (wq *WarehouseQuery) Select(fields ...string) *WarehouseSelect {
	wq.ctx.Fields = append(wq.ctx.Fields, fields...)
	sbuild := &WarehouseSelect{WarehouseQuery: wq}
	sbuild.label = warehouse.Label
	sbuild.flds, sbuild.scan = &wq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a WarehouseSelect configured with the given aggregations.
func (wq *WarehouseQuery) Aggregate(fns ...AggregateFunc) *WarehouseSelect {
	return wq.Select().Aggregate(fns...)
}

func (wq *WarehouseQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range wq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, wq); err != nil {
				return err
			}
		}
	}
	for _, f := range wq.ctx.Fields {
		if !warehouse.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if wq.path != nil {
		prev, err := wq.path(ctx)
		if err != nil {
			return err
		}
		wq.sql = prev
	}
	return nil
}

func (wq *WarehouseQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Warehouse, error) {
	var (
		nodes       = []*Warehouse{}
		_spec       = wq.querySpec()
		loadedTypes = [4]bool{
			wq.withCity != nil,
			wq.withBelongAssetManagers != nil,
			wq.withDutyAssetManagers != nil,
			wq.withAsset != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Warehouse).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Warehouse{config: wq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(wq.modifiers) > 0 {
		_spec.Modifiers = wq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, wq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := wq.withCity; query != nil {
		if err := wq.loadCity(ctx, query, nodes, nil,
			func(n *Warehouse, e *City) { n.Edges.City = e }); err != nil {
			return nil, err
		}
	}
	if query := wq.withBelongAssetManagers; query != nil {
		if err := wq.loadBelongAssetManagers(ctx, query, nodes,
			func(n *Warehouse) { n.Edges.BelongAssetManagers = []*AssetManager{} },
			func(n *Warehouse, e *AssetManager) {
				n.Edges.BelongAssetManagers = append(n.Edges.BelongAssetManagers, e)
			}); err != nil {
			return nil, err
		}
	}
	if query := wq.withDutyAssetManagers; query != nil {
		if err := wq.loadDutyAssetManagers(ctx, query, nodes,
			func(n *Warehouse) { n.Edges.DutyAssetManagers = []*AssetManager{} },
			func(n *Warehouse, e *AssetManager) { n.Edges.DutyAssetManagers = append(n.Edges.DutyAssetManagers, e) }); err != nil {
			return nil, err
		}
	}
	if query := wq.withAsset; query != nil {
		if err := wq.loadAsset(ctx, query, nodes,
			func(n *Warehouse) { n.Edges.Asset = []*Asset{} },
			func(n *Warehouse, e *Asset) { n.Edges.Asset = append(n.Edges.Asset, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (wq *WarehouseQuery) loadCity(ctx context.Context, query *CityQuery, nodes []*Warehouse, init func(*Warehouse), assign func(*Warehouse, *City)) error {
	ids := make([]uint64, 0, len(nodes))
	nodeids := make(map[uint64][]*Warehouse)
	for i := range nodes {
		fk := nodes[i].CityID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(city.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "city_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (wq *WarehouseQuery) loadBelongAssetManagers(ctx context.Context, query *AssetManagerQuery, nodes []*Warehouse, init func(*Warehouse), assign func(*Warehouse, *AssetManager)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[uint64]*Warehouse)
	nids := make(map[uint64]map[*Warehouse]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(warehouse.BelongAssetManagersTable)
		s.Join(joinT).On(s.C(assetmanager.FieldID), joinT.C(warehouse.BelongAssetManagersPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(warehouse.BelongAssetManagersPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(warehouse.BelongAssetManagersPrimaryKey[0]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := uint64(values[0].(*sql.NullInt64).Int64)
				inValue := uint64(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*Warehouse]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*AssetManager](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "belong_asset_managers" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}
func (wq *WarehouseQuery) loadDutyAssetManagers(ctx context.Context, query *AssetManagerQuery, nodes []*Warehouse, init func(*Warehouse), assign func(*Warehouse, *AssetManager)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uint64]*Warehouse)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(assetmanager.FieldWarehouseID)
	}
	query.Where(predicate.AssetManager(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(warehouse.DutyAssetManagersColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.WarehouseID
		if fk == nil {
			return fmt.Errorf(`foreign-key "warehouse_id" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "warehouse_id" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (wq *WarehouseQuery) loadAsset(ctx context.Context, query *AssetQuery, nodes []*Warehouse, init func(*Warehouse), assign func(*Warehouse, *Asset)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uint64]*Warehouse)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(asset.FieldLocationsID)
	}
	query.Where(predicate.Asset(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(warehouse.AssetColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.LocationsID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "locations_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (wq *WarehouseQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := wq.querySpec()
	if len(wq.modifiers) > 0 {
		_spec.Modifiers = wq.modifiers
	}
	_spec.Node.Columns = wq.ctx.Fields
	if len(wq.ctx.Fields) > 0 {
		_spec.Unique = wq.ctx.Unique != nil && *wq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, wq.driver, _spec)
}

func (wq *WarehouseQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(warehouse.Table, warehouse.Columns, sqlgraph.NewFieldSpec(warehouse.FieldID, field.TypeUint64))
	_spec.From = wq.sql
	if unique := wq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if wq.path != nil {
		_spec.Unique = true
	}
	if fields := wq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, warehouse.FieldID)
		for i := range fields {
			if fields[i] != warehouse.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if wq.withCity != nil {
			_spec.Node.AddColumnOnce(warehouse.FieldCityID)
		}
	}
	if ps := wq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := wq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := wq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := wq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (wq *WarehouseQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(wq.driver.Dialect())
	t1 := builder.Table(warehouse.Table)
	columns := wq.ctx.Fields
	if len(columns) == 0 {
		columns = warehouse.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if wq.sql != nil {
		selector = wq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if wq.ctx.Unique != nil && *wq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range wq.modifiers {
		m(selector)
	}
	for _, p := range wq.predicates {
		p(selector)
	}
	for _, p := range wq.order {
		p(selector)
	}
	if offset := wq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := wq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (wq *WarehouseQuery) Modify(modifiers ...func(s *sql.Selector)) *WarehouseSelect {
	wq.modifiers = append(wq.modifiers, modifiers...)
	return wq.Select()
}

type WarehouseQueryWith string

var (
	WarehouseQueryWithCity                WarehouseQueryWith = "City"
	WarehouseQueryWithBelongAssetManagers WarehouseQueryWith = "BelongAssetManagers"
	WarehouseQueryWithDutyAssetManagers   WarehouseQueryWith = "DutyAssetManagers"
	WarehouseQueryWithAsset               WarehouseQueryWith = "Asset"
)

func (wq *WarehouseQuery) With(withEdges ...WarehouseQueryWith) *WarehouseQuery {
	for _, v := range withEdges {
		switch v {
		case WarehouseQueryWithCity:
			wq.WithCity()
		case WarehouseQueryWithBelongAssetManagers:
			wq.WithBelongAssetManagers()
		case WarehouseQueryWithDutyAssetManagers:
			wq.WithDutyAssetManagers()
		case WarehouseQueryWithAsset:
			wq.WithAsset()
		}
	}
	return wq
}

// WarehouseGroupBy is the group-by builder for Warehouse entities.
type WarehouseGroupBy struct {
	selector
	build *WarehouseQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (wgb *WarehouseGroupBy) Aggregate(fns ...AggregateFunc) *WarehouseGroupBy {
	wgb.fns = append(wgb.fns, fns...)
	return wgb
}

// Scan applies the selector query and scans the result into the given value.
func (wgb *WarehouseGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, wgb.build.ctx, ent.OpQueryGroupBy)
	if err := wgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*WarehouseQuery, *WarehouseGroupBy](ctx, wgb.build, wgb, wgb.build.inters, v)
}

func (wgb *WarehouseGroupBy) sqlScan(ctx context.Context, root *WarehouseQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(wgb.fns))
	for _, fn := range wgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*wgb.flds)+len(wgb.fns))
		for _, f := range *wgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*wgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := wgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// WarehouseSelect is the builder for selecting fields of Warehouse entities.
type WarehouseSelect struct {
	*WarehouseQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ws *WarehouseSelect) Aggregate(fns ...AggregateFunc) *WarehouseSelect {
	ws.fns = append(ws.fns, fns...)
	return ws
}

// Scan applies the selector query and scans the result into the given value.
func (ws *WarehouseSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ws.ctx, ent.OpQuerySelect)
	if err := ws.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*WarehouseQuery, *WarehouseSelect](ctx, ws.WarehouseQuery, ws, ws.inters, v)
}

func (ws *WarehouseSelect) sqlScan(ctx context.Context, root *WarehouseQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ws.fns))
	for _, fn := range ws.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ws.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ws.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ws *WarehouseSelect) Modify(modifiers ...func(s *sql.Selector)) *WarehouseSelect {
	ws.modifiers = append(ws.modifiers, modifiers...)
	return ws
}
