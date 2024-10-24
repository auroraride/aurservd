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
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assetscrap"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
)

// AssetScrap is the model entity for the AssetScrap schema.
type AssetScrap struct {
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
	// 报废原因 1:丢失 2:损坏 3:其他
	ReasonType uint8 `json:"reason_type,omitempty"`
	// 报废时间
	ScrapAt time.Time `json:"scrap_at,omitempty"`
	// 操作报废人员ID
	OperateID *uint64 `json:"operate_id,omitempty"`
	// 报废人员角色类型 0:业务后台 1:门店 2:代理 3:运维 4:电柜 5:骑手 6:资产后台
	OperateRoleType *uint8 `json:"operate_role_type,omitempty"`
	// 报废编号
	Sn string `json:"sn,omitempty"`
	// 报废数量
	Num uint `json:"num,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AssetScrapQuery when eager-loading is set.
	Edges        AssetScrapEdges `json:"edges"`
	selectValues sql.SelectValues
}

// AssetScrapEdges holds the relations/edges for other nodes in the graph.
type AssetScrapEdges struct {
	// Manager holds the value of the manager edge.
	Manager *AssetManager `json:"manager,omitempty"`
	// Employee holds the value of the employee edge.
	Employee *Employee `json:"employee,omitempty"`
	// Maintainer holds the value of the maintainer edge.
	Maintainer *Maintainer `json:"maintainer,omitempty"`
	// Agent holds the value of the agent edge.
	Agent *Agent `json:"agent,omitempty"`
	// ScrapDetails holds the value of the scrap_details edge.
	ScrapDetails []*AssetScrapDetails `json:"scrap_details,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [5]bool
}

// ManagerOrErr returns the Manager value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssetScrapEdges) ManagerOrErr() (*AssetManager, error) {
	if e.Manager != nil {
		return e.Manager, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: assetmanager.Label}
	}
	return nil, &NotLoadedError{edge: "manager"}
}

// EmployeeOrErr returns the Employee value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssetScrapEdges) EmployeeOrErr() (*Employee, error) {
	if e.Employee != nil {
		return e.Employee, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: employee.Label}
	}
	return nil, &NotLoadedError{edge: "employee"}
}

// MaintainerOrErr returns the Maintainer value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssetScrapEdges) MaintainerOrErr() (*Maintainer, error) {
	if e.Maintainer != nil {
		return e.Maintainer, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: maintainer.Label}
	}
	return nil, &NotLoadedError{edge: "maintainer"}
}

// AgentOrErr returns the Agent value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssetScrapEdges) AgentOrErr() (*Agent, error) {
	if e.Agent != nil {
		return e.Agent, nil
	} else if e.loadedTypes[3] {
		return nil, &NotFoundError{label: agent.Label}
	}
	return nil, &NotLoadedError{edge: "agent"}
}

// ScrapDetailsOrErr returns the ScrapDetails value or an error if the edge
// was not loaded in eager-loading.
func (e AssetScrapEdges) ScrapDetailsOrErr() ([]*AssetScrapDetails, error) {
	if e.loadedTypes[4] {
		return e.ScrapDetails, nil
	}
	return nil, &NotLoadedError{edge: "scrap_details"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*AssetScrap) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case assetscrap.FieldCreator, assetscrap.FieldLastModifier:
			values[i] = new([]byte)
		case assetscrap.FieldID, assetscrap.FieldReasonType, assetscrap.FieldOperateID, assetscrap.FieldOperateRoleType, assetscrap.FieldNum:
			values[i] = new(sql.NullInt64)
		case assetscrap.FieldRemark, assetscrap.FieldSn:
			values[i] = new(sql.NullString)
		case assetscrap.FieldCreatedAt, assetscrap.FieldUpdatedAt, assetscrap.FieldScrapAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the AssetScrap fields.
func (as *AssetScrap) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case assetscrap.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			as.ID = uint64(value.Int64)
		case assetscrap.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				as.CreatedAt = value.Time
			}
		case assetscrap.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				as.UpdatedAt = value.Time
			}
		case assetscrap.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &as.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case assetscrap.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &as.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case assetscrap.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				as.Remark = value.String
			}
		case assetscrap.FieldReasonType:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field reason_type", values[i])
			} else if value.Valid {
				as.ReasonType = uint8(value.Int64)
			}
		case assetscrap.FieldScrapAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field scrap_at", values[i])
			} else if value.Valid {
				as.ScrapAt = value.Time
			}
		case assetscrap.FieldOperateID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field operate_id", values[i])
			} else if value.Valid {
				as.OperateID = new(uint64)
				*as.OperateID = uint64(value.Int64)
			}
		case assetscrap.FieldOperateRoleType:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field operate_role_type", values[i])
			} else if value.Valid {
				as.OperateRoleType = new(uint8)
				*as.OperateRoleType = uint8(value.Int64)
			}
		case assetscrap.FieldSn:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sn", values[i])
			} else if value.Valid {
				as.Sn = value.String
			}
		case assetscrap.FieldNum:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field num", values[i])
			} else if value.Valid {
				as.Num = uint(value.Int64)
			}
		default:
			as.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the AssetScrap.
