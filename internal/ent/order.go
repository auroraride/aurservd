// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/commission"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/rider"
)

// Order is the model entity for the Order schema.
type Order struct {
	config `json:"-"`
	// ID of the ent.
	ID uint64 `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	// Creator holds the value of the "creator" field.
	// 创建人
	Creator *model.Modifier `json:"creator,omitempty"`
	// LastModifier holds the value of the "last_modifier" field.
	// 最后修改人
	LastModifier *model.Modifier `json:"last_modifier,omitempty"`
	// Remark holds the value of the "remark" field.
	// 备注
	Remark string `json:"remark,omitempty"`
	// RiderID holds the value of the "rider_id" field.
	// 骑手ID
	RiderID uint64 `json:"rider_id,omitempty"`
	// PlanID holds the value of the "plan_id" field.
	// 骑士卡ID
	PlanID uint64 `json:"plan_id,omitempty"`
	// Status holds the value of the "status" field.
	// 订单状态 0未支付 1已支付 2申请退款 3已退款
	Status uint8 `json:"status,omitempty"`
	// Payway holds the value of the "payway" field.
	// 支付方式 1支付宝 2微信
	Payway uint8 `json:"payway,omitempty"`
	// Type holds the value of the "type" field.
	// 订单类型 1新签 2续签 3重签 4更改电池 5救援 6滞纳金
	Type uint `json:"type,omitempty"`
	// OutTradeNo holds the value of the "out_trade_no" field.
	// 交易订单号
	OutTradeNo string `json:"out_trade_no,omitempty"`
	// TradeNo holds the value of the "trade_no" field.
	// 平台订单号
	TradeNo string `json:"trade_no,omitempty"`
	// Amount holds the value of the "amount" field.
	// 支付金额
	Amount float64 `json:"amount,omitempty"`
	// PlanDetail holds the value of the "plan_detail" field.
	// 骑士卡详情
	PlanDetail model.PlanItem `json:"plan_detail,omitempty"`
	// ParentID holds the value of the "parent_id" field.
	// 续签/更改电池接续从属订单ID(上个订单)
	ParentID uint64 `json:"parent_id,omitempty"`
	// Subordinate holds the value of the "subordinate" field.
	// 接续订单属性
	Subordinate model.OrderSubordinate `json:"subordinate,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the OrderQuery when eager-loading is set.
	Edges OrderEdges `json:"edges"`
}

