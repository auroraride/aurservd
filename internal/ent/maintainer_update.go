// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// MaintainerUpdate is the builder for updating Maintainer entities.
type MaintainerUpdate struct {
	config
	hooks     []Hook
	mutation  *MaintainerMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the MaintainerUpdate builder.
func (mu *MaintainerUpdate) Where(ps ...predicate.Maintainer) *MaintainerUpdate {
	mu.mutation.Where(ps...)
	return mu
}

// SetEnable sets the "enable" field.
func (mu *MaintainerUpdate) SetEnable(b bool) *MaintainerUpdate {
	mu.mutation.SetEnable(b)
	return mu
}

// SetNillableEnable sets the "enable" field if the given value is not nil.
func (mu *MaintainerUpdate) SetNillableEnable(b *bool) *MaintainerUpdate {
	if b != nil {
		mu.SetEnable(*b)
	}
	return mu
}

// SetName sets the "name" field.
func (mu *MaintainerUpdate) SetName(s string) *MaintainerUpdate {
	mu.mutation.SetName(s)
	return mu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (mu *MaintainerUpdate) SetNillableName(s *string) *MaintainerUpdate {
	if s != nil {
		mu.SetName(*s)
	}
	return mu
}

// SetPhone sets the "phone" field.
func (mu *MaintainerUpdate) SetPhone(s string) *MaintainerUpdate {
	mu.mutation.SetPhone(s)
	return mu
}

// SetNillablePhone sets the "phone" field if the given value is not nil.
func (mu *MaintainerUpdate) SetNillablePhone(s *string) *MaintainerUpdate {
	if s != nil {
		mu.SetPhone(*s)
	}
	return mu
}

// SetPassword sets the "password" field.
func (mu *MaintainerUpdate) SetPassword(s string) *MaintainerUpdate {
	mu.mutation.SetPassword(s)
	return mu
}

// SetNillablePassword sets the "password" field if the given value is not nil.
func (mu *MaintainerUpdate) SetNillablePassword(s *string) *MaintainerUpdate {
	if s != nil {
		mu.SetPassword(*s)
	}
	return mu
}

// AddCityIDs adds the "cities" edge to the City entity by IDs.
func (mu *MaintainerUpdate) AddCityIDs(ids ...uint64) *MaintainerUpdate {
	mu.mutation.AddCityIDs(ids...)
	return mu
}

// AddCities adds the "cities" edges to the City entity.
func (mu *MaintainerUpdate) AddCities(c ...*City) *MaintainerUpdate {
	ids := make([]uint64, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return mu.AddCityIDs(ids...)
}

// AddAssetIDs adds the "asset" edge to the Asset entity by IDs.
func (mu *MaintainerUpdate) AddAssetIDs(ids ...uint64) *MaintainerUpdate {
	mu.mutation.AddAssetIDs(ids...)
	return mu
}

// AddAsset adds the "asset" edges to the Asset entity.
func (mu *MaintainerUpdate) AddAsset(a ...*Asset) *MaintainerUpdate {
	ids := make([]uint64, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return mu.AddAssetIDs(ids...)
}

// Mutation returns the MaintainerMutation object of the builder.
func (mu *MaintainerUpdate) Mutation() *MaintainerMutation {
	return mu.mutation
}

// ClearCities clears all "cities" edges to the City entity.
func (mu *MaintainerUpdate) ClearCities() *MaintainerUpdate {
	mu.mutation.ClearCities()
	return mu
}

// RemoveCityIDs removes the "cities" edge to City entities by IDs.
func (mu *MaintainerUpdate) RemoveCityIDs(ids ...uint64) *MaintainerUpdate {
	mu.mutation.RemoveCityIDs(ids...)
	return mu
}

// RemoveCities removes "cities" edges to City entities.
func (mu *MaintainerUpdate) RemoveCities(c ...*City) *MaintainerUpdate {
	ids := make([]uint64, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return mu.RemoveCityIDs(ids...)
}

// ClearAsset clears all "asset" edges to the Asset entity.
func (mu *MaintainerUpdate) ClearAsset() *MaintainerUpdate {
	mu.mutation.ClearAsset()
	return mu
}

// RemoveAssetIDs removes the "asset" edge to Asset entities by IDs.
func (mu *MaintainerUpdate) RemoveAssetIDs(ids ...uint64) *MaintainerUpdate {
	mu.mutation.RemoveAssetIDs(ids...)
	return mu
}

// RemoveAsset removes "asset" edges to Asset entities.
func (mu *MaintainerUpdate) RemoveAsset(a ...*Asset) *MaintainerUpdate {
	ids := make([]uint64, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return mu.RemoveAssetIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mu *MaintainerUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, mu.sqlSave, mu.mutation, mu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mu *MaintainerUpdate) SaveX(ctx context.Context) int {
	affected, err := mu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mu *MaintainerUpdate) Exec(ctx context.Context) error {
	_, err := mu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mu *MaintainerUpdate) ExecX(ctx context.Context) {
	if err := mu.Exec(ctx); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (mu *MaintainerUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *MaintainerUpdate {
	mu.modifiers = append(mu.modifiers, modifiers...)
	return mu
}

func (mu *MaintainerUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(maintainer.Table, maintainer.Columns, sqlgraph.NewFieldSpec(maintainer.FieldID, field.TypeUint64))
	if ps := mu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mu.mutation.Enable(); ok {
		_spec.SetField(maintainer.FieldEnable, field.TypeBool, value)
	}
	if value, ok := mu.mutation.Name(); ok {
		_spec.SetField(maintainer.FieldName, field.TypeString, value)
	}
	if value, ok := mu.mutation.Phone(); ok {
		_spec.SetField(maintainer.FieldPhone, field.TypeString, value)
	}
	if value, ok := mu.mutation.Password(); ok {
		_spec.SetField(maintainer.FieldPassword, field.TypeString, value)
	}
	if mu.mutation.CitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   maintainer.CitiesTable,
			Columns: maintainer.CitiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.RemovedCitiesIDs(); len(nodes) > 0 && !mu.mutation.CitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   maintainer.CitiesTable,
			Columns: maintainer.CitiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.CitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   maintainer.CitiesTable,
			Columns: maintainer.CitiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if mu.mutation.AssetCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   maintainer.AssetTable,
			Columns: []string{maintainer.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.RemovedAssetIDs(); len(nodes) > 0 && !mu.mutation.AssetCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   maintainer.AssetTable,
			Columns: []string{maintainer.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.AssetIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   maintainer.AssetTable,
			Columns: []string{maintainer.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(mu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, mu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{maintainer.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	mu.mutation.done = true
	return n, nil
}

// MaintainerUpdateOne is the builder for updating a single Maintainer entity.
type MaintainerUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *MaintainerMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetEnable sets the "enable" field.
func (muo *MaintainerUpdateOne) SetEnable(b bool) *MaintainerUpdateOne {
	muo.mutation.SetEnable(b)
	return muo
}

// SetNillableEnable sets the "enable" field if the given value is not nil.
func (muo *MaintainerUpdateOne) SetNillableEnable(b *bool) *MaintainerUpdateOne {
	if b != nil {
		muo.SetEnable(*b)
	}
	return muo
}

// SetName sets the "name" field.
func (muo *MaintainerUpdateOne) SetName(s string) *MaintainerUpdateOne {
	muo.mutation.SetName(s)
	return muo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (muo *MaintainerUpdateOne) SetNillableName(s *string) *MaintainerUpdateOne {
	if s != nil {
		muo.SetName(*s)
	}
	return muo
}

// SetPhone sets the "phone" field.
func (muo *MaintainerUpdateOne) SetPhone(s string) *MaintainerUpdateOne {
	muo.mutation.SetPhone(s)
	return muo
}

// SetNillablePhone sets the "phone" field if the given value is not nil.
func (muo *MaintainerUpdateOne) SetNillablePhone(s *string) *MaintainerUpdateOne {
	if s != nil {
		muo.SetPhone(*s)
	}
	return muo
}

// SetPassword sets the "password" field.
func (muo *MaintainerUpdateOne) SetPassword(s string) *MaintainerUpdateOne {
	muo.mutation.SetPassword(s)
	return muo
}

// SetNillablePassword sets the "password" field if the given value is not nil.
func (muo *MaintainerUpdateOne) SetNillablePassword(s *string) *MaintainerUpdateOne {
	if s != nil {
		muo.SetPassword(*s)
	}
	return muo
}

// AddCityIDs adds the "cities" edge to the City entity by IDs.
func (muo *MaintainerUpdateOne) AddCityIDs(ids ...uint64) *MaintainerUpdateOne {
	muo.mutation.AddCityIDs(ids...)
	return muo
}

// AddCities adds the "cities" edges to the City entity.
func (muo *MaintainerUpdateOne) AddCities(c ...*City) *MaintainerUpdateOne {
	ids := make([]uint64, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return muo.AddCityIDs(ids...)
}

// AddAssetIDs adds the "asset" edge to the Asset entity by IDs.
func (muo *MaintainerUpdateOne) AddAssetIDs(ids ...uint64) *MaintainerUpdateOne {
	muo.mutation.AddAssetIDs(ids...)
	return muo
}

// AddAsset adds the "asset" edges to the Asset entity.
func (muo *MaintainerUpdateOne) AddAsset(a ...*Asset) *MaintainerUpdateOne {
	ids := make([]uint64, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return muo.AddAssetIDs(ids...)
}

// Mutation returns the MaintainerMutation object of the builder.
func (muo *MaintainerUpdateOne) Mutation() *MaintainerMutation {
	return muo.mutation
}

// ClearCities clears all "cities" edges to the City entity.
func (muo *MaintainerUpdateOne) ClearCities() *MaintainerUpdateOne {
	muo.mutation.ClearCities()
	return muo
}

// RemoveCityIDs removes the "cities" edge to City entities by IDs.
func (muo *MaintainerUpdateOne) RemoveCityIDs(ids ...uint64) *MaintainerUpdateOne {
	muo.mutation.RemoveCityIDs(ids...)
	return muo
}

// RemoveCities removes "cities" edges to City entities.
func (muo *MaintainerUpdateOne) RemoveCities(c ...*City) *MaintainerUpdateOne {
	ids := make([]uint64, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return muo.RemoveCityIDs(ids...)
}

// ClearAsset clears all "asset" edges to the Asset entity.
func (muo *MaintainerUpdateOne) ClearAsset() *MaintainerUpdateOne {
	muo.mutation.ClearAsset()
	return muo
}

// RemoveAssetIDs removes the "asset" edge to Asset entities by IDs.
func (muo *MaintainerUpdateOne) RemoveAssetIDs(ids ...uint64) *MaintainerUpdateOne {
	muo.mutation.RemoveAssetIDs(ids...)
	return muo
}

// RemoveAsset removes "asset" edges to Asset entities.
func (muo *MaintainerUpdateOne) RemoveAsset(a ...*Asset) *MaintainerUpdateOne {
	ids := make([]uint64, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return muo.RemoveAssetIDs(ids...)
}

// Where appends a list predicates to the MaintainerUpdate builder.
func (muo *MaintainerUpdateOne) Where(ps ...predicate.Maintainer) *MaintainerUpdateOne {
	muo.mutation.Where(ps...)
	return muo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (muo *MaintainerUpdateOne) Select(field string, fields ...string) *MaintainerUpdateOne {
	muo.fields = append([]string{field}, fields...)
	return muo
}

// Save executes the query and returns the updated Maintainer entity.
func (muo *MaintainerUpdateOne) Save(ctx context.Context) (*Maintainer, error) {
	return withHooks(ctx, muo.sqlSave, muo.mutation, muo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (muo *MaintainerUpdateOne) SaveX(ctx context.Context) *Maintainer {
	node, err := muo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (muo *MaintainerUpdateOne) Exec(ctx context.Context) error {
	_, err := muo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (muo *MaintainerUpdateOne) ExecX(ctx context.Context) {
	if err := muo.Exec(ctx); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (muo *MaintainerUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *MaintainerUpdateOne {
	muo.modifiers = append(muo.modifiers, modifiers...)
	return muo
}

func (muo *MaintainerUpdateOne) sqlSave(ctx context.Context) (_node *Maintainer, err error) {
	_spec := sqlgraph.NewUpdateSpec(maintainer.Table, maintainer.Columns, sqlgraph.NewFieldSpec(maintainer.FieldID, field.TypeUint64))
	id, ok := muo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Maintainer.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := muo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, maintainer.FieldID)
		for _, f := range fields {
			if !maintainer.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != maintainer.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := muo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := muo.mutation.Enable(); ok {
		_spec.SetField(maintainer.FieldEnable, field.TypeBool, value)
	}
	if value, ok := muo.mutation.Name(); ok {
		_spec.SetField(maintainer.FieldName, field.TypeString, value)
	}
	if value, ok := muo.mutation.Phone(); ok {
		_spec.SetField(maintainer.FieldPhone, field.TypeString, value)
	}
	if value, ok := muo.mutation.Password(); ok {
		_spec.SetField(maintainer.FieldPassword, field.TypeString, value)
	}
	if muo.mutation.CitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   maintainer.CitiesTable,
			Columns: maintainer.CitiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.RemovedCitiesIDs(); len(nodes) > 0 && !muo.mutation.CitiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   maintainer.CitiesTable,
			Columns: maintainer.CitiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.CitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   maintainer.CitiesTable,
			Columns: maintainer.CitiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if muo.mutation.AssetCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   maintainer.AssetTable,
			Columns: []string{maintainer.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.RemovedAssetIDs(); len(nodes) > 0 && !muo.mutation.AssetCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   maintainer.AssetTable,
			Columns: []string{maintainer.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.AssetIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   maintainer.AssetTable,
			Columns: []string{maintainer.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(muo.modifiers...)
	_node = &Maintainer{config: muo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, muo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{maintainer.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	muo.mutation.done = true
	return _node, nil
}
