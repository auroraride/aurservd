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
	"github.com/auroraride/aurservd/internal/ent/assistance"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
)

// Assistance is the model entity for the Assistance schema.
type Assistance struct {
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
	// 门店ID
	StoreID *uint64 `json:"store_id,omitempty"`
	// 骑手ID
	RiderID uint64 `json:"rider_id,omitempty"`
	// SubscribeID holds the value of the "subscribe_id" field.
	SubscribeID uint64 `json:"subscribe_id,omitempty"`
	// 城市ID
	CityID uint64 `json:"city_id,omitempty"`
	// 店员ID
	EmployeeID *uint64 `json:"employee_id,omitempty"`
	// 支付订单
	OrderID *uint64 `json:"order_id,omitempty"`
	// 救援状态 0:待分配 1:已接单/已分配 2:已拒绝 3:救援失败 4:救援成功待支付 5:救援成功已支付
	Status uint8 `json:"status,omitempty"`
	// 经度
	Lng float64 `json:"lng,omitempty"`
	// 纬度
	Lat float64 `json:"lat,omitempty"`
	// 详细地址
	Address string `json:"address,omitempty"`
	// 故障
	Breakdown string `json:"breakdown,omitempty"`
	// 故障描述
	BreakdownDesc string `json:"breakdown_desc,omitempty"`
	// 故障照片
	BreakdownPhotos []string `json:"breakdown_photos,omitempty"`
	// 取消原因
	CancelReason *string `json:"cancel_reason,omitempty"`
	// 取消原因详细描述
	CancelReasonDesc *string `json:"cancel_reason_desc,omitempty"`
	// 救援距离
	Distance float64 `json:"distance,omitempty"`
	// 救援原因
	Reason string `json:"reason,omitempty"`
	// 检测照片
	DetectPhoto string `json:"detect_photo,omitempty"`
	// 与用户合影
	JointPhoto string `json:"joint_photo,omitempty"`
	// 本次救援费用
	Cost float64 `json:"cost,omitempty"`
	// 拒绝原因
	RefusedDesc *string `json:"refused_desc,omitempty"`
	// 支付时间
	PayAt *time.Time `json:"pay_at,omitempty"`
	// 分配时间
	AllocateAt *time.Time `json:"allocate_at,omitempty"`
	// 分配等待时间(s)
	Wait int `json:"wait,omitempty"`
	// 免费理由
	FreeReason *string `json:"free_reason,omitempty"`
	// 失败原因
	FailReason *string `json:"fail_reason,omitempty"`
	// 救援处理时间
	ProcessAt *time.Time `json:"process_at,omitempty"`
	// 救援费用单价 元/公里
	Price float64 `json:"price,omitempty"`
	// 路径导航规划时间 (s)
	NaviDuration int `json:"navi_duration,omitempty"`
	// 路径导航规划坐标组
	NaviPolylines []string `json:"navi_polylines,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AssistanceQuery when eager-loading is set.
	Edges        AssistanceEdges `json:"edges"`
	selectValues sql.SelectValues
}

// AssistanceEdges holds the relations/edges for other nodes in the graph.
type AssistanceEdges struct {
	// Store holds the value of the store edge.
	Store *Store `json:"store,omitempty"`
	// 骑手
	Rider *Rider `json:"rider,omitempty"`
	// Subscribe holds the value of the subscribe edge.
	Subscribe *Subscribe `json:"subscribe,omitempty"`
	// City holds the value of the city edge.
	City *City `json:"city,omitempty"`
	// Order holds the value of the order edge.
	Order *Order `json:"order,omitempty"`
	// Employee holds the value of the employee edge.
	Employee *Employee `json:"employee,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [6]bool
}

// StoreOrErr returns the Store value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssistanceEdges) StoreOrErr() (*Store, error) {
	if e.Store != nil {
		return e.Store, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: store.Label}
	}
	return nil, &NotLoadedError{edge: "store"}
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssistanceEdges) RiderOrErr() (*Rider, error) {
	if e.Rider != nil {
		return e.Rider, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: rider.Label}
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// SubscribeOrErr returns the Subscribe value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssistanceEdges) SubscribeOrErr() (*Subscribe, error) {
	if e.Subscribe != nil {
		return e.Subscribe, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: subscribe.Label}
	}
	return nil, &NotLoadedError{edge: "subscribe"}
}

