// Code generated by liasica and entc, DO NOT EDIT.

package ent

import (
    "context"
    "fmt"
    "sync"
    "errors"
    "time"
    "github.com/auroraride/aurservd/internal/ent/commission"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/predicate"

    "entgo.io/ent"
)


// CommissionMutation represents an operation that mutates the Commission nodes in the graph.
type CommissionMutation struct {
	config
	op              Op
	typ             string
	id              *uint64
	created_at      *time.Time
	updated_at      *time.Time
	deleted_at      *time.Time
	creator         **model.Modifier
	last_modifier   **model.Modifier
	remark          *string
	amount          *float64
	addamount       *float64
	status          *uint8
	addstatus       *int8
	clearedFields   map[string]struct{}
	_order          *uint64
	cleared_order   bool
	employee        *uint64
	clearedemployee bool
	done            bool
	oldValue        func(context.Context) (*Commission, error)
	predicates      []predicate.Commission
}

var _ ent.Mutation = (*CommissionMutation)(nil)

// commissionOption allows management of the mutation configuration using functional options.
type commissionOption func(*CommissionMutation)

// newCommissionMutation creates new mutation for the Commission entity.
func newCommissionMutation(c config, op Op, opts ...commissionOption) *CommissionMutation {
	m := &CommissionMutation{
		config:        c,
		op:            op,
		typ:           TypeCommission,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withCommissionID sets the ID field of the mutation.
func withCommissionID(id uint64) commissionOption {
	return func(m *CommissionMutation) {
		var (
			err   error
			once  sync.Once
			value *Commission
		)
		m.oldValue = func(ctx context.Context) (*Commission, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Commission.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withCommission sets the old Commission of the mutation.
func withCommission(node *Commission) commissionOption {
	return func(m *CommissionMutation) {
		m.oldValue = func(context.Context) (*Commission, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m CommissionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m CommissionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *CommissionMutation) ID() (id uint64, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *CommissionMutation) IDs(ctx context.Context) ([]uint64, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uint64{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Commission.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreatedAt sets the "created_at" field.
func (m *CommissionMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *CommissionMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the Commission entity.
// If the Commission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CommissionMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *CommissionMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the "updated_at" field.
func (m *CommissionMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *CommissionMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the Commission entity.
// If the Commission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CommissionMutation) OldUpdatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdatedAt: %w", err)
	}
	return oldValue.UpdatedAt, nil
}

// ResetUpdatedAt resets all changes to the "updated_at" field.
func (m *CommissionMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetDeletedAt sets the "deleted_at" field.
func (m *CommissionMutation) SetDeletedAt(t time.Time) {
	m.deleted_at = &t
}

// DeletedAt returns the value of the "deleted_at" field in the mutation.
func (m *CommissionMutation) DeletedAt() (r time.Time, exists bool) {
	v := m.deleted_at
	if v == nil {
		return
	}
	return *v, true
}

// OldDeletedAt returns the old "deleted_at" field's value of the Commission entity.
// If the Commission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CommissionMutation) OldDeletedAt(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDeletedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDeletedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDeletedAt: %w", err)
	}
	return oldValue.DeletedAt, nil
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (m *CommissionMutation) ClearDeletedAt() {
	m.deleted_at = nil
	m.clearedFields[commission.FieldDeletedAt] = struct{}{}
}

// DeletedAtCleared returns if the "deleted_at" field was cleared in this mutation.
func (m *CommissionMutation) DeletedAtCleared() bool {
	_, ok := m.clearedFields[commission.FieldDeletedAt]
	return ok
}

// ResetDeletedAt resets all changes to the "deleted_at" field.
func (m *CommissionMutation) ResetDeletedAt() {
	m.deleted_at = nil
	delete(m.clearedFields, commission.FieldDeletedAt)
}

// SetCreator sets the "creator" field.
func (m *CommissionMutation) SetCreator(value *model.Modifier) {
	m.creator = &value
}

// Creator returns the value of the "creator" field in the mutation.
func (m *CommissionMutation) Creator() (r *model.Modifier, exists bool) {
	v := m.creator
	if v == nil {
		return
	}
	return *v, true
}

// OldCreator returns the old "creator" field's value of the Commission entity.
// If the Commission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CommissionMutation) OldCreator(ctx context.Context) (v *model.Modifier, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreator is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreator requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreator: %w", err)
	}
	return oldValue.Creator, nil
}

// ClearCreator clears the value of the "creator" field.
func (m *CommissionMutation) ClearCreator() {
	m.creator = nil
	m.clearedFields[commission.FieldCreator] = struct{}{}
}

// CreatorCleared returns if the "creator" field was cleared in this mutation.
func (m *CommissionMutation) CreatorCleared() bool {
	_, ok := m.clearedFields[commission.FieldCreator]
	return ok
}

// ResetCreator resets all changes to the "creator" field.
func (m *CommissionMutation) ResetCreator() {
	m.creator = nil
	delete(m.clearedFields, commission.FieldCreator)
}

// SetLastModifier sets the "last_modifier" field.
func (m *CommissionMutation) SetLastModifier(value *model.Modifier) {
	m.last_modifier = &value
}

// LastModifier returns the value of the "last_modifier" field in the mutation.
func (m *CommissionMutation) LastModifier() (r *model.Modifier, exists bool) {
	v := m.last_modifier
	if v == nil {
		return
	}
	return *v, true
}

// OldLastModifier returns the old "last_modifier" field's value of the Commission entity.
// If the Commission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CommissionMutation) OldLastModifier(ctx context.Context) (v *model.Modifier, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldLastModifier is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldLastModifier requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldLastModifier: %w", err)
	}
	return oldValue.LastModifier, nil
}

// ClearLastModifier clears the value of the "last_modifier" field.
func (m *CommissionMutation) ClearLastModifier() {
	m.last_modifier = nil
	m.clearedFields[commission.FieldLastModifier] = struct{}{}
}

// LastModifierCleared returns if the "last_modifier" field was cleared in this mutation.
func (m *CommissionMutation) LastModifierCleared() bool {
	_, ok := m.clearedFields[commission.FieldLastModifier]
	return ok
}

// ResetLastModifier resets all changes to the "last_modifier" field.
func (m *CommissionMutation) ResetLastModifier() {
	m.last_modifier = nil
	delete(m.clearedFields, commission.FieldLastModifier)
}

// SetRemark sets the "remark" field.
func (m *CommissionMutation) SetRemark(s string) {
	m.remark = &s
}

// Remark returns the value of the "remark" field in the mutation.
func (m *CommissionMutation) Remark() (r string, exists bool) {
	v := m.remark
	if v == nil {
		return
	}
	return *v, true
}

// OldRemark returns the old "remark" field's value of the Commission entity.
// If the Commission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CommissionMutation) OldRemark(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRemark is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRemark requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRemark: %w", err)
	}
	return oldValue.Remark, nil
}

// ClearRemark clears the value of the "remark" field.
func (m *CommissionMutation) ClearRemark() {
	m.remark = nil
	m.clearedFields[commission.FieldRemark] = struct{}{}
}

// RemarkCleared returns if the "remark" field was cleared in this mutation.
func (m *CommissionMutation) RemarkCleared() bool {
	_, ok := m.clearedFields[commission.FieldRemark]
	return ok
}

// ResetRemark resets all changes to the "remark" field.
func (m *CommissionMutation) ResetRemark() {
	m.remark = nil
	delete(m.clearedFields, commission.FieldRemark)
}

// SetOrderID sets the "order_id" field.
func (m *CommissionMutation) SetOrderID(u uint64) {
	m._order = &u
}

// OrderID returns the value of the "order_id" field in the mutation.
func (m *CommissionMutation) OrderID() (r uint64, exists bool) {
	v := m._order
	if v == nil {
		return
	}
	return *v, true
}

// OldOrderID returns the old "order_id" field's value of the Commission entity.
// If the Commission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CommissionMutation) OldOrderID(ctx context.Context) (v uint64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldOrderID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldOrderID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldOrderID: %w", err)
	}
	return oldValue.OrderID, nil
}

// ResetOrderID resets all changes to the "order_id" field.
func (m *CommissionMutation) ResetOrderID() {
	m._order = nil
}

// SetAmount sets the "amount" field.
func (m *CommissionMutation) SetAmount(f float64) {
	m.amount = &f
	m.addamount = nil
}

// Amount returns the value of the "amount" field in the mutation.
func (m *CommissionMutation) Amount() (r float64, exists bool) {
	v := m.amount
	if v == nil {
		return
	}
	return *v, true
}

// OldAmount returns the old "amount" field's value of the Commission entity.
// If the Commission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CommissionMutation) OldAmount(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldAmount is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldAmount requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldAmount: %w", err)
	}
	return oldValue.Amount, nil
}

