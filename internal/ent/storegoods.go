// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/internal/ent/goods"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/storegoods"
)

// StoreGoods is the model entity for the StoreGoods schema.
type StoreGoods struct {
	config `json:"-"`
	// ID of the ent.
	ID uint64 `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	// GoodsID holds the value of the "goods_id" field.
	GoodsID uint64 `json:"goods_id,omitempty"`
	// StoreID holds the value of the "store_id" field.
	StoreID uint64 `json:"store_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the StoreGoodsQuery when eager-loading is set.
	Edges        StoreGoodsEdges `json:"edges"`
	selectValues sql.SelectValues
}

// StoreGoodsEdges holds the relations/edges for other nodes in the graph.
type StoreGoodsEdges struct {
	// Goods holds the value of the goods edge.
	Goods *Goods `json:"goods,omitempty"`
	// Store holds the value of the store edge.
	Store *Store `json:"store,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// GoodsOrErr returns the Goods value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e StoreGoodsEdges) GoodsOrErr() (*Goods, error) {
	if e.Goods != nil {
		return e.Goods, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: goods.Label}
	}
	return nil, &NotLoadedError{edge: "goods"}
}

// StoreOrErr returns the Store value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e StoreGoodsEdges) StoreOrErr() (*Store, error) {
	if e.Store != nil {
		return e.Store, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: store.Label}
	}
	return nil, &NotLoadedError{edge: "store"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*StoreGoods) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case storegoods.FieldID, storegoods.FieldGoodsID, storegoods.FieldStoreID:
			values[i] = new(sql.NullInt64)
		case storegoods.FieldCreatedAt, storegoods.FieldUpdatedAt, storegoods.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the StoreGoods fields.
func (sg *StoreGoods) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case storegoods.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			sg.ID = uint64(value.Int64)
		case storegoods.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				sg.CreatedAt = value.Time
			}
		case storegoods.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				sg.UpdatedAt = value.Time
			}
		case storegoods.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				sg.DeletedAt = new(time.Time)
				*sg.DeletedAt = value.Time
			}
		case storegoods.FieldGoodsID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field goods_id", values[i])
			} else if value.Valid {
				sg.GoodsID = uint64(value.Int64)
			}
		case storegoods.FieldStoreID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field store_id", values[i])
			} else if value.Valid {
				sg.StoreID = uint64(value.Int64)
			}
		default:
			sg.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the StoreGoods.
// This includes values selected through modifiers, order, etc.
func (sg *StoreGoods) Value(name string) (ent.Value, error) {
	return sg.selectValues.Get(name)
}

// QueryGoods queries the "goods" edge of the StoreGoods entity.
func (sg *StoreGoods) QueryGoods() *GoodsQuery {
	return NewStoreGoodsClient(sg.config).QueryGoods(sg)
}

// QueryStore queries the "store" edge of the StoreGoods entity.
func (sg *StoreGoods) QueryStore() *StoreQuery {
	return NewStoreGoodsClient(sg.config).QueryStore(sg)
}

// Update returns a builder for updating this StoreGoods.
// Note that you need to call StoreGoods.Unwrap() before calling this method if this StoreGoods
// was returned from a transaction, and the transaction was committed or rolled back.
func (sg *StoreGoods) Update() *StoreGoodsUpdateOne {
	return NewStoreGoodsClient(sg.config).UpdateOne(sg)
}

// Unwrap unwraps the StoreGoods entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (sg *StoreGoods) Unwrap() *StoreGoods {
	_tx, ok := sg.config.driver.(*txDriver)
	if !ok {
		panic("ent: StoreGoods is not a transactional entity")
	}
	sg.config.driver = _tx.drv
	return sg
}

// String implements the fmt.Stringer.
func (sg *StoreGoods) String() string {
	var builder strings.Builder
	builder.WriteString("StoreGoods(")
	builder.WriteString(fmt.Sprintf("id=%v, ", sg.ID))
	builder.WriteString("created_at=")
	builder.WriteString(sg.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(sg.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := sg.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("goods_id=")
	builder.WriteString(fmt.Sprintf("%v", sg.GoodsID))
	builder.WriteString(", ")
	builder.WriteString("store_id=")
	builder.WriteString(fmt.Sprintf("%v", sg.StoreID))
	builder.WriteByte(')')
	return builder.String()
}

// StoreGoodsSlice is a parsable slice of StoreGoods.
type StoreGoodsSlice []*StoreGoods