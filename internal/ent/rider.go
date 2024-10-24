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
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/person"
	"github.com/auroraride/aurservd/internal/ent/rider"
)

// Rider is the model entity for the Rider schema.
type Rider struct {
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
	// 站点ID
	StationID *uint64 `json:"station_id,omitempty"`
	// 身份
	PersonID *uint64 `json:"person_id,omitempty"`
	// 骑手姓名
	Name string `json:"name,omitempty"`
	// 身份证号
	IDCardNumber string `json:"id_card_number,omitempty"`
	// 所属企业
	EnterpriseID *uint64 `json:"enterprise_id,omitempty"`
	// 手机号
	Phone string `json:"phone,omitempty"`
	// 紧急联系人
	Contact *model.RiderContact `json:"contact,omitempty"`
	// 登录设备类型: 1iOS 2Android
	DeviceType uint8 `json:"device_type,omitempty"`
	// 最近登录设备
	LastDevice string `json:"last_device,omitempty"`
	// 是否新设备
	IsNewDevice bool `json:"is_new_device,omitempty"`
	// 推送ID
	PushID string `json:"push_id,omitempty"`
	// 最后登录时间
	LastSigninAt *time.Time `json:"last_signin_at,omitempty"`
	// 是否封禁骑手账号
	Blocked bool `json:"blocked,omitempty"`
	// 骑手积分
	Points int64 `json:"points,omitempty"`
	// 换电间隔配置
	ExchangeLimit model.RiderExchangeLimit `json:"exchange_limit,omitempty"`
	// 换电频次配置
	ExchangeFrequency model.RiderExchangeFrequency `json:"exchange_frequency,omitempty"`
	// 加入团签时间
	JoinEnterpriseAt *time.Time `json:"join_enterprise_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the RiderQuery when eager-loading is set.
	Edges        RiderEdges `json:"edges"`
	selectValues sql.SelectValues
}

// RiderEdges holds the relations/edges for other nodes in the graph.
type RiderEdges struct {
	// Station holds the value of the station edge.
	Station *EnterpriseStation `json:"station,omitempty"`
	// Person holds the value of the person edge.
	Person *Person `json:"person,omitempty"`
	// Enterprise holds the value of the enterprise edge.
	Enterprise *Enterprise `json:"enterprise,omitempty"`
	// Contracts holds the value of the contracts edge.
	Contracts []*Contract `json:"contracts,omitempty"`
	// Faults holds the value of the faults edge.
	Faults []*CabinetFault `json:"faults,omitempty"`
	// Orders holds the value of the orders edge.
	Orders []*Order `json:"orders,omitempty"`
	// 换电记录
	Exchanges []*Exchange `json:"exchanges,omitempty"`
	// 订阅
	Subscribes []*Subscribe `json:"subscribes,omitempty"`
	// Asset holds the value of the asset edge.
	Asset []*Asset `json:"asset,omitempty"`
	// Stocks holds the value of the stocks edge.
	Stocks []*Stock `json:"stocks,omitempty"`
	// Followups holds the value of the followups edge.
	Followups []*RiderFollowUp `json:"followups,omitempty"`
	// Battery holds the value of the battery edge.
	Battery *Asset `json:"battery,omitempty"`
	// BatteryFlows holds the value of the battery_flows edge.
	BatteryFlows []*BatteryFlow `json:"battery_flows,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [13]bool
}

// StationOrErr returns the Station value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RiderEdges) StationOrErr() (*EnterpriseStation, error) {
	if e.Station != nil {
		return e.Station, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: enterprisestation.Label}
	}
	return nil, &NotLoadedError{edge: "station"}
}

// PersonOrErr returns the Person value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RiderEdges) PersonOrErr() (*Person, error) {
	if e.Person != nil {
		return e.Person, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: person.Label}
	}
	return nil, &NotLoadedError{edge: "person"}
}

// EnterpriseOrErr returns the Enterprise value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RiderEdges) EnterpriseOrErr() (*Enterprise, error) {
	if e.Enterprise != nil {
		return e.Enterprise, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: enterprise.Label}
	}
	return nil, &NotLoadedError{edge: "enterprise"}
}