// CityOrErr returns the City value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssistanceEdges) CityOrErr() (*City, error) {
	if e.City != nil {
		return e.City, nil
	} else if e.loadedTypes[3] {
		return nil, &NotFoundError{label: city.Label}
	}
	return nil, &NotLoadedError{edge: "city"}
}

// OrderOrErr returns the Order value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssistanceEdges) OrderOrErr() (*Order, error) {
	if e.Order != nil {
		return e.Order, nil
	} else if e.loadedTypes[4] {
		return nil, &NotFoundError{label: order.Label}
	}
	return nil, &NotLoadedError{edge: "order"}
}

// EmployeeOrErr returns the Employee value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssistanceEdges) EmployeeOrErr() (*Employee, error) {
	if e.Employee != nil {
		return e.Employee, nil
	} else if e.loadedTypes[5] {
		return nil, &NotFoundError{label: employee.Label}
	}
	return nil, &NotLoadedError{edge: "employee"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Assistance) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case assistance.FieldCreator, assistance.FieldLastModifier, assistance.FieldBreakdownPhotos, assistance.FieldNaviPolylines:
			values[i] = new([]byte)
		case assistance.FieldLng, assistance.FieldLat, assistance.FieldDistance, assistance.FieldCost, assistance.FieldPrice:
			values[i] = new(sql.NullFloat64)
		case assistance.FieldID, assistance.FieldStoreID, assistance.FieldRiderID, assistance.FieldSubscribeID, assistance.FieldCityID, assistance.FieldEmployeeID, assistance.FieldOrderID, assistance.FieldStatus, assistance.FieldWait, assistance.FieldNaviDuration:
			values[i] = new(sql.NullInt64)
		case assistance.FieldRemark, assistance.FieldAddress, assistance.FieldBreakdown, assistance.FieldBreakdownDesc, assistance.FieldCancelReason, assistance.FieldCancelReasonDesc, assistance.FieldReason, assistance.FieldDetectPhoto, assistance.FieldJointPhoto, assistance.FieldRefusedDesc, assistance.FieldFreeReason, assistance.FieldFailReason:
			values[i] = new(sql.NullString)
		case assistance.FieldCreatedAt, assistance.FieldUpdatedAt, assistance.FieldDeletedAt, assistance.FieldPayAt, assistance.FieldAllocateAt, assistance.FieldProcessAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Assistance fields.
