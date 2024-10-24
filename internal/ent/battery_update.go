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
	"github.com/auroraride/adapter"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/batteryflow"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// BatteryUpdate is the builder for updating Battery entities.
type BatteryUpdate struct {
	config
	hooks     []Hook
	mutation  *BatteryMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the BatteryUpdate builder.
func (bu *BatteryUpdate) Where(ps ...predicate.Battery) *BatteryUpdate {
	bu.mutation.Where(ps...)
	return bu
}

// SetUpdatedAt sets the "updated_at" field.
func (bu *BatteryUpdate) SetUpdatedAt(t time.Time) *BatteryUpdate {
	bu.mutation.SetUpdatedAt(t)
	return bu
}

// SetDeletedAt sets the "deleted_at" field.
func (bu *BatteryUpdate) SetDeletedAt(t time.Time) *BatteryUpdate {
	bu.mutation.SetDeletedAt(t)
	return bu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableDeletedAt(t *time.Time) *BatteryUpdate {
	if t != nil {
		bu.SetDeletedAt(*t)
	}
	return bu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (bu *BatteryUpdate) ClearDeletedAt() *BatteryUpdate {
	bu.mutation.ClearDeletedAt()
	return bu
}

// SetLastModifier sets the "last_modifier" field.
func (bu *BatteryUpdate) SetLastModifier(m *model.Modifier) *BatteryUpdate {
	bu.mutation.SetLastModifier(m)
	return bu
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (bu *BatteryUpdate) ClearLastModifier() *BatteryUpdate {
	bu.mutation.ClearLastModifier()
	return bu
}

// SetRemark sets the "remark" field.
func (bu *BatteryUpdate) SetRemark(s string) *BatteryUpdate {
	bu.mutation.SetRemark(s)
	return bu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableRemark(s *string) *BatteryUpdate {
	if s != nil {
		bu.SetRemark(*s)
	}
	return bu
}

// ClearRemark clears the value of the "remark" field.
func (bu *BatteryUpdate) ClearRemark() *BatteryUpdate {
	bu.mutation.ClearRemark()
	return bu
}

// SetCityID sets the "city_id" field.
func (bu *BatteryUpdate) SetCityID(u uint64) *BatteryUpdate {
	bu.mutation.SetCityID(u)
	return bu
}

// SetNillableCityID sets the "city_id" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableCityID(u *uint64) *BatteryUpdate {
	if u != nil {
		bu.SetCityID(*u)
	}
	return bu
}

// ClearCityID clears the value of the "city_id" field.
func (bu *BatteryUpdate) ClearCityID() *BatteryUpdate {
	bu.mutation.ClearCityID()
	return bu
}

// SetRiderID sets the "rider_id" field.
func (bu *BatteryUpdate) SetRiderID(u uint64) *BatteryUpdate {
	bu.mutation.ResetRiderID()
	bu.mutation.SetRiderID(u)
	return bu
}

// SetNillableRiderID sets the "rider_id" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableRiderID(u *uint64) *BatteryUpdate {
	if u != nil {
		bu.SetRiderID(*u)
	}
	return bu
}

// AddRiderID adds u to the "rider_id" field.
func (bu *BatteryUpdate) AddRiderID(u int64) *BatteryUpdate {
	bu.mutation.AddRiderID(u)
	return bu
}

// ClearRiderID clears the value of the "rider_id" field.
func (bu *BatteryUpdate) ClearRiderID() *BatteryUpdate {
	bu.mutation.ClearRiderID()
	return bu
}

// SetCabinetID sets the "cabinet_id" field.
func (bu *BatteryUpdate) SetCabinetID(u uint64) *BatteryUpdate {
	bu.mutation.ResetCabinetID()
	bu.mutation.SetCabinetID(u)
	return bu
}

// SetNillableCabinetID sets the "cabinet_id" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableCabinetID(u *uint64) *BatteryUpdate {
	if u != nil {
		bu.SetCabinetID(*u)
	}
	return bu
}

// AddCabinetID adds u to the "cabinet_id" field.
func (bu *BatteryUpdate) AddCabinetID(u int64) *BatteryUpdate {
	bu.mutation.AddCabinetID(u)
	return bu
}

// ClearCabinetID clears the value of the "cabinet_id" field.
func (bu *BatteryUpdate) ClearCabinetID() *BatteryUpdate {
	bu.mutation.ClearCabinetID()
	return bu
}

// SetSubscribeID sets the "subscribe_id" field.
func (bu *BatteryUpdate) SetSubscribeID(u uint64) *BatteryUpdate {
	bu.mutation.ResetSubscribeID()
	bu.mutation.SetSubscribeID(u)
	return bu
}

// SetNillableSubscribeID sets the "subscribe_id" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableSubscribeID(u *uint64) *BatteryUpdate {
	if u != nil {
		bu.SetSubscribeID(*u)
	}
	return bu
}

// AddSubscribeID adds u to the "subscribe_id" field.
func (bu *BatteryUpdate) AddSubscribeID(u int64) *BatteryUpdate {
	bu.mutation.AddSubscribeID(u)
	return bu
}

// ClearSubscribeID clears the value of the "subscribe_id" field.
func (bu *BatteryUpdate) ClearSubscribeID() *BatteryUpdate {
	bu.mutation.ClearSubscribeID()
	return bu
}

// SetEnterpriseID sets the "enterprise_id" field.
func (bu *BatteryUpdate) SetEnterpriseID(u uint64) *BatteryUpdate {
	bu.mutation.ResetEnterpriseID()
	bu.mutation.SetEnterpriseID(u)
	return bu
}

// SetNillableEnterpriseID sets the "enterprise_id" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableEnterpriseID(u *uint64) *BatteryUpdate {
	if u != nil {
		bu.SetEnterpriseID(*u)
	}
	return bu
}

