// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/coupon"
	"github.com/auroraride/aurservd/internal/ent/couponassembly"
	"github.com/auroraride/aurservd/internal/ent/coupontemplate"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/rider"
)

// Coupon is the model entity for the Coupon schema.
type Coupon struct {
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
	// AssemblyID holds the value of the "assembly_id" field.
	AssemblyID uint64 `json:"assembly_id,omitempty"`
	// OrderID holds the value of the "order_id" field.
	OrderID *uint64 `json:"order_id,omitempty"`
	// 实际使用骑士卡
	PlanID *uint64 `json:"plan_id,omitempty"`
	// TemplateID holds the value of the "template_id" field.
	TemplateID uint64 `json:"template_id,omitempty"`
	// 名称
	Name string `json:"name,omitempty"`
	// 使用规则
	Rule uint8 `json:"rule,omitempty"`
	// 该券是否可叠加
	Multiple bool `json:"multiple,omitempty"`
	// 金额
	Amount float64 `json:"amount,omitempty"`
	// 券码
	Code string `json:"code,omitempty"`
	// 过期时间
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	// 使用时间
	UsedAt *time.Time `json:"used_at,omitempty"`
	// 有效期规则
	Duration *model.CouponDuration `json:"duration,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CouponQuery when eager-loading is set.
	Edges CouponEdges `json:"edges"`
}

// CouponEdges holds the relations/edges for other nodes in the graph.
type CouponEdges struct {
	// 骑手
	Rider *Rider `json:"rider,omitempty"`
	// Assembly holds the value of the assembly edge.
	Assembly *CouponAssembly `json:"assembly,omitempty"`
	// Order holds the value of the order edge.
	Order *Order `json:"order,omitempty"`
	// Plan holds the value of the plan edge.
	Plan *Plan `json:"plan,omitempty"`
	// Template holds the value of the template edge.
	Template *CouponTemplate `json:"template,omitempty"`
	// Cities holds the value of the cities edge.
	Cities []*City `json:"cities,omitempty"`
	// Plans holds the value of the plans edge.
	Plans []*Plan `json:"plans,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [7]bool
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CouponEdges) RiderOrErr() (*Rider, error) {
	if e.loadedTypes[0] {
		if e.Rider == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: rider.Label}
		}
		return e.Rider, nil
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// AssemblyOrErr returns the Assembly value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CouponEdges) AssemblyOrErr() (*CouponAssembly, error) {
	if e.loadedTypes[1] {
		if e.Assembly == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: couponassembly.Label}
		}
		return e.Assembly, nil
	}
	return nil, &NotLoadedError{edge: "assembly"}
}

// OrderOrErr returns the Order value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CouponEdges) OrderOrErr() (*Order, error) {
	if e.loadedTypes[2] {
		if e.Order == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: order.Label}
		}
		return e.Order, nil
	}
	return nil, &NotLoadedError{edge: "order"}
}

// PlanOrErr returns the Plan value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CouponEdges) PlanOrErr() (*Plan, error) {
	if e.loadedTypes[3] {
		if e.Plan == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: plan.Label}
		}
		return e.Plan, nil
	}
	return nil, &NotLoadedError{edge: "plan"}
}

// TemplateOrErr returns the Template value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CouponEdges) TemplateOrErr() (*CouponTemplate, error) {
	if e.loadedTypes[4] {
		if e.Template == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: coupontemplate.Label}
		}
		return e.Template, nil
	}
	return nil, &NotLoadedError{edge: "template"}
}

// CitiesOrErr returns the Cities value or an error if the edge
// was not loaded in eager-loading.
func (e CouponEdges) CitiesOrErr() ([]*City, error) {
	if e.loadedTypes[5] {
		return e.Cities, nil
	}
	return nil, &NotLoadedError{edge: "cities"}
}

