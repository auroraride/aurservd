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
	"github.com/auroraride/aurservd/internal/ent/assetmaintenance"
	"github.com/auroraride/aurservd/internal/ent/assetmaintenancedetails"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
)

// AssetMaintenanceCreate is the builder for creating a AssetMaintenance entity.
type AssetMaintenanceCreate struct {
	config
	mutation *AssetMaintenanceMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (amc *AssetMaintenanceCreate) SetCreatedAt(t time.Time) *AssetMaintenanceCreate {
	amc.mutation.SetCreatedAt(t)
	return amc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (amc *AssetMaintenanceCreate) SetNillableCreatedAt(t *time.Time) *AssetMaintenanceCreate {
	if t != nil {
		amc.SetCreatedAt(*t)
	}
	return amc
}

// SetUpdatedAt sets the "updated_at" field.
func (amc *AssetMaintenanceCreate) SetUpdatedAt(t time.Time) *AssetMaintenanceCreate {
	amc.mutation.SetUpdatedAt(t)
	return amc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (amc *AssetMaintenanceCreate) SetNillableUpdatedAt(t *time.Time) *AssetMaintenanceCreate {
	if t != nil {
		amc.SetUpdatedAt(*t)
	}
	return amc
}

// SetDeletedAt sets the "deleted_at" field.
func (amc *AssetMaintenanceCreate) SetDeletedAt(t time.Time) *AssetMaintenanceCreate {
	amc.mutation.SetDeletedAt(t)
	return amc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (amc *AssetMaintenanceCreate) SetNillableDeletedAt(t *time.Time) *AssetMaintenanceCreate {
	if t != nil {
		amc.SetDeletedAt(*t)
	}
	return amc
}

// SetCreator sets the "creator" field.
func (amc *AssetMaintenanceCreate) SetCreator(m *model.Modifier) *AssetMaintenanceCreate {
	amc.mutation.SetCreator(m)
	return amc
}

// SetLastModifier sets the "last_modifier" field.
func (amc *AssetMaintenanceCreate) SetLastModifier(m *model.Modifier) *AssetMaintenanceCreate {
	amc.mutation.SetLastModifier(m)
	return amc
}

// SetRemark sets the "remark" field.
func (amc *AssetMaintenanceCreate) SetRemark(s string) *AssetMaintenanceCreate {
	amc.mutation.SetRemark(s)
	return amc
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (amc *AssetMaintenanceCreate) SetNillableRemark(s *string) *AssetMaintenanceCreate {
	if s != nil {
		amc.SetRemark(*s)
	}
	return amc
}

// SetCabinetID sets the "cabinet_id" field.
func (amc *AssetMaintenanceCreate) SetCabinetID(u uint64) *AssetMaintenanceCreate {
	amc.mutation.SetCabinetID(u)
	return amc
}

// SetNillableCabinetID sets the "cabinet_id" field if the given value is not nil.
func (amc *AssetMaintenanceCreate) SetNillableCabinetID(u *uint64) *AssetMaintenanceCreate {
	if u != nil {
		amc.SetCabinetID(*u)
	}
	return amc
}

// SetMaintainerID sets the "maintainer_id" field.
func (amc *AssetMaintenanceCreate) SetMaintainerID(u uint64) *AssetMaintenanceCreate {
	amc.mutation.SetMaintainerID(u)
	return amc
}

// SetNillableMaintainerID sets the "maintainer_id" field if the given value is not nil.
func (amc *AssetMaintenanceCreate) SetNillableMaintainerID(u *uint64) *AssetMaintenanceCreate {
	if u != nil {
		amc.SetMaintainerID(*u)
	}
	return amc
}

// SetReason sets the "reason" field.
func (amc *AssetMaintenanceCreate) SetReason(s string) *AssetMaintenanceCreate {
	amc.mutation.SetReason(s)
	return amc
}

// SetNillableReason sets the "reason" field if the given value is not nil.
func (amc *AssetMaintenanceCreate) SetNillableReason(s *string) *AssetMaintenanceCreate {
	if s != nil {
		amc.SetReason(*s)
	}
	return amc
}

// SetContent sets the "content" field.
func (amc *AssetMaintenanceCreate) SetContent(s string) *AssetMaintenanceCreate {
	amc.mutation.SetContent(s)
	return amc
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (amc *AssetMaintenanceCreate) SetNillableContent(s *string) *AssetMaintenanceCreate {
	if s != nil {
		amc.SetContent(*s)
	}
	return amc
}

// SetStatus sets the "status" field.
func (amc *AssetMaintenanceCreate) SetStatus(u uint8) *AssetMaintenanceCreate {
	amc.mutation.SetStatus(u)
	return amc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (amc *AssetMaintenanceCreate) SetNillableStatus(u *uint8) *AssetMaintenanceCreate {
	if u != nil {
		amc.SetStatus(*u)
	}
	return amc
}

// SetCabinet sets the "cabinet" edge to the Cabinet entity.
func (amc *AssetMaintenanceCreate) SetCabinet(c *Cabinet) *AssetMaintenanceCreate {
	return amc.SetCabinetID(c.ID)
}

// SetMaintainer sets the "maintainer" edge to the Maintainer entity.
func (amc *AssetMaintenanceCreate) SetMaintainer(m *Maintainer) *AssetMaintenanceCreate {
	return amc.SetMaintainerID(m.ID)
}

// AddMaintenanceDetailIDs adds the "maintenance_details" edge to the AssetMaintenanceDetails entity by IDs.
func (amc *AssetMaintenanceCreate) AddMaintenanceDetailIDs(ids ...uint64) *AssetMaintenanceCreate {
	amc.mutation.AddMaintenanceDetailIDs(ids...)
	return amc
}

// AddMaintenanceDetails adds the "maintenance_details" edges to the AssetMaintenanceDetails entity.
func (amc *AssetMaintenanceCreate) AddMaintenanceDetails(a ...*AssetMaintenanceDetails) *AssetMaintenanceCreate {
	ids := make([]uint64, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return amc.AddMaintenanceDetailIDs(ids...)
}

// Mutation returns the AssetMaintenanceMutation object of the builder.
func (amc *AssetMaintenanceCreate) Mutation() *AssetMaintenanceMutation {
	return amc.mutation
}

// Save creates the AssetMaintenance in the database.
func (amc *AssetMaintenanceCreate) Save(ctx context.Context) (*AssetMaintenance, error) {
	if err := amc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, amc.sqlSave, amc.mutation, amc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (amc *AssetMaintenanceCreate) SaveX(ctx context.Context) *AssetMaintenance {
	v, err := amc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (amc *AssetMaintenanceCreate) Exec(ctx context.Context) error {
	_, err := amc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amc *AssetMaintenanceCreate) ExecX(ctx context.Context) {
	if err := amc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (amc *AssetMaintenanceCreate) defaults() error {
	if _, ok := amc.mutation.CreatedAt(); !ok {
		if assetmaintenance.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized assetmaintenance.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := assetmaintenance.DefaultCreatedAt()
		amc.mutation.SetCreatedAt(v)
	}
	if _, ok := amc.mutation.UpdatedAt(); !ok {
		if assetmaintenance.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized assetmaintenance.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := assetmaintenance.DefaultUpdatedAt()
		amc.mutation.SetUpdatedAt(v)
	}
	if _, ok := amc.mutation.Status(); !ok {
		v := assetmaintenance.DefaultStatus
		amc.mutation.SetStatus(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (amc *AssetMaintenanceCreate) check() error {
	if _, ok := amc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "AssetMaintenance.created_at"`)}
	}
	if _, ok := amc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "AssetMaintenance.updated_at"`)}
	}
	if _, ok := amc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "AssetMaintenance.status"`)}
	}
	return nil
}

func (amc *AssetMaintenanceCreate) sqlSave(ctx context.Context) (*AssetMaintenance, error) {
	if err := amc.check(); err != nil {
		return nil, err
	}
	_node, _spec := amc.createSpec()
	if err := sqlgraph.CreateNode(ctx, amc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	amc.mutation.id = &_node.ID
	amc.mutation.done = true
	return _node, nil
}

func (amc *AssetMaintenanceCreate) createSpec() (*AssetMaintenance, *sqlgraph.CreateSpec) {
	var (
		_node = &AssetMaintenance{config: amc.config}
		_spec = sqlgraph.NewCreateSpec(assetmaintenance.Table, sqlgraph.NewFieldSpec(assetmaintenance.FieldID, field.TypeUint64))
	)
	_spec.OnConflict = amc.conflict
	if value, ok := amc.mutation.CreatedAt(); ok {
		_spec.SetField(assetmaintenance.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := amc.mutation.UpdatedAt(); ok {
		_spec.SetField(assetmaintenance.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := amc.mutation.DeletedAt(); ok {
		_spec.SetField(assetmaintenance.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	if value, ok := amc.mutation.Creator(); ok {
		_spec.SetField(assetmaintenance.FieldCreator, field.TypeJSON, value)
		_node.Creator = value
	}
	if value, ok := amc.mutation.LastModifier(); ok {
		_spec.SetField(assetmaintenance.FieldLastModifier, field.TypeJSON, value)
		_node.LastModifier = value
	}
	if value, ok := amc.mutation.Remark(); ok {
		_spec.SetField(assetmaintenance.FieldRemark, field.TypeString, value)
		_node.Remark = value
	}
	if value, ok := amc.mutation.Reason(); ok {
		_spec.SetField(assetmaintenance.FieldReason, field.TypeString, value)
		_node.Reason = value
	}
	if value, ok := amc.mutation.Content(); ok {
		_spec.SetField(assetmaintenance.FieldContent, field.TypeString, value)
		_node.Content = value
	}
	if value, ok := amc.mutation.Status(); ok {
		_spec.SetField(assetmaintenance.FieldStatus, field.TypeUint8, value)
		_node.Status = value
	}
	if nodes := amc.mutation.CabinetIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   assetmaintenance.CabinetTable,
			Columns: []string{assetmaintenance.CabinetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cabinet.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.CabinetID = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := amc.mutation.MaintainerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   assetmaintenance.MaintainerTable,
			Columns: []string{assetmaintenance.MaintainerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(maintainer.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.MaintainerID = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := amc.mutation.MaintenanceDetailsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   assetmaintenance.MaintenanceDetailsTable,
			Columns: []string{assetmaintenance.MaintenanceDetailsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(assetmaintenancedetails.FieldID, field.TypeUint64),
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
//	client.AssetMaintenance.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AssetMaintenanceUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (amc *AssetMaintenanceCreate) OnConflict(opts ...sql.ConflictOption) *AssetMaintenanceUpsertOne {
	amc.conflict = opts
	return &AssetMaintenanceUpsertOne{
		create: amc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.AssetMaintenance.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (amc *AssetMaintenanceCreate) OnConflictColumns(columns ...string) *AssetMaintenanceUpsertOne {
	amc.conflict = append(amc.conflict, sql.ConflictColumns(columns...))
	return &AssetMaintenanceUpsertOne{
		create: amc,
	}
}

type (
	// AssetMaintenanceUpsertOne is the builder for "upsert"-ing
	//  one AssetMaintenance node.
	AssetMaintenanceUpsertOne struct {
		create *AssetMaintenanceCreate
	}

	// AssetMaintenanceUpsert is the "OnConflict" setter.
	AssetMaintenanceUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *AssetMaintenanceUpsert) SetUpdatedAt(v time.Time) *AssetMaintenanceUpsert {
	u.Set(assetmaintenance.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *AssetMaintenanceUpsert) UpdateUpdatedAt() *AssetMaintenanceUpsert {
	u.SetExcluded(assetmaintenance.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *AssetMaintenanceUpsert) SetDeletedAt(v time.Time) *AssetMaintenanceUpsert {
	u.Set(assetmaintenance.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *AssetMaintenanceUpsert) UpdateDeletedAt() *AssetMaintenanceUpsert {
	u.SetExcluded(assetmaintenance.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *AssetMaintenanceUpsert) ClearDeletedAt() *AssetMaintenanceUpsert {
	u.SetNull(assetmaintenance.FieldDeletedAt)
	return u
}

// SetLastModifier sets the "last_modifier" field.
func (u *AssetMaintenanceUpsert) SetLastModifier(v *model.Modifier) *AssetMaintenanceUpsert {
	u.Set(assetmaintenance.FieldLastModifier, v)
	return u
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *AssetMaintenanceUpsert) UpdateLastModifier() *AssetMaintenanceUpsert {
	u.SetExcluded(assetmaintenance.FieldLastModifier)
	return u
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *AssetMaintenanceUpsert) ClearLastModifier() *AssetMaintenanceUpsert {
	u.SetNull(assetmaintenance.FieldLastModifier)
	return u
}

// SetRemark sets the "remark" field.
func (u *AssetMaintenanceUpsert) SetRemark(v string) *AssetMaintenanceUpsert {
	u.Set(assetmaintenance.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *AssetMaintenanceUpsert) UpdateRemark() *AssetMaintenanceUpsert {
	u.SetExcluded(assetmaintenance.FieldRemark)
	return u
}

// ClearRemark clears the value of the "remark" field.
func (u *AssetMaintenanceUpsert) ClearRemark() *AssetMaintenanceUpsert {
	u.SetNull(assetmaintenance.FieldRemark)
	return u
}

// SetCabinetID sets the "cabinet_id" field.
func (u *AssetMaintenanceUpsert) SetCabinetID(v uint64) *AssetMaintenanceUpsert {
	u.Set(assetmaintenance.FieldCabinetID, v)
	return u
}

// UpdateCabinetID sets the "cabinet_id" field to the value that was provided on create.
func (u *AssetMaintenanceUpsert) UpdateCabinetID() *AssetMaintenanceUpsert {
	u.SetExcluded(assetmaintenance.FieldCabinetID)
	return u
}

// ClearCabinetID clears the value of the "cabinet_id" field.
func (u *AssetMaintenanceUpsert) ClearCabinetID() *AssetMaintenanceUpsert {
	u.SetNull(assetmaintenance.FieldCabinetID)
	return u
}

// SetMaintainerID sets the "maintainer_id" field.
func (u *AssetMaintenanceUpsert) SetMaintainerID(v uint64) *AssetMaintenanceUpsert {
	u.Set(assetmaintenance.FieldMaintainerID, v)
	return u
}

// UpdateMaintainerID sets the "maintainer_id" field to the value that was provided on create.
func (u *AssetMaintenanceUpsert) UpdateMaintainerID() *AssetMaintenanceUpsert {
	u.SetExcluded(assetmaintenance.FieldMaintainerID)
	return u
}

// ClearMaintainerID clears the value of the "maintainer_id" field.
func (u *AssetMaintenanceUpsert) ClearMaintainerID() *AssetMaintenanceUpsert {
	u.SetNull(assetmaintenance.FieldMaintainerID)
	return u
}

// SetReason sets the "reason" field.
func (u *AssetMaintenanceUpsert) SetReason(v string) *AssetMaintenanceUpsert {
	u.Set(assetmaintenance.FieldReason, v)
	return u
}

// UpdateReason sets the "reason" field to the value that was provided on create.
func (u *AssetMaintenanceUpsert) UpdateReason() *AssetMaintenanceUpsert {
	u.SetExcluded(assetmaintenance.FieldReason)
	return u
}

// ClearReason clears the value of the "reason" field.
func (u *AssetMaintenanceUpsert) ClearReason() *AssetMaintenanceUpsert {
	u.SetNull(assetmaintenance.FieldReason)
	return u
}

// SetContent sets the "content" field.
func (u *AssetMaintenanceUpsert) SetContent(v string) *AssetMaintenanceUpsert {
	u.Set(assetmaintenance.FieldContent, v)
	return u
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *AssetMaintenanceUpsert) UpdateContent() *AssetMaintenanceUpsert {
	u.SetExcluded(assetmaintenance.FieldContent)
	return u
}

// ClearContent clears the value of the "content" field.
func (u *AssetMaintenanceUpsert) ClearContent() *AssetMaintenanceUpsert {
	u.SetNull(assetmaintenance.FieldContent)
	return u
}

// SetStatus sets the "status" field.
func (u *AssetMaintenanceUpsert) SetStatus(v uint8) *AssetMaintenanceUpsert {
	u.Set(assetmaintenance.FieldStatus, v)
	return u
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *AssetMaintenanceUpsert) UpdateStatus() *AssetMaintenanceUpsert {
	u.SetExcluded(assetmaintenance.FieldStatus)
	return u
}

// AddStatus adds v to the "status" field.
func (u *AssetMaintenanceUpsert) AddStatus(v uint8) *AssetMaintenanceUpsert {
	u.Add(assetmaintenance.FieldStatus, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.AssetMaintenance.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *AssetMaintenanceUpsertOne) UpdateNewValues() *AssetMaintenanceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(assetmaintenance.FieldCreatedAt)
		}
		if _, exists := u.create.mutation.Creator(); exists {
			s.SetIgnore(assetmaintenance.FieldCreator)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.AssetMaintenance.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *AssetMaintenanceUpsertOne) Ignore() *AssetMaintenanceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AssetMaintenanceUpsertOne) DoNothing() *AssetMaintenanceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AssetMaintenanceCreate.OnConflict
// documentation for more info.
func (u *AssetMaintenanceUpsertOne) Update(set func(*AssetMaintenanceUpsert)) *AssetMaintenanceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AssetMaintenanceUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *AssetMaintenanceUpsertOne) SetUpdatedAt(v time.Time) *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertOne) UpdateUpdatedAt() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *AssetMaintenanceUpsertOne) SetDeletedAt(v time.Time) *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertOne) UpdateDeletedAt() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *AssetMaintenanceUpsertOne) ClearDeletedAt() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *AssetMaintenanceUpsertOne) SetLastModifier(v *model.Modifier) *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertOne) UpdateLastModifier() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *AssetMaintenanceUpsertOne) ClearLastModifier() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *AssetMaintenanceUpsertOne) SetRemark(v string) *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertOne) UpdateRemark() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *AssetMaintenanceUpsertOne) ClearRemark() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearRemark()
	})
}

// SetCabinetID sets the "cabinet_id" field.
func (u *AssetMaintenanceUpsertOne) SetCabinetID(v uint64) *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetCabinetID(v)
	})
}

// UpdateCabinetID sets the "cabinet_id" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertOne) UpdateCabinetID() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateCabinetID()
	})
}

// ClearCabinetID clears the value of the "cabinet_id" field.
func (u *AssetMaintenanceUpsertOne) ClearCabinetID() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearCabinetID()
	})
}

// SetMaintainerID sets the "maintainer_id" field.
func (u *AssetMaintenanceUpsertOne) SetMaintainerID(v uint64) *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetMaintainerID(v)
	})
}

// UpdateMaintainerID sets the "maintainer_id" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertOne) UpdateMaintainerID() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateMaintainerID()
	})
}

// ClearMaintainerID clears the value of the "maintainer_id" field.
func (u *AssetMaintenanceUpsertOne) ClearMaintainerID() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearMaintainerID()
	})
}

// SetReason sets the "reason" field.
func (u *AssetMaintenanceUpsertOne) SetReason(v string) *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetReason(v)
	})
}

// UpdateReason sets the "reason" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertOne) UpdateReason() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateReason()
	})
}

// ClearReason clears the value of the "reason" field.
func (u *AssetMaintenanceUpsertOne) ClearReason() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearReason()
	})
}

// SetContent sets the "content" field.
func (u *AssetMaintenanceUpsertOne) SetContent(v string) *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertOne) UpdateContent() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateContent()
	})
}

// ClearContent clears the value of the "content" field.
func (u *AssetMaintenanceUpsertOne) ClearContent() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearContent()
	})
}

// SetStatus sets the "status" field.
func (u *AssetMaintenanceUpsertOne) SetStatus(v uint8) *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetStatus(v)
	})
}

// AddStatus adds v to the "status" field.
func (u *AssetMaintenanceUpsertOne) AddStatus(v uint8) *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.AddStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertOne) UpdateStatus() *AssetMaintenanceUpsertOne {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateStatus()
	})
}

// Exec executes the query.
func (u *AssetMaintenanceUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AssetMaintenanceCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AssetMaintenanceUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *AssetMaintenanceUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *AssetMaintenanceUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// AssetMaintenanceCreateBulk is the builder for creating many AssetMaintenance entities in bulk.
type AssetMaintenanceCreateBulk struct {
	config
	err      error
	builders []*AssetMaintenanceCreate
	conflict []sql.ConflictOption
}

// Save creates the AssetMaintenance entities in the database.
func (amcb *AssetMaintenanceCreateBulk) Save(ctx context.Context) ([]*AssetMaintenance, error) {
	if amcb.err != nil {
		return nil, amcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(amcb.builders))
	nodes := make([]*AssetMaintenance, len(amcb.builders))
	mutators := make([]Mutator, len(amcb.builders))
	for i := range amcb.builders {
		func(i int, root context.Context) {
			builder := amcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AssetMaintenanceMutation)
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
					_, err = mutators[i+1].Mutate(root, amcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = amcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, amcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, amcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (amcb *AssetMaintenanceCreateBulk) SaveX(ctx context.Context) []*AssetMaintenance {
	v, err := amcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (amcb *AssetMaintenanceCreateBulk) Exec(ctx context.Context) error {
	_, err := amcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amcb *AssetMaintenanceCreateBulk) ExecX(ctx context.Context) {
	if err := amcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.AssetMaintenance.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AssetMaintenanceUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (amcb *AssetMaintenanceCreateBulk) OnConflict(opts ...sql.ConflictOption) *AssetMaintenanceUpsertBulk {
	amcb.conflict = opts
	return &AssetMaintenanceUpsertBulk{
		create: amcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.AssetMaintenance.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (amcb *AssetMaintenanceCreateBulk) OnConflictColumns(columns ...string) *AssetMaintenanceUpsertBulk {
	amcb.conflict = append(amcb.conflict, sql.ConflictColumns(columns...))
	return &AssetMaintenanceUpsertBulk{
		create: amcb,
	}
}

// AssetMaintenanceUpsertBulk is the builder for "upsert"-ing
// a bulk of AssetMaintenance nodes.
type AssetMaintenanceUpsertBulk struct {
	create *AssetMaintenanceCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.AssetMaintenance.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *AssetMaintenanceUpsertBulk) UpdateNewValues() *AssetMaintenanceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(assetmaintenance.FieldCreatedAt)
			}
			if _, exists := b.mutation.Creator(); exists {
				s.SetIgnore(assetmaintenance.FieldCreator)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.AssetMaintenance.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *AssetMaintenanceUpsertBulk) Ignore() *AssetMaintenanceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AssetMaintenanceUpsertBulk) DoNothing() *AssetMaintenanceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AssetMaintenanceCreateBulk.OnConflict
// documentation for more info.
func (u *AssetMaintenanceUpsertBulk) Update(set func(*AssetMaintenanceUpsert)) *AssetMaintenanceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AssetMaintenanceUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *AssetMaintenanceUpsertBulk) SetUpdatedAt(v time.Time) *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertBulk) UpdateUpdatedAt() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *AssetMaintenanceUpsertBulk) SetDeletedAt(v time.Time) *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertBulk) UpdateDeletedAt() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *AssetMaintenanceUpsertBulk) ClearDeletedAt() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *AssetMaintenanceUpsertBulk) SetLastModifier(v *model.Modifier) *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertBulk) UpdateLastModifier() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *AssetMaintenanceUpsertBulk) ClearLastModifier() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *AssetMaintenanceUpsertBulk) SetRemark(v string) *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertBulk) UpdateRemark() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *AssetMaintenanceUpsertBulk) ClearRemark() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearRemark()
	})
}

// SetCabinetID sets the "cabinet_id" field.
func (u *AssetMaintenanceUpsertBulk) SetCabinetID(v uint64) *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetCabinetID(v)
	})
}

// UpdateCabinetID sets the "cabinet_id" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertBulk) UpdateCabinetID() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateCabinetID()
	})
}

// ClearCabinetID clears the value of the "cabinet_id" field.
func (u *AssetMaintenanceUpsertBulk) ClearCabinetID() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearCabinetID()
	})
}

// SetMaintainerID sets the "maintainer_id" field.
func (u *AssetMaintenanceUpsertBulk) SetMaintainerID(v uint64) *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetMaintainerID(v)
	})
}

// UpdateMaintainerID sets the "maintainer_id" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertBulk) UpdateMaintainerID() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateMaintainerID()
	})
}

// ClearMaintainerID clears the value of the "maintainer_id" field.
func (u *AssetMaintenanceUpsertBulk) ClearMaintainerID() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearMaintainerID()
	})
}

// SetReason sets the "reason" field.
func (u *AssetMaintenanceUpsertBulk) SetReason(v string) *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetReason(v)
	})
}

// UpdateReason sets the "reason" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertBulk) UpdateReason() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateReason()
	})
}

// ClearReason clears the value of the "reason" field.
func (u *AssetMaintenanceUpsertBulk) ClearReason() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearReason()
	})
}

// SetContent sets the "content" field.
func (u *AssetMaintenanceUpsertBulk) SetContent(v string) *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertBulk) UpdateContent() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateContent()
	})
}

// ClearContent clears the value of the "content" field.
func (u *AssetMaintenanceUpsertBulk) ClearContent() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.ClearContent()
	})
}

// SetStatus sets the "status" field.
func (u *AssetMaintenanceUpsertBulk) SetStatus(v uint8) *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.SetStatus(v)
	})
}

// AddStatus adds v to the "status" field.
func (u *AssetMaintenanceUpsertBulk) AddStatus(v uint8) *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.AddStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *AssetMaintenanceUpsertBulk) UpdateStatus() *AssetMaintenanceUpsertBulk {
	return u.Update(func(s *AssetMaintenanceUpsert) {
		s.UpdateStatus()
	})
}

// Exec executes the query.
func (u *AssetMaintenanceUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the AssetMaintenanceCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AssetMaintenanceCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AssetMaintenanceUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
