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
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/manager"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribealter"
)

// SubscribeAlter is the model entity for the SubscribeAlter schema.
type SubscribeAlter struct {
	config `json:"-"`
	// ID of the ent.
	ID uint64 `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// 创建人
	Creator *model.Modifier `json:"creator,omitempty"`
	// 最后修改人
	LastModifier *model.Modifier `json:"last_modifier,omitempty"`
	// 管理员改动原因/备注
	Remark string `json:"remark,omitempty"`
	// 骑手ID
	RiderID uint64 `json:"rider_id,omitempty"`
	// 管理人ID
	ManagerID *uint64 `json:"manager_id,omitempty"`
	// 企业ID
	EnterpriseID *uint64 `json:"enterprise_id,omitempty"`
	// AgentID holds the value of the "agent_id" field.
	AgentID *uint64 `json:"agent_id,omitempty"`
	// 订阅ID
	SubscribeID uint64 `json:"subscribe_id,omitempty"`
	// 更改天数
	Days int `json:"days,omitempty"`
	// 状态
	Status int `json:"status,omitempty"`
	// 审批时间
	ReviewTime *time.Time `json:"review_time,omitempty"`
	// 订阅预期到期时间
	SubscribeEndAt *time.Time `json:"subscribe_end_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SubscribeAlterQuery when eager-loading is set.
	Edges        SubscribeAlterEdges `json:"edges"`
	selectValues sql.SelectValues
}

