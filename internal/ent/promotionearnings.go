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
	"github.com/auroraride/aurservd/internal/ent/promotioncommission"
	"github.com/auroraride/aurservd/internal/ent/promotionearnings"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/rider"
)

// PromotionEarnings is the model entity for the PromotionEarnings schema.
type PromotionEarnings struct {
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
	// CommissionID holds the value of the "commission_id" field.
	CommissionID uint64 `json:"commission_id,omitempty"`
	// MemberID holds the value of the "member_id" field.
	MemberID uint64 `json:"member_id,omitempty"`
	// 骑手ID
	RiderID uint64 `json:"rider_id,omitempty"`
	// 收益状态 0:未结算 1:已结算 2:已取消
	Status uint8 `json:"status,omitempty"`
	// 收益金额
	Amount float64 `json:"amount,omitempty"`
	// 返佣任务类型
	CommissionRuleKey string `json:"commission_rule_key,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PromotionEarningsQuery when eager-loading is set.
	Edges        PromotionEarningsEdges `json:"edges"`
	selectValues sql.SelectValues
}

// PromotionEarningsEdges holds the relations/edges for other nodes in the graph.
type PromotionEarningsEdges struct {
	// Commission holds the value of the commission edge.
	Commission *PromotionCommission `json:"commission,omitempty"`
	// Member holds the value of the member edge.
	Member *PromotionMember `json:"member,omitempty"`
	// 骑手
	Rider *Rider `json:"rider,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// CommissionOrErr returns the Commission value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PromotionEarningsEdges) CommissionOrErr() (*PromotionCommission, error) {
	if e.loadedTypes[0] {
		if e.Commission == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: promotioncommission.Label}
		}
		return e.Commission, nil
	}
	return nil, &NotLoadedError{edge: "commission"}
}

// MemberOrErr returns the Member value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PromotionEarningsEdges) MemberOrErr() (*PromotionMember, error) {
	if e.loadedTypes[1] {
		if e.Member == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: promotionmember.Label}
		}
		return e.Member, nil
	}
	return nil, &NotLoadedError{edge: "member"}
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PromotionEarningsEdges) RiderOrErr() (*Rider, error) {
	if e.loadedTypes[2] {
		if e.Rider == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: rider.Label}
		}
		return e.Rider, nil
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PromotionEarnings) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case promotionearnings.FieldCreator, promotionearnings.FieldLastModifier:
			values[i] = new([]byte)
		case promotionearnings.FieldAmount:
			values[i] = new(sql.NullFloat64)
		case promotionearnings.FieldID, promotionearnings.FieldCommissionID, promotionearnings.FieldMemberID, promotionearnings.FieldRiderID, promotionearnings.FieldStatus:
			values[i] = new(sql.NullInt64)
		case promotionearnings.FieldRemark, promotionearnings.FieldCommissionRuleKey:
			values[i] = new(sql.NullString)
		case promotionearnings.FieldCreatedAt, promotionearnings.FieldUpdatedAt, promotionearnings.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PromotionEarnings fields.
func (pe *PromotionEarnings) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case promotionearnings.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pe.ID = uint64(value.Int64)
		case promotionearnings.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pe.CreatedAt = value.Time
			}
		case promotionearnings.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				pe.UpdatedAt = value.Time
			}
		case promotionearnings.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				pe.DeletedAt = new(time.Time)
				*pe.DeletedAt = value.Time
			}
		case promotionearnings.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &pe.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case promotionearnings.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &pe.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case promotionearnings.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				pe.Remark = value.String
			}
		case promotionearnings.FieldCommissionID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field commission_id", values[i])
			} else if value.Valid {
				pe.CommissionID = uint64(value.Int64)
			}
		case promotionearnings.FieldMemberID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field member_id", values[i])
			} else if value.Valid {
				pe.MemberID = uint64(value.Int64)
			}
		case promotionearnings.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				pe.RiderID = uint64(value.Int64)
			}
		case promotionearnings.FieldStatus:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				pe.Status = uint8(value.Int64)
			}
		case promotionearnings.FieldAmount:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field amount", values[i])
			} else if value.Valid {
				pe.Amount = value.Float64
			}
		case promotionearnings.FieldCommissionRuleKey:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field commission_rule_key", values[i])
			} else if value.Valid {
				pe.CommissionRuleKey = value.String
			}
		default:
			pe.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the PromotionEarnings.
// This includes values selected through modifiers, order, etc.
func (pe *PromotionEarnings) Value(name string) (ent.Value, error) {
	return pe.selectValues.Get(name)
}

// QueryCommission queries the "commission" edge of the PromotionEarnings entity.
func (pe *PromotionEarnings) QueryCommission() *PromotionCommissionQuery {
	return NewPromotionEarningsClient(pe.config).QueryCommission(pe)
}

// QueryMember queries the "member" edge of the PromotionEarnings entity.
func (pe *PromotionEarnings) QueryMember() *PromotionMemberQuery {
	return NewPromotionEarningsClient(pe.config).QueryMember(pe)
}

// QueryRider queries the "rider" edge of the PromotionEarnings entity.
func (pe *PromotionEarnings) QueryRider() *RiderQuery {
	return NewPromotionEarningsClient(pe.config).QueryRider(pe)
}

// Update returns a builder for updating this PromotionEarnings.
// Note that you need to call PromotionEarnings.Unwrap() before calling this method if this PromotionEarnings
// was returned from a transaction, and the transaction was committed or rolled back.
func (pe *PromotionEarnings) Update() *PromotionEarningsUpdateOne {
	return NewPromotionEarningsClient(pe.config).UpdateOne(pe)
}

// Unwrap unwraps the PromotionEarnings entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pe *PromotionEarnings) Unwrap() *PromotionEarnings {
	_tx, ok := pe.config.driver.(*txDriver)
	if !ok {
		panic("ent: PromotionEarnings is not a transactional entity")
	}
	pe.config.driver = _tx.drv
	return pe
}

// String implements the fmt.Stringer.
func (pe *PromotionEarnings) String() string {
	var builder strings.Builder
	builder.WriteString("PromotionEarnings(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pe.ID))
	builder.WriteString("created_at=")
	builder.WriteString(pe.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(pe.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := pe.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", pe.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", pe.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(pe.Remark)
	builder.WriteString(", ")
	builder.WriteString("commission_id=")
	builder.WriteString(fmt.Sprintf("%v", pe.CommissionID))
	builder.WriteString(", ")
	builder.WriteString("member_id=")
	builder.WriteString(fmt.Sprintf("%v", pe.MemberID))
	builder.WriteString(", ")
	builder.WriteString("rider_id=")
	builder.WriteString(fmt.Sprintf("%v", pe.RiderID))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", pe.Status))
	builder.WriteString(", ")
	builder.WriteString("amount=")
	builder.WriteString(fmt.Sprintf("%v", pe.Amount))
	builder.WriteString(", ")
	builder.WriteString("commission_rule_key=")
	builder.WriteString(pe.CommissionRuleKey)
	builder.WriteByte(')')
	return builder.String()
}

// PromotionEarningsSlice is a parsable slice of PromotionEarnings.
type PromotionEarningsSlice []*PromotionEarnings