// This includes values selected through modifiers, order, etc.
func (as *AssetScrap) Value(name string) (ent.Value, error) {
	return as.selectValues.Get(name)
}

// QueryManager queries the "manager" edge of the AssetScrap entity.
func (as *AssetScrap) QueryManager() *AssetManagerQuery {
	return NewAssetScrapClient(as.config).QueryManager(as)
}

// QueryEmployee queries the "employee" edge of the AssetScrap entity.
func (as *AssetScrap) QueryEmployee() *EmployeeQuery {
	return NewAssetScrapClient(as.config).QueryEmployee(as)
}

// QueryMaintainer queries the "maintainer" edge of the AssetScrap entity.
func (as *AssetScrap) QueryMaintainer() *MaintainerQuery {
	return NewAssetScrapClient(as.config).QueryMaintainer(as)
}

// QueryAgent queries the "agent" edge of the AssetScrap entity.
func (as *AssetScrap) QueryAgent() *AgentQuery {
	return NewAssetScrapClient(as.config).QueryAgent(as)
}

// QueryScrapDetails queries the "scrap_details" edge of the AssetScrap entity.
func (as *AssetScrap) QueryScrapDetails() *AssetScrapDetailsQuery {
	return NewAssetScrapClient(as.config).QueryScrapDetails(as)
}

// Update returns a builder for updating this AssetScrap.
// Note that you need to call AssetScrap.Unwrap() before calling this method if this AssetScrap
// was returned from a transaction, and the transaction was committed or rolled back.
func (as *AssetScrap) Update() *AssetScrapUpdateOne {
	return NewAssetScrapClient(as.config).UpdateOne(as)
}

// Unwrap unwraps the AssetScrap entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (as *AssetScrap) Unwrap() *AssetScrap {
	_tx, ok := as.config.driver.(*txDriver)
	if !ok {
		panic("ent: AssetScrap is not a transactional entity")
	}
	as.config.driver = _tx.drv
	return as
}

// String implements the fmt.Stringer.
func (as *AssetScrap) String() string {
	var builder strings.Builder
	builder.WriteString("AssetScrap(")
	builder.WriteString(fmt.Sprintf("id=%v, ", as.ID))
	builder.WriteString("created_at=")
	builder.WriteString(as.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(as.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", as.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", as.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(as.Remark)
	builder.WriteString(", ")
	builder.WriteString("reason_type=")
	builder.WriteString(fmt.Sprintf("%v", as.ReasonType))
	builder.WriteString(", ")
	builder.WriteString("scrap_at=")
	builder.WriteString(as.ScrapAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := as.OperateID; v != nil {
		builder.WriteString("operate_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := as.OperateRoleType; v != nil {
		builder.WriteString("operate_role_type=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("sn=")
	builder.WriteString(as.Sn)
	builder.WriteString(", ")
	builder.WriteString("num=")
	builder.WriteString(fmt.Sprintf("%v", as.Num))
	builder.WriteByte(')')
	return builder.String()
}

// AssetScraps is a parsable slice of AssetScrap.
type AssetScraps []*AssetScrap