// ContractsOrErr returns the Contracts value or an error if the edge
// was not loaded in eager-loading.
func (e RiderEdges) ContractsOrErr() ([]*Contract, error) {
	if e.loadedTypes[3] {
		return e.Contracts, nil
	}
	return nil, &NotLoadedError{edge: "contracts"}
}

// FaultsOrErr returns the Faults value or an error if the edge
// was not loaded in eager-loading.
func (e RiderEdges) FaultsOrErr() ([]*CabinetFault, error) {
	if e.loadedTypes[4] {
		return e.Faults, nil
	}
	return nil, &NotLoadedError{edge: "faults"}
}

// OrdersOrErr returns the Orders value or an error if the edge
// was not loaded in eager-loading.
func (e RiderEdges) OrdersOrErr() ([]*Order, error) {
	if e.loadedTypes[5] {
		return e.Orders, nil
	}
	return nil, &NotLoadedError{edge: "orders"}
}

// ExchangesOrErr returns the Exchanges value or an error if the edge
// was not loaded in eager-loading.
func (e RiderEdges) ExchangesOrErr() ([]*Exchange, error) {
	if e.loadedTypes[6] {
		return e.Exchanges, nil
	}
	return nil, &NotLoadedError{edge: "exchanges"}
}

// SubscribesOrErr returns the Subscribes value or an error if the edge
// was not loaded in eager-loading.
func (e RiderEdges) SubscribesOrErr() ([]*Subscribe, error) {
	if e.loadedTypes[7] {
		return e.Subscribes, nil
	}
	return nil, &NotLoadedError{edge: "subscribes"}
}

// AssetOrErr returns the Asset value or an error if the edge
// was not loaded in eager-loading.
func (e RiderEdges) AssetOrErr() ([]*Asset, error) {
	if e.loadedTypes[8] {
		return e.Asset, nil
	}
	return nil, &NotLoadedError{edge: "asset"}
}

// StocksOrErr returns the Stocks value or an error if the edge
// was not loaded in eager-loading.
func (e RiderEdges) StocksOrErr() ([]*Stock, error) {
	if e.loadedTypes[9] {
		return e.Stocks, nil
	}
	return nil, &NotLoadedError{edge: "stocks"}
}

// FollowupsOrErr returns the Followups value or an error if the edge
// was not loaded in eager-loading.
func (e RiderEdges) FollowupsOrErr() ([]*RiderFollowUp, error) {
	if e.loadedTypes[10] {
		return e.Followups, nil
	}
	return nil, &NotLoadedError{edge: "followups"}
}

// BatteryOrErr returns the Battery value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RiderEdges) BatteryOrErr() (*Asset, error) {
	if e.Battery != nil {
		return e.Battery, nil
	} else if e.loadedTypes[11] {
		return nil, &NotFoundError{label: asset.Label}
	}
	return nil, &NotLoadedError{edge: "battery"}
}

