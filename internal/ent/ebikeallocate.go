// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/contract"
	"github.com/auroraride/aurservd/internal/ent/ebike"
	"github.com/auroraride/aurservd/internal/ent/ebikeallocate"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
)

// EbikeAllocate is the model entity for the EbikeAllocate schema.
type EbikeAllocate struct {
	config `json:"-"`
	// ID of the ent.
	ID uint64 `json:"id,omitempty"`
	// 店员ID
	EmployeeID uint64 `json:"employee_id,omitempty"`
	// 门店ID
	StoreID uint64 `json:"store_id,omitempty"`
	// EbikeID holds the value of the "ebike_id" field.
	EbikeID uint64 `json:"ebike_id,omitempty"`
	// BrandID holds the value of the "brand_id" field.
	BrandID uint64 `json:"brand_id,omitempty"`
	// SubscribeID holds the value of the "subscribe_id" field.
	SubscribeID uint64 `json:"subscribe_id,omitempty"`
	// 骑手ID
	RiderID uint64 `json:"rider_id,omitempty"`
	// 分配状态
	Status uint8 `json:"status,omitempty"`
	// 电车信息
	Info *model.EbikeAllocate `json:"info,omitempty"`
	// 分配时间
	Time time.Time `json:"time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EbikeAllocateQuery when eager-loading is set.
	Edges EbikeAllocateEdges `json:"edges"`
}

// EbikeAllocateEdges holds the relations/edges for other nodes in the graph.
type EbikeAllocateEdges struct {
	// Employee holds the value of the employee edge.
	Employee *Employee `json:"employee,omitempty"`
	// Store holds the value of the store edge.
	Store *Store `json:"store,omitempty"`
	// Ebike holds the value of the ebike edge.
	Ebike *Ebike `json:"ebike,omitempty"`
	// Brand holds the value of the brand edge.
	Brand *EbikeBrand `json:"brand,omitempty"`
	// Subscribe holds the value of the subscribe edge.
	Subscribe *Subscribe `json:"subscribe,omitempty"`
	// 骑手
	Rider *Rider `json:"rider,omitempty"`
	// Contract holds the value of the contract edge.
	Contract *Contract `json:"contract,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [7]bool
}

// EmployeeOrErr returns the Employee value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EbikeAllocateEdges) EmployeeOrErr() (*Employee, error) {
	if e.loadedTypes[0] {
		if e.Employee == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: employee.Label}
		}
		return e.Employee, nil
	}
	return nil, &NotLoadedError{edge: "employee"}
}

// StoreOrErr returns the Store value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EbikeAllocateEdges) StoreOrErr() (*Store, error) {
	if e.loadedTypes[1] {
		if e.Store == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: store.Label}
		}
		return e.Store, nil
	}
	return nil, &NotLoadedError{edge: "store"}
}

// EbikeOrErr returns the Ebike value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EbikeAllocateEdges) EbikeOrErr() (*Ebike, error) {
	if e.loadedTypes[2] {
		if e.Ebike == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: ebike.Label}
		}
		return e.Ebike, nil
	}
	return nil, &NotLoadedError{edge: "ebike"}
}

// BrandOrErr returns the Brand value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EbikeAllocateEdges) BrandOrErr() (*EbikeBrand, error) {
	if e.loadedTypes[3] {
		if e.Brand == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: ebikebrand.Label}
		}
		return e.Brand, nil
	}
	return nil, &NotLoadedError{edge: "brand"}
}

// SubscribeOrErr returns the Subscribe value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EbikeAllocateEdges) SubscribeOrErr() (*Subscribe, error) {
	if e.loadedTypes[4] {
		if e.Subscribe == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: subscribe.Label}
		}
		return e.Subscribe, nil
	}
	return nil, &NotLoadedError{edge: "subscribe"}
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EbikeAllocateEdges) RiderOrErr() (*Rider, error) {
	if e.loadedTypes[5] {
		if e.Rider == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: rider.Label}
		}
		return e.Rider, nil
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// ContractOrErr returns the Contract value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EbikeAllocateEdges) ContractOrErr() (*Contract, error) {
	if e.loadedTypes[6] {
		if e.Contract == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: contract.Label}
		}
		return e.Contract, nil
	}
	return nil, &NotLoadedError{edge: "contract"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*EbikeAllocate) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case ebikeallocate.FieldInfo:
			values[i] = new([]byte)
		case ebikeallocate.FieldID, ebikeallocate.FieldEmployeeID, ebikeallocate.FieldStoreID, ebikeallocate.FieldEbikeID, ebikeallocate.FieldBrandID, ebikeallocate.FieldSubscribeID, ebikeallocate.FieldRiderID, ebikeallocate.FieldStatus:
			values[i] = new(sql.NullInt64)
		case ebikeallocate.FieldTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type EbikeAllocate", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the EbikeAllocate fields.
