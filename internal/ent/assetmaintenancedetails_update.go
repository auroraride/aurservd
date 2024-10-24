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
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetmaintenance"
	"github.com/auroraride/aurservd/internal/ent/assetmaintenancedetails"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// AssetMaintenanceDetailsUpdate is the builder for updating AssetMaintenanceDetails entities.
type AssetMaintenanceDetailsUpdate struct {
	config
	hooks     []Hook
	mutation  *AssetMaintenanceDetailsMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the AssetMaintenanceDetailsUpdate builder.
func (amdu *AssetMaintenanceDetailsUpdate) Where(ps ...predicate.AssetMaintenanceDetails) *AssetMaintenanceDetailsUpdate {
	amdu.mutation.Where(ps...)
	return amdu
}

// SetUpdatedAt sets the "updated_at" field.
func (amdu *AssetMaintenanceDetailsUpdate) SetUpdatedAt(t time.Time) *AssetMaintenanceDetailsUpdate {
	amdu.mutation.SetUpdatedAt(t)
	return amdu
}

// SetDeletedAt sets the "deleted_at" field.
func (amdu *AssetMaintenanceDetailsUpdate) SetDeletedAt(t time.Time) *AssetMaintenanceDetailsUpdate {
	amdu.mutation.SetDeletedAt(t)
	return amdu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (amdu *AssetMaintenanceDetailsUpdate) SetNillableDeletedAt(t *time.Time) *AssetMaintenanceDetailsUpdate {
	if t != nil {
		amdu.SetDeletedAt(*t)
	}
	return amdu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (amdu *AssetMaintenanceDetailsUpdate) ClearDeletedAt() *AssetMaintenanceDetailsUpdate {
	amdu.mutation.ClearDeletedAt()
	return amdu
}

// SetLastModifier sets the "last_modifier" field.
func (amdu *AssetMaintenanceDetailsUpdate) SetLastModifier(m *model.Modifier) *AssetMaintenanceDetailsUpdate {
	amdu.mutation.SetLastModifier(m)
	return amdu
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (amdu *AssetMaintenanceDetailsUpdate) ClearLastModifier() *AssetMaintenanceDetailsUpdate {
	amdu.mutation.ClearLastModifier()
	return amdu
}

// SetRemark sets the "remark" field.
func (amdu *AssetMaintenanceDetailsUpdate) SetRemark(s string) *AssetMaintenanceDetailsUpdate {
	amdu.mutation.SetRemark(s)
	return amdu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (amdu *AssetMaintenanceDetailsUpdate) SetNillableRemark(s *string) *AssetMaintenanceDetailsUpdate {
	if s != nil {
		amdu.SetRemark(*s)
	}
	return amdu
}

// ClearRemark clears the value of the "remark" field.
func (amdu *AssetMaintenanceDetailsUpdate) ClearRemark() *AssetMaintenanceDetailsUpdate {
	amdu.mutation.ClearRemark()
	return amdu
}

// SetMaterialID sets the "material_id" field.
func (amdu *AssetMaintenanceDetailsUpdate) SetMaterialID(u uint64) *AssetMaintenanceDetailsUpdate {
	amdu.mutation.SetMaterialID(u)
	return amdu
}

// SetNillableMaterialID sets the "material_id" field if the given value is not nil.
func (amdu *AssetMaintenanceDetailsUpdate) SetNillableMaterialID(u *uint64) *AssetMaintenanceDetailsUpdate {
	if u != nil {
		amdu.SetMaterialID(*u)
	}
	return amdu
}

// ClearMaterialID clears the value of the "material_id" field.
func (amdu *AssetMaintenanceDetailsUpdate) ClearMaterialID() *AssetMaintenanceDetailsUpdate {
	amdu.mutation.ClearMaterialID()
	return amdu
}

// SetSn sets the "sn" field.
func (amdu *AssetMaintenanceDetailsUpdate) SetSn(s string) *AssetMaintenanceDetailsUpdate {
	amdu.mutation.SetSn(s)
	return amdu
}

// SetNillableSn sets the "sn" field if the given value is not nil.
func (amdu *AssetMaintenanceDetailsUpdate) SetNillableSn(s *string) *AssetMaintenanceDetailsUpdate {
	if s != nil {
		amdu.SetSn(*s)
	}
	return amdu
}

// ClearSn clears the value of the "sn" field.
func (amdu *AssetMaintenanceDetailsUpdate) ClearSn() *AssetMaintenanceDetailsUpdate {
	amdu.mutation.ClearSn()
	return amdu
}

// SetAssetID sets the "asset_id" field.
func (amdu *AssetMaintenanceDetailsUpdate) SetAssetID(u uint64) *AssetMaintenanceDetailsUpdate {
	amdu.mutation.SetAssetID(u)
	return amdu
}

// SetNillableAssetID sets the "asset_id" field if the given value is not nil.
func (amdu *AssetMaintenanceDetailsUpdate) SetNillableAssetID(u *uint64) *AssetMaintenanceDetailsUpdate {
	if u != nil {
		amdu.SetAssetID(*u)
	}
	return amdu
}

// ClearAssetID clears the value of the "asset_id" field.
func (amdu *AssetMaintenanceDetailsUpdate) ClearAssetID() *AssetMaintenanceDetailsUpdate {
	amdu.mutation.ClearAssetID()
	return amdu
}

// SetMaintenanceID sets the "maintenance_id" field.
func (amdu *AssetMaintenanceDetailsUpdate) SetMaintenanceID(u uint64) *AssetMaintenanceDetailsUpdate {
	amdu.mutation.SetMaintenanceID(u)
	return amdu
}

// SetNillableMaintenanceID sets the "maintenance_id" field if the given value is not nil.
func (amdu *AssetMaintenanceDetailsUpdate) SetNillableMaintenanceID(u *uint64) *AssetMaintenanceDetailsUpdate {
	if u != nil {
		amdu.SetMaintenanceID(*u)
	}
	return amdu
}

// ClearMaintenanceID clears the value of the "maintenance_id" field.
func (amdu *AssetMaintenanceDetailsUpdate) ClearMaintenanceID() *AssetMaintenanceDetailsUpdate {
	amdu.mutation.ClearMaintenanceID()
	return amdu
}

// SetMaterial sets the "material" edge to the Material entity.
func (amdu *AssetMaintenanceDetailsUpdate) SetMaterial(m *Material) *AssetMaintenanceDetailsUpdate {
	return amdu.SetMaterialID(m.ID)
}

// SetAsset sets the "asset" edge to the Asset entity.
func (amdu *AssetMaintenanceDetailsUpdate) SetAsset(a *Asset) *AssetMaintenanceDetailsUpdate {
	return amdu.SetAssetID(a.ID)
}

// SetMaintenance sets the "maintenance" edge to the AssetMaintenance entity.
func (amdu *AssetMaintenanceDetailsUpdate) SetMaintenance(a *AssetMaintenance) *AssetMaintenanceDetailsUpdate {
	return amdu.SetMaintenanceID(a.ID)
}

// Mutation returns the AssetMaintenanceDetailsMutation object of the builder.
func (amdu *AssetMaintenanceDetailsUpdate) Mutation() *AssetMaintenanceDetailsMutation {
	return amdu.mutation
}

// ClearMaterial clears the "material" edge to the Material entity.
func (amdu *AssetMaintenanceDetailsUpdate) ClearMaterial() *AssetMaintenanceDetailsUpdate {
	amdu.mutation.ClearMaterial()
	return amdu
}

// ClearAsset clears the "asset" edge to the Asset entity.
func (amdu *AssetMaintenanceDetailsUpdate) ClearAsset() *AssetMaintenanceDetailsUpdate {
	amdu.mutation.ClearAsset()
	return amdu
}

// ClearMaintenance clears the "maintenance" edge to the AssetMaintenance entity.
func (amdu *AssetMaintenanceDetailsUpdate) ClearMaintenance() *AssetMaintenanceDetailsUpdate {
	amdu.mutation.ClearMaintenance()
	return amdu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (amdu *AssetMaintenanceDetailsUpdate) Save(ctx context.Context) (int, error) {
	if err := amdu.defaults(); err != nil {
		return 0, err
	}
	return withHooks(ctx, amdu.sqlSave, amdu.mutation, amdu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (amdu *AssetMaintenanceDetailsUpdate) SaveX(ctx context.Context) int {
	affected, err := amdu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (amdu *AssetMaintenanceDetailsUpdate) Exec(ctx context.Context) error {
	_, err := amdu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amdu *AssetMaintenanceDetailsUpdate) ExecX(ctx context.Context) {
	if err := amdu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (amdu *AssetMaintenanceDetailsUpdate) defaults() error {
	if _, ok := amdu.mutation.UpdatedAt(); !ok {
		if assetmaintenancedetails.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized assetmaintenancedetails.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := assetmaintenancedetails.UpdateDefaultUpdatedAt()
		amdu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (amdu *AssetMaintenanceDetailsUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *AssetMaintenanceDetailsUpdate {
	amdu.modifiers = append(amdu.modifiers, modifiers...)
	return amdu
}

func (amdu *AssetMaintenanceDetailsUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(assetmaintenancedetails.Table, assetmaintenancedetails.Columns, sqlgraph.NewFieldSpec(assetmaintenancedetails.FieldID, field.TypeUint64))
	if ps := amdu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := amdu.mutation.UpdatedAt(); ok {
		_spec.SetField(assetmaintenancedetails.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := amdu.mutation.DeletedAt(); ok {
		_spec.SetField(assetmaintenancedetails.FieldDeletedAt, field.TypeTime, value)
	}
	if amdu.mutation.DeletedAtCleared() {
		_spec.ClearField(assetmaintenancedetails.FieldDeletedAt, field.TypeTime)
	}
	if amdu.mutation.CreatorCleared() {
		_spec.ClearField(assetmaintenancedetails.FieldCreator, field.TypeJSON)
	}
	if value, ok := amdu.mutation.LastModifier(); ok {
		_spec.SetField(assetmaintenancedetails.FieldLastModifier, field.TypeJSON, value)
	}
	if amdu.mutation.LastModifierCleared() {
		_spec.ClearField(assetmaintenancedetails.FieldLastModifier, field.TypeJSON)
	}
	if value, ok := amdu.mutation.Remark(); ok {
		_spec.SetField(assetmaintenancedetails.FieldRemark, field.TypeString, value)
	}
	if amdu.mutation.RemarkCleared() {
		_spec.ClearField(assetmaintenancedetails.FieldRemark, field.TypeString)
	}
	if value, ok := amdu.mutation.Sn(); ok {
		_spec.SetField(assetmaintenancedetails.FieldSn, field.TypeString, value)
	}
	if amdu.mutation.SnCleared() {
		_spec.ClearField(assetmaintenancedetails.FieldSn, field.TypeString)
	}
	if amdu.mutation.MaterialCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   assetmaintenancedetails.MaterialTable,
			Columns: []string{assetmaintenancedetails.MaterialColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(material.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := amdu.mutation.MaterialIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   assetmaintenancedetails.MaterialTable,
			Columns: []string{assetmaintenancedetails.MaterialColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(material.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if amdu.mutation.AssetCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   assetmaintenancedetails.AssetTable,
			Columns: []string{assetmaintenancedetails.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := amdu.mutation.AssetIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   assetmaintenancedetails.AssetTable,
			Columns: []string{assetmaintenancedetails.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if amdu.mutation.MaintenanceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   assetmaintenancedetails.MaintenanceTable,
			Columns: []string{assetmaintenancedetails.MaintenanceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(assetmaintenance.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := amdu.mutation.MaintenanceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   assetmaintenancedetails.MaintenanceTable,
			Columns: []string{assetmaintenancedetails.MaintenanceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(assetmaintenance.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(amdu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, amdu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{assetmaintenancedetails.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	amdu.mutation.done = true
	return n, nil
}

// AssetMaintenanceDetailsUpdateOne is the builder for updating a single AssetMaintenanceDetails entity.
type AssetMaintenanceDetailsUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *AssetMaintenanceDetailsMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetUpdatedAt(t time.Time) *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.SetUpdatedAt(t)
	return amduo
}

// SetDeletedAt sets the "deleted_at" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetDeletedAt(t time.Time) *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.SetDeletedAt(t)
	return amduo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetNillableDeletedAt(t *time.Time) *AssetMaintenanceDetailsUpdateOne {
	if t != nil {
		amduo.SetDeletedAt(*t)
	}
	return amduo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) ClearDeletedAt() *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.ClearDeletedAt()
	return amduo
}

// SetLastModifier sets the "last_modifier" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetLastModifier(m *model.Modifier) *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.SetLastModifier(m)
	return amduo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) ClearLastModifier() *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.ClearLastModifier()
	return amduo
}

// SetRemark sets the "remark" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetRemark(s string) *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.SetRemark(s)
	return amduo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetNillableRemark(s *string) *AssetMaintenanceDetailsUpdateOne {
	if s != nil {
		amduo.SetRemark(*s)
	}
	return amduo
}

// ClearRemark clears the value of the "remark" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) ClearRemark() *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.ClearRemark()
	return amduo
}

// SetMaterialID sets the "material_id" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetMaterialID(u uint64) *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.SetMaterialID(u)
	return amduo
}

// SetNillableMaterialID sets the "material_id" field if the given value is not nil.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetNillableMaterialID(u *uint64) *AssetMaintenanceDetailsUpdateOne {
	if u != nil {
		amduo.SetMaterialID(*u)
	}
	return amduo
}

// ClearMaterialID clears the value of the "material_id" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) ClearMaterialID() *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.ClearMaterialID()
	return amduo
}

// SetSn sets the "sn" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetSn(s string) *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.SetSn(s)
	return amduo
}

// SetNillableSn sets the "sn" field if the given value is not nil.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetNillableSn(s *string) *AssetMaintenanceDetailsUpdateOne {
	if s != nil {
		amduo.SetSn(*s)
	}
	return amduo
}

// ClearSn clears the value of the "sn" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) ClearSn() *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.ClearSn()
	return amduo
}

// SetAssetID sets the "asset_id" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetAssetID(u uint64) *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.SetAssetID(u)
	return amduo
}

// SetNillableAssetID sets the "asset_id" field if the given value is not nil.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetNillableAssetID(u *uint64) *AssetMaintenanceDetailsUpdateOne {
	if u != nil {
		amduo.SetAssetID(*u)
	}
	return amduo
}

// ClearAssetID clears the value of the "asset_id" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) ClearAssetID() *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.ClearAssetID()
	return amduo
}

// SetMaintenanceID sets the "maintenance_id" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetMaintenanceID(u uint64) *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.SetMaintenanceID(u)
	return amduo
}

// SetNillableMaintenanceID sets the "maintenance_id" field if the given value is not nil.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetNillableMaintenanceID(u *uint64) *AssetMaintenanceDetailsUpdateOne {
	if u != nil {
		amduo.SetMaintenanceID(*u)
	}
	return amduo
}

// ClearMaintenanceID clears the value of the "maintenance_id" field.
func (amduo *AssetMaintenanceDetailsUpdateOne) ClearMaintenanceID() *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.ClearMaintenanceID()
	return amduo
}

// SetMaterial sets the "material" edge to the Material entity.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetMaterial(m *Material) *AssetMaintenanceDetailsUpdateOne {
	return amduo.SetMaterialID(m.ID)
}

// SetAsset sets the "asset" edge to the Asset entity.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetAsset(a *Asset) *AssetMaintenanceDetailsUpdateOne {
	return amduo.SetAssetID(a.ID)
}

// SetMaintenance sets the "maintenance" edge to the AssetMaintenance entity.
func (amduo *AssetMaintenanceDetailsUpdateOne) SetMaintenance(a *AssetMaintenance) *AssetMaintenanceDetailsUpdateOne {
	return amduo.SetMaintenanceID(a.ID)
}

// Mutation returns the AssetMaintenanceDetailsMutation object of the builder.
func (amduo *AssetMaintenanceDetailsUpdateOne) Mutation() *AssetMaintenanceDetailsMutation {
	return amduo.mutation
}

// ClearMaterial clears the "material" edge to the Material entity.
func (amduo *AssetMaintenanceDetailsUpdateOne) ClearMaterial() *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.ClearMaterial()
	return amduo
}

// ClearAsset clears the "asset" edge to the Asset entity.
func (amduo *AssetMaintenanceDetailsUpdateOne) ClearAsset() *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.ClearAsset()
	return amduo
}

// ClearMaintenance clears the "maintenance" edge to the AssetMaintenance entity.
func (amduo *AssetMaintenanceDetailsUpdateOne) ClearMaintenance() *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.ClearMaintenance()
	return amduo
}

// Where appends a list predicates to the AssetMaintenanceDetailsUpdate builder.
func (amduo *AssetMaintenanceDetailsUpdateOne) Where(ps ...predicate.AssetMaintenanceDetails) *AssetMaintenanceDetailsUpdateOne {
	amduo.mutation.Where(ps...)
	return amduo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (amduo *AssetMaintenanceDetailsUpdateOne) Select(field string, fields ...string) *AssetMaintenanceDetailsUpdateOne {
	amduo.fields = append([]string{field}, fields...)
	return amduo
}

// Save executes the query and returns the updated AssetMaintenanceDetails entity.
func (amduo *AssetMaintenanceDetailsUpdateOne) Save(ctx context.Context) (*AssetMaintenanceDetails, error) {
	if err := amduo.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, amduo.sqlSave, amduo.mutation, amduo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (amduo *AssetMaintenanceDetailsUpdateOne) SaveX(ctx context.Context) *AssetMaintenanceDetails {
	node, err := amduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (amduo *AssetMaintenanceDetailsUpdateOne) Exec(ctx context.Context) error {
	_, err := amduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amduo *AssetMaintenanceDetailsUpdateOne) ExecX(ctx context.Context) {
	if err := amduo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (amduo *AssetMaintenanceDetailsUpdateOne) defaults() error {
	if _, ok := amduo.mutation.UpdatedAt(); !ok {
		if assetmaintenancedetails.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized assetmaintenancedetails.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := assetmaintenancedetails.UpdateDefaultUpdatedAt()
		amduo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (amduo *AssetMaintenanceDetailsUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *AssetMaintenanceDetailsUpdateOne {
	amduo.modifiers = append(amduo.modifiers, modifiers...)
	return amduo
}

func (amduo *AssetMaintenanceDetailsUpdateOne) sqlSave(ctx context.Context) (_node *AssetMaintenanceDetails, err error) {
	_spec := sqlgraph.NewUpdateSpec(assetmaintenancedetails.Table, assetmaintenancedetails.Columns, sqlgraph.NewFieldSpec(assetmaintenancedetails.FieldID, field.TypeUint64))
	id, ok := amduo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "AssetMaintenanceDetails.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := amduo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, assetmaintenancedetails.FieldID)
		for _, f := range fields {
			if !assetmaintenancedetails.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != assetmaintenancedetails.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := amduo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := amduo.mutation.UpdatedAt(); ok {
		_spec.SetField(assetmaintenancedetails.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := amduo.mutation.DeletedAt(); ok {
		_spec.SetField(assetmaintenancedetails.FieldDeletedAt, field.TypeTime, value)
	}
	if amduo.mutation.DeletedAtCleared() {
		_spec.ClearField(assetmaintenancedetails.FieldDeletedAt, field.TypeTime)
	}
	if amduo.mutation.CreatorCleared() {
		_spec.ClearField(assetmaintenancedetails.FieldCreator, field.TypeJSON)
	}
	if value, ok := amduo.mutation.LastModifier(); ok {
		_spec.SetField(assetmaintenancedetails.FieldLastModifier, field.TypeJSON, value)
	}
	if amduo.mutation.LastModifierCleared() {
		_spec.ClearField(assetmaintenancedetails.FieldLastModifier, field.TypeJSON)
	}
	if value, ok := amduo.mutation.Remark(); ok {
		_spec.SetField(assetmaintenancedetails.FieldRemark, field.TypeString, value)
	}
	if amduo.mutation.RemarkCleared() {
		_spec.ClearField(assetmaintenancedetails.FieldRemark, field.TypeString)
	}
	if value, ok := amduo.mutation.Sn(); ok {
		_spec.SetField(assetmaintenancedetails.FieldSn, field.TypeString, value)
	}
	if amduo.mutation.SnCleared() {
		_spec.ClearField(assetmaintenancedetails.FieldSn, field.TypeString)
	}
	if amduo.mutation.MaterialCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   assetmaintenancedetails.MaterialTable,
			Columns: []string{assetmaintenancedetails.MaterialColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(material.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := amduo.mutation.MaterialIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   assetmaintenancedetails.MaterialTable,
			Columns: []string{assetmaintenancedetails.MaterialColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(material.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if amduo.mutation.AssetCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   assetmaintenancedetails.AssetTable,
			Columns: []string{assetmaintenancedetails.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := amduo.mutation.AssetIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   assetmaintenancedetails.AssetTable,
			Columns: []string{assetmaintenancedetails.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if amduo.mutation.MaintenanceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   assetmaintenancedetails.MaintenanceTable,
			Columns: []string{assetmaintenancedetails.MaintenanceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(assetmaintenance.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := amduo.mutation.MaintenanceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   assetmaintenancedetails.MaintenanceTable,
			Columns: []string{assetmaintenancedetails.MaintenanceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(assetmaintenance.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(amduo.modifiers...)
	_node = &AssetMaintenanceDetails{config: amduo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, amduo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{assetmaintenancedetails.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	amduo.mutation.done = true
	return _node, nil
}
