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
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterprisebatteryswap"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/stock"
)

// EnterpriseStationCreate is the builder for creating a EnterpriseStation entity.
type EnterpriseStationCreate struct {
	config
	mutation *EnterpriseStationMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (esc *EnterpriseStationCreate) SetCreatedAt(t time.Time) *EnterpriseStationCreate {
	esc.mutation.SetCreatedAt(t)
	return esc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (esc *EnterpriseStationCreate) SetNillableCreatedAt(t *time.Time) *EnterpriseStationCreate {
	if t != nil {
		esc.SetCreatedAt(*t)
	}
	return esc
}

// SetUpdatedAt sets the "updated_at" field.
func (esc *EnterpriseStationCreate) SetUpdatedAt(t time.Time) *EnterpriseStationCreate {
	esc.mutation.SetUpdatedAt(t)
	return esc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (esc *EnterpriseStationCreate) SetNillableUpdatedAt(t *time.Time) *EnterpriseStationCreate {
	if t != nil {
		esc.SetUpdatedAt(*t)
	}
	return esc
}

// SetDeletedAt sets the "deleted_at" field.
func (esc *EnterpriseStationCreate) SetDeletedAt(t time.Time) *EnterpriseStationCreate {
	esc.mutation.SetDeletedAt(t)
	return esc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (esc *EnterpriseStationCreate) SetNillableDeletedAt(t *time.Time) *EnterpriseStationCreate {
	if t != nil {
		esc.SetDeletedAt(*t)
	}
	return esc
}

// SetCreator sets the "creator" field.
func (esc *EnterpriseStationCreate) SetCreator(m *model.Modifier) *EnterpriseStationCreate {
	esc.mutation.SetCreator(m)
	return esc
}

// SetLastModifier sets the "last_modifier" field.
func (esc *EnterpriseStationCreate) SetLastModifier(m *model.Modifier) *EnterpriseStationCreate {
	esc.mutation.SetLastModifier(m)
	return esc
}

// SetRemark sets the "remark" field.
func (esc *EnterpriseStationCreate) SetRemark(s string) *EnterpriseStationCreate {
	esc.mutation.SetRemark(s)
	return esc
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (esc *EnterpriseStationCreate) SetNillableRemark(s *string) *EnterpriseStationCreate {
	if s != nil {
		esc.SetRemark(*s)
	}
	return esc
}

// SetCityID sets the "city_id" field.
func (esc *EnterpriseStationCreate) SetCityID(u uint64) *EnterpriseStationCreate {
	esc.mutation.SetCityID(u)
	return esc
}

// SetNillableCityID sets the "city_id" field if the given value is not nil.
func (esc *EnterpriseStationCreate) SetNillableCityID(u *uint64) *EnterpriseStationCreate {
	if u != nil {
		esc.SetCityID(*u)
	}
	return esc
}

// SetEnterpriseID sets the "enterprise_id" field.
func (esc *EnterpriseStationCreate) SetEnterpriseID(u uint64) *EnterpriseStationCreate {
	esc.mutation.SetEnterpriseID(u)
	return esc
}

// SetName sets the "name" field.
func (esc *EnterpriseStationCreate) SetName(s string) *EnterpriseStationCreate {
	esc.mutation.SetName(s)
	return esc
}

// SetCity sets the "city" edge to the City entity.
func (esc *EnterpriseStationCreate) SetCity(c *City) *EnterpriseStationCreate {
	return esc.SetCityID(c.ID)
}

// SetEnterprise sets the "enterprise" edge to the Enterprise entity.
func (esc *EnterpriseStationCreate) SetEnterprise(e *Enterprise) *EnterpriseStationCreate {
	return esc.SetEnterpriseID(e.ID)
}

// AddAgentIDs adds the "agents" edge to the Agent entity by IDs.
func (esc *EnterpriseStationCreate) AddAgentIDs(ids ...uint64) *EnterpriseStationCreate {
	esc.mutation.AddAgentIDs(ids...)
	return esc
}

// AddAgents adds the "agents" edges to the Agent entity.
func (esc *EnterpriseStationCreate) AddAgents(a ...*Agent) *EnterpriseStationCreate {
	ids := make([]uint64, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return esc.AddAgentIDs(ids...)
}

// AddSwapPutinBatteryIDs adds the "swap_putin_batteries" edge to the EnterpriseBatterySwap entity by IDs.
func (esc *EnterpriseStationCreate) AddSwapPutinBatteryIDs(ids ...uint64) *EnterpriseStationCreate {
	esc.mutation.AddSwapPutinBatteryIDs(ids...)
	return esc
}

// AddSwapPutinBatteries adds the "swap_putin_batteries" edges to the EnterpriseBatterySwap entity.
func (esc *EnterpriseStationCreate) AddSwapPutinBatteries(e ...*EnterpriseBatterySwap) *EnterpriseStationCreate {
	ids := make([]uint64, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return esc.AddSwapPutinBatteryIDs(ids...)
}

// AddSwapPutoutBatteryIDs adds the "swap_putout_batteries" edge to the EnterpriseBatterySwap entity by IDs.
func (esc *EnterpriseStationCreate) AddSwapPutoutBatteryIDs(ids ...uint64) *EnterpriseStationCreate {
	esc.mutation.AddSwapPutoutBatteryIDs(ids...)
	return esc
}

// AddSwapPutoutBatteries adds the "swap_putout_batteries" edges to the EnterpriseBatterySwap entity.
func (esc *EnterpriseStationCreate) AddSwapPutoutBatteries(e ...*EnterpriseBatterySwap) *EnterpriseStationCreate {
	ids := make([]uint64, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return esc.AddSwapPutoutBatteryIDs(ids...)
}

// AddCabinetIDs adds the "cabinets" edge to the Cabinet entity by IDs.
func (esc *EnterpriseStationCreate) AddCabinetIDs(ids ...uint64) *EnterpriseStationCreate {
	esc.mutation.AddCabinetIDs(ids...)
	return esc
}

// AddCabinets adds the "cabinets" edges to the Cabinet entity.
func (esc *EnterpriseStationCreate) AddCabinets(c ...*Cabinet) *EnterpriseStationCreate {
	ids := make([]uint64, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return esc.AddCabinetIDs(ids...)
}

// AddAssetIDs adds the "asset" edge to the Asset entity by IDs.
func (esc *EnterpriseStationCreate) AddAssetIDs(ids ...uint64) *EnterpriseStationCreate {
	esc.mutation.AddAssetIDs(ids...)
	return esc
}

// AddAsset adds the "asset" edges to the Asset entity.
func (esc *EnterpriseStationCreate) AddAsset(a ...*Asset) *EnterpriseStationCreate {
	ids := make([]uint64, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return esc.AddAssetIDs(ids...)
}

// AddStockIDs adds the "stocks" edge to the Stock entity by IDs.
func (esc *EnterpriseStationCreate) AddStockIDs(ids ...uint64) *EnterpriseStationCreate {
	esc.mutation.AddStockIDs(ids...)
	return esc
}

// AddStocks adds the "stocks" edges to the Stock entity.
func (esc *EnterpriseStationCreate) AddStocks(s ...*Stock) *EnterpriseStationCreate {
	ids := make([]uint64, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return esc.AddStockIDs(ids...)
}

// AddRentAssetIDs adds the "rent_asset" edge to the Asset entity by IDs.
func (esc *EnterpriseStationCreate) AddRentAssetIDs(ids ...uint64) *EnterpriseStationCreate {
	esc.mutation.AddRentAssetIDs(ids...)
	return esc
}

// AddRentAsset adds the "rent_asset" edges to the Asset entity.
func (esc *EnterpriseStationCreate) AddRentAsset(a ...*Asset) *EnterpriseStationCreate {
	ids := make([]uint64, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return esc.AddRentAssetIDs(ids...)
}

// Mutation returns the EnterpriseStationMutation object of the builder.
func (esc *EnterpriseStationCreate) Mutation() *EnterpriseStationMutation {
	return esc.mutation
}

// Save creates the EnterpriseStation in the database.
func (esc *EnterpriseStationCreate) Save(ctx context.Context) (*EnterpriseStation, error) {
	if err := esc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, esc.sqlSave, esc.mutation, esc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (esc *EnterpriseStationCreate) SaveX(ctx context.Context) *EnterpriseStation {
	v, err := esc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (esc *EnterpriseStationCreate) Exec(ctx context.Context) error {
	_, err := esc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (esc *EnterpriseStationCreate) ExecX(ctx context.Context) {
	if err := esc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (esc *EnterpriseStationCreate) defaults() error {
	if _, ok := esc.mutation.CreatedAt(); !ok {
		if enterprisestation.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized enterprisestation.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := enterprisestation.DefaultCreatedAt()
		esc.mutation.SetCreatedAt(v)
	}
	if _, ok := esc.mutation.UpdatedAt(); !ok {
		if enterprisestation.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized enterprisestation.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := enterprisestation.DefaultUpdatedAt()
		esc.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (esc *EnterpriseStationCreate) check() error {
	if _, ok := esc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "EnterpriseStation.created_at"`)}
	}
	if _, ok := esc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "EnterpriseStation.updated_at"`)}
	}
	if _, ok := esc.mutation.EnterpriseID(); !ok {
		return &ValidationError{Name: "enterprise_id", err: errors.New(`ent: missing required field "EnterpriseStation.enterprise_id"`)}
	}
	if _, ok := esc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "EnterpriseStation.name"`)}
	}
	if len(esc.mutation.EnterpriseIDs()) == 0 {
		return &ValidationError{Name: "enterprise", err: errors.New(`ent: missing required edge "EnterpriseStation.enterprise"`)}
	}
	return nil
}

func (esc *EnterpriseStationCreate) sqlSave(ctx context.Context) (*EnterpriseStation, error) {
	if err := esc.check(); err != nil {
		return nil, err
	}
	_node, _spec := esc.createSpec()
	if err := sqlgraph.CreateNode(ctx, esc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = uint64(id)
	esc.mutation.id = &_node.ID
	esc.mutation.done = true
	return _node, nil
}

func (esc *EnterpriseStationCreate) createSpec() (*EnterpriseStation, *sqlgraph.CreateSpec) {
	var (
		_node = &EnterpriseStation{config: esc.config}
		_spec = sqlgraph.NewCreateSpec(enterprisestation.Table, sqlgraph.NewFieldSpec(enterprisestation.FieldID, field.TypeUint64))
	)
	_spec.OnConflict = esc.conflict
	if value, ok := esc.mutation.CreatedAt(); ok {
		_spec.SetField(enterprisestation.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := esc.mutation.UpdatedAt(); ok {
		_spec.SetField(enterprisestation.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := esc.mutation.DeletedAt(); ok {
		_spec.SetField(enterprisestation.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	if value, ok := esc.mutation.Creator(); ok {
		_spec.SetField(enterprisestation.FieldCreator, field.TypeJSON, value)
		_node.Creator = value
	}
	if value, ok := esc.mutation.LastModifier(); ok {
		_spec.SetField(enterprisestation.FieldLastModifier, field.TypeJSON, value)
		_node.LastModifier = value
	}
	if value, ok := esc.mutation.Remark(); ok {
		_spec.SetField(enterprisestation.FieldRemark, field.TypeString, value)
		_node.Remark = value
	}
	if value, ok := esc.mutation.Name(); ok {
		_spec.SetField(enterprisestation.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if nodes := esc.mutation.CityIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   enterprisestation.CityTable,
			Columns: []string{enterprisestation.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.CityID = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := esc.mutation.EnterpriseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   enterprisestation.EnterpriseTable,
			Columns: []string{enterprisestation.EnterpriseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(enterprise.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.EnterpriseID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := esc.mutation.AgentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   enterprisestation.AgentsTable,
			Columns: enterprisestation.AgentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agent.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := esc.mutation.SwapPutinBatteriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   enterprisestation.SwapPutinBatteriesTable,
			Columns: []string{enterprisestation.SwapPutinBatteriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(enterprisebatteryswap.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := esc.mutation.SwapPutoutBatteriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   enterprisestation.SwapPutoutBatteriesTable,
			Columns: []string{enterprisestation.SwapPutoutBatteriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(enterprisebatteryswap.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := esc.mutation.CabinetsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   enterprisestation.CabinetsTable,
			Columns: []string{enterprisestation.CabinetsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cabinet.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := esc.mutation.AssetIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   enterprisestation.AssetTable,
			Columns: []string{enterprisestation.AssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := esc.mutation.StocksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   enterprisestation.StocksTable,
			Columns: []string{enterprisestation.StocksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(stock.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := esc.mutation.RentAssetIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   enterprisestation.RentAssetTable,
			Columns: []string{enterprisestation.RentAssetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeUint64),
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
//	client.EnterpriseStation.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnterpriseStationUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (esc *EnterpriseStationCreate) OnConflict(opts ...sql.ConflictOption) *EnterpriseStationUpsertOne {
	esc.conflict = opts
	return &EnterpriseStationUpsertOne{
		create: esc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.EnterpriseStation.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (esc *EnterpriseStationCreate) OnConflictColumns(columns ...string) *EnterpriseStationUpsertOne {
	esc.conflict = append(esc.conflict, sql.ConflictColumns(columns...))
	return &EnterpriseStationUpsertOne{
		create: esc,
	}
}

type (
	// EnterpriseStationUpsertOne is the builder for "upsert"-ing
	//  one EnterpriseStation node.
	EnterpriseStationUpsertOne struct {
		create *EnterpriseStationCreate
	}

	// EnterpriseStationUpsert is the "OnConflict" setter.
	EnterpriseStationUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *EnterpriseStationUpsert) SetUpdatedAt(v time.Time) *EnterpriseStationUpsert {
	u.Set(enterprisestation.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *EnterpriseStationUpsert) UpdateUpdatedAt() *EnterpriseStationUpsert {
	u.SetExcluded(enterprisestation.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *EnterpriseStationUpsert) SetDeletedAt(v time.Time) *EnterpriseStationUpsert {
	u.Set(enterprisestation.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *EnterpriseStationUpsert) UpdateDeletedAt() *EnterpriseStationUpsert {
	u.SetExcluded(enterprisestation.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *EnterpriseStationUpsert) ClearDeletedAt() *EnterpriseStationUpsert {
	u.SetNull(enterprisestation.FieldDeletedAt)
	return u
}

// SetLastModifier sets the "last_modifier" field.
func (u *EnterpriseStationUpsert) SetLastModifier(v *model.Modifier) *EnterpriseStationUpsert {
	u.Set(enterprisestation.FieldLastModifier, v)
	return u
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *EnterpriseStationUpsert) UpdateLastModifier() *EnterpriseStationUpsert {
	u.SetExcluded(enterprisestation.FieldLastModifier)
	return u
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *EnterpriseStationUpsert) ClearLastModifier() *EnterpriseStationUpsert {
	u.SetNull(enterprisestation.FieldLastModifier)
	return u
}

// SetRemark sets the "remark" field.
func (u *EnterpriseStationUpsert) SetRemark(v string) *EnterpriseStationUpsert {
	u.Set(enterprisestation.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *EnterpriseStationUpsert) UpdateRemark() *EnterpriseStationUpsert {
	u.SetExcluded(enterprisestation.FieldRemark)
	return u
}

// ClearRemark clears the value of the "remark" field.
func (u *EnterpriseStationUpsert) ClearRemark() *EnterpriseStationUpsert {
	u.SetNull(enterprisestation.FieldRemark)
	return u
}

// SetCityID sets the "city_id" field.
func (u *EnterpriseStationUpsert) SetCityID(v uint64) *EnterpriseStationUpsert {
	u.Set(enterprisestation.FieldCityID, v)
	return u
}

// UpdateCityID sets the "city_id" field to the value that was provided on create.
func (u *EnterpriseStationUpsert) UpdateCityID() *EnterpriseStationUpsert {
	u.SetExcluded(enterprisestation.FieldCityID)
	return u
}

// ClearCityID clears the value of the "city_id" field.
func (u *EnterpriseStationUpsert) ClearCityID() *EnterpriseStationUpsert {
	u.SetNull(enterprisestation.FieldCityID)
	return u
}

// SetEnterpriseID sets the "enterprise_id" field.
func (u *EnterpriseStationUpsert) SetEnterpriseID(v uint64) *EnterpriseStationUpsert {
	u.Set(enterprisestation.FieldEnterpriseID, v)
	return u
}

// UpdateEnterpriseID sets the "enterprise_id" field to the value that was provided on create.
func (u *EnterpriseStationUpsert) UpdateEnterpriseID() *EnterpriseStationUpsert {
	u.SetExcluded(enterprisestation.FieldEnterpriseID)
	return u
}

// SetName sets the "name" field.
func (u *EnterpriseStationUpsert) SetName(v string) *EnterpriseStationUpsert {
	u.Set(enterprisestation.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *EnterpriseStationUpsert) UpdateName() *EnterpriseStationUpsert {
	u.SetExcluded(enterprisestation.FieldName)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.EnterpriseStation.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *EnterpriseStationUpsertOne) UpdateNewValues() *EnterpriseStationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(enterprisestation.FieldCreatedAt)
		}
		if _, exists := u.create.mutation.Creator(); exists {
			s.SetIgnore(enterprisestation.FieldCreator)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.EnterpriseStation.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *EnterpriseStationUpsertOne) Ignore() *EnterpriseStationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EnterpriseStationUpsertOne) DoNothing() *EnterpriseStationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EnterpriseStationCreate.OnConflict
// documentation for more info.
func (u *EnterpriseStationUpsertOne) Update(set func(*EnterpriseStationUpsert)) *EnterpriseStationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EnterpriseStationUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *EnterpriseStationUpsertOne) SetUpdatedAt(v time.Time) *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *EnterpriseStationUpsertOne) UpdateUpdatedAt() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *EnterpriseStationUpsertOne) SetDeletedAt(v time.Time) *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *EnterpriseStationUpsertOne) UpdateDeletedAt() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *EnterpriseStationUpsertOne) ClearDeletedAt() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *EnterpriseStationUpsertOne) SetLastModifier(v *model.Modifier) *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *EnterpriseStationUpsertOne) UpdateLastModifier() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *EnterpriseStationUpsertOne) ClearLastModifier() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *EnterpriseStationUpsertOne) SetRemark(v string) *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *EnterpriseStationUpsertOne) UpdateRemark() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *EnterpriseStationUpsertOne) ClearRemark() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.ClearRemark()
	})
}

// SetCityID sets the "city_id" field.
func (u *EnterpriseStationUpsertOne) SetCityID(v uint64) *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetCityID(v)
	})
}

// UpdateCityID sets the "city_id" field to the value that was provided on create.
func (u *EnterpriseStationUpsertOne) UpdateCityID() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateCityID()
	})
}

// ClearCityID clears the value of the "city_id" field.
func (u *EnterpriseStationUpsertOne) ClearCityID() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.ClearCityID()
	})
}

// SetEnterpriseID sets the "enterprise_id" field.
func (u *EnterpriseStationUpsertOne) SetEnterpriseID(v uint64) *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetEnterpriseID(v)
	})
}

// UpdateEnterpriseID sets the "enterprise_id" field to the value that was provided on create.
func (u *EnterpriseStationUpsertOne) UpdateEnterpriseID() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateEnterpriseID()
	})
}

// SetName sets the "name" field.
func (u *EnterpriseStationUpsertOne) SetName(v string) *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *EnterpriseStationUpsertOne) UpdateName() *EnterpriseStationUpsertOne {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateName()
	})
}

// Exec executes the query.
func (u *EnterpriseStationUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for EnterpriseStationCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EnterpriseStationUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *EnterpriseStationUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *EnterpriseStationUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// EnterpriseStationCreateBulk is the builder for creating many EnterpriseStation entities in bulk.
type EnterpriseStationCreateBulk struct {
	config
	err      error
	builders []*EnterpriseStationCreate
	conflict []sql.ConflictOption
}

// Save creates the EnterpriseStation entities in the database.
func (escb *EnterpriseStationCreateBulk) Save(ctx context.Context) ([]*EnterpriseStation, error) {
	if escb.err != nil {
		return nil, escb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(escb.builders))
	nodes := make([]*EnterpriseStation, len(escb.builders))
	mutators := make([]Mutator, len(escb.builders))
	for i := range escb.builders {
		func(i int, root context.Context) {
			builder := escb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EnterpriseStationMutation)
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
					_, err = mutators[i+1].Mutate(root, escb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = escb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, escb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, escb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (escb *EnterpriseStationCreateBulk) SaveX(ctx context.Context) []*EnterpriseStation {
	v, err := escb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (escb *EnterpriseStationCreateBulk) Exec(ctx context.Context) error {
	_, err := escb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (escb *EnterpriseStationCreateBulk) ExecX(ctx context.Context) {
	if err := escb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.EnterpriseStation.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnterpriseStationUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (escb *EnterpriseStationCreateBulk) OnConflict(opts ...sql.ConflictOption) *EnterpriseStationUpsertBulk {
	escb.conflict = opts
	return &EnterpriseStationUpsertBulk{
		create: escb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.EnterpriseStation.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (escb *EnterpriseStationCreateBulk) OnConflictColumns(columns ...string) *EnterpriseStationUpsertBulk {
	escb.conflict = append(escb.conflict, sql.ConflictColumns(columns...))
	return &EnterpriseStationUpsertBulk{
		create: escb,
	}
}

// EnterpriseStationUpsertBulk is the builder for "upsert"-ing
// a bulk of EnterpriseStation nodes.
type EnterpriseStationUpsertBulk struct {
	create *EnterpriseStationCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.EnterpriseStation.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *EnterpriseStationUpsertBulk) UpdateNewValues() *EnterpriseStationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(enterprisestation.FieldCreatedAt)
			}
			if _, exists := b.mutation.Creator(); exists {
				s.SetIgnore(enterprisestation.FieldCreator)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.EnterpriseStation.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *EnterpriseStationUpsertBulk) Ignore() *EnterpriseStationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EnterpriseStationUpsertBulk) DoNothing() *EnterpriseStationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EnterpriseStationCreateBulk.OnConflict
// documentation for more info.
func (u *EnterpriseStationUpsertBulk) Update(set func(*EnterpriseStationUpsert)) *EnterpriseStationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EnterpriseStationUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *EnterpriseStationUpsertBulk) SetUpdatedAt(v time.Time) *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *EnterpriseStationUpsertBulk) UpdateUpdatedAt() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *EnterpriseStationUpsertBulk) SetDeletedAt(v time.Time) *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *EnterpriseStationUpsertBulk) UpdateDeletedAt() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *EnterpriseStationUpsertBulk) ClearDeletedAt() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.ClearDeletedAt()
	})
}

// SetLastModifier sets the "last_modifier" field.
func (u *EnterpriseStationUpsertBulk) SetLastModifier(v *model.Modifier) *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetLastModifier(v)
	})
}

// UpdateLastModifier sets the "last_modifier" field to the value that was provided on create.
func (u *EnterpriseStationUpsertBulk) UpdateLastModifier() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateLastModifier()
	})
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (u *EnterpriseStationUpsertBulk) ClearLastModifier() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.ClearLastModifier()
	})
}

// SetRemark sets the "remark" field.
func (u *EnterpriseStationUpsertBulk) SetRemark(v string) *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *EnterpriseStationUpsertBulk) UpdateRemark() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateRemark()
	})
}

// ClearRemark clears the value of the "remark" field.
func (u *EnterpriseStationUpsertBulk) ClearRemark() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.ClearRemark()
	})
}

// SetCityID sets the "city_id" field.
func (u *EnterpriseStationUpsertBulk) SetCityID(v uint64) *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetCityID(v)
	})
}

// UpdateCityID sets the "city_id" field to the value that was provided on create.
func (u *EnterpriseStationUpsertBulk) UpdateCityID() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateCityID()
	})
}

// ClearCityID clears the value of the "city_id" field.
func (u *EnterpriseStationUpsertBulk) ClearCityID() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.ClearCityID()
	})
}

// SetEnterpriseID sets the "enterprise_id" field.
func (u *EnterpriseStationUpsertBulk) SetEnterpriseID(v uint64) *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetEnterpriseID(v)
	})
}

// UpdateEnterpriseID sets the "enterprise_id" field to the value that was provided on create.
func (u *EnterpriseStationUpsertBulk) UpdateEnterpriseID() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateEnterpriseID()
	})
}

// SetName sets the "name" field.
func (u *EnterpriseStationUpsertBulk) SetName(v string) *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *EnterpriseStationUpsertBulk) UpdateName() *EnterpriseStationUpsertBulk {
	return u.Update(func(s *EnterpriseStationUpsert) {
		s.UpdateName()
	})
}

// Exec executes the query.
func (u *EnterpriseStationUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the EnterpriseStationCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for EnterpriseStationCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EnterpriseStationUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
