// Code generated by ent, DO NOT EDIT.

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
	"github.com/auroraride/aurservd/internal/ent/manager"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/riderfollowup"
)

// RiderFollowUpCreate is the builder for creating a RiderFollowUp entity.
type RiderFollowUpCreate struct {
	config
	mutation *RiderFollowUpMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (rfuc *RiderFollowUpCreate) SetCreatedAt(t time.Time) *RiderFollowUpCreate {
	rfuc.mutation.SetCreatedAt(t)
	return rfuc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (rfuc *RiderFollowUpCreate) SetNillableCreatedAt(t *time.Time) *RiderFollowUpCreate {
	if t != nil {
		rfuc.SetCreatedAt(*t)
	}
	return rfuc
}

// SetUpdatedAt sets the "updated_at" field.
func (rfuc *RiderFollowUpCreate) SetUpdatedAt(t time.Time) *RiderFollowUpCreate {
	rfuc.mutation.SetUpdatedAt(t)
	return rfuc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (rfuc *RiderFollowUpCreate) SetNillableUpdatedAt(t *time.Time) *RiderFollowUpCreate {
	if t != nil {
		rfuc.SetUpdatedAt(*t)
	}
	return rfuc
}

// SetDeletedAt sets the "deleted_at" field.
func (rfuc *RiderFollowUpCreate) SetDeletedAt(t time.Time) *RiderFollowUpCreate {
	rfuc.mutation.SetDeletedAt(t)
	return rfuc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (rfuc *RiderFollowUpCreate) SetNillableDeletedAt(t *time.Time) *RiderFollowUpCreate {
	if t != nil {
		rfuc.SetDeletedAt(*t)
	}
	return rfuc
}

// SetCreator sets the "creator" field.
func (rfuc *RiderFollowUpCreate) SetCreator(m *model.Modifier) *RiderFollowUpCreate {
	rfuc.mutation.SetCreator(m)
	return rfuc
}

// SetLastModifier sets the "last_modifier" field.
func (rfuc *RiderFollowUpCreate) SetLastModifier(m *model.Modifier) *RiderFollowUpCreate {
	rfuc.mutation.SetLastModifier(m)
	return rfuc
}

// SetRemark sets the "remark" field.
func (rfuc *RiderFollowUpCreate) SetRemark(s string) *RiderFollowUpCreate {
	rfuc.mutation.SetRemark(s)
	return rfuc
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (rfuc *RiderFollowUpCreate) SetNillableRemark(s *string) *RiderFollowUpCreate {
	if s != nil {
		rfuc.SetRemark(*s)
	}
	return rfuc
}

// SetManagerID sets the "manager_id" field.
func (rfuc *RiderFollowUpCreate) SetManagerID(u uint64) *RiderFollowUpCreate {
	rfuc.mutation.SetManagerID(u)
	return rfuc
}

// SetRiderID sets the "rider_id" field.
func (rfuc *RiderFollowUpCreate) SetRiderID(u uint64) *RiderFollowUpCreate {
	rfuc.mutation.SetRiderID(u)
	return rfuc
}

// SetManager sets the "manager" edge to the Manager entity.
func (rfuc *RiderFollowUpCreate) SetManager(m *Manager) *RiderFollowUpCreate {
	return rfuc.SetManagerID(m.ID)
}

// SetRider sets the "rider" edge to the Rider entity.
func (rfuc *RiderFollowUpCreate) SetRider(r *Rider) *RiderFollowUpCreate {
	return rfuc.SetRiderID(r.ID)
}

// Mutation returns the RiderFollowUpMutation object of the builder.
func (rfuc *RiderFollowUpCreate) Mutation() *RiderFollowUpMutation {
	return rfuc.mutation
}

// Save creates the RiderFollowUp in the database.
func (rfuc *RiderFollowUpCreate) Save(ctx context.Context) (*RiderFollowUp, error) {
	if err := rfuc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, rfuc.sqlSave, rfuc.mutation, rfuc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (rfuc *RiderFollowUpCreate) SaveX(ctx context.Context) *RiderFollowUp {
	v, err := rfuc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rfuc *RiderFollowUpCreate) Exec(ctx context.Context) error {
	_, err := rfuc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rfuc *RiderFollowUpCreate) ExecX(ctx context.Context) {
	if err := rfuc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rfuc *RiderFollowUpCreate) defaults() error {
	if _, ok := rfuc.mutation.CreatedAt(); !ok {
		if riderfollowup.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized riderfollowup.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := riderfollowup.DefaultCreatedAt()
		rfuc.mutation.SetCreatedAt(v)
	}
	if _, ok := rfuc.mutation.UpdatedAt(); !ok {
		if riderfollowup.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized riderfollowup.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := riderfollowup.DefaultUpdatedAt()
		rfuc.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (rfuc *RiderFollowUpCreate) check() error {
	if _, ok := rfuc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "RiderFollowUp.created_at"`)}
	}
	if _, ok := rfuc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "RiderFollowUp.updated_at"`)}
	}
	if _, ok := rfuc.mutation.ManagerID(); !ok {
		return &ValidationError{Name: "manager_id", err: errors.New(`ent: missing required field "RiderFollowUp.manager_id"`)}
	}
	if _, ok := rfuc.mutation.RiderID(); !ok {
		return &ValidationError{Name: "rider_id", err: errors.New(`ent: missing required field "RiderFollowUp.rider_id"`)}
	}
	if len(rfuc.mutation.ManagerIDs()) == 0 {
		return &ValidationError{Name: "manager", err: errors.New(`ent: missing required edge "RiderFollowUp.manager"`)}
	}
	if len(rfuc.mutation.RiderIDs()) == 0 {
		return &ValidationError{Name: "rider", err: errors.New(`ent: missing required edge "RiderFollowUp.rider"`)}
	}
	return nil
}

func (rfuc *RiderFollowUpCreate) sqlSave(ctx context.Context) (*RiderFollowUp, error) {
	if err := rfuc.check(); err != nil {
		return nil, err
	}
	_node, _spec := rfuc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rfuc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	rfuc.mutation.id = &_node.ID
	rfuc.mutation.done = true
	return _node, nil
}

func (rfuc *RiderFollowUpCreate) createSpec() (*RiderFollowUp, *sqlgraph.CreateSpec) {
	var (
		_node = &RiderFollowUp{config: rfuc.config}
		_spec = sqlgraph.NewCreateSpec(riderfollowup.Table, sqlgraph.NewFieldSpec(riderfollowup.FieldID, field.TypeUint64))
	)
	_spec.OnConflict = rfuc.conflict
	if value, ok := rfuc.mutation.CreatedAt(); ok {
		_spec.SetField(riderfollowup.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := rfuc.mutation.UpdatedAt(); ok {
		_spec.SetField(riderfollowup.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := rfuc.mutation.DeletedAt(); ok {
		_spec.SetField(riderfollowup.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	if value, ok := rfuc.mutation.Creator(); ok {
		_spec.SetField(riderfollowup.FieldCreator, field.TypeJSON, value)
		_node.Creator = value
	}
	if value, ok := rfuc.mutation.LastModifier(); ok {
		_spec.SetField(riderfollowup.FieldLastModifier, field.TypeJSON, value)
		_node.LastModifier = value
	}
	if value, ok := rfuc.mutation.Remark(); ok {
		_spec.SetField(riderfollowup.FieldRemark, field.TypeString, value)
		_node.Remark = value
	}
	if nodes := rfuc.mutation.ManagerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   riderfollowup.ManagerTable,
			Columns: []string{riderfollowup.ManagerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manager.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ManagerID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rfuc.mutation.RiderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   riderfollowup.RiderTable,
			Columns: []string{riderfollowup.RiderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(rider.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.RiderID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.RiderFollowUp.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.RiderFollowUpUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (rfuc *RiderFollowUpCreate) OnConflict(opts ...sql.ConflictOption) *RiderFollowUpUpsertOne {
	rfuc.conflict = opts
	return &RiderFollowUpUpsertOne{
		create: rfuc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.RiderFollowUp.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (rfuc *RiderFollowUpCreate) OnConflictColumns(columns ...string) *RiderFollowUpUpsertOne {
	rfuc.conflict = append(rfuc.conflict, sql.ConflictColumns(columns...))
	return &RiderFollowUpUpsertOne{
		create: rfuc,
	}
}

type (
	// RiderFollowUpUpsertOne is the builder for "upsert"-ing
	//  one RiderFollowUp node.
	RiderFollowUpUpsertOne struct {
		create *RiderFollowUpCreate
	}

	// RiderFollowUpUpsert is the "OnConflict" setter.
	RiderFollowUpUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *RiderFollowUpUpsert) SetUpdatedAt(v time.Time) *RiderFollowUpUpsert {
	u.Set(riderfollowup.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *RiderFollowUpUpsert) UpdateUpdatedAt() *RiderFollowUpUpsert {
	u.SetExcluded(riderfollowup.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *RiderFollowUpUpsert) SetDeletedAt(v time.Time) *RiderFollowUpUpsert {
	u.Set(riderfollowup.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *RiderFollowUpUpsert) UpdateDeletedAt() *RiderFollowUpUpsert {
	u.SetExcluded(riderfollowup.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *RiderFollowUpUpsert) ClearDeletedAt() *RiderFollowUpUpsert {
	u.SetNull(riderfollowup.FieldDeletedAt)
	return u
}

// SetLastModifier sets the "last_modifier" field.
func (u *RiderFollowUpUpsert) SetLastModifier(v *model.Modifier) *RiderFollowUpUpsert {
	u.Set(riderfollowup.FieldLastModifier, v)
	return u
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *RiderFollowUpUpsert) UpdateLastModifier() *RiderFollowUpUpsert {
	u.SetExcluded(riderfollowup.FieldLastModifier)
	return u
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *RiderFollowUpUpsert) ClearLastModifier() *RiderFollowUpUpsert {
	u.SetNull(riderfollowup.FieldLastModifier)
	return u
}

// SetRemark sets the "remark" field.
func (u *RiderFollowUpUpsert) SetRemark(v string) *RiderFollowUpUpsert {
	u.Set(riderfollowup.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *RiderFollowUpUpsert) UpdateRemark() *RiderFollowUpUpsert {
	u.SetExcluded(riderfollowup.FieldRemark)
	return u
}

// ClearRemark clears the value of the "remark" field.
func (u *RiderFollowUpUpsert) ClearRemark() *RiderFollowUpUpsert {
	u.SetNull(riderfollowup.FieldRemark)
	return u
}

// SetManagerID sets the "manager_id" field.
func (u *RiderFollowUpUpsert) SetManagerID(v uint64) *RiderFollowUpUpsert {
	u.Set(riderfollowup.FieldManagerID, v)
	return u
}

// UpdateManagerID sets the "manager_id" field to the value that was provided on create.
func (u *RiderFollowUpUpsert) UpdateManagerID() *RiderFollowUpUpsert {
	u.SetExcluded(riderfollowup.FieldManagerID)
	return u
}

// SetRiderID sets the "rider_id" field.
func (u *RiderFollowUpUpsert) SetRiderID(v uint64) *RiderFollowUpUpsert {
	u.Set(riderfollowup.FieldRiderID, v)
	return u
}

// UpdateRiderID sets the "rider_id" field to the value that was provided on create.
func (u *RiderFollowUpUpsert) UpdateRiderID() *RiderFollowUpUpsert {
	u.SetExcluded(riderfollowup.FieldRiderID)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.RiderFollowUp.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *RiderFollowUpUpsertOne) UpdateNewValues() *RiderFollowUpUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(riderfollowup.FieldCreatedAt)
		}
		if _, exists := u.create.mutation.Creator(); exists {
			s.SetIgnore(riderfollowup.FieldCreator)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.RiderFollowUp.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *RiderFollowUpUpsertOne) Ignore() *RiderFollowUpUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *RiderFollowUpUpsertOne) DoNothing() *RiderFollowUpUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the RiderFollowUpCreate.OnConflict
// documentation for more info.
func (u *RiderFollowUpUpsertOne) Update(set func(*RiderFollowUpUpsert)) *RiderFollowUpUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&RiderFollowUpUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *RiderFollowUpUpsertOne) SetUpdatedAt(v time.Time) *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *RiderFollowUpUpsertOne) UpdateUpdatedAt() *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *RiderFollowUpUpsertOne) SetDeletedAt(v time.Time) *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *RiderFollowUpUpsertOne) UpdateDeletedAt() *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *RiderFollowUpUpsertOne) ClearDeletedAt() *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *RiderFollowUpUpsertOne) SetLastModifier(v *model.Modifier) *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *RiderFollowUpUpsertOne) UpdateLastModifier() *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *RiderFollowUpUpsertOne) ClearLastModifier() *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *RiderFollowUpUpsertOne) SetRemark(v string) *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *RiderFollowUpUpsertOne) UpdateRemark() *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *RiderFollowUpUpsertOne) ClearRemark() *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.ClearRemark()
	})
}

// SetManagerID sets the "manager_id" field.
func (u *RiderFollowUpUpsertOne) SetManagerID(v uint64) *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetManagerID(v)
	})
}

// UpdateManagerID sets the "manager_id" field to the value that was provided on create.
func (u *RiderFollowUpUpsertOne) UpdateManagerID() *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateManagerID()
	})
}