func (ea *EbikeAllocate) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case ebikeallocate.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ea.ID = uint64(value.Int64)
		case ebikeallocate.FieldEmployeeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field employee_id", values[i])
			} else if value.Valid {
				ea.EmployeeID = uint64(value.Int64)
			}
		case ebikeallocate.FieldStoreID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field store_id", values[i])
			} else if value.Valid {
				ea.StoreID = uint64(value.Int64)
			}
		case ebikeallocate.FieldEbikeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field ebike_id", values[i])
			} else if value.Valid {
				ea.EbikeID = uint64(value.Int64)
			}
		case ebikeallocate.FieldBrandID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field brand_id", values[i])
			} else if value.Valid {
				ea.BrandID = uint64(value.Int64)
			}
		case ebikeallocate.FieldSubscribeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field subscribe_id", values[i])
			} else if value.Valid {
				ea.SubscribeID = uint64(value.Int64)
			}
		case ebikeallocate.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				ea.RiderID = uint64(value.Int64)
			}
		case ebikeallocate.FieldStatus:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				ea.Status = uint8(value.Int64)
			}
		case ebikeallocate.FieldInfo:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field info", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ea.Info); err != nil {
					return fmt.Errorf("unmarshal field info: %w", err)
				}
			}
		case ebikeallocate.FieldTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field time", values[i])
			} else if value.Valid {
				ea.Time = value.Time
			}
		}
	}
	return nil
}

// QueryEmployee queries the "employee" edge of the EbikeAllocate entity.
func (ea *EbikeAllocate) QueryEmployee() *EmployeeQuery {
	return (&EbikeAllocateClient{config: ea.config}).QueryEmployee(ea)
}

// QueryStore queries the "store" edge of the EbikeAllocate entity.
func (ea *EbikeAllocate) QueryStore() *StoreQuery {
	return (&EbikeAllocateClient{config: ea.config}).QueryStore(ea)
}

// QueryEbike queries the "ebike" edge of the EbikeAllocate entity.
func (ea *EbikeAllocate) QueryEbike() *EbikeQuery {
	return (&EbikeAllocateClient{config: ea.config}).QueryEbike(ea)
}

// QueryBrand queries the "brand" edge of the EbikeAllocate entity.
func (ea *EbikeAllocate) QueryBrand() *EbikeBrandQuery {
	return (&EbikeAllocateClient{config: ea.config}).QueryBrand(ea)
}

// QuerySubscribe queries the "subscribe" edge of the EbikeAllocate entity.
func (ea *EbikeAllocate) QuerySubscribe() *SubscribeQuery {
	return (&EbikeAllocateClient{config: ea.config}).QuerySubscribe(ea)
}

// QueryRider queries the "rider" edge of the EbikeAllocate entity.
func (ea *EbikeAllocate) QueryRider() *RiderQuery {
	return (&EbikeAllocateClient{config: ea.config}).QueryRider(ea)
}

// QueryContract queries the "contract" edge of the EbikeAllocate entity.
func (ea *EbikeAllocate) QueryContract() *ContractQuery {
	return (&EbikeAllocateClient{config: ea.config}).QueryContract(ea)
}

// Update returns a builder for updating this EbikeAllocate.
// Note that you need to call EbikeAllocate.Unwrap() before calling this method if this EbikeAllocate
// was returned from a transaction, and the transaction was committed or rolled back.
func (ea *EbikeAllocate) Update() *EbikeAllocateUpdateOne {
	return (&EbikeAllocateClient{config: ea.config}).UpdateOne(ea)
}

// Unwrap unwraps the EbikeAllocate entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ea *EbikeAllocate) Unwrap() *EbikeAllocate {
	_tx, ok := ea.config.driver.(*txDriver)
	if !ok {
		panic("ent: EbikeAllocate is not a transactional entity")
	}
	ea.config.driver = _tx.drv
	return ea
}

// String implements the fmt.Stringer.
func (ea *EbikeAllocate) String() string {
	var builder strings.Builder
	builder.WriteString("EbikeAllocate(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ea.ID))
	builder.WriteString("employee_id=")
	builder.WriteString(fmt.Sprintf("%v", ea.EmployeeID))
	builder.WriteString(", ")
	builder.WriteString("store_id=")
	builder.WriteString(fmt.Sprintf("%v", ea.StoreID))
	builder.WriteString(", ")
	builder.WriteString("ebike_id=")
	builder.WriteString(fmt.Sprintf("%v", ea.EbikeID))
	builder.WriteString(", ")
	builder.WriteString("brand_id=")
	builder.WriteString(fmt.Sprintf("%v", ea.BrandID))
	builder.WriteString(", ")
	builder.WriteString("subscribe_id=")
	builder.WriteString(fmt.Sprintf("%v", ea.SubscribeID))
	builder.WriteString(", ")
	builder.WriteString("rider_id=")
	builder.WriteString(fmt.Sprintf("%v", ea.RiderID))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", ea.Status))
	builder.WriteString(", ")
	builder.WriteString("info=")
	builder.WriteString(fmt.Sprintf("%v", ea.Info))
	builder.WriteString(", ")
	builder.WriteString("time=")
	builder.WriteString(ea.Time.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// EbikeAllocates is a parsable slice of EbikeAllocate.
type EbikeAllocates []*EbikeAllocate

func (ea EbikeAllocates) config(cfg config) {
	for _i := range ea {
		ea[_i].config = cfg
	}
}