// AddAmount adds f to the "amount" field.
func (m *CommissionMutation) AddAmount(f float64) {
	if m.addamount != nil {
		*m.addamount += f
	} else {
		m.addamount = &f
	}
}

// AddedAmount returns the value that was added to the "amount" field in this mutation.
func (m *CommissionMutation) AddedAmount() (r float64, exists bool) {
	v := m.addamount
	if v == nil {
		return
	}
	return *v, true
}

// ResetAmount resets all changes to the "amount" field.
func (m *CommissionMutation) ResetAmount() {
	m.amount = nil
	m.addamount = nil
}

// SetStatus sets the "status" field.
func (m *CommissionMutation) SetStatus(u uint8) {
	m.status = &u
	m.addstatus = nil
}

// Status returns the value of the "status" field in the mutation.
func (m *CommissionMutation) Status() (r uint8, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// OldStatus returns the old "status" field's value of the Commission entity.
// If the Commission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CommissionMutation) OldStatus(ctx context.Context) (v uint8, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatus is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatus requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatus: %w", err)
	}
	return oldValue.Status, nil
}

// AddStatus adds u to the "status" field.
func (m *CommissionMutation) AddStatus(u int8) {
	if m.addstatus != nil {
		*m.addstatus += u
	} else {
		m.addstatus = &u
	}
}

