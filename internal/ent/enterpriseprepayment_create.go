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
	"github.com/auroraride/aurservd/internal/ent/enterpriseprepayment"
)

// EnterprisePrepaymentCreate is the builder for creating a EnterprisePrepayment entity.
type EnterprisePrepaymentCreate struct {
	config
	mutation *EnterprisePrepaymentMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (epc *EnterprisePrepaymentCreate) SetCreatedAt(t time.Time) *EnterprisePrepaymentCreate {
	epc.mutation.SetCreatedAt(t)
	return epc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (epc *EnterprisePrepaymentCreate) SetNillableCreatedAt(t *time.Time) *EnterprisePrepaymentCreate {
	if t != nil {
		epc.SetCreatedAt(*t)
	}
	return epc
}

// SetUpdatedAt sets the "updated_at" field.
func (epc *EnterprisePrepaymentCreate) SetUpdatedAt(t time.Time) *EnterprisePrepaymentCreate {
	epc.mutation.SetUpdatedAt(t)
	return epc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (epc *EnterprisePrepaymentCreate) SetNillableUpdatedAt(t *time.Time) *EnterprisePrepaymentCreate {
	if t != nil {
		epc.SetUpdatedAt(*t)
	}
	return epc
}

// SetDeletedAt sets the "deleted_at" field.
func (epc *EnterprisePrepaymentCreate) SetDeletedAt(t time.Time) *EnterprisePrepaymentCreate {
	epc.mutation.SetDeletedAt(t)
	return epc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (epc *EnterprisePrepaymentCreate) SetNillableDeletedAt(t *time.Time) *EnterprisePrepaymentCreate {
	if t != nil {
		epc.SetDeletedAt(*t)
	}
	return epc
}

// SetCreator sets the "creator" field.
func (epc *EnterprisePrepaymentCreate) SetCreator(m *model.Modifier) *EnterprisePrepaymentCreate {
	epc.mutation.SetCreator(m)
	return epc
}

// SetLastModifier sets the "last_modifier" field.
func (epc *EnterprisePrepaymentCreate) SetLastModifier(m *model.Modifier) *EnterprisePrepaymentCreate {
	epc.mutation.SetLastModifier(m)
	return epc
}

// SetRemark sets the "remark" field.
func (epc *EnterprisePrepaymentCreate) SetRemark(s string) *EnterprisePrepaymentCreate {
	epc.mutation.SetRemark(s)
	return epc
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (epc *EnterprisePrepaymentCreate) SetNillableRemark(s *string) *EnterprisePrepaymentCreate {
	if s != nil {
		epc.SetRemark(*s)
	}
	return epc
}

// SetEnterpriseID sets the "enterprise_id" field.
func (epc *EnterprisePrepaymentCreate) SetEnterpriseID(u uint64) *EnterprisePrepaymentCreate {
	epc.mutation.SetEnterpriseID(u)
	return epc
}

// SetAmount sets the "amount" field.
func (epc *EnterprisePrepaymentCreate) SetAmount(f float64) *EnterprisePrepaymentCreate {
	epc.mutation.SetAmount(f)
	return epc
}

// SetEnterprise sets the "enterprise" edge to the Enterprise entity.
func (epc *EnterprisePrepaymentCreate) SetEnterprise(e *Enterprise) *EnterprisePrepaymentCreate {
	return epc.SetEnterpriseID(e.ID)
}

// Mutation returns the EnterprisePrepaymentMutation object of the builder.
func (epc *EnterprisePrepaymentCreate) Mutation() *EnterprisePrepaymentMutation {
	return epc.mutation
}

// Save creates the EnterprisePrepayment in the database.
func (epc *EnterprisePrepaymentCreate) Save(ctx context.Context) (*EnterprisePrepayment, error) {
	var (
		err  error
		node *EnterprisePrepayment
	)
	if err := epc.defaults(); err != nil {
		return nil, err
	}
	if len(epc.hooks) == 0 {
		if err = epc.check(); err != nil {
			return nil, err
		}
		node, err = epc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EnterprisePrepaymentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = epc.check(); err != nil {
				return nil, err
			}
			epc.mutation = mutation
			if node, err = epc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(epc.hooks) - 1; i >= 0; i-- {
			if epc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = epc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, epc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*EnterprisePrepayment)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from EnterprisePrepaymentMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (epc *EnterprisePrepaymentCreate) SaveX(ctx context.Context) *EnterprisePrepayment {
	v, err := epc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (epc *EnterprisePrepaymentCreate) Exec(ctx context.Context) error {
	_, err := epc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epc *EnterprisePrepaymentCreate) ExecX(ctx context.Context) {
	if err := epc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (epc *EnterprisePrepaymentCreate) defaults() error {
	if _, ok := epc.mutation.CreatedAt(); !ok {
		if enterpriseprepayment.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized enterpriseprepayment.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := enterpriseprepayment.DefaultCreatedAt()
		epc.mutation.SetCreatedAt(v)
	}
	if _, ok := epc.mutation.UpdatedAt(); !ok {
		if enterpriseprepayment.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized enterpriseprepayment.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := enterpriseprepayment.DefaultUpdatedAt()
		epc.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (epc *EnterprisePrepaymentCreate) check() error {
	if _, ok := epc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "EnterprisePrepayment.created_at"`)}
	}
	if _, ok := epc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "EnterprisePrepayment.updated_at"`)}
	}
	if _, ok := epc.mutation.EnterpriseID(); !ok {
		return &ValidationError{Name: "enterprise_id", err: errors.New(`ent: missing required field "EnterprisePrepayment.enterprise_id"`)}
	}
	if _, ok := epc.mutation.Amount(); !ok {
		return &ValidationError{Name: "amount", err: errors.New(`ent: missing required field "EnterprisePrepayment.amount"`)}
	}
	if _, ok := epc.mutation.EnterpriseID(); !ok {
		return &ValidationError{Name: "enterprise", err: errors.New(`ent: missing required edge "EnterprisePrepayment.enterprise"`)}
	}
	return nil
}

func (epc *EnterprisePrepaymentCreate) sqlSave(ctx context.Context) (*EnterprisePrepayment, error) {
	_node, _spec := epc.createSpec()
	if err := sqlgraph.CreateNode(ctx, epc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	return _node, nil
}

func (epc *EnterprisePrepaymentCreate) createSpec() (*EnterprisePrepayment, *sqlgraph.CreateSpec) {
	var (
		_node = &EnterprisePrepayment{config: epc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: enterpriseprepayment.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: enterpriseprepayment.FieldID,
			},
		}
	)
	_spec.OnConflict = epc.conflict
	if value, ok := epc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: enterpriseprepayment.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := epc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: enterpriseprepayment.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := epc.mutation.DeletedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: enterpriseprepayment.FieldDeletedAt,
		})
		_node.DeletedAt = &value
	}
	if value, ok := epc.mutation.Creator(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: enterpriseprepayment.FieldCreator,
		})
		_node.Creator = value
	}
	if value, ok := epc.mutation.LastModifier(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: enterpriseprepayment.FieldLastModifier,
		})
		_node.LastModifier = value
	}
	if value, ok := epc.mutation.Remark(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: enterpriseprepayment.FieldRemark,
		})
		_node.Remark = value
	}
	if value, ok := epc.mutation.Amount(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: enterpriseprepayment.FieldAmount,
		})
		_node.Amount = value
	}
	if nodes := epc.mutation.EnterpriseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   enterpriseprepayment.EnterpriseTable,
			Columns: []string{enterpriseprepayment.EnterpriseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: enterprise.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.EnterpriseID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.EnterprisePrepayment.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnterprisePrepaymentUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
//
func (epc *EnterprisePrepaymentCreate) OnConflict(opts ...sql.ConflictOption) *EnterprisePrepaymentUpsertOne {
	epc.conflict = opts
	return &EnterprisePrepaymentUpsertOne{
		create: epc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.EnterprisePrepayment.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (epc *EnterprisePrepaymentCreate) OnConflictColumns(columns ...string) *EnterprisePrepaymentUpsertOne {
	epc.conflict = append(epc.conflict, sql.ConflictColumns(columns...))
	return &EnterprisePrepaymentUpsertOne{
		create: epc,
	}
}

type (
	// EnterprisePrepaymentUpsertOne is the builder for "upsert"-ing
	//  one EnterprisePrepayment node.
	EnterprisePrepaymentUpsertOne struct {
		create *EnterprisePrepaymentCreate
	}

	// EnterprisePrepaymentUpsert is the "OnConflict" setter.
	EnterprisePrepaymentUpsert struct {
		*sql.UpdateSet
	}
)

// SetCreatedAt sets the "created_at" field.
func (u *EnterprisePrepaymentUpsert) SetCreatedAt(v time.Time) *EnterprisePrepaymentUpsert {
	u.Set(enterpriseprepayment.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsert) UpdateCreatedAt() *EnterprisePrepaymentUpsert {
	u.SetExcluded(enterpriseprepayment.FieldCreatedAt)
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *EnterprisePrepaymentUpsert) SetUpdatedAt(v time.Time) *EnterprisePrepaymentUpsert {
	u.Set(enterpriseprepayment.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsert) UpdateUpdatedAt() *EnterprisePrepaymentUpsert {
	u.SetExcluded(enterpriseprepayment.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *EnterprisePrepaymentUpsert) SetDeletedAt(v time.Time) *EnterprisePrepaymentUpsert {
	u.Set(enterpriseprepayment.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsert) UpdateDeletedAt() *EnterprisePrepaymentUpsert {
	u.SetExcluded(enterpriseprepayment.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *EnterprisePrepaymentUpsert) ClearDeletedAt() *EnterprisePrepaymentUpsert {
	u.SetNull(enterpriseprepayment.FieldDeletedAt)
	return u
}

// SetCreator sets the "creator" field.
func (u *EnterprisePrepaymentUpsert) SetCreator(v *model.Modifier) *EnterprisePrepaymentUpsert {
	u.Set(enterpriseprepayment.FieldCreator, v)
	return u
}

// UpdateCreator sets the "creator" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsert) UpdateCreator() *EnterprisePrepaymentUpsert {
	u.SetExcluded(enterpriseprepayment.FieldCreator)
	return u
}

// ClearCreator clears the value of the "creator" field.
func (u *EnterprisePrepaymentUpsert) ClearCreator() *EnterprisePrepaymentUpsert {
	u.SetNull(enterpriseprepayment.FieldCreator)
	return u
}

// SetLastModifier sets the "last_modifier" field.
func (u *EnterprisePrepaymentUpsert) SetLastModifier(v *model.Modifier) *EnterprisePrepaymentUpsert {
	u.Set(enterpriseprepayment.FieldLastModifier, v)
	return u
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsert) UpdateLastModifier() *EnterprisePrepaymentUpsert {
	u.SetExcluded(enterpriseprepayment.FieldLastModifier)
	return u
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *EnterprisePrepaymentUpsert) ClearLastModifier() *EnterprisePrepaymentUpsert {
	u.SetNull(enterpriseprepayment.FieldLastModifier)
	return u
}

// SetRemark sets the "remark" field.
func (u *EnterprisePrepaymentUpsert) SetRemark(v string) *EnterprisePrepaymentUpsert {
	u.Set(enterpriseprepayment.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsert) UpdateRemark() *EnterprisePrepaymentUpsert {
	u.SetExcluded(enterpriseprepayment.FieldRemark)
	return u
}

// ClearRemark clears the value of the "remark" field.
func (u *EnterprisePrepaymentUpsert) ClearRemark() *EnterprisePrepaymentUpsert {
	u.SetNull(enterpriseprepayment.FieldRemark)
	return u
}

// SetEnterpriseID sets the "enterprise_id" field.
func (u *EnterprisePrepaymentUpsert) SetEnterpriseID(v uint64) *EnterprisePrepaymentUpsert {
	u.Set(enterpriseprepayment.FieldEnterpriseID, v)
	return u
}

// UpdateEnterpriseID sets the "enterprise_id" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsert) UpdateEnterpriseID() *EnterprisePrepaymentUpsert {
	u.SetExcluded(enterpriseprepayment.FieldEnterpriseID)
	return u
}

// SetAmount sets the "amount" field.
func (u *EnterprisePrepaymentUpsert) SetAmount(v float64) *EnterprisePrepaymentUpsert {
	u.Set(enterpriseprepayment.FieldAmount, v)
	return u
}

// UpdateAmount sets the "amount" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsert) UpdateAmount() *EnterprisePrepaymentUpsert {
	u.SetExcluded(enterpriseprepayment.FieldAmount)
	return u
}

// AddAmount adds v to the "amount" field.
func (u *EnterprisePrepaymentUpsert) AddAmount(v float64) *EnterprisePrepaymentUpsert {
	u.Add(enterpriseprepayment.FieldAmount, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.EnterprisePrepayment.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
//
func (u *EnterprisePrepaymentUpsertOne) UpdateNewValues() *EnterprisePrepaymentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(enterpriseprepayment.FieldCreatedAt)
		}
		if _, exists := u.create.mutation.Creator(); exists {
			s.SetIgnore(enterpriseprepayment.FieldCreator)
		}
		if _, exists := u.create.mutation.Amount(); exists {
			s.SetIgnore(enterpriseprepayment.FieldAmount)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.EnterprisePrepayment.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *EnterprisePrepaymentUpsertOne) Ignore() *EnterprisePrepaymentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EnterprisePrepaymentUpsertOne) DoNothing() *EnterprisePrepaymentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EnterprisePrepaymentCreate.OnConflict
// documentation for more info.
func (u *EnterprisePrepaymentUpsertOne) Update(set func(*EnterprisePrepaymentUpsert)) *EnterprisePrepaymentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EnterprisePrepaymentUpsert{UpdateSet: update})
	}))
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *EnterprisePrepaymentUpsertOne) SetCreatedAt(v time.Time) *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertOne) UpdateCreatedAt() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *EnterprisePrepaymentUpsertOne) SetUpdatedAt(v time.Time) *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertOne) UpdateUpdatedAt() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *EnterprisePrepaymentUpsertOne) SetDeletedAt(v time.Time) *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertOne) UpdateDeletedAt() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *EnterprisePrepaymentUpsertOne) ClearDeletedAt() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.ClearDeletedAt()
	})
}

