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
	"github.com/auroraride/aurservd/internal/ent/manager"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/riderfollowup"
)

// RiderFollowUp is the model entity for the RiderFollowUp schema.
type RiderFollowUp struct {
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
	// 管理人ID
	ManagerID uint64 `json:"manager_id,omitempty"`
	// 骑手ID
	RiderID uint64 `json:"rider_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the RiderFollowUpQuery when eager-loading is set.
	Edges        RiderFollowUpEdges `json:"edges"`
	selectValues sql.SelectValues
}

// RiderFollowUpEdges holds the relations/edges for other nodes in the graph.
type RiderFollowUpEdges struct {
	// Manager holds the value of the manager edge.
	Manager *Manager `json:"manager,omitempty"`
	// Rider holds the value of the rider edge.
	Rider *Rider `json:"rider,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ManagerOrErr returns the Manager value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RiderFollowUpEdges) ManagerOrErr() (*Manager, error) {
	if e.Manager != nil {
		return e.Manager, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: manager.Label}
	}
	return nil, &NotLoadedError{edge: "manager"}
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RiderFollowUpEdges) RiderOrErr() (*Rider, error) {
	if e.Rider != nil {
		return e.Rider, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: rider.Label}
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*RiderFollowUp) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case riderfollowup.FieldCreator, riderfollowup.FieldLastModifier:
			values[i] = new([]byte)
		case riderfollowup.FieldID, riderfollowup.FieldManagerID, riderfollowup.FieldRiderID:
			values[i] = new(sql.NullInt64)
		case riderfollowup.FieldRemark:
			values[i] = new(sql.NullString)
		case riderfollowup.FieldCreatedAt, riderfollowup.FieldUpdatedAt, riderfollowup.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the RiderFollowUp fields.
func (rfu *RiderFollowUp) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case riderfollowup.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			rfu.ID = uint64(value.Int64)
		case riderfollowup.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				rfu.CreatedAt = value.Time
			}
		case riderfollowup.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				rfu.UpdatedAt = value.Time
			}
		case riderfollowup.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				rfu.DeletedAt = new(time.Time)
				*rfu.DeletedAt = value.Time
			}
		case riderfollowup.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &rfu.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case riderfollowup.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &rfu.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case riderfollowup.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				rfu.Remark = value.String
			}
		case riderfollowup.FieldManagerID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field manager_id", values[i])
			} else if value.Valid {
				rfu.ManagerID = uint64(value.Int64)
			}
		case riderfollowup.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				rfu.RiderID = uint64(value.Int64)
			}
		default:
			rfu.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the RiderFollowUp.
// This includes values selected through modifiers, order, etc.
func (rfu *RiderFollowUp) Value(name string) (ent.Value, error) {
	return rfu.selectValues.Get(name)
}

// QueryManager queries the "manager" edge of the RiderFollowUp entity.
func (rfu *RiderFollowUp) QueryManager() *ManagerQuery {
	return NewRiderFollowUpClient(rfu.config).QueryManager(rfu)
}

// QueryRider queries the "rider" edge of the RiderFollowUp entity.
func (rfu *RiderFollowUp) QueryRider() *RiderQuery {
	return NewRiderFollowUpClient(rfu.config).QueryRider(rfu)
}

// Update returns a builder for updating this RiderFollowUp.
// Note that you need to call RiderFollowUp.Unwrap() before calling this method if this RiderFollowUp
// was returned from a transaction, and the transaction was committed or rolled back.
func (rfu *RiderFollowUp) Update() *RiderFollowUpUpdateOne {
	return NewRiderFollowUpClient(rfu.config).UpdateOne(rfu)
}

// Unwrap unwraps the RiderFollowUp entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (rfu *RiderFollowUp) Unwrap() *RiderFollowUp {
	_tx, ok := rfu.config.driver.(*txDriver)
	if !ok {
		panic("ent: RiderFollowUp is not a transactional entity")
	}
	rfu.config.driver = _tx.drv
	return rfu
}

// String implements the fmt.Stringer.
func (rfu *RiderFollowUp) String() string {
	var builder strings.Builder
	builder.WriteString("RiderFollowUp(")
	builder.WriteString(fmt.Sprintf("id=%v, ", rfu.ID))
	builder.WriteString("created_at=")
	builder.WriteString(rfu.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(rfu.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := rfu.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", rfu.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", rfu.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(rfu.Remark)
	builder.WriteString(", ")
	builder.WriteString("manager_id=")
	builder.WriteString(fmt.Sprintf("%v", rfu.ManagerID))
	builder.WriteString(", ")
	builder.WriteString("rider_id=")
	builder.WriteString(fmt.Sprintf("%v", rfu.RiderID))
	builder.WriteByte(')')
	return builder.String()
}

// RiderFollowUps is a parsable slice of RiderFollowUp.
type RiderFollowUps []*RiderFollowUp
