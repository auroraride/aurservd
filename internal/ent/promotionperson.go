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
	"github.com/auroraride/aurservd/internal/ent/promotionperson"
)

// PromotionPerson is the model entity for the PromotionPerson schema.
type PromotionPerson struct {
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
	// 认证状态 0未认证 1已认证 2认证失败
	Status uint8 `json:"status,omitempty"`
	// 真实姓名
	Name string `json:"name,omitempty"`
	// 证件号码
	IDCardNumber string `json:"id_card_number,omitempty"`
	// 地址
	Address string `json:"address,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PromotionPersonQuery when eager-loading is set.
	Edges        PromotionPersonEdges `json:"edges"`
	selectValues sql.SelectValues
}

// PromotionPersonEdges holds the relations/edges for other nodes in the graph.
type PromotionPersonEdges struct {
	// Member holds the value of the member edge.
	Member []*PromotionMember `json:"member,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// MemberOrErr returns the Member value or an error if the edge
// was not loaded in eager-loading.
func (e PromotionPersonEdges) MemberOrErr() ([]*PromotionMember, error) {
	if e.loadedTypes[0] {
		return e.Member, nil
	}
	return nil, &NotLoadedError{edge: "member"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PromotionPerson) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case promotionperson.FieldCreator, promotionperson.FieldLastModifier:
			values[i] = new([]byte)
		case promotionperson.FieldID, promotionperson.FieldStatus:
			values[i] = new(sql.NullInt64)
		case promotionperson.FieldRemark, promotionperson.FieldName, promotionperson.FieldIDCardNumber, promotionperson.FieldAddress:
			values[i] = new(sql.NullString)
		case promotionperson.FieldCreatedAt, promotionperson.FieldUpdatedAt, promotionperson.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PromotionPerson fields.
func (pp *PromotionPerson) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case promotionperson.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pp.ID = uint64(value.Int64)
		case promotionperson.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pp.CreatedAt = value.Time
			}
		case promotionperson.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				pp.UpdatedAt = value.Time
			}
		case promotionperson.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				pp.DeletedAt = new(time.Time)
				*pp.DeletedAt = value.Time
			}
		case promotionperson.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &pp.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case promotionperson.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &pp.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case promotionperson.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				pp.Remark = value.String
			}
		case promotionperson.FieldStatus:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				pp.Status = uint8(value.Int64)
			}
		case promotionperson.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				pp.Name = value.String
			}
		case promotionperson.FieldIDCardNumber:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id_card_number", values[i])
			} else if value.Valid {
				pp.IDCardNumber = value.String
			}
		case promotionperson.FieldAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field address", values[i])
			} else if value.Valid {
				pp.Address = value.String
			}
		default:
			pp.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the PromotionPerson.
// This includes values selected through modifiers, order, etc.
func (pp *PromotionPerson) Value(name string) (ent.Value, error) {
	return pp.selectValues.Get(name)
}

// QueryMember queries the "member" edge of the PromotionPerson entity.
func (pp *PromotionPerson) QueryMember() *PromotionMemberQuery {
	return NewPromotionPersonClient(pp.config).QueryMember(pp)
}

// Update returns a builder for updating this PromotionPerson.
// Note that you need to call PromotionPerson.Unwrap() before calling this method if this PromotionPerson
// was returned from a transaction, and the transaction was committed or rolled back.
func (pp *PromotionPerson) Update() *PromotionPersonUpdateOne {
	return NewPromotionPersonClient(pp.config).UpdateOne(pp)
}

// Unwrap unwraps the PromotionPerson entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pp *PromotionPerson) Unwrap() *PromotionPerson {
	_tx, ok := pp.config.driver.(*txDriver)
	if !ok {
		panic("ent: PromotionPerson is not a transactional entity")
	}
	pp.config.driver = _tx.drv
	return pp
}

// String implements the fmt.Stringer.
func (pp *PromotionPerson) String() string {
	var builder strings.Builder
	builder.WriteString("PromotionPerson(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pp.ID))
	builder.WriteString("created_at=")
	builder.WriteString(pp.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(pp.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := pp.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", pp.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", pp.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(pp.Remark)
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", pp.Status))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(pp.Name)
	builder.WriteString(", ")
	builder.WriteString("id_card_number=")
	builder.WriteString(pp.IDCardNumber)
	builder.WriteString(", ")
	builder.WriteString("address=")
	builder.WriteString(pp.Address)
	builder.WriteByte(')')
	return builder.String()
}

// PromotionPersons is a parsable slice of PromotionPerson.
type PromotionPersons []*PromotionPerson