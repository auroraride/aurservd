// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribepause"
	"github.com/auroraride/aurservd/internal/ent/subscribesuspend"
)

// SubscribeSuspend is the model entity for the SubscribeSuspend schema.
type SubscribeSuspend struct {
	config `json:"-"`
	// ID of the ent.
	ID uint64 `json:"id,omitempty"`
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
	// RiderID holds the value of the "rider_id" field.
	// 骑手ID
	RiderID uint64 `json:"rider_id,omitempty"`
	// SubscribeID holds the value of the "subscribe_id" field.
	// 订阅ID
	SubscribeID uint64 `json:"subscribe_id,omitempty"`
	// PauseID holds the value of the "pause_id" field.
	// 寄存ID
	PauseID uint64 `json:"pause_id,omitempty"`
	// Days holds the value of the "days" field.
	// 暂停天数
	Days int `json:"days,omitempty"`
	// StartAt holds the value of the "start_at" field.
	// 开始时间
	StartAt time.Time `json:"start_at,omitempty"`
	// EndAt holds the value of the "end_at" field.
	// 结束时间
	EndAt time.Time `json:"end_at,omitempty"`
	// EndReason holds the value of the "end_reason" field.
	// 结束理由
	EndReason string `json:"end_reason,omitempty"`
	// EndModifier holds the value of the "end_modifier" field.
	// 继续计费管理员信息
	EndModifier *model.Modifier `json:"end_modifier,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SubscribeSuspendQuery when eager-loading is set.
	Edges SubscribeSuspendEdges `json:"edges"`
}

// SubscribeSuspendEdges holds the relations/edges for other nodes in the graph.
type SubscribeSuspendEdges struct {
	// City holds the value of the city edge.
	City *City `json:"city,omitempty"`
	// Rider holds the value of the rider edge.
	Rider *Rider `json:"rider,omitempty"`
	// Subscribe holds the value of the subscribe edge.
	Subscribe *Subscribe `json:"subscribe,omitempty"`
	// Pause holds the value of the pause edge.
	Pause *SubscribePause `json:"pause,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// CityOrErr returns the City value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SubscribeSuspendEdges) CityOrErr() (*City, error) {
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

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SubscribeSuspendEdges) RiderOrErr() (*Rider, error) {
	if e.loadedTypes[1] {
		if e.Rider == nil {
			// The edge rider was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: rider.Label}
		}
		return e.Rider, nil
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// SubscribeOrErr returns the Subscribe value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SubscribeSuspendEdges) SubscribeOrErr() (*Subscribe, error) {
	if e.loadedTypes[2] {
		if e.Subscribe == nil {
			// The edge subscribe was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: subscribe.Label}
		}
		return e.Subscribe, nil
	}
	return nil, &NotLoadedError{edge: "subscribe"}
}

// PauseOrErr returns the Pause value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SubscribeSuspendEdges) PauseOrErr() (*SubscribePause, error) {
	if e.loadedTypes[3] {
		if e.Pause == nil {
			// The edge pause was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: subscribepause.Label}
		}
		return e.Pause, nil
	}
	return nil, &NotLoadedError{edge: "pause"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*SubscribeSuspend) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case subscribesuspend.FieldCreator, subscribesuspend.FieldLastModifier, subscribesuspend.FieldEndModifier:
			values[i] = new([]byte)
		case subscribesuspend.FieldID, subscribesuspend.FieldCityID, subscribesuspend.FieldRiderID, subscribesuspend.FieldSubscribeID, subscribesuspend.FieldPauseID, subscribesuspend.FieldDays:
			values[i] = new(sql.NullInt64)
		case subscribesuspend.FieldRemark, subscribesuspend.FieldEndReason:
			values[i] = new(sql.NullString)
		case subscribesuspend.FieldStartAt, subscribesuspend.FieldEndAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type SubscribeSuspend", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SubscribeSuspend fields.
func (ss *SubscribeSuspend) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case subscribesuspend.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ss.ID = uint64(value.Int64)
		case subscribesuspend.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ss.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case subscribesuspend.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ss.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case subscribesuspend.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				ss.Remark = value.String
			}
		case subscribesuspend.FieldCityID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field city_id", values[i])
			} else if value.Valid {
				ss.CityID = uint64(value.Int64)
			}
		case subscribesuspend.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				ss.RiderID = uint64(value.Int64)
			}
		case subscribesuspend.FieldSubscribeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field subscribe_id", values[i])
			} else if value.Valid {
				ss.SubscribeID = uint64(value.Int64)
			}
		case subscribesuspend.FieldPauseID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field pause_id", values[i])
			} else if value.Valid {
				ss.PauseID = uint64(value.Int64)
			}
		case subscribesuspend.FieldDays:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field days", values[i])
			} else if value.Valid {
				ss.Days = int(value.Int64)
			}
		case subscribesuspend.FieldStartAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field start_at", values[i])
			} else if value.Valid {
				ss.StartAt = value.Time
			}
		case subscribesuspend.FieldEndAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field end_at", values[i])
			} else if value.Valid {
				ss.EndAt = value.Time
			}
		case subscribesuspend.FieldEndReason:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field end_reason", values[i])
			} else if value.Valid {
				ss.EndReason = value.String
			}
		case subscribesuspend.FieldEndModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field end_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ss.EndModifier); err != nil {
					return fmt.Errorf("unmarshal field end_modifier: %w", err)
				}
			}
		}
	}
	return nil
}

