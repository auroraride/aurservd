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
	"github.com/auroraride/aurservd/internal/ent/questioncategory"
)

// QuestionCategory is the model entity for the QuestionCategory schema.
type QuestionCategory struct {
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
	Sort uint64 `json:"sort,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the QuestionCategoryQuery when eager-loading is set.
	Edges        QuestionCategoryEdges `json:"edges"`
	selectValues sql.SelectValues
}

// QuestionCategoryEdges holds the relations/edges for other nodes in the graph.
type QuestionCategoryEdges struct {
	// 问题列表
	Questions []*Question `json:"questions,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// QuestionsOrErr returns the Questions value or an error if the edge
// was not loaded in eager-loading.
func (e QuestionCategoryEdges) QuestionsOrErr() ([]*Question, error) {
	if e.loadedTypes[0] {
		return e.Questions, nil
	}
	return nil, &NotLoadedError{edge: "questions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*QuestionCategory) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case questioncategory.FieldCreator, questioncategory.FieldLastModifier:
			values[i] = new([]byte)
		case questioncategory.FieldID, questioncategory.FieldSort:
			values[i] = new(sql.NullInt64)
		case questioncategory.FieldRemark, questioncategory.FieldName:
			values[i] = new(sql.NullString)
		case questioncategory.FieldCreatedAt, questioncategory.FieldUpdatedAt, questioncategory.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the QuestionCategory fields.
func (qc *QuestionCategory) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case questioncategory.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			qc.ID = uint64(value.Int64)
		case questioncategory.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				qc.CreatedAt = value.Time
			}
		case questioncategory.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				qc.UpdatedAt = value.Time
			}
		case questioncategory.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				qc.DeletedAt = new(time.Time)
				*qc.DeletedAt = value.Time
			}
		case questioncategory.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &qc.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case questioncategory.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &qc.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case questioncategory.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				qc.Remark = value.String
			}
		case questioncategory.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				qc.Name = value.String
			}
		case questioncategory.FieldSort:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field sort", values[i])
			} else if value.Valid {
				qc.Sort = uint64(value.Int64)
			}
		default:
			qc.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the QuestionCategory.
// This includes values selected through modifiers, order, etc.
func (qc *QuestionCategory) Value(name string) (ent.Value, error) {
	return qc.selectValues.Get(name)
}

// QueryQuestions queries the "questions" edge of the QuestionCategory entity.
func (qc *QuestionCategory) QueryQuestions() *QuestionQuery {
	return NewQuestionCategoryClient(qc.config).QueryQuestions(qc)
}

// Update returns a builder for updating this QuestionCategory.
// Note that you need to call QuestionCategory.Unwrap() before calling this method if this QuestionCategory
// was returned from a transaction, and the transaction was committed or rolled back.
func (qc *QuestionCategory) Update() *QuestionCategoryUpdateOne {
	return NewQuestionCategoryClient(qc.config).UpdateOne(qc)
}

// Unwrap unwraps the QuestionCategory entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (qc *QuestionCategory) Unwrap() *QuestionCategory {
	_tx, ok := qc.config.driver.(*txDriver)
	if !ok {
		panic("ent: QuestionCategory is not a transactional entity")
	}
	qc.config.driver = _tx.drv
	return qc
}

// String implements the fmt.Stringer.
func (qc *QuestionCategory) String() string {
	var builder strings.Builder
	builder.WriteString("QuestionCategory(")
	builder.WriteString(fmt.Sprintf("id=%v, ", qc.ID))
	builder.WriteString("created_at=")
	builder.WriteString(qc.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(qc.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := qc.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", qc.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", qc.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(qc.Remark)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(qc.Name)
	builder.WriteString(", ")
	builder.WriteString("sort=")
	builder.WriteString(fmt.Sprintf("%v", qc.Sort))
	builder.WriteByte(')')
	return builder.String()
}

// QuestionCategories is a parsable slice of QuestionCategory.
type QuestionCategories []*QuestionCategory
