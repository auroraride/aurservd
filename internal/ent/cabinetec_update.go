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
	"github.com/auroraride/aurservd/internal/ent/cabinetec"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// CabinetEcUpdate is the builder for updating CabinetEc entities.
type CabinetEcUpdate struct {
	config
	hooks     []Hook
	mutation  *CabinetEcMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the CabinetEcUpdate builder.
func (ceu *CabinetEcUpdate) Where(ps ...predicate.CabinetEc) *CabinetEcUpdate {
	ceu.mutation.Where(ps...)
	return ceu
}

// SetUpdatedAt sets the "updated_at" field.
func (ceu *CabinetEcUpdate) SetUpdatedAt(t time.Time) *CabinetEcUpdate {
	ceu.mutation.SetUpdatedAt(t)
	return ceu
}

// SetDeletedAt sets the "deleted_at" field.
func (ceu *CabinetEcUpdate) SetDeletedAt(t time.Time) *CabinetEcUpdate {
	ceu.mutation.SetDeletedAt(t)
	return ceu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (ceu *CabinetEcUpdate) SetNillableDeletedAt(t *time.Time) *CabinetEcUpdate {
	if t != nil {
		ceu.SetDeletedAt(*t)
	}
	return ceu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (ceu *CabinetEcUpdate) ClearDeletedAt() *CabinetEcUpdate {
	ceu.mutation.ClearDeletedAt()
	return ceu
}

// SetSerial sets the "serial" field.
func (ceu *CabinetEcUpdate) SetSerial(s string) *CabinetEcUpdate {
	ceu.mutation.SetSerial(s)
	return ceu
}

// SetNillableSerial sets the "serial" field if the given value is not nil.
func (ceu *CabinetEcUpdate) SetNillableSerial(s *string) *CabinetEcUpdate {
	if s != nil {
		ceu.SetSerial(*s)
	}
	return ceu
}

// SetDate sets the "date" field.
func (ceu *CabinetEcUpdate) SetDate(t time.Time) *CabinetEcUpdate {
	ceu.mutation.SetDate(t)
	return ceu
}

// SetNillableDate sets the "date" field if the given value is not nil.
func (ceu *CabinetEcUpdate) SetNillableDate(t *time.Time) *CabinetEcUpdate {
	if t != nil {
		ceu.SetDate(*t)
	}
	return ceu
}

// SetStart sets the "start" field.
func (ceu *CabinetEcUpdate) SetStart(f float64) *CabinetEcUpdate {
	ceu.mutation.ResetStart()
	ceu.mutation.SetStart(f)
	return ceu
}

// SetNillableStart sets the "start" field if the given value is not nil.
func (ceu *CabinetEcUpdate) SetNillableStart(f *float64) *CabinetEcUpdate {
	if f != nil {
		ceu.SetStart(*f)
	}
	return ceu
}

// AddStart adds f to the "start" field.
func (ceu *CabinetEcUpdate) AddStart(f float64) *CabinetEcUpdate {
	ceu.mutation.AddStart(f)
	return ceu
}

// SetEnd sets the "end" field.
func (ceu *CabinetEcUpdate) SetEnd(f float64) *CabinetEcUpdate {
	ceu.mutation.ResetEnd()
	ceu.mutation.SetEnd(f)
	return ceu
}

// SetNillableEnd sets the "end" field if the given value is not nil.
func (ceu *CabinetEcUpdate) SetNillableEnd(f *float64) *CabinetEcUpdate {
	if f != nil {
		ceu.SetEnd(*f)
	}
	return ceu
}

// AddEnd adds f to the "end" field.
func (ceu *CabinetEcUpdate) AddEnd(f float64) *CabinetEcUpdate {
	ceu.mutation.AddEnd(f)
	return ceu
}

// ClearEnd clears the value of the "end" field.
func (ceu *CabinetEcUpdate) ClearEnd() *CabinetEcUpdate {
	ceu.mutation.ClearEnd()
	return ceu
}

// SetTotal sets the "total" field.
func (ceu *CabinetEcUpdate) SetTotal(f float64) *CabinetEcUpdate {
	ceu.mutation.ResetTotal()
	ceu.mutation.SetTotal(f)
	return ceu
}

// SetNillableTotal sets the "total" field if the given value is not nil.
func (ceu *CabinetEcUpdate) SetNillableTotal(f *float64) *CabinetEcUpdate {
	if f != nil {
		ceu.SetTotal(*f)
	}
	return ceu
}

// AddTotal adds f to the "total" field.
func (ceu *CabinetEcUpdate) AddTotal(f float64) *CabinetEcUpdate {
	ceu.mutation.AddTotal(f)
	return ceu
}

// Mutation returns the CabinetEcMutation object of the builder.
func (ceu *CabinetEcUpdate) Mutation() *CabinetEcMutation {
	return ceu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ceu *CabinetEcUpdate) Save(ctx context.Context) (int, error) {
	ceu.defaults()
	return withHooks(ctx, ceu.sqlSave, ceu.mutation, ceu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ceu *CabinetEcUpdate) SaveX(ctx context.Context) int {
	affected, err := ceu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ceu *CabinetEcUpdate) Exec(ctx context.Context) error {
	_, err := ceu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ceu *CabinetEcUpdate) ExecX(ctx context.Context) {
	if err := ceu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ceu *CabinetEcUpdate) defaults() {
	if _, ok := ceu.mutation.UpdatedAt(); !ok {
		v := cabinetec.UpdateDefaultUpdatedAt()
		ceu.mutation.SetUpdatedAt(v)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ceu *CabinetEcUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *CabinetEcUpdate {
	ceu.modifiers = append(ceu.modifiers, modifiers...)
	return ceu
}

func (ceu *CabinetEcUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(cabinetec.Table, cabinetec.Columns, sqlgraph.NewFieldSpec(cabinetec.FieldID, field.TypeUint64))
	if ps := ceu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ceu.mutation.UpdatedAt(); ok {
		_spec.SetField(cabinetec.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ceu.mutation.DeletedAt(); ok {
		_spec.SetField(cabinetec.FieldDeletedAt, field.TypeTime, value)
	}
	if ceu.mutation.DeletedAtCleared() {
		_spec.ClearField(cabinetec.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := ceu.mutation.Serial(); ok {
		_spec.SetField(cabinetec.FieldSerial, field.TypeString, value)
	}
	if value, ok := ceu.mutation.Date(); ok {
		_spec.SetField(cabinetec.FieldDate, field.TypeTime, value)
	}
	if value, ok := ceu.mutation.Start(); ok {
		_spec.SetField(cabinetec.FieldStart, field.TypeFloat64, value)
	}
	if value, ok := ceu.mutation.AddedStart(); ok {
		_spec.AddField(cabinetec.FieldStart, field.TypeFloat64, value)
	}
	if value, ok := ceu.mutation.End(); ok {
		_spec.SetField(cabinetec.FieldEnd, field.TypeFloat64, value)
	}
	if value, ok := ceu.mutation.AddedEnd(); ok {
		_spec.AddField(cabinetec.FieldEnd, field.TypeFloat64, value)
	}
	if ceu.mutation.EndCleared() {
		_spec.ClearField(cabinetec.FieldEnd, field.TypeFloat64)
	}
	if value, ok := ceu.mutation.Total(); ok {
		_spec.SetField(cabinetec.FieldTotal, field.TypeFloat64, value)
	}
	if value, ok := ceu.mutation.AddedTotal(); ok {
		_spec.AddField(cabinetec.FieldTotal, field.TypeFloat64, value)
	}
	_spec.AddModifiers(ceu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, ceu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cabinetec.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ceu.mutation.done = true
	return n, nil
}

// CabinetEcUpdateOne is the builder for updating a single CabinetEc entity.
type CabinetEcUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *CabinetEcMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (ceuo *CabinetEcUpdateOne) SetUpdatedAt(t time.Time) *CabinetEcUpdateOne {
	ceuo.mutation.SetUpdatedAt(t)
	return ceuo
}

// SetDeletedAt sets the "deleted_at" field.
func (ceuo *CabinetEcUpdateOne) SetDeletedAt(t time.Time) *CabinetEcUpdateOne {
	ceuo.mutation.SetDeletedAt(t)
	return ceuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (ceuo *CabinetEcUpdateOne) SetNillableDeletedAt(t *time.Time) *CabinetEcUpdateOne {
	if t != nil {
		ceuo.SetDeletedAt(*t)
	}
	return ceuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (ceuo *CabinetEcUpdateOne) ClearDeletedAt() *CabinetEcUpdateOne {
	ceuo.mutation.ClearDeletedAt()
	return ceuo
}

// SetSerial sets the "serial" field.
func (ceuo *CabinetEcUpdateOne) SetSerial(s string) *CabinetEcUpdateOne {
	ceuo.mutation.SetSerial(s)
	return ceuo
}

// SetNillableSerial sets the "serial" field if the given value is not nil.
func (ceuo *CabinetEcUpdateOne) SetNillableSerial(s *string) *CabinetEcUpdateOne {
	if s != nil {
		ceuo.SetSerial(*s)
	}
	return ceuo
}

// SetDate sets the "date" field.
func (ceuo *CabinetEcUpdateOne) SetDate(t time.Time) *CabinetEcUpdateOne {
	ceuo.mutation.SetDate(t)
	return ceuo
}

// SetNillableDate sets the "date" field if the given value is not nil.
func (ceuo *CabinetEcUpdateOne) SetNillableDate(t *time.Time) *CabinetEcUpdateOne {
	if t != nil {
		ceuo.SetDate(*t)
	}
	return ceuo
}

// SetStart sets the "start" field.
func (ceuo *CabinetEcUpdateOne) SetStart(f float64) *CabinetEcUpdateOne {
	ceuo.mutation.ResetStart()
	ceuo.mutation.SetStart(f)
	return ceuo
}

// SetNillableStart sets the "start" field if the given value is not nil.
func (ceuo *CabinetEcUpdateOne) SetNillableStart(f *float64) *CabinetEcUpdateOne {
	if f != nil {
		ceuo.SetStart(*f)
	}
	return ceuo
}

// AddStart adds f to the "start" field.
func (ceuo *CabinetEcUpdateOne) AddStart(f float64) *CabinetEcUpdateOne {
	ceuo.mutation.AddStart(f)
	return ceuo
}

// SetEnd sets the "end" field.
func (ceuo *CabinetEcUpdateOne) SetEnd(f float64) *CabinetEcUpdateOne {
	ceuo.mutation.ResetEnd()
	ceuo.mutation.SetEnd(f)
	return ceuo
}

// SetNillableEnd sets the "end" field if the given value is not nil.
func (ceuo *CabinetEcUpdateOne) SetNillableEnd(f *float64) *CabinetEcUpdateOne {
	if f != nil {
		ceuo.SetEnd(*f)
	}
	return ceuo
}

// AddEnd adds f to the "end" field.
func (ceuo *CabinetEcUpdateOne) AddEnd(f float64) *CabinetEcUpdateOne {
	ceuo.mutation.AddEnd(f)
	return ceuo
}

// ClearEnd clears the value of the "end" field.
func (ceuo *CabinetEcUpdateOne) ClearEnd() *CabinetEcUpdateOne {
	ceuo.mutation.ClearEnd()
	return ceuo
}

// SetTotal sets the "total" field.
func (ceuo *CabinetEcUpdateOne) SetTotal(f float64) *CabinetEcUpdateOne {
	ceuo.mutation.ResetTotal()
	ceuo.mutation.SetTotal(f)
	return ceuo
}

// SetNillableTotal sets the "total" field if the given value is not nil.
func (ceuo *CabinetEcUpdateOne) SetNillableTotal(f *float64) *CabinetEcUpdateOne {
	if f != nil {
		ceuo.SetTotal(*f)
	}
	return ceuo
}

// AddTotal adds f to the "total" field.
func (ceuo *CabinetEcUpdateOne) AddTotal(f float64) *CabinetEcUpdateOne {
	ceuo.mutation.AddTotal(f)
	return ceuo
}

// Mutation returns the CabinetEcMutation object of the builder.
func (ceuo *CabinetEcUpdateOne) Mutation() *CabinetEcMutation {
	return ceuo.mutation
}

// Where appends a list predicates to the CabinetEcUpdate builder.
func (ceuo *CabinetEcUpdateOne) Where(ps ...predicate.CabinetEc) *CabinetEcUpdateOne {
	ceuo.mutation.Where(ps...)
	return ceuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ceuo *CabinetEcUpdateOne) Select(field string, fields ...string) *CabinetEcUpdateOne {
	ceuo.fields = append([]string{field}, fields...)
	return ceuo
}

// Save executes the query and returns the updated CabinetEc entity.
func (ceuo *CabinetEcUpdateOne) Save(ctx context.Context) (*CabinetEc, error) {
	ceuo.defaults()
	return withHooks(ctx, ceuo.sqlSave, ceuo.mutation, ceuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ceuo *CabinetEcUpdateOne) SaveX(ctx context.Context) *CabinetEc {
	node, err := ceuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ceuo *CabinetEcUpdateOne) Exec(ctx context.Context) error {
	_, err := ceuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ceuo *CabinetEcUpdateOne) ExecX(ctx context.Context) {
	if err := ceuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ceuo *CabinetEcUpdateOne) defaults() {
	if _, ok := ceuo.mutation.UpdatedAt(); !ok {
		v := cabinetec.UpdateDefaultUpdatedAt()
		ceuo.mutation.SetUpdatedAt(v)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ceuo *CabinetEcUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *CabinetEcUpdateOne {
	ceuo.modifiers = append(ceuo.modifiers, modifiers...)
	return ceuo
}

func (ceuo *CabinetEcUpdateOne) sqlSave(ctx context.Context) (_node *CabinetEc, err error) {
	_spec := sqlgraph.NewUpdateSpec(cabinetec.Table, cabinetec.Columns, sqlgraph.NewFieldSpec(cabinetec.FieldID, field.TypeUint64))
	id, ok := ceuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "CabinetEc.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ceuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cabinetec.FieldID)
		for _, f := range fields {
			if !cabinetec.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != cabinetec.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ceuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ceuo.mutation.UpdatedAt(); ok {
		_spec.SetField(cabinetec.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ceuo.mutation.DeletedAt(); ok {
		_spec.SetField(cabinetec.FieldDeletedAt, field.TypeTime, value)
	}
	if ceuo.mutation.DeletedAtCleared() {
		_spec.ClearField(cabinetec.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := ceuo.mutation.Serial(); ok {
		_spec.SetField(cabinetec.FieldSerial, field.TypeString, value)
	}
	if value, ok := ceuo.mutation.Date(); ok {
		_spec.SetField(cabinetec.FieldDate, field.TypeTime, value)
	}
	if value, ok := ceuo.mutation.Start(); ok {
		_spec.SetField(cabinetec.FieldStart, field.TypeFloat64, value)
	}
	if value, ok := ceuo.mutation.AddedStart(); ok {
		_spec.AddField(cabinetec.FieldStart, field.TypeFloat64, value)
	}
	if value, ok := ceuo.mutation.End(); ok {
		_spec.SetField(cabinetec.FieldEnd, field.TypeFloat64, value)
	}
	if value, ok := ceuo.mutation.AddedEnd(); ok {
		_spec.AddField(cabinetec.FieldEnd, field.TypeFloat64, value)
	}
	if ceuo.mutation.EndCleared() {
		_spec.ClearField(cabinetec.FieldEnd, field.TypeFloat64)
	}
	if value, ok := ceuo.mutation.Total(); ok {
		_spec.SetField(cabinetec.FieldTotal, field.TypeFloat64, value)
	}
	if value, ok := ceuo.mutation.AddedTotal(); ok {
		_spec.AddField(cabinetec.FieldTotal, field.TypeFloat64, value)
	}
	_spec.AddModifiers(ceuo.modifiers...)
	_node = &CabinetEc{config: ceuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ceuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cabinetec.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ceuo.mutation.done = true
	return _node, nil
}