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
	"github.com/auroraride/aurservd/internal/ent/promotionreferralsprogress"
	"github.com/auroraride/aurservd/internal/ent/rider"
)

// PromotionReferralsProgress is the model entity for the PromotionReferralsProgress schema.
type PromotionReferralsProgress struct {
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
	RiderID *uint64 `json:"rider_id,omitempty"`
	// 推广者id
	ReferringMemberID uint64 `json:"referring_member_id,omitempty"`
	// 被推广者ID<骑手>
	ReferredMemberID uint64 `json:"referred_member_id,omitempty"`
	// 姓名
	Name string `json:"name,omitempty"`
	// 状态  0: 邀请中 1:邀请成功 2:邀请失败
	Status *uint8 `json:"status,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PromotionReferralsProgressQuery when eager-loading is set.
	Edges        PromotionReferralsProgressEdges `json:"edges"`
	selectValues sql.SelectValues
}

// PromotionReferralsProgressEdges holds the relations/edges for other nodes in the graph.
type PromotionReferralsProgressEdges struct {
	// 骑手
	Rider *Rider `json:"rider,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PromotionReferralsProgressEdges) RiderOrErr() (*Rider, error) {
	if e.Rider != nil {
		return e.Rider, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: rider.Label}
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PromotionReferralsProgress) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case promotionreferralsprogress.FieldCreator, promotionreferralsprogress.FieldLastModifier:
			values[i] = new([]byte)
		case promotionreferralsprogress.FieldID, promotionreferralsprogress.FieldRiderID, promotionreferralsprogress.FieldReferringMemberID, promotionreferralsprogress.FieldReferredMemberID, promotionreferralsprogress.FieldStatus:
			values[i] = new(sql.NullInt64)
		case promotionreferralsprogress.FieldRemark, promotionreferralsprogress.FieldName:
			values[i] = new(sql.NullString)
		case promotionreferralsprogress.FieldCreatedAt, promotionreferralsprogress.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PromotionReferralsProgress fields.
func (prp *PromotionReferralsProgress) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case promotionreferralsprogress.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			prp.ID = uint64(value.Int64)
		case promotionreferralsprogress.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				prp.CreatedAt = value.Time
			}
		case promotionreferralsprogress.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				prp.UpdatedAt = value.Time
			}
		case promotionreferralsprogress.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &prp.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case promotionreferralsprogress.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &prp.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case promotionreferralsprogress.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				prp.Remark = value.String
			}
		case promotionreferralsprogress.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				prp.RiderID = new(uint64)
				*prp.RiderID = uint64(value.Int64)
			}
		case promotionreferralsprogress.FieldReferringMemberID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field referring_member_id", values[i])
			} else if value.Valid {
				prp.ReferringMemberID = uint64(value.Int64)
			}
		case promotionreferralsprogress.FieldReferredMemberID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field referred_member_id", values[i])
			} else if value.Valid {
				prp.ReferredMemberID = uint64(value.Int64)
			}
		case promotionreferralsprogress.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				prp.Name = value.String
			}
		case promotionreferralsprogress.FieldStatus:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				prp.Status = new(uint8)
				*prp.Status = uint8(value.Int64)
			}
		default:
			prp.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the PromotionReferralsProgress.
// This includes values selected through modifiers, order, etc.
func (prp *PromotionReferralsProgress) Value(name string) (ent.Value, error) {
	return prp.selectValues.Get(name)
}

// QueryRider queries the "rider" edge of the PromotionReferralsProgress entity.
func (prp *PromotionReferralsProgress) QueryRider() *RiderQuery {
	return NewPromotionReferralsProgressClient(prp.config).QueryRider(prp)
}

// Update returns a builder for updating this PromotionReferralsProgress.
// Note that you need to call PromotionReferralsProgress.Unwrap() before calling this method if this PromotionReferralsProgress
// was returned from a transaction, and the transaction was committed or rolled back.
func (prp *PromotionReferralsProgress) Update() *PromotionReferralsProgressUpdateOne {
	return NewPromotionReferralsProgressClient(prp.config).UpdateOne(prp)
}

// Unwrap unwraps the PromotionReferralsProgress entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (prp *PromotionReferralsProgress) Unwrap() *PromotionReferralsProgress {
	_tx, ok := prp.config.driver.(*txDriver)
	if !ok {
		panic("ent: PromotionReferralsProgress is not a transactional entity")
	}
	prp.config.driver = _tx.drv
	return prp
}

// String implements the fmt.Stringer.
func (prp *PromotionReferralsProgress) String() string {
	var builder strings.Builder
	builder.WriteString("PromotionReferralsProgress(")
	builder.WriteString(fmt.Sprintf("id=%v, ", prp.ID))
	builder.WriteString("created_at=")
	builder.WriteString(prp.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(prp.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", prp.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", prp.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(prp.Remark)
	builder.WriteString(", ")
	if v := prp.RiderID; v != nil {
		builder.WriteString("rider_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("referring_member_id=")
	builder.WriteString(fmt.Sprintf("%v", prp.ReferringMemberID))
	builder.WriteString(", ")
	builder.WriteString("referred_member_id=")
	builder.WriteString(fmt.Sprintf("%v", prp.ReferredMemberID))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(prp.Name)
	builder.WriteString(", ")
	if v := prp.Status; v != nil {
		builder.WriteString("status=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteByte(')')
	return builder.String()
}

// PromotionReferralsProgresses is a parsable slice of PromotionReferralsProgress.
type PromotionReferralsProgresses []*PromotionReferralsProgress
