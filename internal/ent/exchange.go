// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/exchange"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
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
	// 创建人
	Creator *model.Modifier `json:"creator,omitempty"`
	// 最后修改人
	LastModifier *model.Modifier `json:"last_modifier,omitempty"`
	// 管理员改动原因/备注
	Remark string `json:"remark,omitempty"`
	// SubscribeID holds the value of the "subscribe_id" field.
	SubscribeID uint64 `json:"subscribe_id,omitempty"`
	// 城市ID
	CityID uint64 `json:"city_id,omitempty"`
	// 门店ID
	StoreID *uint64 `json:"store_id,omitempty"`
	// 企业ID
	EnterpriseID *uint64 `json:"enterprise_id,omitempty"`
	// 站点ID
	StationID *uint64 `json:"station_id,omitempty"`
	// 骑手ID
	RiderID uint64 `json:"rider_id,omitempty"`
	// 店员ID
	EmployeeID *uint64 `json:"employee_id,omitempty"`
	// UUID holds the value of the "uuid" field.
	UUID string `json:"uuid,omitempty"`
	// 电柜ID
	CabinetID uint64 `json:"cabinet_id,omitempty"`
	// 是否成功
	Success bool `json:"success,omitempty"`
	// 电池型号
	Model string `json:"model,omitempty"`
	// 是否备用方案
	Alternative bool `json:"alternative,omitempty"`
	// 换电开始时间
	StartAt time.Time `json:"start_at,omitempty"`
	// 换电结束时间
	FinishAt time.Time `json:"finish_at,omitempty"`
	// 换电耗时(s)
	Duration int `json:"duration,omitempty"`
	// 骑手当前电池编号
	RiderBattery *string `json:"rider_battery,omitempty"`
	// 放入电池编号
	PutinBattery *string `json:"putin_battery,omitempty"`
	// 取出电池编号
	PutoutBattery *string `json:"putout_battery,omitempty"`
	// 电柜信息
	CabinetInfo *model.ExchangeCabinetInfo `json:"cabinet_info,omitempty"`
	// 空仓信息
	Empty *model.BinInfo `json:"empty,omitempty"`
	// 满仓信息
	Fully *model.BinInfo `json:"fully,omitempty"`
	// 步骤信息
	Steps []*model.ExchangeStepInfo `json:"steps,omitempty"`
	// 消息
	Message string `json:"message,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ExchangeQuery when eager-loading is set.
	Edges        ExchangeEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ExchangeEdges holds the relations/edges for other nodes in the graph.
type ExchangeEdges struct {
	// Subscribe holds the value of the subscribe edge.
	Subscribe *Subscribe `json:"subscribe,omitempty"`
	// City holds the value of the city edge.
	City *City `json:"city,omitempty"`
	// Store holds the value of the store edge.
	Store *Store `json:"store,omitempty"`
	// Enterprise holds the value of the enterprise edge.
	Enterprise *Enterprise `json:"enterprise,omitempty"`
	// Station holds the value of the station edge.
	Station *EnterpriseStation `json:"station,omitempty"`
	// Cabinet holds the value of the cabinet edge.
	Cabinet *Cabinet `json:"cabinet,omitempty"`
	// Rider holds the value of the rider edge.
	Rider *Rider `json:"rider,omitempty"`
	// Employee holds the value of the employee edge.
	Employee *Employee `json:"employee,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [8]bool
}

// SubscribeOrErr returns the Subscribe value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) SubscribeOrErr() (*Subscribe, error) {
	if e.Subscribe != nil {
		return e.Subscribe, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: subscribe.Label}
	}
	return nil, &NotLoadedError{edge: "subscribe"}
}

// CityOrErr returns the City value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) CityOrErr() (*City, error) {
	if e.City != nil {
		return e.City, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: city.Label}
	}
	return nil, &NotLoadedError{edge: "city"}
}

// StoreOrErr returns the Store value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) StoreOrErr() (*Store, error) {
	if e.Store != nil {
		return e.Store, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: store.Label}
	}
	return nil, &NotLoadedError{edge: "store"}
}

// EnterpriseOrErr returns the Enterprise value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) EnterpriseOrErr() (*Enterprise, error) {
	if e.Enterprise != nil {
		return e.Enterprise, nil
	} else if e.loadedTypes[3] {
		return nil, &NotFoundError{label: enterprise.Label}
	}
	return nil, &NotLoadedError{edge: "enterprise"}
}

// StationOrErr returns the Station value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) StationOrErr() (*EnterpriseStation, error) {
	if e.Station != nil {
		return e.Station, nil
	} else if e.loadedTypes[4] {
		return nil, &NotFoundError{label: enterprisestation.Label}
	}
	return nil, &NotLoadedError{edge: "station"}
}