// QueryCity queries the "city" edge of the SubscribeSuspend entity.
func (ss *SubscribeSuspend) QueryCity() *CityQuery {
	return (&SubscribeSuspendClient{config: ss.config}).QueryCity(ss)
}

// QueryRider queries the "rider" edge of the SubscribeSuspend entity.
func (ss *SubscribeSuspend) QueryRider() *RiderQuery {
	return (&SubscribeSuspendClient{config: ss.config}).QueryRider(ss)
}

// QuerySubscribe queries the "subscribe" edge of the SubscribeSuspend entity.
func (ss *SubscribeSuspend) QuerySubscribe() *SubscribeQuery {
	return (&SubscribeSuspendClient{config: ss.config}).QuerySubscribe(ss)
}

// QueryPause queries the "pause" edge of the SubscribeSuspend entity.
func (ss *SubscribeSuspend) QueryPause() *SubscribePauseQuery {
	return (&SubscribeSuspendClient{config: ss.config}).QueryPause(ss)
}

// Update returns a builder for updating this SubscribeSuspend.
// Note that you need to call SubscribeSuspend.Unwrap() before calling this method if this SubscribeSuspend
// was returned from a transaction, and the transaction was committed or rolled back.
func (ss *SubscribeSuspend) Update() *SubscribeSuspendUpdateOne {
	return (&SubscribeSuspendClient{config: ss.config}).UpdateOne(ss)
}

// Unwrap unwraps the SubscribeSuspend entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ss *SubscribeSuspend) Unwrap() *SubscribeSuspend {
	tx, ok := ss.config.driver.(*txDriver)
	if !ok {
		panic("ent: SubscribeSuspend is not a transactional entity")
	}
	ss.config.driver = tx.drv
	return ss
}

// String implements the fmt.Stringer.
func (ss *SubscribeSuspend) String() string {
	var builder strings.Builder
	builder.WriteString("SubscribeSuspend(")
	builder.WriteString(fmt.Sprintf("id=%v", ss.ID))
	builder.WriteString(", creator=")
	builder.WriteString(fmt.Sprintf("%v", ss.Creator))
	builder.WriteString(", last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", ss.LastModifier))
	builder.WriteString(", remark=")
	builder.WriteString(ss.Remark)
	builder.WriteString(", city_id=")
	builder.WriteString(fmt.Sprintf("%v", ss.CityID))
	builder.WriteString(", rider_id=")
	builder.WriteString(fmt.Sprintf("%v", ss.RiderID))
	builder.WriteString(", subscribe_id=")
	builder.WriteString(fmt.Sprintf("%v", ss.SubscribeID))
	builder.WriteString(", pause_id=")
	builder.WriteString(fmt.Sprintf("%v", ss.PauseID))
	builder.WriteString(", days=")
	builder.WriteString(fmt.Sprintf("%v", ss.Days))
	builder.WriteString(", start_at=")
	builder.WriteString(ss.StartAt.Format(time.ANSIC))
	builder.WriteString(", end_at=")
	builder.WriteString(ss.EndAt.Format(time.ANSIC))
	builder.WriteString(", end_reason=")
	builder.WriteString(ss.EndReason)
	builder.WriteString(", end_modifier=")
	builder.WriteString(fmt.Sprintf("%v", ss.EndModifier))
	builder.WriteByte(')')
	return builder.String()
}

// SubscribeSuspends is a parsable slice of SubscribeSuspend.
type SubscribeSuspends []*SubscribeSuspend

func (ss SubscribeSuspends) config(cfg config) {
	for _i := range ss {
		ss[_i].config = cfg
	}
}