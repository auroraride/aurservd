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
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/rider"
)

// EnterpriseCreate is the builder for creating a Enterprise entity.
type EnterpriseCreate struct {
	config
	mutation *EnterpriseMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (ec *EnterpriseCreate) SetCreatedAt(t time.Time) *EnterpriseCreate {
	ec.mutation.SetCreatedAt(t)
	return ec
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (ec *EnterpriseCreate) SetNillableCreatedAt(t *time.Time) *EnterpriseCreate {
	if t != nil {
		ec.SetCreatedAt(*t)
	}
	return ec
}

// SetUpdatedAt sets the "updated_at" field.
func (ec *EnterpriseCreate) SetUpdatedAt(t time.Time) *EnterpriseCreate {
	ec.mutation.SetUpdatedAt(t)
	return ec
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (ec *EnterpriseCreate) SetNillableUpdatedAt(t *time.Time) *EnterpriseCreate {
	if t != nil {
		ec.SetUpdatedAt(*t)
	}
	return ec
}

// SetDeletedAt sets the "deleted_at" field.
func (ec *EnterpriseCreate) SetDeletedAt(t time.Time) *EnterpriseCreate {
	ec.mutation.SetDeletedAt(t)
	return ec
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (ec *EnterpriseCreate) SetNillableDeletedAt(t *time.Time) *EnterpriseCreate {
	if t != nil {
		ec.SetDeletedAt(*t)
	}
	return ec
}

// SetCreator sets the "creator" field.
func (ec *EnterpriseCreate) SetCreator(m *model.Modifier) *EnterpriseCreate {
	ec.mutation.SetCreator(m)
	return ec
}

// SetLastModifier sets the "last_modifier" field.
func (ec *EnterpriseCreate) SetLastModifier(m *model.Modifier) *EnterpriseCreate {
	ec.mutation.SetLastModifier(m)
	return ec
}

// SetRemark sets the "remark" field.
func (ec *EnterpriseCreate) SetRemark(s string) *EnterpriseCreate {
	ec.mutation.SetRemark(s)
	return ec
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (ec *EnterpriseCreate) SetNillableRemark(s *string) *EnterpriseCreate {
	if s != nil {
		ec.SetRemark(*s)
	}
	return ec
}

// SetName sets the "name" field.
func (ec *EnterpriseCreate) SetName(s string) *EnterpriseCreate {
	ec.mutation.SetName(s)
	return ec
}

// AddRiderIDs adds the "riders" edge to the Rider entity by IDs.
func (ec *EnterpriseCreate) AddRiderIDs(ids ...uint64) *EnterpriseCreate {
	ec.mutation.AddRiderIDs(ids...)
	return ec
}

// AddRiders adds the "riders" edges to the Rider entity.
func (ec *EnterpriseCreate) AddRiders(r ...*Rider) *EnterpriseCreate {
	ids := make([]uint64, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ec.AddRiderIDs(ids...)
}

// Mutation returns the EnterpriseMutation object of the builder.
func (ec *EnterpriseCreate) Mutation() *EnterpriseMutation {
	return ec.mutation
}

// Save creates the Enterprise in the database.
func (ec *EnterpriseCreate) Save(ctx context.Context) (*Enterprise, error) {
	var (
		err  error
		node *Enterprise
	)
	ec.defaults()
	if len(ec.hooks) == 0 {
		if err = ec.check(); err != nil {
			return nil, err
		}
		node, err = ec.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EnterpriseMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ec.check(); err != nil {
				return nil, err
			}
			ec.mutation = mutation
			if node, err = ec.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(ec.hooks) - 1; i >= 0; i-- {
			if ec.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ec.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, ec.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Enterprise)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from EnterpriseMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ec *EnterpriseCreate) SaveX(ctx context.Context) *Enterprise {
	v, err := ec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ec *EnterpriseCreate) Exec(ctx context.Context) error {
	_, err := ec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ec *EnterpriseCreate) ExecX(ctx context.Context) {
	if err := ec.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ec *EnterpriseCreate) defaults() {
	if _, ok := ec.mutation.CreatedAt(); !ok {
		v := enterprise.DefaultCreatedAt()
		ec.mutation.SetCreatedAt(v)
	}
	if _, ok := ec.mutation.UpdatedAt(); !ok {
		v := enterprise.DefaultUpdatedAt()
		ec.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ec *EnterpriseCreate) check() error {
	if _, ok := ec.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Enterprise.created_at"`)}
	}
	if _, ok := ec.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Enterprise.updated_at"`)}
	}
	if _, ok := ec.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Enterprise.name"`)}
	}
	return nil
}

func (ec *EnterpriseCreate) sqlSave(ctx context.Context) (*Enterprise, error) {
	_node, _spec := ec.createSpec()
	if err := sqlgraph.CreateNode(ctx, ec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	return _node, nil
}

func (ec *EnterpriseCreate) createSpec() (*Enterprise, *sqlgraph.CreateSpec) {
	var (
		_node = &Enterprise{config: ec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: enterprise.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: enterprise.FieldID,
			},
		}
	)
	_spec.OnConflict = ec.conflict
	if value, ok := ec.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: enterprise.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := ec.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: enterprise.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := ec.mutation.DeletedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: enterprise.FieldDeletedAt,
		})
		_node.DeletedAt = &value
	}
	if value, ok := ec.mutation.Creator(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: enterprise.FieldCreator,
		})
		_node.Creator = value
	}
	if value, ok := ec.mutation.LastModifier(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: enterprise.FieldLastModifier,
		})
		_node.LastModifier = value
	}
	if value, ok := ec.mutation.Remark(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: enterprise.FieldRemark,
		})
		_node.Remark = value
	}
	if value, ok := ec.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: enterprise.FieldName,
		})
		_node.Name = value
	}
	if nodes := ec.mutation.RidersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   enterprise.RidersTable,
			Columns: []string{enterprise.RidersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: rider.FieldID,
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
//	client.Enterprise.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnterpriseUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
//
func (ec *EnterpriseCreate) OnConflict(opts ...sql.ConflictOption) *EnterpriseUpsertOne {
	ec.conflict = opts
	return &EnterpriseUpsertOne{
		create: ec,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Enterprise.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (ec *EnterpriseCreate) OnConflictColumns(columns ...string) *EnterpriseUpsertOne {
	ec.conflict = append(ec.conflict, sql.ConflictColumns(columns...))
	return &EnterpriseUpsertOne{
		create: ec,
	}
}

type (
	// EnterpriseUpsertOne is the builder for "upsert"-ing
	//  one Enterprise node.
	EnterpriseUpsertOne struct {
		create *EnterpriseCreate
	}

	// EnterpriseUpsert is the "OnConflict" setter.
	EnterpriseUpsert struct {
		*sql.UpdateSet
	}
)

// SetCreatedAt sets the "created_at" field.
func (u *EnterpriseUpsert) SetCreatedAt(v time.Time) *EnterpriseUpsert {
	u.Set(enterprise.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *EnterpriseUpsert) UpdateCreatedAt() *EnterpriseUpsert {
	u.SetExcluded(enterprise.FieldCreatedAt)
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *EnterpriseUpsert) SetUpdatedAt(v time.Time) *EnterpriseUpsert {
	u.Set(enterprise.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *EnterpriseUpsert) UpdateUpdatedAt() *EnterpriseUpsert {
	u.SetExcluded(enterprise.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *EnterpriseUpsert) SetDeletedAt(v time.Time) *EnterpriseUpsert {
	u.Set(enterprise.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *EnterpriseUpsert) UpdateDeletedAt() *EnterpriseUpsert {
	u.SetExcluded(enterprise.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *EnterpriseUpsert) ClearDeletedAt() *EnterpriseUpsert {
	u.SetNull(enterprise.FieldDeletedAt)
	return u
}

// SetCreator sets the "creator" field.
func (u *EnterpriseUpsert) SetCreator(v *model.Modifier) *EnterpriseUpsert {
	u.Set(enterprise.FieldCreator, v)
	return u
}

// UpdateCreator sets the "creator" field to the value that was provided on create.
func (u *EnterpriseUpsert) UpdateCreator() *EnterpriseUpsert {
	u.SetExcluded(enterprise.FieldCreator)
	return u
}

// ClearCreator clears the value of the "creator" field.
func (u *EnterpriseUpsert) ClearCreator() *EnterpriseUpsert {
	u.SetNull(enterprise.FieldCreator)
	return u
}

// SetLastModifier sets the "last_modifier" field.
func (u *EnterpriseUpsert) SetLastModifier(v *model.Modifier) *EnterpriseUpsert {
	u.Set(enterprise.FieldLastModifier, v)
	return u
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *EnterpriseUpsert) UpdateLastModifier() *EnterpriseUpsert {
	u.SetExcluded(enterprise.FieldLastModifier)
	return u
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *EnterpriseUpsert) ClearLastModifier() *EnterpriseUpsert {
	u.SetNull(enterprise.FieldLastModifier)
	return u
}

// SetRemark sets the "remark" field.
func (u *EnterpriseUpsert) SetRemark(v string) *EnterpriseUpsert {
	u.Set(enterprise.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *EnterpriseUpsert) UpdateRemark() *EnterpriseUpsert {
	u.SetExcluded(enterprise.FieldRemark)
	return u
}

// ClearRemark clears the value of the "remark" field.
func (u *EnterpriseUpsert) ClearRemark() *EnterpriseUpsert {
	u.SetNull(enterprise.FieldRemark)
	return u
}

// SetName sets the "name" field.
func (u *EnterpriseUpsert) SetName(v string) *EnterpriseUpsert {
	u.Set(enterprise.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *EnterpriseUpsert) UpdateName() *EnterpriseUpsert {
	u.SetExcluded(enterprise.FieldName)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Enterprise.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
//
func (u *EnterpriseUpsertOne) UpdateNewValues() *EnterpriseUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(enterprise.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.Enterprise.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *EnterpriseUpsertOne) Ignore() *EnterpriseUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EnterpriseUpsertOne) DoNothing() *EnterpriseUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EnterpriseCreate.OnConflict
// documentation for more info.
func (u *EnterpriseUpsertOne) Update(set func(*EnterpriseUpsert)) *EnterpriseUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EnterpriseUpsert{UpdateSet: update})
	}))
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *EnterpriseUpsertOne) SetCreatedAt(v time.Time) *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *EnterpriseUpsertOne) UpdateCreatedAt() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *EnterpriseUpsertOne) SetUpdatedAt(v time.Time) *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *EnterpriseUpsertOne) UpdateUpdatedAt() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *EnterpriseUpsertOne) SetDeletedAt(v time.Time) *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *EnterpriseUpsertOne) UpdateDeletedAt() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *EnterpriseUpsertOne) ClearDeletedAt() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.ClearDeletedAt()
	})
}