// CabinetOrErr returns the Cabinet value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) CabinetOrErr() (*Cabinet, error) {
	if e.Cabinet != nil {
		return e.Cabinet, nil
	} else if e.loadedTypes[5] {
		return nil, &NotFoundError{label: cabinet.Label}
	}
	return nil, &NotLoadedError{edge: "cabinet"}
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) RiderOrErr() (*Rider, error) {
	if e.Rider != nil {
		return e.Rider, nil
	} else if e.loadedTypes[6] {
		return nil, &NotFoundError{label: rider.Label}
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// EmployeeOrErr returns the Employee value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ExchangeEdges) EmployeeOrErr() (*Employee, error) {
	if e.Employee != nil {
		return e.Employee, nil
	} else if e.loadedTypes[7] {
		return nil, &NotFoundError{label: employee.Label}
	}
	return nil, &NotLoadedError{edge: "employee"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Exchange) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case exchange.FieldCreator, exchange.FieldLastModifier, exchange.FieldCabinetInfo, exchange.FieldEmpty, exchange.FieldFully, exchange.FieldSteps:
			values[i] = new([]byte)
		case exchange.FieldSuccess, exchange.FieldAlternative:
			values[i] = new(sql.NullBool)
		case exchange.FieldID, exchange.FieldSubscribeID, exchange.FieldCityID, exchange.FieldStoreID, exchange.FieldEnterpriseID, exchange.FieldStationID, exchange.FieldRiderID, exchange.FieldEmployeeID, exchange.FieldCabinetID, exchange.FieldDuration:
			values[i] = new(sql.NullInt64)
		case exchange.FieldRemark, exchange.FieldUUID, exchange.FieldModel, exchange.FieldRiderBattery, exchange.FieldPutinBattery, exchange.FieldPutoutBattery, exchange.FieldMessage:
			values[i] = new(sql.NullString)
		case exchange.FieldCreatedAt, exchange.FieldUpdatedAt, exchange.FieldDeletedAt, exchange.FieldStartAt, exchange.FieldFinishAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Exchange fields.
func (e *Exchange) assignValues(columns []string, values []any) error {
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
		case exchange.FieldSubscribeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field subscribe_id", values[i])
			} else if value.Valid {
				e.SubscribeID = uint64(value.Int64)
			}
		case exchange.FieldCityID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field city_id", values[i])
			} else if value.Valid {
				e.CityID = uint64(value.Int64)
			}
		case exchange.FieldStoreID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field store_id", values[i])
			} else if value.Valid {
				e.StoreID = new(uint64)
				*e.StoreID = uint64(value.Int64)
			}
		case exchange.FieldEnterpriseID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field enterprise_id", values[i])
			} else if value.Valid {
				e.EnterpriseID = new(uint64)
				*e.EnterpriseID = uint64(value.Int64)
			}
		case exchange.FieldStationID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field station_id", values[i])
			} else if value.Valid {
				e.StationID = new(uint64)
				*e.StationID = uint64(value.Int64)
			}
		case exchange.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				e.RiderID = uint64(value.Int64)
			}
		case exchange.FieldEmployeeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field employee_id", values[i])
			} else if value.Valid {
				e.EmployeeID = new(uint64)
				*e.EmployeeID = uint64(value.Int64)
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
		case exchange.FieldModel:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field model", values[i])
			} else if value.Valid {
				e.Model = value.String
			}
		case exchange.FieldAlternative:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field alternative", values[i])
			} else if value.Valid {
				e.Alternative = value.Bool
			}
		case exchange.FieldStartAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field start_at", values[i])
			} else if value.Valid {
				e.StartAt = value.Time
			}
		case exchange.FieldFinishAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field finish_at", values[i])
			} else if value.Valid {
				e.FinishAt = value.Time
			}
		case exchange.FieldDuration:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field duration", values[i])
			} else if value.Valid {
				e.Duration = int(value.Int64)
			}
		case exchange.FieldRiderBattery:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field rider_battery", values[i])
			} else if value.Valid {
				e.RiderBattery = new(string)
				*e.RiderBattery = value.String
			}
		case exchange.FieldPutinBattery:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field putin_battery", values[i])
			} else if value.Valid {
				e.PutinBattery = new(string)
				*e.PutinBattery = value.String
			}
		case exchange.FieldPutoutBattery:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field putout_battery", values[i])
			} else if value.Valid {
				e.PutoutBattery = new(string)
				*e.PutoutBattery = value.String
			}
		case exchange.FieldCabinetInfo:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field cabinet_info", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.CabinetInfo); err != nil {
					return fmt.Errorf("unmarshal field cabinet_info: %w", err)
				}
			}
		case exchange.FieldEmpty:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field empty", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.Empty); err != nil {
					return fmt.Errorf("unmarshal field empty: %w", err)
				}
			}
		case exchange.FieldFully:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field fully", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.Fully); err != nil {
					return fmt.Errorf("unmarshal field fully: %w", err)
				}
			}
		case exchange.FieldSteps:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field steps", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.Steps); err != nil {
					return fmt.Errorf("unmarshal field steps: %w", err)
				}
			}
		case exchange.FieldMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field message", values[i])
			} else if value.Valid {
				e.Message = value.String
			}
		default:
			e.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Exchange.
// This includes values selected through modifiers, order, etc.
func (e *Exchange) Value(name string) (ent.Value, error) {
	return e.selectValues.Get(name)
}