// SetCreator sets the "creator" field.
func (u *EnterprisePrepaymentUpsertOne) SetCreator(v *model.Modifier) *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetCreator(v)
	})
}

// UpdateCreator sets the "creator" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertOne) UpdateCreator() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateCreator()
	})
}

// ClearCreator clears the value of the "creator" field.
func (u *EnterprisePrepaymentUpsertOne) ClearCreator() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.ClearCreator()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *EnterprisePrepaymentUpsertOne) SetLastModifier(v *model.Modifier) *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertOne) UpdateLastModifier() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *EnterprisePrepaymentUpsertOne) ClearLastModifier() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *EnterprisePrepaymentUpsertOne) SetRemark(v string) *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertOne) UpdateRemark() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *EnterprisePrepaymentUpsertOne) ClearRemark() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.ClearRemark()
	})
}

// SetEnterpriseID sets the "enterprise_id" field.
func (u *EnterprisePrepaymentUpsertOne) SetEnterpriseID(v uint64) *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetEnterpriseID(v)
	})
}

// UpdateEnterpriseID sets the "enterprise_id" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertOne) UpdateEnterpriseID() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateEnterpriseID()
	})
}

// SetAmount sets the "amount" field.
func (u *EnterprisePrepaymentUpsertOne) SetAmount(v float64) *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetAmount(v)
	})
}