// SetCreator sets the "creator" field.
func (u *EnterpriseUpsertOne) SetCreator(v *model.Modifier) *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetCreator(v)
	})
}

// UpdateCreator sets the "creator" field to the value that was provided on create.
func (u *EnterpriseUpsertOne) UpdateCreator() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateCreator()
	})
}

// ClearCreator clears the value of the "creator" field.
func (u *EnterpriseUpsertOne) ClearCreator() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.ClearCreator()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *EnterpriseUpsertOne) SetLastModifier(v *model.Modifier) *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *EnterpriseUpsertOne) UpdateLastModifier() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *EnterpriseUpsertOne) ClearLastModifier() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *EnterpriseUpsertOne) SetRemark(v string) *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *EnterpriseUpsertOne) UpdateRemark() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *EnterpriseUpsertOne) ClearRemark() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.ClearRemark()
	})
}

// SetName sets the "name" field.
func (u *EnterpriseUpsertOne) SetName(v string) *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *EnterpriseUpsertOne) UpdateName() *EnterpriseUpsertOne {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateName()
	})
}

// Exec executes the query.
func (u *EnterpriseUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for EnterpriseCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EnterpriseUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *EnterpriseUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *EnterpriseUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// EnterpriseCreateBulk is the builder for creating many Enterprise entities in bulk.
type EnterpriseCreateBulk struct {
	config
	builders []*EnterpriseCreate
	conflict []sql.ConflictOption
}

// Save creates the Enterprise entities in the database.
func (ecb *EnterpriseCreateBulk) Save(ctx context.Context) ([]*Enterprise, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ecb.builders))
	nodes := make([]*Enterprise, len(ecb.builders))
	mutators := make([]Mutator, len(ecb.builders))
	for i := range ecb.builders {
		func(i int, root context.Context) {
			builder := ecb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EnterpriseMutation)
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
					_, err = mutators[i+1].Mutate(root, ecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ecb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ecb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ecb *EnterpriseCreateBulk) SaveX(ctx context.Context) []*Enterprise {
	v, err := ecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ecb *EnterpriseCreateBulk) Exec(ctx context.Context) error {
	_, err := ecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecb *EnterpriseCreateBulk) ExecX(ctx context.Context) {
	if err := ecb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Enterprise.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnterpriseUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
//
func (ecb *EnterpriseCreateBulk) OnConflict(opts ...sql.ConflictOption) *EnterpriseUpsertBulk {
	ecb.conflict = opts
	return &EnterpriseUpsertBulk{
		create: ecb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Enterprise.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (ecb *EnterpriseCreateBulk) OnConflictColumns(columns ...string) *EnterpriseUpsertBulk {
	ecb.conflict = append(ecb.conflict, sql.ConflictColumns(columns...))
	return &EnterpriseUpsertBulk{
		create: ecb,
	}
}

// EnterpriseUpsertBulk is the builder for "upsert"-ing
// a bulk of Enterprise nodes.
type EnterpriseUpsertBulk struct {
	create *EnterpriseCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Enterprise.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
//
func (u *EnterpriseUpsertBulk) UpdateNewValues() *EnterpriseUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(enterprise.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Enterprise.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
//
func (u *EnterpriseUpsertBulk) Ignore() *EnterpriseUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EnterpriseUpsertBulk) DoNothing() *EnterpriseUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EnterpriseCreateBulk.OnConflict
// documentation for more info.
func (u *EnterpriseUpsertBulk) Update(set func(*EnterpriseUpsert)) *EnterpriseUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EnterpriseUpsert{UpdateSet: update})
	}))
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *EnterpriseUpsertBulk) SetCreatedAt(v time.Time) *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *EnterpriseUpsertBulk) UpdateCreatedAt() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *EnterpriseUpsertBulk) SetUpdatedAt(v time.Time) *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *EnterpriseUpsertBulk) UpdateUpdatedAt() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *EnterpriseUpsertBulk) SetDeletedAt(v time.Time) *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *EnterpriseUpsertBulk) UpdateDeletedAt() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *EnterpriseUpsertBulk) ClearDeletedAt() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.ClearDeletedAt()
	})
}

