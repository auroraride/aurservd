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
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/statement"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
)

// StatementUpdate is the builder for updating Statement entities.
type StatementUpdate struct {
	config
	hooks    []Hook
	mutation *StatementMutation
}

// Where appends a list predicates to the StatementUpdate builder.
func (su *StatementUpdate) Where(ps ...predicate.Statement) *StatementUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetUpdatedAt sets the "updated_at" field.
func (su *StatementUpdate) SetUpdatedAt(t time.Time) *StatementUpdate {
	su.mutation.SetUpdatedAt(t)
	return su
}

// SetDeletedAt sets the "deleted_at" field.
func (su *StatementUpdate) SetDeletedAt(t time.Time) *StatementUpdate {
	su.mutation.SetDeletedAt(t)
	return su
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (su *StatementUpdate) SetNillableDeletedAt(t *time.Time) *StatementUpdate {
	if t != nil {
		su.SetDeletedAt(*t)
	}
	return su
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (su *StatementUpdate) ClearDeletedAt() *StatementUpdate {
	su.mutation.ClearDeletedAt()
	return su
}

// SetLastModifier sets the "last_modifier" field.
func (su *StatementUpdate) SetLastModifier(m *model.Modifier) *StatementUpdate {
	su.mutation.SetLastModifier(m)
	return su
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (su *StatementUpdate) ClearLastModifier() *StatementUpdate {
	su.mutation.ClearLastModifier()
	return su
}

// SetRemark sets the "remark" field.
func (su *StatementUpdate) SetRemark(s string) *StatementUpdate {
	su.mutation.SetRemark(s)
	return su
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (su *StatementUpdate) SetNillableRemark(s *string) *StatementUpdate {
	if s != nil {
		su.SetRemark(*s)
	}
	return su
}

// ClearRemark clears the value of the "remark" field.
func (su *StatementUpdate) ClearRemark() *StatementUpdate {
	su.mutation.ClearRemark()
	return su
}

// SetEnterpriseID sets the "enterprise_id" field.
func (su *StatementUpdate) SetEnterpriseID(u uint64) *StatementUpdate {
	su.mutation.SetEnterpriseID(u)
	return su
}

// SetCost sets the "cost" field.
func (su *StatementUpdate) SetCost(f float64) *StatementUpdate {
	su.mutation.ResetCost()
	su.mutation.SetCost(f)
	return su
}

// SetNillableCost sets the "cost" field if the given value is not nil.
func (su *StatementUpdate) SetNillableCost(f *float64) *StatementUpdate {
	if f != nil {
		su.SetCost(*f)
	}
	return su
}

// AddCost adds f to the "cost" field.
func (su *StatementUpdate) AddCost(f float64) *StatementUpdate {
	su.mutation.AddCost(f)
	return su
}

// SetAmount sets the "amount" field.
func (su *StatementUpdate) SetAmount(f float64) *StatementUpdate {
	su.mutation.ResetAmount()
	su.mutation.SetAmount(f)
	return su
}

// SetNillableAmount sets the "amount" field if the given value is not nil.
func (su *StatementUpdate) SetNillableAmount(f *float64) *StatementUpdate {
	if f != nil {
		su.SetAmount(*f)
	}
	return su
}

// AddAmount adds f to the "amount" field.
func (su *StatementUpdate) AddAmount(f float64) *StatementUpdate {
	su.mutation.AddAmount(f)
	return su
}

// SetBalance sets the "balance" field.
func (su *StatementUpdate) SetBalance(f float64) *StatementUpdate {
	su.mutation.ResetBalance()
	su.mutation.SetBalance(f)
	return su
}

// SetNillableBalance sets the "balance" field if the given value is not nil.
func (su *StatementUpdate) SetNillableBalance(f *float64) *StatementUpdate {
	if f != nil {
		su.SetBalance(*f)
	}
	return su
}

// AddBalance adds f to the "balance" field.
func (su *StatementUpdate) AddBalance(f float64) *StatementUpdate {
	su.mutation.AddBalance(f)
	return su
}

// SetSettledAt sets the "settled_at" field.
func (su *StatementUpdate) SetSettledAt(t time.Time) *StatementUpdate {
	su.mutation.SetSettledAt(t)
	return su
}

// SetNillableSettledAt sets the "settled_at" field if the given value is not nil.
func (su *StatementUpdate) SetNillableSettledAt(t *time.Time) *StatementUpdate {
	if t != nil {
		su.SetSettledAt(*t)
	}
	return su
}

// ClearSettledAt clears the value of the "settled_at" field.
func (su *StatementUpdate) ClearSettledAt() *StatementUpdate {
	su.mutation.ClearSettledAt()
	return su
}

// SetDays sets the "days" field.
func (su *StatementUpdate) SetDays(i int) *StatementUpdate {
	su.mutation.ResetDays()
	su.mutation.SetDays(i)
	return su
}

// SetNillableDays sets the "days" field if the given value is not nil.
func (su *StatementUpdate) SetNillableDays(i *int) *StatementUpdate {
	if i != nil {
		su.SetDays(*i)
	}
	return su
}

// AddDays adds i to the "days" field.
func (su *StatementUpdate) AddDays(i int) *StatementUpdate {
	su.mutation.AddDays(i)
	return su
}

// SetRiderNumber sets the "rider_number" field.
func (su *StatementUpdate) SetRiderNumber(i int) *StatementUpdate {
	su.mutation.ResetRiderNumber()
	su.mutation.SetRiderNumber(i)
	return su
}

// SetNillableRiderNumber sets the "rider_number" field if the given value is not nil.
func (su *StatementUpdate) SetNillableRiderNumber(i *int) *StatementUpdate {
	if i != nil {
		su.SetRiderNumber(*i)
	}
	return su
}

// AddRiderNumber adds i to the "rider_number" field.
func (su *StatementUpdate) AddRiderNumber(i int) *StatementUpdate {
	su.mutation.AddRiderNumber(i)
	return su
}

// AddSubscribeIDs adds the "subscribes" edge to the Subscribe entity by IDs.
func (su *StatementUpdate) AddSubscribeIDs(ids ...uint64) *StatementUpdate {
	su.mutation.AddSubscribeIDs(ids...)
	return su
}

// AddSubscribes adds the "subscribes" edges to the Subscribe entity.
func (su *StatementUpdate) AddSubscribes(s ...*Subscribe) *StatementUpdate {
	ids := make([]uint64, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddSubscribeIDs(ids...)
}

// SetEnterprise sets the "enterprise" edge to the Enterprise entity.
func (su *StatementUpdate) SetEnterprise(e *Enterprise) *StatementUpdate {
	return su.SetEnterpriseID(e.ID)
}

// Mutation returns the StatementMutation object of the builder.
func (su *StatementUpdate) Mutation() *StatementMutation {
	return su.mutation
}

// ClearSubscribes clears all "subscribes" edges to the Subscribe entity.
func (su *StatementUpdate) ClearSubscribes() *StatementUpdate {
	su.mutation.ClearSubscribes()
	return su
}

// RemoveSubscribeIDs removes the "subscribes" edge to Subscribe entities by IDs.
func (su *StatementUpdate) RemoveSubscribeIDs(ids ...uint64) *StatementUpdate {
	su.mutation.RemoveSubscribeIDs(ids...)
	return su
}

// RemoveSubscribes removes "subscribes" edges to Subscribe entities.
func (su *StatementUpdate) RemoveSubscribes(s ...*Subscribe) *StatementUpdate {
	ids := make([]uint64, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveSubscribeIDs(ids...)
}

// ClearEnterprise clears the "enterprise" edge to the Enterprise entity.
func (su *StatementUpdate) ClearEnterprise() *StatementUpdate {
	su.mutation.ClearEnterprise()
	return su
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *StatementUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if err := su.defaults(); err != nil {
		return 0, err
	}
	if len(su.hooks) == 0 {
		if err = su.check(); err != nil {
			return 0, err
		}
		affected, err = su.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*StatementMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = su.check(); err != nil {
				return 0, err
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
func (su *StatementUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *StatementUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *StatementUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (su *StatementUpdate) defaults() error {
	if _, ok := su.mutation.UpdatedAt(); !ok {
		if statement.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized statement.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := statement.UpdateDefaultUpdatedAt()
		su.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (su *StatementUpdate) check() error {
	if _, ok := su.mutation.EnterpriseID(); su.mutation.EnterpriseCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Statement.enterprise"`)
	}
	return nil
}

func (su *StatementUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   statement.Table,
			Columns: statement.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: statement.FieldID,
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
			Column: statement.FieldUpdatedAt,
		})
	}
	if value, ok := su.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: statement.FieldDeletedAt,
		})
	}
	if su.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: statement.FieldDeletedAt,
		})
	}
	if su.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: statement.FieldCreator,
		})
	}
	if value, ok := su.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: statement.FieldLastModifier,
		})
	}
	if su.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: statement.FieldLastModifier,
		})
	}
	if value, ok := su.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: statement.FieldRemark,
		})
	}
	if su.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: statement.FieldRemark,
		})
	}
	if value, ok := su.mutation.Cost(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldCost,
		})
	}
	if value, ok := su.mutation.AddedCost(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldCost,
		})
	}
	if value, ok := su.mutation.Amount(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldAmount,
		})
	}
	if value, ok := su.mutation.AddedAmount(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldAmount,
		})
	}
	if value, ok := su.mutation.Balance(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldBalance,
		})
	}
	if value, ok := su.mutation.AddedBalance(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldBalance,
		})
	}
	if value, ok := su.mutation.SettledAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: statement.FieldSettledAt,
		})
	}
	if su.mutation.SettledAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: statement.FieldSettledAt,
		})
	}
	if value, ok := su.mutation.Days(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: statement.FieldDays,
		})
	}
	if value, ok := su.mutation.AddedDays(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: statement.FieldDays,
		})
	}
	if value, ok := su.mutation.RiderNumber(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: statement.FieldRiderNumber,
		})
	}
	if value, ok := su.mutation.AddedRiderNumber(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: statement.FieldRiderNumber,
		})
	}
	if su.mutation.SubscribesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubscribesTable,
			Columns: []string{statement.SubscribesColumn},
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
	if nodes := su.mutation.RemovedSubscribesIDs(); len(nodes) > 0 && !su.mutation.SubscribesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubscribesTable,
			Columns: []string{statement.SubscribesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.SubscribesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubscribesTable,
			Columns: []string{statement.SubscribesColumn},
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
	if su.mutation.EnterpriseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   statement.EnterpriseTable,
			Columns: []string{statement.EnterpriseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: enterprise.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.EnterpriseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   statement.EnterpriseTable,
			Columns: []string{statement.EnterpriseColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{statement.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// StatementUpdateOne is the builder for updating a single Statement entity.
type StatementUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *StatementMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (suo *StatementUpdateOne) SetUpdatedAt(t time.Time) *StatementUpdateOne {
	suo.mutation.SetUpdatedAt(t)
	return suo
}

// SetDeletedAt sets the "deleted_at" field.
func (suo *StatementUpdateOne) SetDeletedAt(t time.Time) *StatementUpdateOne {
	suo.mutation.SetDeletedAt(t)
	return suo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (suo *StatementUpdateOne) SetNillableDeletedAt(t *time.Time) *StatementUpdateOne {
	if t != nil {
		suo.SetDeletedAt(*t)
	}
	return suo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (suo *StatementUpdateOne) ClearDeletedAt() *StatementUpdateOne {
	suo.mutation.ClearDeletedAt()
	return suo
}

// SetLastModifier sets the "last_modifier" field.
func (suo *StatementUpdateOne) SetLastModifier(m *model.Modifier) *StatementUpdateOne {
	suo.mutation.SetLastModifier(m)
	return suo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (suo *StatementUpdateOne) ClearLastModifier() *StatementUpdateOne {
	suo.mutation.ClearLastModifier()
	return suo
}

// SetRemark sets the "remark" field.
func (suo *StatementUpdateOne) SetRemark(s string) *StatementUpdateOne {
	suo.mutation.SetRemark(s)
	return suo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (suo *StatementUpdateOne) SetNillableRemark(s *string) *StatementUpdateOne {
	if s != nil {
		suo.SetRemark(*s)
	}
	return suo
}

// ClearRemark clears the value of the "remark" field.
func (suo *StatementUpdateOne) ClearRemark() *StatementUpdateOne {
	suo.mutation.ClearRemark()
	return suo
}

// SetEnterpriseID sets the "enterprise_id" field.
func (suo *StatementUpdateOne) SetEnterpriseID(u uint64) *StatementUpdateOne {
	suo.mutation.SetEnterpriseID(u)
	return suo
}

// SetCost sets the "cost" field.
func (suo *StatementUpdateOne) SetCost(f float64) *StatementUpdateOne {
	suo.mutation.ResetCost()
	suo.mutation.SetCost(f)
	return suo
}

// SetNillableCost sets the "cost" field if the given value is not nil.
func (suo *StatementUpdateOne) SetNillableCost(f *float64) *StatementUpdateOne {
	if f != nil {
		suo.SetCost(*f)
	}
	return suo
}

// AddCost adds f to the "cost" field.
func (suo *StatementUpdateOne) AddCost(f float64) *StatementUpdateOne {
	suo.mutation.AddCost(f)
	return suo
}

// SetAmount sets the "amount" field.
func (suo *StatementUpdateOne) SetAmount(f float64) *StatementUpdateOne {
	suo.mutation.ResetAmount()
	suo.mutation.SetAmount(f)
	return suo
}

// SetNillableAmount sets the "amount" field if the given value is not nil.
func (suo *StatementUpdateOne) SetNillableAmount(f *float64) *StatementUpdateOne {
	if f != nil {
		suo.SetAmount(*f)
	}
	return suo
}

// AddAmount adds f to the "amount" field.
func (suo *StatementUpdateOne) AddAmount(f float64) *StatementUpdateOne {
	suo.mutation.AddAmount(f)
	return suo
}

// SetBalance sets the "balance" field.
func (suo *StatementUpdateOne) SetBalance(f float64) *StatementUpdateOne {
	suo.mutation.ResetBalance()
	suo.mutation.SetBalance(f)
	return suo
}

// SetNillableBalance sets the "balance" field if the given value is not nil.
func (suo *StatementUpdateOne) SetNillableBalance(f *float64) *StatementUpdateOne {
	if f != nil {
		suo.SetBalance(*f)
	}
	return suo
}

// AddBalance adds f to the "balance" field.
func (suo *StatementUpdateOne) AddBalance(f float64) *StatementUpdateOne {
	suo.mutation.AddBalance(f)
	return suo
}

// SetSettledAt sets the "settled_at" field.
func (suo *StatementUpdateOne) SetSettledAt(t time.Time) *StatementUpdateOne {
	suo.mutation.SetSettledAt(t)
	return suo
}

// SetNillableSettledAt sets the "settled_at" field if the given value is not nil.
func (suo *StatementUpdateOne) SetNillableSettledAt(t *time.Time) *StatementUpdateOne {
	if t != nil {
		suo.SetSettledAt(*t)
	}
	return suo
}

// ClearSettledAt clears the value of the "settled_at" field.
func (suo *StatementUpdateOne) ClearSettledAt() *StatementUpdateOne {
	suo.mutation.ClearSettledAt()
	return suo
}

// SetDays sets the "days" field.
func (suo *StatementUpdateOne) SetDays(i int) *StatementUpdateOne {
	suo.mutation.ResetDays()
	suo.mutation.SetDays(i)
	return suo
}

// SetNillableDays sets the "days" field if the given value is not nil.
func (suo *StatementUpdateOne) SetNillableDays(i *int) *StatementUpdateOne {
	if i != nil {
		suo.SetDays(*i)
	}
	return suo
}

// AddDays adds i to the "days" field.
func (suo *StatementUpdateOne) AddDays(i int) *StatementUpdateOne {
	suo.mutation.AddDays(i)
	return suo
}

// SetRiderNumber sets the "rider_number" field.
func (suo *StatementUpdateOne) SetRiderNumber(i int) *StatementUpdateOne {
	suo.mutation.ResetRiderNumber()
	suo.mutation.SetRiderNumber(i)
	return suo
}

// SetNillableRiderNumber sets the "rider_number" field if the given value is not nil.
func (suo *StatementUpdateOne) SetNillableRiderNumber(i *int) *StatementUpdateOne {
	if i != nil {
		suo.SetRiderNumber(*i)
	}
	return suo
}

// AddRiderNumber adds i to the "rider_number" field.
func (suo *StatementUpdateOne) AddRiderNumber(i int) *StatementUpdateOne {
	suo.mutation.AddRiderNumber(i)
	return suo
}

// AddSubscribeIDs adds the "subscribes" edge to the Subscribe entity by IDs.
func (suo *StatementUpdateOne) AddSubscribeIDs(ids ...uint64) *StatementUpdateOne {
	suo.mutation.AddSubscribeIDs(ids...)
	return suo
}

// AddSubscribes adds the "subscribes" edges to the Subscribe entity.
func (suo *StatementUpdateOne) AddSubscribes(s ...*Subscribe) *StatementUpdateOne {
	ids := make([]uint64, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddSubscribeIDs(ids...)
}

// SetEnterprise sets the "enterprise" edge to the Enterprise entity.
func (suo *StatementUpdateOne) SetEnterprise(e *Enterprise) *StatementUpdateOne {
	return suo.SetEnterpriseID(e.ID)
}

// Mutation returns the StatementMutation object of the builder.
func (suo *StatementUpdateOne) Mutation() *StatementMutation {
	return suo.mutation
}

// ClearSubscribes clears all "subscribes" edges to the Subscribe entity.
func (suo *StatementUpdateOne) ClearSubscribes() *StatementUpdateOne {
	suo.mutation.ClearSubscribes()
	return suo
}

// RemoveSubscribeIDs removes the "subscribes" edge to Subscribe entities by IDs.
func (suo *StatementUpdateOne) RemoveSubscribeIDs(ids ...uint64) *StatementUpdateOne {
	suo.mutation.RemoveSubscribeIDs(ids...)
	return suo
}

// RemoveSubscribes removes "subscribes" edges to Subscribe entities.
func (suo *StatementUpdateOne) RemoveSubscribes(s ...*Subscribe) *StatementUpdateOne {
	ids := make([]uint64, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveSubscribeIDs(ids...)
}

// ClearEnterprise clears the "enterprise" edge to the Enterprise entity.
func (suo *StatementUpdateOne) ClearEnterprise() *StatementUpdateOne {
	suo.mutation.ClearEnterprise()
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *StatementUpdateOne) Select(field string, fields ...string) *StatementUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Statement entity.
func (suo *StatementUpdateOne) Save(ctx context.Context) (*Statement, error) {
	var (
		err  error
		node *Statement
	)
	if err := suo.defaults(); err != nil {
		return nil, err
	}
	if len(suo.hooks) == 0 {
		if err = suo.check(); err != nil {
			return nil, err
		}
		node, err = suo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*StatementMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = suo.check(); err != nil {
				return nil, err
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
		nv, ok := v.(*Statement)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from StatementMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (suo *StatementUpdateOne) SaveX(ctx context.Context) *Statement {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *StatementUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *StatementUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (suo *StatementUpdateOne) defaults() error {
	if _, ok := suo.mutation.UpdatedAt(); !ok {
		if statement.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized statement.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := statement.UpdateDefaultUpdatedAt()
		suo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (suo *StatementUpdateOne) check() error {
	if _, ok := suo.mutation.EnterpriseID(); suo.mutation.EnterpriseCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Statement.enterprise"`)
	}
	return nil
}

func (suo *StatementUpdateOne) sqlSave(ctx context.Context) (_node *Statement, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   statement.Table,
			Columns: statement.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: statement.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Statement.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, statement.FieldID)
		for _, f := range fields {
			if !statement.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != statement.FieldID {
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
			Column: statement.FieldUpdatedAt,
		})
	}
	if value, ok := suo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: statement.FieldDeletedAt,
		})
	}
	if suo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: statement.FieldDeletedAt,
		})
	}
	if suo.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: statement.FieldCreator,
		})
	}
	if value, ok := suo.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: statement.FieldLastModifier,
		})
	}
	if suo.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: statement.FieldLastModifier,
		})
	}
	if value, ok := suo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: statement.FieldRemark,
		})
	}
	if suo.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: statement.FieldRemark,
		})
	}
	if value, ok := suo.mutation.Cost(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldCost,
		})
	}
	if value, ok := suo.mutation.AddedCost(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldCost,
		})
	}
	if value, ok := suo.mutation.Amount(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldAmount,
		})
	}
	if value, ok := suo.mutation.AddedAmount(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldAmount,
		})
	}
	if value, ok := suo.mutation.Balance(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldBalance,
		})
	}
	if value, ok := suo.mutation.AddedBalance(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: statement.FieldBalance,
		})
	}
	if value, ok := suo.mutation.SettledAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: statement.FieldSettledAt,
		})
	}
	if suo.mutation.SettledAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: statement.FieldSettledAt,
		})
	}
	if value, ok := suo.mutation.Days(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: statement.FieldDays,
		})
	}
	if value, ok := suo.mutation.AddedDays(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: statement.FieldDays,
		})
	}
	if value, ok := suo.mutation.RiderNumber(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: statement.FieldRiderNumber,
		})
	}
	if value, ok := suo.mutation.AddedRiderNumber(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: statement.FieldRiderNumber,
		})
	}
	if suo.mutation.SubscribesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubscribesTable,
			Columns: []string{statement.SubscribesColumn},
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
	if nodes := suo.mutation.RemovedSubscribesIDs(); len(nodes) > 0 && !suo.mutation.SubscribesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubscribesTable,
			Columns: []string{statement.SubscribesColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.SubscribesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubscribesTable,
			Columns: []string{statement.SubscribesColumn},
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
	if suo.mutation.EnterpriseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   statement.EnterpriseTable,
			Columns: []string{statement.EnterpriseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: enterprise.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.EnterpriseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   statement.EnterpriseTable,
			Columns: []string{statement.EnterpriseColumn},
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Statement{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{statement.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}