// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/exchange"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
)

// Exchange is the model entity for the Exchange schema.
type Exchange struct {
	config `json:"-"`
	// ID of the ent.
	ID uint64 `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	// Creator holds the value of the "creator" field.
	// 创建人
	Creator *model.Modifier `json:"creator,omitempty"`
	// LastModifier holds the value of the "last_modifier" field.
	// 最后修改人
	LastModifier *model.Modifier `json:"last_modifier,omitempty"`
	// Remark holds the value of the "remark" field.
	// 管理员改动原因/备注
	Remark string `json:"remark,omitempty"`
	// CityID holds the value of the "city_id" field.
	// 城市ID
	CityID uint64 `json:"city_id,omitempty"`
	// EmployeeID holds the value of the "employee_id" field.
	// 操作店员ID
	EmployeeID uint64 `json:"employee_id,omitempty"`
	// StoreID holds the value of the "store_id" field.
	StoreID uint64 `json:"store_id,omitempty"`
	// RiderID holds the value of the "rider_id" field.
	// 骑手ID
	RiderID uint64 `json:"rider_id,omitempty"`
	// UUID holds the value of the "uuid" field.
	UUID string `json:"uuid,omitempty"`
	// CabinetID holds the value of the "cabinet_id" field.
	// 电柜ID
	CabinetID uint64 `json:"cabinet_id,omitempty"`
	// Success holds the value of the "success" field.
	// 是否成功
	Success bool `json:"success,omitempty"`
	// Detail holds the value of the "detail" field.
	// 电柜换电信息
	Detail *model.ExchangeCabinet `json:"detail,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ExchangeQuery when eager-loading is set.
	Edges ExchangeEdges `json:"edges"`
}

