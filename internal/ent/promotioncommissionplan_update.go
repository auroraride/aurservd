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
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/promotioncommission"
	"github.com/auroraride/aurservd/internal/ent/promotioncommissionplan"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
)

// PromotionCommissionPlanUpdate is the builder for updating PromotionCommissionPlan entities.
type PromotionCommissionPlanUpdate struct {
	config
	hooks     []Hook
	mutation  *PromotionCommissionPlanMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the PromotionCommissionPlanUpdate builder.
func (pcpu *PromotionCommissionPlanUpdate) Where(ps ...predicate.PromotionCommissionPlan) *PromotionCommissionPlanUpdate {
	pcpu.mutation.Where(ps...)
	return pcpu
}

// SetUpdatedAt sets the "updated_at" field.
func (pcpu *PromotionCommissionPlanUpdate) SetUpdatedAt(t time.Time) *PromotionCommissionPlanUpdate {
	pcpu.mutation.SetUpdatedAt(t)
	return pcpu
}

// SetDeletedAt sets the "deleted_at" field.
func (pcpu *PromotionCommissionPlanUpdate) SetDeletedAt(t time.Time) *PromotionCommissionPlanUpdate {
	pcpu.mutation.SetDeletedAt(t)
	return pcpu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (pcpu *PromotionCommissionPlanUpdate) SetNillableDeletedAt(t *time.Time) *PromotionCommissionPlanUpdate {
	if t != nil {
		pcpu.SetDeletedAt(*t)
	}
	return pcpu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (pcpu *PromotionCommissionPlanUpdate) ClearDeletedAt() *PromotionCommissionPlanUpdate {
	pcpu.mutation.ClearDeletedAt()
	return pcpu
}

// SetMemberID sets the "member_id" field.
func (pcpu *PromotionCommissionPlanUpdate) SetMemberID(u uint64) *PromotionCommissionPlanUpdate {
	pcpu.mutation.SetMemberID(u)
	return pcpu
}

// SetNillableMemberID sets the "member_id" field if the given value is not nil.
func (pcpu *PromotionCommissionPlanUpdate) SetNillableMemberID(u *uint64) *PromotionCommissionPlanUpdate {
	if u != nil {
		pcpu.SetMemberID(*u)
	}
	return pcpu
}

// ClearMemberID clears the value of the "member_id" field.
func (pcpu *PromotionCommissionPlanUpdate) ClearMemberID() *PromotionCommissionPlanUpdate {
	pcpu.mutation.ClearMemberID()
	return pcpu
}

// SetCommissionID sets the "commission_id" field.
func (pcpu *PromotionCommissionPlanUpdate) SetCommissionID(u uint64) *PromotionCommissionPlanUpdate {
	pcpu.mutation.SetCommissionID(u)
	return pcpu
}

// SetNillableCommissionID sets the "commission_id" field if the given value is not nil.
func (pcpu *PromotionCommissionPlanUpdate) SetNillableCommissionID(u *uint64) *PromotionCommissionPlanUpdate {
	if u != nil {
		pcpu.SetCommissionID(*u)
	}
	return pcpu
}

// ClearCommissionID clears the value of the "commission_id" field.
func (pcpu *PromotionCommissionPlanUpdate) ClearCommissionID() *PromotionCommissionPlanUpdate {
	pcpu.mutation.ClearCommissionID()
	return pcpu
}

// SetPlanID sets the "plan_id" field.
func (pcpu *PromotionCommissionPlanUpdate) SetPlanID(u uint64) *PromotionCommissionPlanUpdate {
	pcpu.mutation.SetPlanID(u)
	return pcpu
}

// SetNillablePlanID sets the "plan_id" field if the given value is not nil.
func (pcpu *PromotionCommissionPlanUpdate) SetNillablePlanID(u *uint64) *PromotionCommissionPlanUpdate {
	if u != nil {
		pcpu.SetPlanID(*u)
	}
	return pcpu
}

// ClearPlanID clears the value of the "plan_id" field.
func (pcpu *PromotionCommissionPlanUpdate) ClearPlanID() *PromotionCommissionPlanUpdate {
	pcpu.mutation.ClearPlanID()
	return pcpu
}

// SetMember sets the "member" edge to the PromotionMember entity.
func (pcpu *PromotionCommissionPlanUpdate) SetMember(p *PromotionMember) *PromotionCommissionPlanUpdate {
	return pcpu.SetMemberID(p.ID)
}

// SetPromotionCommissionID sets the "promotion_commission" edge to the PromotionCommission entity by ID.
func (pcpu *PromotionCommissionPlanUpdate) SetPromotionCommissionID(id uint64) *PromotionCommissionPlanUpdate {
	pcpu.mutation.SetPromotionCommissionID(id)
	return pcpu
}

// SetNillablePromotionCommissionID sets the "promotion_commission" edge to the PromotionCommission entity by ID if the given value is not nil.
func (pcpu *PromotionCommissionPlanUpdate) SetNillablePromotionCommissionID(id *uint64) *PromotionCommissionPlanUpdate {
	if id != nil {
		pcpu = pcpu.SetPromotionCommissionID(*id)
	}
	return pcpu
}

// SetPromotionCommission sets the "promotion_commission" edge to the PromotionCommission entity.
func (pcpu *PromotionCommissionPlanUpdate) SetPromotionCommission(p *PromotionCommission) *PromotionCommissionPlanUpdate {
	return pcpu.SetPromotionCommissionID(p.ID)
}

// SetPlan sets the "plan" edge to the Plan entity.
func (pcpu *PromotionCommissionPlanUpdate) SetPlan(p *Plan) *PromotionCommissionPlanUpdate {
	return pcpu.SetPlanID(p.ID)
}

// Mutation returns the PromotionCommissionPlanMutation object of the builder.
func (pcpu *PromotionCommissionPlanUpdate) Mutation() *PromotionCommissionPlanMutation {
	return pcpu.mutation
}

// ClearMember clears the "member" edge to the PromotionMember entity.
func (pcpu *PromotionCommissionPlanUpdate) ClearMember() *PromotionCommissionPlanUpdate {
	pcpu.mutation.ClearMember()
	return pcpu
}

// ClearPromotionCommission clears the "promotion_commission" edge to the PromotionCommission entity.
func (pcpu *PromotionCommissionPlanUpdate) ClearPromotionCommission() *PromotionCommissionPlanUpdate {
	pcpu.mutation.ClearPromotionCommission()
	return pcpu
}

// ClearPlan clears the "plan" edge to the Plan entity.
func (pcpu *PromotionCommissionPlanUpdate) ClearPlan() *PromotionCommissionPlanUpdate {
	pcpu.mutation.ClearPlan()
	return pcpu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pcpu *PromotionCommissionPlanUpdate) Save(ctx context.Context) (int, error) {
	pcpu.defaults()
	return withHooks(ctx, pcpu.sqlSave, pcpu.mutation, pcpu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pcpu *PromotionCommissionPlanUpdate) SaveX(ctx context.Context) int {
	affected, err := pcpu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pcpu *PromotionCommissionPlanUpdate) Exec(ctx context.Context) error {
	_, err := pcpu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pcpu *PromotionCommissionPlanUpdate) ExecX(ctx context.Context) {
	if err := pcpu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pcpu *PromotionCommissionPlanUpdate) defaults() {
	if _, ok := pcpu.mutation.UpdatedAt(); !ok {
		v := promotioncommissionplan.UpdateDefaultUpdatedAt()
		pcpu.mutation.SetUpdatedAt(v)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (pcpu *PromotionCommissionPlanUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PromotionCommissionPlanUpdate {
	pcpu.modifiers = append(pcpu.modifiers, modifiers...)
	return pcpu
}

func (pcpu *PromotionCommissionPlanUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(promotioncommissionplan.Table, promotioncommissionplan.Columns, sqlgraph.NewFieldSpec(promotioncommissionplan.FieldID, field.TypeUint64))
	if ps := pcpu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pcpu.mutation.UpdatedAt(); ok {
		_spec.SetField(promotioncommissionplan.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := pcpu.mutation.DeletedAt(); ok {
		_spec.SetField(promotioncommissionplan.FieldDeletedAt, field.TypeTime, value)
	}
	if pcpu.mutation.DeletedAtCleared() {
		_spec.ClearField(promotioncommissionplan.FieldDeletedAt, field.TypeTime)
	}
	if pcpu.mutation.MemberCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   promotioncommissionplan.MemberTable,
			Columns: []string{promotioncommissionplan.MemberColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(promotionmember.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pcpu.mutation.MemberIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   promotioncommissionplan.MemberTable,
			Columns: []string{promotioncommissionplan.MemberColumn},
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
	if pcpu.mutation.PromotionCommissionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   promotioncommissionplan.PromotionCommissionTable,
			Columns: []string{promotioncommissionplan.PromotionCommissionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(promotioncommission.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pcpu.mutation.PromotionCommissionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   promotioncommissionplan.PromotionCommissionTable,
			Columns: []string{promotioncommissionplan.PromotionCommissionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(promotioncommission.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pcpu.mutation.PlanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   promotioncommissionplan.PlanTable,
			Columns: []string{promotioncommissionplan.PlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(plan.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pcpu.mutation.PlanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   promotioncommissionplan.PlanTable,
			Columns: []string{promotioncommissionplan.PlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(plan.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(pcpu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, pcpu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{promotioncommissionplan.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pcpu.mutation.done = true
	return n, nil
}

// PromotionCommissionPlanUpdateOne is the builder for updating a single PromotionCommissionPlan entity.
type PromotionCommissionPlanUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *PromotionCommissionPlanMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetUpdatedAt(t time.Time) *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.SetUpdatedAt(t)
	return pcpuo
}

// SetDeletedAt sets the "deleted_at" field.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetDeletedAt(t time.Time) *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.SetDeletedAt(t)
	return pcpuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetNillableDeletedAt(t *time.Time) *PromotionCommissionPlanUpdateOne {
	if t != nil {
		pcpuo.SetDeletedAt(*t)
	}
	return pcpuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (pcpuo *PromotionCommissionPlanUpdateOne) ClearDeletedAt() *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.ClearDeletedAt()
	return pcpuo
}

// SetMemberID sets the "member_id" field.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetMemberID(u uint64) *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.SetMemberID(u)
	return pcpuo
}

// SetNillableMemberID sets the "member_id" field if the given value is not nil.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetNillableMemberID(u *uint64) *PromotionCommissionPlanUpdateOne {
	if u != nil {
		pcpuo.SetMemberID(*u)
	}
	return pcpuo
}

// ClearMemberID clears the value of the "member_id" field.
func (pcpuo *PromotionCommissionPlanUpdateOne) ClearMemberID() *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.ClearMemberID()
	return pcpuo
}

// SetCommissionID sets the "commission_id" field.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetCommissionID(u uint64) *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.SetCommissionID(u)
	return pcpuo
}

// SetNillableCommissionID sets the "commission_id" field if the given value is not nil.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetNillableCommissionID(u *uint64) *PromotionCommissionPlanUpdateOne {
	if u != nil {
		pcpuo.SetCommissionID(*u)
	}
	return pcpuo
}

// ClearCommissionID clears the value of the "commission_id" field.
func (pcpuo *PromotionCommissionPlanUpdateOne) ClearCommissionID() *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.ClearCommissionID()
	return pcpuo
}

// SetPlanID sets the "plan_id" field.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetPlanID(u uint64) *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.SetPlanID(u)
	return pcpuo
}

// SetNillablePlanID sets the "plan_id" field if the given value is not nil.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetNillablePlanID(u *uint64) *PromotionCommissionPlanUpdateOne {
	if u != nil {
		pcpuo.SetPlanID(*u)
	}
	return pcpuo
}

// ClearPlanID clears the value of the "plan_id" field.
func (pcpuo *PromotionCommissionPlanUpdateOne) ClearPlanID() *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.ClearPlanID()
	return pcpuo
}

// SetMember sets the "member" edge to the PromotionMember entity.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetMember(p *PromotionMember) *PromotionCommissionPlanUpdateOne {
	return pcpuo.SetMemberID(p.ID)
}

// SetPromotionCommissionID sets the "promotion_commission" edge to the PromotionCommission entity by ID.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetPromotionCommissionID(id uint64) *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.SetPromotionCommissionID(id)
	return pcpuo
}

// SetNillablePromotionCommissionID sets the "promotion_commission" edge to the PromotionCommission entity by ID if the given value is not nil.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetNillablePromotionCommissionID(id *uint64) *PromotionCommissionPlanUpdateOne {
	if id != nil {
		pcpuo = pcpuo.SetPromotionCommissionID(*id)
	}
	return pcpuo
}

// SetPromotionCommission sets the "promotion_commission" edge to the PromotionCommission entity.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetPromotionCommission(p *PromotionCommission) *PromotionCommissionPlanUpdateOne {
	return pcpuo.SetPromotionCommissionID(p.ID)
}

// SetPlan sets the "plan" edge to the Plan entity.
func (pcpuo *PromotionCommissionPlanUpdateOne) SetPlan(p *Plan) *PromotionCommissionPlanUpdateOne {
	return pcpuo.SetPlanID(p.ID)
}

// Mutation returns the PromotionCommissionPlanMutation object of the builder.
func (pcpuo *PromotionCommissionPlanUpdateOne) Mutation() *PromotionCommissionPlanMutation {
	return pcpuo.mutation
}

// ClearMember clears the "member" edge to the PromotionMember entity.
func (pcpuo *PromotionCommissionPlanUpdateOne) ClearMember() *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.ClearMember()
	return pcpuo
}

// ClearPromotionCommission clears the "promotion_commission" edge to the PromotionCommission entity.
func (pcpuo *PromotionCommissionPlanUpdateOne) ClearPromotionCommission() *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.ClearPromotionCommission()
	return pcpuo
}

// ClearPlan clears the "plan" edge to the Plan entity.
func (pcpuo *PromotionCommissionPlanUpdateOne) ClearPlan() *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.ClearPlan()
	return pcpuo
}

// Where appends a list predicates to the PromotionCommissionPlanUpdate builder.
func (pcpuo *PromotionCommissionPlanUpdateOne) Where(ps ...predicate.PromotionCommissionPlan) *PromotionCommissionPlanUpdateOne {
	pcpuo.mutation.Where(ps...)
	return pcpuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (pcpuo *PromotionCommissionPlanUpdateOne) Select(field string, fields ...string) *PromotionCommissionPlanUpdateOne {
	pcpuo.fields = append([]string{field}, fields...)
	return pcpuo
}

// Save executes the query and returns the updated PromotionCommissionPlan entity.
func (pcpuo *PromotionCommissionPlanUpdateOne) Save(ctx context.Context) (*PromotionCommissionPlan, error) {
	pcpuo.defaults()
	return withHooks(ctx, pcpuo.sqlSave, pcpuo.mutation, pcpuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pcpuo *PromotionCommissionPlanUpdateOne) SaveX(ctx context.Context) *PromotionCommissionPlan {
	node, err := pcpuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (pcpuo *PromotionCommissionPlanUpdateOne) Exec(ctx context.Context) error {
	_, err := pcpuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pcpuo *PromotionCommissionPlanUpdateOne) ExecX(ctx context.Context) {
	if err := pcpuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pcpuo *PromotionCommissionPlanUpdateOne) defaults() {
	if _, ok := pcpuo.mutation.UpdatedAt(); !ok {
		v := promotioncommissionplan.UpdateDefaultUpdatedAt()
		pcpuo.mutation.SetUpdatedAt(v)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (pcpuo *PromotionCommissionPlanUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PromotionCommissionPlanUpdateOne {
	pcpuo.modifiers = append(pcpuo.modifiers, modifiers...)
	return pcpuo
}

func (pcpuo *PromotionCommissionPlanUpdateOne) sqlSave(ctx context.Context) (_node *PromotionCommissionPlan, err error) {
	_spec := sqlgraph.NewUpdateSpec(promotioncommissionplan.Table, promotioncommissionplan.Columns, sqlgraph.NewFieldSpec(promotioncommissionplan.FieldID, field.TypeUint64))
	id, ok := pcpuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "PromotionCommissionPlan.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := pcpuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, promotioncommissionplan.FieldID)
		for _, f := range fields {
			if !promotioncommissionplan.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != promotioncommissionplan.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := pcpuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pcpuo.mutation.UpdatedAt(); ok {
		_spec.SetField(promotioncommissionplan.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := pcpuo.mutation.DeletedAt(); ok {
		_spec.SetField(promotioncommissionplan.FieldDeletedAt, field.TypeTime, value)
	}
	if pcpuo.mutation.DeletedAtCleared() {
		_spec.ClearField(promotioncommissionplan.FieldDeletedAt, field.TypeTime)
	}
	if pcpuo.mutation.MemberCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   promotioncommissionplan.MemberTable,
			Columns: []string{promotioncommissionplan.MemberColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(promotionmember.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pcpuo.mutation.MemberIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   promotioncommissionplan.MemberTable,
			Columns: []string{promotioncommissionplan.MemberColumn},
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
	if pcpuo.mutation.PromotionCommissionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   promotioncommissionplan.PromotionCommissionTable,
			Columns: []string{promotioncommissionplan.PromotionCommissionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(promotioncommission.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pcpuo.mutation.PromotionCommissionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   promotioncommissionplan.PromotionCommissionTable,
			Columns: []string{promotioncommissionplan.PromotionCommissionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(promotioncommission.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pcpuo.mutation.PlanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   promotioncommissionplan.PlanTable,
			Columns: []string{promotioncommissionplan.PlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(plan.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pcpuo.mutation.PlanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   promotioncommissionplan.PlanTable,
			Columns: []string{promotioncommissionplan.PlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(plan.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(pcpuo.modifiers...)
	_node = &PromotionCommissionPlan{config: pcpuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, pcpuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{promotioncommissionplan.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	pcpuo.mutation.done = true
	return _node, nil
}
