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
	"github.com/auroraride/aurservd/internal/ent/branch"
	"github.com/auroraride/aurservd/internal/ent/branchcontract"
)

// BranchContract is the model entity for the BranchContract schema.
type BranchContract struct {
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
	// 网点ID
	BranchID uint64 `json:"branch_id,omitempty"`
	// 房东姓名
	LandlordName string `json:"landlord_name,omitempty"`
	// 房东身份证
	IDCardNumber string `json:"id_card_number,omitempty"`
	// 房东手机号
	Phone string `json:"phone,omitempty"`
	// 房东银行卡号
	BankNumber string `json:"bank_number,omitempty"`
	// 押金
	Pledge float64 `json:"pledge,omitempty"`
	// 租金
	Rent float64 `json:"rent,omitempty"`
	// 租期月数
	Lease uint `json:"lease,omitempty"`
	// 电费押金
	ElectricityPledge float64 `json:"electricity_pledge,omitempty"`
	// 电费单价
	Electricity string `json:"electricity,omitempty"`
	// 网点面积
	Area float64 `json:"area,omitempty"`
	// 租期开始时间
	StartTime time.Time `json:"start_time,omitempty"`
	// 租期结束时间
	EndTime time.Time `json:"end_time,omitempty"`
	// 合同文件
	File string `json:"file,omitempty"`
	// 底单
	Sheets []string `json:"sheets,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the BranchContractQuery when eager-loading is set.
	Edges        BranchContractEdges `json:"edges"`
	selectValues sql.SelectValues
}

// BranchContractEdges holds the relations/edges for other nodes in the graph.
type BranchContractEdges struct {
	// Branch holds the value of the branch edge.
	Branch *Branch `json:"branch,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// BranchOrErr returns the Branch value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BranchContractEdges) BranchOrErr() (*Branch, error) {
	if e.Branch != nil {
		return e.Branch, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: branch.Label}
	}
	return nil, &NotLoadedError{edge: "branch"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*BranchContract) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case branchcontract.FieldCreator, branchcontract.FieldLastModifier, branchcontract.FieldSheets:
			values[i] = new([]byte)
		case branchcontract.FieldPledge, branchcontract.FieldRent, branchcontract.FieldElectricityPledge, branchcontract.FieldArea:
			values[i] = new(sql.NullFloat64)
		case branchcontract.FieldID, branchcontract.FieldBranchID, branchcontract.FieldLease:
			values[i] = new(sql.NullInt64)
		case branchcontract.FieldRemark, branchcontract.FieldLandlordName, branchcontract.FieldIDCardNumber, branchcontract.FieldPhone, branchcontract.FieldBankNumber, branchcontract.FieldElectricity, branchcontract.FieldFile:
			values[i] = new(sql.NullString)
		case branchcontract.FieldCreatedAt, branchcontract.FieldUpdatedAt, branchcontract.FieldDeletedAt, branchcontract.FieldStartTime, branchcontract.FieldEndTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the BranchContract fields.
func (bc *BranchContract) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case branchcontract.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			bc.ID = uint64(value.Int64)
		case branchcontract.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				bc.CreatedAt = value.Time
			}
		case branchcontract.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				bc.UpdatedAt = value.Time
			}
		case branchcontract.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				bc.DeletedAt = new(time.Time)
				*bc.DeletedAt = value.Time
			}
		case branchcontract.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &bc.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case branchcontract.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &bc.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case branchcontract.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				bc.Remark = value.String
			}
		case branchcontract.FieldBranchID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field branch_id", values[i])
			} else if value.Valid {
				bc.BranchID = uint64(value.Int64)
			}
		case branchcontract.FieldLandlordName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field landlord_name", values[i])
			} else if value.Valid {
				bc.LandlordName = value.String
			}
		case branchcontract.FieldIDCardNumber:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id_card_number", values[i])
			} else if value.Valid {
				bc.IDCardNumber = value.String
			}
		case branchcontract.FieldPhone:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field phone", values[i])
			} else if value.Valid {
				bc.Phone = value.String
			}
		case branchcontract.FieldBankNumber:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field bank_number", values[i])
			} else if value.Valid {
				bc.BankNumber = value.String
			}
		case branchcontract.FieldPledge:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field pledge", values[i])
			} else if value.Valid {
				bc.Pledge = value.Float64
			}
		case branchcontract.FieldRent:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field rent", values[i])
			} else if value.Valid {
				bc.Rent = value.Float64
			}
		case branchcontract.FieldLease:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field lease", values[i])
			} else if value.Valid {
				bc.Lease = uint(value.Int64)
			}
		case branchcontract.FieldElectricityPledge:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field electricity_pledge", values[i])
			} else if value.Valid {
				bc.ElectricityPledge = value.Float64
			}
		case branchcontract.FieldElectricity:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field electricity", values[i])
			} else if value.Valid {
				bc.Electricity = value.String
			}
		case branchcontract.FieldArea:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field area", values[i])
			} else if value.Valid {
				bc.Area = value.Float64
			}
		case branchcontract.FieldStartTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field start_time", values[i])
			} else if value.Valid {
				bc.StartTime = value.Time
			}
		case branchcontract.FieldEndTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field end_time", values[i])
			} else if value.Valid {
				bc.EndTime = value.Time
			}
		case branchcontract.FieldFile:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field file", values[i])
			} else if value.Valid {
				bc.File = value.String
			}
		case branchcontract.FieldSheets:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field sheets", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &bc.Sheets); err != nil {
					return fmt.Errorf("unmarshal field sheets: %w", err)
				}
			}
		default:
			bc.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the BranchContract.