// BatteryFlowsOrErr returns the BatteryFlows value or an error if the edge
// was not loaded in eager-loading.
func (e RiderEdges) BatteryFlowsOrErr() ([]*BatteryFlow, error) {
	if e.loadedTypes[12] {
		return e.BatteryFlows, nil
	}
	return nil, &NotLoadedError{edge: "battery_flows"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Rider) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case rider.FieldCreator, rider.FieldLastModifier, rider.FieldContact, rider.FieldExchangeLimit, rider.FieldExchangeFrequency:
			values[i] = new([]byte)
		case rider.FieldIsNewDevice, rider.FieldBlocked:
			values[i] = new(sql.NullBool)
		case rider.FieldID, rider.FieldStationID, rider.FieldPersonID, rider.FieldEnterpriseID, rider.FieldDeviceType, rider.FieldPoints:
			values[i] = new(sql.NullInt64)
		case rider.FieldRemark, rider.FieldName, rider.FieldIDCardNumber, rider.FieldPhone, rider.FieldLastDevice, rider.FieldPushID:
			values[i] = new(sql.NullString)
		case rider.FieldCreatedAt, rider.FieldUpdatedAt, rider.FieldDeletedAt, rider.FieldLastSigninAt, rider.FieldJoinEnterpriseAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Rider fields.
func (r *Rider) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case rider.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			r.ID = uint64(value.Int64)
		case rider.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				r.CreatedAt = value.Time
			}
		case rider.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				r.UpdatedAt = value.Time
			}
		case rider.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				r.DeletedAt = new(time.Time)
				*r.DeletedAt = value.Time
			}
		case rider.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &r.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case rider.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &r.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case rider.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				r.Remark = value.String
			}
		case rider.FieldStationID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field station_id", values[i])
			} else if value.Valid {
				r.StationID = new(uint64)
				*r.StationID = uint64(value.Int64)
			}
		case rider.FieldPersonID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field person_id", values[i])
			} else if value.Valid {
				r.PersonID = new(uint64)
				*r.PersonID = uint64(value.Int64)
			}
		case rider.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				r.Name = value.String
			}
		case rider.FieldIDCardNumber:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id_card_number", values[i])
			} else if value.Valid {
				r.IDCardNumber = value.String
			}
		case rider.FieldEnterpriseID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field enterprise_id", values[i])
			} else if value.Valid {
				r.EnterpriseID = new(uint64)
				*r.EnterpriseID = uint64(value.Int64)
			}
		case rider.FieldPhone:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field phone", values[i])
			} else if value.Valid {
				r.Phone = value.String
			}
		case rider.FieldContact:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field contact", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &r.Contact); err != nil {
					return fmt.Errorf("unmarshal field contact: %w", err)
				}
			}
		case rider.FieldDeviceType:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field device_type", values[i])
			} else if value.Valid {
				r.DeviceType = uint8(value.Int64)
			}
		case rider.FieldLastDevice:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field last_device", values[i])
			} else if value.Valid {
				r.LastDevice = value.String
			}
		case rider.FieldIsNewDevice:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_new_device", values[i])
			} else if value.Valid {
				r.IsNewDevice = value.Bool
			}
		case rider.FieldPushID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field push_id", values[i])
			} else if value.Valid {
				r.PushID = value.String
			}
		case rider.FieldLastSigninAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field last_signin_at", values[i])
			} else if value.Valid {
				r.LastSigninAt = new(time.Time)
				*r.LastSigninAt = value.Time
			}
		case rider.FieldBlocked:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field blocked", values[i])
			} else if value.Valid {
				r.Blocked = value.Bool
			}
		case rider.FieldPoints:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field points", values[i])
			} else if value.Valid {
				r.Points = value.Int64
			}
		case rider.FieldExchangeLimit:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field exchange_limit", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &r.ExchangeLimit); err != nil {
					return fmt.Errorf("unmarshal field exchange_limit: %w", err)
				}
			}
		case rider.FieldExchangeFrequency:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field exchange_frequency", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &r.ExchangeFrequency); err != nil {
					return fmt.Errorf("unmarshal field exchange_frequency: %w", err)
				}
			}
		case rider.FieldJoinEnterpriseAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field join_enterprise_at", values[i])
			} else if value.Valid {
				r.JoinEnterpriseAt = new(time.Time)
				*r.JoinEnterpriseAt = value.Time
			}
		default:
			r.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Rider.
// This includes values selected through modifiers, order, etc.
func (r *Rider) Value(name string) (ent.Value, error) {
	return r.selectValues.Get(name)
}

// QueryStation queries the "station" edge of the Rider entity.
func (r *Rider) QueryStation() *EnterpriseStationQuery {
	return NewRiderClient(r.config).QueryStation(r)
}

// QueryPerson queries the "person" edge of the Rider entity.
func (r *Rider) QueryPerson() *PersonQuery {
	return NewRiderClient(r.config).QueryPerson(r)
}

// QueryEnterprise queries the "enterprise" edge of the Rider entity.
func (r *Rider) QueryEnterprise() *EnterpriseQuery {
	return NewRiderClient(r.config).QueryEnterprise(r)
}

// QueryContracts queries the "contracts" edge of the Rider entity.
func (r *Rider) QueryContracts() *ContractQuery {
	return NewRiderClient(r.config).QueryContracts(r)
}

// QueryFaults queries the "faults" edge of the Rider entity.
func (r *Rider) QueryFaults() *CabinetFaultQuery {
	return NewRiderClient(r.config).QueryFaults(r)
}

// QueryOrders queries the "orders" edge of the Rider entity.
func (r *Rider) QueryOrders() *OrderQuery {
	return NewRiderClient(r.config).QueryOrders(r)
}

// QueryExchanges queries the "exchanges" edge of the Rider entity.
func (r *Rider) QueryExchanges() *ExchangeQuery {
	return NewRiderClient(r.config).QueryExchanges(r)
}

