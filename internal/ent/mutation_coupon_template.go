// Code generated by liasica and entc, DO NOT EDIT.

package ent

import (
    "context"
    "fmt"
    "sync"
    "errors"
    "time"
    "github.com/auroraride/aurservd/internal/ent/coupontemplate"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/predicate"

    "entgo.io/ent"
)


// CouponTemplateMutation represents an operation that mutates the CouponTemplate nodes in the graph.
type CouponTemplateMutation struct {
	config
	op             Op
	typ            string
	id             *uint64
	created_at     *time.Time
	updated_at     *time.Time
	creator        **model.Modifier
	last_modifier  **model.Modifier
	remark         *string
	enable         *bool
	name           *string
	meta           **model.CouponTemplateMeta
	clearedFields  map[string]struct{}
	coupons        map[uint64]struct{}
	removedcoupons map[uint64]struct{}
	clearedcoupons bool
	done           bool
	oldValue       func(context.Context) (*CouponTemplate, error)
	predicates     []predicate.CouponTemplate
}

var _ ent.Mutation = (*CouponTemplateMutation)(nil)

// coupontemplateOption allows management of the mutation configuration using functional options.
type coupontemplateOption func(*CouponTemplateMutation)

// newCouponTemplateMutation creates new mutation for the CouponTemplate entity.
func newCouponTemplateMutation(c config, op Op, opts ...coupontemplateOption) *CouponTemplateMutation {
	m := &CouponTemplateMutation{
		config:        c,
		op:            op,
		typ:           TypeCouponTemplate,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withCouponTemplateID sets the ID field of the mutation.
func withCouponTemplateID(id uint64) coupontemplateOption {
	return func(m *CouponTemplateMutation) {
		var (
			err   error
			once  sync.Once
			value *CouponTemplate
		)
		m.oldValue = func(ctx context.Context) (*CouponTemplate, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().CouponTemplate.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withCouponTemplate sets the old CouponTemplate of the mutation.
func withCouponTemplate(node *CouponTemplate) coupontemplateOption {
	return func(m *CouponTemplateMutation) {
		m.oldValue = func(context.Context) (*CouponTemplate, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m CouponTemplateMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m CouponTemplateMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *CouponTemplateMutation) ID() (id uint64, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *CouponTemplateMutation) IDs(ctx context.Context) ([]uint64, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uint64{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().CouponTemplate.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreatedAt sets the "created_at" field.
func (m *CouponTemplateMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *CouponTemplateMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the CouponTemplate entity.
// If the CouponTemplate object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CouponTemplateMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
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
func (m *CouponTemplateMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the "updated_at" field.
func (m *CouponTemplateMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *CouponTemplateMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the CouponTemplate entity.
// If the CouponTemplate object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CouponTemplateMutation) OldUpdatedAt(ctx context.Context) (v time.Time, err error) {
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
func (m *CouponTemplateMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetCreator sets the "creator" field.
func (m *CouponTemplateMutation) SetCreator(value *model.Modifier) {
	m.creator = &value
}

// Creator returns the value of the "creator" field in the mutation.
func (m *CouponTemplateMutation) Creator() (r *model.Modifier, exists bool) {
	v := m.creator
	if v == nil {
		return
	}
	return *v, true
}

// OldCreator returns the old "creator" field's value of the CouponTemplate entity.
// If the CouponTemplate object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CouponTemplateMutation) OldCreator(ctx context.Context) (v *model.Modifier, err error) {
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
func (m *CouponTemplateMutation) ClearCreator() {
	m.creator = nil
	m.clearedFields[coupontemplate.FieldCreator] = struct{}{}
}

// CreatorCleared returns if the "creator" field was cleared in this mutation.
func (m *CouponTemplateMutation) CreatorCleared() bool {
	_, ok := m.clearedFields[coupontemplate.FieldCreator]
	return ok
}

// ResetCreator resets all changes to the "creator" field.
func (m *CouponTemplateMutation) ResetCreator() {
	m.creator = nil
	delete(m.clearedFields, coupontemplate.FieldCreator)
}

// SetLastModifier sets the "last_modifier" field.
func (m *CouponTemplateMutation) SetLastModifier(value *model.Modifier) {
	m.last_modifier = &value
}

// LastModifier returns the value of the "last_modifier" field in the mutation.
func (m *CouponTemplateMutation) LastModifier() (r *model.Modifier, exists bool) {
	v := m.last_modifier
	if v == nil {
		return
	}
	return *v, true
}

// OldLastModifier returns the old "last_modifier" field's value of the CouponTemplate entity.
// If the CouponTemplate object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CouponTemplateMutation) OldLastModifier(ctx context.Context) (v *model.Modifier, err error) {
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
func (m *CouponTemplateMutation) ClearLastModifier() {
	m.last_modifier = nil
	m.clearedFields[coupontemplate.FieldLastModifier] = struct{}{}
}

// LastModifierCleared returns if the "last_modifier" field was cleared in this mutation.
func (m *CouponTemplateMutation) LastModifierCleared() bool {
	_, ok := m.clearedFields[coupontemplate.FieldLastModifier]
	return ok
}

// ResetLastModifier resets all changes to the "last_modifier" field.
func (m *CouponTemplateMutation) ResetLastModifier() {
	m.last_modifier = nil
	delete(m.clearedFields, coupontemplate.FieldLastModifier)
}

// SetRemark sets the "remark" field.
func (m *CouponTemplateMutation) SetRemark(s string) {
	m.remark = &s
}

// Remark returns the value of the "remark" field in the mutation.
func (m *CouponTemplateMutation) Remark() (r string, exists bool) {
	v := m.remark
	if v == nil {
		return
	}
	return *v, true
}

// OldRemark returns the old "remark" field's value of the CouponTemplate entity.
// If the CouponTemplate object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CouponTemplateMutation) OldRemark(ctx context.Context) (v string, err error) {
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
func (m *CouponTemplateMutation) ClearRemark() {
	m.remark = nil
	m.clearedFields[coupontemplate.FieldRemark] = struct{}{}
}

// RemarkCleared returns if the "remark" field was cleared in this mutation.
func (m *CouponTemplateMutation) RemarkCleared() bool {
	_, ok := m.clearedFields[coupontemplate.FieldRemark]
	return ok
}

// ResetRemark resets all changes to the "remark" field.
func (m *CouponTemplateMutation) ResetRemark() {
	m.remark = nil
	delete(m.clearedFields, coupontemplate.FieldRemark)
}

// SetEnable sets the "enable" field.
func (m *CouponTemplateMutation) SetEnable(b bool) {
	m.enable = &b
}

// Enable returns the value of the "enable" field in the mutation.
func (m *CouponTemplateMutation) Enable() (r bool, exists bool) {
	v := m.enable
	if v == nil {
		return
	}
	return *v, true
}

// OldEnable returns the old "enable" field's value of the CouponTemplate entity.
// If the CouponTemplate object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CouponTemplateMutation) OldEnable(ctx context.Context) (v bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEnable is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEnable requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEnable: %w", err)
	}
	return oldValue.Enable, nil
}

// ResetEnable resets all changes to the "enable" field.
func (m *CouponTemplateMutation) ResetEnable() {
	m.enable = nil
}

// SetName sets the "name" field.
func (m *CouponTemplateMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *CouponTemplateMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the CouponTemplate entity.
// If the CouponTemplate object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CouponTemplateMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *CouponTemplateMutation) ResetName() {
	m.name = nil
}

// SetMeta sets the "meta" field.
func (m *CouponTemplateMutation) SetMeta(mtm *model.CouponTemplateMeta) {
	m.meta = &mtm
}

// Meta returns the value of the "meta" field in the mutation.
func (m *CouponTemplateMutation) Meta() (r *model.CouponTemplateMeta, exists bool) {
	v := m.meta
	if v == nil {
		return
	}
	return *v, true
}

// OldMeta returns the old "meta" field's value of the CouponTemplate entity.
// If the CouponTemplate object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CouponTemplateMutation) OldMeta(ctx context.Context) (v *model.CouponTemplateMeta, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldMeta is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldMeta requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldMeta: %w", err)
	}
	return oldValue.Meta, nil
}

// ResetMeta resets all changes to the "meta" field.
func (m *CouponTemplateMutation) ResetMeta() {
	m.meta = nil
}

// AddCouponIDs adds the "coupons" edge to the Coupon entity by ids.
func (m *CouponTemplateMutation) AddCouponIDs(ids ...uint64) {
	if m.coupons == nil {
		m.coupons = make(map[uint64]struct{})
	}
	for i := range ids {
		m.coupons[ids[i]] = struct{}{}
	}
}

// ClearCoupons clears the "coupons" edge to the Coupon entity.
func (m *CouponTemplateMutation) ClearCoupons() {
	m.clearedcoupons = true
}

// CouponsCleared reports if the "coupons" edge to the Coupon entity was cleared.
func (m *CouponTemplateMutation) CouponsCleared() bool {
	return m.clearedcoupons
}

// RemoveCouponIDs removes the "coupons" edge to the Coupon entity by IDs.
func (m *CouponTemplateMutation) RemoveCouponIDs(ids ...uint64) {
	if m.removedcoupons == nil {
		m.removedcoupons = make(map[uint64]struct{})
	}
	for i := range ids {
		delete(m.coupons, ids[i])
		m.removedcoupons[ids[i]] = struct{}{}
	}
}

// RemovedCoupons returns the removed IDs of the "coupons" edge to the Coupon entity.
func (m *CouponTemplateMutation) RemovedCouponsIDs() (ids []uint64) {
	for id := range m.removedcoupons {
		ids = append(ids, id)
	}
	return
}

// CouponsIDs returns the "coupons" edge IDs in the mutation.
func (m *CouponTemplateMutation) CouponsIDs() (ids []uint64) {
	for id := range m.coupons {
		ids = append(ids, id)
	}
	return
}

// ResetCoupons resets all changes to the "coupons" edge.
func (m *CouponTemplateMutation) ResetCoupons() {
	m.coupons = nil
	m.clearedcoupons = false
	m.removedcoupons = nil
}

// Where appends a list predicates to the CouponTemplateMutation builder.
func (m *CouponTemplateMutation) Where(ps ...predicate.CouponTemplate) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *CouponTemplateMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (CouponTemplate).
func (m *CouponTemplateMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *CouponTemplateMutation) Fields() []string {
	fields := make([]string, 0, 8)
	if m.created_at != nil {
		fields = append(fields, coupontemplate.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, coupontemplate.FieldUpdatedAt)
	}
	if m.creator != nil {
		fields = append(fields, coupontemplate.FieldCreator)
	}
	if m.last_modifier != nil {
		fields = append(fields, coupontemplate.FieldLastModifier)
	}
	if m.remark != nil {
		fields = append(fields, coupontemplate.FieldRemark)
	}
	if m.enable != nil {
		fields = append(fields, coupontemplate.FieldEnable)
	}
	if m.name != nil {
		fields = append(fields, coupontemplate.FieldName)
	}
	if m.meta != nil {
		fields = append(fields, coupontemplate.FieldMeta)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *CouponTemplateMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case coupontemplate.FieldCreatedAt:
		return m.CreatedAt()
	case coupontemplate.FieldUpdatedAt:
		return m.UpdatedAt()
	case coupontemplate.FieldCreator:
		return m.Creator()
	case coupontemplate.FieldLastModifier:
		return m.LastModifier()
	case coupontemplate.FieldRemark:
		return m.Remark()
	case coupontemplate.FieldEnable:
		return m.Enable()
	case coupontemplate.FieldName:
		return m.Name()
	case coupontemplate.FieldMeta:
		return m.Meta()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *CouponTemplateMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case coupontemplate.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case coupontemplate.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	case coupontemplate.FieldCreator:
		return m.OldCreator(ctx)
	case coupontemplate.FieldLastModifier:
		return m.OldLastModifier(ctx)
	case coupontemplate.FieldRemark:
		return m.OldRemark(ctx)
	case coupontemplate.FieldEnable:
		return m.OldEnable(ctx)
	case coupontemplate.FieldName:
		return m.OldName(ctx)
	case coupontemplate.FieldMeta:
		return m.OldMeta(ctx)
	}
	return nil, fmt.Errorf("unknown CouponTemplate field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *CouponTemplateMutation) SetField(name string, value ent.Value) error {
	switch name {
	case coupontemplate.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case coupontemplate.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case coupontemplate.FieldCreator:
		v, ok := value.(*model.Modifier)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreator(v)
		return nil
	case coupontemplate.FieldLastModifier:
		v, ok := value.(*model.Modifier)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLastModifier(v)
		return nil
	case coupontemplate.FieldRemark:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRemark(v)
		return nil
	case coupontemplate.FieldEnable:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnable(v)
		return nil
	case coupontemplate.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case coupontemplate.FieldMeta:
		v, ok := value.(*model.CouponTemplateMeta)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMeta(v)
		return nil
	}
	return fmt.Errorf("unknown CouponTemplate field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *CouponTemplateMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *CouponTemplateMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *CouponTemplateMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown CouponTemplate numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *CouponTemplateMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(coupontemplate.FieldCreator) {
		fields = append(fields, coupontemplate.FieldCreator)
	}
	if m.FieldCleared(coupontemplate.FieldLastModifier) {
		fields = append(fields, coupontemplate.FieldLastModifier)
	}
	if m.FieldCleared(coupontemplate.FieldRemark) {
		fields = append(fields, coupontemplate.FieldRemark)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *CouponTemplateMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *CouponTemplateMutation) ClearField(name string) error {
	switch name {
	case coupontemplate.FieldCreator:
		m.ClearCreator()
		return nil
	case coupontemplate.FieldLastModifier:
		m.ClearLastModifier()
		return nil
	case coupontemplate.FieldRemark:
		m.ClearRemark()
		return nil
	}
	return fmt.Errorf("unknown CouponTemplate nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *CouponTemplateMutation) ResetField(name string) error {
	switch name {
	case coupontemplate.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case coupontemplate.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case coupontemplate.FieldCreator:
		m.ResetCreator()
		return nil
	case coupontemplate.FieldLastModifier:
		m.ResetLastModifier()
		return nil
	case coupontemplate.FieldRemark:
		m.ResetRemark()
		return nil
	case coupontemplate.FieldEnable:
		m.ResetEnable()
		return nil
	case coupontemplate.FieldName:
		m.ResetName()
		return nil
	case coupontemplate.FieldMeta:
		m.ResetMeta()
		return nil
	}
	return fmt.Errorf("unknown CouponTemplate field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *CouponTemplateMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.coupons != nil {
		edges = append(edges, coupontemplate.EdgeCoupons)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *CouponTemplateMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case coupontemplate.EdgeCoupons:
		ids := make([]ent.Value, 0, len(m.coupons))
		for id := range m.coupons {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *CouponTemplateMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedcoupons != nil {
		edges = append(edges, coupontemplate.EdgeCoupons)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *CouponTemplateMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case coupontemplate.EdgeCoupons:
		ids := make([]ent.Value, 0, len(m.removedcoupons))
		for id := range m.removedcoupons {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *CouponTemplateMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedcoupons {
		edges = append(edges, coupontemplate.EdgeCoupons)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *CouponTemplateMutation) EdgeCleared(name string) bool {
	switch name {
	case coupontemplate.EdgeCoupons:
		return m.clearedcoupons
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *CouponTemplateMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown CouponTemplate unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *CouponTemplateMutation) ResetEdge(name string) error {
	switch name {
	case coupontemplate.EdgeCoupons:
		m.ResetCoupons()
		return nil
	}
	return fmt.Errorf("unknown CouponTemplate edge %s", name)
}
