// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/adapter"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/city"
)

// Battery is the model entity for the Battery schema.
type Battery struct {
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
	// 城市ID
	CityID *uint64 `json:"city_id,omitempty"`
	// 骑手ID
	RiderID *uint64 `json:"rider_id,omitempty"`
	// 电柜ID
	CabinetID *uint64 `json:"cabinet_id,omitempty"`
	// 订阅ID
	SubscribeID *uint64 `json:"subscribe_id,omitempty"`
	// 所属团签
	EnterpriseID *uint64 `json:"enterprise_id,omitempty"`
	// 所属站点Id
	StationID *uint64 `json:"station_id,omitempty"`
	// 电池编号
	Sn string `json:"sn,omitempty"`
	// 品牌
	Brand adapter.BatteryBrand `json:"brand,omitempty"`
	// 是否启用
	Enable bool `json:"enable,omitempty"`
	// 电池型号
	Model string `json:"model,omitempty"`
	// 所在智能柜仓位序号
	Ordinal *int `json:"ordinal,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the BatteryQuery when eager-loading is set.
	Edges        BatteryEdges `json:"edges"`
	selectValues sql.SelectValues
}

// BatteryEdges holds the relations/edges for other nodes in the graph.
type BatteryEdges struct {
	// City holds the value of the city edge.
	City *City `json:"city,omitempty"`
	// 流转记录
	Flows []*BatteryFlow `json:"flows,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// CityOrErr returns the City value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BatteryEdges) CityOrErr() (*City, error) {
	if e.City != nil {
		return e.City, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: city.Label}
	}
	return nil, &NotLoadedError{edge: "city"}
}

// FlowsOrErr returns the Flows value or an error if the edge
// was not loaded in eager-loading.
func (e BatteryEdges) FlowsOrErr() ([]*BatteryFlow, error) {
	if e.loadedTypes[1] {
		return e.Flows, nil
	}
	return nil, &NotLoadedError{edge: "flows"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Battery) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case battery.FieldCreator, battery.FieldLastModifier:
			values[i] = new([]byte)
		case battery.FieldBrand:
			values[i] = new(adapter.BatteryBrand)
		case battery.FieldEnable:
			values[i] = new(sql.NullBool)
		case battery.FieldID, battery.FieldCityID, battery.FieldRiderID, battery.FieldCabinetID, battery.FieldSubscribeID, battery.FieldEnterpriseID, battery.FieldStationID, battery.FieldOrdinal:
			values[i] = new(sql.NullInt64)
		case battery.FieldRemark, battery.FieldSn, battery.FieldModel:
			values[i] = new(sql.NullString)
		case battery.FieldCreatedAt, battery.FieldUpdatedAt, battery.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Battery fields.
func (b *Battery) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case battery.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			b.ID = uint64(value.Int64)
		case battery.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				b.CreatedAt = value.Time
			}
		case battery.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				b.UpdatedAt = value.Time
			}
		case battery.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				b.DeletedAt = new(time.Time)
				*b.DeletedAt = value.Time
			}
		case battery.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &b.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case battery.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &b.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case battery.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				b.Remark = value.String
			}
		case battery.FieldCityID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field city_id", values[i])
			} else if value.Valid {
				b.CityID = new(uint64)
				*b.CityID = uint64(value.Int64)
			}
		case battery.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				b.RiderID = new(uint64)
				*b.RiderID = uint64(value.Int64)
			}
		case battery.FieldCabinetID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field cabinet_id", values[i])
			} else if value.Valid {
				b.CabinetID = new(uint64)
				*b.CabinetID = uint64(value.Int64)
			}
		case battery.FieldSubscribeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field subscribe_id", values[i])
			} else if value.Valid {
				b.SubscribeID = new(uint64)
				*b.SubscribeID = uint64(value.Int64)
			}
		case battery.FieldEnterpriseID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field enterprise_id", values[i])
			} else if value.Valid {
				b.EnterpriseID = new(uint64)
				*b.EnterpriseID = uint64(value.Int64)
			}
		case battery.FieldStationID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field station_id", values[i])
			} else if value.Valid {
				b.StationID = new(uint64)
				*b.StationID = uint64(value.Int64)
			}
		case battery.FieldSn:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sn", values[i])
			} else if value.Valid {
				b.Sn = value.String
			}
		case battery.FieldBrand:
			if value, ok := values[i].(*adapter.BatteryBrand); !ok {
				return fmt.Errorf("unexpected type %T for field brand", values[i])
			} else if value != nil {
				b.Brand = *value
			}
		case battery.FieldEnable:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field enable", values[i])
			} else if value.Valid {
				b.Enable = value.Bool
			}
		case battery.FieldModel:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field model", values[i])
			} else if value.Valid {
				b.Model = value.String
			}
		case battery.FieldOrdinal:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field ordinal", values[i])
			} else if value.Valid {
				b.Ordinal = new(int)
				*b.Ordinal = int(value.Int64)
			}
		default:
			b.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Battery.
// This includes values selected through modifiers, order, etc.
func (b *Battery) Value(name string) (ent.Value, error) {
	return b.selectValues.Get(name)
}

// QueryCity queries the "city" edge of the Battery entity.
func (b *Battery) QueryCity() *CityQuery {
	return NewBatteryClient(b.config).QueryCity(b)
}

// QueryFlows queries the "flows" edge of the Battery entity.
func (b *Battery) QueryFlows() *BatteryFlowQuery {
	return NewBatteryClient(b.config).QueryFlows(b)
}

// Update returns a builder for updating this Battery.
// Note that you need to call Battery.Unwrap() before calling this method if this Battery
// was returned from a transaction, and the transaction was committed or rolled back.
func (b *Battery) Update() *BatteryUpdateOne {
	return NewBatteryClient(b.config).UpdateOne(b)
}

// Unwrap unwraps the Battery entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (b *Battery) Unwrap() *Battery {
	_tx, ok := b.config.driver.(*txDriver)
	if !ok {
		panic("ent: Battery is not a transactional entity")
	}
	b.config.driver = _tx.drv
	return b
}

// String implements the fmt.Stringer.
func (b *Battery) String() string {
	var builder strings.Builder
	builder.WriteString("Battery(")
	builder.WriteString(fmt.Sprintf("id=%v, ", b.ID))
	builder.WriteString("created_at=")
	builder.WriteString(b.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(b.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := b.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", b.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", b.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(b.Remark)
	builder.WriteString(", ")
	if v := b.CityID; v != nil {
		builder.WriteString("city_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := b.RiderID; v != nil {
		builder.WriteString("rider_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := b.CabinetID; v != nil {
		builder.WriteString("cabinet_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := b.SubscribeID; v != nil {
		builder.WriteString("subscribe_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := b.EnterpriseID; v != nil {
		builder.WriteString("enterprise_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := b.StationID; v != nil {
		builder.WriteString("station_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("sn=")
	builder.WriteString(b.Sn)
	builder.WriteString(", ")
	builder.WriteString("brand=")
	builder.WriteString(fmt.Sprintf("%v", b.Brand))
	builder.WriteString(", ")
	builder.WriteString("enable=")
	builder.WriteString(fmt.Sprintf("%v", b.Enable))
	builder.WriteString(", ")
	builder.WriteString("model=")
	builder.WriteString(b.Model)
	builder.WriteString(", ")
	if v := b.Ordinal; v != nil {
		builder.WriteString("ordinal=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Batteries is a parsable slice of Battery.
type Batteries []*Battery
