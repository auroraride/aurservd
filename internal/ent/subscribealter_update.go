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
	"github.com/auroraride/aurservd/internal/ent/manager"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribealter"
)

// SubscribeAlterUpdate is the builder for updating SubscribeAlter entities.
type SubscribeAlterUpdate struct {
	config
	hooks    []Hook
	mutation *SubscribeAlterMutation
}

// Where appends a list predicates to the SubscribeAlterUpdate builder.
func (sau *SubscribeAlterUpdate) Where(ps ...predicate.SubscribeAlter) *SubscribeAlterUpdate {
	sau.mutation.Where(ps...)
	return sau
}

// SetUpdatedAt sets the "updated_at" field.
func (sau *SubscribeAlterUpdate) SetUpdatedAt(t time.Time) *SubscribeAlterUpdate {
	sau.mutation.SetUpdatedAt(t)
	return sau
}

// SetDeletedAt sets the "deleted_at" field.
func (sau *SubscribeAlterUpdate) SetDeletedAt(t time.Time) *SubscribeAlterUpdate {
	sau.mutation.SetDeletedAt(t)
	return sau
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (sau *SubscribeAlterUpdate) SetNillableDeletedAt(t *time.Time) *SubscribeAlterUpdate {
	if t != nil {
		sau.SetDeletedAt(*t)
	}
	return sau
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (sau *SubscribeAlterUpdate) ClearDeletedAt() *SubscribeAlterUpdate {
	sau.mutation.ClearDeletedAt()
	return sau
}

// SetLastModifier sets the "last_modifier" field.
func (sau *SubscribeAlterUpdate) SetLastModifier(m *model.Modifier) *SubscribeAlterUpdate {
	sau.mutation.SetLastModifier(m)
	return sau
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (sau *SubscribeAlterUpdate) ClearLastModifier() *SubscribeAlterUpdate {
	sau.mutation.ClearLastModifier()
	return sau
}

// SetRemark sets the "remark" field.
func (sau *SubscribeAlterUpdate) SetRemark(s string) *SubscribeAlterUpdate {
	sau.mutation.SetRemark(s)
	return sau
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (sau *SubscribeAlterUpdate) SetNillableRemark(s *string) *SubscribeAlterUpdate {
	if s != nil {
		sau.SetRemark(*s)
	}
	return sau
}

// ClearRemark clears the value of the "remark" field.
func (sau *SubscribeAlterUpdate) ClearRemark() *SubscribeAlterUpdate {
	sau.mutation.ClearRemark()
	return sau
}

// SetRiderID sets the "rider_id" field.
func (sau *SubscribeAlterUpdate) SetRiderID(u uint64) *SubscribeAlterUpdate {
	sau.mutation.SetRiderID(u)
	return sau
}

// SetManagerID sets the "manager_id" field.
func (sau *SubscribeAlterUpdate) SetManagerID(u uint64) *SubscribeAlterUpdate {
	sau.mutation.SetManagerID(u)
	return sau
}

// SetSubscribeID sets the "subscribe_id" field.
func (sau *SubscribeAlterUpdate) SetSubscribeID(u uint64) *SubscribeAlterUpdate {
	sau.mutation.SetSubscribeID(u)
	return sau
}

// SetDays sets the "days" field.
func (sau *SubscribeAlterUpdate) SetDays(i int) *SubscribeAlterUpdate {
	sau.mutation.ResetDays()
	sau.mutation.SetDays(i)
	return sau
}

// AddDays adds i to the "days" field.
func (sau *SubscribeAlterUpdate) AddDays(i int) *SubscribeAlterUpdate {
	sau.mutation.AddDays(i)
	return sau
}

// SetReason sets the "reason" field.
func (sau *SubscribeAlterUpdate) SetReason(s string) *SubscribeAlterUpdate {
	sau.mutation.SetReason(s)
	return sau
}

// SetRider sets the "rider" edge to the Rider entity.
func (sau *SubscribeAlterUpdate) SetRider(r *Rider) *SubscribeAlterUpdate {
	return sau.SetRiderID(r.ID)
}

// SetManager sets the "manager" edge to the Manager entity.
func (sau *SubscribeAlterUpdate) SetManager(m *Manager) *SubscribeAlterUpdate {
	return sau.SetManagerID(m.ID)
}

// SetSubscribe sets the "subscribe" edge to the Subscribe entity.
func (sau *SubscribeAlterUpdate) SetSubscribe(s *Subscribe) *SubscribeAlterUpdate {
	return sau.SetSubscribeID(s.ID)
}

// Mutation returns the SubscribeAlterMutation object of the builder.
func (sau *SubscribeAlterUpdate) Mutation() *SubscribeAlterMutation {
	return sau.mutation
}

// ClearRider clears the "rider" edge to the Rider entity.
func (sau *SubscribeAlterUpdate) ClearRider() *SubscribeAlterUpdate {
	sau.mutation.ClearRider()
	return sau
}

// ClearManager clears the "manager" edge to the Manager entity.
func (sau *SubscribeAlterUpdate) ClearManager() *SubscribeAlterUpdate {
	sau.mutation.ClearManager()
	return sau
}

// ClearSubscribe clears the "subscribe" edge to the Subscribe entity.
func (sau *SubscribeAlterUpdate) ClearSubscribe() *SubscribeAlterUpdate {
	sau.mutation.ClearSubscribe()
	return sau
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (sau *SubscribeAlterUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if err := sau.defaults(); err != nil {
		return 0, err
	}
	if len(sau.hooks) == 0 {
		if err = sau.check(); err != nil {
			return 0, err
		}
		affected, err = sau.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SubscribeAlterMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = sau.check(); err != nil {
				return 0, err
			}
			sau.mutation = mutation
			affected, err = sau.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(sau.hooks) - 1; i >= 0; i-- {
			if sau.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = sau.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sau.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (sau *SubscribeAlterUpdate) SaveX(ctx context.Context) int {
	affected, err := sau.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (sau *SubscribeAlterUpdate) Exec(ctx context.Context) error {
	_, err := sau.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sau *SubscribeAlterUpdate) ExecX(ctx context.Context) {
	if err := sau.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sau *SubscribeAlterUpdate) defaults() error {
	if _, ok := sau.mutation.UpdatedAt(); !ok {
		if subscribealter.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized subscribealter.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := subscribealter.UpdateDefaultUpdatedAt()
		sau.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (sau *SubscribeAlterUpdate) check() error {
	if _, ok := sau.mutation.RiderID(); sau.mutation.RiderCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SubscribeAlter.rider"`)
	}
	if _, ok := sau.mutation.ManagerID(); sau.mutation.ManagerCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SubscribeAlter.manager"`)
	}
	if _, ok := sau.mutation.SubscribeID(); sau.mutation.SubscribeCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SubscribeAlter.subscribe"`)
	}
	return nil
}

func (sau *SubscribeAlterUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   subscribealter.Table,
			Columns: subscribealter.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: subscribealter.FieldID,
			},
		},
	}
	if ps := sau.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sau.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribealter.FieldUpdatedAt,
		})
	}
	if value, ok := sau.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribealter.FieldDeletedAt,
		})
	}
	if sau.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: subscribealter.FieldDeletedAt,
		})
	}
	if sau.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: subscribealter.FieldCreator,
		})
	}
	if value, ok := sau.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: subscribealter.FieldLastModifier,
		})
	}
	if sau.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: subscribealter.FieldLastModifier,
		})
	}
	if value, ok := sau.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: subscribealter.FieldRemark,
		})
	}
	if sau.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: subscribealter.FieldRemark,
		})
	}
	if value, ok := sau.mutation.Days(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: subscribealter.FieldDays,
		})
	}
	if value, ok := sau.mutation.AddedDays(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: subscribealter.FieldDays,
		})
	}
	if value, ok := sau.mutation.Reason(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: subscribealter.FieldReason,
		})
	}
	if sau.mutation.RiderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribealter.RiderTable,
			Columns: []string{subscribealter.RiderColumn},
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
	if nodes := sau.mutation.RiderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribealter.RiderTable,
			Columns: []string{subscribealter.RiderColumn},
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
	if sau.mutation.ManagerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribealter.ManagerTable,
			Columns: []string{subscribealter.ManagerColumn},
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
	if nodes := sau.mutation.ManagerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribealter.ManagerTable,
			Columns: []string{subscribealter.ManagerColumn},
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
	if sau.mutation.SubscribeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   subscribealter.SubscribeTable,
			Columns: []string{subscribealter.SubscribeColumn},
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
	if nodes := sau.mutation.SubscribeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   subscribealter.SubscribeTable,
			Columns: []string{subscribealter.SubscribeColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, sau.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{subscribealter.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// SubscribeAlterUpdateOne is the builder for updating a single SubscribeAlter entity.
type SubscribeAlterUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SubscribeAlterMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (sauo *SubscribeAlterUpdateOne) SetUpdatedAt(t time.Time) *SubscribeAlterUpdateOne {
	sauo.mutation.SetUpdatedAt(t)
	return sauo
}

// SetDeletedAt sets the "deleted_at" field.
func (sauo *SubscribeAlterUpdateOne) SetDeletedAt(t time.Time) *SubscribeAlterUpdateOne {
	sauo.mutation.SetDeletedAt(t)
	return sauo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (sauo *SubscribeAlterUpdateOne) SetNillableDeletedAt(t *time.Time) *SubscribeAlterUpdateOne {
	if t != nil {
		sauo.SetDeletedAt(*t)
	}
	return sauo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (sauo *SubscribeAlterUpdateOne) ClearDeletedAt() *SubscribeAlterUpdateOne {
	sauo.mutation.ClearDeletedAt()
	return sauo
}

// SetLastModifier sets the "last_modifier" field.
func (sauo *SubscribeAlterUpdateOne) SetLastModifier(m *model.Modifier) *SubscribeAlterUpdateOne {
	sauo.mutation.SetLastModifier(m)
	return sauo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (sauo *SubscribeAlterUpdateOne) ClearLastModifier() *SubscribeAlterUpdateOne {
	sauo.mutation.ClearLastModifier()
	return sauo
}

// SetRemark sets the "remark" field.
func (sauo *SubscribeAlterUpdateOne) SetRemark(s string) *SubscribeAlterUpdateOne {
	sauo.mutation.SetRemark(s)
	return sauo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (sauo *SubscribeAlterUpdateOne) SetNillableRemark(s *string) *SubscribeAlterUpdateOne {
	if s != nil {
		sauo.SetRemark(*s)
	}
	return sauo
}

// ClearRemark clears the value of the "remark" field.
func (sauo *SubscribeAlterUpdateOne) ClearRemark() *SubscribeAlterUpdateOne {
	sauo.mutation.ClearRemark()
	return sauo
}

// SetRiderID sets the "rider_id" field.
func (sauo *SubscribeAlterUpdateOne) SetRiderID(u uint64) *SubscribeAlterUpdateOne {
	sauo.mutation.SetRiderID(u)
	return sauo
}

// SetManagerID sets the "manager_id" field.
func (sauo *SubscribeAlterUpdateOne) SetManagerID(u uint64) *SubscribeAlterUpdateOne {
	sauo.mutation.SetManagerID(u)
	return sauo
}

// SetSubscribeID sets the "subscribe_id" field.
func (sauo *SubscribeAlterUpdateOne) SetSubscribeID(u uint64) *SubscribeAlterUpdateOne {
	sauo.mutation.SetSubscribeID(u)
	return sauo
}

// SetDays sets the "days" field.
func (sauo *SubscribeAlterUpdateOne) SetDays(i int) *SubscribeAlterUpdateOne {
	sauo.mutation.ResetDays()
	sauo.mutation.SetDays(i)
	return sauo
}

// AddDays adds i to the "days" field.
func (sauo *SubscribeAlterUpdateOne) AddDays(i int) *SubscribeAlterUpdateOne {
	sauo.mutation.AddDays(i)
	return sauo
}

// SetReason sets the "reason" field.
func (sauo *SubscribeAlterUpdateOne) SetReason(s string) *SubscribeAlterUpdateOne {
	sauo.mutation.SetReason(s)
	return sauo
}

// SetRider sets the "rider" edge to the Rider entity.
func (sauo *SubscribeAlterUpdateOne) SetRider(r *Rider) *SubscribeAlterUpdateOne {
	return sauo.SetRiderID(r.ID)
}

// SetManager sets the "manager" edge to the Manager entity.
func (sauo *SubscribeAlterUpdateOne) SetManager(m *Manager) *SubscribeAlterUpdateOne {
	return sauo.SetManagerID(m.ID)
}

// SetSubscribe sets the "subscribe" edge to the Subscribe entity.
func (sauo *SubscribeAlterUpdateOne) SetSubscribe(s *Subscribe) *SubscribeAlterUpdateOne {
	return sauo.SetSubscribeID(s.ID)
}

// Mutation returns the SubscribeAlterMutation object of the builder.
func (sauo *SubscribeAlterUpdateOne) Mutation() *SubscribeAlterMutation {
	return sauo.mutation
}

// ClearRider clears the "rider" edge to the Rider entity.
func (sauo *SubscribeAlterUpdateOne) ClearRider() *SubscribeAlterUpdateOne {
	sauo.mutation.ClearRider()
	return sauo
}

// ClearManager clears the "manager" edge to the Manager entity.
func (sauo *SubscribeAlterUpdateOne) ClearManager() *SubscribeAlterUpdateOne {
	sauo.mutation.ClearManager()
	return sauo
}

// ClearSubscribe clears the "subscribe" edge to the Subscribe entity.
func (sauo *SubscribeAlterUpdateOne) ClearSubscribe() *SubscribeAlterUpdateOne {
	sauo.mutation.ClearSubscribe()
	return sauo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (sauo *SubscribeAlterUpdateOne) Select(field string, fields ...string) *SubscribeAlterUpdateOne {
	sauo.fields = append([]string{field}, fields...)
	return sauo
}

// Save executes the query and returns the updated SubscribeAlter entity.
func (sauo *SubscribeAlterUpdateOne) Save(ctx context.Context) (*SubscribeAlter, error) {
	var (
		err  error
		node *SubscribeAlter
	)
	if err := sauo.defaults(); err != nil {
		return nil, err
	}
	if len(sauo.hooks) == 0 {
		if err = sauo.check(); err != nil {
			return nil, err
		}
		node, err = sauo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SubscribeAlterMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = sauo.check(); err != nil {
				return nil, err
			}
			sauo.mutation = mutation
			node, err = sauo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(sauo.hooks) - 1; i >= 0; i-- {
			if sauo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = sauo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, sauo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*SubscribeAlter)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from SubscribeAlterMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (sauo *SubscribeAlterUpdateOne) SaveX(ctx context.Context) *SubscribeAlter {
	node, err := sauo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (sauo *SubscribeAlterUpdateOne) Exec(ctx context.Context) error {
	_, err := sauo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sauo *SubscribeAlterUpdateOne) ExecX(ctx context.Context) {
	if err := sauo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sauo *SubscribeAlterUpdateOne) defaults() error {
	if _, ok := sauo.mutation.UpdatedAt(); !ok {
		if subscribealter.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized subscribealter.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := subscribealter.UpdateDefaultUpdatedAt()
		sauo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (sauo *SubscribeAlterUpdateOne) check() error {
	if _, ok := sauo.mutation.RiderID(); sauo.mutation.RiderCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SubscribeAlter.rider"`)
	}
	if _, ok := sauo.mutation.ManagerID(); sauo.mutation.ManagerCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SubscribeAlter.manager"`)
	}
	if _, ok := sauo.mutation.SubscribeID(); sauo.mutation.SubscribeCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "SubscribeAlter.subscribe"`)
	}
	return nil
}

func (sauo *SubscribeAlterUpdateOne) sqlSave(ctx context.Context) (_node *SubscribeAlter, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   subscribealter.Table,
			Columns: subscribealter.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: subscribealter.FieldID,
			},
		},
	}
	id, ok := sauo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "SubscribeAlter.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := sauo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, subscribealter.FieldID)
		for _, f := range fields {
			if !subscribealter.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != subscribealter.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := sauo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sauo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribealter.FieldUpdatedAt,
		})
	}
	if value, ok := sauo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscribealter.FieldDeletedAt,
		})
	}
	if sauo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: subscribealter.FieldDeletedAt,
		})
	}
	if sauo.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: subscribealter.FieldCreator,
		})
	}
	if value, ok := sauo.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: subscribealter.FieldLastModifier,
		})
	}
	if sauo.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: subscribealter.FieldLastModifier,
		})
	}
	if value, ok := sauo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: subscribealter.FieldRemark,
		})
	}
	if sauo.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: subscribealter.FieldRemark,
		})
	}
	if value, ok := sauo.mutation.Days(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: subscribealter.FieldDays,
		})
	}
	if value, ok := sauo.mutation.AddedDays(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: subscribealter.FieldDays,
		})
	}
	if value, ok := sauo.mutation.Reason(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: subscribealter.FieldReason,
		})
	}
	if sauo.mutation.RiderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribealter.RiderTable,
			Columns: []string{subscribealter.RiderColumn},
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
	if nodes := sauo.mutation.RiderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribealter.RiderTable,
			Columns: []string{subscribealter.RiderColumn},
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
	if sauo.mutation.ManagerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribealter.ManagerTable,
			Columns: []string{subscribealter.ManagerColumn},
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
	if nodes := sauo.mutation.ManagerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscribealter.ManagerTable,
			Columns: []string{subscribealter.ManagerColumn},
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
	if sauo.mutation.SubscribeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   subscribealter.SubscribeTable,
			Columns: []string{subscribealter.SubscribeColumn},
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
	if nodes := sauo.mutation.SubscribeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   subscribealter.SubscribeTable,
			Columns: []string{subscribealter.SubscribeColumn},
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
	_node = &SubscribeAlter{config: sauo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, sauo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{subscribealter.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}