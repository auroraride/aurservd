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
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/promotionperson"
)

// PromotionPersonUpdate is the builder for updating PromotionPerson entities.
type PromotionPersonUpdate struct {
	config
	hooks     []Hook
	mutation  *PromotionPersonMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the PromotionPersonUpdate builder.
func (ppu *PromotionPersonUpdate) Where(ps ...predicate.PromotionPerson) *PromotionPersonUpdate {
	ppu.mutation.Where(ps...)
	return ppu
}

// SetUpdatedAt sets the "updated_at" field.
func (ppu *PromotionPersonUpdate) SetUpdatedAt(t time.Time) *PromotionPersonUpdate {
	ppu.mutation.SetUpdatedAt(t)
	return ppu
}

// SetDeletedAt sets the "deleted_at" field.
func (ppu *PromotionPersonUpdate) SetDeletedAt(t time.Time) *PromotionPersonUpdate {
	ppu.mutation.SetDeletedAt(t)
	return ppu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (ppu *PromotionPersonUpdate) SetNillableDeletedAt(t *time.Time) *PromotionPersonUpdate {
	if t != nil {
		ppu.SetDeletedAt(*t)
	}
	return ppu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (ppu *PromotionPersonUpdate) ClearDeletedAt() *PromotionPersonUpdate {
	ppu.mutation.ClearDeletedAt()
	return ppu
}

// SetLastModifier sets the "last_modifier" field.
func (ppu *PromotionPersonUpdate) SetLastModifier(m *model.Modifier) *PromotionPersonUpdate {
	ppu.mutation.SetLastModifier(m)
	return ppu
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (ppu *PromotionPersonUpdate) ClearLastModifier() *PromotionPersonUpdate {
	ppu.mutation.ClearLastModifier()
	return ppu
}

// SetRemark sets the "remark" field.
func (ppu *PromotionPersonUpdate) SetRemark(s string) *PromotionPersonUpdate {
	ppu.mutation.SetRemark(s)
	return ppu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (ppu *PromotionPersonUpdate) SetNillableRemark(s *string) *PromotionPersonUpdate {
	if s != nil {
		ppu.SetRemark(*s)
	}
	return ppu
}

// ClearRemark clears the value of the "remark" field.
func (ppu *PromotionPersonUpdate) ClearRemark() *PromotionPersonUpdate {
	ppu.mutation.ClearRemark()
	return ppu
}

// SetStatus sets the "status" field.
func (ppu *PromotionPersonUpdate) SetStatus(u uint8) *PromotionPersonUpdate {
	ppu.mutation.ResetStatus()
	ppu.mutation.SetStatus(u)
	return ppu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (ppu *PromotionPersonUpdate) SetNillableStatus(u *uint8) *PromotionPersonUpdate {
	if u != nil {
		ppu.SetStatus(*u)
	}
	return ppu
}

// AddStatus adds u to the "status" field.
func (ppu *PromotionPersonUpdate) AddStatus(u int8) *PromotionPersonUpdate {
	ppu.mutation.AddStatus(u)
	return ppu
}

// SetName sets the "name" field.
func (ppu *PromotionPersonUpdate) SetName(s string) *PromotionPersonUpdate {
	ppu.mutation.SetName(s)
	return ppu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ppu *PromotionPersonUpdate) SetNillableName(s *string) *PromotionPersonUpdate {
	if s != nil {
		ppu.SetName(*s)
	}
	return ppu
}

// ClearName clears the value of the "name" field.
func (ppu *PromotionPersonUpdate) ClearName() *PromotionPersonUpdate {
	ppu.mutation.ClearName()
	return ppu
}

// SetIDCardNumber sets the "id_card_number" field.
func (ppu *PromotionPersonUpdate) SetIDCardNumber(s string) *PromotionPersonUpdate {
	ppu.mutation.SetIDCardNumber(s)
	return ppu
}

// SetNillableIDCardNumber sets the "id_card_number" field if the given value is not nil.
func (ppu *PromotionPersonUpdate) SetNillableIDCardNumber(s *string) *PromotionPersonUpdate {
	if s != nil {
		ppu.SetIDCardNumber(*s)
	}
	return ppu
}

// ClearIDCardNumber clears the value of the "id_card_number" field.
func (ppu *PromotionPersonUpdate) ClearIDCardNumber() *PromotionPersonUpdate {
	ppu.mutation.ClearIDCardNumber()
	return ppu
}

// SetAddress sets the "address" field.
func (ppu *PromotionPersonUpdate) SetAddress(s string) *PromotionPersonUpdate {
	ppu.mutation.SetAddress(s)
	return ppu
}

// SetNillableAddress sets the "address" field if the given value is not nil.
func (ppu *PromotionPersonUpdate) SetNillableAddress(s *string) *PromotionPersonUpdate {
	if s != nil {
		ppu.SetAddress(*s)
	}
	return ppu
}

// ClearAddress clears the value of the "address" field.
func (ppu *PromotionPersonUpdate) ClearAddress() *PromotionPersonUpdate {
	ppu.mutation.ClearAddress()
	return ppu
}

// AddMemberIDs adds the "member" edge to the PromotionMember entity by IDs.
func (ppu *PromotionPersonUpdate) AddMemberIDs(ids ...uint64) *PromotionPersonUpdate {
	ppu.mutation.AddMemberIDs(ids...)
	return ppu
}

// AddMember adds the "member" edges to the PromotionMember entity.
func (ppu *PromotionPersonUpdate) AddMember(p ...*PromotionMember) *PromotionPersonUpdate {
	ids := make([]uint64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ppu.AddMemberIDs(ids...)
}

// Mutation returns the PromotionPersonMutation object of the builder.
func (ppu *PromotionPersonUpdate) Mutation() *PromotionPersonMutation {
	return ppu.mutation
}

// ClearMember clears all "member" edges to the PromotionMember entity.
func (ppu *PromotionPersonUpdate) ClearMember() *PromotionPersonUpdate {
	ppu.mutation.ClearMember()
	return ppu
}

// RemoveMemberIDs removes the "member" edge to PromotionMember entities by IDs.
func (ppu *PromotionPersonUpdate) RemoveMemberIDs(ids ...uint64) *PromotionPersonUpdate {
	ppu.mutation.RemoveMemberIDs(ids...)
	return ppu
}

// RemoveMember removes "member" edges to PromotionMember entities.
func (ppu *PromotionPersonUpdate) RemoveMember(p ...*PromotionMember) *PromotionPersonUpdate {
	ids := make([]uint64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ppu.RemoveMemberIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ppu *PromotionPersonUpdate) Save(ctx context.Context) (int, error) {
	if err := ppu.defaults(); err != nil {
		return 0, err
	}
	return withHooks(ctx, ppu.sqlSave, ppu.mutation, ppu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ppu *PromotionPersonUpdate) SaveX(ctx context.Context) int {
	affected, err := ppu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ppu *PromotionPersonUpdate) Exec(ctx context.Context) error {
	_, err := ppu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ppu *PromotionPersonUpdate) ExecX(ctx context.Context) {
	if err := ppu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ppu *PromotionPersonUpdate) defaults() error {
	if _, ok := ppu.mutation.UpdatedAt(); !ok {
		if promotionperson.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized promotionperson.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := promotionperson.UpdateDefaultUpdatedAt()
		ppu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (ppu *PromotionPersonUpdate) check() error {
	if v, ok := ppu.mutation.Name(); ok {
		if err := promotionperson.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "PromotionPerson.name": %w`, err)}
		}
	}
	if v, ok := ppu.mutation.IDCardNumber(); ok {
		if err := promotionperson.IDCardNumberValidator(v); err != nil {
			return &ValidationError{Name: "id_card_number", err: fmt.Errorf(`ent: validator failed for field "PromotionPerson.id_card_number": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ppu *PromotionPersonUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PromotionPersonUpdate {
	ppu.modifiers = append(ppu.modifiers, modifiers...)
	return ppu
}

func (ppu *PromotionPersonUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ppu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(promotionperson.Table, promotionperson.Columns, sqlgraph.NewFieldSpec(promotionperson.FieldID, field.TypeUint64))
	if ps := ppu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ppu.mutation.UpdatedAt(); ok {
		_spec.SetField(promotionperson.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ppu.mutation.DeletedAt(); ok {
		_spec.SetField(promotionperson.FieldDeletedAt, field.TypeTime, value)
	}
	if ppu.mutation.DeletedAtCleared() {
		_spec.ClearField(promotionperson.FieldDeletedAt, field.TypeTime)
	}
	if ppu.mutation.CreatorCleared() {
		_spec.ClearField(promotionperson.FieldCreator, field.TypeJSON)
	}
	if value, ok := ppu.mutation.LastModifier(); ok {
		_spec.SetField(promotionperson.FieldLastModifier, field.TypeJSON, value)
	}
	if ppu.mutation.LastModifierCleared() {
		_spec.ClearField(promotionperson.FieldLastModifier, field.TypeJSON)
	}
	if value, ok := ppu.mutation.Remark(); ok {
		_spec.SetField(promotionperson.FieldRemark, field.TypeString, value)
	}
	if ppu.mutation.RemarkCleared() {
		_spec.ClearField(promotionperson.FieldRemark, field.TypeString)
	}
	if value, ok := ppu.mutation.Status(); ok {
		_spec.SetField(promotionperson.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := ppu.mutation.AddedStatus(); ok {
		_spec.AddField(promotionperson.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := ppu.mutation.Name(); ok {
		_spec.SetField(promotionperson.FieldName, field.TypeString, value)
	}
	if ppu.mutation.NameCleared() {
		_spec.ClearField(promotionperson.FieldName, field.TypeString)
	}
	if value, ok := ppu.mutation.IDCardNumber(); ok {
		_spec.SetField(promotionperson.FieldIDCardNumber, field.TypeString, value)
	}
	if ppu.mutation.IDCardNumberCleared() {
		_spec.ClearField(promotionperson.FieldIDCardNumber, field.TypeString)
	}
	if value, ok := ppu.mutation.Address(); ok {
		_spec.SetField(promotionperson.FieldAddress, field.TypeString, value)
	}
	if ppu.mutation.AddressCleared() {
		_spec.ClearField(promotionperson.FieldAddress, field.TypeString)
	}
	if ppu.mutation.MemberCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ppu.mutation.RemovedMemberIDs(); len(nodes) > 0 && !ppu.mutation.MemberCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ppu.mutation.MemberIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(ppu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, ppu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{promotionperson.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ppu.mutation.done = true
	return n, nil
}

// PromotionPersonUpdateOne is the builder for updating a single PromotionPerson entity.
type PromotionPersonUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *PromotionPersonMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (ppuo *PromotionPersonUpdateOne) SetUpdatedAt(t time.Time) *PromotionPersonUpdateOne {
	ppuo.mutation.SetUpdatedAt(t)
	return ppuo
}

// SetDeletedAt sets the "deleted_at" field.
func (ppuo *PromotionPersonUpdateOne) SetDeletedAt(t time.Time) *PromotionPersonUpdateOne {
	ppuo.mutation.SetDeletedAt(t)
	return ppuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (ppuo *PromotionPersonUpdateOne) SetNillableDeletedAt(t *time.Time) *PromotionPersonUpdateOne {
	if t != nil {
		ppuo.SetDeletedAt(*t)
	}
	return ppuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (ppuo *PromotionPersonUpdateOne) ClearDeletedAt() *PromotionPersonUpdateOne {
	ppuo.mutation.ClearDeletedAt()
	return ppuo
}

// SetLastModifier sets the "last_modifier" field.
func (ppuo *PromotionPersonUpdateOne) SetLastModifier(m *model.Modifier) *PromotionPersonUpdateOne {
	ppuo.mutation.SetLastModifier(m)
	return ppuo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (ppuo *PromotionPersonUpdateOne) ClearLastModifier() *PromotionPersonUpdateOne {
	ppuo.mutation.ClearLastModifier()
	return ppuo
}

// SetRemark sets the "remark" field.
func (ppuo *PromotionPersonUpdateOne) SetRemark(s string) *PromotionPersonUpdateOne {
	ppuo.mutation.SetRemark(s)
	return ppuo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (ppuo *PromotionPersonUpdateOne) SetNillableRemark(s *string) *PromotionPersonUpdateOne {
	if s != nil {
		ppuo.SetRemark(*s)
	}
	return ppuo
}

// ClearRemark clears the value of the "remark" field.
func (ppuo *PromotionPersonUpdateOne) ClearRemark() *PromotionPersonUpdateOne {
	ppuo.mutation.ClearRemark()
	return ppuo
}

// SetStatus sets the "status" field.
func (ppuo *PromotionPersonUpdateOne) SetStatus(u uint8) *PromotionPersonUpdateOne {
	ppuo.mutation.ResetStatus()
	ppuo.mutation.SetStatus(u)
	return ppuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (ppuo *PromotionPersonUpdateOne) SetNillableStatus(u *uint8) *PromotionPersonUpdateOne {
	if u != nil {
		ppuo.SetStatus(*u)
	}
	return ppuo
}

// AddStatus adds u to the "status" field.
func (ppuo *PromotionPersonUpdateOne) AddStatus(u int8) *PromotionPersonUpdateOne {
	ppuo.mutation.AddStatus(u)
	return ppuo
}

// SetName sets the "name" field.
func (ppuo *PromotionPersonUpdateOne) SetName(s string) *PromotionPersonUpdateOne {
	ppuo.mutation.SetName(s)
	return ppuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ppuo *PromotionPersonUpdateOne) SetNillableName(s *string) *PromotionPersonUpdateOne {
	if s != nil {
		ppuo.SetName(*s)
	}
	return ppuo
}

// ClearName clears the value of the "name" field.
func (ppuo *PromotionPersonUpdateOne) ClearName() *PromotionPersonUpdateOne {
	ppuo.mutation.ClearName()
	return ppuo
}

// SetIDCardNumber sets the "id_card_number" field.
func (ppuo *PromotionPersonUpdateOne) SetIDCardNumber(s string) *PromotionPersonUpdateOne {
	ppuo.mutation.SetIDCardNumber(s)
	return ppuo
}

// SetNillableIDCardNumber sets the "id_card_number" field if the given value is not nil.
func (ppuo *PromotionPersonUpdateOne) SetNillableIDCardNumber(s *string) *PromotionPersonUpdateOne {
	if s != nil {
		ppuo.SetIDCardNumber(*s)
	}
	return ppuo
}

// ClearIDCardNumber clears the value of the "id_card_number" field.
func (ppuo *PromotionPersonUpdateOne) ClearIDCardNumber() *PromotionPersonUpdateOne {
	ppuo.mutation.ClearIDCardNumber()
	return ppuo
}

// SetAddress sets the "address" field.
func (ppuo *PromotionPersonUpdateOne) SetAddress(s string) *PromotionPersonUpdateOne {
	ppuo.mutation.SetAddress(s)
	return ppuo
}

// SetNillableAddress sets the "address" field if the given value is not nil.
func (ppuo *PromotionPersonUpdateOne) SetNillableAddress(s *string) *PromotionPersonUpdateOne {
	if s != nil {
		ppuo.SetAddress(*s)
	}
	return ppuo
}

// ClearAddress clears the value of the "address" field.
func (ppuo *PromotionPersonUpdateOne) ClearAddress() *PromotionPersonUpdateOne {
	ppuo.mutation.ClearAddress()
	return ppuo
}

// AddMemberIDs adds the "member" edge to the PromotionMember entity by IDs.
func (ppuo *PromotionPersonUpdateOne) AddMemberIDs(ids ...uint64) *PromotionPersonUpdateOne {
	ppuo.mutation.AddMemberIDs(ids...)
	return ppuo
}

// AddMember adds the "member" edges to the PromotionMember entity.
func (ppuo *PromotionPersonUpdateOne) AddMember(p ...*PromotionMember) *PromotionPersonUpdateOne {
	ids := make([]uint64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ppuo.AddMemberIDs(ids...)
}

// Mutation returns the PromotionPersonMutation object of the builder.
func (ppuo *PromotionPersonUpdateOne) Mutation() *PromotionPersonMutation {
	return ppuo.mutation
}

// ClearMember clears all "member" edges to the PromotionMember entity.
func (ppuo *PromotionPersonUpdateOne) ClearMember() *PromotionPersonUpdateOne {
	ppuo.mutation.ClearMember()
	return ppuo
}

// RemoveMemberIDs removes the "member" edge to PromotionMember entities by IDs.
func (ppuo *PromotionPersonUpdateOne) RemoveMemberIDs(ids ...uint64) *PromotionPersonUpdateOne {
	ppuo.mutation.RemoveMemberIDs(ids...)
	return ppuo
}

// RemoveMember removes "member" edges to PromotionMember entities.
func (ppuo *PromotionPersonUpdateOne) RemoveMember(p ...*PromotionMember) *PromotionPersonUpdateOne {
	ids := make([]uint64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ppuo.RemoveMemberIDs(ids...)
}

// Where appends a list predicates to the PromotionPersonUpdate builder.
func (ppuo *PromotionPersonUpdateOne) Where(ps ...predicate.PromotionPerson) *PromotionPersonUpdateOne {
	ppuo.mutation.Where(ps...)
	return ppuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ppuo *PromotionPersonUpdateOne) Select(field string, fields ...string) *PromotionPersonUpdateOne {
	ppuo.fields = append([]string{field}, fields...)
	return ppuo
}

// Save executes the query and returns the updated PromotionPerson entity.
func (ppuo *PromotionPersonUpdateOne) Save(ctx context.Context) (*PromotionPerson, error) {
	if err := ppuo.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, ppuo.sqlSave, ppuo.mutation, ppuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ppuo *PromotionPersonUpdateOne) SaveX(ctx context.Context) *PromotionPerson {
	node, err := ppuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ppuo *PromotionPersonUpdateOne) Exec(ctx context.Context) error {
	_, err := ppuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ppuo *PromotionPersonUpdateOne) ExecX(ctx context.Context) {
	if err := ppuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ppuo *PromotionPersonUpdateOne) defaults() error {
	if _, ok := ppuo.mutation.UpdatedAt(); !ok {
		if promotionperson.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized promotionperson.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := promotionperson.UpdateDefaultUpdatedAt()
		ppuo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (ppuo *PromotionPersonUpdateOne) check() error {
	if v, ok := ppuo.mutation.Name(); ok {
		if err := promotionperson.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "PromotionPerson.name": %w`, err)}
		}
	}
	if v, ok := ppuo.mutation.IDCardNumber(); ok {
		if err := promotionperson.IDCardNumberValidator(v); err != nil {
			return &ValidationError{Name: "id_card_number", err: fmt.Errorf(`ent: validator failed for field "PromotionPerson.id_card_number": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ppuo *PromotionPersonUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PromotionPersonUpdateOne {
	ppuo.modifiers = append(ppuo.modifiers, modifiers...)
	return ppuo
}

func (ppuo *PromotionPersonUpdateOne) sqlSave(ctx context.Context) (_node *PromotionPerson, err error) {
	if err := ppuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(promotionperson.Table, promotionperson.Columns, sqlgraph.NewFieldSpec(promotionperson.FieldID, field.TypeUint64))
	id, ok := ppuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "PromotionPerson.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ppuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, promotionperson.FieldID)
		for _, f := range fields {
			if !promotionperson.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != promotionperson.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ppuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ppuo.mutation.UpdatedAt(); ok {
		_spec.SetField(promotionperson.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ppuo.mutation.DeletedAt(); ok {
		_spec.SetField(promotionperson.FieldDeletedAt, field.TypeTime, value)
	}
	if ppuo.mutation.DeletedAtCleared() {
		_spec.ClearField(promotionperson.FieldDeletedAt, field.TypeTime)
	}
	if ppuo.mutation.CreatorCleared() {
		_spec.ClearField(promotionperson.FieldCreator, field.TypeJSON)
	}
	if value, ok := ppuo.mutation.LastModifier(); ok {
		_spec.SetField(promotionperson.FieldLastModifier, field.TypeJSON, value)
	}
	if ppuo.mutation.LastModifierCleared() {
		_spec.ClearField(promotionperson.FieldLastModifier, field.TypeJSON)
	}
	if value, ok := ppuo.mutation.Remark(); ok {
		_spec.SetField(promotionperson.FieldRemark, field.TypeString, value)
	}
	if ppuo.mutation.RemarkCleared() {
		_spec.ClearField(promotionperson.FieldRemark, field.TypeString)
	}
	if value, ok := ppuo.mutation.Status(); ok {
		_spec.SetField(promotionperson.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := ppuo.mutation.AddedStatus(); ok {
		_spec.AddField(promotionperson.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := ppuo.mutation.Name(); ok {
		_spec.SetField(promotionperson.FieldName, field.TypeString, value)
	}
	if ppuo.mutation.NameCleared() {
		_spec.ClearField(promotionperson.FieldName, field.TypeString)
	}
	if value, ok := ppuo.mutation.IDCardNumber(); ok {
		_spec.SetField(promotionperson.FieldIDCardNumber, field.TypeString, value)
	}
	if ppuo.mutation.IDCardNumberCleared() {
		_spec.ClearField(promotionperson.FieldIDCardNumber, field.TypeString)
	}
	if value, ok := ppuo.mutation.Address(); ok {
		_spec.SetField(promotionperson.FieldAddress, field.TypeString, value)
	}
	if ppuo.mutation.AddressCleared() {
		_spec.ClearField(promotionperson.FieldAddress, field.TypeString)
	}
	if ppuo.mutation.MemberCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ppuo.mutation.RemovedMemberIDs(); len(nodes) > 0 && !ppuo.mutation.MemberCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ppuo.mutation.MemberIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(ppuo.modifiers...)
	_node = &PromotionPerson{config: ppuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ppuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{promotionperson.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ppuo.mutation.done = true
	return _node, nil
}