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
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterpriseprice"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// EnterprisePriceUpdate is the builder for updating EnterprisePrice entities.
type EnterprisePriceUpdate struct {
	config
	hooks    []Hook
	mutation *EnterprisePriceMutation
}

// Where appends a list predicates to the EnterprisePriceUpdate builder.
func (epu *EnterprisePriceUpdate) Where(ps ...predicate.EnterprisePrice) *EnterprisePriceUpdate {
	epu.mutation.Where(ps...)
	return epu
}

// SetUpdatedAt sets the "updated_at" field.
func (epu *EnterprisePriceUpdate) SetUpdatedAt(t time.Time) *EnterprisePriceUpdate {
	epu.mutation.SetUpdatedAt(t)
	return epu
}

// SetDeletedAt sets the "deleted_at" field.
func (epu *EnterprisePriceUpdate) SetDeletedAt(t time.Time) *EnterprisePriceUpdate {
	epu.mutation.SetDeletedAt(t)
	return epu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (epu *EnterprisePriceUpdate) SetNillableDeletedAt(t *time.Time) *EnterprisePriceUpdate {
	if t != nil {
		epu.SetDeletedAt(*t)
	}
	return epu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (epu *EnterprisePriceUpdate) ClearDeletedAt() *EnterprisePriceUpdate {
	epu.mutation.ClearDeletedAt()
	return epu
}

// SetLastModifier sets the "last_modifier" field.
func (epu *EnterprisePriceUpdate) SetLastModifier(m *model.Modifier) *EnterprisePriceUpdate {
	epu.mutation.SetLastModifier(m)
	return epu
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (epu *EnterprisePriceUpdate) ClearLastModifier() *EnterprisePriceUpdate {
	epu.mutation.ClearLastModifier()
	return epu
}

// SetRemark sets the "remark" field.
func (epu *EnterprisePriceUpdate) SetRemark(s string) *EnterprisePriceUpdate {
	epu.mutation.SetRemark(s)
	return epu
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (epu *EnterprisePriceUpdate) SetNillableRemark(s *string) *EnterprisePriceUpdate {
	if s != nil {
		epu.SetRemark(*s)
	}
	return epu
}

// ClearRemark clears the value of the "remark" field.
func (epu *EnterprisePriceUpdate) ClearRemark() *EnterprisePriceUpdate {
	epu.mutation.ClearRemark()
	return epu
}

// SetCityID sets the "city_id" field.
func (epu *EnterprisePriceUpdate) SetCityID(u uint64) *EnterprisePriceUpdate {
	epu.mutation.SetCityID(u)
	return epu
}

// SetEnterpriseID sets the "enterprise_id" field.
func (epu *EnterprisePriceUpdate) SetEnterpriseID(u uint64) *EnterprisePriceUpdate {
	epu.mutation.SetEnterpriseID(u)
	return epu
}

// SetPrice sets the "price" field.
func (epu *EnterprisePriceUpdate) SetPrice(f float64) *EnterprisePriceUpdate {
	epu.mutation.ResetPrice()
	epu.mutation.SetPrice(f)
	return epu
}

// AddPrice adds f to the "price" field.
func (epu *EnterprisePriceUpdate) AddPrice(f float64) *EnterprisePriceUpdate {
	epu.mutation.AddPrice(f)
	return epu
}

// SetVoltage sets the "voltage" field.
func (epu *EnterprisePriceUpdate) SetVoltage(f float64) *EnterprisePriceUpdate {
	epu.mutation.ResetVoltage()
	epu.mutation.SetVoltage(f)
	return epu
}

// AddVoltage adds f to the "voltage" field.
func (epu *EnterprisePriceUpdate) AddVoltage(f float64) *EnterprisePriceUpdate {
	epu.mutation.AddVoltage(f)
	return epu
}

// SetCity sets the "city" edge to the City entity.
func (epu *EnterprisePriceUpdate) SetCity(c *City) *EnterprisePriceUpdate {
	return epu.SetCityID(c.ID)
}

// SetEnterprise sets the "enterprise" edge to the Enterprise entity.
func (epu *EnterprisePriceUpdate) SetEnterprise(e *Enterprise) *EnterprisePriceUpdate {
	return epu.SetEnterpriseID(e.ID)
}

// Mutation returns the EnterprisePriceMutation object of the builder.
func (epu *EnterprisePriceUpdate) Mutation() *EnterprisePriceMutation {
	return epu.mutation
}

// ClearCity clears the "city" edge to the City entity.
func (epu *EnterprisePriceUpdate) ClearCity() *EnterprisePriceUpdate {
	epu.mutation.ClearCity()
	return epu
}

// ClearEnterprise clears the "enterprise" edge to the Enterprise entity.
func (epu *EnterprisePriceUpdate) ClearEnterprise() *EnterprisePriceUpdate {
	epu.mutation.ClearEnterprise()
	return epu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (epu *EnterprisePriceUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if err := epu.defaults(); err != nil {
		return 0, err
	}
	if len(epu.hooks) == 0 {
		if err = epu.check(); err != nil {
			return 0, err
		}
		affected, err = epu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EnterprisePriceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = epu.check(); err != nil {
				return 0, err
			}
			epu.mutation = mutation
			affected, err = epu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(epu.hooks) - 1; i >= 0; i-- {
			if epu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = epu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, epu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (epu *EnterprisePriceUpdate) SaveX(ctx context.Context) int {
	affected, err := epu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (epu *EnterprisePriceUpdate) Exec(ctx context.Context) error {
	_, err := epu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epu *EnterprisePriceUpdate) ExecX(ctx context.Context) {
	if err := epu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (epu *EnterprisePriceUpdate) defaults() error {
	if _, ok := epu.mutation.UpdatedAt(); !ok {
		if enterpriseprice.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized enterpriseprice.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := enterpriseprice.UpdateDefaultUpdatedAt()
		epu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (epu *EnterprisePriceUpdate) check() error {
	if _, ok := epu.mutation.CityID(); epu.mutation.CityCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "EnterprisePrice.city"`)
	}
	if _, ok := epu.mutation.EnterpriseID(); epu.mutation.EnterpriseCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "EnterprisePrice.enterprise"`)
	}
	return nil
}

func (epu *EnterprisePriceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   enterpriseprice.Table,
			Columns: enterpriseprice.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: enterpriseprice.FieldID,
			},
		},
	}
	if ps := epu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := epu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: enterpriseprice.FieldUpdatedAt,
		})
	}
	if value, ok := epu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: enterpriseprice.FieldDeletedAt,
		})
	}
	if epu.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: enterpriseprice.FieldDeletedAt,
		})
	}
	if epu.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: enterpriseprice.FieldCreator,
		})
	}
	if value, ok := epu.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: enterpriseprice.FieldLastModifier,
		})
	}
	if epu.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: enterpriseprice.FieldLastModifier,
		})
	}
	if value, ok := epu.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: enterpriseprice.FieldRemark,
		})
	}
	if epu.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: enterpriseprice.FieldRemark,
		})
	}
	if value, ok := epu.mutation.Price(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: enterpriseprice.FieldPrice,
		})
	}
	if value, ok := epu.mutation.AddedPrice(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: enterpriseprice.FieldPrice,
		})
	}
	if value, ok := epu.mutation.Voltage(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: enterpriseprice.FieldVoltage,
		})
	}
	if value, ok := epu.mutation.AddedVoltage(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: enterpriseprice.FieldVoltage,
		})
	}
	if epu.mutation.CityCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   enterpriseprice.CityTable,
			Columns: []string{enterpriseprice.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: city.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epu.mutation.CityIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   enterpriseprice.CityTable,
			Columns: []string{enterpriseprice.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: city.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epu.mutation.EnterpriseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   enterpriseprice.EnterpriseTable,
			Columns: []string{enterpriseprice.EnterpriseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: enterprise.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epu.mutation.EnterpriseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   enterpriseprice.EnterpriseTable,
			Columns: []string{enterpriseprice.EnterpriseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: enterprise.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, epu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{enterpriseprice.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// EnterprisePriceUpdateOne is the builder for updating a single EnterprisePrice entity.
type EnterprisePriceUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EnterprisePriceMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (epuo *EnterprisePriceUpdateOne) SetUpdatedAt(t time.Time) *EnterprisePriceUpdateOne {
	epuo.mutation.SetUpdatedAt(t)
	return epuo
}

// SetDeletedAt sets the "deleted_at" field.
func (epuo *EnterprisePriceUpdateOne) SetDeletedAt(t time.Time) *EnterprisePriceUpdateOne {
	epuo.mutation.SetDeletedAt(t)
	return epuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (epuo *EnterprisePriceUpdateOne) SetNillableDeletedAt(t *time.Time) *EnterprisePriceUpdateOne {
	if t != nil {
		epuo.SetDeletedAt(*t)
	}
	return epuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (epuo *EnterprisePriceUpdateOne) ClearDeletedAt() *EnterprisePriceUpdateOne {
	epuo.mutation.ClearDeletedAt()
	return epuo
}

// SetLastModifier sets the "last_modifier" field.
func (epuo *EnterprisePriceUpdateOne) SetLastModifier(m *model.Modifier) *EnterprisePriceUpdateOne {
	epuo.mutation.SetLastModifier(m)
	return epuo
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (epuo *EnterprisePriceUpdateOne) ClearLastModifier() *EnterprisePriceUpdateOne {
	epuo.mutation.ClearLastModifier()
	return epuo
}

// SetRemark sets the "remark" field.
func (epuo *EnterprisePriceUpdateOne) SetRemark(s string) *EnterprisePriceUpdateOne {
	epuo.mutation.SetRemark(s)
	return epuo
}

// SetNillableRemark sets the "remark" field if the given value is not nil.
func (epuo *EnterprisePriceUpdateOne) SetNillableRemark(s *string) *EnterprisePriceUpdateOne {
	if s != nil {
		epuo.SetRemark(*s)
	}
	return epuo
}

// ClearRemark clears the value of the "remark" field.
func (epuo *EnterprisePriceUpdateOne) ClearRemark() *EnterprisePriceUpdateOne {
	epuo.mutation.ClearRemark()
	return epuo
}

// SetCityID sets the "city_id" field.
func (epuo *EnterprisePriceUpdateOne) SetCityID(u uint64) *EnterprisePriceUpdateOne {
	epuo.mutation.SetCityID(u)
	return epuo
}

// SetEnterpriseID sets the "enterprise_id" field.
func (epuo *EnterprisePriceUpdateOne) SetEnterpriseID(u uint64) *EnterprisePriceUpdateOne {
	epuo.mutation.SetEnterpriseID(u)
	return epuo
}

// SetPrice sets the "price" field.
func (epuo *EnterprisePriceUpdateOne) SetPrice(f float64) *EnterprisePriceUpdateOne {
	epuo.mutation.ResetPrice()
	epuo.mutation.SetPrice(f)
	return epuo
}

// AddPrice adds f to the "price" field.
func (epuo *EnterprisePriceUpdateOne) AddPrice(f float64) *EnterprisePriceUpdateOne {
	epuo.mutation.AddPrice(f)
	return epuo
}

// SetVoltage sets the "voltage" field.
func (epuo *EnterprisePriceUpdateOne) SetVoltage(f float64) *EnterprisePriceUpdateOne {
	epuo.mutation.ResetVoltage()
	epuo.mutation.SetVoltage(f)
	return epuo
}

// AddVoltage adds f to the "voltage" field.
func (epuo *EnterprisePriceUpdateOne) AddVoltage(f float64) *EnterprisePriceUpdateOne {
	epuo.mutation.AddVoltage(f)
	return epuo
}

// SetCity sets the "city" edge to the City entity.
func (epuo *EnterprisePriceUpdateOne) SetCity(c *City) *EnterprisePriceUpdateOne {
	return epuo.SetCityID(c.ID)
}

// SetEnterprise sets the "enterprise" edge to the Enterprise entity.
func (epuo *EnterprisePriceUpdateOne) SetEnterprise(e *Enterprise) *EnterprisePriceUpdateOne {
	return epuo.SetEnterpriseID(e.ID)
}

// Mutation returns the EnterprisePriceMutation object of the builder.
func (epuo *EnterprisePriceUpdateOne) Mutation() *EnterprisePriceMutation {
	return epuo.mutation
}

// ClearCity clears the "city" edge to the City entity.
func (epuo *EnterprisePriceUpdateOne) ClearCity() *EnterprisePriceUpdateOne {
	epuo.mutation.ClearCity()
	return epuo
}

// ClearEnterprise clears the "enterprise" edge to the Enterprise entity.
func (epuo *EnterprisePriceUpdateOne) ClearEnterprise() *EnterprisePriceUpdateOne {
	epuo.mutation.ClearEnterprise()
	return epuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (epuo *EnterprisePriceUpdateOne) Select(field string, fields ...string) *EnterprisePriceUpdateOne {
	epuo.fields = append([]string{field}, fields...)
	return epuo
}

// Save executes the query and returns the updated EnterprisePrice entity.
func (epuo *EnterprisePriceUpdateOne) Save(ctx context.Context) (*EnterprisePrice, error) {
	var (
		err  error
		node *EnterprisePrice
	)
	if err := epuo.defaults(); err != nil {
		return nil, err
	}
	if len(epuo.hooks) == 0 {
		if err = epuo.check(); err != nil {
			return nil, err
		}
		node, err = epuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EnterprisePriceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = epuo.check(); err != nil {
				return nil, err
			}
			epuo.mutation = mutation
			node, err = epuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(epuo.hooks) - 1; i >= 0; i-- {
			if epuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = epuo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, epuo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*EnterprisePrice)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from EnterprisePriceMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (epuo *EnterprisePriceUpdateOne) SaveX(ctx context.Context) *EnterprisePrice {
	node, err := epuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (epuo *EnterprisePriceUpdateOne) Exec(ctx context.Context) error {
	_, err := epuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epuo *EnterprisePriceUpdateOne) ExecX(ctx context.Context) {
	if err := epuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (epuo *EnterprisePriceUpdateOne) defaults() error {
	if _, ok := epuo.mutation.UpdatedAt(); !ok {
		if enterpriseprice.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized enterpriseprice.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := enterpriseprice.UpdateDefaultUpdatedAt()
		epuo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (epuo *EnterprisePriceUpdateOne) check() error {
	if _, ok := epuo.mutation.CityID(); epuo.mutation.CityCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "EnterprisePrice.city"`)
	}
	if _, ok := epuo.mutation.EnterpriseID(); epuo.mutation.EnterpriseCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "EnterprisePrice.enterprise"`)
	}
	return nil
}

func (epuo *EnterprisePriceUpdateOne) sqlSave(ctx context.Context) (_node *EnterprisePrice, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   enterpriseprice.Table,
			Columns: enterpriseprice.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: enterpriseprice.FieldID,
			},
		},
	}
	id, ok := epuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "EnterprisePrice.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := epuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, enterpriseprice.FieldID)
		for _, f := range fields {
			if !enterpriseprice.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != enterpriseprice.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := epuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := epuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: enterpriseprice.FieldUpdatedAt,
		})
	}
	if value, ok := epuo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: enterpriseprice.FieldDeletedAt,
		})
	}
	if epuo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: enterpriseprice.FieldDeletedAt,
		})
	}
	if epuo.mutation.CreatorCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: enterpriseprice.FieldCreator,
		})
	}
	if value, ok := epuo.mutation.LastModifier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: enterpriseprice.FieldLastModifier,
		})
	}
	if epuo.mutation.LastModifierCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: enterpriseprice.FieldLastModifier,
		})
	}
	if value, ok := epuo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: enterpriseprice.FieldRemark,
		})
	}
	if epuo.mutation.RemarkCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: enterpriseprice.FieldRemark,
		})
	}
	if value, ok := epuo.mutation.Price(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: enterpriseprice.FieldPrice,
		})
	}
	if value, ok := epuo.mutation.AddedPrice(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: enterpriseprice.FieldPrice,
		})
	}
	if value, ok := epuo.mutation.Voltage(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: enterpriseprice.FieldVoltage,
		})
	}
	if value, ok := epuo.mutation.AddedVoltage(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: enterpriseprice.FieldVoltage,
		})
	}
	if epuo.mutation.CityCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   enterpriseprice.CityTable,
			Columns: []string{enterpriseprice.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: city.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epuo.mutation.CityIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   enterpriseprice.CityTable,
			Columns: []string{enterpriseprice.CityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: city.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epuo.mutation.EnterpriseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   enterpriseprice.EnterpriseTable,
			Columns: []string{enterpriseprice.EnterpriseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: enterprise.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epuo.mutation.EnterpriseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   enterpriseprice.EnterpriseTable,
			Columns: []string{enterpriseprice.EnterpriseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: enterprise.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &EnterprisePrice{config: epuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, epuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{enterpriseprice.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}