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
)

// AssetMaintenanceDetailsCreate is the builder for creating a AssetMaintenanceDetails entity.
type AssetMaintenanceDetailsCreate struct {
	config
	mutation *AssetMaintenanceDetailsMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (amdc *AssetMaintenanceDetailsCreate) SetCreatedAt(t time.Time) *AssetMaintenanceDetailsCreate {
	amdc.mutation.SetCreatedAt(t)
	return amdc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (amdc *AssetMaintenanceDetailsCreate) SetNillableCreatedAt(t *time.Time) *AssetMaintenanceDetailsCreate {
	if t != nil {
		amdc.SetCreatedAt(*t)
	}
	return amdc
}

// SetUpdatedAt sets the "updated_at" field.
func (amdc *AssetMaintenanceDetailsCreate) SetUpdatedAt(t time.Time) *AssetMaintenanceDetailsCreate {
	amdc.mutation.SetUpdatedAt(t)
	return amdc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (amdc *AssetMaintenanceDetailsCreate) SetNillableUpdatedAt(t *time.Time) *AssetMaintenanceDetailsCreate {
	if t != nil {
		amdc.SetUpdatedAt(*t)
	}
	return amdc
}

// SetDeletedAt sets the "deleted_at" field.
func (amdc *AssetMaintenanceDetailsCreate) SetDeletedAt(t time.Time) *AssetMaintenanceDetailsCreate {
	amdc.mutation.SetDeletedAt(t)
	return amdc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (amdc *AssetMaintenanceDetailsCreate) SetNillableDeletedAt(t *time.Time) *AssetMaintenanceDetailsCreate {
	if t != nil {
		amdc.SetDeletedAt(*t)
	}
	return amdc
}

// SetCreator sets the "creator" field.
func (amdc *AssetMaintenanceDetailsCreate) SetCreator(m *model.Modifier) *AssetMaintenanceDetailsCreate {
	amdc.mutation.SetCreator(m)
	return amdc
}

// SetLastModifier sets the "last_modifier" field.
func (amdc *AssetMaintenanceDetailsCreate) SetLastModifier(m *model.Modifier) *AssetMaintenanceDetailsCreate {
	amdc.mutation.SetLastModifier(m)
	return amdc
}

// SetRemark sets the "remark" field.
func (amdc *AssetMaintenanceDetailsCreate) SetRemark(s string) *AssetMaintenanceDetailsCreate {
	amdc.mutation.SetRemark(s)
	return amdc
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (amdc *AssetMaintenanceDetailsCreate) SetNillableRemark(s *string) *AssetMaintenanceDetailsCreate {
	if s != nil {
		amdc.SetRemark(*s)
	}
	return amdc
}

// SetMaterialID sets the "material_id" field.
func (amdc *AssetMaintenanceDetailsCreate) SetMaterialID(u uint64) *AssetMaintenanceDetailsCreate {
	amdc.mutation.SetMaterialID(u)
	return amdc
}

// SetNillableMaterialID sets the "material_id" field if the given value is not nil.
func (amdc *AssetMaintenanceDetailsCreate) SetNillableMaterialID(u *uint64) *AssetMaintenanceDetailsCreate {
	if u != nil {
		amdc.SetMaterialID(*u)
	}
	return amdc
}

// SetAssetID sets the "asset_id" field.
func (amdc *AssetMaintenanceDetailsCreate) SetAssetID(u uint64) *AssetMaintenanceDetailsCreate {
	amdc.mutation.SetAssetID(u)
	return amdc
}

// SetNillableAssetID sets the "asset_id" field if the given value is not nil.
func (amdc *AssetMaintenanceDetailsCreate) SetNillableAssetID(u *uint64) *AssetMaintenanceDetailsCreate {
	if u != nil {
		amdc.SetAssetID(*u)
	}
	return amdc
}

// SetMaintenanceID sets the "maintenance_id" field.
func (amdc *AssetMaintenanceDetailsCreate) SetMaintenanceID(u uint64) *AssetMaintenanceDetailsCreate {
	amdc.mutation.SetMaintenanceID(u)
	return amdc
}

// SetNillableMaintenanceID sets the "maintenance_id" field if the given value is not nil.
func (amdc *AssetMaintenanceDetailsCreate) SetNillableMaintenanceID(u *uint64) *AssetMaintenanceDetailsCreate {
	if u != nil {
		amdc.SetMaintenanceID(*u)
	}
	return amdc
}

// SetMaterial sets the "material" edge to the Material entity.
func (amdc *AssetMaintenanceDetailsCreate) SetMaterial(m *Material) *AssetMaintenanceDetailsCreate {
	return amdc.SetMaterialID(m.ID)
}

// SetAsset sets the "asset" edge to the Asset entity.
func (amdc *AssetMaintenanceDetailsCreate) SetAsset(a *Asset) *AssetMaintenanceDetailsCreate {
	return amdc.SetAssetID(a.ID)
}

// SetMaintenance sets the "maintenance" edge to the AssetMaintenance entity.
func (amdc *AssetMaintenanceDetailsCreate) SetMaintenance(a *AssetMaintenance) *AssetMaintenanceDetailsCreate {
	return amdc.SetMaintenanceID(a.ID)
}

// Mutation returns the AssetMaintenanceDetailsMutation object of the builder.
func (amdc *AssetMaintenanceDetailsCreate) Mutation() *AssetMaintenanceDetailsMutation {
	return amdc.mutation
}

// Save creates the AssetMaintenanceDetails in the database.
func (amdc *AssetMaintenanceDetailsCreate) Save(ctx context.Context) (*AssetMaintenanceDetails, error) {
	if err := amdc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, amdc.sqlSave, amdc.mutation, amdc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (amdc *AssetMaintenanceDetailsCreate) SaveX(ctx context.Context) *AssetMaintenanceDetails {
	v, err := amdc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (amdc *AssetMaintenanceDetailsCreate) Exec(ctx context.Context) error {
	_, err := amdc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amdc *AssetMaintenanceDetailsCreate) ExecX(ctx context.Context) {
	if err := amdc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (amdc *AssetMaintenanceDetailsCreate) defaults() error {
	if _, ok := amdc.mutation.CreatedAt(); !ok {
		if assetmaintenancedetails.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized assetmaintenancedetails.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := assetmaintenancedetails.DefaultCreatedAt()
		amdc.mutation.SetCreatedAt(v)
	}
	if _, ok := amdc.mutation.UpdatedAt(); !ok {
		if assetmaintenancedetails.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized assetmaintenancedetails.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := assetmaintenancedetails.DefaultUpdatedAt()
		amdc.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (amdc *AssetMaintenanceDetailsCreate) check() error {
	if _, ok := amdc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "AssetMaintenanceDetails.created_at"`)}
	}
	if _, ok := amdc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "AssetMaintenanceDetails.updated_at"`)}
	}
	return nil
}

func (amdc *AssetMaintenanceDetailsCreate) sqlSave(ctx context.Context) (*AssetMaintenanceDetails, error) {
	if err := amdc.check(); err != nil {
		return nil, err
	}
	_node, _spec := amdc.createSpec()
	if err := sqlgraph.CreateNode(ctx, amdc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	amdc.mutation.id = &_node.ID
	amdc.mutation.done = true
	return _node, nil
}

func (amdc *AssetMaintenanceDetailsCreate) createSpec() (*AssetMaintenanceDetails, *sqlgraph.CreateSpec) {
	var (
		_node = &AssetMaintenanceDetails{config: amdc.config}
		_spec = sqlgraph.NewCreateSpec(assetmaintenancedetails.Table, sqlgraph.NewFieldSpec(assetmaintenancedetails.FieldID, field.TypeUint64))
	)
	_spec.OnConflict = amdc.conflict
	if value, ok := amdc.mutation.CreatedAt(); ok {
		_spec.SetField(assetmaintenancedetails.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := amdc.mutation.UpdatedAt(); ok {
		_spec.SetField(assetmaintenancedetails.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := amdc.mutation.DeletedAt(); ok {
		_spec.SetField(assetmaintenancedetails.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	if value, ok := amdc.mutation.Creator(); ok {
		_spec.SetField(assetmaintenancedetails.FieldCreator, field.TypeJSON, value)
		_node.Creator = value
	}
	if value, ok := amdc.mutation.LastModifier(); ok {
		_spec.SetField(assetmaintenancedetails.FieldLastModifier, field.TypeJSON, value)
		_node.LastModifier = value
	}
	if value, ok := amdc.mutation.Remark(); ok {
		_spec.SetField(assetmaintenancedetails.FieldRemark, field.TypeString, value)
		_node.Remark = value
	}
	if nodes := amdc.mutation.MaterialIDs(); len(nodes) > 0 {
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
		_node.MaterialID = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := amdc.mutation.AssetIDs(); len(nodes) > 0 {
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
		_node.AssetID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := amdc.mutation.MaintenanceIDs(); len(nodes) > 0 {
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
		_node.MaintenanceID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.AssetMaintenanceDetails.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AssetMaintenanceDetailsUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (amdc *AssetMaintenanceDetailsCreate) OnConflict(opts ...sql.ConflictOption) *AssetMaintenanceDetailsUpsertOne {
	amdc.conflict = opts
	return &AssetMaintenanceDetailsUpsertOne{
		create: amdc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.AssetMaintenanceDetails.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (amdc *AssetMaintenanceDetailsCreate) OnConflictColumns(columns ...string) *AssetMaintenanceDetailsUpsertOne {
	amdc.conflict = append(amdc.conflict, sql.ConflictColumns(columns...))
	return &AssetMaintenanceDetailsUpsertOne{
		create: amdc,
	}
}

type (
	// AssetMaintenanceDetailsUpsertOne is the builder for "upsert"-ing
	//  one AssetMaintenanceDetails node.
	AssetMaintenanceDetailsUpsertOne struct {
		create *AssetMaintenanceDetailsCreate
	}

	// AssetMaintenanceDetailsUpsert is the "OnConflict" setter.
	AssetMaintenanceDetailsUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *AssetMaintenanceDetailsUpsert) SetUpdatedAt(v time.Time) *AssetMaintenanceDetailsUpsert {
	u.Set(assetmaintenancedetails.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsert) UpdateUpdatedAt() *AssetMaintenanceDetailsUpsert {
	u.SetExcluded(assetmaintenancedetails.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *AssetMaintenanceDetailsUpsert) SetDeletedAt(v time.Time) *AssetMaintenanceDetailsUpsert {
	u.Set(assetmaintenancedetails.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsert) UpdateDeletedAt() *AssetMaintenanceDetailsUpsert {
	u.SetExcluded(assetmaintenancedetails.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *AssetMaintenanceDetailsUpsert) ClearDeletedAt() *AssetMaintenanceDetailsUpsert {
	u.SetNull(assetmaintenancedetails.FieldDeletedAt)
	return u
}

// SetLastModifier sets the "last_modifier" field.
func (u *AssetMaintenanceDetailsUpsert) SetLastModifier(v *model.Modifier) *AssetMaintenanceDetailsUpsert {
	u.Set(assetmaintenancedetails.FieldLastModifier, v)
	return u
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsert) UpdateLastModifier() *AssetMaintenanceDetailsUpsert {
	u.SetExcluded(assetmaintenancedetails.FieldLastModifier)
	return u
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *AssetMaintenanceDetailsUpsert) ClearLastModifier() *AssetMaintenanceDetailsUpsert {
	u.SetNull(assetmaintenancedetails.FieldLastModifier)
	return u
}

// SetRemark sets the "remark" field.
func (u *AssetMaintenanceDetailsUpsert) SetRemark(v string) *AssetMaintenanceDetailsUpsert {
	u.Set(assetmaintenancedetails.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsert) UpdateRemark() *AssetMaintenanceDetailsUpsert {
	u.SetExcluded(assetmaintenancedetails.FieldRemark)
	return u
}

// ClearRemark clears the value of the "remark" field.
func (u *AssetMaintenanceDetailsUpsert) ClearRemark() *AssetMaintenanceDetailsUpsert {
	u.SetNull(assetmaintenancedetails.FieldRemark)
	return u
}

// SetMaterialID sets the "material_id" field.
func (u *AssetMaintenanceDetailsUpsert) SetMaterialID(v uint64) *AssetMaintenanceDetailsUpsert {
	u.Set(assetmaintenancedetails.FieldMaterialID, v)
	return u
}

// UpdateMaterialID sets the "material_id" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsert) UpdateMaterialID() *AssetMaintenanceDetailsUpsert {
	u.SetExcluded(assetmaintenancedetails.FieldMaterialID)
	return u
}

// ClearMaterialID clears the value of the "material_id" field.
func (u *AssetMaintenanceDetailsUpsert) ClearMaterialID() *AssetMaintenanceDetailsUpsert {
	u.SetNull(assetmaintenancedetails.FieldMaterialID)
	return u
}

// SetAssetID sets the "asset_id" field.
func (u *AssetMaintenanceDetailsUpsert) SetAssetID(v uint64) *AssetMaintenanceDetailsUpsert {
	u.Set(assetmaintenancedetails.FieldAssetID, v)
	return u
}

// UpdateAssetID sets the "asset_id" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsert) UpdateAssetID() *AssetMaintenanceDetailsUpsert {
	u.SetExcluded(assetmaintenancedetails.FieldAssetID)
	return u
}

// ClearAssetID clears the value of the "asset_id" field.
func (u *AssetMaintenanceDetailsUpsert) ClearAssetID() *AssetMaintenanceDetailsUpsert {
	u.SetNull(assetmaintenancedetails.FieldAssetID)
	return u
}

// SetMaintenanceID sets the "maintenance_id" field.
func (u *AssetMaintenanceDetailsUpsert) SetMaintenanceID(v uint64) *AssetMaintenanceDetailsUpsert {
	u.Set(assetmaintenancedetails.FieldMaintenanceID, v)
	return u
}

// UpdateMaintenanceID sets the "maintenance_id" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsert) UpdateMaintenanceID() *AssetMaintenanceDetailsUpsert {
	u.SetExcluded(assetmaintenancedetails.FieldMaintenanceID)
	return u
}

// ClearMaintenanceID clears the value of the "maintenance_id" field.
func (u *AssetMaintenanceDetailsUpsert) ClearMaintenanceID() *AssetMaintenanceDetailsUpsert {
	u.SetNull(assetmaintenancedetails.FieldMaintenanceID)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.AssetMaintenanceDetails.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *AssetMaintenanceDetailsUpsertOne) UpdateNewValues() *AssetMaintenanceDetailsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(assetmaintenancedetails.FieldCreatedAt)
		}
		if _, exists := u.create.mutation.Creator(); exists {
			s.SetIgnore(assetmaintenancedetails.FieldCreator)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.AssetMaintenanceDetails.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *AssetMaintenanceDetailsUpsertOne) Ignore() *AssetMaintenanceDetailsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AssetMaintenanceDetailsUpsertOne) DoNothing() *AssetMaintenanceDetailsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AssetMaintenanceDetailsCreate.OnConflict
// documentation for more info.
func (u *AssetMaintenanceDetailsUpsertOne) Update(set func(*AssetMaintenanceDetailsUpsert)) *AssetMaintenanceDetailsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AssetMaintenanceDetailsUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *AssetMaintenanceDetailsUpsertOne) SetUpdatedAt(v time.Time) *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertOne) UpdateUpdatedAt() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *AssetMaintenanceDetailsUpsertOne) SetDeletedAt(v time.Time) *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertOne) UpdateDeletedAt() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *AssetMaintenanceDetailsUpsertOne) ClearDeletedAt() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *AssetMaintenanceDetailsUpsertOne) SetLastModifier(v *model.Modifier) *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertOne) UpdateLastModifier() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *AssetMaintenanceDetailsUpsertOne) ClearLastModifier() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *AssetMaintenanceDetailsUpsertOne) SetRemark(v string) *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertOne) UpdateRemark() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *AssetMaintenanceDetailsUpsertOne) ClearRemark() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearRemark()
	})
}

// SetMaterialID sets the "material_id" field.
func (u *AssetMaintenanceDetailsUpsertOne) SetMaterialID(v uint64) *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetMaterialID(v)
	})
}

// UpdateMaterialID sets the "material_id" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertOne) UpdateMaterialID() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateMaterialID()
	})
}

