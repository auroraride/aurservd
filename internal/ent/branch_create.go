// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/branch"
	"github.com/auroraride/aurservd/internal/ent/branchcontract"
)

// BranchCreate is the builder for creating a Branch entity.
type BranchCreate struct {
	config
	mutation *BranchMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (bc *BranchCreate) SetCreatedAt(t time.Time) *BranchCreate {
	bc.mutation.SetCreatedAt(t)
	return bc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (bc *BranchCreate) SetNillableCreatedAt(t *time.Time) *BranchCreate {
	if t != nil {
		bc.SetCreatedAt(*t)
	}
	return bc
}

// SetUpdatedAt sets the "updated_at" field.
func (bc *BranchCreate) SetUpdatedAt(t time.Time) *BranchCreate {
	bc.mutation.SetUpdatedAt(t)
	return bc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (bc *BranchCreate) SetNillableUpdatedAt(t *time.Time) *BranchCreate {
	if t != nil {
		bc.SetUpdatedAt(*t)
	}
	return bc
}

// SetDeletedAt sets the "deleted_at" field.
func (bc *BranchCreate) SetDeletedAt(t time.Time) *BranchCreate {
	bc.mutation.SetDeletedAt(t)
	return bc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (bc *BranchCreate) SetNillableDeletedAt(t *time.Time) *BranchCreate {
	if t != nil {
		bc.SetDeletedAt(*t)
	}
	return bc
}

// SetCreator sets the "creator" field.
func (bc *BranchCreate) SetCreator(m *model.Modifier) *BranchCreate {
	bc.mutation.SetCreator(m)
	return bc
}

// SetLastModifier sets the "last_modifier" field.
func (bc *BranchCreate) SetLastModifier(m *model.Modifier) *BranchCreate {
	bc.mutation.SetLastModifier(m)
	return bc
}

// SetRemark sets the "remark" field.
func (bc *BranchCreate) SetRemark(s string) *BranchCreate {
	bc.mutation.SetRemark(s)
	return bc
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (bc *BranchCreate) SetNillableRemark(s *string) *BranchCreate {
	if s != nil {
		bc.SetRemark(*s)
	}
	return bc
}

// SetCityID sets the "city_id" field.
func (bc *BranchCreate) SetCityID(u uint64) *BranchCreate {
	bc.mutation.SetCityID(u)
	return bc
}

// SetName sets the "name" field.
func (bc *BranchCreate) SetName(s string) *BranchCreate {
	bc.mutation.SetName(s)
	return bc
}

// SetLng sets the "lng" field.
func (bc *BranchCreate) SetLng(f float64) *BranchCreate {
	bc.mutation.SetLng(f)
	return bc
}

// SetLat sets the "lat" field.
func (bc *BranchCreate) SetLat(f float64) *BranchCreate {
	bc.mutation.SetLat(f)
	return bc
}

// SetAddress sets the "address" field.
func (bc *BranchCreate) SetAddress(s string) *BranchCreate {
	bc.mutation.SetAddress(s)
	return bc
}

// SetPhotos sets the "photos" field.
func (bc *BranchCreate) SetPhotos(s []string) *BranchCreate {
	bc.mutation.SetPhotos(s)
	return bc
}

// AddContractIDs adds the "contracts" edge to the BranchContract entity by IDs.
func (bc *BranchCreate) AddContractIDs(ids ...uint64) *BranchCreate {
	bc.mutation.AddContractIDs(ids...)
	return bc
}

// AddContracts adds the "contracts" edges to the BranchContract entity.
func (bc *BranchCreate) AddContracts(b ...*BranchContract) *BranchCreate {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return bc.AddContractIDs(ids...)
}

// Mutation returns the BranchMutation object of the builder.
func (bc *BranchCreate) Mutation() *BranchMutation {
	return bc.mutation
}

// Save creates the Branch in the database.
func (bc *BranchCreate) Save(ctx context.Context) (*Branch, error) {
	var (
		err  error
		node *Branch
	)
	bc.defaults()
	if len(bc.hooks) == 0 {
		if err = bc.check(); err != nil {
			return nil, err
		}
		node, err = bc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*BranchMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = bc.check(); err != nil {
				return nil, err
			}
			bc.mutation = mutation
			if node, err = bc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(bc.hooks) - 1; i >= 0; i-- {
			if bc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = bc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, bc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (bc *BranchCreate) SaveX(ctx context.Context) *Branch {
	v, err := bc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bc *BranchCreate) Exec(ctx context.Context) error {
	_, err := bc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bc *BranchCreate) ExecX(ctx context.Context) {
	if err := bc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (bc *BranchCreate) defaults() {
	if _, ok := bc.mutation.CreatedAt(); !ok {
		v := branch.DefaultCreatedAt()
		bc.mutation.SetCreatedAt(v)
	}
	if _, ok := bc.mutation.UpdatedAt(); !ok {
		v := branch.DefaultUpdatedAt()
		bc.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bc *BranchCreate) check() error {
	if _, ok := bc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Branch.created_at"`)}
	}
	if _, ok := bc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Branch.updated_at"`)}
	}
	if _, ok := bc.mutation.CityID(); !ok {
		return &ValidationError{Name: "city_id", err: errors.New(`ent: missing required field "Branch.city_id"`)}
	}
	if _, ok := bc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Branch.name"`)}
	}
	if _, ok := bc.mutation.Lng(); !ok {
		return &ValidationError{Name: "lng", err: errors.New(`ent: missing required field "Branch.lng"`)}
	}
	if _, ok := bc.mutation.Lat(); !ok {
		return &ValidationError{Name: "lat", err: errors.New(`ent: missing required field "Branch.lat"`)}
	}
	if _, ok := bc.mutation.Address(); !ok {
		return &ValidationError{Name: "address", err: errors.New(`ent: missing required field "Branch.address"`)}
	}
	if _, ok := bc.mutation.Photos(); !ok {
		return &ValidationError{Name: "photos", err: errors.New(`ent: missing required field "Branch.photos"`)}
	}
	return nil
}

func (bc *BranchCreate) sqlSave(ctx context.Context) (*Branch, error) {
	_node, _spec := bc.createSpec()
	if err := sqlgraph.CreateNode(ctx, bc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	return _node, nil
}

func (bc *BranchCreate) createSpec() (*Branch, *sqlgraph.CreateSpec) {
	var (
		_node = &Branch{config: bc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: branch.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: branch.FieldID,
			},
		}
	)
	_spec.OnConflict = bc.conflict
	if value, ok := bc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: branch.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := bc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: branch.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := bc.mutation.DeletedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: branch.FieldDeletedAt,
		})
		_node.DeletedAt = &value
	}
	if value, ok := bc.mutation.Creator(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: branch.FieldCreator,
		})
		_node.Creator = value
	}
	if value, ok := bc.mutation.LastModifier(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: branch.FieldLastModifier,
		})
		_node.LastModifier = value
	}
	if value, ok := bc.mutation.Remark(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: branch.FieldRemark,
		})
		_node.Remark = &value
	}
	if value, ok := bc.mutation.CityID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint64,
			Value:  value,
			Column: branch.FieldCityID,
		})
		_node.CityID = value
	}
	if value, ok := bc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: branch.FieldName,
		})
		_node.Name = value
	}
	if value, ok := bc.mutation.Lng(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: branch.FieldLng,
		})
		_node.Lng = value
	}
	if value, ok := bc.mutation.Lat(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: branch.FieldLat,
		})
		_node.Lat = value
	}
	if value, ok := bc.mutation.Address(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: branch.FieldAddress,
		})
		_node.Address = value
	}
	if value, ok := bc.mutation.Photos(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: branch.FieldPhotos,
		})
		_node.Photos = value
	}
	if nodes := bc.mutation.ContractsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   branch.ContractsTable,
			Columns: []string{branch.ContractsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branchcontract.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Branch.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.BranchUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
//
func (bc *BranchCreate) OnConflict(opts ...sql.ConflictOption) *BranchUpsertOne {
	bc.conflict = opts
	return &BranchUpsertOne{
		create: bc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Branch.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (bc *BranchCreate) OnConflictColumns(columns ...string) *BranchUpsertOne {
	bc.conflict = append(bc.conflict, sql.ConflictColumns(columns...))
	return &BranchUpsertOne{
		create: bc,
	}
}

type (
	// BranchUpsertOne is the builder for "upsert"-ing
	//  one Branch node.
	BranchUpsertOne struct {
		create *BranchCreate
	}

	// BranchUpsert is the "OnConflict" setter.
	BranchUpsert struct {
		*sql.UpdateSet
	}
)

// SetCreatedAt sets the "created_at" field.
func (u *BranchUpsert) SetCreatedAt(v time.Time) *BranchUpsert {
	u.Set(branch.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *BranchUpsert) UpdateCreatedAt() *BranchUpsert {
	u.SetExcluded(branch.FieldCreatedAt)
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *BranchUpsert) SetUpdatedAt(v time.Time) *BranchUpsert {
	u.Set(branch.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *BranchUpsert) UpdateUpdatedAt() *BranchUpsert {
	u.SetExcluded(branch.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *BranchUpsert) SetDeletedAt(v time.Time) *BranchUpsert {
	u.Set(branch.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *BranchUpsert) UpdateDeletedAt() *BranchUpsert {
	u.SetExcluded(branch.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *BranchUpsert) ClearDeletedAt() *BranchUpsert {
	u.SetNull(branch.FieldDeletedAt)
	return u
}

// SetCreator sets the "creator" field.
func (u *BranchUpsert) SetCreator(v *model.Modifier) *BranchUpsert {
	u.Set(branch.FieldCreator, v)
	return u
}

// UpdateCreator sets the "creator" field to the value that was provided on create.
func (u *BranchUpsert) UpdateCreator() *BranchUpsert {
	u.SetExcluded(branch.FieldCreator)
	return u
}

// ClearCreator clears the value of the "creator" field.
func (u *BranchUpsert) ClearCreator() *BranchUpsert {
	u.SetNull(branch.FieldCreator)
	return u
}

// SetLastModifier sets the "last_modifier" field.
func (u *BranchUpsert) SetLastModifier(v *model.Modifier) *BranchUpsert {
	u.Set(branch.FieldLastModifier, v)
	return u
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *BranchUpsert) UpdateLastModifier() *BranchUpsert {
	u.SetExcluded(branch.FieldLastModifier)
	return u
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *BranchUpsert) ClearLastModifier() *BranchUpsert {
	u.SetNull(branch.FieldLastModifier)
	return u
}

// SetRemark sets the "remark" field.
func (u *BranchUpsert) SetRemark(v string) *BranchUpsert {
	u.Set(branch.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *BranchUpsert) UpdateRemark() *BranchUpsert {
	u.SetExcluded(branch.FieldRemark)
	return u
}

// ClearRemark clears the value of the "remark" field.
func (u *BranchUpsert) ClearRemark() *BranchUpsert {
	u.SetNull(branch.FieldRemark)
	return u
}

// SetCityID sets the "city_id" field.
func (u *BranchUpsert) SetCityID(v uint64) *BranchUpsert {
	u.Set(branch.FieldCityID, v)
	return u
}

// UpdateCityID sets the "city_id" field to the value that was provided on create.
func (u *BranchUpsert) UpdateCityID() *BranchUpsert {
	u.SetExcluded(branch.FieldCityID)
	return u
}

// AddCityID adds v to the "city_id" field.
func (u *BranchUpsert) AddCityID(v uint64) *BranchUpsert {
	u.Add(branch.FieldCityID, v)
	return u
}

// SetName sets the "name" field.
func (u *BranchUpsert) SetName(v string) *BranchUpsert {
	u.Set(branch.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *BranchUpsert) UpdateName() *BranchUpsert {
	u.SetExcluded(branch.FieldName)
	return u
}

// SetLng sets the "lng" field.
func (u *BranchUpsert) SetLng(v float64) *BranchUpsert {
	u.Set(branch.FieldLng, v)
	return u
}

// UpdateLng sets the "lng" field to the value that was provided on create.
func (u *BranchUpsert) UpdateLng() *BranchUpsert {
	u.SetExcluded(branch.FieldLng)
	return u
}

// AddLng adds v to the "lng" field.
func (u *BranchUpsert) AddLng(v float64) *BranchUpsert {
	u.Add(branch.FieldLng, v)
	return u
}

// SetLat sets the "lat" field.
func (u *BranchUpsert) SetLat(v float64) *BranchUpsert {
	u.Set(branch.FieldLat, v)
	return u
}

// UpdateLat sets the "lat" field to the value that was provided on create.
func (u *BranchUpsert) UpdateLat() *BranchUpsert {
	u.SetExcluded(branch.FieldLat)
	return u
}

// AddLat adds v to the "lat" field.
func (u *BranchUpsert) AddLat(v float64) *BranchUpsert {
	u.Add(branch.FieldLat, v)
	return u
}

// SetAddress sets the "address" field.
func (u *BranchUpsert) SetAddress(v string) *BranchUpsert {
	u.Set(branch.FieldAddress, v)
	return u
}

// UpdateAddress sets the "address" field to the value that was provided on create.
func (u *BranchUpsert) UpdateAddress() *BranchUpsert {
	u.SetExcluded(branch.FieldAddress)
	return u
}

// SetPhotos sets the "photos" field.
func (u *BranchUpsert) SetPhotos(v []string) *BranchUpsert {
	u.Set(branch.FieldPhotos, v)
	return u
}

// UpdatePhotos sets the "photos" field to the value that was provided on create.
func (u *BranchUpsert) UpdatePhotos() *BranchUpsert {
	u.SetExcluded(branch.FieldPhotos)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Branch.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
//
func (u *BranchUpsertOne) UpdateNewValues() *BranchUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(branch.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.Branch.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *BranchUpsertOne) Ignore() *BranchUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *BranchUpsertOne) DoNothing() *BranchUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the BranchCreate.OnConflict
// documentation for more info.
func (u *BranchUpsertOne) Update(set func(*BranchUpsert)) *BranchUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&BranchUpsert{UpdateSet: update})
	}))
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *BranchUpsertOne) SetCreatedAt(v time.Time) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateCreatedAt() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *BranchUpsertOne) SetUpdatedAt(v time.Time) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateUpdatedAt() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *BranchUpsertOne) SetDeletedAt(v time.Time) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateDeletedAt() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *BranchUpsertOne) ClearDeletedAt() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.ClearDeletedAt()
	})
}

