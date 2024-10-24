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
	"github.com/auroraride/aurservd/internal/ent/promotionsetting"
)

// PromotionSetting is the model entity for the PromotionSetting schema.
type PromotionSetting struct {
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
	// 标题
	Title string `json:"title,omitempty"`
	// 内容
	Content string `json:"content,omitempty"`
	// 设置项
	Key          string `json:"key,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PromotionSetting) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case promotionsetting.FieldCreator, promotionsetting.FieldLastModifier:
			values[i] = new([]byte)
		case promotionsetting.FieldID:
			values[i] = new(sql.NullInt64)
		case promotionsetting.FieldRemark, promotionsetting.FieldTitle, promotionsetting.FieldContent, promotionsetting.FieldKey:
			values[i] = new(sql.NullString)
		case promotionsetting.FieldCreatedAt, promotionsetting.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PromotionSetting fields.
func (ps *PromotionSetting) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case promotionsetting.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ps.ID = uint64(value.Int64)
		case promotionsetting.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				ps.CreatedAt = value.Time
			}
		case promotionsetting.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				ps.UpdatedAt = value.Time
			}
		case promotionsetting.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ps.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case promotionsetting.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ps.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case promotionsetting.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				ps.Remark = value.String
			}
		case promotionsetting.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				ps.Title = value.String
			}
		case promotionsetting.FieldContent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content", values[i])
			} else if value.Valid {
				ps.Content = value.String
			}
		case promotionsetting.FieldKey:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field key", values[i])
			} else if value.Valid {
				ps.Key = value.String
			}
		default:
			ps.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the PromotionSetting.
// This includes values selected through modifiers, order, etc.
func (ps *PromotionSetting) Value(name string) (ent.Value, error) {
	return ps.selectValues.Get(name)
}

// Update returns a builder for updating this PromotionSetting.
// Note that you need to call PromotionSetting.Unwrap() before calling this method if this PromotionSetting
// was returned from a transaction, and the transaction was committed or rolled back.
func (ps *PromotionSetting) Update() *PromotionSettingUpdateOne {
	return NewPromotionSettingClient(ps.config).UpdateOne(ps)
}

// Unwrap unwraps the PromotionSetting entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ps *PromotionSetting) Unwrap() *PromotionSetting {
	_tx, ok := ps.config.driver.(*txDriver)
	if !ok {
		panic("ent: PromotionSetting is not a transactional entity")
	}
	ps.config.driver = _tx.drv
	return ps
}

// String implements the fmt.Stringer.
func (ps *PromotionSetting) String() string {
	var builder strings.Builder
	builder.WriteString("PromotionSetting(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ps.ID))
	builder.WriteString("created_at=")
	builder.WriteString(ps.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(ps.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", ps.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", ps.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(ps.Remark)
	builder.WriteString(", ")
	builder.WriteString("title=")
	builder.WriteString(ps.Title)
	builder.WriteString(", ")
	builder.WriteString("content=")
	builder.WriteString(ps.Content)
	builder.WriteString(", ")
	builder.WriteString("key=")
	builder.WriteString(ps.Key)
	builder.WriteByte(')')
	return builder.String()
}

// PromotionSettings is a parsable slice of PromotionSetting.
type PromotionSettings []*PromotionSetting