// ClearMaterialID clears the value of the "material_id" field.
func (u *AssetMaintenanceDetailsUpsertOne) ClearMaterialID() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearMaterialID()
	})
}

// SetAssetID sets the "asset_id" field.
func (u *AssetMaintenanceDetailsUpsertOne) SetAssetID(v uint64) *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetAssetID(v)
	})
}

// UpdateAssetID sets the "asset_id" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertOne) UpdateAssetID() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateAssetID()
	})
}

// ClearAssetID clears the value of the "asset_id" field.
func (u *AssetMaintenanceDetailsUpsertOne) ClearAssetID() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearAssetID()
	})
}

// SetMaintenanceID sets the "maintenance_id" field.
func (u *AssetMaintenanceDetailsUpsertOne) SetMaintenanceID(v uint64) *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetMaintenanceID(v)
	})
}

// UpdateMaintenanceID sets the "maintenance_id" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertOne) UpdateMaintenanceID() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateMaintenanceID()
	})
}

// ClearMaintenanceID clears the value of the "maintenance_id" field.
func (u *AssetMaintenanceDetailsUpsertOne) ClearMaintenanceID() *AssetMaintenanceDetailsUpsertOne {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearMaintenanceID()
	})
}

// Exec executes the query.
func (u *AssetMaintenanceDetailsUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AssetMaintenanceDetailsCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AssetMaintenanceDetailsUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *AssetMaintenanceDetailsUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *AssetMaintenanceDetailsUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// AssetMaintenanceDetailsCreateBulk is the builder for creating many AssetMaintenanceDetails entities in bulk.
type AssetMaintenanceDetailsCreateBulk struct {
	config
	err      error
	builders []*AssetMaintenanceDetailsCreate
	conflict []sql.ConflictOption
}

// Save creates the AssetMaintenanceDetails entities in the database.
func (amdcb *AssetMaintenanceDetailsCreateBulk) Save(ctx context.Context) ([]*AssetMaintenanceDetails, error) {
	if amdcb.err != nil {
		return nil, amdcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(amdcb.builders))
	nodes := make([]*AssetMaintenanceDetails, len(amdcb.builders))
	mutators := make([]Mutator, len(amdcb.builders))
	for i := range amdcb.builders {
		func(i int, root context.Context) {
			builder := amdcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AssetMaintenanceDetailsMutation)
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
					_, err = mutators[i+1].Mutate(root, amdcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = amdcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, amdcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, amdcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (amdcb *AssetMaintenanceDetailsCreateBulk) SaveX(ctx context.Context) []*AssetMaintenanceDetails {
	v, err := amdcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (amdcb *AssetMaintenanceDetailsCreateBulk) Exec(ctx context.Context) error {
	_, err := amdcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amdcb *AssetMaintenanceDetailsCreateBulk) ExecX(ctx context.Context) {
	if err := amdcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.AssetMaintenanceDetails.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AssetMaintenanceDetailsUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (amdcb *AssetMaintenanceDetailsCreateBulk) OnConflict(opts ...sql.ConflictOption) *AssetMaintenanceDetailsUpsertBulk {
	amdcb.conflict = opts
	return &AssetMaintenanceDetailsUpsertBulk{
		create: amdcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.AssetMaintenanceDetails.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (amdcb *AssetMaintenanceDetailsCreateBulk) OnConflictColumns(columns ...string) *AssetMaintenanceDetailsUpsertBulk {
	amdcb.conflict = append(amdcb.conflict, sql.ConflictColumns(columns...))
	return &AssetMaintenanceDetailsUpsertBulk{
		create: amdcb,
	}
}

// AssetMaintenanceDetailsUpsertBulk is the builder for "upsert"-ing
// a bulk of AssetMaintenanceDetails nodes.
type AssetMaintenanceDetailsUpsertBulk struct {
	create *AssetMaintenanceDetailsCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.AssetMaintenanceDetails.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *AssetMaintenanceDetailsUpsertBulk) UpdateNewValues() *AssetMaintenanceDetailsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(assetmaintenancedetails.FieldCreatedAt)
			}
			if _, exists := b.mutation.Creator(); exists {
				s.SetIgnore(assetmaintenancedetails.FieldCreator)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.AssetMaintenanceDetails.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *AssetMaintenanceDetailsUpsertBulk) Ignore() *AssetMaintenanceDetailsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AssetMaintenanceDetailsUpsertBulk) DoNothing() *AssetMaintenanceDetailsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AssetMaintenanceDetailsCreateBulk.OnConflict
// documentation for more info.
func (u *AssetMaintenanceDetailsUpsertBulk) Update(set func(*AssetMaintenanceDetailsUpsert)) *AssetMaintenanceDetailsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AssetMaintenanceDetailsUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *AssetMaintenanceDetailsUpsertBulk) SetUpdatedAt(v time.Time) *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertBulk) UpdateUpdatedAt() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *AssetMaintenanceDetailsUpsertBulk) SetDeletedAt(v time.Time) *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertBulk) UpdateDeletedAt() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *AssetMaintenanceDetailsUpsertBulk) ClearDeletedAt() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *AssetMaintenanceDetailsUpsertBulk) SetLastModifier(v *model.Modifier) *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertBulk) UpdateLastModifier() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *AssetMaintenanceDetailsUpsertBulk) ClearLastModifier() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *AssetMaintenanceDetailsUpsertBulk) SetRemark(v string) *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertBulk) UpdateRemark() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *AssetMaintenanceDetailsUpsertBulk) ClearRemark() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearRemark()
	})
}