// QuerySubscribes queries the "subscribes" edge of the Rider entity.
func (r *Rider) QuerySubscribes() *SubscribeQuery {
	return NewRiderClient(r.config).QuerySubscribes(r)
}

// QueryAsset queries the "asset" edge of the Rider entity.
func (r *Rider) QueryAsset() *AssetQuery {
	return NewRiderClient(r.config).QueryAsset(r)
}

// QueryStocks queries the "stocks" edge of the Rider entity.
func (r *Rider) QueryStocks() *StockQuery {
	return NewRiderClient(r.config).QueryStocks(r)
}

// QueryFollowups queries the "followups" edge of the Rider entity.
func (r *Rider) QueryFollowups() *RiderFollowUpQuery {
	return NewRiderClient(r.config).QueryFollowups(r)
}

// QueryBattery queries the "battery" edge of the Rider entity.
func (r *Rider) QueryBattery() *AssetQuery {
	return NewRiderClient(r.config).QueryBattery(r)
}

// QueryBatteryFlows queries the "battery_flows" edge of the Rider entity.
func (r *Rider) QueryBatteryFlows() *BatteryFlowQuery {
	return NewRiderClient(r.config).QueryBatteryFlows(r)
}

// Update returns a builder for updating this Rider.
// Note that you need to call Rider.Unwrap() before calling this method if this Rider
// was returned from a transaction, and the transaction was committed or rolled back.
func (r *Rider) Update() *RiderUpdateOne {
	return NewRiderClient(r.config).UpdateOne(r)
}

// Unwrap unwraps the Rider entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (r *Rider) Unwrap() *Rider {
	_tx, ok := r.config.driver.(*txDriver)
	if !ok {
		panic("ent: Rider is not a transactional entity")
	}
	r.config.driver = _tx.drv
	return r
}

// String implements the fmt.Stringer.
func (r *Rider) String() string {
	var builder strings.Builder
	builder.WriteString("Rider(")
	builder.WriteString(fmt.Sprintf("id=%v, ", r.ID))
	builder.WriteString("created_at=")
	builder.WriteString(r.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(r.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := r.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", r.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", r.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(r.Remark)
	builder.WriteString(", ")
	if v := r.StationID; v != nil {
		builder.WriteString("station_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := r.PersonID; v != nil {
		builder.WriteString("person_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(r.Name)
	builder.WriteString(", ")
	builder.WriteString("id_card_number=")
	builder.WriteString(r.IDCardNumber)
	builder.WriteString(", ")
	if v := r.EnterpriseID; v != nil {
		builder.WriteString("enterprise_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("phone=")
	builder.WriteString(r.Phone)
	builder.WriteString(", ")
	builder.WriteString("contact=")
	builder.WriteString(fmt.Sprintf("%v", r.Contact))
	builder.WriteString(", ")
	builder.WriteString("device_type=")
	builder.WriteString(fmt.Sprintf("%v", r.DeviceType))
	builder.WriteString(", ")
	builder.WriteString("last_device=")
	builder.WriteString(r.LastDevice)
	builder.WriteString(", ")
	builder.WriteString("is_new_device=")
	builder.WriteString(fmt.Sprintf("%v", r.IsNewDevice))
	builder.WriteString(", ")
	builder.WriteString("push_id=")
	builder.WriteString(r.PushID)
	builder.WriteString(", ")
	if v := r.LastSigninAt; v != nil {
		builder.WriteString("last_signin_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("blocked=")
	builder.WriteString(fmt.Sprintf("%v", r.Blocked))
	builder.WriteString(", ")
	builder.WriteString("points=")
	builder.WriteString(fmt.Sprintf("%v", r.Points))
	builder.WriteString(", ")
	builder.WriteString("exchange_limit=")
	builder.WriteString(fmt.Sprintf("%v", r.ExchangeLimit))
	builder.WriteString(", ")
	builder.WriteString("exchange_frequency=")
	builder.WriteString(fmt.Sprintf("%v", r.ExchangeFrequency))
	builder.WriteString(", ")
	if v := r.JoinEnterpriseAt; v != nil {
		builder.WriteString("join_enterprise_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Riders is a parsable slice of Rider.
type Riders []*Rider
