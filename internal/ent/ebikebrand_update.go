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
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// EbikeBrandUpdate is the builder for updating EbikeBrand entities.
type EbikeBrandUpdate struct {
	config
	hooks     []Hook
	mutation  *EbikeBrandMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the EbikeBrandUpdate builder.
func (ebu *EbikeBrandUpdate) Where(ps ...predicate.EbikeBrand) *EbikeBrandUpdate {
	ebu.mutation.Where(ps...)
	return ebu
}

// SetUpdatedAt sets the "updated_at" field.
func (ebu *EbikeBrandUpdate) SetUpdatedAt(t time.Time) *EbikeBrandUpdate {
	ebu.mutation.SetUpdatedAt(t)
	return ebu
}

// SetDeletedAt sets the "deleted_at" field.
func (ebu *EbikeBrandUpdate) SetDeletedAt(t time.Time) *EbikeBrandUpdate {
	ebu.mutation.SetDeletedAt(t)
	return ebu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (ebu *EbikeBrandUpdate) SetNillableDeletedAt(t *time.Time) *EbikeBrandUpdate {
	if t != nil {
		ebu.SetDeletedAt(*t)
	}
	return ebu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (ebu *EbikeBrandUpdate) ClearDeletedAt() *EbikeBrandUpdate {
	ebu.mutation.ClearDeletedAt()
	return ebu
}

// SetLastModifier sets the "last_modifier" field.
func (ebu *EbikeBrandUpdate) SetLastModifier(m *model.Modifier) *EbikeBrandUpdate {
	ebu.mutation.SetLastModifier(m)
	return ebu
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (ebu *EbikeBrandUpdate) ClearLastModifier() *EbikeBrandUpdate {
	ebu.mutation.ClearLastModifier()
	return ebu
}

// SetRemark sets the "remark" field.
func (ebu *EbikeBrandUpdate) SetRemark(s string) *EbikeBrandUpdate {
	ebu.mutation.SetRemark(s)
	return ebu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (ebu *EbikeBrandUpdate) SetNillableRemark(s *string) *EbikeBrandUpdate {
	if s != nil {
		ebu.SetRemark(*s)
	}
	return ebu
}

// ClearRemark clears the value of the "remark" field.
func (ebu *EbikeBrandUpdate) ClearRemark() *EbikeBrandUpdate {
	ebu.mutation.ClearRemark()
	return ebu
}

// SetName sets the "name" field.
func (ebu *EbikeBrandUpdate) SetName(s string) *EbikeBrandUpdate {
	ebu.mutation.SetName(s)
	return ebu
}

// SetCover sets the "cover" field.
func (ebu *EbikeBrandUpdate) SetCover(s string) *EbikeBrandUpdate {
	ebu.mutation.SetCover(s)
	return ebu
}

// Mutation returns the EbikeBrandMutation object of the builder.
func (ebu *EbikeBrandUpdate) Mutation() *EbikeBrandMutation {
	return ebu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ebu *EbikeBrandUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if err := ebu.defaults(); err != nil {
		return 0, err
	}
	if len(ebu.hooks) == 0 {
		affected, err = ebu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EbikeBrandMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ebu.mutation = mutation
			affected, err = ebu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ebu.hooks) - 1; i >= 0; i-- {
			if ebu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ebu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ebu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (ebu *EbikeBrandUpdate) SaveX(ctx context.Context) int {
	affected, err := ebu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ebu *EbikeBrandUpdate) Exec(ctx context.Context) error {
	_, err := ebu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ebu *EbikeBrandUpdate) ExecX(ctx context.Context) {
	if err := ebu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ebu *EbikeBrandUpdate) defaults() error {
	if _, ok := ebu.mutation.UpdatedAt(); !ok {
		if ebikebrand.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized ebikebrand.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := ebikebrand.UpdateDefaultUpdatedAt()
		ebu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ebu *EbikeBrandUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *EbikeBrandUpdate {
	ebu.modifiers = append(ebu.modifiers, modifiers...)
	return ebu
}

func (ebu *EbikeBrandUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   ebikebrand.Table,
			Columns: ebikebrand.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: ebikebrand.FieldID,
			},
		},
	}
	if ps := ebu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ebu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: ebikebrand.FieldUpdatedAt,
		})
	}
	if value, ok := ebu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: ebikebrand.FieldDeletedAt,
		})
	}
	if ebu.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: ebikebrand.FieldDeletedAt,
		})
	}
	if ebu.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: ebikebrand.FieldCreator,
		})
	}
	if value, ok := ebu.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: ebikebrand.FieldLastModifier,
		})
	}
	if ebu.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: ebikebrand.FieldLastModifier,
		})
	}
	if value, ok := ebu.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: ebikebrand.FieldRemark,
		})
	}
	if ebu.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: ebikebrand.FieldRemark,
		})
	}
	if value, ok := ebu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: ebikebrand.FieldName,
		})
	}
	if value, ok := ebu.mutation.Cover(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: ebikebrand.FieldCover,
		})
	}
	_spec.Modifiers = ebu.modifiers
	if n, err = sqlgraph.UpdateNodes(ctx, ebu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{ebikebrand.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	return n, nil
}

// EbikeBrandUpdateOne is the builder for updating a single EbikeBrand entity.
type EbikeBrandUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *EbikeBrandMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (ebuo *EbikeBrandUpdateOne) SetUpdatedAt(t time.Time) *EbikeBrandUpdateOne {
	ebuo.mutation.SetUpdatedAt(t)
	return ebuo
}

// SetDeletedAt sets the "deleted_at" field.
func (ebuo *EbikeBrandUpdateOne) SetDeletedAt(t time.Time) *EbikeBrandUpdateOne {
	ebuo.mutation.SetDeletedAt(t)
	return ebuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (ebuo *EbikeBrandUpdateOne) SetNillableDeletedAt(t *time.Time) *EbikeBrandUpdateOne {
	if t != nil {
		ebuo.SetDeletedAt(*t)
	}
	return ebuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (ebuo *EbikeBrandUpdateOne) ClearDeletedAt() *EbikeBrandUpdateOne {
	ebuo.mutation.ClearDeletedAt()
	return ebuo
}

// SetLastModifier sets the "last_modifier" field.
func (ebuo *EbikeBrandUpdateOne) SetLastModifier(m *model.Modifier) *EbikeBrandUpdateOne {
	ebuo.mutation.SetLastModifier(m)
	return ebuo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (ebuo *EbikeBrandUpdateOne) ClearLastModifier() *EbikeBrandUpdateOne {
	ebuo.mutation.ClearLastModifier()
	return ebuo
}

// SetRemark sets the "remark" field.
func (ebuo *EbikeBrandUpdateOne) SetRemark(s string) *EbikeBrandUpdateOne {
	ebuo.mutation.SetRemark(s)
	return ebuo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (ebuo *EbikeBrandUpdateOne) SetNillableRemark(s *string) *EbikeBrandUpdateOne {
	if s != nil {
		ebuo.SetRemark(*s)
	}
	return ebuo
}

// ClearRemark clears the value of the "remark" field.
func (ebuo *EbikeBrandUpdateOne) ClearRemark() *EbikeBrandUpdateOne {
	ebuo.mutation.ClearRemark()
	return ebuo
}

// SetName sets the "name" field.
func (ebuo *EbikeBrandUpdateOne) SetName(s string) *EbikeBrandUpdateOne {
	ebuo.mutation.SetName(s)
	return ebuo
}

// SetCover sets the "cover" field.
func (ebuo *EbikeBrandUpdateOne) SetCover(s string) *EbikeBrandUpdateOne {
	ebuo.mutation.SetCover(s)
	return ebuo
}

// Mutation returns the EbikeBrandMutation object of the builder.
func (ebuo *EbikeBrandUpdateOne) Mutation() *EbikeBrandMutation {
	return ebuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ebuo *EbikeBrandUpdateOne) Select(field string, fields ...string) *EbikeBrandUpdateOne {
	ebuo.fields = append([]string{field}, fields...)
	return ebuo
}

// Save executes the query and returns the updated EbikeBrand entity.
func (ebuo *EbikeBrandUpdateOne) Save(ctx context.Context) (*EbikeBrand, error) {
	var (
		err  error
		node *EbikeBrand
	)
	if err := ebuo.defaults(); err != nil {
		return nil, err
	}
	if len(ebuo.hooks) == 0 {
		node, err = ebuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EbikeBrandMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ebuo.mutation = mutation
			node, err = ebuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ebuo.hooks) - 1; i >= 0; i-- {
			if ebuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ebuo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, ebuo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*EbikeBrand)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from EbikeBrandMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (ebuo *EbikeBrandUpdateOne) SaveX(ctx context.Context) *EbikeBrand {
	node, err := ebuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ebuo *EbikeBrandUpdateOne) Exec(ctx context.Context) error {
	_, err := ebuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ebuo *EbikeBrandUpdateOne) ExecX(ctx context.Context) {
	if err := ebuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ebuo *EbikeBrandUpdateOne) defaults() error {
	if _, ok := ebuo.mutation.UpdatedAt(); !ok {
		if ebikebrand.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized ebikebrand.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := ebikebrand.UpdateDefaultUpdatedAt()
		ebuo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ebuo *EbikeBrandUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *EbikeBrandUpdateOne {
	ebuo.modifiers = append(ebuo.modifiers, modifiers...)
	return ebuo
}

func (ebuo *EbikeBrandUpdateOne) sqlSave(ctx context.Context) (_node *EbikeBrand, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   ebikebrand.Table,
			Columns: ebikebrand.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: ebikebrand.FieldID,
			},
		},
	}
	id, ok := ebuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "EbikeBrand.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ebuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, ebikebrand.FieldID)
		for _, f := range fields {
			if !ebikebrand.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != ebikebrand.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ebuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ebuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: ebikebrand.FieldUpdatedAt,
		})
	}
	if value, ok := ebuo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: ebikebrand.FieldDeletedAt,
		})
	}
	if ebuo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: ebikebrand.FieldDeletedAt,
		})
	}
	if ebuo.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: ebikebrand.FieldCreator,
		})
	}
	if value, ok := ebuo.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: ebikebrand.FieldLastModifier,
		})
	}
	if ebuo.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: ebikebrand.FieldLastModifier,
		})
	}
	if value, ok := ebuo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: ebikebrand.FieldRemark,
		})
	}
	if ebuo.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: ebikebrand.FieldRemark,
		})
	}
	if value, ok := ebuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: ebikebrand.FieldName,
		})
	}
	if value, ok := ebuo.mutation.Cover(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: ebikebrand.FieldCover,
		})
	}
	_spec.Modifiers = ebuo.modifiers
	_node = &EbikeBrand{config: ebuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ebuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{ebikebrand.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}