// SetMaterialID sets the "material_id" field.
func (u *AssetMaintenanceDetailsUpsertBulk) SetMaterialID(v uint64) *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetMaterialID(v)
	})
}

// UpdateMaterialID sets the "material_id" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertBulk) UpdateMaterialID() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateMaterialID()
	})
}

// ClearMaterialID clears the value of the "material_id" field.
func (u *AssetMaintenanceDetailsUpsertBulk) ClearMaterialID() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearMaterialID()
	})
}

// SetAssetID sets the "asset_id" field.
func (u *AssetMaintenanceDetailsUpsertBulk) SetAssetID(v uint64) *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetAssetID(v)
	})
}

// UpdateAssetID sets the "asset_id" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertBulk) UpdateAssetID() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateAssetID()
	})
}

// ClearAssetID clears the value of the "asset_id" field.
func (u *AssetMaintenanceDetailsUpsertBulk) ClearAssetID() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearAssetID()
	})
}

// SetMaintenanceID sets the "maintenance_id" field.
func (u *AssetMaintenanceDetailsUpsertBulk) SetMaintenanceID(v uint64) *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.SetMaintenanceID(v)
	})
}

// UpdateMaintenanceID sets the "maintenance_id" field to the value that was provided on create.
func (u *AssetMaintenanceDetailsUpsertBulk) UpdateMaintenanceID() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.UpdateMaintenanceID()
	})
}

// ClearMaintenanceID clears the value of the "maintenance_id" field.
func (u *AssetMaintenanceDetailsUpsertBulk) ClearMaintenanceID() *AssetMaintenanceDetailsUpsertBulk {
	return u.Update(func(s *AssetMaintenanceDetailsUpsert) {
		s.ClearMaintenanceID()
	})
}

// Exec executes the query.
func (u *AssetMaintenanceDetailsUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the AssetMaintenanceDetailsCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AssetMaintenanceDetailsCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AssetMaintenanceDetailsUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}