// SetRiderID sets the "rider_id" field.
func (u *RiderFollowUpUpsertOne) SetRiderID(v uint64) *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetRiderID(v)
	})
}

// UpdateRiderID sets the "rider_id" field to the value that was provided on create.
func (u *RiderFollowUpUpsertOne) UpdateRiderID() *RiderFollowUpUpsertOne {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateRiderID()
	})
}

// Exec executes the query.
func (u *RiderFollowUpUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for RiderFollowUpCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *RiderFollowUpUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *RiderFollowUpUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *RiderFollowUpUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// RiderFollowUpCreateBulk is the builder for creating many RiderFollowUp entities in bulk.
type RiderFollowUpCreateBulk struct {
	config
	err      error
	builders []*RiderFollowUpCreate
	conflict []sql.ConflictOption
}

// Save creates the RiderFollowUp entities in the database.
func (rfucb *RiderFollowUpCreateBulk) Save(ctx context.Context) ([]*RiderFollowUp, error) {
	if rfucb.err != nil {
		return nil, rfucb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(rfucb.builders))
	nodes := make([]*RiderFollowUp, len(rfucb.builders))
	mutators := make([]Mutator, len(rfucb.builders))
	for i := range rfucb.builders {
		func(i int, root context.Context) {
			builder := rfucb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RiderFollowUpMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, rfucb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = rfucb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rfucb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = uint64(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, rfucb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rfucb *RiderFollowUpCreateBulk) SaveX(ctx context.Context) []*RiderFollowUp {
	v, err := rfucb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rfucb *RiderFollowUpCreateBulk) Exec(ctx context.Context) error {
	_, err := rfucb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rfucb *RiderFollowUpCreateBulk) ExecX(ctx context.Context) {
	if err := rfucb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.RiderFollowUp.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.RiderFollowUpUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (rfucb *RiderFollowUpCreateBulk) OnConflict(opts ...sql.ConflictOption) *RiderFollowUpUpsertBulk {
	rfucb.conflict = opts
	return &RiderFollowUpUpsertBulk{
		create: rfucb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.RiderFollowUp.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (rfucb *RiderFollowUpCreateBulk) OnConflictColumns(columns ...string) *RiderFollowUpUpsertBulk {
	rfucb.conflict = append(rfucb.conflict, sql.ConflictColumns(columns...))
	return &RiderFollowUpUpsertBulk{
		create: rfucb,
	}
}

// RiderFollowUpUpsertBulk is the builder for "upsert"-ing
// a bulk of RiderFollowUp nodes.
type RiderFollowUpUpsertBulk struct {
	create *RiderFollowUpCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.RiderFollowUp.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *RiderFollowUpUpsertBulk) UpdateNewValues() *RiderFollowUpUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(riderfollowup.FieldCreatedAt)
			}
			if _, exists := b.mutation.Creator(); exists {
				s.SetIgnore(riderfollowup.FieldCreator)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.RiderFollowUp.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *RiderFollowUpUpsertBulk) Ignore() *RiderFollowUpUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *RiderFollowUpUpsertBulk) DoNothing() *RiderFollowUpUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the RiderFollowUpCreateBulk.OnConflict
// documentation for more info.
func (u *RiderFollowUpUpsertBulk) Update(set func(*RiderFollowUpUpsert)) *RiderFollowUpUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&RiderFollowUpUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *RiderFollowUpUpsertBulk) SetUpdatedAt(v time.Time) *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *RiderFollowUpUpsertBulk) UpdateUpdatedAt() *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *RiderFollowUpUpsertBulk) SetDeletedAt(v time.Time) *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *RiderFollowUpUpsertBulk) UpdateDeletedAt() *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *RiderFollowUpUpsertBulk) ClearDeletedAt() *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *RiderFollowUpUpsertBulk) SetLastModifier(v *model.Modifier) *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *RiderFollowUpUpsertBulk) UpdateLastModifier() *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *RiderFollowUpUpsertBulk) ClearLastModifier() *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *RiderFollowUpUpsertBulk) SetRemark(v string) *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *RiderFollowUpUpsertBulk) UpdateRemark() *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *RiderFollowUpUpsertBulk) ClearRemark() *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.ClearRemark()
	})
}

// SetManagerID sets the "manager_id" field.
func (u *RiderFollowUpUpsertBulk) SetManagerID(v uint64) *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetManagerID(v)
	})
}

// UpdateManagerID sets the "manager_id" field to the value that was provided on create.
func (u *RiderFollowUpUpsertBulk) UpdateManagerID() *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateManagerID()
	})
}

// SetRiderID sets the "rider_id" field.
func (u *RiderFollowUpUpsertBulk) SetRiderID(v uint64) *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.SetRiderID(v)
	})
}

// UpdateRiderID sets the "rider_id" field to the value that was provided on create.
func (u *RiderFollowUpUpsertBulk) UpdateRiderID() *RiderFollowUpUpsertBulk {
	return u.Update(func(s *RiderFollowUpUpsert) {
		s.UpdateRiderID()
	})
}

// Exec executes the query.
func (u *RiderFollowUpUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the RiderFollowUpCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for RiderFollowUpCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *RiderFollowUpUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
