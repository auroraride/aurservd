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
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/promotionperson"
)

// PromotionPersonCreate is the builder for creating a PromotionPerson entity.
type PromotionPersonCreate struct {
	config
	mutation *PromotionPersonMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (ppc *PromotionPersonCreate) SetCreatedAt(t time.Time) *PromotionPersonCreate {
	ppc.mutation.SetCreatedAt(t)
	return ppc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (ppc *PromotionPersonCreate) SetNillableCreatedAt(t *time.Time) *PromotionPersonCreate {
	if t != nil {
		ppc.SetCreatedAt(*t)
	}
	return ppc
}

// SetUpdatedAt sets the "updated_at" field.
func (ppc *PromotionPersonCreate) SetUpdatedAt(t time.Time) *PromotionPersonCreate {
	ppc.mutation.SetUpdatedAt(t)
	return ppc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (ppc *PromotionPersonCreate) SetNillableUpdatedAt(t *time.Time) *PromotionPersonCreate {
	if t != nil {
		ppc.SetUpdatedAt(*t)
	}
	return ppc
}

// SetDeletedAt sets the "deleted_at" field.
func (ppc *PromotionPersonCreate) SetDeletedAt(t time.Time) *PromotionPersonCreate {
	ppc.mutation.SetDeletedAt(t)
	return ppc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (ppc *PromotionPersonCreate) SetNillableDeletedAt(t *time.Time) *PromotionPersonCreate {
	if t != nil {
		ppc.SetDeletedAt(*t)
	}
	return ppc
}

// SetCreator sets the "creator" field.
func (ppc *PromotionPersonCreate) SetCreator(m *model.Modifier) *PromotionPersonCreate {
	ppc.mutation.SetCreator(m)
	return ppc
}

// SetLastModifier sets the "last_modifier" field.
func (ppc *PromotionPersonCreate) SetLastModifier(m *model.Modifier) *PromotionPersonCreate {
	ppc.mutation.SetLastModifier(m)
	return ppc
}

// SetRemark sets the "remark" field.
func (ppc *PromotionPersonCreate) SetRemark(s string) *PromotionPersonCreate {
	ppc.mutation.SetRemark(s)
	return ppc
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (ppc *PromotionPersonCreate) SetNillableRemark(s *string) *PromotionPersonCreate {
	if s != nil {
		ppc.SetRemark(*s)
	}
	return ppc
}

// SetStatus sets the "status" field.
func (ppc *PromotionPersonCreate) SetStatus(u uint8) *PromotionPersonCreate {
	ppc.mutation.SetStatus(u)
	return ppc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (ppc *PromotionPersonCreate) SetNillableStatus(u *uint8) *PromotionPersonCreate {
	if u != nil {
		ppc.SetStatus(*u)
	}
	return ppc
}

// SetName sets the "name" field.
func (ppc *PromotionPersonCreate) SetName(s string) *PromotionPersonCreate {
	ppc.mutation.SetName(s)
	return ppc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ppc *PromotionPersonCreate) SetNillableName(s *string) *PromotionPersonCreate {
	if s != nil {
		ppc.SetName(*s)
	}
	return ppc
}

// SetIDCardNumber sets the "id_card_number" field.
func (ppc *PromotionPersonCreate) SetIDCardNumber(s string) *PromotionPersonCreate {
	ppc.mutation.SetIDCardNumber(s)
	return ppc
}

// SetNillableIDCardNumber sets the "id_card_number" field if the given value is not nil.
func (ppc *PromotionPersonCreate) SetNillableIDCardNumber(s *string) *PromotionPersonCreate {
	if s != nil {
		ppc.SetIDCardNumber(*s)
	}
	return ppc
}

// SetAddress sets the "address" field.
func (ppc *PromotionPersonCreate) SetAddress(s string) *PromotionPersonCreate {
	ppc.mutation.SetAddress(s)
	return ppc
}

// SetNillableAddress sets the "address" field if the given value is not nil.
func (ppc *PromotionPersonCreate) SetNillableAddress(s *string) *PromotionPersonCreate {
	if s != nil {
		ppc.SetAddress(*s)
	}
	return ppc
}

// AddMemberIDs adds the "member" edge to the PromotionMember entity by IDs.
func (ppc *PromotionPersonCreate) AddMemberIDs(ids ...uint64) *PromotionPersonCreate {
	ppc.mutation.AddMemberIDs(ids...)
	return ppc
}

// AddMember adds the "member" edges to the PromotionMember entity.
func (ppc *PromotionPersonCreate) AddMember(p ...*PromotionMember) *PromotionPersonCreate {
	ids := make([]uint64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ppc.AddMemberIDs(ids...)
}

// Mutation returns the PromotionPersonMutation object of the builder.
func (ppc *PromotionPersonCreate) Mutation() *PromotionPersonMutation {
	return ppc.mutation
}

// Save creates the PromotionPerson in the database.
func (ppc *PromotionPersonCreate) Save(ctx context.Context) (*PromotionPerson, error) {
	if err := ppc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, ppc.sqlSave, ppc.mutation, ppc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ppc *PromotionPersonCreate) SaveX(ctx context.Context) *PromotionPerson {
	v, err := ppc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ppc *PromotionPersonCreate) Exec(ctx context.Context) error {
	_, err := ppc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ppc *PromotionPersonCreate) ExecX(ctx context.Context) {
	if err := ppc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ppc *PromotionPersonCreate) defaults() error {
	if _, ok := ppc.mutation.CreatedAt(); !ok {
		if promotionperson.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized promotionperson.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := promotionperson.DefaultCreatedAt()
		ppc.mutation.SetCreatedAt(v)
	}
	if _, ok := ppc.mutation.UpdatedAt(); !ok {
		if promotionperson.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized promotionperson.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := promotionperson.DefaultUpdatedAt()
		ppc.mutation.SetUpdatedAt(v)
	}
	if _, ok := ppc.mutation.Status(); !ok {
		v := promotionperson.DefaultStatus
		ppc.mutation.SetStatus(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (ppc *PromotionPersonCreate) check() error {
	if _, ok := ppc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "PromotionPerson.created_at"`)}
	}
	if _, ok := ppc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "PromotionPerson.updated_at"`)}
	}
	if _, ok := ppc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "PromotionPerson.status"`)}
	}
	if v, ok := ppc.mutation.Name(); ok {
		if err := promotionperson.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "PromotionPerson.name": %w`, err)}
		}
	}
	if v, ok := ppc.mutation.IDCardNumber(); ok {
		if err := promotionperson.IDCardNumberValidator(v); err != nil {
			return &ValidationError{Name: "id_card_number", err: fmt.Errorf(`ent: validator failed for field "PromotionPerson.id_card_number": %w`, err)}
		}
	}
	return nil
}

func (ppc *PromotionPersonCreate) sqlSave(ctx context.Context) (*PromotionPerson, error) {
	if err := ppc.check(); err != nil {
		return nil, err
	}
	_node, _spec := ppc.createSpec()
	if err := sqlgraph.CreateNode(ctx, ppc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	ppc.mutation.id = &_node.ID
	ppc.mutation.done = true
	return _node, nil
}

func (ppc *PromotionPersonCreate) createSpec() (*PromotionPerson, *sqlgraph.CreateSpec) {
	var (
		_node = &PromotionPerson{config: ppc.config}
		_spec = sqlgraph.NewCreateSpec(promotionperson.Table, sqlgraph.NewFieldSpec(promotionperson.FieldID, field.TypeUint64))
	)
	_spec.OnConflict = ppc.conflict
	if value, ok := ppc.mutation.CreatedAt(); ok {
		_spec.SetField(promotionperson.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := ppc.mutation.UpdatedAt(); ok {
		_spec.SetField(promotionperson.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := ppc.mutation.DeletedAt(); ok {
		_spec.SetField(promotionperson.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	if value, ok := ppc.mutation.Creator(); ok {
		_spec.SetField(promotionperson.FieldCreator, field.TypeJSON, value)
		_node.Creator = value
	}
	if value, ok := ppc.mutation.LastModifier(); ok {
		_spec.SetField(promotionperson.FieldLastModifier, field.TypeJSON, value)
		_node.LastModifier = value
	}
	if value, ok := ppc.mutation.Remark(); ok {
		_spec.SetField(promotionperson.FieldRemark, field.TypeString, value)
		_node.Remark = value
	}
	if value, ok := ppc.mutation.Status(); ok {
		_spec.SetField(promotionperson.FieldStatus, field.TypeUint8, value)
		_node.Status = value
	}
	if value, ok := ppc.mutation.Name(); ok {
		_spec.SetField(promotionperson.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := ppc.mutation.IDCardNumber(); ok {
		_spec.SetField(promotionperson.FieldIDCardNumber, field.TypeString, value)
		_node.IDCardNumber = value
	}
	if value, ok := ppc.mutation.Address(); ok {
		_spec.SetField(promotionperson.FieldAddress, field.TypeString, value)
		_node.Address = value
	}
	if nodes := ppc.mutation.MemberIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   promotionperson.MemberTable,
			Columns: []string{promotionperson.MemberColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(promotionmember.FieldID, field.TypeUint64),
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
//	client.PromotionPerson.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PromotionPersonUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (ppc *PromotionPersonCreate) OnConflict(opts ...sql.ConflictOption) *PromotionPersonUpsertOne {
	ppc.conflict = opts
	return &PromotionPersonUpsertOne{
		create: ppc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.PromotionPerson.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ppc *PromotionPersonCreate) OnConflictColumns(columns ...string) *PromotionPersonUpsertOne {
	ppc.conflict = append(ppc.conflict, sql.ConflictColumns(columns...))
	return &PromotionPersonUpsertOne{
		create: ppc,
	}
}

type (
	// PromotionPersonUpsertOne is the builder for "upsert"-ing
	//  one PromotionPerson node.
	PromotionPersonUpsertOne struct {
		create *PromotionPersonCreate
	}

	// PromotionPersonUpsert is the "OnConflict" setter.
	PromotionPersonUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *PromotionPersonUpsert) SetUpdatedAt(v time.Time) *PromotionPersonUpsert {
	u.Set(promotionperson.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *PromotionPersonUpsert) UpdateUpdatedAt() *PromotionPersonUpsert {
	u.SetExcluded(promotionperson.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *PromotionPersonUpsert) SetDeletedAt(v time.Time) *PromotionPersonUpsert {
	u.Set(promotionperson.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *PromotionPersonUpsert) UpdateDeletedAt() *PromotionPersonUpsert {
	u.SetExcluded(promotionperson.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *PromotionPersonUpsert) ClearDeletedAt() *PromotionPersonUpsert {
	u.SetNull(promotionperson.FieldDeletedAt)
	return u
}

// SetLastModifier sets the "last_modifier" field.
func (u *PromotionPersonUpsert) SetLastModifier(v *model.Modifier) *PromotionPersonUpsert {
	u.Set(promotionperson.FieldLastModifier, v)
	return u
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *PromotionPersonUpsert) UpdateLastModifier() *PromotionPersonUpsert {
	u.SetExcluded(promotionperson.FieldLastModifier)
	return u
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *PromotionPersonUpsert) ClearLastModifier() *PromotionPersonUpsert {
	u.SetNull(promotionperson.FieldLastModifier)
	return u
}

// SetRemark sets the "remark" field.
func (u *PromotionPersonUpsert) SetRemark(v string) *PromotionPersonUpsert {
	u.Set(promotionperson.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *PromotionPersonUpsert) UpdateRemark() *PromotionPersonUpsert {
	u.SetExcluded(promotionperson.FieldRemark)
	return u
}

// ClearRemark clears the value of the "remark" field.
func (u *PromotionPersonUpsert) ClearRemark() *PromotionPersonUpsert {
	u.SetNull(promotionperson.FieldRemark)
	return u
}

// SetStatus sets the "status" field.
func (u *PromotionPersonUpsert) SetStatus(v uint8) *PromotionPersonUpsert {
	u.Set(promotionperson.FieldStatus, v)
	return u
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *PromotionPersonUpsert) UpdateStatus() *PromotionPersonUpsert {
	u.SetExcluded(promotionperson.FieldStatus)
	return u
}

// AddStatus adds v to the "status" field.
func (u *PromotionPersonUpsert) AddStatus(v uint8) *PromotionPersonUpsert {
	u.Add(promotionperson.FieldStatus, v)
	return u
}

// SetName sets the "name" field.
func (u *PromotionPersonUpsert) SetName(v string) *PromotionPersonUpsert {
	u.Set(promotionperson.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *PromotionPersonUpsert) UpdateName() *PromotionPersonUpsert {
	u.SetExcluded(promotionperson.FieldName)
	return u
}

// ClearName clears the value of the "name" field.
func (u *PromotionPersonUpsert) ClearName() *PromotionPersonUpsert {
	u.SetNull(promotionperson.FieldName)
	return u
}

// SetIDCardNumber sets the "id_card_number" field.
func (u *PromotionPersonUpsert) SetIDCardNumber(v string) *PromotionPersonUpsert {
	u.Set(promotionperson.FieldIDCardNumber, v)
	return u
}

// UpdateIDCardNumber sets the "id_card_number" field to the value that was provided on create.
func (u *PromotionPersonUpsert) UpdateIDCardNumber() *PromotionPersonUpsert {
	u.SetExcluded(promotionperson.FieldIDCardNumber)
	return u
}

// ClearIDCardNumber clears the value of the "id_card_number" field.
func (u *PromotionPersonUpsert) ClearIDCardNumber() *PromotionPersonUpsert {
	u.SetNull(promotionperson.FieldIDCardNumber)
	return u
}

// SetAddress sets the "address" field.
func (u *PromotionPersonUpsert) SetAddress(v string) *PromotionPersonUpsert {
	u.Set(promotionperson.FieldAddress, v)
	return u
}

// UpdateAddress sets the "address" field to the value that was provided on create.
func (u *PromotionPersonUpsert) UpdateAddress() *PromotionPersonUpsert {
	u.SetExcluded(promotionperson.FieldAddress)
	return u
}

// ClearAddress clears the value of the "address" field.
func (u *PromotionPersonUpsert) ClearAddress() *PromotionPersonUpsert {
	u.SetNull(promotionperson.FieldAddress)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.PromotionPerson.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *PromotionPersonUpsertOne) UpdateNewValues() *PromotionPersonUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(promotionperson.FieldCreatedAt)
		}
		if _, exists := u.create.mutation.Creator(); exists {
			s.SetIgnore(promotionperson.FieldCreator)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.PromotionPerson.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *PromotionPersonUpsertOne) Ignore() *PromotionPersonUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PromotionPersonUpsertOne) DoNothing() *PromotionPersonUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PromotionPersonCreate.OnConflict
// documentation for more info.
func (u *PromotionPersonUpsertOne) Update(set func(*PromotionPersonUpsert)) *PromotionPersonUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PromotionPersonUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *PromotionPersonUpsertOne) SetUpdatedAt(v time.Time) *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *PromotionPersonUpsertOne) UpdateUpdatedAt() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *PromotionPersonUpsertOne) SetDeletedAt(v time.Time) *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *PromotionPersonUpsertOne) UpdateDeletedAt() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *PromotionPersonUpsertOne) ClearDeletedAt() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *PromotionPersonUpsertOne) SetLastModifier(v *model.Modifier) *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *PromotionPersonUpsertOne) UpdateLastModifier() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *PromotionPersonUpsertOne) ClearLastModifier() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *PromotionPersonUpsertOne) SetRemark(v string) *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *PromotionPersonUpsertOne) UpdateRemark() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *PromotionPersonUpsertOne) ClearRemark() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearRemark()
	})
}

// SetStatus sets the "status" field.
func (u *PromotionPersonUpsertOne) SetStatus(v uint8) *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetStatus(v)
	})
}

// AddStatus adds v to the "status" field.
func (u *PromotionPersonUpsertOne) AddStatus(v uint8) *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.AddStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *PromotionPersonUpsertOne) UpdateStatus() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateStatus()
	})
}