// PlansOrErr returns the Plans value or an error if the edge
// was not loaded in eager-loading.
func (e CouponEdges) PlansOrErr() ([]*Plan, error) {
	if e.loadedTypes[6] {
		return e.Plans, nil
	}
	return nil, &NotLoadedError{edge: "plans"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Coupon) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case coupon.FieldCreator, coupon.FieldLastModifier, coupon.FieldDuration:
			values[i] = new([]byte)
		case coupon.FieldMultiple:
			values[i] = new(sql.NullBool)
		case coupon.FieldAmount:
			values[i] = new(sql.NullFloat64)
		case coupon.FieldID, coupon.FieldRiderID, coupon.FieldAssemblyID, coupon.FieldOrderID, coupon.FieldPlanID, coupon.FieldTemplateID, coupon.FieldRule:
			values[i] = new(sql.NullInt64)
		case coupon.FieldRemark, coupon.FieldName, coupon.FieldCode:
			values[i] = new(sql.NullString)
		case coupon.FieldCreatedAt, coupon.FieldUpdatedAt, coupon.FieldExpiresAt, coupon.FieldUsedAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Coupon", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Coupon fields.
func (c *Coupon) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case coupon.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = uint64(value.Int64)
		case coupon.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				c.CreatedAt = value.Time
			}
		case coupon.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				c.UpdatedAt = value.Time
			}
		case coupon.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case coupon.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case coupon.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				c.Remark = value.String
			}
		case coupon.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				c.RiderID = new(uint64)
				*c.RiderID = uint64(value.Int64)
			}
		case coupon.FieldAssemblyID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field assembly_id", values[i])
			} else if value.Valid {
				c.AssemblyID = uint64(value.Int64)
			}
		case coupon.FieldOrderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field order_id", values[i])
			} else if value.Valid {
				c.OrderID = new(uint64)
				*c.OrderID = uint64(value.Int64)
			}
		case coupon.FieldPlanID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field plan_id", values[i])
			} else if value.Valid {
				c.PlanID = new(uint64)
				*c.PlanID = uint64(value.Int64)
			}
		case coupon.FieldTemplateID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field template_id", values[i])
			} else if value.Valid {
				c.TemplateID = uint64(value.Int64)
			}
		case coupon.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case coupon.FieldRule:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rule", values[i])
			} else if value.Valid {
				c.Rule = uint8(value.Int64)
			}
		case coupon.FieldMultiple:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field multiple", values[i])
			} else if value.Valid {
				c.Multiple = value.Bool
			}
		case coupon.FieldAmount:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field amount", values[i])
			} else if value.Valid {
				c.Amount = value.Float64
			}
		case coupon.FieldCode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field code", values[i])
			} else if value.Valid {
				c.Code = value.String
			}
		case coupon.FieldExpiresAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field expires_at", values[i])
			} else if value.Valid {
				c.ExpiresAt = value.Time
			}
		case coupon.FieldUsedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field used_at", values[i])
			} else if value.Valid {
				c.UsedAt = new(time.Time)
				*c.UsedAt = value.Time
			}
		case coupon.FieldDuration:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field duration", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.Duration); err != nil {
					return fmt.Errorf("unmarshal field duration: %w", err)
				}
			}
		}
	}
	return nil
}

// QueryRider queries the "rider" edge of the Coupon entity.
func (c *Coupon) QueryRider() *RiderQuery {
	return (&CouponClient{config: c.config}).QueryRider(c)
}

// QueryAssembly queries the "assembly" edge of the Coupon entity.
func (c *Coupon) QueryAssembly() *CouponAssemblyQuery {
	return (&CouponClient{config: c.config}).QueryAssembly(c)
}

// QueryOrder queries the "order" edge of the Coupon entity.
func (c *Coupon) QueryOrder() *OrderQuery {
	return (&CouponClient{config: c.config}).QueryOrder(c)
}

// QueryPlan queries the "plan" edge of the Coupon entity.
func (c *Coupon) QueryPlan() *PlanQuery {
	return (&CouponClient{config: c.config}).QueryPlan(c)
}

// QueryTemplate queries the "template" edge of the Coupon entity.
func (c *Coupon) QueryTemplate() *CouponTemplateQuery {
	return (&CouponClient{config: c.config}).QueryTemplate(c)
}

// QueryCities queries the "cities" edge of the Coupon entity.
func (c *Coupon) QueryCities() *CityQuery {
	return (&CouponClient{config: c.config}).QueryCities(c)
}

// QueryPlans queries the "plans" edge of the Coupon entity.
func (c *Coupon) QueryPlans() *PlanQuery {
	return (&CouponClient{config: c.config}).QueryPlans(c)
}

// Update returns a builder for updating this Coupon.
// Note that you need to call Coupon.Unwrap() before calling this method if this Coupon
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Coupon) Update() *CouponUpdateOne {
	return (&CouponClient{config: c.config}).UpdateOne(c)
}

// Unwrap unwraps the Coupon entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Coupon) Unwrap() *Coupon {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Coupon is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Coupon) String() string {
	var builder strings.Builder
	builder.WriteString("Coupon(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("created_at=")
	builder.WriteString(c.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(c.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", c.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", c.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(c.Remark)
	builder.WriteString(", ")
	if v := c.RiderID; v != nil {
		builder.WriteString("rider_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("assembly_id=")
	builder.WriteString(fmt.Sprintf("%v", c.AssemblyID))
	builder.WriteString(", ")
	if v := c.OrderID; v != nil {
		builder.WriteString("order_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := c.PlanID; v != nil {
		builder.WriteString("plan_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("template_id=")
	builder.WriteString(fmt.Sprintf("%v", c.TemplateID))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(c.Name)
	builder.WriteString(", ")
	builder.WriteString("rule=")
	builder.WriteString(fmt.Sprintf("%v", c.Rule))
	builder.WriteString(", ")
	builder.WriteString("multiple=")
	builder.WriteString(fmt.Sprintf("%v", c.Multiple))
	builder.WriteString(", ")
	builder.WriteString("amount=")
	builder.WriteString(fmt.Sprintf("%v", c.Amount))
	builder.WriteString(", ")
	builder.WriteString("code=")
	builder.WriteString(c.Code)
	builder.WriteString(", ")
	builder.WriteString("expires_at=")
	builder.WriteString(c.ExpiresAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := c.UsedAt; v != nil {
		builder.WriteString("used_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("duration=")
	builder.WriteString(fmt.Sprintf("%v", c.Duration))
	builder.WriteByte(')')
	return builder.String()
}

// Coupons is a parsable slice of Coupon.
type Coupons []*Coupon

func (c Coupons) config(cfg config) {
	for _i := range c {
		c[_i].config = cfg
	}
}