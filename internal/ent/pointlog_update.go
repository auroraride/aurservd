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
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/pointlog"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
)

// PointLogUpdate is the builder for updating PointLog entities.
type PointLogUpdate struct {
	config
	hooks     []Hook
	mutation  *PointLogMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the PointLogUpdate builder.
func (plu *PointLogUpdate) Where(ps ...predicate.PointLog) *PointLogUpdate {
	plu.mutation.Where(ps...)
	return plu
}

// SetUpdatedAt sets the "updated_at" field.
func (plu *PointLogUpdate) SetUpdatedAt(t time.Time) *PointLogUpdate {
	plu.mutation.SetUpdatedAt(t)
	return plu
}

// SetRiderID sets the "rider_id" field.
func (plu *PointLogUpdate) SetRiderID(u uint64) *PointLogUpdate {
	plu.mutation.SetRiderID(u)
	return plu
}

// SetOrderID sets the "order_id" field.
func (plu *PointLogUpdate) SetOrderID(u uint64) *PointLogUpdate {
	plu.mutation.SetOrderID(u)
	return plu
}

// SetNillableOrderID sets the "order_id" field if the given value is not nil.
func (plu *PointLogUpdate) SetNillableOrderID(u *uint64) *PointLogUpdate {
	if u != nil {
		plu.SetOrderID(*u)
	}
	return plu
}

// ClearOrderID clears the value of the "order_id" field.
func (plu *PointLogUpdate) ClearOrderID() *PointLogUpdate {
	plu.mutation.ClearOrderID()
	return plu
}

// SetModifier sets the "modifier" field.
func (plu *PointLogUpdate) SetModifier(m *model.Modifier) *PointLogUpdate {
	plu.mutation.SetModifier(m)
	return plu
}

// ClearModifier clears the value of the "modifier" field.
func (plu *PointLogUpdate) ClearModifier() *PointLogUpdate {
	plu.mutation.ClearModifier()
	return plu
}

// SetEmployeeInfo sets the "employee_info" field.
func (plu *PointLogUpdate) SetEmployeeInfo(m *model.Employee) *PointLogUpdate {
	plu.mutation.SetEmployeeInfo(m)
	return plu
}

// ClearEmployeeInfo clears the value of the "employee_info" field.
func (plu *PointLogUpdate) ClearEmployeeInfo() *PointLogUpdate {
	plu.mutation.ClearEmployeeInfo()
	return plu
}

// SetType sets the "type" field.
func (plu *PointLogUpdate) SetType(u uint8) *PointLogUpdate {
	plu.mutation.ResetType()
	plu.mutation.SetType(u)
	return plu
}

// AddType adds u to the "type" field.
func (plu *PointLogUpdate) AddType(u int8) *PointLogUpdate {
	plu.mutation.AddType(u)
	return plu
}

// SetPoints sets the "points" field.
func (plu *PointLogUpdate) SetPoints(i int64) *PointLogUpdate {
	plu.mutation.ResetPoints()
	plu.mutation.SetPoints(i)
	return plu
}

// AddPoints adds i to the "points" field.
func (plu *PointLogUpdate) AddPoints(i int64) *PointLogUpdate {
	plu.mutation.AddPoints(i)
	return plu
}

// SetAfter sets the "after" field.
func (plu *PointLogUpdate) SetAfter(i int64) *PointLogUpdate {
	plu.mutation.ResetAfter()
	plu.mutation.SetAfter(i)
	return plu
}

// AddAfter adds i to the "after" field.
func (plu *PointLogUpdate) AddAfter(i int64) *PointLogUpdate {
	plu.mutation.AddAfter(i)
	return plu
}

// SetReason sets the "reason" field.
func (plu *PointLogUpdate) SetReason(s string) *PointLogUpdate {
	plu.mutation.SetReason(s)
	return plu
}

// SetNillableReason sets the "reason" field if the given value is not nil.
func (plu *PointLogUpdate) SetNillableReason(s *string) *PointLogUpdate {
	if s != nil {
		plu.SetReason(*s)
	}
	return plu
}

// ClearReason clears the value of the "reason" field.
func (plu *PointLogUpdate) ClearReason() *PointLogUpdate {
	plu.mutation.ClearReason()
	return plu
}

// SetAttach sets the "attach" field.
func (plu *PointLogUpdate) SetAttach(mla *model.PointLogAttach) *PointLogUpdate {
	plu.mutation.SetAttach(mla)
	return plu
}

// ClearAttach clears the value of the "attach" field.
func (plu *PointLogUpdate) ClearAttach() *PointLogUpdate {
	plu.mutation.ClearAttach()
	return plu
}

// SetRider sets the "rider" edge to the Rider entity.
func (plu *PointLogUpdate) SetRider(r *Rider) *PointLogUpdate {
	return plu.SetRiderID(r.ID)
}

// SetOrder sets the "order" edge to the Order entity.
func (plu *PointLogUpdate) SetOrder(o *Order) *PointLogUpdate {
	return plu.SetOrderID(o.ID)
}

// Mutation returns the PointLogMutation object of the builder.
func (plu *PointLogUpdate) Mutation() *PointLogMutation {
	return plu.mutation
}

// ClearRider clears the "rider" edge to the Rider entity.
func (plu *PointLogUpdate) ClearRider() *PointLogUpdate {
	plu.mutation.ClearRider()
	return plu
}

// ClearOrder clears the "order" edge to the Order entity.
func (plu *PointLogUpdate) ClearOrder() *PointLogUpdate {
	plu.mutation.ClearOrder()
	return plu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (plu *PointLogUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if err := plu.defaults(); err != nil {
		return 0, err
	}
	if len(plu.hooks) == 0 {
		if err = plu.check(); err != nil {
			return 0, err
		}
		affected, err = plu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PointLogMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = plu.check(); err != nil {
				return 0, err
			}
			plu.mutation = mutation
			affected, err = plu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(plu.hooks) - 1; i >= 0; i-- {
			if plu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = plu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, plu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (plu *PointLogUpdate) SaveX(ctx context.Context) int {
	affected, err := plu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (plu *PointLogUpdate) Exec(ctx context.Context) error {
	_, err := plu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (plu *PointLogUpdate) ExecX(ctx context.Context) {
	if err := plu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (plu *PointLogUpdate) defaults() error {
	if _, ok := plu.mutation.UpdatedAt(); !ok {
		if pointlog.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized pointlog.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := pointlog.UpdateDefaultUpdatedAt()
		plu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (plu *PointLogUpdate) check() error {
	if _, ok := plu.mutation.RiderID(); plu.mutation.RiderCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "PointLog.rider"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (plu *PointLogUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PointLogUpdate {
	plu.modifiers = append(plu.modifiers, modifiers...)
	return plu
}

func (plu *PointLogUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   pointlog.Table,
			Columns: pointlog.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: pointlog.FieldID,
			},
		},
	}
	if ps := plu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := plu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: pointlog.FieldUpdatedAt,
		})
	}
	if value, ok := plu.mutation.Modifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: pointlog.FieldModifier,
		})
	}
	if plu.mutation.ModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: pointlog.FieldModifier,
		})
	}
	if value, ok := plu.mutation.EmployeeInfo(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: pointlog.FieldEmployeeInfo,
		})
	}
	if plu.mutation.EmployeeInfoCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: pointlog.FieldEmployeeInfo,
		})
	}
	if value, ok := plu.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: pointlog.FieldType,
		})
	}
	if value, ok := plu.mutation.AddedType(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: pointlog.FieldType,
		})
	}
	if value, ok := plu.mutation.Points(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: pointlog.FieldPoints,
		})
	}
	if value, ok := plu.mutation.AddedPoints(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: pointlog.FieldPoints,
		})
	}
	if value, ok := plu.mutation.After(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: pointlog.FieldAfter,
		})
	}
	if value, ok := plu.mutation.AddedAfter(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: pointlog.FieldAfter,
		})
	}
	if value, ok := plu.mutation.Reason(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: pointlog.FieldReason,
		})
	}
	if plu.mutation.ReasonCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: pointlog.FieldReason,
		})
	}
	if value, ok := plu.mutation.Attach(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: pointlog.FieldAttach,
		})
	}
	if plu.mutation.AttachCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: pointlog.FieldAttach,
		})
	}
	if plu.mutation.RiderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   pointlog.RiderTable,
			Columns: []string{pointlog.RiderColumn},
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
	if nodes := plu.mutation.RiderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   pointlog.RiderTable,
			Columns: []string{pointlog.RiderColumn},
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
	if plu.mutation.OrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   pointlog.OrderTable,
			Columns: []string{pointlog.OrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: order.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := plu.mutation.OrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   pointlog.OrderTable,
			Columns: []string{pointlog.OrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: order.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Modifiers = plu.modifiers
	if n, err = sqlgraph.UpdateNodes(ctx, plu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{pointlog.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	return n, nil
}

// PointLogUpdateOne is the builder for updating a single PointLog entity.
type PointLogUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *PointLogMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (pluo *PointLogUpdateOne) SetUpdatedAt(t time.Time) *PointLogUpdateOne {
	pluo.mutation.SetUpdatedAt(t)
	return pluo
}

// SetRiderID sets the "rider_id" field.
func (pluo *PointLogUpdateOne) SetRiderID(u uint64) *PointLogUpdateOne {
	pluo.mutation.SetRiderID(u)
	return pluo
}

// SetOrderID sets the "order_id" field.
func (pluo *PointLogUpdateOne) SetOrderID(u uint64) *PointLogUpdateOne {
	pluo.mutation.SetOrderID(u)
	return pluo
}

// SetNillableOrderID sets the "order_id" field if the given value is not nil.
func (pluo *PointLogUpdateOne) SetNillableOrderID(u *uint64) *PointLogUpdateOne {
	if u != nil {
		pluo.SetOrderID(*u)
	}
	return pluo
}

// ClearOrderID clears the value of the "order_id" field.
func (pluo *PointLogUpdateOne) ClearOrderID() *PointLogUpdateOne {
	pluo.mutation.ClearOrderID()
	return pluo
}

// SetModifier sets the "modifier" field.
func (pluo *PointLogUpdateOne) SetModifier(m *model.Modifier) *PointLogUpdateOne {
	pluo.mutation.SetModifier(m)
	return pluo
}

// ClearModifier clears the value of the "modifier" field.
func (pluo *PointLogUpdateOne) ClearModifier() *PointLogUpdateOne {
	pluo.mutation.ClearModifier()
	return pluo
}

// SetEmployeeInfo sets the "employee_info" field.
func (pluo *PointLogUpdateOne) SetEmployeeInfo(m *model.Employee) *PointLogUpdateOne {
	pluo.mutation.SetEmployeeInfo(m)
	return pluo
}

// ClearEmployeeInfo clears the value of the "employee_info" field.
func (pluo *PointLogUpdateOne) ClearEmployeeInfo() *PointLogUpdateOne {
	pluo.mutation.ClearEmployeeInfo()
	return pluo
}

// SetType sets the "type" field.
func (pluo *PointLogUpdateOne) SetType(u uint8) *PointLogUpdateOne {
	pluo.mutation.ResetType()
	pluo.mutation.SetType(u)
	return pluo
}

// AddType adds u to the "type" field.
func (pluo *PointLogUpdateOne) AddType(u int8) *PointLogUpdateOne {
	pluo.mutation.AddType(u)
	return pluo
}

// SetPoints sets the "points" field.
func (pluo *PointLogUpdateOne) SetPoints(i int64) *PointLogUpdateOne {
	pluo.mutation.ResetPoints()
	pluo.mutation.SetPoints(i)
	return pluo
}

// AddPoints adds i to the "points" field.
func (pluo *PointLogUpdateOne) AddPoints(i int64) *PointLogUpdateOne {
	pluo.mutation.AddPoints(i)
	return pluo
}

// SetAfter sets the "after" field.
func (pluo *PointLogUpdateOne) SetAfter(i int64) *PointLogUpdateOne {
	pluo.mutation.ResetAfter()
	pluo.mutation.SetAfter(i)
	return pluo
}

// AddAfter adds i to the "after" field.
func (pluo *PointLogUpdateOne) AddAfter(i int64) *PointLogUpdateOne {
	pluo.mutation.AddAfter(i)
	return pluo
}

// SetReason sets the "reason" field.
func (pluo *PointLogUpdateOne) SetReason(s string) *PointLogUpdateOne {
	pluo.mutation.SetReason(s)
	return pluo
}

// SetNillableReason sets the "reason" field if the given value is not nil.
func (pluo *PointLogUpdateOne) SetNillableReason(s *string) *PointLogUpdateOne {
	if s != nil {
		pluo.SetReason(*s)
	}
	return pluo
}

// ClearReason clears the value of the "reason" field.
func (pluo *PointLogUpdateOne) ClearReason() *PointLogUpdateOne {
	pluo.mutation.ClearReason()
	return pluo
}

// SetAttach sets the "attach" field.
func (pluo *PointLogUpdateOne) SetAttach(mla *model.PointLogAttach) *PointLogUpdateOne {
	pluo.mutation.SetAttach(mla)
	return pluo
}

// ClearAttach clears the value of the "attach" field.
func (pluo *PointLogUpdateOne) ClearAttach() *PointLogUpdateOne {
	pluo.mutation.ClearAttach()
	return pluo
}

// SetRider sets the "rider" edge to the Rider entity.
func (pluo *PointLogUpdateOne) SetRider(r *Rider) *PointLogUpdateOne {
	return pluo.SetRiderID(r.ID)
}

// SetOrder sets the "order" edge to the Order entity.
func (pluo *PointLogUpdateOne) SetOrder(o *Order) *PointLogUpdateOne {
	return pluo.SetOrderID(o.ID)
}

// Mutation returns the PointLogMutation object of the builder.
func (pluo *PointLogUpdateOne) Mutation() *PointLogMutation {
	return pluo.mutation
}

// ClearRider clears the "rider" edge to the Rider entity.
func (pluo *PointLogUpdateOne) ClearRider() *PointLogUpdateOne {
	pluo.mutation.ClearRider()
	return pluo
}

// ClearOrder clears the "order" edge to the Order entity.
func (pluo *PointLogUpdateOne) ClearOrder() *PointLogUpdateOne {
	pluo.mutation.ClearOrder()
	return pluo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (pluo *PointLogUpdateOne) Select(field string, fields ...string) *PointLogUpdateOne {
	pluo.fields = append([]string{field}, fields...)
	return pluo
}

// Save executes the query and returns the updated PointLog entity.
func (pluo *PointLogUpdateOne) Save(ctx context.Context) (*PointLog, error) {
	var (
		err  error
		node *PointLog
	)
	if err := pluo.defaults(); err != nil {
		return nil, err
	}
	if len(pluo.hooks) == 0 {
		if err = pluo.check(); err != nil {
			return nil, err
		}
		node, err = pluo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PointLogMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = pluo.check(); err != nil {
				return nil, err
			}
			pluo.mutation = mutation
			node, err = pluo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(pluo.hooks) - 1; i >= 0; i-- {
			if pluo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = pluo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, pluo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*PointLog)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from PointLogMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (pluo *PointLogUpdateOne) SaveX(ctx context.Context) *PointLog {
	node, err := pluo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (pluo *PointLogUpdateOne) Exec(ctx context.Context) error {
	_, err := pluo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pluo *PointLogUpdateOne) ExecX(ctx context.Context) {
	if err := pluo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pluo *PointLogUpdateOne) defaults() error {
	if _, ok := pluo.mutation.UpdatedAt(); !ok {
		if pointlog.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized pointlog.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := pointlog.UpdateDefaultUpdatedAt()
		pluo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (pluo *PointLogUpdateOne) check() error {
	if _, ok := pluo.mutation.RiderID(); pluo.mutation.RiderCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "PointLog.rider"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (pluo *PointLogUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PointLogUpdateOne {
	pluo.modifiers = append(pluo.modifiers, modifiers...)
	return pluo
}

func (pluo *PointLogUpdateOne) sqlSave(ctx context.Context) (_node *PointLog, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   pointlog.Table,
			Columns: pointlog.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: pointlog.FieldID,
			},
		},
	}
	id, ok := pluo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "PointLog.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := pluo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, pointlog.FieldID)
		for _, f := range fields {
			if !pointlog.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != pointlog.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := pluo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pluo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: pointlog.FieldUpdatedAt,
		})
	}
	if value, ok := pluo.mutation.Modifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: pointlog.FieldModifier,
		})
	}
	if pluo.mutation.ModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: pointlog.FieldModifier,
		})
	}
	if value, ok := pluo.mutation.EmployeeInfo(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: pointlog.FieldEmployeeInfo,
		})
	}
	if pluo.mutation.EmployeeInfoCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: pointlog.FieldEmployeeInfo,
		})
	}
	if value, ok := pluo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: pointlog.FieldType,
		})
	}
	if value, ok := pluo.mutation.AddedType(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint8,
			Value:  value,
			Column: pointlog.FieldType,
		})
	}
	if value, ok := pluo.mutation.Points(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: pointlog.FieldPoints,
		})
	}
	if value, ok := pluo.mutation.AddedPoints(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: pointlog.FieldPoints,
		})
	}
	if value, ok := pluo.mutation.After(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: pointlog.FieldAfter,
		})
	}
	if value, ok := pluo.mutation.AddedAfter(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: pointlog.FieldAfter,
		})
	}
	if value, ok := pluo.mutation.Reason(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: pointlog.FieldReason,
		})
	}
	if pluo.mutation.ReasonCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: pointlog.FieldReason,
		})
	}
	if value, ok := pluo.mutation.Attach(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: pointlog.FieldAttach,
		})
	}
	if pluo.mutation.AttachCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: pointlog.FieldAttach,
		})
	}
	if pluo.mutation.RiderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   pointlog.RiderTable,
			Columns: []string{pointlog.RiderColumn},
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
	if nodes := pluo.mutation.RiderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   pointlog.RiderTable,
			Columns: []string{pointlog.RiderColumn},
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
	if pluo.mutation.OrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   pointlog.OrderTable,
			Columns: []string{pointlog.OrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: order.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pluo.mutation.OrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   pointlog.OrderTable,
			Columns: []string{pointlog.OrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: order.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Modifiers = pluo.modifiers
	_node = &PointLog{config: pluo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, pluo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{pointlog.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}