// SetName sets the "name" field.
func (u *PromotionPersonUpsertOne) SetName(v string) *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *PromotionPersonUpsertOne) UpdateName() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateName()
	})
}

// ClearName clears the value of the "name" field.
func (u *PromotionPersonUpsertOne) ClearName() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearName()
	})
}

// SetIDCardNumber sets the "id_card_number" field.
func (u *PromotionPersonUpsertOne) SetIDCardNumber(v string) *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetIDCardNumber(v)
	})
}

// UpdateIDCardNumber sets the "id_card_number" field to the value that was provided on create.
func (u *PromotionPersonUpsertOne) UpdateIDCardNumber() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateIDCardNumber()
	})
}

// ClearIDCardNumber clears the value of the "id_card_number" field.
func (u *PromotionPersonUpsertOne) ClearIDCardNumber() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearIDCardNumber()
	})
}

// SetAddress sets the "address" field.
func (u *PromotionPersonUpsertOne) SetAddress(v string) *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetAddress(v)
	})
}

// UpdateAddress sets the "address" field to the value that was provided on create.
func (u *PromotionPersonUpsertOne) UpdateAddress() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateAddress()
	})
}

// ClearAddress clears the value of the "address" field.
func (u *PromotionPersonUpsertOne) ClearAddress() *PromotionPersonUpsertOne {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearAddress()
	})
}