// ExchangeEdges holds the relations/edges for other nodes in the graph.
type ExchangeEdges struct {
	// City holds the value of the city edge.
	City *City `json:"city,omitempty"`
	// Employee holds the value of the employee edge.
	Employee *Employee `json:"employee,omitempty"`
	// Store holds the value of the store edge.
	Store *Store `json:"store,omitempty"`
	// Cabinet holds the value of the cabinet edge.
	Cabinet *Cabinet `json:"cabinet,omitempty"`
	// Rider holds the value of the rider edge.
	Rider *Rider `json:"rider,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [5]bool
}

// CityOrErr returns the City value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) CityOrErr() (*City, error) {
	if e.loadedTypes[0] {
		if e.City == nil {
			// The edge city was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: city.Label}
		}
		return e.City, nil
	}
	return nil, &NotLoadedError{edge: "city"}
}

// EmployeeOrErr returns the Employee value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) EmployeeOrErr() (*Employee, error) {
	if e.loadedTypes[1] {
		if e.Employee == nil {
			// The edge employee was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: employee.Label}
		}
		return e.Employee, nil
	}
	return nil, &NotLoadedError{edge: "employee"}
}

// StoreOrErr returns the Store value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) StoreOrErr() (*Store, error) {
	if e.loadedTypes[2] {
		if e.Store == nil {
			// The edge store was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: store.Label}
		}
		return e.Store, nil
	}
	return nil, &NotLoadedError{edge: "store"}
}

// CabinetOrErr returns the Cabinet value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) CabinetOrErr() (*Cabinet, error) {
	if e.loadedTypes[3] {
		if e.Cabinet == nil {
			// The edge cabinet was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: cabinet.Label}
		}
		return e.Cabinet, nil
	}
	return nil, &NotLoadedError{edge: "cabinet"}
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) RiderOrErr() (*Rider, error) {
	if e.loadedTypes[4] {
		if e.Rider == nil {
			// The edge rider was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: rider.Label}
		}
		return e.Rider, nil
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Exchange) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case exchange.FieldCreator, exchange.FieldLastModifier, exchange.FieldDetail:
			values[i] = new([]byte)
		case exchange.FieldSuccess:
			values[i] = new(sql.NullBool)
		case exchange.FieldID, exchange.FieldCityID, exchange.FieldEmployeeID, exchange.FieldStoreID, exchange.FieldRiderID, exchange.FieldCabinetID:
			values[i] = new(sql.NullInt64)
		case exchange.FieldRemark, exchange.FieldUUID:
			values[i] = new(sql.NullString)
		case exchange.FieldCreatedAt, exchange.FieldUpdatedAt, exchange.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Exchange", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Exchange fields.
func (e *Exchange) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case exchange.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			e.ID = uint64(value.Int64)
		case exchange.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				e.CreatedAt = value.Time
			}
		case exchange.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				e.UpdatedAt = value.Time
			}
		case exchange.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				e.DeletedAt = new(time.Time)
				*e.DeletedAt = value.Time
			}
		case exchange.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case exchange.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case exchange.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				e.Remark = value.String
			}
		case exchange.FieldCityID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field city_id", values[i])
			} else if value.Valid {
				e.CityID = uint64(value.Int64)
			}
		case exchange.FieldEmployeeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field employee_id", values[i])
			} else if value.Valid {
				e.EmployeeID = uint64(value.Int64)
			}
		case exchange.FieldStoreID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field store_id", values[i])
			} else if value.Valid {
				e.StoreID = uint64(value.Int64)
			}
		case exchange.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				e.RiderID = uint64(value.Int64)
			}
		case exchange.FieldUUID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field uuid", values[i])
			} else if value.Valid {
				e.UUID = value.String
			}
		case exchange.FieldCabinetID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field cabinet_id", values[i])
			} else if value.Valid {
				e.CabinetID = uint64(value.Int64)
			}
		case exchange.FieldSuccess:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field success", values[i])
			} else if value.Valid {
				e.Success = value.Bool
			}
		case exchange.FieldDetail:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field detail", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.Detail); err != nil {
					return fmt.Errorf("unmarshal field detail: %w", err)
				}
			}
		}
	}
	return nil
}

// QueryCity queries the "city" edge of the Exchange entity.
func (e *Exchange) QueryCity() *CityQuery {
	return (&ExchangeClient{config: e.config}).QueryCity(e)
}

// QueryEmployee queries the "employee" edge of the Exchange entity.
func (e *Exchange) QueryEmployee() *EmployeeQuery {
	return (&ExchangeClient{config: e.config}).QueryEmployee(e)
}

// QueryStore queries the "store" edge of the Exchange entity.
func (e *Exchange) QueryStore() *StoreQuery {
	return (&ExchangeClient{config: e.config}).QueryStore(e)
}

// QueryCabinet queries the "cabinet" edge of the Exchange entity.
func (e *Exchange) QueryCabinet() *CabinetQuery {
	return (&ExchangeClient{config: e.config}).QueryCabinet(e)
}

// QueryRider queries the "rider" edge of the Exchange entity.
func (e *Exchange) QueryRider() *RiderQuery {
	return (&ExchangeClient{config: e.config}).QueryRider(e)
}

// Update returns a builder for updating this Exchange.
// Note that you need to call Exchange.Unwrap() before calling this method if this Exchange
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Exchange) Update() *ExchangeUpdateOne {
	return (&ExchangeClient{config: e.config}).UpdateOne(e)
}

// Unwrap unwraps the Exchange entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Exchange) Unwrap() *Exchange {
	tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Exchange is not a transactional entity")
	}
	e.config.driver = tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Exchange) String() string {
	var builder strings.Builder
	builder.WriteString("Exchange(")
	builder.WriteString(fmt.Sprintf("id=%v", e.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(e.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(e.UpdatedAt.Format(time.ANSIC))
	if v := e.DeletedAt; v != nil {
		builder.WriteString(", deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", creator=")
	builder.WriteString(fmt.Sprintf("%v", e.Creator))
	builder.WriteString(", last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", e.LastModifier))
	builder.WriteString(", remark=")
	builder.WriteString(e.Remark)
	builder.WriteString(", city_id=")
	builder.WriteString(fmt.Sprintf("%v", e.CityID))
	builder.WriteString(", employee_id=")
	builder.WriteString(fmt.Sprintf("%v", e.EmployeeID))
	builder.WriteString(", store_id=")
	builder.WriteString(fmt.Sprintf("%v", e.StoreID))
	builder.WriteString(", rider_id=")
	builder.WriteString(fmt.Sprintf("%v", e.RiderID))
	builder.WriteString(", uuid=")
	builder.WriteString(e.UUID)
	builder.WriteString(", cabinet_id=")
	builder.WriteString(fmt.Sprintf("%v", e.CabinetID))
	builder.WriteString(", success=")
	builder.WriteString(fmt.Sprintf("%v", e.Success))
	builder.WriteString(", detail=")
	builder.WriteString(fmt.Sprintf("%v", e.Detail))
	builder.WriteByte(')')
	return builder.String()
}

// Exchanges is a parsable slice of Exchange.
type Exchanges []*Exchange

func (e Exchanges) config(cfg config) {
	for _i := range e {
		e[_i].config = cfg
	}
}