// This includes values selected through modifiers, order, etc.
func (bc *BranchContract) Value(name string) (ent.Value, error) {
	return bc.selectValues.Get(name)
}

// QueryBranch queries the "branch" edge of the BranchContract entity.
func (bc *BranchContract) QueryBranch() *BranchQuery {
	return NewBranchContractClient(bc.config).QueryBranch(bc)
}

// Update returns a builder for updating this BranchContract.
// Note that you need to call BranchContract.Unwrap() before calling this method if this BranchContract
// was returned from a transaction, and the transaction was committed or rolled back.
func (bc *BranchContract) Update() *BranchContractUpdateOne {
	return NewBranchContractClient(bc.config).UpdateOne(bc)
}

// Unwrap unwraps the BranchContract entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (bc *BranchContract) Unwrap() *BranchContract {
	_tx, ok := bc.config.driver.(*txDriver)
	if !ok {
		panic("ent: BranchContract is not a transactional entity")
	}
	bc.config.driver = _tx.drv
	return bc
}

// String implements the fmt.Stringer.
func (bc *BranchContract) String() string {
	var builder strings.Builder
	builder.WriteString("BranchContract(")
	builder.WriteString(fmt.Sprintf("id=%v, ", bc.ID))
	builder.WriteString("created_at=")
	builder.WriteString(bc.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(bc.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := bc.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", bc.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", bc.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(bc.Remark)
	builder.WriteString(", ")
	builder.WriteString("branch_id=")
	builder.WriteString(fmt.Sprintf("%v", bc.BranchID))
	builder.WriteString(", ")
	builder.WriteString("landlord_name=")
	builder.WriteString(bc.LandlordName)
	builder.WriteString(", ")
	builder.WriteString("id_card_number=")
	builder.WriteString(bc.IDCardNumber)
	builder.WriteString(", ")
	builder.WriteString("phone=")
	builder.WriteString(bc.Phone)
	builder.WriteString(", ")
	builder.WriteString("bank_number=")
	builder.WriteString(bc.BankNumber)
	builder.WriteString(", ")
	builder.WriteString("pledge=")
	builder.WriteString(fmt.Sprintf("%v", bc.Pledge))
	builder.WriteString(", ")
	builder.WriteString("rent=")
	builder.WriteString(fmt.Sprintf("%v", bc.Rent))
	builder.WriteString(", ")
	builder.WriteString("lease=")
	builder.WriteString(fmt.Sprintf("%v", bc.Lease))
	builder.WriteString(", ")
	builder.WriteString("electricity_pledge=")
	builder.WriteString(fmt.Sprintf("%v", bc.ElectricityPledge))
	builder.WriteString(", ")
	builder.WriteString("electricity=")
	builder.WriteString(bc.Electricity)
	builder.WriteString(", ")
	builder.WriteString("area=")
	builder.WriteString(fmt.Sprintf("%v", bc.Area))
	builder.WriteString(", ")
	builder.WriteString("start_time=")
	builder.WriteString(bc.StartTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("end_time=")
	builder.WriteString(bc.EndTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("file=")
	builder.WriteString(bc.File)
	builder.WriteString(", ")
	builder.WriteString("sheets=")
	builder.WriteString(fmt.Sprintf("%v", bc.Sheets))
	builder.WriteByte(')')
	return builder.String()
}

// BranchContracts is a parsable slice of BranchContract.
type BranchContracts []*BranchContract