// Exec executes the query.
func (u *PromotionPersonUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for PromotionPersonCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PromotionPersonUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *PromotionPersonUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *PromotionPersonUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// PromotionPersonCreateBulk is the builder for creating many PromotionPerson entities in bulk.
type PromotionPersonCreateBulk struct {
	config
	builders []*PromotionPersonCreate
	conflict []sql.ConflictOption
}

// Save creates the PromotionPerson entities in the database.
func (ppcb *PromotionPersonCreateBulk) Save(ctx context.Context) ([]*PromotionPerson, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ppcb.builders))
	nodes := make([]*PromotionPerson, len(ppcb.builders))
	mutators := make([]Mutator, len(ppcb.builders))
	for i := range ppcb.builders {
		func(i int, root context.Context) {
			builder := ppcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PromotionPersonMutation)
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
					_, err = mutators[i+1].Mutate(root, ppcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ppcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ppcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ppcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ppcb *PromotionPersonCreateBulk) SaveX(ctx context.Context) []*PromotionPerson {
	v, err := ppcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ppcb *PromotionPersonCreateBulk) Exec(ctx context.Context) error {
	_, err := ppcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ppcb *PromotionPersonCreateBulk) ExecX(ctx context.Context) {
	if err := ppcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.PromotionPerson.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PromotionPersonUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (ppcb *PromotionPersonCreateBulk) OnConflict(opts ...sql.ConflictOption) *PromotionPersonUpsertBulk {
	ppcb.conflict = opts
	return &PromotionPersonUpsertBulk{
		create: ppcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.PromotionPerson.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ppcb *PromotionPersonCreateBulk) OnConflictColumns(columns ...string) *PromotionPersonUpsertBulk {
	ppcb.conflict = append(ppcb.conflict, sql.ConflictColumns(columns...))
	return &PromotionPersonUpsertBulk{
		create: ppcb,
	}
}

// PromotionPersonUpsertBulk is the builder for "upsert"-ing
// a bulk of PromotionPerson nodes.
type PromotionPersonUpsertBulk struct {
	create *PromotionPersonCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.PromotionPerson.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *PromotionPersonUpsertBulk) UpdateNewValues() *PromotionPersonUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(promotionperson.FieldCreatedAt)
			}
			if _, exists := b.mutation.Creator(); exists {
				s.SetIgnore(promotionperson.FieldCreator)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.PromotionPerson.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *PromotionPersonUpsertBulk) Ignore() *PromotionPersonUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PromotionPersonUpsertBulk) DoNothing() *PromotionPersonUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PromotionPersonCreateBulk.OnConflict
// documentation for more info.
func (u *PromotionPersonUpsertBulk) Update(set func(*PromotionPersonUpsert)) *PromotionPersonUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PromotionPersonUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *PromotionPersonUpsertBulk) SetUpdatedAt(v time.Time) *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *PromotionPersonUpsertBulk) UpdateUpdatedAt() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *PromotionPersonUpsertBulk) SetDeletedAt(v time.Time) *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *PromotionPersonUpsertBulk) UpdateDeletedAt() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *PromotionPersonUpsertBulk) ClearDeletedAt() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *PromotionPersonUpsertBulk) SetLastModifier(v *model.Modifier) *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *PromotionPersonUpsertBulk) UpdateLastModifier() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *PromotionPersonUpsertBulk) ClearLastModifier() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *PromotionPersonUpsertBulk) SetRemark(v string) *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *PromotionPersonUpsertBulk) UpdateRemark() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *PromotionPersonUpsertBulk) ClearRemark() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearRemark()
	})
}

