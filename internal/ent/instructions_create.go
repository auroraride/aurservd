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
	"github.com/auroraride/aurservd/internal/ent/instructions"
)

// InstructionsCreate is the builder for creating a Instructions entity.
type InstructionsCreate struct {
	config
	mutation *InstructionsMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (ic *InstructionsCreate) SetCreatedAt(t time.Time) *InstructionsCreate {
	ic.mutation.SetCreatedAt(t)
	return ic
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (ic *InstructionsCreate) SetNillableCreatedAt(t *time.Time) *InstructionsCreate {
	if t != nil {
		ic.SetCreatedAt(*t)
	}
	return ic
}

// SetUpdatedAt sets the "updated_at" field.
func (ic *InstructionsCreate) SetUpdatedAt(t time.Time) *InstructionsCreate {
	ic.mutation.SetUpdatedAt(t)
	return ic
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (ic *InstructionsCreate) SetNillableUpdatedAt(t *time.Time) *InstructionsCreate {
	if t != nil {
		ic.SetUpdatedAt(*t)
	}
	return ic
}

// SetDeletedAt sets the "deleted_at" field.
func (ic *InstructionsCreate) SetDeletedAt(t time.Time) *InstructionsCreate {
	ic.mutation.SetDeletedAt(t)
	return ic
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (ic *InstructionsCreate) SetNillableDeletedAt(t *time.Time) *InstructionsCreate {
	if t != nil {
		ic.SetDeletedAt(*t)
	}
	return ic
}

// SetCreator sets the "creator" field.
func (ic *InstructionsCreate) SetCreator(m *model.Modifier) *InstructionsCreate {
	ic.mutation.SetCreator(m)
	return ic
}

// SetLastModifier sets the "last_modifier" field.
func (ic *InstructionsCreate) SetLastModifier(m *model.Modifier) *InstructionsCreate {
	ic.mutation.SetLastModifier(m)
	return ic
}

// SetRemark sets the "remark" field.
func (ic *InstructionsCreate) SetRemark(s string) *InstructionsCreate {
	ic.mutation.SetRemark(s)
	return ic
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (ic *InstructionsCreate) SetNillableRemark(s *string) *InstructionsCreate {
	if s != nil {
		ic.SetRemark(*s)
	}
	return ic
}

// SetTitle sets the "title" field.
func (ic *InstructionsCreate) SetTitle(s string) *InstructionsCreate {
	ic.mutation.SetTitle(s)
	return ic
}

// SetContent sets the "content" field.
func (ic *InstructionsCreate) SetContent(i *interface{}) *InstructionsCreate {
	ic.mutation.SetContent(i)
	return ic
}

// SetKey sets the "key" field.
func (ic *InstructionsCreate) SetKey(s string) *InstructionsCreate {
	ic.mutation.SetKey(s)
	return ic
}

// Mutation returns the InstructionsMutation object of the builder.
func (ic *InstructionsCreate) Mutation() *InstructionsMutation {
	return ic.mutation
}

// Save creates the Instructions in the database.
func (ic *InstructionsCreate) Save(ctx context.Context) (*Instructions, error) {
	if err := ic.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, ic.sqlSave, ic.mutation, ic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ic *InstructionsCreate) SaveX(ctx context.Context) *Instructions {
	v, err := ic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ic *InstructionsCreate) Exec(ctx context.Context) error {
	_, err := ic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ic *InstructionsCreate) ExecX(ctx context.Context) {
	if err := ic.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ic *InstructionsCreate) defaults() error {
	if _, ok := ic.mutation.CreatedAt(); !ok {
		if instructions.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized instructions.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := instructions.DefaultCreatedAt()
		ic.mutation.SetCreatedAt(v)
	}
	if _, ok := ic.mutation.UpdatedAt(); !ok {
		if instructions.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized instructions.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := instructions.DefaultUpdatedAt()
		ic.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (ic *InstructionsCreate) check() error {
	if _, ok := ic.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Instructions.created_at"`)}
	}
	if _, ok := ic.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Instructions.updated_at"`)}
	}
	if _, ok := ic.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New(`ent: missing required field "Instructions.title"`)}
	}
	if _, ok := ic.mutation.Content(); !ok {
		return &ValidationError{Name: "content", err: errors.New(`ent: missing required field "Instructions.content"`)}
	}
	if _, ok := ic.mutation.Key(); !ok {
		return &ValidationError{Name: "key", err: errors.New(`ent: missing required field "Instructions.key"`)}
	}
	return nil
}

func (ic *InstructionsCreate) sqlSave(ctx context.Context) (*Instructions, error) {
	if err := ic.check(); err != nil {
		return nil, err
	}
	_node, _spec := ic.createSpec()
	if err := sqlgraph.CreateNode(ctx, ic.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	ic.mutation.id = &_node.ID
	ic.mutation.done = true
	return _node, nil
}

func (ic *InstructionsCreate) createSpec() (*Instructions, *sqlgraph.CreateSpec) {
	var (
		_node = &Instructions{config: ic.config}
		_spec = sqlgraph.NewCreateSpec(instructions.Table, sqlgraph.NewFieldSpec(instructions.FieldID, field.TypeUint64))
	)
	_spec.OnConflict = ic.conflict
	if value, ok := ic.mutation.CreatedAt(); ok {
		_spec.SetField(instructions.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := ic.mutation.UpdatedAt(); ok {
		_spec.SetField(instructions.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := ic.mutation.DeletedAt(); ok {
		_spec.SetField(instructions.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	if value, ok := ic.mutation.Creator(); ok {
		_spec.SetField(instructions.FieldCreator, field.TypeJSON, value)
		_node.Creator = value
	}
	if value, ok := ic.mutation.LastModifier(); ok {
		_spec.SetField(instructions.FieldLastModifier, field.TypeJSON, value)
		_node.LastModifier = value
	}
	if value, ok := ic.mutation.Remark(); ok {
		_spec.SetField(instructions.FieldRemark, field.TypeString, value)
		_node.Remark = value
	}
	if value, ok := ic.mutation.Title(); ok {
		_spec.SetField(instructions.FieldTitle, field.TypeString, value)
		_node.Title = value
	}
	if value, ok := ic.mutation.Content(); ok {
		_spec.SetField(instructions.FieldContent, field.TypeJSON, value)
		_node.Content = value
	}
	if value, ok := ic.mutation.Key(); ok {
		_spec.SetField(instructions.FieldKey, field.TypeString, value)
		_node.Key = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Instructions.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.InstructionsUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (ic *InstructionsCreate) OnConflict(opts ...sql.ConflictOption) *InstructionsUpsertOne {
	ic.conflict = opts
	return &InstructionsUpsertOne{
		create: ic,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Instructions.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ic *InstructionsCreate) OnConflictColumns(columns ...string) *InstructionsUpsertOne {
	ic.conflict = append(ic.conflict, sql.ConflictColumns(columns...))
	return &InstructionsUpsertOne{
		create: ic,
	}
}

type (
	// InstructionsUpsertOne is the builder for "upsert"-ing
	//  one Instructions node.
	InstructionsUpsertOne struct {
		create *InstructionsCreate
	}

	// InstructionsUpsert is the "OnConflict" setter.
	InstructionsUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *InstructionsUpsert) SetUpdatedAt(v time.Time) *InstructionsUpsert {
	u.Set(instructions.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *InstructionsUpsert) UpdateUpdatedAt() *InstructionsUpsert {
	u.SetExcluded(instructions.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *InstructionsUpsert) SetDeletedAt(v time.Time) *InstructionsUpsert {
	u.Set(instructions.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *InstructionsUpsert) UpdateDeletedAt() *InstructionsUpsert {
	u.SetExcluded(instructions.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *InstructionsUpsert) ClearDeletedAt() *InstructionsUpsert {
	u.SetNull(instructions.FieldDeletedAt)
	return u
}

// SetLastModifier sets the "last_modifier" field.
func (u *InstructionsUpsert) SetLastModifier(v *model.Modifier) *InstructionsUpsert {
	u.Set(instructions.FieldLastModifier, v)
	return u
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *InstructionsUpsert) UpdateLastModifier() *InstructionsUpsert {
	u.SetExcluded(instructions.FieldLastModifier)
	return u
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *InstructionsUpsert) ClearLastModifier() *InstructionsUpsert {
	u.SetNull(instructions.FieldLastModifier)
	return u
}

// SetRemark sets the "remark" field.
func (u *InstructionsUpsert) SetRemark(v string) *InstructionsUpsert {
	u.Set(instructions.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *InstructionsUpsert) UpdateRemark() *InstructionsUpsert {
	u.SetExcluded(instructions.FieldRemark)
	return u
}

// ClearRemark clears the value of the "remark" field.
func (u *InstructionsUpsert) ClearRemark() *InstructionsUpsert {
	u.SetNull(instructions.FieldRemark)
	return u
}

// SetTitle sets the "title" field.
func (u *InstructionsUpsert) SetTitle(v string) *InstructionsUpsert {
	u.Set(instructions.FieldTitle, v)
	return u
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *InstructionsUpsert) UpdateTitle() *InstructionsUpsert {
	u.SetExcluded(instructions.FieldTitle)
	return u
}

// SetContent sets the "content" field.
func (u *InstructionsUpsert) SetContent(v *interface{}) *InstructionsUpsert {
	u.Set(instructions.FieldContent, v)
	return u
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *InstructionsUpsert) UpdateContent() *InstructionsUpsert {
	u.SetExcluded(instructions.FieldContent)
	return u
}

// SetKey sets the "key" field.
func (u *InstructionsUpsert) SetKey(v string) *InstructionsUpsert {
	u.Set(instructions.FieldKey, v)
	return u
}

// UpdateKey sets the "key" field to the value that was provided on create.
func (u *InstructionsUpsert) UpdateKey() *InstructionsUpsert {
	u.SetExcluded(instructions.FieldKey)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Instructions.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *InstructionsUpsertOne) UpdateNewValues() *InstructionsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(instructions.FieldCreatedAt)
		}
		if _, exists := u.create.mutation.Creator(); exists {
			s.SetIgnore(instructions.FieldCreator)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Instructions.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *InstructionsUpsertOne) Ignore() *InstructionsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *InstructionsUpsertOne) DoNothing() *InstructionsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the InstructionsCreate.OnConflict
// documentation for more info.
func (u *InstructionsUpsertOne) Update(set func(*InstructionsUpsert)) *InstructionsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&InstructionsUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *InstructionsUpsertOne) SetUpdatedAt(v time.Time) *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *InstructionsUpsertOne) UpdateUpdatedAt() *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *InstructionsUpsertOne) SetDeletedAt(v time.Time) *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *InstructionsUpsertOne) UpdateDeletedAt() *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *InstructionsUpsertOne) ClearDeletedAt() *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *InstructionsUpsertOne) SetLastModifier(v *model.Modifier) *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *InstructionsUpsertOne) UpdateLastModifier() *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *InstructionsUpsertOne) ClearLastModifier() *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *InstructionsUpsertOne) SetRemark(v string) *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *InstructionsUpsertOne) UpdateRemark() *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *InstructionsUpsertOne) ClearRemark() *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.ClearRemark()
	})
}

// SetTitle sets the "title" field.
func (u *InstructionsUpsertOne) SetTitle(v string) *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetTitle(v)
	})
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *InstructionsUpsertOne) UpdateTitle() *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateTitle()
	})
}

