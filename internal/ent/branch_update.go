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
	"github.com/auroraride/aurservd/internal/ent/branch"
	"github.com/auroraride/aurservd/internal/ent/branchcontract"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// BranchUpdate is the builder for updating Branch entities.
type BranchUpdate struct {
	config
	hooks    []Hook
	mutation *BranchMutation
}

// Where appends a list predicates to the BranchUpdate builder.
func (bu *BranchUpdate) Where(ps ...predicate.Branch) *BranchUpdate {
	bu.mutation.Where(ps...)
	return bu
}

// SetUpdatedAt sets the "updated_at" field.
func (bu *BranchUpdate) SetUpdatedAt(t time.Time) *BranchUpdate {
	bu.mutation.SetUpdatedAt(t)
	return bu
}

// SetDeletedAt sets the "deleted_at" field.
func (bu *BranchUpdate) SetDeletedAt(t time.Time) *BranchUpdate {
	bu.mutation.SetDeletedAt(t)
	return bu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (bu *BranchUpdate) SetNillableDeletedAt(t *time.Time) *BranchUpdate {
	if t != nil {
		bu.SetDeletedAt(*t)
	}
	return bu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (bu *BranchUpdate) ClearDeletedAt() *BranchUpdate {
	bu.mutation.ClearDeletedAt()
	return bu
}

// SetCreator sets the "creator" field.
func (bu *BranchUpdate) SetCreator(m *model.Modifier) *BranchUpdate {
	bu.mutation.SetCreator(m)
	return bu
}

// ClearCreator clears the value of the "creator" field.
func (bu *BranchUpdate) ClearCreator() *BranchUpdate {
	bu.mutation.ClearCreator()
	return bu
}

// SetLastModifier sets the "last_modifier" field.
func (bu *BranchUpdate) SetLastModifier(m *model.Modifier) *BranchUpdate {
	bu.mutation.SetLastModifier(m)
	return bu
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (bu *BranchUpdate) ClearLastModifier() *BranchUpdate {
	bu.mutation.ClearLastModifier()
	return bu
}

// SetRemark sets the "remark" field.
func (bu *BranchUpdate) SetRemark(s string) *BranchUpdate {
	bu.mutation.SetRemark(s)
	return bu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (bu *BranchUpdate) SetNillableRemark(s *string) *BranchUpdate {
	if s != nil {
		bu.SetRemark(*s)
	}
	return bu
}

// ClearRemark clears the value of the "remark" field.
func (bu *BranchUpdate) ClearRemark() *BranchUpdate {
	bu.mutation.ClearRemark()
	return bu
}

// SetCityID sets the "city_id" field.
func (bu *BranchUpdate) SetCityID(u uint64) *BranchUpdate {
	bu.mutation.ResetCityID()
	bu.mutation.SetCityID(u)
	return bu
}

// AddCityID adds u to the "city_id" field.
func (bu *BranchUpdate) AddCityID(u int64) *BranchUpdate {
	bu.mutation.AddCityID(u)
	return bu
}

// SetName sets the "name" field.
func (bu *BranchUpdate) SetName(s string) *BranchUpdate {
	bu.mutation.SetName(s)
	return bu
}

// SetLng sets the "lng" field.
func (bu *BranchUpdate) SetLng(f float64) *BranchUpdate {
	bu.mutation.ResetLng()
	bu.mutation.SetLng(f)
	return bu
}

// AddLng adds f to the "lng" field.
func (bu *BranchUpdate) AddLng(f float64) *BranchUpdate {
	bu.mutation.AddLng(f)
	return bu
}

// SetLat sets the "lat" field.
func (bu *BranchUpdate) SetLat(f float64) *BranchUpdate {
	bu.mutation.ResetLat()
	bu.mutation.SetLat(f)
	return bu
}

// AddLat adds f to the "lat" field.
func (bu *BranchUpdate) AddLat(f float64) *BranchUpdate {
	bu.mutation.AddLat(f)
	return bu
}

// SetAddress sets the "address" field.
func (bu *BranchUpdate) SetAddress(s string) *BranchUpdate {
	bu.mutation.SetAddress(s)
	return bu
}

// SetPhotos sets the "photos" field.
func (bu *BranchUpdate) SetPhotos(s []string) *BranchUpdate {
	bu.mutation.SetPhotos(s)
	return bu
}

// AddContractIDs adds the "contracts" edge to the BranchContract entity by IDs.
func (bu *BranchUpdate) AddContractIDs(ids ...uint64) *BranchUpdate {
	bu.mutation.AddContractIDs(ids...)
	return bu
}

// AddContracts adds the "contracts" edges to the BranchContract entity.
func (bu *BranchUpdate) AddContracts(b ...*BranchContract) *BranchUpdate {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return bu.AddContractIDs(ids...)
}

// Mutation returns the BranchMutation object of the builder.
func (bu *BranchUpdate) Mutation() *BranchMutation {
	return bu.mutation
}

// ClearContracts clears all "contracts" edges to the BranchContract entity.
func (bu *BranchUpdate) ClearContracts() *BranchUpdate {
	bu.mutation.ClearContracts()
	return bu
}

// RemoveContractIDs removes the "contracts" edge to BranchContract entities by IDs.
func (bu *BranchUpdate) RemoveContractIDs(ids ...uint64) *BranchUpdate {
	bu.mutation.RemoveContractIDs(ids...)
	return bu
}

// RemoveContracts removes "contracts" edges to BranchContract entities.
func (bu *BranchUpdate) RemoveContracts(b ...*BranchContract) *BranchUpdate {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return bu.RemoveContractIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (bu *BranchUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	bu.defaults()
	if len(bu.hooks) == 0 {
		affected, err = bu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*BranchMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			bu.mutation = mutation
			affected, err = bu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(bu.hooks) - 1; i >= 0; i-- {
			if bu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = bu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, bu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (bu *BranchUpdate) SaveX(ctx context.Context) int {
	affected, err := bu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (bu *BranchUpdate) Exec(ctx context.Context) error {
	_, err := bu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bu *BranchUpdate) ExecX(ctx context.Context) {
	if err := bu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (bu *BranchUpdate) defaults() {
	if _, ok := bu.mutation.UpdatedAt(); !ok {
		v := branch.UpdateDefaultUpdatedAt()
		bu.mutation.SetUpdatedAt(v)
	}
}

func (bu *BranchUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   branch.Table,
			Columns: branch.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: branch.FieldID,
			},
		},
	}
	if ps := bu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := bu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: branch.FieldUpdatedAt,
		})
	}
	if value, ok := bu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: branch.FieldDeletedAt,
		})
	}
	if bu.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: branch.FieldDeletedAt,
		})
	}
	if value, ok := bu.mutation.Creator(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: branch.FieldCreator,
		})
	}
	if bu.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: branch.FieldCreator,
		})
	}
	if value, ok := bu.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: branch.FieldLastModifier,
		})
	}
	if bu.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: branch.FieldLastModifier,
		})
	}
	if value, ok := bu.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: branch.FieldRemark,
		})
	}
	if bu.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: branch.FieldRemark,
		})
	}
	if value, ok := bu.mutation.CityID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint64,
			Value:  value,
			Column: branch.FieldCityID,
		})
	}
	if value, ok := bu.mutation.AddedCityID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint64,
			Value:  value,
			Column: branch.FieldCityID,
		})
	}
	if value, ok := bu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: branch.FieldName,
		})
	}
	if value, ok := bu.mutation.Lng(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: branch.FieldLng,
		})
	}
	if value, ok := bu.mutation.AddedLng(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: branch.FieldLng,
		})
	}
	if value, ok := bu.mutation.Lat(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: branch.FieldLat,
		})
	}
	if value, ok := bu.mutation.AddedLat(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: branch.FieldLat,
		})
	}
	if value, ok := bu.mutation.Address(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: branch.FieldAddress,
		})
	}
	if value, ok := bu.mutation.Photos(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: branch.FieldPhotos,
		})
	}
	if bu.mutation.ContractsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   branch.ContractsTable,
			Columns: []string{branch.ContractsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branchcontract.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.RemovedContractsIDs(); len(nodes) > 0 && !bu.mutation.ContractsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   branch.ContractsTable,
			Columns: []string{branch.ContractsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branchcontract.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.ContractsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   branch.ContractsTable,
			Columns: []string{branch.ContractsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branchcontract.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, bu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{branch.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// BranchUpdateOne is the builder for updating a single Branch entity.
type BranchUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *BranchMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (buo *BranchUpdateOne) SetUpdatedAt(t time.Time) *BranchUpdateOne {
	buo.mutation.SetUpdatedAt(t)
	return buo
}

// SetDeletedAt sets the "deleted_at" field.
func (buo *BranchUpdateOne) SetDeletedAt(t time.Time) *BranchUpdateOne {
	buo.mutation.SetDeletedAt(t)
	return buo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (buo *BranchUpdateOne) SetNillableDeletedAt(t *time.Time) *BranchUpdateOne {
	if t != nil {
		buo.SetDeletedAt(*t)
	}
	return buo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (buo *BranchUpdateOne) ClearDeletedAt() *BranchUpdateOne {
	buo.mutation.ClearDeletedAt()
	return buo
}

// SetCreator sets the "creator" field.
func (buo *BranchUpdateOne) SetCreator(m *model.Modifier) *BranchUpdateOne {
	buo.mutation.SetCreator(m)
	return buo
}

// ClearCreator clears the value of the "creator" field.
func (buo *BranchUpdateOne) ClearCreator() *BranchUpdateOne {
	buo.mutation.ClearCreator()
	return buo
}

// SetLastModifier sets the "last_modifier" field.
func (buo *BranchUpdateOne) SetLastModifier(m *model.Modifier) *BranchUpdateOne {
	buo.mutation.SetLastModifier(m)
	return buo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (buo *BranchUpdateOne) ClearLastModifier() *BranchUpdateOne {
	buo.mutation.ClearLastModifier()
	return buo
}

// SetRemark sets the "remark" field.
func (buo *BranchUpdateOne) SetRemark(s string) *BranchUpdateOne {
	buo.mutation.SetRemark(s)
	return buo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (buo *BranchUpdateOne) SetNillableRemark(s *string) *BranchUpdateOne {
	if s != nil {
		buo.SetRemark(*s)
	}
	return buo
}

// ClearRemark clears the value of the "remark" field.
func (buo *BranchUpdateOne) ClearRemark() *BranchUpdateOne {
	buo.mutation.ClearRemark()
	return buo
}

// SetCityID sets the "city_id" field.
func (buo *BranchUpdateOne) SetCityID(u uint64) *BranchUpdateOne {
	buo.mutation.ResetCityID()
	buo.mutation.SetCityID(u)
	return buo
}

// AddCityID adds u to the "city_id" field.
func (buo *BranchUpdateOne) AddCityID(u int64) *BranchUpdateOne {
	buo.mutation.AddCityID(u)
	return buo
}

// SetName sets the "name" field.
func (buo *BranchUpdateOne) SetName(s string) *BranchUpdateOne {
	buo.mutation.SetName(s)
	return buo
}

// SetLng sets the "lng" field.
func (buo *BranchUpdateOne) SetLng(f float64) *BranchUpdateOne {
	buo.mutation.ResetLng()
	buo.mutation.SetLng(f)
	return buo
}

// AddLng adds f to the "lng" field.
func (buo *BranchUpdateOne) AddLng(f float64) *BranchUpdateOne {
	buo.mutation.AddLng(f)
	return buo
}

// SetLat sets the "lat" field.
func (buo *BranchUpdateOne) SetLat(f float64) *BranchUpdateOne {
	buo.mutation.ResetLat()
	buo.mutation.SetLat(f)
	return buo
}

// AddLat adds f to the "lat" field.
func (buo *BranchUpdateOne) AddLat(f float64) *BranchUpdateOne {
	buo.mutation.AddLat(f)
	return buo
}

// SetAddress sets the "address" field.
func (buo *BranchUpdateOne) SetAddress(s string) *BranchUpdateOne {
	buo.mutation.SetAddress(s)
	return buo
}

// SetPhotos sets the "photos" field.
func (buo *BranchUpdateOne) SetPhotos(s []string) *BranchUpdateOne {
	buo.mutation.SetPhotos(s)
	return buo
}

// AddContractIDs adds the "contracts" edge to the BranchContract entity by IDs.
func (buo *BranchUpdateOne) AddContractIDs(ids ...uint64) *BranchUpdateOne {
	buo.mutation.AddContractIDs(ids...)
	return buo
}

// AddContracts adds the "contracts" edges to the BranchContract entity.
func (buo *BranchUpdateOne) AddContracts(b ...*BranchContract) *BranchUpdateOne {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return buo.AddContractIDs(ids...)
}

// Mutation returns the BranchMutation object of the builder.
func (buo *BranchUpdateOne) Mutation() *BranchMutation {
	return buo.mutation
}

// ClearContracts clears all "contracts" edges to the BranchContract entity.
func (buo *BranchUpdateOne) ClearContracts() *BranchUpdateOne {
	buo.mutation.ClearContracts()
	return buo
}

// RemoveContractIDs removes the "contracts" edge to BranchContract entities by IDs.
func (buo *BranchUpdateOne) RemoveContractIDs(ids ...uint64) *BranchUpdateOne {
	buo.mutation.RemoveContractIDs(ids...)
	return buo
}

// RemoveContracts removes "contracts" edges to BranchContract entities.
func (buo *BranchUpdateOne) RemoveContracts(b ...*BranchContract) *BranchUpdateOne {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return buo.RemoveContractIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (buo *BranchUpdateOne) Select(field string, fields ...string) *BranchUpdateOne {
	buo.fields = append([]string{field}, fields...)
	return buo
}

// Save executes the query and returns the updated Branch entity.
func (buo *BranchUpdateOne) Save(ctx context.Context) (*Branch, error) {
	var (
		err  error
		node *Branch
	)
	buo.defaults()
	if len(buo.hooks) == 0 {
		node, err = buo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*BranchMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			buo.mutation = mutation
			node, err = buo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(buo.hooks) - 1; i >= 0; i-- {
			if buo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = buo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, buo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (buo *BranchUpdateOne) SaveX(ctx context.Context) *Branch {
	node, err := buo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (buo *BranchUpdateOne) Exec(ctx context.Context) error {
	_, err := buo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (buo *BranchUpdateOne) ExecX(ctx context.Context) {
	if err := buo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (buo *BranchUpdateOne) defaults() {
	if _, ok := buo.mutation.UpdatedAt(); !ok {
		v := branch.UpdateDefaultUpdatedAt()
		buo.mutation.SetUpdatedAt(v)
	}
}

func (buo *BranchUpdateOne) sqlSave(ctx context.Context) (_node *Branch, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   branch.Table,
			Columns: branch.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: branch.FieldID,
			},
		},
	}
	id, ok := buo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Branch.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := buo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, branch.FieldID)
		for _, f := range fields {
			if !branch.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != branch.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := buo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := buo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: branch.FieldUpdatedAt,
		})
	}
	if value, ok := buo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: branch.FieldDeletedAt,
		})
	}
	if buo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: branch.FieldDeletedAt,
		})
	}
	if value, ok := buo.mutation.Creator(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: branch.FieldCreator,
		})
	}
	if buo.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: branch.FieldCreator,
		})
	}
	if value, ok := buo.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: branch.FieldLastModifier,
		})
	}
	if buo.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: branch.FieldLastModifier,
		})
	}
	if value, ok := buo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: branch.FieldRemark,
		})
	}
	if buo.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: branch.FieldRemark,
		})
	}
	if value, ok := buo.mutation.CityID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint64,
			Value:  value,
			Column: branch.FieldCityID,
		})
	}
	if value, ok := buo.mutation.AddedCityID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint64,
			Value:  value,
			Column: branch.FieldCityID,
		})
	}
	if value, ok := buo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: branch.FieldName,
		})
	}
	if value, ok := buo.mutation.Lng(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: branch.FieldLng,
		})
	}
	if value, ok := buo.mutation.AddedLng(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: branch.FieldLng,
		})
	}
	if value, ok := buo.mutation.Lat(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: branch.FieldLat,
		})
	}
	if value, ok := buo.mutation.AddedLat(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: branch.FieldLat,
		})
	}
	if value, ok := buo.mutation.Address(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: branch.FieldAddress,
		})
	}
	if value, ok := buo.mutation.Photos(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: branch.FieldPhotos,
		})
	}
	if buo.mutation.ContractsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   branch.ContractsTable,
			Columns: []string{branch.ContractsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branchcontract.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.RemovedContractsIDs(); len(nodes) > 0 && !buo.mutation.ContractsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   branch.ContractsTable,
			Columns: []string{branch.ContractsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branchcontract.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.ContractsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   branch.ContractsTable,
			Columns: []string{branch.ContractsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branchcontract.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Branch{config: buo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, buo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{branch.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}