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
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/exception"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/store"
)

// ExceptionUpdate is the builder for updating Exception entities.
type ExceptionUpdate struct {
	config
	hooks    []Hook
	mutation *ExceptionMutation
}

// Where appends a list predicates to the ExceptionUpdate builder.
func (eu *ExceptionUpdate) Where(ps ...predicate.Exception) *ExceptionUpdate {
	eu.mutation.Where(ps...)
	return eu
}

// SetUpdatedAt sets the "updated_at" field.
func (eu *ExceptionUpdate) SetUpdatedAt(t time.Time) *ExceptionUpdate {
	eu.mutation.SetUpdatedAt(t)
	return eu
}

// SetDeletedAt sets the "deleted_at" field.
func (eu *ExceptionUpdate) SetDeletedAt(t time.Time) *ExceptionUpdate {
	eu.mutation.SetDeletedAt(t)
	return eu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (eu *ExceptionUpdate) SetNillableDeletedAt(t *time.Time) *ExceptionUpdate {
	if t != nil {
		eu.SetDeletedAt(*t)
	}
	return eu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (eu *ExceptionUpdate) ClearDeletedAt() *ExceptionUpdate {
	eu.mutation.ClearDeletedAt()
	return eu
}

// SetLastModifier sets the "last_modifier" field.
func (eu *ExceptionUpdate) SetLastModifier(m *model.Modifier) *ExceptionUpdate {
	eu.mutation.SetLastModifier(m)
	return eu
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (eu *ExceptionUpdate) ClearLastModifier() *ExceptionUpdate {
	eu.mutation.ClearLastModifier()
	return eu
}

// SetRemark sets the "remark" field.
func (eu *ExceptionUpdate) SetRemark(s string) *ExceptionUpdate {
	eu.mutation.SetRemark(s)
	return eu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (eu *ExceptionUpdate) SetNillableRemark(s *string) *ExceptionUpdate {
	if s != nil {
		eu.SetRemark(*s)
	}
	return eu
}

// ClearRemark clears the value of the "remark" field.
func (eu *ExceptionUpdate) ClearRemark() *ExceptionUpdate {
	eu.mutation.ClearRemark()
	return eu
}

// SetCityID sets the "city_id" field.
func (eu *ExceptionUpdate) SetCityID(u uint64) *ExceptionUpdate {
	eu.mutation.SetCityID(u)
	return eu
}

// SetEmployeeID sets the "employee_id" field.
func (eu *ExceptionUpdate) SetEmployeeID(u uint64) *ExceptionUpdate {
	eu.mutation.SetEmployeeID(u)
	return eu
}

// SetStatus sets the "status" field.
func (eu *ExceptionUpdate) SetStatus(u uint8) *ExceptionUpdate {
	eu.mutation.ResetStatus()
	eu.mutation.SetStatus(u)
	return eu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (eu *ExceptionUpdate) SetNillableStatus(u *uint8) *ExceptionUpdate {
	if u != nil {
		eu.SetStatus(*u)
	}
	return eu
}

// AddStatus adds u to the "status" field.
func (eu *ExceptionUpdate) AddStatus(u int8) *ExceptionUpdate {
	eu.mutation.AddStatus(u)
	return eu
}

// SetStoreID sets the "store_id" field.
func (eu *ExceptionUpdate) SetStoreID(u uint64) *ExceptionUpdate {
	eu.mutation.SetStoreID(u)
	return eu
}

// SetName sets the "name" field.
func (eu *ExceptionUpdate) SetName(s string) *ExceptionUpdate {
	eu.mutation.SetName(s)
	return eu
}

// SetVoltage sets the "voltage" field.
func (eu *ExceptionUpdate) SetVoltage(f float64) *ExceptionUpdate {
	eu.mutation.ResetVoltage()
	eu.mutation.SetVoltage(f)
	return eu
}

// SetNillableVoltage sets the "voltage" field if the given value is not nil.
func (eu *ExceptionUpdate) SetNillableVoltage(f *float64) *ExceptionUpdate {
	if f != nil {
		eu.SetVoltage(*f)
	}
	return eu
}

// AddVoltage adds f to the "voltage" field.
func (eu *ExceptionUpdate) AddVoltage(f float64) *ExceptionUpdate {
	eu.mutation.AddVoltage(f)
	return eu
}

// ClearVoltage clears the value of the "voltage" field.
func (eu *ExceptionUpdate) ClearVoltage() *ExceptionUpdate {
	eu.mutation.ClearVoltage()
	return eu
}

// SetReason sets the "reason" field.
func (eu *ExceptionUpdate) SetReason(s string) *ExceptionUpdate {
	eu.mutation.SetReason(s)
	return eu
}

// SetDescription sets the "description" field.
func (eu *ExceptionUpdate) SetDescription(s string) *ExceptionUpdate {
	eu.mutation.SetDescription(s)
	return eu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (eu *ExceptionUpdate) SetNillableDescription(s *string) *ExceptionUpdate {
	if s != nil {
		eu.SetDescription(*s)
	}
	return eu
}

// ClearDescription clears the value of the "description" field.
func (eu *ExceptionUpdate) ClearDescription() *ExceptionUpdate {
	eu.mutation.ClearDescription()
	return eu
}

// SetAttachments sets the "attachments" field.
func (eu *ExceptionUpdate) SetAttachments(s []string) *ExceptionUpdate {
	eu.mutation.SetAttachments(s)
	return eu
}

// ClearAttachments clears the value of the "attachments" field.
func (eu *ExceptionUpdate) ClearAttachments() *ExceptionUpdate {
	eu.mutation.ClearAttachments()
	return eu
}

// SetCity sets the "city" edge to the City entity.
func (eu *ExceptionUpdate) SetCity(c *City) *ExceptionUpdate {
	return eu.SetCityID(c.ID)
}

// SetEmployee sets the "employee" edge to the Employee entity.
func (eu *ExceptionUpdate) SetEmployee(e *Employee) *ExceptionUpdate {
	return eu.SetEmployeeID(e.ID)
}

// SetStore sets the "store" edge to the Store entity.
func (eu *ExceptionUpdate) SetStore(s *Store) *ExceptionUpdate {
	return eu.SetStoreID(s.ID)
}

// Mutation returns the ExceptionMutation object of the builder.
func (eu *ExceptionUpdate) Mutation() *ExceptionMutation {
	return eu.mutation
}

// ClearCity clears the "city" edge to the City entity.
func (eu *ExceptionUpdate) ClearCity() *ExceptionUpdate {
	eu.mutation.ClearCity()
	return eu
}

// ClearEmployee clears the "employee" edge to the Employee entity.
func (eu *ExceptionUpdate) ClearEmployee() *ExceptionUpdate {
	eu.mutation.ClearEmployee()
	return eu
}

// ClearStore clears the "store" edge to the Store entity.
func (eu *ExceptionUpdate) ClearStore() *ExceptionUpdate {
	eu.mutation.ClearStore()
	return eu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eu *ExceptionUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if err := eu.defaults(); err != nil {
		return 0, err
	}
	if len(eu.hooks) == 0 {
		if err = eu.check(); err != nil {
			return 0, err
		}
		affected, err = eu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ExceptionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = eu.check(); err != nil {
				return 0, err
			}
			eu.mutation = mutation
			affected, err = eu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(eu.hooks) - 1; i >= 0; i-- {
			if eu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = eu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, eu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (eu *ExceptionUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *ExceptionUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *ExceptionUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (eu *ExceptionUpdate) defaults() error {
	if _, ok := eu.mutation.UpdatedAt(); !ok {
		if exception.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized exception.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := exception.UpdateDefaultUpdatedAt()
		eu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (eu *ExceptionUpdate) check() error {
	if _, ok := eu.mutation.CityID(); eu.mutation.CityCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Exception.city"`)
	}
	if _, ok := eu.mutation.EmployeeID(); eu.mutation.EmployeeCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Exception.employee"`)
	}
	if _, ok := eu.mutation.StoreID(); eu.mutation.StoreCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Exception.store"`)
	}
	return nil
}

func (eu *ExceptionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   exception.Table,
			Columns: exception.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: exception.FieldID,
			},
		},
	}
	if ps := eu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: exception.FieldUpdatedAt,
		})
	}
	if value, ok := eu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: exception.FieldDeletedAt,
		})
	}
	if eu.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: exception.FieldDeletedAt,
		})
	}
	if eu.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: exception.FieldCreator,
		})
	}
	if value, ok := eu.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: exception.FieldLastModifier,
		})
	}
	if eu.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: exception.FieldLastModifier,
		})
	}
	if value, ok := eu.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: exception.FieldRemark,
		})
	}
	if eu.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: exception.FieldRemark,
		})
	}
	if value, ok := eu.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: exception.FieldStatus,
		})
	}
	if value, ok := eu.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: exception.FieldStatus,
		})
	}
	if value, ok := eu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: exception.FieldName,
		})
	}
	if value, ok := eu.mutation.Voltage(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: exception.FieldVoltage,
		})
	}
	if value, ok := eu.mutation.AddedVoltage(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: exception.FieldVoltage,
		})
	}
	if eu.mutation.VoltageCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: exception.FieldVoltage,
		})
	}
	if value, ok := eu.mutation.Reason(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: exception.FieldReason,
		})
	}
	if value, ok := eu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: exception.FieldDescription,
		})
	}
	if eu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: exception.FieldDescription,
		})
	}
	if value, ok := eu.mutation.Attachments(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: exception.FieldAttachments,
		})
	}
	if eu.mutation.AttachmentsCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: exception.FieldAttachments,
		})
	}
	if eu.mutation.CityCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   exception.CityTable,
			Columns: []string{exception.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: city.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.CityIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   exception.CityTable,
			Columns: []string{exception.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: city.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.mutation.EmployeeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   exception.EmployeeTable,
			Columns: []string{exception.EmployeeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: employee.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.EmployeeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   exception.EmployeeTable,
			Columns: []string{exception.EmployeeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: employee.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.mutation.StoreCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   exception.StoreTable,
			Columns: []string{exception.StoreColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: store.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.StoreIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   exception.StoreTable,
			Columns: []string{exception.StoreColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: store.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{exception.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// ExceptionUpdateOne is the builder for updating a single Exception entity.
type ExceptionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ExceptionMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (euo *ExceptionUpdateOne) SetUpdatedAt(t time.Time) *ExceptionUpdateOne {
	euo.mutation.SetUpdatedAt(t)
	return euo
}

// SetDeletedAt sets the "deleted_at" field.
func (euo *ExceptionUpdateOne) SetDeletedAt(t time.Time) *ExceptionUpdateOne {
	euo.mutation.SetDeletedAt(t)
	return euo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (euo *ExceptionUpdateOne) SetNillableDeletedAt(t *time.Time) *ExceptionUpdateOne {
	if t != nil {
		euo.SetDeletedAt(*t)
	}
	return euo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (euo *ExceptionUpdateOne) ClearDeletedAt() *ExceptionUpdateOne {
	euo.mutation.ClearDeletedAt()
	return euo
}

// SetLastModifier sets the "last_modifier" field.
func (euo *ExceptionUpdateOne) SetLastModifier(m *model.Modifier) *ExceptionUpdateOne {
	euo.mutation.SetLastModifier(m)
	return euo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (euo *ExceptionUpdateOne) ClearLastModifier() *ExceptionUpdateOne {
	euo.mutation.ClearLastModifier()
	return euo
}

// SetRemark sets the "remark" field.
func (euo *ExceptionUpdateOne) SetRemark(s string) *ExceptionUpdateOne {
	euo.mutation.SetRemark(s)
	return euo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (euo *ExceptionUpdateOne) SetNillableRemark(s *string) *ExceptionUpdateOne {
	if s != nil {
		euo.SetRemark(*s)
	}
	return euo
}

// ClearRemark clears the value of the "remark" field.
func (euo *ExceptionUpdateOne) ClearRemark() *ExceptionUpdateOne {
	euo.mutation.ClearRemark()
	return euo
}

// SetCityID sets the "city_id" field.
func (euo *ExceptionUpdateOne) SetCityID(u uint64) *ExceptionUpdateOne {
	euo.mutation.SetCityID(u)
	return euo
}

// SetEmployeeID sets the "employee_id" field.
func (euo *ExceptionUpdateOne) SetEmployeeID(u uint64) *ExceptionUpdateOne {
	euo.mutation.SetEmployeeID(u)
	return euo
}

// SetStatus sets the "status" field.
func (euo *ExceptionUpdateOne) SetStatus(u uint8) *ExceptionUpdateOne {
	euo.mutation.ResetStatus()
	euo.mutation.SetStatus(u)
	return euo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (euo *ExceptionUpdateOne) SetNillableStatus(u *uint8) *ExceptionUpdateOne {
	if u != nil {
		euo.SetStatus(*u)
	}
	return euo
}

// AddStatus adds u to the "status" field.
func (euo *ExceptionUpdateOne) AddStatus(u int8) *ExceptionUpdateOne {
	euo.mutation.AddStatus(u)
	return euo
}

// SetStoreID sets the "store_id" field.
func (euo *ExceptionUpdateOne) SetStoreID(u uint64) *ExceptionUpdateOne {
	euo.mutation.SetStoreID(u)
	return euo
}

// SetName sets the "name" field.
func (euo *ExceptionUpdateOne) SetName(s string) *ExceptionUpdateOne {
	euo.mutation.SetName(s)
	return euo
}

// SetVoltage sets the "voltage" field.
func (euo *ExceptionUpdateOne) SetVoltage(f float64) *ExceptionUpdateOne {
	euo.mutation.ResetVoltage()
	euo.mutation.SetVoltage(f)
	return euo
}

// SetNillableVoltage sets the "voltage" field if the given value is not nil.
func (euo *ExceptionUpdateOne) SetNillableVoltage(f *float64) *ExceptionUpdateOne {
	if f != nil {
		euo.SetVoltage(*f)
	}
	return euo
}

// AddVoltage adds f to the "voltage" field.
func (euo *ExceptionUpdateOne) AddVoltage(f float64) *ExceptionUpdateOne {
	euo.mutation.AddVoltage(f)
	return euo
}

// ClearVoltage clears the value of the "voltage" field.
func (euo *ExceptionUpdateOne) ClearVoltage() *ExceptionUpdateOne {
	euo.mutation.ClearVoltage()
	return euo
}

// SetReason sets the "reason" field.
func (euo *ExceptionUpdateOne) SetReason(s string) *ExceptionUpdateOne {
	euo.mutation.SetReason(s)
	return euo
}

// SetDescription sets the "description" field.
func (euo *ExceptionUpdateOne) SetDescription(s string) *ExceptionUpdateOne {
	euo.mutation.SetDescription(s)
	return euo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (euo *ExceptionUpdateOne) SetNillableDescription(s *string) *ExceptionUpdateOne {
	if s != nil {
		euo.SetDescription(*s)
	}
	return euo
}

// ClearDescription clears the value of the "description" field.
func (euo *ExceptionUpdateOne) ClearDescription() *ExceptionUpdateOne {
	euo.mutation.ClearDescription()
	return euo
}

// SetAttachments sets the "attachments" field.
func (euo *ExceptionUpdateOne) SetAttachments(s []string) *ExceptionUpdateOne {
	euo.mutation.SetAttachments(s)
	return euo
}

// ClearAttachments clears the value of the "attachments" field.
func (euo *ExceptionUpdateOne) ClearAttachments() *ExceptionUpdateOne {
	euo.mutation.ClearAttachments()
	return euo
}

// SetCity sets the "city" edge to the City entity.
func (euo *ExceptionUpdateOne) SetCity(c *City) *ExceptionUpdateOne {
	return euo.SetCityID(c.ID)
}

// SetEmployee sets the "employee" edge to the Employee entity.
func (euo *ExceptionUpdateOne) SetEmployee(e *Employee) *ExceptionUpdateOne {
	return euo.SetEmployeeID(e.ID)
}

// SetStore sets the "store" edge to the Store entity.
func (euo *ExceptionUpdateOne) SetStore(s *Store) *ExceptionUpdateOne {
	return euo.SetStoreID(s.ID)
}

// Mutation returns the ExceptionMutation object of the builder.
func (euo *ExceptionUpdateOne) Mutation() *ExceptionMutation {
	return euo.mutation
}

// ClearCity clears the "city" edge to the City entity.
func (euo *ExceptionUpdateOne) ClearCity() *ExceptionUpdateOne {
	euo.mutation.ClearCity()
	return euo
}

// ClearEmployee clears the "employee" edge to the Employee entity.
func (euo *ExceptionUpdateOne) ClearEmployee() *ExceptionUpdateOne {
	euo.mutation.ClearEmployee()
	return euo
}

// ClearStore clears the "store" edge to the Store entity.
func (euo *ExceptionUpdateOne) ClearStore() *ExceptionUpdateOne {
	euo.mutation.ClearStore()
	return euo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (euo *ExceptionUpdateOne) Select(field string, fields ...string) *ExceptionUpdateOne {
	euo.fields = append([]string{field}, fields...)
	return euo
}

// Save executes the query and returns the updated Exception entity.
func (euo *ExceptionUpdateOne) Save(ctx context.Context) (*Exception, error) {
	var (
		err  error
		node *Exception
	)
	if err := euo.defaults(); err != nil {
		return nil, err
	}
	if len(euo.hooks) == 0 {
		if err = euo.check(); err != nil {
			return nil, err
		}
		node, err = euo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ExceptionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = euo.check(); err != nil {
				return nil, err
			}
			euo.mutation = mutation
			node, err = euo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(euo.hooks) - 1; i >= 0; i-- {
			if euo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = euo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, euo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Exception)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from ExceptionMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (euo *ExceptionUpdateOne) SaveX(ctx context.Context) *Exception {
	node, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (euo *ExceptionUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *ExceptionUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (euo *ExceptionUpdateOne) defaults() error {
	if _, ok := euo.mutation.UpdatedAt(); !ok {
		if exception.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized exception.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := exception.UpdateDefaultUpdatedAt()
		euo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (euo *ExceptionUpdateOne) check() error {
	if _, ok := euo.mutation.CityID(); euo.mutation.CityCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Exception.city"`)
	}
	if _, ok := euo.mutation.EmployeeID(); euo.mutation.EmployeeCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Exception.employee"`)
	}
	if _, ok := euo.mutation.StoreID(); euo.mutation.StoreCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Exception.store"`)
	}
	return nil
}

func (euo *ExceptionUpdateOne) sqlSave(ctx context.Context) (_node *Exception, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   exception.Table,
			Columns: exception.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: exception.FieldID,
			},
		},
	}
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Exception.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := euo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, exception.FieldID)
		for _, f := range fields {
			if !exception.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != exception.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := euo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := euo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: exception.FieldUpdatedAt,
		})
	}
	if value, ok := euo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: exception.FieldDeletedAt,
		})
	}
	if euo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: exception.FieldDeletedAt,
		})
	}
	if euo.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: exception.FieldCreator,
		})
	}
	if value, ok := euo.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: exception.FieldLastModifier,
		})
	}
	if euo.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: exception.FieldLastModifier,
		})
	}
	if value, ok := euo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: exception.FieldRemark,
		})
	}
	if euo.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: exception.FieldRemark,
		})
	}
	if value, ok := euo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: exception.FieldStatus,
		})
	}
	if value, ok := euo.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: exception.FieldStatus,
		})
	}
	if value, ok := euo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: exception.FieldName,
		})
	}
	if value, ok := euo.mutation.Voltage(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: exception.FieldVoltage,
		})
	}
	if value, ok := euo.mutation.AddedVoltage(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: exception.FieldVoltage,
		})
	}
	if euo.mutation.VoltageCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: exception.FieldVoltage,
		})
	}
	if value, ok := euo.mutation.Reason(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: exception.FieldReason,
		})
	}
	if value, ok := euo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: exception.FieldDescription,
		})
	}
	if euo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: exception.FieldDescription,
		})
	}
	if value, ok := euo.mutation.Attachments(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: exception.FieldAttachments,
		})
	}
	if euo.mutation.AttachmentsCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: exception.FieldAttachments,
		})
	}
	if euo.mutation.CityCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   exception.CityTable,
			Columns: []string{exception.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: city.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.CityIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   exception.CityTable,
			Columns: []string{exception.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: city.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.mutation.EmployeeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   exception.EmployeeTable,
			Columns: []string{exception.EmployeeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: employee.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.EmployeeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   exception.EmployeeTable,
			Columns: []string{exception.EmployeeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: employee.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.mutation.StoreCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   exception.StoreTable,
			Columns: []string{exception.StoreColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: store.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.StoreIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   exception.StoreTable,
			Columns: []string{exception.StoreColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: store.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Exception{config: euo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{exception.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}