// AddAmount adds v to the "amount" field.
func (u *EnterprisePrepaymentUpsertOne) AddAmount(v float64) *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.AddAmount(v)
	})
}

// UpdateAmount sets the "amount" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertOne) UpdateAmount() *EnterprisePrepaymentUpsertOne {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateAmount()
	})
}

// Exec executes the query.
func (u *EnterprisePrepaymentUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for EnterprisePrepaymentCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EnterprisePrepaymentUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *EnterprisePrepaymentUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *EnterprisePrepaymentUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// EnterprisePrepaymentCreateBulk is the builder for creating many EnterprisePrepayment entities in bulk.
type EnterprisePrepaymentCreateBulk struct {
	config
	builders []*EnterprisePrepaymentCreate
	conflict []sql.ConflictOption
}

// Save creates the EnterprisePrepayment entities in the database.
func (epcb *EnterprisePrepaymentCreateBulk) Save(ctx context.Context) ([]*EnterprisePrepayment, error) {
	specs := make([]*sqlgraph.CreateSpec, len(epcb.builders))
	nodes := make([]*EnterprisePrepayment, len(epcb.builders))
	mutators := make([]Mutator, len(epcb.builders))
	for i := range epcb.builders {
		func(i int, root context.Context) {
			builder := epcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EnterprisePrepaymentMutation)
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
					_, err = mutators[i+1].Mutate(root, epcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = epcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, epcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, epcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (epcb *EnterprisePrepaymentCreateBulk) SaveX(ctx context.Context) []*EnterprisePrepayment {
	v, err := epcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (epcb *EnterprisePrepaymentCreateBulk) Exec(ctx context.Context) error {
	_, err := epcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epcb *EnterprisePrepaymentCreateBulk) ExecX(ctx context.Context) {
	if err := epcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.EnterprisePrepayment.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnterprisePrepaymentUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
//
func (epcb *EnterprisePrepaymentCreateBulk) OnConflict(opts ...sql.ConflictOption) *EnterprisePrepaymentUpsertBulk {
	epcb.conflict = opts
	return &EnterprisePrepaymentUpsertBulk{
		create: epcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.EnterprisePrepayment.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (epcb *EnterprisePrepaymentCreateBulk) OnConflictColumns(columns ...string) *EnterprisePrepaymentUpsertBulk {
	epcb.conflict = append(epcb.conflict, sql.ConflictColumns(columns...))
	return &EnterprisePrepaymentUpsertBulk{
		create: epcb,
	}
}

// EnterprisePrepaymentUpsertBulk is the builder for "upsert"-ing
// a bulk of EnterprisePrepayment nodes.
type EnterprisePrepaymentUpsertBulk struct {
	create *EnterprisePrepaymentCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.EnterprisePrepayment.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
//
func (u *EnterprisePrepaymentUpsertBulk) UpdateNewValues() *EnterprisePrepaymentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(enterpriseprepayment.FieldCreatedAt)
			}
			if _, exists := b.mutation.Creator(); exists {
				s.SetIgnore(enterpriseprepayment.FieldCreator)
			}
			if _, exists := b.mutation.Amount(); exists {
				s.SetIgnore(enterpriseprepayment.FieldAmount)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.EnterprisePrepayment.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
//
func (u *EnterprisePrepaymentUpsertBulk) Ignore() *EnterprisePrepaymentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EnterprisePrepaymentUpsertBulk) DoNothing() *EnterprisePrepaymentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EnterprisePrepaymentCreateBulk.OnConflict
// documentation for more info.
func (u *EnterprisePrepaymentUpsertBulk) Update(set func(*EnterprisePrepaymentUpsert)) *EnterprisePrepaymentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EnterprisePrepaymentUpsert{UpdateSet: update})
	}))
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *EnterprisePrepaymentUpsertBulk) SetCreatedAt(v time.Time) *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertBulk) UpdateCreatedAt() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *EnterprisePrepaymentUpsertBulk) SetUpdatedAt(v time.Time) *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertBulk) UpdateUpdatedAt() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *EnterprisePrepaymentUpsertBulk) SetDeletedAt(v time.Time) *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertBulk) UpdateDeletedAt() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *EnterprisePrepaymentUpsertBulk) ClearDeletedAt() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.ClearDeletedAt()
	})
}