func (a *Assistance) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case assistance.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			a.ID = uint64(value.Int64)
		case assistance.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				a.CreatedAt = value.Time
			}
		case assistance.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				a.UpdatedAt = value.Time
			}
		case assistance.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				a.DeletedAt = new(time.Time)
				*a.DeletedAt = value.Time
			}
		case assistance.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &a.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case assistance.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &a.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case assistance.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				a.Remark = value.String
			}
		case assistance.FieldStoreID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field store_id", values[i])
			} else if value.Valid {
				a.StoreID = new(uint64)
				*a.StoreID = uint64(value.Int64)
			}
		case assistance.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				a.RiderID = uint64(value.Int64)
			}
		case assistance.FieldSubscribeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field subscribe_id", values[i])
			} else if value.Valid {
				a.SubscribeID = uint64(value.Int64)
			}
		case assistance.FieldCityID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field city_id", values[i])
			} else if value.Valid {
				a.CityID = uint64(value.Int64)
			}
		case assistance.FieldEmployeeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field employee_id", values[i])
			} else if value.Valid {
				a.EmployeeID = new(uint64)
				*a.EmployeeID = uint64(value.Int64)
			}
		case assistance.FieldOrderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field order_id", values[i])
			} else if value.Valid {
				a.OrderID = new(uint64)
				*a.OrderID = uint64(value.Int64)
			}
		case assistance.FieldStatus:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				a.Status = uint8(value.Int64)
			}
		case assistance.FieldLng:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field lng", values[i])
			} else if value.Valid {
				a.Lng = value.Float64
			}
		case assistance.FieldLat:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field lat", values[i])
			} else if value.Valid {
				a.Lat = value.Float64
			}
		case assistance.FieldAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field address", values[i])
			} else if value.Valid {
				a.Address = value.String
			}
		case assistance.FieldBreakdown:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field breakdown", values[i])
			} else if value.Valid {
				a.Breakdown = value.String
			}
		case assistance.FieldBreakdownDesc:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field breakdown_desc", values[i])
			} else if value.Valid {
				a.BreakdownDesc = value.String
			}
		case assistance.FieldBreakdownPhotos:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field breakdown_photos", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &a.BreakdownPhotos); err != nil {
					return fmt.Errorf("unmarshal field breakdown_photos: %w", err)
				}
			}
		case assistance.FieldCancelReason:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field cancel_reason", values[i])
			} else if value.Valid {
				a.CancelReason = new(string)
				*a.CancelReason = value.String
			}
		case assistance.FieldCancelReasonDesc:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field cancel_reason_desc", values[i])
			} else if value.Valid {
				a.CancelReasonDesc = new(string)
				*a.CancelReasonDesc = value.String
			}
		case assistance.FieldDistance:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field distance", values[i])
			} else if value.Valid {
				a.Distance = value.Float64
			}
		case assistance.FieldReason:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field reason", values[i])
			} else if value.Valid {
				a.Reason = value.String
			}
		case assistance.FieldDetectPhoto:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field detect_photo", values[i])
			} else if value.Valid {
				a.DetectPhoto = value.String
			}
		case assistance.FieldJointPhoto:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field joint_photo", values[i])
			} else if value.Valid {
				a.JointPhoto = value.String
			}
		case assistance.FieldCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field cost", values[i])
			} else if value.Valid {
				a.Cost = value.Float64
			}
		case assistance.FieldRefusedDesc:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field refused_desc", values[i])
			} else if value.Valid {
				a.RefusedDesc = new(string)
				*a.RefusedDesc = value.String
			}
		case assistance.FieldPayAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field pay_at", values[i])
			} else if value.Valid {
				a.PayAt = new(time.Time)
				*a.PayAt = value.Time
			}
		case assistance.FieldAllocateAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field allocate_at", values[i])
			} else if value.Valid {
				a.AllocateAt = new(time.Time)
				*a.AllocateAt = value.Time
			}
		case assistance.FieldWait:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field wait", values[i])
			} else if value.Valid {
				a.Wait = int(value.Int64)
			}
		case assistance.FieldFreeReason:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field free_reason", values[i])
			} else if value.Valid {
				a.FreeReason = new(string)
				*a.FreeReason = value.String
			}
		case assistance.FieldFailReason:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field fail_reason", values[i])
			} else if value.Valid {
				a.FailReason = new(string)
				*a.FailReason = value.String
			}
		case assistance.FieldProcessAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field process_at", values[i])
			} else if value.Valid {
				a.ProcessAt = new(time.Time)
				*a.ProcessAt = value.Time
			}
		case assistance.FieldPrice:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field price", values[i])
			} else if value.Valid {
				a.Price = value.Float64
			}
		case assistance.FieldNaviDuration:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field navi_duration", values[i])
			} else if value.Valid {
				a.NaviDuration = int(value.Int64)
			}
		case assistance.FieldNaviPolylines:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field navi_polylines", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &a.NaviPolylines); err != nil {
					return fmt.Errorf("unmarshal field navi_polylines: %w", err)
				}
			}
		default:
			a.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Assistance.
// This includes values selected through modifiers, order, etc.
func (a *Assistance) Value(name string) (ent.Value, error) {
	return a.selectValues.Get(name)
}

// QueryStore queries the "store" edge of the Assistance entity.
func (a *Assistance) QueryStore() *StoreQuery {
	return NewAssistanceClient(a.config).QueryStore(a)
}

// QueryRider queries the "rider" edge of the Assistance entity.
func (a *Assistance) QueryRider() *RiderQuery {
	return NewAssistanceClient(a.config).QueryRider(a)
}

// QuerySubscribe queries the "subscribe" edge of the Assistance entity.
func (a *Assistance) QuerySubscribe() *SubscribeQuery {
	return NewAssistanceClient(a.config).QuerySubscribe(a)
}

// QueryCity queries the "city" edge of the Assistance entity.
func (a *Assistance) QueryCity() *CityQuery {
	return NewAssistanceClient(a.config).QueryCity(a)
}

// QueryOrder queries the "order" edge of the Assistance entity.
func (a *Assistance) QueryOrder() *OrderQuery {
	return NewAssistanceClient(a.config).QueryOrder(a)
}