// SetCreator sets the "creator" field.
func (u *EnterpriseUpsertBulk) SetCreator(v *model.Modifier) *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetCreator(v)
	})
}

// UpdateCreator sets the "creator" field to the value that was provided on create.
func (u *EnterpriseUpsertBulk) UpdateCreator() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateCreator()
	})
}

// ClearCreator clears the value of the "creator" field.
func (u *EnterpriseUpsertBulk) ClearCreator() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.ClearCreator()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *EnterpriseUpsertBulk) SetLastModifier(v *model.Modifier) *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *EnterpriseUpsertBulk) UpdateLastModifier() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *EnterpriseUpsertBulk) ClearLastModifier() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *EnterpriseUpsertBulk) SetRemark(v string) *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *EnterpriseUpsertBulk) UpdateRemark() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *EnterpriseUpsertBulk) ClearRemark() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.ClearRemark()
	})
}

// SetName sets the "name" field.
func (u *EnterpriseUpsertBulk) SetName(v string) *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *EnterpriseUpsertBulk) UpdateName() *EnterpriseUpsertBulk {
	return u.Update(func(s *EnterpriseUpsert) {
		s.UpdateName()
	})
}

// Exec executes the query.
func (u *EnterpriseUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the EnterpriseCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for EnterpriseCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EnterpriseUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}