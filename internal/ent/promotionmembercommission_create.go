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
	"github.com/auroraride/aurservd/internal/ent/promotioncommission"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/promotionmembercommission"
)

// PromotionMemberCommissionCreate is the builder for creating a PromotionMemberCommission entity.
type PromotionMemberCommissionCreate struct {
	config
	mutation *PromotionMemberCommissionMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (pmcc *PromotionMemberCommissionCreate) SetCreatedAt(t time.Time) *PromotionMemberCommissionCreate {
	pmcc.mutation.SetCreatedAt(t)
	return pmcc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (pmcc *PromotionMemberCommissionCreate) SetNillableCreatedAt(t *time.Time) *PromotionMemberCommissionCreate {
	if t != nil {
		pmcc.SetCreatedAt(*t)
	}
	return pmcc
}

// SetUpdatedAt sets the "updated_at" field.
func (pmcc *PromotionMemberCommissionCreate) SetUpdatedAt(t time.Time) *PromotionMemberCommissionCreate {
	pmcc.mutation.SetUpdatedAt(t)
	return pmcc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (pmcc *PromotionMemberCommissionCreate) SetNillableUpdatedAt(t *time.Time) *PromotionMemberCommissionCreate {
	if t != nil {
		pmcc.SetUpdatedAt(*t)
	}
	return pmcc
}

// SetDeletedAt sets the "deleted_at" field.
func (pmcc *PromotionMemberCommissionCreate) SetDeletedAt(t time.Time) *PromotionMemberCommissionCreate {
	pmcc.mutation.SetDeletedAt(t)
	return pmcc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (pmcc *PromotionMemberCommissionCreate) SetNillableDeletedAt(t *time.Time) *PromotionMemberCommissionCreate {
	if t != nil {
		pmcc.SetDeletedAt(*t)
	}
	return pmcc
}

// SetCommissionID sets the "commission_id" field.
func (pmcc *PromotionMemberCommissionCreate) SetCommissionID(u uint64) *PromotionMemberCommissionCreate {
	pmcc.mutation.SetCommissionID(u)
	return pmcc
}

// SetMemberID sets the "member_id" field.
func (pmcc *PromotionMemberCommissionCreate) SetMemberID(u uint64) *PromotionMemberCommissionCreate {
	pmcc.mutation.SetMemberID(u)
	return pmcc
}

// SetNillableMemberID sets the "member_id" field if the given value is not nil.
func (pmcc *PromotionMemberCommissionCreate) SetNillableMemberID(u *uint64) *PromotionMemberCommissionCreate {
	if u != nil {
		pmcc.SetMemberID(*u)
	}
	return pmcc
}

// SetCommission sets the "commission" edge to the PromotionCommission entity.
func (pmcc *PromotionMemberCommissionCreate) SetCommission(p *PromotionCommission) *PromotionMemberCommissionCreate {
	return pmcc.SetCommissionID(p.ID)
}

// SetMember sets the "member" edge to the PromotionMember entity.
func (pmcc *PromotionMemberCommissionCreate) SetMember(p *PromotionMember) *PromotionMemberCommissionCreate {
	return pmcc.SetMemberID(p.ID)
}

// Mutation returns the PromotionMemberCommissionMutation object of the builder.
func (pmcc *PromotionMemberCommissionCreate) Mutation() *PromotionMemberCommissionMutation {
	return pmcc.mutation
}

// Save creates the PromotionMemberCommission in the database.
func (pmcc *PromotionMemberCommissionCreate) Save(ctx context.Context) (*PromotionMemberCommission, error) {
	pmcc.defaults()
	return withHooks(ctx, pmcc.sqlSave, pmcc.mutation, pmcc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (pmcc *PromotionMemberCommissionCreate) SaveX(ctx context.Context) *PromotionMemberCommission {
	v, err := pmcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pmcc *PromotionMemberCommissionCreate) Exec(ctx context.Context) error {
	_, err := pmcc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pmcc *PromotionMemberCommissionCreate) ExecX(ctx context.Context) {
	if err := pmcc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pmcc *PromotionMemberCommissionCreate) defaults() {
	if _, ok := pmcc.mutation.CreatedAt(); !ok {
		v := promotionmembercommission.DefaultCreatedAt()
		pmcc.mutation.SetCreatedAt(v)
	}
	if _, ok := pmcc.mutation.UpdatedAt(); !ok {
		v := promotionmembercommission.DefaultUpdatedAt()
		pmcc.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pmcc *PromotionMemberCommissionCreate) check() error {
	if _, ok := pmcc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "PromotionMemberCommission.created_at"`)}
	}
	if _, ok := pmcc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "PromotionMemberCommission.updated_at"`)}
	}
	if _, ok := pmcc.mutation.CommissionID(); !ok {
		return &ValidationError{Name: "commission_id", err: errors.New(`ent: missing required field "PromotionMemberCommission.commission_id"`)}
	}
	if len(pmcc.mutation.CommissionIDs()) == 0 {
		return &ValidationError{Name: "commission", err: errors.New(`ent: missing required edge "PromotionMemberCommission.commission"`)}
	}
	return nil
}

func (pmcc *PromotionMemberCommissionCreate) sqlSave(ctx context.Context) (*PromotionMemberCommission, error) {
	if err := pmcc.check(); err != nil {
		return nil, err
	}
	_node, _spec := pmcc.createSpec()
	if err := sqlgraph.CreateNode(ctx, pmcc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	pmcc.mutation.id = &_node.ID
	pmcc.mutation.done = true
	return _node, nil
}

func (pmcc *PromotionMemberCommissionCreate) createSpec() (*PromotionMemberCommission, *sqlgraph.CreateSpec) {
	var (
		_node = &PromotionMemberCommission{config: pmcc.config}
		_spec = sqlgraph.NewCreateSpec(promotionmembercommission.Table, sqlgraph.NewFieldSpec(promotionmembercommission.FieldID, field.TypeUint64))
	)
	_spec.OnConflict = pmcc.conflict
	if value, ok := pmcc.mutation.CreatedAt(); ok {
		_spec.SetField(promotionmembercommission.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := pmcc.mutation.UpdatedAt(); ok {
		_spec.SetField(promotionmembercommission.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := pmcc.mutation.DeletedAt(); ok {
		_spec.SetField(promotionmembercommission.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	if nodes := pmcc.mutation.CommissionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   promotionmembercommission.CommissionTable,
			Columns: []string{promotionmembercommission.CommissionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(promotioncommission.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.CommissionID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pmcc.mutation.MemberIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   promotionmembercommission.MemberTable,
			Columns: []string{promotionmembercommission.MemberColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(promotionmember.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.MemberID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.PromotionMemberCommission.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PromotionMemberCommissionUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (pmcc *PromotionMemberCommissionCreate) OnConflict(opts ...sql.ConflictOption) *PromotionMemberCommissionUpsertOne {
	pmcc.conflict = opts
	return &PromotionMemberCommissionUpsertOne{
		create: pmcc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.PromotionMemberCommission.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (pmcc *PromotionMemberCommissionCreate) OnConflictColumns(columns ...string) *PromotionMemberCommissionUpsertOne {
	pmcc.conflict = append(pmcc.conflict, sql.ConflictColumns(columns...))
	return &PromotionMemberCommissionUpsertOne{
		create: pmcc,
	}
}

type (
	// PromotionMemberCommissionUpsertOne is the builder for "upsert"-ing
	//  one PromotionMemberCommission node.
	PromotionMemberCommissionUpsertOne struct {
		create *PromotionMemberCommissionCreate
	}

	// PromotionMemberCommissionUpsert is the "OnConflict" setter.
	PromotionMemberCommissionUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *PromotionMemberCommissionUpsert) SetUpdatedAt(v time.Time) *PromotionMemberCommissionUpsert {
	u.Set(promotionmembercommission.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsert) UpdateUpdatedAt() *PromotionMemberCommissionUpsert {
	u.SetExcluded(promotionmembercommission.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *PromotionMemberCommissionUpsert) SetDeletedAt(v time.Time) *PromotionMemberCommissionUpsert {
	u.Set(promotionmembercommission.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsert) UpdateDeletedAt() *PromotionMemberCommissionUpsert {
	u.SetExcluded(promotionmembercommission.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *PromotionMemberCommissionUpsert) ClearDeletedAt() *PromotionMemberCommissionUpsert {
	u.SetNull(promotionmembercommission.FieldDeletedAt)
	return u
}

// SetCommissionID sets the "commission_id" field.
func (u *PromotionMemberCommissionUpsert) SetCommissionID(v uint64) *PromotionMemberCommissionUpsert {
	u.Set(promotionmembercommission.FieldCommissionID, v)
	return u
}

// UpdateCommissionID sets the "commission_id" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsert) UpdateCommissionID() *PromotionMemberCommissionUpsert {
	u.SetExcluded(promotionmembercommission.FieldCommissionID)
	return u
}

// SetMemberID sets the "member_id" field.
func (u *PromotionMemberCommissionUpsert) SetMemberID(v uint64) *PromotionMemberCommissionUpsert {
	u.Set(promotionmembercommission.FieldMemberID, v)
	return u
}

// UpdateMemberID sets the "member_id" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsert) UpdateMemberID() *PromotionMemberCommissionUpsert {
	u.SetExcluded(promotionmembercommission.FieldMemberID)
	return u
}

// ClearMemberID clears the value of the "member_id" field.
func (u *PromotionMemberCommissionUpsert) ClearMemberID() *PromotionMemberCommissionUpsert {
	u.SetNull(promotionmembercommission.FieldMemberID)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.PromotionMemberCommission.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *PromotionMemberCommissionUpsertOne) UpdateNewValues() *PromotionMemberCommissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(promotionmembercommission.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.PromotionMemberCommission.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *PromotionMemberCommissionUpsertOne) Ignore() *PromotionMemberCommissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PromotionMemberCommissionUpsertOne) DoNothing() *PromotionMemberCommissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PromotionMemberCommissionCreate.OnConflict
// documentation for more info.
func (u *PromotionMemberCommissionUpsertOne) Update(set func(*PromotionMemberCommissionUpsert)) *PromotionMemberCommissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PromotionMemberCommissionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *PromotionMemberCommissionUpsertOne) SetUpdatedAt(v time.Time) *PromotionMemberCommissionUpsertOne {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsertOne) UpdateUpdatedAt() *PromotionMemberCommissionUpsertOne {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *PromotionMemberCommissionUpsertOne) SetDeletedAt(v time.Time) *PromotionMemberCommissionUpsertOne {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsertOne) UpdateDeletedAt() *PromotionMemberCommissionUpsertOne {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *PromotionMemberCommissionUpsertOne) ClearDeletedAt() *PromotionMemberCommissionUpsertOne {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.ClearDeletedAt()
	})
}

// SetCommissionID sets the "commission_id" field.
func (u *PromotionMemberCommissionUpsertOne) SetCommissionID(v uint64) *PromotionMemberCommissionUpsertOne {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.SetCommissionID(v)
	})
}

// UpdateCommissionID sets the "commission_id" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsertOne) UpdateCommissionID() *PromotionMemberCommissionUpsertOne {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.UpdateCommissionID()
	})
}

// SetMemberID sets the "member_id" field.
func (u *PromotionMemberCommissionUpsertOne) SetMemberID(v uint64) *PromotionMemberCommissionUpsertOne {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.SetMemberID(v)
	})
}

// UpdateMemberID sets the "member_id" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsertOne) UpdateMemberID() *PromotionMemberCommissionUpsertOne {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.UpdateMemberID()
	})
}

// ClearMemberID clears the value of the "member_id" field.
func (u *PromotionMemberCommissionUpsertOne) ClearMemberID() *PromotionMemberCommissionUpsertOne {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.ClearMemberID()
	})
}

// Exec executes the query.
func (u *PromotionMemberCommissionUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for PromotionMemberCommissionCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PromotionMemberCommissionUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *PromotionMemberCommissionUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *PromotionMemberCommissionUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// PromotionMemberCommissionCreateBulk is the builder for creating many PromotionMemberCommission entities in bulk.
type PromotionMemberCommissionCreateBulk struct {
	config
	err      error
	builders []*PromotionMemberCommissionCreate
	conflict []sql.ConflictOption
}

// Save creates the PromotionMemberCommission entities in the database.
func (pmccb *PromotionMemberCommissionCreateBulk) Save(ctx context.Context) ([]*PromotionMemberCommission, error) {
	if pmccb.err != nil {
		return nil, pmccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(pmccb.builders))
	nodes := make([]*PromotionMemberCommission, len(pmccb.builders))
	mutators := make([]Mutator, len(pmccb.builders))
	for i := range pmccb.builders {
		func(i int, root context.Context) {
			builder := pmccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PromotionMemberCommissionMutation)
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
					_, err = mutators[i+1].Mutate(root, pmccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = pmccb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pmccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, pmccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (pmccb *PromotionMemberCommissionCreateBulk) SaveX(ctx context.Context) []*PromotionMemberCommission {
	v, err := pmccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pmccb *PromotionMemberCommissionCreateBulk) Exec(ctx context.Context) error {
	_, err := pmccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pmccb *PromotionMemberCommissionCreateBulk) ExecX(ctx context.Context) {
	if err := pmccb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.PromotionMemberCommission.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PromotionMemberCommissionUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (pmccb *PromotionMemberCommissionCreateBulk) OnConflict(opts ...sql.ConflictOption) *PromotionMemberCommissionUpsertBulk {
	pmccb.conflict = opts
	return &PromotionMemberCommissionUpsertBulk{
		create: pmccb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.PromotionMemberCommission.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (pmccb *PromotionMemberCommissionCreateBulk) OnConflictColumns(columns ...string) *PromotionMemberCommissionUpsertBulk {
	pmccb.conflict = append(pmccb.conflict, sql.ConflictColumns(columns...))
	return &PromotionMemberCommissionUpsertBulk{
		create: pmccb,
	}
}

// PromotionMemberCommissionUpsertBulk is the builder for "upsert"-ing
// a bulk of PromotionMemberCommission nodes.
type PromotionMemberCommissionUpsertBulk struct {
	create *PromotionMemberCommissionCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.PromotionMemberCommission.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *PromotionMemberCommissionUpsertBulk) UpdateNewValues() *PromotionMemberCommissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(promotionmembercommission.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.PromotionMemberCommission.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *PromotionMemberCommissionUpsertBulk) Ignore() *PromotionMemberCommissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PromotionMemberCommissionUpsertBulk) DoNothing() *PromotionMemberCommissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PromotionMemberCommissionCreateBulk.OnConflict
// documentation for more info.
func (u *PromotionMemberCommissionUpsertBulk) Update(set func(*PromotionMemberCommissionUpsert)) *PromotionMemberCommissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PromotionMemberCommissionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *PromotionMemberCommissionUpsertBulk) SetUpdatedAt(v time.Time) *PromotionMemberCommissionUpsertBulk {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsertBulk) UpdateUpdatedAt() *PromotionMemberCommissionUpsertBulk {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *PromotionMemberCommissionUpsertBulk) SetDeletedAt(v time.Time) *PromotionMemberCommissionUpsertBulk {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsertBulk) UpdateDeletedAt() *PromotionMemberCommissionUpsertBulk {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *PromotionMemberCommissionUpsertBulk) ClearDeletedAt() *PromotionMemberCommissionUpsertBulk {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.ClearDeletedAt()
	})
}

// SetCommissionID sets the "commission_id" field.
func (u *PromotionMemberCommissionUpsertBulk) SetCommissionID(v uint64) *PromotionMemberCommissionUpsertBulk {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.SetCommissionID(v)
	})
}

// UpdateCommissionID sets the "commission_id" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsertBulk) UpdateCommissionID() *PromotionMemberCommissionUpsertBulk {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.UpdateCommissionID()
	})
}

// SetMemberID sets the "member_id" field.
func (u *PromotionMemberCommissionUpsertBulk) SetMemberID(v uint64) *PromotionMemberCommissionUpsertBulk {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.SetMemberID(v)
	})
}

// UpdateMemberID sets the "member_id" field to the value that was provided on create.
func (u *PromotionMemberCommissionUpsertBulk) UpdateMemberID() *PromotionMemberCommissionUpsertBulk {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.UpdateMemberID()
	})
}

// ClearMemberID clears the value of the "member_id" field.
func (u *PromotionMemberCommissionUpsertBulk) ClearMemberID() *PromotionMemberCommissionUpsertBulk {
	return u.Update(func(s *PromotionMemberCommissionUpsert) {
		s.ClearMemberID()
	})
}

// Exec executes the query.
func (u *PromotionMemberCommissionUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the PromotionMemberCommissionCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for PromotionMemberCommissionCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PromotionMemberCommissionUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