// SetCreator sets the "creator" field.
func (u *BranchUpsertOne) SetCreator(v *model.Modifier) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetCreator(v)
	})
}

// UpdateCreator sets the "creator" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateCreator() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCreator()
	})
}

// ClearCreator clears the value of the "creator" field.
func (u *BranchUpsertOne) ClearCreator() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.ClearCreator()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *BranchUpsertOne) SetLastModifier(v *model.Modifier) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateLastModifier() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *BranchUpsertOne) ClearLastModifier() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *BranchUpsertOne) SetRemark(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateRemark() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *BranchUpsertOne) ClearRemark() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.ClearRemark()
	})
}

// SetCityID sets the "city_id" field.
func (u *BranchUpsertOne) SetCityID(v uint64) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetCityID(v)
	})
}

// AddCityID adds v to the "city_id" field.
func (u *BranchUpsertOne) AddCityID(v uint64) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.AddCityID(v)
	})
}

// UpdateCityID sets the "city_id" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateCityID() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCityID()
	})
}

// SetName sets the "name" field.
func (u *BranchUpsertOne) SetName(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateName() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateName()
	})
}

// SetLng sets the "lng" field.
func (u *BranchUpsertOne) SetLng(v float64) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetLng(v)
	})
}

// AddLng adds v to the "lng" field.
func (u *BranchUpsertOne) AddLng(v float64) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.AddLng(v)
	})
}