// AddedStatus returns the value that was added to the "status" field in this mutation.
func (m *CommissionMutation) AddedStatus() (r int8, exists bool) {
	v := m.addstatus
	if v == nil {
		return
	}
	return *v, true
}

// ResetStatus resets all changes to the "status" field.
func (m *CommissionMutation) ResetStatus() {
	m.status = nil
	m.addstatus = nil
}

// SetEmployeeID sets the "employee_id" field.
func (m *CommissionMutation) SetEmployeeID(u uint64) {
	m.employee = &u
}

// EmployeeID returns the value of the "employee_id" field in the mutation.
func (m *CommissionMutation) EmployeeID() (r uint64, exists bool) {
	v := m.employee
	if v == nil {
		return
	}
	return *v, true
}

// OldEmployeeID returns the old "employee_id" field's value of the Commission entity.
// If the Commission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CommissionMutation) OldEmployeeID(ctx context.Context) (v *uint64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEmployeeID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEmployeeID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEmployeeID: %w", err)
	}
	return oldValue.EmployeeID, nil
}

// ClearEmployeeID clears the value of the "employee_id" field.
func (m *CommissionMutation) ClearEmployeeID() {
	m.employee = nil
	m.clearedFields[commission.FieldEmployeeID] = struct{}{}
}

// EmployeeIDCleared returns if the "employee_id" field was cleared in this mutation.
func (m *CommissionMutation) EmployeeIDCleared() bool {
	_, ok := m.clearedFields[commission.FieldEmployeeID]
	return ok
}

// ResetEmployeeID resets all changes to the "employee_id" field.
func (m *CommissionMutation) ResetEmployeeID() {
	m.employee = nil
	delete(m.clearedFields, commission.FieldEmployeeID)
}

// ClearOrder clears the "order" edge to the Order entity.
func (m *CommissionMutation) ClearOrder() {
	m.cleared_order = true
}

// OrderCleared reports if the "order" edge to the Order entity was cleared.
func (m *CommissionMutation) OrderCleared() bool {
	return m.cleared_order
}