// SetContent sets the "content" field.
func (u *InstructionsUpsertOne) SetContent(v *interface{}) *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *InstructionsUpsertOne) UpdateContent() *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateContent()
	})
}

// SetKey sets the "key" field.
func (u *InstructionsUpsertOne) SetKey(v string) *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetKey(v)
	})
}

// UpdateKey sets the "key" field to the value that was provided on create.
func (u *InstructionsUpsertOne) UpdateKey() *InstructionsUpsertOne {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateKey()
	})
}

// Exec executes the query.
func (u *InstructionsUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for InstructionsCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *InstructionsUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *InstructionsUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *InstructionsUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// InstructionsCreateBulk is the builder for creating many Instructions entities in bulk.
type InstructionsCreateBulk struct {
	config
	err      error
	builders []*InstructionsCreate
	conflict []sql.ConflictOption
}

// Save creates the Instructions entities in the database.
func (icb *InstructionsCreateBulk) Save(ctx context.Context) ([]*Instructions, error) {
	if icb.err != nil {
		return nil, icb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(icb.builders))
	nodes := make([]*Instructions, len(icb.builders))
	mutators := make([]Mutator, len(icb.builders))
	for i := range icb.builders {
		func(i int, root context.Context) {
			builder := icb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*InstructionsMutation)
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
					_, err = mutators[i+1].Mutate(root, icb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = icb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, icb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, icb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (icb *InstructionsCreateBulk) SaveX(ctx context.Context) []*Instructions {
	v, err := icb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (icb *InstructionsCreateBulk) Exec(ctx context.Context) error {
	_, err := icb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (icb *InstructionsCreateBulk) ExecX(ctx context.Context) {
	if err := icb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Instructions.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.InstructionsUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (icb *InstructionsCreateBulk) OnConflict(opts ...sql.ConflictOption) *InstructionsUpsertBulk {
	icb.conflict = opts
	return &InstructionsUpsertBulk{
		create: icb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Instructions.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (icb *InstructionsCreateBulk) OnConflictColumns(columns ...string) *InstructionsUpsertBulk {
	icb.conflict = append(icb.conflict, sql.ConflictColumns(columns...))
	return &InstructionsUpsertBulk{
		create: icb,
	}
}

// InstructionsUpsertBulk is the builder for "upsert"-ing
// a bulk of Instructions nodes.
type InstructionsUpsertBulk struct {
	create *InstructionsCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Instructions.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *InstructionsUpsertBulk) UpdateNewValues() *InstructionsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(instructions.FieldCreatedAt)
			}
			if _, exists := b.mutation.Creator(); exists {
				s.SetIgnore(instructions.FieldCreator)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Instructions.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *InstructionsUpsertBulk) Ignore() *InstructionsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *InstructionsUpsertBulk) DoNothing() *InstructionsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the InstructionsCreateBulk.OnConflict
// documentation for more info.
func (u *InstructionsUpsertBulk) Update(set func(*InstructionsUpsert)) *InstructionsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&InstructionsUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *InstructionsUpsertBulk) SetUpdatedAt(v time.Time) *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *InstructionsUpsertBulk) UpdateUpdatedAt() *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *InstructionsUpsertBulk) SetDeletedAt(v time.Time) *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *InstructionsUpsertBulk) UpdateDeletedAt() *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *InstructionsUpsertBulk) ClearDeletedAt() *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *InstructionsUpsertBulk) SetLastModifier(v *model.Modifier) *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *InstructionsUpsertBulk) UpdateLastModifier() *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *InstructionsUpsertBulk) ClearLastModifier() *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *InstructionsUpsertBulk) SetRemark(v string) *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *InstructionsUpsertBulk) UpdateRemark() *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *InstructionsUpsertBulk) ClearRemark() *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.ClearRemark()
	})
}

// SetTitle sets the "title" field.
func (u *InstructionsUpsertBulk) SetTitle(v string) *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetTitle(v)
	})
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *InstructionsUpsertBulk) UpdateTitle() *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateTitle()
	})
}

// SetContent sets the "content" field.
func (u *InstructionsUpsertBulk) SetContent(v *interface{}) *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *InstructionsUpsertBulk) UpdateContent() *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateContent()
	})
}

// SetKey sets the "key" field.
func (u *InstructionsUpsertBulk) SetKey(v string) *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.SetKey(v)
	})
}

// UpdateKey sets the "key" field to the value that was provided on create.
func (u *InstructionsUpsertBulk) UpdateKey() *InstructionsUpsertBulk {
	return u.Update(func(s *InstructionsUpsert) {
		s.UpdateKey()
	})
}

// Exec executes the query.
func (u *InstructionsUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the InstructionsCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for InstructionsCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *InstructionsUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}