// UpdateLng sets the "lng" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateLng() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateLng()
	})
}

// SetLat sets the "lat" field.
func (u *BranchUpsertOne) SetLat(v float64) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetLat(v)
	})
}

// AddLat adds v to the "lat" field.
func (u *BranchUpsertOne) AddLat(v float64) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.AddLat(v)
	})
}

// UpdateLat sets the "lat" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateLat() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateLat()
	})
}

// SetAddress sets the "address" field.
func (u *BranchUpsertOne) SetAddress(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetAddress(v)
	})
}

// UpdateAddress sets the "address" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateAddress() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateAddress()
	})
}

// SetPhotos sets the "photos" field.
func (u *BranchUpsertOne) SetPhotos(v []string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetPhotos(v)
	})
}

// UpdatePhotos sets the "photos" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdatePhotos() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdatePhotos()
	})
}

// Exec executes the query.
func (u *BranchUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for BranchCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *BranchUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *BranchUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *BranchUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// BranchCreateBulk is the builder for creating many Branch entities in bulk.
type BranchCreateBulk struct {
	config
	builders []*BranchCreate
	conflict []sql.ConflictOption
}

// Save creates the Branch entities in the database.
func (bcb *BranchCreateBulk) Save(ctx context.Context) ([]*Branch, error) {
	specs := make([]*sqlgraph.CreateSpec, len(bcb.builders))
	nodes := make([]*Branch, len(bcb.builders))
	mutators := make([]Mutator, len(bcb.builders))
	for i := range bcb.builders {
		func(i int, root context.Context) {
			builder := bcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*BranchMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, bcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = bcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, bcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = uint64(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, bcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (bcb *BranchCreateBulk) SaveX(ctx context.Context) []*Branch {
	v, err := bcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bcb *BranchCreateBulk) Exec(ctx context.Context) error {
	_, err := bcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bcb *BranchCreateBulk) ExecX(ctx context.Context) {
	if err := bcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Branch.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.BranchUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
//
func (bcb *BranchCreateBulk) OnConflict(opts ...sql.ConflictOption) *BranchUpsertBulk {
	bcb.conflict = opts
	return &BranchUpsertBulk{
		create: bcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Branch.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (bcb *BranchCreateBulk) OnConflictColumns(columns ...string) *BranchUpsertBulk {
	bcb.conflict = append(bcb.conflict, sql.ConflictColumns(columns...))
	return &BranchUpsertBulk{
		create: bcb,
	}
}

// BranchUpsertBulk is the builder for "upsert"-ing
// a bulk of Branch nodes.
type BranchUpsertBulk struct {
	create *BranchCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Branch.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
//
func (u *BranchUpsertBulk) UpdateNewValues() *BranchUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(branch.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Branch.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
//
func (u *BranchUpsertBulk) Ignore() *BranchUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *BranchUpsertBulk) DoNothing() *BranchUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the BranchCreateBulk.OnConflict
// documentation for more info.
func (u *BranchUpsertBulk) Update(set func(*BranchUpsert)) *BranchUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&BranchUpsert{UpdateSet: update})
	}))
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *BranchUpsertBulk) SetCreatedAt(v time.Time) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateCreatedAt() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *BranchUpsertBulk) SetUpdatedAt(v time.Time) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateUpdatedAt() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *BranchUpsertBulk) SetDeletedAt(v time.Time) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateDeletedAt() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *BranchUpsertBulk) ClearDeletedAt() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.ClearDeletedAt()
	})
}

