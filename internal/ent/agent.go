// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
)

// Agent is the model entity for the Agent schema.
type Agent struct {
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
	// 企业ID
	EnterpriseID uint64 `json:"enterprise_id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Phone holds the value of the "phone" field.
	Phone string `json:"phone,omitempty"`
	// Password holds the value of the "password" field.
	Password string `json:"password,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AgentQuery when eager-loading is set.
	Edges AgentEdges `json:"edges"`
}

// AgentEdges holds the relations/edges for other nodes in the graph.
type AgentEdges struct {
	// Enterprise holds the value of the enterprise edge.
	Enterprise *Enterprise `json:"enterprise,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// EnterpriseOrErr returns the Enterprise value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AgentEdges) EnterpriseOrErr() (*Enterprise, error) {
	if e.loadedTypes[0] {
		if e.Enterprise == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: enterprise.Label}
		}
		return e.Enterprise, nil
	}
	return nil, &NotLoadedError{edge: "enterprise"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Agent) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case agent.FieldCreator, agent.FieldLastModifier:
			values[i] = new([]byte)
		case agent.FieldID, agent.FieldEnterpriseID:
			values[i] = new(sql.NullInt64)
		case agent.FieldRemark, agent.FieldName, agent.FieldPhone, agent.FieldPassword:
			values[i] = new(sql.NullString)
		case agent.FieldCreatedAt, agent.FieldUpdatedAt, agent.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Agent", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Agent fields.
func (a *Agent) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case agent.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			a.ID = uint64(value.Int64)
		case agent.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				a.CreatedAt = value.Time
			}
		case agent.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				a.UpdatedAt = value.Time
			}
		case agent.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				a.DeletedAt = new(time.Time)
				*a.DeletedAt = value.Time
			}
		case agent.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &a.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case agent.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &a.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case agent.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				a.Remark = value.String
			}
		case agent.FieldEnterpriseID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field enterprise_id", values[i])
			} else if value.Valid {
				a.EnterpriseID = uint64(value.Int64)
			}
		case agent.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				a.Name = value.String
			}
		case agent.FieldPhone:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field phone", values[i])
			} else if value.Valid {
				a.Phone = value.String
			}
		case agent.FieldPassword:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field password", values[i])
			} else if value.Valid {
				a.Password = value.String
			}
		}
	}
	return nil
}

// QueryEnterprise queries the "enterprise" edge of the Agent entity.
func (a *Agent) QueryEnterprise() *EnterpriseQuery {
	return (&AgentClient{config: a.config}).QueryEnterprise(a)
}

// Update returns a builder for updating this Agent.
// Note that you need to call Agent.Unwrap() before calling this method if this Agent
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Agent) Update() *AgentUpdateOne {
	return (&AgentClient{config: a.config}).UpdateOne(a)
}

// Unwrap unwraps the Agent entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Agent) Unwrap() *Agent {
	_tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Agent is not a transactional entity")
	}
	a.config.driver = _tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Agent) String() string {
	var builder strings.Builder
	builder.WriteString("Agent(")
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
	builder.WriteString("enterprise_id=")
	builder.WriteString(fmt.Sprintf("%v", a.EnterpriseID))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(a.Name)
	builder.WriteString(", ")
	builder.WriteString("phone=")
	builder.WriteString(a.Phone)
	builder.WriteString(", ")
	builder.WriteString("password=")
	builder.WriteString(a.Password)
	builder.WriteByte(')')
	return builder.String()
}

// Agents is a parsable slice of Agent.
type Agents []*Agent

func (a Agents) config(cfg config) {
	for _i := range a {
		a[_i].config = cfg
	}
}