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
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribepause"
)

// SubscribePauseUpdate is the builder for updating SubscribePause entities.
type SubscribePauseUpdate struct {
	config
	hooks    []Hook
	mutation *SubscribePauseMutation
}

// Where appends a list predicates to the SubscribePauseUpdate builder.
func (spu *SubscribePauseUpdate) Where(ps ...predicate.SubscribePause) *SubscribePauseUpdate {
	spu.mutation.Where(ps...)
	return spu
}

// SetUpdatedAt sets the "updated_at" field.
func (spu *SubscribePauseUpdate) SetUpdatedAt(t time.Time) *SubscribePauseUpdate {
	spu.mutation.SetUpdatedAt(t)
	return spu
}

// SetDeletedAt sets the "deleted_at" field.
func (spu *SubscribePauseUpdate) SetDeletedAt(t time.Time) *SubscribePauseUpdate {
	spu.mutation.SetDeletedAt(t)
	return spu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (spu *SubscribePauseUpdate) SetNillableDeletedAt(t *time.Time) *SubscribePauseUpdate {
	if t != nil {
		spu.SetDeletedAt(*t)
	}
	return spu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (spu *SubscribePauseUpdate) ClearDeletedAt() *SubscribePauseUpdate {
	spu.mutation.ClearDeletedAt()
	return spu
}

// SetLastModifier sets the "last_modifier" field.
func (spu *SubscribePauseUpdate) SetLastModifier(m *model.Modifier) *SubscribePauseUpdate {
	spu.mutation.SetLastModifier(m)
	return spu
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (spu *SubscribePauseUpdate) ClearLastModifier() *SubscribePauseUpdate {
	spu.mutation.ClearLastModifier()
	return spu
}

// SetRemark sets the "remark" field.
func (spu *SubscribePauseUpdate) SetRemark(s string) *SubscribePauseUpdate {
	spu.mutation.SetRemark(s)
	return spu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (spu *SubscribePauseUpdate) SetNillableRemark(s *string) *SubscribePauseUpdate {
	if s != nil {
		spu.SetRemark(*s)
	}
	return spu
}

// ClearRemark clears the value of the "remark" field.
func (spu *SubscribePauseUpdate) ClearRemark() *SubscribePauseUpdate {
	spu.mutation.ClearRemark()
	return spu
}

// SetRiderID sets the "rider_id" field.
func (spu *SubscribePauseUpdate) SetRiderID(u uint64) *SubscribePauseUpdate {
	spu.mutation.SetRiderID(u)
	return spu
}

// SetEmployeeID sets the "employee_id" field.
func (spu *SubscribePauseUpdate) SetEmployeeID(u uint64) *SubscribePauseUpdate {
	spu.mutation.SetEmployeeID(u)
	return spu
}

// SetNillableEmployeeID sets the "employee_id" field if the given value is not nil.
func (spu *SubscribePauseUpdate) SetNillableEmployeeID(u *uint64) *SubscribePauseUpdate {
	if u != nil {
		spu.SetEmployeeID(*u)
	}
	return spu
}

// ClearEmployeeID clears the value of the "employee_id" field.
func (spu *SubscribePauseUpdate) ClearEmployeeID() *SubscribePauseUpdate {
	spu.mutation.ClearEmployeeID()
	return spu
}

// SetSubscribeID sets the "subscribe_id" field.
func (spu *SubscribePauseUpdate) SetSubscribeID(u uint64) *SubscribePauseUpdate {
	spu.mutation.SetSubscribeID(u)
	return spu
}

// SetStartAt sets the "start_at" field.
func (spu *SubscribePauseUpdate) SetStartAt(t time.Time) *SubscribePauseUpdate {
	spu.mutation.SetStartAt(t)
	return spu
}

// SetEndAt sets the "end_at" field.
func (spu *SubscribePauseUpdate) SetEndAt(t time.Time) *SubscribePauseUpdate {
	spu.mutation.SetEndAt(t)
	return spu
}

// SetNillableEndAt sets the "end_at" field if the given value is not nil.
func (spu *SubscribePauseUpdate) SetNillableEndAt(t *time.Time) *SubscribePauseUpdate {
	if t != nil {
		spu.SetEndAt(*t)
	}
	return spu
}

// ClearEndAt clears the value of the "end_at" field.
func (spu *SubscribePauseUpdate) ClearEndAt() *SubscribePauseUpdate {
	spu.mutation.ClearEndAt()
	return spu
}

// SetDays sets the "days" field.
func (spu *SubscribePauseUpdate) SetDays(i int) *SubscribePauseUpdate {
	spu.mutation.ResetDays()
	spu.mutation.SetDays(i)
	return spu
}

// SetNillableDays sets the "days" field if the given value is not nil.
func (spu *SubscribePauseUpdate) SetNillableDays(i *int) *SubscribePauseUpdate {
	if i != nil {
		spu.SetDays(*i)
	}
	return spu
}

// AddDays adds i to the "days" field.
func (spu *SubscribePauseUpdate) AddDays(i int) *SubscribePauseUpdate {
	spu.mutation.AddDays(i)
	return spu
}

// ClearDays clears the value of the "days" field.
func (spu *SubscribePauseUpdate) ClearDays() *SubscribePauseUpdate {
	spu.mutation.ClearDays()
	return spu
}

// SetRider sets the "rider" edge to the Rider entity.
func (spu *SubscribePauseUpdate) SetRider(r *Rider) *SubscribePauseUpdate {
	return spu.SetRiderID(r.ID)
}

// SetEmployee sets the "employee" edge to the Employee entity.
func (spu *SubscribePauseUpdate) SetEmployee(e *Employee) *SubscribePauseUpdate {
	return spu.SetEmployeeID(e.ID)
}

// SetSubscribe sets the "subscribe" edge to the Subscribe entity.
func (spu *SubscribePauseUpdate) SetSubscribe(s *Subscribe) *SubscribePauseUpdate {
	return spu.SetSubscribeID(s.ID)
}

// Mutation returns the SubscribePauseMutation object of the builder.
func (spu *SubscribePauseUpdate) Mutation() *SubscribePauseMutation {
	return spu.mutation
}

// ClearRider clears the "rider" edge to the Rider entity.
func (spu *SubscribePauseUpdate) ClearRider() *SubscribePauseUpdate {
	spu.mutation.ClearRider()
	return spu
}

// ClearEmployee clears the "employee" edge to the Employee entity.
func (spu *SubscribePauseUpdate) ClearEmployee() *SubscribePauseUpdate {
	spu.mutation.ClearEmployee()
	return spu
}

// ClearSubscribe clears the "subscribe" edge to the Subscribe entity.
func (spu *SubscribePauseUpdate) ClearSubscribe() *SubscribePauseUpdate {
	spu.mutation.ClearSubscribe()
	return spu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (spu *SubscribePauseUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if err := spu.defaults(); err != nil {
		return 0, err
	}
	if len(spu.hooks) == 0 {
		if err = spu.check(); err != nil {
			return 0, err
		}
		affected, err = spu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SubscribePauseMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = spu.check(); err != nil {
				return 0, err
			}
			spu.mutation = mutation
			affected, err = spu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(spu.hooks) - 1; i >= 0; i-- {
			if spu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = spu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, spu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (spu *SubscribePauseUpdate) SaveX(ctx context.Context) int {
	affected, err := spu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (spu *SubscribePauseUpdate) Exec(ctx context.Context) error {
	_, err := spu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (spu *SubscribePauseUpdate) ExecX(ctx context.Context) {
	if err := spu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (spu *SubscribePauseUpdate) defaults() error {
	if _, ok := spu.mutation.UpdatedAt(); !ok {
		if subscribepause.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized subscribepause.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := subscribepause.UpdateDefaultUpdatedAt()
		spu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (spu *SubscribePauseUpdate) check() error {
	if _, ok := spu.mutation.RiderID(); spu.mutation.RiderCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SubscribePause.rider"`)
	}
	if _, ok := spu.mutation.SubscribeID(); spu.mutation.SubscribeCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SubscribePause.subscribe"`)
	}
	return nil
}

func (spu *SubscribePauseUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   subscribepause.Table,
			Columns: subscribepause.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: subscribepause.FieldID,
			},
		},
	}
	if ps := spu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := spu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribepause.FieldUpdatedAt,
		})
	}
	if value, ok := spu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribepause.FieldDeletedAt,
		})
	}
	if spu.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: subscribepause.FieldDeletedAt,
		})
	}
	if spu.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: subscribepause.FieldCreator,
		})
	}
	if value, ok := spu.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: subscribepause.FieldLastModifier,
		})
	}
	if spu.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: subscribepause.FieldLastModifier,
		})
	}
	if value, ok := spu.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: subscribepause.FieldRemark,
		})
	}
	if spu.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: subscribepause.FieldRemark,
		})
	}
	if value, ok := spu.mutation.StartAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribepause.FieldStartAt,
		})
	}
	if value, ok := spu.mutation.EndAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribepause.FieldEndAt,
		})
	}
	if spu.mutation.EndAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: subscribepause.FieldEndAt,
		})
	}
	if value, ok := spu.mutation.Days(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: subscribepause.FieldDays,
		})
	}
	if value, ok := spu.mutation.AddedDays(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: subscribepause.FieldDays,
		})
	}
	if spu.mutation.DaysCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: subscribepause.FieldDays,
		})
	}
	if spu.mutation.RiderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribepause.RiderTable,
			Columns: []string{subscribepause.RiderColumn},
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
	if nodes := spu.mutation.RiderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribepause.RiderTable,
			Columns: []string{subscribepause.RiderColumn},
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
	if spu.mutation.EmployeeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribepause.EmployeeTable,
			Columns: []string{subscribepause.EmployeeColumn},
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
	if nodes := spu.mutation.EmployeeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribepause.EmployeeTable,
			Columns: []string{subscribepause.EmployeeColumn},
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
	if spu.mutation.SubscribeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   subscribepause.SubscribeTable,
			Columns: []string{subscribepause.SubscribeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: subscribe.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := spu.mutation.SubscribeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   subscribepause.SubscribeTable,
			Columns: []string{subscribepause.SubscribeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: subscribe.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, spu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{subscribepause.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// SubscribePauseUpdateOne is the builder for updating a single SubscribePause entity.
type SubscribePauseUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SubscribePauseMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (spuo *SubscribePauseUpdateOne) SetUpdatedAt(t time.Time) *SubscribePauseUpdateOne {
	spuo.mutation.SetUpdatedAt(t)
	return spuo
}

// SetDeletedAt sets the "deleted_at" field.
func (spuo *SubscribePauseUpdateOne) SetDeletedAt(t time.Time) *SubscribePauseUpdateOne {
	spuo.mutation.SetDeletedAt(t)
	return spuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (spuo *SubscribePauseUpdateOne) SetNillableDeletedAt(t *time.Time) *SubscribePauseUpdateOne {
	if t != nil {
		spuo.SetDeletedAt(*t)
	}
	return spuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (spuo *SubscribePauseUpdateOne) ClearDeletedAt() *SubscribePauseUpdateOne {
	spuo.mutation.ClearDeletedAt()
	return spuo
}

// SetLastModifier sets the "last_modifier" field.
func (spuo *SubscribePauseUpdateOne) SetLastModifier(m *model.Modifier) *SubscribePauseUpdateOne {
	spuo.mutation.SetLastModifier(m)
	return spuo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (spuo *SubscribePauseUpdateOne) ClearLastModifier() *SubscribePauseUpdateOne {
	spuo.mutation.ClearLastModifier()
	return spuo
}

// SetRemark sets the "remark" field.
func (spuo *SubscribePauseUpdateOne) SetRemark(s string) *SubscribePauseUpdateOne {
	spuo.mutation.SetRemark(s)
	return spuo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (spuo *SubscribePauseUpdateOne) SetNillableRemark(s *string) *SubscribePauseUpdateOne {
	if s != nil {
		spuo.SetRemark(*s)
	}
	return spuo
}

// ClearRemark clears the value of the "remark" field.
func (spuo *SubscribePauseUpdateOne) ClearRemark() *SubscribePauseUpdateOne {
	spuo.mutation.ClearRemark()
	return spuo
}

// SetRiderID sets the "rider_id" field.
func (spuo *SubscribePauseUpdateOne) SetRiderID(u uint64) *SubscribePauseUpdateOne {
	spuo.mutation.SetRiderID(u)
	return spuo
}

// SetEmployeeID sets the "employee_id" field.
func (spuo *SubscribePauseUpdateOne) SetEmployeeID(u uint64) *SubscribePauseUpdateOne {
	spuo.mutation.SetEmployeeID(u)
	return spuo
}

// SetNillableEmployeeID sets the "employee_id" field if the given value is not nil.
func (spuo *SubscribePauseUpdateOne) SetNillableEmployeeID(u *uint64) *SubscribePauseUpdateOne {
	if u != nil {
		spuo.SetEmployeeID(*u)
	}
	return spuo
}

// ClearEmployeeID clears the value of the "employee_id" field.
func (spuo *SubscribePauseUpdateOne) ClearEmployeeID() *SubscribePauseUpdateOne {
	spuo.mutation.ClearEmployeeID()
	return spuo
}

// SetSubscribeID sets the "subscribe_id" field.
func (spuo *SubscribePauseUpdateOne) SetSubscribeID(u uint64) *SubscribePauseUpdateOne {
	spuo.mutation.SetSubscribeID(u)
	return spuo
}

// SetStartAt sets the "start_at" field.
func (spuo *SubscribePauseUpdateOne) SetStartAt(t time.Time) *SubscribePauseUpdateOne {
	spuo.mutation.SetStartAt(t)
	return spuo
}

// SetEndAt sets the "end_at" field.
func (spuo *SubscribePauseUpdateOne) SetEndAt(t time.Time) *SubscribePauseUpdateOne {
	spuo.mutation.SetEndAt(t)
	return spuo
}

// SetNillableEndAt sets the "end_at" field if the given value is not nil.
func (spuo *SubscribePauseUpdateOne) SetNillableEndAt(t *time.Time) *SubscribePauseUpdateOne {
	if t != nil {
		spuo.SetEndAt(*t)
	}
	return spuo
}

// ClearEndAt clears the value of the "end_at" field.
func (spuo *SubscribePauseUpdateOne) ClearEndAt() *SubscribePauseUpdateOne {
	spuo.mutation.ClearEndAt()
	return spuo
}

// SetDays sets the "days" field.
func (spuo *SubscribePauseUpdateOne) SetDays(i int) *SubscribePauseUpdateOne {
	spuo.mutation.ResetDays()
	spuo.mutation.SetDays(i)
	return spuo
}

// SetNillableDays sets the "days" field if the given value is not nil.
func (spuo *SubscribePauseUpdateOne) SetNillableDays(i *int) *SubscribePauseUpdateOne {
	if i != nil {
		spuo.SetDays(*i)
	}
	return spuo
}

// AddDays adds i to the "days" field.
func (spuo *SubscribePauseUpdateOne) AddDays(i int) *SubscribePauseUpdateOne {
	spuo.mutation.AddDays(i)
	return spuo
}

// ClearDays clears the value of the "days" field.
func (spuo *SubscribePauseUpdateOne) ClearDays() *SubscribePauseUpdateOne {
	spuo.mutation.ClearDays()
	return spuo
}

// SetRider sets the "rider" edge to the Rider entity.
func (spuo *SubscribePauseUpdateOne) SetRider(r *Rider) *SubscribePauseUpdateOne {
	return spuo.SetRiderID(r.ID)
}

// SetEmployee sets the "employee" edge to the Employee entity.
func (spuo *SubscribePauseUpdateOne) SetEmployee(e *Employee) *SubscribePauseUpdateOne {
	return spuo.SetEmployeeID(e.ID)
}

// SetSubscribe sets the "subscribe" edge to the Subscribe entity.
func (spuo *SubscribePauseUpdateOne) SetSubscribe(s *Subscribe) *SubscribePauseUpdateOne {
	return spuo.SetSubscribeID(s.ID)
}

// Mutation returns the SubscribePauseMutation object of the builder.
func (spuo *SubscribePauseUpdateOne) Mutation() *SubscribePauseMutation {
	return spuo.mutation
}

// ClearRider clears the "rider" edge to the Rider entity.
func (spuo *SubscribePauseUpdateOne) ClearRider() *SubscribePauseUpdateOne {
	spuo.mutation.ClearRider()
	return spuo
}

// ClearEmployee clears the "employee" edge to the Employee entity.
func (spuo *SubscribePauseUpdateOne) ClearEmployee() *SubscribePauseUpdateOne {
	spuo.mutation.ClearEmployee()
	return spuo
}

// ClearSubscribe clears the "subscribe" edge to the Subscribe entity.
func (spuo *SubscribePauseUpdateOne) ClearSubscribe() *SubscribePauseUpdateOne {
	spuo.mutation.ClearSubscribe()
	return spuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (spuo *SubscribePauseUpdateOne) Select(field string, fields ...string) *SubscribePauseUpdateOne {
	spuo.fields = append([]string{field}, fields...)
	return spuo
}

// Save executes the query and returns the updated SubscribePause entity.
func (spuo *SubscribePauseUpdateOne) Save(ctx context.Context) (*SubscribePause, error) {
	var (
		err  error
		node *SubscribePause
	)
	if err := spuo.defaults(); err != nil {
		return nil, err
	}
	if len(spuo.hooks) == 0 {
		if err = spuo.check(); err != nil {
			return nil, err
		}
		node, err = spuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SubscribePauseMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = spuo.check(); err != nil {
				return nil, err
			}
			spuo.mutation = mutation
			node, err = spuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(spuo.hooks) - 1; i >= 0; i-- {
			if spuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = spuo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, spuo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*SubscribePause)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from SubscribePauseMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (spuo *SubscribePauseUpdateOne) SaveX(ctx context.Context) *SubscribePause {
	node, err := spuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (spuo *SubscribePauseUpdateOne) Exec(ctx context.Context) error {
	_, err := spuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (spuo *SubscribePauseUpdateOne) ExecX(ctx context.Context) {
	if err := spuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (spuo *SubscribePauseUpdateOne) defaults() error {
	if _, ok := spuo.mutation.UpdatedAt(); !ok {
		if subscribepause.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized subscribepause.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := subscribepause.UpdateDefaultUpdatedAt()
		spuo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (spuo *SubscribePauseUpdateOne) check() error {
	if _, ok := spuo.mutation.RiderID(); spuo.mutation.RiderCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SubscribePause.rider"`)
	}
	if _, ok := spuo.mutation.SubscribeID(); spuo.mutation.SubscribeCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SubscribePause.subscribe"`)
	}
	return nil
}

func (spuo *SubscribePauseUpdateOne) sqlSave(ctx context.Context) (_node *SubscribePause, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   subscribepause.Table,
			Columns: subscribepause.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: subscribepause.FieldID,
			},
		},
	}
	id, ok := spuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "SubscribePause.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := spuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, subscribepause.FieldID)
		for _, f := range fields {
			if !subscribepause.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != subscribepause.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := spuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := spuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribepause.FieldUpdatedAt,
		})
	}
	if value, ok := spuo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribepause.FieldDeletedAt,
		})
	}
	if spuo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: subscribepause.FieldDeletedAt,
		})
	}
	if spuo.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: subscribepause.FieldCreator,
		})
	}
	if value, ok := spuo.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: subscribepause.FieldLastModifier,
		})
	}
	if spuo.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: subscribepause.FieldLastModifier,
		})
	}
	if value, ok := spuo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: subscribepause.FieldRemark,
		})
	}
	if spuo.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: subscribepause.FieldRemark,
		})
	}
	if value, ok := spuo.mutation.StartAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribepause.FieldStartAt,
		})
	}
	if value, ok := spuo.mutation.EndAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribepause.FieldEndAt,
		})
	}
	if spuo.mutation.EndAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: subscribepause.FieldEndAt,
		})
	}
	if value, ok := spuo.mutation.Days(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: subscribepause.FieldDays,
		})
	}
	if value, ok := spuo.mutation.AddedDays(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: subscribepause.FieldDays,
		})
	}
	if spuo.mutation.DaysCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: subscribepause.FieldDays,
		})
	}
	if spuo.mutation.RiderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribepause.RiderTable,
			Columns: []string{subscribepause.RiderColumn},
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
	if nodes := spuo.mutation.RiderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribepause.RiderTable,
			Columns: []string{subscribepause.RiderColumn},
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
	if spuo.mutation.EmployeeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribepause.EmployeeTable,
			Columns: []string{subscribepause.EmployeeColumn},
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
	if nodes := spuo.mutation.EmployeeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribepause.EmployeeTable,
			Columns: []string{subscribepause.EmployeeColumn},
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
	if spuo.mutation.SubscribeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   subscribepause.SubscribeTable,
			Columns: []string{subscribepause.SubscribeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: subscribe.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := spuo.mutation.SubscribeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   subscribepause.SubscribeTable,
			Columns: []string{subscribepause.SubscribeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: subscribe.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &SubscribePause{config: spuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, spuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{subscribepause.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}