// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/internal/ent/cabinetec"
)

// CabinetEc is the model entity for the CabinetEc schema.
type CabinetEc struct {
	config `json:"-"`
	// ID of the ent.
	ID uint64 `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	// 电柜原始编号
	Serial string `json:"serial,omitempty"`
	// 日期
	Date time.Time `json:"date,omitempty"`
	// 开始电量
	Start float64 `json:"start,omitempty"`
	// 结束电量
	End float64 `json:"end,omitempty"`
	// 耗电量
	Total        float64 `json:"total,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*CabinetEc) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case cabinetec.FieldStart, cabinetec.FieldEnd, cabinetec.FieldTotal:
			values[i] = new(sql.NullFloat64)
		case cabinetec.FieldID:
			values[i] = new(sql.NullInt64)
		case cabinetec.FieldSerial:
			values[i] = new(sql.NullString)
		case cabinetec.FieldCreatedAt, cabinetec.FieldUpdatedAt, cabinetec.FieldDeletedAt, cabinetec.FieldDate:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CabinetEc fields.
func (ce *CabinetEc) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case cabinetec.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ce.ID = uint64(value.Int64)
		case cabinetec.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				ce.CreatedAt = value.Time
			}
		case cabinetec.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				ce.UpdatedAt = value.Time
			}
		case cabinetec.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				ce.DeletedAt = new(time.Time)
				*ce.DeletedAt = value.Time
			}
		case cabinetec.FieldSerial:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field serial", values[i])
			} else if value.Valid {
				ce.Serial = value.String
			}
		case cabinetec.FieldDate:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date", values[i])
			} else if value.Valid {
				ce.Date = value.Time
			}
		case cabinetec.FieldStart:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field start", values[i])
			} else if value.Valid {
				ce.Start = value.Float64
			}
		case cabinetec.FieldEnd:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field end", values[i])
			} else if value.Valid {
				ce.End = value.Float64
			}
		case cabinetec.FieldTotal:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field total", values[i])
			} else if value.Valid {
				ce.Total = value.Float64
			}
		default:
			ce.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the CabinetEc.
// This includes values selected through modifiers, order, etc.
func (ce *CabinetEc) Value(name string) (ent.Value, error) {
	return ce.selectValues.Get(name)
}

// Update returns a builder for updating this CabinetEc.
// Note that you need to call CabinetEc.Unwrap() before calling this method if this CabinetEc
// was returned from a transaction, and the transaction was committed or rolled back.
func (ce *CabinetEc) Update() *CabinetEcUpdateOne {
	return NewCabinetEcClient(ce.config).UpdateOne(ce)
}

// Unwrap unwraps the CabinetEc entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ce *CabinetEc) Unwrap() *CabinetEc {
	_tx, ok := ce.config.driver.(*txDriver)
	if !ok {
		panic("ent: CabinetEc is not a transactional entity")
	}
	ce.config.driver = _tx.drv
	return ce
}

// String implements the fmt.Stringer.
func (ce *CabinetEc) String() string {
	var builder strings.Builder
	builder.WriteString("CabinetEc(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ce.ID))
	builder.WriteString("created_at=")
	builder.WriteString(ce.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(ce.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := ce.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("serial=")
	builder.WriteString(ce.Serial)
	builder.WriteString(", ")
	builder.WriteString("date=")
	builder.WriteString(ce.Date.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("start=")
	builder.WriteString(fmt.Sprintf("%v", ce.Start))
	builder.WriteString(", ")
	builder.WriteString("end=")
	builder.WriteString(fmt.Sprintf("%v", ce.End))
	builder.WriteString(", ")
	builder.WriteString("total=")
	builder.WriteString(fmt.Sprintf("%v", ce.Total))
	builder.WriteByte(')')
	return builder.String()
}

// CabinetEcs is a parsable slice of CabinetEc.
type CabinetEcs []*CabinetEc