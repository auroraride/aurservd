// Code generated by liasica and entc, DO NOT EDIT.

package ent

import (
    "context"
    "fmt"
    "sync"
    "errors"
    "time"
    "github.com/auroraride/aurservd/internal/ent/planintroduce"
    "github.com/auroraride/aurservd/internal/ent/predicate"

    "entgo.io/ent"
)


// PlanIntroduceMutation represents an operation that mutates the PlanIntroduce nodes in the graph.
type PlanIntroduceMutation struct {
	config
	op            Op
	typ           string
	id            *uint64
	created_at    *time.Time
	updated_at    *time.Time
	model         *string
	image         *string
	clearedFields map[string]struct{}
	brand         *uint64
	clearedbrand  bool
	done          bool
	oldValue      func(context.Context) (*PlanIntroduce, error)
	predicates    []predicate.PlanIntroduce
}

var _ ent.Mutation = (*PlanIntroduceMutation)(nil)

// planintroduceOption allows management of the mutation configuration using functional options.
type planintroduceOption func(*PlanIntroduceMutation)

// newPlanIntroduceMutation creates new mutation for the PlanIntroduce entity.
func newPlanIntroduceMutation(c config, op Op, opts ...planintroduceOption) *PlanIntroduceMutation {
	m := &PlanIntroduceMutation{
		config:        c,
		op:            op,
		typ:           TypePlanIntroduce,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withPlanIntroduceID sets the ID field of the mutation.
func withPlanIntroduceID(id uint64) planintroduceOption {
	return func(m *PlanIntroduceMutation) {
		var (
			err   error
			once  sync.Once
			value *PlanIntroduce
		)
		m.oldValue = func(ctx context.Context) (*PlanIntroduce, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().PlanIntroduce.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withPlanIntroduce sets the old PlanIntroduce of the mutation.
func withPlanIntroduce(node *PlanIntroduce) planintroduceOption {
	return func(m *PlanIntroduceMutation) {
		m.oldValue = func(context.Context) (*PlanIntroduce, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m PlanIntroduceMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m PlanIntroduceMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *PlanIntroduceMutation) ID() (id uint64, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *PlanIntroduceMutation) IDs(ctx context.Context) ([]uint64, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uint64{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().PlanIntroduce.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreatedAt sets the "created_at" field.
func (m *PlanIntroduceMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *PlanIntroduceMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the PlanIntroduce entity.
// If the PlanIntroduce object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PlanIntroduceMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
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
func (m *PlanIntroduceMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the "updated_at" field.
func (m *PlanIntroduceMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *PlanIntroduceMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the PlanIntroduce entity.
// If the PlanIntroduce object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PlanIntroduceMutation) OldUpdatedAt(ctx context.Context) (v time.Time, err error) {
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
func (m *PlanIntroduceMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetBrandID sets the "brand_id" field.
func (m *PlanIntroduceMutation) SetBrandID(u uint64) {
	m.brand = &u
}

// BrandID returns the value of the "brand_id" field in the mutation.
func (m *PlanIntroduceMutation) BrandID() (r uint64, exists bool) {
	v := m.brand
	if v == nil {
		return
	}
	return *v, true
}

// OldBrandID returns the old "brand_id" field's value of the PlanIntroduce entity.
// If the PlanIntroduce object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PlanIntroduceMutation) OldBrandID(ctx context.Context) (v *uint64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldBrandID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldBrandID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldBrandID: %w", err)
	}
	return oldValue.BrandID, nil
}

// ClearBrandID clears the value of the "brand_id" field.
func (m *PlanIntroduceMutation) ClearBrandID() {
	m.brand = nil
	m.clearedFields[planintroduce.FieldBrandID] = struct{}{}
}

// BrandIDCleared returns if the "brand_id" field was cleared in this mutation.
func (m *PlanIntroduceMutation) BrandIDCleared() bool {
	_, ok := m.clearedFields[planintroduce.FieldBrandID]
	return ok
}

// ResetBrandID resets all changes to the "brand_id" field.
func (m *PlanIntroduceMutation) ResetBrandID() {
	m.brand = nil
	delete(m.clearedFields, planintroduce.FieldBrandID)
}

// SetModel sets the "model" field.
func (m *PlanIntroduceMutation) SetModel(s string) {
	m.model = &s
}

// Model returns the value of the "model" field in the mutation.
func (m *PlanIntroduceMutation) Model() (r string, exists bool) {
	v := m.model
	if v == nil {
		return
	}
	return *v, true
}

// OldModel returns the old "model" field's value of the PlanIntroduce entity.
// If the PlanIntroduce object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PlanIntroduceMutation) OldModel(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldModel is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldModel requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldModel: %w", err)
	}
	return oldValue.Model, nil
}

// ResetModel resets all changes to the "model" field.
func (m *PlanIntroduceMutation) ResetModel() {
	m.model = nil
}

// SetImage sets the "image" field.
func (m *PlanIntroduceMutation) SetImage(s string) {
	m.image = &s
}

// Image returns the value of the "image" field in the mutation.
func (m *PlanIntroduceMutation) Image() (r string, exists bool) {
	v := m.image
	if v == nil {
		return
	}
	return *v, true
}

// OldImage returns the old "image" field's value of the PlanIntroduce entity.
// If the PlanIntroduce object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PlanIntroduceMutation) OldImage(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldImage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldImage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldImage: %w", err)
	}
	return oldValue.Image, nil
}

// ResetImage resets all changes to the "image" field.
func (m *PlanIntroduceMutation) ResetImage() {
	m.image = nil
}

// ClearBrand clears the "brand" edge to the EbikeBrand entity.
func (m *PlanIntroduceMutation) ClearBrand() {
	m.clearedbrand = true
}

// BrandCleared reports if the "brand" edge to the EbikeBrand entity was cleared.
func (m *PlanIntroduceMutation) BrandCleared() bool {
	return m.BrandIDCleared() || m.clearedbrand
}

// BrandIDs returns the "brand" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// BrandID instead. It exists only for internal usage by the builders.
func (m *PlanIntroduceMutation) BrandIDs() (ids []uint64) {
	if id := m.brand; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetBrand resets all changes to the "brand" edge.
func (m *PlanIntroduceMutation) ResetBrand() {
	m.brand = nil
	m.clearedbrand = false
}

// Where appends a list predicates to the PlanIntroduceMutation builder.
func (m *PlanIntroduceMutation) Where(ps ...predicate.PlanIntroduce) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *PlanIntroduceMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (PlanIntroduce).
func (m *PlanIntroduceMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *PlanIntroduceMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.created_at != nil {
		fields = append(fields, planintroduce.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, planintroduce.FieldUpdatedAt)
	}
	if m.brand != nil {
		fields = append(fields, planintroduce.FieldBrandID)
	}
	if m.model != nil {
		fields = append(fields, planintroduce.FieldModel)
	}
	if m.image != nil {
		fields = append(fields, planintroduce.FieldImage)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *PlanIntroduceMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case planintroduce.FieldCreatedAt:
		return m.CreatedAt()
	case planintroduce.FieldUpdatedAt:
		return m.UpdatedAt()
	case planintroduce.FieldBrandID:
		return m.BrandID()
	case planintroduce.FieldModel:
		return m.Model()
	case planintroduce.FieldImage:
		return m.Image()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *PlanIntroduceMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case planintroduce.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case planintroduce.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	case planintroduce.FieldBrandID:
		return m.OldBrandID(ctx)
	case planintroduce.FieldModel:
		return m.OldModel(ctx)
	case planintroduce.FieldImage:
		return m.OldImage(ctx)
	}
	return nil, fmt.Errorf("unknown PlanIntroduce field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *PlanIntroduceMutation) SetField(name string, value ent.Value) error {
	switch name {
	case planintroduce.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case planintroduce.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case planintroduce.FieldBrandID:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBrandID(v)
		return nil
	case planintroduce.FieldModel:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetModel(v)
		return nil
	case planintroduce.FieldImage:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetImage(v)
		return nil
	}
	return fmt.Errorf("unknown PlanIntroduce field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *PlanIntroduceMutation) AddedFields() []string {
	var fields []string
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *PlanIntroduceMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *PlanIntroduceMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown PlanIntroduce numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *PlanIntroduceMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(planintroduce.FieldBrandID) {
		fields = append(fields, planintroduce.FieldBrandID)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *PlanIntroduceMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *PlanIntroduceMutation) ClearField(name string) error {
	switch name {
	case planintroduce.FieldBrandID:
		m.ClearBrandID()
		return nil
	}
	return fmt.Errorf("unknown PlanIntroduce nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *PlanIntroduceMutation) ResetField(name string) error {
	switch name {
	case planintroduce.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case planintroduce.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case planintroduce.FieldBrandID:
		m.ResetBrandID()
		return nil
	case planintroduce.FieldModel:
		m.ResetModel()
		return nil
	case planintroduce.FieldImage:
		m.ResetImage()
		return nil
	}
	return fmt.Errorf("unknown PlanIntroduce field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *PlanIntroduceMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.brand != nil {
		edges = append(edges, planintroduce.EdgeBrand)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *PlanIntroduceMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case planintroduce.EdgeBrand:
		if id := m.brand; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *PlanIntroduceMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *PlanIntroduceMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *PlanIntroduceMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedbrand {
		edges = append(edges, planintroduce.EdgeBrand)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *PlanIntroduceMutation) EdgeCleared(name string) bool {
	switch name {
	case planintroduce.EdgeBrand:
		return m.clearedbrand
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *PlanIntroduceMutation) ClearEdge(name string) error {
	switch name {
	case planintroduce.EdgeBrand:
		m.ClearBrand()
		return nil
	}
	return fmt.Errorf("unknown PlanIntroduce unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *PlanIntroduceMutation) ResetEdge(name string) error {
	switch name {
	case planintroduce.EdgeBrand:
		m.ResetBrand()
		return nil
	}
	return fmt.Errorf("unknown PlanIntroduce edge %s", name)
}
