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
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/branch"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// CabinetUpdate is the builder for updating Cabinet entities.
type CabinetUpdate struct {
	config
	hooks    []Hook
	mutation *CabinetMutation
}

// Where appends a list predicates to the CabinetUpdate builder.
func (cu *CabinetUpdate) Where(ps ...predicate.Cabinet) *CabinetUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetUpdatedAt sets the "updated_at" field.
func (cu *CabinetUpdate) SetUpdatedAt(t time.Time) *CabinetUpdate {
	cu.mutation.SetUpdatedAt(t)
	return cu
}

// SetDeletedAt sets the "deleted_at" field.
func (cu *CabinetUpdate) SetDeletedAt(t time.Time) *CabinetUpdate {
	cu.mutation.SetDeletedAt(t)
	return cu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (cu *CabinetUpdate) SetNillableDeletedAt(t *time.Time) *CabinetUpdate {
	if t != nil {
		cu.SetDeletedAt(*t)
	}
	return cu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (cu *CabinetUpdate) ClearDeletedAt() *CabinetUpdate {
	cu.mutation.ClearDeletedAt()
	return cu
}

// SetCreator sets the "creator" field.
func (cu *CabinetUpdate) SetCreator(m *model.Modifier) *CabinetUpdate {
	cu.mutation.SetCreator(m)
	return cu
}

// ClearCreator clears the value of the "creator" field.
func (cu *CabinetUpdate) ClearCreator() *CabinetUpdate {
	cu.mutation.ClearCreator()
	return cu
}

// SetLastModifier sets the "last_modifier" field.
func (cu *CabinetUpdate) SetLastModifier(m *model.Modifier) *CabinetUpdate {
	cu.mutation.SetLastModifier(m)
	return cu
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (cu *CabinetUpdate) ClearLastModifier() *CabinetUpdate {
	cu.mutation.ClearLastModifier()
	return cu
}

// SetRemark sets the "remark" field.
func (cu *CabinetUpdate) SetRemark(s string) *CabinetUpdate {
	cu.mutation.SetRemark(s)
	return cu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (cu *CabinetUpdate) SetNillableRemark(s *string) *CabinetUpdate {
	if s != nil {
		cu.SetRemark(*s)
	}
	return cu
}

// ClearRemark clears the value of the "remark" field.
func (cu *CabinetUpdate) ClearRemark() *CabinetUpdate {
	cu.mutation.ClearRemark()
	return cu
}

// SetBranchID sets the "branch_id" field.
func (cu *CabinetUpdate) SetBranchID(u uint64) *CabinetUpdate {
	cu.mutation.SetBranchID(u)
	return cu
}

// SetNillableBranchID sets the "branch_id" field if the given value is not nil.
func (cu *CabinetUpdate) SetNillableBranchID(u *uint64) *CabinetUpdate {
	if u != nil {
		cu.SetBranchID(*u)
	}
	return cu
}

// ClearBranchID clears the value of the "branch_id" field.
func (cu *CabinetUpdate) ClearBranchID() *CabinetUpdate {
	cu.mutation.ClearBranchID()
	return cu
}

// SetSn sets the "sn" field.
func (cu *CabinetUpdate) SetSn(s string) *CabinetUpdate {
	cu.mutation.SetSn(s)
	return cu
}

// SetBrand sets the "brand" field.
func (cu *CabinetUpdate) SetBrand(s string) *CabinetUpdate {
	cu.mutation.SetBrand(s)
	return cu
}

// SetSerial sets the "serial" field.
func (cu *CabinetUpdate) SetSerial(s string) *CabinetUpdate {
	cu.mutation.SetSerial(s)
	return cu
}

// SetName sets the "name" field.
func (cu *CabinetUpdate) SetName(s string) *CabinetUpdate {
	cu.mutation.SetName(s)
	return cu
}

// SetDoors sets the "doors" field.
func (cu *CabinetUpdate) SetDoors(u uint) *CabinetUpdate {
	cu.mutation.ResetDoors()
	cu.mutation.SetDoors(u)
	return cu
}

// AddDoors adds u to the "doors" field.
func (cu *CabinetUpdate) AddDoors(u int) *CabinetUpdate {
	cu.mutation.AddDoors(u)
	return cu
}

// SetStatus sets the "status" field.
func (cu *CabinetUpdate) SetStatus(u uint) *CabinetUpdate {
	cu.mutation.ResetStatus()
	cu.mutation.SetStatus(u)
	return cu
}

// AddStatus adds u to the "status" field.
func (cu *CabinetUpdate) AddStatus(u int) *CabinetUpdate {
	cu.mutation.AddStatus(u)
	return cu
}

// SetModels sets the "models" field.
func (cu *CabinetUpdate) SetModels(mm []model.BatteryModel) *CabinetUpdate {
	cu.mutation.SetModels(mm)
	return cu
}

// SetHealth sets the "health" field.
func (cu *CabinetUpdate) SetHealth(u uint) *CabinetUpdate {
	cu.mutation.ResetHealth()
	cu.mutation.SetHealth(u)
	return cu
}

// AddHealth adds u to the "health" field.
func (cu *CabinetUpdate) AddHealth(u int) *CabinetUpdate {
	cu.mutation.AddHealth(u)
	return cu
}

// SetBin sets the "bin" field.
func (cu *CabinetUpdate) SetBin(mb []model.CabinetBin) *CabinetUpdate {
	cu.mutation.SetBin(mb)
	return cu
}

// ClearBin clears the value of the "bin" field.
func (cu *CabinetUpdate) ClearBin() *CabinetUpdate {
	cu.mutation.ClearBin()
	return cu
}

// SetBatteryNum sets the "battery_num" field.
func (cu *CabinetUpdate) SetBatteryNum(u uint) *CabinetUpdate {
	cu.mutation.ResetBatteryNum()
	cu.mutation.SetBatteryNum(u)
	return cu
}

// SetNillableBatteryNum sets the "battery_num" field if the given value is not nil.
func (cu *CabinetUpdate) SetNillableBatteryNum(u *uint) *CabinetUpdate {
	if u != nil {
		cu.SetBatteryNum(*u)
	}
	return cu
}

// AddBatteryNum adds u to the "battery_num" field.
func (cu *CabinetUpdate) AddBatteryNum(u int) *CabinetUpdate {
	cu.mutation.AddBatteryNum(u)
	return cu
}

// SetBatteryFullNum sets the "battery_full_num" field.
func (cu *CabinetUpdate) SetBatteryFullNum(u uint) *CabinetUpdate {
	cu.mutation.ResetBatteryFullNum()
	cu.mutation.SetBatteryFullNum(u)
	return cu
}

// SetNillableBatteryFullNum sets the "battery_full_num" field if the given value is not nil.
func (cu *CabinetUpdate) SetNillableBatteryFullNum(u *uint) *CabinetUpdate {
	if u != nil {
		cu.SetBatteryFullNum(*u)
	}
	return cu
}

// AddBatteryFullNum adds u to the "battery_full_num" field.
func (cu *CabinetUpdate) AddBatteryFullNum(u int) *CabinetUpdate {
	cu.mutation.AddBatteryFullNum(u)
	return cu
}

// SetBranch sets the "branch" edge to the Branch entity.
func (cu *CabinetUpdate) SetBranch(b *Branch) *CabinetUpdate {
	return cu.SetBranchID(b.ID)
}

// AddBmIDs adds the "bms" edge to the BatteryModel entity by IDs.
func (cu *CabinetUpdate) AddBmIDs(ids ...uint64) *CabinetUpdate {
	cu.mutation.AddBmIDs(ids...)
	return cu
}

// AddBms adds the "bms" edges to the BatteryModel entity.
func (cu *CabinetUpdate) AddBms(b ...*BatteryModel) *CabinetUpdate {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return cu.AddBmIDs(ids...)
}

// Mutation returns the CabinetMutation object of the builder.
func (cu *CabinetUpdate) Mutation() *CabinetMutation {
	return cu.mutation
}

// ClearBranch clears the "branch" edge to the Branch entity.
func (cu *CabinetUpdate) ClearBranch() *CabinetUpdate {
	cu.mutation.ClearBranch()
	return cu
}

// ClearBms clears all "bms" edges to the BatteryModel entity.
func (cu *CabinetUpdate) ClearBms() *CabinetUpdate {
	cu.mutation.ClearBms()
	return cu
}

// RemoveBmIDs removes the "bms" edge to BatteryModel entities by IDs.
func (cu *CabinetUpdate) RemoveBmIDs(ids ...uint64) *CabinetUpdate {
	cu.mutation.RemoveBmIDs(ids...)
	return cu
}

// RemoveBms removes "bms" edges to BatteryModel entities.
func (cu *CabinetUpdate) RemoveBms(b ...*BatteryModel) *CabinetUpdate {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return cu.RemoveBmIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *CabinetUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	cu.defaults()
	if len(cu.hooks) == 0 {
		affected, err = cu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CabinetMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			cu.mutation = mutation
			affected, err = cu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(cu.hooks) - 1; i >= 0; i-- {
			if cu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = cu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CabinetUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CabinetUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CabinetUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cu *CabinetUpdate) defaults() {
	if _, ok := cu.mutation.UpdatedAt(); !ok {
		v := cabinet.UpdateDefaultUpdatedAt()
		cu.mutation.SetUpdatedAt(v)
	}
}

func (cu *CabinetUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cabinet.Table,
			Columns: cabinet.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: cabinet.FieldID,
			},
		},
	}
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cabinet.FieldUpdatedAt,
		})
	}
	if value, ok := cu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cabinet.FieldDeletedAt,
		})
	}
	if cu.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: cabinet.FieldDeletedAt,
		})
	}
	if value, ok := cu.mutation.Creator(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: cabinet.FieldCreator,
		})
	}
	if cu.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: cabinet.FieldCreator,
		})
	}
	if value, ok := cu.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: cabinet.FieldLastModifier,
		})
	}
	if cu.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: cabinet.FieldLastModifier,
		})
	}
	if value, ok := cu.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cabinet.FieldRemark,
		})
	}
	if cu.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: cabinet.FieldRemark,
		})
	}
	if value, ok := cu.mutation.Sn(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cabinet.FieldSn,
		})
	}
	if value, ok := cu.mutation.Brand(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cabinet.FieldBrand,
		})
	}
	if value, ok := cu.mutation.Serial(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cabinet.FieldSerial,
		})
	}
	if value, ok := cu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cabinet.FieldName,
		})
	}
	if value, ok := cu.mutation.Doors(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldDoors,
		})
	}
	if value, ok := cu.mutation.AddedDoors(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldDoors,
		})
	}
	if value, ok := cu.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldStatus,
		})
	}
	if value, ok := cu.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldStatus,
		})
	}
	if value, ok := cu.mutation.Models(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: cabinet.FieldModels,
		})
	}
	if value, ok := cu.mutation.Health(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldHealth,
		})
	}
	if value, ok := cu.mutation.AddedHealth(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldHealth,
		})
	}
	if value, ok := cu.mutation.Bin(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: cabinet.FieldBin,
		})
	}
	if cu.mutation.BinCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: cabinet.FieldBin,
		})
	}
	if value, ok := cu.mutation.BatteryNum(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldBatteryNum,
		})
	}
	if value, ok := cu.mutation.AddedBatteryNum(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldBatteryNum,
		})
	}
	if value, ok := cu.mutation.BatteryFullNum(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldBatteryFullNum,
		})
	}
	if value, ok := cu.mutation.AddedBatteryFullNum(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldBatteryFullNum,
		})
	}
	if cu.mutation.BranchCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   cabinet.BranchTable,
			Columns: []string{cabinet.BranchColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branch.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.BranchIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   cabinet.BranchTable,
			Columns: []string{cabinet.BranchColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branch.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cu.mutation.BmsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cabinet.BmsTable,
			Columns: cabinet.BmsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: batterymodel.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedBmsIDs(); len(nodes) > 0 && !cu.mutation.BmsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cabinet.BmsTable,
			Columns: cabinet.BmsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: batterymodel.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.BmsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cabinet.BmsTable,
			Columns: cabinet.BmsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: batterymodel.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cabinet.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// CabinetUpdateOne is the builder for updating a single Cabinet entity.
type CabinetUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *CabinetMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (cuo *CabinetUpdateOne) SetUpdatedAt(t time.Time) *CabinetUpdateOne {
	cuo.mutation.SetUpdatedAt(t)
	return cuo
}

// SetDeletedAt sets the "deleted_at" field.
func (cuo *CabinetUpdateOne) SetDeletedAt(t time.Time) *CabinetUpdateOne {
	cuo.mutation.SetDeletedAt(t)
	return cuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (cuo *CabinetUpdateOne) SetNillableDeletedAt(t *time.Time) *CabinetUpdateOne {
	if t != nil {
		cuo.SetDeletedAt(*t)
	}
	return cuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (cuo *CabinetUpdateOne) ClearDeletedAt() *CabinetUpdateOne {
	cuo.mutation.ClearDeletedAt()
	return cuo
}

// SetCreator sets the "creator" field.
func (cuo *CabinetUpdateOne) SetCreator(m *model.Modifier) *CabinetUpdateOne {
	cuo.mutation.SetCreator(m)
	return cuo
}

// ClearCreator clears the value of the "creator" field.
func (cuo *CabinetUpdateOne) ClearCreator() *CabinetUpdateOne {
	cuo.mutation.ClearCreator()
	return cuo
}

// SetLastModifier sets the "last_modifier" field.
func (cuo *CabinetUpdateOne) SetLastModifier(m *model.Modifier) *CabinetUpdateOne {
	cuo.mutation.SetLastModifier(m)
	return cuo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (cuo *CabinetUpdateOne) ClearLastModifier() *CabinetUpdateOne {
	cuo.mutation.ClearLastModifier()
	return cuo
}

// SetRemark sets the "remark" field.
func (cuo *CabinetUpdateOne) SetRemark(s string) *CabinetUpdateOne {
	cuo.mutation.SetRemark(s)
	return cuo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (cuo *CabinetUpdateOne) SetNillableRemark(s *string) *CabinetUpdateOne {
	if s != nil {
		cuo.SetRemark(*s)
	}
	return cuo
}

// ClearRemark clears the value of the "remark" field.
func (cuo *CabinetUpdateOne) ClearRemark() *CabinetUpdateOne {
	cuo.mutation.ClearRemark()
	return cuo
}

// SetBranchID sets the "branch_id" field.
func (cuo *CabinetUpdateOne) SetBranchID(u uint64) *CabinetUpdateOne {
	cuo.mutation.SetBranchID(u)
	return cuo
}

// SetNillableBranchID sets the "branch_id" field if the given value is not nil.
func (cuo *CabinetUpdateOne) SetNillableBranchID(u *uint64) *CabinetUpdateOne {
	if u != nil {
		cuo.SetBranchID(*u)
	}
	return cuo
}

// ClearBranchID clears the value of the "branch_id" field.
func (cuo *CabinetUpdateOne) ClearBranchID() *CabinetUpdateOne {
	cuo.mutation.ClearBranchID()
	return cuo
}

// SetSn sets the "sn" field.
func (cuo *CabinetUpdateOne) SetSn(s string) *CabinetUpdateOne {
	cuo.mutation.SetSn(s)
	return cuo
}

// SetBrand sets the "brand" field.
func (cuo *CabinetUpdateOne) SetBrand(s string) *CabinetUpdateOne {
	cuo.mutation.SetBrand(s)
	return cuo
}

// SetSerial sets the "serial" field.
func (cuo *CabinetUpdateOne) SetSerial(s string) *CabinetUpdateOne {
	cuo.mutation.SetSerial(s)
	return cuo
}

// SetName sets the "name" field.
func (cuo *CabinetUpdateOne) SetName(s string) *CabinetUpdateOne {
	cuo.mutation.SetName(s)
	return cuo
}

// SetDoors sets the "doors" field.
func (cuo *CabinetUpdateOne) SetDoors(u uint) *CabinetUpdateOne {
	cuo.mutation.ResetDoors()
	cuo.mutation.SetDoors(u)
	return cuo
}

// AddDoors adds u to the "doors" field.
func (cuo *CabinetUpdateOne) AddDoors(u int) *CabinetUpdateOne {
	cuo.mutation.AddDoors(u)
	return cuo
}

// SetStatus sets the "status" field.
func (cuo *CabinetUpdateOne) SetStatus(u uint) *CabinetUpdateOne {
	cuo.mutation.ResetStatus()
	cuo.mutation.SetStatus(u)
	return cuo
}

// AddStatus adds u to the "status" field.
func (cuo *CabinetUpdateOne) AddStatus(u int) *CabinetUpdateOne {
	cuo.mutation.AddStatus(u)
	return cuo
}

// SetModels sets the "models" field.
func (cuo *CabinetUpdateOne) SetModels(mm []model.BatteryModel) *CabinetUpdateOne {
	cuo.mutation.SetModels(mm)
	return cuo
}

// SetHealth sets the "health" field.
func (cuo *CabinetUpdateOne) SetHealth(u uint) *CabinetUpdateOne {
	cuo.mutation.ResetHealth()
	cuo.mutation.SetHealth(u)
	return cuo
}

// AddHealth adds u to the "health" field.
func (cuo *CabinetUpdateOne) AddHealth(u int) *CabinetUpdateOne {
	cuo.mutation.AddHealth(u)
	return cuo
}

// SetBin sets the "bin" field.
func (cuo *CabinetUpdateOne) SetBin(mb []model.CabinetBin) *CabinetUpdateOne {
	cuo.mutation.SetBin(mb)
	return cuo
}

// ClearBin clears the value of the "bin" field.
func (cuo *CabinetUpdateOne) ClearBin() *CabinetUpdateOne {
	cuo.mutation.ClearBin()
	return cuo
}

// SetBatteryNum sets the "battery_num" field.
func (cuo *CabinetUpdateOne) SetBatteryNum(u uint) *CabinetUpdateOne {
	cuo.mutation.ResetBatteryNum()
	cuo.mutation.SetBatteryNum(u)
	return cuo
}

// SetNillableBatteryNum sets the "battery_num" field if the given value is not nil.
func (cuo *CabinetUpdateOne) SetNillableBatteryNum(u *uint) *CabinetUpdateOne {
	if u != nil {
		cuo.SetBatteryNum(*u)
	}
	return cuo
}

// AddBatteryNum adds u to the "battery_num" field.
func (cuo *CabinetUpdateOne) AddBatteryNum(u int) *CabinetUpdateOne {
	cuo.mutation.AddBatteryNum(u)
	return cuo
}

// SetBatteryFullNum sets the "battery_full_num" field.
func (cuo *CabinetUpdateOne) SetBatteryFullNum(u uint) *CabinetUpdateOne {
	cuo.mutation.ResetBatteryFullNum()
	cuo.mutation.SetBatteryFullNum(u)
	return cuo
}

// SetNillableBatteryFullNum sets the "battery_full_num" field if the given value is not nil.
func (cuo *CabinetUpdateOne) SetNillableBatteryFullNum(u *uint) *CabinetUpdateOne {
	if u != nil {
		cuo.SetBatteryFullNum(*u)
	}
	return cuo
}

// AddBatteryFullNum adds u to the "battery_full_num" field.
func (cuo *CabinetUpdateOne) AddBatteryFullNum(u int) *CabinetUpdateOne {
	cuo.mutation.AddBatteryFullNum(u)
	return cuo
}

// SetBranch sets the "branch" edge to the Branch entity.
func (cuo *CabinetUpdateOne) SetBranch(b *Branch) *CabinetUpdateOne {
	return cuo.SetBranchID(b.ID)
}

// AddBmIDs adds the "bms" edge to the BatteryModel entity by IDs.
func (cuo *CabinetUpdateOne) AddBmIDs(ids ...uint64) *CabinetUpdateOne {
	cuo.mutation.AddBmIDs(ids...)
	return cuo
}

// AddBms adds the "bms" edges to the BatteryModel entity.
func (cuo *CabinetUpdateOne) AddBms(b ...*BatteryModel) *CabinetUpdateOne {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return cuo.AddBmIDs(ids...)
}

// Mutation returns the CabinetMutation object of the builder.
func (cuo *CabinetUpdateOne) Mutation() *CabinetMutation {
	return cuo.mutation
}

// ClearBranch clears the "branch" edge to the Branch entity.
func (cuo *CabinetUpdateOne) ClearBranch() *CabinetUpdateOne {
	cuo.mutation.ClearBranch()
	return cuo
}

// ClearBms clears all "bms" edges to the BatteryModel entity.
func (cuo *CabinetUpdateOne) ClearBms() *CabinetUpdateOne {
	cuo.mutation.ClearBms()
	return cuo
}

// RemoveBmIDs removes the "bms" edge to BatteryModel entities by IDs.
func (cuo *CabinetUpdateOne) RemoveBmIDs(ids ...uint64) *CabinetUpdateOne {
	cuo.mutation.RemoveBmIDs(ids...)
	return cuo
}

// RemoveBms removes "bms" edges to BatteryModel entities.
func (cuo *CabinetUpdateOne) RemoveBms(b ...*BatteryModel) *CabinetUpdateOne {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return cuo.RemoveBmIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *CabinetUpdateOne) Select(field string, fields ...string) *CabinetUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Cabinet entity.
func (cuo *CabinetUpdateOne) Save(ctx context.Context) (*Cabinet, error) {
	var (
		err  error
		node *Cabinet
	)
	cuo.defaults()
	if len(cuo.hooks) == 0 {
		node, err = cuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CabinetMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			cuo.mutation = mutation
			node, err = cuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(cuo.hooks) - 1; i >= 0; i-- {
			if cuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = cuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CabinetUpdateOne) SaveX(ctx context.Context) *Cabinet {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *CabinetUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CabinetUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cuo *CabinetUpdateOne) defaults() {
	if _, ok := cuo.mutation.UpdatedAt(); !ok {
		v := cabinet.UpdateDefaultUpdatedAt()
		cuo.mutation.SetUpdatedAt(v)
	}
}

func (cuo *CabinetUpdateOne) sqlSave(ctx context.Context) (_node *Cabinet, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cabinet.Table,
			Columns: cabinet.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: cabinet.FieldID,
			},
		},
	}
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Cabinet.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cabinet.FieldID)
		for _, f := range fields {
			if !cabinet.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != cabinet.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cabinet.FieldUpdatedAt,
		})
	}
	if value, ok := cuo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cabinet.FieldDeletedAt,
		})
	}
	if cuo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: cabinet.FieldDeletedAt,
		})
	}
	if value, ok := cuo.mutation.Creator(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: cabinet.FieldCreator,
		})
	}
	if cuo.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: cabinet.FieldCreator,
		})
	}
	if value, ok := cuo.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: cabinet.FieldLastModifier,
		})
	}
	if cuo.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: cabinet.FieldLastModifier,
		})
	}
	if value, ok := cuo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cabinet.FieldRemark,
		})
	}
	if cuo.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: cabinet.FieldRemark,
		})
	}
	if value, ok := cuo.mutation.Sn(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cabinet.FieldSn,
		})
	}
	if value, ok := cuo.mutation.Brand(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cabinet.FieldBrand,
		})
	}
	if value, ok := cuo.mutation.Serial(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cabinet.FieldSerial,
		})
	}
	if value, ok := cuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: cabinet.FieldName,
		})
	}
	if value, ok := cuo.mutation.Doors(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldDoors,
		})
	}
	if value, ok := cuo.mutation.AddedDoors(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldDoors,
		})
	}
	if value, ok := cuo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldStatus,
		})
	}
	if value, ok := cuo.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldStatus,
		})
	}
	if value, ok := cuo.mutation.Models(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: cabinet.FieldModels,
		})
	}
	if value, ok := cuo.mutation.Health(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldHealth,
		})
	}
	if value, ok := cuo.mutation.AddedHealth(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldHealth,
		})
	}
	if value, ok := cuo.mutation.Bin(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: cabinet.FieldBin,
		})
	}
	if cuo.mutation.BinCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: cabinet.FieldBin,
		})
	}
	if value, ok := cuo.mutation.BatteryNum(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldBatteryNum,
		})
	}
	if value, ok := cuo.mutation.AddedBatteryNum(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldBatteryNum,
		})
	}
	if value, ok := cuo.mutation.BatteryFullNum(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldBatteryFullNum,
		})
	}
	if value, ok := cuo.mutation.AddedBatteryFullNum(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint,
			Value:  value,
			Column: cabinet.FieldBatteryFullNum,
		})
	}
	if cuo.mutation.BranchCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   cabinet.BranchTable,
			Columns: []string{cabinet.BranchColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branch.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.BranchIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   cabinet.BranchTable,
			Columns: []string{cabinet.BranchColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: branch.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cuo.mutation.BmsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cabinet.BmsTable,
			Columns: cabinet.BmsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: batterymodel.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedBmsIDs(); len(nodes) > 0 && !cuo.mutation.BmsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cabinet.BmsTable,
			Columns: cabinet.BmsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: batterymodel.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.BmsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cabinet.BmsTable,
			Columns: cabinet.BmsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: batterymodel.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Cabinet{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cabinet.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}