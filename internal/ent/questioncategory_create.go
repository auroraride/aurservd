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
	"github.com/auroraride/aurservd/internal/ent/question"
	"github.com/auroraride/aurservd/internal/ent/questioncategory"
)

// QuestionCategoryCreate is the builder for creating a QuestionCategory entity.
type QuestionCategoryCreate struct {
	config
	mutation *QuestionCategoryMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (qcc *QuestionCategoryCreate) SetCreatedAt(t time.Time) *QuestionCategoryCreate {
	qcc.mutation.SetCreatedAt(t)
	return qcc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (qcc *QuestionCategoryCreate) SetNillableCreatedAt(t *time.Time) *QuestionCategoryCreate {
	if t != nil {
		qcc.SetCreatedAt(*t)
	}
	return qcc
}

// SetUpdatedAt sets the "updated_at" field.
func (qcc *QuestionCategoryCreate) SetUpdatedAt(t time.Time) *QuestionCategoryCreate {
	qcc.mutation.SetUpdatedAt(t)
	return qcc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (qcc *QuestionCategoryCreate) SetNillableUpdatedAt(t *time.Time) *QuestionCategoryCreate {
	if t != nil {
		qcc.SetUpdatedAt(*t)
	}
	return qcc
}

// SetDeletedAt sets the "deleted_at" field.
func (qcc *QuestionCategoryCreate) SetDeletedAt(t time.Time) *QuestionCategoryCreate {
	qcc.mutation.SetDeletedAt(t)
	return qcc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (qcc *QuestionCategoryCreate) SetNillableDeletedAt(t *time.Time) *QuestionCategoryCreate {
	if t != nil {
		qcc.SetDeletedAt(*t)
	}
	return qcc
}

// SetCreator sets the "creator" field.
func (qcc *QuestionCategoryCreate) SetCreator(m *model.Modifier) *QuestionCategoryCreate {
	qcc.mutation.SetCreator(m)
	return qcc
}

// SetLastModifier sets the "last_modifier" field.
func (qcc *QuestionCategoryCreate) SetLastModifier(m *model.Modifier) *QuestionCategoryCreate {
	qcc.mutation.SetLastModifier(m)
	return qcc
}

// SetRemark sets the "remark" field.
func (qcc *QuestionCategoryCreate) SetRemark(s string) *QuestionCategoryCreate {
	qcc.mutation.SetRemark(s)
	return qcc
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (qcc *QuestionCategoryCreate) SetNillableRemark(s *string) *QuestionCategoryCreate {
	if s != nil {
		qcc.SetRemark(*s)
	}
	return qcc
}

// SetName sets the "name" field.
func (qcc *QuestionCategoryCreate) SetName(s string) *QuestionCategoryCreate {
	qcc.mutation.SetName(s)
	return qcc
}

// SetSort sets the "sort" field.
func (qcc *QuestionCategoryCreate) SetSort(u uint64) *QuestionCategoryCreate {
	qcc.mutation.SetSort(u)
	return qcc
}

// AddQuestionIDs adds the "questions" edge to the Question entity by IDs.
func (qcc *QuestionCategoryCreate) AddQuestionIDs(ids ...uint64) *QuestionCategoryCreate {
	qcc.mutation.AddQuestionIDs(ids...)
	return qcc
}

// AddQuestions adds the "questions" edges to the Question entity.
func (qcc *QuestionCategoryCreate) AddQuestions(q ...*Question) *QuestionCategoryCreate {
	ids := make([]uint64, len(q))
	for i := range q {
		ids[i] = q[i].ID
	}
	return qcc.AddQuestionIDs(ids...)
}

// Mutation returns the QuestionCategoryMutation object of the builder.
func (qcc *QuestionCategoryCreate) Mutation() *QuestionCategoryMutation {
	return qcc.mutation
}

// Save creates the QuestionCategory in the database.
func (qcc *QuestionCategoryCreate) Save(ctx context.Context) (*QuestionCategory, error) {
	if err := qcc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, qcc.sqlSave, qcc.mutation, qcc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (qcc *QuestionCategoryCreate) SaveX(ctx context.Context) *QuestionCategory {
	v, err := qcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (qcc *QuestionCategoryCreate) Exec(ctx context.Context) error {
	_, err := qcc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (qcc *QuestionCategoryCreate) ExecX(ctx context.Context) {
	if err := qcc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (qcc *QuestionCategoryCreate) defaults() error {
	if _, ok := qcc.mutation.CreatedAt(); !ok {
		if questioncategory.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized questioncategory.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := questioncategory.DefaultCreatedAt()
		qcc.mutation.SetCreatedAt(v)
	}
	if _, ok := qcc.mutation.UpdatedAt(); !ok {
		if questioncategory.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized questioncategory.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := questioncategory.DefaultUpdatedAt()
		qcc.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (qcc *QuestionCategoryCreate) check() error {
	if _, ok := qcc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "QuestionCategory.created_at"`)}
	}
	if _, ok := qcc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "QuestionCategory.updated_at"`)}
	}
	if _, ok := qcc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "QuestionCategory.name"`)}
	}
	if _, ok := qcc.mutation.Sort(); !ok {
		return &ValidationError{Name: "sort", err: errors.New(`ent: missing required field "QuestionCategory.sort"`)}
	}
	return nil
}

func (qcc *QuestionCategoryCreate) sqlSave(ctx context.Context) (*QuestionCategory, error) {
	if err := qcc.check(); err != nil {
		return nil, err
	}
	_node, _spec := qcc.createSpec()
	if err := sqlgraph.CreateNode(ctx, qcc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	qcc.mutation.id = &_node.ID
	qcc.mutation.done = true
	return _node, nil
}

func (qcc *QuestionCategoryCreate) createSpec() (*QuestionCategory, *sqlgraph.CreateSpec) {
	var (
		_node = &QuestionCategory{config: qcc.config}
		_spec = sqlgraph.NewCreateSpec(questioncategory.Table, sqlgraph.NewFieldSpec(questioncategory.FieldID, field.TypeUint64))
	)
	_spec.OnConflict = qcc.conflict
	if value, ok := qcc.mutation.CreatedAt(); ok {
		_spec.SetField(questioncategory.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := qcc.mutation.UpdatedAt(); ok {
		_spec.SetField(questioncategory.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := qcc.mutation.DeletedAt(); ok {
		_spec.SetField(questioncategory.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	if value, ok := qcc.mutation.Creator(); ok {
		_spec.SetField(questioncategory.FieldCreator, field.TypeJSON, value)
		_node.Creator = value
	}
	if value, ok := qcc.mutation.LastModifier(); ok {
		_spec.SetField(questioncategory.FieldLastModifier, field.TypeJSON, value)
		_node.LastModifier = value
	}
	if value, ok := qcc.mutation.Remark(); ok {
		_spec.SetField(questioncategory.FieldRemark, field.TypeString, value)
		_node.Remark = value
	}
	if value, ok := qcc.mutation.Name(); ok {
		_spec.SetField(questioncategory.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := qcc.mutation.Sort(); ok {
		_spec.SetField(questioncategory.FieldSort, field.TypeUint64, value)
		_node.Sort = value
	}
	if nodes := qcc.mutation.QuestionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   questioncategory.QuestionsTable,
			Columns: []string{questioncategory.QuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(question.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.QuestionCategory.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.QuestionCategoryUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (qcc *QuestionCategoryCreate) OnConflict(opts ...sql.ConflictOption) *QuestionCategoryUpsertOne {
	qcc.conflict = opts
	return &QuestionCategoryUpsertOne{
		create: qcc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.QuestionCategory.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (qcc *QuestionCategoryCreate) OnConflictColumns(columns ...string) *QuestionCategoryUpsertOne {
	qcc.conflict = append(qcc.conflict, sql.ConflictColumns(columns...))
	return &QuestionCategoryUpsertOne{
		create: qcc,
	}
}

type (
	// QuestionCategoryUpsertOne is the builder for "upsert"-ing
	//  one QuestionCategory node.
	QuestionCategoryUpsertOne struct {
		create *QuestionCategoryCreate
	}

	// QuestionCategoryUpsert is the "OnConflict" setter.
	QuestionCategoryUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *QuestionCategoryUpsert) SetUpdatedAt(v time.Time) *QuestionCategoryUpsert {
	u.Set(questioncategory.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *QuestionCategoryUpsert) UpdateUpdatedAt() *QuestionCategoryUpsert {
	u.SetExcluded(questioncategory.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *QuestionCategoryUpsert) SetDeletedAt(v time.Time) *QuestionCategoryUpsert {
	u.Set(questioncategory.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *QuestionCategoryUpsert) UpdateDeletedAt() *QuestionCategoryUpsert {
	u.SetExcluded(questioncategory.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *QuestionCategoryUpsert) ClearDeletedAt() *QuestionCategoryUpsert {
	u.SetNull(questioncategory.FieldDeletedAt)
	return u
}

// SetLastModifier sets the "last_modifier" field.
func (u *QuestionCategoryUpsert) SetLastModifier(v *model.Modifier) *QuestionCategoryUpsert {
	u.Set(questioncategory.FieldLastModifier, v)
	return u
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *QuestionCategoryUpsert) UpdateLastModifier() *QuestionCategoryUpsert {
	u.SetExcluded(questioncategory.FieldLastModifier)
	return u
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *QuestionCategoryUpsert) ClearLastModifier() *QuestionCategoryUpsert {
	u.SetNull(questioncategory.FieldLastModifier)
	return u
}

// SetRemark sets the "remark" field.
func (u *QuestionCategoryUpsert) SetRemark(v string) *QuestionCategoryUpsert {
	u.Set(questioncategory.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *QuestionCategoryUpsert) UpdateRemark() *QuestionCategoryUpsert {
	u.SetExcluded(questioncategory.FieldRemark)
	return u
}

// ClearRemark clears the value of the "remark" field.
func (u *QuestionCategoryUpsert) ClearRemark() *QuestionCategoryUpsert {
	u.SetNull(questioncategory.FieldRemark)
	return u
}

// SetName sets the "name" field.
func (u *QuestionCategoryUpsert) SetName(v string) *QuestionCategoryUpsert {
	u.Set(questioncategory.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *QuestionCategoryUpsert) UpdateName() *QuestionCategoryUpsert {
	u.SetExcluded(questioncategory.FieldName)
	return u
}

// SetSort sets the "sort" field.
func (u *QuestionCategoryUpsert) SetSort(v uint64) *QuestionCategoryUpsert {
	u.Set(questioncategory.FieldSort, v)
	return u
}

// UpdateSort sets the "sort" field to the value that was provided on create.
func (u *QuestionCategoryUpsert) UpdateSort() *QuestionCategoryUpsert {
	u.SetExcluded(questioncategory.FieldSort)
	return u
}

// AddSort adds v to the "sort" field.
func (u *QuestionCategoryUpsert) AddSort(v uint64) *QuestionCategoryUpsert {
	u.Add(questioncategory.FieldSort, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.QuestionCategory.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *QuestionCategoryUpsertOne) UpdateNewValues() *QuestionCategoryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(questioncategory.FieldCreatedAt)
		}
		if _, exists := u.create.mutation.Creator(); exists {
			s.SetIgnore(questioncategory.FieldCreator)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.QuestionCategory.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *QuestionCategoryUpsertOne) Ignore() *QuestionCategoryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *QuestionCategoryUpsertOne) DoNothing() *QuestionCategoryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the QuestionCategoryCreate.OnConflict
// documentation for more info.
func (u *QuestionCategoryUpsertOne) Update(set func(*QuestionCategoryUpsert)) *QuestionCategoryUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&QuestionCategoryUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *QuestionCategoryUpsertOne) SetUpdatedAt(v time.Time) *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *QuestionCategoryUpsertOne) UpdateUpdatedAt() *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *QuestionCategoryUpsertOne) SetDeletedAt(v time.Time) *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *QuestionCategoryUpsertOne) UpdateDeletedAt() *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *QuestionCategoryUpsertOne) ClearDeletedAt() *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *QuestionCategoryUpsertOne) SetLastModifier(v *model.Modifier) *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *QuestionCategoryUpsertOne) UpdateLastModifier() *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *QuestionCategoryUpsertOne) ClearLastModifier() *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *QuestionCategoryUpsertOne) SetRemark(v string) *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *QuestionCategoryUpsertOne) UpdateRemark() *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *QuestionCategoryUpsertOne) ClearRemark() *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.ClearRemark()
	})
}

// SetName sets the "name" field.
func (u *QuestionCategoryUpsertOne) SetName(v string) *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *QuestionCategoryUpsertOne) UpdateName() *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateName()
	})
}

// SetSort sets the "sort" field.
func (u *QuestionCategoryUpsertOne) SetSort(v uint64) *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetSort(v)
	})
}

// AddSort adds v to the "sort" field.
func (u *QuestionCategoryUpsertOne) AddSort(v uint64) *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.AddSort(v)
	})
}

// UpdateSort sets the "sort" field to the value that was provided on create.
func (u *QuestionCategoryUpsertOne) UpdateSort() *QuestionCategoryUpsertOne {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateSort()
	})
}

// Exec executes the query.
func (u *QuestionCategoryUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for QuestionCategoryCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *QuestionCategoryUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *QuestionCategoryUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *QuestionCategoryUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// QuestionCategoryCreateBulk is the builder for creating many QuestionCategory entities in bulk.
type QuestionCategoryCreateBulk struct {
	config
	err      error
	builders []*QuestionCategoryCreate
	conflict []sql.ConflictOption
}

// Save creates the QuestionCategory entities in the database.
func (qccb *QuestionCategoryCreateBulk) Save(ctx context.Context) ([]*QuestionCategory, error) {
	if qccb.err != nil {
		return nil, qccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(qccb.builders))
	nodes := make([]*QuestionCategory, len(qccb.builders))
	mutators := make([]Mutator, len(qccb.builders))
	for i := range qccb.builders {
		func(i int, root context.Context) {
			builder := qccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*QuestionCategoryMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, qccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = qccb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, qccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = uint64(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, qccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (qccb *QuestionCategoryCreateBulk) SaveX(ctx context.Context) []*QuestionCategory {
	v, err := qccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (qccb *QuestionCategoryCreateBulk) Exec(ctx context.Context) error {
	_, err := qccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (qccb *QuestionCategoryCreateBulk) ExecX(ctx context.Context) {
	if err := qccb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.QuestionCategory.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.QuestionCategoryUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (qccb *QuestionCategoryCreateBulk) OnConflict(opts ...sql.ConflictOption) *QuestionCategoryUpsertBulk {
	qccb.conflict = opts
	return &QuestionCategoryUpsertBulk{
		create: qccb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.QuestionCategory.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (qccb *QuestionCategoryCreateBulk) OnConflictColumns(columns ...string) *QuestionCategoryUpsertBulk {
	qccb.conflict = append(qccb.conflict, sql.ConflictColumns(columns...))
	return &QuestionCategoryUpsertBulk{
		create: qccb,
	}
}

// QuestionCategoryUpsertBulk is the builder for "upsert"-ing
// a bulk of QuestionCategory nodes.
type QuestionCategoryUpsertBulk struct {
	create *QuestionCategoryCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.QuestionCategory.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *QuestionCategoryUpsertBulk) UpdateNewValues() *QuestionCategoryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(questioncategory.FieldCreatedAt)
			}
			if _, exists := b.mutation.Creator(); exists {
				s.SetIgnore(questioncategory.FieldCreator)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.QuestionCategory.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *QuestionCategoryUpsertBulk) Ignore() *QuestionCategoryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *QuestionCategoryUpsertBulk) DoNothing() *QuestionCategoryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the QuestionCategoryCreateBulk.OnConflict
// documentation for more info.
func (u *QuestionCategoryUpsertBulk) Update(set func(*QuestionCategoryUpsert)) *QuestionCategoryUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&QuestionCategoryUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *QuestionCategoryUpsertBulk) SetUpdatedAt(v time.Time) *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *QuestionCategoryUpsertBulk) UpdateUpdatedAt() *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *QuestionCategoryUpsertBulk) SetDeletedAt(v time.Time) *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *QuestionCategoryUpsertBulk) UpdateDeletedAt() *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *QuestionCategoryUpsertBulk) ClearDeletedAt() *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *QuestionCategoryUpsertBulk) SetLastModifier(v *model.Modifier) *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *QuestionCategoryUpsertBulk) UpdateLastModifier() *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *QuestionCategoryUpsertBulk) ClearLastModifier() *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *QuestionCategoryUpsertBulk) SetRemark(v string) *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *QuestionCategoryUpsertBulk) UpdateRemark() *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *QuestionCategoryUpsertBulk) ClearRemark() *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.ClearRemark()
	})
}

// SetName sets the "name" field.
func (u *QuestionCategoryUpsertBulk) SetName(v string) *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *QuestionCategoryUpsertBulk) UpdateName() *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateName()
	})
}

// SetSort sets the "sort" field.
func (u *QuestionCategoryUpsertBulk) SetSort(v uint64) *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.SetSort(v)
	})
}

// AddSort adds v to the "sort" field.
func (u *QuestionCategoryUpsertBulk) AddSort(v uint64) *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.AddSort(v)
	})
}

// UpdateSort sets the "sort" field to the value that was provided on create.
func (u *QuestionCategoryUpsertBulk) UpdateSort() *QuestionCategoryUpsertBulk {
	return u.Update(func(s *QuestionCategoryUpsert) {
		s.UpdateSort()
	})
}

// Exec executes the query.
func (u *QuestionCategoryUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the QuestionCategoryCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for QuestionCategoryCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *QuestionCategoryUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