// OrderEdges holds the relations/edges for other nodes in the graph.
type OrderEdges struct {
	// Rider holds the value of the rider edge.
	Rider *Rider `json:"rider,omitempty"`
	// Plan holds the value of the plan edge.
	Plan *Plan `json:"plan,omitempty"`
	// Commission holds the value of the commission edge.
	Commission *Commission `json:"commission,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OrderEdges) RiderOrErr() (*Rider, error) {
	if e.loadedTypes[0] {
		if e.Rider == nil {
			// The edge rider was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: rider.Label}
		}
		return e.Rider, nil
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// PlanOrErr returns the Plan value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OrderEdges) PlanOrErr() (*Plan, error) {
	if e.loadedTypes[1] {
		if e.Plan == nil {
			// The edge plan was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: plan.Label}
		}
		return e.Plan, nil
	}
	return nil, &NotLoadedError{edge: "plan"}
}

// CommissionOrErr returns the Commission value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OrderEdges) CommissionOrErr() (*Commission, error) {
	if e.loadedTypes[2] {
		if e.Commission == nil {
			// The edge commission was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: commission.Label}
		}
		return e.Commission, nil
	}
	return nil, &NotLoadedError{edge: "commission"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Order) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case order.FieldCreator, order.FieldLastModifier, order.FieldPlanDetail, order.FieldSubordinate:
			values[i] = new([]byte)
		case order.FieldAmount:
			values[i] = new(sql.NullFloat64)
		case order.FieldID, order.FieldRiderID, order.FieldPlanID, order.FieldStatus, order.FieldPayway, order.FieldType, order.FieldParentID:
			values[i] = new(sql.NullInt64)
		case order.FieldRemark, order.FieldOutTradeNo, order.FieldTradeNo:
			values[i] = new(sql.NullString)
		case order.FieldCreatedAt, order.FieldUpdatedAt, order.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Order", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Order fields.
func (o *Order) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case order.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			o.ID = uint64(value.Int64)
		case order.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				o.CreatedAt = value.Time
			}
		case order.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				o.UpdatedAt = value.Time
			}
		case order.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				o.DeletedAt = new(time.Time)
				*o.DeletedAt = value.Time
			}
		case order.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &o.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case order.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &o.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case order.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				o.Remark = value.String
			}
		case order.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				o.RiderID = uint64(value.Int64)
			}
		case order.FieldPlanID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field plan_id", values[i])
			} else if value.Valid {
				o.PlanID = uint64(value.Int64)
			}
		case order.FieldStatus:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				o.Status = uint8(value.Int64)
			}
		case order.FieldPayway:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field payway", values[i])
			} else if value.Valid {
				o.Payway = uint8(value.Int64)
			}
		case order.FieldType:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				o.Type = uint(value.Int64)
			}
		case order.FieldOutTradeNo:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field out_trade_no", values[i])
			} else if value.Valid {
				o.OutTradeNo = value.String
			}
		case order.FieldTradeNo:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field trade_no", values[i])
			} else if value.Valid {
				o.TradeNo = value.String
			}
		case order.FieldAmount:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field amount", values[i])
			} else if value.Valid {
				o.Amount = value.Float64
			}
		case order.FieldPlanDetail:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field plan_detail", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &o.PlanDetail); err != nil {
					return fmt.Errorf("unmarshal field plan_detail: %w", err)
				}
			}
		case order.FieldParentID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field parent_id", values[i])
			} else if value.Valid {
				o.ParentID = uint64(value.Int64)
			}
		case order.FieldSubordinate:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field subordinate", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &o.Subordinate); err != nil {
					return fmt.Errorf("unmarshal field subordinate: %w", err)
				}
			}
		}
	}
	return nil
}

// QueryRider queries the "rider" edge of the Order entity.
func (o *Order) QueryRider() *RiderQuery {
	return (&OrderClient{config: o.config}).QueryRider(o)
}

// QueryPlan queries the "plan" edge of the Order entity.
func (o *Order) QueryPlan() *PlanQuery {
	return (&OrderClient{config: o.config}).QueryPlan(o)
}

// QueryCommission queries the "commission" edge of the Order entity.
func (o *Order) QueryCommission() *CommissionQuery {
	return (&OrderClient{config: o.config}).QueryCommission(o)
}

// Update returns a builder for updating this Order.
// Note that you need to call Order.Unwrap() before calling this method if this Order
// was returned from a transaction, and the transaction was committed or rolled back.
func (o *Order) Update() *OrderUpdateOne {
	return (&OrderClient{config: o.config}).UpdateOne(o)
}

// Unwrap unwraps the Order entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (o *Order) Unwrap() *Order {
	tx, ok := o.config.driver.(*txDriver)
	if !ok {
		panic("ent: Order is not a transactional entity")
	}
	o.config.driver = tx.drv
	return o
}

// String implements the fmt.Stringer.
func (o *Order) String() string {
	var builder strings.Builder
	builder.WriteString("Order(")
	builder.WriteString(fmt.Sprintf("id=%v", o.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(o.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(o.UpdatedAt.Format(time.ANSIC))
	if v := o.DeletedAt; v != nil {
		builder.WriteString(", deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", creator=")
	builder.WriteString(fmt.Sprintf("%v", o.Creator))
	builder.WriteString(", last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", o.LastModifier))
	builder.WriteString(", remark=")
	builder.WriteString(o.Remark)
	builder.WriteString(", rider_id=")
	builder.WriteString(fmt.Sprintf("%v", o.RiderID))
	builder.WriteString(", plan_id=")
	builder.WriteString(fmt.Sprintf("%v", o.PlanID))
	builder.WriteString(", status=")
	builder.WriteString(fmt.Sprintf("%v", o.Status))
	builder.WriteString(", payway=")
	builder.WriteString(fmt.Sprintf("%v", o.Payway))
	builder.WriteString(", type=")
	builder.WriteString(fmt.Sprintf("%v", o.Type))
	builder.WriteString(", out_trade_no=")
	builder.WriteString(o.OutTradeNo)
	builder.WriteString(", trade_no=")
	builder.WriteString(o.TradeNo)
	builder.WriteString(", amount=")
	builder.WriteString(fmt.Sprintf("%v", o.Amount))
	builder.WriteString(", plan_detail=")
	builder.WriteString(fmt.Sprintf("%v", o.PlanDetail))
	builder.WriteString(", parent_id=")
	builder.WriteString(fmt.Sprintf("%v", o.ParentID))
	builder.WriteString(", subordinate=")
	builder.WriteString(fmt.Sprintf("%v", o.Subordinate))
	builder.WriteByte(')')
	return builder.String()
}

// Orders is a parsable slice of Order.
type Orders []*Order

func (o Orders) config(cfg config) {
	for _i := range o {
		o[_i].config = cfg
	}
}