// OrderIDs returns the "order" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// OrderID instead. It exists only for internal usage by the builders.
func (m *CommissionMutation) OrderIDs() (ids []uint64) {
	if id := m._order; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetOrder resets all changes to the "order" edge.
func (m *CommissionMutation) ResetOrder() {
	m._order = nil
	m.cleared_order = false
}

// ClearEmployee clears the "employee" edge to the Employee entity.
func (m *CommissionMutation) ClearEmployee() {
	m.clearedemployee = true
}

// EmployeeCleared reports if the "employee" edge to the Employee entity was cleared.
func (m *CommissionMutation) EmployeeCleared() bool {
	return m.EmployeeIDCleared() || m.clearedemployee
}

// EmployeeIDs returns the "employee" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// EmployeeID instead. It exists only for internal usage by the builders.
func (m *CommissionMutation) EmployeeIDs() (ids []uint64) {
	if id := m.employee; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEmployee resets all changes to the "employee" edge.
func (m *CommissionMutation) ResetEmployee() {
	m.employee = nil
	m.clearedemployee = false
}

// Where appends a list predicates to the CommissionMutation builder.
func (m *CommissionMutation) Where(ps ...predicate.Commission) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *CommissionMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Commission).
func (m *CommissionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *CommissionMutation) Fields() []string {
	fields := make([]string, 0, 10)
	if m.created_at != nil {
		fields = append(fields, commission.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, commission.FieldUpdatedAt)
	}
	if m.deleted_at != nil {
		fields = append(fields, commission.FieldDeletedAt)
	}
	if m.creator != nil {
		fields = append(fields, commission.FieldCreator)
	}
	if m.last_modifier != nil {
		fields = append(fields, commission.FieldLastModifier)
	}
	if m.remark != nil {
		fields = append(fields, commission.FieldRemark)
	}
	if m._order != nil {
		fields = append(fields, commission.FieldOrderID)
	}
	if m.amount != nil {
		fields = append(fields, commission.FieldAmount)
	}
	if m.status != nil {
		fields = append(fields, commission.FieldStatus)
	}
	if m.employee != nil {
		fields = append(fields, commission.FieldEmployeeID)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *CommissionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case commission.FieldCreatedAt:
		return m.CreatedAt()
	case commission.FieldUpdatedAt:
		return m.UpdatedAt()
	case commission.FieldDeletedAt:
		return m.DeletedAt()
	case commission.FieldCreator:
		return m.Creator()
	case commission.FieldLastModifier:
		return m.LastModifier()
	case commission.FieldRemark:
		return m.Remark()
	case commission.FieldOrderID:
		return m.OrderID()
	case commission.FieldAmount:
		return m.Amount()
	case commission.FieldStatus:
		return m.Status()
	case commission.FieldEmployeeID:
		return m.EmployeeID()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *CommissionMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case commission.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case commission.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	case commission.FieldDeletedAt:
		return m.OldDeletedAt(ctx)
	case commission.FieldCreator:
		return m.OldCreator(ctx)
	case commission.FieldLastModifier:
		return m.OldLastModifier(ctx)
	case commission.FieldRemark:
		return m.OldRemark(ctx)
	case commission.FieldOrderID:
		return m.OldOrderID(ctx)
	case commission.FieldAmount:
		return m.OldAmount(ctx)
	case commission.FieldStatus:
		return m.OldStatus(ctx)
	case commission.FieldEmployeeID:
		return m.OldEmployeeID(ctx)
	}
	return nil, fmt.Errorf("unknown Commission field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *CommissionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case commission.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case commission.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case commission.FieldDeletedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDeletedAt(v)
		return nil
	case commission.FieldCreator:
		v, ok := value.(*model.Modifier)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreator(v)
		return nil
	case commission.FieldLastModifier:
		v, ok := value.(*model.Modifier)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLastModifier(v)
		return nil
	case commission.FieldRemark:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRemark(v)
		return nil
	case commission.FieldOrderID:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetOrderID(v)
		return nil
	case commission.FieldAmount:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetAmount(v)
		return nil
	case commission.FieldStatus:
		v, ok := value.(uint8)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case commission.FieldEmployeeID:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEmployeeID(v)
		return nil
	}
	return fmt.Errorf("unknown Commission field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *CommissionMutation) AddedFields() []string {
	var fields []string
	if m.addamount != nil {
		fields = append(fields, commission.FieldAmount)
	}
	if m.addstatus != nil {
		fields = append(fields, commission.FieldStatus)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *CommissionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case commission.FieldAmount:
		return m.AddedAmount()
	case commission.FieldStatus:
		return m.AddedStatus()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *CommissionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case commission.FieldAmount:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddAmount(v)
		return nil
	case commission.FieldStatus:
		v, ok := value.(int8)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddStatus(v)
		return nil
	}
	return fmt.Errorf("unknown Commission numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *CommissionMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(commission.FieldDeletedAt) {
		fields = append(fields, commission.FieldDeletedAt)
	}
	if m.FieldCleared(commission.FieldCreator) {
		fields = append(fields, commission.FieldCreator)
	}
	if m.FieldCleared(commission.FieldLastModifier) {
		fields = append(fields, commission.FieldLastModifier)
	}
	if m.FieldCleared(commission.FieldRemark) {
		fields = append(fields, commission.FieldRemark)
	}
	if m.FieldCleared(commission.FieldEmployeeID) {
		fields = append(fields, commission.FieldEmployeeID)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *CommissionMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *CommissionMutation) ClearField(name string) error {
	switch name {
	case commission.FieldDeletedAt:
		m.ClearDeletedAt()
		return nil
	case commission.FieldCreator:
		m.ClearCreator()
		return nil
	case commission.FieldLastModifier:
		m.ClearLastModifier()
		return nil
	case commission.FieldRemark:
		m.ClearRemark()
		return nil
	case commission.FieldEmployeeID:
		m.ClearEmployeeID()
		return nil
	}
	return fmt.Errorf("unknown Commission nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *CommissionMutation) ResetField(name string) error {
	switch name {
	case commission.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case commission.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case commission.FieldDeletedAt:
		m.ResetDeletedAt()
		return nil
	case commission.FieldCreator:
		m.ResetCreator()
		return nil
	case commission.FieldLastModifier:
		m.ResetLastModifier()
		return nil
	case commission.FieldRemark:
		m.ResetRemark()
		return nil
	case commission.FieldOrderID:
		m.ResetOrderID()
		return nil
	case commission.FieldAmount:
		m.ResetAmount()
		return nil
	case commission.FieldStatus:
		m.ResetStatus()
		return nil
	case commission.FieldEmployeeID:
		m.ResetEmployeeID()
		return nil
	}
	return fmt.Errorf("unknown Commission field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *CommissionMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m._order != nil {
		edges = append(edges, commission.EdgeOrder)
	}
	if m.employee != nil {
		edges = append(edges, commission.EdgeEmployee)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *CommissionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case commission.EdgeOrder:
		if id := m._order; id != nil {
			return []ent.Value{*id}
		}
	case commission.EdgeEmployee:
		if id := m.employee; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *CommissionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *CommissionMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *CommissionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.cleared_order {
		edges = append(edges, commission.EdgeOrder)
	}
	if m.clearedemployee {
		edges = append(edges, commission.EdgeEmployee)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *CommissionMutation) EdgeCleared(name string) bool {
	switch name {
	case commission.EdgeOrder:
		return m.cleared_order
	case commission.EdgeEmployee:
		return m.clearedemployee
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *CommissionMutation) ClearEdge(name string) error {
	switch name {
	case commission.EdgeOrder:
		m.ClearOrder()
		return nil
	case commission.EdgeEmployee:
		m.ClearEmployee()
		return nil
	}
	return fmt.Errorf("unknown Commission unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *CommissionMutation) ResetEdge(name string) error {
	switch name {
	case commission.EdgeOrder:
		m.ResetOrder()
		return nil
	case commission.EdgeEmployee:
		m.ResetEmployee()
		return nil
	}
	return fmt.Errorf("unknown Commission edge %s", name)
}