// AddEnterpriseID adds u to the "enterprise_id" field.
func (bu *BatteryUpdate) AddEnterpriseID(u int64) *BatteryUpdate {
	bu.mutation.AddEnterpriseID(u)
	return bu
}

// ClearEnterpriseID clears the value of the "enterprise_id" field.
func (bu *BatteryUpdate) ClearEnterpriseID() *BatteryUpdate {
	bu.mutation.ClearEnterpriseID()
	return bu
}

// SetStationID sets the "station_id" field.
func (bu *BatteryUpdate) SetStationID(u uint64) *BatteryUpdate {
	bu.mutation.ResetStationID()
	bu.mutation.SetStationID(u)
	return bu
}

// SetNillableStationID sets the "station_id" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableStationID(u *uint64) *BatteryUpdate {
	if u != nil {
		bu.SetStationID(*u)
	}
	return bu
}

// AddStationID adds u to the "station_id" field.
func (bu *BatteryUpdate) AddStationID(u int64) *BatteryUpdate {
	bu.mutation.AddStationID(u)
	return bu
}

// ClearStationID clears the value of the "station_id" field.
func (bu *BatteryUpdate) ClearStationID() *BatteryUpdate {
	bu.mutation.ClearStationID()
	return bu
}

// SetSn sets the "sn" field.
func (bu *BatteryUpdate) SetSn(s string) *BatteryUpdate {
	bu.mutation.SetSn(s)
	return bu
}

// SetNillableSn sets the "sn" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableSn(s *string) *BatteryUpdate {
	if s != nil {
		bu.SetSn(*s)
	}
	return bu
}

// SetBrand sets the "brand" field.
func (bu *BatteryUpdate) SetBrand(ab adapter.BatteryBrand) *BatteryUpdate {
	bu.mutation.SetBrand(ab)
	return bu
}

// SetNillableBrand sets the "brand" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableBrand(ab *adapter.BatteryBrand) *BatteryUpdate {
	if ab != nil {
		bu.SetBrand(*ab)
	}
	return bu
}

// SetEnable sets the "enable" field.
func (bu *BatteryUpdate) SetEnable(b bool) *BatteryUpdate {
	bu.mutation.SetEnable(b)
	return bu
}

// SetNillableEnable sets the "enable" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableEnable(b *bool) *BatteryUpdate {
	if b != nil {
		bu.SetEnable(*b)
	}
	return bu
}

// SetModel sets the "model" field.
func (bu *BatteryUpdate) SetModel(s string) *BatteryUpdate {
	bu.mutation.SetModel(s)
	return bu
}

// SetNillableModel sets the "model" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableModel(s *string) *BatteryUpdate {
	if s != nil {
		bu.SetModel(*s)
	}
	return bu
}

// SetOrdinal sets the "ordinal" field.
func (bu *BatteryUpdate) SetOrdinal(i int) *BatteryUpdate {
	bu.mutation.ResetOrdinal()
	bu.mutation.SetOrdinal(i)
	return bu
}