// QueryEmployee queries the "employee" edge of the Assistance entity.
func (a *Assistance) QueryEmployee() *EmployeeQuery {
	return NewAssistanceClient(a.config).QueryEmployee(a)
}

// Update returns a builder for updating this Assistance.
// Note that you need to call Assistance.Unwrap() before calling this method if this Assistance
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Assistance) Update() *AssistanceUpdateOne {
	return NewAssistanceClient(a.config).UpdateOne(a)
}

// Unwrap unwraps the Assistance entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Assistance) Unwrap() *Assistance {
	_tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Assistance is not a transactional entity")
	}
	a.config.driver = _tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Assistance) String() string {
	var builder strings.Builder
	builder.WriteString("Assistance(")
	builder.WriteString(fmt.Sprintf("id=%v, ", a.ID))
	builder.WriteString("created_at=")
	builder.WriteString(a.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(a.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := a.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", a.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", a.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(a.Remark)
	builder.WriteString(", ")
	if v := a.StoreID; v != nil {
		builder.WriteString("store_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("rider_id=")
	builder.WriteString(fmt.Sprintf("%v", a.RiderID))
	builder.WriteString(", ")
	builder.WriteString("subscribe_id=")
	builder.WriteString(fmt.Sprintf("%v", a.SubscribeID))
	builder.WriteString(", ")
	builder.WriteString("city_id=")
	builder.WriteString(fmt.Sprintf("%v", a.CityID))
	builder.WriteString(", ")
	if v := a.EmployeeID; v != nil {
		builder.WriteString("employee_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := a.OrderID; v != nil {
		builder.WriteString("order_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", a.Status))
	builder.WriteString(", ")
	builder.WriteString("lng=")
	builder.WriteString(fmt.Sprintf("%v", a.Lng))
	builder.WriteString(", ")
	builder.WriteString("lat=")
	builder.WriteString(fmt.Sprintf("%v", a.Lat))
	builder.WriteString(", ")
	builder.WriteString("address=")
	builder.WriteString(a.Address)
	builder.WriteString(", ")
	builder.WriteString("breakdown=")
	builder.WriteString(a.Breakdown)
	builder.WriteString(", ")
	builder.WriteString("breakdown_desc=")
	builder.WriteString(a.BreakdownDesc)
	builder.WriteString(", ")
	builder.WriteString("breakdown_photos=")
	builder.WriteString(fmt.Sprintf("%v", a.BreakdownPhotos))
	builder.WriteString(", ")
	if v := a.CancelReason; v != nil {
		builder.WriteString("cancel_reason=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := a.CancelReasonDesc; v != nil {
		builder.WriteString("cancel_reason_desc=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("distance=")
	builder.WriteString(fmt.Sprintf("%v", a.Distance))
	builder.WriteString(", ")
	builder.WriteString("reason=")
	builder.WriteString(a.Reason)
	builder.WriteString(", ")
	builder.WriteString("detect_photo=")
	builder.WriteString(a.DetectPhoto)
	builder.WriteString(", ")
	builder.WriteString("joint_photo=")
	builder.WriteString(a.JointPhoto)
	builder.WriteString(", ")
	builder.WriteString("cost=")
	builder.WriteString(fmt.Sprintf("%v", a.Cost))
	builder.WriteString(", ")
	if v := a.RefusedDesc; v != nil {
		builder.WriteString("refused_desc=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := a.PayAt; v != nil {
		builder.WriteString("pay_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := a.AllocateAt; v != nil {
		builder.WriteString("allocate_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("wait=")
	builder.WriteString(fmt.Sprintf("%v", a.Wait))
	builder.WriteString(", ")
	if v := a.FreeReason; v != nil {
		builder.WriteString("free_reason=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := a.FailReason; v != nil {
		builder.WriteString("fail_reason=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := a.ProcessAt; v != nil {
		builder.WriteString("process_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("price=")
	builder.WriteString(fmt.Sprintf("%v", a.Price))
	builder.WriteString(", ")
	builder.WriteString("navi_duration=")
	builder.WriteString(fmt.Sprintf("%v", a.NaviDuration))
	builder.WriteString(", ")
	builder.WriteString("navi_polylines=")
	builder.WriteString(fmt.Sprintf("%v", a.NaviPolylines))
	builder.WriteByte(')')
	return builder.String()
}

// Assistances is a parsable slice of Assistance.
type Assistances []*Assistance
