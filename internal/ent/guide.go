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
	"github.com/auroraride/aurservd/internal/ent/guide"
)

// Guide is the model entity for the Guide schema.
type Guide struct {
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
	// 名称
	Name string `json:"name,omitempty"`
	// 排序
	Sort uint8 `json:"sort,omitempty"`
	// 答案
	Answer       string `json:"answer,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Guide) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case guide.FieldCreator, guide.FieldLastModifier:
			values[i] = new([]byte)
		case guide.FieldID, guide.FieldSort:
			values[i] = new(sql.NullInt64)
		case guide.FieldRemark, guide.FieldName, guide.FieldAnswer:
			values[i] = new(sql.NullString)
		case guide.FieldCreatedAt, guide.FieldUpdatedAt, guide.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Guide fields.
func (gu *Guide) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case guide.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			gu.ID = uint64(value.Int64)
		case guide.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				gu.CreatedAt = value.Time
			}
		case guide.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				gu.UpdatedAt = value.Time
			}
		case guide.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				gu.DeletedAt = new(time.Time)
				*gu.DeletedAt = value.Time
			}
		case guide.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &gu.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case guide.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &gu.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case guide.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				gu.Remark = value.String
			}
		case guide.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				gu.Name = value.String
			}
		case guide.FieldSort:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field sort", values[i])
			} else if value.Valid {
				gu.Sort = uint8(value.Int64)
			}
		case guide.FieldAnswer:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field answer", values[i])
			} else if value.Valid {
				gu.Answer = value.String
			}
		default:
			gu.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Guide.
// This includes values selected through modifiers, order, etc.
func (gu *Guide) Value(name string) (ent.Value, error) {
	return gu.selectValues.Get(name)
}

// Update returns a builder for updating this Guide.
// Note that you need to call Guide.Unwrap() before calling this method if this Guide
// was returned from a transaction, and the transaction was committed or rolled back.
func (gu *Guide) Update() *GuideUpdateOne {
	return NewGuideClient(gu.config).UpdateOne(gu)
}

// Unwrap unwraps the Guide entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (gu *Guide) Unwrap() *Guide {
	_tx, ok := gu.config.driver.(*txDriver)
	if !ok {
		panic("ent: Guide is not a transactional entity")
	}
	gu.config.driver = _tx.drv
	return gu
}

// String implements the fmt.Stringer.
func (gu *Guide) String() string {
	var builder strings.Builder
	builder.WriteString("Guide(")
	builder.WriteString(fmt.Sprintf("id=%v, ", gu.ID))
	builder.WriteString("created_at=")
	builder.WriteString(gu.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(gu.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := gu.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", gu.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", gu.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(gu.Remark)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(gu.Name)
	builder.WriteString(", ")
	builder.WriteString("sort=")
	builder.WriteString(fmt.Sprintf("%v", gu.Sort))
	builder.WriteString(", ")
	builder.WriteString("answer=")
	builder.WriteString(gu.Answer)
	builder.WriteByte(')')
	return builder.String()
}

// Guides is a parsable slice of Guide.
type Guides []*Guide