// SetNillableOrdinal sets the "ordinal" field if the given value is not nil.
func (bu *BatteryUpdate) SetNillableOrdinal(i *int) *BatteryUpdate {
	if i != nil {
		bu.SetOrdinal(*i)
	}
	return bu
}

// AddOrdinal adds i to the "ordinal" field.
func (bu *BatteryUpdate) AddOrdinal(i int) *BatteryUpdate {
	bu.mutation.AddOrdinal(i)
	return bu
}

// ClearOrdinal clears the value of the "ordinal" field.
func (bu *BatteryUpdate) ClearOrdinal() *BatteryUpdate {
	bu.mutation.ClearOrdinal()
	return bu
}

// SetCity sets the "city" edge to the City entity.
func (bu *BatteryUpdate) SetCity(c *City) *BatteryUpdate {
	return bu.SetCityID(c.ID)
}

// AddFlowIDs adds the "flows" edge to the BatteryFlow entity by IDs.
func (bu *BatteryUpdate) AddFlowIDs(ids ...uint64) *BatteryUpdate {
	bu.mutation.AddFlowIDs(ids...)
	return bu
}

// AddFlows adds the "flows" edges to the BatteryFlow entity.
func (bu *BatteryUpdate) AddFlows(b ...*BatteryFlow) *BatteryUpdate {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return bu.AddFlowIDs(ids...)
}

// Mutation returns the BatteryMutation object of the builder.
func (bu *BatteryUpdate) Mutation() *BatteryMutation {
	return bu.mutation
}

// ClearCity clears the "city" edge to the City entity.
func (bu *BatteryUpdate) ClearCity() *BatteryUpdate {
	bu.mutation.ClearCity()
	return bu
}

// ClearFlows clears all "flows" edges to the BatteryFlow entity.
func (bu *BatteryUpdate) ClearFlows() *BatteryUpdate {
	bu.mutation.ClearFlows()
	return bu
}

// RemoveFlowIDs removes the "flows" edge to BatteryFlow entities by IDs.
func (bu *BatteryUpdate) RemoveFlowIDs(ids ...uint64) *BatteryUpdate {
	bu.mutation.RemoveFlowIDs(ids...)
	return bu
}

