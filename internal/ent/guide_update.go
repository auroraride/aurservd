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
	"github.com/auroraride/aurservd/internal/ent/guide"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// GuideUpdate is the builder for updating Guide entities.
type GuideUpdate struct {
	config
	hooks     []Hook
	mutation  *GuideMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the GuideUpdate builder.
func (gu *GuideUpdate) Where(ps ...predicate.Guide) *GuideUpdate {
	gu.mutation.Where(ps...)
	return gu
}

// SetUpdatedAt sets the "updated_at" field.
func (gu *GuideUpdate) SetUpdatedAt(t time.Time) *GuideUpdate {
	gu.mutation.SetUpdatedAt(t)
	return gu
}

// SetDeletedAt sets the "deleted_at" field.
func (gu *GuideUpdate) SetDeletedAt(t time.Time) *GuideUpdate {
	gu.mutation.SetDeletedAt(t)
	return gu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (gu *GuideUpdate) SetNillableDeletedAt(t *time.Time) *GuideUpdate {
	if t != nil {
		gu.SetDeletedAt(*t)
	}
	return gu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (gu *GuideUpdate) ClearDeletedAt() *GuideUpdate {
	gu.mutation.ClearDeletedAt()
	return gu
}

// SetLastModifier sets the "last_modifier" field.
func (gu *GuideUpdate) SetLastModifier(m *model.Modifier) *GuideUpdate {
	gu.mutation.SetLastModifier(m)
	return gu
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (gu *GuideUpdate) ClearLastModifier() *GuideUpdate {
	gu.mutation.ClearLastModifier()
	return gu
}

// SetRemark sets the "remark" field.
func (gu *GuideUpdate) SetRemark(s string) *GuideUpdate {
	gu.mutation.SetRemark(s)
	return gu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (gu *GuideUpdate) SetNillableRemark(s *string) *GuideUpdate {
	if s != nil {
		gu.SetRemark(*s)
	}
	return gu
}

// ClearRemark clears the value of the "remark" field.
func (gu *GuideUpdate) ClearRemark() *GuideUpdate {
	gu.mutation.ClearRemark()
	return gu
}

// SetName sets the "name" field.
func (gu *GuideUpdate) SetName(s string) *GuideUpdate {
	gu.mutation.SetName(s)
	return gu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (gu *GuideUpdate) SetNillableName(s *string) *GuideUpdate {
	if s != nil {
		gu.SetName(*s)
	}
	return gu
}

// SetSort sets the "sort" field.
func (gu *GuideUpdate) SetSort(u uint8) *GuideUpdate {
	gu.mutation.ResetSort()
	gu.mutation.SetSort(u)
	return gu
}

// SetNillableSort sets the "sort" field if the given value is not nil.
func (gu *GuideUpdate) SetNillableSort(u *uint8) *GuideUpdate {
	if u != nil {
		gu.SetSort(*u)
	}
	return gu
}

// AddSort adds u to the "sort" field.
func (gu *GuideUpdate) AddSort(u int8) *GuideUpdate {
	gu.mutation.AddSort(u)
	return gu
}

// SetAnswer sets the "answer" field.
func (gu *GuideUpdate) SetAnswer(s string) *GuideUpdate {
	gu.mutation.SetAnswer(s)
	return gu
}

// SetNillableAnswer sets the "answer" field if the given value is not nil.
func (gu *GuideUpdate) SetNillableAnswer(s *string) *GuideUpdate {
	if s != nil {
		gu.SetAnswer(*s)
	}
	return gu
}

// Mutation returns the GuideMutation object of the builder.
func (gu *GuideUpdate) Mutation() *GuideMutation {
	return gu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (gu *GuideUpdate) Save(ctx context.Context) (int, error) {
	if err := gu.defaults(); err != nil {
		return 0, err
	}
	return withHooks(ctx, gu.sqlSave, gu.mutation, gu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (gu *GuideUpdate) SaveX(ctx context.Context) int {
	affected, err := gu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (gu *GuideUpdate) Exec(ctx context.Context) error {
	_, err := gu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gu *GuideUpdate) ExecX(ctx context.Context) {
	if err := gu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gu *GuideUpdate) defaults() error {
	if _, ok := gu.mutation.UpdatedAt(); !ok {
		if guide.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized guide.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := guide.UpdateDefaultUpdatedAt()
		gu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (gu *GuideUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *GuideUpdate {
	gu.modifiers = append(gu.modifiers, modifiers...)
	return gu
}

func (gu *GuideUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(guide.Table, guide.Columns, sqlgraph.NewFieldSpec(guide.FieldID, field.TypeUint64))
	if ps := gu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := gu.mutation.UpdatedAt(); ok {
		_spec.SetField(guide.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := gu.mutation.DeletedAt(); ok {
		_spec.SetField(guide.FieldDeletedAt, field.TypeTime, value)
	}
	if gu.mutation.DeletedAtCleared() {
		_spec.ClearField(guide.FieldDeletedAt, field.TypeTime)
	}
	if gu.mutation.CreatorCleared() {
		_spec.ClearField(guide.FieldCreator, field.TypeJSON)
	}
	if value, ok := gu.mutation.LastModifier(); ok {
		_spec.SetField(guide.FieldLastModifier, field.TypeJSON, value)
	}
	if gu.mutation.LastModifierCleared() {
		_spec.ClearField(guide.FieldLastModifier, field.TypeJSON)
	}
	if value, ok := gu.mutation.Remark(); ok {
		_spec.SetField(guide.FieldRemark, field.TypeString, value)
	}
	if gu.mutation.RemarkCleared() {
		_spec.ClearField(guide.FieldRemark, field.TypeString)
	}
	if value, ok := gu.mutation.Name(); ok {
		_spec.SetField(guide.FieldName, field.TypeString, value)
	}
	if value, ok := gu.mutation.Sort(); ok {
		_spec.SetField(guide.FieldSort, field.TypeUint8, value)
	}
	if value, ok := gu.mutation.AddedSort(); ok {
		_spec.AddField(guide.FieldSort, field.TypeUint8, value)
	}
	if value, ok := gu.mutation.Answer(); ok {
		_spec.SetField(guide.FieldAnswer, field.TypeString, value)
	}
	_spec.AddModifiers(gu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, gu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{guide.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	gu.mutation.done = true
	return n, nil
}

// GuideUpdateOne is the builder for updating a single Guide entity.
type GuideUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *GuideMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (guo *GuideUpdateOne) SetUpdatedAt(t time.Time) *GuideUpdateOne {
	guo.mutation.SetUpdatedAt(t)
	return guo
}

// SetDeletedAt sets the "deleted_at" field.
func (guo *GuideUpdateOne) SetDeletedAt(t time.Time) *GuideUpdateOne {
	guo.mutation.SetDeletedAt(t)
	return guo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (guo *GuideUpdateOne) SetNillableDeletedAt(t *time.Time) *GuideUpdateOne {
	if t != nil {
		guo.SetDeletedAt(*t)
	}
	return guo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (guo *GuideUpdateOne) ClearDeletedAt() *GuideUpdateOne {
	guo.mutation.ClearDeletedAt()
	return guo
}

// SetLastModifier sets the "last_modifier" field.
func (guo *GuideUpdateOne) SetLastModifier(m *model.Modifier) *GuideUpdateOne {
	guo.mutation.SetLastModifier(m)
	return guo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (guo *GuideUpdateOne) ClearLastModifier() *GuideUpdateOne {
	guo.mutation.ClearLastModifier()
	return guo
}

// SetRemark sets the "remark" field.
func (guo *GuideUpdateOne) SetRemark(s string) *GuideUpdateOne {
	guo.mutation.SetRemark(s)
	return guo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (guo *GuideUpdateOne) SetNillableRemark(s *string) *GuideUpdateOne {
	if s != nil {
		guo.SetRemark(*s)
	}
	return guo
}

// ClearRemark clears the value of the "remark" field.
func (guo *GuideUpdateOne) ClearRemark() *GuideUpdateOne {
	guo.mutation.ClearRemark()
	return guo
}

// SetName sets the "name" field.
func (guo *GuideUpdateOne) SetName(s string) *GuideUpdateOne {
	guo.mutation.SetName(s)
	return guo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (guo *GuideUpdateOne) SetNillableName(s *string) *GuideUpdateOne {
	if s != nil {
		guo.SetName(*s)
	}
	return guo
}

// SetSort sets the "sort" field.
func (guo *GuideUpdateOne) SetSort(u uint8) *GuideUpdateOne {
	guo.mutation.ResetSort()
	guo.mutation.SetSort(u)
	return guo
}

// SetNillableSort sets the "sort" field if the given value is not nil.
func (guo *GuideUpdateOne) SetNillableSort(u *uint8) *GuideUpdateOne {
	if u != nil {
		guo.SetSort(*u)
	}
	return guo
}

// AddSort adds u to the "sort" field.
func (guo *GuideUpdateOne) AddSort(u int8) *GuideUpdateOne {
	guo.mutation.AddSort(u)
	return guo
}

// SetAnswer sets the "answer" field.
func (guo *GuideUpdateOne) SetAnswer(s string) *GuideUpdateOne {
	guo.mutation.SetAnswer(s)
	return guo
}

// SetNillableAnswer sets the "answer" field if the given value is not nil.
func (guo *GuideUpdateOne) SetNillableAnswer(s *string) *GuideUpdateOne {
	if s != nil {
		guo.SetAnswer(*s)
	}
	return guo
}

// Mutation returns the GuideMutation object of the builder.
func (guo *GuideUpdateOne) Mutation() *GuideMutation {
	return guo.mutation
}

// Where appends a list predicates to the GuideUpdate builder.
func (guo *GuideUpdateOne) Where(ps ...predicate.Guide) *GuideUpdateOne {
	guo.mutation.Where(ps...)
	return guo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (guo *GuideUpdateOne) Select(field string, fields ...string) *GuideUpdateOne {
	guo.fields = append([]string{field}, fields...)
	return guo
}

// Save executes the query and returns the updated Guide entity.
func (guo *GuideUpdateOne) Save(ctx context.Context) (*Guide, error) {
	if err := guo.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, guo.sqlSave, guo.mutation, guo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (guo *GuideUpdateOne) SaveX(ctx context.Context) *Guide {
	node, err := guo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (guo *GuideUpdateOne) Exec(ctx context.Context) error {
	_, err := guo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (guo *GuideUpdateOne) ExecX(ctx context.Context) {
	if err := guo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (guo *GuideUpdateOne) defaults() error {
	if _, ok := guo.mutation.UpdatedAt(); !ok {
		if guide.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized guide.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := guide.UpdateDefaultUpdatedAt()
		guo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (guo *GuideUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *GuideUpdateOne {
	guo.modifiers = append(guo.modifiers, modifiers...)
	return guo
}

func (guo *GuideUpdateOne) sqlSave(ctx context.Context) (_node *Guide, err error) {
	_spec := sqlgraph.NewUpdateSpec(guide.Table, guide.Columns, sqlgraph.NewFieldSpec(guide.FieldID, field.TypeUint64))
	id, ok := guo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Guide.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := guo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, guide.FieldID)
		for _, f := range fields {
			if !guide.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != guide.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := guo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := guo.mutation.UpdatedAt(); ok {
		_spec.SetField(guide.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := guo.mutation.DeletedAt(); ok {
		_spec.SetField(guide.FieldDeletedAt, field.TypeTime, value)
	}
	if guo.mutation.DeletedAtCleared() {
		_spec.ClearField(guide.FieldDeletedAt, field.TypeTime)
	}
	if guo.mutation.CreatorCleared() {
		_spec.ClearField(guide.FieldCreator, field.TypeJSON)
	}
	if value, ok := guo.mutation.LastModifier(); ok {
		_spec.SetField(guide.FieldLastModifier, field.TypeJSON, value)
	}
	if guo.mutation.LastModifierCleared() {
		_spec.ClearField(guide.FieldLastModifier, field.TypeJSON)
	}
	if value, ok := guo.mutation.Remark(); ok {
		_spec.SetField(guide.FieldRemark, field.TypeString, value)
	}
	if guo.mutation.RemarkCleared() {
		_spec.ClearField(guide.FieldRemark, field.TypeString)
	}
	if value, ok := guo.mutation.Name(); ok {
		_spec.SetField(guide.FieldName, field.TypeString, value)
	}
	if value, ok := guo.mutation.Sort(); ok {
		_spec.SetField(guide.FieldSort, field.TypeUint8, value)
	}
	if value, ok := guo.mutation.AddedSort(); ok {
		_spec.AddField(guide.FieldSort, field.TypeUint8, value)
	}
	if value, ok := guo.mutation.Answer(); ok {
		_spec.SetField(guide.FieldAnswer, field.TypeString, value)
	}
	_spec.AddModifiers(guo.modifiers...)
	_node = &Guide{config: guo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, guo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{guide.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	guo.mutation.done = true
	return _node, nil
}