// SetStatus sets the "status" field.
func (u *PromotionPersonUpsertBulk) SetStatus(v uint8) *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetStatus(v)
	})
}

// AddStatus adds v to the "status" field.
func (u *PromotionPersonUpsertBulk) AddStatus(v uint8) *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.AddStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *PromotionPersonUpsertBulk) UpdateStatus() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateStatus()
	})
}

// SetName sets the "name" field.
func (u *PromotionPersonUpsertBulk) SetName(v string) *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *PromotionPersonUpsertBulk) UpdateName() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateName()
	})
}

// ClearName clears the value of the "name" field.
func (u *PromotionPersonUpsertBulk) ClearName() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearName()
	})
}

// SetIDCardNumber sets the "id_card_number" field.
func (u *PromotionPersonUpsertBulk) SetIDCardNumber(v string) *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetIDCardNumber(v)
	})
}

// UpdateIDCardNumber sets the "id_card_number" field to the value that was provided on create.
func (u *PromotionPersonUpsertBulk) UpdateIDCardNumber() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateIDCardNumber()
	})
}

// ClearIDCardNumber clears the value of the "id_card_number" field.
func (u *PromotionPersonUpsertBulk) ClearIDCardNumber() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearIDCardNumber()
	})
}

// SetAddress sets the "address" field.
func (u *PromotionPersonUpsertBulk) SetAddress(v string) *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.SetAddress(v)
	})
}

// UpdateAddress sets the "address" field to the value that was provided on create.
func (u *PromotionPersonUpsertBulk) UpdateAddress() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.UpdateAddress()
	})
}

// ClearAddress clears the value of the "address" field.
func (u *PromotionPersonUpsertBulk) ClearAddress() *PromotionPersonUpsertBulk {
	return u.Update(func(s *PromotionPersonUpsert) {
		s.ClearAddress()
	})
}

// Exec executes the query.
func (u *PromotionPersonUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the PromotionPersonCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for PromotionPersonCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PromotionPersonUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}