// RemoveFlows removes "flows" edges to BatteryFlow entities.
func (bu *BatteryUpdate) RemoveFlows(b ...*BatteryFlow) *BatteryUpdate {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return bu.RemoveFlowIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (bu *BatteryUpdate) Save(ctx context.Context) (int, error) {
	if err := bu.defaults(); err != nil {
		return 0, err
	}
	return withHooks(ctx, bu.sqlSave, bu.mutation, bu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (bu *BatteryUpdate) SaveX(ctx context.Context) int {
	affected, err := bu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (bu *BatteryUpdate) Exec(ctx context.Context) error {
	_, err := bu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bu *BatteryUpdate) ExecX(ctx context.Context) {
	if err := bu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (bu *BatteryUpdate) defaults() error {
	if _, ok := bu.mutation.UpdatedAt(); !ok {
		if battery.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized battery.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := battery.UpdateDefaultUpdatedAt()
		bu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (bu *BatteryUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *BatteryUpdate {
	bu.modifiers = append(bu.modifiers, modifiers...)
	return bu
}

func (bu *BatteryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(battery.Table, battery.Columns, sqlgraph.NewFieldSpec(battery.FieldID, field.TypeUint64))
	if ps := bu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := bu.mutation.UpdatedAt(); ok {
		_spec.SetField(battery.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := bu.mutation.DeletedAt(); ok {
		_spec.SetField(battery.FieldDeletedAt, field.TypeTime, value)
	}
	if bu.mutation.DeletedAtCleared() {
		_spec.ClearField(battery.FieldDeletedAt, field.TypeTime)
	}
	if bu.mutation.CreatorCleared() {
		_spec.ClearField(battery.FieldCreator, field.TypeJSON)
	}
	if value, ok := bu.mutation.LastModifier(); ok {
		_spec.SetField(battery.FieldLastModifier, field.TypeJSON, value)
	}
	if bu.mutation.LastModifierCleared() {
		_spec.ClearField(battery.FieldLastModifier, field.TypeJSON)
	}
	if value, ok := bu.mutation.Remark(); ok {
		_spec.SetField(battery.FieldRemark, field.TypeString, value)
	}
	if bu.mutation.RemarkCleared() {
		_spec.ClearField(battery.FieldRemark, field.TypeString)
	}
	if value, ok := bu.mutation.RiderID(); ok {
		_spec.SetField(battery.FieldRiderID, field.TypeUint64, value)
	}
	if value, ok := bu.mutation.AddedRiderID(); ok {
		_spec.AddField(battery.FieldRiderID, field.TypeUint64, value)
	}
	if bu.mutation.RiderIDCleared() {
		_spec.ClearField(battery.FieldRiderID, field.TypeUint64)
	}
	if value, ok := bu.mutation.CabinetID(); ok {
		_spec.SetField(battery.FieldCabinetID, field.TypeUint64, value)
	}
	if value, ok := bu.mutation.AddedCabinetID(); ok {
		_spec.AddField(battery.FieldCabinetID, field.TypeUint64, value)
	}
	if bu.mutation.CabinetIDCleared() {
		_spec.ClearField(battery.FieldCabinetID, field.TypeUint64)
	}
	if value, ok := bu.mutation.SubscribeID(); ok {
		_spec.SetField(battery.FieldSubscribeID, field.TypeUint64, value)
	}
	if value, ok := bu.mutation.AddedSubscribeID(); ok {
		_spec.AddField(battery.FieldSubscribeID, field.TypeUint64, value)
	}
	if bu.mutation.SubscribeIDCleared() {
		_spec.ClearField(battery.FieldSubscribeID, field.TypeUint64)
	}
	if value, ok := bu.mutation.EnterpriseID(); ok {
		_spec.SetField(battery.FieldEnterpriseID, field.TypeUint64, value)
	}
	if value, ok := bu.mutation.AddedEnterpriseID(); ok {
		_spec.AddField(battery.FieldEnterpriseID, field.TypeUint64, value)
	}
	if bu.mutation.EnterpriseIDCleared() {
		_spec.ClearField(battery.FieldEnterpriseID, field.TypeUint64)
	}
	if value, ok := bu.mutation.StationID(); ok {
		_spec.SetField(battery.FieldStationID, field.TypeUint64, value)
	}
	if value, ok := bu.mutation.AddedStationID(); ok {
		_spec.AddField(battery.FieldStationID, field.TypeUint64, value)
	}
	if bu.mutation.StationIDCleared() {
		_spec.ClearField(battery.FieldStationID, field.TypeUint64)
	}
	if value, ok := bu.mutation.Sn(); ok {
		_spec.SetField(battery.FieldSn, field.TypeString, value)
	}
	if value, ok := bu.mutation.Brand(); ok {
		_spec.SetField(battery.FieldBrand, field.TypeOther, value)
	}
	if value, ok := bu.mutation.Enable(); ok {
		_spec.SetField(battery.FieldEnable, field.TypeBool, value)
	}
	if value, ok := bu.mutation.Model(); ok {
		_spec.SetField(battery.FieldModel, field.TypeString, value)
	}
	if value, ok := bu.mutation.Ordinal(); ok {
		_spec.SetField(battery.FieldOrdinal, field.TypeInt, value)
	}
	if value, ok := bu.mutation.AddedOrdinal(); ok {
		_spec.AddField(battery.FieldOrdinal, field.TypeInt, value)
	}
	if bu.mutation.OrdinalCleared() {
		_spec.ClearField(battery.FieldOrdinal, field.TypeInt)
	}
	if bu.mutation.CityCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   battery.CityTable,
			Columns: []string{battery.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.CityIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   battery.CityTable,
			Columns: []string{battery.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if bu.mutation.FlowsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   battery.FlowsTable,
			Columns: []string{battery.FlowsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(batteryflow.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.RemovedFlowsIDs(); len(nodes) > 0 && !bu.mutation.FlowsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   battery.FlowsTable,
			Columns: []string{battery.FlowsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(batteryflow.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.FlowsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   battery.FlowsTable,
			Columns: []string{battery.FlowsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(batteryflow.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(bu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, bu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{battery.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	bu.mutation.done = true
	return n, nil
}

// BatteryUpdateOne is the builder for updating a single Battery entity.
type BatteryUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *BatteryMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (buo *BatteryUpdateOne) SetUpdatedAt(t time.Time) *BatteryUpdateOne {
	buo.mutation.SetUpdatedAt(t)
	return buo
}

// SetDeletedAt sets the "deleted_at" field.
func (buo *BatteryUpdateOne) SetDeletedAt(t time.Time) *BatteryUpdateOne {
	buo.mutation.SetDeletedAt(t)
	return buo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableDeletedAt(t *time.Time) *BatteryUpdateOne {
	if t != nil {
		buo.SetDeletedAt(*t)
	}
	return buo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (buo *BatteryUpdateOne) ClearDeletedAt() *BatteryUpdateOne {
	buo.mutation.ClearDeletedAt()
	return buo
}

// SetLastModifier sets the "last_modifier" field.
func (buo *BatteryUpdateOne) SetLastModifier(m *model.Modifier) *BatteryUpdateOne {
	buo.mutation.SetLastModifier(m)
	return buo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (buo *BatteryUpdateOne) ClearLastModifier() *BatteryUpdateOne {
	buo.mutation.ClearLastModifier()
	return buo
}

// SetRemark sets the "remark" field.
func (buo *BatteryUpdateOne) SetRemark(s string) *BatteryUpdateOne {
	buo.mutation.SetRemark(s)
	return buo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableRemark(s *string) *BatteryUpdateOne {
	if s != nil {
		buo.SetRemark(*s)
	}
	return buo
}

// ClearRemark clears the value of the "remark" field.
func (buo *BatteryUpdateOne) ClearRemark() *BatteryUpdateOne {
	buo.mutation.ClearRemark()
	return buo
}

// SetCityID sets the "city_id" field.
func (buo *BatteryUpdateOne) SetCityID(u uint64) *BatteryUpdateOne {
	buo.mutation.SetCityID(u)
	return buo
}

// SetNillableCityID sets the "city_id" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableCityID(u *uint64) *BatteryUpdateOne {
	if u != nil {
		buo.SetCityID(*u)
	}
	return buo
}

// ClearCityID clears the value of the "city_id" field.
func (buo *BatteryUpdateOne) ClearCityID() *BatteryUpdateOne {
	buo.mutation.ClearCityID()
	return buo
}

// SetRiderID sets the "rider_id" field.
func (buo *BatteryUpdateOne) SetRiderID(u uint64) *BatteryUpdateOne {
	buo.mutation.ResetRiderID()
	buo.mutation.SetRiderID(u)
	return buo
}

// SetNillableRiderID sets the "rider_id" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableRiderID(u *uint64) *BatteryUpdateOne {
	if u != nil {
		buo.SetRiderID(*u)
	}
	return buo
}

// AddRiderID adds u to the "rider_id" field.
func (buo *BatteryUpdateOne) AddRiderID(u int64) *BatteryUpdateOne {
	buo.mutation.AddRiderID(u)
	return buo
}

// ClearRiderID clears the value of the "rider_id" field.
func (buo *BatteryUpdateOne) ClearRiderID() *BatteryUpdateOne {
	buo.mutation.ClearRiderID()
	return buo
}

// SetCabinetID sets the "cabinet_id" field.
func (buo *BatteryUpdateOne) SetCabinetID(u uint64) *BatteryUpdateOne {
	buo.mutation.ResetCabinetID()
	buo.mutation.SetCabinetID(u)
	return buo
}

// SetNillableCabinetID sets the "cabinet_id" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableCabinetID(u *uint64) *BatteryUpdateOne {
	if u != nil {
		buo.SetCabinetID(*u)
	}
	return buo
}

// AddCabinetID adds u to the "cabinet_id" field.
func (buo *BatteryUpdateOne) AddCabinetID(u int64) *BatteryUpdateOne {
	buo.mutation.AddCabinetID(u)
	return buo
}

// ClearCabinetID clears the value of the "cabinet_id" field.
func (buo *BatteryUpdateOne) ClearCabinetID() *BatteryUpdateOne {
	buo.mutation.ClearCabinetID()
	return buo
}

// SetSubscribeID sets the "subscribe_id" field.
func (buo *BatteryUpdateOne) SetSubscribeID(u uint64) *BatteryUpdateOne {
	buo.mutation.ResetSubscribeID()
	buo.mutation.SetSubscribeID(u)
	return buo
}

// SetNillableSubscribeID sets the "subscribe_id" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableSubscribeID(u *uint64) *BatteryUpdateOne {
	if u != nil {
		buo.SetSubscribeID(*u)
	}
	return buo
}

// AddSubscribeID adds u to the "subscribe_id" field.
func (buo *BatteryUpdateOne) AddSubscribeID(u int64) *BatteryUpdateOne {
	buo.mutation.AddSubscribeID(u)
	return buo
}

// ClearSubscribeID clears the value of the "subscribe_id" field.
func (buo *BatteryUpdateOne) ClearSubscribeID() *BatteryUpdateOne {
	buo.mutation.ClearSubscribeID()
	return buo
}

// SetEnterpriseID sets the "enterprise_id" field.
func (buo *BatteryUpdateOne) SetEnterpriseID(u uint64) *BatteryUpdateOne {
	buo.mutation.ResetEnterpriseID()
	buo.mutation.SetEnterpriseID(u)
	return buo
}

// SetNillableEnterpriseID sets the "enterprise_id" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableEnterpriseID(u *uint64) *BatteryUpdateOne {
	if u != nil {
		buo.SetEnterpriseID(*u)
	}
	return buo
}

// AddEnterpriseID adds u to the "enterprise_id" field.
func (buo *BatteryUpdateOne) AddEnterpriseID(u int64) *BatteryUpdateOne {
	buo.mutation.AddEnterpriseID(u)
	return buo
}

// ClearEnterpriseID clears the value of the "enterprise_id" field.
func (buo *BatteryUpdateOne) ClearEnterpriseID() *BatteryUpdateOne {
	buo.mutation.ClearEnterpriseID()
	return buo
}

// SetStationID sets the "station_id" field.
func (buo *BatteryUpdateOne) SetStationID(u uint64) *BatteryUpdateOne {
	buo.mutation.ResetStationID()
	buo.mutation.SetStationID(u)
	return buo
}

// SetNillableStationID sets the "station_id" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableStationID(u *uint64) *BatteryUpdateOne {
	if u != nil {
		buo.SetStationID(*u)
	}
	return buo
}

// AddStationID adds u to the "station_id" field.
func (buo *BatteryUpdateOne) AddStationID(u int64) *BatteryUpdateOne {
	buo.mutation.AddStationID(u)
	return buo
}

// ClearStationID clears the value of the "station_id" field.
func (buo *BatteryUpdateOne) ClearStationID() *BatteryUpdateOne {
	buo.mutation.ClearStationID()
	return buo
}

// SetSn sets the "sn" field.
func (buo *BatteryUpdateOne) SetSn(s string) *BatteryUpdateOne {
	buo.mutation.SetSn(s)
	return buo
}

// SetNillableSn sets the "sn" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableSn(s *string) *BatteryUpdateOne {
	if s != nil {
		buo.SetSn(*s)
	}
	return buo
}

// SetBrand sets the "brand" field.
func (buo *BatteryUpdateOne) SetBrand(ab adapter.BatteryBrand) *BatteryUpdateOne {
	buo.mutation.SetBrand(ab)
	return buo
}

// SetNillableBrand sets the "brand" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableBrand(ab *adapter.BatteryBrand) *BatteryUpdateOne {
	if ab != nil {
		buo.SetBrand(*ab)
	}
	return buo
}

// SetEnable sets the "enable" field.
func (buo *BatteryUpdateOne) SetEnable(b bool) *BatteryUpdateOne {
	buo.mutation.SetEnable(b)
	return buo
}

// SetNillableEnable sets the "enable" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableEnable(b *bool) *BatteryUpdateOne {
	if b != nil {
		buo.SetEnable(*b)
	}
	return buo
}

// SetModel sets the "model" field.
func (buo *BatteryUpdateOne) SetModel(s string) *BatteryUpdateOne {
	buo.mutation.SetModel(s)
	return buo
}

// SetNillableModel sets the "model" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableModel(s *string) *BatteryUpdateOne {
	if s != nil {
		buo.SetModel(*s)
	}
	return buo
}

// SetOrdinal sets the "ordinal" field.
func (buo *BatteryUpdateOne) SetOrdinal(i int) *BatteryUpdateOne {
	buo.mutation.ResetOrdinal()
	buo.mutation.SetOrdinal(i)
	return buo
}

// SetNillableOrdinal sets the "ordinal" field if the given value is not nil.
func (buo *BatteryUpdateOne) SetNillableOrdinal(i *int) *BatteryUpdateOne {
	if i != nil {
		buo.SetOrdinal(*i)
	}
	return buo
}

// AddOrdinal adds i to the "ordinal" field.
func (buo *BatteryUpdateOne) AddOrdinal(i int) *BatteryUpdateOne {
	buo.mutation.AddOrdinal(i)
	return buo
}

// ClearOrdinal clears the value of the "ordinal" field.
func (buo *BatteryUpdateOne) ClearOrdinal() *BatteryUpdateOne {
	buo.mutation.ClearOrdinal()
	return buo
}

// SetCity sets the "city" edge to the City entity.
func (buo *BatteryUpdateOne) SetCity(c *City) *BatteryUpdateOne {
	return buo.SetCityID(c.ID)
}

// AddFlowIDs adds the "flows" edge to the BatteryFlow entity by IDs.
func (buo *BatteryUpdateOne) AddFlowIDs(ids ...uint64) *BatteryUpdateOne {
	buo.mutation.AddFlowIDs(ids...)
	return buo
}

// AddFlows adds the "flows" edges to the BatteryFlow entity.
func (buo *BatteryUpdateOne) AddFlows(b ...*BatteryFlow) *BatteryUpdateOne {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return buo.AddFlowIDs(ids...)
}

// Mutation returns the BatteryMutation object of the builder.
func (buo *BatteryUpdateOne) Mutation() *BatteryMutation {
	return buo.mutation
}

// ClearCity clears the "city" edge to the City entity.
func (buo *BatteryUpdateOne) ClearCity() *BatteryUpdateOne {
	buo.mutation.ClearCity()
	return buo
}

// ClearFlows clears all "flows" edges to the BatteryFlow entity.
func (buo *BatteryUpdateOne) ClearFlows() *BatteryUpdateOne {
	buo.mutation.ClearFlows()
	return buo
}

// RemoveFlowIDs removes the "flows" edge to BatteryFlow entities by IDs.
func (buo *BatteryUpdateOne) RemoveFlowIDs(ids ...uint64) *BatteryUpdateOne {
	buo.mutation.RemoveFlowIDs(ids...)
	return buo
}

// RemoveFlows removes "flows" edges to BatteryFlow entities.
func (buo *BatteryUpdateOne) RemoveFlows(b ...*BatteryFlow) *BatteryUpdateOne {
	ids := make([]uint64, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return buo.RemoveFlowIDs(ids...)
}

// Where appends a list predicates to the BatteryUpdate builder.
func (buo *BatteryUpdateOne) Where(ps ...predicate.Battery) *BatteryUpdateOne {
	buo.mutation.Where(ps...)
	return buo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (buo *BatteryUpdateOne) Select(field string, fields ...string) *BatteryUpdateOne {
	buo.fields = append([]string{field}, fields...)
	return buo
}

// Save executes the query and returns the updated Battery entity.
func (buo *BatteryUpdateOne) Save(ctx context.Context) (*Battery, error) {
	if err := buo.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, buo.sqlSave, buo.mutation, buo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (buo *BatteryUpdateOne) SaveX(ctx context.Context) *Battery {
	node, err := buo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (buo *BatteryUpdateOne) Exec(ctx context.Context) error {
	_, err := buo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (buo *BatteryUpdateOne) ExecX(ctx context.Context) {
	if err := buo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (buo *BatteryUpdateOne) defaults() error {
	if _, ok := buo.mutation.UpdatedAt(); !ok {
		if battery.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized battery.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := battery.UpdateDefaultUpdatedAt()
		buo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (buo *BatteryUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *BatteryUpdateOne {
	buo.modifiers = append(buo.modifiers, modifiers...)
	return buo
}

func (buo *BatteryUpdateOne) sqlSave(ctx context.Context) (_node *Battery, err error) {
	_spec := sqlgraph.NewUpdateSpec(battery.Table, battery.Columns, sqlgraph.NewFieldSpec(battery.FieldID, field.TypeUint64))
	id, ok := buo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Battery.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := buo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, battery.FieldID)
		for _, f := range fields {
			if !battery.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != battery.FieldID {
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
		_spec.SetField(battery.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := buo.mutation.DeletedAt(); ok {
		_spec.SetField(battery.FieldDeletedAt, field.TypeTime, value)
	}
	if buo.mutation.DeletedAtCleared() {
		_spec.ClearField(battery.FieldDeletedAt, field.TypeTime)
	}
	if buo.mutation.CreatorCleared() {
		_spec.ClearField(battery.FieldCreator, field.TypeJSON)
	}
	if value, ok := buo.mutation.LastModifier(); ok {
		_spec.SetField(battery.FieldLastModifier, field.TypeJSON, value)
	}
	if buo.mutation.LastModifierCleared() {
		_spec.ClearField(battery.FieldLastModifier, field.TypeJSON)
	}
	if value, ok := buo.mutation.Remark(); ok {
		_spec.SetField(battery.FieldRemark, field.TypeString, value)
	}
	if buo.mutation.RemarkCleared() {
		_spec.ClearField(battery.FieldRemark, field.TypeString)
	}
	if value, ok := buo.mutation.RiderID(); ok {
		_spec.SetField(battery.FieldRiderID, field.TypeUint64, value)
	}
	if value, ok := buo.mutation.AddedRiderID(); ok {
		_spec.AddField(battery.FieldRiderID, field.TypeUint64, value)
	}
	if buo.mutation.RiderIDCleared() {
		_spec.ClearField(battery.FieldRiderID, field.TypeUint64)
	}
	if value, ok := buo.mutation.CabinetID(); ok {
		_spec.SetField(battery.FieldCabinetID, field.TypeUint64, value)
	}
	if value, ok := buo.mutation.AddedCabinetID(); ok {
		_spec.AddField(battery.FieldCabinetID, field.TypeUint64, value)
	}
	if buo.mutation.CabinetIDCleared() {
		_spec.ClearField(battery.FieldCabinetID, field.TypeUint64)
	}
	if value, ok := buo.mutation.SubscribeID(); ok {
		_spec.SetField(battery.FieldSubscribeID, field.TypeUint64, value)
	}
	if value, ok := buo.mutation.AddedSubscribeID(); ok {
		_spec.AddField(battery.FieldSubscribeID, field.TypeUint64, value)
	}
	if buo.mutation.SubscribeIDCleared() {
		_spec.ClearField(battery.FieldSubscribeID, field.TypeUint64)
	}
	if value, ok := buo.mutation.EnterpriseID(); ok {
		_spec.SetField(battery.FieldEnterpriseID, field.TypeUint64, value)
	}
	if value, ok := buo.mutation.AddedEnterpriseID(); ok {
		_spec.AddField(battery.FieldEnterpriseID, field.TypeUint64, value)
	}
	if buo.mutation.EnterpriseIDCleared() {
		_spec.ClearField(battery.FieldEnterpriseID, field.TypeUint64)
	}
	if value, ok := buo.mutation.StationID(); ok {
		_spec.SetField(battery.FieldStationID, field.TypeUint64, value)
	}
	if value, ok := buo.mutation.AddedStationID(); ok {
		_spec.AddField(battery.FieldStationID, field.TypeUint64, value)
	}
	if buo.mutation.StationIDCleared() {
		_spec.ClearField(battery.FieldStationID, field.TypeUint64)
	}
	if value, ok := buo.mutation.Sn(); ok {
		_spec.SetField(battery.FieldSn, field.TypeString, value)
	}
	if value, ok := buo.mutation.Brand(); ok {
		_spec.SetField(battery.FieldBrand, field.TypeOther, value)
	}
	if value, ok := buo.mutation.Enable(); ok {
		_spec.SetField(battery.FieldEnable, field.TypeBool, value)
	}
	if value, ok := buo.mutation.Model(); ok {
		_spec.SetField(battery.FieldModel, field.TypeString, value)
	}
	if value, ok := buo.mutation.Ordinal(); ok {
		_spec.SetField(battery.FieldOrdinal, field.TypeInt, value)
	}
	if value, ok := buo.mutation.AddedOrdinal(); ok {
		_spec.AddField(battery.FieldOrdinal, field.TypeInt, value)
	}
	if buo.mutation.OrdinalCleared() {
		_spec.ClearField(battery.FieldOrdinal, field.TypeInt)
	}
	if buo.mutation.CityCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   battery.CityTable,
			Columns: []string{battery.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.CityIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   battery.CityTable,
			Columns: []string{battery.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(city.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if buo.mutation.FlowsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   battery.FlowsTable,
			Columns: []string{battery.FlowsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(batteryflow.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.RemovedFlowsIDs(); len(nodes) > 0 && !buo.mutation.FlowsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   battery.FlowsTable,
			Columns: []string{battery.FlowsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(batteryflow.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.FlowsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   battery.FlowsTable,
			Columns: []string{battery.FlowsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(batteryflow.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(buo.modifiers...)
	_node = &Battery{config: buo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, buo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{battery.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	buo.mutation.done = true
	return _node, nil
}