// SetCreator sets the "creator" field.
func (u *EnterprisePrepaymentUpsertBulk) SetCreator(v *model.Modifier) *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetCreator(v)
	})
}

// UpdateCreator sets the "creator" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertBulk) UpdateCreator() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateCreator()
	})
}

// ClearCreator clears the value of the "creator" field.
func (u *EnterprisePrepaymentUpsertBulk) ClearCreator() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.ClearCreator()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *EnterprisePrepaymentUpsertBulk) SetLastModifier(v *model.Modifier) *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertBulk) UpdateLastModifier() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *EnterprisePrepaymentUpsertBulk) ClearLastModifier() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *EnterprisePrepaymentUpsertBulk) SetRemark(v string) *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertBulk) UpdateRemark() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *EnterprisePrepaymentUpsertBulk) ClearRemark() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.ClearRemark()
	})
}

// SetEnterpriseID sets the "enterprise_id" field.
func (u *EnterprisePrepaymentUpsertBulk) SetEnterpriseID(v uint64) *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetEnterpriseID(v)
	})
}

// UpdateEnterpriseID sets the "enterprise_id" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertBulk) UpdateEnterpriseID() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateEnterpriseID()
	})
}

// SetAmount sets the "amount" field.
func (u *EnterprisePrepaymentUpsertBulk) SetAmount(v float64) *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.SetAmount(v)
	})
}

// AddAmount adds v to the "amount" field.
func (u *EnterprisePrepaymentUpsertBulk) AddAmount(v float64) *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.AddAmount(v)
	})
}

// UpdateAmount sets the "amount" field to the value that was provided on create.
func (u *EnterprisePrepaymentUpsertBulk) UpdateAmount() *EnterprisePrepaymentUpsertBulk {
	return u.Update(func(s *EnterprisePrepaymentUpsert) {
		s.UpdateAmount()
	})
}

// Exec executes the query.
func (u *EnterprisePrepaymentUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the EnterprisePrepaymentCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for EnterprisePrepaymentCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EnterprisePrepaymentUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}