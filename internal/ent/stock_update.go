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
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/manager"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/stock"
	"github.com/auroraride/aurservd/internal/ent/store"
)

// StockUpdate is the builder for updating Stock entities.
type StockUpdate struct {
	config
	hooks    []Hook
	mutation *StockMutation
}

// Where appends a list predicates to the StockUpdate builder.
func (su *StockUpdate) Where(ps ...predicate.Stock) *StockUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetUpdatedAt sets the "updated_at" field.
func (su *StockUpdate) SetUpdatedAt(t time.Time) *StockUpdate {
	su.mutation.SetUpdatedAt(t)
	return su
}

// SetDeletedAt sets the "deleted_at" field.
func (su *StockUpdate) SetDeletedAt(t time.Time) *StockUpdate {
	su.mutation.SetDeletedAt(t)
	return su
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (su *StockUpdate) SetNillableDeletedAt(t *time.Time) *StockUpdate {
	if t != nil {
		su.SetDeletedAt(*t)
	}
	return su
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (su *StockUpdate) ClearDeletedAt() *StockUpdate {
	su.mutation.ClearDeletedAt()
	return su
}

// SetLastModifier sets the "last_modifier" field.
func (su *StockUpdate) SetLastModifier(m *model.Modifier) *StockUpdate {
	su.mutation.SetLastModifier(m)
	return su
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (su *StockUpdate) ClearLastModifier() *StockUpdate {
	su.mutation.ClearLastModifier()
	return su
}

// SetRemark sets the "remark" field.
func (su *StockUpdate) SetRemark(s string) *StockUpdate {
	su.mutation.SetRemark(s)
	return su
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (su *StockUpdate) SetNillableRemark(s *string) *StockUpdate {
	if s != nil {
		su.SetRemark(*s)
	}
	return su
}

// ClearRemark clears the value of the "remark" field.
func (su *StockUpdate) ClearRemark() *StockUpdate {
	su.mutation.ClearRemark()
	return su
}

// SetManagerID sets the "manager_id" field.
func (su *StockUpdate) SetManagerID(u uint64) *StockUpdate {
	su.mutation.SetManagerID(u)
	return su
}

// SetNillableManagerID sets the "manager_id" field if the given value is not nil.
func (su *StockUpdate) SetNillableManagerID(u *uint64) *StockUpdate {
	if u != nil {
		su.SetManagerID(*u)
	}
	return su
}

// ClearManagerID clears the value of the "manager_id" field.
func (su *StockUpdate) ClearManagerID() *StockUpdate {
	su.mutation.ClearManagerID()
	return su
}

// SetSn sets the "sn" field.
func (su *StockUpdate) SetSn(s string) *StockUpdate {
	su.mutation.SetSn(s)
	return su
}

// SetType sets the "type" field.
func (su *StockUpdate) SetType(u uint8) *StockUpdate {
	su.mutation.ResetType()
	su.mutation.SetType(u)
	return su
}

// SetNillableType sets the "type" field if the given value is not nil.
func (su *StockUpdate) SetNillableType(u *uint8) *StockUpdate {
	if u != nil {
		su.SetType(*u)
	}
	return su
}

// AddType adds u to the "type" field.
func (su *StockUpdate) AddType(u int8) *StockUpdate {
	su.mutation.AddType(u)
	return su
}

// SetStoreID sets the "store_id" field.
func (su *StockUpdate) SetStoreID(u uint64) *StockUpdate {
	su.mutation.SetStoreID(u)
	return su
}

// SetNillableStoreID sets the "store_id" field if the given value is not nil.
func (su *StockUpdate) SetNillableStoreID(u *uint64) *StockUpdate {
	if u != nil {
		su.SetStoreID(*u)
	}
	return su
}

// ClearStoreID clears the value of the "store_id" field.
func (su *StockUpdate) ClearStoreID() *StockUpdate {
	su.mutation.ClearStoreID()
	return su
}

// SetRiderID sets the "rider_id" field.
func (su *StockUpdate) SetRiderID(u uint64) *StockUpdate {
	su.mutation.SetRiderID(u)
	return su
}

// SetNillableRiderID sets the "rider_id" field if the given value is not nil.
func (su *StockUpdate) SetNillableRiderID(u *uint64) *StockUpdate {
	if u != nil {
		su.SetRiderID(*u)
	}
	return su
}

// ClearRiderID clears the value of the "rider_id" field.
func (su *StockUpdate) ClearRiderID() *StockUpdate {
	su.mutation.ClearRiderID()
	return su
}

// SetEmployeeID sets the "employee_id" field.
func (su *StockUpdate) SetEmployeeID(u uint64) *StockUpdate {
	su.mutation.SetEmployeeID(u)
	return su
}

// SetNillableEmployeeID sets the "employee_id" field if the given value is not nil.
func (su *StockUpdate) SetNillableEmployeeID(u *uint64) *StockUpdate {
	if u != nil {
		su.SetEmployeeID(*u)
	}
	return su
}

// ClearEmployeeID clears the value of the "employee_id" field.
func (su *StockUpdate) ClearEmployeeID() *StockUpdate {
	su.mutation.ClearEmployeeID()
	return su
}

// SetName sets the "name" field.
func (su *StockUpdate) SetName(s string) *StockUpdate {
	su.mutation.SetName(s)
	return su
}

// SetVoltage sets the "voltage" field.
func (su *StockUpdate) SetVoltage(f float64) *StockUpdate {
	su.mutation.ResetVoltage()
	su.mutation.SetVoltage(f)
	return su
}

// SetNillableVoltage sets the "voltage" field if the given value is not nil.
func (su *StockUpdate) SetNillableVoltage(f *float64) *StockUpdate {
	if f != nil {
		su.SetVoltage(*f)
	}
	return su
}

// AddVoltage adds f to the "voltage" field.
func (su *StockUpdate) AddVoltage(f float64) *StockUpdate {
	su.mutation.AddVoltage(f)
	return su
}

// ClearVoltage clears the value of the "voltage" field.
func (su *StockUpdate) ClearVoltage() *StockUpdate {
	su.mutation.ClearVoltage()
	return su
}

// SetManager sets the "manager" edge to the Manager entity.
func (su *StockUpdate) SetManager(m *Manager) *StockUpdate {
	return su.SetManagerID(m.ID)
}

// SetStore sets the "store" edge to the Store entity.
func (su *StockUpdate) SetStore(s *Store) *StockUpdate {
	return su.SetStoreID(s.ID)
}

// SetRider sets the "rider" edge to the Rider entity.
func (su *StockUpdate) SetRider(r *Rider) *StockUpdate {
	return su.SetRiderID(r.ID)
}

// SetEmployee sets the "employee" edge to the Employee entity.
func (su *StockUpdate) SetEmployee(e *Employee) *StockUpdate {
	return su.SetEmployeeID(e.ID)
}

// Mutation returns the StockMutation object of the builder.
func (su *StockUpdate) Mutation() *StockMutation {
	return su.mutation
}

// ClearManager clears the "manager" edge to the Manager entity.
func (su *StockUpdate) ClearManager() *StockUpdate {
	su.mutation.ClearManager()
	return su
}

// ClearStore clears the "store" edge to the Store entity.
func (su *StockUpdate) ClearStore() *StockUpdate {
	su.mutation.ClearStore()
	return su
}

// ClearRider clears the "rider" edge to the Rider entity.
func (su *StockUpdate) ClearRider() *StockUpdate {
	su.mutation.ClearRider()
	return su
}

// ClearEmployee clears the "employee" edge to the Employee entity.
func (su *StockUpdate) ClearEmployee() *StockUpdate {
	su.mutation.ClearEmployee()
	return su
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *StockUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if err := su.defaults(); err != nil {
		return 0, err
	}
	if len(su.hooks) == 0 {
		affected, err = su.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*StockMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			su.mutation = mutation
			affected, err = su.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(su.hooks) - 1; i >= 0; i-- {
			if su.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = su.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, su.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (su *StockUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *StockUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *StockUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (su *StockUpdate) defaults() error {
	if _, ok := su.mutation.UpdatedAt(); !ok {
		if stock.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized stock.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := stock.UpdateDefaultUpdatedAt()
		su.mutation.SetUpdatedAt(v)
	}
	return nil
}

func (su *StockUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   stock.Table,
			Columns: stock.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: stock.FieldID,
			},
		},
	}
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: stock.FieldUpdatedAt,
		})
	}
	if value, ok := su.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: stock.FieldDeletedAt,
		})
	}
	if su.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: stock.FieldDeletedAt,
		})
	}
	if su.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: stock.FieldCreator,
		})
	}
	if value, ok := su.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: stock.FieldLastModifier,
		})
	}
	if su.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: stock.FieldLastModifier,
		})
	}
	if value, ok := su.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: stock.FieldRemark,
		})
	}
	if su.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: stock.FieldRemark,
		})
	}
	if value, ok := su.mutation.Sn(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: stock.FieldSn,
		})
	}
	if value, ok := su.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: stock.FieldType,
		})
	}
	if value, ok := su.mutation.AddedType(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: stock.FieldType,
		})
	}
	if value, ok := su.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: stock.FieldName,
		})
	}
	if value, ok := su.mutation.Voltage(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: stock.FieldVoltage,
		})
	}
	if value, ok := su.mutation.AddedVoltage(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: stock.FieldVoltage,
		})
	}
	if su.mutation.VoltageCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: stock.FieldVoltage,
		})
	}
	if su.mutation.ManagerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   stock.ManagerTable,
			Columns: []string{stock.ManagerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: manager.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.ManagerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   stock.ManagerTable,
			Columns: []string{stock.ManagerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: manager.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.StoreCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.StoreTable,
			Columns: []string{stock.StoreColumn},
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
	if nodes := su.mutation.StoreIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.StoreTable,
			Columns: []string{stock.StoreColumn},
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
	if su.mutation.RiderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.RiderTable,
			Columns: []string{stock.RiderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: rider.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RiderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.RiderTable,
			Columns: []string{stock.RiderColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.EmployeeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.EmployeeTable,
			Columns: []string{stock.EmployeeColumn},
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
	if nodes := su.mutation.EmployeeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.EmployeeTable,
			Columns: []string{stock.EmployeeColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{stock.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// StockUpdateOne is the builder for updating a single Stock entity.
type StockUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *StockMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (suo *StockUpdateOne) SetUpdatedAt(t time.Time) *StockUpdateOne {
	suo.mutation.SetUpdatedAt(t)
	return suo
}

// SetDeletedAt sets the "deleted_at" field.
func (suo *StockUpdateOne) SetDeletedAt(t time.Time) *StockUpdateOne {
	suo.mutation.SetDeletedAt(t)
	return suo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (suo *StockUpdateOne) SetNillableDeletedAt(t *time.Time) *StockUpdateOne {
	if t != nil {
		suo.SetDeletedAt(*t)
	}
	return suo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (suo *StockUpdateOne) ClearDeletedAt() *StockUpdateOne {
	suo.mutation.ClearDeletedAt()
	return suo
}

// SetLastModifier sets the "last_modifier" field.
func (suo *StockUpdateOne) SetLastModifier(m *model.Modifier) *StockUpdateOne {
	suo.mutation.SetLastModifier(m)
	return suo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (suo *StockUpdateOne) ClearLastModifier() *StockUpdateOne {
	suo.mutation.ClearLastModifier()
	return suo
}

// SetRemark sets the "remark" field.
func (suo *StockUpdateOne) SetRemark(s string) *StockUpdateOne {
	suo.mutation.SetRemark(s)
	return suo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (suo *StockUpdateOne) SetNillableRemark(s *string) *StockUpdateOne {
	if s != nil {
		suo.SetRemark(*s)
	}
	return suo
}

// ClearRemark clears the value of the "remark" field.
func (suo *StockUpdateOne) ClearRemark() *StockUpdateOne {
	suo.mutation.ClearRemark()
	return suo
}

// SetManagerID sets the "manager_id" field.
func (suo *StockUpdateOne) SetManagerID(u uint64) *StockUpdateOne {
	suo.mutation.SetManagerID(u)
	return suo
}

// SetNillableManagerID sets the "manager_id" field if the given value is not nil.
func (suo *StockUpdateOne) SetNillableManagerID(u *uint64) *StockUpdateOne {
	if u != nil {
		suo.SetManagerID(*u)
	}
	return suo
}

// ClearManagerID clears the value of the "manager_id" field.
func (suo *StockUpdateOne) ClearManagerID() *StockUpdateOne {
	suo.mutation.ClearManagerID()
	return suo
}

// SetSn sets the "sn" field.
func (suo *StockUpdateOne) SetSn(s string) *StockUpdateOne {
	suo.mutation.SetSn(s)
	return suo
}

// SetType sets the "type" field.
func (suo *StockUpdateOne) SetType(u uint8) *StockUpdateOne {
	suo.mutation.ResetType()
	suo.mutation.SetType(u)
	return suo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (suo *StockUpdateOne) SetNillableType(u *uint8) *StockUpdateOne {
	if u != nil {
		suo.SetType(*u)
	}
	return suo
}

// AddType adds u to the "type" field.
func (suo *StockUpdateOne) AddType(u int8) *StockUpdateOne {
	suo.mutation.AddType(u)
	return suo
}

// SetStoreID sets the "store_id" field.
func (suo *StockUpdateOne) SetStoreID(u uint64) *StockUpdateOne {
	suo.mutation.SetStoreID(u)
	return suo
}

// SetNillableStoreID sets the "store_id" field if the given value is not nil.
func (suo *StockUpdateOne) SetNillableStoreID(u *uint64) *StockUpdateOne {
	if u != nil {
		suo.SetStoreID(*u)
	}
	return suo
}

// ClearStoreID clears the value of the "store_id" field.
func (suo *StockUpdateOne) ClearStoreID() *StockUpdateOne {
	suo.mutation.ClearStoreID()
	return suo
}

// SetRiderID sets the "rider_id" field.
func (suo *StockUpdateOne) SetRiderID(u uint64) *StockUpdateOne {
	suo.mutation.SetRiderID(u)
	return suo
}

// SetNillableRiderID sets the "rider_id" field if the given value is not nil.
func (suo *StockUpdateOne) SetNillableRiderID(u *uint64) *StockUpdateOne {
	if u != nil {
		suo.SetRiderID(*u)
	}
	return suo
}

// ClearRiderID clears the value of the "rider_id" field.
func (suo *StockUpdateOne) ClearRiderID() *StockUpdateOne {
	suo.mutation.ClearRiderID()
	return suo
}

// SetEmployeeID sets the "employee_id" field.
func (suo *StockUpdateOne) SetEmployeeID(u uint64) *StockUpdateOne {
	suo.mutation.SetEmployeeID(u)
	return suo
}

// SetNillableEmployeeID sets the "employee_id" field if the given value is not nil.
func (suo *StockUpdateOne) SetNillableEmployeeID(u *uint64) *StockUpdateOne {
	if u != nil {
		suo.SetEmployeeID(*u)
	}
	return suo
}

// ClearEmployeeID clears the value of the "employee_id" field.
func (suo *StockUpdateOne) ClearEmployeeID() *StockUpdateOne {
	suo.mutation.ClearEmployeeID()
	return suo
}

// SetName sets the "name" field.
func (suo *StockUpdateOne) SetName(s string) *StockUpdateOne {
	suo.mutation.SetName(s)
	return suo
}

// SetVoltage sets the "voltage" field.
func (suo *StockUpdateOne) SetVoltage(f float64) *StockUpdateOne {
	suo.mutation.ResetVoltage()
	suo.mutation.SetVoltage(f)
	return suo
}

// SetNillableVoltage sets the "voltage" field if the given value is not nil.
func (suo *StockUpdateOne) SetNillableVoltage(f *float64) *StockUpdateOne {
	if f != nil {
		suo.SetVoltage(*f)
	}
	return suo
}

// AddVoltage adds f to the "voltage" field.
func (suo *StockUpdateOne) AddVoltage(f float64) *StockUpdateOne {
	suo.mutation.AddVoltage(f)
	return suo
}

// ClearVoltage clears the value of the "voltage" field.
func (suo *StockUpdateOne) ClearVoltage() *StockUpdateOne {
	suo.mutation.ClearVoltage()
	return suo
}

// SetManager sets the "manager" edge to the Manager entity.
func (suo *StockUpdateOne) SetManager(m *Manager) *StockUpdateOne {
	return suo.SetManagerID(m.ID)
}

// SetStore sets the "store" edge to the Store entity.
func (suo *StockUpdateOne) SetStore(s *Store) *StockUpdateOne {
	return suo.SetStoreID(s.ID)
}

// SetRider sets the "rider" edge to the Rider entity.
func (suo *StockUpdateOne) SetRider(r *Rider) *StockUpdateOne {
	return suo.SetRiderID(r.ID)
}

// SetEmployee sets the "employee" edge to the Employee entity.
func (suo *StockUpdateOne) SetEmployee(e *Employee) *StockUpdateOne {
	return suo.SetEmployeeID(e.ID)
}

// Mutation returns the StockMutation object of the builder.
func (suo *StockUpdateOne) Mutation() *StockMutation {
	return suo.mutation
}

// ClearManager clears the "manager" edge to the Manager entity.
func (suo *StockUpdateOne) ClearManager() *StockUpdateOne {
	suo.mutation.ClearManager()
	return suo
}

// ClearStore clears the "store" edge to the Store entity.
func (suo *StockUpdateOne) ClearStore() *StockUpdateOne {
	suo.mutation.ClearStore()
	return suo
}

// ClearRider clears the "rider" edge to the Rider entity.
func (suo *StockUpdateOne) ClearRider() *StockUpdateOne {
	suo.mutation.ClearRider()
	return suo
}

// ClearEmployee clears the "employee" edge to the Employee entity.
func (suo *StockUpdateOne) ClearEmployee() *StockUpdateOne {
	suo.mutation.ClearEmployee()
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *StockUpdateOne) Select(field string, fields ...string) *StockUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Stock entity.
func (suo *StockUpdateOne) Save(ctx context.Context) (*Stock, error) {
	var (
		err  error
		node *Stock
	)
	if err := suo.defaults(); err != nil {
		return nil, err
	}
	if len(suo.hooks) == 0 {
		node, err = suo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*StockMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			suo.mutation = mutation
			node, err = suo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(suo.hooks) - 1; i >= 0; i-- {
			if suo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = suo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, suo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Stock)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from StockMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (suo *StockUpdateOne) SaveX(ctx context.Context) *Stock {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *StockUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *StockUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (suo *StockUpdateOne) defaults() error {
	if _, ok := suo.mutation.UpdatedAt(); !ok {
		if stock.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized stock.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := stock.UpdateDefaultUpdatedAt()
		suo.mutation.SetUpdatedAt(v)
	}
	return nil
}

func (suo *StockUpdateOne) sqlSave(ctx context.Context) (_node *Stock, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   stock.Table,
			Columns: stock.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: stock.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Stock.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, stock.FieldID)
		for _, f := range fields {
			if !stock.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != stock.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: stock.FieldUpdatedAt,
		})
	}
	if value, ok := suo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: stock.FieldDeletedAt,
		})
	}
	if suo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: stock.FieldDeletedAt,
		})
	}
	if suo.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: stock.FieldCreator,
		})
	}
	if value, ok := suo.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: stock.FieldLastModifier,
		})
	}
	if suo.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: stock.FieldLastModifier,
		})
	}
	if value, ok := suo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: stock.FieldRemark,
		})
	}
	if suo.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: stock.FieldRemark,
		})
	}
	if value, ok := suo.mutation.Sn(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: stock.FieldSn,
		})
	}
	if value, ok := suo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: stock.FieldType,
		})
	}
	if value, ok := suo.mutation.AddedType(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: stock.FieldType,
		})
	}
	if value, ok := suo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: stock.FieldName,
		})
	}
	if value, ok := suo.mutation.Voltage(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: stock.FieldVoltage,
		})
	}
	if value, ok := suo.mutation.AddedVoltage(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: stock.FieldVoltage,
		})
	}
	if suo.mutation.VoltageCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Column: stock.FieldVoltage,
		})
	}
	if suo.mutation.ManagerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   stock.ManagerTable,
			Columns: []string{stock.ManagerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: manager.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.ManagerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   stock.ManagerTable,
			Columns: []string{stock.ManagerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: manager.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.StoreCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.StoreTable,
			Columns: []string{stock.StoreColumn},
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
	if nodes := suo.mutation.StoreIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.StoreTable,
			Columns: []string{stock.StoreColumn},
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
	if suo.mutation.RiderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.RiderTable,
			Columns: []string{stock.RiderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: rider.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RiderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.RiderTable,
			Columns: []string{stock.RiderColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.EmployeeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.EmployeeTable,
			Columns: []string{stock.EmployeeColumn},
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
	if nodes := suo.mutation.EmployeeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   stock.EmployeeTable,
			Columns: []string{stock.EmployeeColumn},
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
	_node = &Stock{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{stock.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}