// SubscribeAlterEdges holds the relations/edges for other nodes in the graph.
type SubscribeAlterEdges struct {
	// 骑手
	Rider *Rider `json:"rider,omitempty"`
	// Manager holds the value of the manager edge.
	Manager *Manager `json:"manager,omitempty"`
	// Enterprise holds the value of the enterprise edge.
	Enterprise *Enterprise `json:"enterprise,omitempty"`
	// Agent holds the value of the agent edge.
	Agent *Agent `json:"agent,omitempty"`
	// 订阅
	Subscribe *Subscribe `json:"subscribe,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [5]bool
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SubscribeAlterEdges) RiderOrErr() (*Rider, error) {
	if e.Rider != nil {
		return e.Rider, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: rider.Label}
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// ManagerOrErr returns the Manager value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SubscribeAlterEdges) ManagerOrErr() (*Manager, error) {
	if e.Manager != nil {
		return e.Manager, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: manager.Label}
	}
	return nil, &NotLoadedError{edge: "manager"}
}

// EnterpriseOrErr returns the Enterprise value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SubscribeAlterEdges) EnterpriseOrErr() (*Enterprise, error) {
	if e.Enterprise != nil {
		return e.Enterprise, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: enterprise.Label}
	}
	return nil, &NotLoadedError{edge: "enterprise"}
}

// AgentOrErr returns the Agent value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SubscribeAlterEdges) AgentOrErr() (*Agent, error) {
	if e.Agent != nil {
		return e.Agent, nil
	} else if e.loadedTypes[3] {
		return nil, &NotFoundError{label: agent.Label}
	}
	return nil, &NotLoadedError{edge: "agent"}
}

// SubscribeOrErr returns the Subscribe value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SubscribeAlterEdges) SubscribeOrErr() (*Subscribe, error) {
	if e.Subscribe != nil {
		return e.Subscribe, nil
	} else if e.loadedTypes[4] {
		return nil, &NotFoundError{label: subscribe.Label}
	}
	return nil, &NotLoadedError{edge: "subscribe"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*SubscribeAlter) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case subscribealter.FieldCreator, subscribealter.FieldLastModifier:
			values[i] = new([]byte)
		case subscribealter.FieldID, subscribealter.FieldRiderID, subscribealter.FieldManagerID, subscribealter.FieldEnterpriseID, subscribealter.FieldAgentID, subscribealter.FieldSubscribeID, subscribealter.FieldDays, subscribealter.FieldStatus:
			values[i] = new(sql.NullInt64)
		case subscribealter.FieldRemark:
			values[i] = new(sql.NullString)
		case subscribealter.FieldCreatedAt, subscribealter.FieldUpdatedAt, subscribealter.FieldReviewTime, subscribealter.FieldSubscribeEndAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SubscribeAlter fields.
func (sa *SubscribeAlter) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case subscribealter.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			sa.ID = uint64(value.Int64)
		case subscribealter.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				sa.CreatedAt = value.Time
			}
		case subscribealter.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				sa.UpdatedAt = value.Time
			}
		case subscribealter.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &sa.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case subscribealter.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &sa.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case subscribealter.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				sa.Remark = value.String
			}
		case subscribealter.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				sa.RiderID = uint64(value.Int64)
			}
		case subscribealter.FieldManagerID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field manager_id", values[i])
			} else if value.Valid {
				sa.ManagerID = new(uint64)
				*sa.ManagerID = uint64(value.Int64)
			}
		case subscribealter.FieldEnterpriseID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field enterprise_id", values[i])
			} else if value.Valid {
				sa.EnterpriseID = new(uint64)
				*sa.EnterpriseID = uint64(value.Int64)
			}
		case subscribealter.FieldAgentID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field agent_id", values[i])
			} else if value.Valid {
				sa.AgentID = new(uint64)
				*sa.AgentID = uint64(value.Int64)
			}
		case subscribealter.FieldSubscribeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field subscribe_id", values[i])
			} else if value.Valid {
				sa.SubscribeID = uint64(value.Int64)
			}
		case subscribealter.FieldDays:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field days", values[i])
			} else if value.Valid {
				sa.Days = int(value.Int64)
			}
		case subscribealter.FieldStatus:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				sa.Status = int(value.Int64)
			}
		case subscribealter.FieldReviewTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field review_time", values[i])
			} else if value.Valid {
				sa.ReviewTime = new(time.Time)
				*sa.ReviewTime = value.Time
			}
		case subscribealter.FieldSubscribeEndAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field subscribe_end_at", values[i])
			} else if value.Valid {
				sa.SubscribeEndAt = new(time.Time)
				*sa.SubscribeEndAt = value.Time
			}
		default:
			sa.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the SubscribeAlter.
// This includes values selected through modifiers, order, etc.
func (sa *SubscribeAlter) Value(name string) (ent.Value, error) {
	return sa.selectValues.Get(name)
}

// QueryRider queries the "rider" edge of the SubscribeAlter entity.
func (sa *SubscribeAlter) QueryRider() *RiderQuery {
	return NewSubscribeAlterClient(sa.config).QueryRider(sa)
}

// QueryManager queries the "manager" edge of the SubscribeAlter entity.
func (sa *SubscribeAlter) QueryManager() *ManagerQuery {
	return NewSubscribeAlterClient(sa.config).QueryManager(sa)
}

// QueryEnterprise queries the "enterprise" edge of the SubscribeAlter entity.
func (sa *SubscribeAlter) QueryEnterprise() *EnterpriseQuery {
	return NewSubscribeAlterClient(sa.config).QueryEnterprise(sa)
}

// QueryAgent queries the "agent" edge of the SubscribeAlter entity.
func (sa *SubscribeAlter) QueryAgent() *AgentQuery {
	return NewSubscribeAlterClient(sa.config).QueryAgent(sa)
}

// QuerySubscribe queries the "subscribe" edge of the SubscribeAlter entity.
func (sa *SubscribeAlter) QuerySubscribe() *SubscribeQuery {
	return NewSubscribeAlterClient(sa.config).QuerySubscribe(sa)
}

// Update returns a builder for updating this SubscribeAlter.
// Note that you need to call SubscribeAlter.Unwrap() before calling this method if this SubscribeAlter
// was returned from a transaction, and the transaction was committed or rolled back.
func (sa *SubscribeAlter) Update() *SubscribeAlterUpdateOne {
	return NewSubscribeAlterClient(sa.config).UpdateOne(sa)
}

// Unwrap unwraps the SubscribeAlter entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (sa *SubscribeAlter) Unwrap() *SubscribeAlter {
	_tx, ok := sa.config.driver.(*txDriver)
	if !ok {
		panic("ent: SubscribeAlter is not a transactional entity")
	}
	sa.config.driver = _tx.drv
	return sa
}

// String implements the fmt.Stringer.
func (sa *SubscribeAlter) String() string {
	var builder strings.Builder
	builder.WriteString("SubscribeAlter(")
	builder.WriteString(fmt.Sprintf("id=%v, ", sa.ID))
	builder.WriteString("created_at=")
	builder.WriteString(sa.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(sa.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", sa.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", sa.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(sa.Remark)
	builder.WriteString(", ")
	builder.WriteString("rider_id=")
	builder.WriteString(fmt.Sprintf("%v", sa.RiderID))
	builder.WriteString(", ")
	if v := sa.ManagerID; v != nil {
		builder.WriteString("manager_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := sa.EnterpriseID; v != nil {
		builder.WriteString("enterprise_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := sa.AgentID; v != nil {
		builder.WriteString("agent_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("subscribe_id=")
	builder.WriteString(fmt.Sprintf("%v", sa.SubscribeID))
	builder.WriteString(", ")
	builder.WriteString("days=")
	builder.WriteString(fmt.Sprintf("%v", sa.Days))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", sa.Status))
	builder.WriteString(", ")
	if v := sa.ReviewTime; v != nil {
		builder.WriteString("review_time=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := sa.SubscribeEndAt; v != nil {
		builder.WriteString("subscribe_end_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// SubscribeAlters is a parsable slice of SubscribeAlter.
type SubscribeAlters []*SubscribeAlter