// SetCreator sets the "creator" field.
func (u *BranchUpsertBulk) SetCreator(v *model.Modifier) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetCreator(v)
	})
}

// UpdateCreator sets the "creator" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateCreator() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCreator()
	})
}

// ClearCreator clears the value of the "creator" field.
func (u *BranchUpsertBulk) ClearCreator() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.ClearCreator()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *BranchUpsertBulk) SetLastModifier(v *model.Modifier) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateLastModifier() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *BranchUpsertBulk) ClearLastModifier() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *BranchUpsertBulk) SetRemark(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateRemark() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *BranchUpsertBulk) ClearRemark() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.ClearRemark()
	})
}

// SetCityID sets the "city_id" field.
func (u *BranchUpsertBulk) SetCityID(v uint64) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetCityID(v)
	})
}

// AddCityID adds v to the "city_id" field.
func (u *BranchUpsertBulk) AddCityID(v uint64) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.AddCityID(v)
	})
}

// UpdateCityID sets the "city_id" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateCityID() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCityID()
	})
}

// SetName sets the "name" field.
func (u *BranchUpsertBulk) SetName(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateName() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateName()
	})
}

// SetLng sets the "lng" field.
func (u *BranchUpsertBulk) SetLng(v float64) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetLng(v)
	})
}

// AddLng adds v to the "lng" field.
func (u *BranchUpsertBulk) AddLng(v float64) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.AddLng(v)
	})
}

// UpdateLng sets the "lng" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateLng() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateLng()
	})
}

// SetLat sets the "lat" field.
func (u *BranchUpsertBulk) SetLat(v float64) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetLat(v)
	})
}

// AddLat adds v to the "lat" field.
func (u *BranchUpsertBulk) AddLat(v float64) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.AddLat(v)
	})
}

// UpdateLat sets the "lat" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateLat() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateLat()
	})
}

// SetAddress sets the "address" field.
func (u *BranchUpsertBulk) SetAddress(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetAddress(v)
	})
}

// UpdateAddress sets the "address" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateAddress() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateAddress()
	})
}

// SetPhotos sets the "photos" field.
func (u *BranchUpsertBulk) SetPhotos(v []string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetPhotos(v)
	})
}

// UpdatePhotos sets the "photos" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdatePhotos() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdatePhotos()
	})
}

// Exec executes the query.
func (u *BranchUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the BranchCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for BranchCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *BranchUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}