// QuerySubscribe queries the "subscribe" edge of the Exchange entity.
func (e *Exchange) QuerySubscribe() *SubscribeQuery {
	return NewExchangeClient(e.config).QuerySubscribe(e)
}

// QueryCity queries the "city" edge of the Exchange entity.
func (e *Exchange) QueryCity() *CityQuery {
	return NewExchangeClient(e.config).QueryCity(e)
}

// QueryStore queries the "store" edge of the Exchange entity.
func (e *Exchange) QueryStore() *StoreQuery {
	return NewExchangeClient(e.config).QueryStore(e)
}

// QueryEnterprise queries the "enterprise" edge of the Exchange entity.
func (e *Exchange) QueryEnterprise() *EnterpriseQuery {
	return NewExchangeClient(e.config).QueryEnterprise(e)
}

// QueryStation queries the "station" edge of the Exchange entity.
func (e *Exchange) QueryStation() *EnterpriseStationQuery {
	return NewExchangeClient(e.config).QueryStation(e)
}

// QueryCabinet queries the "cabinet" edge of the Exchange entity.
func (e *Exchange) QueryCabinet() *CabinetQuery {
	return NewExchangeClient(e.config).QueryCabinet(e)
}

// QueryRider queries the "rider" edge of the Exchange entity.
func (e *Exchange) QueryRider() *RiderQuery {
	return NewExchangeClient(e.config).QueryRider(e)
}

// QueryEmployee queries the "employee" edge of the Exchange entity.
func (e *Exchange) QueryEmployee() *EmployeeQuery {
	return NewExchangeClient(e.config).QueryEmployee(e)
}

// Update returns a builder for updating this Exchange.
// Note that you need to call Exchange.Unwrap() before calling this method if this Exchange
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Exchange) Update() *ExchangeUpdateOne {
	return NewExchangeClient(e.config).UpdateOne(e)
}

// Unwrap unwraps the Exchange entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Exchange) Unwrap() *Exchange {
	_tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Exchange is not a transactional entity")
	}
	e.config.driver = _tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Exchange) String() string {
	var builder strings.Builder
	builder.WriteString("Exchange(")
	builder.WriteString(fmt.Sprintf("id=%v, ", e.ID))
	builder.WriteString("created_at=")
	builder.WriteString(e.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(e.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := e.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", e.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", e.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(e.Remark)
	builder.WriteString(", ")
	builder.WriteString("subscribe_id=")
	builder.WriteString(fmt.Sprintf("%v", e.SubscribeID))
	builder.WriteString(", ")
	builder.WriteString("city_id=")
	builder.WriteString(fmt.Sprintf("%v", e.CityID))
	builder.WriteString(", ")
	if v := e.StoreID; v != nil {
		builder.WriteString("store_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := e.EnterpriseID; v != nil {
		builder.WriteString("enterprise_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := e.StationID; v != nil {
		builder.WriteString("station_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("rider_id=")
	builder.WriteString(fmt.Sprintf("%v", e.RiderID))
	builder.WriteString(", ")
	if v := e.EmployeeID; v != nil {
		builder.WriteString("employee_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("uuid=")
	builder.WriteString(e.UUID)
	builder.WriteString(", ")
	builder.WriteString("cabinet_id=")
	builder.WriteString(fmt.Sprintf("%v", e.CabinetID))
	builder.WriteString(", ")
	builder.WriteString("success=")
	builder.WriteString(fmt.Sprintf("%v", e.Success))
	builder.WriteString(", ")
	builder.WriteString("model=")
	builder.WriteString(e.Model)
	builder.WriteString(", ")
	builder.WriteString("alternative=")
	builder.WriteString(fmt.Sprintf("%v", e.Alternative))
	builder.WriteString(", ")
	builder.WriteString("start_at=")
	builder.WriteString(e.StartAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("finish_at=")
	builder.WriteString(e.FinishAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("duration=")
	builder.WriteString(fmt.Sprintf("%v", e.Duration))
	builder.WriteString(", ")
	if v := e.RiderBattery; v != nil {
		builder.WriteString("rider_battery=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := e.PutinBattery; v != nil {
		builder.WriteString("putin_battery=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := e.PutoutBattery; v != nil {
		builder.WriteString("putout_battery=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("cabinet_info=")
	builder.WriteString(fmt.Sprintf("%v", e.CabinetInfo))
	builder.WriteString(", ")
	builder.WriteString("empty=")
	builder.WriteString(fmt.Sprintf("%v", e.Empty))
	builder.WriteString(", ")
	builder.WriteString("fully=")
	builder.WriteString(fmt.Sprintf("%v", e.Fully))
	builder.WriteString(", ")
	builder.WriteString("steps=")
	builder.WriteString(fmt.Sprintf("%v", e.Steps))
	builder.WriteString(", ")
	builder.WriteString("message=")
	builder.WriteString(e.Message)
	builder.WriteByte(')')
	return builder.String()
}

// Exchanges is a parsable slice of Exchange.
type Exchanges []*Exchange
